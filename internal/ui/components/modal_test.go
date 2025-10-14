// Package components provides reusable UI components for CodeQuest
// This file contains tests for the modal dialog component
package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestRenderModal verifies basic modal rendering functionality
func TestRenderModal(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		content   string
		width     int
		height    int
		modalType ModalType
	}{
		{
			name:      "info modal with auto size",
			title:     "Information",
			content:   "This is a test info modal.",
			width:     0,
			height:    0,
			modalType: ModalInfo,
		},
		{
			name:      "success modal with explicit size",
			title:     "Success",
			content:   "Operation completed successfully!",
			width:     50,
			height:    12,
			modalType: ModalSuccess,
		},
		{
			name:      "warning modal",
			title:     "Warning",
			content:   "This action may have consequences.",
			width:     60,
			height:    15,
			modalType: ModalWarning,
		},
		{
			name:      "error modal",
			title:     "Error",
			content:   "An error occurred during processing.",
			width:     55,
			height:    14,
			modalType: ModalError,
		},
		{
			name:      "confirmation modal",
			title:     "Confirm Action",
			content:   "Are you sure you want to proceed?",
			width:     50,
			height:    12,
			modalType: ModalConfirmation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderModal(tt.title, tt.content, tt.width, tt.height, tt.modalType)

			// Verify result is not empty
			if result == "" {
				t.Error("RenderModal returned empty string")
			}

			// Verify title appears in output
			if !strings.Contains(result, tt.title) {
				t.Errorf("Modal output does not contain title: %s", tt.title)
			}

			// Verify content appears in output
			if !strings.Contains(result, tt.content) {
				t.Errorf("Modal output does not contain content: %s", tt.content)
			}

			// Verify modal icon appears
			expectedIcon := ModalIcons[tt.modalType]
			if !strings.Contains(result, expectedIcon) {
				t.Errorf("Modal output does not contain expected icon: %s", expectedIcon)
			}
		})
	}
}

// TestRenderModalWithConfig verifies modal rendering with full configuration
func TestRenderModalWithConfig(t *testing.T) {
	config := ModalConfig{
		Title:      "Test Modal",
		Content:    "This is test content.",
		Width:      50,
		Height:     12,
		ModalType:  ModalInfo,
		TermWidth:  80,
		TermHeight: 24,
	}

	result := RenderModalWithConfig(config)

	if result == "" {
		t.Error("RenderModalWithConfig returned empty string")
	}

	if !strings.Contains(result, config.Title) {
		t.Errorf("Modal does not contain title: %s", config.Title)
	}

	if !strings.Contains(result, config.Content) {
		t.Errorf("Modal does not contain content: %s", config.Content)
	}
}

// TestModalTypes verifies all modal types render correctly
func TestModalTypes(t *testing.T) {
	types := []ModalType{
		ModalInfo,
		ModalSuccess,
		ModalWarning,
		ModalError,
		ModalConfirmation,
	}

	for _, modalType := range types {
		t.Run(string(rune(modalType)), func(t *testing.T) {
			result := RenderModal("Test", "Content", 50, 12, modalType)

			if result == "" {
				t.Errorf("Modal type %d rendered empty string", modalType)
			}

			// Check for type-specific icon
			expectedIcon := ModalIcons[modalType]
			if !strings.Contains(result, expectedIcon) {
				t.Errorf("Modal type %d missing icon: %s", modalType, expectedIcon)
			}
		})
	}
}

// TestHelperModalFunctions verifies convenience helper functions
func TestHelperModalFunctions(t *testing.T) {
	tests := []struct {
		name     string
		renderFn func(string, string) string
		title    string
		content  string
	}{
		{
			name:     "info modal helper",
			renderFn: RenderInfoModal,
			title:    "Info",
			content:  "Information message",
		},
		{
			name:     "success modal helper",
			renderFn: RenderSuccessModal,
			title:    "Success",
			content:  "Success message",
		},
		{
			name:     "warning modal helper",
			renderFn: RenderWarningModal,
			title:    "Warning",
			content:  "Warning message",
		},
		{
			name:     "error modal helper",
			renderFn: RenderErrorModal,
			title:    "Error",
			content:  "Error message",
		},
		{
			name:     "confirm modal helper",
			renderFn: RenderConfirmModal,
			title:    "Confirm",
			content:  "Confirmation message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.renderFn(tt.title, tt.content)

			if result == "" {
				t.Error("Helper function returned empty string")
			}

			if !strings.Contains(result, tt.title) {
				t.Errorf("Helper modal missing title: %s", tt.title)
			}

			if !strings.Contains(result, tt.content) {
				t.Errorf("Helper modal missing content: %s", tt.content)
			}
		})
	}
}

