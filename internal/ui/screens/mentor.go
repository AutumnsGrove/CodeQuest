// Package screens provides screen rendering functions for CodeQuest UI.
// This file implements the Mentor screen, providing AI-powered coding help and mentorship.
package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// Message represents a single message in the conversation history.
// Used to track conversation between user and AI mentor.
type Message struct {
	Sender    string    // "user", "crush", "mods", "claude", "system"
	Content   string    // Message text
	Timestamp time.Time // When sent
}

// AIProviderStatus represents the current status of AI providers.
type AIProviderStatus struct {
	Primary   string // Primary provider name (e.g., "Crush")
	Available bool   // Whether any provider is available
	Fallback  string // Fallback provider if primary unavailable
}

// RenderMentor renders the mentor screen with conversation history and input field.
// This is the main AI assistance screen where players can ask questions and get help.
//
// Features:
//   - Conversation history display (scrollable if needed)
//   - User messages (right-aligned, cyan theme)
//   - AI responses (left-aligned, lavender theme)
//   - System messages (centered, dim theme)
//   - Input field for typing questions
//   - AI provider status indicator
//   - Example prompts when conversation is empty
//   - Responsive layout for different terminal sizes
//   - Nil-safe: handles empty conversation and nil character gracefully
//
// Layout Structure:
//   - Header: Screen title with character info
//   - Status bar: AI provider status
//   - Conversation area: Message history (scrollable)
//   - Input area: Text input with placeholder and hints
//   - Footer: Key bindings and navigation help
//
// Parameters:
//   - character: Player character (for header display, nil-safe)
//   - inputText: Current text in the input field
//   - conversation: Slice of messages in conversation history
//   - width: Terminal width in characters
//   - height: Terminal height in characters
//
// Returns:
//   - string: Rendered mentor screen UI
func RenderMentor(character *game.Character, inputText string, conversation []Message, width, height int) string {
	// Render header
	header := renderMentorHeader(character, width)

	// Render AI provider status
	status := renderAIStatus(getAIProviderStatus(), width)

	// Render conversation or example prompts
	var conversationArea string
	if len(conversation) == 0 {
		conversationArea = renderExamplePrompts(width, height-20)
	} else {
		conversationArea = renderConversation(conversation, width, height-20)
	}

	// Render input area
	inputArea := renderInputArea(inputText, width)

	// Render footer with key bindings
	footer := renderMentorFooter(width)

	// Assemble screen
	screen := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		status,
		"",
		conversationArea,
		"",
		inputArea,
		"",
		footer,
	)

	return screen
}

// renderMentorHeader creates a header for the Mentor screen.
// Similar to Quest Board header but with Mentor branding.
func renderMentorHeader(char *game.Character, width int) string {
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
		return style.Render("🧙 AI Mentor")
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
	centerSection := centerStyle.Render("[AI Mentor]")

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

// renderAIStatus renders the AI provider status bar.
// Shows which AI provider is active and available.
func renderAIStatus(status AIProviderStatus, width int) string {
	var statusText string
	var statusColor lipgloss.Color

	if status.Available {
		statusText = fmt.Sprintf("🤖 AI Mentor - Powered by %s", status.Primary)
		statusColor = ColorSuccess
	} else {
		statusText = fmt.Sprintf("🤖 AI Mentor - Offline (will use %s when available)", status.Fallback)
		statusColor = ColorWarning
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(statusColor).
		Bold(true)

	statusLabel := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Render("Status: ")

	var statusIcon string
	if status.Available {
		statusIcon = lipgloss.NewStyle().Foreground(ColorSuccess).Render("✅ Online")
	} else {
		statusIcon = lipgloss.NewStyle().Foreground(ColorWarning).Render("⚠️  Offline")
	}

	// Build status line
	line1 := statusStyle.Render(statusText)
	line2 := statusLabel + statusIcon

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		line1,
		line2,
	)

	return content
}

// renderConversation renders the conversation history.
// Displays messages with appropriate styling based on sender.
func renderConversation(messages []Message, width, maxHeight int) string {
	title := SubtitleStyle.Render("Conversation")

	// Calculate available height for messages
	maxMessages := (maxHeight - 2) / 4 // Approximate: each message takes ~4 lines

	// If we have more messages than can fit, show only the most recent
	startIdx := 0
	if len(messages) > maxMessages {
		startIdx = len(messages) - maxMessages
	}

	// Render each message
	messageStrings := make([]string, 0)
	for i := startIdx; i < len(messages); i++ {
		msg := messages[i]
		rendered := renderMessage(msg, width-8)
		messageStrings = append(messageStrings, rendered)
	}

	// Join messages vertically
	messagesContent := lipgloss.JoinVertical(lipgloss.Left, messageStrings...)

	// Wrap in box
	boxContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		messagesContent,
	)

	return BoxStyle.Width(width - 4).Render(boxContent)
}

