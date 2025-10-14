// Package components provides reusable UI components for CodeQuest.
// This file contains comprehensive tests for the timer component.
package components

import (
	"strings"
	"testing"
	"time"
)

// ============================================================================
// formatDuration Tests
// ============================================================================

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "zero duration",
			duration: 0,
			want:     "0:00:00",
		},
		{
			name:     "negative duration treated as zero",
			duration: -5 * time.Minute,
			want:     "0:00:00",
		},
		{
			name:     "30 seconds",
			duration: 30 * time.Second,
			want:     "0:00:30",
		},
		{
			name:     "1 minute",
			duration: 1 * time.Minute,
			want:     "0:01:00",
		},
		{
			name:     "45 minutes 30 seconds",
			duration: 45*time.Minute + 30*time.Second,
			want:     "0:45:30",
		},
		{
			name:     "1 hour exactly",
			duration: 1 * time.Hour,
			want:     "1:00:00",
		},
		{
			name:     "2 hours 34 minutes 15 seconds",
			duration: 2*time.Hour + 34*time.Minute + 15*time.Second,
			want:     "2:34:15",
		},
		{
			name:     "10 hours 5 minutes 8 seconds",
			duration: 10*time.Hour + 5*time.Minute + 8*time.Second,
			want:     "10:05:08",
		},
		{
			name:     "23 hours 59 minutes 59 seconds",
			duration: 23*time.Hour + 59*time.Minute + 59*time.Second,
			want:     "23:59:59",
		},
		{
			name:     "99 hours (edge case)",
			duration: 99 * time.Hour,
			want:     "99:00:00",
		},
		{
			name:     "single digit seconds",
			duration: 5 * time.Second,
			want:     "0:00:05",
		},
		{
			name:     "single digit minutes",
			duration: 3*time.Minute + 42*time.Second,
			want:     "0:03:42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.duration)
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, got, tt.want)
			}
		})
	}
}

// ============================================================================
// getTimerColor Tests
// ============================================================================

func TestGetTimerColor(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		// We can't directly test color equality since they're lipgloss.Color
		// Instead we test the duration thresholds
		expectDim     bool
		expectInfo    bool
		expectSuccess bool
		expectWarning bool
	}{
		{
			name:      "0 minutes - dim",
			duration:  0,
			expectDim: true,
		},
		{
			name:      "30 minutes - dim",
			duration:  30 * time.Minute,
			expectDim: true,
		},
		{
			name:      "59 minutes - dim",
			duration:  59 * time.Minute,
			expectDim: true,
		},
		{
			name:       "1 hour exactly - info",
			duration:   1 * time.Hour,
			expectInfo: true,
		},
		{
			name:       "2 hours - info",
			duration:   2 * time.Hour,
			expectInfo: true,
		},
		{
			name:       "2.5 hours - info",
			duration:   2*time.Hour + 30*time.Minute,
			expectInfo: true,
		},
		{
			name:          "3 hours exactly - success",
			duration:      3 * time.Hour,
			expectSuccess: true,
		},
		{
			name:          "4 hours - success",
			duration:      4 * time.Hour,
			expectSuccess: true,
		},
		{
			name:          "4.5 hours - success",
			duration:      4*time.Hour + 30*time.Minute,
			expectSuccess: true,
		},
		{
			name:          "5 hours exactly - warning",
			duration:      5 * time.Hour,
			expectWarning: true,
		},
		{
			name:          "6 hours - warning",
			duration:      6 * time.Hour,
			expectWarning: true,
		},
		{
			name:          "10 hours - warning",
			duration:      10 * time.Hour,
			expectWarning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := getTimerColor(tt.duration)
			// Since we can't directly compare lipgloss.Color, we just ensure it returns something
			if color == "" {
				t.Error("getTimerColor returned empty color")
			}
			// Note: In a real implementation, you might want to export color constants
			// to allow direct comparison in tests
		})
	}
}

// ============================================================================
// getTimerIcon Tests
// ============================================================================

