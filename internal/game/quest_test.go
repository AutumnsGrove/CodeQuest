package game

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestNewQuest tests the quest creation function
func TestNewQuest(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		description   string
		questType     QuestType
		target        int
		xpReward      int
		requiredLevel int
	}{
		{
			name:          "commit quest",
			title:         "Make 5 Commits",
			description:   "Commit your code 5 times",
			questType:     QuestTypeCommit,
			target:        5,
			xpReward:      100,
			requiredLevel: 1,
		},
		{
			name:          "lines quest",
			title:         "Add 100 Lines",
			description:   "Add 100 lines of code",
			questType:     QuestTypeLines,
			target:        100,
			xpReward:      200,
			requiredLevel: 5,
		},
		{
			name:          "zero target",
			title:         "Test Quest",
			description:   "A quest with zero target",
			questType:     QuestTypeCommit,
			target:        0,
			xpReward:      50,
			requiredLevel: 1,
		},
		{
			name:          "high level requirement",
			title:         "Advanced Quest",
			description:   "For high level characters",
			questType:     QuestTypeTests,
			target:        50,
			xpReward:      1000,
			requiredLevel: 50,
		},
		{
			name:          "empty strings",
			title:         "",
			description:   "",
			questType:     QuestTypeCommit,
			target:        10,
			xpReward:      75,
			requiredLevel: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quest := NewQuest(tt.title, tt.description, tt.questType, tt.target, tt.xpReward, tt.requiredLevel)

			// Verify ID is a valid UUID
			if _, err := uuid.Parse(quest.ID); err != nil {
				t.Errorf("NewQuest() ID is not a valid UUID: %v", err)
			}

			// Verify basic fields
			if quest.Title != tt.title {
				t.Errorf("NewQuest() Title = %v, want %v", quest.Title, tt.title)
			}
			if quest.Description != tt.description {
				t.Errorf("NewQuest() Description = %v, want %v", quest.Description, tt.description)
			}
			if quest.Type != tt.questType {
				t.Errorf("NewQuest() Type = %v, want %v", quest.Type, tt.questType)
			}

			// Verify requirements
			if quest.RequiredLevel != tt.requiredLevel {
				t.Errorf("NewQuest() RequiredLevel = %v, want %v", quest.RequiredLevel, tt.requiredLevel)
			}
			if quest.Prerequisites == nil {
				t.Errorf("NewQuest() Prerequisites is nil, should be empty slice")
			}
			if len(quest.Prerequisites) != 0 {
				t.Errorf("NewQuest() Prerequisites length = %v, want 0", len(quest.Prerequisites))
			}

			// Verify objectives
			if quest.Target != tt.target {
				t.Errorf("NewQuest() Target = %v, want %v", quest.Target, tt.target)
			}
			if quest.Current != 0 {
				t.Errorf("NewQuest() Current = %v, want 0", quest.Current)
			}

			// Verify rewards
			if quest.XPReward != tt.xpReward {
				t.Errorf("NewQuest() XPReward = %v, want %v", quest.XPReward, tt.xpReward)
			}
			if quest.UnlocksSkills == nil {
				t.Errorf("NewQuest() UnlocksSkills is nil, should be empty slice")
			}
			if quest.UnlocksQuests == nil {
				t.Errorf("NewQuest() UnlocksQuests is nil, should be empty slice")
			}

			// Verify tracking fields are empty
			if quest.GitRepo != "" {
				t.Errorf("NewQuest() GitRepo = %v, want empty string", quest.GitRepo)
			}
			if quest.GitBaseSHA != "" {
				t.Errorf("NewQuest() GitBaseSHA = %v, want empty string", quest.GitBaseSHA)
			}

			// Verify status
			if quest.Status != QuestAvailable {
				t.Errorf("NewQuest() Status = %v, want %v", quest.Status, QuestAvailable)
			}
			if quest.StartedAt != nil {
				t.Errorf("NewQuest() StartedAt = %v, want nil", quest.StartedAt)
			}
			if quest.CompletedAt != nil {
				t.Errorf("NewQuest() CompletedAt = %v, want nil", quest.CompletedAt)
			}
			if quest.Progress != 0.0 {
				t.Errorf("NewQuest() Progress = %v, want 0.0", quest.Progress)
			}
		})
	}
}

