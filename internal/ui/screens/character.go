// Package screens provides screen rendering functions for CodeQuest UI.
// This file implements the Character screen, displaying detailed character stats and history.
package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// RenderCharacter renders the complete character sheet screen.
// This is a detailed view of all character information, statistics, and history.
//
// The character screen shows:
//   - Header with character name and level
//   - Core stats (CodePower, Wisdom, Agility) with visual stat bars
//   - XP progress bar with detailed breakdown
//   - Streak information (current and longest)
//   - Lifetime statistics (commits, lines, quests)
//   - Session history (today's activity)
//   - Future: Achievements section (post-MVP)
//
// Layout Structure:
//   - Header: Screen title with character info
//   - Left panel: Stats and progression
//   - Right panel: History and achievements
//   - Footer: Key bindings
//
// Parameters:
//   - character: Player character to display (nil-safe)
//   - width: Terminal width in characters
//   - height: Terminal height in characters
//
// Returns:
//   - string: Rendered character screen UI
func RenderCharacter(character *game.Character, width, height int) string {
	// Handle nil character gracefully
	if character == nil {
		return renderNoCharacterScreen(width, height)
	}

	// Render header (inline to avoid import cycle)
	header := renderCharacterHeader(character, width)

	// Determine layout based on terminal width
	useWideLayout := width > 100

	var content string
	if useWideLayout {
		content = renderCharacterWide(character, width, height)
	} else {
		content = renderCharacterNarrow(character, width, height)
	}

	// Render footer with key bindings
	footer := renderCharacterFooter(width)

	// Assemble screen
	screen := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		"",
		footer,
	)

	return screen
}

// renderCharacterWide renders character screen with side-by-side panels for wide terminals.
func renderCharacterWide(character *game.Character, width, height int) string {
	// Split width into two columns (55% left, 45% right)
	leftWidth := int(float64(width) * 0.53)
	rightWidth := width - leftWidth - 2 // Account for spacing

	// Render left panel: Core stats and progression
	leftPanel := renderStatsPanel(character, leftWidth)

	// Render right panel: History and activity
	rightPanel := renderHistoryPanel(character, rightWidth)

	// Join panels horizontally
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)

	return panels
}

// renderCharacterNarrow renders character screen with stacked panels for narrow terminals.
func renderCharacterNarrow(character *game.Character, width, height int) string {
	// Full width for each panel
	panelWidth := width

	// Render panels vertically
	statsPanel := renderStatsPanel(character, panelWidth)
	historyPanel := renderHistoryPanel(character, panelWidth)

	// Stack all panels
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		statsPanel,
		"",
		historyPanel,
	)

	return content
}

// renderStatsPanel renders the left panel with character stats and progression.
func renderStatsPanel(character *game.Character, width int) string {
	sections := make([]string, 0)

	// Character Identity Section
	identitySection := renderIdentitySection(character)
	sections = append(sections, identitySection)

	// XP Progress Section
	xpSection := renderXPSection(character, width)
	sections = append(sections, xpSection)

	// Core Stats Section (using StatBar component)
	statsSection := renderCoreStatsSection(character, width)
	sections = append(sections, statsSection)

	// Streak Section
	streakSection := renderStreakSectionDetailed(character)
	sections = append(sections, streakSection)

	// Join all sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		sections...,
	)

	// Wrap in box
	return BoxStyle.Width(width - 4).Render(content)
}

// renderHistoryPanel renders the right panel with history and activity.
func renderHistoryPanel(character *game.Character, width int) string {
	sections := make([]string, 0)

	// Today's Activity Section
	todaySection := renderTodayActivityDetailed(character)
	sections = append(sections, todaySection)

	// Lifetime Statistics Section
	lifetimeSection := renderLifetimeStatsDetailed(character)
	sections = append(sections, lifetimeSection)

	// Achievements Section (placeholder)
	achievementsSection := renderAchievementsPlaceholder()
	sections = append(sections, achievementsSection)

	// Join all sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		sections...,
	)

	// Wrap in box
	return BoxStyle.Width(width - 4).Render(content)
}

// ============================================================================
// Section Rendering Functions
// ============================================================================

// renderIdentitySection renders character identity information.
func renderIdentitySection(character *game.Character) string {
	title := SubtitleStyle.Render("⚔️ Character")

	// Character name with emphasis
	nameLabel := StatLabelStyle.Render("Name: ")
	nameValue := BoldTextStyle.Render(character.Name)
	name := nameLabel + nameValue

	// Level display with badge
	levelLabel := StatLabelStyle.Render("Level: ")
	levelValue := lipgloss.NewStyle().
		Foreground(ColorLevel).
		Bold(true).
		Render(fmt.Sprintf("%d", character.Level))
	level := levelLabel + levelValue

	// Created date
	createdLabel := StatLabelStyle.Render("Created: ")
	createdValue := DimTextStyle.Render(formatDate(character.CreatedAt))
	created := createdLabel + createdValue

	// Days played
	daysSinceCreation := int(time.Since(character.CreatedAt).Hours() / 24)
	daysLabel := StatLabelStyle.Render("Days Played: ")
	daysValue := StatValueStyle.Render(fmt.Sprintf("%d", daysSinceCreation))
	days := daysLabel + daysValue

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		name,
		level,
		created,
		days,
	)
}

