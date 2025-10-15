// Package components provides reusable UI components for CodeQuest
// This file implements the modal dialog component for confirmations, notifications, and detailed views
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color constants (duplicated here to avoid import cycle with ui package)
var (
	modalColorPrimary = lipgloss.Color("205") // Pink/Magenta
	modalColorAccent  = lipgloss.Color("86")  // Cyan
	modalColorSuccess = lipgloss.Color("42")  // Green
	modalColorWarning = lipgloss.Color("214") // Orange
	modalColorError   = lipgloss.Color("196") // Red
	modalColorInfo    = lipgloss.Color("69")  // Blue
	modalColorDim     = lipgloss.Color("240") // Gray
	modalColorBright  = lipgloss.Color("15")  // White
)

// ModalType defines the type of modal dialog to display
type ModalType int

const (
	// ModalInfo displays informational messages with a blue/cyan theme
	ModalInfo ModalType = iota
	// ModalSuccess displays success confirmations with a green theme
	ModalSuccess
	// ModalWarning displays warning messages with an orange theme
	ModalWarning
	// ModalError displays error messages with a red theme
	ModalError
	// ModalConfirmation displays Yes/No confirmations with button hints
	ModalConfirmation
)

// ModalConfig holds configuration options for modal rendering
type ModalConfig struct {
	Title      string    // Title text displayed in the title bar
	Content    string    // Main content text (supports multi-line)
	Width      int       // Width of the modal (0 = auto-size)
	Height     int       // Height of the modal (0 = auto-size)
	ModalType  ModalType // Type of modal (affects styling and icon)
	TermWidth  int       // Terminal width for centering
	TermHeight int       // Terminal height for centering
}

// ModalIcons maps modal types to their display icons
var ModalIcons = map[ModalType]string{
	ModalInfo:         "ℹ️",
	ModalSuccess:      "✅",
	ModalWarning:      "⚠️",
	ModalError:        "❌",
	ModalConfirmation: "❓",
}

// RenderModal creates a modal dialog overlay for any screen.
// This is the main entry point for rendering modals with automatic sizing.
//
// The modal includes:
//   - Centered overlay with backdrop effect
//   - Title bar with type-specific icon and styling
//   - Multi-line content with word wrapping
//   - Footer with action hints based on modal type
//
// Parameters:
//   - title: The modal title (displayed in the title bar)
//   - content: The modal content (supports multi-line text)
//   - width: Desired width in characters (0 for auto-size based on content)
//   - height: Desired height in characters (0 for auto-size based on content)
//   - modalType: Type of modal (affects styling, icon, and action hints)
//
// Returns:
//   - string: The rendered modal with backdrop, ready to overlay on screen content
func RenderModal(title, content string, width, height int, modalType ModalType) string {
	config := ModalConfig{
		Title:      title,
		Content:    content,
		Width:      width,
		Height:     height,
		ModalType:  modalType,
		TermWidth:  80, // Default terminal width
		TermHeight: 24, // Default terminal height
	}

	return RenderModalWithConfig(config)
}

// RenderModalWithConfig renders a modal with full configuration control.
// This allows for precise control over all modal aspects including terminal dimensions.
//
// Parameters:
//   - config: Complete modal configuration
//
// Returns:
//   - string: The rendered modal with backdrop
func RenderModalWithConfig(config ModalConfig) string {
	// Auto-calculate dimensions if not specified
	if config.Width == 0 {
		config.Width = calculateModalWidth(config.Content)
	}
	if config.Height == 0 {
		config.Height = calculateModalHeight(config.Content, config.Width)
	}

	// Enforce minimum and maximum constraints
	config.Width = clamp(config.Width, 30, config.TermWidth-10)
	config.Height = clamp(config.Height, 8, config.TermHeight-6)

	// Build modal sections
	titleBar := renderModalTitleBar(config.Title, config.ModalType, config.Width)
	contentArea := renderModalContent(config.Content, config.Width, config.Height)
	footer := renderModalFooter(config.ModalType, config.Width)

	// Combine sections
	modalContent := titleBar + "\n" + contentArea + "\n" + footer

	// Apply modal styling with border
	styledModal := styleModalBox(modalContent, config.Width, config.ModalType)

	// Center the modal on screen with backdrop
	return renderModalWithBackdrop(styledModal, config.TermWidth, config.TermHeight)
}

// renderModalTitleBar creates the title bar section of the modal.
// Includes the modal type icon and title text with appropriate styling.
//
// Parameters:
//   - title: The title text
//   - modalType: Type of modal (determines icon and color)
//   - width: Width for the title bar
//
// Returns:
//   - string: Styled title bar
func renderModalTitleBar(title string, modalType ModalType, width int) string {
	// Get the icon for this modal type
	icon := ModalIcons[modalType]

	// Get the color for this modal type
	color := getModalColor(modalType)

	// Create title style
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(color).
		Width(width - 4) // Account for padding

	titleText := icon + " " + title
	return titleStyle.Render(titleText)
}

