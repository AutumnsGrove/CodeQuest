// Package components provides reusable UI components for CodeQuest screens.
// This file contains tests for the statbar component.
package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// TestRenderStatBar verifies the stat bar renders correctly with different inputs.
func TestRenderStatBar(t *testing.T) {
	tests := []struct {
		name      string
		char      *game.Character
		width     int
		wantEmpty bool
	}{
		{
			name:      "normal character with default width",
			char:      game.NewCharacter("TestHero"),
			width:     80,
			wantEmpty: false,
		},
		{
			name:      "nil character",
			char:      nil,
			width:     80,
			wantEmpty: false, // Should render error message
		},
		{
			name:      "narrow width",
			char:      game.NewCharacter("TestHero"),
			width:     50,
			wantEmpty: false,
		},
		{
			name:      "minimum width",
			char:      game.NewCharacter("TestHero"),
			width:     40,
			wantEmpty: false,
		},
		{
			name:      "below minimum width (should enforce minimum)",
			char:      game.NewCharacter("TestHero"),
			width:     20,
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderStatBar(tt.char, tt.width)

			if tt.wantEmpty {
				if got != "" {
					t.Errorf("RenderStatBar() should be empty, got %q", got)
				}
				return
			}

			// Verify we got some output
			if got == "" {
				t.Error("RenderStatBar() returned empty string")
			}

			// If character provided, verify character name appears
			if tt.char != nil {
				if !strings.Contains(got, tt.char.Name) {
					t.Errorf("RenderStatBar() missing character name %q", tt.char.Name)
				}

				// Verify level appears
				if !strings.Contains(got, "Lv.") {
					t.Error("RenderStatBar() missing level indicator")
				}

				// Verify XP label appears
				if !strings.Contains(got, "XP:") {
					t.Error("RenderStatBar() missing XP label")
				}
			} else {
				// Nil character should show error message
				if !strings.Contains(got, "No character") && !strings.Contains(got, "unavailable") {
					t.Error("RenderStatBar() with nil character should show error message")
				}
			}
		})
	}
}

