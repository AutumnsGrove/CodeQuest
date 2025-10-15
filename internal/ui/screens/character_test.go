package screens

import (
	"strings"
	"testing"
	"time"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// TestRenderCharacter tests the main character screen rendering function.
func TestRenderCharacter(t *testing.T) {
	tests := []struct {
		name      string
		character *game.Character
		width     int
		height    int
		wantEmpty bool
		wantErr   bool
	}{
		{
			name:      "normal character wide layout",
			character: createTestCharacter(),
			width:     120,
			height:    40,
			wantEmpty: false,
			wantErr:   false,
		},
		{
			name:      "normal character narrow layout",
			character: createTestCharacter(),
			width:     80,
			height:    40,
			wantEmpty: false,
			wantErr:   false,
		},
		{
			name:      "nil character",
			character: nil,
			width:     80,
			height:    40,
			wantEmpty: false,
			wantErr:   false, // Should render error message, not panic
		},
		{
			name:      "small terminal",
			character: createTestCharacter(),
			width:     40,
			height:    20,
			wantEmpty: false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderCharacter(tt.character, tt.width, tt.height)

			// Should never return empty string
			if result == "" && !tt.wantEmpty {
				t.Error("RenderCharacter() returned empty string")
			}

			// Should contain character name if not nil
			if tt.character != nil {
				if !strings.Contains(result, tt.character.Name) {
					t.Errorf("RenderCharacter() should contain character name %q", tt.character.Name)
				}
			}

			// Should contain "Character Sheet" in header
			if !strings.Contains(result, "Character Sheet") {
				t.Error("RenderCharacter() should contain 'Character Sheet' in header")
			}
		})
	}
}

// TestRenderCharacterWide tests the wide layout rendering.
func TestRenderCharacterWide(t *testing.T) {
	char := createTestCharacter()
	result := renderCharacterWide(char, 120, 40)

	if result == "" {
		t.Error("renderCharacterWide() returned empty string")
	}

	// Should contain character stats
	expectedStrings := []string{
		char.Name,
		"Level",
		"XP",
		"CodePower",
		"Wisdom",
		"Agility",
		"Today's Activity",
		"Lifetime Statistics",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderCharacterWide() should contain %q", expected)
		}
	}
}

// TestRenderCharacterNarrow tests the narrow layout rendering.
func TestRenderCharacterNarrow(t *testing.T) {
	char := createTestCharacter()
	result := renderCharacterNarrow(char, 80, 40)

	if result == "" {
		t.Error("renderCharacterNarrow() returned empty string")
	}

	// Should contain character stats
	if !strings.Contains(result, char.Name) {
		t.Errorf("renderCharacterNarrow() should contain character name %q", char.Name)
	}
}

// TestRenderStatsPanel tests the stats panel rendering.
func TestRenderStatsPanel(t *testing.T) {
	char := createTestCharacter()
	result := renderStatsPanel(char, 60)

	if result == "" {
		t.Error("renderStatsPanel() returned empty string")
	}

	// Should contain all stat sections
	expectedStrings := []string{
		"Character",
		"Experience",
		"Core Stats",
		"Streaks",
		char.Name,
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderStatsPanel() should contain %q", expected)
		}
	}
}

// TestRenderHistoryPanel tests the history panel rendering.
func TestRenderHistoryPanel(t *testing.T) {
	char := createTestCharacter()
	result := renderHistoryPanel(char, 60)

	if result == "" {
		t.Error("renderHistoryPanel() returned empty string")
	}

	// Should contain history sections
	expectedStrings := []string{
		"Today's Activity",
		"Lifetime Statistics",
		"Achievements",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderHistoryPanel() should contain %q", expected)
		}
	}
}

// TestRenderIdentitySection tests character identity rendering.
func TestRenderIdentitySection(t *testing.T) {
	char := createTestCharacter()
	result := renderIdentitySection(char)

	if result == "" {
		t.Error("renderIdentitySection() returned empty string")
	}

	// Should contain identity information
	expectedStrings := []string{
		"Character",
		"Name:",
		char.Name,
		"Level:",
		"Created:",
		"Days Played:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderIdentitySection() should contain %q", expected)
		}
	}
}

// TestRenderXPSection tests XP section rendering.
func TestRenderXPSection(t *testing.T) {
	char := createTestCharacter()
	result := renderXPSection(char, 60)

	if result == "" {
		t.Error("renderXPSection() returned empty string")
	}

	// Should contain XP information
	expectedStrings := []string{
		"Experience",
		"XP:",
		"To Next Level:",
		"Progress:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderXPSection() should contain %q", expected)
		}
	}
}

// TestRenderCoreStatsSection tests core stats section rendering.
func TestRenderCoreStatsSection(t *testing.T) {
	char := createTestCharacter()
	result := renderCoreStatsSection(char, 60)

	if result == "" {
		t.Error("renderCoreStatsSection() returned empty string")
	}

	// Should contain stat names
	expectedStrings := []string{
		"Core Stats",
		"CodePower",
		"Wisdom",
		"Agility",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderCoreStatsSection() should contain %q", expected)
		}
	}
}

