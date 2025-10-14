// Package components provides reusable UI components for CodeQuest screens.
// This file contains tests for the header component.
package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/AutumnsGrove/codequest/internal/game"
)

// TestRenderHeader verifies the header renders correctly with different inputs.
func TestRenderHeader(t *testing.T) {
	tests := []struct {
		name       string
		screenName string
		char       *game.Character
		width      int
		wantEmpty  bool
	}{
		{
			name:       "normal width with character",
			screenName: "Dashboard",
			char:       game.NewCharacter("TestHero"),
			width:      80,
			wantEmpty:  false,
		},
		{
			name:       "normal width without character",
			screenName: "Dashboard",
			char:       nil,
			width:      80,
			wantEmpty:  false,
		},
		{
			name:       "narrow width (minimal mode)",
			screenName: "Dashboard",
			char:       game.NewCharacter("TestHero"),
			width:      30,
			wantEmpty:  false,
		},
		{
			name:       "very narrow width",
			screenName: "Dashboard",
			char:       nil,
			width:      20,
			wantEmpty:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderHeader(tt.screenName, tt.char, tt.width)

			if tt.wantEmpty {
				if got != "" {
					t.Errorf("RenderHeader() should be empty, got %q", got)
				}
				return
			}

			// Verify we got some output
			if got == "" {
				t.Error("RenderHeader() returned empty string")
			}

			// Verify the screen name appears in the output
			if !strings.Contains(got, tt.screenName) {
				t.Errorf("RenderHeader() missing screen name %q in output", tt.screenName)
			}

			// Verify CodeQuest title appears
			if !strings.Contains(got, "CodeQuest") {
				t.Error("RenderHeader() missing CodeQuest title")
			}

			// If character provided, verify character name appears
			if tt.char != nil {
				if !strings.Contains(got, tt.char.Name) {
					t.Errorf("RenderHeader() missing character name %q", tt.char.Name)
				}
			}
		})
	}
}

// TestRenderLeftSection verifies the left section renders correctly.
func TestRenderLeftSection(t *testing.T) {
	got := renderLeftSection()

	// Should not be empty
	if got == "" {
		t.Error("renderLeftSection() returned empty string")
	}

	// Should contain icon and title
	if !strings.Contains(got, "CodeQuest") {
		t.Error("renderLeftSection() missing CodeQuest title")
	}
}

// TestRenderCenterSection verifies the center section renders correctly.
func TestRenderCenterSection(t *testing.T) {
	tests := []struct {
		name       string
		screenName string
		want       string
	}{
		{
			name:       "dashboard",
			screenName: "Dashboard",
			want:       "Dashboard",
		},
		{
			name:       "quest board",
			screenName: "Quest Board",
			want:       "Quest Board",
		},
		{
			name:       "settings",
			screenName: "Settings",
			want:       "Settings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderCenterSection(tt.screenName)

			// Should contain the screen name
			if !strings.Contains(got, tt.want) {
				t.Errorf("renderCenterSection() missing %q", tt.want)
			}

			// Should contain brackets
			if !strings.Contains(got, "[") || !strings.Contains(got, "]") {
				t.Error("renderCenterSection() missing brackets")
			}
		})
	}
}

// TestRenderRightSection verifies the right section renders correctly.
func TestRenderRightSection(t *testing.T) {
	tests := []struct {
		name string
		char *game.Character
		want string
	}{
		{
			name: "with character",
			char: game.NewCharacter("Hero"),
			want: "Hero",
		},
		{
			name: "nil character",
			char: nil,
			want: "No Character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderRightSection(tt.char)

			// Should contain expected text
			if !strings.Contains(got, tt.want) {
				t.Errorf("renderRightSection() missing %q in output", tt.want)
			}

			// If character exists, should contain level
			if tt.char != nil {
				if !strings.Contains(got, "Lvl") {
					t.Error("renderRightSection() missing level indicator")
				}
			}
		})
	}
}

// TestJoinHeaderSections verifies sections are joined correctly.
func TestJoinHeaderSections(t *testing.T) {
	tests := []struct {
		name   string
		left   string
		center string
		right  string
		width  int
	}{
		{
			name:   "normal width",
			left:   "Left",
			center: "Center",
			right:  "Right",
			width:  80,
		},
		{
			name:   "narrow width",
			left:   "L",
			center: "C",
			right:  "R",
			width:  20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := joinHeaderSections(tt.left, tt.center, tt.right, tt.width)

			// Should contain all sections
			if !strings.Contains(got, tt.left) {
				t.Error("joinHeaderSections() missing left section")
			}
			if !strings.Contains(got, tt.center) {
				t.Error("joinHeaderSections() missing center section")
			}
			if !strings.Contains(got, tt.right) {
				t.Error("joinHeaderSections() missing right section")
			}
		})
	}
}

// TestWrapHeader verifies the header is wrapped correctly.
func TestWrapHeader(t *testing.T) {
	content := "Test Header Content"
	width := 80

	got := wrapHeader(content, width)

	// Should contain the content
	if !strings.Contains(got, content) {
		t.Error("wrapHeader() missing content")
	}

	// Should not be empty
	if got == "" {
		t.Error("wrapHeader() returned empty string")
	}
}

// TestRenderMinimalHeader verifies minimal header for narrow terminals.
func TestRenderMinimalHeader(t *testing.T) {
	tests := []struct {
		name       string
		screenName string
		width      int
	}{
		{
			name:       "very narrow",
			screenName: "Dashboard",
			width:      20,
		},
		{
			name:       "minimal width",
			screenName: "Settings",
			width:      30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderMinimalHeader(tt.screenName, tt.width)

			// Should contain screen name
			if !strings.Contains(got, tt.screenName) {
				t.Errorf("renderMinimalHeader() missing screen name %q", tt.screenName)
			}

			// Should not be empty
			if got == "" {
				t.Error("renderMinimalHeader() returned empty string")
			}
		})
	}
}

// BenchmarkRenderHeader benchmarks the header rendering performance.
func BenchmarkRenderHeader(b *testing.B) {
	char := game.NewCharacter("BenchHero")
	screenName := "Dashboard"
	width := 80

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderHeader(screenName, char, width)
	}
}

// TestHeaderWidth verifies the header respects width constraints.
func TestHeaderWidth(t *testing.T) {
	char := game.NewCharacter("TestHero")
	screenName := "Dashboard"

	tests := []struct {
		name  string
		width int
	}{
		{"80 columns", 80},
		{"120 columns", 120},
		{"40 columns", 40},
		{"30 columns", 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderHeader(screenName, char, tt.width)

			// Measure the actual width (accounting for ANSI codes)
			actualWidth := lipgloss.Width(got)

			// The actual width should be <= the requested width
			// (allowing some margin for borders and padding)
			if actualWidth > tt.width+10 {
				t.Errorf("RenderHeader() width %d exceeds requested %d (allowing 10 char margin)",
					actualWidth, tt.width)
			}
		})
	}
}
