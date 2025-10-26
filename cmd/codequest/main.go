// Package main is the entry point for the CodeQuest application.
// It initializes all components and launches the Bubble Tea TUI.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/AutumnsGrove/codequest/internal/config"
	"github.com/AutumnsGrove/codequest/internal/game"
	"github.com/AutumnsGrove/codequest/internal/storage"
	"github.com/AutumnsGrove/codequest/internal/ui"
	"github.com/AutumnsGrove/codequest/internal/watcher"
)

// Version information (set during build via ldflags)
var (
	Version    = "v0.1.0-beta" // Application version
	BuildTime  = "unknown"     // Build timestamp
	CommitHash = "unknown"     // Git commit hash
)

// CLI flags
var (
	showVersion = flag.Bool("version", false, "Show version information and exit")
	showHelp    = flag.Bool("help", false, "Show help message and exit")
)

func main() {
	// Step 1: Parse CLI flags
	flag.Parse()

	// Handle --version flag
	if *showVersion {
		fmt.Printf("CodeQuest %s\n", Version)
		fmt.Printf("Built: %s\n", BuildTime)
		fmt.Printf("Commit: %s\n", CommitHash)
		os.Exit(0)
	}

	// Handle --help flag
	if *showHelp {
		showHelpMessage()
		os.Exit(0)
	}

	// Step 2: Load or create default configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to load configuration: %v\n", err)
		fmt.Fprintf(os.Stderr, "   Config file: ~/.config/codequest/config.toml\n")
		os.Exit(1)
	}

	// Step 3: Initialize Skate storage (graceful error if missing)
	storageClient, err := storage.NewSkateClient()
	if err != nil {
		showSkateInstallInstructions()
		os.Exit(1)
	}

	// Step 4: Load or create character
	character, err := storageClient.LoadCharacter()
	if err != nil {
		// First run - create new character
		character = promptForCharacterCreation(cfg)
		if err := storageClient.SaveCharacter(character); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to save new character: %v\n", err)
			os.Exit(1)
		}
	}

	// Step 5: Load quests
	quests, err := storageClient.LoadQuests()
	if err != nil {
		// Non-fatal - start with empty quests
		quests = []*game.Quest{}
	}

	// Step 6: Create EventBus and register GameEventHandler
	eventBus := game.NewEventBus()

	// Create and start game event handler
	gameHandler, err := game.NewGameEventHandler(character, quests, eventBus, storageClient, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create game event handler: %v\n", err)
		os.Exit(1)
	}

	if err := gameHandler.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to start game event handler: %v\n", err)
		os.Exit(1)
	}

	// Step 7: Create context for managing watcher lifecycle
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start GitWatcher with context
	watcherManager, err := watcher.NewWatcherManager(eventBus, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create watcher manager: %v\n", err)
		os.Exit(1)
	}

	if err := watcherManager.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Warning: Failed to start git watcher: %v\n", err)
		fmt.Fprintf(os.Stderr, "   The app will run without automatic commit tracking.\n")
		// Non-fatal - continue without git watching
	}

	// Step 8: SessionTracker is initialized inside ui.NewModel()
	// (already handled by the UI layer)

	// Step 9: Create Bubble Tea Model
	model := ui.NewModel(storageClient, cfg, Version)

	// Step 10: Setup graceful shutdown
	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Spawn a goroutine to handle shutdown signals
	go func() {
		<-sigChan
		// Signal received - initiate graceful shutdown
		cancel() // Cancel watcher context

		// Stop game event handler
		if err := gameHandler.Stop(); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Warning: Error stopping game handler: %v\n", err)
		}

		// Stop watcher manager
		if err := watcherManager.Stop(); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Warning: Error stopping watcher: %v\n", err)
		}

		// Model cleanup is handled by the Bubble Tea program
		// Exit cleanly
		os.Exit(0)
	}()

	// Step 11: Run Bubble Tea program
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program and handle any errors
	if _, err := program.Run(); err != nil {
		// Ensure cleanup on error
		cancel()
		gameHandler.Stop()
		watcherManager.Stop()

		fmt.Fprintf(os.Stderr, "❌ Error running CodeQuest: %v\n", err)
		os.Exit(1)
	}

	// Normal exit - cleanup
	cancel()
	gameHandler.Stop()
	watcherManager.Stop()
	model.Cleanup()
}

// showHelpMessage displays usage information and available commands.
func showHelpMessage() {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("246"))

	fmt.Println(titleStyle.Render("CodeQuest - Gamified Developer Productivity RPG"))
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  codequest [flags]")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println("  --version    Show version information")
	fmt.Println("  --help       Show this help message")
	fmt.Println()
	fmt.Println("KEYBOARD SHORTCUTS:")
	fmt.Println("  Dashboard:")
	fmt.Println("    Q - Quest Board")
	fmt.Println("    C - Character Sheet")
	fmt.Println("    M - AI Mentor")
	fmt.Println("    S - Settings")
	fmt.Println()
	fmt.Println("  Global:")
	fmt.Println("    Alt+D - Dashboard")
	fmt.Println("    Alt+M - AI Mentor")
	fmt.Println("    Alt+S - Settings")
	fmt.Println("    Ctrl+T - Toggle session timer")
	fmt.Println("    Ctrl+S - Save game state")
	fmt.Println("    Ctrl+C - Quit")
	fmt.Println("    ? - Help overlay")
	fmt.Println()
	fmt.Println(helpStyle.Render("For more information, visit: https://github.com/AutumnsGrove/codequest"))
}

// showSkateInstallInstructions displays helpful instructions for installing Skate.
func showSkateInstallInstructions() {
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196"))

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117"))

	codeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("220")).
		Bold(true)

	fmt.Println(errorStyle.Render("❌ Skate not found"))
	fmt.Println()
	fmt.Println(infoStyle.Render("CodeQuest requires Skate for data persistence."))
	fmt.Println(infoStyle.Render("Skate is a simple key-value store from Charmbracelet."))
	fmt.Println()
	fmt.Println("To install Skate:")
	fmt.Println()
	fmt.Println("  " + codeStyle.Render("brew install charmbracelet/tap/skate"))
	fmt.Println()
	fmt.Println("After installation, run CodeQuest again.")
	fmt.Println()
	fmt.Println(infoStyle.Render("Learn more: https://github.com/charmbracelet/skate"))
}

// promptForCharacterCreation prompts the user to create a new character.
// This is called on first run when no saved character exists.
// It uses simple terminal I/O (not Bubble Tea) for the initial setup.
//
// Parameters:
//   - cfg: Application configuration (for default character name)
//
// Returns:
//   - *game.Character: A newly created character
func promptForCharacterCreation(cfg *config.Config) *game.Character {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117"))

	fmt.Println(titleStyle.Render("🎮 Welcome to CodeQuest!"))
	fmt.Println()
	fmt.Println("This appears to be your first time running CodeQuest.")
	fmt.Println("Let's create your character to get started!")
	fmt.Println()

	// Prompt for character name
	defaultName := cfg.Character.Name
	if defaultName == "" {
		defaultName = "CodeWarrior"
	}

	fmt.Print(promptStyle.Render(fmt.Sprintf("Enter your character name [%s]: ", defaultName)))

	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(input)

	// Use default name if no input provided
	if input == "" {
		input = defaultName
	}

	// Create character
	character := game.NewCharacter(input)

	fmt.Println()
	fmt.Printf("✓ Character '%s' created!\n", character.Name)
	fmt.Println()
	fmt.Println("Starting your coding adventure...")
	fmt.Println()

	return character
}
