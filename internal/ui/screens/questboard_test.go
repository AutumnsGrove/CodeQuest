package screens

import (
	"strings"
	"testing"
	"time"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// TestRenderQuestBoard tests the main RenderQuestBoard function.
func TestRenderQuestBoard(t *testing.T) {
	// Create test character
	character := game.NewCharacter("TestHero")

	// Create test quests
	quests := []*game.Quest{
		createTestQuest("Quest 1", game.QuestAvailable),
		createTestQuest("Quest 2", game.QuestActive),
		createTestQuest("Quest 3", game.QuestCompleted),
	}

	tests := []struct {
		name          string
		character     *game.Character
		quests        []*game.Quest
		selectedIndex int
		filter        QuestFilter
		width         int
		height        int
		wantContains  []string
	}{
		{
			name:          "renders all quests",
			character:     character,
			quests:        quests,
			selectedIndex: -1,
			filter:        FilterAll,
			width:         80,
			height:        40,
			wantContains: []string{
				"Quest Board",
				"Quest 1",
				"Quest 2",
				"Quest 3",
			},
		},
		{
			name:          "filters available quests",
			character:     character,
			quests:        quests,
			selectedIndex: -1,
			filter:        FilterAvailable,
			width:         80,
			height:        40,
			wantContains: []string{
				"Quest 1",
			},
		},
		{
			name:          "filters active quests",
			character:     character,
			quests:        quests,
			selectedIndex: -1,
			filter:        FilterActive,
			width:         80,
			height:        40,
			wantContains: []string{
				"Quest 2",
			},
		},
		{
			name:          "filters completed quests",
			character:     character,
			quests:        quests,
			selectedIndex: -1,
			filter:        FilterCompleted,
			width:         80,
			height:        40,
			wantContains: []string{
				"Quest 3",
			},
		},
		{
			name:          "highlights selected quest",
			character:     character,
			quests:        quests,
			selectedIndex: 0,
			filter:        FilterAll,
			width:         80,
			height:        40,
			wantContains: []string{
				"▶",
			},
		},
		{
			name:          "handles nil character",
			character:     nil,
			quests:        quests,
			selectedIndex: -1,
			filter:        FilterAll,
			width:         80,
			height:        40,
			wantContains: []string{
				"Quest Board",
			},
		},
		{
			name:          "handles empty quest list",
			character:     character,
			quests:        []*game.Quest{},
			selectedIndex: -1,
			filter:        FilterAll,
			width:         80,
			height:        40,
			wantContains: []string{
				"No quests available",
			},
		},
		{
			name:          "handles narrow terminal",
			character:     character,
			quests:        quests,
			selectedIndex: -1,
			filter:        FilterAll,
			width:         40,
			height:        24,
			wantContains: []string{
				"Quest Board",
			},
		},
		{
			name:          "handles small terminal",
			character:     character,
			quests:        quests,
			selectedIndex: -1,
			filter:        FilterAll,
			width:         20,
			height:        10,
			wantContains: []string{
				"Quest Board",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := RenderQuestBoard(tt.character, tt.quests, tt.selectedIndex, tt.filter, tt.width, tt.height)

			// Check that output is not empty
			if output == "" {
				t.Error("RenderQuestBoard returned empty string")
			}

			// Check for expected content
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("RenderQuestBoard output does not contain %q", want)
				}
			}
		})
	}
}

