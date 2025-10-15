// Package watcher provides Git repository monitoring and integration with the game event system.
// This file implements the WatcherManager, which bridges GitWatcher and EventBus.
package watcher

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/AutumnsGrove/codequest/internal/config"
	"github.com/AutumnsGrove/codequest/internal/game"
)

// WatcherManager manages multiple GitWatcher instances and integrates them with the game's EventBus.
// It handles the lifecycle of watchers, converts commit events to game events, and provides
// dynamic repository management.
//
// Architecture:
//   - Each watched repository gets its own GitWatcher instance
//   - Each watcher runs in its own goroutine listening for commits
//   - CommitEvent objects are converted to game.Event and published to EventBus
//   - Error handling is centralized with logging
//
// Thread Safety:
// WatcherManager uses sync.RWMutex to protect the watchers map from concurrent access.
// All public methods are thread-safe and can be called from multiple goroutines.
//
// Lifecycle:
//  1. Create with NewWatcherManager()
//  2. Start() to begin monitoring configured repositories
//  3. AddRepository()/RemoveRepository() to dynamically manage watched repos
//  4. Stop() to cleanly shutdown all watchers
//
// Example Usage:
//
//	cfg, _ := config.Load()
//	eventBus := game.NewEventBus()
//	manager, err := NewWatcherManager(eventBus, cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	if err := manager.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	defer manager.Stop()
type WatcherManager struct {
	// Core dependencies
	eventBus *game.EventBus // Event bus for publishing game events
	config   *config.Config // Application configuration

	// Watcher registry
	// Map of repository path -> GitWatcher instance
	watchers map[string]*GitWatcher

	// Context management
	// Map of repository path -> cancel function for that watcher's goroutine
	cancelFuncs map[string]context.CancelFunc

	// Thread safety
	mu sync.RWMutex // Protects watchers and cancelFuncs maps

	// State tracking
	running   bool       // Whether the manager is currently running
	runningMu sync.Mutex // Protects running flag
}

// NewWatcherManager creates a new WatcherManager instance.
// It initializes the watcher registry but does NOT start monitoring yet.
// Call Start() to begin watching repositories.
//
// The manager will be configured to watch repositories based on config.Git.WatchPaths.
// These paths support ~ expansion (e.g., "~/projects" becomes "/home/user/projects").
//
// Parameters:
//   - eventBus: The game's event bus for publishing commit events
//   - config: Application configuration containing Git settings
//
// Returns:
//   - *WatcherManager: Configured manager instance ready to start
//   - error: Configuration validation error
//
// Errors:
//   - eventBus is nil
//   - config is nil
//
// Example:
//
//	cfg := config.DefaultConfig()
//	bus := game.NewEventBus()
//	manager, err := NewWatcherManager(bus, cfg)
//	if err != nil {
//	    log.Fatalf("Failed to create watcher manager: %v", err)
//	}
func NewWatcherManager(eventBus *game.EventBus, config *config.Config) (*WatcherManager, error) {
	if eventBus == nil {
		return nil, fmt.Errorf("event bus cannot be nil")
	}
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	return &WatcherManager{
		eventBus:    eventBus,
		config:      config,
		watchers:    make(map[string]*GitWatcher),
		cancelFuncs: make(map[string]context.CancelFunc),
		running:     false,
	}, nil
}

