package ai

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/AutumnsGrove/codequest/internal/config"
)

// ModsProvider implements the AIProvider interface using the Mods CLI.
// Mods executes local LLM models and is used for code review and offline fallback.
type ModsProvider struct {
	executable  string
	model       string
	rateLimiter *RateLimiter
	priority    int
	config      *config.AIConfig
}

// NewModsProvider creates a new Mods provider.
// It checks for the mods executable and configures the model to use.
func NewModsProvider(cfg *config.AIConfig) *ModsProvider {
	// Try to find mods executable
	executable, err := exec.LookPath("mods")
	if err != nil {
		// Mods not found, use empty string to mark as unavailable
		executable = ""
	}

	return &ModsProvider{
		executable:  executable,
		model:       cfg.Review.ModelPrimary,
		rateLimiter: NewRateLimiter(10, time.Minute), // 10 requests per minute (conservative for local)
		priority:    2,                                // Second priority after Crush
		config:      cfg,
	}
}

// Ask implements the AIProvider interface.
// It executes the mods CLI with the given prompt and returns the response.
func (m *ModsProvider) Ask(ctx context.Context, req *Request) (*Response, error) {
	// Check if mods is available
	if m.executable == "" {
		return nil, fmt.Errorf("%w: mods executable not found", ErrProviderUnavailable)
	}

	// Check rate limit
	if !m.rateLimiter.Allow() {
		return nil, ErrRateLimited
	}

	// Select model based on complexity or use default
	model := m.selectModel(req.Complexity)

	// Build prompt with context if provided
	prompt := req.Prompt
	if req.Context != "" {
		prompt = fmt.Sprintf("Context:\n%s\n\nQuestion: %s", req.Context, req.Prompt)
	}

	// Build command arguments
	args := []string{
		"--model", model,
	}

	// Add max tokens if specified
	if req.MaxTokens > 0 {
		args = append(args, "--max-tokens", fmt.Sprintf("%d", req.MaxTokens))
	}

	// Add temperature if specified
	if req.Temperature > 0 {
		args = append(args, "--temp", fmt.Sprintf("%.2f", req.Temperature))
	}

	// Add the prompt as the final argument
	args = append(args, prompt)

	// Create command with context
	cmd := exec.CommandContext(ctx, m.executable, args...)

	// Execute command
	startTime := time.Now()
	output, err := cmd.CombinedOutput()
	latency := time.Since(startTime)

	if err != nil {
		// Check if context was cancelled or deadline exceeded
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, ErrProviderTimeout
		}
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil, ctx.Err()
		}

		// Check for specific error patterns
		errStr := string(output)
		if strings.Contains(errStr, "model not found") || strings.Contains(errStr, "model does not exist") {
			return nil, fmt.Errorf("%w: model %s not found", ErrProviderError, model)
		}

		return nil, fmt.Errorf("%w: mods execution failed: %v: %s", ErrProviderError, err, string(output))
	}

	// Parse output
	content := strings.TrimSpace(string(output))
	if content == "" {
		return nil, fmt.Errorf("%w: empty response from mods", ErrProviderError)
	}

	// Build response
	// Note: Mods doesn't provide token usage, so we estimate based on output length
	estimatedTokens := len(content) / 4 // Rough estimate: 4 chars per token

	return &Response{
		Content:    content,
		Provider:   m.GetName(),
		TokensUsed: estimatedTokens,
		Latency:    latency,
		Cached:     false,
		Error:      nil,
		Metadata: map[string]string{
			"model":      model,
			"executable": m.executable,
		},
	}, nil
}

// IsAvailable implements the AIProvider interface.
// It checks if the mods executable exists and is accessible.
func (m *ModsProvider) IsAvailable(ctx context.Context) bool {
	// Check if executable was found during initialization
	if m.executable == "" {
		return false
	}

	// Verify the executable still exists
	_, err := exec.LookPath("mods")
	return err == nil
}

// GetName implements the AIProvider interface.
func (m *ModsProvider) GetName() string {
	return "Mods"
}

// GetPriority implements the AIProvider interface.
// Mods has priority 2 (second to Crush) as the local fallback.
func (m *ModsProvider) GetPriority() int {
	return m.priority
}

// GetRateLimiter implements the AIProvider interface.
func (m *ModsProvider) GetRateLimiter() *RateLimiter {
	return m.rateLimiter
}

// selectModel chooses the appropriate model based on request complexity.
// For complex queries, use the primary model; for simple ones, use the fallback.
func (m *ModsProvider) selectModel(complexity string) string {
	if complexity == "simple" && m.config.Review.ModelFallback != "" {
		return m.config.Review.ModelFallback
	}
	// Default to primary model for unspecified or "complex" complexity
	return m.config.Review.ModelPrimary
}

// SetModel updates the model to use (useful for runtime configuration changes)
func (m *ModsProvider) SetModel(model string) {
	m.model = model
}

// ReviewCode is a convenience method for code review operations.
// It sends a code diff to Mods with a code review prompt.
func (m *ModsProvider) ReviewCode(ctx context.Context, diff string) (*Response, error) {
	prompt := fmt.Sprintf(`Review this code diff and provide:
1. Overall quality score (0-100)
2. Key strengths
3. Areas for improvement
4. Go idioms and best practices
5. Suggested bonus XP (0-50 based on quality)

Code Diff:
%s

Provide your review in a clear, structured format.`, diff)

	req := &Request{
		Prompt:      prompt,
		Context:     "",
		MaxTokens:   1000,
		Temperature: 0.7,
		Complexity:  "complex", // Code reviews are complex tasks
		Metadata: map[string]string{
			"type": "code_review",
		},
	}

	return m.Ask(ctx, req)
}

// AskOffline is a convenience method for offline Crush queries routed through Mods.
// It formats the prompt for mentor-style questions.
func (m *ModsProvider) AskOffline(ctx context.Context, question string, contextInfo string, complexity string) (*Response, error) {
	// Add mentor context to the prompt
	systemContext := "You are Crush, a helpful AI mentor for CodeQuest. Provide clear, concise answers to help developers with their coding questions."

	prompt := fmt.Sprintf("%s\n\nQuestion: %s", systemContext, question)

	req := &Request{
		Prompt:      prompt,
		Context:     contextInfo,
		MaxTokens:   800,
		Temperature: 0.7,
		Complexity:  complexity,
		Metadata: map[string]string{
			"type": "offline_mentor",
		},
	}

	return m.Ask(ctx, req)
}
