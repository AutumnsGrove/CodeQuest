// Package integration provides helpers for integration testing CodeQuest MVP flows.
// This package contains mock implementations and test utilities for end-to-end testing.
package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/AutumnsGrove/codequest/internal/game"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// mockStorage is a simple in-memory storage implementation for testing.
// This avoids dependencies on Skate CLI during integration tests.
type mockStorage struct {
	character           *game.Character
	quests              []*game.Quest
	SaveCharacterCalled bool
	SaveQuestsCalled    bool
	LoadCharacterCalled bool
	LoadQuestsCalled    bool
}

// SaveCharacter stores the character in memory.
func (m *mockStorage) SaveCharacter(char *game.Character) error {
	m.SaveCharacterCalled = true
	if char == nil {
		return fmt.Errorf("cannot save nil character")
	}
	// Deep copy to avoid test interference
	charCopy := *char
	m.character = &charCopy
	return nil
}

// LoadCharacter retrieves the stored character.
func (m *mockStorage) LoadCharacter() (*game.Character, error) {
	m.LoadCharacterCalled = true
	if m.character == nil {
		return nil, fmt.Errorf("no character saved")
	}
	// Deep copy to avoid test interference
	charCopy := *m.character
	return &charCopy, nil
}

// SaveQuests stores the quest list in memory.
func (m *mockStorage) SaveQuests(quests []*game.Quest) error {
	m.SaveQuestsCalled = true
	if quests == nil {
		return fmt.Errorf("cannot save nil quests list (use empty slice instead)")
	}
	// Deep copy to avoid test interference
	m.quests = make([]*game.Quest, len(quests))
	for i, q := range quests {
		questCopy := *q
		m.quests[i] = &questCopy
	}
	return nil
}

// LoadQuests retrieves the stored quests.
func (m *mockStorage) LoadQuests() ([]*game.Quest, error) {
	m.LoadQuestsCalled = true
	if m.quests == nil {
		return []*game.Quest{}, nil
	}
	// Deep copy to avoid test interference
	quests := make([]*game.Quest, len(m.quests))
	for i, q := range m.quests {
		questCopy := *q
		quests[i] = &questCopy
	}
	return quests, nil
}

// DeleteCharacter removes the character from storage.
func (m *mockStorage) DeleteCharacter() error {
	m.character = nil
	return nil
}

// DeleteQuests removes all quests from storage.
func (m *mockStorage) DeleteQuests() error {
	m.quests = nil
	return nil
}

// CharacterExists checks if a character is stored.
func (m *mockStorage) CharacterExists() bool {
	return m.character != nil
}

// Reset clears all stored data and call flags.
func (m *mockStorage) Reset() {
	m.character = nil
	m.quests = nil
	m.SaveCharacterCalled = false
	m.SaveQuestsCalled = false
	m.LoadCharacterCalled = false
	m.LoadQuestsCalled = false
}

// createTestRepo creates a temporary git repository for testing.
// The repository is created in a temp directory and cleaned up by t.Cleanup().
//
// Returns:
//   - string: Path to the repository
//   - *git.Repository: The go-git repository object
func createTestRepo(t *testing.T) (string, *git.Repository) {
	t.Helper()

	// Create temp directory
	tmpDir := t.TempDir()

	// Initialize git repository
	repo, err := git.PlainInit(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Create initial commit to establish repo
	makeCommit(t, tmpDir, "README.md", "# Test Repository\n", "Initial commit")

	return tmpDir, repo
}

// makeCommit creates a commit in the test repository.
// This simulates a real git commit with the specified file content.
//
// Parameters:
//   - t: Test instance
//   - repoPath: Path to the repository
//   - filename: File to create/modify
//   - content: Content to write to the file
//   - message: Commit message
//
// Returns:
//   - string: The commit SHA hash
func makeCommit(t *testing.T, repoPath, filename, content, message string) string {
	t.Helper()

	// Write file content
	filePath := filepath.Join(repoPath, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file %s: %v", filename, err)
	}

	// Open repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}

	// Get worktree
	w, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %v", err)
	}

	// Stage the file
	if _, err := w.Add(filename); err != nil {
		t.Fatalf("Failed to add file %s: %v", filename, err)
	}

	// Create commit
	commit, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@codequest.dev",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	return commit.String()
}

