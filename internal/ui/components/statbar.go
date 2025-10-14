// Package components provides reusable UI components for CodeQuest
// This file implements the stat bar component that displays character statistics
package components

import (
	"fmt"
	"strings"

	"github.com/AutumnsGrove/codequest/internal/game"
	"github.com/AutumnsGrove/codequest/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

// StatBarConfig holds configuration options for the stat bar rendering
type StatBarConfig struct {
	Width           int  // Total width available for the stat bar
	ShowSessionInfo bool // Whether to show session stats (commits/lines today)
	ShowStreak      bool // Whether to show streak counter
	Compact         bool // If true, use a more compact layout
}

// DefaultStatBarConfig returns sensible defaults for stat bar configuration
func DefaultStatBarConfig() StatBarConfig {
	return StatBarConfig{
		Width:           80,
		ShowSessionInfo: true,
		ShowStreak:      true,
		Compact:         false,
	}
}

// RenderStatBar renders the complete character stat bar with progress bars and stats.
// This is the main entry point for displaying character information in the UI.
//
// The stat bar includes:
//   - Character name and level with XP progress bar
//   - RPG stats (CodePower, Wisdom, Agility)
//   - Daily streak counter
//   - Session stats (commits and lines added today)
//
// Parameters:
//   - char: The character whose stats to display (nil-safe)
//   - width: The total width available for the stat bar (minimum 40 characters)
//
// Returns:
//   - string: The rendered stat bar as a formatted string
func RenderStatBar(char *game.Character, width int) string {
	// Handle nil character gracefully
	if char == nil {
		return renderNilCharacterStatBar(width)
	}

	config := DefaultStatBarConfig()
	config.Width = width
	return RenderStatBarWithConfig(char, config)
}

// RenderStatBarWithConfig renders the stat bar with custom configuration.
// This allows for more control over what information is displayed and how.
//
// Parameters:
//   - char: The character whose stats to display
//   - config: Configuration options for rendering
//
// Returns:
//   - string: The rendered stat bar as a formatted string
func RenderStatBarWithConfig(char *game.Character, config StatBarConfig) string {
	if config.Width < 40 {
		config.Width = 40 // Enforce minimum width
	}

	var sections []string

	// Section 1: Name, Level, and XP Bar
	headerSection := renderHeaderSection(char, config.Width)
	sections = append(sections, headerSection)

	// Section 2: RPG Stats (CodePower, Wisdom, Agility)
	statsSection := renderRPGStats(char, config.Width)
	sections = append(sections, statsSection)

	// Section 3: Streak Counter (if enabled)
	if config.ShowStreak {
		streakSection := renderStreakSection(char, config.Width)
		sections = append(sections, streakSection)
	}

	// Section 4: Session Stats (if enabled)
	if config.ShowSessionInfo {
		sessionSection := renderSessionStats(char, config.Width)
		sections = append(sections, sessionSection)
	}

	// Join all sections with newlines
	return strings.Join(sections, "\n")
}

// renderHeaderSection renders the character name, level, and XP progress bar.
// This is the top section of the stat bar showing core progression.
func renderHeaderSection(char *game.Character, width int) string {
	// Character name with level badge
	nameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.ColorPrimary)

	levelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ui.ColorLevel).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	nameAndLevel := nameStyle.Render(char.Name) + " " +
		levelStyle.Render(fmt.Sprintf("Lv.%d", char.Level))

	// XP Progress Bar
	// Calculate available width for the progress bar (reserve space for text)
	barWidth := width - 30 // Reserve space for XP text
	if barWidth < 20 {
		barWidth = 20 // Minimum bar width
	}

	xpBar := ui.RenderProgressBar(char.XP, char.XPToNextLevel, barWidth, "xp")
	xpLabel := ui.StatLabelStyle.Render("XP: ")

	// Combine name/level with XP bar
	line1 := nameAndLevel
	line2 := xpLabel + xpBar

	return line1 + "\n" + line2
}

// renderRPGStats renders the character's RPG statistics (CodePower, Wisdom, Agility).
// These stats are displayed in a horizontal layout with icons and colors.
func renderRPGStats(char *game.Character, width int) string {
	// Create styled stat displays with icons
	codePowerIcon := "⚔️"
	wisdomIcon := "📖"
	agilityIcon := "⚡"

	// Stat styles
	statNameStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMuted).
		Bold(false)

	statValueStyle := lipgloss.NewStyle().
		Foreground(ui.ColorPrimary).
		Bold(true)

	// Format each stat with icon, name, and value
	codePowerStat := fmt.Sprintf("%s %s %s",
		codePowerIcon,
		statNameStyle.Render("CodePower:"),
		statValueStyle.Render(fmt.Sprintf("%d", char.CodePower)))

	wisdomStat := fmt.Sprintf("%s %s %s",
		wisdomIcon,
		statNameStyle.Render("Wisdom:"),
		statValueStyle.Render(fmt.Sprintf("%d", char.Wisdom)))

	agilityStat := fmt.Sprintf("%s %s %s",
		agilityIcon,
		statNameStyle.Render("Agility:"),
		statValueStyle.Render(fmt.Sprintf("%d", char.Agility)))

	// Join stats horizontally with spacing
	statsLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		codePowerStat+"  ",
		wisdomStat+"  ",
		agilityStat,
	)

	return statsLine
}