// renderModalContent creates the content area of the modal with word wrapping.
// Supports multi-line content and ensures proper text flow.
//
// Parameters:
//   - content: The content text (can be multi-line)
//   - width: Available width for content
//   - height: Available height for content
//
// Returns:
//   - string: Formatted content with proper wrapping
func renderModalContent(content string, width, height int) string {
	// Word wrap the content to fit width
	wrappedContent := wordWrap(content, width-6) // Account for padding and borders

	// Split into lines
	lines := strings.Split(wrappedContent, "\n")

	// Limit to available height
	maxLines := height - 4 // Reserve space for title and footer
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		// Add ellipsis to indicate truncation
		if len(lines) > 0 {
			lines[len(lines)-1] += "..."
		}
	}

	// Style the content
	contentStyle := lipgloss.NewStyle().
		Foreground(modalColorBright).
		Width(width - 6).
		Align(lipgloss.Left)

	return contentStyle.Render(strings.Join(lines, "\n"))
}

// renderModalFooter creates the footer section with action hints.
// The footer content depends on the modal type.
//
// Parameters:
//   - modalType: Type of modal (determines action hints)
//   - width: Width for the footer
//
// Returns:
//   - string: Styled footer with action hints
func renderModalFooter(modalType ModalType, width int) string {
	var actionHint string

	switch modalType {
	case ModalConfirmation:
		actionHint = "[Y] Yes  [N] No"
	case ModalError, ModalWarning:
		actionHint = "[Enter] OK  [Esc] Close"
	default:
		actionHint = "[Enter] Continue"
	}

	// Create separator line
	separator := strings.Repeat("─", width-4)
	separatorStyle := lipgloss.NewStyle().
		Foreground(modalColorDim)

	// Style the action hint
	hintStyle := lipgloss.NewStyle().
		Foreground(modalColorAccent).
		Bold(true).
		Width(width - 4).
		Align(lipgloss.Center)

	return separatorStyle.Render(separator) + "\n" + hintStyle.Render(actionHint)
}

// styleModalBox applies the border and styling to the modal content.
// Different modal types get different border colors.
//
// Parameters:
//   - content: The modal content to style
//   - width: Width of the modal
//   - modalType: Type of modal (affects border color)
//
// Returns:
//   - string: Modal content with border and styling applied
func styleModalBox(content string, width int, modalType ModalType) string {
	borderColor := getModalColor(modalType)

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(width)

	return style.Render(content)
}

