package screens

import (
	"strings"
	"testing"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// TestRenderSettings tests the main settings screen rendering function.
func TestRenderSettings(t *testing.T) {
	tests := []struct {
		name      string
		character *game.Character
		width     int
		height    int
		wantEmpty bool
	}{
		{
			name:      "normal settings screen",
			character: createTestCharacter(),
			width:     100,
			height:    40,
			wantEmpty: false,
		},
		{
			name:      "narrow terminal",
			character: createTestCharacter(),
			width:     60,
			height:    30,
			wantEmpty: false,
		},
		{
			name:      "nil character",
			character: nil,
			width:     80,
			height:    40,
			wantEmpty: false, // Should still render
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderSettings(tt.character, tt.width, tt.height)

			// Should never return empty string
			if result == "" && !tt.wantEmpty {
				t.Error("RenderSettings() returned empty string")
			}

			// Should contain "Settings" in header
			if !strings.Contains(result, "Settings") {
				t.Error("RenderSettings() should contain 'Settings' in header")
			}

			// Should contain all setting categories
			categories := []string{
				"Game Settings",
				"UI Settings",
				"AI Settings",
				"Git Settings",
				"Debug Settings",
			}

			for _, category := range categories {
				if !strings.Contains(result, category) {
					t.Errorf("RenderSettings() should contain category %q", category)
				}
			}
		})
	}
}

// TestRenderSettingsPanel tests the settings panel rendering.
func TestRenderSettingsPanel(t *testing.T) {
	result := renderSettingsPanel(100)

	if result == "" {
		t.Error("renderSettingsPanel() returned empty string")
	}

	// Should contain all categories
	expectedStrings := []string{
		"Game Settings",
		"UI Settings",
		"AI Settings",
		"Git Settings",
		"Debug Settings",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderSettingsPanel() should contain %q", expected)
		}
	}
}

// TestRenderGameSettings tests game settings section rendering.
func TestRenderGameSettings(t *testing.T) {
	result := renderGameSettings()

	if result == "" {
		t.Error("renderGameSettings() returned empty string")
	}

	// Should contain game setting items
	expectedStrings := []string{
		"Game Settings",
		"Difficulty:",
		"XP Multiplier:",
		"Auto-save:",
		"Quest Notifications:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderGameSettings() should contain %q", expected)
		}
	}

	// Should show current values
	if !strings.Contains(result, "Normal") {
		t.Error("renderGameSettings() should show difficulty value")
	}

	if !strings.Contains(result, "1.0x") {
		t.Error("renderGameSettings() should show XP multiplier value")
	}
}

// TestRenderUISettings tests UI settings section rendering.
func TestRenderUISettings(t *testing.T) {
	result := renderUISettings()

	if result == "" {
		t.Error("renderUISettings() returned empty string")
	}

	// Should contain UI setting items
	expectedStrings := []string{
		"UI Settings",
		"Color Theme:",
		"Animations:",
		"Compact Mode:",
		"Show Help Hints:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderUISettings() should contain %q", expected)
		}
	}

	// Should show current theme
	if !strings.Contains(result, "Default") {
		t.Error("renderUISettings() should show theme value")
	}
}

// TestRenderAISettings tests AI settings section rendering.
func TestRenderAISettings(t *testing.T) {
	result := renderAISettings()

	if result == "" {
		t.Error("renderAISettings() returned empty string")
	}

	// Should contain AI setting items
	expectedStrings := []string{
		"AI Settings",
		"Primary Provider:",
		"Status:",
		"Fallback Provider:",
		"Rate Limiting:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderAISettings() should contain %q", expected)
		}
	}

	// Should show provider name
	if !strings.Contains(result, "Crush") {
		t.Error("renderAISettings() should show primary provider")
	}

	// Should show status
	if !strings.Contains(result, "Online") {
		t.Error("renderAISettings() should show provider status")
	}
}

// TestRenderGitSettings tests Git settings section rendering.
func TestRenderGitSettings(t *testing.T) {
	result := renderGitSettings()

	if result == "" {
		t.Error("renderGitSettings() returned empty string")
	}

	// Should contain Git setting items
	expectedStrings := []string{
		"Git Settings",
		"Auto-detect commits:",
		"Repository:",
		"Watch mode:",
		"Commit XP formula:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderGitSettings() should contain %q", expected)
		}
	}

	// Should show auto-detect status
	if !strings.Contains(result, "Enabled") {
		t.Error("renderGitSettings() should show auto-detect status")
	}
}

// TestRenderDebugSettings tests debug settings section rendering.
func TestRenderDebugSettings(t *testing.T) {
	result := renderDebugSettings()

	if result == "" {
		t.Error("renderDebugSettings() returned empty string")
	}

	// Should contain debug setting items
	expectedStrings := []string{
		"Debug Settings",
		"Log Level:",
		"Developer Mode:",
		"Show Debug UI:",
		"Performance Monitor:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderDebugSettings() should contain %q", expected)
		}
	}

	// Should show log level
	if !strings.Contains(result, "Info") {
		t.Error("renderDebugSettings() should show log level")
	}
}