// renderMessage renders a single message based on sender type.
func renderMessage(msg Message, maxWidth int) string {
	timestamp := formatTime(msg.Timestamp)

	switch msg.Sender {
	case "user":
		return renderUserMessage(msg.Content, timestamp, maxWidth)
	case "crush", "mods", "claude":
		return renderAIMessage(msg.Sender, msg.Content, timestamp, maxWidth)
	case "system":
		return renderSystemMessage(msg.Content, timestamp, maxWidth)
	default:
		return renderSystemMessage(msg.Content, timestamp, maxWidth)
	}
}

// renderUserMessage renders a user message (right-aligned, cyan theme).
func renderUserMessage(content, timestamp string, maxWidth int) string {
	// Wrap content if too long
	wrappedContent := wrapText(content, maxWidth-20)

	// User message style (cyan, right-aligned)
	contentStyle := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Align(lipgloss.Right).
		Width(maxWidth)

	timeStyle := lipgloss.NewStyle().
		Foreground(ColorDim).
		Align(lipgloss.Right).
		Width(maxWidth)

	senderStyle := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Align(lipgloss.Right).
		Width(maxWidth)

	message := contentStyle.Render(wrappedContent)
	sender := senderStyle.Render("[You]")
	time := timeStyle.Render(timestamp)

	return lipgloss.JoinVertical(
		lipgloss.Right,
		message,
		sender,
		time,
		"", // Spacing
	)
}

// renderAIMessage renders an AI response (left-aligned, lavender theme).
func renderAIMessage(sender, content, timestamp string, maxWidth int) string {
	// Wrap content if too long
	wrappedContent := wrapText(content, maxWidth-20)

	// AI message style (lavender/purple, left-aligned)
	contentStyle := lipgloss.NewStyle().
		Foreground(ColorMagic).
		Width(maxWidth)

	timeStyle := lipgloss.NewStyle().
		Foreground(ColorDim).
		Width(maxWidth)

	// Capitalize sender name
	senderName := strings.Title(sender)
	senderStyle := lipgloss.NewStyle().
		Foreground(ColorMagic).
		Bold(true)

	message := contentStyle.Render(wrappedContent)
	senderLabel := senderStyle.Render(fmt.Sprintf("[%s]", senderName))
	time := timeStyle.Render(timestamp)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		senderLabel,
		message,
		time,
		"", // Spacing
	)
}

// renderSystemMessage renders a system message (centered, dim theme).
func renderSystemMessage(content, timestamp string, maxWidth int) string {
	// System message style (gray, centered)
	contentStyle := lipgloss.NewStyle().
		Foreground(ColorDim).
		Italic(true).
		Align(lipgloss.Center).
		Width(maxWidth)

	timeStyle := lipgloss.NewStyle().
		Foreground(ColorDim).
		Align(lipgloss.Center).
		Width(maxWidth)

	message := contentStyle.Render(content)
	time := timeStyle.Render(timestamp)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		message,
		time,
		"", // Spacing
	)
}

