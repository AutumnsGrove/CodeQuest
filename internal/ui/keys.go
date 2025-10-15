package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines all key bindings for the application.
// This struct contains every keyboard shortcut used throughout CodeQuest,
// organized by their function. Keys can be reconfigured at runtime.
type KeyMap struct {
	// Navigation keys - Basic movement
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Tab   key.Binding

	// Action keys - Common interactions
	Enter key.Binding
	Esc   key.Binding
	Space key.Binding

	// Dashboard screen shortcuts (single keys - safe in non-input mode)
	DashboardQuests    key.Binding
	DashboardCharacter key.Binding
	DashboardInventory key.Binding
	DashboardMentor    key.Binding
	DashboardSettings  key.Binding
	DashboardHelpKey   key.Binding

	// Help overlay key (works from any screen)
	HelpOverlay key.Binding

	// Global shortcuts (modifiers required - safe in input mode)
	GlobalDashboard key.Binding
	GlobalMentor    key.Binding
	GlobalSettings  key.Binding
	GlobalHelp      key.Binding
	GlobalTimer     key.Binding
	GlobalQuit      key.Binding

	// Special function keys
	CommandPalette key.Binding
	Save           key.Binding
	Cancel         key.Binding
}

// NewKeyMap creates a new KeyMap with default bindings.
// This follows Bubble Tea conventions and provides both standard arrow keys
// and vim-style hjkl navigation for power users.
func NewKeyMap() *KeyMap {
	return &KeyMap{
		// Navigation - Arrow keys + vim-style alternatives
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next section"),
		),

		// Action keys
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/accept"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back/cancel"),
		),
		Space: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle/action"),
		),

		// Dashboard screen shortcuts (single keys)
		// These are only active when on the dashboard screen (non-input mode)
		DashboardQuests: key.NewBinding(
			key.WithKeys("q", "Q"),
			key.WithHelp("Q", "quest board"),
		),
		DashboardCharacter: key.NewBinding(
			key.WithKeys("c", "C"),
			key.WithHelp("C", "character sheet"),
		),
		DashboardInventory: key.NewBinding(
			key.WithKeys("i", "I"),
			key.WithHelp("I", "inventory/skills"),
		),
		DashboardMentor: key.NewBinding(
			key.WithKeys("m", "M"),
			key.WithHelp("M", "mentor"),
		),
		DashboardSettings: key.NewBinding(
			key.WithKeys("s", "S"),
			key.WithHelp("S", "settings"),
		),
		DashboardHelpKey: key.NewBinding(
			key.WithKeys("h", "H", "?"),
			key.WithHelp("H/?", "help"),
		),

		// Help overlay key (works from any screen)
		HelpOverlay: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "show help"),
		),

		// Global shortcuts (modifiers required - safe everywhere)
		// These work from any screen, including text input screens
		GlobalDashboard: key.NewBinding(
			key.WithKeys("alt+q"),
			key.WithHelp("alt+Q", "return to dashboard"),
		),
		GlobalMentor: key.NewBinding(
			key.WithKeys("alt+m"),
			key.WithHelp("alt+M", "quick mentor help"),
		),
		GlobalSettings: key.NewBinding(
			key.WithKeys("alt+s"),
			key.WithHelp("alt+S", "settings"),
		),
		GlobalHelp: key.NewBinding(
			key.WithKeys("alt+h", "alt+?"),
			key.WithHelp("alt+H", "help overlay"),
		),
		GlobalTimer: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+T", "toggle session timer"),
		),
		GlobalQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+C", "quit application"),
		),

		// Special function keys
		CommandPalette: key.NewBinding(
			key.WithKeys("alt+/", "ctrl+p"),
			key.WithHelp("alt+/", "command palette"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+S", "save"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc", "ctrl+g"),
			key.WithHelp("esc", "cancel"),
		),
	}
}

// ShortHelp returns a slice of key bindings to show in the short help view.
// This is typically displayed at the bottom of screens in a compact format.
// It shows the most essential keybinds that users need frequently.
func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.GlobalQuit,
		k.GlobalHelp,
	}
}

