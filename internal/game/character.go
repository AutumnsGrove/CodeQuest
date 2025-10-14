// Package game contains the core game logic for CodeQuest
package game

import (
	"time"

	"github.com/google/uuid"
)

// Character represents the player in the game world.
// This is the core model that tracks all player progress, stats, and achievements.
type Character struct {
	// Identity - Basic character information
	ID        string    `json:"id"`         // Unique identifier (UUID)
	Name      string    `json:"name"`       // Player's chosen character name
	CreatedAt time.Time `json:"created_at"` // When the character was created

	// Core Stats - Primary progression metrics
	Level         int `json:"level"`            // Current character level
	XP            int `json:"xp"`               // Current experience points
	XPToNextLevel int `json:"xp_to_next_level"` // XP required to reach next level

	// RPG Stats - Character attributes that affect gameplay
	CodePower int `json:"code_power"` // Increases commit quality bonus (damage output)
	Wisdom    int `json:"wisdom"`     // Increases XP gain (experience multiplier)
	Agility   int `json:"agility"`    // Faster quest completion bonuses (speed)

	// Progress Tracking - Lifetime statistics
	TotalCommits      int       `json:"total_commits"`       // All-time commit count
	TotalLinesAdded   int       `json:"total_lines_added"`   // All-time lines of code added
	TotalLinesRemoved int       `json:"total_lines_removed"` // All-time lines of code removed
	QuestsCompleted   int       `json:"quests_completed"`    // Total quests completed
	CurrentStreak     int       `json:"current_streak"`      // Consecutive days of activity
	LongestStreak     int       `json:"longest_streak"`      // Best streak ever achieved
	LastActiveDate    time.Time `json:"last_active_date"`    // Last day the player was active

	// Session Stats - Today's activity (resets daily)
	TodayCommits     int           `json:"today_commits"`      // Commits made today
	TodayLinesAdded  int           `json:"today_lines_added"`  // Lines added today
	TodaySessionTime time.Duration `json:"today_session_time"` // Time spent coding today
}

// NewCharacter creates a new character with starting stats.
// This initializes a fresh character with base values suitable for level 1.
//
// Parameters:
//   - name: The character's display name
//
// Returns:
//   - *Character: A pointer to the newly created character
func NewCharacter(name string) *Character {
	now := time.Now()
	startingLevel := 1

	return &Character{
		// Identity
		ID:        generateID(),
		Name:      name,
		CreatedAt: now,

		// Core Stats
		Level:         startingLevel,
		XP:            0,
		XPToNextLevel: calculateXPForLevel(startingLevel), // Placeholder formula

		// RPG Stats - Starting values for a level 1 character
		CodePower: 10, // Base power
		Wisdom:    10, // Base wisdom
		Agility:   10, // Base agility

		// Progress Tracking - All zeroed out for new character
		TotalCommits:      0,
		TotalLinesAdded:   0,
		TotalLinesRemoved: 0,
		QuestsCompleted:   0,
		CurrentStreak:     0,
		LongestStreak:     0,
		LastActiveDate:    now,

		// Session Stats - Start with clean slate
		TodayCommits:     0,
		TodayLinesAdded:  0,
		TodaySessionTime: 0,
	}
}

// AddXP adds experience points and handles level-ups.
// This method can trigger multiple level-ups if enough XP is added at once.
//
// Parameters:
//   - amount: The amount of XP to add (must be positive)
//
// Returns:
//   - bool: true if at least one level-up occurred, false otherwise
func (c *Character) AddXP(amount int) bool {
	if amount <= 0 {
		return false
	}

	c.XP += amount
	leveledUp := false

	// Keep leveling up while we have enough XP
	// This handles cases where a single action gives enough XP for multiple levels
	for c.XP >= c.XPToNextLevel {
		// Subtract the XP cost of the level we just gained
		c.XP -= c.XPToNextLevel
		c.Level++
		leveledUp = true

		// Recalculate XP needed for the new level
		c.XPToNextLevel = calculateXPForLevel(c.Level)

		// On level up, grant small stat increases
		// This makes leveling feel rewarding beyond just the level number
		c.CodePower++
		c.Wisdom++
		c.Agility++
	}

	return leveledUp
}

// UpdateStreak updates the daily activity streak counter.
// This should be called whenever the player performs an activity (like making a commit).
// It maintains both the current streak and tracks the longest streak achieved.
func (c *Character) UpdateStreak() {
	now := time.Now()
	today := truncateToDay(now)
	lastActive := truncateToDay(c.LastActiveDate)

	// Calculate the difference in days
	daysDiff := int(today.Sub(lastActive).Hours() / 24)

	switch {
	case daysDiff == 0:
		// Already active today
		// If streak is 0 (brand new character or after reset), start it
		if c.CurrentStreak == 0 {
			c.CurrentStreak = 1
		}
		// Otherwise no change - already counted today

	case daysDiff == 1:
		// Active yesterday, increment streak
		c.CurrentStreak++

	default:
		// Missed a day (or more), reset streak to 1
		c.CurrentStreak = 1
	}

	// Update longest streak if current streak is now the best
	if c.CurrentStreak > c.LongestStreak {
		c.LongestStreak = c.CurrentStreak
	}

	// Update last active date to today
	c.LastActiveDate = now
}

// ResetDailyStats resets all today's statistics to zero.
// This should be called at the start of each new day.
func (c *Character) ResetDailyStats() {
	c.TodayCommits = 0
	c.TodayLinesAdded = 0
	c.TodaySessionTime = 0
}

// IsToday checks if the given time is the same day as now.
// This is useful for determining if daily stats should be reset.
//
// Parameters:
//   - t: The time to check
//
// Returns:
//   - bool: true if t is today, false otherwise
func (c *Character) IsToday(t time.Time) bool {
	now := time.Now()
	return truncateToDay(now).Equal(truncateToDay(t))
}

// generateID creates a unique identifier for the character using UUID v4.
// This ensures each character has a globally unique ID.
//
// Returns:
//   - string: A UUID string in the format "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
func generateID() string {
	return uuid.New().String()
}

// calculateXPForLevel calculates the XP required to reach the next level.
// This wraps the engine's CalculateXPForLevel function for internal use.
//
// The XP curve uses a polynomial formula that creates smooth, engaging progression:
// - Early levels (1-10): Fast progression to hook new players
// - Mid levels (10-50): Steady, satisfying growth
// - Late levels (50-100): Challenging but achievable
//
// See engine.go for detailed progression curve and balance rationale.
//
// Parameters:
//   - level: The current level
//
// Returns:
//   - int: The amount of XP needed to level up from this level
func calculateXPForLevel(level int) int {
	return CalculateXPForLevel(level)
}

// truncateToDay truncates a time to midnight (start of day) in local timezone.
// This is used for comparing dates without considering the time of day.
//
// Parameters:
//   - t: The time to truncate
//
// Returns:
//   - time.Time: The time truncated to midnight
func truncateToDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
