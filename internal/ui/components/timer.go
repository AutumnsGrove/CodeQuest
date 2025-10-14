// Package components provides reusable UI components for CodeQuest.
// This file implements the timer component for displaying coding session time.
package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/AutumnsGrove/codequest/internal/ui"
)

// Timer display modes for different contexts
const (
	TimerModeInline  = "inline"  // Compact badge for headers (e.g., "⏱ 2:34:15")
	TimerModeCard    = "card"    // Full card with hints and decoration
	TimerModeMinimal = "minimal" // Just the time (e.g., "2:34:15")
)

// TimerConfig holds configuration options for timer rendering
type TimerConfig struct {
	Width     int    // Total width for card mode (ignored for inline/minimal)
	Mode      string // Display mode: "inline", "card", or "minimal"
	IsRunning bool   // Whether the timer is actively running
}

// DefaultTimerConfig returns sensible defaults for timer configuration
func DefaultTimerConfig() TimerConfig {
	return TimerConfig{
		Width:     30,
		Mode:      TimerModeInline,
		IsRunning: false,
	}
}

// ============================================================================
// Public API - Main Rendering Functions
// ============================================================================

// RenderTimer renders the timer in the default (inline) mode.
// This is a convenience function for the most common use case.
//
// Parameters:
//   - duration: The elapsed session time (nil-safe, zero duration shows "00:00:00")
//   - isRunning: Whether the timer is actively running
//
// Returns:
//   - string: The rendered timer as a styled string
func RenderTimer(duration time.Duration, isRunning bool) string {
	config := DefaultTimerConfig()
	config.IsRunning = isRunning
	return RenderTimerWithConfig(duration, config)
}

// RenderTimerWithConfig renders the timer with custom configuration.
// This allows for full control over display mode and appearance.
//
// Parameters:
//   - duration: The elapsed session time (zero duration shows "00:00:00")
//   - config: Configuration options for rendering
//
// Returns:
//   - string: The rendered timer in the specified mode
func RenderTimerWithConfig(duration time.Duration, config TimerConfig) string {
	// Handle negative duration (shouldn't happen, but be defensive)
	if duration < 0 {
		duration = 0
	}

	switch config.Mode {
	case TimerModeCard:
		return renderTimerCard(duration, config.IsRunning, config.Width)
	case TimerModeMinimal:
		return renderTimerMinimal(duration)
	case TimerModeInline:
		fallthrough
	default:
		return renderTimerInline(duration, config.IsRunning)
	}
}

// RenderInlineTimer renders a compact inline timer badge.
// This is useful for dashboard headers or status bars.
//
// Example output: "⏱ 2:34:15" or "🔴 2:34:15" (when running)
//
// Parameters:
//   - duration: The elapsed session time
//   - isRunning: Whether the timer is actively running
//
// Returns:
//   - string: A compact one-line timer badge
func RenderInlineTimer(duration time.Duration, isRunning bool) string {
	return renderTimerInline(duration, isRunning)
}

// RenderTimerCard renders a full timer card with decoration and hints.
// This is suitable for dedicated timer widgets or larger displays.
//
// Parameters:
//   - duration: The elapsed session time
//   - isRunning: Whether the timer is actively running
//   - width: The total width for the card (minimum 25 characters)
//
// Returns:
//   - string: A full card with border, time, and keyboard hints
func RenderTimerCard(duration time.Duration, isRunning bool, width int) string {
	return renderTimerCard(duration, isRunning, width)
}

// RenderMinimalTimer renders just the time without any icons or decoration.
// This is useful when you need only the raw formatted time display.
//
// Example output: "2:34:15" or "0:45:30"
//
// Parameters:
//   - duration: The elapsed session time
//
// Returns:
//   - string: The formatted time string (HH:MM:SS or H:MM:SS)
func RenderMinimalTimer(duration time.Duration) string {
	return renderTimerMinimal(duration)
}

// ============================================================================
// Internal Rendering Functions
// ============================================================================

// renderTimerInline creates a compact inline timer badge with icon.
func renderTimerInline(duration time.Duration, isRunning bool) string {
	timeStr := formatDuration(duration)
	icon := getTimerIcon(isRunning)
	color := getTimerColor(duration)

	// Style the icon and time together
	style := lipgloss.NewStyle().
		Foreground(color).
		Bold(true)

	return style.Render(icon + " " + timeStr)
}

