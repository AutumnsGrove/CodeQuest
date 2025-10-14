// Package game contains the quest system for CodeQuest.
// Quests represent coding challenges that players can complete for XP rewards.
package game

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// QuestStatus represents the current state of a quest in its lifecycle.
type QuestStatus string

const (
	QuestAvailable QuestStatus = "available" // Quest can be started
	QuestActive    QuestStatus = "active"    // Quest is currently in progress
	QuestCompleted QuestStatus = "completed" // Quest has been successfully finished
	QuestFailed    QuestStatus = "failed"    // Quest was abandoned or failed
)

// QuestType categorizes quests by their completion criteria.
// Different quest types track different developer activities.
type QuestType string

const (
	QuestTypeCommit    QuestType = "commit"    // Make N commits
	QuestTypeLines     QuestType = "lines"     // Add/modify N lines of code
	QuestTypeTests     QuestType = "tests"     // Add N test cases (post-MVP)
	QuestTypePR        QuestType = "pr"        // Create/merge pull request (post-MVP)
	QuestTypeRefactor  QuestType = "refactor"  // Refactor code (post-MVP)
	QuestTypeDaily     QuestType = "daily"     // Daily quest (post-MVP)
	QuestTypeStreak    QuestType = "streak"    // Maintain N-day streak (post-MVP)
)

// Quest represents a coding task or challenge that players can accept and complete.
// Quests are the primary way players earn large XP rewards and unlock new content.
type Quest struct {
	// Identity - Basic quest information
	ID          string    `json:"id"`          // Unique identifier (UUID)
	Title       string    `json:"title"`       // Display name for the quest
	Description string    `json:"description"` // Detailed description of what to do
	Type        QuestType `json:"type"`        // Type of quest (commit, lines, etc.)

	// Requirements - Prerequisites to start the quest
	RequiredLevel int      `json:"required_level"` // Minimum character level required
	Prerequisites []string `json:"prerequisites"`  // Quest IDs that must be completed first

	// Objectives - What the player needs to accomplish
	Target  int `json:"target"`  // Target count (e.g., 5 commits, 100 lines)
	Current int `json:"current"` // Current progress toward target

	// Rewards - What the player earns upon completion
	XPReward      int      `json:"xp_reward"`                 // Base XP awarded (before multipliers)
	UnlocksSkills []string `json:"unlocks_skills,omitempty"`  // Skills unlocked (post-MVP)
	UnlocksQuests []string `json:"unlocks_quests,omitempty"`  // Quests unlocked (post-MVP)

	// Tracking - Git repository context for the quest
	GitRepo    string `json:"git_repo,omitempty"`     // Path to the git repository
	GitBaseSHA string `json:"git_base_sha,omitempty"` // Starting commit SHA

	// Status - Current state and progress
	Status      QuestStatus `json:"status"`                  // Current lifecycle state
	StartedAt   *time.Time  `json:"started_at,omitempty"`    // When quest was started
	CompletedAt *time.Time  `json:"completed_at,omitempty"`  // When quest was completed
	Progress    float64     `json:"progress"`                // Progress percentage (0.0 to 1.0)
}

// NewQuest creates a new quest with the given parameters.
// This initializes a quest in the "available" state, ready to be started.
//
// Parameters:
//   - title: Display name for the quest
//   - description: What the player needs to do
//   - questType: Type of quest (commit, lines, etc.)
//   - target: Target count to reach
//   - xpReward: Base XP reward for completion
//   - requiredLevel: Minimum level needed to start
//
// Returns:
//   - *Quest: A pointer to the newly created quest
func NewQuest(title, description string, questType QuestType, target, xpReward, requiredLevel int) *Quest {
	return &Quest{
		// Identity
		ID:          generateQuestID(),
		Title:       title,
		Description: description,
		Type:        questType,

		// Requirements
		RequiredLevel: requiredLevel,
		Prerequisites: []string{}, // Empty by default

		// Objectives
		Target:  target,
		Current: 0,

		// Rewards
		XPReward:      xpReward,
		UnlocksSkills: []string{}, // Empty by default
		UnlocksQuests: []string{}, // Empty by default

		// Tracking
		GitRepo:    "",
		GitBaseSHA: "",

		// Status
		Status:      QuestAvailable,
		StartedAt:   nil,
		CompletedAt: nil,
		Progress:    0.0,
	}
}

