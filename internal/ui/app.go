// Package ui provides the terminal user interface for CodeQuest.
// This file contains the main Bubble Tea application model that orchestrates
// all UI screens and manages application state.
package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
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

	// Mentor screen state
	mentorInputText   string             // Current text in mentor input field
	mentorConversation []screens.Message // Conversation history with AI mentor

	// Terminal dimensions - Updated on window resize
	width  int // Terminal width in characters
	height int // Terminal height in characters

	// Status - Application state
	loading bool   // True during async operations (load, save)
	err     error  // Most recent error (if any)

	// Timer state - Session timer (stub for now)
	timerRunning bool // Whether the session timer is active

	// Help overlay state
	showingHelp bool // Whether the help overlay is currently displayed

	// Notification system - Real-time event notifications
	notifications      []Notification // Queue of pending notifications
	currentNotification *Notification  // Currently displayed notification (nil if none)
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

		// Mentor screen state
		mentorInputText:    "",
		mentorConversation: []screens.Message{},

		// Dimensions - Will be set on first WindowSizeMsg
		width:  80,  // Default width
		height: 24,  // Default height

		// Status
		loading: true, // Start in loading state
		err:     nil,

		// Timer
		timerRunning: false,

		// Help overlay
		showingHelp: false,

		// Notifications
		notifications:      []Notification{},
		currentNotification: nil,
	}
}

// Init is the Bubble Tea initialization method.
// It returns commands to load character and quests from storage,
// and subscribes to game events for real-time UI updates.
//
// If loading fails (e.g., first run), it will create a new character.
//
// Returns:
//   - tea.Cmd: Commands to load data from storage and listen for events
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadCharacterCmd(m.storage),
		loadQuestsCmd(m.storage),
		listenForGameEvents(m.eventBus), // Subscribe to game events
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

	// ========================================================================
	// Game Event Messages - Real-time updates from EventBus
	// ========================================================================

	// Commit detected - Show XP gain notification
	case commitDetectedMsg:
		// Add XP gain notification
		notification := Notification{
			Message:   fmt.Sprintf("+%d XP from commit!", msg.xpAwarded),
			Type:      NotificationSuccess,
			Duration:  3 * time.Second,
			Timestamp: time.Now(),
		}
		m.addNotification(notification)

		// Reload character data to reflect new XP and continue listening
		return m, tea.Batch(
			loadCharacterCmd(m.storage),
			m.showNextNotification(),
			listenForGameEvents(m.eventBus), // Keep listening for more events
		)

	// Level up - Show celebration and reload character
	case levelUpMsg:
		// Add level-up notification with celebration
		notification := Notification{
			Message:   fmt.Sprintf("⚡ LEVEL UP! ⚡\nYou are now Level %d!", msg.newLevel),
			Type:      NotificationLevelUp,
			Duration:  5 * time.Second,
			Timestamp: time.Now(),
		}
		m.addNotification(notification)

		// Reload character data to reflect new level and continue listening
		return m, tea.Batch(
			loadCharacterCmd(m.storage),
			m.showNextNotification(),
			listenForGameEvents(m.eventBus), // Keep listening for more events
		)

	// Quest completed - Show completion notification and reload quests
	case questCompleteMsg:
		// Add quest completion notification
		notification := Notification{
			Message:   fmt.Sprintf("✓ QUEST COMPLETE!\n%s\n+%d XP", msg.questName, msg.xpAwarded),
			Type:      NotificationQuestComplete,
			Duration:  4 * time.Second,
			Timestamp: time.Now(),
		}
		m.addNotification(notification)

		// Reload both character (XP changed) and quests (status changed), and continue listening
		return m, tea.Batch(
			loadCharacterCmd(m.storage),
			loadQuestsCmd(m.storage),
			m.showNextNotification(),
			listenForGameEvents(m.eventBus), // Keep listening for more events
		)

	// Quest started - Reload quests to reflect new active quest
	case questStartMsg:
		// Add quest start notification
		notification := Notification{
			Message:   fmt.Sprintf("Quest Started: %s", msg.questName),
			Type:      NotificationInfo,
			Duration:  3 * time.Second,
			Timestamp: time.Now(),
		}
		m.addNotification(notification)

		// Reload quests and continue listening
		return m, tea.Batch(
			loadQuestsCmd(m.storage),
			m.showNextNotification(),
			listenForGameEvents(m.eventBus), // Keep listening for more events
		)

	// Notification dismissed - Show next notification if any
	case notificationDismissedMsg:
		m.currentNotification = nil
		return m, m.showNextNotification()
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

	// Get the main screen content
	var mainContent string
	switch m.currentScreen {
	case ScreenDashboard:
		mainContent = m.viewDashboard()
	case ScreenQuestBoard:
		mainContent = m.viewQuestBoard()
	case ScreenCharacter:
		mainContent = m.viewCharacter()
	case ScreenMentor:
		mainContent = m.viewMentor()
	case ScreenSettings:
		mainContent = m.viewSettings()
	default:
		mainContent = "Unknown screen"
	}

	// If help overlay is showing, render it on top
	if m.showingHelp {
		return m.viewHelpOverlay(mainContent)
	}

	// If there's a current notification, render it on top of main content
	if m.currentNotification != nil {
		return m.viewWithNotification(mainContent)
	}

	return mainContent
}

