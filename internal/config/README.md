# Configuration Package

The `config` package provides comprehensive configuration management for CodeQuest, including loading, saving, validating, and managing application settings from TOML files.

## Overview

Configuration is stored in `~/.config/codequest/config.toml` and supports multiple sections for different aspects of the application:

- **character**: Player character settings
- **game**: Game mechanics and behavior
- **ui**: User interface preferences
- **tracking**: Activity tracking configuration
- **ai**: AI provider settings (mentor and code review)
- **git**: Git integration settings
- **github**: GitHub integration settings
- **keybinds**: Keyboard shortcut mappings
- **debug**: Debugging and logging configuration

## Usage

### Loading Configuration

```go
package main

import (
    "log"
    "github.com/AutumnsGrove/codequest/internal/config"
)

func main() {
    // Load config from ~/.config/codequest/config.toml
    // If the file doesn't exist, it creates it with default values
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Validate the configuration
    if err := cfg.Validate(); err != nil {
        log.Fatalf("Invalid config: %v", err)
    }

    // Use configuration values
    log.Printf("Character: %s", cfg.Character.Name)
    log.Printf("Difficulty: %s", cfg.Game.Difficulty)
    log.Printf("Theme: %s", cfg.UI.Theme)
}
```

### Modifying and Saving Configuration

```go
// Modify configuration
cfg.Character.Name = "CodeNinja"
cfg.Game.Difficulty = "hard"
cfg.UI.Theme = "light"

// Validate before saving
if err := cfg.Validate(); err != nil {
    log.Fatalf("Invalid config: %v", err)
}

// Save to disk
if err := cfg.Save(); err != nil {
    log.Fatalf("Failed to save config: %v", err)
}
```

### Using Default Configuration

```go
// Get a config with default values (doesn't save to disk)
cfg := config.DefaultConfig()

// Modify as needed
cfg.Character.Name = "MyHero"

// Save when ready
if err := cfg.Save(); err != nil {
    log.Fatalf("Failed to save config: %v", err)
}
```

### Path Expansion

The package provides utilities to expand `~` in paths to the user's home directory:

```go
// Expand a single path
expandedPath, err := config.ExpandPath("~/projects/myapp")
if err != nil {
    log.Fatal(err)
}
// expandedPath = "/home/user/projects/myapp"

// Expand multiple paths
paths := []string{"~/projects", "~/code", "/absolute/path"}
expandedPaths, err := config.ExpandPaths(paths)
if err != nil {
    log.Fatal(err)
}
```

## Configuration Structure

### Example TOML File

```toml
[character]
name = "CodeWarrior"

[game]
auto_start_quests = false
show_tips = true
difficulty = "normal"  # Options: easy, normal, hard

[ui]
theme = "dark"  # Options: dark, light, auto
show_animations = true
compact_mode = false
show_keybind_hints = true

[tracking]
session_timer_enabled = true
session_hotkey = "ctrl+t"
wakatime_enabled = false

[ai.mentor]
provider = "crush"  # Options: crush, mods, claude-code
model_complex = "openrouter/kimi/k2-0925"
model_simple = "openrouter/deepseek/glm-4.5-air"
model_complex_offline = "qwen3:30b"
model_simple_offline = "qwen3:4b"
temperature = 0.7

[ai.review]
provider = "mods"  # Options: crush, mods, claude-code
model_primary = "qwen3:30b"
model_fallback = "qwen3:4b"
auto_review = true
bonus_xp_enabled = true

[git]
auto_detect_repos = true
watch_paths = ["~/projects"]

[github]
enabled = false

[keybinds]
dashboard_quests = "q"
dashboard_character = "c"
global_timer = "ctrl+t"

[debug]
enabled = false
log_level = "info"  # Options: debug, info, warn, error
log_file = ""  # Empty means no file logging
```

## Validation

The `Validate()` method checks all configuration values for validity:

- **game.difficulty**: Must be "easy", "normal", or "hard"
- **ui.theme**: Must be "dark", "light", or "auto"
- **ai.mentor.provider**: Must be "crush", "mods", or "claude-code"
- **ai.review.provider**: Must be "crush", "mods", or "claude-code"
- **ai.mentor.temperature**: Must be between 0 and 2
- **debug.log_level**: Must be "debug", "info", "warn", or "error"
- **character.name**: Must not be empty

If validation fails, it returns a `ValidationError` with details about the invalid field:

```go
err := cfg.Validate()
if err != nil {
    if verr, ok := err.(config.ValidationError); ok {
        log.Printf("Invalid field: %s", verr.Field)
        log.Printf("Invalid value: %v", verr.Value)
        log.Printf("Error message: %s", verr.Message)
    }
}
```

## Default Values

The package provides sensible defaults for all configuration options:

- Character name: "CodeWarrior"
- Game difficulty: "normal"
- UI theme: "dark"
- AI mentor provider: "crush"
- Temperature: 0.7
- Session timer: enabled
- Auto-start quests: disabled
- Show tips: enabled
- Animations: enabled

See `defaults.go` for the complete list of default values.

## Testing

The package includes comprehensive tests with >80% coverage:

```bash
# Run tests
go test ./internal/config

# Run tests with coverage
go test ./internal/config -cover

# View detailed coverage
go test ./internal/config -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Notes

- API keys are NOT stored in the config file. Use the storage layer (Skate) for secure key storage.
- The config directory is created automatically if it doesn't exist.
- All paths in the config file can use `~` to refer to the home directory.
- Configuration is validated separately from loading - always call `Validate()` after `Load()`.
- The package uses `github.com/BurntSushi/toml` for TOML parsing.

## Error Handling

All functions return descriptive errors that can be wrapped with additional context:

```go
cfg, err := config.Load()
if err != nil {
    return fmt.Errorf("loading config: %w", err)
}

if err := cfg.Validate(); err != nil {
    return fmt.Errorf("validating config: %w", err)
}
```
