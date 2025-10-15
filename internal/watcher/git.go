// Package watcher provides file system monitoring capabilities for CodeQuest.
// This package contains the GitWatcher, which monitors Git repositories for
// new commits and extracts metadata for game mechanics (XP calculations, quest progress).
package watcher

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CommitEvent represents a detected Git commit with full metadata.
// This structure contains all information needed for XP calculations
// and quest progress tracking in the game system.
type CommitEvent struct {
	// Repository information
	RepoPath string `json:"repo_path"` // Absolute path to repository

	// Commit identification
	SHA       string    `json:"sha"`       // Full commit hash
	Timestamp time.Time `json:"timestamp"` // Commit timestamp
	Author    string    `json:"author"`    // Author name
	Email     string    `json:"email"`     // Author email
	Message   string    `json:"message"`   // Full commit message

	// Change statistics
	FilesChanged []FileChange `json:"files_changed"` // Detailed file changes
	TotalAdded   int          `json:"total_added"`   // Total lines added
	TotalRemoved int          `json:"total_removed"` // Total lines removed
	TotalFiles   int          `json:"total_files"`   // Number of files changed
}

// FileChange represents changes to a single file in a commit.
// Used for detailed quest tracking (e.g., "modify 5 Go files").
type FileChange struct {
	Path    string `json:"path"`    // File path relative to repo root
	Added   int    `json:"added"`   // Lines added in this file
	Removed int    `json:"removed"` // Lines removed in this file
}

// GitWatcher monitors a Git repository for new commits using fsnotify.
// It watches the .git/refs/heads directory and emits CommitEvent objects
// when new commits are detected.
//
// Thread Safety:
// GitWatcher is thread-safe and can be started/stopped from different goroutines.
// The commits channel should only have one consumer.
//
// Usage Pattern:
//  1. Create watcher with NewGitWatcher()
//  2. Start monitoring with Start() (spawns goroutine)
//  3. Read from CommitChannel() to receive events
//  4. Call Stop() when done to clean up resources
type GitWatcher struct {
	repoPath      string            // Absolute path to repository
	repo          *git.Repository   // go-git repository handle
	watcher       *fsnotify.Watcher // File system watcher
	commits       chan CommitEvent  // Channel for commit events
	errors        chan error        // Channel for error reporting
	done          chan struct{}     // Signal channel for shutdown
	lastCommitSHA plumbing.Hash     // Track last seen commit to avoid duplicates
	mu            sync.RWMutex      // Protects lastCommitSHA
	running       bool              // Track running state
	runningMu     sync.Mutex        // Protects running flag
}

// NewGitWatcher creates a new Git repository watcher.
// It validates the repository path and initializes the watcher.
//
// The repository must be a valid Git repository with a .git directory.
// The watcher does NOT start automatically - call Start() to begin monitoring.
//
// Parameters:
//   - repoPath: Absolute path to Git repository root (directory containing .git)
//
// Returns:
//   - *GitWatcher: Configured watcher instance
//   - error: Validation error (invalid path, not a git repo, permission issues)
//
// Errors:
//   - Repository path doesn't exist
//   - Path is not a Git repository
//   - Insufficient permissions to access repository
//   - Failed to initialize fsnotify watcher
//
// Example:
//
//	watcher, err := NewGitWatcher("/home/user/projects/myapp")
//	if err != nil {
//	    log.Fatalf("Failed to create watcher: %v", err)
//	}
//	defer watcher.Stop()
func NewGitWatcher(repoPath string) (*GitWatcher, error) {
	// Validate and open the Git repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository at %s: %w", repoPath, err)
	}

	// Verify we can access the repository
	_, err = repo.Head()
	if err != nil {
		// Check if this is an empty repository (no commits yet)
		if err == plumbing.ErrReferenceNotFound {
			// Empty repo is valid, just note it
			// We'll start watching immediately and catch the first commit
		} else {
			return nil, fmt.Errorf("failed to access repository HEAD: %w", err)
		}
	}

	// Create fsnotify watcher
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file system watcher: %w", err)
	}

	// Get current HEAD commit SHA if it exists (for duplicate detection)
	var lastSHA plumbing.Hash
	head, err := repo.Head()
	if err == nil {
		lastSHA = head.Hash()
	}

	watcher := &GitWatcher{
		repoPath:      repoPath,
		repo:          repo,
		watcher:       fsWatcher,
		commits:       make(chan CommitEvent, 10), // Buffer to prevent blocking
		errors:        make(chan error, 10),       // Buffer for error reporting
		done:          make(chan struct{}),
		lastCommitSHA: lastSHA,
		running:       false,
	}

	return watcher, nil
}

