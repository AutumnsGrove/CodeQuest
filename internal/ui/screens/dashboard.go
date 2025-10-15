// Package screens provides screen rendering functions for CodeQuest UI.
// This file implements the Dashboard screen, the main hub of the application.
package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// Style references - imported from ui package via ui.Styles
// To avoid circular dependencies, we define styles locally
var (
	// Color Palette
	ColorPrimary   = lipgloss.Color("205") // Pink/Magenta
	ColorSecondary = lipgloss.Color("63")  // Purple
	ColorAccent    = lipgloss.Color("86")  // Cyan
	ColorSuccess   = lipgloss.Color("42")  // Green
	ColorWarning   = lipgloss.Color("214") // Orange
	ColorError     = lipgloss.Color("196") // Red
	ColorInfo      = lipgloss.Color("69")  // Blue
	ColorDim       = lipgloss.Color("240") // Gray
	ColorBright    = lipgloss.Color("15")  // White
	ColorMuted     = lipgloss.Color("243") // Light gray
	ColorXP        = lipgloss.Color("226") // Gold/Yellow
	ColorLevel     = lipgloss.Color("93")  // Yellow-Orange
	ColorQuest     = lipgloss.Color("111") // Light Blue
	ColorMagic     = lipgloss.Color("177") // Lavender

	// Text Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			MarginTop(1).
			MarginBottom(0)

	HeadingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	TextStyle = lipgloss.NewStyle().
			Foreground(ColorBright)

	BoldTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBright)

	DimTextStyle = lipgloss.NewStyle().
			Foreground(ColorDim)

	MutedTextStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Italic(true)

	ErrorTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorError)

	SuccessTextStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorSuccess)

	WarningTextStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorWarning)

	InfoTextStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)

	// Box Styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(1, 2)

	BoxStyleFocused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			Bold(true)

	BoxStyleDim = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDim).
			Padding(1, 2)

	// Label and value styles
	StatLabelStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Bold(true)

	StatValueStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	KeybindStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	KeybindDescStyle = lipgloss.NewStyle().
				Foreground(ColorBright)

	// Progress bar styles
	XPBarStyle = lipgloss.NewStyle().
			Foreground(ColorXP).
			Bold(true)

	QuestProgressBarStyle = lipgloss.NewStyle().
				Foreground(ColorQuest).
				Bold(true)

	ProgressBarEmptyStyle = lipgloss.NewStyle().
				Foreground(ColorDim)
)

// RenderDashboard renders the main dashboard screen.
// This is the primary screen players see, showing character overview, active quest, and quick actions.
//
// The dashboard layout is responsive and adapts to terminal dimensions:
//   - Wide terminals (>100 cols): Side-by-side panels with timer
//   - Narrow terminals (≤100 cols): Stacked vertical layout
//
// Parameters:
//   - character: Player character data (nil-safe)
//   - quests: All quests (used to find active quest)
//   - width: Terminal width in characters
//   - height: Terminal height in characters
//
// Returns:
//   - string: Rendered dashboard UI
//
// Note: The timer display is currently a placeholder showing character's TodaySessionTime.
// Full timer state management will be added by Subagent 32.
func RenderDashboard(character *game.Character, quests []*game.Quest, width, height int) string {
	// Handle nil character gracefully
	if character == nil {
		return renderNoCharacter(width, height)
	}

	// Determine layout based on terminal width
	useWideLayout := width > 100

	if useWideLayout {
		return renderDashboardWide(character, quests, width, height)
	}
	return renderDashboardNarrow(character, quests, width, height)
}

