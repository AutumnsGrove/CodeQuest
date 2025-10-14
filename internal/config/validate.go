package config

import (
	"fmt"
	"strings"
)

// ValidationError represents a configuration validation error.
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("config validation error [%s]: %s (value: %v)", e.Field, e.Message, e.Value)
}

// Validate checks that all configuration values are valid.
// It returns a ValidationError if any value is invalid, or nil if all values are valid.
func (c *Config) Validate() error {
	// Validate Game.Difficulty
	validDifficulties := []string{"easy", "normal", "hard"}
	if !contains(validDifficulties, c.Game.Difficulty) {
		return ValidationError{
			Field:   "game.difficulty",
			Value:   c.Game.Difficulty,
			Message: fmt.Sprintf("must be one of: %s", strings.Join(validDifficulties, ", ")),
		}
	}

	// Validate UI.Theme
	validThemes := []string{"dark", "light", "auto"}
	if !contains(validThemes, c.UI.Theme) {
		return ValidationError{
			Field:   "ui.theme",
			Value:   c.UI.Theme,
			Message: fmt.Sprintf("must be one of: %s", strings.Join(validThemes, ", ")),
		}
	}

	// Validate AI.Mentor.Provider
	validAIProviders := []string{"crush", "mods", "claude-code"}
	if !contains(validAIProviders, c.AI.Mentor.Provider) {
		return ValidationError{
			Field:   "ai.mentor.provider",
			Value:   c.AI.Mentor.Provider,
			Message: fmt.Sprintf("must be one of: %s", strings.Join(validAIProviders, ", ")),
		}
	}

	// Validate AI.Review.Provider
	if !contains(validAIProviders, c.AI.Review.Provider) {
		return ValidationError{
			Field:   "ai.review.provider",
			Value:   c.AI.Review.Provider,
			Message: fmt.Sprintf("must be one of: %s", strings.Join(validAIProviders, ", ")),
		}
	}

	// Validate AI.Mentor.Temperature (should be between 0 and 2)
	if c.AI.Mentor.Temperature < 0 || c.AI.Mentor.Temperature > 2 {
		return ValidationError{
			Field:   "ai.mentor.temperature",
			Value:   c.AI.Mentor.Temperature,
			Message: "must be between 0 and 2",
		}
	}

	// Validate Debug.LogLevel
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, c.Debug.LogLevel) {
		return ValidationError{
			Field:   "debug.log_level",
			Value:   c.Debug.LogLevel,
			Message: fmt.Sprintf("must be one of: %s", strings.Join(validLogLevels, ", ")),
		}
	}

	// Validate Character.Name (must not be empty)
	if strings.TrimSpace(c.Character.Name) == "" {
		return ValidationError{
			Field:   "character.name",
			Value:   c.Character.Name,
			Message: "must not be empty",
		}
	}

	return nil
}

// contains checks if a slice contains a specific string.
// This is a helper function for validation.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
