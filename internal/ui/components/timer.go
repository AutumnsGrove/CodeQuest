// Package components provides reusable UI components for CodeQuest.
// This file implements the timer component for displaying coding session time.
package components

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AutumnsGrove/codequest/internal/watcher"
)

// Color constants (duplicated here to avoid import cycle with ui package)
var (
	timerColorSuccess = lipgloss.Color("42")  // Green
	timerColorWarning = lipgloss.Color("214") // Orange
	timerColorInfo    = lipgloss.Color("69")  // Blue
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
		Foreground(colorAccent).
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
		Foreground(colorDim).
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
			Foreground(timerColorWarning).
			Bold(true).
			Align(lipgloss.Center)
		breakReminder = "\n" + breakStyle.Render("⚠ Take a break!")
	}

	// Combine all sections
	content := title + "\n\n" + timeDisplay + "\n\n" + hint + breakReminder

	// Wrap in a box with appropriate border color
	borderColor := colorAccent
	if isRunning {
		borderColor = timerColorSuccess // Green border when running
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width-4). // Account for padding
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
		return timerColorWarning // 5+ hours: Orange (take a break!)
	case hours >= 3.0:
		return timerColorSuccess // 3-5 hours: Green (great session)
	case hours >= 1.0:
		return timerColorInfo // 1-3 hours: Cyan (good work)
	default:
		return colorDim // 0-1 hour: Gray (warming up) - use from header.go
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

// ============================================================================
// Bubble Tea Timer Component
// ============================================================================

// Timer is a Bubble Tea component that displays session time and integrates
// with SessionTracker for live updates.
type Timer struct {
	tracker   *watcher.SessionTracker // SessionTracker for elapsed time
	elapsed   time.Duration           // Cached elapsed time
	isRunning bool                    // Whether timer is currently running
	width     int                     // Available width for display
}

// timerTickMsg is sent every second to update the timer display.
type timerTickMsg time.Time

// TimerTickMsg is the exported version of timerTickMsg for external use.
type TimerTickMsg time.Time

// NewTimer creates a new Timer component connected to SessionTracker.
//
// Parameters:
//   - tracker: The SessionTracker to use for time data (required)
//
// Returns:
//   - Timer: Initialized timer component
func NewTimer(tracker *watcher.SessionTracker) Timer {
	return Timer{
		tracker:   tracker,
		elapsed:   tracker.GetElapsed(),
		isRunning: tracker.GetState() == watcher.SessionRunning,
		width:     80, // Default width
	}
}

// Update handles Bubble Tea messages for the timer component.
// Responds to timerTickMsg to update elapsed time display.
//
// Parameters:
//   - msg: The Bubble Tea message
//
// Returns:
//   - Timer: Updated timer component
//   - tea.Cmd: Command to schedule next tick (if running)
func (t Timer) Update(msg tea.Msg) (Timer, tea.Cmd) {
	switch msg.(type) {
	case timerTickMsg, TimerTickMsg:
		// Update elapsed time from tracker
		t.elapsed = t.tracker.GetElapsed()
		t.isRunning = t.tracker.GetState() == watcher.SessionRunning

		// Request next tick if still running
		if t.isRunning {
			return t, timerTick()
		}
	}

	return t, nil
}

// View renders the timer in inline mode (suitable for footers/headers).
// Uses the RenderInlineTimer function for consistent styling.
//
// Returns:
//   - string: Rendered timer display
func (t Timer) View() string {
	return RenderInlineTimer(t.elapsed, t.isRunning)
}

// SetWidth updates the available width for the timer display.
// This is used for responsive layout adjustments.
//
// Parameters:
//   - width: Available width in characters
func (t *Timer) SetWidth(width int) {
	t.width = width
}

// GetElapsed returns the current elapsed time.
// This is useful for displaying formatted time elsewhere.
//
// Returns:
//   - time.Duration: Current elapsed session time
func (t Timer) GetElapsed() time.Duration {
	return t.elapsed
}

// IsRunning returns whether the timer is currently running.
//
// Returns:
//   - bool: true if timer is actively running
func (t Timer) IsRunning() bool {
	return t.isRunning
}

// timerTick returns a Bubble Tea command that sends a TimerTickMsg
// after 1 second. This creates the timer update loop.
//
// Returns:
//   - tea.Cmd: Command to schedule next tick
func timerTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TimerTickMsg(t)
	})
}

// TimerTick is the exported version of timerTick for starting the timer
// from the main app initialization.
//
// Returns:
//   - tea.Cmd: Command to schedule first tick
func TimerTick() tea.Cmd {
	return timerTick()
}
