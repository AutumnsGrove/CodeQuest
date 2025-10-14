package game

import (
	"math"
	"testing"
)

// TestCalculateXPForLevel tests the level progression curve
func TestCalculateXPForLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		wantXP   int
		wantDesc string
	}{
		{"level 1", 1, 110, "early game fast progression"},
		{"level 2", 2, 240, "early game progression"},
		{"level 5", 5, 750, "mid-early game"},
		{"level 10", 10, 2000, "reached double digits"},
		{"level 25", 25, 8750, "mid game"},
		{"level 50", 50, 30000, "late mid game"},
		{"level 75", 75, 63750, "late game"},
		{"level 99", 99, 107910, "near max level"},
		{"level 100", 100, maxLevelXPRequired, "max level cap"},
		{"level 101", 101, maxLevelXPRequired, "beyond max level"},
		{"level 999", 999, maxLevelXPRequired, "way beyond max"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateXPForLevel(tt.level)
			if got != tt.wantXP {
				t.Errorf("CalculateXPForLevel(%v) = %v, want %v (%s)",
					tt.level, got, tt.wantXP, tt.wantDesc)
			}
		})
	}
}

// TestCalculateXPForLevel_Progression tests that XP curve is always increasing
func TestCalculateXPForLevel_Progression(t *testing.T) {
	prevXP := 0
	for level := 1; level < maxLevel; level++ {
		xp := CalculateXPForLevel(level)

		// Each level should require more XP than the previous
		if xp <= prevXP {
			t.Errorf("CalculateXPForLevel(%v) = %v, should be greater than previous level XP %v",
				level, xp, prevXP)
		}

		// XP should never be negative or zero
		if xp <= 0 {
			t.Errorf("CalculateXPForLevel(%v) = %v, should be positive", level, xp)
		}

		prevXP = xp
	}
}

// TestCalculateCommitXP tests commit XP calculation
func TestCalculateCommitXP(t *testing.T) {
	tests := []struct {
		name         string
		linesAdded   int
		linesRemoved int
		wantXP       int
	}{
		{
			name:         "empty commit",
			linesAdded:   0,
			linesRemoved: 0,
			wantXP:       baseCommitXP, // 10
		},
		{
			name:         "small commit (10 lines added)",
			linesAdded:   10,
			linesRemoved: 0,
			wantXP:       baseCommitXP + 10, // 20
		},
		{
			name:         "small commit (5 added, 5 removed)",
			linesAdded:   5,
			linesRemoved: 5,
			wantXP:       baseCommitXP + 10, // 20
		},
		{
			name:         "medium commit (30 lines)",
			linesAdded:   20,
			linesRemoved: 10,
			wantXP:       baseCommitXP + 30, // 40
		},
		{
			name:         "at cap (50 lines)",
			linesAdded:   50,
			linesRemoved: 0,
			wantXP:       baseCommitXP + maxLinesBonus, // 60
		},
		{
			name:         "beyond cap (100 lines added)",
			linesAdded:   100,
			linesRemoved: 0,
			wantXP:       baseCommitXP + maxLinesBonus, // 60 (capped)
		},
		{
			name:         "beyond cap (500 total)",
			linesAdded:   300,
			linesRemoved: 200,
			wantXP:       baseCommitXP + maxLinesBonus, // 60 (capped)
		},
		{
			name:         "large refactor (1000 added, 1000 removed)",
			linesAdded:   1000,
			linesRemoved: 1000,
			wantXP:       baseCommitXP + maxLinesBonus, // 60 (capped)
		},
		{
			name:         "negative lines added (invalid)",
			linesAdded:   -10,
			linesRemoved: 0,
			wantXP:       baseCommitXP, // 10 (negative clamped to 0)
		},
		{
			name:         "negative lines removed (invalid)",
			linesAdded:   0,
			linesRemoved: -10,
			wantXP:       baseCommitXP, // 10 (negative clamped to 0)
		},
		{
			name:         "both negative (invalid)",
			linesAdded:   -50,
			linesRemoved: -50,
			wantXP:       baseCommitXP, // 10 (both clamped to 0)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateCommitXP(tt.linesAdded, tt.linesRemoved)
			if got != tt.wantXP {
				t.Errorf("CalculateCommitXP(%v, %v) = %v, want %v",
					tt.linesAdded, tt.linesRemoved, got, tt.wantXP)
			}

			// Verify XP is always at least baseCommitXP
			if got < baseCommitXP {
				t.Errorf("CalculateCommitXP(%v, %v) = %v, should never be less than baseCommitXP (%v)",
					tt.linesAdded, tt.linesRemoved, got, baseCommitXP)
			}

			// Verify XP never exceeds base + cap
			maxPossible := baseCommitXP + maxLinesBonus
			if got > maxPossible {
				t.Errorf("CalculateCommitXP(%v, %v) = %v, should never exceed %v",
					tt.linesAdded, tt.linesRemoved, got, maxPossible)
			}
		})
	}
}

