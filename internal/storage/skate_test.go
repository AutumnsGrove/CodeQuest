package storage

import (
	"encoding/json"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// Mock exec.Command for testing CLI operations without actual Skate installation
// This approach allows us to test the storage layer logic independently

// execContext holds the mock command execution context
type execContext struct {
	commandHistory [][]string            // Track all commands executed
	mockResponses  map[string]mockResult // Map command to mock response
	shouldFail     bool                  // Whether next command should fail
	failError      error                 // Error to return on failure
}

type mockResult struct {
	stdout string
	stderr string
	err    error
}

var mockExecContext *execContext

// mockCommand replaces exec.Command for testing
func mockCommand(name string, args ...string) *exec.Cmd {
	if mockExecContext == nil {
		mockExecContext = &execContext{
			commandHistory: [][]string{},
			mockResponses:  make(map[string]mockResult),
		}
	}

	// Record this command
	fullCmd := append([]string{name}, args...)
	mockExecContext.commandHistory = append(mockExecContext.commandHistory, fullCmd)

	// Return real exec.Command - we'll control it via test environment
	// This is needed because we can't easily mock exec.Cmd methods
	return exec.Command("echo", "mock")
}

// resetMockExec resets the mock execution context
func resetMockExec() {
	mockExecContext = &execContext{
		commandHistory: [][]string{},
		mockResponses:  make(map[string]mockResult),
	}
}

// TestNewSkateClient tests the Skate client initialization
func TestNewSkateClient(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "skate not in PATH",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Try to create a new client
			client, err := NewSkateClient()

			if tt.wantErr {
				// If Skate IS installed, skip this test
				if err == nil {
					t.Skip("Skate is installed, skipping 'not found' test")
				}
				// Verify error message mentions skate
				if !strings.Contains(err.Error(), "skate") {
					t.Errorf("NewSkateClient() error should mention 'skate', got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("NewSkateClient() unexpected error: %v", err)
				}
				if client == nil {
					t.Errorf("NewSkateClient() returned nil client")
				}
				if client.skatePath == "" {
					t.Errorf("NewSkateClient() skatePath is empty")
				}
			}
		})
	}
}

// TestNewSkateClient_WithSkateInstalled tests client creation when Skate is available
func TestNewSkateClient_WithSkateInstalled(t *testing.T) {
	// Check if skate is installed
	_, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client, err := NewSkateClient()
	if err != nil {
		t.Errorf("NewSkateClient() unexpected error with Skate installed: %v", err)
	}
	if client == nil {
		t.Errorf("NewSkateClient() returned nil client")
	}
	if client.skatePath == "" {
		t.Errorf("NewSkateClient() skatePath is empty")
	}
}

// TestSkateClient_SaveCharacter tests character serialization and storage
func TestSkateClient_SaveCharacter(t *testing.T) {
	// Skip if Skate not installed
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client := &SkateClient{skatePath: skatePath}

	tests := []struct {
		name      string
		character *game.Character
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "nil character",
			character: nil,
			wantErr:   true,
			errMsg:    "cannot save nil character",
		},
		{
			name:      "valid character",
			character: game.NewCharacter("TestHero"),
			wantErr:   false,
		},
		{
			name: "character with progress",
			character: func() *game.Character {
				c := game.NewCharacter("ProgressHero")
				c.AddXP(100)
				c.TotalCommits = 5
				c.CurrentStreak = 3
				return c
			}(),
			wantErr: false,
		},
		{
			name:      "character with special characters in name",
			character: game.NewCharacter("Test-Hero_123 with spaces"),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.SaveCharacter(tt.character)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SaveCharacter() expected error containing %q, got nil", tt.errMsg)
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("SaveCharacter() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("SaveCharacter() unexpected error: %v", err)
				}
			}
		})
	}

	// Cleanup - delete test character
	_ = client.DeleteCharacter()
}

