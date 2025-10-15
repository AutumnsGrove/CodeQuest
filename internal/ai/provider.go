package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/AutumnsGrove/codequest/internal/config"
)

// Common errors returned by AI providers
var (
	ErrNoProvidersAvailable = errors.New("no AI providers available")
	ErrRateLimited          = errors.New("rate limit exceeded")
	ErrProviderTimeout      = errors.New("provider request timeout")
	ErrInvalidRequest       = errors.New("invalid request parameters")
	ErrProviderError        = errors.New("provider returned error")
	ErrProviderUnavailable  = errors.New("provider unavailable")
)

// AIProvider defines the interface all AI providers must implement.
// Each provider (Crush, Mods, Claude Code) should implement this interface
// to allow seamless fallback and hot-swapping.
type AIProvider interface {
	// Ask sends a question to the AI and returns the response.
	// The context can be used for cancellation and timeout control.
	Ask(ctx context.Context, request *Request) (*Response, error)

	// IsAvailable checks if the provider is currently available.
	// This can check API connectivity, local model availability, etc.
	IsAvailable(ctx context.Context) bool

	// GetName returns the provider's name (e.g., "Crush", "Mods", "Claude").
	GetName() string

	// GetPriority returns provider priority where lower = higher priority.
	// Example: Crush=1, Mods=2, Claude=3
	GetPriority() int

	// GetRateLimiter returns the rate limiter for this provider.
	// Returns nil if rate limiting is not needed.
	GetRateLimiter() *RateLimiter
}

// Request encapsulates an AI query with all necessary context.
type Request struct {
	// Prompt is the user's question or instruction
	Prompt string

	// Context provides optional context like code snippets, error messages, etc.
	Context string

	// MaxTokens is the maximum response length (0 = provider default)
	MaxTokens int

	// Temperature controls creativity (0.0-1.0, higher = more creative)
	Temperature float64

	// Metadata stores additional key-value pairs for provider-specific options
	Metadata map[string]string

	// Complexity hints at query complexity: "simple" or "complex"
	// Providers can use this to select appropriate models
	Complexity string
}

// Response encapsulates an AI response with metadata.
type Response struct {
	// Content is the AI's answer
	Content string

	// Provider is the name of the provider that answered
	Provider string

	// TokensUsed is the number of tokens consumed
	TokensUsed int

	// Latency is the time taken to get the response
	Latency time.Duration

	// Cached indicates if the response came from cache
	Cached bool

	// Error is any error encountered (nil on success)
	Error error

	// Metadata stores provider-specific response metadata
	Metadata map[string]string
}

// AIManager manages multiple AI providers with fallback chain support.
// It handles provider registration, selection, fallback logic, and caching.
type AIManager struct {
	// providers is the list of registered providers, sorted by priority
	providers []AIProvider

	// config is the application configuration
	config *config.Config

	// cache is the optional response cache
	cache *ResponseCache

	// mu protects the providers slice for thread-safe operations
	mu sync.RWMutex

	// stats tracks usage statistics per provider
	stats map[string]*ProviderStats
}

// ProviderStats tracks statistics for a provider.
type ProviderStats struct {
	TotalRequests   int64
	SuccessfulCalls int64
	FailedCalls     int64
	TotalLatency    time.Duration
	LastSuccess     time.Time
	LastFailure     time.Time
	mu              sync.Mutex
}

// NewAIManager creates a new AI manager with the given configuration.
// Call RegisterProvider to add providers before using.
func NewAIManager(cfg *config.Config) *AIManager {
	return &AIManager{
		providers: make([]AIProvider, 0),
		config:    cfg,
		cache:     NewResponseCache(15 * time.Minute), // 15-minute cache TTL
		stats:     make(map[string]*ProviderStats),
	}
}

// RegisterProvider adds a provider to the manager.
// Providers are automatically sorted by priority after registration.
func (m *AIManager) RegisterProvider(provider AIProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Add provider to the list
	m.providers = append(m.providers, provider)

	// Initialize stats for this provider
	if _, exists := m.stats[provider.GetName()]; !exists {
		m.stats[provider.GetName()] = &ProviderStats{}
	}

	// Sort providers by priority (lower number = higher priority)
	m.sortProviders()
}

// sortProviders sorts the provider list by priority.
// Must be called with m.mu held.
func (m *AIManager) sortProviders() {
	// Simple bubble sort since we have a small number of providers
	n := len(m.providers)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if m.providers[j].GetPriority() > m.providers[j+1].GetPriority() {
				m.providers[j], m.providers[j+1] = m.providers[j+1], m.providers[j]
			}
		}
	}
}

