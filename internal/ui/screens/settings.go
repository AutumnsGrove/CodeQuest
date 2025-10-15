// Package screens provides screen rendering functions for CodeQuest UI.
// This file implements the Settings screen, displaying configuration options.
package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// SettingsCategory represents different categories of settings.
type SettingsCategory int

const (
	// CategoryGame represents game-related settings
	CategoryGame SettingsCategory = iota
	// CategoryUI represents UI/display settings
	CategoryUI
	// CategoryAI represents AI provider settings
	CategoryAI
	// CategoryGit represents Git integration settings
	CategoryGit
	// CategoryDebug represents debug/developer settings
	CategoryDebug
)

// RenderSettings renders the complete settings screen.
// This screen displays all configuration options in a read-only format.
//
// The settings screen shows:
//   - Header with character info
//   - Settings categories (Game, UI, AI, Git, Debug)
//   - Current values for all configuration options
//   - Navigation hints for changing settings (future feature)
//
// Layout Structure:
//   - Header: Screen title with character info
//   - Main panel: All settings grouped by category
//   - Footer: Key bindings
//
// Note: Settings modification is not yet implemented (read-only for MVP).
// Future versions will add interactive editing capabilities.
//
// Parameters:
//   - character: Player character (for header display)
//   - width: Terminal width in characters
//   - height: Terminal height in characters
//
// Returns:
//   - string: Rendered settings screen UI
func RenderSettings(character *game.Character, width, height int) string {
	// Render header (inline to avoid import cycle)
	header := renderSettingsHeader(character, width)

	// Render main settings panel
	settingsPanel := renderSettingsPanel(width)

	// Render footer with key bindings
	footer := renderSettingsFooter(width)

	// Assemble screen
	screen := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		settingsPanel,
		"",
		footer,
	)

	return screen
}

// renderSettingsPanel renders the main settings panel with all categories.
func renderSettingsPanel(width int) string {
	sections := make([]string, 0)

	// Game Settings Section
	gameSection := renderGameSettings()
	sections = append(sections, gameSection)

	// UI Settings Section
	uiSection := renderUISettings()
	sections = append(sections, uiSection)

	// AI Settings Section
	aiSection := renderAISettings()
	sections = append(sections, aiSection)

	// Git Settings Section
	gitSection := renderGitSettings()
	sections = append(sections, gitSection)

	// Debug Settings Section
	debugSection := renderDebugSettings()
	sections = append(sections, debugSection)

	// Join all sections with spacing
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		sections...,
	)

	// Wrap in box
	return BoxStyle.Width(width - 4).Render(content)
}

// ============================================================================
// Settings Category Rendering Functions
// ============================================================================

// renderGameSettings renders game-related settings.
func renderGameSettings() string {
	title := SubtitleStyle.Render("🎮 Game Settings")

	// Difficulty setting (placeholder values)
	difficultyLabel := StatLabelStyle.Render("Difficulty: ")
	difficultyValue := StatValueStyle.Render("Normal")
	difficulty := difficultyLabel + difficultyValue

	// XP multiplier setting
	xpMultiplierLabel := StatLabelStyle.Render("XP Multiplier: ")
	xpMultiplierValue := StatValueStyle.Render("1.0x")
	xpMultiplier := xpMultiplierLabel + xpMultiplierValue

	// Auto-save setting
	autoSaveLabel := StatLabelStyle.Render("Auto-save: ")
	autoSaveValue := SuccessTextStyle.Render("Enabled ✓")
	autoSave := autoSaveLabel + autoSaveValue

	// Quest notifications
	questNotifLabel := StatLabelStyle.Render("Quest Notifications: ")
	questNotifValue := SuccessTextStyle.Render("Enabled ✓")
	questNotif := questNotifLabel + questNotifValue

	hint := MutedTextStyle.Render("  (Game settings can be modified in config file)")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		difficulty,
		xpMultiplier,
		autoSave,
		questNotif,
		"",
		hint,
	)
}

