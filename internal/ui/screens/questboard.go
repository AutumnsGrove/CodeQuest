// Package screens provides screen rendering functions for CodeQuest UI.
// This file implements the Quest Board screen, displaying all quests with filtering and navigation.
package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// QuestFilter represents filter options for displaying quests.
// Used to show only quests matching a specific status.
type QuestFilter int

const (
	// FilterAll shows all quests regardless of status
	FilterAll QuestFilter = iota
	// FilterAvailable shows only available quests
	FilterAvailable
	// FilterActive shows only active quests
	FilterActive
	// FilterCompleted shows only completed quests
	FilterCompleted
)

// RenderQuestBoard renders the quest board screen with all quests.
// This is the main quest management screen where players can browse and select quests.
//
// Features:
//   - Lists all quests grouped by status (Available, Active, Completed)
//   - Supports quest filtering by status
//   - Highlights selected quest for keyboard navigation
//   - Shows quest details: title, description, type badge, progress, XP reward
//   - Scrollable list if quests exceed screen height
//   - Responsive layout for different terminal sizes
//   - Nil-safe: handles empty quest list gracefully
//
// Layout Structure:
//   - Header: Screen title with character info
//   - Filter tabs: Quick filter by status
//   - Quest sections: Grouped by status (Available, Active, Completed)
//   - Footer: Key bindings and navigation help
//
// Parameters:
//   - character: Player character (for header display)
//   - quests: All quests to display
//   - selectedIndex: Index of currently selected quest (-1 for none)
//   - filter: Current filter setting
//   - width: Terminal width in characters
//   - height: Terminal height in characters
//
// Returns:
//   - string: Rendered quest board UI
func RenderQuestBoard(character *game.Character, quests []*game.Quest, selectedIndex int, filter QuestFilter, width, height int) string {
	// Render header (inline to avoid import cycle)
	header := renderQuestBoardHeader(character, width)

	// Filter quests based on current filter
	filteredQuests := filterQuests(quests, filter)

	// Render filter tabs
	filterTabs := renderFilterTabs(filter, quests, width)

	// Render quest list
	var questList string
	if len(filteredQuests) == 0 {
		questList = renderEmptyQuestList(filter, width)
	} else {
		questList = renderQuestList(filteredQuests, selectedIndex, width, height-15)
	}

	// Render footer with key bindings
	footer := renderQuestBoardFooter(width)

	// Assemble screen
	screen := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		filterTabs,
		"",
		questList,
		"",
		footer,
	)

	return screen
}

// filterQuests filters quests based on the selected filter.
func filterQuests(quests []*game.Quest, filter QuestFilter) []*game.Quest {
	if filter == FilterAll {
		return quests
	}

	filtered := make([]*game.Quest, 0)
	for _, quest := range quests {
		switch filter {
		case FilterAvailable:
			if quest.Status == game.QuestAvailable {
				filtered = append(filtered, quest)
			}
		case FilterActive:
			if quest.Status == game.QuestActive {
				filtered = append(filtered, quest)
			}
		case FilterCompleted:
			if quest.Status == game.QuestCompleted {
				filtered = append(filtered, quest)
			}
		}
	}

	return filtered
}

// renderFilterTabs renders the filter tabs showing quest counts by status.
func renderFilterTabs(currentFilter QuestFilter, quests []*game.Quest, width int) string {
	// Count quests by status
	availableCount := 0
	activeCount := 0
	completedCount := 0
	totalCount := len(quests)

	for _, quest := range quests {
		switch quest.Status {
		case game.QuestAvailable:
			availableCount++
		case game.QuestActive:
			activeCount++
		case game.QuestCompleted:
			completedCount++
		}
	}

	// Create tab styles
	activeTabStyle := lipgloss.NewStyle().
		Foreground(ColorBright).
		Background(ColorPrimary).
		Bold(true).
		Padding(0, 2)

	inactiveTabStyle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Background(ColorDim).
		Padding(0, 2)

	// Render each tab
	allTab := fmt.Sprintf("All (%d)", totalCount)
	availableTab := fmt.Sprintf("📋 Available (%d)", availableCount)
	activeTab := fmt.Sprintf("⚡ Active (%d)", activeCount)
	completedTab := fmt.Sprintf("✅ Completed (%d)", completedCount)

	// Apply active/inactive styles
	var tabs []string
	if currentFilter == FilterAll {
		tabs = append(tabs, activeTabStyle.Render(allTab))
	} else {
		tabs = append(tabs, inactiveTabStyle.Render(allTab))
	}

	if currentFilter == FilterAvailable {
		tabs = append(tabs, activeTabStyle.Render(availableTab))
	} else {
		tabs = append(tabs, inactiveTabStyle.Render(availableTab))
	}

	if currentFilter == FilterActive {
		tabs = append(tabs, activeTabStyle.Render(activeTab))
	} else {
		tabs = append(tabs, inactiveTabStyle.Render(activeTab))
	}

	if currentFilter == FilterCompleted {
		tabs = append(tabs, activeTabStyle.Render(completedTab))
	} else {
		tabs = append(tabs, inactiveTabStyle.Render(completedTab))
	}

	// Join tabs horizontally with spacing
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)

	return tabBar
}