// Ask sends a request through the fallback chain until a provider succeeds.
// It tries each provider in priority order, checking availability and rate limits.
func (m *AIManager) Ask(ctx context.Context, req *Request) (*Response, error) {
	// Validate request
	if err := m.validateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	// Check cache first
	if cacheKey := m.getCacheKey(req); cacheKey != "" {
		if cached, found := m.cache.Get(cacheKey); found {
			cached.Cached = true
			return cached, nil
		}
	}

	// Try each provider in priority order
	m.mu.RLock()
	providers := make([]AIProvider, len(m.providers))
	copy(providers, m.providers)
	m.mu.RUnlock()

	var lastErr error
	for _, provider := range providers {
		// Check if provider is available
		if !provider.IsAvailable(ctx) {
			lastErr = fmt.Errorf("%s: %w", provider.GetName(), ErrProviderUnavailable)
			continue
		}

		// Check rate limit
		if limiter := provider.GetRateLimiter(); limiter != nil {
			if !limiter.Allow() {
				lastErr = fmt.Errorf("%s: %w", provider.GetName(), ErrRateLimited)
				continue
			}
		}

		// Attempt the request
		startTime := time.Now()
		resp, err := provider.Ask(ctx, req)
		latency := time.Since(startTime)

		// Record stats
		m.recordStats(provider.GetName(), err == nil, latency)

		if err != nil {
			lastErr = fmt.Errorf("%s: %w", provider.GetName(), err)
			continue
		}

		// Success! Fill in metadata and cache
		resp.Provider = provider.GetName()
		resp.Latency = latency

		// Cache the response
		if cacheKey := m.getCacheKey(req); cacheKey != "" {
			m.cache.Set(cacheKey, resp)
		}

		return resp, nil
	}

	// All providers failed
	if lastErr != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoProvidersAvailable, lastErr)
	}
	return nil, ErrNoProvidersAvailable
}

// AskSpecific sends a request to a specific provider by name.
// This bypasses the fallback chain and targets a single provider.
func (m *AIManager) AskSpecific(ctx context.Context, providerName string, req *Request) (*Response, error) {
	// Validate request
	if err := m.validateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	// Find the provider
	m.mu.RLock()
	var targetProvider AIProvider
	for _, provider := range m.providers {
		if provider.GetName() == providerName {
			targetProvider = provider
			break
		}
	}
	m.mu.RUnlock()

	if targetProvider == nil {
		return nil, fmt.Errorf("provider %q not found", providerName)
	}

	// Check availability
	if !targetProvider.IsAvailable(ctx) {
		return nil, fmt.Errorf("%s: %w", providerName, ErrProviderUnavailable)
	}

	// Check rate limit
	if limiter := targetProvider.GetRateLimiter(); limiter != nil {
		if !limiter.Allow() {
			return nil, fmt.Errorf("%s: %w", providerName, ErrRateLimited)
		}
	}

	// Make the request
	startTime := time.Now()
	resp, err := targetProvider.Ask(ctx, req)
	latency := time.Since(startTime)

	// Record stats
	m.recordStats(providerName, err == nil, latency)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", providerName, err)
	}

	// Fill in metadata
	resp.Provider = providerName
	resp.Latency = latency

	return resp, nil
}

// GetAvailableProviders returns a list of currently available provider names.
func (m *AIManager) GetAvailableProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	available := make([]string, 0, len(m.providers))
	for _, provider := range m.providers {
		if provider.IsAvailable(ctx) {
			available = append(available, provider.GetName())
		}
	}
	return available
}

// HealthCheck checks the availability of all registered providers.
// Returns a map of provider name to availability status.
func (m *AIManager) HealthCheck(ctx context.Context) map[string]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	health := make(map[string]bool)
	for _, provider := range m.providers {
		health[provider.GetName()] = provider.IsAvailable(ctx)
	}
	return health
}

// GetStats returns statistics for a specific provider.
func (m *AIManager) GetStats(providerName string) *ProviderStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if stats, exists := m.stats[providerName]; exists {
		stats.mu.Lock()
		defer stats.mu.Unlock()
		// Return a copy to prevent external modification
		return &ProviderStats{
			TotalRequests:   stats.TotalRequests,
			SuccessfulCalls: stats.SuccessfulCalls,
			FailedCalls:     stats.FailedCalls,
			TotalLatency:    stats.TotalLatency,
			LastSuccess:     stats.LastSuccess,
			LastFailure:     stats.LastFailure,
		}
	}
	return nil
}