// TestFilterQuests tests the quest filtering logic.
func TestFilterQuests(t *testing.T) {
	quests := []*game.Quest{
		createTestQuest("Available 1", game.QuestAvailable),
		createTestQuest("Available 2", game.QuestAvailable),
		createTestQuest("Active 1", game.QuestActive),
		createTestQuest("Active 2", game.QuestActive),
		createTestQuest("Completed 1", game.QuestCompleted),
		createTestQuest("Failed 1", game.QuestFailed),
	}

	tests := []struct {
		name       string
		quests     []*game.Quest
		filter     QuestFilter
		wantCount  int
		wantTitles []string
	}{
		{
			name:       "filter all shows all quests",
			quests:     quests,
			filter:     FilterAll,
			wantCount:  6,
			wantTitles: []string{"Available 1", "Available 2", "Active 1", "Active 2", "Completed 1", "Failed 1"},
		},
		{
			name:       "filter available shows only available",
			quests:     quests,
			filter:     FilterAvailable,
			wantCount:  2,
			wantTitles: []string{"Available 1", "Available 2"},
		},
		{
			name:       "filter active shows only active",
			quests:     quests,
			filter:     FilterActive,
			wantCount:  2,
			wantTitles: []string{"Active 1", "Active 2"},
		},
		{
			name:       "filter completed shows only completed",
			quests:     quests,
			filter:     FilterCompleted,
			wantCount:  1,
			wantTitles: []string{"Completed 1"},
		},
		{
			name:       "handles empty quest list",
			quests:     []*game.Quest{},
			filter:     FilterAll,
			wantCount:  0,
			wantTitles: []string{},
		},
		{
			name:       "handles nil quest list",
			quests:     nil,
			filter:     FilterAll,
			wantCount:  0,
			wantTitles: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := filterQuests(tt.quests, tt.filter)

			if len(filtered) != tt.wantCount {
				t.Errorf("filterQuests() returned %d quests, want %d", len(filtered), tt.wantCount)
			}

			// Check quest titles
			for i, wantTitle := range tt.wantTitles {
				if i >= len(filtered) {
					t.Errorf("filterQuests() missing quest %q", wantTitle)
					continue
				}
				if filtered[i].Title != wantTitle {
					t.Errorf("filterQuests() quest[%d].Title = %q, want %q", i, filtered[i].Title, wantTitle)
				}
			}
		})
	}
}

// TestRenderFilterTabs tests the filter tab rendering.
func TestRenderFilterTabs(t *testing.T) {
	quests := []*game.Quest{
		createTestQuest("Available", game.QuestAvailable),
		createTestQuest("Active", game.QuestActive),
		createTestQuest("Completed", game.QuestCompleted),
	}

	tests := []struct {
		name         string
		filter       QuestFilter
		quests       []*game.Quest
		width        int
		wantContains []string
	}{
		{
			name:   "shows all filter counts",
			filter: FilterAll,
			quests: quests,
			width:  80,
			wantContains: []string{
				"All (3)",
				"Available (1)",
				"Active (1)",
				"Completed (1)",
			},
		},
		{
			name:         "shows correct counts with empty list",
			filter:       FilterAll,
			quests:       []*game.Quest{},
			width:        80,
			wantContains: []string{"All (0)", "Available (0)", "Active (0)", "Completed (0)"},
		},
		{
			name:   "highlights all tab when selected",
			filter: FilterAll,
			quests: quests,
			width:  80,
			wantContains: []string{
				"All",
			},
		},
		{
			name:   "highlights available tab when selected",
			filter: FilterAvailable,
			quests: quests,
			width:  80,
			wantContains: []string{
				"Available",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := renderFilterTabs(tt.filter, tt.quests, tt.width)

			if output == "" {
				t.Error("renderFilterTabs returned empty string")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("renderFilterTabs output does not contain %q", want)
				}
			}
		})
	}
}

// TestRenderQuestList tests quest list rendering.
func TestRenderQuestList(t *testing.T) {
	quests := []*game.Quest{
		createTestQuest("Quest 1", game.QuestAvailable),
		createTestQuest("Quest 2", game.QuestActive),
		createTestQuest("Quest 3", game.QuestCompleted),
	}

	tests := []struct {
		name          string
		quests        []*game.Quest
		selectedIndex int
		width         int
		maxHeight     int
		wantContains  []string
	}{
		{
			name:          "renders all quests",
			quests:        quests,
			selectedIndex: -1,
			width:         80,
			maxHeight:     30,
			wantContains: []string{
				"Quest 1",
				"Quest 2",
				"Quest 3",
				"Available Quests",
				"Active Quests",
				"Completed Quests",
			},
		},
		{
			name:          "highlights selected quest",
			quests:        quests,
			selectedIndex: 1,
			width:         80,
			maxHeight:     30,
			wantContains:  []string{"▶"},
		},
		{
			name:          "handles narrow width",
			quests:        quests,
			selectedIndex: -1,
			width:         40,
			maxHeight:     20,
			wantContains:  []string{"Quest 1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := renderQuestList(tt.quests, tt.selectedIndex, tt.width, tt.maxHeight)

			if output == "" {
				t.Error("renderQuestList returned empty string")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("renderQuestList output does not contain %q", want)
				}
			}
		})
	}
}

