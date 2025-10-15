package watcher

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/AutumnsGrove/codequest/internal/config"
	"github.com/AutumnsGrove/codequest/internal/game"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// mockEventBus is a thread-safe mock implementation of game.EventBus for testing.
// It records all published events for later verification.
type mockEventBus struct {
	events []game.Event
	mu     sync.Mutex
}

// newMockEventBus creates a new mock event bus.
func newMockEventBus() *mockEventBus {
	return &mockEventBus{
		events: make([]game.Event, 0),
	}
}

// PublishAsync records an event (mimics game.EventBus.PublishAsync).
func (m *mockEventBus) PublishAsync(event game.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
}

// Subscribe is a no-op for the mock (not needed for testing WatcherManager).
func (m *mockEventBus) Subscribe(eventType game.EventType, handler game.EventHandler) {
	// Not implemented for mock
}

// GetEvents returns a copy of all recorded events (thread-safe).
func (m *mockEventBus) GetEvents() []game.Event {
	m.mu.Lock()
	defer m.mu.Unlock()
	eventsCopy := make([]game.Event, len(m.events))
	copy(eventsCopy, m.events)
	return eventsCopy
}

// ClearEvents removes all recorded events.
func (m *mockEventBus) ClearEvents() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = make([]game.Event, 0)
}

// EventCount returns the number of recorded events.
func (m *mockEventBus) EventCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.events)
}

// createTestRepoWithInitialCommit creates a test repo for integration testing.
// This is similar to the helper in git_test.go but duplicated to keep tests independent.
func createTestRepoWithInitialCommit(t *testing.T) (repoPath string) {
	t.Helper()

	tmpDir := t.TempDir()

	repo, err := git.PlainInit(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to init test repo: %v", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %v", err)
	}

	// Create initial file and commit
	testFile := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test\n"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	return tmpDir
}

// makeCommitInRepo creates a new commit in an existing repository.
func makeCommitInRepo(t *testing.T, repoPath, message string, files map[string]string) {
	t.Helper()

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %v", err)
	}

	for filename, content := range files {
		filePath := filepath.Join(repoPath, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", filename, err)
		}
		if _, err := worktree.Add(filename); err != nil {
			t.Fatalf("Failed to stage file %s: %v", filename, err)
		}
	}

	_, err = worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
}

// TestNewWatcherManager tests the WatcherManager constructor.
func TestNewWatcherManager(t *testing.T) {
	t.Parallel()

	t.Run("valid inputs", func(t *testing.T) {
		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if manager == nil {
			t.Fatal("Expected manager but got nil")
		}
		if manager.eventBus != eventBus {
			t.Error("EventBus not set correctly")
		}
		if manager.config != cfg {
			t.Error("Config not set correctly")
		}
		if manager.watchers == nil {
			t.Error("Watchers map not initialized")
		}
		if manager.cancelFuncs == nil {
			t.Error("CancelFuncs map not initialized")
		}
		if manager.running {
			t.Error("Manager should not be running after construction")
		}
	})

	t.Run("nil EventBus", func(t *testing.T) {
		cfg := &config.Config{}
		manager, err := NewWatcherManager(nil, cfg)
		if err == nil {
			t.Error("Expected error for nil EventBus")
		}
		if manager != nil {
			t.Error("Expected nil manager when EventBus is nil")
		}
		if err != nil && err.Error() != "event bus cannot be nil" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	t.Run("nil Config", func(t *testing.T) {
		eventBus := game.NewEventBus()
		manager, err := NewWatcherManager(eventBus, nil)
		if err == nil {
			t.Error("Expected error for nil Config")
		}
		if manager != nil {
			t.Error("Expected nil manager when Config is nil")
		}
		if err != nil && err.Error() != "config cannot be nil" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})
}

// TestWatcherManager_Lifecycle tests Start/Stop behavior.
func TestWatcherManager_Lifecycle(t *testing.T) {
	t.Run("Start initializes configured repositories", func(t *testing.T) {
		repo1 := createTestRepoWithInitialCommit(t)
		repo2 := createTestRepoWithInitialCommit(t)

		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{repo1, repo2},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}
		defer manager.Stop()

		if !manager.IsRunning() {
			t.Error("Manager should be running after Start()")
		}

		// Verify both repositories are being watched
		repos := manager.GetWatchedRepositories()
		if len(repos) != 2 {
			t.Errorf("Expected 2 repositories, got %d", len(repos))
		}

		if !manager.IsWatching(repo1) {
			t.Errorf("Manager should be watching %s", repo1)
		}
		if !manager.IsWatching(repo2) {
			t.Errorf("Manager should be watching %s", repo2)
		}
	})

	t.Run("Stop shuts down all watchers", func(t *testing.T) {
		repo1 := createTestRepoWithInitialCommit(t)

		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{repo1},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}

		if err := manager.Stop(); err != nil {
			t.Errorf("Failed to stop manager: %v", err)
		}

		if manager.IsRunning() {
			t.Error("Manager should not be running after Stop()")
		}

		// Verify repositories are no longer watched
		repos := manager.GetWatchedRepositories()
		if len(repos) != 0 {
			t.Errorf("Expected 0 repositories after Stop(), got %d", len(repos))
		}
	})

	t.Run("Start is idempotent", func(t *testing.T) {
		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Start twice
		if err := manager.Start(ctx); err != nil {
			t.Fatalf("First Start() failed: %v", err)
		}

		if err := manager.Start(ctx); err != nil {
			t.Errorf("Second Start() should be no-op, got error: %v", err)
		}

		if !manager.IsRunning() {
			t.Error("Manager should still be running")
		}

		manager.Stop()
	})

	t.Run("Stop is idempotent", func(t *testing.T) {
		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}

		// Stop twice
		if err := manager.Stop(); err != nil {
			t.Errorf("First Stop() failed: %v", err)
		}

		if err := manager.Stop(); err != nil {
			t.Errorf("Second Stop() should be no-op, got error: %v", err)
		}

		if manager.IsRunning() {
			t.Error("Manager should not be running")
		}
	})
}

