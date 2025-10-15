// Package watcher provides file system monitoring capabilities for CodeQuest.
// This file implements session time tracking for coding sessions.
package watcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// SessionState represents the current state of the session timer.
type SessionState int

const (
	// SessionStopped indicates the timer is not running and has been reset.
	SessionStopped SessionState = iota

	// SessionRunning indicates the timer is actively tracking time.
	SessionRunning

	// SessionPaused indicates the timer is paused but retains accumulated time.
	SessionPaused
)

// String returns a human-readable representation of the session state.
func (s SessionState) String() string {
	switch s {
	case SessionStopped:
		return "stopped"
	case SessionRunning:
		return "running"
	case SessionPaused:
		return "paused"
	default:
		return "unknown"
	}
}

// SessionTracker tracks coding session time with start/pause/resume/stop functionality.
// It updates the character's TodaySessionTime and persists state to survive app restarts.
//
// Thread Safety:
// SessionTracker is fully thread-safe and can be called from multiple goroutines.
// All public methods use mutex protection for concurrent access.
//
// Persistence:
// The tracker saves its state to Skate (encrypted key-value store) every 60 seconds
// and on every state change (pause, stop). This allows sessions to survive app crashes
// and restarts.
//
// Usage Pattern:
//  1. Create tracker with NewSessionTracker()
//  2. Call Start() to begin or resume a session
//  3. Call Pause() to temporarily stop tracking
//  4. Call Resume() to continue from pause
//  5. Call Stop() to end the session and perform final update
//  6. Call GetElapsed() at any time to check current duration
//
// Example:
//
//	tracker := NewSessionTracker(character, storage)
//	if err := tracker.Start(); err != nil {
//	    log.Fatal(err)
//	}
//	defer tracker.Stop()
//
//	// Later...
//	elapsed := tracker.GetElapsed()
//	fmt.Printf("Session time: %v\n", elapsed)
type SessionTracker struct {
	// Start time of current session (when timer started running)
	// For paused sessions, this is adjusted backward by totalElapsed
	// to maintain continuity when resumed.
	startTime time.Time

	// When the session was paused (zero value if not paused)
	pausedAt time.Time

	// Total elapsed time in current session
	// Updated when pausing and continuously calculated when running
	totalElapsed time.Duration

	// Current state of the session (stopped, running, paused)
	state SessionState

	// Ticker for periodic character updates (every 60 seconds)
	ticker *time.Ticker

	// Character to update with session time
	character *game.Character

	// Storage interface for persistence (optional, can be nil)
	storage Storage

	// Channel to signal shutdown of update loop
	stopChan chan struct{}

	// Mutex for thread safety
	mu sync.Mutex
}

// Storage defines the interface for persisting session state.
// This abstraction allows for different storage backends (Skate, JSON files, etc.).
type Storage interface {
	// SaveCharacter persists the character data
	SaveCharacter(char *game.Character) error
}

// sessionStateData is the internal structure for persistence.
// This is what gets serialized to JSON and stored in Skate.
type sessionStateData struct {
	StartTime    time.Time     `json:"start_time"`     // When session started (adjusted for pause)
	TotalElapsed time.Duration `json:"total_elapsed"`  // Accumulated time
	State        string        `json:"state"`          // "stopped", "running", or "paused"
	SavedAt      time.Time     `json:"saved_at"`       // When state was saved
	CharacterID  string        `json:"character_id"`   // Character this session belongs to
}

// NewSessionTracker creates a new session tracker for the given character.
// It attempts to load any previously saved session state from Skate.
//
// If a previous session exists and was running/paused, it will be restored.
// If loading fails or no previous session exists, starts fresh.
//
// Parameters:
//   - char: Character to track session time for (updates TodaySessionTime)
//   - storage: Storage backend for character persistence (can be nil)
//
// Returns:
//   - *SessionTracker: Initialized tracker (state loaded if available)
//
// Example:
//
//	tracker := NewSessionTracker(character, storageBackend)
//	// Automatically tries to resume previous session if it exists
func NewSessionTracker(char *game.Character, storage Storage) *SessionTracker {
	tracker := &SessionTracker{
		character: char,
		storage:   storage,
		state:     SessionStopped,
		stopChan:  make(chan struct{}),
	}

	// Try to load previous session state
	// Errors are silently ignored - we just start fresh
	_ = tracker.loadState()

	return tracker
}

