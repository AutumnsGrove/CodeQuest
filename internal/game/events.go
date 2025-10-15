// Package game contains the event system for CodeQuest.
// This implements a publish/subscribe pattern that decouples game components,
// allowing different parts of the system to react to events without direct dependencies.
package game

import (
	"sync"
	"time"
)

// EventType defines different event types that can occur in the game.
// Each event type represents a significant game state change that other
// components might want to react to.
type EventType string

const (
	// EventCommit is fired when the player makes a git commit.
	// Data fields:
	//   - "sha": string - Commit SHA hash
	//   - "message": string - Commit message
	//   - "files_changed": int - Number of files changed
	//   - "lines_added": int - Lines of code added
	//   - "lines_removed": int - Lines of code removed
	EventCommit EventType = "commit"

	// EventLevelUp is fired when the character gains a level.
	// Data fields:
	//   - "old_level": int - Previous level
	//   - "new_level": int - New level
	//   - "character_id": string - Character UUID
	EventLevelUp EventType = "level_up"

	// EventQuestStart is fired when a quest becomes active.
	// Data fields:
	//   - "quest_id": string - Quest UUID
	//   - "quest_title": string - Quest display name
	//   - "quest_type": string - Quest type (commit, lines, etc.)
	EventQuestStart EventType = "quest_start"

	// EventQuestDone is fired when a quest is completed.
	// Data fields:
	//   - "quest_id": string - Quest UUID
	//   - "quest_title": string - Quest display name
	//   - "xp_reward": int - XP awarded
	EventQuestDone EventType = "quest_done"

	// EventSkillUnlock is fired when the player unlocks a new skill (post-MVP).
	// Data fields:
	//   - "skill_id": string - Skill identifier
	//   - "skill_name": string - Skill display name
	EventSkillUnlock EventType = "skill_unlock"

	// EventAchievement is fired when the player earns an achievement (post-MVP).
	// Data fields:
	//   - "achievement_id": string - Achievement identifier
	//   - "achievement_name": string - Achievement display name
	EventAchievement EventType = "achievement"
)

// Event represents something that happened in the game.
// Events carry data about the occurrence and can be subscribed to by multiple handlers.
//
// The Data map allows flexible payloads without creating separate event structs.
// Handlers should type-assert or type-check values from the Data map.
type Event struct {
	Type      EventType              `json:"type"`      // Type of event
	Timestamp time.Time              `json:"timestamp"` // When the event occurred
	Data      map[string]interface{} `json:"data"`      // Flexible payload
}

// EventHandler is a function that handles an event.
// Handlers should be fast and non-blocking. If heavy processing is needed,
// spawn a goroutine inside the handler.
//
// Errors are ignored in the current implementation, but future versions
// might collect and log them.
type EventHandler func(Event)

// EventBus manages event publishing and subscription.
// It provides a thread-safe pub/sub system for decoupling game components.
//
// Usage pattern:
//  1. Components subscribe to event types they care about
//  2. Components publish events when state changes occur
//  3. EventBus notifies all subscribers synchronously
//
// Thread Safety:
// EventBus uses sync.RWMutex for concurrent access protection.
// - Subscribe/Unsubscribe use write locks (exclusive)
// - Publish uses read locks (multiple publishers can read handlers simultaneously)
type EventBus struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus creates a new event bus with an empty handler registry.
// This is the factory function for creating EventBus instances.
//
// Returns:
//   - *EventBus: A pointer to the newly created event bus
//
// Example:
//
//	bus := NewEventBus()
//	bus.Subscribe(EventCommit, func(e Event) {
//	    fmt.Printf("Commit detected: %s\n", e.Data["message"])
//	})
//	bus.Publish(Event{
//	    Type: EventCommit,
//	    Timestamp: time.Now(),
//	    Data: map[string]interface{}{
//	        "message": "feat: Add awesome feature",
//	        "sha": "abc123",
//	    },
//	})
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Subscribe registers a handler function for a specific event type.
// The handler will be called whenever an event of that type is published.
//
// Multiple handlers can subscribe to the same event type. They will all
// be called in the order they were registered.
//
// Thread Safety:
// This method acquires a write lock, so it's safe to call from multiple
// goroutines, but it blocks other Subscribe/Unsubscribe/Publish calls.
//
// Parameters:
//   - eventType: The type of event to listen for
//   - handler: The function to call when the event occurs
//
// Example:
//
//	bus.Subscribe(EventLevelUp, func(e Event) {
//	    oldLevel := e.Data["old_level"].(int)
//	    newLevel := e.Data["new_level"].(int)
//	    fmt.Printf("Leveled up from %d to %d!\n", oldLevel, newLevel)
//	})
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// Initialize the handler slice if this is the first subscriber for this event type
	if eb.handlers[eventType] == nil {
		eb.handlers[eventType] = make([]EventHandler, 0)
	}

	// Append the handler to the list
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Publish sends an event to all registered handlers for that event type.
// Handlers are called synchronously in the order they were registered.
//
// This is a blocking call - it won't return until all handlers have finished.
// If you need non-blocking behavior, use PublishAsync() instead (when implemented).
//
// Thread Safety:
// This method acquires a read lock to get the handler list, then releases it
// before calling handlers. This allows multiple publishers to dispatch events
// concurrently while still protecting the handlers map from concurrent modification.
//
// Parameters:
//   - event: The event to publish
//
// Example:
//
//	bus.Publish(Event{
//	    Type: EventQuestDone,
//	    Timestamp: time.Now(),
//	    Data: map[string]interface{}{
//	        "quest_id": "abc-123",
//	        "quest_title": "First Steps",
//	        "xp_reward": 100,
//	    },
//	})
func (eb *EventBus) Publish(event Event) {
	// Acquire read lock to get handlers
	eb.mu.RLock()
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	// Call each handler synchronously
	// We iterate over a copy of the handler slice, so it's safe even if
	// handlers are added/removed during event dispatch
	for _, handler := range handlers {
		handler(event)
	}
}

