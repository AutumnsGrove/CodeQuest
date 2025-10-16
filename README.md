# CodeQuest 🎮⚔️

> Transform your coding sessions into an epic RPG adventure

CodeQuest is a terminal-based gamified developer productivity tool that turns your Git commits into XP, your coding sessions into quests, and your progress into an RPG character journey. Built with Go and the beautiful Charmbracelet ecosystem.

**⚠️ STATUS: IN DEVELOPMENT** - Core systems implemented, main application wiring in progress. See [DEVELOPMENT_STATUS.md](DEVELOPMENT_STATUS.md) for details.

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/AutumnsGrove/codequest)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## ✨ Features

- 🎯 **Quest System**: Turn your coding tasks into epic quests
- ⚡ **Real-time XP**: Earn experience points from every commit
- 📊 **Character Progression**: Level up and increase your stats (CodePower, Wisdom, Agility)
- 🤖 **AI Mentor**: Get coding help from Crush, Mods, or Claude
- ⏱️ **Session Tracking**: Monitor your coding time with Ctrl+T
- 🔥 **Daily Streaks**: Track consecutive days of activity
- 📈 **Beautiful Dashboard**: TUI showing all your stats and progress
- 💾 **Auto-save**: All progress persists between sessions

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** (for building from source)
- **Git** (for commit tracking)
- **Skate** (for data persistence) - Install: `brew install charmbracelet/tap/skate`