// TestRenderStatBarWithConfig verifies custom configuration options work correctly.
func TestRenderStatBarWithConfig(t *testing.T) {
	char := game.NewCharacter("ConfigTest")
	char.CodePower = 15
	char.Wisdom = 12
	char.Agility = 10

	tests := []struct {
		name   string
		config StatBarConfig
		verify func(t *testing.T, output string)
	}{
		{
			name: "default config",
			config: StatBarConfig{
				Width:           80,
				ShowSessionInfo: true,
				ShowStreak:      true,
				Compact:         false,
			},
			verify: func(t *testing.T, output string) {
				// Should include all sections
				if !strings.Contains(output, "ConfigTest") {
					t.Error("Missing character name")
				}
				if !strings.Contains(output, "Streak:") {
					t.Error("Missing streak section")
				}
				if !strings.Contains(output, "Today:") {
					t.Error("Missing session info section")
				}
			},
		},
		{
			name: "no session info",
			config: StatBarConfig{
				Width:           80,
				ShowSessionInfo: false,
				ShowStreak:      true,
				Compact:         false,
			},
			verify: func(t *testing.T, output string) {
				// Should NOT include session info
				if strings.Contains(output, "Today:") {
					t.Error("Should not show session info when disabled")
				}
				// Should still include streak
				if !strings.Contains(output, "Streak:") {
					t.Error("Missing streak section")
				}
			},
		},
		{
			name: "no streak",
			config: StatBarConfig{
				Width:           80,
				ShowSessionInfo: true,
				ShowStreak:      false,
				Compact:         false,
			},
			verify: func(t *testing.T, output string) {
				// Should NOT include streak
				if strings.Contains(output, "Streak:") {
					t.Error("Should not show streak when disabled")
				}
				// Should still include session info
				if !strings.Contains(output, "Today:") {
					t.Error("Missing session info section")
				}
			},
		},
		{
			name: "minimal config",
			config: StatBarConfig{
				Width:           80,
				ShowSessionInfo: false,
				ShowStreak:      false,
				Compact:         true,
			},
			verify: func(t *testing.T, output string) {
				// Should only show basic info
				if !strings.Contains(output, "ConfigTest") {
					t.Error("Missing character name")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderStatBarWithConfig(char, tt.config)

			if got == "" {
				t.Error("RenderStatBarWithConfig() returned empty string")
			}

			tt.verify(t, got)
		})
	}
}

// TestRenderHeaderSection verifies the header section renders correctly.
func TestRenderHeaderSection(t *testing.T) {
	char := game.NewCharacter("HeaderTest")
	char.Level = 5
	char.XP = 50
	char.XPToNextLevel = 100

	got := renderHeaderSection(char, 80)

	// Should not be empty
	if got == "" {
		t.Error("renderHeaderSection() returned empty string")
	}

	// Should contain character name
	if !strings.Contains(got, "HeaderTest") {
		t.Error("renderHeaderSection() missing character name")
	}

	// Should contain level
	if !strings.Contains(got, "Lv.5") {
		t.Error("renderHeaderSection() missing correct level")
	}

	// Should contain XP label
	if !strings.Contains(got, "XP:") {
		t.Error("renderHeaderSection() missing XP label")
	}
}

// TestRenderRPGStats verifies RPG stats (CodePower, Wisdom, Agility) render correctly.
func TestRenderRPGStats(t *testing.T) {
	char := game.NewCharacter("StatsTest")
	char.CodePower = 15
	char.Wisdom = 12
	char.Agility = 10

	got := renderRPGStats(char, 80)

	// Should not be empty
	if got == "" {
		t.Error("renderRPGStats() returned empty string")
	}

	// Should contain CodePower
	if !strings.Contains(got, "CodePower:") {
		t.Error("renderRPGStats() missing CodePower label")
	}
	if !strings.Contains(got, "15") {
		t.Error("renderRPGStats() missing CodePower value")
	}

	// Should contain Wisdom
	if !strings.Contains(got, "Wisdom:") {
		t.Error("renderRPGStats() missing Wisdom label")
	}
	if !strings.Contains(got, "12") {
		t.Error("renderRPGStats() missing Wisdom value")
	}

	// Should contain Agility
	if !strings.Contains(got, "Agility:") {
		t.Error("renderRPGStats() missing Agility label")
	}
	if !strings.Contains(got, "10") {
		t.Error("renderRPGStats() missing Agility value")
	}

	// Should contain stat icons
	if !strings.Contains(got, "⚔️") {
		t.Error("renderRPGStats() missing CodePower icon")
	}
	if !strings.Contains(got, "📖") {
		t.Error("renderRPGStats() missing Wisdom icon")
	}
	if !strings.Contains(got, "⚡") {
		t.Error("renderRPGStats() missing Agility icon")
	}
}

// TestRenderStreakSection verifies streak counter renders correctly.
func TestRenderStreakSection(t *testing.T) {
	tests := []struct {
		name          string
		currentStreak int
		longestStreak int
		verifyStreak  string
		verifyBest    string
	}{
		{
			name:          "no streak",
			currentStreak: 0,
			longestStreak: 0,
			verifyStreak:  "0 days",
			verifyBest:    "0 days",
		},
		{
			name:          "active streak",
			currentStreak: 5,
			longestStreak: 10,
			verifyStreak:  "5 days",
			verifyBest:    "10 days",
		},
		{
			name:          "at personal best",
			currentStreak: 15,
			longestStreak: 15,
			verifyStreak:  "15 days",
			verifyBest:    "15 days",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := game.NewCharacter("StreakTest")
			char.CurrentStreak = tt.currentStreak
			char.LongestStreak = tt.longestStreak

			got := renderStreakSection(char, 80)

			// Should not be empty
			if got == "" {
				t.Error("renderStreakSection() returned empty string")
			}

			// Verify streak label
			if !strings.Contains(got, "Streak:") {
				t.Error("renderStreakSection() missing Streak label")
			}

			// Verify best label
			if !strings.Contains(got, "Best:") {
				t.Error("renderStreakSection() missing Best label")
			}

			// Verify values
			if !strings.Contains(got, tt.verifyStreak) {
				t.Errorf("renderStreakSection() missing current streak %q", tt.verifyStreak)
			}
			if !strings.Contains(got, tt.verifyBest) {
				t.Errorf("renderStreakSection() missing best streak %q", tt.verifyBest)
			}

			// Verify icons
			if !strings.Contains(got, "🔥") {
				t.Error("renderStreakSection() missing fire icon")
			}
			if !strings.Contains(got, "🏆") {
				t.Error("renderStreakSection() missing trophy icon")
			}
		})
	}
}