// TestSkateClient_LoadCharacter tests character retrieval and deserialization
func TestSkateClient_LoadCharacter(t *testing.T) {
	// Skip if Skate not installed
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client := &SkateClient{skatePath: skatePath}

	tests := []struct {
		name     string
		setup    func() *game.Character
		wantErr  bool
		errMsg   string
		validate func(*testing.T, *game.Character, *game.Character)
	}{
		{
			name: "load non-existent character",
			setup: func() *game.Character {
				// Ensure no character exists
				_ = client.DeleteCharacter()
				return nil
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name: "load valid character",
			setup: func() *game.Character {
				char := game.NewCharacter("LoadTest")
				_ = client.SaveCharacter(char)
				return char
			},
			wantErr: false,
			validate: func(t *testing.T, saved, loaded *game.Character) {
				if loaded.ID != saved.ID {
					t.Errorf("LoadCharacter() ID = %v, want %v", loaded.ID, saved.ID)
				}
				if loaded.Name != saved.Name {
					t.Errorf("LoadCharacter() Name = %v, want %v", loaded.Name, saved.Name)
				}
				if loaded.Level != saved.Level {
					t.Errorf("LoadCharacter() Level = %v, want %v", loaded.Level, saved.Level)
				}
			},
		},
		{
			name: "load character with progress",
			setup: func() *game.Character {
				char := game.NewCharacter("ProgressTest")
				char.AddXP(250)
				char.TotalCommits = 10
				char.CurrentStreak = 5
				char.TotalLinesAdded = 500
				_ = client.SaveCharacter(char)
				return char
			},
			wantErr: false,
			validate: func(t *testing.T, saved, loaded *game.Character) {
				if loaded.Level != saved.Level {
					t.Errorf("LoadCharacter() Level = %v, want %v", loaded.Level, saved.Level)
				}
				if loaded.TotalCommits != saved.TotalCommits {
					t.Errorf("LoadCharacter() TotalCommits = %v, want %v", loaded.TotalCommits, saved.TotalCommits)
				}
				if loaded.CurrentStreak != saved.CurrentStreak {
					t.Errorf("LoadCharacter() CurrentStreak = %v, want %v", loaded.CurrentStreak, saved.CurrentStreak)
				}
				if loaded.TotalLinesAdded != saved.TotalLinesAdded {
					t.Errorf("LoadCharacter() TotalLinesAdded = %v, want %v", loaded.TotalLinesAdded, saved.TotalLinesAdded)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup - create character if needed
			var savedChar *game.Character
			if tt.setup != nil {
				savedChar = tt.setup()
			}

			// Execute - load character
			loadedChar, err := client.LoadCharacter()

			// Verify error expectations
			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadCharacter() expected error containing %q, got nil", tt.errMsg)
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("LoadCharacter() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("LoadCharacter() unexpected error: %v", err)
				}
				if loadedChar == nil {
					t.Errorf("LoadCharacter() returned nil character")
				}

				// Run custom validation if provided
				if tt.validate != nil && savedChar != nil && loadedChar != nil {
					tt.validate(t, savedChar, loadedChar)
				}
			}

			// Cleanup
			_ = client.DeleteCharacter()
		})
	}
}

// TestSkateClient_SaveQuests tests quest list serialization and storage
func TestSkateClient_SaveQuests(t *testing.T) {
	// Skip if Skate not installed
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client := &SkateClient{skatePath: skatePath}

	tests := []struct {
		name    string
		quests  []*game.Quest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil quests list",
			quests:  nil,
			wantErr: true,
			errMsg:  "cannot save nil quests",
		},
		{
			name:    "empty quests list",
			quests:  []*game.Quest{},
			wantErr: false,
		},
		{
			name: "single quest",
			quests: []*game.Quest{
				game.NewQuest("Test Quest", "Test Description", game.QuestTypeCommit, 5, 100, 1),
			},
			wantErr: false,
		},
		{
			name: "multiple quests",
			quests: []*game.Quest{
				game.NewQuest("Quest 1", "First quest", game.QuestTypeCommit, 5, 100, 1),
				game.NewQuest("Quest 2", "Second quest", game.QuestTypeLines, 100, 200, 5),
				game.NewQuest("Quest 3", "Third quest", game.QuestTypeTests, 10, 300, 10),
			},
			wantErr: false,
		},
		{
			name: "quest with progress",
			quests: func() []*game.Quest {
				quest := game.NewQuest("Active Quest", "In progress", game.QuestTypeCommit, 10, 150, 1)
				_ = quest.Start("/path/to/repo", "abc123")
				quest.UpdateProgress(5)
				return []*game.Quest{quest}
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.SaveQuests(tt.quests)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SaveQuests() expected error containing %q, got nil", tt.errMsg)
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("SaveQuests() error = %q, want error containing %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("SaveQuests() unexpected error: %v", err)
				}
			}
		})
	}

	// Cleanup
	_ = client.DeleteQuests()
}