// renderQuestList renders the list of quests with the selected one highlighted.
func renderQuestList(quests []*game.Quest, selectedIndex int, width, maxHeight int) string {
	// Group quests by status
	availableQuests := make([]*game.Quest, 0)
	activeQuests := make([]*game.Quest, 0)
	completedQuests := make([]*game.Quest, 0)

	for _, quest := range quests {
		switch quest.Status {
		case game.QuestAvailable:
			availableQuests = append(availableQuests, quest)
		case game.QuestActive:
			activeQuests = append(activeQuests, quest)
		case game.QuestCompleted:
			completedQuests = append(completedQuests, quest)
		}
	}

	// Render each section
	sections := make([]string, 0)

	if len(availableQuests) > 0 {
		sections = append(sections, renderQuestSection("📋 Available Quests", availableQuests, selectedIndex, 0, width))
	}

	if len(activeQuests) > 0 {
		offset := len(availableQuests)
		sections = append(sections, renderQuestSection("⚡ Active Quests", activeQuests, selectedIndex, offset, width))
	}

	if len(completedQuests) > 0 {
		offset := len(availableQuests) + len(activeQuests)
		sections = append(sections, renderQuestSection("✅ Completed Quests", completedQuests, selectedIndex, offset, width))
	}

	// Join sections vertically
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderQuestSection renders a section of quests with a title.
func renderQuestSection(title string, quests []*game.Quest, selectedIndex, offset int, width int) string {
	sectionTitle := SubtitleStyle.Render(title)

	questCards := make([]string, 0)
	for i, quest := range quests {
		globalIndex := offset + i
		isSelected := globalIndex == selectedIndex
		card := renderQuestCard(quest, isSelected, width-4)
		questCards = append(questCards, card)
	}

	questList := lipgloss.JoinVertical(lipgloss.Left, questCards...)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		sectionTitle,
		"",
		questList,
		"",
	)
}

// renderQuestCard renders a single quest card.
func renderQuestCard(quest *game.Quest, selected bool, width int) string {
	// Choose style based on selection
	cardStyle := BoxStyle
	if selected {
		cardStyle = BoxStyleFocused
	}

	// Selection indicator
	indicator := "  "
	if selected {
		indicator = "▶ "
	}

	// Quest title with type badge
	questTitle := BoldTextStyle.Render(quest.Title)
	typeBadge := renderQuestTypeBadge(quest.Type)
	header := indicator + questTitle + " " + typeBadge

	// Description (truncate if too long)
	description := quest.Description
	maxDescLen := width - 10
	if len(description) > maxDescLen {
		description = description[:maxDescLen-3] + "..."
	}
	desc := TextStyle.Render(description)

	// Status-specific content
	var statusContent string
	switch quest.Status {
	case game.QuestAvailable:
		statusContent = renderAvailableQuestInfo(quest)
	case game.QuestActive:
		// Ensure barWidth is at least 10 to prevent negative values
		barWidth := width - 20
		if barWidth < 10 {
			barWidth = 10
		}
		statusContent = renderActiveQuestInfo(quest, barWidth)
	case game.QuestCompleted:
		statusContent = renderCompletedQuestInfo(quest)
	default:
		statusContent = ""
	}

	// Assemble card content
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		desc,
		"",
		statusContent,
	)

	return cardStyle.Width(width - 4).Render(content)
}