**Optional (for AI mentor):**
- **Crush** (OpenRouter): Get API key at [openrouter.ai](https://openrouter.ai/keys)
- **Mods** (local): Install via `brew install charmbracelet/tap/mods`

### Installation

#### Option 1: From Source (Recommended)

```bash
# Clone the repository
git clone https://github.com/AutumnsGrove/codequest.git
cd codequest

# Build the binary
make build

# Run CodeQuest
./build/codequest
```

#### Option 2: Go Install

```bash
go install github.com/AutumnsGrove/codequest/cmd/codequest@latest
```

### First Run

**Note:** The application is currently being wired up. The binary builds successfully, but the interactive TUI is not yet launched. See [DEVELOPMENT_STATUS.md](DEVELOPMENT_STATUS.md) for implementation progress.

When the application is complete, launching CodeQuest will:

1. **Create Your Character** - Choose your name
2. **Accept a Quest** - Browse the Quest Board (press `Q`)
3. **Start Coding** - Make commits in your Git repository
4. **Watch Your Progress** - Earn XP and level up!

```bash
# Planned workflow (not yet functional)
codequest                    # Will launch the game (currently shows placeholder)

# In another terminal, work on your project
cd ~/projects/my-project
git add .
git commit -m "feat: Add awesome feature"

# Return to CodeQuest to see your XP and quest progress update!
```

## 📖 Configuration

CodeQuest creates a config file at `~/.config/codequest/config.toml` on first run.

### Basic Configuration

```toml
[character]
name = "YourName"

[game]
difficulty = "normal"  # easy, normal, hard

[git]
auto_detect_repos = true
watch_paths = ["~/projects"]  # Directories to watch for commits

[ai.mentor]
provider = "crush"
model_complex = "openrouter/kimi/k2-0925"
model_simple = "openrouter/deepseek/glm-4.5-air"
temperature = 0.7

[ai.review]
provider = "mods"
model_primary = "qwen3:30b"
auto_review = true
```

### AI Provider Setup

CodeQuest supports three AI providers with automatic fallback: **Crush** (OpenRouter) → **Mods** (Local) → **Claude** (Anthropic API).

#### Crush (OpenRouter) - Online Models

Crush provides access to various online AI models:

```bash
# Store API key securely in Skate (never in plaintext)
skate set codequest.openrouter_api_key "YOUR_API_KEY_HERE"
```

Get your key at [openrouter.ai/keys](https://openrouter.ai/keys)

**Example models in config:**
```toml
[ai.mentor]
provider = "crush"
model_complex = "openrouter/kimi/k2-0925"           # For complex queries
model_simple = "openrouter/deepseek/glm-4.5-air"    # For quick questions
temperature = 0.7

# Other OpenRouter model options:
# - "openrouter/anthropic/claude-3.5-sonnet"
# - "openrouter/openai/gpt-4-turbo"
# - "openrouter/google/gemini-pro"
```

#### Mods (Local LLM) - Offline Models

For offline AI assistance with locally-run models:

```bash
# Install Mods
brew install charmbracelet/tap/mods

# Configure your preferred local model
mods --settings
```

**Example models in config:**
```toml
[ai.mentor]
provider = "mods"
model_complex_offline = "qwen3:30b"    # Larger model for deep analysis
model_simple_offline = "qwen3:4b"      # Smaller model for quick help

# Other local model options (if installed via Ollama):
# - "llama3.1:70b"
# - "codellama:34b"
# - "mistral:7b"
# - "deepseek-coder:33b"
```

#### Claude (Anthropic API) - Fallback

For Claude Code integration and complex tasks:

```bash
# Store API key securely in Skate
skate set codequest.anthropic_api_key "YOUR_API_KEY_HERE"
```

Get your key at [console.anthropic.com](https://console.anthropic.com)

**Example configuration:**
```toml
[ai.mentor]
provider = "claude"
model_complex = "claude-sonnet-4-5-20250929"     # Latest Sonnet 4.5
model_simple = "claude-haiku-4-5-20251001"       # Latest Haiku 4.5 (new!)

# Alternative models:
# - "claude-3-5-sonnet-20241022" (Previous Sonnet 3.5)
# - "claude-3-5-haiku-20241022" (Previous Haiku 3.5)
```

**Note**: API keys are stored in Skate's encrypted storage, never in the config file.

## 🎮 Usage

**Note:** The interactive TUI is currently being integrated. The features below describe the planned user experience based on implemented internal packages.

### Navigation (Planned)

- **Arrow Keys** or **h/j/k/l**: Navigate screens
- **Enter**: Select/Confirm
- **Esc**: Go back
- **?**: Show help
- **q**: Quit (from Dashboard)

### Screens (Planned)

- **Dashboard** (`d`): Overview of character, quests, and stats
- **Quest Board** (`q`): Browse and manage quests
- **Character** (`c`): View detailed character stats
- **Mentor** (`m`): Chat with AI for coding help
- **Settings** (`s`): Adjust configuration

### Global Hotkeys (Planned)

- **Ctrl+T**: Pause/Resume session timer (works anywhere)
- **Ctrl+C**: Quit application
- **?**: Toggle help overlay

### Workflows

#### Starting a Coding Session

1. Launch CodeQuest: `codequest`
2. Session timer starts automatically
3. Navigate to Quest Board (`q`) and accept a quest
4. Start coding in your repository
5. Commits automatically award XP and update quest progress

#### Getting AI Help

1. Press `m` to open Mentor screen
2. Type your question and press Enter
3. AI responds using Crush → Mods → Claude fallback chain
4. Chat history persists between sessions

#### Tracking Progress

- Dashboard shows real-time stats
- Session timer displays in footer (Ctrl+T to pause/resume)
- Quest progress updates on every commit
- Level-up notifications appear automatically
- Daily streak tracking encourages consistency

## 🛠️ Development

### Building from Source

```bash
# Download dependencies
make deps

# Build the application
make build

# Run tests
make test

# Run with coverage
make coverage

# Format and lint
make fmt vet lint
```

### Running Tests

```bash
# All tests
make test

# With coverage report
make coverage-html

# Integration tests only
go test ./test/integration -v

# Short tests (quick check)
make test-short
```

### Project Structure

```
codequest/
├── cmd/codequest/          # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── game/               # Game logic (character, quests, XP)
│   ├── storage/            # Data persistence (Skate)
│   ├── watcher/            # Git & session tracking
│   ├── ai/                 # AI provider integrations
│   └── ui/                 # Bubble Tea TUI
│       ├── screens/        # Dashboard, Quest Board, Mentor, etc.
│       └── components/     # Header, StatBar, Modal, Timer
└── test/integration/       # End-to-end tests
```

### Available Make Targets

```bash
make help           # Show all available commands
make build          # Build binary
make test           # Run all tests
make coverage       # Generate coverage report
make run            # Build and run
make install        # Install globally
make clean          # Clean build artifacts
make fmt            # Format code
make vet            # Run go vet
make lint           # Run linter
```

## 🎯 How It Works

### XP System

- **Commits**: 10-60 XP (base + lines bonus, capped)
- **Difficulty**: Easy +20%, Normal 1.0x, Hard -20%
- **Wisdom Bonus**: 1% per point above 10
- **Level Progression**: Polynomial curve (L1→2: 110 XP, L10→11: 2000 XP)

### Quest Types

- **Commit Quest**: Make N commits
- **Lines Quest**: Add/modify N lines of code
- **More types**: Tests, PR, refactoring (post-MVP)

### Character Stats

- **CodePower**: Increases commit quality bonus
- **Wisdom**: Increases XP gain
- **Agility**: Faster quest completion bonuses

All stats increase by +1 per level-up.

## 🔧 Troubleshooting

### "skate: command not found"

Install Skate for data persistence:
```bash
brew install charmbracelet/tap/skate
```

### "Failed to load configuration"

Create the config directory manually:
```bash
mkdir -p ~/.config/codequest
```

The application will generate a default config on next run.

### "No AI providers available"

Ensure at least one provider is configured:

- **Crush**: Set API key with `skate set codequest.openrouter_api_key "YOUR_KEY"`
- **Mods**: Install with `brew install charmbracelet/tap/mods`

You can also use CodeQuest without AI features - they're optional!

### Git commits not detected

Verify repository path in config:
```toml
[tracking]
repository_paths = ["~/your/project/path"]
```

Ensure you're working in a Git repository:
```bash
cd ~/your/project
git status  # Should show repository info
```

### Session timer not updating

Press **Ctrl+T** to start the timer if it's paused.

### Build fails with "missing go.mod"

Initialize the Go module:
```bash
make init
make deps
```

## 🧑‍💻 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new features
5. Run `make check` to ensure all tests pass
6. Follow [commit conventions](GIT_COMMIT_STYLE_GUIDE.md)
7. Submit a pull request

### Development Guidelines

- Follow [CLAUDE.md](CLAUDE.md) for Go standards and patterns
- Maintain >80% test coverage on core packages
- Document all exported functions and types
- Use table-driven tests
- Handle errors explicitly

## 📊 Tech Stack

- **Go 1.21+**: Core language
- **Bubble Tea**: TUI framework
- **Lip Gloss**: Styling and layout
- **Bubbles**: TUI components
- **go-git**: Git repository interactions
- **fsnotify**: File system watching
- **Skate**: Key-value storage
- **Cobra**: CLI framework

## 📚 Documentation

- [CODEQUEST_SPEC.md](CODEQUEST_SPEC.md) - Complete technical specification
- [CLAUDE.md](CLAUDE.md) - AI development guide and Go standards
- [GIT_COMMIT_STYLE_GUIDE.md](GIT_COMMIT_STYLE_GUIDE.md) - Commit message conventions
- [DEVELOPMENT_STATUS.md](DEVELOPMENT_STATUS.md) - Current development progress

## 🗺️ Roadmap

### MVP (v0.1.0) - In Progress 🚧

**Core Systems (Implemented):**
- ✅ Character system with XP/leveling (`internal/game/character.go`)
- ✅ Quest system (commits, lines) (`internal/game/quest.go`)
- ✅ TUI components with Bubble Tea (`internal/ui/`)
- ✅ Git activity monitoring (`internal/watcher/git.go`)
- ✅ Data persistence with Skate (`internal/storage/skate.go`)
- ✅ AI provider interfaces (`internal/ai/`)
- ✅ Session timer tracking (`internal/watcher/session.go`)
- ✅ Comprehensive test suite (>80% coverage on core packages)

**Integration (In Progress):**
- 🚧 Wire main.go to launch Bubble Tea application
- 🚧 Connect all components into working application
- 🚧 End-to-end testing of complete workflow

### Post-MVP Features

- [ ] Advanced quest types (tests, PR, refactoring)
- [ ] Skill tree system
- [ ] Achievement system
- [ ] GitHub API integration
- [ ] WakaTime integration
- [ ] Enhanced UI with animations
- [ ] Code review with bonus XP
- [ ] Quest generation with AI
- [ ] Multiplayer/guild features
- [ ] Web dashboard

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details.

## 🙏 Credits

Built with ❤️ using the [Charmbracelet](https://charm.sh) ecosystem.

Special thanks to:
- **Charmbracelet team** for amazing TUI tools
- **OpenRouter** for AI API access
- **Go community** for excellent libraries
- **Claude Code** for AI-assisted development

## 💬 Support

- **Issues**: [GitHub Issues](https://github.com/AutumnsGrove/codequest/issues)
- **Discussions**: [GitHub Discussions](https://github.com/AutumnsGrove/codequest/discussions)
- **Documentation**: See [docs/](docs/) for additional guides

## 🎉 Getting Started

Ready to transform your coding into an adventure?

```bash
# Install CodeQuest
git clone https://github.com/AutumnsGrove/codequest.git
cd codequest
make build

# Launch your quest
./build/codequest

# Start coding and watch your character grow!
```

---

**Ready to level up your coding?** 🚀

Start your quest: `codequest`
