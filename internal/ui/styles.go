// Package ui provides the terminal user interface styling for CodeQuest
// This file defines a comprehensive Lip Gloss styling system for consistent UI appearance
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ============================================================================
// Color Palette
// ============================================================================
// Colors chosen for both light and dark terminal compatibility with RPG aesthetic

var (
	// Primary colors - Main brand and emphasis
	ColorPrimary   = lipgloss.Color("205") // Pink/Magenta - Main accent
	ColorSecondary = lipgloss.Color("63")  // Purple - Secondary accent

	// Accent colors - Highlights and focus
	ColorAccent   = lipgloss.Color("86")  // Cyan - Interactive elements
	ColorAccentAlt = lipgloss.Color("39") // Bright Cyan - Hover states

	// Status colors - Semantic meaning
	ColorSuccess = lipgloss.Color("42")  // Green - Completed, success
	ColorWarning = lipgloss.Color("214") // Orange - Warnings, attention
	ColorError   = lipgloss.Color("196") // Red - Errors, failed
	ColorInfo    = lipgloss.Color("69")  // Blue - Information

	// Neutral colors - Text and backgrounds
	ColorDim    = lipgloss.Color("240") // Gray - Dim/inactive text
	ColorBright = lipgloss.Color("15")  // White - Bright text
	ColorMuted  = lipgloss.Color("243") // Light gray - Secondary text

	// RPG-specific colors
	ColorXP     = lipgloss.Color("226") // Gold/Yellow - XP and rewards
	ColorLevel  = lipgloss.Color("93")  // Yellow-Orange - Level indicators
	ColorQuest  = lipgloss.Color("111") // Light Blue - Quest markers
	ColorMagic  = lipgloss.Color("177") // Lavender - Magic/special effects
)

// ============================================================================
// Common Text Styles
// ============================================================================

var (
	// TitleStyle - Large titles and headers (screens, sections)
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1).
			Padding(0, 1)

	// SubtitleStyle - Section headers and subtitles
	SubtitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			MarginTop(1).
			MarginBottom(0)

	// HeadingStyle - Small headings within sections
	HeadingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	// TextStyle - Normal body text
	TextStyle = lipgloss.NewStyle().
			Foreground(ColorBright)

	// BoldTextStyle - Emphasized text
	BoldTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBright)

	// DimTextStyle - De-emphasized, secondary text
	DimTextStyle = lipgloss.NewStyle().
			Foreground(ColorDim)

	// MutedTextStyle - Muted text for labels and hints
	MutedTextStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)

	// ErrorTextStyle - Error messages
	ErrorTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorError)

	// SuccessTextStyle - Success messages
	SuccessTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSuccess)

	// WarningTextStyle - Warning messages
	WarningTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorWarning)

	// InfoTextStyle - Informational messages
	InfoTextStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)
)

// ============================================================================
// Border and Container Styles
// ============================================================================

var (
	// BoxStyle - Standard box with rounded border
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(1, 2)

	// BoxStyleFocused - Box style when focused/selected
	BoxStyleFocused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			Bold(true)

	// BoxStyleDim - Box style for inactive/background elements
	BoxStyleDim = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDim).
			Padding(1, 2)

	// NormalBorderBox - Box with normal square borders
	NormalBorderBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 2)

	// ThickBorderBox - Box with thick borders for emphasis
	ThickBorderBox = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	// DoubleBorderBox - Box with double borders for special sections
	DoubleBorderBox = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorAccent).
			Padding(1, 2)

	// PanelStyle - Panel for grouping content without heavy borders
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), false, false, false, true). // Left border only
			BorderForeground(ColorAccent).
			PaddingLeft(2).
			MarginBottom(1)
)

// ============================================================================
// Progress Bar Styles
// ============================================================================

var (
	// XPBarStyle - Style for experience point progress bars
	XPBarStyle = lipgloss.NewStyle().
			Foreground(ColorXP).
			Bold(true)

	// QuestProgressBarStyle - Style for quest progress bars
	QuestProgressBarStyle = lipgloss.NewStyle().
				Foreground(ColorQuest).
				Bold(true)

	// HealthBarStyle - Style for health/status bars
	HealthBarStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	// ProgressBarEmptyStyle - Style for empty portion of progress bars
	ProgressBarEmptyStyle = lipgloss.NewStyle().
				Foreground(ColorDim)
)