// renderDashboardWide renders dashboard with side-by-side panels for wide terminals.
func renderDashboardWide(character *game.Character, quests []*game.Quest, width, height int) string {
	// Split width into two columns (60% left, 40% right)
	leftWidth := int(float64(width) * 0.58)
	rightWidth := width - leftWidth - 2 // Account for spacing

	// Render left panel: Character overview and stats
	leftPanel := renderCharacterPanel(character, leftWidth)

	// Render right panel: Active quest, today's stats, and timer
	activeQuest := findActiveQuest(quests)
	rightPanel := renderActivityPanel(character, activeQuest, rightWidth)

	// Join panels horizontally
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		rightPanel,
	)

	// Render timer section (inline display)
	timerSection := renderTimerSection(character)

	// Render quick actions menu at bottom
	quickActions := renderQuickActions(width)

	// Join vertically with spacing
	dashboard := lipgloss.JoinVertical(
		lipgloss.Left,
		panels,
		"",
		timerSection,
		"",
		quickActions,
	)

	return dashboard
}

// renderDashboardNarrow renders dashboard with stacked panels for narrow terminals.
func renderDashboardNarrow(character *game.Character, quests []*game.Quest, width, height int) string {
	// Full width for each panel
	panelWidth := width

	// Render panels vertically
	charPanel := renderCharacterPanel(character, panelWidth)
	activeQuest := findActiveQuest(quests)
	activityPanel := renderActivityPanel(character, activeQuest, panelWidth)
	timerSection := renderTimerSection(character)
	quickActions := renderQuickActions(width)

	// Stack all panels
	dashboard := lipgloss.JoinVertical(
		lipgloss.Left,
		charPanel,
		"",
		activityPanel,
		"",
		timerSection,
		"",
		quickActions,
	)

	return dashboard
}

// renderCharacterPanel renders the character overview panel.
// Shows name, level, XP progress bar, and core stats.
func renderCharacterPanel(character *game.Character, width int) string {
	// Title with icon
	title := renderTitle("Character", "⚔️")

	// Character name with emphasis
	nameLabel := StatLabelStyle.Render("Name: ")
	nameValue := BoldTextStyle.Render(character.Name)
	name := nameLabel + nameValue

	// Level display
	levelLabel := StatLabelStyle.Render("Level: ")
	levelValue := lipgloss.NewStyle().
		Foreground(ColorLevel).
		Bold(true).
		Render(fmt.Sprintf("%d", character.Level))
	level := levelLabel + levelValue

	// XP progress bar
	xpLabel := StatLabelStyle.Render("XP: ")
	xpBar := renderProgressBar(
		character.XP,
		character.XPToNextLevel,
		width-25, // Account for label and padding
		"xp",
	)
	xp := xpLabel + xpBar

	// Core stats in a grid layout
	statsTitle := SubtitleStyle.Render("Core Stats")

	codePowerLabel := StatLabelStyle.Render("Code Power: ")
	codePowerValue := StatValueStyle.Render(fmt.Sprintf("%d", character.CodePower))
	codePower := codePowerLabel + codePowerValue + MutedTextStyle.Render(" (commit quality)")

	wisdomLabel := StatLabelStyle.Render("Wisdom: ")
	wisdomValue := StatValueStyle.Render(fmt.Sprintf("%d", character.Wisdom))
	wisdom := wisdomLabel + wisdomValue + MutedTextStyle.Render(" (XP multiplier)")

	agilityLabel := StatLabelStyle.Render("Agility: ")
	agilityValue := StatValueStyle.Render(fmt.Sprintf("%d", character.Agility))
	agility := agilityLabel + agilityValue + MutedTextStyle.Render(" (quest speed)")

	// Streak display
	streakLabel := StatLabelStyle.Render("Current Streak: ")
	streakValue := lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true).
		Render(fmt.Sprintf("%d days 🔥", character.CurrentStreak))
	streak := streakLabel + streakValue

	longestLabel := StatLabelStyle.Render("Longest Streak: ")
	longestValue := lipgloss.NewStyle().
		Foreground(ColorXP).
		Bold(true).
		Render(fmt.Sprintf("%d days 🏆", character.LongestStreak))
	longest := longestLabel + longestValue

	// Lifetime stats
	lifetimeTitle := SubtitleStyle.Render("Lifetime Stats")

	totalCommitsLabel := StatLabelStyle.Render("Total Commits: ")
	totalCommitsValue := StatValueStyle.Render(fmt.Sprintf("%d", character.TotalCommits))
	totalCommits := totalCommitsLabel + totalCommitsValue

	totalLinesLabel := StatLabelStyle.Render("Lines Added: ")
	totalLinesValue := StatValueStyle.Render(fmt.Sprintf("%d", character.TotalLinesAdded))
	totalLines := totalLinesLabel + totalLinesValue

	questsLabel := StatLabelStyle.Render("Quests Completed: ")
	questsValue := StatValueStyle.Render(fmt.Sprintf("%d", character.QuestsCompleted))
	totalQuests := questsLabel + questsValue

	// Assemble panel content
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		name,
		level,
		xp,
		"",
		statsTitle,
		codePower,
		wisdom,
		agility,
		"",
		streak,
		longest,
		"",
		lifetimeTitle,
		totalCommits,
		totalLines,
		totalQuests,
	)

	// Wrap in box
	return BoxStyle.Width(width - 4).Render(content)
}

