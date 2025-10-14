package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
)

// TestDefaultConfig verifies that DefaultConfig returns a valid configuration.
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Check character defaults
	if cfg.Character.Name != "CodeWarrior" {
		t.Errorf("expected character name 'CodeWarrior', got %q", cfg.Character.Name)
	}

	// Check game defaults
	if cfg.Game.Difficulty != "normal" {
		t.Errorf("expected difficulty 'normal', got %q", cfg.Game.Difficulty)
	}
	if cfg.Game.AutoStartQuests != false {
		t.Error("expected auto_start_quests to be false")
	}
	if cfg.Game.ShowTips != true {
		t.Error("expected show_tips to be true")
	}

	// Check UI defaults
	if cfg.UI.Theme != "dark" {
		t.Errorf("expected theme 'dark', got %q", cfg.UI.Theme)
	}
	if cfg.UI.ShowAnimations != true {
		t.Error("expected show_animations to be true")
	}

	// Check AI defaults
	if cfg.AI.Mentor.Provider != "crush" {
		t.Errorf("expected AI mentor provider 'crush', got %q", cfg.AI.Mentor.Provider)
	}
	if cfg.AI.Mentor.Temperature != 0.7 {
		t.Errorf("expected temperature 0.7, got %f", cfg.AI.Mentor.Temperature)
	}

	// Validate that defaults are valid
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config should be valid, got error: %v", err)
	}
}

// TestValidate_ValidConfig tests validation with valid configurations.
func TestValidate_ValidConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
	}{
		{
			name: "default config",
			cfg:  DefaultConfig(),
		},
		{
			name: "easy difficulty",
			cfg: &Config{
				Character: CharacterConfig{Name: "Test"},
				Game:      GameConfig{Difficulty: "easy"},
				UI:        UIConfig{Theme: "light"},
				AI: AIConfig{
					Mentor: AIMentorConfig{Provider: "mods", Temperature: 0.5},
					Review: AIReviewConfig{Provider: "crush"},
				},
				Debug: DebugConfig{LogLevel: "debug"},
			},
		},
		{
			name: "hard difficulty",
			cfg: &Config{
				Character: CharacterConfig{Name: "Warrior"},
				Game:      GameConfig{Difficulty: "hard"},
				UI:        UIConfig{Theme: "auto"},
				AI: AIConfig{
					Mentor: AIMentorConfig{Provider: "claude-code", Temperature: 1.5},
					Review: AIReviewConfig{Provider: "mods"},
				},
				Debug: DebugConfig{LogLevel: "error"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); err != nil {
				t.Errorf("expected valid config, got error: %v", err)
			}
		})
	}
}

