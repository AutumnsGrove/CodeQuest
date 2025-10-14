package game

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestNewCharacter tests the character creation function
func TestNewCharacter(t *testing.T) {
	tests := []struct {
		name     string
		charName string
		wantName string
	}{
		{"normal name", "TestHero", "TestHero"},
		{"empty name", "", ""},
		{"long name", "VeryLongCharacterNameThatGoesOnAndOn", "VeryLongCharacterNameThatGoesOnAndOn"},
		{"special characters", "Test-Hero_123", "Test-Hero_123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := NewCharacter(tt.charName)

			// Verify name
			if char.Name != tt.wantName {
				t.Errorf("NewCharacter() name = %v, want %v", char.Name, tt.wantName)
			}

			// Verify ID is a valid UUID
			if _, err := uuid.Parse(char.ID); err != nil {
				t.Errorf("NewCharacter() ID is not a valid UUID: %v", err)
			}

			// Verify starting stats
			if char.Level != 1 {
				t.Errorf("NewCharacter() Level = %v, want 1", char.Level)
			}
			if char.XP != 0 {
				t.Errorf("NewCharacter() XP = %v, want 0", char.XP)
			}
			if char.XPToNextLevel != calculateXPForLevel(1) {
				t.Errorf("NewCharacter() XPToNextLevel = %v, want %v", char.XPToNextLevel, calculateXPForLevel(1))
			}

			// Verify RPG stats
			if char.CodePower != 10 {
				t.Errorf("NewCharacter() CodePower = %v, want 10", char.CodePower)
			}
			if char.Wisdom != 10 {
				t.Errorf("NewCharacter() Wisdom = %v, want 10", char.Wisdom)
			}
			if char.Agility != 10 {
				t.Errorf("NewCharacter() Agility = %v, want 10", char.Agility)
			}

			// Verify progress tracking is zeroed
			if char.TotalCommits != 0 {
				t.Errorf("NewCharacter() TotalCommits = %v, want 0", char.TotalCommits)
			}
			if char.TotalLinesAdded != 0 {
				t.Errorf("NewCharacter() TotalLinesAdded = %v, want 0", char.TotalLinesAdded)
			}
			if char.TotalLinesRemoved != 0 {
				t.Errorf("NewCharacter() TotalLinesRemoved = %v, want 0", char.TotalLinesRemoved)
			}
			if char.QuestsCompleted != 0 {
				t.Errorf("NewCharacter() QuestsCompleted = %v, want 0", char.QuestsCompleted)
			}
			if char.CurrentStreak != 0 {
				t.Errorf("NewCharacter() CurrentStreak = %v, want 0", char.CurrentStreak)
			}
			if char.LongestStreak != 0 {
				t.Errorf("NewCharacter() LongestStreak = %v, want 0", char.LongestStreak)
			}

			// Verify session stats are zeroed
			if char.TodayCommits != 0 {
				t.Errorf("NewCharacter() TodayCommits = %v, want 0", char.TodayCommits)
			}
			if char.TodayLinesAdded != 0 {
				t.Errorf("NewCharacter() TodayLinesAdded = %v, want 0", char.TodayLinesAdded)
			}
			if char.TodaySessionTime != 0 {
				t.Errorf("NewCharacter() TodaySessionTime = %v, want 0", char.TodaySessionTime)
			}

			// Verify timestamps are recent (within last second)
			now := time.Now()
			if now.Sub(char.CreatedAt) > time.Second {
				t.Errorf("NewCharacter() CreatedAt is too old: %v", char.CreatedAt)
			}
			if now.Sub(char.LastActiveDate) > time.Second {
				t.Errorf("NewCharacter() LastActiveDate is too old: %v", char.LastActiveDate)
			}
		})
	}
}