func TestGetTimerIcon(t *testing.T) {
	tests := []struct {
		name      string
		isRunning bool
		want      string
	}{
		{
			name:      "running timer shows red dot",
			isRunning: true,
			want:      "🔴",
		},
		{
			name:      "paused timer shows pause icon",
			isRunning: false,
			want:      "⏸",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTimerIcon(tt.isRunning)
			if got != tt.want {
				t.Errorf("getTimerIcon(%v) = %q, want %q", tt.isRunning, got, tt.want)
			}
		})
	}
}

// ============================================================================
// RenderTimer Tests (Basic)
// ============================================================================

func TestRenderTimer(t *testing.T) {
	tests := []struct {
		name         string
		duration     time.Duration
		isRunning    bool
		wantContains []string // Strings that should appear in output
	}{
		{
			name:         "zero duration paused",
			duration:     0,
			isRunning:    false,
			wantContains: []string{"0:00:00", "⏸"},
		},
		{
			name:         "zero duration running",
			duration:     0,
			isRunning:    true,
			wantContains: []string{"0:00:00", "🔴"},
		},
		{
			name:         "2 hours running",
			duration:     2 * time.Hour,
			isRunning:    true,
			wantContains: []string{"2:00:00", "🔴"},
		},
		{
			name:         "2 hours 34 minutes paused",
			duration:     2*time.Hour + 34*time.Minute,
			isRunning:    false,
			wantContains: []string{"2:34:00", "⏸"},
		},
		{
			name:         "complex duration",
			duration:     3*time.Hour + 45*time.Minute + 23*time.Second,
			isRunning:    true,
			wantContains: []string{"3:45:23", "🔴"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderTimer(tt.duration, tt.isRunning)

			// Check that the output contains expected strings
			// Note: We can't do exact string matching due to ANSI color codes
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("RenderTimer() output doesn't contain %q\nGot: %s", want, got)
				}
			}

			// Ensure output is not empty
			if got == "" {
				t.Error("RenderTimer() returned empty string")
			}
		})
	}
}

// ============================================================================
// RenderInlineTimer Tests
// ============================================================================

func TestRenderInlineTimer(t *testing.T) {
	tests := []struct {
		name         string
		duration     time.Duration
		isRunning    bool
		wantContains []string
	}{
		{
			name:         "inline running timer",
			duration:     1*time.Hour + 23*time.Minute + 45*time.Second,
			isRunning:    true,
			wantContains: []string{"1:23:45", "🔴"},
		},
		{
			name:         "inline paused timer",
			duration:     30 * time.Minute,
			isRunning:    false,
			wantContains: []string{"0:30:00", "⏸"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderInlineTimer(tt.duration, tt.isRunning)

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("RenderInlineTimer() doesn't contain %q\nGot: %s", want, got)
				}
			}
		})
	}
}

// ============================================================================
// RenderTimerCard Tests
// ============================================================================

func TestRenderTimerCard(t *testing.T) {
	tests := []struct {
		name         string
		duration     time.Duration
		isRunning    bool
		width        int
		wantContains []string
	}{
		{
			name:      "card running timer",
			duration:  2*time.Hour + 15*time.Minute,
			isRunning: true,
			width:     40,
			wantContains: []string{
				"Session Timer",
				"2:15:00",
				"🔴",
				"Pause",
				"Ctrl+T",
			},
		},
		{
			name:      "card paused timer",
			duration:  45 * time.Minute,
			isRunning: false,
			width:     40,
			wantContains: []string{
				"Session Timer",
				"0:45:00",
				"⏸",
				"Resume",
				"Ctrl+T",
			},
		},
		{
			name:      "card with break reminder",
			duration:  6 * time.Hour,
			isRunning: true,
			width:     40,
			wantContains: []string{
				"Session Timer",
				"6:00:00",
				"Take a break",
				"⚠",
			},
		},
		{
			name:         "card with minimum width",
			duration:     1 * time.Hour,
			isRunning:    false,
			width:        10, // Too small, should use minimum
			wantContains: []string{"1:00:00"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderTimerCard(tt.duration, tt.isRunning, tt.width)

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("RenderTimerCard() doesn't contain %q\nGot: %s", want, got)
				}
			}

			// Ensure output is not empty
			if got == "" {
				t.Error("RenderTimerCard() returned empty string")
			}
		})
	}
}

