// Package screens provides screen rendering functions for CodeQuest UI.
// This file tests the Mentor screen rendering functions.
package screens

import (
	"strings"
	"testing"
	"time"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// TestRenderMentor tests the main mentor screen rendering function.
func TestRenderMentor(t *testing.T) {
	t.Run("renders with nil character", func(t *testing.T) {
		result := RenderMentor(nil, "", []Message{}, 80, 24)

		if result == "" {
			t.Error("Expected non-empty output for nil character")
		}

		// Should still render without crashing
		if !strings.Contains(result, "AI Mentor") {
			t.Error("Expected 'AI Mentor' in output")
		}
	})

	t.Run("renders with empty conversation", func(t *testing.T) {
		char := game.NewCharacter("TestPlayer")
		result := RenderMentor(char, "", []Message{}, 80, 24)

		if result == "" {
			t.Error("Expected non-empty output")
		}

		// Should show example prompts
		if !strings.Contains(result, "Need inspiration") {
			t.Error("Expected example prompts when conversation is empty")
		}

		// Should show some example questions
		if !strings.Contains(result, "optimize this function") {
			t.Error("Expected example question in output")
		}
	})

	t.Run("renders with conversation history", func(t *testing.T) {
		char := game.NewCharacter("TestPlayer")
		messages := []Message{
			{
				Sender:    "system",
				Content:   "Session started",
				Timestamp: time.Now().Add(-10 * time.Minute),
			},
			{
				Sender:    "user",
				Content:   "How do I optimize this loop?",
				Timestamp: time.Now().Add(-5 * time.Minute),
			},
			{
				Sender:    "crush",
				Content:   "You can optimize by using a map instead of nested loops.",
				Timestamp: time.Now(),
			},
		}

		result := RenderMentor(char, "", messages, 80, 24)

		if result == "" {
			t.Error("Expected non-empty output")
		}

		// Should contain conversation elements
		if !strings.Contains(result, "Conversation") {
			t.Error("Expected 'Conversation' header")
		}

		// Should not show example prompts when conversation exists
		if strings.Contains(result, "Need inspiration") {
			t.Error("Should not show example prompts when conversation exists")
		}
	})

	t.Run("renders with input text", func(t *testing.T) {
		char := game.NewCharacter("TestPlayer")
		inputText := "What is the best way to structure my code?"

		result := RenderMentor(char, inputText, []Message{}, 80, 24)

		if result == "" {
			t.Error("Expected non-empty output")
		}

		// Should show the input text
		if !strings.Contains(result, inputText) {
			t.Error("Expected input text to be displayed")
		}
	})

	t.Run("renders with small terminal", func(t *testing.T) {
		char := game.NewCharacter("TestPlayer")
		result := RenderMentor(char, "", []Message{}, 40, 20)

		if result == "" {
			t.Error("Expected non-empty output for small terminal")
		}

		// Should still render without panicking
		if !strings.Contains(result, "AI Mentor") {
			t.Error("Expected 'AI Mentor' in output")
		}
	})

	t.Run("renders with large terminal", func(t *testing.T) {
		char := game.NewCharacter("TestPlayer")
		result := RenderMentor(char, "", []Message{}, 120, 40)

		if result == "" {
			t.Error("Expected non-empty output for large terminal")
		}
	})
}

// TestRenderMentorHeader tests the header rendering.
func TestRenderMentorHeader(t *testing.T) {
	t.Run("renders header with character", func(t *testing.T) {
		char := game.NewCharacter("TestHero")
		char.Level = 5

		result := renderMentorHeader(char, 80)

		// Due to ANSI styling codes, we check for character name presence
		if !strings.Contains(result, "TestHero") {
			t.Error("Expected character name in header")
		}

		// Check for screen name
		if !strings.Contains(result, "AI Mentor") {
			t.Error("Expected screen name in header")
		}

		// Result should not be empty
		if result == "" {
			t.Error("Expected non-empty header")
		}
	})

	t.Run("renders header without character", func(t *testing.T) {
		result := renderMentorHeader(nil, 80)

		// Check for screen name
		if !strings.Contains(result, "AI Mentor") {
			t.Error("Expected screen name in header")
		}

		// Result should not be empty
		if result == "" {
			t.Error("Expected non-empty header")
		}
	})

	t.Run("renders minimal header for narrow width", func(t *testing.T) {
		char := game.NewCharacter("TestHero")
		result := renderMentorHeader(char, 30)

		if !strings.Contains(result, "AI Mentor") {
			t.Error("Expected screen name even in minimal header")
		}
	})
}

// TestRenderAIStatus tests the AI provider status rendering.
func TestRenderAIStatus(t *testing.T) {
	t.Run("renders online status", func(t *testing.T) {
		status := AIProviderStatus{
			Primary:   "Crush",
			Available: true,
			Fallback:  "Mods",
		}

		result := renderAIStatus(status, 80)

		if !strings.Contains(result, "Crush") {
			t.Error("Expected provider name in status")
		}

		if !strings.Contains(result, "Online") {
			t.Error("Expected 'Online' indicator")
		}
	})

	t.Run("renders offline status", func(t *testing.T) {
		status := AIProviderStatus{
			Primary:   "Crush",
			Available: false,
			Fallback:  "Mods",
		}

		result := renderAIStatus(status, 80)

		if !strings.Contains(result, "Offline") {
			t.Error("Expected 'Offline' indicator when not available")
		}

		if !strings.Contains(result, "Mods") {
			t.Error("Expected fallback provider name")
		}
	})
}

// TestRenderConversation tests conversation history rendering.
func TestRenderConversation(t *testing.T) {
	t.Run("renders empty conversation", func(t *testing.T) {
		result := renderConversation([]Message{}, 80, 20)

		if !strings.Contains(result, "Conversation") {
			t.Error("Expected 'Conversation' header")
		}
	})

	t.Run("renders multiple messages", func(t *testing.T) {
		messages := []Message{
			{
				Sender:    "user",
				Content:   "Hello AI",
				Timestamp: time.Now(),
			},
			{
				Sender:    "crush",
				Content:   "Hello! How can I help?",
				Timestamp: time.Now(),
			},
		}

		result := renderConversation(messages, 80, 20)

		if !strings.Contains(result, "Hello AI") {
			t.Error("Expected user message content")
		}

		if !strings.Contains(result, "Hello! How can I help?") {
			t.Error("Expected AI response content")
		}
	})

	t.Run("handles long conversation with scrolling", func(t *testing.T) {
		// Create many messages
		messages := make([]Message, 20)
		for i := 0; i < 20; i++ {
			messages[i] = Message{
				Sender:    "user",
				Content:   "Message " + string(rune('A'+i)),
				Timestamp: time.Now(),
			}
		}

		result := renderConversation(messages, 80, 20)

		// Should render without panicking
		if result == "" {
			t.Error("Expected non-empty output for long conversation")
		}
	})
}

// TestRenderMessage tests individual message rendering.
func TestRenderMessage(t *testing.T) {
	timestamp := time.Now()

	t.Run("renders user message", func(t *testing.T) {
		msg := Message{
			Sender:    "user",
			Content:   "Test user message",
			Timestamp: timestamp,
		}

		result := renderMessage(msg, 80)

		if !strings.Contains(result, "Test user message") {
			t.Error("Expected message content")
		}

		if !strings.Contains(result, "[You]") {
			t.Error("Expected user indicator")
		}
	})

	t.Run("renders crush message", func(t *testing.T) {
		msg := Message{
			Sender:    "crush",
			Content:   "Test AI response",
			Timestamp: timestamp,
		}

		result := renderMessage(msg, 80)

		if !strings.Contains(result, "Test AI response") {
			t.Error("Expected message content")
		}

		if !strings.Contains(result, "[Crush]") {
			t.Error("Expected AI sender indicator")
		}
	})

	t.Run("renders mods message", func(t *testing.T) {
		msg := Message{
			Sender:    "mods",
			Content:   "Local AI response",
			Timestamp: timestamp,
		}

		result := renderMessage(msg, 80)

		if !strings.Contains(result, "Local AI response") {
			t.Error("Expected message content")
		}

		if !strings.Contains(result, "[Mods]") {
			t.Error("Expected Mods sender indicator")
		}
	})

	t.Run("renders claude message", func(t *testing.T) {
		msg := Message{
			Sender:    "claude",
			Content:   "Claude AI response",
			Timestamp: timestamp,
		}

		result := renderMessage(msg, 80)

		if !strings.Contains(result, "Claude AI response") {
			t.Error("Expected message content")
		}

		if !strings.Contains(result, "[Claude]") {
			t.Error("Expected Claude sender indicator")
		}
	})

	t.Run("renders system message", func(t *testing.T) {
		msg := Message{
			Sender:    "system",
			Content:   "Session started",
			Timestamp: timestamp,
		}

		result := renderMessage(msg, 80)

		if !strings.Contains(result, "Session started") {
			t.Error("Expected message content")
		}
	})
}

// TestRenderExamplePrompts tests example prompts rendering.
func TestRenderExamplePrompts(t *testing.T) {
	t.Run("renders example prompts", func(t *testing.T) {
		result := renderExamplePrompts(80, 20)

		if !strings.Contains(result, "Need inspiration") {
			t.Error("Expected 'Need inspiration' header")
		}

		if !strings.Contains(result, "optimize this function") {
			t.Error("Expected example prompt")
		}

		if !strings.Contains(result, "Welcome to the AI Mentor") {
			t.Error("Expected welcome message")
		}
	})

	t.Run("renders with different widths", func(t *testing.T) {
		result := renderExamplePrompts(60, 20)

		if result == "" {
			t.Error("Expected non-empty output")
		}
	})
}

// TestRenderInputArea tests input field rendering.
func TestRenderInputArea(t *testing.T) {
	t.Run("renders empty input with placeholder", func(t *testing.T) {
		result := renderInputArea("", 80)

		if !strings.Contains(result, "Ask anything about your code") {
			t.Error("Expected placeholder text when input is empty")
		}

		if !strings.Contains(result, "[Enter]") {
			t.Error("Expected Enter key hint")
		}
	})

	t.Run("renders input with text", func(t *testing.T) {
		inputText := "How do I fix this bug?"
		result := renderInputArea(inputText, 80)

		if !strings.Contains(result, inputText) {
			t.Error("Expected input text to be displayed")
		}

		// Should show cursor indicator
		if !strings.Contains(result, "█") {
			t.Error("Expected cursor indicator")
		}
	})

	t.Run("renders with different widths", func(t *testing.T) {
		result := renderInputArea("test", 40)

		if result == "" {
			t.Error("Expected non-empty output for narrow width")
		}
	})
}

// TestRenderMentorFooter tests footer rendering.
func TestRenderMentorFooter(t *testing.T) {
	t.Run("renders footer with key bindings", func(t *testing.T) {
		result := renderMentorFooter(80)

		if !strings.Contains(result, "[Enter]") {
			t.Error("Expected Enter key binding")
		}

		if !strings.Contains(result, "[Esc]") {
			t.Error("Expected Esc key binding")
		}

		if !strings.Contains(result, "[Ctrl+C]") {
			t.Error("Expected Ctrl+C key binding")
		}
	})
}

// TestWrapText tests text wrapping utility.
func TestWrapText(t *testing.T) {
	t.Run("wraps long text", func(t *testing.T) {
		text := "This is a very long line of text that should be wrapped to fit within the specified width"
		result := wrapText(text, 30)

		lines := strings.Split(result, "\n")
		if len(lines) <= 1 {
			t.Error("Expected text to be wrapped into multiple lines")
		}

		// Check that no line exceeds width
		for _, line := range lines {
			if len(line) > 30 {
				t.Errorf("Line exceeds width: %d > 30", len(line))
			}
		}
	})

	t.Run("preserves existing line breaks", func(t *testing.T) {
		text := "Line 1\nLine 2\nLine 3"
		result := wrapText(text, 50)

		if !strings.Contains(result, "Line 1") ||
			!strings.Contains(result, "Line 2") ||
			!strings.Contains(result, "Line 3") {
			t.Error("Expected existing line breaks to be preserved")
		}
	})

	t.Run("handles empty text", func(t *testing.T) {
		result := wrapText("", 30)

		if result != "" {
			t.Error("Expected empty result for empty input")
		}
	})

	t.Run("handles text shorter than width", func(t *testing.T) {
		text := "Short text"
		result := wrapText(text, 50)

		if !strings.Contains(result, text) {
			t.Error("Expected text to be unchanged when shorter than width")
		}
	})

	t.Run("handles zero or negative width", func(t *testing.T) {
		text := "Test text"
		result := wrapText(text, 0)

		if result == "" {
			t.Error("Expected non-empty result even with zero width")
		}

		result = wrapText(text, -10)
		if result == "" {
			t.Error("Expected non-empty result even with negative width")
		}
	})
}

// TestFormatTime tests timestamp formatting.
func TestFormatTime(t *testing.T) {
	t.Run("formats today's time", func(t *testing.T) {
		now := time.Now()
		result := formatTime(now)

		// Should show time only (e.g., "3:04 PM")
		if !strings.Contains(result, "M") {
			t.Error("Expected AM/PM indicator for today's time")
		}

		// Should not include year
		if strings.Contains(result, "2024") || strings.Contains(result, "2025") {
			t.Error("Should not include year for today's time")
		}
	})

	t.Run("formats this year's date", func(t *testing.T) {
		thisYear := time.Date(time.Now().Year(), 6, 15, 14, 30, 0, 0, time.Local)
		result := formatTime(thisYear)

		// Should include month and day
		if !strings.Contains(result, "Jun") {
			t.Error("Expected month in format")
		}

		if !strings.Contains(result, "15") {
			t.Error("Expected day in format")
		}
	})

	t.Run("formats old date", func(t *testing.T) {
		oldDate := time.Date(2020, 3, 10, 14, 30, 0, 0, time.Local)
		result := formatTime(oldDate)

		// Should include full date with year
		if !strings.Contains(result, "2020") {
			t.Error("Expected year in format for old date")
		}

		if !strings.Contains(result, "Mar") {
			t.Error("Expected month in format")
		}
	})
}

// TestGetAIProviderStatus tests AI provider status retrieval.
func TestGetAIProviderStatus(t *testing.T) {
	t.Run("returns status", func(t *testing.T) {
		status := getAIProviderStatus()

		if status.Primary == "" {
			t.Error("Expected primary provider name")
		}

		if status.Fallback == "" {
			t.Error("Expected fallback provider name")
		}

		// Currently should be offline (placeholder)
		if status.Available {
			t.Log("Note: AI provider is available (may be placeholder)")
		}
	})
}

// TestMessageStruct tests the Message struct.
func TestMessageStruct(t *testing.T) {
	t.Run("creates message", func(t *testing.T) {
		msg := Message{
			Sender:    "user",
			Content:   "Test message",
			Timestamp: time.Now(),
		}

		if msg.Sender != "user" {
			t.Error("Expected sender to be set")
		}

		if msg.Content != "Test message" {
			t.Error("Expected content to be set")
		}

		if msg.Timestamp.IsZero() {
			t.Error("Expected timestamp to be set")
		}
	})
}

// TestAIProviderStatusStruct tests the AIProviderStatus struct.
func TestAIProviderStatusStruct(t *testing.T) {
	t.Run("creates status", func(t *testing.T) {
		status := AIProviderStatus{
			Primary:   "Crush",
			Available: true,
			Fallback:  "Mods",
		}

		if status.Primary != "Crush" {
			t.Error("Expected primary to be set")
		}

		if !status.Available {
			t.Error("Expected available to be set")
		}

		if status.Fallback != "Mods" {
			t.Error("Expected fallback to be set")
		}
	})
}

// Benchmark tests
func BenchmarkRenderMentor(b *testing.B) {
	char := game.NewCharacter("BenchPlayer")
	messages := []Message{
		{Sender: "user", Content: "Test question", Timestamp: time.Now()},
		{Sender: "crush", Content: "Test response", Timestamp: time.Now()},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RenderMentor(char, "test input", messages, 80, 24)
	}
}

func BenchmarkRenderConversation(b *testing.B) {
	messages := make([]Message, 10)
	for i := 0; i < 10; i++ {
		messages[i] = Message{
			Sender:    "user",
			Content:   "Benchmark message " + string(rune('A'+i)),
			Timestamp: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderConversation(messages, 80, 20)
	}
}

func BenchmarkWrapText(b *testing.B) {
	text := "This is a very long line of text that needs to be wrapped to fit within the specified width constraint for optimal display in a terminal interface"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wrapText(text, 50)
	}
}