// makeCommitWithLines creates a commit with a specified number of lines.
// This is useful for testing quest progress based on lines of code.
//
// Parameters:
//   - t: Test instance
//   - repoPath: Path to the repository
//   - filename: File to create/modify
//   - lines: Number of lines to write
//   - message: Commit message
//
// Returns:
//   - string: The commit SHA hash
func makeCommitWithLines(t *testing.T, repoPath, filename string, lines int, message string) string {
	t.Helper()

	// Generate content with specified number of lines
	var content strings.Builder
	for i := 1; i <= lines; i++ {
		content.WriteString(fmt.Sprintf("line %d\n", i))
	}

	return makeCommit(t, repoPath, filename, content.String(), message)
}

// getCommitStats extracts commit statistics (files changed, lines added/removed).
// This simulates what the git watcher would extract from a real commit.
//
// Parameters:
//   - t: Test instance
//   - repoPath: Path to the repository
//   - commitSHA: Commit hash to analyze
//
// Returns:
//   - filesChanged: Number of files changed
//   - linesAdded: Lines of code added
//   - linesRemoved: Lines of code removed
func getCommitStats(t *testing.T, repoPath, commitSHA string) (filesChanged, linesAdded, linesRemoved int) {
	t.Helper()

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}

	// Get commit object
	commit, err := repo.CommitObject(plumbing.NewHash(commitSHA))
	if err != nil {
		t.Fatalf("Failed to get commit object: %v", err)
	}

	// Get parent commit
	parent, err := commit.Parent(0)
	if err != nil {
		// No parent (initial commit)
		// For initial commit, count all content as added
		tree, _ := commit.Tree()
		filesChanged = 0
		linesAdded = 0
		tree.Files().ForEach(func(f *object.File) error {
			filesChanged++
			content, _ := f.Contents()
			linesAdded += strings.Count(content, "\n")
			return nil
		})
		return filesChanged, linesAdded, 0
	}

	// Get diff between parent and current commit
	parentTree, err := parent.Tree()
	if err != nil {
		t.Fatalf("Failed to get parent tree: %v", err)
	}

	currentTree, err := commit.Tree()
	if err != nil {
		t.Fatalf("Failed to get current tree: %v", err)
	}

	// Get changes
	changes, err := parentTree.Diff(currentTree)
	if err != nil {
		t.Fatalf("Failed to get diff: %v", err)
	}

	// Count stats using git diff
	filesChanged = len(changes)

	// Use git CLI to get accurate line counts
	cmd := exec.Command("git", "show", "--numstat", "--format=", commitSHA)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		// Fallback to approximate counts
		return filesChanged, 10, 0
	}

	// Parse numstat output: "added\tremoved\tfilename"
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			var added, removed int
			fmt.Sscanf(parts[0], "%d", &added)
			fmt.Sscanf(parts[1], "%d", &removed)
			linesAdded += added
			linesRemoved += removed
		}
	}

	return filesChanged, linesAdded, linesRemoved
}

// waitForCondition polls a condition function until it returns true or times out.
// This is useful for testing asynchronous behavior like event processing.
//
// Parameters:
//   - timeout: Maximum time to wait
//   - interval: How often to check the condition
//   - condition: Function that returns true when condition is met
//
// Returns:
//   - bool: true if condition was met, false if timeout occurred
func waitForCondition(timeout, interval time.Duration, condition func() bool) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(interval)
	}
	return false
}

// assertWithinRange checks if a value is within an expected range.
// This is useful for testing numeric values that might have slight variations.
func assertWithinRange(t *testing.T, name string, got, min, max int) {
	t.Helper()
	if got < min || got > max {
		t.Errorf("%s = %d, want between %d and %d", name, got, min, max)
	}
}

// assertDurationWithin checks if a duration is within an expected range.
func assertDurationWithin(t *testing.T, name string, got, expected, tolerance time.Duration) {
	t.Helper()
	diff := got - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		t.Errorf("%s = %v, want %v (±%v)", name, got, expected, tolerance)
	}
}