// FullHelp returns a multi-column layout of key bindings for the full help view.
// This is displayed when the user requests detailed help (Alt+H or ?).
// Keys are organized into logical groups for easier scanning.
func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// Column 1: Navigation
		{k.Up, k.Down, k.Left, k.Right, k.Tab},
		// Column 2: Actions
		{k.Enter, k.Space, k.Esc},
		// Column 3: Screen navigation (dashboard)
		{k.DashboardQuests, k.DashboardCharacter, k.DashboardInventory},
		// Column 4: More screens
		{k.DashboardMentor, k.DashboardSettings, k.DashboardHelpKey},
		// Column 5: Global shortcuts
		{k.GlobalDashboard, k.GlobalMentor, k.GlobalSettings, k.GlobalTimer},
		// Column 6: Special functions
		{k.CommandPalette, k.Save, k.GlobalQuit, k.GlobalHelp},
	}
}

// DashboardHelp returns key bindings specific to the dashboard screen.
// Shows only the shortcuts that work when on the dashboard.
func (k *KeyMap) DashboardHelp() []key.Binding {
	return []key.Binding{
		k.DashboardQuests,
		k.DashboardCharacter,
		k.DashboardInventory,
		k.DashboardMentor,
		k.DashboardSettings,
		k.DashboardHelpKey,
		k.GlobalTimer,
		k.GlobalQuit,
	}
}

// QuestBoardHelp returns key bindings specific to the quest board screen.
// Uses global shortcuts (Alt+ modifiers) to avoid conflicts with text input.
func (k *KeyMap) QuestBoardHelp() []key.Binding {
	return []key.Binding{
		k.Up,
		k.Down,
		k.Enter,
		k.GlobalDashboard,
		k.GlobalMentor,
		k.Esc,
	}
}

// CharacterHelp returns key bindings specific to the character sheet screen.
func (k *KeyMap) CharacterHelp() []key.Binding {
	return []key.Binding{
		k.Up,
		k.Down,
		k.Tab,
		k.GlobalDashboard,
		k.Esc,
	}
}

// MentorHelp returns key bindings specific to the mentor/AI screen.
// Most keys are input-safe (Alt+ modifiers) since user may be typing questions.
func (k *KeyMap) MentorHelp() []key.Binding {
	return []key.Binding{
		k.Enter,
		k.GlobalDashboard,
		k.GlobalSettings,
		k.Esc,
	}
}

// SettingsHelp returns key bindings specific to the settings screen.
func (k *KeyMap) SettingsHelp() []key.Binding {
	return []key.Binding{
		k.Up,
		k.Down,
		k.Tab,
		k.Space,
		k.Enter,
		k.Save,
		k.Esc,
	}
}

// RenderDashboardHelp formats the dashboard help text for display.
// Returns a formatted string showing all available dashboard shortcuts.
func (k *KeyMap) RenderDashboardHelp() string {
	return RenderKeybind("Q", "Quests") + "  " +
		RenderKeybind("C", "Character") + "  " +
		RenderKeybind("I", "Inventory") + "  " +
		RenderKeybind("M", "Mentor") + "\n" +
		RenderKeybind("S", "Settings") + "  " +
		RenderKeybind("H", "Help") + "  " +
		RenderKeybind("Ctrl+T", "Timer") + "  " +
		RenderKeybind("Esc", "Exit")
}

// RenderQuestBoardHelp formats the quest board help text for display.
func (k *KeyMap) RenderQuestBoardHelp() string {
	return RenderKeybind("Alt+Q", "Dashboard") + "  " +
		RenderKeybind("Alt+M", "Mentor") + "  " +
		RenderKeybind("Alt+S", "Settings") + "\n" +
		RenderKeybind("↑↓", "Navigate") + "  " +
		RenderKeybind("Enter", "Accept") + "  " +
		RenderKeybind("Esc", "Back")
}

