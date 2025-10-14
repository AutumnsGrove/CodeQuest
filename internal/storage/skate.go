// Package storage provides data persistence for CodeQuest using Skate KV store.
// Skate is a key-value store from Charm that provides encrypted, cloud-synced storage.
package storage

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// Skate key names used for storing game data
const (
	KeyCharacter = "codequest.character" // Character data storage key
	KeyQuests    = "codequest.quests"    // Quests list storage key
)

// SkateClient provides a wrapper around the Skate CLI for data persistence.
// It handles JSON serialization and CLI interaction for saving/loading game data.
type SkateClient struct {
	// skatePath is the path to the skate binary (default: "skate" in PATH)
	skatePath string
}

// NewSkateClient creates a new Skate storage client.
// It verifies that the Skate CLI is available in the system PATH.
//
// Returns:
//   - *SkateClient: A new Skate client instance
//   - error: An error if Skate is not installed or not found in PATH
func NewSkateClient() (*SkateClient, error) {
	// Check if skate is installed and available in PATH
	skatePath, err := exec.LookPath("skate")
	if err != nil {
		return nil, fmt.Errorf("skate CLI not found in PATH: %w (install from https://github.com/charmbracelet/skate)", err)
	}

	return &SkateClient{
		skatePath: skatePath,
	}, nil
}

// SaveCharacter persists a character to Skate storage.
// The character is serialized to JSON before being stored.
//
// Parameters:
//   - character: The character to save (must not be nil)
//
// Returns:
//   - error: An error if serialization or storage fails
func (s *SkateClient) SaveCharacter(character *game.Character) error {
	if character == nil {
		return fmt.Errorf("cannot save nil character")
	}

	// Marshal character to JSON
	jsonData, err := json.Marshal(character)
	if err != nil {
		return fmt.Errorf("failed to marshal character to JSON: %w", err)
	}

	// Store in Skate using: skate set <key> <value>
	if err := s.setKey(KeyCharacter, string(jsonData)); err != nil {
		return fmt.Errorf("failed to save character to Skate: %w", err)
	}

	return nil
}

// LoadCharacter retrieves a character from Skate storage.
// The stored JSON is deserialized into a Character struct.
//
// Returns:
//   - *game.Character: The loaded character
//   - error: An error if the character doesn't exist, or if retrieval/deserialization fails
func (s *SkateClient) LoadCharacter() (*game.Character, error) {
	// Retrieve from Skate using: skate get <key>
	jsonData, err := s.getKey(KeyCharacter)
	if err != nil {
		return nil, fmt.Errorf("failed to load character from Skate: %w", err)
	}

	// Unmarshal JSON to Character struct
	var character game.Character
	if err := json.Unmarshal([]byte(jsonData), &character); err != nil {
		return nil, fmt.Errorf("failed to unmarshal character JSON: %w", err)
	}

	return &character, nil
}

// SaveQuests persists a list of quests to Skate storage.
// The quest list is serialized to JSON before being stored.
//
// Parameters:
//   - quests: The list of quests to save (can be empty, but not nil)
//
// Returns:
//   - error: An error if serialization or storage fails
func (s *SkateClient) SaveQuests(quests []*game.Quest) error {
	if quests == nil {
		return fmt.Errorf("cannot save nil quests list (use empty slice instead)")
	}

	// Marshal quests to JSON
	jsonData, err := json.Marshal(quests)
	if err != nil {
		return fmt.Errorf("failed to marshal quests to JSON: %w", err)
	}

	// Store in Skate using: skate set <key> <value>
	if err := s.setKey(KeyQuests, string(jsonData)); err != nil {
		return fmt.Errorf("failed to save quests to Skate: %w", err)
	}

	return nil
}