// renderXPSection renders XP progress with detailed breakdown.
func renderXPSection(character *game.Character, width int) string {
	title := SubtitleStyle.Render("📊 Experience")

	// XP progress bar
	xpLabel := StatLabelStyle.Render("XP: ")
	barWidth := width - 30
	if barWidth < 20 {
		barWidth = 20
	}
	xpBar := renderProgressBar(
		character.XP,
		character.XPToNextLevel,
		barWidth,
		"xp",
	)
	xpProgress := xpLabel + xpBar

	// Next level XP requirement
	remainingLabel := StatLabelStyle.Render("To Next Level: ")
	remainingValue := StatValueStyle.Render(
		fmt.Sprintf("%d XP", character.XPToNextLevel-character.XP),
	)
	remaining := remainingLabel + remainingValue

	// Percentage to next level
	percentage := 0.0
	if character.XPToNextLevel > 0 {
		percentage = (float64(character.XP) / float64(character.XPToNextLevel)) * 100
	}
	percentLabel := StatLabelStyle.Render("Progress: ")
	percentValue := lipgloss.NewStyle().
		Foreground(ColorXP).
		Bold(true).
		Render(fmt.Sprintf("%.1f%%", percentage))
	percent := percentLabel + percentValue

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		xpProgress,
		remaining,
		percent,
	)
}

// renderCoreStatsSection renders RPG stats inline (avoiding import cycle).
func renderCoreStatsSection(character *game.Character, width int) string {
	title := SubtitleStyle.Render("⚡ Core Stats")

	// Render stats inline instead of using StatBar component to avoid import cycle
	// Create styled stat displays with icons
	codePowerIcon := "⚔️"
	wisdomIcon := "📖"
	agilityIcon := "⚡"

	// Stat styles
	statNameStyle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Bold(false)

	statValueStyle := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	// Format each stat with icon, name, and value
	codePowerStat := fmt.Sprintf("%s %s %s",
		codePowerIcon,
		statNameStyle.Render("CodePower:"),
		statValueStyle.Render(fmt.Sprintf("%d", character.CodePower)))

	wisdomStat := fmt.Sprintf("%s %s %s",
		wisdomIcon,
		statNameStyle.Render("Wisdom:"),
		statValueStyle.Render(fmt.Sprintf("%d", character.Wisdom)))

	agilityStat := fmt.Sprintf("%s %s %s",
		agilityIcon,
		statNameStyle.Render("Agility:"),
		statValueStyle.Render(fmt.Sprintf("%d", character.Agility)))

	// Join stats horizontally with spacing
	statsLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		codePowerStat+"  ",
		wisdomStat+"  ",
		agilityStat,
	)

	// Add explanations for each stat
	codePowerDesc := MutedTextStyle.Render("  ├─ Increases commit quality bonus")
	wisdomDesc := MutedTextStyle.Render("  ├─ Multiplies XP gain from activities")
	agilityDesc := MutedTextStyle.Render("  └─ Speeds up quest completion")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		statsLine,
		"",
		codePowerDesc,
		wisdomDesc,
		agilityDesc,
	)
}

// renderStreakSectionDetailed renders detailed streak information.
func renderStreakSectionDetailed(character *game.Character) string {
	title := SubtitleStyle.Render("🔥 Streaks")

	// Current streak with icon
	currentLabel := StatLabelStyle.Render("Current: ")
	currentValue := lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true).
		Render(fmt.Sprintf("%d days 🔥", character.CurrentStreak))
	current := currentLabel + currentValue

	// Longest streak with trophy
	longestLabel := StatLabelStyle.Render("Longest: ")
	longestValue := lipgloss.NewStyle().
		Foreground(ColorXP).
		Bold(true).
		Render(fmt.Sprintf("%d days 🏆", character.LongestStreak))
	longest := longestLabel + longestValue

	// Last active
	lastActiveLabel := StatLabelStyle.Render("Last Active: ")
	lastActiveValue := DimTextStyle.Render(formatDate(character.LastActiveDate))
	lastActive := lastActiveLabel + lastActiveValue

	// Motivational message based on streak
	var motivation string
	if character.CurrentStreak == 0 {
		motivation = InfoTextStyle.Render("💡 Start your coding journey today!")
	} else if character.CurrentStreak < 3 {
		motivation = InfoTextStyle.Render("🚀 Keep the momentum going!")
	} else if character.CurrentStreak < 7 {
		motivation = SuccessTextStyle.Render("🔥 You're on fire!")
	} else if character.CurrentStreak < 30 {
		motivation = SuccessTextStyle.Render("⚡ Incredible consistency!")
	} else {
		motivation = SuccessTextStyle.Render("🌟 LEGENDARY STREAK!")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		current,
		longest,
		lastActive,
		"",
		motivation,
	)
}