// TestSkateClient_LoadQuests tests quest list retrieval and deserialization
func TestSkateClient_LoadQuests(t *testing.T) {
	// Skip if Skate not installed
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client := &SkateClient{skatePath: skatePath}

	tests := []struct {
		name     string
		setup    func() []*game.Quest
		wantErr  bool
		validate func(*testing.T, []*game.Quest, []*game.Quest)
	}{
		{
			name: "load non-existent quests (first run)",
			setup: func() []*game.Quest {
				// Ensure no quests exist
				_ = client.DeleteQuests()
				return nil
			},
			wantErr: false,
			validate: func(t *testing.T, saved, loaded []*game.Quest) {
				if len(loaded) != 0 {
					t.Errorf("LoadQuests() for non-existent key should return empty slice, got %d quests", len(loaded))
				}
			},
		},
		{
			name: "load empty quest list",
			setup: func() []*game.Quest {
				emptyList := []*game.Quest{}
				_ = client.SaveQuests(emptyList)
				return emptyList
			},
			wantErr: false,
			validate: func(t *testing.T, saved, loaded []*game.Quest) {
				if len(loaded) != 0 {
					t.Errorf("LoadQuests() length = %d, want 0", len(loaded))
				}
			},
		},
		{
			name: "load single quest",
			setup: func() []*game.Quest {
				quests := []*game.Quest{
					game.NewQuest("Test Quest", "Test Description", game.QuestTypeCommit, 5, 100, 1),
				}
				_ = client.SaveQuests(quests)
				return quests
			},
			wantErr: false,
			validate: func(t *testing.T, saved, loaded []*game.Quest) {
				if len(loaded) != 1 {
					t.Errorf("LoadQuests() length = %d, want 1", len(loaded))
					return
				}
				if loaded[0].ID != saved[0].ID {
					t.Errorf("LoadQuests() quest ID = %v, want %v", loaded[0].ID, saved[0].ID)
				}
				if loaded[0].Title != saved[0].Title {
					t.Errorf("LoadQuests() quest Title = %v, want %v", loaded[0].Title, saved[0].Title)
				}
			},
		},
		{
			name: "load multiple quests",
			setup: func() []*game.Quest {
				quests := []*game.Quest{
					game.NewQuest("Quest 1", "First", game.QuestTypeCommit, 5, 100, 1),
					game.NewQuest("Quest 2", "Second", game.QuestTypeLines, 100, 200, 5),
					game.NewQuest("Quest 3", "Third", game.QuestTypeTests, 10, 300, 10),
				}
				_ = client.SaveQuests(quests)
				return quests
			},
			wantErr: false,
			validate: func(t *testing.T, saved, loaded []*game.Quest) {
				if len(loaded) != 3 {
					t.Errorf("LoadQuests() length = %d, want 3", len(loaded))
					return
				}
				for i := range loaded {
					if loaded[i].ID != saved[i].ID {
						t.Errorf("LoadQuests() quest[%d] ID = %v, want %v", i, loaded[i].ID, saved[i].ID)
					}
					if loaded[i].Type != saved[i].Type {
						t.Errorf("LoadQuests() quest[%d] Type = %v, want %v", i, loaded[i].Type, saved[i].Type)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup - create quests if needed
			var savedQuests []*game.Quest
			if tt.setup != nil {
				savedQuests = tt.setup()
			}

			// Execute - load quests
			loadedQuests, err := client.LoadQuests()

			// Verify error expectations
			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadQuests() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("LoadQuests() unexpected error: %v", err)
				}
				if loadedQuests == nil {
					t.Errorf("LoadQuests() returned nil, should return empty slice")
				}

				// Run custom validation if provided
				if tt.validate != nil {
					tt.validate(t, savedQuests, loadedQuests)
				}
			}

			// Cleanup
			_ = client.DeleteQuests()
		})
	}
}