// Start begins monitoring all configured Git repositories.
// It spawns a goroutine for each repository to listen for commit events.
//
// The method:
//  1. Expands ~ in configured watch paths
//  2. Creates GitWatcher instances for each repository
//  3. Spawns goroutines to listen for commit events
//  4. Converts and publishes events to the EventBus
//
// Start() is idempotent - calling it multiple times is safe (subsequent calls are no-ops).
// The context controls the manager's lifecycle - when cancelled, all watchers stop.
//
// Parameters:
//   - ctx: Context for cancellation and lifecycle control
//
// Returns:
//   - error: Startup error (failed to create watchers, invalid paths, etc.)
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	if err := manager.Start(ctx); err != nil {
//	    log.Fatalf("Failed to start watcher manager: %v", err)
//	}
//
//	// Manager is now running, listening for commits
//	// Commits will be published as game events to the EventBus
func (wm *WatcherManager) Start(ctx context.Context) error {
	wm.runningMu.Lock()
	if wm.running {
		wm.runningMu.Unlock()
		return nil // Already running
	}
	wm.running = true
	wm.runningMu.Unlock()

	// Expand configured watch paths (handle ~ expansion)
	watchPaths, err := config.ExpandPaths(wm.config.Git.WatchPaths)
	if err != nil {
		wm.runningMu.Lock()
		wm.running = false
		wm.runningMu.Unlock()
		return fmt.Errorf("failed to expand watch paths: %w", err)
	}

	// Start watching each configured repository
	for _, repoPath := range watchPaths {
		if err := wm.AddRepository(repoPath); err != nil {
			// Log error but continue with other repositories
			log.Printf("Warning: Failed to add repository %s: %v", repoPath, err)
		}
	}

	// Spawn a monitoring goroutine that handles global context cancellation
	go wm.monitorContext(ctx)

	return nil
}

// monitorContext watches the global context and stops all watchers when it's cancelled.
// This goroutine ensures clean shutdown when the application context is cancelled.
func (wm *WatcherManager) monitorContext(ctx context.Context) {
	<-ctx.Done()
	// Context cancelled, stop all watchers
	if err := wm.Stop(); err != nil {
		log.Printf("Error stopping watcher manager on context cancellation: %v", err)
	}
}

// Stop gracefully shuts down all Git watchers.
// It cancels all watcher goroutines and waits for them to finish.
//
// This method is idempotent - it's safe to call multiple times.
// After calling Stop(), the manager cannot be restarted (create a new instance instead).
//
// Returns:
//   - error: Shutdown error (failed to stop one or more watchers)
//
// Example:
//
//	if err := manager.Stop(); err != nil {
//	    log.Printf("Error stopping watcher manager: %v", err)
//	}
func (wm *WatcherManager) Stop() error {
	wm.runningMu.Lock()
	defer wm.runningMu.Unlock()

	if !wm.running {
		return nil // Already stopped
	}

	// Cancel all watcher goroutines
	wm.mu.Lock()
	for _, cancelFunc := range wm.cancelFuncs {
		cancelFunc()
	}

	// Stop all GitWatcher instances
	var lastErr error
	for repoPath, watcher := range wm.watchers {
		if err := watcher.Stop(); err != nil {
			log.Printf("Error stopping watcher for %s: %v", repoPath, err)
			lastErr = err
		}
	}

	// Clear the maps
	wm.watchers = make(map[string]*GitWatcher)
	wm.cancelFuncs = make(map[string]context.CancelFunc)
	wm.mu.Unlock()

	wm.running = false
	return lastErr
}

// AddRepository dynamically adds a new repository to watch.
// This allows repositories to be added at runtime without restarting the manager.
//
// If the repository is already being watched, this is a no-op.
// The repository path should be absolute. Use config.ExpandPath() if needed.
//
// Parameters:
//   - repoPath: Absolute path to the Git repository
//
// Returns:
//   - error: Failed to create or start watcher (invalid path, not a git repo, etc.)
//
// Example:
//
//	if err := manager.AddRepository("/home/user/projects/myapp"); err != nil {
//	    log.Printf("Failed to add repository: %v", err)
//	}
func (wm *WatcherManager) AddRepository(repoPath string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// Check if already watching
	if _, exists := wm.watchers[repoPath]; exists {
		return nil // Already watching, no-op
	}

	// Create new GitWatcher
	watcher, err := NewGitWatcher(repoPath)
	if err != nil {
		return fmt.Errorf("failed to create git watcher for %s: %w", repoPath, err)
	}

	// Create a context for this watcher
	ctx, cancel := context.WithCancel(context.Background())

	// Start the watcher
	if err := watcher.Start(ctx); err != nil {
		cancel()
		return fmt.Errorf("failed to start watcher for %s: %w", repoPath, err)
	}

	// Store watcher and cancel function
	wm.watchers[repoPath] = watcher
	wm.cancelFuncs[repoPath] = cancel

	// Spawn goroutine to listen for commits from this watcher
	go wm.listenForCommits(ctx, repoPath, watcher)

	// Spawn goroutine to listen for errors from this watcher
	go wm.listenForErrors(ctx, repoPath, watcher)

	log.Printf("Now watching repository: %s", repoPath)
	return nil
}

// RemoveRepository stops watching a repository and removes it from the manager.
// This allows dynamic removal of repositories at runtime.
//
// If the repository is not being watched, this is a no-op.
//
// Parameters:
//   - repoPath: Absolute path to the Git repository
//
// Returns:
//   - error: Failed to stop watcher
//
// Example:
//
//	if err := manager.RemoveRepository("/home/user/projects/old-project"); err != nil {
//	    log.Printf("Failed to remove repository: %v", err)
//	}
func (wm *WatcherManager) RemoveRepository(repoPath string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// Check if we're watching this repo
	watcher, exists := wm.watchers[repoPath]
	if !exists {
		return nil // Not watching, no-op
	}

	// Cancel the watcher's context (stops goroutines)
	if cancelFunc, ok := wm.cancelFuncs[repoPath]; ok {
		cancelFunc()
		delete(wm.cancelFuncs, repoPath)
	}

	// Stop the watcher
	if err := watcher.Stop(); err != nil {
		return fmt.Errorf("failed to stop watcher for %s: %w", repoPath, err)
	}

	// Remove from map
	delete(wm.watchers, repoPath)

	log.Printf("Stopped watching repository: %s", repoPath)
	return nil
}

// listenForCommits runs in a goroutine and listens for commit events from a GitWatcher.
// It converts CommitEvent objects to game.Event objects and publishes them to the EventBus.
//
// This goroutine exits when:
//   - The context is cancelled
//   - The watcher's commit channel is closed
func (wm *WatcherManager) listenForCommits(ctx context.Context, repoPath string, watcher *GitWatcher) {
	for {
		select {
		case <-ctx.Done():
			// Context cancelled, exit
			return

		case commitEvent, ok := <-watcher.CommitChannel():
			if !ok {
				// Commit channel closed, watcher stopped
				return
			}

			// Convert watcher.CommitEvent to game.Event
			gameEvent := wm.convertCommitToEvent(commitEvent)

			// Publish to EventBus asynchronously (non-blocking)
			wm.eventBus.PublishAsync(gameEvent)

			// Log the commit for debugging
			if wm.config.Debug.Enabled {
				log.Printf("Commit detected in %s: %s by %s (+%d -%d lines)",
					repoPath,
					commitEvent.SHA[:7],
					commitEvent.Author,
					commitEvent.TotalAdded,
					commitEvent.TotalRemoved,
				)
			}
		}
	}
}

// listenForErrors runs in a goroutine and listens for errors from a GitWatcher.
// It logs errors with context about which repository generated them.
//
// This goroutine exits when:
//   - The context is cancelled
//   - The watcher's error channel is closed
func (wm *WatcherManager) listenForErrors(ctx context.Context, repoPath string, watcher *GitWatcher) {
	for {
		select {
		case <-ctx.Done():
			// Context cancelled, exit
			return

		case err, ok := <-watcher.ErrorChannel():
			if !ok {
				// Error channel closed, watcher stopped
				return
			}

			// Log error with repository context
			log.Printf("Watcher error for %s: %v", repoPath, err)
		}
	}
}

// convertCommitToEvent transforms a watcher.CommitEvent into a game.Event.
// This mapping extracts relevant commit metadata and formats it for game logic consumption.
//
// The game.Event.Data map includes:
//   - "sha": string - Full commit hash
//   - "message": string - Commit message
//   - "author": string - Author name
//   - "email": string - Author email
//   - "timestamp": time.Time - Commit timestamp
//   - "files_changed": int - Number of files changed
//   - "lines_added": int - Total lines added
//   - "lines_removed": int - Total lines removed
//   - "repo_path": string - Absolute repository path
//   - "file_details": []FileChange - Per-file change details
//
// This data can be used by game logic handlers to:
//   - Calculate XP rewards (based on lines changed)
//   - Update quest progress (e.g., "commit 3 times")
//   - Track statistics (commits per day, languages used, etc.)
//
// Parameters:
//   - commit: The CommitEvent from GitWatcher
//
// Returns:
//   - game.Event: The converted event ready for EventBus
func (wm *WatcherManager) convertCommitToEvent(commit CommitEvent) game.Event {
	return game.Event{
		Type:      game.EventCommit,
		Timestamp: time.Now(), // Use current time for event processing
		Data: map[string]interface{}{
			// Commit identification
			"sha":       commit.SHA,
			"message":   commit.Message,
			"author":    commit.Author,
			"email":     commit.Email,
			"timestamp": commit.Timestamp, // Original commit timestamp

			// Change statistics (for XP calculation)
			"files_changed": commit.TotalFiles,
			"lines_added":   commit.TotalAdded,
			"lines_removed": commit.TotalRemoved,

			// Repository context
			"repo_path": commit.RepoPath,

			// Detailed file changes (for advanced quest tracking)
			"file_details": commit.FilesChanged,
		},
	}
}

// GetWatchedRepositories returns a list of all currently watched repository paths.
// This is useful for status displays and debugging.
//
// Thread Safety:
// This method acquires a read lock and is safe to call concurrently.
//
// Returns:
//   - []string: List of absolute paths to watched repositories
//
// Example:
//
//	repos := manager.GetWatchedRepositories()
//	fmt.Printf("Watching %d repositories:\n", len(repos))
//	for _, repo := range repos {
//	    fmt.Printf("  - %s\n", repo)
//	}
func (wm *WatcherManager) GetWatchedRepositories() []string {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	repos := make([]string, 0, len(wm.watchers))
	for repoPath := range wm.watchers {
		repos = append(repos, repoPath)
	}
	return repos
}

// IsRunning returns whether the manager is currently active.
// This is useful for status checks and ensuring proper lifecycle management.
//
// Thread Safety:
// This method is thread-safe and can be called concurrently.
//
// Returns:
//   - bool: true if manager is running, false otherwise
func (wm *WatcherManager) IsRunning() bool {
	wm.runningMu.Lock()
	defer wm.runningMu.Unlock()
	return wm.running
}

// IsWatching checks if a specific repository is currently being watched.
// This is useful before calling AddRepository() to avoid duplicates.
//
// Thread Safety:
// This method acquires a read lock and is safe to call concurrently.
//
// Parameters:
//   - repoPath: Absolute path to the repository
//
// Returns:
//   - bool: true if repository is being watched, false otherwise
//
// Example:
//
//	if !manager.IsWatching("/home/user/projects/myapp") {
//	    manager.AddRepository("/home/user/projects/myapp")
//	}
func (wm *WatcherManager) IsWatching(repoPath string) bool {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	_, exists := wm.watchers[repoPath]
	return exists
}