// renderActivityPanel renders the activity panel showing active quest and today's stats.
func renderActivityPanel(character *game.Character, activeQuest *game.Quest, width int) string {
	// Active quest section
	var questSection string
	if activeQuest != nil {
		questSection = renderActiveQuestCard(activeQuest, width)
	} else {
		questSection = renderNoActiveQuest(width)
	}

	// Today's activity section
	todaySection := renderTodayActivity(character, width)

	// Combine sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		questSection,
		"",
		todaySection,
	)

	return content
}

// renderActiveQuestCard renders the active quest card with progress.
func renderActiveQuestCard(quest *game.Quest, width int) string {
	title := renderTitle("Active Quest", "📋")

	// Quest title with type badge
	questTitle := BoldTextStyle.Render(quest.Title)
	typeBadge := renderQuestTypeBadge(quest.Type)
	header := questTitle + " " + typeBadge

	// Quest description (truncate if too long)
	description := quest.Description
	maxDescLen := width - 10
	if len(description) > maxDescLen {
		description = description[:maxDescLen-3] + "..."
	}
	desc := TextStyle.Render(description)

	// Progress bar
	progressLabel := StatLabelStyle.Render("Progress: ")
	progressBar := renderProgressBar(
		quest.Current,
		quest.Target,
		width-30, // Account for label
		"quest",
	)
	progress := progressLabel + progressBar

	// Time tracking
	timeLabel := StatLabelStyle.Render("Started: ")
	var timeValue string
	if quest.StartedAt != nil {
		elapsed := time.Since(*quest.StartedAt)
		timeValue = MutedTextStyle.Render(formatDuration(elapsed) + " ago")
	} else {
		timeValue = MutedTextStyle.Render("Unknown")
	}
	timeInfo := timeLabel + timeValue

	// XP reward
	rewardLabel := StatLabelStyle.Render("Reward: ")
	rewardValue := lipgloss.NewStyle().
		Foreground(ColorXP).
		Bold(true).
		Render(fmt.Sprintf("%d XP ⭐", quest.XPReward))
	reward := rewardLabel + rewardValue

	// Assemble card
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		header,
		desc,
		"",
		progress,
		timeInfo,
		reward,
	)

	// Wrap in focused box to emphasize
	return BoxStyleFocused.Width(width - 4).Render(content)
}

// renderNoActiveQuest renders a message when no quest is active.
func renderNoActiveQuest(width int) string {
	title := renderTitle("Active Quest", "📋")

	message := MutedTextStyle.Render("No active quest")
	hint := InfoTextStyle.Render("Press [Q] to browse the Quest Board and start a new quest!")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		message,
		"",
		hint,
	)

	return BoxStyleDim.Width(width - 4).Render(content)
}

