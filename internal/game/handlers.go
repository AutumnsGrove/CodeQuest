// Package game contains the event handler system for CodeQuest.
// This file implements handlers that process game events and update game state.
package game

import (
	"fmt"
	"log"
	"sync"

	"github.com/AutumnsGrove/codequest/internal/config"
)

// Storage defines the interface for persisting game data.
// This breaks the circular dependency between game and storage packages.
// The storage.SkateClient implements this interface.
type Storage interface {
	SaveCharacter(character *Character) error
	LoadCharacter() (*Character, error)
	SaveQuests(quests []*Quest) error
	LoadQuests() ([]*Quest, error)
}

// GameEventHandler processes game events and updates character and quest state.
// It subscribes to the EventBus and handles commit events by:
//  1. Calculating XP rewards using the game engine
//  2. Applying difficulty and wisdom multipliers
//  3. Awarding XP to the character (triggering level-ups if applicable)
//  4. Updating active quest progress
//  5. Publishing secondary events (EventLevelUp, EventQuestDone)
//  6. Persisting state changes to storage
//
// Thread Safety:
// The handler uses a mutex to protect concurrent access to character and quest data.
// Multiple events can be processed safely without race conditions.
type GameEventHandler struct {
	// Core dependencies
	character *Character     // Player character (mutable state)
	quests    []*Quest       // Active and available quests (mutable state)
	eventBus  *EventBus      // Event system for subscribing/publishing
	storage   Storage        // Persistence layer (interface for flexibility)
	config    *config.Config // Game configuration (difficulty, etc.)

	// Thread safety
	mu sync.Mutex // Protects character and quests from concurrent access

	// State management
	running bool // Indicates if handler is active
}

// NewGameEventHandler creates a new event handler with the given dependencies.
// This initializes the handler but does NOT start event processing.
// Call Start() to begin subscribing to events.
//
// Parameters:
//   - character: The player character (must not be nil)
//   - quests: Initial quest list (can be empty, but not nil)
//   - eventBus: The event bus to subscribe to (must not be nil)
//   - storage: Storage implementation for persistence (must not be nil)
//   - config: Game configuration (must not be nil)
//
// Returns:
//   - *GameEventHandler: A new handler instance ready to be started
//   - error: An error if any parameters are invalid
func NewGameEventHandler(
	character *Character,
	quests []*Quest,
	eventBus *EventBus,
	storage Storage,
	config *config.Config,
) (*GameEventHandler, error) {
	// Validate parameters
	if character == nil {
		return nil, fmt.Errorf("character cannot be nil")
	}
	if quests == nil {
		return nil, fmt.Errorf("quests cannot be nil (use empty slice)")
	}
	if eventBus == nil {
		return nil, fmt.Errorf("eventBus cannot be nil")
	}
	if storage == nil {
		return nil, fmt.Errorf("storage cannot be nil")
	}
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	return &GameEventHandler{
		character: character,
		quests:    quests,
		eventBus:  eventBus,
		storage:   storage,
		config:    config,
		running:   false,
	}, nil
}

// Start begins processing events from the EventBus.
// This subscribes the handler to EventCommit and starts processing.
// Call Stop() to unsubscribe and halt processing.
//
// Returns:
//   - error: An error if the handler is already running
func (h *GameEventHandler) Start() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.running {
		return fmt.Errorf("handler is already running")
	}

	// Subscribe to commit events
	h.eventBus.Subscribe(EventCommit, h.handleCommitEvent)

	h.running = true
	log.Println("GameEventHandler started - subscribing to commit events")

	return nil
}

// Stop halts event processing and unsubscribes from the EventBus.
// After calling Stop(), the handler will no longer process events.
// Call Start() again to resume processing.
//
// Returns:
//   - error: An error if the handler is not running
func (h *GameEventHandler) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.running {
		return fmt.Errorf("handler is not running")
	}

	// Unsubscribe from all commit event handlers
	h.eventBus.UnsubscribeAll(EventCommit)

	h.running = false
	log.Println("GameEventHandler stopped - unsubscribed from commit events")

	return nil
}

