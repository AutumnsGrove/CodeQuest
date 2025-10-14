// Package ui provides the terminal user interface for CodeQuest.
// This file contains the main Bubble Tea application model that orchestrates
// all UI screens and manages application state.
package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/AutumnsGrove/codequest/internal/game"
	"github.com/AutumnsGrove/codequest/internal/storage"
	"github.com/AutumnsGrove/codequest/internal/ui/screens"
)

// Screen represents different application screens/views.
// The app switches between these screens based on user navigation.
type Screen int

const (
	// ScreenDashboard is the main overview screen showing character stats and quick actions
	ScreenDashboard Screen = iota

	// ScreenQuestBoard displays available, active, and completed quests
	ScreenQuestBoard

	// ScreenCharacter shows detailed character sheet with stats, progress, and history
	ScreenCharacter

	// ScreenMentor is the AI assistant screen for getting help and guidance
	ScreenMentor

	// ScreenSettings allows configuration of app settings and preferences
	ScreenSettings
)

// Model is the main Bubble Tea model for the CodeQuest application.
// It manages all application state and coordinates between different screens.
//
// Architecture:
//   - Game State: character, quests, eventBus (persistent)
//   - UI State: currentScreen, keys, dimensions (transient)
//   - Storage: storage client for save/load operations
//   - Status: loading, error handling
type Model struct {
	// Game State - Core game data
	character *game.Character  // Player character (nil if not loaded)
	quests    []*game.Quest    // All quests (active, completed, available)
	eventBus  *game.EventBus   // Event system for game events

	// Storage - Data persistence
	storage *storage.SkateClient // Skate KV store client

	// UI State - Current view and controls
	currentScreen Screen   // Which screen is currently displayed
	keys          *KeyMap  // Key bindings for navigation

	// Quest Board state
	questBoardSelectedIndex int                  // Currently selected quest index
	questBoardFilter        screens.QuestFilter  // Current quest filter

	// Terminal dimensions - Updated on window resize
	width  int // Terminal width in characters
	height int // Terminal height in characters

	// Status - Application state
	loading bool   // True during async operations (load, save)
	err     error  // Most recent error (if any)

	// Timer state - Session timer (stub for now)
	timerRunning bool // Whether the session timer is active
}

// NewModel creates and initializes a new application model.
// This is the factory function for creating the Model instance.
//
// Parameters:
//   - storageClient: Skate storage client for save/load operations (required)
//
// Returns:
//   - *Model: A new Model instance ready for Bubble Tea initialization
//
// The model starts on the Dashboard screen with loading state active.
// Init() will attempt to load saved data from storage.
func NewModel(storageClient *storage.SkateClient) *Model {
	return &Model{
		// Game State - Will be loaded in Init()
		character: nil,
		quests:    []*game.Quest{},
		eventBus:  game.NewEventBus(),

		// Storage
		storage: storageClient,

		// UI State - Start on dashboard
		currentScreen: ScreenDashboard,
		keys:          NewKeyMap(),

		// Quest Board state
		questBoardSelectedIndex: 0,
		questBoardFilter:        screens.FilterAll,

		// Dimensions - Will be set on first WindowSizeMsg
		width:  80,  // Default width
		height: 24,  // Default height

		// Status
		loading: true, // Start in loading state
		err:     nil,

		// Timer
		timerRunning: false,
	}
}

// Init is the Bubble Tea initialization method.
// It returns commands to load character and quests from storage.
//
// If loading fails (e.g., first run), it will create a new character.
//
// Returns:
//   - tea.Cmd: Commands to load data from storage
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadCharacterCmd(m.storage),
		loadQuestsCmd(m.storage),
	)
}