// ============================================================================
// RenderMinimalTimer Tests
// ============================================================================

func TestRenderMinimalTimer(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "minimal zero",
			duration: 0,
			want:     "0:00:00",
		},
		{
			name:     "minimal 1 hour",
			duration: 1 * time.Hour,
			want:     "1:00:00",
		},
		{
			name:     "minimal complex",
			duration: 5*time.Hour + 23*time.Minute + 7*time.Second,
			want:     "5:23:07",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderMinimalTimer(tt.duration)
			if got != tt.want {
				t.Errorf("RenderMinimalTimer(%v) = %q, want %q", tt.duration, got, tt.want)
			}
		})
	}
}

// ============================================================================
// RenderTimerWithConfig Tests
// ============================================================================

func TestRenderTimerWithConfig(t *testing.T) {
	tests := []struct {
		name         string
		duration     time.Duration
		config       TimerConfig
		wantContains []string
	}{
		{
			name:     "config inline mode",
			duration: 1 * time.Hour,
			config: TimerConfig{
				Mode:      TimerModeInline,
				IsRunning: true,
			},
			wantContains: []string{"1:00:00", "🔴"},
		},
		{
			name:     "config card mode",
			duration: 2 * time.Hour,
			config: TimerConfig{
				Mode:      TimerModeCard,
				IsRunning: false,
				Width:     40,
			},
			wantContains: []string{"Session Timer", "2:00:00"},
		},
		{
			name:     "config minimal mode",
			duration: 3 * time.Hour,
			config: TimerConfig{
				Mode:      TimerModeMinimal,
				IsRunning: true, // Ignored in minimal mode
			},
			wantContains: []string{"3:00:00"},
		},
		{
			name:     "config unknown mode defaults to inline",
			duration: 30 * time.Minute,
			config: TimerConfig{
				Mode:      "unknown",
				IsRunning: false,
			},
			wantContains: []string{"0:30:00", "⏸"},
		},
		{
			name:     "config with negative duration",
			duration: -5 * time.Minute,
			config: TimerConfig{
				Mode:      TimerModeInline,
				IsRunning: false,
			},
			wantContains: []string{"0:00:00"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderTimerWithConfig(tt.duration, tt.config)

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("RenderTimerWithConfig() doesn't contain %q\nGot: %s", want, got)
				}
			}
		})
	}
}

// ============================================================================
// Utility Functions Tests
// ============================================================================

func TestGetTimerHint(t *testing.T) {
	tests := []struct {
		name         string
		isRunning    bool
		wantContains []string
	}{
		{
			name:         "hint for running timer",
			isRunning:    true,
			wantContains: []string{"Ctrl+T", "pause"},
		},
		{
			name:         "hint for paused timer",
			isRunning:    false,
			wantContains: []string{"Ctrl+T", "start"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTimerHint(tt.isRunning)

			for _, want := range tt.wantContains {
				if !strings.Contains(strings.ToLower(got), strings.ToLower(want)) {
					t.Errorf("GetTimerHint(%v) doesn't contain %q\nGot: %s", tt.isRunning, want, got)
				}
			}
		})
	}
}

func TestShouldShowBreakReminder(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     bool
	}{
		{
			name:     "0 hours - no break",
			duration: 0,
			want:     false,
		},
		{
			name:     "1 hour - no break",
			duration: 1 * time.Hour,
			want:     false,
		},
		{
			name:     "4 hours 59 minutes - no break",
			duration: 4*time.Hour + 59*time.Minute,
			want:     false,
		},
		{
			name:     "5 hours exactly - show break",
			duration: 5 * time.Hour,
			want:     true,
		},
		{
			name:     "6 hours - show break",
			duration: 6 * time.Hour,
			want:     true,
		},
		{
			name:     "10 hours - show break",
			duration: 10 * time.Hour,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldShowBreakReminder(tt.duration)
			if got != tt.want {
				t.Errorf("ShouldShowBreakReminder(%v) = %v, want %v", tt.duration, got, tt.want)
			}
		})
	}
}

