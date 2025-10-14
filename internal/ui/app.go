// Package ui handles the terminal user interface for CodeQuest
package ui

// Model represents the main application model for Bubble Tea
type Model struct {
	// TODO: Add fields as per specification
	width  int
	height int
	ready  bool
}

// NewModel creates a new UI model
func NewModel() Model {
	return Model{
		ready: false,
	}
}

// TODO: Implement Bubble Tea interface methods (Init, Update, View)