// RenderCharacterHelp formats the character sheet help text for display.
func (k *KeyMap) RenderCharacterHelp() string {
	return RenderKeybind("Alt+Q", "Dashboard") + "  " +
		RenderKeybind("Alt+M", "Mentor") + "\n" +
		RenderKeybind("↑↓", "Navigate") + "  " +
		RenderKeybind("Tab", "Next Section") + "  " +
		RenderKeybind("Esc", "Back")
}

// RenderMentorHelp formats the mentor screen help text for display.
func (k *KeyMap) RenderMentorHelp() string {
	return RenderKeybind("Alt+Q", "Dashboard") + "  " +
		RenderKeybind("Alt+S", "Settings") + "\n" +
		RenderKeybind("Enter", "Send") + "  " +
		RenderKeybind("Esc", "Back")
}

// RenderSettingsHelp formats the settings screen help text for display.
func (k *KeyMap) RenderSettingsHelp() string {
	return RenderKeybind("↑↓", "Navigate") + "  " +
		RenderKeybind("Space", "Toggle") + "  " +
		RenderKeybind("Enter", "Edit") + "\n" +
		RenderKeybind("Ctrl+S", "Save") + "  " +
		RenderKeybind("Esc", "Cancel")
}

// EnableDashboardKeys enables dashboard-specific single-key shortcuts.
// Call this when transitioning to the dashboard screen.
func (k *KeyMap) EnableDashboardKeys() {
	k.DashboardQuests.SetEnabled(true)
	k.DashboardCharacter.SetEnabled(true)
	k.DashboardInventory.SetEnabled(true)
	k.DashboardMentor.SetEnabled(true)
	k.DashboardSettings.SetEnabled(true)
	k.DashboardHelpKey.SetEnabled(true)
}

// DisableDashboardKeys disables dashboard-specific single-key shortcuts.
// Call this when transitioning away from the dashboard to prevent
// accidental triggers during text input on other screens.
func (k *KeyMap) DisableDashboardKeys() {
	k.DashboardQuests.SetEnabled(false)
	k.DashboardCharacter.SetEnabled(false)
	k.DashboardInventory.SetEnabled(false)
	k.DashboardMentor.SetEnabled(false)
	k.DashboardSettings.SetEnabled(false)
	k.DashboardHelpKey.SetEnabled(false)
}

// EnableAllKeys enables all key bindings.
// Useful for resetting key state or during initialization.
func (k *KeyMap) EnableAllKeys() {
	k.Up.SetEnabled(true)
	k.Down.SetEnabled(true)
	k.Left.SetEnabled(true)
	k.Right.SetEnabled(true)
	k.Tab.SetEnabled(true)
	k.Enter.SetEnabled(true)
	k.Esc.SetEnabled(true)
	k.Space.SetEnabled(true)

	k.EnableDashboardKeys()

	k.GlobalDashboard.SetEnabled(true)
	k.GlobalMentor.SetEnabled(true)
	k.GlobalSettings.SetEnabled(true)
	k.GlobalHelp.SetEnabled(true)
	k.GlobalTimer.SetEnabled(true)
	k.GlobalQuit.SetEnabled(true)

	k.CommandPalette.SetEnabled(true)
	k.Save.SetEnabled(true)
	k.Cancel.SetEnabled(true)
}

// DisableAllKeys disables all key bindings.
// Useful for modal dialogs or special input modes.
func (k *KeyMap) DisableAllKeys() {
	k.Up.SetEnabled(false)
	k.Down.SetEnabled(false)
	k.Left.SetEnabled(false)
	k.Right.SetEnabled(false)
	k.Tab.SetEnabled(false)
	k.Enter.SetEnabled(false)
	k.Esc.SetEnabled(false)
	k.Space.SetEnabled(false)

	k.DisableDashboardKeys()

	k.GlobalDashboard.SetEnabled(false)
	k.GlobalMentor.SetEnabled(false)
	k.GlobalSettings.SetEnabled(false)
	k.GlobalHelp.SetEnabled(false)
	k.GlobalTimer.SetEnabled(false)
	k.GlobalQuit.SetEnabled(false)

	k.CommandPalette.SetEnabled(false)
	k.Save.SetEnabled(false)
	k.Cancel.SetEnabled(false)
}