// TestQuest_IsAvailable tests the quest availability check
func TestQuest_IsAvailable(t *testing.T) {
	tests := []struct {
		name          string
		questStatus   QuestStatus
		requiredLevel int
		charLevel     int
		want          bool
	}{
		{
			name:          "available quest, level met",
			questStatus:   QuestAvailable,
			requiredLevel: 1,
			charLevel:     1,
			want:          true,
		},
		{
			name:          "available quest, over-leveled",
			questStatus:   QuestAvailable,
			requiredLevel: 5,
			charLevel:     10,
			want:          true,
		},
		{
			name:          "available quest, under-leveled",
			questStatus:   QuestAvailable,
			requiredLevel: 10,
			charLevel:     5,
			want:          false,
		},
		{
			name:          "active quest",
			questStatus:   QuestActive,
			requiredLevel: 1,
			charLevel:     5,
			want:          false,
		},
		{
			name:          "completed quest",
			questStatus:   QuestCompleted,
			requiredLevel: 1,
			charLevel:     5,
			want:          false,
		},
		{
			name:          "failed quest",
			questStatus:   QuestFailed,
			requiredLevel: 1,
			charLevel:     5,
			want:          false,
		},
		{
			name:          "level exactly met",
			questStatus:   QuestAvailable,
			requiredLevel: 25,
			charLevel:     25,
			want:          true,
		},
		{
			name:          "level one below requirement",
			questStatus:   QuestAvailable,
			requiredLevel: 10,
			charLevel:     9,
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quest := NewQuest("Test Quest", "Test Description", QuestTypeCommit, 5, 100, tt.requiredLevel)
			quest.Status = tt.questStatus

			char := NewCharacter("TestHero")
			char.Level = tt.charLevel

			got := quest.IsAvailable(char)
			if got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestQuest_Start tests starting a quest
func TestQuest_Start(t *testing.T) {
	tests := []struct {
		name       string
		status     QuestStatus
		repoPath   string
		baseSHA    string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:     "start available quest",
			status:   QuestAvailable,
			repoPath: "/path/to/repo",
			baseSHA:  "abc123",
			wantErr:  false,
		},
		{
			name:     "start available quest no repo",
			status:   QuestAvailable,
			repoPath: "",
			baseSHA:  "",
			wantErr:  false,
		},
		{
			name:       "start active quest",
			status:     QuestActive,
			repoPath:   "/path/to/repo",
			baseSHA:    "abc123",
			wantErr:    true,
			wantErrMsg: "not available",
		},
		{
			name:       "start completed quest",
			status:     QuestCompleted,
			repoPath:   "/path/to/repo",
			baseSHA:    "abc123",
			wantErr:    true,
			wantErrMsg: "not available",
		},
		{
			name:       "start failed quest",
			status:     QuestFailed,
			repoPath:   "/path/to/repo",
			baseSHA:    "abc123",
			wantErr:    true,
			wantErrMsg: "not available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quest := NewQuest("Test Quest", "Test Description", QuestTypeCommit, 5, 100, 1)
			quest.Status = tt.status

			err := quest.Start(tt.repoPath, tt.baseSHA)

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Errorf("Start() error = nil, want error containing %q", tt.wantErrMsg)
				} else if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("Start() error = %q, want error containing %q", err.Error(), tt.wantErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Start() unexpected error = %v", err)
				}

				// Verify quest state after successful start
				if quest.Status != QuestActive {
					t.Errorf("Start() Status = %v, want %v", quest.Status, QuestActive)
				}
				if quest.StartedAt == nil {
					t.Errorf("Start() StartedAt is nil, should be set")
				} else {
					// Verify StartedAt is recent (within last second)
					if time.Since(*quest.StartedAt) > time.Second {
						t.Errorf("Start() StartedAt = %v, should be recent", quest.StartedAt)
					}
				}
				if quest.GitRepo != tt.repoPath {
					t.Errorf("Start() GitRepo = %v, want %v", quest.GitRepo, tt.repoPath)
				}
				if quest.GitBaseSHA != tt.baseSHA {
					t.Errorf("Start() GitBaseSHA = %v, want %v", quest.GitBaseSHA, tt.baseSHA)
				}
				if quest.Current != 0 {
					t.Errorf("Start() Current = %v, want 0", quest.Current)
				}
				if quest.Progress != 0.0 {
					t.Errorf("Start() Progress = %v, want 0.0", quest.Progress)
				}
			}
		})
	}
}

