package config

// DefaultConfig returns a Config struct populated with sensible default values.
// These defaults are used when creating a new config file or when specific
// values are not provided in an existing config file.
func DefaultConfig() *Config {
	return &Config{
		Character: CharacterConfig{
			Name: "CodeWarrior",
		},
		Game: GameConfig{
			AutoStartQuests: false,
			ShowTips:        true,
			Difficulty:      "normal", // easy, normal, hard
		},
		UI: UIConfig{
			Theme:            "dark", // dark, light, auto
			ShowAnimations:   true,
			CompactMode:      false,
			ShowKeybindHints: true,
		},
		Tracking: TrackingConfig{
			SessionTimerEnabled: true,
			SessionHotkey:       "ctrl+t",
			WakatimeEnabled:     false,
		},
		AI: AIConfig{
			Mentor: AIMentorConfig{
				Provider:            "crush",
				ModelComplex:        "openrouter/kimi/k2-0925",
				ModelSimple:         "openrouter/deepseek/glm-4.5-air",
				ModelComplexOffline: "qwen3:30b",
				ModelSimpleOffline:  "qwen3:4b",
				Temperature:         0.7,
			},
			Review: AIReviewConfig{
				Provider:       "mods",
				ModelPrimary:   "qwen3:30b",
				ModelFallback:  "qwen3:4b",
				AutoReview:     true,
				BonusXPEnabled: true,
			},
		},
		Git: GitConfig{
			AutoDetectRepos: true,
			WatchPaths:      []string{"~/projects"},
		},
		Github: GithubConfig{
			Enabled: false,
		},
		Keybinds: KeybindsConfig{
			DashboardQuests:    "q",
			DashboardCharacter: "c",
			GlobalTimer:        "ctrl+t",
		},
		Debug: DebugConfig{
			Enabled:  false,
			LogLevel: "info", // debug, info, warn, error
			LogFile:  "",     // empty means no file logging
		},
	}
}