// Start begins a new session or resumes from a previously saved state.
// If called while already running, this is a no-op and returns nil.
//
// Behavior:
//   - If stopped: Starts fresh session with current time
//   - If paused: Use Resume() instead
//   - If running: No-op, returns nil
//
// The ticker starts immediately and will update character stats every 60 seconds.
//
// Returns:
//   - error: Only if unable to load previous state (non-fatal)
//
// Example:
//
//	if err := tracker.Start(); err != nil {
//	    log.Printf("Warning: Could not load previous session: %v", err)
//	}
func (s *SessionTracker) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If already running, do nothing
	if s.state == SessionRunning {
		return nil
	}

	// If paused, user should call Resume() instead
	if s.state == SessionPaused {
		return errors.New("session is paused, use Resume() to continue")
	}

	// Start fresh session
	s.startTime = time.Now()
	s.totalElapsed = 0
	s.pausedAt = time.Time{} // Zero value
	s.state = SessionRunning

	// Start update ticker (60 second intervals)
	s.ticker = time.NewTicker(60 * time.Second)
	go s.updateLoop()

	// Save initial state
	_ = s.saveState()

	return nil
}

// Pause pauses the current session, preserving accumulated time.
// The timer stops tracking until Resume() is called.
//
// Returns:
//   - error: If session is not running or if save fails
//
// Example:
//
//	if err := tracker.Pause(); err != nil {
//	    log.Printf("Failed to pause: %v", err)
//	}
func (s *SessionTracker) Pause() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != SessionRunning {
		return fmt.Errorf("cannot pause: session is %s", s.state)
	}

	// Capture current elapsed time
	s.totalElapsed = time.Since(s.startTime)
	s.pausedAt = time.Now()
	s.state = SessionPaused

	// Stop ticker (but don't close stopChan - that's for full shutdown)
	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}

	// Perform final update before pausing
	s.updateCharacterTimeUnsafe()

	// Save state
	return s.saveState()
}

// Resume continues a paused session.
// The timer picks up where it left off, maintaining accumulated time.
//
// Returns:
//   - error: If session is not paused
//
// Example:
//
//	if err := tracker.Resume(); err != nil {
//	    log.Printf("Failed to resume: %v", err)
//	}
func (s *SessionTracker) Resume() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != SessionPaused {
		return fmt.Errorf("cannot resume: session is %s", s.state)
	}

	// Resume by adjusting start time backward by accumulated time
	// This makes GetElapsed() calculation seamless
	s.startTime = time.Now().Add(-s.totalElapsed)
	s.pausedAt = time.Time{} // Zero value
	s.state = SessionRunning

	// Restart ticker
	s.ticker = time.NewTicker(60 * time.Second)
	go s.updateLoop()

	// Save state
	return s.saveState()
}

// Stop ends the current session and performs a final character update.
// This is the proper way to end a session - it saves all data and cleans up resources.
//
// After calling Stop(), the tracker is in SessionStopped state and cannot be resumed.
// To start a new session, call Start() again.
//
// Returns:
//   - error: If final save fails (session is still stopped regardless)
//
// Example:
//
//	defer tracker.Stop()
func (s *SessionTracker) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If already stopped, no-op
	if s.state == SessionStopped {
		return nil
	}

	// Stop ticker if running
	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}

	// Signal update loop to stop (if running)
	select {
	case s.stopChan <- struct{}{}:
	default:
		// Channel might be full or goroutine already stopped, that's fine
	}

	// Perform final character update
	if s.state == SessionRunning || s.state == SessionPaused {
		s.updateCharacterTimeUnsafe()
	}

	// Reset to stopped state
	s.state = SessionStopped
	s.totalElapsed = 0
	s.startTime = time.Time{}
	s.pausedAt = time.Time{}

	// Save final state (cleared state)
	return s.saveState()
}

// GetElapsed returns the current elapsed session time.
// This works correctly regardless of session state.
//
// Returns:
//   - time.Duration: Elapsed time (0 if stopped, accumulated if paused, calculated if running)
//
// Thread Safety:
// This method is thread-safe and can be called concurrently.
//
// Example:
//
//	elapsed := tracker.GetElapsed()
//	fmt.Printf("Session time: %v\n", elapsed)
func (s *SessionTracker) GetElapsed() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.state {
	case SessionStopped:
		return 0
	case SessionPaused:
		return s.totalElapsed
	case SessionRunning:
		return time.Since(s.startTime)
	default:
		return 0
	}
}

// GetState returns the current session state.
//
// Returns:
//   - SessionState: Current state (stopped, running, or paused)
//
// Thread Safety:
// This method is thread-safe and can be called concurrently.
func (s *SessionTracker) GetState() SessionState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.state
}

// updateLoop runs in a goroutine and performs periodic character updates.
// It ticks every 60 seconds while the session is running.
func (s *SessionTracker) updateLoop() {
	for {
		select {
		case <-s.ticker.C:
			// Periodic update (every 60 seconds)
			s.mu.Lock()
			if s.state == SessionRunning {
				s.updateCharacterTimeUnsafe()
				_ = s.saveState()
			}
			s.mu.Unlock()

		case <-s.stopChan:
			// Shutdown signal
			return
		}
	}
}