// TestValidate_InvalidConfig tests validation with invalid configurations.
func TestValidate_InvalidConfig(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		wantField string
	}{
		{
			name: "invalid difficulty",
			cfg: &Config{
				Character: CharacterConfig{Name: "Test"},
				Game:      GameConfig{Difficulty: "impossible"},
				UI:        UIConfig{Theme: "dark"},
				AI: AIConfig{
					Mentor: AIMentorConfig{Provider: "crush", Temperature: 0.7},
					Review: AIReviewConfig{Provider: "mods"},
				},
				Debug: DebugConfig{LogLevel: "info"},
			},
			wantField: "game.difficulty",
		},
		{
			name: "invalid theme",
			cfg: &Config{
				Character: CharacterConfig{Name: "Test"},
				Game:      GameConfig{Difficulty: "normal"},
				UI:        UIConfig{Theme: "rainbow"},
				AI: AIConfig{
					Mentor: AIMentorConfig{Provider: "crush", Temperature: 0.7},
					Review: AIReviewConfig{Provider: "mods"},
				},
				Debug: DebugConfig{LogLevel: "info"},
			},
			wantField: "ui.theme",
		},
		{
			name: "invalid AI provider",
			cfg: &Config{
				Character: CharacterConfig{Name: "Test"},
				Game:      GameConfig{Difficulty: "normal"},
				UI:        UIConfig{Theme: "dark"},
				AI: AIConfig{
					Mentor: AIMentorConfig{Provider: "gpt-4", Temperature: 0.7},
					Review: AIReviewConfig{Provider: "mods"},
				},
				Debug: DebugConfig{LogLevel: "info"},
			},
			wantField: "ai.mentor.provider",
		},
		{
			name: "temperature too high",
			cfg: &Config{
				Character: CharacterConfig{Name: "Test"},
				Game:      GameConfig{Difficulty: "normal"},
				UI:        UIConfig{Theme: "dark"},
				AI: AIConfig{
					Mentor: AIMentorConfig{Provider: "crush", Temperature: 3.0},
					Review: AIReviewConfig{Provider: "mods"},
				},
				Debug: DebugConfig{LogLevel: "info"},
			},
			wantField: "ai.mentor.temperature",
		},
		{
			name: "temperature negative",
			cfg: &Config{
				Character: CharacterConfig{Name: "Test"},
				Game:      GameConfig{Difficulty: "normal"},
				UI:        UIConfig{Theme: "dark"},
				AI: AIConfig{
					Mentor: AIMentorConfig{Provider: "crush", Temperature: -0.5},
					Review: AIReviewConfig{Provider: "mods"},
				},
				Debug: DebugConfig{LogLevel: "info"},
			},
			wantField: "ai.mentor.temperature",
		},
		{
			name: "invalid log level",
			cfg: &Config{
				Character: CharacterConfig{Name: "Test"},
				Game:      GameConfig{Difficulty: "normal"},
				UI:        UIConfig{Theme: "dark"},
				AI: AIConfig{
					Mentor: AIMentorConfig{Provider: "crush", Temperature: 0.7},
					Review: AIReviewConfig{Provider: "mods"},
				},
				Debug: DebugConfig{LogLevel: "trace"},
			},
			wantField: "debug.log_level",
		},
		{
			name: "empty character name",
			cfg: &Config{
				Character: CharacterConfig{Name: ""},
				Game:      GameConfig{Difficulty: "normal"},
				UI:        UIConfig{Theme: "dark"},
				AI: AIConfig{
					Mentor: AIMentorConfig{Provider: "crush", Temperature: 0.7},
					Review: AIReviewConfig{Provider: "mods"},
				},
				Debug: DebugConfig{LogLevel: "info"},
			},
			wantField: "character.name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if err == nil {
				t.Error("expected validation error, got nil")
				return
			}

			verr, ok := err.(ValidationError)
			if !ok {
				t.Errorf("expected ValidationError, got %T", err)
				return
			}

			if verr.Field != tt.wantField {
				t.Errorf("expected error for field %q, got %q", tt.wantField, verr.Field)
			}
		})
	}
}