// PublishAsync sends an event to all handlers asynchronously.
// Each handler is called in its own goroutine, allowing parallel processing.
//
// This is useful when handlers might do heavy processing and you don't want
// to block the publisher. However, be aware that:
//  1. Error handling becomes more difficult
//  2. Handler execution order is not guaranteed
//  3. You need to manage goroutine lifecycle
//
// Thread Safety:
// Like Publish(), this uses a read lock to get the handler list.
// Each handler runs in its own goroutine, so they can execute in parallel.
//
// Parameters:
//   - event: The event to publish
//
// Example:
//
//	// Heavy processing in handler won't block the publisher
//	bus.PublishAsync(Event{
//	    Type: EventCommit,
//	    Timestamp: time.Now(),
//	    Data: map[string]interface{}{
//	        "sha": "abc123",
//	    },
//	})
//	// This code continues immediately, handlers run in background
func (eb *EventBus) PublishAsync(event Event) {
	// Acquire read lock to get handlers
	eb.mu.RLock()
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	// Call each handler in its own goroutine
	for _, handler := range handlers {
		go handler(event)
	}
}

// Unsubscribe removes a specific handler from an event type.
// This is challenging because Go doesn't have function equality,
// so for MVP we'll keep this simple by providing UnsubscribeAll instead.
//
// Note: This is intentionally not implemented for MVP because:
//  1. Function comparison is not possible in Go (handlers can't be compared)
//  2. Most use cases don't need granular unsubscribe
//  3. UnsubscribeAll() or clearing the entire bus is usually sufficient
//
// Post-MVP solution: Use handler IDs or named handlers for targeted removal.

// UnsubscribeAll removes all handlers for a specific event type.
// This is useful when resetting game state or cleaning up resources.
//
// Thread Safety:
// This method acquires a write lock, blocking all other operations.
//
// Parameters:
//   - eventType: The event type to clear handlers for
//
// Example:
//
//	// Remove all level-up handlers
//	bus.UnsubscribeAll(EventLevelUp)
func (eb *EventBus) UnsubscribeAll(eventType EventType) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// Remove all handlers for this event type
	delete(eb.handlers, eventType)
}

// Clear removes all handlers for all event types.
// This resets the EventBus to its initial state.
//
// Thread Safety:
// This method acquires a write lock, blocking all other operations.
//
// Example:
//
//	// Reset the entire event system
//	bus.Clear()
func (eb *EventBus) Clear() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// Recreate the handlers map
	eb.handlers = make(map[EventType][]EventHandler)
}

// HandlerCount returns the number of handlers registered for a specific event type.
// This is useful for debugging and testing.
//
// Thread Safety:
// This method acquires a read lock.
//
// Parameters:
//   - eventType: The event type to count handlers for
//
// Returns:
//   - int: Number of handlers registered
//
// Example:
//
//	count := bus.HandlerCount(EventCommit)
//	fmt.Printf("There are %d commit handlers\n", count)
func (eb *EventBus) HandlerCount(eventType EventType) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	return len(eb.handlers[eventType])
}

// Helper functions for creating common events

// NewCommitEvent creates a commit event with the given data.
// This is a convenience function to ensure consistent event structure.
//
// Parameters:
//   - sha: Commit SHA hash
//   - message: Commit message
//   - filesChanged: Number of files changed
//   - linesAdded: Lines of code added
//   - linesRemoved: Lines of code removed
//
// Returns:
//   - Event: The constructed commit event
func NewCommitEvent(sha, message string, filesChanged, linesAdded, linesRemoved int) Event {
	return Event{
		Type:      EventCommit,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"sha":           sha,
			"message":       message,
			"files_changed": filesChanged,
			"lines_added":   linesAdded,
			"lines_removed": linesRemoved,
		},
	}
}

// NewLevelUpEvent creates a level-up event with the given data.
//
// Parameters:
//   - characterID: Character UUID
//   - oldLevel: Previous level
//   - newLevel: New level
//
// Returns:
//   - Event: The constructed level-up event
func NewLevelUpEvent(characterID string, oldLevel, newLevel int) Event {
	return Event{
		Type:      EventLevelUp,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"character_id": characterID,
			"old_level":    oldLevel,
			"new_level":    newLevel,
		},
	}
}

// NewQuestStartEvent creates a quest start event.
//
// Parameters:
//   - questID: Quest UUID
//   - questTitle: Quest display name
//   - questType: Quest type (commit, lines, etc.)
//
// Returns:
//   - Event: The constructed quest start event
func NewQuestStartEvent(questID, questTitle, questType string) Event {
	return Event{
		Type:      EventQuestStart,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"quest_id":    questID,
			"quest_title": questTitle,
			"quest_type":  questType,
		},
	}
}

// NewQuestDoneEvent creates a quest completion event.
//
// Parameters:
//   - questID: Quest UUID
//   - questTitle: Quest display name
//   - xpReward: XP awarded for completion
//
// Returns:
//   - Event: The constructed quest done event
func NewQuestDoneEvent(questID, questTitle string, xpReward int) Event {
	return Event{
		Type:      EventQuestDone,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"quest_id":    questID,
			"quest_title": questTitle,
			"xp_reward":   xpReward,
		},
	}
}