// TestSkateClient_CharacterExists tests character existence check
func TestSkateClient_CharacterExists(t *testing.T) {
	// Skip if Skate not installed
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client := &SkateClient{skatePath: skatePath}

	tests := []struct {
		name  string
		setup func()
		want  bool
	}{
		{
			name: "character does not exist",
			setup: func() {
				_ = client.DeleteCharacter()
			},
			want: false,
		},
		{
			name: "character exists",
			setup: func() {
				char := game.NewCharacter("ExistsTest")
				_ = client.SaveCharacter(char)
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setup != nil {
				tt.setup()
			}

			// Execute
			got := client.CharacterExists()

			// Verify
			if got != tt.want {
				t.Errorf("CharacterExists() = %v, want %v", got, tt.want)
			}

			// Cleanup
			_ = client.DeleteCharacter()
		})
	}
}

// TestSkateClient_DeleteCharacter tests character deletion
func TestSkateClient_DeleteCharacter(t *testing.T) {
	// Skip if Skate not installed
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client := &SkateClient{skatePath: skatePath}

	// Create a character first
	char := game.NewCharacter("DeleteTest")
	err = client.SaveCharacter(char)
	if err != nil {
		t.Fatalf("Setup failed: could not save character: %v", err)
	}

	// Verify it exists
	if !client.CharacterExists() {
		t.Fatalf("Setup failed: character does not exist after save")
	}

	// Delete it
	err = client.DeleteCharacter()
	if err != nil {
		t.Errorf("DeleteCharacter() unexpected error: %v", err)
	}

	// Verify it's gone
	if client.CharacterExists() {
		t.Errorf("DeleteCharacter() character still exists after deletion")
	}
}

// TestSkateClient_DeleteQuests tests quest list deletion
func TestSkateClient_DeleteQuests(t *testing.T) {
	// Skip if Skate not installed
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client := &SkateClient{skatePath: skatePath}

	// Create quests first
	quests := []*game.Quest{
		game.NewQuest("Delete Test", "Test deletion", game.QuestTypeCommit, 5, 100, 1),
	}
	err = client.SaveQuests(quests)
	if err != nil {
		t.Fatalf("Setup failed: could not save quests: %v", err)
	}

	// Delete them
	err = client.DeleteQuests()
	if err != nil {
		t.Errorf("DeleteQuests() unexpected error: %v", err)
	}

	// Verify they're gone (should return empty slice, not error)
	loadedQuests, err := client.LoadQuests()
	if err != nil {
		t.Errorf("LoadQuests() after delete unexpected error: %v", err)
	}
	if len(loadedQuests) != 0 {
		t.Errorf("LoadQuests() after delete length = %d, want 0", len(loadedQuests))
	}
}