// TestWatcherManager_DynamicRepositoryManagement tests AddRepository/RemoveRepository.
func TestWatcherManager_DynamicRepositoryManagement(t *testing.T) {
	t.Run("AddRepository creates new watcher", func(t *testing.T) {
		repo := createTestRepoWithInitialCommit(t)

		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}
		defer manager.Stop()

		// Add repository
		if err := manager.AddRepository(repo); err != nil {
			t.Fatalf("Failed to add repository: %v", err)
		}

		if !manager.IsWatching(repo) {
			t.Error("Manager should be watching added repository")
		}

		repos := manager.GetWatchedRepositories()
		if len(repos) != 1 {
			t.Errorf("Expected 1 repository, got %d", len(repos))
		}
	})

	t.Run("AddRepository duplicate path is no-op", func(t *testing.T) {
		repo := createTestRepoWithInitialCommit(t)

		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{repo},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}
		defer manager.Stop()

		// Add same repository again
		if err := manager.AddRepository(repo); err != nil {
			t.Errorf("AddRepository should be no-op for duplicate, got error: %v", err)
		}

		repos := manager.GetWatchedRepositories()
		if len(repos) != 1 {
			t.Errorf("Expected 1 repository (no duplicate), got %d", len(repos))
		}
	})

	t.Run("RemoveRepository stops watcher", func(t *testing.T) {
		repo := createTestRepoWithInitialCommit(t)

		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{repo},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}
		defer manager.Stop()

		// Remove repository
		if err := manager.RemoveRepository(repo); err != nil {
			t.Fatalf("Failed to remove repository: %v", err)
		}

		if manager.IsWatching(repo) {
			t.Error("Manager should not be watching removed repository")
		}

		repos := manager.GetWatchedRepositories()
		if len(repos) != 0 {
			t.Errorf("Expected 0 repositories after removal, got %d", len(repos))
		}
	})

	t.Run("RemoveRepository non-existent path is no-op", func(t *testing.T) {
		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}
		defer manager.Stop()

		// Remove non-existent repository
		if err := manager.RemoveRepository("/does/not/exist"); err != nil {
			t.Errorf("RemoveRepository should be no-op for non-existent path, got error: %v", err)
		}
	})
}