// Update is the Bubble Tea update method.
// It handles all messages (key presses, window size changes, custom events)
// and returns an updated model and optional commands.
//
// Message Flow:
//  1. KeyMsg: User keyboard input (navigation, actions, quit)
//  2. WindowSizeMsg: Terminal resize events
//  3. characterLoadedMsg: Character loaded from storage
//  4. questsLoadedMsg: Quests loaded from storage
//  5. errorMsg: Error occurred during async operation
//
// Parameters:
//   - msg: The message to handle
//
// Returns:
//   - tea.Model: Updated model
//   - tea.Cmd: Optional command to execute
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Key press handling
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	// Window size changes
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	// Character loaded from storage
	case characterLoadedMsg:
		m.character = msg.character
		m.loading = false
		return m, nil

	// Quests loaded from storage
	case questsLoadedMsg:
		m.quests = msg.quests
		return m, nil

	// Error occurred
	case errorMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	}

	return m, nil
}

// View is the Bubble Tea view method.
// It renders the current screen based on the currentScreen value.
//
// Returns:
//   - string: The rendered UI to display
func (m Model) View() string {
	// Show loading screen if still loading
	if m.loading {
		return m.viewLoading()
	}

	// Show error screen if error occurred
	if m.err != nil {
		return m.viewError()
	}

	// Route to appropriate screen based on currentScreen
	switch m.currentScreen {
	case ScreenDashboard:
		return m.viewDashboard()
	case ScreenQuestBoard:
		return m.viewQuestBoard()
	case ScreenCharacter:
		return m.viewCharacter()
	case ScreenMentor:
		return m.viewMentor()
	case ScreenSettings:
		return m.viewSettings()
	default:
		return "Unknown screen"
	}
}

// handleKeyPress handles keyboard input and routes to appropriate handlers.
//
// Priority order:
//  1. Global keys (Ctrl+C quit, Alt+ modifiers)
//  2. Screen-specific keys (Q, C, M, S on dashboard)
//
// Parameters:
//   - msg: The key press message
//
// Returns:
//   - tea.Model: Updated model
//   - tea.Cmd: Optional command
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit (Ctrl+C)
	if key.Matches(msg, m.keys.GlobalQuit) {
		return m, tea.Quit
	}

	// Global save (Ctrl+S)
	if key.Matches(msg, m.keys.Save) {
		return m, m.saveStateCmd()
	}

	// Global timer toggle (Ctrl+T)
	if key.Matches(msg, m.keys.GlobalTimer) {
		m.timerRunning = !m.timerRunning
		return m, nil
	}

	// Global navigation (Alt+ modifiers - work from any screen)
	if key.Matches(msg, m.keys.GlobalDashboard) {
		return m.switchScreen(ScreenDashboard)
	}
	if key.Matches(msg, m.keys.GlobalMentor) {
		return m.switchScreen(ScreenMentor)
	}
	if key.Matches(msg, m.keys.GlobalSettings) {
		return m.switchScreen(ScreenSettings)
	}

	// Dashboard-specific keys (only active on dashboard)
	if m.currentScreen == ScreenDashboard {
		if key.Matches(msg, m.keys.DashboardQuests) {
			return m.switchScreen(ScreenQuestBoard)
		}
		if key.Matches(msg, m.keys.DashboardCharacter) {
			return m.switchScreen(ScreenCharacter)
		}
		if key.Matches(msg, m.keys.DashboardMentor) {
			return m.switchScreen(ScreenMentor)
		}
		if key.Matches(msg, m.keys.DashboardSettings) {
			return m.switchScreen(ScreenSettings)
		}
	}

	// Quest Board specific keys
	if m.currentScreen == ScreenQuestBoard {
		return m.handleQuestBoardKeys(msg)
	}

	// Escape key - return to dashboard from any screen
	if key.Matches(msg, m.keys.Esc) && m.currentScreen != ScreenDashboard {
		return m.switchScreen(ScreenDashboard)
	}

	return m, nil
}

// switchScreen changes the current screen and updates key bindings.
//
// When switching to dashboard, it enables dashboard-specific single-key shortcuts.
// When leaving dashboard, it disables them to prevent conflicts with text input.
//
// Parameters:
//   - screen: The screen to switch to
//
// Returns:
//   - tea.Model: Updated model
//   - tea.Cmd: Optional command (nil for now)
func (m Model) switchScreen(screen Screen) (tea.Model, tea.Cmd) {
	m.currentScreen = screen

	// Reset quest board state when switching to it
	if screen == ScreenQuestBoard {
		m.questBoardSelectedIndex = 0
		m.questBoardFilter = screens.FilterAll
	}

	// Enable/disable dashboard keys based on screen
	if screen == ScreenDashboard {
		m.keys.EnableDashboardKeys()
	} else {
		m.keys.DisableDashboardKeys()
	}

	return m, nil
}