// TestRenderStreakSectionDetailed tests streak section rendering.
func TestRenderStreakSectionDetailed(t *testing.T) {
	char := createTestCharacter()
	char.CurrentStreak = 5
	char.LongestStreak = 10

	result := renderStreakSectionDetailed(char)

	if result == "" {
		t.Error("renderStreakSectionDetailed() returned empty string")
	}

	// Should contain streak information
	expectedStrings := []string{
		"Streaks",
		"Current:",
		"5 days",
		"Longest:",
		"10 days",
		"Last Active:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderStreakSectionDetailed() should contain %q", expected)
		}
	}
}

// TestRenderTodayActivityDetailed tests today's activity rendering.
func TestRenderTodayActivityDetailed(t *testing.T) {
	char := createTestCharacter()
	char.TodayCommits = 5
	char.TodayLinesAdded = 250
	char.TodaySessionTime = 2 * time.Hour

	result := renderTodayActivityDetailed(char)

	if result == "" {
		t.Error("renderTodayActivityDetailed() returned empty string")
	}

	// Should contain today's stats
	expectedStrings := []string{
		"Today's Activity",
		"Commits:",
		"5",
		"Lines Added:",
		"250",
		"Session Time:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderTodayActivityDetailed() should contain %q", expected)
		}
	}
}

// TestRenderLifetimeStatsDetailed tests lifetime stats rendering.
func TestRenderLifetimeStatsDetailed(t *testing.T) {
	char := createTestCharacter()
	char.TotalCommits = 100
	char.TotalLinesAdded = 5000
	char.TotalLinesRemoved = 1000
	char.QuestsCompleted = 15

	result := renderLifetimeStatsDetailed(char)

	if result == "" {
		t.Error("renderLifetimeStatsDetailed() returned empty string")
	}

	// Should contain lifetime stats
	expectedStrings := []string{
		"Lifetime Statistics",
		"Total Commits:",
		"100",
		"Lines Added:",
		"5000",
		"Lines Removed:",
		"1000",
		"Quests Completed:",
		"15",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderLifetimeStatsDetailed() should contain %q", expected)
		}
	}
}

// TestRenderAchievementsPlaceholder tests achievements placeholder rendering.
func TestRenderAchievementsPlaceholder(t *testing.T) {
	result := renderAchievementsPlaceholder()

	if result == "" {
		t.Error("renderAchievementsPlaceholder() returned empty string")
	}

	// Should contain placeholder text
	expectedStrings := []string{
		"Achievements",
		"Coming soon",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderAchievementsPlaceholder() should contain %q", expected)
		}
	}
}

// TestRenderNoCharacterScreen tests the nil character handling.
func TestRenderNoCharacterScreen(t *testing.T) {
	result := renderNoCharacterScreen(80, 40)

	if result == "" {
		t.Error("renderNoCharacterScreen() returned empty string")
	}

	// Should contain error message
	expectedStrings := []string{
		"No character loaded",
		"restart",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderNoCharacterScreen() should contain %q", expected)
		}
	}
}

// TestRenderCharacterFooter tests footer rendering.
func TestRenderCharacterFooter(t *testing.T) {
	result := renderCharacterFooter(80)

	if result == "" {
		t.Error("renderCharacterFooter() returned empty string")
	}

	// Should contain key bindings
	expectedStrings := []string{
		"Alt+Q",
		"Dashboard",
		"Alt+M",
		"Mentor",
		"Esc",
		"Back",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderCharacterFooter() should contain %q", expected)
		}
	}
}

// TestFormatDate tests date formatting function.
func TestFormatDate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "today",
			input:    now,
			expected: "Today",
		},
		{
			name:     "yesterday",
			input:    now.Add(-24 * time.Hour),
			expected: "Yesterday",
		},
		{
			name:     "2 days ago",
			input:    now.Add(-48 * time.Hour),
			expected: "2 days ago",
		},
		{
			name:     "6 days ago",
			input:    now.Add(-6 * 24 * time.Hour),
			expected: "6 days ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDate(tt.input)
			if result != tt.expected {
				t.Errorf("formatDate() = %q, want %q", result, tt.expected)
			}
		})
	}

	// Test date more than a week ago (should show actual date)
	oldDate := now.Add(-30 * 24 * time.Hour)
	result := formatDate(oldDate)
	if !strings.Contains(result, ",") {
		t.Errorf("formatDate() for old date should contain date format with comma, got %q", result)
	}
}

// TestTruncateToDay tests day truncation function.
func TestTruncateToDay(t *testing.T) {
	// Create a time with hours, minutes, seconds
	testTime := time.Date(2025, 3, 15, 14, 30, 45, 0, time.UTC)

	result := truncateToDay(testTime)

	// Should be midnight of the same day
	expected := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)

	if !result.Equal(expected) {
		t.Errorf("truncateToDay() = %v, want %v", result, expected)
	}

	// Verify time components are zero
	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
		t.Errorf("truncateToDay() should set time to midnight, got %v", result)
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// createTestCharacter creates a character for testing with preset values.
func createTestCharacter() *game.Character {
	char := game.NewCharacter("TestHero")

	// Set some test values
	char.Level = 5
	char.XP = 450
	char.XPToNextLevel = 600
	char.CodePower = 15
	char.Wisdom = 12
	char.Agility = 14

	char.TotalCommits = 50
	char.TotalLinesAdded = 2500
	char.TotalLinesRemoved = 500
	char.QuestsCompleted = 8

	char.CurrentStreak = 3
	char.LongestStreak = 7

	char.TodayCommits = 3
	char.TodayLinesAdded = 150
	char.TodaySessionTime = 1*time.Hour + 30*time.Minute

	return char
}