// TestWatcherManager_EventPublishing tests event conversion and publishing.
func TestWatcherManager_EventPublishing(t *testing.T) {
	t.Run("Commits trigger EventBus.PublishAsync", func(t *testing.T) {
		repo := createTestRepoWithInitialCommit(t)

		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{repo},
			},
			Debug: config.DebugConfig{
				Enabled: false,
			},
		}

		eventBus := game.NewEventBus()
		var publishedEvents []game.Event
		var mu sync.Mutex

		// Subscribe to commit events
		eventBus.Subscribe(game.EventCommit, func(e game.Event) {
			mu.Lock()
			defer mu.Unlock()
			publishedEvents = append(publishedEvents, e)
		})

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}
		defer manager.Stop()

		// Wait for watcher to initialize
		time.Sleep(200 * time.Millisecond)

		// Make a commit
		commitMsg := "feat: Test feature"
		files := map[string]string{
			"test.go": "package main\n",
		}
		makeCommitInRepo(t, repo, commitMsg, files)

		// Wait for event to be published
		timeout := time.After(5 * time.Second)
		for {
			mu.Lock()
			count := len(publishedEvents)
			mu.Unlock()

			if count > 0 {
				break
			}

			select {
			case <-timeout:
				t.Fatal("Timeout waiting for commit event to be published")
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}

		mu.Lock()
		defer mu.Unlock()

		if len(publishedEvents) == 0 {
			t.Fatal("No events published")
		}

		event := publishedEvents[0]
		if event.Type != game.EventCommit {
			t.Errorf("Expected event type EventCommit, got %v", event.Type)
		}

		// Verify event data
		if sha, ok := event.Data["sha"].(string); !ok || sha == "" {
			t.Error("Event should contain non-empty SHA")
		}
		if msg, ok := event.Data["message"].(string); !ok || !strings.Contains(msg, commitMsg) {
			t.Errorf("Event message should contain %q", commitMsg)
		}
		if author, ok := event.Data["author"].(string); !ok || author != "Test User" {
			t.Errorf("Expected author 'Test User', got %v", author)
		}
		if repoPath, ok := event.Data["repo_path"].(string); !ok || repoPath != repo {
			t.Errorf("Expected repo_path %q, got %v", repo, repoPath)
		}
	})

	t.Run("Event data conversion is correct", func(t *testing.T) {
		// Test the convertCommitToEvent method directly
		eventBus := game.NewEventBus()
		cfg := &config.Config{}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		commitEvent := CommitEvent{
			RepoPath:  "/test/repo",
			SHA:       "abc123def456",
			Timestamp: time.Now(),
			Author:    "Test Author",
			Email:     "test@example.com",
			Message:   "Test commit message",
			FilesChanged: []FileChange{
				{Path: "file1.go", Added: 10, Removed: 2},
				{Path: "file2.go", Added: 5, Removed: 0},
			},
			TotalAdded:   15,
			TotalRemoved: 2,
			TotalFiles:   2,
		}

		gameEvent := manager.convertCommitToEvent(commitEvent)

		if gameEvent.Type != game.EventCommit {
			t.Errorf("Expected EventCommit, got %v", gameEvent.Type)
		}

		// Verify all data fields
		if gameEvent.Data["sha"] != commitEvent.SHA {
			t.Errorf("SHA mismatch: expected %q, got %v", commitEvent.SHA, gameEvent.Data["sha"])
		}
		if gameEvent.Data["message"] != commitEvent.Message {
			t.Errorf("Message mismatch")
		}
		if gameEvent.Data["author"] != commitEvent.Author {
			t.Errorf("Author mismatch")
		}
		if gameEvent.Data["email"] != commitEvent.Email {
			t.Errorf("Email mismatch")
		}
		if gameEvent.Data["files_changed"] != commitEvent.TotalFiles {
			t.Errorf("Files changed mismatch")
		}
		if gameEvent.Data["lines_added"] != commitEvent.TotalAdded {
			t.Errorf("Lines added mismatch")
		}
		if gameEvent.Data["lines_removed"] != commitEvent.TotalRemoved {
			t.Errorf("Lines removed mismatch")
		}
		if gameEvent.Data["repo_path"] != commitEvent.RepoPath {
			t.Errorf("Repo path mismatch")
		}

		// Verify file details
		fileDetails, ok := gameEvent.Data["file_details"].([]FileChange)
		if !ok {
			t.Fatal("file_details should be []FileChange")
		}
		if len(fileDetails) != 2 {
			t.Errorf("Expected 2 file details, got %d", len(fileDetails))
		}
	})
}

