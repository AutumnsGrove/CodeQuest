// Package game contains the XP calculation engine for CodeQuest.
// This file implements the progression system that determines XP rewards and level thresholds.
package game

import "math"

// XP Engine Constants
// These values control the overall progression pace and balance.
const (
	// Level Progression - XP curve parameters
	// The formula: baseXP * level * (1 + level/levelScaleFactor)
	// This creates a smooth polynomial curve that starts gentle and gradually steepens
	baseXPPerLevel     = 100     // Base XP multiplier - level 1→2 costs 100 XP
	levelScaleFactor   = 10.0    // Controls curve steepness - lower = steeper, higher = gentler
	maxLevel           = 100     // Level cap to prevent integer overflow and maintain balance
	maxLevelXPRequired = 1000000 // XP required for max level (prevents overflow)

	// Commit XP Rewards - How much XP you earn from commits
	baseCommitXP      = 10 // Minimum XP per commit (even empty commits count!)
	xpPerLinesChanged = 1  // XP per line changed (added + removed)
	maxLinesBonus     = 50 // Cap on lines bonus to prevent farming (max 60 XP per commit)

	// Quest XP Rewards - XP awarded for completing quests by difficulty
	questXPSimple = 50   // Simple quests (e.g., "Make 3 commits")
	questXPMedium = 150  // Medium quests (e.g., "Add 100 lines")
	questXPHard   = 300  // Hard quests (e.g., "7 day streak")
	questXPEpic   = 1000 // Epic quests (e.g., "Reach level 10")

	// Difficulty Multipliers - Scale all XP based on game difficulty
	// These multiply ALL XP gains (commits, quests, bonuses)
	difficultyEasy   = 1.2 // Easy mode: 20% more XP (faster progression for beginners)
	difficultyNormal = 1.0 // Normal mode: Standard XP rates (balanced for most players)
	difficultyHard   = 0.8 // Hard mode: 20% less XP (challenging for veterans)

	// Wisdom Stat Bonus - Character's Wisdom stat increases XP gain
	// Formula: 1 + (wisdom - 10) * 0.01
	// At 10 wisdom (starting): 0% bonus
	// At 20 wisdom: +10% bonus
	// At 50 wisdom: +40% bonus
	// At 100 wisdom: +90% bonus (nearly double XP!)
	wisdomBaseValue  = 10   // Starting wisdom value (no bonus)
	wisdomBonusPer10 = 0.01 // 1% bonus per wisdom point above base
)

// Quest difficulty types for XP calculation
const (
	QuestDifficultySimple = "simple"
	QuestDifficultyMedium = "medium"
	QuestDifficultyHard   = "hard"
	QuestDifficultyEpic   = "epic"
)

// Game difficulty types (from config)
const (
	DifficultyEasy   = "easy"
	DifficultyNormal = "normal"
	DifficultyHard   = "hard"
)

// CalculateXPForLevel calculates the total XP required to advance from the given level to the next.
// This implements a polynomial progression curve that balances fast early levels with challenging late game.
//
// Progression Examples:
//   - Level 1→2:   100 XP  (~10 commits)
//   - Level 5→6:   600 XP  (~40 commits)
//   - Level 10→11: 2100 XP (~140 commits)
//   - Level 25→26: 10625 XP (~700 commits)
//   - Level 50→51: 35500 XP (~2370 commits)
//   - Level 75→76: 69625 XP (~4640 commits)
//   - Level 99→100: 117910 XP (~7860 commits)
//
// The formula: baseXP * level * (1 + level/levelScaleFactor)
// This creates a curve where:
//   - Early levels (1-10): Fast progression to hook new players
//   - Mid levels (10-50): Steady, satisfying growth
//   - Late levels (50-100): Challenging but achievable with dedication
//
// Parameters:
//   - level: The current level (1-100)
//
// Returns:
//   - int: The XP required to level up from this level
func CalculateXPForLevel(level int) int {
	// Enforce level cap to prevent overflow and maintain game balance
	if level >= maxLevel {
		return maxLevelXPRequired
	}

	// Apply polynomial formula for smooth scaling
	// Using float64 for precision, then converting to int
	xpFloat := float64(baseXPPerLevel) * float64(level) * (1.0 + float64(level)/levelScaleFactor)
	xp := int(math.Round(xpFloat))

	// Safety check to prevent negative or zero values
	if xp < baseXPPerLevel {
		return baseXPPerLevel
	}

	return xp
}