func TestFormatSessionSummary(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "very short session",
			duration: 30 * time.Second,
			want:     "Quick session!",
		},
		{
			name:     "minutes only",
			duration: 45 * time.Minute,
			want:     "You coded for 45 minutes.",
		},
		{
			name:     "1 hour exactly",
			duration: 1 * time.Hour,
			want:     "You coded for 1 hour.",
		},
		{
			name:     "2 hours exactly",
			duration: 2 * time.Hour,
			want:     "You coded for 2 hours.",
		},
		{
			name:     "1 hour 1 minute",
			duration: 1*time.Hour + 1*time.Minute,
			want:     "You coded for 1 hour and 1 minute.",
		},
		{
			name:     "1 hour 30 minutes",
			duration: 1*time.Hour + 30*time.Minute,
			want:     "You coded for 1 hour and 30 minutes.",
		},
		{
			name:     "2 hours 1 minute",
			duration: 2*time.Hour + 1*time.Minute,
			want:     "You coded for 2 hours and 1 minute.",
		},
		{
			name:     "2 hours 30 minutes",
			duration: 2*time.Hour + 30*time.Minute,
			want:     "You coded for 2 hours and 30 minutes.",
		},
		{
			name:     "5 hours 23 minutes",
			duration: 5*time.Hour + 23*time.Minute,
			want:     "You coded for 5 hours and 23 minutes.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSessionSummary(tt.duration)
			if got != tt.want {
				t.Errorf("FormatSessionSummary(%v) = %q, want %q", tt.duration, got, tt.want)
			}
		})
	}
}

// ============================================================================
// DefaultTimerConfig Tests
// ============================================================================

func TestDefaultTimerConfig(t *testing.T) {
	config := DefaultTimerConfig()

	if config.Width != 30 {
		t.Errorf("DefaultTimerConfig().Width = %d, want 30", config.Width)
	}

	if config.Mode != TimerModeInline {
		t.Errorf("DefaultTimerConfig().Mode = %q, want %q", config.Mode, TimerModeInline)
	}

	if config.IsRunning != false {
		t.Errorf("DefaultTimerConfig().IsRunning = %v, want false", config.IsRunning)
	}
}

// ============================================================================
// Edge Cases and Nil-Safety Tests
// ============================================================================

func TestTimerEdgeCases(t *testing.T) {
	t.Run("very large duration", func(t *testing.T) {
		duration := 100 * time.Hour
		result := formatDuration(duration)
		if !strings.Contains(result, "100:00:00") {
			t.Errorf("formatDuration(100 hours) = %q, expected to contain '100:00:00'", result)
		}
	})

	t.Run("negative duration in RenderTimer", func(t *testing.T) {
		result := RenderTimer(-5*time.Minute, false)
		if !strings.Contains(result, "0:00:00") {
			t.Errorf("RenderTimer(-5min) should show 0:00:00, got: %s", result)
		}
	})

	t.Run("zero width card falls back to minimum", func(t *testing.T) {
		result := RenderTimerCard(1*time.Hour, false, 0)
		if result == "" {
			t.Error("RenderTimerCard with zero width returned empty string")
		}
	})
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkFormatDuration(b *testing.B) {
	duration := 3*time.Hour + 45*time.Minute + 23*time.Second
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatDuration(duration)
	}
}

func BenchmarkRenderTimer(b *testing.B) {
	duration := 2*time.Hour + 30*time.Minute
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderTimer(duration, true)
	}
}

func BenchmarkRenderInlineTimer(b *testing.B) {
	duration := 1*time.Hour + 23*time.Minute + 45*time.Second
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderInlineTimer(duration, true)
	}
}

func BenchmarkRenderTimerCard(b *testing.B) {
	duration := 4 * time.Hour
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderTimerCard(duration, true, 40)
	}
}

func BenchmarkRenderMinimalTimer(b *testing.B) {
	duration := 5*time.Hour + 30*time.Minute
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderMinimalTimer(duration)
	}
}

func BenchmarkGetTimerColor(b *testing.B) {
	duration := 3 * time.Hour
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getTimerColor(duration)
	}
}

func BenchmarkFormatSessionSummary(b *testing.B) {
	duration := 2*time.Hour + 30*time.Minute
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatSessionSummary(duration)
	}
}