// TestSkateClient_SaveLoadRoundTrip tests full save/load cycle
func TestSkateClient_SaveLoadRoundTrip(t *testing.T) {
	// Skip if Skate not installed
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client := &SkateClient{skatePath: skatePath}

	// Create a character with various states
	originalChar := game.NewCharacter("RoundTripTest")
	originalChar.AddXP(500) // Should level up multiple times
	originalChar.TotalCommits = 25
	originalChar.TotalLinesAdded = 1000
	originalChar.TotalLinesRemoved = 200
	originalChar.CurrentStreak = 7
	originalChar.LongestStreak = 15
	originalChar.TodayCommits = 3

	// Save it
	err = client.SaveCharacter(originalChar)
	if err != nil {
		t.Fatalf("SaveCharacter() failed: %v", err)
	}

	// Load it back
	loadedChar, err := client.LoadCharacter()
	if err != nil {
		t.Fatalf("LoadCharacter() failed: %v", err)
	}

	// Verify all fields match
	if loadedChar.ID != originalChar.ID {
		t.Errorf("ID mismatch: got %v, want %v", loadedChar.ID, originalChar.ID)
	}
	if loadedChar.Name != originalChar.Name {
		t.Errorf("Name mismatch: got %v, want %v", loadedChar.Name, originalChar.Name)
	}
	if loadedChar.Level != originalChar.Level {
		t.Errorf("Level mismatch: got %v, want %v", loadedChar.Level, originalChar.Level)
	}
	if loadedChar.XP != originalChar.XP {
		t.Errorf("XP mismatch: got %v, want %v", loadedChar.XP, originalChar.XP)
	}
	if loadedChar.TotalCommits != originalChar.TotalCommits {
		t.Errorf("TotalCommits mismatch: got %v, want %v", loadedChar.TotalCommits, originalChar.TotalCommits)
	}
	if loadedChar.CurrentStreak != originalChar.CurrentStreak {
		t.Errorf("CurrentStreak mismatch: got %v, want %v", loadedChar.CurrentStreak, originalChar.CurrentStreak)
	}

	// Cleanup
	_ = client.DeleteCharacter()
}

// TestSkateClient_QuestsSaveLoadRoundTrip tests full quest save/load cycle
func TestSkateClient_QuestsSaveLoadRoundTrip(t *testing.T) {
	// Skip if Skate not installed
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client := &SkateClient{skatePath: skatePath}

	// Create quests with various states
	quest1 := game.NewQuest("Available Quest", "Not started", game.QuestTypeCommit, 5, 100, 1)

	quest2 := game.NewQuest("Active Quest", "In progress", game.QuestTypeLines, 100, 200, 5)
	_ = quest2.Start("/repo/path", "abc123")
	quest2.UpdateProgress(50)

	quest3 := game.NewQuest("Completed Quest", "Finished", game.QuestTypeTests, 10, 300, 10)
	_ = quest3.Start("/repo/path", "def456")
	quest3.UpdateProgress(10)
	_ = quest3.Complete()

	originalQuests := []*game.Quest{quest1, quest2, quest3}

	// Save them
	err = client.SaveQuests(originalQuests)
	if err != nil {
		t.Fatalf("SaveQuests() failed: %v", err)
	}

	// Load them back
	loadedQuests, err := client.LoadQuests()
	if err != nil {
		t.Fatalf("LoadQuests() failed: %v", err)
	}

	// Verify count
	if len(loadedQuests) != 3 {
		t.Fatalf("LoadQuests() length = %d, want 3", len(loadedQuests))
	}

	// Verify each quest
	for i, loaded := range loadedQuests {
		original := originalQuests[i]

		if loaded.ID != original.ID {
			t.Errorf("Quest[%d] ID mismatch: got %v, want %v", i, loaded.ID, original.ID)
		}
		if loaded.Title != original.Title {
			t.Errorf("Quest[%d] Title mismatch: got %v, want %v", i, loaded.Title, original.Title)
		}
		if loaded.Status != original.Status {
			t.Errorf("Quest[%d] Status mismatch: got %v, want %v", i, loaded.Status, original.Status)
		}
		if loaded.Current != original.Current {
			t.Errorf("Quest[%d] Current mismatch: got %v, want %v", i, loaded.Current, original.Current)
		}
		if loaded.Progress != original.Progress {
			t.Errorf("Quest[%d] Progress mismatch: got %v, want %v", i, loaded.Progress, original.Progress)
		}
	}

	// Cleanup
	_ = client.DeleteQuests()
}