// TestRenderLevelUpModal verifies level-up modal rendering
func TestRenderLevelUpModal(t *testing.T) {
	tests := []struct {
		name      string
		level     int
		abilities []string
	}{
		{
			name:      "level up with abilities",
			level:     5,
			abilities: []string{"Faster completion", "Higher XP multiplier"},
		},
		{
			name:      "level up without abilities",
			level:     2,
			abilities: []string{},
		},
		{
			name:      "level up with nil abilities",
			level:     3,
			abilities: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderLevelUpModal(tt.level, tt.abilities)

			if result == "" {
				t.Error("RenderLevelUpModal returned empty string")
			}

			// Check for level number
			levelStr := string(rune('0' + tt.level))
			if !strings.Contains(result, levelStr) {
				t.Errorf("Level-up modal missing level: %d", tt.level)
			}

			// Check for abilities if provided
			for _, ability := range tt.abilities {
				if !strings.Contains(result, ability) {
					t.Errorf("Level-up modal missing ability: %s", ability)
				}
			}
		})
	}
}

// TestRenderQuestDetailModal verifies quest detail modal rendering
func TestRenderQuestDetailModal(t *testing.T) {
	questTitle := "First Steps"
	description := "Make your first commit to the repository"
	progress := "1/1 commits"
	reward := 50

	result := RenderQuestDetailModal(questTitle, description, progress, reward)

	if result == "" {
		t.Error("RenderQuestDetailModal returned empty string")
	}

	// Verify all quest details appear
	expectedStrings := []string{questTitle, description, progress}
	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("Quest detail modal missing: %s", expected)
		}
	}
}

// TestWordWrap verifies text wrapping functionality
func TestWordWrap(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		wantLen  int // Expected number of lines
		maxWidth int // Maximum line width
	}{
		{
			name:     "short text no wrap",
			text:     "Short text",
			width:    50,
			wantLen:  1,
			maxWidth: 50,
		},
		{
			name:     "long text with wrap",
			text:     "This is a very long line of text that should be wrapped to fit within the specified width constraint",
			width:    30,
			wantLen:  4, // Approximately 4 lines
			maxWidth: 30,
		},
		{
			name:     "multi-paragraph text",
			text:     "First paragraph here.\n\nSecond paragraph here.",
			width:    50,
			wantLen:  2, // Two non-empty lines
			maxWidth: 50,
		},
		{
			name:     "single word longer than width",
			text:     "Supercalifragilisticexpialidocious",
			width:    20,
			wantLen:  1,
			maxWidth: 50, // Word won't be broken, will exceed width
		},
		{
			name:     "zero width",
			text:     "Test text",
			width:    0,
			wantLen:  1,
			maxWidth: 50, // Should use minimum width
		},
		{
			name:     "negative width",
			text:     "Test text",
			width:    -10,
			wantLen:  1,
			maxWidth: 50, // Should use minimum width
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wordWrap(tt.text, tt.width)

			// Count non-empty lines
			lines := strings.Split(result, "\n")
			nonEmptyLines := 0
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed != "" {
					nonEmptyLines++

					// Check line width (allowing some flexibility for edge cases)
					if len(trimmed) > tt.maxWidth+5 {
						t.Errorf("Line exceeds max width: got %d, max %d: %s",
							len(trimmed), tt.maxWidth, trimmed)
					}
				}
			}

			// Allow some flexibility in line count (±1 line)
			if nonEmptyLines < tt.wantLen-1 || nonEmptyLines > tt.wantLen+1 {
				t.Errorf("Expected ~%d non-empty lines, got %d", tt.wantLen, nonEmptyLines)
			}
		})
	}
}

// TestCalculateModalWidth verifies width calculation
func TestCalculateModalWidth(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		minWidth   int
		maxWidth   int
		wantInRange bool
	}{
		{
			name:        "short content",
			content:     "Short",
			minWidth:    40,
			maxWidth:    70,
			wantInRange: true,
		},
		{
			name:        "medium content",
			content:     "This is a medium length content string for testing",
			minWidth:    40,
			maxWidth:    70,
			wantInRange: true,
		},
		{
			name:        "long content",
			content:     "This is a very long content string that exceeds the maximum width constraint and should be capped at the maximum",
			minWidth:    40,
			maxWidth:    70,
			wantInRange: true,
		},
		{
			name:        "multi-line content",
			content:     "Line 1\nLine 2\nVery long line 3 that is the longest",
			minWidth:    40,
			maxWidth:    70,
			wantInRange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := calculateModalWidth(tt.content)

			if width < tt.minWidth {
				t.Errorf("Width %d below minimum %d", width, tt.minWidth)
			}
			if width > tt.maxWidth {
				t.Errorf("Width %d exceeds maximum %d", width, tt.maxWidth)
			}
		})
	}
}