// renderModalWithBackdrop centers the modal on screen with a backdrop effect.
// Creates the overlay appearance by adding padding/margins.
//
// Parameters:
//   - modal: The styled modal content
//   - termWidth: Terminal width
//   - termHeight: Terminal height
//
// Returns:
//   - string: Modal centered on screen with backdrop
func renderModalWithBackdrop(modal string, termWidth, termHeight int) string {
	// Calculate vertical spacing to center modal
	modalHeight := lipgloss.Height(modal)
	topPadding := (termHeight - modalHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Add backdrop effect with dimmed background indication
	backdropLine := strings.Repeat("░", termWidth)
	backdropStyle := lipgloss.NewStyle().
		Foreground(modalColorDim).
		Faint(true)

	// Build backdrop with modal centered
	var result strings.Builder

	// Top backdrop lines
	for i := 0; i < topPadding; i++ {
		result.WriteString(backdropStyle.Render(backdropLine) + "\n")
	}

	// Centered modal
	centeredModal := lipgloss.Place(termWidth, modalHeight, lipgloss.Center, lipgloss.Top, modal)
	result.WriteString(centeredModal)

	// Bottom backdrop lines
	bottomPadding := termHeight - topPadding - modalHeight
	for i := 0; i < bottomPadding; i++ {
		result.WriteString("\n" + backdropStyle.Render(backdropLine))
	}

	return result.String()
}

// getModalColor returns the primary color for a modal type.
// Used for borders, titles, and accents.
//
// Parameters:
//   - modalType: The type of modal
//
// Returns:
//   - lipgloss.Color: The color to use for this modal type
func getModalColor(modalType ModalType) lipgloss.Color {
	switch modalType {
	case ModalInfo:
		return modalColorInfo
	case ModalSuccess:
		return modalColorSuccess
	case ModalWarning:
		return modalColorWarning
	case ModalError:
		return modalColorError
	case ModalConfirmation:
		return modalColorAccent
	default:
		return modalColorPrimary
	}
}

// ============================================================================
// Helper Modal Rendering Functions
// ============================================================================

// RenderInfoModal creates an informational modal with default sizing.
// Convenience wrapper for ModalInfo type.
//
// Parameters:
//   - title: Modal title
//   - content: Modal content
//
// Returns:
//   - string: Rendered info modal
func RenderInfoModal(title, content string) string {
	return RenderModal(title, content, 0, 0, ModalInfo)
}

// RenderSuccessModal creates a success confirmation modal with default sizing.
// Convenience wrapper for ModalSuccess type.
//
// Parameters:
//   - title: Modal title
//   - content: Modal content
//
// Returns:
//   - string: Rendered success modal
func RenderSuccessModal(title, content string) string {
	return RenderModal(title, content, 0, 0, ModalSuccess)
}

// RenderWarningModal creates a warning modal with default sizing.
// Convenience wrapper for ModalWarning type.
//
// Parameters:
//   - title: Modal title
//   - content: Modal content
//
// Returns:
//   - string: Rendered warning modal
func RenderWarningModal(title, content string) string {
	return RenderModal(title, content, 0, 0, ModalWarning)
}

// RenderErrorModal creates an error modal with default sizing.
// Convenience wrapper for ModalError type.
//
// Parameters:
//   - title: Modal title
//   - content: Modal content
//
// Returns:
//   - string: Rendered error modal
func RenderErrorModal(title, content string) string {
	return RenderModal(title, content, 0, 0, ModalError)
}

// RenderConfirmModal creates a confirmation modal with Yes/No options.
// Convenience wrapper for ModalConfirmation type.
//
// Parameters:
//   - title: Modal title
//   - content: Modal content
//
// Returns:
//   - string: Rendered confirmation modal
func RenderConfirmModal(title, content string) string {
	return RenderModal(title, content, 0, 0, ModalConfirmation)
}

// RenderLevelUpModal creates a special level-up notification modal.
// This is a themed modal for character level progression.
//
// Parameters:
//   - level: The new level reached
//   - abilities: Slice of new abilities unlocked
//
// Returns:
//   - string: Rendered level-up modal
func RenderLevelUpModal(level int, abilities []string) string {
	title := "Level Up!"

	content := fmt.Sprintf("You've reached Level %d!\n\n", level)

	if len(abilities) > 0 {
		content += "New abilities unlocked:\n"
		for _, ability := range abilities {
			content += fmt.Sprintf("- %s\n", ability)
		}
	} else {
		content += "Keep up the great work!"
	}

	return RenderSuccessModal(title, content)
}

// RenderQuestDetailModal creates a detailed quest information modal.
// Displays quest information in a formatted modal.
//
// Parameters:
//   - questTitle: The quest title
//   - description: Quest description
//   - progress: Current progress (e.g., "3/5 commits")
//   - reward: XP reward amount
//
// Returns:
//   - string: Rendered quest detail modal
func RenderQuestDetailModal(questTitle, description, progress string, reward int) string {
	content := fmt.Sprintf("%s\n\nProgress: %s\nReward: %d XP",
		description, progress, reward)

	return RenderInfoModal(questTitle, content)
}

// ============================================================================
// Utility Functions
// ============================================================================

// calculateModalWidth determines optimal width based on content.
// Analyzes content to find appropriate width for readability.
//
// Parameters:
//   - content: The content to analyze
//
// Returns:
//   - int: Recommended width in characters
func calculateModalWidth(content string) int {
	lines := strings.Split(content, "\n")
	maxLen := 0

	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Add padding and reasonable constraints
	width := maxLen + 10
	if width < 40 {
		width = 40
	}
	if width > 70 {
		width = 70
	}

	return width
}

// calculateModalHeight determines optimal height based on content and width.
// Accounts for word wrapping to calculate required height.
//
// Parameters:
//   - content: The content to analyze
//   - width: The width constraint
//
// Returns:
//   - int: Recommended height in characters
func calculateModalHeight(content string, width int) int {
	wrapped := wordWrap(content, width-6)
	lines := strings.Split(wrapped, "\n")

	// Content height + title + footer + padding
	height := len(lines) + 6

	if height < 10 {
		height = 10
	}
	if height > 20 {
		height = 20
	}

	return height
}

// wordWrap wraps text to fit within a specified width.
// Preserves existing line breaks and adds new breaks as needed.
//
// Parameters:
//   - text: The text to wrap
//   - width: Maximum width per line
//
// Returns:
//   - string: Word-wrapped text
func wordWrap(text string, width int) string {
	if width <= 0 {
		width = 40 // Minimum reasonable width
	}

	var result strings.Builder
	paragraphs := strings.Split(text, "\n")

	for i, paragraph := range paragraphs {
		if i > 0 {
			result.WriteString("\n")
		}

		// Skip empty lines
		if strings.TrimSpace(paragraph) == "" {
			continue
		}

		words := strings.Fields(paragraph)
		lineLength := 0

		for j, word := range words {
			wordLen := len(word)

			// Check if adding this word exceeds width
			if lineLength+wordLen > width {
				if lineLength > 0 {
					result.WriteString("\n")
					lineLength = 0
				}
			}

			// Add space before word (except at line start)
			if lineLength > 0 {
				result.WriteString(" ")
				lineLength++
			}

			result.WriteString(word)
			lineLength += wordLen

			// Add space after word if not the last word
			if j < len(words)-1 && lineLength < width {
				// Space will be added before next word
			}
		}
	}

	return result.String()
}

// clamp restricts a value to a specified range.
// Ensures value is between min and max inclusive.
//
// Parameters:
//   - value: The value to clamp
//   - min: Minimum allowed value
//   - max: Maximum allowed value
//
// Returns:
//   - int: Clamped value
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