// renderTodayActivityDetailed renders detailed today's activity section.
func renderTodayActivityDetailed(character *game.Character) string {
	title := SubtitleStyle.Render("📊 Today's Activity")

	// Commits today
	commitsLabel := StatLabelStyle.Render("Commits: ")
	commitsValue := StatValueStyle.Render(fmt.Sprintf("%d", character.TodayCommits))
	commitsIcon := "💾"
	commits := commitsIcon + " " + commitsLabel + commitsValue

	// Lines added today
	linesLabel := StatLabelStyle.Render("Lines Added: ")
	linesValue := StatValueStyle.Render(fmt.Sprintf("%d", character.TodayLinesAdded))
	linesIcon := "📝"
	lines := linesIcon + " " + linesLabel + linesValue

	// Session time today
	sessionLabel := StatLabelStyle.Render("Session Time: ")
	sessionValue := StatValueStyle.Render(formatDuration(character.TodaySessionTime))
	sessionIcon := "⏱"
	session := sessionIcon + " " + sessionLabel + sessionValue

	// Productivity rating
	productivityLabel := StatLabelStyle.Render("Productivity: ")
	var productivityValue string
	var productivityColor lipgloss.Color
	if character.TodayCommits == 0 {
		productivityValue = "Not started"
		productivityColor = ColorDim
	} else if character.TodayCommits < 3 {
		productivityValue = "Getting started ⭐"
		productivityColor = ColorInfo
	} else if character.TodayCommits < 5 {
		productivityValue = "Productive ⭐⭐"
		productivityColor = ColorSuccess
	} else if character.TodayCommits < 10 {
		productivityValue = "Highly productive ⭐⭐⭐"
		productivityColor = ColorXP
	} else {
		productivityValue = "LEGENDARY ⭐⭐⭐⭐"
		productivityColor = ColorPrimary
	}
	productivity := productivityLabel + lipgloss.NewStyle().
		Foreground(productivityColor).
		Bold(true).
		Render(productivityValue)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		commits,
		lines,
		session,
		"",
		productivity,
	)
}

// renderLifetimeStatsDetailed renders detailed lifetime statistics.
func renderLifetimeStatsDetailed(character *game.Character) string {
	title := SubtitleStyle.Render("📈 Lifetime Statistics")

	// Total commits
	totalCommitsLabel := StatLabelStyle.Render("Total Commits: ")
	totalCommitsValue := StatValueStyle.Render(fmt.Sprintf("%d", character.TotalCommits))
	totalCommitsIcon := "💾"
	totalCommits := totalCommitsIcon + " " + totalCommitsLabel + totalCommitsValue

	// Total lines added
	totalLinesLabel := StatLabelStyle.Render("Lines Added: ")
	totalLinesValue := StatValueStyle.Render(fmt.Sprintf("%d", character.TotalLinesAdded))
	totalLinesIcon := "📝"
	totalLines := totalLinesIcon + " " + totalLinesLabel + totalLinesValue

	// Total lines removed
	totalRemovedLabel := StatLabelStyle.Render("Lines Removed: ")
	totalRemovedValue := StatValueStyle.Render(fmt.Sprintf("%d", character.TotalLinesRemoved))
	totalRemovedIcon := "🗑"
	totalRemoved := totalRemovedIcon + " " + totalRemovedLabel + totalRemovedValue

	// Net lines (added - removed)
	netLines := character.TotalLinesAdded - character.TotalLinesRemoved
	netLabel := StatLabelStyle.Render("Net Lines: ")
	netValue := StatValueStyle.Render(fmt.Sprintf("%d", netLines))
	netIcon := "📊"
	net := netIcon + " " + netLabel + netValue

	// Quests completed
	questsLabel := StatLabelStyle.Render("Quests Completed: ")
	questsValue := StatValueStyle.Render(fmt.Sprintf("%d", character.QuestsCompleted))
	questsIcon := "✅"
	quests := questsIcon + " " + questsLabel + questsValue

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		totalCommits,
		totalLines,
		totalRemoved,
		net,
		"",
		quests,
	)
}

// renderAchievementsPlaceholder renders a placeholder for achievements (post-MVP).
func renderAchievementsPlaceholder() string {
	title := SubtitleStyle.Render("🏆 Achievements")

	message := MutedTextStyle.Render("Coming soon!")
	hint := InfoTextStyle.Render("Unlock achievements by completing quests and reaching milestones.")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		message,
		hint,
	)
}