// TestRenderSessionStats verifies today's session stats render correctly.
func TestRenderSessionStats(t *testing.T) {
	tests := []struct {
		name            string
		todayCommits    int
		todayLinesAdded int
		verifyCommits   string
		verifyLines     string
	}{
		{
			name:            "no activity",
			todayCommits:    0,
			todayLinesAdded: 0,
			verifyCommits:   "0 commits",
			verifyLines:     "+0",
		},
		{
			name:            "some activity",
			todayCommits:    3,
			todayLinesAdded: 150,
			verifyCommits:   "3 commits",
			verifyLines:     "+150",
		},
		{
			name:            "high activity",
			todayCommits:    10,
			todayLinesAdded: 500,
			verifyCommits:   "10 commits",
			verifyLines:     "+500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := game.NewCharacter("SessionTest")
			char.TodayCommits = tt.todayCommits
			char.TodayLinesAdded = tt.todayLinesAdded

			got := renderSessionStats(char, 80)

			// Should not be empty
			if got == "" {
				t.Error("renderSessionStats() returned empty string")
			}

			// Verify labels
			if !strings.Contains(got, "Today:") {
				t.Error("renderSessionStats() missing Today label")
			}
			if !strings.Contains(got, "Lines:") {
				t.Error("renderSessionStats() missing Lines label")
			}

			// Verify values
			if !strings.Contains(got, tt.verifyCommits) {
				t.Errorf("renderSessionStats() missing commits %q", tt.verifyCommits)
			}
			if !strings.Contains(got, tt.verifyLines) {
				t.Errorf("renderSessionStats() missing lines %q", tt.verifyLines)
			}

			// Verify icons
			if !strings.Contains(got, "💾") {
				t.Error("renderSessionStats() missing commit icon")
			}
			if !strings.Contains(got, "📝") {
				t.Error("renderSessionStats() missing code icon")
			}
		})
	}
}

// TestRenderCompactStatBar verifies compact mode renders correctly.
func TestRenderCompactStatBar(t *testing.T) {
	tests := []struct {
		name  string
		char  *game.Character
		width int
	}{
		{
			name:  "normal character",
			char:  game.NewCharacter("CompactTest"),
			width: 60,
		},
		{
			name:  "nil character",
			char:  nil,
			width: 60,
		},
		{
			name:  "narrow width",
			char:  game.NewCharacter("CompactTest"),
			width: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderCompactStatBar(tt.char, tt.width)

			// Should not be empty
			if got == "" {
				t.Error("RenderCompactStatBar() returned empty string")
			}

			if tt.char != nil {
				// Should include character name
				if !strings.Contains(got, tt.char.Name) {
					t.Error("RenderCompactStatBar() missing character name")
				}

				// Should include streak (compact mode still shows this)
				if !strings.Contains(got, "Streak:") {
					t.Error("RenderCompactStatBar() missing streak")
				}

				// Should NOT include session info (disabled in compact mode)
				if strings.Contains(got, "Today:") {
					t.Error("RenderCompactStatBar() should not show session info in compact mode")
				}
			} else {
				// Nil character should show error
				if !strings.Contains(got, "No character") && !strings.Contains(got, "unavailable") {
					t.Error("RenderCompactStatBar() with nil character should show error")
				}
			}
		})
	}
}