// validateRequest checks if a request is valid.
func (m *AIManager) validateRequest(req *Request) error {
	if req == nil {
		return errors.New("request is nil")
	}
	if req.Prompt == "" {
		return errors.New("prompt is empty")
	}
	if req.Temperature < 0 || req.Temperature > 1 {
		return errors.New("temperature must be between 0 and 1")
	}
	return nil
}

// getCacheKey generates a cache key for a request.
// Returns empty string if caching should be disabled for this request.
func (m *AIManager) getCacheKey(req *Request) string {
	// Don't cache if metadata says not to
	if req.Metadata != nil && req.Metadata["no_cache"] == "true" {
		return ""
	}

	// Generate hash of prompt + context
	hasher := sha256.New()
	hasher.Write([]byte(req.Prompt))
	hasher.Write([]byte(req.Context))
	hasher.Write([]byte(req.Complexity))
	return hex.EncodeToString(hasher.Sum(nil))
}

// recordStats updates statistics for a provider.
func (m *AIManager) recordStats(providerName string, success bool, latency time.Duration) {
	m.mu.RLock()
	stats, exists := m.stats[providerName]
	m.mu.RUnlock()

	if !exists {
		return
	}

	stats.mu.Lock()
	defer stats.mu.Unlock()

	stats.TotalRequests++
	stats.TotalLatency += latency

	if success {
		stats.SuccessfulCalls++
		stats.LastSuccess = time.Now()
	} else {
		stats.FailedCalls++
		stats.LastFailure = time.Now()
	}
}

// RateLimiter implements a token bucket rate limiter.
// It allows a certain number of requests within a time window.
type RateLimiter struct {
	// requests is the maximum number of requests allowed
	requests int

	// window is the time window for rate limiting
	window time.Duration

	// tokens is the current number of available tokens
	tokens int

	// lastRefill is the last time tokens were refilled
	lastRefill time.Time

	// mu protects the limiter state
	mu sync.Mutex
}

// NewRateLimiter creates a new rate limiter.
// Example: NewRateLimiter(10, 1*time.Minute) allows 10 requests per minute.
func NewRateLimiter(requests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:   requests,
		window:     window,
		tokens:     requests,
		lastRefill: time.Now(),
	}
}

// Allow returns true if a request is allowed, false if rate limited.
// This is a non-blocking check.
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.refill()

	if r.tokens > 0 {
		r.tokens--
		return true
	}
	return false
}

// Wait blocks until a token is available or the context is cancelled.
// Returns an error if the context is cancelled before a token becomes available.
func (r *RateLimiter) Wait(ctx context.Context) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if r.Allow() {
				return nil
			}
		}
	}
}

// Reset resets the rate limiter to full capacity.
func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tokens = r.requests
	r.lastRefill = time.Now()
}

// refill adds tokens based on elapsed time since last refill.
// Must be called with r.mu held.
func (r *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(r.lastRefill)

	// If enough time has passed, refill tokens
	if elapsed >= r.window {
		r.tokens = r.requests
		r.lastRefill = now
	}
}

// ResponseCache caches AI responses to avoid redundant API calls.
type ResponseCache struct {
	// cache stores cached responses by key
	cache map[string]*CachedResponse

	// ttl is how long cached responses remain valid
	ttl time.Duration

	// mu protects the cache
	mu sync.RWMutex
}

// CachedResponse wraps a response with timestamp information.
type CachedResponse struct {
	response  *Response
	timestamp time.Time
}

// NewResponseCache creates a new response cache with the given TTL.
func NewResponseCache(ttl time.Duration) *ResponseCache {
	return &ResponseCache{
		cache: make(map[string]*CachedResponse),
		ttl:   ttl,
	}
}

// Get retrieves a cached response if it exists and is still valid.
func (c *ResponseCache) Get(key string) (*Response, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	// Check if cache entry is still valid
	if time.Since(cached.timestamp) > c.ttl {
		return nil, false
	}

	return cached.response, true
}

// Set stores a response in the cache.
func (c *ResponseCache) Set(key string, response *Response) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &CachedResponse{
		response:  response,
		timestamp: time.Now(),
	}
}

// Invalidate removes a specific cache entry.
func (c *ResponseCache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}

// Clear removes all cache entries.
func (c *ResponseCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*CachedResponse)
}

// CleanExpired removes expired cache entries.
// Should be called periodically to prevent memory bloat.
func (c *ResponseCache) CleanExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, cached := range c.cache {
		if now.Sub(cached.timestamp) > c.ttl {
			delete(c.cache, key)
		}
	}
}