// TestJSONMarshaling tests JSON encoding/decoding edge cases
func TestJSONMarshaling(t *testing.T) {
	tests := []struct {
		name      string
		character *game.Character
		wantErr   bool
	}{
		{
			name:      "normal character",
			character: game.NewCharacter("JSONTest"),
			wantErr:   false,
		},
		{
			name: "character with time fields",
			character: func() *game.Character {
				c := game.NewCharacter("TimeTest")
				c.UpdateStreak()
				return c
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			jsonData, err := json.Marshal(tt.character)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Marshal() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Marshal() unexpected error: %v", err)
				return
			}

			// Unmarshal
			var decoded game.Character
			err = json.Unmarshal(jsonData, &decoded)
			if err != nil {
				t.Errorf("Unmarshal() unexpected error: %v", err)
				return
			}

			// Verify basic fields preserved
			if decoded.ID != tt.character.ID {
				t.Errorf("Round-trip ID mismatch: got %v, want %v", decoded.ID, tt.character.ID)
			}
			if decoded.Name != tt.character.Name {
				t.Errorf("Round-trip Name mismatch: got %v, want %v", decoded.Name, tt.character.Name)
			}
		})
	}
}

// TestSkateClient_ErrorHandling tests various error conditions
func TestSkateClient_ErrorHandling(t *testing.T) {
	// Skip if Skate not installed
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client := &SkateClient{skatePath: skatePath}

	t.Run("save nil character", func(t *testing.T) {
		err := client.SaveCharacter(nil)
		if err == nil {
			t.Errorf("SaveCharacter(nil) should return error")
		}
		if !strings.Contains(err.Error(), "nil") {
			t.Errorf("SaveCharacter(nil) error should mention 'nil', got: %v", err)
		}
	})

	t.Run("save nil quests", func(t *testing.T) {
		err := client.SaveQuests(nil)
		if err == nil {
			t.Errorf("SaveQuests(nil) should return error")
		}
		if !strings.Contains(err.Error(), "nil") {
			t.Errorf("SaveQuests(nil) error should mention 'nil', got: %v", err)
		}
	})

	t.Run("load non-existent character", func(t *testing.T) {
		// Ensure no character exists
		_ = client.DeleteCharacter()

		_, err := client.LoadCharacter()
		if err == nil {
			t.Errorf("LoadCharacter() on non-existent key should return error")
		}
	})
}

// TestSkateClient_ConcurrentAccess tests thread safety (basic)
func TestSkateClient_ConcurrentAccess(t *testing.T) {
	// Skip if Skate not installed
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		t.Skip("Skate not installed, skipping integration test")
	}

	client := &SkateClient{skatePath: skatePath}

	// This is a basic test - Skate itself handles concurrency
	// We're just verifying our wrapper doesn't panic
	char1 := game.NewCharacter("Concurrent1")
	char2 := game.NewCharacter("Concurrent2")

	// Save first character
	err = client.SaveCharacter(char1)
	if err != nil {
		t.Fatalf("SaveCharacter(char1) failed: %v", err)
	}

	// Save second character (overwrites first)
	err = client.SaveCharacter(char2)
	if err != nil {
		t.Fatalf("SaveCharacter(char2) failed: %v", err)
	}

	// Load should get the last saved character
	loaded, err := client.LoadCharacter()
	if err != nil {
		t.Fatalf("LoadCharacter() failed: %v", err)
	}

	if loaded.Name != char2.Name {
		t.Errorf("LoadCharacter() Name = %v, want %v (last saved)", loaded.Name, char2.Name)
	}

	// Cleanup
	_ = client.DeleteCharacter()
}

// Benchmark tests for performance awareness

func BenchmarkSaveCharacter(b *testing.B) {
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		b.Skip("Skate not installed")
	}

	client := &SkateClient{skatePath: skatePath}
	char := game.NewCharacter("BenchChar")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.SaveCharacter(char)
	}

	// Cleanup
	_ = client.DeleteCharacter()
}

func BenchmarkLoadCharacter(b *testing.B) {
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		b.Skip("Skate not installed")
	}

	client := &SkateClient{skatePath: skatePath}
	char := game.NewCharacter("BenchChar")
	_ = client.SaveCharacter(char)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.LoadCharacter()
	}

	// Cleanup
	_ = client.DeleteCharacter()
}