// renderExamplePrompts renders helpful example questions when conversation is empty.
func renderExamplePrompts(width, maxHeight int) string {
	title := SubtitleStyle.Render("Need inspiration? Try asking:")

	examples := []string{
		"• \"How can I optimize this function?\"",
		"• \"Explain what this error means\"",
		"• \"What's the best way to structure this code?\"",
		"• \"Help me debug this issue\"",
		"• \"Review my recent commit\"",
		"• \"Suggest improvements for my code\"",
	}

	exampleStyle := lipgloss.NewStyle().
		Foreground(ColorInfo).
		PaddingLeft(2)

	exampleStrings := make([]string, len(examples))
	for i, ex := range examples {
		exampleStrings[i] = exampleStyle.Render(ex)
	}

	examplesContent := lipgloss.JoinVertical(lipgloss.Left, exampleStrings...)

	welcomeStyle := lipgloss.NewStyle().
		Foreground(ColorMagic).
		Bold(true).
		Align(lipgloss.Center).
		Width(width - 8)

	welcome := welcomeStyle.Render("👋 Welcome to the AI Mentor!")

	descStyle := lipgloss.NewStyle().
		Foreground(ColorBright).
		Align(lipgloss.Center).
		Width(width - 8)

	description := descStyle.Render("Ask me anything about your code, quests, or coding in general.")

	// Combine content
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		welcome,
		"",
		description,
		"",
		title,
		"",
		examplesContent,
	)

	return BoxStyle.Width(width - 4).Render(content)
}

// renderInputArea renders the input field for typing questions.
func renderInputArea(inputText string, width int) string {
	title := SubtitleStyle.Render("Your Question")

	// Input field
	var displayText string
	if inputText == "" {
		displayText = DimTextStyle.Render("Ask anything about your code...")
	} else {
		displayText = TextStyle.Render(inputText + "█") // Cursor indicator
	}

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(0, 1).
		Width(width - 10)

	inputField := inputStyle.Render(displayText)

	// Character counter (placeholder for future limit)
	// counterStyle := lipgloss.NewStyle().
	// 	Foreground(ColorMuted).
	// 	Align(lipgloss.Right)
	// counter := counterStyle.Render(fmt.Sprintf("%d/%d", len(inputText), 500))

	// Hints
	hintStyle := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	hints := hintStyle.Render("[Enter]") + " Send  " +
		hintStyle.Render("[Esc]") + " Back to Dashboard"

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		inputField,
		"",
		hints,
	)

	return content
}

// renderMentorFooter renders the footer with key bindings.
func renderMentorFooter(width int) string {
	// Key bindings
	enterKey := renderKeybind("Enter", "Send Message")
	escKey := renderKeybind("Esc", "Back")
	ctrlC := renderKeybind("Ctrl+C", "Quit")

	keybinds := lipgloss.JoinHorizontal(
		lipgloss.Left,
		enterKey,
		"  ",
		escKey,
		"  ",
		ctrlC,
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

// getAIProviderStatus returns the current AI provider status.
// TODO: This will be implemented in Subagents 28-30 with actual provider logic.
// For now, returns a placeholder status.
func getAIProviderStatus() AIProviderStatus {
	// Placeholder: return "coming soon" status
	return AIProviderStatus{
		Primary:   "Crush",
		Available: false, // Will be true when AI is integrated
		Fallback:  "Mods (local)",
	}
}

// wrapText wraps text to fit within a specified width.
// Preserves existing line breaks and adds new breaks as needed.
func wrapText(text string, width int) string {
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

			// No need to track if last word
			_ = j
		}
	}

	return result.String()
}

// formatTime formats a timestamp for display in messages.
func formatTime(t time.Time) string {
	// If the message is from today, show just the time
	now := time.Now()
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return t.Format("3:04 PM")
	}

	// If from this year, show month and day
	if t.Year() == now.Year() {
		return t.Format("Jan 2, 3:04 PM")
	}

	// Otherwise show full date
	return t.Format("Jan 2, 2006 3:04 PM")
}