// TestQuest_UpdateProgress tests quest progress tracking
func TestQuest_UpdateProgress(t *testing.T) {
	tests := []struct {
		name           string
		status         QuestStatus
		target         int
		initialCurrent int
		amount         int
		wantCurrent    int
		wantProgress   float64
	}{
		{
			name:           "add progress to active quest",
			status:         QuestActive,
			target:         10,
			initialCurrent: 0,
			amount:         3,
			wantCurrent:    3,
			wantProgress:   0.3,
		},
		{
			name:           "add progress multiple times",
			status:         QuestActive,
			target:         100,
			initialCurrent: 50,
			amount:         25,
			wantCurrent:    75,
			wantProgress:   0.75,
		},
		{
			name:           "reach target exactly",
			status:         QuestActive,
			target:         10,
			initialCurrent: 8,
			amount:         2,
			wantCurrent:    10,
			wantProgress:   1.0,
		},
		{
			name:           "exceed target (should clamp)",
			status:         QuestActive,
			target:         10,
			initialCurrent: 8,
			amount:         5,
			wantCurrent:    10,
			wantProgress:   1.0,
		},
		{
			name:           "zero amount (no change)",
			status:         QuestActive,
			target:         10,
			initialCurrent: 5,
			amount:         0,
			wantCurrent:    5,
			wantProgress:   0.5,
		},
		{
			name:           "negative amount (no change)",
			status:         QuestActive,
			target:         10,
			initialCurrent: 5,
			amount:         -3,
			wantCurrent:    5,
			wantProgress:   0.5,
		},
		{
			name:           "update available quest (no change)",
			status:         QuestAvailable,
			target:         10,
			initialCurrent: 0,
			amount:         5,
			wantCurrent:    0,
			wantProgress:   0.0,
		},
		{
			name:           "update completed quest (no change)",
			status:         QuestCompleted,
			target:         10,
			initialCurrent: 10,
			amount:         5,
			wantCurrent:    10,
			wantProgress:   1.0,
		},
		{
			name:           "zero target edge case",
			status:         QuestActive,
			target:         0,
			initialCurrent: 0,
			amount:         5,
			wantCurrent:    0,
			wantProgress:   1.0,
		},
		{
			name:           "large numbers",
			status:         QuestActive,
			target:         10000,
			initialCurrent: 5000,
			amount:         2500,
			wantCurrent:    7500,
			wantProgress:   0.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quest := NewQuest("Test Quest", "Test Description", QuestTypeCommit, tt.target, 100, 1)
			quest.Status = tt.status
			quest.Current = tt.initialCurrent
			if tt.target > 0 {
				quest.Progress = float64(tt.initialCurrent) / float64(tt.target)
			}

			quest.UpdateProgress(tt.amount)

			// Verify current progress
			if quest.Current != tt.wantCurrent {
				t.Errorf("UpdateProgress() Current = %v, want %v", quest.Current, tt.wantCurrent)
			}

			// Verify progress percentage (with small tolerance for float comparison)
			tolerance := 0.001
			if abs(quest.Progress-tt.wantProgress) > tolerance {
				t.Errorf("UpdateProgress() Progress = %v, want %v", quest.Progress, tt.wantProgress)
			}

			// Verify progress is always between 0.0 and 1.0
			if quest.Progress < 0.0 || quest.Progress > 1.0 {
				t.Errorf("UpdateProgress() Progress = %v, should be between 0.0 and 1.0", quest.Progress)
			}
		})
	}
}