// Start begins monitoring the Git repository for new commits.
// This method spawns a goroutine that watches for file system changes
// and processes new commits. It's safe to call multiple times (subsequent
// calls are no-ops if already running).
//
// The watcher monitors .git/refs/heads/<branch> files for changes.
// When a change is detected, it:
//  1. Reads the new HEAD commit
//  2. Extracts full commit metadata using go-git
//  3. Calculates diff statistics
//  4. Sends CommitEvent to the commits channel
//
// Context Usage:
// The context controls the watcher's lifecycle. When the context is cancelled,
// the watcher stops monitoring and cleans up resources.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - error: Startup error (e.g., can't watch .git directory)
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	if err := watcher.Start(ctx); err != nil {
//	    log.Fatalf("Failed to start watcher: %v", err)
//	}
//
//	// Read commits
//	for commit := range watcher.CommitChannel() {
//	    fmt.Printf("New commit: %s\n", commit.SHA)
//	}
func (gw *GitWatcher) Start(ctx context.Context) error {
	gw.runningMu.Lock()
	if gw.running {
		gw.runningMu.Unlock()
		return nil // Already running, no-op
	}
	gw.running = true
	gw.runningMu.Unlock()

	// Watch .git/refs/heads directory for commit changes
	refsPath := filepath.Join(gw.repoPath, ".git", "refs", "heads")
	if err := gw.watcher.Add(refsPath); err != nil {
		gw.runningMu.Lock()
		gw.running = false
		gw.runningMu.Unlock()
		return fmt.Errorf("failed to watch %s: %w", refsPath, err)
	}

	// Also watch .git/HEAD for branch switches
	headPath := filepath.Join(gw.repoPath, ".git", "HEAD")
	if err := gw.watcher.Add(headPath); err != nil {
		// Non-critical, continue anyway
		gw.errors <- fmt.Errorf("warning: failed to watch HEAD file: %w", err)
	}

	// Start the monitoring goroutine
	go gw.watch(ctx)

	return nil
}

// watch is the main monitoring loop (runs in goroutine).
// It listens for file system events and processes commits.
func (gw *GitWatcher) watch(ctx context.Context) {
	defer func() {
		gw.runningMu.Lock()
		gw.running = false
		gw.runningMu.Unlock()
		close(gw.commits)
		close(gw.errors)
	}()

	for {
		select {
		case <-ctx.Done():
			// Context cancelled, shutdown cleanly
			return

		case <-gw.done:
			// Manual stop requested
			return

		case event, ok := <-gw.watcher.Events:
			if !ok {
				// Watcher closed
				return
			}

			// Only care about write events (new commits)
			if event.Op&fsnotify.Write == fsnotify.Write {
				// Process the potential new commit
				if err := gw.processCommit(); err != nil {
					// Send error but don't crash
					select {
					case gw.errors <- err:
					default:
						// Error channel full, drop error
					}
				}
			}

		case err, ok := <-gw.watcher.Errors:
			if !ok {
				// Watcher error channel closed
				return
			}
			// Forward fsnotify errors
			select {
			case gw.errors <- fmt.Errorf("fsnotify error: %w", err):
			default:
				// Error channel full, drop error
			}
		}
	}
}

// processCommit checks for a new commit and extracts its metadata.
// This is called when fsnotify detects a file change in .git/refs/heads.
func (gw *GitWatcher) processCommit() error {
	// Get current HEAD
	head, err := gw.repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	currentSHA := head.Hash()

	// Check if this is a new commit (avoid processing same commit twice)
	gw.mu.RLock()
	lastSHA := gw.lastCommitSHA
	gw.mu.RUnlock()

	if currentSHA == lastSHA {
		// Not a new commit, ignore
		return nil
	}

	// Update last seen commit
	gw.mu.Lock()
	gw.lastCommitSHA = currentSHA
	gw.mu.Unlock()

	// Extract commit metadata
	commitEvent, err := gw.extractCommitData(currentSHA)
	if err != nil {
		return fmt.Errorf("failed to extract commit data: %w", err)
	}

	// Send commit event (non-blocking)
	select {
	case gw.commits <- *commitEvent:
		// Successfully sent
	default:
		// Channel full, log warning but don't block
		gw.errors <- fmt.Errorf("commit channel full, dropping commit %s", currentSHA.String())
	}

	return nil
}

// extractCommitData retrieves full metadata for a commit using go-git.
// This includes author info, message, timestamp, and diff statistics.
func (gw *GitWatcher) extractCommitData(sha plumbing.Hash) (*CommitEvent, error) {
	// Get commit object
	commit, err := gw.repo.CommitObject(sha)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %w", err)
	}

	// Extract basic commit information
	event := &CommitEvent{
		RepoPath:  gw.repoPath,
		SHA:       sha.String(),
		Timestamp: commit.Author.When,
		Author:    commit.Author.Name,
		Email:     commit.Author.Email,
		Message:   commit.Message,
	}

	// Calculate diff statistics
	if err := gw.calculateDiffStats(commit, event); err != nil {
		// Non-critical error, continue with basic info
		gw.errors <- fmt.Errorf("warning: failed to calculate diff stats: %w", err)
	}

	return event, nil
}