// CalculateCommitXP calculates the base XP reward for making a commit.
// This is the primary way players earn XP during normal development.
//
// The calculation:
//  1. Base XP: 10 points (every commit counts!)
//  2. Lines bonus: 1 XP per line changed (added + removed)
//  3. Lines bonus is capped at 50 XP to prevent "farming" large auto-generated changes
//
// Examples:
//   - Empty commit (0 lines): 10 XP
//   - Small commit (20 lines): 30 XP (10 base + 20 lines)
//   - Medium commit (50 lines): 60 XP (10 base + 50 capped)
//   - Large commit (500 lines): 60 XP (10 base + 50 capped)
//
// This design encourages:
//   - Regular commits (base XP always awarded)
//   - Meaningful changes (lines bonus)
//   - Good practices (cap prevents gaming the system)
//
// Parameters:
//   - linesAdded: Number of lines added in the commit
//   - linesRemoved: Number of lines removed in the commit
//
// Returns:
//   - int: The base XP earned from this commit (before multipliers)
func CalculateCommitXP(linesAdded, linesRemoved int) int {
	// Ensure non-negative values
	if linesAdded < 0 {
		linesAdded = 0
	}
	if linesRemoved < 0 {
		linesRemoved = 0
	}

	// Calculate total lines changed
	totalLinesChanged := linesAdded + linesRemoved

	// Calculate lines bonus with cap
	linesBonus := totalLinesChanged * xpPerLinesChanged
	if linesBonus > maxLinesBonus {
		linesBonus = maxLinesBonus
	}

	return baseCommitXP + linesBonus
}

// ApplyDifficultyMultiplier scales XP based on the game difficulty setting.
// This affects ALL XP gains: commits, quests, bonuses, everything.
//
// Multipliers:
//   - Easy (1.2x): For beginners or casual play - faster progression
//   - Normal (1.0x): Balanced for most players - standard progression
//   - Hard (0.8x): For veterans seeking a challenge - slower progression
//
// Examples:
//   - 50 XP on Easy: 60 XP (50 * 1.2)
//   - 50 XP on Normal: 50 XP (50 * 1.0)
//   - 50 XP on Hard: 40 XP (50 * 0.8)
//
// Parameters:
//   - baseXP: The XP amount before difficulty scaling
//   - difficulty: The game difficulty ("easy", "normal", "hard")
//
// Returns:
//   - int: The XP after applying difficulty multiplier
func ApplyDifficultyMultiplier(baseXP int, difficulty string) int {
	var multiplier float64

	switch difficulty {
	case DifficultyEasy:
		multiplier = difficultyEasy
	case DifficultyHard:
		multiplier = difficultyHard
	case DifficultyNormal:
		fallthrough // Normal is default
	default:
		multiplier = difficultyNormal
	}

	// Apply multiplier and round to nearest integer
	adjustedXP := float64(baseXP) * multiplier
	return int(math.Round(adjustedXP))
}

// ApplyWisdomBonus applies the character's Wisdom stat bonus to XP gains.
// Wisdom is one of the three core RPG stats and directly increases XP acquisition rate.
//
// The formula: baseXP * (1 + (wisdom - 10) * 0.01)
//
// Wisdom acts as an XP multiplier:
//   - 10 Wisdom (starting): 1.00x (no bonus)
//   - 15 Wisdom: 1.05x (+5% more XP)
//   - 20 Wisdom: 1.10x (+10% more XP)
//   - 30 Wisdom: 1.20x (+20% more XP)
//   - 50 Wisdom: 1.40x (+40% more XP)
//   - 100 Wisdom: 1.90x (+90% more XP - nearly double!)
//
// Since wisdom increases by 1 per level, this creates a positive feedback loop:
// Higher level → More wisdom → Faster XP gain → Faster leveling
//
// This is intentional design to:
//  1. Make leveling feel progressively more rewarding
//  2. Offset the increasing XP requirements at higher levels
//  3. Reward long-term commitment to the game
//
// Examples:
//   - 50 XP with 10 wisdom: 50 XP (no bonus)
//   - 50 XP with 20 wisdom: 55 XP (+10% bonus)
//   - 50 XP with 50 wisdom: 70 XP (+40% bonus)
//
// Parameters:
//   - baseXP: The XP amount before wisdom bonus
//   - wisdom: The character's Wisdom stat value
//
// Returns:
//   - int: The XP after applying wisdom bonus
func ApplyWisdomBonus(baseXP int, wisdom int) int {
	// Calculate wisdom bonus percentage
	// Formula: 1 + (wisdom - base) * bonusPerPoint
	wisdomMultiplier := 1.0 + float64(wisdom-wisdomBaseValue)*wisdomBonusPer10

	// Ensure multiplier is never negative (minimum 1.0x)
	if wisdomMultiplier < 1.0 {
		wisdomMultiplier = 1.0
	}

	// Apply multiplier and round to nearest integer
	bonusXP := float64(baseXP) * wisdomMultiplier
	return int(math.Round(bonusXP))
}