// TestWatcherManager_StatusQueries tests GetWatchedRepositories, IsRunning, IsWatching.
func TestWatcherManager_StatusQueries(t *testing.T) {
	t.Run("GetWatchedRepositories returns correct paths", func(t *testing.T) {
		repo1 := createTestRepoWithInitialCommit(t)
		repo2 := createTestRepoWithInitialCommit(t)

		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{repo1, repo2},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}
		defer manager.Stop()

		repos := manager.GetWatchedRepositories()
		if len(repos) != 2 {
			t.Errorf("Expected 2 repositories, got %d", len(repos))
		}

		// Check that both repos are in the list
		repoMap := make(map[string]bool)
		for _, r := range repos {
			repoMap[r] = true
		}

		if !repoMap[repo1] {
			t.Errorf("Expected %s in watched repositories", repo1)
		}
		if !repoMap[repo2] {
			t.Errorf("Expected %s in watched repositories", repo2)
		}
	})

	t.Run("IsRunning reflects manager state", func(t *testing.T) {
		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		if manager.IsRunning() {
			t.Error("Manager should not be running before Start()")
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}

		if !manager.IsRunning() {
			t.Error("Manager should be running after Start()")
		}

		if err := manager.Stop(); err != nil {
			t.Fatalf("Failed to stop manager: %v", err)
		}

		if manager.IsRunning() {
			t.Error("Manager should not be running after Stop()")
		}
	})

	t.Run("IsWatching detects active watchers", func(t *testing.T) {
		repo := createTestRepoWithInitialCommit(t)

		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}
		defer manager.Stop()

		if manager.IsWatching(repo) {
			t.Error("Should not be watching repository before AddRepository()")
		}

		if err := manager.AddRepository(repo); err != nil {
			t.Fatalf("Failed to add repository: %v", err)
		}

		if !manager.IsWatching(repo) {
			t.Error("Should be watching repository after AddRepository()")
		}

		if err := manager.RemoveRepository(repo); err != nil {
			t.Fatalf("Failed to remove repository: %v", err)
		}

		if manager.IsWatching(repo) {
			t.Error("Should not be watching repository after RemoveRepository()")
		}
	})
}

// TestWatcherManager_ThreadSafety tests concurrent operations on WatcherManager.
func TestWatcherManager_ThreadSafety(t *testing.T) {
	t.Run("Concurrent AddRepository/RemoveRepository", func(t *testing.T) {
		// Create multiple test repositories
		repos := make([]string, 5)
		for i := 0; i < 5; i++ {
			repos[i] = createTestRepoWithInitialCommit(t)
		}

		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}
		defer manager.Stop()

		var wg sync.WaitGroup

		// Concurrently add and remove repositories
		for i := 0; i < 10; i++ {
			wg.Add(2)

			go func(idx int) {
				defer wg.Done()
				repo := repos[idx%len(repos)]
				_ = manager.AddRepository(repo)
			}(i)

			go func(idx int) {
				defer wg.Done()
				repo := repos[idx%len(repos)]
				_ = manager.RemoveRepository(repo)
			}(i)
		}

		wg.Wait()
		// If we get here without panic, thread safety is maintained
	})

	t.Run("Concurrent status queries", func(t *testing.T) {
		repo := createTestRepoWithInitialCommit(t)

		eventBus := game.NewEventBus()
		cfg := &config.Config{
			Git: config.GitConfig{
				WatchPaths: []string{repo},
			},
		}

		manager, err := NewWatcherManager(eventBus, cfg)
		if err != nil {
			t.Fatalf("Failed to create manager: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := manager.Start(ctx); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}
		defer manager.Stop()

		var wg sync.WaitGroup

		// Concurrently query status
		for i := 0; i < 100; i++ {
			wg.Add(3)

			go func() {
				defer wg.Done()
				_ = manager.GetWatchedRepositories()
			}()

			go func() {
				defer wg.Done()
				_ = manager.IsRunning()
			}()

			go func() {
				defer wg.Done()
				_ = manager.IsWatching(repo)
			}()
		}

		wg.Wait()
		// If we get here without panic, thread safety is maintained
	})
}