// TestRenderQuestCard tests individual quest card rendering.
func TestRenderQuestCard(t *testing.T) {
	tests := []struct {
		name         string
		quest        *game.Quest
		selected     bool
		width        int
		wantContains []string
	}{
		{
			name:         "renders available quest",
			quest:        createTestQuest("Test Quest", game.QuestAvailable),
			selected:     false,
			width:        60,
			wantContains: []string{"Test Quest", "Reward:", "XP", "Required Level:"},
		},
		{
			name:         "renders active quest",
			quest:        createActiveTestQuest("Active Quest"),
			selected:     false,
			width:        60,
			wantContains: []string{"Active Quest", "Progress:", "Started:"},
		},
		{
			name:         "renders completed quest",
			quest:        createCompletedTestQuest("Completed Quest"),
			selected:     false,
			width:        60,
			wantContains: []string{"Completed Quest", "XP Earned:", "Completed:"},
		},
		{
			name:         "shows selection indicator",
			quest:        createTestQuest("Selected", game.QuestAvailable),
			selected:     true,
			width:        60,
			wantContains: []string{"▶", "Selected"},
		},
		{
			name:         "truncates long description",
			quest:        createLongDescriptionQuest(),
			selected:     false,
			width:        40,
			wantContains: []string{"..."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := renderQuestCard(tt.quest, tt.selected, tt.width)

			if output == "" {
				t.Error("renderQuestCard returned empty string")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("renderQuestCard output does not contain %q", want)
				}
			}
		})
	}
}

// TestRenderEmptyQuestList tests empty quest list rendering.
func TestRenderEmptyQuestList(t *testing.T) {
	tests := []struct {
		name         string
		filter       QuestFilter
		width        int
		wantContains []string
	}{
		{
			name:         "shows message for no quests with all filter",
			filter:       FilterAll,
			width:        80,
			wantContains: []string{"No quests available"},
		},
		{
			name:         "shows message for no available quests",
			filter:       FilterAvailable,
			width:        80,
			wantContains: []string{"No available quests"},
		},
		{
			name:         "shows message for no active quests",
			filter:       FilterActive,
			width:        80,
			wantContains: []string{"No active quests"},
		},
		{
			name:         "shows message for no completed quests",
			filter:       FilterCompleted,
			width:        80,
			wantContains: []string{"No completed quests"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := renderEmptyQuestList(tt.filter, tt.width)

			if output == "" {
				t.Error("renderEmptyQuestList returned empty string")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("renderEmptyQuestList output does not contain %q", want)
				}
			}
		})
	}
}

// TestRenderQuestBoardFooter tests footer rendering.
func TestRenderQuestBoardFooter(t *testing.T) {
	tests := []struct {
		name         string
		width        int
		wantContains []string
	}{
		{
			name:         "shows all key bindings",
			width:        80,
			wantContains: []string{"Navigate", "Start/View", "Filter", "Back"},
		},
		{
			name:         "handles narrow width",
			width:        40,
			wantContains: []string{"Navigate"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := renderQuestBoardFooter(tt.width)

			if output == "" {
				t.Error("renderQuestBoardFooter returned empty string")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("renderQuestBoardFooter output does not contain %q", want)
				}
			}
		})
	}
}

// TestRenderQuestSection tests quest section rendering.
func TestRenderQuestSection(t *testing.T) {
	quests := []*game.Quest{
		createTestQuest("Quest 1", game.QuestAvailable),
		createTestQuest("Quest 2", game.QuestAvailable),
	}

	tests := []struct {
		name          string
		title         string
		quests        []*game.Quest
		selectedIndex int
		offset        int
		width         int
		wantContains  []string
	}{
		{
			name:          "renders section title and quests",
			title:         "Test Section",
			quests:        quests,
			selectedIndex: -1,
			offset:        0,
			width:         80,
			wantContains:  []string{"Test Section", "Quest 1", "Quest 2"},
		},
		{
			name:          "highlights selected quest with offset",
			title:         "Test Section",
			quests:        quests,
			selectedIndex: 1,
			offset:        1,
			width:         80,
			wantContains:  []string{"▶"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := renderQuestSection(tt.title, tt.quests, tt.selectedIndex, tt.offset, tt.width)

			if output == "" {
				t.Error("renderQuestSection returned empty string")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("renderQuestSection output does not contain %q", want)
				}
			}
		})
	}
}

