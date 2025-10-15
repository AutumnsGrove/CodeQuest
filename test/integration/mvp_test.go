// Package integration provides end-to-end integration tests for CodeQuest MVP.
// These tests validate complete user workflows from Git commits to XP gains,
// quest completion, and data persistence.
package integration

import (
	"testing"
	"time"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// TestMVP_CharacterCreation tests the basic character creation flow.
// This validates that users can create a character with proper initialization.
func TestMVP_CharacterCreation(t *testing.T) {
	// Create character
	char := game.NewCharacter("TestHero")

	// Verify character was created with correct initial state
	if char.Name != "TestHero" {
		t.Errorf("Character name = %s, want TestHero", char.Name)
	}

	if char.Level != 1 {
		t.Errorf("Character level = %d, want 1", char.Level)
	}

	if char.XP != 0 {
		t.Errorf("Character XP = %d, want 0", char.XP)
	}

	if char.CodePower != 10 || char.Wisdom != 10 || char.Agility != 10 {
		t.Errorf("Character stats = %d/%d/%d, want 10/10/10",
			char.CodePower, char.Wisdom, char.Agility)
	}

	// Verify ID is generated and not empty
	if char.ID == "" {
		t.Error("Character ID should not be empty")
	}

	// Verify timestamps are recent
	if time.Since(char.CreatedAt) > time.Second {
		t.Errorf("Character CreatedAt is too old: %v", char.CreatedAt)
	}
}

// TestMVP_CommitFlow tests the end-to-end commit detection and XP award flow.
// This simulates: Git commit → Event publication → XP award → Level up → Persistence
func TestMVP_CommitFlow(t *testing.T) {
	// Setup
	char := game.NewCharacter("TestHero")
	storage := &mockStorage{}
	eventBus := game.NewEventBus()

	// Track if XP was awarded
	xpAwarded := false
	initialLevel := char.Level

	// Subscribe to commit events
	eventBus.Subscribe(game.EventCommit, func(e game.Event) {
		// Extract commit data
		linesAdded, _ := e.Data["lines_added"].(int)
		linesRemoved, _ := e.Data["lines_removed"].(int)

		// Award XP based on commit
		baseXP := 10
		linesXP := (linesAdded + linesRemoved) / 2
		if linesXP > 50 {
			linesXP = 50 // Cap bonus
		}
		totalXP := baseXP + linesXP

		// Add XP to character
		char.AddXP(totalXP)
		char.TotalCommits++
		char.TotalLinesAdded += linesAdded
		char.TotalLinesRemoved += linesRemoved

		xpAwarded = true

		// Save character after XP award
		if err := storage.SaveCharacter(char); err != nil {
			t.Errorf("Failed to save character: %v", err)
		}
	})

	// Simulate a commit event
	commitEvent := game.NewCommitEvent(
		"abc123",
		"feat: Add awesome feature",
		3,   // files changed
		100, // lines added
		10,  // lines removed
	)

	eventBus.Publish(commitEvent)

	// Verify XP was awarded
	if !xpAwarded {
		t.Error("Expected XP to be awarded after commit event")
	}

	// Verify character received XP
	if char.XP == 0 {
		t.Error("Character should have received XP")
	}

	// Verify commit was tracked
	if char.TotalCommits != 1 {
		t.Errorf("Character TotalCommits = %d, want 1", char.TotalCommits)
	}

	if char.TotalLinesAdded != 100 {
		t.Errorf("Character TotalLinesAdded = %d, want 100", char.TotalLinesAdded)
	}

	// Verify character was saved
	if !storage.SaveCharacterCalled {
		t.Error("Expected character to be saved after XP award")
	}

	// Verify character can be loaded
	loadedChar, err := storage.LoadCharacter()
	if err != nil {
		t.Fatalf("Failed to load character: %v", err)
	}

	if loadedChar.TotalCommits != 1 {
		t.Errorf("Loaded character TotalCommits = %d, want 1", loadedChar.TotalCommits)
	}

	// Test level-up scenario - publish enough commits to level up
	// Level 1→2 requires 110 XP total
	// First commit gave 65 XP (10 base + 55 from 110 lines)
	// We need at least 45 more XP to level up

	// Publish another commit with enough XP to level up
	bigCommitEvent := game.NewCommitEvent(
		"def456",
		"feat: Major feature",
		10,  // files changed
		100, // lines added (will award 60 XP: 10 base + 50 capped bonus)
		0,   // lines removed
	)

	eventBus.Publish(bigCommitEvent)

	// After 2 commits (65 + 60 = 125 XP), should have leveled up to 2
	if char.Level <= initialLevel {
		t.Errorf("Character should have leveled up, got level %d, expected > %d (total XP earned should be enough)",
			char.Level, initialLevel)
	}
}

// TestMVP_QuestLifecycle tests the complete quest workflow.
// This validates: Quest creation → Start → Progress → Complete → Reward
func TestMVP_QuestLifecycle(t *testing.T) {
	// Setup
	char := game.NewCharacter("TestHero")
	storage := &mockStorage{}
	eventBus := game.NewEventBus()

	// Create a commit quest
	quest := game.NewQuest(
		"First Commits",
		"Make your first 3 commits to learn the basics",
		game.QuestTypeCommit,
		3,  // Target: 3 commits
		50, // Reward: 50 XP
		1,  // Required level: 1
	)

	// Verify quest is available
	if !quest.IsAvailable(char) {
		t.Error("Quest should be available for level 1 character")
	}

	// Start the quest
	if err := quest.Start("", ""); err != nil {
		t.Fatalf("Failed to start quest: %v", err)
	}

	// Verify quest is now active
	if quest.Status != game.QuestActive {
		t.Errorf("Quest status = %s, want %s", quest.Status, game.QuestActive)
	}

	if quest.StartedAt == nil {
		t.Error("Quest StartedAt should be set after starting")
	}

	// Save initial quest state
	if err := storage.SaveQuests([]*game.Quest{quest}); err != nil {
		t.Fatalf("Failed to save quests: %v", err)
	}

	// Subscribe to commit events to update quest progress
	eventBus.Subscribe(game.EventCommit, func(e game.Event) {
		if quest.Status == game.QuestActive && quest.Type == game.QuestTypeCommit {
			quest.UpdateProgress(1) // Each commit adds 1 to progress

			// Check if quest is complete
			if quest.CheckCompletion() {
				if err := quest.Complete(); err != nil {
					t.Errorf("Failed to complete quest: %v", err)
				}

				// Award quest XP to character
				char.AddXP(quest.XPReward)
				char.QuestsCompleted++

				// Publish quest done event
				eventBus.Publish(game.NewQuestDoneEvent(
					quest.ID,
					quest.Title,
					quest.XPReward,
				))

				// Save updated state
				storage.SaveCharacter(char)
				storage.SaveQuests([]*game.Quest{quest})
			}
		}
	})

	// Track quest completion
	questCompleted := false
	eventBus.Subscribe(game.EventQuestDone, func(e game.Event) {
		questCompleted = true
	})

	// Simulate 3 commits
	for i := 1; i <= 3; i++ {
		eventBus.Publish(game.NewCommitEvent(
			"commit"+string(rune(i)),
			"Commit "+string(rune(i)),
			1,
			10,
			0,
		))

		// After each commit, check quest progress
		expectedProgress := float64(i) / 3.0
		if quest.Progress != expectedProgress {
			t.Errorf("After commit %d, quest progress = %.2f, want %.2f",
				i, quest.Progress, expectedProgress)
		}
	}

	// Verify quest is completed
	if quest.Status != game.QuestCompleted {
		t.Errorf("Quest status = %s, want %s", quest.Status, game.QuestCompleted)
	}

	if !questCompleted {
		t.Error("Quest completion event should have been fired")
	}

	if quest.CompletedAt == nil {
		t.Error("Quest CompletedAt should be set after completion")
	}

	// Verify character received quest reward
	if char.QuestsCompleted != 1 {
		t.Errorf("Character QuestsCompleted = %d, want 1", char.QuestsCompleted)
	}

	// Verify XP reward was added
	if char.XP < quest.XPReward {
		t.Errorf("Character XP = %d, should be at least %d (quest reward)",
			char.XP, quest.XPReward)
	}

	// Verify quest state persisted
	loadedQuests, err := storage.LoadQuests()
	if err != nil {
		t.Fatalf("Failed to load quests: %v", err)
	}

	if len(loadedQuests) != 1 {
		t.Errorf("Loaded %d quests, want 1", len(loadedQuests))
	}

	if loadedQuests[0].Status != game.QuestCompleted {
		t.Errorf("Loaded quest status = %s, want %s",
			loadedQuests[0].Status, game.QuestCompleted)
	}
}

// TestMVP_QuestProgress tests quest progress tracking for lines-based quests.
func TestMVP_QuestProgress(t *testing.T) {
	// Create a lines quest
	quest := game.NewQuest(
		"Code Sprint",
		"Write 500 lines of code",
		game.QuestTypeLines,
		500, // Target: 500 lines
		100, // Reward: 100 XP
		1,   // Required level: 1
	)

	// Start the quest
	if err := quest.Start("", ""); err != nil {
		t.Fatalf("Failed to start quest: %v", err)
	}

	// Update progress incrementally
	tests := []struct {
		lines            int
		expectedCurrent  int
		expectedProgress float64
		shouldComplete   bool
	}{
		{100, 100, 0.2, false},
		{150, 250, 0.5, false},
		{200, 450, 0.9, false},
		{50, 500, 1.0, true},
		{50, 500, 1.0, true}, // Extra lines don't overshoot
	}

	for i, tt := range tests {
		quest.UpdateProgress(tt.lines)

		if quest.Current != tt.expectedCurrent {
			t.Errorf("Step %d: quest current = %d, want %d",
				i, quest.Current, tt.expectedCurrent)
		}

		// Allow small floating point tolerance
		if quest.Progress < tt.expectedProgress-0.01 || quest.Progress > tt.expectedProgress+0.01 {
			t.Errorf("Step %d: quest progress = %.2f, want %.2f",
				i, quest.Progress, tt.expectedProgress)
		}

		isComplete := quest.CheckCompletion()
		if isComplete != tt.shouldComplete {
			t.Errorf("Step %d: quest completion = %v, want %v",
				i, isComplete, tt.shouldComplete)
		}
	}

	// Complete the quest
	if err := quest.Complete(); err != nil {
		t.Errorf("Failed to complete quest: %v", err)
	}

	// Verify progress is exactly 1.0
	if quest.Progress != 1.0 {
		t.Errorf("Completed quest progress = %.2f, want 1.0", quest.Progress)
	}
}

// TestMVP_LevelUpFlow tests the level-up progression system.
func TestMVP_LevelUpFlow(t *testing.T) {
	char := game.NewCharacter("TestHero")
	eventBus := game.NewEventBus()

	// Track level-up events
	levelUps := 0
	eventBus.Subscribe(game.EventLevelUp, func(e game.Event) {
		levelUps++
		oldLevel, _ := e.Data["old_level"].(int)
		newLevel, _ := e.Data["new_level"].(int)
		t.Logf("Level up: %d → %d", oldLevel, newLevel)
	})

	// Simulate gaining XP
	initialStats := map[string]int{
		"CodePower": char.CodePower,
		"Wisdom":    char.Wisdom,
		"Agility":   char.Agility,
	}

	// Add enough XP to level up from 1 to 2
	xpForLevelTwo := game.CalculateXPForLevel(1) // 110 XP
	oldLevel := char.Level
	leveledUp := char.AddXP(xpForLevelTwo)

	if !leveledUp {
		t.Error("Should have leveled up")
	}

	if char.Level != 2 {
		t.Errorf("Character level = %d, want 2", char.Level)
	}

	// Publish level-up event
	if leveledUp {
		eventBus.Publish(game.NewLevelUpEvent(char.ID, oldLevel, char.Level))
	}

	// Verify stats increased
	if char.CodePower != initialStats["CodePower"]+1 {
		t.Errorf("CodePower = %d, want %d",
			char.CodePower, initialStats["CodePower"]+1)
	}

	// Verify level-up event was fired
	if levelUps != 1 {
		t.Errorf("Level-up events = %d, want 1", levelUps)
	}

	// Test multi-level-up
	xpForMultipleLevels := game.CalculateXPForLevel(2) + game.CalculateXPForLevel(3) + 50
	oldLevel = char.Level
	leveledUp = char.AddXP(xpForMultipleLevels)

	if !leveledUp {
		t.Error("Should have leveled up multiple times")
	}

	if char.Level < 4 {
		t.Errorf("Character level = %d, want at least 4 after multi-level XP", char.Level)
	}
}

// TestMVP_SessionTracking tests session time tracking functionality.
func TestMVP_SessionTracking(t *testing.T) {
	char := game.NewCharacter("TestHero")

	// Simulate a coding session
	sessionStart := time.Now()

	// Simulate some work over 2 seconds
	time.Sleep(2 * time.Second)

	sessionEnd := time.Now()
	sessionDuration := sessionEnd.Sub(sessionStart)

	// Update character's session time
	char.TodaySessionTime += sessionDuration

	// Verify session time is tracked
	if char.TodaySessionTime < 2*time.Second {
		t.Errorf("Session time = %v, want at least 2s", char.TodaySessionTime)
	}

	// Test daily stats reset
	char.TodayCommits = 5
	char.TodayLinesAdded = 100

	char.ResetDailyStats()

	// Verify daily stats were reset
	if char.TodayCommits != 0 {
		t.Errorf("After reset, TodayCommits = %d, want 0", char.TodayCommits)
	}

	if char.TodaySessionTime != 0 {
		t.Errorf("After reset, TodaySessionTime = %v, want 0", char.TodaySessionTime)
	}
}

// TestMVP_StreakTracking tests daily streak tracking.
func TestMVP_StreakTracking(t *testing.T) {
	char := game.NewCharacter("TestHero")

	// First activity - start streak
	char.UpdateStreak()

	if char.CurrentStreak != 1 {
		t.Errorf("After first activity, streak = %d, want 1", char.CurrentStreak)
	}

	if char.LongestStreak != 1 {
		t.Errorf("After first activity, longest streak = %d, want 1", char.LongestStreak)
	}

	// Same day activity - no change
	char.UpdateStreak()

	if char.CurrentStreak != 1 {
		t.Errorf("Same day activity should not change streak, got %d", char.CurrentStreak)
	}

	// Simulate next day activity
	char.LastActiveDate = time.Now().AddDate(0, 0, -1) // Yesterday
	char.UpdateStreak()

	if char.CurrentStreak != 2 {
		t.Errorf("Consecutive day activity, streak = %d, want 2", char.CurrentStreak)
	}

	if char.LongestStreak != 2 {
		t.Errorf("After 2 days, longest streak = %d, want 2", char.LongestStreak)
	}

	// Simulate missed days - reset streak
	char.LastActiveDate = time.Now().AddDate(0, 0, -5) // 5 days ago
	char.UpdateStreak()

	if char.CurrentStreak != 1 {
		t.Errorf("After missed days, streak = %d, want 1 (reset)", char.CurrentStreak)
	}

	if char.LongestStreak != 2 {
		t.Errorf("Longest streak should remain = %d, got %d", 2, char.LongestStreak)
	}
}

// TestMVP_PersistenceFlow tests data persistence between sessions.
func TestMVP_PersistenceFlow(t *testing.T) {
	storage := &mockStorage{}

	// Session 1: Create character and quest
	char := game.NewCharacter("TestHero")
	char.Level = 5
	char.XP = 250
	char.TotalCommits = 42

	quest := game.NewQuest(
		"Test Quest",
		"A test quest",
		game.QuestTypeCommit,
		10,
		100,
		1,
	)
	quest.Start("", "")
	quest.UpdateProgress(5)

	// Save data
	if err := storage.SaveCharacter(char); err != nil {
		t.Fatalf("Failed to save character: %v", err)
	}

	if err := storage.SaveQuests([]*game.Quest{quest}); err != nil {
		t.Fatalf("Failed to save quests: %v", err)
	}

	// Session 2: Load data (simulating new session)
	loadedChar, err := storage.LoadCharacter()
	if err != nil {
		t.Fatalf("Failed to load character: %v", err)
	}

	// Verify character data persisted
	if loadedChar.Name != "TestHero" {
		t.Errorf("Loaded character name = %s, want TestHero", loadedChar.Name)
	}

	if loadedChar.Level != 5 {
		t.Errorf("Loaded character level = %d, want 5", loadedChar.Level)
	}

	if loadedChar.XP != 250 {
		t.Errorf("Loaded character XP = %d, want 250", loadedChar.XP)
	}

	if loadedChar.TotalCommits != 42 {
		t.Errorf("Loaded character TotalCommits = %d, want 42", loadedChar.TotalCommits)
	}

	// Load quests
	loadedQuests, err := storage.LoadQuests()
	if err != nil {
		t.Fatalf("Failed to load quests: %v", err)
	}

	if len(loadedQuests) != 1 {
		t.Fatalf("Loaded %d quests, want 1", len(loadedQuests))
	}

	loadedQuest := loadedQuests[0]
	if loadedQuest.Title != "Test Quest" {
		t.Errorf("Loaded quest title = %s, want Test Quest", loadedQuest.Title)
	}

	if loadedQuest.Current != 5 {
		t.Errorf("Loaded quest progress = %d, want 5", loadedQuest.Current)
	}

	if loadedQuest.Status != game.QuestActive {
		t.Errorf("Loaded quest status = %s, want %s", loadedQuest.Status, game.QuestActive)
	}
}

// TestMVP_EventBusConcurrency tests thread-safe event bus operations.
func TestMVP_EventBusConcurrency(t *testing.T) {
	eventBus := game.NewEventBus()

	// Track event counts
	commitCount := 0
	levelUpCount := 0

	// Subscribe handlers
	eventBus.Subscribe(game.EventCommit, func(e game.Event) {
		commitCount++
		time.Sleep(10 * time.Millisecond) // Simulate processing
	})

	eventBus.Subscribe(game.EventLevelUp, func(e game.Event) {
		levelUpCount++
		time.Sleep(10 * time.Millisecond)
	})

	// Publish events concurrently
	done := make(chan bool)
	numEvents := 10

	go func() {
		for i := 0; i < numEvents; i++ {
			eventBus.Publish(game.NewCommitEvent("sha"+string(rune(i)), "msg", 1, 10, 0))
		}
		done <- true
	}()

	go func() {
		for i := 0; i < numEvents; i++ {
			eventBus.Publish(game.NewLevelUpEvent("char-id", i, i+1))
		}
		done <- true
	}()

	// Wait for all events to be published
	<-done
	<-done

	// Give handlers time to process
	time.Sleep(200 * time.Millisecond)

	// Verify all events were handled
	if commitCount != numEvents {
		t.Errorf("Commit event count = %d, want %d", commitCount, numEvents)
	}

	if levelUpCount != numEvents {
		t.Errorf("Level-up event count = %d, want %d", levelUpCount, numEvents)
	}
}

// TestMVP_QuestPrerequisites tests quest availability based on level requirements.
func TestMVP_QuestPrerequisites(t *testing.T) {
	// Create characters at different levels
	level1Char := game.NewCharacter("Newbie")
	level5Char := game.NewCharacter("Intermediate")
	level5Char.Level = 5

	// Use character to avoid unused variable warning
	_ = level1Char
	_ = level5Char

	// Create quests with different requirements
	starterQuest := game.NewQuest(
		"Starter Quest",
		"For beginners",
		game.QuestTypeCommit,
		5,
		50,
		1, // Level 1 required
	)

	advancedQuest := game.NewQuest(
		"Advanced Quest",
		"For experienced developers",
		game.QuestTypeCommit,
		10,
		200,
		5, // Level 5 required
	)

	// Test availability
	if !starterQuest.IsAvailable(level1Char) {
		t.Error("Starter quest should be available for level 1 character")
	}

	if advancedQuest.IsAvailable(level1Char) {
		t.Error("Advanced quest should NOT be available for level 1 character")
	}

	if !starterQuest.IsAvailable(level5Char) {
		t.Error("Starter quest should be available for level 5 character")
	}

	if !advancedQuest.IsAvailable(level5Char) {
		t.Error("Advanced quest should be available for level 5 character")
	}

	// Test that started quests are no longer available
	starterQuest.Start("", "")
	if starterQuest.IsAvailable(level1Char) {
		t.Error("Started quest should not be available")
	}
}

// TestMVP_MultipleQuestsParallel tests tracking multiple active quests simultaneously.
func TestMVP_MultipleQuestsParallel(t *testing.T) {
	// Character creation (not actively used in this test but represents the player)
	_ = game.NewCharacter("MultiTasker")
	eventBus := game.NewEventBus()

	// Create multiple quests
	commitQuest := game.NewQuest("Commit Quest", "Make commits", game.QuestTypeCommit, 5, 50, 1)
	linesQuest := game.NewQuest("Lines Quest", "Write code", game.QuestTypeLines, 100, 75, 1)

	// Start both quests
	commitQuest.Start("", "")
	linesQuest.Start("", "")

	quests := []*game.Quest{commitQuest, linesQuest}

	// Subscribe to commit events and update all relevant quests
	eventBus.Subscribe(game.EventCommit, func(e game.Event) {
		linesAdded, _ := e.Data["lines_added"].(int)

		for _, quest := range quests {
			if quest.Status != game.QuestActive {
				continue
			}

			switch quest.Type {
			case game.QuestTypeCommit:
				quest.UpdateProgress(1)
			case game.QuestTypeLines:
				quest.UpdateProgress(linesAdded)
			}

			// Complete if ready
			if quest.CheckCompletion() {
				quest.Complete()
			}
		}
	})

	// Simulate 5 commits with 25 lines each (total 125 lines)
	for i := 0; i < 5; i++ {
		eventBus.Publish(game.NewCommitEvent(
			"sha"+string(rune(i)),
			"Commit "+string(rune(i)),
			1,
			25,
			0,
		))
	}

	// Verify commit quest completed
	if commitQuest.Status != game.QuestCompleted {
		t.Errorf("Commit quest status = %s, want %s", commitQuest.Status, game.QuestCompleted)
	}

	// Verify lines quest completed (5 commits × 25 lines = 125 lines > 100 target)
	if linesQuest.Status != game.QuestCompleted {
		t.Errorf("Lines quest status = %s, want %s", linesQuest.Status, game.QuestCompleted)
	}

	// Verify progress values
	if commitQuest.Current != 5 {
		t.Errorf("Commit quest progress = %d, want 5", commitQuest.Current)
	}

	if linesQuest.Current < 100 {
		t.Errorf("Lines quest progress = %d, want at least 100", linesQuest.Current)
	}
}