// handleKeyPress handles keyboard input and routes to appropriate handlers.
//
// Priority order:
//  1. Help overlay (if showing, Esc to close)
//  2. Global keys (Ctrl+C quit, ? for help, Alt+ modifiers)
//  3. Screen-specific keys (Q, C, M, S on dashboard)
//
// Parameters:
//   - msg: The key press message
//
// Returns:
//   - tea.Model: Updated model
//   - tea.Cmd: Optional command
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If help overlay is showing, only handle Esc to close it
	if m.showingHelp {
		if key.Matches(msg, m.keys.Esc) {
			m.showingHelp = false
			return m, nil
		}
		// Ignore all other keys when help is showing
		return m, nil
	}

	// Global quit (Ctrl+C)
	if key.Matches(msg, m.keys.GlobalQuit) {
		return m, tea.Quit
	}

	// Global help overlay (? key) - works from any screen except dashboard (dashboard uses ? for help)
	if m.currentScreen != ScreenDashboard && key.Matches(msg, m.keys.HelpOverlay) {
		m.showingHelp = true
		return m, nil
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
// Delegates to screens.RenderCharacter for full implementation.
func (m Model) viewCharacter() string {
	return screens.RenderCharacter(m.character, m.width, m.height)
}

// viewMentor renders the mentor/AI assistant screen.
// Delegates to screens.RenderMentor for full implementation.
func (m Model) viewMentor() string {
	return screens.RenderMentor(
		m.character,
		m.mentorInputText,
		m.mentorConversation,
		m.width,
		m.height,
	)
}

// viewSettings renders the settings screen.
// Delegates to screens.RenderSettings for full implementation.
func (m Model) viewSettings() string {
	return screens.RenderSettings(m.character, m.width, m.height)
}

// viewHelpOverlay renders the help overlay on top of the main content.
// The overlay shows all available keyboard shortcuts for the current screen.
func (m Model) viewHelpOverlay(mainContent string) string {
	// Create help content based on current screen
	var helpTitle string
	var helpBindings []key.Binding

	switch m.currentScreen {
	case ScreenDashboard:
		helpTitle = "Dashboard Help"
		helpBindings = m.keys.DashboardHelp()
	case ScreenQuestBoard:
		helpTitle = "Quest Board Help"
		helpBindings = m.keys.QuestBoardHelp()
	case ScreenCharacter:
		helpTitle = "Character Sheet Help"
		helpBindings = m.keys.CharacterHelp()
	case ScreenMentor:
		helpTitle = "Mentor Help"
		helpBindings = m.keys.MentorHelp()
	case ScreenSettings:
		helpTitle = "Settings Help"
		helpBindings = m.keys.SettingsHelp()
	default:
		helpTitle = "Help"
		helpBindings = m.keys.ShortHelp()
	}

	// Build help text from bindings
	helpLines := make([]string, 0)
	helpLines = append(helpLines, TitleStyle.Render(helpTitle))
	helpLines = append(helpLines, "")

	for _, binding := range helpBindings {
		keys := binding.Help().Key
		desc := binding.Help().Desc
		line := RenderKeybind(keys, desc)
		helpLines = append(helpLines, line)
	}

	helpLines = append(helpLines, "")
	helpLines = append(helpLines, MutedTextStyle.Render("Press Esc to close this help overlay"))

	helpContent := lipgloss.JoinVertical(lipgloss.Left, helpLines...)

	// Wrap in modal style
	helpBox := ModalStyle.Render(helpContent)

	// Center the help box over the main content
	overlay := PlaceInCenter(m.width, m.height, helpBox)

	return overlay
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
// Game Event Messages - Bridge between EventBus and Bubble Tea
// ============================================================================
// These messages convert game.Event types into Bubble Tea messages
// for real-time UI updates.

// commitDetectedMsg is sent when a git commit is detected.
// The UI can show XP gain notifications and update character stats.
type commitDetectedMsg struct {
	sha          string
	message      string
	xpAwarded    int
	linesAdded   int
	linesRemoved int
}

// levelUpMsg is sent when the character gains a level.
// The UI can show a celebratory level-up animation.
type levelUpMsg struct {
	characterID string
	oldLevel    int
	newLevel    int
}

// questCompleteMsg is sent when a quest is completed.
// The UI can show quest completion notification and update quest list.
type questCompleteMsg struct {
	questID   string
	questName string
	xpAwarded int
}

// questStartMsg is sent when a quest is started.
// The UI can update the quest board to reflect the new active quest.
type questStartMsg struct {
	questID   string
	questName string
	questType string
}

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

// ============================================================================
// Notification System
// ============================================================================

// NotificationType defines different types of notifications with distinct styling.
type NotificationType int

const (
	// NotificationInfo - General information (blue)
	NotificationInfo NotificationType = iota
	// NotificationSuccess - Success messages (green)
	NotificationSuccess
	// NotificationWarning - Warnings (orange)
	NotificationWarning
	// NotificationError - Errors (red)
	NotificationError
	// NotificationLevelUp - Level-up celebration (gold)
	NotificationLevelUp
	// NotificationQuestComplete - Quest completion (cyan)
	NotificationQuestComplete
)

// Notification represents a temporary message to display to the user.
// Notifications auto-dismiss after their duration expires.
type Notification struct {
	Message   string           // The message to display
	Type      NotificationType // Visual style of notification
	Duration  time.Duration    // How long to display (0 = manual dismiss)
	Timestamp time.Time        // When the notification was created
}

// notificationDismissedMsg is sent when a notification's timer expires.
type notificationDismissedMsg struct{}

// addNotification adds a notification to the queue.
// If no notification is currently showing, it will be displayed immediately.
func (m *Model) addNotification(notification Notification) {
	m.notifications = append(m.notifications, notification)
}

// showNextNotification displays the next notification in the queue.
// If there's already a notification showing, this does nothing.
// If the queue is empty, this does nothing.
//
// Returns a tea.Cmd that will dismiss the notification after its duration.
func (m *Model) showNextNotification() tea.Cmd {
	// Don't show if there's already a notification
	if m.currentNotification != nil {
		return nil
	}

	// Don't show if queue is empty
	if len(m.notifications) == 0 {
		return nil
	}

	// Pop the first notification from the queue
	m.currentNotification = &m.notifications[0]
	m.notifications = m.notifications[1:]

	// Return a command to dismiss after duration
	if m.currentNotification.Duration > 0 {
		return tea.Tick(m.currentNotification.Duration, func(t time.Time) tea.Msg {
			return notificationDismissedMsg{}
		})
	}

	return nil
}

// viewWithNotification renders the main content with a notification overlay.
// The notification appears at the top-center of the screen.
func (m Model) viewWithNotification(mainContent string) string {
	if m.currentNotification == nil {
		return mainContent
	}

	// Render the notification box
	notificationBox := m.renderNotification(*m.currentNotification)

	// Center the notification horizontally
	centeredNotification := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render(notificationBox)

	// Place notification at the top of the screen
	// Since Lip Gloss doesn't support true overlays, we'll just prepend it
	result := lipgloss.JoinVertical(
		lipgloss.Left,
		"", // Empty line for spacing
		centeredNotification,
		"",
		mainContent,
	)

	return result
}

// renderNotification renders a single notification with appropriate styling.
func (m Model) renderNotification(n Notification) string {
	// Choose style based on notification type
	var style lipgloss.Style
	var icon string

	switch n.Type {
	case NotificationLevelUp:
		style = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorXP).
			Foreground(ColorXP).
			Background(lipgloss.Color("235")).
			Bold(true).
			Padding(1, 2)
		icon = "⚡"

	case NotificationQuestComplete:
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSuccess).
			Foreground(ColorSuccess).
			Background(lipgloss.Color("235")).
			Bold(true).
			Padding(1, 2)
		icon = "✓"

	case NotificationSuccess:
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSuccess).
			Foreground(ColorBright).
			Background(lipgloss.Color("235")).
			Padding(0, 2)
		icon = "✓"

	case NotificationInfo:
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorInfo).
			Foreground(ColorBright).
			Background(lipgloss.Color("235")).
			Padding(0, 2)
		icon = "ℹ"

	case NotificationWarning:
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorWarning).
			Foreground(ColorWarning).
			Background(lipgloss.Color("235")).
			Padding(0, 2)
		icon = "⚠"

	case NotificationError:
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorError).
			Foreground(ColorError).
			Background(lipgloss.Color("235")).
			Bold(true).
			Padding(0, 2)
		icon = "✗"

	default:
		style = NotificationStyle
		icon = "•"
	}

	// Format the message with icon
	content := icon + " " + n.Message

	return style.Render(content)
}