// TestRenderStatBadge verifies stat badge renders correctly.
func TestRenderStatBadge(t *testing.T) {
	tests := []struct {
		name string
		char *game.Character
	}{
		{
			name: "normal character",
			char: game.NewCharacter("BadgeTest"),
		},
		{
			name: "nil character",
			char: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderStatBadge(tt.char)

			// Should not be empty
			if got == "" {
				t.Error("RenderStatBadge() returned empty string")
			}

			if tt.char != nil {
				// Should include level
				if !strings.Contains(got, "Lv.") {
					t.Error("RenderStatBadge() missing level indicator")
				}

				// Should include XP info
				if !strings.Contains(got, "XP:") {
					t.Error("RenderStatBadge() missing XP label")
				}

				// Should include streak
				if !strings.Contains(got, "🔥") {
					t.Error("RenderStatBadge() missing streak icon")
				}
			} else {
				// Nil character should show placeholder
				if !strings.Contains(got, "No character") {
					t.Error("RenderStatBadge() with nil character should show placeholder")
				}
			}
		})
	}
}

// TestRenderNilCharacterStatBar verifies nil character handling.
func TestRenderNilCharacterStatBar(t *testing.T) {
	tests := []struct {
		name  string
		width int
	}{
		{
			name:  "wide width",
			width: 80,
		},
		{
			name:  "narrow width",
			width: 40,
		},
		{
			name:  "very narrow width",
			width: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderNilCharacterStatBar(tt.width)

			// Should not be empty
			if got == "" {
				t.Error("renderNilCharacterStatBar() returned empty string")
			}

			// Should contain error indicator
			if !strings.Contains(got, "⚠️") {
				t.Error("renderNilCharacterStatBar() missing warning icon")
			}

			// Should contain helpful message
			if !strings.Contains(got, "No character") && !strings.Contains(got, "unavailable") {
				t.Error("renderNilCharacterStatBar() missing error message")
			}
		})
	}
}

// TestStatBarWidth verifies the stat bar respects width constraints.
func TestStatBarWidth(t *testing.T) {
	char := game.NewCharacter("WidthTest")

	tests := []struct {
		name  string
		width int
	}{
		{"80 columns", 80},
		{"120 columns", 120},
		{"60 columns", 60},
		{"40 columns", 40},
		{"30 columns", 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderStatBar(char, tt.width)

			// Measure the actual width (accounting for ANSI codes)
			lines := strings.Split(got, "\n")
			for i, line := range lines {
				actualWidth := lipgloss.Width(line)

				// The actual width should be reasonable relative to requested width
				// Allow generous margin for borders, padding, and ANSI codes
				if actualWidth > tt.width+20 {
					t.Errorf("Line %d width %d exceeds requested %d by too much (allowing 20 char margin)",
						i, actualWidth, tt.width)
				}
			}
		})
	}
}

// TestDefaultStatBarConfig verifies default configuration is sensible.
func TestDefaultStatBarConfig(t *testing.T) {
	config := DefaultStatBarConfig()

	if config.Width != 80 {
		t.Errorf("DefaultStatBarConfig().Width = %d, want 80", config.Width)
	}

	if !config.ShowSessionInfo {
		t.Error("DefaultStatBarConfig().ShowSessionInfo should be true")
	}

	if !config.ShowStreak {
		t.Error("DefaultStatBarConfig().ShowStreak should be true")
	}

	if config.Compact {
		t.Error("DefaultStatBarConfig().Compact should be false")
	}
}

// BenchmarkRenderStatBar benchmarks the stat bar rendering performance.
func BenchmarkRenderStatBar(b *testing.B) {
	char := game.NewCharacter("BenchHero")
	width := 80

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderStatBar(char, width)
	}
}

// BenchmarkRenderCompactStatBar benchmarks compact stat bar rendering.
func BenchmarkRenderCompactStatBar(b *testing.B) {
	char := game.NewCharacter("BenchHero")
	width := 60

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderCompactStatBar(char, width)
	}
}

// BenchmarkRenderStatBadge benchmarks stat badge rendering.
func BenchmarkRenderStatBadge(b *testing.B) {
	char := game.NewCharacter("BenchHero")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderStatBadge(char)
	}
}