// renderAvailableQuestInfo renders info for available quests.
func renderAvailableQuestInfo(quest *game.Quest) string {
	// XP reward
	rewardLabel := StatLabelStyle.Render("Reward: ")
	rewardValue := lipgloss.NewStyle().
		Foreground(ColorXP).
		Bold(true).
		Render(fmt.Sprintf("%d XP ⭐", quest.XPReward))
	reward := rewardLabel + rewardValue

	// Required level
	levelLabel := StatLabelStyle.Render("Required Level: ")
	levelValue := StatValueStyle.Render(fmt.Sprintf("%d", quest.RequiredLevel))
	level := levelLabel + levelValue

	// Start hint
	hint := InfoTextStyle.Render("Press [Enter] to start this quest")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		reward,
		level,
		"",
		hint,
	)
}

// renderActiveQuestInfo renders info for active quests.
func renderActiveQuestInfo(quest *game.Quest, barWidth int) string {
	// Progress bar
	progressLabel := StatLabelStyle.Render("Progress: ")
	progressBar := renderProgressBar(
		quest.Current,
		quest.Target,
		barWidth,
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

	return lipgloss.JoinVertical(
		lipgloss.Left,
		progress,
		timeInfo,
		reward,
	)
}

// renderCompletedQuestInfo renders info for completed quests.
func renderCompletedQuestInfo(quest *game.Quest) string {
	// XP earned
	xpLabel := StatLabelStyle.Render("XP Earned: ")
	xpValue := lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true).
		Render(fmt.Sprintf("%d XP ✓", quest.XPReward))
	xpEarned := xpLabel + xpValue

	// Completion time
	timeLabel := StatLabelStyle.Render("Completed: ")
	var timeValue string
	if quest.CompletedAt != nil {
		elapsed := time.Since(*quest.CompletedAt)
		timeValue = MutedTextStyle.Render(formatDuration(elapsed) + " ago")
	} else {
		timeValue = MutedTextStyle.Render("Unknown")
	}
	timeInfo := timeLabel + timeValue

	// Duration (if available)
	var duration string
	if quest.StartedAt != nil && quest.CompletedAt != nil {
		totalDuration := quest.CompletedAt.Sub(*quest.StartedAt)
		durationLabel := StatLabelStyle.Render("Duration: ")
		durationValue := MutedTextStyle.Render(formatDuration(totalDuration))
		duration = durationLabel + durationValue
	}

	if duration != "" {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			xpEarned,
			timeInfo,
			duration,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		xpEarned,
		timeInfo,
	)
}

// renderEmptyQuestList renders a message when no quests match the filter.
func renderEmptyQuestList(filter QuestFilter, width int) string {
	var message string
	var hint string

	switch filter {
	case FilterAll:
		message = "No quests available"
		hint = "Quest system is initializing. Check back soon!"
	case FilterAvailable:
		message = "No available quests"
		hint = "Complete your active quests to unlock more!"
	case FilterActive:
		message = "No active quests"
		hint = "Browse available quests and start your adventure!"
	case FilterCompleted:
		message = "No completed quests"
		hint = "Complete quests to build your achievement history!"
	default:
		message = "No quests found"
		hint = "Try changing the filter or check back later."
	}

	messageStyled := MutedTextStyle.Render(message)
	hintStyled := InfoTextStyle.Render(hint)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		messageStyled,
		"",
		hintStyled,
	)

	return BoxStyleDim.Width(width - 8).Render(content)
}

// renderQuestBoardFooter renders the footer with key bindings.
func renderQuestBoardFooter(width int) string {
	// Key bindings
	upDown := renderKeybind("↑/↓", "Navigate")
	enter := renderKeybind("Enter", "Start/View")
	filter := renderKeybind("F", "Filter")
	esc := renderKeybind("Esc", "Back")

	keybinds := lipgloss.JoinHorizontal(
		lipgloss.Left,
		upDown,
		"  ",
		enter,
		"  ",
		filter,
		"  ",
		esc,
	)

	// Center the footer
	footer := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(keybinds)

	return footer
}

// renderQuestBoardHeader creates a header for the Quest Board screen.
// Simplified version to avoid import cycle with components package.
func renderQuestBoardHeader(char *game.Character, width int) string {
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
			Width(width - 4).
			Padding(0, 1).
			MarginBottom(1)
		return style.Render("🎮 Quest Board")
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
	centerSection := centerStyle.Render("[Quest Board]")

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
		Width(width - 4).
		Padding(0, 1).
		MarginBottom(1)

	return style.Render(content)
}