// renderStreakSection renders the daily activity streak counter.
// Shows both current streak and longest streak achieved.
func renderStreakSection(char *game.Character, width int) string {
	streakIcon := "🔥"
	trophyIcon := "🏆"

	labelStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMuted)

	// Current streak with fire icon
	currentStreakStyle := lipgloss.NewStyle().
		Foreground(ui.ColorWarning). // Orange/gold color for fire theme
		Bold(true)

	currentStreakText := fmt.Sprintf("%s %s %s",
		streakIcon,
		labelStyle.Render("Streak:"),
		currentStreakStyle.Render(fmt.Sprintf("%d days", char.CurrentStreak)))

	// Longest streak with trophy icon
	longestStreakStyle := lipgloss.NewStyle().
		Foreground(ui.ColorXP).
		Bold(true)

	longestStreakText := fmt.Sprintf("%s %s %s",
		trophyIcon,
		labelStyle.Render("Best:"),
		longestStreakStyle.Render(fmt.Sprintf("%d days", char.LongestStreak)))

	// Join streaks horizontally
	streakLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		currentStreakText+"  ",
		longestStreakText,
	)

	return streakLine
}

// renderSessionStats renders today's session statistics.
// Shows commits and lines of code added today.
func renderSessionStats(char *game.Character, width int) string {
	commitIcon := "💾"
	codeIcon := "📝"

	labelStyle := lipgloss.NewStyle().
		Foreground(ui.ColorMuted)

	valueStyle := lipgloss.NewStyle().
		Foreground(ui.ColorInfo).
		Bold(true)

	// Today's commits
	commitsText := fmt.Sprintf("%s %s %s",
		commitIcon,
		labelStyle.Render("Today:"),
		valueStyle.Render(fmt.Sprintf("%d commits", char.TodayCommits)))

	// Today's lines added
	linesText := fmt.Sprintf("%s %s %s",
		codeIcon,
		labelStyle.Render("Lines:"),
		valueStyle.Render(fmt.Sprintf("+%d", char.TodayLinesAdded)))

	// Join session stats horizontally
	sessionLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		commitsText+"  ",
		linesText,
	)

	return sessionLine
}

// RenderCompactStatBar renders a more compact version of the stat bar.
// Useful for smaller terminal windows or when space is limited.
// This version shows only essential information: level, XP, and streak.
//
// Parameters:
//   - char: The character whose stats to display (nil-safe)
//   - width: The total width available for the stat bar
//
// Returns:
//   - string: The rendered compact stat bar
func RenderCompactStatBar(char *game.Character, width int) string {
	// Handle nil character gracefully
	if char == nil {
		return renderNilCharacterStatBar(width)
	}

	config := DefaultStatBarConfig()
	config.Width = width
	config.Compact = true
	config.ShowSessionInfo = false // Hide session info in compact mode

	// Just show name, level, XP bar, and streak
	headerSection := renderHeaderSection(char, width)
	streakSection := renderStreakSection(char, width)

	return headerSection + "\n" + streakSection
}

// RenderStatBadge renders a small inline stat badge.
// Useful for displaying in headers or alongside other content.
//
// Parameters:
//   - char: The character whose stats to display (nil-safe)
//
// Returns:
//   - string: A compact one-line stat badge
func RenderStatBadge(char *game.Character) string {
	// Handle nil character gracefully
	if char == nil {
		noCharStyle := lipgloss.NewStyle().
			Foreground(ui.ColorDim).
			Italic(true)
		return noCharStyle.Render("No character")
	}

	badgeStyle := lipgloss.NewStyle().
		Foreground(ui.ColorBright).
		Background(ui.ColorSecondary).
		Padding(0, 1).
		Bold(true)

	levelStyle := lipgloss.NewStyle().
		Foreground(ui.ColorLevel).
		Bold(true)

	streakStyle := lipgloss.NewStyle().
		Foreground(ui.ColorWarning).
		Bold(true)

	badge := fmt.Sprintf("%s | %s XP: %d/%d | %s %d🔥",
		levelStyle.Render(fmt.Sprintf("Lv.%d", char.Level)),
		ui.StatLabelStyle.Render(""),
		char.XP,
		char.XPToNextLevel,
		streakStyle.Render(""),
		char.CurrentStreak,
	)

	return badgeStyle.Render(badge)
}

// renderNilCharacterStatBar renders a placeholder message when character is nil.
// This prevents crashes and provides useful feedback to the user.
//
// Parameters:
//   - width: The available width for rendering
//
// Returns:
//   - string: A styled error message indicating no character is loaded
func renderNilCharacterStatBar(width int) string {
	errorStyle := lipgloss.NewStyle().
		Foreground(ui.ColorError).
		Bold(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(ui.ColorDim).
		Italic(true)

	message := errorStyle.Render("⚠️  No character loaded")
	hint := hintStyle.Render("Stats unavailable")

	// Center the message if width is large enough
	if width > 50 {
		content := message + "\n" + hint
		return lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(content)
	}

	// Otherwise just return the simple message
	return message
}