// TestQuest_CheckCompletion tests completion detection
func TestQuest_CheckCompletion(t *testing.T) {
	tests := []struct {
		name    string
		status  QuestStatus
		target  int
		current int
		want    bool
	}{
		{
			name:    "active quest, target reached",
			status:  QuestActive,
			target:  10,
			current: 10,
			want:    true,
		},
		{
			name:    "active quest, target exceeded",
			status:  QuestActive,
			target:  10,
			current: 15,
			want:    true,
		},
		{
			name:    "active quest, not yet complete",
			status:  QuestActive,
			target:  10,
			current: 9,
			want:    false,
		},
		{
			name:    "active quest, no progress",
			status:  QuestActive,
			target:  10,
			current: 0,
			want:    false,
		},
		{
			name:    "available quest (not started)",
			status:  QuestAvailable,
			target:  10,
			current: 0,
			want:    false,
		},
		{
			name:    "completed quest",
			status:  QuestCompleted,
			target:  10,
			current: 10,
			want:    false,
		},
		{
			name:    "failed quest",
			status:  QuestFailed,
			target:  10,
			current: 5,
			want:    false,
		},
		{
			name:    "zero target (edge case)",
			status:  QuestActive,
			target:  0,
			current: 0,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quest := NewQuest("Test Quest", "Test Description", QuestTypeCommit, tt.target, 100, 1)
			quest.Status = tt.status
			quest.Current = tt.current

			got := quest.CheckCompletion()
			if got != tt.want {
				t.Errorf("CheckCompletion() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestQuest_Complete tests marking a quest as complete
func TestQuest_Complete(t *testing.T) {
	tests := []struct {
		name       string
		status     QuestStatus
		target     int
		current    int
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "complete active quest at target",
			status:  QuestActive,
			target:  10,
			current: 10,
			wantErr: false,
		},
		{
			name:    "complete active quest over target",
			status:  QuestActive,
			target:  10,
			current: 15,
			wantErr: false,
		},
		{
			name:       "complete active quest under target",
			status:     QuestActive,
			target:     10,
			current:    9,
			wantErr:    true,
			wantErrMsg: "objectives not met",
		},
		{
			name:       "complete available quest",
			status:     QuestAvailable,
			target:     10,
			current:    10,
			wantErr:    true,
			wantErrMsg: "not active",
		},
		{
			name:       "complete completed quest",
			status:     QuestCompleted,
			target:     10,
			current:    10,
			wantErr:    true,
			wantErrMsg: "not active",
		},
		{
			name:       "complete failed quest",
			status:     QuestFailed,
			target:     10,
			current:    5,
			wantErr:    true,
			wantErrMsg: "not active",
		},
		{
			name:    "complete quest with zero target",
			status:  QuestActive,
			target:  0,
			current: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quest := NewQuest("Test Quest", "Test Description", QuestTypeCommit, tt.target, 100, 1)
			quest.Status = tt.status
			quest.Current = tt.current

			err := quest.Complete()

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Errorf("Complete() error = nil, want error containing %q", tt.wantErrMsg)
				} else if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("Complete() error = %q, want error containing %q", err.Error(), tt.wantErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Complete() unexpected error = %v", err)
				}

				// Verify quest state after successful completion
				if quest.Status != QuestCompleted {
					t.Errorf("Complete() Status = %v, want %v", quest.Status, QuestCompleted)
				}
				if quest.CompletedAt == nil {
					t.Errorf("Complete() CompletedAt is nil, should be set")
				} else {
					// Verify CompletedAt is recent (within last second)
					if time.Since(*quest.CompletedAt) > time.Second {
						t.Errorf("Complete() CompletedAt = %v, should be recent", quest.CompletedAt)
					}
				}
				if quest.Progress != 1.0 {
					t.Errorf("Complete() Progress = %v, want 1.0", quest.Progress)
				}
			}
		})
	}
}

// TestQuest_Fail tests failing a quest
func TestQuest_Fail(t *testing.T) {
	tests := []struct {
		name       string
		status     QuestStatus
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "fail active quest",
			status:  QuestActive,
			wantErr: false,
		},
		{
			name:       "fail available quest",
			status:     QuestAvailable,
			wantErr:    true,
			wantErrMsg: "not active",
		},
		{
			name:       "fail completed quest",
			status:     QuestCompleted,
			wantErr:    true,
			wantErrMsg: "not active",
		},
		{
			name:       "fail failed quest",
			status:     QuestFailed,
			wantErr:    true,
			wantErrMsg: "not active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quest := NewQuest("Test Quest", "Test Description", QuestTypeCommit, 10, 100, 1)
			quest.Status = tt.status
			if tt.status == QuestActive {
				now := time.Now()
				quest.StartedAt = &now
			}

			err := quest.Fail()

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Errorf("Fail() error = nil, want error containing %q", tt.wantErrMsg)
				} else if tt.wantErrMsg != "" && !contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("Fail() error = %q, want error containing %q", err.Error(), tt.wantErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Fail() unexpected error = %v", err)
				}

				// Verify quest state after successful fail
				if quest.Status != QuestFailed {
					t.Errorf("Fail() Status = %v, want %v", quest.Status, QuestFailed)
				}
			}
		})
	}
}