// renderUISettings renders UI/display settings.
func renderUISettings() string {
	title := SubtitleStyle.Render("🎨 UI Settings")

	// Theme setting
	themeLabel := StatLabelStyle.Render("Color Theme: ")
	themeValue := StatValueStyle.Render("Default (Charmbracelet)")
	theme := themeLabel + themeValue

	// Animations setting
	animationsLabel := StatLabelStyle.Render("Animations: ")
	animationsValue := SuccessTextStyle.Render("Enabled ✓")
	animations := animationsLabel + animationsValue

	// Compact mode
	compactLabel := StatLabelStyle.Render("Compact Mode: ")
	compactValue := DimTextStyle.Render("Disabled")
	compact := compactLabel + compactValue

	// Show help hints
	helpHintsLabel := StatLabelStyle.Render("Show Help Hints: ")
	helpHintsValue := SuccessTextStyle.Render("Enabled ✓")
	helpHints := helpHintsLabel + helpHintsValue

	hint := MutedTextStyle.Render("  (UI customization coming in future updates)")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		theme,
		animations,
		compact,
		helpHints,
		"",
		hint,
	)
}

// renderAISettings renders AI provider settings.
func renderAISettings() string {
	title := SubtitleStyle.Render("🤖 AI Settings")

	// Primary provider
	primaryLabel := StatLabelStyle.Render("Primary Provider: ")
	primaryValue := StatValueStyle.Render("Crush (Anthropic)")
	primary := primaryLabel + primaryValue

	// Provider status
	statusLabel := StatLabelStyle.Render("Status: ")
	statusValue := SuccessTextStyle.Render("✓ Online")
	status := statusLabel + statusValue

	// Fallback provider
	fallbackLabel := StatLabelStyle.Render("Fallback Provider: ")
	fallbackValue := InfoTextStyle.Render("Mods (local)")
	fallback := fallbackLabel + fallbackValue

	// Rate limiting
	rateLimitLabel := StatLabelStyle.Render("Rate Limiting: ")
	rateLimitValue := SuccessTextStyle.Render("Enabled ✓")
	rateLimit := rateLimitLabel + rateLimitValue

	// Auto-mentor
	autoMentorLabel := StatLabelStyle.Render("Auto-mentor on errors: ")
	autoMentorValue := DimTextStyle.Render("Disabled")
	autoMentor := autoMentorLabel + autoMentorValue

	hint := MutedTextStyle.Render("  (API keys configured in secrets.json)")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		primary,
		status,
		fallback,
		rateLimit,
		autoMentor,
		"",
		hint,
	)
}

// renderGitSettings renders Git integration settings.
func renderGitSettings() string {
	title := SubtitleStyle.Render("📁 Git Settings")

	// Auto-detect commits
	autoDetectLabel := StatLabelStyle.Render("Auto-detect commits: ")
	autoDetectValue := SuccessTextStyle.Render("Enabled ✓")
	autoDetect := autoDetectLabel + autoDetectValue

	// Repository path (placeholder)
	repoLabel := StatLabelStyle.Render("Repository: ")
	repoValue := DimTextStyle.Render("Auto-detected from current directory")
	repo := repoLabel + repoValue

	// Watch mode
	watchLabel := StatLabelStyle.Render("Watch mode: ")
	watchValue := SuccessTextStyle.Render("Active ✓")
	watch := watchLabel + watchValue

	// Commit XP calculation
	commitXPLabel := StatLabelStyle.Render("Commit XP formula: ")
	commitXPValue := StatValueStyle.Render("Lines-based with quality bonus")
	commitXP := commitXPLabel + commitXPValue

	hint := MutedTextStyle.Render("  (Git watcher runs in background)")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		autoDetect,
		repo,
		watch,
		commitXP,
		"",
		hint,
	)
}

// renderDebugSettings renders debug/developer settings.
func renderDebugSettings() string {
	title := SubtitleStyle.Render("🐛 Debug Settings")

	// Log level
	logLevelLabel := StatLabelStyle.Render("Log Level: ")
	logLevelValue := InfoTextStyle.Render("Info")
	logLevel := logLevelLabel + logLevelValue

	// Dev mode
	devModeLabel := StatLabelStyle.Render("Developer Mode: ")
	devModeValue := DimTextStyle.Render("Disabled")
	devMode := devModeLabel + devModeValue

	// Debug UI
	debugUILabel := StatLabelStyle.Render("Show Debug UI: ")
	debugUIValue := DimTextStyle.Render("Disabled")
	debugUI := debugUILabel + debugUIValue

	// Performance monitoring
	perfLabel := StatLabelStyle.Render("Performance Monitor: ")
	perfValue := DimTextStyle.Render("Disabled")
	perf := perfLabel + perfValue

	hint := MutedTextStyle.Render("  (Enable dev mode with CODEQUEST_DEBUG=1)")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		logLevel,
		devMode,
		debugUI,
		perf,
		"",
		hint,
	)
}