// TestCharacter_AddXP tests the XP addition and level-up logic
func TestCharacter_AddXP(t *testing.T) {
	tests := []struct {
		name          string
		initialLevel  int
		initialXP     int
		xpToAdd       int
		wantLevel     int
		wantLeveledUp bool
		wantXPRange   [2]int // Min and max XP after adding (handles overflow)
	}{
		{
			name:          "zero XP",
			initialLevel:  1,
			initialXP:     0,
			xpToAdd:       0,
			wantLevel:     1,
			wantLeveledUp: false,
			wantXPRange:   [2]int{0, 0},
		},
		{
			name:          "negative XP",
			initialLevel:  1,
			initialXP:     0,
			xpToAdd:       -50,
			wantLevel:     1,
			wantLeveledUp: false,
			wantXPRange:   [2]int{0, 0},
		},
		{
			name:          "small XP gain no level up",
			initialLevel:  1,
			initialXP:     0,
			xpToAdd:       50,
			wantLevel:     1,
			wantLeveledUp: false,
			wantXPRange:   [2]int{50, 50},
		},
		{
			name:          "exact XP for level up",
			initialLevel:  1,
			initialXP:     0,
			xpToAdd:       110, // Level 1→2 requires 110 XP
			wantLevel:     2,
			wantLeveledUp: true,
			wantXPRange:   [2]int{0, 0},
		},
		{
			name:          "XP with overflow into next level",
			initialLevel:  1,
			initialXP:     0,
			xpToAdd:       150, // 110 for level up, 40 remaining
			wantLevel:     2,
			wantLeveledUp: true,
			wantXPRange:   [2]int{40, 40},
		},
		{
			name:          "multiple level ups at once",
			initialLevel:  1,
			initialXP:     0,
			xpToAdd:       500, // Should level up multiple times
			wantLevel:     3,   // Level 1→2 (110) + Level 2→3 (240) = 350, 150 remaining
			wantLeveledUp: true,
			wantXPRange:   [2]int{150, 150},
		},
		{
			name:          "level up from existing XP",
			initialLevel:  1,
			initialXP:     100,
			xpToAdd:       20, // 100 + 20 = 120, should level up with 10 remaining
			wantLevel:     2,
			wantLeveledUp: true,
			wantXPRange:   [2]int{10, 10},
		},
		{
			name:          "high level XP gain",
			initialLevel:  10,
			initialXP:     0,
			xpToAdd:       100,
			wantLevel:     10,
			wantLeveledUp: false,
			wantXPRange:   [2]int{100, 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := NewCharacter("TestHero")
			char.Level = tt.initialLevel
			char.XP = tt.initialXP
			char.XPToNextLevel = calculateXPForLevel(tt.initialLevel)

			// Store initial stats to verify level-up bonuses
			initialCodePower := char.CodePower
			initialWisdom := char.Wisdom
			initialAgility := char.Agility

			gotLeveledUp := char.AddXP(tt.xpToAdd)

			// Verify level-up flag
			if gotLeveledUp != tt.wantLeveledUp {
				t.Errorf("AddXP() leveledUp = %v, want %v", gotLeveledUp, tt.wantLeveledUp)
			}

			// Verify final level
			if char.Level != tt.wantLevel {
				t.Errorf("AddXP() Level = %v, want %v", char.Level, tt.wantLevel)
			}

			// Verify XP is in expected range
			if char.XP < tt.wantXPRange[0] || char.XP > tt.wantXPRange[1] {
				t.Errorf("AddXP() XP = %v, want between %v and %v", char.XP, tt.wantXPRange[0], tt.wantXPRange[1])
			}

			// Verify XPToNextLevel is set correctly for new level
			expectedXPToNext := calculateXPForLevel(char.Level)
			if char.XPToNextLevel != expectedXPToNext {
				t.Errorf("AddXP() XPToNextLevel = %v, want %v", char.XPToNextLevel, expectedXPToNext)
			}

			// Verify stat increases on level up
			levelsGained := tt.wantLevel - tt.initialLevel
			if gotLeveledUp {
				if char.CodePower != initialCodePower+levelsGained {
					t.Errorf("AddXP() CodePower = %v, want %v", char.CodePower, initialCodePower+levelsGained)
				}
				if char.Wisdom != initialWisdom+levelsGained {
					t.Errorf("AddXP() Wisdom = %v, want %v", char.Wisdom, initialWisdom+levelsGained)
				}
				if char.Agility != initialAgility+levelsGained {
					t.Errorf("AddXP() Agility = %v, want %v", char.Agility, initialAgility+levelsGained)
				}
			} else {
				// No level up, stats should be unchanged
				if char.CodePower != initialCodePower {
					t.Errorf("AddXP() CodePower changed without level up: %v, want %v", char.CodePower, initialCodePower)
				}
				if char.Wisdom != initialWisdom {
					t.Errorf("AddXP() Wisdom changed without level up: %v, want %v", char.Wisdom, initialWisdom)
				}
				if char.Agility != initialAgility {
					t.Errorf("AddXP() Agility changed without level up: %v, want %v", char.Agility, initialAgility)
				}
			}
		})
	}
}