// TestQuest_Reset tests resetting a quest
func TestQuest_Reset(t *testing.T) {
	tests := []struct {
		name   string
		status QuestStatus
	}{
		{"reset available quest", QuestAvailable},
		{"reset active quest", QuestActive},
		{"reset completed quest", QuestCompleted},
		{"reset failed quest", QuestFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quest := NewQuest("Test Quest", "Test Description", QuestTypeCommit, 10, 100, 1)
			quest.Status = tt.status
			quest.Current = 7
			quest.Progress = 0.7
			now := time.Now()
			quest.StartedAt = &now
			quest.CompletedAt = &now
			quest.GitRepo = "/some/repo"
			quest.GitBaseSHA = "abc123"

			quest.Reset()

			// Verify all progress and state is reset
			if quest.Status != QuestAvailable {
				t.Errorf("Reset() Status = %v, want %v", quest.Status, QuestAvailable)
			}
			if quest.Current != 0 {
				t.Errorf("Reset() Current = %v, want 0", quest.Current)
			}
			if quest.Progress != 0.0 {
				t.Errorf("Reset() Progress = %v, want 0.0", quest.Progress)
			}
			if quest.StartedAt != nil {
				t.Errorf("Reset() StartedAt = %v, want nil", quest.StartedAt)
			}
			if quest.CompletedAt != nil {
				t.Errorf("Reset() CompletedAt = %v, want nil", quest.CompletedAt)
			}
			if quest.GitRepo != "" {
				t.Errorf("Reset() GitRepo = %v, want empty string", quest.GitRepo)
			}
			if quest.GitBaseSHA != "" {
				t.Errorf("Reset() GitBaseSHA = %v, want empty string", quest.GitBaseSHA)
			}

			// Verify quest identity is preserved
			if quest.ID == "" {
				t.Errorf("Reset() ID is empty, should be preserved")
			}
			if quest.Title != "Test Quest" {
				t.Errorf("Reset() Title changed, should be preserved")
			}
			if quest.Target != 10 {
				t.Errorf("Reset() Target changed, should be preserved")
			}
		})
	}
}

// TestGenerateQuestID tests UUID generation for quests
func TestGenerateQuestID(t *testing.T) {
	// Generate multiple IDs and verify they're valid UUIDs and unique
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateQuestID()

		// Verify it's a valid UUID
		if _, err := uuid.Parse(id); err != nil {
			t.Errorf("generateQuestID() produced invalid UUID: %v, error: %v", id, err)
		}

		// Verify it's unique
		if ids[id] {
			t.Errorf("generateQuestID() produced duplicate UUID: %v", id)
		}
		ids[id] = true
	}

	// Verify we generated 100 unique IDs
	if len(ids) != 100 {
		t.Errorf("generateQuestID() uniqueness check failed: got %v unique IDs, want 100", len(ids))
	}
}