// handleQuestBoardKeys handles keyboard input specific to the Quest Board screen.
//
// Supports:
//   - Up/Down: Navigate quest list
//   - Enter: Start/view selected quest (placeholder for now)
//   - F: Cycle through filters
//
// Parameters:
//   - msg: The key press message
//
// Returns:
//   - tea.Model: Updated model
//   - tea.Cmd: Optional command
func (m Model) handleQuestBoardKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Filter quests based on current filter to get the correct count
	filteredQuests := m.getFilteredQuests()
	maxIndex := len(filteredQuests) - 1

	// Up arrow - select previous quest
	if key.Matches(msg, m.keys.Up) {
		if m.questBoardSelectedIndex > 0 {
			m.questBoardSelectedIndex--
		}
		return m, nil
	}

	// Down arrow - select next quest
	if key.Matches(msg, m.keys.Down) {
		if maxIndex >= 0 && m.questBoardSelectedIndex < maxIndex {
			m.questBoardSelectedIndex++
		}
		return m, nil
	}

	// F key - cycle through filters
	if msg.String() == "f" || msg.String() == "F" {
		m.questBoardFilter = (m.questBoardFilter + 1) % 4
		m.questBoardSelectedIndex = 0 // Reset selection when filter changes
		return m, nil
	}

	// Enter key - start/view quest (placeholder)
	if key.Matches(msg, m.keys.Enter) {
		// TODO: Implement quest start/view logic
		// For now, just do nothing
		return m, nil
	}

	return m, nil
}

// getFilteredQuests returns quests filtered by the current filter.
func (m Model) getFilteredQuests() []*game.Quest {
	if m.questBoardFilter == screens.FilterAll {
		return m.quests
	}

	filtered := make([]*game.Quest, 0)
	for _, quest := range m.quests {
		switch m.questBoardFilter {
		case screens.FilterAvailable:
			if quest.Status == game.QuestAvailable {
				filtered = append(filtered, quest)
			}
		case screens.FilterActive:
			if quest.Status == game.QuestActive {
				filtered = append(filtered, quest)
			}
		case screens.FilterCompleted:
			if quest.Status == game.QuestCompleted {
				filtered = append(filtered, quest)
			}
		}
	}
	return filtered
}

// saveStateCmd returns a command to save current game state to storage.
func (m Model) saveStateCmd() tea.Cmd {
	return func() tea.Msg {
		// Save character
		if m.character != nil {
			if err := m.storage.SaveCharacter(m.character); err != nil {
				return errorMsg{err: fmt.Errorf("failed to save character: %w", err)}
			}
		}

		// Save quests
		if err := m.storage.SaveQuests(m.quests); err != nil {
			return errorMsg{err: fmt.Errorf("failed to save quests: %w", err)}
		}

		return saveCompletedMsg{}
	}
}

// ============================================================================
// Screen Rendering Methods (Placeholder implementations)
// ============================================================================
// These will be enhanced by later subagents (15-21) with detailed screens.
// For now, they return simple placeholder text to verify routing works.

// viewLoading renders the loading screen.
func (m Model) viewLoading() string {
	frame := 0 // TODO: Animate this with a tick command
	spinner := RenderLoadingSpinner(frame, "Loading CodeQuest...")
	return PlaceInCenter(m.width, m.height, spinner)
}

// viewError renders the error screen.
func (m Model) viewError() string {
	errorText := ErrorTextStyle.Render("Error: " + m.err.Error())
	hint := MutedTextStyle.Render("\nPress Ctrl+C to quit")
	return PlaceInCenter(m.width, m.height, errorText+hint)
}

// viewDashboard renders the dashboard screen.
// Delegates to screens.RenderDashboard for full implementation.
func (m Model) viewDashboard() string {
	return screens.RenderDashboard(m.character, m.quests, m.width, m.height)
}