// TestRenderSettingsFooter tests footer rendering.
func TestRenderSettingsFooter(t *testing.T) {
	result := renderSettingsFooter(80)

	if result == "" {
		t.Error("renderSettingsFooter() returned empty string")
	}

	// Should contain info message
	if !strings.Contains(result, "read-only") {
		t.Error("renderSettingsFooter() should contain read-only info")
	}

	// Should contain key bindings
	expectedStrings := []string{
		"Alt+Q",
		"Dashboard",
		"Esc",
		"Back",
		"?",
		"Help",
		"Ctrl+S",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("renderSettingsFooter() should contain %q", expected)
		}
	}
}

// TestGetSettingsByCategory tests filtering settings by category.
func TestGetSettingsByCategory(t *testing.T) {
	tests := []struct {
		name     string
		category SettingsCategory
		wantMin  int // Minimum expected items
	}{
		{
			name:     "game settings",
			category: CategoryGame,
			wantMin:  2, // At least difficulty and XP multiplier
		},
		{
			name:     "UI settings",
			category: CategoryUI,
			wantMin:  2, // At least theme and animations
		},
		{
			name:     "AI settings",
			category: CategoryAI,
			wantMin:  1, // At least primary provider
		},
		{
			name:     "Git settings",
			category: CategoryGit,
			wantMin:  1, // At least auto-detect
		},
		{
			name:     "Debug settings",
			category: CategoryDebug,
			wantMin:  1, // At least log level
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := getSettingsByCategory(tt.category)

			if len(items) < tt.wantMin {
				t.Errorf("getSettingsByCategory(%v) returned %d items, want at least %d",
					tt.category, len(items), tt.wantMin)
			}

			// All items should belong to the requested category
			for _, item := range items {
				if item.Category != tt.category {
					t.Errorf("getSettingsByCategory(%v) returned item with category %v",
						tt.category, item.Category)
				}
			}
		})
	}
}

// TestRenderSettingItem tests individual setting item rendering.
func TestRenderSettingItem(t *testing.T) {
	testItem := SettingItem{
		Label:       "Test Setting",
		Key:         "test.setting",
		Value:       "test_value",
		ValueType:   "string",
		Description: "A test setting",
		Category:    CategoryGame,
	}

	tests := []struct {
		name     string
		item     SettingItem
		selected bool
	}{
		{
			name:     "unselected item",
			item:     testItem,
			selected: false,
		},
		{
			name:     "selected item",
			item:     testItem,
			selected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderSettingItem(tt.item, tt.selected)

			if result == "" {
				t.Error("renderSettingItem() returned empty string")
			}

			// Should contain label
			if !strings.Contains(result, tt.item.Label) {
				t.Errorf("renderSettingItem() should contain label %q", tt.item.Label)
			}

			// Selected items should have indicator
			if tt.selected && !strings.Contains(result, "▶") {
				t.Error("renderSettingItem() for selected item should contain selection indicator")
			}
		})
	}
}

// TestSettingsItems validates the settings items structure.
func TestSettingsItems(t *testing.T) {
	if len(settingsItems) == 0 {
		t.Error("settingsItems should not be empty")
	}

	// Verify each item has required fields
	for i, item := range settingsItems {
		if item.Label == "" {
			t.Errorf("settingsItems[%d] has empty Label", i)
		}

		if item.Key == "" {
			t.Errorf("settingsItems[%d] has empty Key", i)
		}

		if item.ValueType == "" {
			t.Errorf("settingsItems[%d] has empty ValueType", i)
		}

		// If type is "choice", must have choices
		if item.ValueType == "choice" && len(item.Choices) == 0 {
			t.Errorf("settingsItems[%d] has type 'choice' but no Choices", i)
		}
	}

	// Verify we have items for each category
	categories := []SettingsCategory{
		CategoryGame,
		CategoryUI,
		CategoryAI,
		CategoryGit,
		CategoryDebug,
	}

	for _, cat := range categories {
		items := getSettingsByCategory(cat)
		if len(items) == 0 {
			t.Errorf("No settings items found for category %v", cat)
		}
	}
}

// TestSettingItemTypes validates different setting value types.
func TestSettingItemTypes(t *testing.T) {
	// Find examples of each type
	var boolItem, stringItem, intItem, floatItem, choiceItem *SettingItem

	for i := range settingsItems {
		item := &settingsItems[i]
		switch item.ValueType {
		case "bool":
			if boolItem == nil {
				boolItem = item
			}
		case "string":
			if stringItem == nil {
				stringItem = item
			}
		case "int":
			if intItem == nil {
				intItem = item
			}
		case "float":
			if floatItem == nil {
				floatItem = item
			}
		case "choice":
			if choiceItem == nil {
				choiceItem = item
			}
		}
	}

	// Verify we have examples of all types
	if boolItem == nil {
		t.Error("No boolean setting items found")
	}

	if choiceItem == nil {
		t.Error("No choice setting items found")
	}

	// Verify choice items have valid choices
	if choiceItem != nil && len(choiceItem.Choices) == 0 {
		t.Error("Choice item should have at least one choice")
	}
}