// TestCharacter_UpdateStreak tests the streak tracking functionality
func TestCharacter_UpdateStreak(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	twoDaysAgo := now.AddDate(0, 0, -2)
	weekAgo := now.AddDate(0, 0, -7)

	tests := []struct {
		name               string
		lastActiveDate     time.Time
		currentStreak      int
		longestStreak      int
		wantCurrentStreak  int
		wantLongestStreak  int
		wantLastActiveDate time.Time // We'll check if it's today
	}{
		{
			name:               "first activity ever (streak 0)",
			lastActiveDate:     now,
			currentStreak:      0,
			longestStreak:      0,
			wantCurrentStreak:  1,
			wantLongestStreak:  1,
			wantLastActiveDate: now,
		},
		{
			name:               "already active today",
			lastActiveDate:     now,
			currentStreak:      5,
			longestStreak:      10,
			wantCurrentStreak:  5, // No change
			wantLongestStreak:  10,
			wantLastActiveDate: now,
		},
		{
			name:               "active yesterday (increment)",
			lastActiveDate:     yesterday,
			currentStreak:      5,
			longestStreak:      10,
			wantCurrentStreak:  6,
			wantLongestStreak:  10,
			wantLastActiveDate: now,
		},
		{
			name:               "active yesterday (new longest)",
			lastActiveDate:     yesterday,
			currentStreak:      9,
			longestStreak:      9,
			wantCurrentStreak:  10,
			wantLongestStreak:  10,
			wantLastActiveDate: now,
		},
		{
			name:               "missed one day (reset to 1)",
			lastActiveDate:     twoDaysAgo,
			currentStreak:      5,
			longestStreak:      10,
			wantCurrentStreak:  1,
			wantLongestStreak:  10,
			wantLastActiveDate: now,
		},
		{
			name:               "missed many days (reset to 1)",
			lastActiveDate:     weekAgo,
			currentStreak:      20,
			longestStreak:      20,
			wantCurrentStreak:  1,
			wantLongestStreak:  20,
			wantLastActiveDate: now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := NewCharacter("TestHero")
			char.LastActiveDate = tt.lastActiveDate
			char.CurrentStreak = tt.currentStreak
			char.LongestStreak = tt.longestStreak

			char.UpdateStreak()

			// Verify current streak
			if char.CurrentStreak != tt.wantCurrentStreak {
				t.Errorf("UpdateStreak() CurrentStreak = %v, want %v", char.CurrentStreak, tt.wantCurrentStreak)
			}

			// Verify longest streak
			if char.LongestStreak != tt.wantLongestStreak {
				t.Errorf("UpdateStreak() LongestStreak = %v, want %v", char.LongestStreak, tt.wantLongestStreak)
			}

			// Verify LastActiveDate is today
			if !char.IsToday(char.LastActiveDate) {
				t.Errorf("UpdateStreak() LastActiveDate should be today, got %v", char.LastActiveDate)
			}
		})
	}
}

// TestCharacter_ResetDailyStats tests daily stats reset
func TestCharacter_ResetDailyStats(t *testing.T) {
	char := NewCharacter("TestHero")

	// Set some daily stats
	char.TodayCommits = 10
	char.TodayLinesAdded = 500
	char.TodaySessionTime = 2 * time.Hour

	// Reset them
	char.ResetDailyStats()

	// Verify all are zeroed
	if char.TodayCommits != 0 {
		t.Errorf("ResetDailyStats() TodayCommits = %v, want 0", char.TodayCommits)
	}
	if char.TodayLinesAdded != 0 {
		t.Errorf("ResetDailyStats() TodayLinesAdded = %v, want 0", char.TodayLinesAdded)
	}
	if char.TodaySessionTime != 0 {
		t.Errorf("ResetDailyStats() TodaySessionTime = %v, want 0", char.TodaySessionTime)
	}

	// Verify other stats are unchanged
	if char.Level != 1 {
		t.Errorf("ResetDailyStats() changed Level to %v, should be unchanged", char.Level)
	}
	if char.TotalCommits != 0 {
		t.Errorf("ResetDailyStats() changed TotalCommits to %v, should be unchanged", char.TotalCommits)
	}
}