// renderSettingsFooter renders the footer with key bindings.
func renderSettingsFooter(width int) string {
	// Info message about settings modification
	infoMsg := InfoTextStyle.Render("ℹ️  Settings are currently read-only. Editing will be added in a future update.")

	// Key bindings
	dashboard := renderKeybind("Alt+Q", "Dashboard")
	esc := renderKeybind("Esc", "Back")
	help := renderKeybind("?", "Help")
	save := renderKeybind("Ctrl+S", "Save (future)")

	keybinds := lipgloss.JoinHorizontal(
		lipgloss.Left,
		dashboard,
		"  ",
		esc,
		"  ",
		help,
		"  ",
		save,
	)

	// Combine info and keybinds
	footer := lipgloss.JoinVertical(
		lipgloss.Center,
		infoMsg,
		"",
		keybinds,
	)

	// Center the footer
	centered := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(footer)

	return centered
}

// ============================================================================
// Future: Interactive Settings Functions (Post-MVP)
// ============================================================================

// These functions are stubs for future implementation when settings
// modification is added in a later phase.

// SettingItem represents a configurable setting item.
// This will be used for interactive settings editing in the future.
type SettingItem struct {
	Label       string           // Display label
	Key         string           // Config key
	Value       interface{}      // Current value
	ValueType   string           // "bool", "string", "int", "float", "choice"
	Choices     []string         // For "choice" type
	Description string           // Help text
	Category    SettingsCategory // Which category this belongs to
}

// Example settings structure for future use:
var settingsItems = []SettingItem{
	{
		Label:       "Difficulty",
		Key:         "game.difficulty",
		Value:       "Normal",
		ValueType:   "choice",
		Choices:     []string{"Easy", "Normal", "Hard", "Expert"},
		Description: "Game difficulty affects XP requirements and quest complexity",
		Category:    CategoryGame,
	},
	{
		Label:       "XP Multiplier",
		Key:         "game.xp_multiplier",
		Value:       1.0,
		ValueType:   "float",
		Description: "Multiplier for all XP gains",
		Category:    CategoryGame,
	},
	{
		Label:       "Auto-save",
		Key:         "game.auto_save",
		Value:       true,
		ValueType:   "bool",
		Description: "Automatically save progress",
		Category:    CategoryGame,
	},
	{
		Label:       "Color Theme",
		Key:         "ui.theme",
		Value:       "default",
		ValueType:   "choice",
		Choices:     []string{"default", "dark", "light", "solarized"},
		Description: "UI color theme",
		Category:    CategoryUI,
	},
	{
		Label:       "Animations",
		Key:         "ui.animations",
		Value:       true,
		ValueType:   "bool",
		Description: "Enable UI animations and transitions",
		Category:    CategoryUI,
	},
	{
		Label:       "Primary AI Provider",
		Key:         "ai.primary_provider",
		Value:       "crush",
		ValueType:   "choice",
		Choices:     []string{"crush", "mods", "claude-code", "none"},
		Description: "Primary AI assistant provider",
		Category:    CategoryAI,
	},
	{
		Label:       "Auto-detect commits",
		Key:         "git.auto_detect",
		Value:       true,
		ValueType:   "bool",
		Description: "Automatically detect and track Git commits",
		Category:    CategoryGit,
	},
	{
		Label:       "Log Level",
		Key:         "debug.log_level",
		Value:       "info",
		ValueType:   "choice",
		Choices:     []string{"debug", "info", "warn", "error"},
		Description: "Logging verbosity level",
		Category:    CategoryDebug,
	},
}

// getSettingsByCategory returns all settings for a given category.
// This will be used for interactive settings screens in the future.
func getSettingsByCategory(category SettingsCategory) []SettingItem {
	filtered := make([]SettingItem, 0)
	for _, item := range settingsItems {
		if item.Category == category {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// renderSettingItem renders a single setting item for interactive editing.
// This is a stub for future implementation.
func renderSettingItem(item SettingItem, selected bool) string {
	// Future implementation will render editable setting items
	// For now, just return a placeholder
	indicator := "  "
	if selected {
		indicator = "▶ "
	}

	label := StatLabelStyle.Render(item.Label + ": ")
	value := fmt.Sprintf("%v", item.Value)

	return indicator + label + value
}

// renderSettingsHeader creates a header for the Settings screen (inline to avoid import cycle).
func renderSettingsHeader(char *game.Character, width int) string {
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
		return style.Render("🎮 Settings")
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
	centerSection := centerStyle.Render("[Settings]")

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