// ============================================================================
// Event Bus Integration - Bridge to Bubble Tea
// ============================================================================

// listenForGameEvents subscribes to the EventBus and converts game events
// to Bubble Tea messages. This creates a bridge between the game's EventBus
// and Bubble Tea's message system.
//
// Architecture:
//  1. Subscribe to EventBus events with handlers
//  2. Handlers convert game.Event to Bubble Tea messages via channels
//  3. A goroutine listens to the channel and returns messages to Bubble Tea
//
// Thread Safety:
// - EventBus handlers run in the publisher's goroutine (synchronously)
// - We use a buffered channel to prevent blocking the EventBus
// - The channel-to-message conversion runs in a dedicated goroutine
//
// Parameters:
//   - eventBus: The game event bus to subscribe to
//
// Returns:
//   - tea.Cmd: A command that waits for the next game event
func listenForGameEvents(eventBus *game.EventBus) tea.Cmd {
	return func() tea.Msg {
		// Create a buffered channel for game events
		eventChan := make(chan game.Event, 10)

		// Subscribe to all relevant event types
		// Each subscription adds a handler that sends events to our channel
		eventBus.Subscribe(game.EventCommit, func(e game.Event) {
			select {
			case eventChan <- e:
				// Event sent successfully
			default:
				// Channel full, drop event (prevents blocking game logic)
			}
		})

		eventBus.Subscribe(game.EventLevelUp, func(e game.Event) {
			select {
			case eventChan <- e:
			default:
			}
		})

		eventBus.Subscribe(game.EventQuestDone, func(e game.Event) {
			select {
			case eventChan <- e:
			default:
			}
		})

		eventBus.Subscribe(game.EventQuestStart, func(e game.Event) {
			select {
			case eventChan <- e:
			default:
			}
		})

		// Wait for the first event from the channel
		// This blocks until an event is received
		event := <-eventChan

		// Convert the game event to a Bubble Tea message
		return convertEventToMessage(event)
	}
}