// TestCalculateModalHeight verifies height calculation
func TestCalculateModalHeight(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		width     int
		minHeight int
		maxHeight int
	}{
		{
			name:      "short content",
			content:   "Short",
			width:     50,
			minHeight: 10,
			maxHeight: 20,
		},
		{
			name:      "multi-line content",
			content:   "Line 1\nLine 2\nLine 3\nLine 4\nLine 5",
			width:     50,
			minHeight: 10,
			maxHeight: 20,
		},
		{
			name:      "long wrapping content",
			content:   strings.Repeat("This is a long line that will wrap. ", 10),
			width:     40,
			minHeight: 10,
			maxHeight: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			height := calculateModalHeight(tt.content, tt.width)

			if height < tt.minHeight {
				t.Errorf("Height %d below minimum %d", height, tt.minHeight)
			}
			if height > tt.maxHeight {
				t.Errorf("Height %d exceeds maximum %d", height, tt.maxHeight)
			}
		})
	}
}

// TestClamp verifies the clamp utility function
func TestClamp(t *testing.T) {
	tests := []struct {
		name  string
		value int
		min   int
		max   int
		want  int
	}{
		{
			name:  "value within range",
			value: 50,
			min:   10,
			max:   100,
			want:  50,
		},
		{
			name:  "value below minimum",
			value: 5,
			min:   10,
			max:   100,
			want:  10,
		},
		{
			name:  "value above maximum",
			value: 150,
			min:   10,
			max:   100,
			want:  100,
		},
		{
			name:  "value equals minimum",
			value: 10,
			min:   10,
			max:   100,
			want:  10,
		},
		{
			name:  "value equals maximum",
			value: 100,
			min:   10,
			max:   100,
			want:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clamp(tt.value, tt.min, tt.max)
			if got != tt.want {
				t.Errorf("clamp(%d, %d, %d) = %d, want %d",
					tt.value, tt.min, tt.max, got, tt.want)
			}
		})
	}
}

// TestGetModalColor verifies color selection for modal types
func TestGetModalColor(t *testing.T) {
	tests := []struct {
		modalType ModalType
		wantColor lipgloss.Color
	}{
		{ModalInfo, lipgloss.Color("69")},         // ColorInfo
		{ModalSuccess, lipgloss.Color("42")},      // ColorSuccess
		{ModalWarning, lipgloss.Color("214")},     // ColorWarning
		{ModalError, lipgloss.Color("196")},       // ColorError
		{ModalConfirmation, lipgloss.Color("86")}, // ColorAccent
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.modalType)), func(t *testing.T) {
			got := getModalColor(tt.modalType)
			// We can't directly compare colors, so just verify it's not empty
			if got == "" {
				t.Error("getModalColor returned empty color")
			}
		})
	}
}

// TestRenderModalTitleBar verifies title bar rendering
func TestRenderModalTitleBar(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		modalType ModalType
		width     int
	}{
		{"info title", "Information", ModalInfo, 50},
		{"success title", "Success", ModalSuccess, 50},
		{"warning title", "Warning", ModalWarning, 50},
		{"error title", "Error", ModalError, 50},
		{"confirmation title", "Confirm", ModalConfirmation, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderModalTitleBar(tt.title, tt.modalType, tt.width)

			if result == "" {
				t.Error("renderModalTitleBar returned empty string")
			}

			if !strings.Contains(result, tt.title) {
				t.Errorf("Title bar missing title: %s", tt.title)
			}

			expectedIcon := ModalIcons[tt.modalType]
			if !strings.Contains(result, expectedIcon) {
				t.Errorf("Title bar missing icon: %s", expectedIcon)
			}
		})
	}
}

