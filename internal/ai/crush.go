package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AutumnsGrove/codequest/internal/config"
)

// CrushProvider implements the AIProvider interface using OpenRouter API.
// Crush serves as the primary in-game mentor for quick help and guidance.
type CrushProvider struct {
	apiKey      string
	baseURL     string
	modelOnline string
	httpClient  *http.Client
	rateLimiter *RateLimiter
	priority    int
	config      *config.AIConfig
}

// openRouterRequest represents the request format for OpenRouter API
type openRouterRequest struct {
	Model       string                 `json:"model"`
	Messages    []openRouterMessage    `json:"messages"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// openRouterMessage represents a message in the OpenRouter chat format
type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openRouterResponse represents the response from OpenRouter API
type openRouterResponse struct {
	ID      string             `json:"id"`
	Choices []openRouterChoice `json:"choices"`
	Usage   openRouterUsage    `json:"usage"`
	Error   *openRouterError   `json:"error,omitempty"`
}

// openRouterChoice represents a completion choice
type openRouterChoice struct {
	Message      openRouterMessage `json:"message"`
	FinishReason string            `json:"finish_reason"`
}

// openRouterUsage tracks token usage
type openRouterUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// openRouterError represents an error from OpenRouter
type openRouterError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
}

// NewCrushProvider creates a new Crush provider using OpenRouter API.
// It requires an API key from the config and sets up HTTP client with timeouts.
func NewCrushProvider(apiKey string, cfg *config.AIConfig) *CrushProvider {
	// Select appropriate model based on complexity
	// Default to complex model - we'll switch based on request
	modelOnline := cfg.Mentor.ModelComplex

	return &CrushProvider{
		apiKey:      apiKey,
		baseURL:     "https://openrouter.ai/api/v1",
		modelOnline: modelOnline,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		rateLimiter: NewRateLimiter(20, time.Minute), // 20 requests per minute
		priority:    1,                               // Highest priority
		config:      cfg,
	}
}

// Ask implements the AIProvider interface.
// It sends a request to OpenRouter API and returns the response.
func (c *CrushProvider) Ask(ctx context.Context, req *Request) (*Response, error) {
	// Check rate limit
	if !c.rateLimiter.Allow() {
		return nil, ErrRateLimited
	}

	// Select model based on complexity
	model := c.selectModel(req.Complexity)

	// Build system message
	systemMsg := "You are Crush, a helpful AI mentor for CodeQuest. Provide clear, concise answers to help developers with their coding questions."

	// Combine prompt and context
	userContent := req.Prompt
	if req.Context != "" {
		userContent = fmt.Sprintf("Context:\n%s\n\nQuestion: %s", req.Context, req.Prompt)
	}

	// Build OpenRouter request
	orReq := openRouterRequest{
		Model: model,
		Messages: []openRouterMessage{
			{Role: "system", Content: systemMsg},
			{Role: "user", Content: userContent},
		},
		Temperature: req.Temperature,
	}

	// Set max tokens if specified
	if req.MaxTokens > 0 {
		orReq.MaxTokens = req.MaxTokens
	}

	// Add metadata
	orReq.Metadata = map[string]interface{}{
		"app": "CodeQuest",
	}

	// Marshal request
	reqBody, err := json.Marshal(orReq)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, fmt.Errorf("creating HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://github.com/AutumnsGrove/codequest")
	httpReq.Header.Set("X-Title", "CodeQuest")
	httpReq.Header.Set("User-Agent", "CodeQuest/1.0")

	// Send request
	startTime := time.Now()
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		// Check if context deadline exceeded
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrProviderTimeout
		}
		return nil, fmt.Errorf("sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode == 429 {
		return nil, ErrRateLimited
	}
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("authentication failed: invalid API key")
	}
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("%w: OpenRouter server error (status %d)", ErrProviderError, resp.StatusCode)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%w: unexpected status code %d: %s", ErrProviderError, resp.StatusCode, string(body))
	}

	// Parse response
	var orResp openRouterResponse
	if err := json.Unmarshal(body, &orResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	// Check for API error
	if orResp.Error != nil {
		return nil, fmt.Errorf("%w: %s", ErrProviderError, orResp.Error.Message)
	}

	// Validate response
	if len(orResp.Choices) == 0 {
		return nil, fmt.Errorf("%w: no choices in response", ErrProviderError)
	}

	// Extract content
	content := strings.TrimSpace(orResp.Choices[0].Message.Content)
	if content == "" {
		return nil, fmt.Errorf("%w: empty response content", ErrProviderError)
	}

	// Build response
	return &Response{
		Content:    content,
		Provider:   c.GetName(),
		TokensUsed: orResp.Usage.TotalTokens,
		Latency:    time.Since(startTime),
		Cached:     false,
		Error:      nil,
		Metadata: map[string]string{
			"model":         model,
			"finish_reason": orResp.Choices[0].FinishReason,
		},
	}, nil
}

// IsAvailable implements the AIProvider interface.
// It checks if the provider can be used (API key is set and internet is available).
func (c *CrushProvider) IsAvailable(ctx context.Context) bool {
	// Check if API key is set
	if c.apiKey == "" {
		return false
	}

	// Optional: Ping OpenRouter API to verify connectivity
	// For MVP, we'll just check if the key exists
	// In the future, we could add a health check endpoint call

	return true
}

// GetName implements the AIProvider interface.
func (c *CrushProvider) GetName() string {
	return "Crush"
}

// GetPriority implements the AIProvider interface.
// Crush has priority 1 (highest) as the primary mentor.
func (c *CrushProvider) GetPriority() int {
	return c.priority
}

// GetRateLimiter implements the AIProvider interface.
func (c *CrushProvider) GetRateLimiter() *RateLimiter {
	return c.rateLimiter
}

// selectModel chooses the appropriate model based on request complexity.
// "simple" queries use the lighter model, "complex" uses the more powerful one.
func (c *CrushProvider) selectModel(complexity string) string {
	if complexity == "simple" {
		return c.config.Mentor.ModelSimple
	}
	// Default to complex model for unspecified or "complex" complexity
	return c.config.Mentor.ModelComplex
}

// SetAPIKey updates the API key (useful for runtime configuration changes)
func (c *CrushProvider) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}