// TestApplyDifficultyMultiplier tests difficulty scaling
func TestApplyDifficultyMultiplier(t *testing.T) {
	tests := []struct {
		name       string
		baseXP     int
		difficulty string
		wantXP     int
	}{
		// Easy mode (1.2x)
		{"easy - 50 XP", 50, DifficultyEasy, 60},
		{"easy - 100 XP", 100, DifficultyEasy, 120},
		{"easy - 0 XP", 0, DifficultyEasy, 0},

		// Normal mode (1.0x)
		{"normal - 50 XP", 50, DifficultyNormal, 50},
		{"normal - 100 XP", 100, DifficultyNormal, 100},
		{"normal - 0 XP", 0, DifficultyNormal, 0},

		// Hard mode (0.8x)
		{"hard - 50 XP", 50, DifficultyHard, 40},
		{"hard - 100 XP", 100, DifficultyHard, 80},
		{"hard - 0 XP", 0, DifficultyHard, 0},

		// Unknown/invalid difficulty (defaults to normal)
		{"invalid - 50 XP", 50, "invalid", 50},
		{"empty - 100 XP", 100, "", 100},

		// Rounding tests
		{"easy - rounds up", 51, DifficultyEasy, 61},    // 51 * 1.2 = 61.2 → 61
		{"hard - rounds down", 51, DifficultyHard, 41},  // 51 * 0.8 = 40.8 → 41
		{"hard - rounds up", 56, DifficultyHard, 45},    // 56 * 0.8 = 44.8 → 45
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyDifficultyMultiplier(tt.baseXP, tt.difficulty)
			if got != tt.wantXP {
				t.Errorf("ApplyDifficultyMultiplier(%v, %v) = %v, want %v",
					tt.baseXP, tt.difficulty, got, tt.wantXP)
			}
		})
	}
}

// TestApplyWisdomBonus tests wisdom-based XP scaling
func TestApplyWisdomBonus(t *testing.T) {
	tests := []struct {
		name   string
		baseXP int
		wisdom int
		wantXP int
	}{
		// Base wisdom (10) - no bonus
		{"10 wisdom - 50 XP", 50, 10, 50},
		{"10 wisdom - 100 XP", 100, 10, 100},

		// Below base wisdom - should still give 1.0x (minimum)
		{"5 wisdom - 50 XP", 50, 5, 50},
		{"0 wisdom - 100 XP", 100, 0, 100},

		// Typical wisdom progression
		{"15 wisdom - 100 XP", 100, 15, 105}, // +5% = 105
		{"20 wisdom - 100 XP", 100, 20, 110}, // +10% = 110
		{"30 wisdom - 100 XP", 100, 30, 120}, // +20% = 120
		{"50 wisdom - 100 XP", 100, 50, 140}, // +40% = 140

		// High wisdom
		{"100 wisdom - 100 XP", 100, 100, 190}, // +90% = 190
		{"100 wisdom - 50 XP", 50, 100, 95},    // +90% = 95

		// Edge cases
		{"0 XP regardless of wisdom", 0, 50, 0},
		{"negative wisdom", 100, -10, 100}, // Should clamp to 1.0x

		// Rounding tests
		{"15 wisdom - 51 XP", 51, 15, 54},  // 51 * 1.05 = 53.55 → 54
		{"20 wisdom - 47 XP", 47, 20, 52},  // 47 * 1.10 = 51.7 → 52
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyWisdomBonus(tt.baseXP, tt.wisdom)
			if got != tt.wantXP {
				t.Errorf("ApplyWisdomBonus(%v, %v) = %v, want %v",
					tt.baseXP, tt.wisdom, got, tt.wantXP)
			}

			// Wisdom should never reduce XP
			if got < tt.baseXP {
				t.Errorf("ApplyWisdomBonus(%v, %v) = %v, wisdom should never reduce XP below base %v",
					tt.baseXP, tt.wisdom, got, tt.baseXP)
			}
		})
	}
}