// ============================================================================
// Status Indicator Styles
// ============================================================================

var (
	// StatusActiveStyle - Active/in-progress quest status
	StatusActiveStyle = lipgloss.NewStyle().
				Foreground(ColorSuccess).
				Bold(true)

	// StatusCompletedStyle - Completed quest status
	StatusCompletedStyle = lipgloss.NewStyle().
				Foreground(ColorInfo).
				Bold(true)

	// StatusFailedStyle - Failed quest status
	StatusFailedStyle = lipgloss.NewStyle().
				Foreground(ColorError).
				Bold(true)

	// StatusPendingStyle - Pending/available quest status
	StatusPendingStyle = lipgloss.NewStyle().
				Foreground(ColorWarning).
				Bold(true)

	// StatusLockedStyle - Locked/unavailable quest status
	StatusLockedStyle = lipgloss.NewStyle().
				Foreground(ColorDim).
				Strikethrough(true)
)

// ============================================================================
// Interactive Element Styles
// ============================================================================

var (
	// ButtonStyle - Default button style
	ButtonStyle = lipgloss.NewStyle().
			Foreground(ColorBright).
			Background(ColorSecondary).
			Padding(0, 2).
			MarginRight(1)

	// ButtonFocusedStyle - Button style when focused
	ButtonFocusedStyle = lipgloss.NewStyle().
				Foreground(ColorBright).
				Background(ColorPrimary).
				Padding(0, 2).
				MarginRight(1).
				Bold(true)

	// InputStyle - Text input field style
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(0, 1)

	// InputFocusedStyle - Text input when focused
	InputFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(0, 1).
				Bold(true)

	// SelectedItemStyle - Selected list item
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true).
				PaddingLeft(2)

	// UnselectedItemStyle - Unselected list item
	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorBright).
				PaddingLeft(2)
)

// ============================================================================
// Special UI Element Styles
// ============================================================================

var (
	// KeybindStyle - Keybind hint style (e.g., [Q] for Quit)
	KeybindStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	// KeybindDescStyle - Keybind description style
	KeybindDescStyle = lipgloss.NewStyle().
				Foreground(ColorBright)

	// StatLabelStyle - Style for stat labels (HP, XP, Level, etc.)
	StatLabelStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Bold(true)

	// StatValueStyle - Style for stat values
	StatValueStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	// TimerStyle - Session timer display
	TimerStyle = lipgloss.NewStyle().
			Foreground(ColorInfo).
			Bold(true)

	// NotificationStyle - Notification/toast messages
	NotificationStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorSuccess).
				Foreground(ColorBright).
				Padding(0, 2).
				MarginTop(1)

	// ModalBackdropStyle - Modal dialog backdrop/overlay
	ModalBackdropStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("0")).
				Foreground(ColorDim)

	// ModalStyle - Modal dialog box
	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(ColorPrimary).
			Background(lipgloss.Color("235")). // Dark background
			Foreground(ColorBright).
			Padding(2, 4).
			Width(60)
)

// ============================================================================
// Helper Functions - Rendering Utilities
// ============================================================================

// RenderTitle renders a styled title with optional icon
func RenderTitle(text string, icon string) string {
	if icon != "" {
		return TitleStyle.Render(icon + " " + text)
	}
	return TitleStyle.Render(text)
}

// RenderSubtitle renders a styled subtitle
func RenderSubtitle(text string) string {
	return SubtitleStyle.Render(text)
}

// RenderKeybind formats a keybind hint (e.g., "[Q] Quit")
func RenderKeybind(key, description string) string {
	return KeybindStyle.Render("["+key+"]") + " " + KeybindDescStyle.Render(description)
}

// RenderStat formats a stat with label and value (e.g., "Level: 5")
func RenderStat(label, value string) string {
	return StatLabelStyle.Render(label+": ") + StatValueStyle.Render(value)
}