// handleCommitEvent processes a commit event and updates game state.
// This is the main event processing pipeline:
//  1. Extract commit data (lines added/removed, files changed, etc.)
//  2. Calculate base XP from lines changed
//  3. Apply difficulty multiplier
//  4. Apply wisdom bonus
//  5. Award XP to character (handle level-ups)
//  6. Update active quest progress
//  7. Persist changes to storage
//
// This method is called automatically when EventCommit is published to the EventBus.
//
// Parameters:
//   - event: The commit event to process
func (h *GameEventHandler) handleCommitEvent(event Event) {
	// Lock for thread safety
	h.mu.Lock()
	defer h.mu.Unlock()

	// Extract commit data from event
	linesAdded, linesRemoved, sha, message, err := h.extractCommitData(event)
	if err != nil {
		log.Printf("ERROR: Failed to extract commit data: %v", err)
		return
	}

	log.Printf("Processing commit %s: %s (lines: +%d -%d)",
		sha[:7], message, linesAdded, linesRemoved)

	// Calculate base XP from commit
	baseXP := CalculateCommitXP(linesAdded, linesRemoved)
	log.Printf("  Base XP: %d", baseXP)

	// Apply difficulty multiplier
	difficulty := h.config.Game.Difficulty
	xpWithDifficulty := ApplyDifficultyMultiplier(baseXP, difficulty)
	log.Printf("  After difficulty (%s): %d XP", difficulty, xpWithDifficulty)

	// Apply wisdom bonus
	wisdom := h.character.Wisdom
	finalXP := ApplyWisdomBonus(xpWithDifficulty, wisdom)
	log.Printf("  After wisdom bonus (wisdom=%d): %d XP", wisdom, finalXP)

	// Award XP to character (handles level-ups automatically)
	oldLevel := h.character.Level
	leveledUp := h.character.AddXP(finalXP)

	// Update character statistics
	h.character.TotalCommits++
	h.character.TotalLinesAdded += linesAdded
	h.character.TotalLinesRemoved += linesRemoved
	h.character.TodayCommits++
	h.character.TodayLinesAdded += linesAdded
	h.character.UpdateStreak()

	log.Printf("  Awarded %d XP to %s (Level %d, %d/%d XP)",
		finalXP, h.character.Name, h.character.Level,
		h.character.XP, h.character.XPToNextLevel)

	// Publish level-up event if leveling occurred
	if leveledUp {
		newLevel := h.character.Level
		log.Printf("  LEVEL UP! %s reached level %d!", h.character.Name, newLevel)

		levelUpEvent := NewLevelUpEvent(h.character.ID, oldLevel, newLevel)
		h.eventBus.Publish(levelUpEvent)
	}

	// Update quest progress for all active quests
	if err := h.updateQuestProgress(linesAdded, linesRemoved); err != nil {
		log.Printf("ERROR: Failed to update quest progress: %v", err)
	}

	// Persist all state changes
	if err := h.saveState(); err != nil {
		log.Printf("ERROR: Failed to save state: %v", err)
	}
}

// extractCommitData extracts the relevant data from a commit event.
// This handles type assertions and validation of the event data map.
//
// Parameters:
//   - event: The commit event to extract data from
//
// Returns:
//   - linesAdded: Number of lines added in the commit
//   - linesRemoved: Number of lines removed in the commit
//   - sha: Commit SHA hash
//   - message: Commit message
//   - error: An error if data extraction fails
func (h *GameEventHandler) extractCommitData(event Event) (int, int, string, string, error) {
	// Verify event type
	if event.Type != EventCommit {
		return 0, 0, "", "", fmt.Errorf("invalid event type: expected %s, got %s",
			EventCommit, event.Type)
	}

	// Get data map (already the correct type)
	data := event.Data

	// Extract lines_added
	linesAdded, ok := data["lines_added"].(int)
	if !ok {
		return 0, 0, "", "", fmt.Errorf("missing or invalid 'lines_added' field")
	}

	// Extract lines_removed
	linesRemoved, ok := data["lines_removed"].(int)
	if !ok {
		return 0, 0, "", "", fmt.Errorf("missing or invalid 'lines_removed' field")
	}

	// Extract sha
	sha, ok := data["sha"].(string)
	if !ok {
		return 0, 0, "", "", fmt.Errorf("missing or invalid 'sha' field")
	}

	// Extract message
	message, ok := data["message"].(string)
	if !ok {
		return 0, 0, "", "", fmt.Errorf("missing or invalid 'message' field")
	}

	// Validate non-negative values
	if linesAdded < 0 {
		linesAdded = 0
	}
	if linesRemoved < 0 {
		linesRemoved = 0
	}

	return linesAdded, linesRemoved, sha, message, nil
}

