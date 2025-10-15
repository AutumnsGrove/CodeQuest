package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the complete application configuration.
// It is loaded from ~/.config/codequest/config.toml and validated on startup.
type Config struct {
	Character CharacterConfig `toml:"character"`
	Game      GameConfig      `toml:"game"`
	UI        UIConfig        `toml:"ui"`
	Tracking  TrackingConfig  `toml:"tracking"`
	AI        AIConfig        `toml:"ai"`
	Git       GitConfig       `toml:"git"`
	Github    GithubConfig    `toml:"github"`
	Keybinds  KeybindsConfig  `toml:"keybinds"`
	Debug     DebugConfig     `toml:"debug"`
}

// CharacterConfig contains character-specific settings.
type CharacterConfig struct {
	Name string `toml:"name"`
}

// GameConfig contains game mechanics settings.
type GameConfig struct {
	AutoStartQuests bool   `toml:"auto_start_quests"`
	ShowTips        bool   `toml:"show_tips"`
	Difficulty      string `toml:"difficulty"` // easy, normal, hard
}

// UIConfig contains user interface preferences.
type UIConfig struct {
	Theme            string `toml:"theme"` // dark, light, auto
	ShowAnimations   bool   `toml:"show_animations"`
	CompactMode      bool   `toml:"compact_mode"`
	ShowKeybindHints bool   `toml:"show_keybind_hints"`
}

// TrackingConfig contains activity tracking settings.
type TrackingConfig struct {
	SessionTimerEnabled bool   `toml:"session_timer_enabled"`
	SessionHotkey       string `toml:"session_hotkey"`
	WakatimeEnabled     bool   `toml:"wakatime_enabled"`
}

// AIConfig contains all AI-related configuration.
type AIConfig struct {
	Mentor AIMentorConfig `toml:"mentor"`
	Review AIReviewConfig `toml:"review"`
}

// AIMentorConfig contains AI mentor (Crush) settings.
type AIMentorConfig struct {
	Provider            string  `toml:"provider"`
	ModelComplex        string  `toml:"model_complex"`
	ModelSimple         string  `toml:"model_simple"`
	ModelComplexOffline string  `toml:"model_complex_offline"`
	ModelSimpleOffline  string  `toml:"model_simple_offline"`
	Temperature         float64 `toml:"temperature"`
}

// AIReviewConfig contains AI code review (Mods) settings.
type AIReviewConfig struct {
	Provider       string `toml:"provider"`
	ModelPrimary   string `toml:"model_primary"`
	ModelFallback  string `toml:"model_fallback"`
	AutoReview     bool   `toml:"auto_review"`
	BonusXPEnabled bool   `toml:"bonus_xp_enabled"`
}

// GitConfig contains Git integration settings.
type GitConfig struct {
	AutoDetectRepos bool     `toml:"auto_detect_repos"`
	WatchPaths      []string `toml:"watch_paths"`
}

// GithubConfig contains GitHub integration settings.
type GithubConfig struct {
	Enabled bool `toml:"enabled"`
}

// KeybindsConfig contains keyboard shortcut mappings.
type KeybindsConfig struct {
	DashboardQuests    string `toml:"dashboard_quests"`
	DashboardCharacter string `toml:"dashboard_character"`
	GlobalTimer        string `toml:"global_timer"`
}

// DebugConfig contains debugging and logging settings.
type DebugConfig struct {
	Enabled  bool   `toml:"enabled"`
	LogLevel string `toml:"log_level"` // debug, info, warn, error
	LogFile  string `toml:"log_file"`  // empty means no file logging
}

// ConfigPath returns the full path to the config file.
// It expands ~ to the user's home directory.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, ".config", "codequest", "config.toml"), nil
}

// Load reads the config file from the standard location.
// If the file doesn't exist, it creates it with default values.
// The returned Config is not validated - call Validate() separately.
func Load() (*Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, fmt.Errorf("determining config path: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config doesn't exist, create it with defaults
		cfg := DefaultConfig()
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("creating default config: %w", err)
		}
		return cfg, nil
	}

	// Config exists, load it
	cfg := &Config{}
	if _, err := toml.DecodeFile(configPath, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return cfg, nil
}

// Save writes the current config to the standard config file location.
// It creates the config directory if it doesn't exist.
func (c *Config) Save() error {
	configPath, err := ConfigPath()
	if err != nil {
		return fmt.Errorf("determining config path: %w", err)
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	// Create or truncate the config file
	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("creating config file: %w", err)
	}
	defer f.Close()

	// Encode the config to TOML
	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("encoding config to TOML: %w", err)
	}

	return nil
}

// ExpandPath expands ~ in a path to the user's home directory.
// If the path doesn't start with ~, it returns the path unchanged.
func ExpandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	if len(path) == 1 {
		return home, nil
	}

	// Replace ~ with home directory
	return filepath.Join(home, path[1:]), nil
}

// ExpandPaths expands ~ in all paths in a slice.
func ExpandPaths(paths []string) ([]string, error) {
	expanded := make([]string, len(paths))
	for i, path := range paths {
		exp, err := ExpandPath(path)
		if err != nil {
			return nil, fmt.Errorf("expanding path %q: %w", path, err)
		}
		expanded[i] = exp
	}
	return expanded, nil
}