// TestCalculateQuestReward tests quest XP rewards
func TestCalculateQuestReward(t *testing.T) {
	tests := []struct {
		name            string
		questDifficulty string
		wantXP          int
	}{
		{"simple quest", QuestDifficultySimple, questXPSimple},     // 50
		{"medium quest", QuestDifficultyMedium, questXPMedium},     // 150
		{"hard quest", QuestDifficultyHard, questXPHard},           // 300
		{"epic quest", QuestDifficultyEpic, questXPEpic},           // 1000
		{"invalid difficulty", "invalid", questXPSimple},           // defaults to simple
		{"empty difficulty", "", questXPSimple},                    // defaults to simple
		{"case sensitive", "SIMPLE", questXPSimple},                // exact match or default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateQuestReward(tt.questDifficulty)
			if got != tt.wantXP {
				t.Errorf("CalculateQuestReward(%v) = %v, want %v",
					tt.questDifficulty, got, tt.wantXP)
			}

			// Quest rewards should always be positive
			if got <= 0 {
				t.Errorf("CalculateQuestReward(%v) = %v, should be positive",
					tt.questDifficulty, got)
			}
		})
	}
}

// TestGetLevelFromXP tests level calculation from total XP
func TestGetLevelFromXP(t *testing.T) {
	tests := []struct {
		name          string
		totalXP       int
		wantLevel     int
		wantRemaining int
	}{
		{"no XP", 0, 1, 0},
		{"50 XP (partial level 1)", 50, 1, 50},
		{"110 XP (exact level 2)", 110, 2, 0},
		{"150 XP (level 2 with overflow)", 150, 2, 40},
		{"350 XP (exact level 3)", 350, 3, 0},
		{"500 XP (mid level 3)", 500, 3, 150},
		{"10000 XP (higher level)", 10000, 11, 650},
		{"100000 XP (very high level)", 100000, 27, 2890},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLevel, gotRemaining := GetLevelFromXP(tt.totalXP)

			if gotLevel != tt.wantLevel {
				t.Errorf("GetLevelFromXP(%v) level = %v, want %v",
					tt.totalXP, gotLevel, tt.wantLevel)
			}

			if gotRemaining != tt.wantRemaining {
				t.Errorf("GetLevelFromXP(%v) remaining = %v, want %v",
					tt.totalXP, gotRemaining, tt.wantRemaining)
			}

			// Verify remaining XP is less than next level requirement
			xpForNextLevel := CalculateXPForLevel(gotLevel)
			if gotRemaining >= xpForNextLevel {
				t.Errorf("GetLevelFromXP(%v) remaining %v >= next level requirement %v",
					tt.totalXP, gotRemaining, xpForNextLevel)
			}

			// Verify remaining is non-negative
			if gotRemaining < 0 {
				t.Errorf("GetLevelFromXP(%v) remaining = %v, should be non-negative",
					tt.totalXP, gotRemaining)
			}
		})
	}
}