// LoadQuests retrieves a list of quests from Skate storage.
// The stored JSON is deserialized into a slice of Quest pointers.
//
// Returns:
//   - []*game.Quest: The loaded quests (empty slice if no quests exist)
//   - error: An error if retrieval or deserialization fails
func (s *SkateClient) LoadQuests() ([]*game.Quest, error) {
	// Retrieve from Skate using: skate get <key>
	jsonData, err := s.getKey(KeyQuests)
	if err != nil {
		// If key doesn't exist, return empty slice instead of error
		// This is expected on first run
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no such key") {
			return []*game.Quest{}, nil
		}
		return nil, fmt.Errorf("failed to load quests from Skate: %w", err)
	}

	// Unmarshal JSON to Quest slice
	var quests []*game.Quest
	if err := json.Unmarshal([]byte(jsonData), &quests); err != nil {
		return nil, fmt.Errorf("failed to unmarshal quests JSON: %w", err)
	}

	// Return empty slice if quests is nil (defensive)
	if quests == nil {
		return []*game.Quest{}, nil
	}

	return quests, nil
}

// DeleteCharacter removes the character from Skate storage.
// This is useful for starting fresh or resetting progress.
//
// Returns:
//   - error: An error if deletion fails
func (s *SkateClient) DeleteCharacter() error {
	if err := s.deleteKey(KeyCharacter); err != nil {
		return fmt.Errorf("failed to delete character from Skate: %w", err)
	}
	return nil
}

// DeleteQuests removes all quests from Skate storage.
// This is useful for starting fresh or resetting progress.
//
// Returns:
//   - error: An error if deletion fails
func (s *SkateClient) DeleteQuests() error {
	if err := s.deleteKey(KeyQuests); err != nil {
		return fmt.Errorf("failed to delete quests from Skate: %w", err)
	}
	return nil
}

// CharacterExists checks if a character is stored in Skate.
// This is useful for determining if this is a first run.
//
// Returns:
//   - bool: true if a character exists, false otherwise
func (s *SkateClient) CharacterExists() bool {
	_, err := s.getKey(KeyCharacter)
	return err == nil
}

// setKey stores a value in Skate using the CLI.
// Executes: skate set <key> <value>
//
// Parameters:
//   - key: The key to store the value under
//   - value: The value to store (should be JSON string)
//
// Returns:
//   - error: An error if the CLI command fails
func (s *SkateClient) setKey(key, value string) error {
	// Execute: skate set <key> <value>
	cmd := exec.Command(s.skatePath, "set", key, value)

	// Capture both stdout and stderr for error reporting
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("skate set failed: %w (output: %s)", err, string(output))
	}

	return nil
}

// getKey retrieves a value from Skate using the CLI.
// Executes: skate get <key>
//
// Parameters:
//   - key: The key to retrieve
//
// Returns:
//   - string: The stored value (trimmed of whitespace)
//   - error: An error if the key doesn't exist or the CLI command fails
func (s *SkateClient) getKey(key string) (string, error) {
	// Execute: skate get <key>
	cmd := exec.Command(s.skatePath, "get", key)

	// Capture output
	output, err := cmd.Output()
	if err != nil {
		// Check if it's a "not found" error
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)
			if strings.Contains(stderr, "not found") || strings.Contains(stderr, "no such key") {
				return "", fmt.Errorf("key %s not found in Skate", key)
			}
			return "", fmt.Errorf("skate get failed: %w (stderr: %s)", err, stderr)
		}
		return "", fmt.Errorf("skate get failed: %w", err)
	}

	// Trim whitespace and return
	return strings.TrimSpace(string(output)), nil
}

// deleteKey removes a value from Skate using the CLI.
// Executes: skate delete <key>
//
// Parameters:
//   - key: The key to delete
//
// Returns:
//   - error: An error if the CLI command fails
func (s *SkateClient) deleteKey(key string) error {
	// Execute: skate delete <key>
	cmd := exec.Command(s.skatePath, "delete", key)

	// Capture both stdout and stderr for error reporting
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Deleting a non-existent key might not be an error for some use cases
		// but we'll report it anyway
		return fmt.Errorf("skate delete failed: %w (output: %s)", err, string(output))
	}

	return nil
}