// viewQuestBoard renders the quest board screen.
// Delegates to screens.RenderQuestBoard for full implementation.
func (m Model) viewQuestBoard() string {
	return screens.RenderQuestBoard(
		m.character,
		m.quests,
		m.questBoardSelectedIndex,
		m.questBoardFilter,
		m.width,
		m.height,
	)
}

// viewCharacter renders the character sheet screen.
// TODO: Subagent 17 will implement full character sheet with stats, history, etc.
func (m Model) viewCharacter() string {
	title := RenderTitle("Character Sheet", "⚔️")

	var charInfo string
	if m.character != nil {
		charInfo = fmt.Sprintf(
			"Name: %s\n"+
				"Level: %d\n"+
				"XP: %d/%d\n"+
				"Code Power: %d\n"+
				"Wisdom: %d\n"+
				"Agility: %d\n\n",
			m.character.Name,
			m.character.Level,
			m.character.XP,
			m.character.XPToNextLevel,
			m.character.CodePower,
			m.character.Wisdom,
			m.character.Agility,
		)
	} else {
		charInfo = "No character loaded.\n\n"
	}

	help := m.keys.RenderCharacterHelp()

	content := title + "\n\n" + charInfo + help

	return BoxStyle.Render(content)
}

// viewMentor renders the mentor/AI assistant screen.
// TODO: Subagent 18 will implement full mentor screen with chat interface.
func (m Model) viewMentor() string {
	title := RenderTitle("AI Mentor", "🧙")

	mentorInfo := "AI Mentor is coming soon!\n\n" +
		"This screen will provide:\n" +
		"- Quest guidance and hints\n" +
		"- Coding tips and best practices\n" +
		"- Progress insights\n\n"

	help := m.keys.RenderMentorHelp()

	content := title + "\n\n" + mentorInfo + help

	return BoxStyle.Render(content)
}

// viewSettings renders the settings screen.
// TODO: Subagent 19 will implement full settings screen with configuration options.
func (m Model) viewSettings() string {
	title := RenderTitle("Settings", "⚙️")

	settingsInfo := "Settings screen coming soon!\n\n" +
		"Configure:\n" +
		"- AI provider preferences\n" +
		"- Notification settings\n" +
		"- Theme and display options\n\n"

	help := m.keys.RenderSettingsHelp()

	content := title + "\n\n" + settingsInfo + help

	return BoxStyle.Render(content)
}

// ============================================================================
// Custom Messages for Async Operations
// ============================================================================

// characterLoadedMsg is sent when character loading completes.
type characterLoadedMsg struct {
	character *game.Character
}

// questsLoadedMsg is sent when quest loading completes.
type questsLoadedMsg struct {
	quests []*game.Quest
}

// errorMsg is sent when an error occurs.
type errorMsg struct {
	err error
}

// saveCompletedMsg is sent when save operation completes successfully.
type saveCompletedMsg struct{}

// ============================================================================
// Async Commands for Storage Operations
// ============================================================================

// loadCharacterCmd loads the character from storage asynchronously.
func loadCharacterCmd(storage *storage.SkateClient) tea.Cmd {
	return func() tea.Msg {
		// Try to load existing character
		character, err := storage.LoadCharacter()
		if err != nil {
			// If character doesn't exist (first run), create a new one
			character = game.NewCharacter("Adventurer")

			// Save the new character
			if saveErr := storage.SaveCharacter(character); saveErr != nil {
				return errorMsg{err: fmt.Errorf("failed to create new character: %w", saveErr)}
			}
		}

		return characterLoadedMsg{character: character}
	}
}

// loadQuestsCmd loads quests from storage asynchronously.
func loadQuestsCmd(storage *storage.SkateClient) tea.Cmd {
	return func() tea.Msg {
		quests, err := storage.LoadQuests()
		if err != nil {
			// Return empty quests on error (non-fatal)
			return questsLoadedMsg{quests: []*game.Quest{}}
		}

		return questsLoadedMsg{quests: quests}
	}
}