// TestGetLevelFromXP_MaxLevel tests max level capping
func TestGetLevelFromXP_MaxLevel(t *testing.T) {
	// XP way beyond what's needed for max level
	hugeXP := 10000000

	level, remaining := GetLevelFromXP(hugeXP)

	// Should cap at max level
	if level > maxLevel {
		t.Errorf("GetLevelFromXP(%v) level = %v, should not exceed maxLevel %v",
			hugeXP, level, maxLevel)
	}

	// Remaining should be positive (all the overflow XP)
	if remaining < 0 {
		t.Errorf("GetLevelFromXP(%v) remaining = %v, should be non-negative",
			hugeXP, remaining)
	}
}

// TestGetProgressToNextLevel tests progress percentage calculation
func TestGetProgressToNextLevel(t *testing.T) {
	tests := []struct {
		name         string
		currentXP    int
		xpNeeded     int
		wantProgress float64
	}{
		{"no progress", 0, 100, 0.0},
		{"25% progress", 25, 100, 0.25},
		{"50% progress", 50, 100, 0.5},
		{"75% progress", 75, 100, 0.75},
		{"99% progress", 99, 100, 0.99},
		{"100% progress", 100, 100, 1.0},
		{"over 100% (should cap)", 150, 100, 1.0},
		{"xp needed is 0 (edge case)", 50, 0, 1.0},
		{"negative current (edge case)", -10, 100, 0.0},
		{"negative needed (edge case)", 50, -100, 1.0},
		{"both zero", 0, 0, 1.0},
		{"realistic example", 45, 200, 0.225},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetProgressToNextLevel(tt.currentXP, tt.xpNeeded)

			// Allow small floating point tolerance
			tolerance := 0.001
			if math.Abs(got-tt.wantProgress) > tolerance {
				t.Errorf("GetProgressToNextLevel(%v, %v) = %v, want %v",
					tt.currentXP, tt.xpNeeded, got, tt.wantProgress)
			}

			// Progress should always be between 0.0 and 1.0
			if got < 0.0 || got > 1.0 {
				t.Errorf("GetProgressToNextLevel(%v, %v) = %v, should be between 0.0 and 1.0",
					tt.currentXP, tt.xpNeeded, got)
			}
		})
	}
}

// TestGetTotalXPForLevel tests cumulative XP calculation
func TestGetTotalXPForLevel(t *testing.T) {
	tests := []struct {
		name        string
		targetLevel int
		wantXP      int
	}{
		{"level 1 (starting)", 1, 0},
		{"level 2", 2, 110},                                                 // 110
		{"level 3", 3, 350},                                                 // 110 + 240
		{"level 4", 4, 740},                                                 // 110 + 240 + 390
		{"level 10", 10, 7350},                                              // sum of levels 1-9
		{"level 0 (invalid)", 0, 0},
		{"level -5 (invalid)", -5, 0},
		{"level 101 (beyond max)", 101, GetTotalXPForLevel(maxLevel)},
		{"level 999 (way beyond max)", 999, GetTotalXPForLevel(maxLevel)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTotalXPForLevel(tt.targetLevel)
			if got != tt.wantXP {
				t.Errorf("GetTotalXPForLevel(%v) = %v, want %v",
					tt.targetLevel, got, tt.wantXP)
			}

			// Total XP should never be negative
			if got < 0 {
				t.Errorf("GetTotalXPForLevel(%v) = %v, should be non-negative",
					tt.targetLevel, got)
			}
		})
	}
}

// TestGetTotalXPForLevel_Progression tests that total XP is always increasing
func TestGetTotalXPForLevel_Progression(t *testing.T) {
	prevTotal := 0
	for level := 1; level <= 20; level++ {
		total := GetTotalXPForLevel(level)

		// Total XP should increase with each level (or stay same for level 1)
		if total < prevTotal {
			t.Errorf("GetTotalXPForLevel(%v) = %v, should be >= previous total %v",
				level, total, prevTotal)
		}

		prevTotal = total
	}
}