// TestExpandPath tests path expansion with ~ for home directory.
func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde alone",
			input:    "~",
			expected: home,
		},
		{
			name:     "tilde with path",
			input:    "~/projects",
			expected: filepath.Join(home, "projects"),
		},
		{
			name:     "no tilde",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "empty path",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExpandPath(tt.input)
			if err != nil {
				t.Errorf("ExpandPath failed: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestExpandPaths tests batch path expansion.
func TestExpandPaths(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	input := []string{"~/projects", "/absolute/path", "~/code"}
	expected := []string{
		filepath.Join(home, "projects"),
		"/absolute/path",
		filepath.Join(home, "code"),
	}

	result, err := ExpandPaths(input)
	if err != nil {
		t.Fatalf("ExpandPaths failed: %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d paths, got %d", len(expected), len(result))
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("path %d: expected %q, got %q", i, expected[i], result[i])
		}
	}
}

// TestSaveAndLoad tests saving and loading configuration.
func TestSaveAndLoad(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "codequest-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override the config path for testing
	testConfigPath := filepath.Join(tmpDir, "config.toml")

	// Create a test config
	testCfg := DefaultConfig()
	testCfg.Character.Name = "TestHero"
	testCfg.Game.Difficulty = "hard"
	testCfg.UI.Theme = "light"

	// Save to the test path
	testConfigDir := filepath.Dir(testConfigPath)
	if err := os.MkdirAll(testConfigDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	f, err := os.Create(testConfigPath)
	if err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(testCfg); err != nil {
		f.Close()
		t.Fatalf("failed to encode config: %v", err)
	}
	f.Close()

	// Load from the test path
	loadedCfg := &Config{}
	if _, err := toml.DecodeFile(testConfigPath, loadedCfg); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Verify loaded values match saved values
	if loadedCfg.Character.Name != "TestHero" {
		t.Errorf("expected character name 'TestHero', got %q", loadedCfg.Character.Name)
	}
	if loadedCfg.Game.Difficulty != "hard" {
		t.Errorf("expected difficulty 'hard', got %q", loadedCfg.Game.Difficulty)
	}
	if loadedCfg.UI.Theme != "light" {
		t.Errorf("expected theme 'light', got %q", loadedCfg.UI.Theme)
	}
}

// TestConfigPath tests that ConfigPath returns a valid path.
func TestConfigPath(t *testing.T) {
	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath failed: %v", err)
	}

	if path == "" {
		t.Error("ConfigPath returned empty string")
	}

	// Should end with config.toml
	if filepath.Base(path) != "config.toml" {
		t.Errorf("expected path to end with 'config.toml', got %q", path)
	}

	// Should contain .config/codequest
	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path, got %q", path)
	}
}

// TestSaveMethod tests the Save method on Config.
func TestSaveMethod(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "codequest-save-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Temporarily override HOME to use temp directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create a test config
	testCfg := DefaultConfig()
	testCfg.Character.Name = "SaveTestHero"
	testCfg.Game.Difficulty = "easy"

	// Save the config
	if err := testCfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify the file was created
	configPath := filepath.Join(tmpDir, ".config", "codequest", "config.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("config file was not created at %s", configPath)
	}

	// Load it back and verify
	loadedCfg := &Config{}
	if _, err := toml.DecodeFile(configPath, loadedCfg); err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}

	if loadedCfg.Character.Name != "SaveTestHero" {
		t.Errorf("expected character name 'SaveTestHero', got %q", loadedCfg.Character.Name)
	}
	if loadedCfg.Game.Difficulty != "easy" {
		t.Errorf("expected difficulty 'easy', got %q", loadedCfg.Game.Difficulty)
	}
}

// TestLoadMethod tests the Load method with various scenarios.
func TestLoadMethod(t *testing.T) {
	t.Run("creates config with defaults when missing", func(t *testing.T) {
		// Create a temporary directory for the test
		tmpDir, err := os.MkdirTemp("", "codequest-load-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Temporarily override HOME to use temp directory
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpDir)
		defer os.Setenv("HOME", originalHome)

		// Load config (file doesn't exist)
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		// Should have default values
		if cfg.Character.Name != "CodeWarrior" {
			t.Errorf("expected default character name 'CodeWarrior', got %q", cfg.Character.Name)
		}

		// Config file should now exist
		configPath := filepath.Join(tmpDir, ".config", "codequest", "config.toml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("config file was not created")
		}
	})

	t.Run("loads existing config", func(t *testing.T) {
		// Create a temporary directory for the test
		tmpDir, err := os.MkdirTemp("", "codequest-load-existing-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Temporarily override HOME to use temp directory
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpDir)
		defer os.Setenv("HOME", originalHome)

		// Create a config file first
		testCfg := DefaultConfig()
		testCfg.Character.Name = "ExistingHero"
		testCfg.Game.Difficulty = "hard"

		if err := testCfg.Save(); err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		// Load the config
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		// Should have the saved values
		if cfg.Character.Name != "ExistingHero" {
			t.Errorf("expected character name 'ExistingHero', got %q", cfg.Character.Name)
		}
		if cfg.Game.Difficulty != "hard" {
			t.Errorf("expected difficulty 'hard', got %q", cfg.Game.Difficulty)
		}
	})
}

// TestValidationErrorMessage tests the Error method on ValidationError.
func TestValidationErrorMessage(t *testing.T) {
	verr := ValidationError{
		Field:   "test.field",
		Value:   "invalid",
		Message: "must be valid",
	}

	errMsg := verr.Error()
	expectedMsg := "config validation error [test.field]: must be valid (value: invalid)"

	if errMsg != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, errMsg)
	}
}