// renderNoCharacterScreen renders a message when no character is loaded.
func renderNoCharacterScreen(width, height int) string {
	// Render header even for nil character
	header := renderCharacterHeader(nil, width)

	message := ErrorTextStyle.Render("⚠️  No character loaded")
	hint := MutedTextStyle.Render("This shouldn't happen. Please restart CodeQuest.")
	quit := InfoTextStyle.Render("Press Ctrl+C to quit or Esc to return to dashboard")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		message,
		"",
		hint,
		"",
		quit,
	)

	centered := placeInCenter(width, height-5, content) // Leave room for header

	// Combine header and centered content
	return header + "\n" + centered
}

// renderCharacterFooter renders the footer with key bindings.
func renderCharacterFooter(width int) string {
	// Key bindings
	dashboard := renderKeybind("Alt+Q", "Dashboard")
	mentor := renderKeybind("Alt+M", "Mentor")
	esc := renderKeybind("Esc", "Back")
	help := renderKeybind("?", "Help")

	keybinds := lipgloss.JoinHorizontal(
		lipgloss.Left,
		dashboard,
		"  ",
		mentor,
		"  ",
		esc,
		"  ",
		help,
	)

	// Center the footer
	footer := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(keybinds)

	return footer
}

// ============================================================================
// Helper Functions
// ============================================================================

// formatDate formats a time.Time to a human-readable date string.
func formatDate(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	// If it's today, show relative time
	if truncateToDay(now).Equal(truncateToDay(t)) {
		return "Today"
	}

	// If it's yesterday
	if truncateToDay(now).Add(-24 * time.Hour).Equal(truncateToDay(t)) {
		return "Yesterday"
	}

	// If it's within the last week, show "X days ago"
	if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}

	// Otherwise show the actual date
	return t.Format("Jan 2, 2006")
}

// truncateToDay truncates a time to midnight (start of day).
func truncateToDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// renderCharacterHeader creates a header for the Character screen (inline to avoid import cycle).
func renderCharacterHeader(char *game.Character, width int) string {
	// Colors for header
	colorPrimary := lipgloss.Color("205") // Pink/Magenta
	colorAccent := lipgloss.Color("86")   // Cyan
	colorLevel := lipgloss.Color("93")    // Yellow-Orange
	colorBright := lipgloss.Color("15")   // White
	colorDim := lipgloss.Color("240")     // Gray

	// If width is too small, render a minimal header
	if width < 40 {
		style := lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			Border(lipgloss.RoundedBorder(), false, false, true, false).
			BorderForeground(colorAccent).
			Width(width-4).
			Padding(0, 1).
			MarginBottom(1)
		return style.Render("🎮 Character Sheet")
	}

	// Left section: CodeQuest title
	leftStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary)
	leftSection := leftStyle.Render("🎮 CodeQuest")

	// Center section: Screen name
	centerStyle := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true)
	centerSection := centerStyle.Render("[Character Sheet]")

	// Right section: Character info
	var rightSection string
	if char == nil {
		dimStyle := lipgloss.NewStyle().
			Foreground(colorDim).
			Italic(true)
		rightSection = dimStyle.Render("No Character")
	} else {
		nameStyle := lipgloss.NewStyle().
			Foreground(colorBright).
			Bold(true)
		levelStyle := lipgloss.NewStyle().
			Foreground(colorLevel).
			Bold(true)
		name := nameStyle.Render(char.Name)
		level := levelStyle.Render(fmt.Sprintf("Lvl %d", char.Level))
		rightSection = name + " " + level
	}

	// Calculate spacing
	leftWidth := lipgloss.Width(leftSection)
	centerWidth := lipgloss.Width(centerSection)
	rightWidth := lipgloss.Width(rightSection)
	contentWidth := leftWidth + centerWidth + rightWidth
	totalSpacing := width - contentWidth - 4

	// If content is too wide, just join with spaces
	var content string
	if totalSpacing < 2 {
		content = leftSection + " " + centerSection + " " + rightSection
	} else {
		leftCenterSpacing := totalSpacing / 2
		centerRightSpacing := totalSpacing - leftCenterSpacing
		if leftCenterSpacing < 1 {
			leftCenterSpacing = 1
		}
		if centerRightSpacing < 1 {
			centerRightSpacing = 1
		}
		spacer1 := strings.Repeat(" ", leftCenterSpacing)
		spacer2 := strings.Repeat(" ", centerRightSpacing)
		content = leftSection + spacer1 + centerSection + spacer2 + rightSection
	}

	// Wrap in styled box
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, false, true, false).
		BorderForeground(colorAccent).
		Width(width-4).
		Padding(0, 1).
		MarginBottom(1)

	return style.Render(content)
}