// TestXPEngine_Integration tests the full XP calculation pipeline
func TestXPEngine_Integration(t *testing.T) {
	// Simulate a player making commits and leveling up

	// Start at level 1
	currentLevel := 1
	currentXP := 0
	xpNeeded := CalculateXPForLevel(currentLevel)
	wisdom := 10

	// Player makes a medium commit (30 lines) on normal difficulty
	commitXP := CalculateCommitXP(20, 10) // 40 base XP
	commitXP = ApplyDifficultyMultiplier(commitXP, DifficultyNormal)
	commitXP = ApplyWisdomBonus(commitXP, wisdom)

	// Should be 40 XP
	if commitXP != 40 {
		t.Errorf("Integration test: commit XP = %v, want 40", commitXP)
	}

	// Add to character
	currentXP += commitXP
	progress := GetProgressToNextLevel(currentXP, xpNeeded)

	// Progress should be 40/110 ≈ 0.364
	expectedProgress := 40.0 / 110.0
	if math.Abs(progress-expectedProgress) > 0.01 {
		t.Errorf("Integration test: progress = %v, want ~%v", progress, expectedProgress)
	}

	// Complete a simple quest (50 XP)
	questXP := CalculateQuestReward(QuestDifficultySimple)
	questXP = ApplyDifficultyMultiplier(questXP, DifficultyNormal)
	questXP = ApplyWisdomBonus(questXP, wisdom)

	// Should be 50 XP
	if questXP != 50 {
		t.Errorf("Integration test: quest XP = %v, want 50", questXP)
	}

	currentXP += questXP // Now at 90 XP

	// Should still be level 1
	level, remaining := GetLevelFromXP(currentXP)
	if level != 1 {
		t.Errorf("Integration test: level = %v, want 1 (not enough XP yet)", level)
	}
	if remaining != 90 {
		t.Errorf("Integration test: remaining XP = %v, want 90", remaining)
	}

	// Make another commit to level up (need 20 more XP)
	commitXP2 := CalculateCommitXP(20, 0) // 30 base XP
	currentXP += commitXP2                 // Now at 120 XP, should level up

	level, remaining = GetLevelFromXP(currentXP)
	if level != 2 {
		t.Errorf("Integration test: level = %v, want 2 (should have leveled up)", level)
	}
	if remaining != 10 {
		t.Errorf("Integration test: remaining XP = %v, want 10", remaining)
	}
}

// TestXPEngine_DifficultyImpact tests how difficulty affects progression
func TestXPEngine_DifficultyImpact(t *testing.T) {
	baseXP := 100

	easyXP := ApplyDifficultyMultiplier(baseXP, DifficultyEasy)
	normalXP := ApplyDifficultyMultiplier(baseXP, DifficultyNormal)
	hardXP := ApplyDifficultyMultiplier(baseXP, DifficultyHard)

	// Easy should give more XP
	if easyXP <= normalXP {
		t.Errorf("Easy difficulty (%v) should give more XP than normal (%v)", easyXP, normalXP)
	}

	// Hard should give less XP
	if hardXP >= normalXP {
		t.Errorf("Hard difficulty (%v) should give less XP than normal (%v)", hardXP, normalXP)
	}

	// Verify ratios
	if easyXP != 120 {
		t.Errorf("Easy XP = %v, want 120 (20%% more)", easyXP)
	}
	if normalXP != 100 {
		t.Errorf("Normal XP = %v, want 100 (no change)", normalXP)
	}
	if hardXP != 80 {
		t.Errorf("Hard XP = %v, want 80 (20%% less)", hardXP)
	}
}

// TestXPEngine_WisdomScaling tests wisdom's impact on XP gain
func TestXPEngine_WisdomScaling(t *testing.T) {
	baseXP := 100

	// Test wisdom growth over character lifetime
	wisdomLevels := []struct {
		wisdom   int
		expected int
	}{
		{10, 100},  // Starting wisdom: 0% bonus
		{20, 110},  // +10 wisdom: +10% bonus
		{30, 120},  // +20 wisdom: +20% bonus
		{50, 140},  // +40 wisdom: +40% bonus
		{100, 190}, // +90 wisdom: +90% bonus
	}

	for _, wl := range wisdomLevels {
		got := ApplyWisdomBonus(baseXP, wl.wisdom)
		if got != wl.expected {
			t.Errorf("Wisdom %v: got %v XP, want %v", wl.wisdom, got, wl.expected)
		}
	}
}