// updateQuestProgress updates progress for all active quests that track commits or lines.
// This checks each quest's type and updates progress accordingly:
//   - QuestTypeCommit: Increment progress by 1 (one commit completed)
//   - QuestTypeLines: Increment progress by total lines changed
//
// If a quest is completed during this update, a EventQuestDone event is published.
//
// Parameters:
//   - linesAdded: Lines added in the commit
//   - linesRemoved: Lines removed in the commit
//
// Returns:
//   - error: An error if quest updates fail
func (h *GameEventHandler) updateQuestProgress(linesAdded, linesRemoved int) error {
	totalLinesChanged := linesAdded + linesRemoved

	for _, quest := range h.quests {
		// Only process active quests
		if quest.Status != QuestActive {
			continue
		}

		// Update progress based on quest type
		oldProgress := quest.Current
		switch quest.Type {
		case QuestTypeCommit:
			// Commit quest: increment by 1 for each commit
			quest.UpdateProgress(1)
			if quest.Current > oldProgress {
				log.Printf("  Quest '%s': %d/%d commits (%d%%)",
					quest.Title, quest.Current, quest.Target,
					int(quest.Progress*100))
			}

		case QuestTypeLines:
			// Lines quest: increment by total lines changed
			quest.UpdateProgress(totalLinesChanged)
			if quest.Current > oldProgress {
				log.Printf("  Quest '%s': %d/%d lines (%d%%)",
					quest.Title, quest.Current, quest.Target,
					int(quest.Progress*100))
			}

		default:
			// Other quest types not handled yet (tests, PRs, refactors, etc.)
			continue
		}

		// Check if quest was completed by this update
		if quest.CheckCompletion() {
			// Mark quest as complete
			if err := quest.Complete(); err != nil {
				log.Printf("ERROR: Failed to complete quest %s: %v", quest.ID, err)
				continue
			}

			// Increment character's quests completed counter
			h.character.QuestsCompleted++

			// Award quest completion XP (with multipliers)
			questXP := quest.XPReward
			questXPWithDifficulty := ApplyDifficultyMultiplier(questXP, h.config.Game.Difficulty)
			finalQuestXP := ApplyWisdomBonus(questXPWithDifficulty, h.character.Wisdom)

			oldLevel := h.character.Level
			leveledUp := h.character.AddXP(finalQuestXP)

			log.Printf("  QUEST COMPLETE! '%s' - Awarded %d XP", quest.Title, finalQuestXP)

			// Check for level-up from quest reward
			if leveledUp {
				newLevel := h.character.Level
				log.Printf("  LEVEL UP! %s reached level %d from quest reward!",
					h.character.Name, newLevel)

				levelUpEvent := NewLevelUpEvent(h.character.ID, oldLevel, newLevel)
				h.eventBus.Publish(levelUpEvent)
			}

			// Publish quest completion event
			questDoneEvent := NewQuestDoneEvent(quest.ID, quest.Title, finalQuestXP)
			h.eventBus.Publish(questDoneEvent)
		}
	}

	return nil
}

// saveState persists the current character and quest state to storage.
// This should be called after any state-modifying operations to ensure
// progress is not lost.
//
// Returns:
//   - error: An error if persistence fails
func (h *GameEventHandler) saveState() error {
	// Save character
	if err := h.storage.SaveCharacter(h.character); err != nil {
		return fmt.Errorf("saving character: %w", err)
	}

	// Save quests
	if err := h.storage.SaveQuests(h.quests); err != nil {
		return fmt.Errorf("saving quests: %w", err)
	}

	return nil
}

// GetCharacter returns a copy of the current character state.
// This is thread-safe and returns a snapshot of the character at call time.
//
// Returns:
//   - *Character: A pointer to the character (read-only access)
func (h *GameEventHandler) GetCharacter() *Character {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.character
}

// GetQuests returns a copy of the current quest list.
// This is thread-safe and returns a snapshot of quests at call time.
//
// Returns:
//   - []*Quest: A slice of quest pointers (read-only access)
func (h *GameEventHandler) GetQuests() []*Quest {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.quests
}

// AddQuest adds a new quest to the handler's quest list.
// This is useful for dynamically adding quests during gameplay.
//
// Parameters:
//   - quest: The quest to add (must not be nil)
//
// Returns:
//   - error: An error if the quest is invalid
func (h *GameEventHandler) AddQuest(quest *Quest) error {
	if quest == nil {
		return fmt.Errorf("cannot add nil quest")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.quests = append(h.quests, quest)

	// Persist updated quest list
	if err := h.storage.SaveQuests(h.quests); err != nil {
		return fmt.Errorf("saving quests after add: %w", err)
	}

	return nil
}

// StartQuest activates a quest by ID if it's available and the character qualifies.
// This marks the quest as active and publishes a EventQuestStart event.
//
// Parameters:
//   - questID: The UUID of the quest to start
//   - repoPath: Git repository path (optional, empty string if not needed)
//   - baseSHA: Starting commit SHA (optional, empty string if not needed)
//
// Returns:
//   - error: An error if the quest cannot be started
func (h *GameEventHandler) StartQuest(questID, repoPath, baseSHA string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Find the quest by ID
	var targetQuest *Quest
	for _, quest := range h.quests {
		if quest.ID == questID {
			targetQuest = quest
			break
		}
	}

	if targetQuest == nil {
		return fmt.Errorf("quest not found: %s", questID)
	}

	// Check if quest is available
	if !targetQuest.IsAvailable(h.character) {
		return fmt.Errorf("quest '%s' is not available (status: %s, required level: %d, character level: %d)",
			targetQuest.Title, targetQuest.Status, targetQuest.RequiredLevel, h.character.Level)
	}

	// Start the quest
	if err := targetQuest.Start(repoPath, baseSHA); err != nil {
		return fmt.Errorf("starting quest: %w", err)
	}

	log.Printf("Started quest: '%s' (Type: %s, Target: %d)",
		targetQuest.Title, targetQuest.Type, targetQuest.Target)

	// Persist updated quest state
	if err := h.storage.SaveQuests(h.quests); err != nil {
		return fmt.Errorf("saving quests after start: %w", err)
	}

	// Publish quest start event
	questStartEvent := NewQuestStartEvent(targetQuest.ID, targetQuest.Title, string(targetQuest.Type))
	h.eventBus.Publish(questStartEvent)

	return nil
}