// IsAvailable checks if the quest can be started by the given character.
// A quest is available if:
//  1. The quest status is "available" (not already started/completed)
//  2. The character meets the minimum level requirement
//  3. All prerequisite quests have been completed
//
// Parameters:
//   - character: The character attempting to start the quest
//
// Returns:
//   - bool: true if the quest can be started, false otherwise
func (q *Quest) IsAvailable(character *Character) bool {
	// Quest must be in available status
	if q.Status != QuestAvailable {
		return false
	}

	// Character must meet level requirement
	if character.Level < q.RequiredLevel {
		return false
	}

	// TODO: Check prerequisites when quest system is more complete
	// For MVP, we'll skip prerequisite checking

	return true
}

// Start begins the quest, marking it as active.
// This should be called when a player accepts a quest.
//
// Parameters:
//   - repoPath: Path to the git repository (optional, empty string if not needed)
//   - baseSHA: Starting commit SHA (optional, empty string if not needed)
//
// Returns:
//   - error: An error if the quest cannot be started
func (q *Quest) Start(repoPath, baseSHA string) error {
	// Verify quest is available
	if q.Status != QuestAvailable {
		return fmt.Errorf("quest %s is not available (current status: %s)", q.ID, q.Status)
	}

	// Set quest to active
	now := time.Now()
	q.Status = QuestActive
	q.StartedAt = &now
	q.GitRepo = repoPath
	q.GitBaseSHA = baseSHA
	q.Current = 0
	q.Progress = 0.0

	return nil
}

// UpdateProgress increments the quest's progress by the given amount.
// This should be called whenever the player makes progress toward the quest objective
// (e.g., makes a commit for a commit quest, adds lines for a lines quest).
//
// The progress is automatically clamped to not exceed the target.
//
// Parameters:
//   - amount: The amount to add to current progress
func (q *Quest) UpdateProgress(amount int) {
	// Only update if quest is active
	if q.Status != QuestActive {
		return
	}

	// Ensure amount is positive
	if amount <= 0 {
		return
	}

	// Add to current progress
	q.Current += amount

	// Clamp to target (don't overshoot)
	if q.Current > q.Target {
		q.Current = q.Target
	}

	// Recalculate progress percentage
	if q.Target > 0 {
		q.Progress = float64(q.Current) / float64(q.Target)
	} else {
		q.Progress = 1.0 // Prevent division by zero
	}

	// Clamp progress to 0.0-1.0 range
	if q.Progress < 0.0 {
		q.Progress = 0.0
	}
	if q.Progress > 1.0 {
		q.Progress = 1.0
	}
}

// CheckCompletion determines if the quest has been completed.
// A quest is complete when the current progress reaches or exceeds the target.
//
// Returns:
//   - bool: true if the quest objectives are met, false otherwise
func (q *Quest) CheckCompletion() bool {
	// Quest must be active to be completed
	if q.Status != QuestActive {
		return false
	}

	// Check if target has been reached
	return q.Current >= q.Target
}

// Complete marks the quest as finished and records the completion time.
// This should be called after CheckCompletion() returns true.
//
// Returns:
//   - error: An error if the quest cannot be completed
func (q *Quest) Complete() error {
	// Verify quest is active
	if q.Status != QuestActive {
		return fmt.Errorf("quest %s is not active (current status: %s)", q.ID, q.Status)
	}

	// Verify objectives are met
	if q.Current < q.Target {
		return fmt.Errorf("quest %s objectives not met (%d/%d)", q.ID, q.Current, q.Target)
	}

	// Mark as completed
	now := time.Now()
	q.Status = QuestCompleted
	q.CompletedAt = &now
	q.Progress = 1.0

	return nil
}

// Fail marks the quest as failed or abandoned.
// This might be used for time-limited quests or when a player wants to abandon a quest.
//
// Returns:
//   - error: An error if the quest cannot be failed
func (q *Quest) Fail() error {
	// Can only fail an active quest
	if q.Status != QuestActive {
		return fmt.Errorf("quest %s is not active (current status: %s)", q.ID, q.Status)
	}

	// Mark as failed
	q.Status = QuestFailed

	return nil
}

// Reset resets the quest to its initial "available" state.
// This clears all progress and allows the quest to be started again.
// Useful for daily quests or repeatable quests.
func (q *Quest) Reset() {
	q.Status = QuestAvailable
	q.Current = 0
	q.Progress = 0.0
	q.StartedAt = nil
	q.CompletedAt = nil
	q.GitRepo = ""
	q.GitBaseSHA = ""
}

// generateQuestID creates a unique identifier for a quest using UUID v4.
//
// Returns:
//   - string: A UUID string in the format "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
func generateQuestID() string {
	return uuid.New().String()
}