// TestCharacter_IsToday tests the date comparison helper
func TestCharacter_IsToday(t *testing.T) {
	now := time.Now()
	today := now
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)
	todayMorning := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEvening := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	tests := []struct {
		name string
		date time.Time
		want bool
	}{
		{"now", today, true},
		{"yesterday", yesterday, false},
		{"tomorrow", tomorrow, false},
		{"today morning", todayMorning, true},
		{"today evening", todayEvening, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := NewCharacter("TestHero")
			if got := char.IsToday(tt.date); got != tt.want {
				t.Errorf("IsToday() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTruncateToDay tests the day truncation helper
func TestTruncateToDay(t *testing.T) {
	tests := []struct {
		name  string
		input time.Time
	}{
		{"morning", time.Date(2024, 1, 15, 9, 30, 45, 0, time.UTC)},
		{"noon", time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)},
		{"evening", time.Date(2024, 1, 15, 23, 59, 59, 0, time.UTC)},
		{"midnight", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateToDay(tt.input)

			// Should be midnight
			if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
				t.Errorf("truncateToDay() time = %v:%v:%v, want 00:00:00",
					result.Hour(), result.Minute(), result.Second())
			}

			// Should be same date
			if result.Year() != tt.input.Year() ||
				result.Month() != tt.input.Month() ||
				result.Day() != tt.input.Day() {
				t.Errorf("truncateToDay() date = %v, want %v",
					result.Format("2006-01-02"), tt.input.Format("2006-01-02"))
			}

			// Should preserve location
			if result.Location() != tt.input.Location() {
				t.Errorf("truncateToDay() location = %v, want %v",
					result.Location(), tt.input.Location())
			}
		})
	}
}

// TestGenerateID tests UUID generation
func TestGenerateID(t *testing.T) {
	// Generate multiple IDs and verify they're valid UUIDs and unique
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateID()

		// Verify it's a valid UUID
		if _, err := uuid.Parse(id); err != nil {
			t.Errorf("generateID() produced invalid UUID: %v, error: %v", id, err)
		}

		// Verify it's unique
		if ids[id] {
			t.Errorf("generateID() produced duplicate UUID: %v", id)
		}
		ids[id] = true
	}

	// Verify we generated 100 unique IDs
	if len(ids) != 100 {
		t.Errorf("generateID() uniqueness check failed: got %v unique IDs, want 100", len(ids))
	}
}

// TestCharacter_MultiLevelUp tests gaining multiple levels in one AddXP call
func TestCharacter_MultiLevelUp(t *testing.T) {
	char := NewCharacter("TestHero")

	// Level 1→2: 110 XP
	// Level 2→3: 240 XP
	// Level 3→4: 390 XP
	// Total: 740 XP

	// Add enough XP to go from level 1 to level 4
	xpForThreeLevels := calculateXPForLevel(1) + calculateXPForLevel(2) + calculateXPForLevel(3)
	leveledUp := char.AddXP(xpForThreeLevels)

	// Verify we leveled up
	if !leveledUp {
		t.Errorf("AddXP() with multi-level XP should return true")
	}

	// Verify we're at level 4
	if char.Level != 4 {
		t.Errorf("AddXP() Level = %v, want 4", char.Level)
	}

	// Verify stats increased by 3 (3 level-ups)
	if char.CodePower != 13 {
		t.Errorf("AddXP() CodePower = %v, want 13 (10 + 3 levels)", char.CodePower)
	}
	if char.Wisdom != 13 {
		t.Errorf("AddXP() Wisdom = %v, want 13 (10 + 3 levels)", char.Wisdom)
	}
	if char.Agility != 13 {
		t.Errorf("AddXP() Agility = %v, want 13 (10 + 3 levels)", char.Agility)
	}

	// Verify XP is 0 (exact amount for levels)
	if char.XP != 0 {
		t.Errorf("AddXP() XP = %v, want 0 (exact level-up)", char.XP)
	}

	// Verify XPToNextLevel is set for level 4
	if char.XPToNextLevel != calculateXPForLevel(4) {
		t.Errorf("AddXP() XPToNextLevel = %v, want %v", char.XPToNextLevel, calculateXPForLevel(4))
	}
}