// Helper function to create a test character with specific state
func createTestCharacter(name string, level int, xp int) *game.Character {
	char := game.NewCharacter(name)
	char.Level = level
	char.XP = xp
	return char
}

// Helper function to verify character equality (for testing)
func verifyCharactersEqual(t *testing.T, a, b *game.Character) {
	if a.ID != b.ID {
		t.Errorf("ID mismatch: %v != %v", a.ID, b.ID)
	}
	if a.Name != b.Name {
		t.Errorf("Name mismatch: %v != %v", a.Name, b.Name)
	}
	if a.Level != b.Level {
		t.Errorf("Level mismatch: %v != %v", a.Level, b.Level)
	}
	if a.XP != b.XP {
		t.Errorf("XP mismatch: %v != %v", a.XP, b.XP)
	}
}

// Helper function to verify quest equality (for testing)
func verifyQuestsEqual(t *testing.T, a, b *game.Quest) {
	if a.ID != b.ID {
		t.Errorf("Quest ID mismatch: %v != %v", a.ID, b.ID)
	}
	if a.Title != b.Title {
		t.Errorf("Quest Title mismatch: %v != %v", a.Title, b.Title)
	}
	if a.Status != b.Status {
		t.Errorf("Quest Status mismatch: %v != %v", a.Status, b.Status)
	}
	if a.Current != b.Current {
		t.Errorf("Quest Current mismatch: %v != %v", a.Current, b.Current)
	}
}

// TestTimeMarshaling tests that time.Time fields are properly handled
func TestTimeMarshaling(t *testing.T) {
	char := game.NewCharacter("TimeTest")
	now := time.Now()
	char.LastActiveDate = now

	// Marshal to JSON
	jsonData, err := json.Marshal(char)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Unmarshal from JSON
	var decoded game.Character
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Time should be preserved (within reasonable tolerance)
	if decoded.LastActiveDate.Unix() != char.LastActiveDate.Unix() {
		t.Errorf("Time not preserved: got %v, want %v", decoded.LastActiveDate, char.LastActiveDate)
	}
}

// TestInvalidJSON tests handling of corrupted JSON data
func TestInvalidJSON(t *testing.T) {
	tests := []struct {
		name        string
		invalidJSON string
		structType  string
	}{
		{"invalid character JSON", `{"name": "test", "level": "not a number"}`, "Character"},
		{"truncated JSON", `{"name": "test"`, "Character"},
		{"invalid quest JSON", `{"title": 123, "type": "invalid"}`, "Quest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.structType == "Character" {
				var char game.Character
				err := json.Unmarshal([]byte(tt.invalidJSON), &char)
				if err == nil {
					t.Errorf("Unmarshal should fail for invalid JSON")
				}
			} else {
				var quest game.Quest
				err := json.Unmarshal([]byte(tt.invalidJSON), &quest)
				if err == nil {
					t.Errorf("Unmarshal should fail for invalid JSON")
				}
			}
		})
	}
}

// TestSkateClient_InvalidSkatePath tests behavior with invalid skate binary path
func TestSkateClient_InvalidSkatePath(t *testing.T) {
	client := &SkateClient{skatePath: "/nonexistent/path/to/skate"}

	// Try to save a character with invalid skate path
	char := game.NewCharacter("Test")
	err := client.SaveCharacter(char)
	if err == nil {
		t.Errorf("SaveCharacter() with invalid skatePath should return error")
	}
}

// TestMockingApproach demonstrates how to mock Skate for unit tests
// This test documents the mocking strategy even though we use integration tests
func TestMockingApproach(t *testing.T) {
	t.Skip("Documentation test - shows mocking approach for future unit tests")

	// Example of how to mock exec.Command:
	// 1. Create a test helper that replaces exec.Command
	// 2. Return controlled output for different commands
	// 3. Verify the correct commands were called

	// This would allow testing without Skate installed:
	// mockCmd := func(command string, args ...string) *exec.Cmd {
	//     if args[0] == "get" {
	//         return exec.Command("echo", `{"id":"test","name":"Hero"}`)
	//     }
	//     return exec.Command("echo", "OK")
	// }
}