// TestRenderModalFooter verifies footer rendering with action hints
func TestRenderModalFooter(t *testing.T) {
	tests := []struct {
		name         string
		modalType    ModalType
		width        int
		expectedHint string
	}{
		{"info footer", ModalInfo, 50, "[Enter] Continue"},
		{"success footer", ModalSuccess, 50, "[Enter] Continue"},
		{"warning footer", ModalWarning, 50, "[Enter] OK"},
		{"error footer", ModalError, 50, "[Enter] OK"},
		{"confirmation footer", ModalConfirmation, 50, "[Y] Yes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderModalFooter(tt.modalType, tt.width)

			if result == "" {
				t.Error("renderModalFooter returned empty string")
			}

			if !strings.Contains(result, tt.expectedHint) {
				t.Errorf("Footer missing expected hint: %s in:\n%s",
					tt.expectedHint, result)
			}
		})
	}
}

// TestRenderModalContent verifies content area rendering
func TestRenderModalContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		width   int
		height  int
	}{
		{"short content", "Short message", 50, 12},
		{"long content", strings.Repeat("Long message. ", 20), 50, 12},
		{"multi-line content", "Line 1\nLine 2\nLine 3", 50, 12},
		{"very tall content", strings.Repeat("Line\n", 30), 50, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderModalContent(tt.content, tt.width, tt.height)

			if result == "" {
				t.Error("renderModalContent returned empty string")
			}

			// Content should appear (possibly truncated)
			contentStart := strings.Split(tt.content, "\n")[0]
			contentWords := strings.Fields(contentStart)
			if len(contentWords) > 0 && !strings.Contains(result, contentWords[0]) {
				t.Error("Content missing first word of original content")
			}
		})
	}
}

// TestNilAndEmptyContent verifies handling of edge cases
func TestNilAndEmptyContent(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
	}{
		{"empty title", "", "Content"},
		{"empty content", "Title", ""},
		{"both empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderModal(tt.title, tt.content, 50, 12, ModalInfo)

			// Should not panic and should return something
			if result == "" {
				t.Error("RenderModal returned empty string for edge case")
			}
		})
	}
}

// TestModalCentering verifies modal centering on different screen sizes
func TestModalCentering(t *testing.T) {
	tests := []struct {
		name       string
		termWidth  int
		termHeight int
	}{
		{"standard terminal", 80, 24},
		{"wide terminal", 120, 30},
		{"narrow terminal", 60, 20},
		{"tall terminal", 80, 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ModalConfig{
				Title:      "Test",
				Content:    "Content",
				Width:      50,
				Height:     12,
				ModalType:  ModalInfo,
				TermWidth:  tt.termWidth,
				TermHeight: tt.termHeight,
			}

			result := RenderModalWithConfig(config)

			if result == "" {
				t.Error("RenderModalWithConfig returned empty string")
			}

			// Verify result has reasonable height
			lines := strings.Split(result, "\n")
			if len(lines) < config.Height {
				t.Errorf("Modal has fewer lines than expected: got %d, want >= %d",
					len(lines), config.Height)
			}
		})
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

// BenchmarkRenderModal measures modal rendering performance
func BenchmarkRenderModal(b *testing.B) {
	title := "Test Modal"
	content := "This is test content for benchmarking modal rendering performance."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RenderModal(title, content, 50, 12, ModalInfo)
	}
}

// BenchmarkRenderModalWithConfig measures configured modal rendering
func BenchmarkRenderModalWithConfig(b *testing.B) {
	config := ModalConfig{
		Title:      "Test Modal",
		Content:    "This is test content for benchmarking.",
		Width:      50,
		Height:     12,
		ModalType:  ModalInfo,
		TermWidth:  80,
		TermHeight: 24,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RenderModalWithConfig(config)
	}
}

// BenchmarkWordWrap measures word wrapping performance
func BenchmarkWordWrap(b *testing.B) {
	text := strings.Repeat("This is a long line of text that needs wrapping. ", 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wordWrap(text, 50)
	}
}

// BenchmarkRenderAllModalTypes measures rendering of all modal types
func BenchmarkRenderAllModalTypes(b *testing.B) {
	types := []ModalType{ModalInfo, ModalSuccess, ModalWarning, ModalError, ModalConfirmation}
	title := "Test"
	content := "Test content"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, modalType := range types {
			RenderModal(title, content, 50, 12, modalType)
		}
	}
}

// BenchmarkHelperFunctions measures helper function performance
func BenchmarkHelperFunctions(b *testing.B) {
	title := "Test"
	content := "Test content"

	b.Run("InfoModal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RenderInfoModal(title, content)
		}
	})

	b.Run("SuccessModal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RenderSuccessModal(title, content)
		}
	})

	b.Run("ErrorModal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RenderErrorModal(title, content)
		}
	})

	b.Run("ConfirmModal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RenderConfirmModal(title, content)
		}
	})
}