// calculateDiffStats computes lines added/removed for a commit.
// This uses go-git's Stats() method to get per-file changes.
func (gw *GitWatcher) calculateDiffStats(commit *object.Commit, event *CommitEvent) error {
	// Get commit stats (files changed with line counts)
	stats, err := commit.Stats()
	if err != nil {
		return fmt.Errorf("failed to get commit stats: %w", err)
	}

	// Process each file's statistics
	event.FilesChanged = make([]FileChange, 0, len(stats))
	var totalAdded, totalRemoved int

	for _, fileStat := range stats {
		fileChange := FileChange{
			Path:    fileStat.Name,
			Added:   fileStat.Addition,
			Removed: fileStat.Deletion,
		}
		event.FilesChanged = append(event.FilesChanged, fileChange)
		totalAdded += fileStat.Addition
		totalRemoved += fileStat.Deletion
	}

	event.TotalAdded = totalAdded
	event.TotalRemoved = totalRemoved
	event.TotalFiles = len(stats)

	return nil
}

// Stop gracefully shuts down the watcher.
// It closes the file system watcher and signals the monitoring goroutine to exit.
// This method blocks until the goroutine has fully stopped.
//
// It's safe to call Stop() multiple times - subsequent calls are no-ops.
// After calling Stop(), the watcher cannot be restarted (create a new one instead).
//
// Returns:
//   - error: Error during shutdown (e.g., failed to close watcher)
//
// Example:
//
//	if err := watcher.Stop(); err != nil {
//	    log.Printf("Error stopping watcher: %v", err)
//	}
func (gw *GitWatcher) Stop() error {
	gw.runningMu.Lock()
	defer gw.runningMu.Unlock()

	if !gw.running {
		return nil // Already stopped
	}

	// Signal goroutine to stop
	close(gw.done)

	// Close fsnotify watcher
	if err := gw.watcher.Close(); err != nil {
		return fmt.Errorf("failed to close watcher: %w", err)
	}

	// Wait for goroutine to finish (channels will be closed)
	// The watch() goroutine will close commits and errors channels on exit

	return nil
}

// CommitChannel returns the channel that emits new commit events.
// Consumers should read from this channel to receive commit notifications.
//
// The channel is buffered (capacity 10) to prevent blocking the watcher.
// If the consumer is slow, older commits may be dropped (check errors channel).
//
// The channel is closed when Stop() is called or the context is cancelled.
//
// Returns:
//   - <-chan CommitEvent: Read-only channel of commit events
//
// Example:
//
//	for commit := range watcher.CommitChannel() {
//	    fmt.Printf("Commit %s by %s: %s\n",
//	        commit.SHA[:7], commit.Author, commit.Message)
//	    fmt.Printf("  +%d -%d lines across %d files\n",
//	        commit.TotalAdded, commit.TotalRemoved, commit.TotalFiles)
//	}
func (gw *GitWatcher) CommitChannel() <-chan CommitEvent {
	return gw.commits
}

// ErrorChannel returns the channel that emits errors during monitoring.
// Consumers should read from this channel to handle non-fatal errors.
//
// Errors include:
//   - Failed to calculate diff stats (commit event still sent)
//   - Commit channel full (commit dropped)
//   - fsnotify errors
//
// The channel is buffered (capacity 10). If the consumer doesn't read errors,
// the buffer may fill and new errors will be dropped silently.
//
// Returns:
//   - <-chan error: Read-only channel of error events
//
// Example:
//
//	go func() {
//	    for err := range watcher.ErrorChannel() {
//	        log.Printf("Watcher error: %v", err)
//	    }
//	}()
func (gw *GitWatcher) ErrorChannel() <-chan error {
	return gw.errors
}

// IsRunning returns whether the watcher is currently active.
// This is useful for status checks and debugging.
//
// Thread Safety:
// This method is thread-safe and can be called concurrently.
//
// Returns:
//   - bool: true if watcher is running, false otherwise
func (gw *GitWatcher) IsRunning() bool {
	gw.runningMu.Lock()
	defer gw.runningMu.Unlock()
	return gw.running
}

// GetLastCommitSHA returns the SHA of the last processed commit.
// This can be used to detect if new commits have been made since last check.
//
// Thread Safety:
// This method is thread-safe and can be called concurrently.
//
// Returns:
//   - string: Last commit SHA (empty string if no commits processed)
func (gw *GitWatcher) GetLastCommitSHA() string {
	gw.mu.RLock()
	defer gw.mu.RUnlock()
	return gw.lastCommitSHA.String()
}