// CalculateQuestReward calculates the XP reward for completing a quest.
// Quest XP is awarded in addition to any commit XP earned while working on the quest.
//
// Quest Difficulty Tiers:
//   - Simple (50 XP): Quick objectives like "Make 3 commits"
//   - Medium (150 XP): Moderate challenges like "Add 100 lines"
//   - Hard (300 XP): Significant goals like "Maintain 7 day streak"
//   - Epic (1000 XP): Major achievements like "Reach level 25"
//
// Note: Quest XP should still be scaled by difficulty multiplier and wisdom bonus
// when actually awarded to the player. This function returns the base reward only.
//
// Parameters:
//   - questDifficulty: The difficulty tier of the quest
//
// Returns:
//   - int: The base XP reward for completing this quest
func CalculateQuestReward(questDifficulty string) int {
	switch questDifficulty {
	case QuestDifficultySimple:
		return questXPSimple
	case QuestDifficultyMedium:
		return questXPMedium
	case QuestDifficultyHard:
		return questXPHard
	case QuestDifficultyEpic:
		return questXPEpic
	default:
		// Default to simple if unknown difficulty
		return questXPSimple
	}
}

// GetLevelFromXP calculates what level a character should be based on total XP earned.
// This is useful for validating character data or recalculating level after XP adjustments.
//
// The function iterates through levels, subtracting the XP cost of each level until
// there isn't enough XP remaining for the next level.
//
// Example:
//   - 0 XP: Level 1
//   - 150 XP: Level 2 (100 for level 2, 50 remaining)
//   - 350 XP: Level 3 (100 for level 2, 210 for level 3, 40 remaining)
//
// Parameters:
//   - totalXP: The total amount of XP accumulated
//
// Returns:
//   - level: The level corresponding to that XP amount
//   - remainingXP: XP left over after reaching that level
func GetLevelFromXP(totalXP int) (level int, remainingXP int) {
	level = 1
	remainingXP = totalXP

	// Keep leveling up while there's enough XP for the next level
	for level < maxLevel {
		xpNeeded := CalculateXPForLevel(level)
		if remainingXP >= xpNeeded {
			remainingXP -= xpNeeded
			level++
		} else {
			break
		}
	}

	return level, remainingXP
}

// GetProgressToNextLevel calculates the progress percentage toward the next level.
// This returns a value between 0.0 and 1.0 suitable for progress bars.
//
// Examples:
//   - 0 XP of 100 needed: 0.0 (0%)
//   - 50 XP of 100 needed: 0.5 (50%)
//   - 99 XP of 100 needed: 0.99 (99%)
//   - 100+ XP of 100 needed: 1.0 (100% - ready to level!)
//
// Parameters:
//   - currentXP: Current XP toward next level
//   - xpNeeded: Total XP required for next level
//
// Returns:
//   - float64: Progress as a decimal between 0.0 and 1.0
func GetProgressToNextLevel(currentXP, xpNeeded int) float64 {
	// Handle edge cases
	if xpNeeded <= 0 {
		return 1.0 // Already at or past the threshold
	}
	if currentXP <= 0 {
		return 0.0 // No progress yet
	}

	// Calculate percentage
	progress := float64(currentXP) / float64(xpNeeded)

	// Clamp to 0.0-1.0 range
	if progress < 0.0 {
		return 0.0
	}
	if progress > 1.0 {
		return 1.0
	}

	return progress
}

// GetTotalXPForLevel calculates the cumulative XP needed to reach a specific level from level 1.
// This is useful for leaderboards, statistics, or comparing character progression.
//
// Example:
//   - Level 1: 0 XP (starting point)
//   - Level 2: 100 XP (sum of level 1→2)
//   - Level 3: 310 XP (sum of 1→2 + 2→3)
//   - Level 10: 5500 XP (sum of all levels 1→10)
//
// Parameters:
//   - targetLevel: The level to calculate cumulative XP for
//
// Returns:
//   - int: Total XP needed to reach that level from level 1
func GetTotalXPForLevel(targetLevel int) int {
	if targetLevel <= 1 {
		return 0
	}

	// Cap at max level
	if targetLevel > maxLevel {
		targetLevel = maxLevel
	}

	totalXP := 0
	for level := 1; level < targetLevel; level++ {
		totalXP += CalculateXPForLevel(level)
	}

	return totalXP
}