// TestQuest_FullLifecycle tests a complete quest lifecycle
func TestQuest_FullLifecycle(t *testing.T) {
	// Create a character
	char := NewCharacter("TestHero")
	char.Level = 5

	// Create a quest
	quest := NewQuest("Make 3 Commits", "Commit your code 3 times", QuestTypeCommit, 3, 150, 5)

	// Step 1: Check initial state
	if quest.Status != QuestAvailable {
		t.Errorf("Initial Status = %v, want %v", quest.Status, QuestAvailable)
	}

	// Step 2: Check availability
	if !quest.IsAvailable(char) {
		t.Errorf("Quest should be available for character level %v", char.Level)
	}

	// Step 3: Start the quest
	err := quest.Start("/path/to/repo", "abc123")
	if err != nil {
		t.Errorf("Start() unexpected error = %v", err)
	}
	if quest.Status != QuestActive {
		t.Errorf("After Start() Status = %v, want %v", quest.Status, QuestActive)
	}

	// Step 4: Make some progress
	quest.UpdateProgress(1) // 1/3 commits
	if quest.Current != 1 {
		t.Errorf("After first update Current = %v, want 1", quest.Current)
	}
	if quest.CheckCompletion() {
		t.Errorf("CheckCompletion() = true, want false (not done yet)")
	}

	quest.UpdateProgress(1) // 2/3 commits
	if quest.Current != 2 {
		t.Errorf("After second update Current = %v, want 2", quest.Current)
	}
	if quest.CheckCompletion() {
		t.Errorf("CheckCompletion() = true, want false (not done yet)")
	}

	quest.UpdateProgress(1) // 3/3 commits
	if quest.Current != 3 {
		t.Errorf("After third update Current = %v, want 3", quest.Current)
	}

	// Step 5: Check completion
	if !quest.CheckCompletion() {
		t.Errorf("CheckCompletion() = false, want true (should be complete)")
	}

	// Step 6: Mark as complete
	err = quest.Complete()
	if err != nil {
		t.Errorf("Complete() unexpected error = %v", err)
	}
	if quest.Status != QuestCompleted {
		t.Errorf("After Complete() Status = %v, want %v", quest.Status, QuestCompleted)
	}
	if quest.Progress != 1.0 {
		t.Errorf("After Complete() Progress = %v, want 1.0", quest.Progress)
	}

	// Step 7: Verify can't start completed quest
	err = quest.Start("/another/repo", "def456")
	if err == nil {
		t.Errorf("Start() on completed quest should return error")
	}

	// Step 8: Reset and verify can start again
	quest.Reset()
	if quest.Status != QuestAvailable {
		t.Errorf("After Reset() Status = %v, want %v", quest.Status, QuestAvailable)
	}
	if quest.Current != 0 {
		t.Errorf("After Reset() Current = %v, want 0", quest.Current)
	}

	err = quest.Start("/new/repo", "ghi789")
	if err != nil {
		t.Errorf("Start() after reset unexpected error = %v", err)
	}
}

// TestQuest_AbandonScenario tests abandoning a quest
func TestQuest_AbandonScenario(t *testing.T) {
	quest := NewQuest("Difficult Quest", "A quest that might be too hard", QuestTypeCommit, 100, 500, 1)

	// Start the quest
	quest.Start("/repo", "sha")

	// Make some progress
	quest.UpdateProgress(10)

	// Player decides to abandon
	err := quest.Fail()
	if err != nil {
		t.Errorf("Fail() unexpected error = %v", err)
	}

	// Verify quest is failed
	if quest.Status != QuestFailed {
		t.Errorf("After Fail() Status = %v, want %v", quest.Status, QuestFailed)
	}

	// Verify progress is preserved (for potential analytics)
	if quest.Current != 10 {
		t.Errorf("After Fail() Current = %v, want 10 (should preserve progress)", quest.Current)
	}

	// Can't complete a failed quest
	err = quest.Complete()
	if err == nil {
		t.Errorf("Complete() on failed quest should return error")
	}

	// Reset allows starting again
	quest.Reset()
	err = quest.Start("/repo", "sha")
	if err != nil {
		t.Errorf("Start() after reset from failed state unexpected error = %v", err)
	}
}

// TestQuest_ProgressClamping tests that progress can't exceed target
func TestQuest_ProgressClamping(t *testing.T) {
	quest := NewQuest("Test Quest", "Test Description", QuestTypeCommit, 10, 100, 1)
	quest.Start("/repo", "sha")

	// Add way more than target
	quest.UpdateProgress(100)

	// Should clamp to target
	if quest.Current != 10 {
		t.Errorf("UpdateProgress() with overshoot Current = %v, want 10 (clamped to target)", quest.Current)
	}
	if quest.Progress != 1.0 {
		t.Errorf("UpdateProgress() with overshoot Progress = %v, want 1.0", quest.Progress)
	}

	// Additional updates shouldn't increase further
	quest.UpdateProgress(50)
	if quest.Current != 10 {
		t.Errorf("UpdateProgress() after clamping Current = %v, want 10 (should stay clamped)", quest.Current)
	}
}

// Helper functions

// contains checks if a string contains a substring (case-sensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