// RenderStatus renders a status indicator with appropriate styling
func RenderStatus(status string) string {
	switch strings.ToLower(status) {
	case "active", "in-progress", "ongoing":
		return StatusActiveStyle.Render("● " + status)
	case "completed", "done", "finished":
		return StatusCompletedStyle.Render("✓ " + status)
	case "failed", "error":
		return StatusFailedStyle.Render("✗ " + status)
	case "pending", "available":
		return StatusPendingStyle.Render("○ " + status)
	case "locked", "unavailable":
		return StatusLockedStyle.Render("🔒 " + status)
	default:
		return TextStyle.Render(status)
	}
}

// ============================================================================
// Progress Bar Rendering
// ============================================================================

// RenderProgressBar creates a progress bar with percentage
// width: total width of the bar in characters
// current: current progress value
// total: total/max value
// barType: "xp" for XP bars, "quest" for quest progress, "health" for health bars
func RenderProgressBar(current, total, width int, barType string) string {
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
	case "health":
		filledStyle = HealthBarStyle
		emptyStyle = ProgressBarEmptyStyle
		fillChar = "♥"
		emptyChar = "♡"
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

// RenderXPBar renders an XP progress bar with level information
func RenderXPBar(currentXP, xpToNextLevel, level int, width int) string {
	label := StatLabelStyle.Render(fmt.Sprintf("Level %d ", level))
	bar := RenderProgressBar(currentXP, xpToNextLevel, width-20, "xp")
	return label + bar
}

// ============================================================================
// Box Rendering with Title
// ============================================================================

// BoxWithTitle wraps content in a styled box with a title header
func BoxWithTitle(title, content string, focused bool) string {
	style := BoxStyle
	if focused {
		style = BoxStyleFocused
	}

	titleStyled := HeadingStyle.Render("┤ " + title + " ├")
	boxContent := titleStyled + "\n\n" + content

	return style.Render(boxContent)
}

// BoxWithTitleAndIcon wraps content with title and icon
func BoxWithTitleAndIcon(title, icon, content string, focused bool) string {
	style := BoxStyle
	if focused {
		style = BoxStyleFocused
	}

	titleStyled := HeadingStyle.Render("┤ " + icon + " " + title + " ├")
	boxContent := titleStyled + "\n\n" + content

	return style.Render(boxContent)
}

// ============================================================================
// Layout Helpers
// ============================================================================

// JoinHorizontal joins strings horizontally with optional spacing
func JoinHorizontal(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}

// JoinVertical joins strings vertically
func JoinVertical(strs ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, strs...)
}

// PlaceInCenter centers content in a given width and height
func PlaceInCenter(width, height int, content string) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// ============================================================================
// Responsive Sizing Helpers
// ============================================================================

// AdaptiveWidth returns a width-adjusted style
func AdaptiveWidth(style lipgloss.Style, width int) lipgloss.Style {
	return style.Width(width)
}

// AdaptiveHeight returns a height-adjusted style
func AdaptiveHeight(style lipgloss.Style, height int) lipgloss.Style {
	return style.Height(height)
}

// ============================================================================
// Special Effects
// ============================================================================

// RenderLevelUpMessage creates a special styled level up notification
func RenderLevelUpMessage(level int) string {
	message := fmt.Sprintf("⚡ LEVEL UP! ⚡\n\nYou are now Level %d", level)

	style := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorXP).
		Foreground(ColorXP).
		Bold(true).
		Padding(2, 4).
		Align(lipgloss.Center)

	return style.Render(message)
}

// RenderQuestCompleteMessage creates a quest completion notification
func RenderQuestCompleteMessage(questTitle string, xpReward int) string {
	message := fmt.Sprintf("✓ QUEST COMPLETE!\n\n%s\n\n+%d XP", questTitle, xpReward)

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSuccess).
		Foreground(ColorSuccess).
		Bold(true).
		Padding(2, 4).
		Align(lipgloss.Center)

	return style.Render(message)
}

// RenderLoadingSpinner returns a styled loading indicator
func RenderLoadingSpinner(frame int, text string) string {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinner := spinners[frame%len(spinners)]

	spinnerStyle := lipgloss.NewStyle().Foreground(ColorAccent)
	return spinnerStyle.Render(spinner) + " " + MutedTextStyle.Render(text)
}