// renderTodayActivity renders today's activity statistics.
func renderTodayActivity(character *game.Character, width int) string {
	title := renderTitle("Today's Activity", "📊")

	// Commits today
	commitsLabel := StatLabelStyle.Render("Commits: ")
	commitsValue := StatValueStyle.Render(fmt.Sprintf("%d", character.TodayCommits))
	commits := commitsLabel + commitsValue

	// Lines added today
	linesLabel := StatLabelStyle.Render("Lines Added: ")
	linesValue := StatValueStyle.Render(fmt.Sprintf("%d", character.TodayLinesAdded))
	lines := linesLabel + linesValue

	// Session time today
	sessionLabel := StatLabelStyle.Render("Session Time: ")
	sessionValue := StatValueStyle.Render(formatDuration(character.TodaySessionTime))
	session := sessionLabel + sessionValue

	// Motivational message based on activity
	var motivation string
	if character.TodayCommits == 0 {
		motivation = WarningTextStyle.Render("💡 No commits yet today. Time to code!")
	} else if character.TodayCommits < 3 {
		motivation = InfoTextStyle.Render("🚀 Great start! Keep the momentum going!")
	} else if character.TodayCommits < 5 {
		motivation = SuccessTextStyle.Render("🔥 On fire! You're crushing it today!")
	} else {
		motivation = SuccessTextStyle.Render("⚡ LEGENDARY! Amazing productivity today!")
	}

	// Assemble content
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		commits,
		lines,
		session,
		"",
		motivation,
	)

	return BoxStyle.Width(width - 4).Render(content)
}

// renderTimerSection renders the session timer display (inline to avoid import cycle).
// This shows the current coding session time.
//
// Note: Currently displays character's TodaySessionTime as a static value.
// Full timer state management (start/stop/tick) will be added by Subagent 32.
func renderTimerSection(character *game.Character) string {
	if character == nil {
		return ""
	}

	// Format the duration inline (avoiding import cycle with components)
	duration := character.TodaySessionTime
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	timeStr := fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)

	// Color coding based on duration
	var color lipgloss.Color
	hoursDur := duration.Hours()
	switch {
	case hoursDur >= 5.0:
		color = ColorWarning // Orange
	case hoursDur >= 3.0:
		color = ColorSuccess // Green
	case hoursDur >= 1.0:
		color = ColorInfo // Cyan
	default:
		color = ColorDim // Gray
	}

	// Timer display with icon (paused for now)
	timerStyle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true)
	icon := "⏸" // Pause symbol (not running yet)
	timerDisplay := timerStyle.Render(icon + " " + timeStr)

	// Create a small badge-style display
	label := StatLabelStyle.Render("Session: ")
	hint := MutedTextStyle.Render(" (Press Ctrl+T to start/stop timer)")

	content := label + timerDisplay + hint

	// Center the timer section
	centered := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(content)

	return centered
}

// renderQuickActions renders the quick actions menu at the bottom of dashboard.
func renderQuickActions(width int) string {
	title := HeadingStyle.Render("Quick Actions")

	// Key bindings with descriptions
	questKey := renderKeybind("Q", "Quest Board")
	charKey := renderKeybind("C", "Character Sheet")
	mentorKey := renderKeybind("M", "AI Mentor")
	settingsKey := renderKeybind("S", "Settings")

	// Save and quit keys
	saveKey := renderKeybind("Ctrl+S", "Save")
	quitKey := renderKeybind("Ctrl+C", "Quit")

	// Layout keys in rows
	row1 := lipgloss.JoinHorizontal(
		lipgloss.Left,
		questKey,
		"  ",
		charKey,
		"  ",
		mentorKey,
		"  ",
		settingsKey,
	)

	row2 := lipgloss.JoinHorizontal(
		lipgloss.Left,
		saveKey,
		"  ",
		quitKey,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		row1,
		row2,
	)

	// Center the quick actions
	centered := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(content)

	return centered
}

