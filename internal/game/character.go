// Package game contains the core game logic for CodeQuest
package game

import "time"

// Character represents the player in the game world
type Character struct {
	// Identity
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`

	// Core Stats
	Level         int `json:"level"`
	XP            int `json:"xp"`
	XPToNextLevel int `json:"xp_to_next_level"`

	// TODO: Add more fields as per specification
}

// NewCharacter creates a new character with starting stats
func NewCharacter(name string) *Character {
	return &Character{
		ID:            generateID(),
		Name:          name,
		CreatedAt:     time.Now(),
		Level:         1,
		XP:            0,
		XPToNextLevel: 100,
	}
}

// AddXP adds experience points and checks for level up
func (c *Character) AddXP(amount int) bool {
	// TODO: Implement XP addition and level up logic
	c.XP += amount
	return false // TODO: Return true if leveled up
}

// generateID creates a unique identifier for the character
func generateID() string {
	// TODO: Implement proper ID generation
	return "char_" + time.Now().Format("20060102150405")
}