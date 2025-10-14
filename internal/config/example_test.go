package config_test

import (
	"fmt"
	"log"

	"github.com/AutumnsGrove/codequest/internal/config"
)

// ExampleLoad demonstrates loading and using configuration.
func ExampleLoad() {
	// Load configuration from ~/.config/codequest/config.toml
	// If it doesn't exist, it will be created with default values
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	// Use the configuration
	fmt.Printf("Character: %s\n", cfg.Character.Name)
	fmt.Printf("Difficulty: %s\n", cfg.Game.Difficulty)
}

// ExampleDefaultConfig demonstrates creating a config with defaults.
func ExampleDefaultConfig() {
	// Get default configuration
	cfg := config.DefaultConfig()

	// Defaults are already set
	fmt.Printf("Default character name: %s\n", cfg.Character.Name)
	fmt.Printf("Default theme: %s\n", cfg.UI.Theme)
	fmt.Printf("Default difficulty: %s\n", cfg.Game.Difficulty)

	// Output:
	// Default character name: CodeWarrior
	// Default theme: dark
	// Default difficulty: normal
}

// ExampleConfig_Validate demonstrates configuration validation.
func ExampleConfig_Validate() {
	cfg := config.DefaultConfig()

	// This is valid
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Configuration is valid")
	}

	// Make it invalid
	cfg.Game.Difficulty = "impossible"
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Output:
	// Configuration is valid
	// Error: config validation error [game.difficulty]: must be one of: easy, normal, hard (value: impossible)
}

// ExampleExpandPath demonstrates path expansion with tilde.
func ExampleExpandPath() {
	// Expand paths with ~ to home directory
	expanded, err := config.ExpandPath("~/projects/myapp")
	if err != nil {
		log.Fatal(err)
	}

	// The path is now absolute
	fmt.Printf("Expanded path starts with /: %v\n", expanded[0] == '/')

	// Absolute paths are unchanged
	absolute := "/usr/local/bin"
	result, _ := config.ExpandPath(absolute)
	fmt.Printf("Absolute unchanged: %v\n", result == absolute)

	// Output:
	// Expanded path starts with /: true
	// Absolute unchanged: true
}

// ExampleExpandPaths demonstrates batch path expansion.
func ExampleExpandPaths() {
	paths := []string{"~/projects", "/absolute/path", "~/code"}
	expanded, err := config.ExpandPaths(paths)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Number of paths: %d\n", len(expanded))
	fmt.Printf("All paths are absolute: %v\n",
		len(expanded[0]) > 0 && expanded[0][0] == '/' &&
		expanded[1][0] == '/' &&
		len(expanded[2]) > 0 && expanded[2][0] == '/')

	// Output:
	// Number of paths: 3
	// All paths are absolute: true
}