// renderNoCharacter renders a message when no character is loaded.
func renderNoCharacter(width, height int) string {
	message := ErrorTextStyle.Render("⚠️  No character loaded")
	hint := MutedTextStyle.Render("This shouldn't happen. Please restart CodeQuest.")
	quit := InfoTextStyle.Render("Press Ctrl+C to quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		message,
		"",
		hint,
		"",
		quit,
	)

	return placeInCenter(width, height, content)
}

// ============================================================================
// Helper Functions - Rendering Utilities
// ============================================================================

// renderTitle renders a styled title with optional icon
func renderTitle(text string, icon string) string {
	if icon != "" {
		return TitleStyle.Render(icon + " " + text)
	}
	return TitleStyle.Render(text)
}

// renderKeybind formats a keybind hint (e.g., "[Q] Quit")
func renderKeybind(key, description string) string {
	return KeybindStyle.Render("["+key+"]") + " " + KeybindDescStyle.Render(description)
}

// renderProgressBar creates a progress bar with percentage
func renderProgressBar(current, total, width int, barType string) string {
	if total == 0 {
		total = 1 // Prevent division by zero
	}

	percentage := float64(current) / float64(total)
	if percentage > 1.0 {
		percentage = 1.0
	}

	filledWidth := int(float64(width) * percentage)
	emptyWidth := width - filledWidth

	// Choose style based on bar type
	var filledStyle, emptyStyle lipgloss.Style
	var fillChar, emptyChar string

	switch barType {
	case "xp":
		filledStyle = XPBarStyle
		emptyStyle = ProgressBarEmptyStyle
		fillChar = "█"
		emptyChar = "░"
	case "quest":
		filledStyle = QuestProgressBarStyle
		emptyStyle = ProgressBarEmptyStyle
		fillChar = "▰"
		emptyChar = "▱"
	default:
		filledStyle = XPBarStyle
		emptyStyle = ProgressBarEmptyStyle
		fillChar = "█"
		emptyChar = "░"
	}

	filled := strings.Repeat(fillChar, filledWidth)
	empty := strings.Repeat(emptyChar, emptyWidth)

	bar := filledStyle.Render(filled) + emptyStyle.Render(empty)
	percentText := fmt.Sprintf(" %d/%d (%.0f%%)", current, total, percentage*100)

	return bar + DimTextStyle.Render(percentText)
}

// placeInCenter centers content in a given width and height
func placeInCenter(width, height int, content string) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// ============================================================================
// Helper Functions - Business Logic
// ============================================================================

// findActiveQuest finds the first active quest from the quest list.
func findActiveQuest(quests []*game.Quest) *game.Quest {
	for _, quest := range quests {
		if quest.Status == game.QuestActive {
			return quest
		}
	}
	return nil
}

// renderQuestTypeBadge renders a styled badge for quest types.
func renderQuestTypeBadge(questType game.QuestType) string {
	var badge string
	var color lipgloss.Color

	switch questType {
	case game.QuestTypeCommit:
		badge = "COMMIT"
		color = ColorSuccess
	case game.QuestTypeLines:
		badge = "LINES"
		color = ColorInfo
	case game.QuestTypeTests:
		badge = "TESTS"
		color = ColorAccent
	case game.QuestTypePR:
		badge = "PR"
		color = ColorPrimary
	case game.QuestTypeRefactor:
		badge = "REFACTOR"
		color = ColorWarning
	case game.QuestTypeDaily:
		badge = "DAILY"
		color = ColorXP
	case game.QuestTypeStreak:
		badge = "STREAK"
		color = ColorMagic
	default:
		badge = "QUEST"
		color = ColorDim
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(color).
		Bold(true).
		Padding(0, 1)

	return style.Render(badge)
}

// formatDuration formats a duration into a human-readable string.
func formatDuration(d time.Duration) string {
	// Handle zero duration
	if d == 0 {
		return "0s"
	}

	// Calculate components
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	// Build string based on what components are non-zero
	var parts []string

	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}