// waitForNextEvent returns a command that waits for the next game event.
// This should be called after processing each event to keep listening.
func waitForNextEvent(eventChan <-chan game.Event) tea.Cmd {
	return func() tea.Msg {
		event := <-eventChan
		return convertEventToMessage(event)
	}
}

// convertEventToMessage converts a game.Event to the appropriate Bubble Tea message.
// This is the translation layer between the game's event system and the UI.
//
// Parameters:
//   - event: The game event to convert
//
// Returns:
//   - tea.Msg: The corresponding Bubble Tea message
func convertEventToMessage(event game.Event) tea.Msg {
	switch event.Type {
	case game.EventCommit:
		// Extract commit event data
		sha, _ := event.Data["sha"].(string)
		message, _ := event.Data["message"].(string)
		linesAdded, _ := event.Data["lines_added"].(int)
		linesRemoved, _ := event.Data["lines_removed"].(int)

		// Calculate XP awarded (simplified - actual XP comes from handler)
		// For display purposes, we'll estimate it
		// In reality, the handler already calculated and applied it
		// We're just showing a notification, so approximate XP is fine
		baseXP := (linesAdded + linesRemoved) * 5
		if baseXP < 10 {
			baseXP = 10 // Minimum XP
		}

		return commitDetectedMsg{
			sha:          sha,
			message:      message,
			xpAwarded:    baseXP,
			linesAdded:   linesAdded,
			linesRemoved: linesRemoved,
		}

	case game.EventLevelUp:
		// Extract level-up event data
		characterID, _ := event.Data["character_id"].(string)
		oldLevel, _ := event.Data["old_level"].(int)
		newLevel, _ := event.Data["new_level"].(int)

		return levelUpMsg{
			characterID: characterID,
			oldLevel:    oldLevel,
			newLevel:    newLevel,
		}

	case game.EventQuestDone:
		// Extract quest completion event data
		questID, _ := event.Data["quest_id"].(string)
		questTitle, _ := event.Data["quest_title"].(string)
		xpReward, _ := event.Data["xp_reward"].(int)

		return questCompleteMsg{
			questID:   questID,
			questName: questTitle,
			xpAwarded: xpReward,
		}

	case game.EventQuestStart:
		// Extract quest start event data
		questID, _ := event.Data["quest_id"].(string)
		questTitle, _ := event.Data["quest_title"].(string)
		questType, _ := event.Data["quest_type"].(string)

		return questStartMsg{
			questID:   questID,
			questName: questTitle,
			questType: questType,
		}

	default:
		// Unknown event type - return nil message
		return nil
	}
}

// Note on Event Bridge Design:
//
// This implementation uses a simplified approach where listenForGameEvents
// returns after receiving ONE event. To continue listening, the Update method
// needs to call listenForGameEvents again after processing each event.
//
// This is the correct Bubble Tea pattern:
// 1. Init() calls listenForGameEvents()
// 2. First event arrives, Update() is called
// 3. Update() processes the event and returns listenForGameEvents() as a cmd
// 4. Next event arrives, Update() is called again
// 5. Repeat...
//
// This ensures events are processed one at a time in the main Bubble Tea loop,
// maintaining thread safety without complex synchronization.
