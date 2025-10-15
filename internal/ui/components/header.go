// Package components provides reusable UI components for CodeQuest screens.
// This file implements the header component that appears at the top of every screen.
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// Color constants (copied from ui package to avoid import cycle)
var (
	colorPrimary = lipgloss.Color("205") // Pink/Magenta
	colorAccent  = lipgloss.Color("86")  // Cyan
	colorLevel   = lipgloss.Color("93")  // Yellow-Orange
	colorBright  = lipgloss.Color("15")  // White
	colorDim     = lipgloss.Color("240") // Gray
)

// RenderHeader creates a consistent header for all screens.
// The header includes the CodeQuest title, current screen name, and character info.
//
// Layout:
//
//	┌────────────────────────────────────────────────┐
//	│ 🎮 CodeQuest          [Dashboard]    Player Lvl 5 │
//	└────────────────────────────────────────────────┘
//
// Parameters:
//   - screenName: The name of the current screen (e.g., "Dashboard", "Quest Board")
//   - char: The character to display info for (can be nil)
//   - width: The total available width in characters
//
// Returns:
//   - string: The rendered header with proper width and styling
func RenderHeader(screenName string, char *game.Character, width int) string {
	// If width is too small, render a minimal header
	if width < 40 {
		return renderMinimalHeader(screenName, width)
	}

	// Left section: CodeQuest title with icon
	leftSection := renderLeftSection()

	// Center section: Current screen indicator
	centerSection := renderCenterSection(screenName)

	// Right section: Character name and level
	rightSection := renderRightSection(char)

	// Join sections with proper spacing
	header := joinHeaderSections(leftSection, centerSection, rightSection, width)

	// Wrap in a box for visual separation
	return wrapHeader(header, width)
}

// renderLeftSection creates the left part of the header with the CodeQuest branding.
//
// Returns:
//   - string: Styled "🎮 CodeQuest" text
func renderLeftSection() string {
	icon := "🎮"
	title := "CodeQuest"

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary)

	return style.Render(icon + " " + title)
}

// renderCenterSection creates the center part of the header with the screen name.
//
// Parameters:
//   - screenName: The name of the current screen
//
// Returns:
//   - string: Styled "[Screen Name]" text
func renderCenterSection(screenName string) string {
	style := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true)

	return style.Render("[" + screenName + "]")
}

// renderRightSection creates the right part of the header with character info.
//
// Parameters:
//   - char: The character to display (nil-safe)
//
// Returns:
//   - string: Styled "Name Lvl X" text, or empty if no character
func renderRightSection(char *game.Character) string {
	if char == nil {
		// No character loaded yet
		style := lipgloss.NewStyle().
			Foreground(colorDim).
			Italic(true)
		return style.Render("No Character")
	}

	// Character name style
	nameStyle := lipgloss.NewStyle().
		Foreground(colorBright).
		Bold(true)

	// Level indicator style
	levelStyle := lipgloss.NewStyle().
		Foreground(colorLevel).
		Bold(true)

	name := nameStyle.Render(char.Name)
	level := levelStyle.Render(fmt.Sprintf("Lvl %d", char.Level))

	// Join with a space
	return name + " " + level
}

// joinHeaderSections combines left, center, and right sections with proper spacing.
// This uses Lip Gloss's horizontal join to distribute space evenly.
//
// Parameters:
//   - left: Left section content
//   - center: Center section content
//   - right: Right section content
//   - width: Total available width
//
// Returns:
//   - string: Sections joined with proper spacing
func joinHeaderSections(left, center, right string, width int) string {
	// Calculate the width of each section's rendered content (strip ANSI codes for measurement)
	leftWidth := lipgloss.Width(left)
	centerWidth := lipgloss.Width(center)
	rightWidth := lipgloss.Width(right)

	// Calculate total content width
	contentWidth := leftWidth + centerWidth + rightWidth

	// If content is too wide, render compactly
	if contentWidth >= width-4 {
		// Just join with single spaces
		return left + " " + center + " " + right
	}

	// Calculate spacing needed
	// We want: left [spaces] center [spaces] right
	totalSpacing := width - contentWidth - 4 // -4 for padding

	// Split spacing between left-center and center-right
	leftCenterSpacing := totalSpacing / 2
	centerRightSpacing := totalSpacing - leftCenterSpacing

	// Ensure at least 1 space
	if leftCenterSpacing < 1 {
		leftCenterSpacing = 1
	}
	if centerRightSpacing < 1 {
		centerRightSpacing = 1
	}

	// Build the header with calculated spacing
	spacer1 := strings.Repeat(" ", leftCenterSpacing)
	spacer2 := strings.Repeat(" ", centerRightSpacing)

	return left + spacer1 + center + spacer2 + right
}

// wrapHeader wraps the header content in a styled box.
//
// Parameters:
//   - content: The header content to wrap
//   - width: The total width
//
// Returns:
//   - string: Content wrapped in a styled border
func wrapHeader(content string, width int) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, false, true, false). // Bottom border only
		BorderForeground(colorAccent).
		Width(width-4). // Account for padding
		Padding(0, 1).
		MarginBottom(1)

	return style.Render(content)
}

// renderMinimalHeader creates a compact header for very narrow terminals.
// This is used when width < 40 characters.
//
// Parameters:
//   - screenName: The current screen name
//   - width: The available width
//
// Returns:
//   - string: A minimal header with CodeQuest and screen name
func renderMinimalHeader(screenName string, width int) string {
	style := lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		Border(lipgloss.RoundedBorder(), false, false, true, false).
		BorderForeground(colorAccent).
		Width(width-4).
		Padding(0, 1).
		MarginBottom(1)

	// Include "CodeQuest" in minimal header for consistency
	return style.Render("🎮 CodeQuest - " + screenName)
}