// renderTimerCard creates a full card display with border and hints.
func renderTimerCard(duration time.Duration, isRunning bool, width int) string {
	// Enforce minimum width
	if width < 25 {
		width = 25
	}

	timeStr := formatDuration(duration)
	icon := getTimerIcon(isRunning)
	color := getTimerColor(duration)

	// Create the title
	titleStyle := lipgloss.NewStyle().
		Foreground(ui.ColorAccent).
		Bold(true)
	title := titleStyle.Render(icon + " Session Timer")

	// Create the large time display
	timeStyle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Align(lipgloss.Center)
	timeDisplay := timeStyle.Render(timeStr)

	// Create the hint based on state
	hintStyle := lipgloss.NewStyle().
		Foreground(ui.ColorDim).
		Italic(true).
		Align(lipgloss.Center)

	var hint string
	if isRunning {
		hint = hintStyle.Render("[Ctrl+T] Pause")
	} else {
		hint = hintStyle.Render("[Ctrl+T] Resume")
	}

	// Add optional break reminder for long sessions
	var breakReminder string
	if duration >= 5*time.Hour {
		breakStyle := lipgloss.NewStyle().
			Foreground(ui.ColorWarning).
			Bold(true).
			Align(lipgloss.Center)
		breakReminder = "\n" + breakStyle.Render("⚠ Take a break!")
	}

	// Combine all sections
	content := title + "\n\n" + timeDisplay + "\n\n" + hint + breakReminder

	// Wrap in a box with appropriate border color
	borderColor := ui.ColorAccent
	if isRunning {
		borderColor = ui.ColorSuccess // Green border when running
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width - 4). // Account for padding
		Padding(1, 2)

	return boxStyle.Render(content)
}

// renderTimerMinimal renders just the formatted time without decoration.
func renderTimerMinimal(duration time.Duration) string {
	return formatDuration(duration)
}

// ============================================================================
// Helper Functions
// ============================================================================

// formatDuration converts a time.Duration to HH:MM:SS or H:MM:SS format.
// This function handles durations up to 99 hours and 59 minutes.
//
// Examples:
//   - 0 seconds -> "0:00:00"
//   - 45 minutes 30 seconds -> "0:45:30"
//   - 2 hours 34 minutes 15 seconds -> "2:34:15"
//   - 10 hours 5 minutes 8 seconds -> "10:05:08"
//
// Parameters:
//   - d: The duration to format (negative durations are treated as zero)
//
// Returns:
//   - string: The formatted time string
func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	// Format as H:MM:SS (no leading zero for hours unless >= 10)
	return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
}

// getTimerColor returns the appropriate color based on session duration.
// This provides visual feedback about session length.
//
// Color coding:
//   - 0-1 hour: Dim (just getting started)
//   - 1-3 hours: Cyan/Info (good session)
//   - 3-5 hours: Green/Success (productive session)
//   - 5+ hours: Orange/Warning (time for a break!)
//
// Parameters:
//   - d: The session duration
//
// Returns:
//   - lipgloss.Color: The appropriate color for the duration
func getTimerColor(d time.Duration) lipgloss.Color {
	hours := d.Hours()

	switch {
	case hours >= 5.0:
		return ui.ColorWarning // 5+ hours: Orange (take a break!)
	case hours >= 3.0:
		return ui.ColorSuccess // 3-5 hours: Green (great session)
	case hours >= 1.0:
		return ui.ColorInfo // 1-3 hours: Cyan (good work)
	default:
		return ui.ColorDim // 0-1 hour: Gray (warming up)
	}
}

// getTimerIcon returns the appropriate icon based on timer state.
// This provides immediate visual feedback about whether the timer is running.
//
// Icons:
//   - Running: 🔴 (red dot, implies active/recording)
//   - Paused: ⏸ (pause symbol)
//
// Parameters:
//   - isRunning: Whether the timer is actively running
//
// Returns:
//   - string: The appropriate emoji icon
func getTimerIcon(isRunning bool) string {
	if isRunning {
		return "🔴" // Red dot for active/running
	}
	return "⏸" // Pause symbol for paused
}

// ============================================================================
// Additional Utility Functions
// ============================================================================

// GetTimerHint returns a helpful hint string based on timer state.
// This can be used to provide contextual help to users.
//
// Parameters:
//   - isRunning: Whether the timer is actively running
//
// Returns:
//   - string: A hint message about how to control the timer
func GetTimerHint(isRunning bool) string {
	if isRunning {
		return "Press Ctrl+T to pause the timer"
	}
	return "Press Ctrl+T to start the timer"
}

// ShouldShowBreakReminder returns true if the session duration suggests a break.
// This can be used to trigger break reminder notifications.
//
// Parameters:
//   - duration: The current session duration
//
// Returns:
//   - bool: true if duration >= 5 hours (time for a break)
func ShouldShowBreakReminder(duration time.Duration) bool {
	return duration >= 5*time.Hour
}

// FormatSessionSummary creates a summary string of the session time.
// This is useful for notifications or completion messages.
//
// Example: "Great session! You coded for 2 hours and 34 minutes."
//
// Parameters:
//   - duration: The total session duration
//
// Returns:
//   - string: A human-friendly summary message
func FormatSessionSummary(duration time.Duration) string {
	if duration < time.Minute {
		return "Quick session!"
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	if hours == 0 {
		return fmt.Sprintf("You coded for %d minutes.", minutes)
	}

	if minutes == 0 {
		if hours == 1 {
			return "You coded for 1 hour."
		}
		return fmt.Sprintf("You coded for %d hours.", hours)
	}

	// Both hours and minutes
	if hours == 1 {
		if minutes == 1 {
			return "You coded for 1 hour and 1 minute."
		}
		return fmt.Sprintf("You coded for 1 hour and %d minutes.", minutes)
	}

	if minutes == 1 {
		return fmt.Sprintf("You coded for %d hours and 1 minute.", hours)
	}

	return fmt.Sprintf("You coded for %d hours and %d minutes.", hours, minutes)
}