// TestRenderAvailableQuestInfo tests available quest info rendering.
func TestRenderAvailableQuestInfo(t *testing.T) {
	quest := createTestQuest("Test", game.QuestAvailable)
	output := renderAvailableQuestInfo(quest)

	if output == "" {
		t.Error("renderAvailableQuestInfo returned empty string")
	}

	wantContains := []string{"Reward:", "XP", "Required Level:", "Enter"}
	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("renderAvailableQuestInfo output does not contain %q", want)
		}
	}
}

// TestRenderActiveQuestInfo tests active quest info rendering.
func TestRenderActiveQuestInfo(t *testing.T) {
	quest := createActiveTestQuest("Test")
	output := renderActiveQuestInfo(quest, 40)

	if output == "" {
		t.Error("renderActiveQuestInfo returned empty string")
	}

	wantContains := []string{"Progress:", "Started:", "Reward:"}
	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("renderActiveQuestInfo output does not contain %q", want)
		}
	}
}

// TestRenderCompletedQuestInfo tests completed quest info rendering.
func TestRenderCompletedQuestInfo(t *testing.T) {
	quest := createCompletedTestQuest("Test")
	output := renderCompletedQuestInfo(quest)

	if output == "" {
		t.Error("renderCompletedQuestInfo returned empty string")
	}

	wantContains := []string{"XP Earned:", "Completed:", "Duration:"}
	for _, want := range wantContains {
		if !strings.Contains(output, want) {
			t.Errorf("renderCompletedQuestInfo output does not contain %q", want)
		}
	}
}

// BenchmarkRenderQuestBoard benchmarks the main rendering function.
func BenchmarkRenderQuestBoard(b *testing.B) {
	character := game.NewCharacter("BenchHero")
	quests := make([]*game.Quest, 10)
	for i := 0; i < 10; i++ {
		quests[i] = createTestQuest("Quest", game.QuestAvailable)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RenderQuestBoard(character, quests, 0, FilterAll, 80, 40)
	}
}

// BenchmarkFilterQuests benchmarks quest filtering.
func BenchmarkFilterQuests(b *testing.B) {
	quests := make([]*game.Quest, 100)
	for i := 0; i < 100; i++ {
		status := game.QuestAvailable
		if i%3 == 0 {
			status = game.QuestActive
		} else if i%3 == 1 {
			status = game.QuestCompleted
		}
		quests[i] = createTestQuest("Quest", status)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filterQuests(quests, FilterActive)
	}
}

// Helper functions for creating test quests

// createTestQuest creates a basic test quest with the given title and status.
func createTestQuest(title string, status game.QuestStatus) *game.Quest {
	quest := game.NewQuest(
		title,
		"Test quest description",
		game.QuestTypeCommit,
		5,    // target
		100,  // xpReward
		1,    // requiredLevel
	)
	quest.Status = status
	return quest
}

// createActiveTestQuest creates an active quest with progress.
func createActiveTestQuest(title string) *game.Quest {
	quest := createTestQuest(title, game.QuestActive)
	now := time.Now().Add(-2 * time.Hour)
	quest.StartedAt = &now
	quest.Current = 3
	quest.Progress = 0.6
	return quest
}

// createCompletedTestQuest creates a completed quest.
func createCompletedTestQuest(title string) *game.Quest {
	quest := createTestQuest(title, game.QuestCompleted)
	started := time.Now().Add(-5 * time.Hour)
	completed := time.Now().Add(-1 * time.Hour)
	quest.StartedAt = &started
	quest.CompletedAt = &completed
	quest.Current = 5
	quest.Progress = 1.0
	return quest
}

// createLongDescriptionQuest creates a quest with a very long description.
func createLongDescriptionQuest() *game.Quest {
	quest := game.NewQuest(
		"Long Quest",
		"This is a very long description that should be truncated when displayed in a narrow terminal. It contains many words and should exceed the maximum description length allowed in the UI rendering.",
		game.QuestTypeCommit,
		5,
		100,
		1,
	)
	return quest
}