// updateCharacterTimeUnsafe updates the character's TodaySessionTime.
// This method MUST be called with s.mu held (Lock/Unlock).
//
// It calculates current elapsed time and updates the character's session stats.
// If a storage backend is available, it persists the character data.
func (s *SessionTracker) updateCharacterTimeUnsafe() {
	// Calculate current elapsed time based on state
	var elapsed time.Duration
	switch s.state {
	case SessionRunning:
		elapsed = time.Since(s.startTime)
	case SessionPaused:
		elapsed = s.totalElapsed
	default:
		elapsed = 0
	}

	// Update character's session time
	s.character.TodaySessionTime = elapsed

	// Persist character data if storage available
	if s.storage != nil {
		_ = s.storage.SaveCharacter(s.character)
	}
}

// saveState persists the current session state to Skate.
// This allows sessions to survive app restarts and crashes.
//
// Storage format: JSON in Skate with key "codequest_session_state"
//
// Returns:
//   - error: If Skate command fails or JSON encoding fails
func (s *SessionTracker) saveState() error {
	// Build state data structure
	stateData := sessionStateData{
		StartTime:    s.startTime,
		TotalElapsed: s.totalElapsed,
		State:        s.state.String(),
		SavedAt:      time.Now(),
		CharacterID:  s.character.ID,
	}

	// If paused, save totalElapsed; if running, calculate current elapsed
	if s.state == SessionRunning {
		stateData.TotalElapsed = time.Since(s.startTime)
	}

	// Serialize to JSON
	data, err := json.Marshal(stateData)
	if err != nil {
		return fmt.Errorf("failed to marshal session state: %w", err)
	}

	// Save to Skate
	cmd := exec.Command("skate", "set", "codequest_session_state", string(data))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to save session state to skate: %w", err)
	}

	return nil
}

// loadState attempts to restore session state from Skate.
// If successful, it resumes the previous session where it left off.
//
// Returns:
//   - error: If Skate read fails, JSON parsing fails, or no saved state exists
//
// Note: This is called automatically by NewSessionTracker().
func (s *SessionTracker) loadState() error {
	// Retrieve from Skate
	cmd := exec.Command("skate", "get", "codequest_session_state")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("no previous session state found: %w", err)
	}

	// Parse JSON
	var stateData sessionStateData
	if err := json.Unmarshal(output, &stateData); err != nil {
		return fmt.Errorf("failed to parse session state: %w", err)
	}

	// Verify this state belongs to current character
	if stateData.CharacterID != s.character.ID {
		return fmt.Errorf("session state belongs to different character")
	}

	// Restore state based on what was saved
	switch stateData.State {
	case "running":
		// Resume running session
		// Adjust start time to account for time passed since save
		timeSinceSave := time.Since(stateData.SavedAt)
		s.startTime = stateData.StartTime.Add(-timeSinceSave)
		s.totalElapsed = stateData.TotalElapsed + timeSinceSave
		s.state = SessionRunning
		s.pausedAt = time.Time{}

		// Start ticker
		s.ticker = time.NewTicker(60 * time.Second)
		go s.updateLoop()

	case "paused":
		// Resume paused session
		s.startTime = stateData.StartTime
		s.totalElapsed = stateData.TotalElapsed
		s.pausedAt = stateData.SavedAt // Use save time as pause time
		s.state = SessionPaused
		// Don't start ticker - session is paused

	case "stopped":
		// Session was stopped, start fresh
		s.state = SessionStopped
		s.totalElapsed = 0

	default:
		return fmt.Errorf("unknown session state: %s", stateData.State)
	}

	return nil
}

// ClearSavedState removes any saved session state from Skate.
// This is useful for debugging or resetting to a clean state.
//
// Returns:
//   - error: If Skate delete command fails
func (s *SessionTracker) ClearSavedState() error {
	cmd := exec.Command("skate", "delete", "codequest_session_state")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clear session state: %w", err)
	}
	return nil
}

// FormatElapsed returns a human-friendly representation of elapsed time.
// This is useful for displaying session duration to users.
//
// Format: "2h 34m" or "45m 30s" or "1h 0m" (minutes omitted if 0 and >1 hour)
//
// Returns:
//   - string: Formatted duration string
func (s *SessionTracker) FormatElapsed() string {
	elapsed := s.GetElapsed()

	hours := int(elapsed.Hours())
	minutes := int(elapsed.Minutes()) % 60
	seconds := int(elapsed.Seconds()) % 60

	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh 0m", hours)
	}

	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}

	return fmt.Sprintf("%ds", seconds)
}
