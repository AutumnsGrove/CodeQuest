# 🎮 CodeQuest

> **Transform your coding into an RPG adventure!**

CodeQuest is a terminal-based gamified developer productivity tool that turns your daily programming work into an epic RPG adventure. Built with Go and the beautiful Charmbracelet ecosystem, every commit earns XP, every bug fix is a quest, and your real development progress drives your character's growth.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)
![Status](https://img.shields.io/badge/Status-Pre--Development-orange?style=flat-square)

## ✨ Features

### 🎮 RPG Mechanics
- **Character Progression** - Level up as you code
- **Quest System** - Turn tasks into adventures
- **Combat Mode** - Real-time coding session tracking
- **Skills & Achievements** - Unlock abilities and earn badges
- **Daily Challenges** - Keep your coding streak alive

### 🤖 AI-Powered Assistance
- **Crush Mentor** - In-game AI companion for help and guidance
- **Code Review** - Automatic code quality feedback with bonus XP
- **Quest Generation** - AI-created challenges based on your project

### 🎨 Beautiful Terminal UI
- Powered by Bubble Tea, Lip Gloss, and Bubbles
- Smooth animations and transitions
- Responsive, adaptive layouts
- Rich colors and styling

### 📊 Productivity Tracking
- Git activity monitoring
- Session time tracking
- WakaTime integration (optional)
- GitHub API integration
- Comprehensive stats and analytics

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Git
- Terminal with 256 color support
- Ollama (for local AI features)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/codequest.git
cd codequest

# Install dependencies
go mod download

# Build the application
make build

# Run CodeQuest
./codequest
```

### First Run

When you launch CodeQuest for the first time:

1. **Create Your Character** - Choose a name and start your journey
2. **Accept Your First Quest** - Complete a simple coding task
3. **Make Commits** - Watch your XP grow with each commit
4. **Level Up** - Unlock new features and abilities

## 🎯 How It Works

CodeQuest monitors your Git activity and rewards you for productive coding:

```
Make a Commit → Earn XP → Complete Quests → Level Up → Unlock Features
```

### Example Session

```bash
# Start CodeQuest
codequest

# In another terminal, work on your project
git add .
git commit -m "feat: Add new feature"

# Return to CodeQuest to see:
# 🎉 +50 XP earned!
# ⚔️ Quest Progress: 3/5 commits
# 📈 Level 2 (450/500 XP)
```

## 🛠️ Configuration

CodeQuest can be configured via `~/.config/codequest/config.toml`:

```toml
[character]
name = "CodeWarrior"

[game]
difficulty = "normal"  # easy, normal, hard
auto_start_quests = false

[ai]
mentor_provider = "crush"  # crush, mods, claude-code
auto_review = true

[tracking]
session_timer_enabled = true
wakatime_enabled = false
```

## 🤝 AI Providers

CodeQuest supports multiple AI providers for different features:

### Crush (Primary Mentor)
- Online: Kimi K2/GLM-4.5 via OpenRouter
- Offline: Qwen3 via Ollama
- Used for in-game help and guidance

### Mods (Code Review)
- Local Qwen3 models
- Automatic code review after commits
- Bonus XP for clean code

### Claude Code (Advanced)
- Quest generation
- Complex assistance
- Backup when local models unavailable

## 🎮 Keybindings

### Dashboard (Main Screen)
- `Q` - Quest Board
- `C` - Character Sheet
- `I` - Inventory
- `M` - Mentor Help
- `S` - Settings
- `T` - Toggle Timer
- `Esc` - Exit

### Other Screens
- `Alt+Q` - Return to Dashboard
- `Alt+M` - Quick Mentor Help
- `Alt+S` - Settings
- `Ctrl+T` - Global Timer Toggle

## 📚 Documentation

- [CODEQUEST_SPEC.md](CODEQUEST_SPEC.md) - Full technical specification
- [CLAUDE.md](CLAUDE.md) - AI development guide
- [GIT_COMMIT_STYLE_GUIDE.md](GIT_COMMIT_STYLE_GUIDE.md) - Commit conventions
- [docs/](docs/) - Additional documentation (coming soon)

## 🗺️ Roadmap

### MVP (Current Focus)
- [x] Project setup and structure
- [ ] Character system with XP/leveling
- [ ] Basic quest types (commits, lines)
- [ ] Simple TUI with Bubble Tea
- [ ] Git activity monitoring
- [ ] Data persistence with Skate
- [ ] Basic AI mentor integration

### Post-MVP Features
- [ ] Advanced quest types
- [ ] Skill tree system
- [ ] Achievement system
- [ ] Multiplayer/guild features
- [ ] Web dashboard
- [ ] Mobile companion app
- [ ] IDE integrations

## 🧑‍💻 Development

### Building from Source

```bash
# Run tests
make test

# Run with hot reload
make dev

# Build for production
make build

# Install globally
make install

# Clean build artifacts
make clean
```

### Project Structure

```
codequest/
├── cmd/codequest/      # Application entry point
├── internal/           # Private application code
│   ├── game/          # Core game logic
│   ├── ui/            # Bubble Tea TUI
│   ├── storage/       # Data persistence
│   ├── ai/            # AI integrations
│   └── config/        # Configuration
├── data/              # Static game data
├── docs/              # Documentation
└── test/              # Integration tests
```

### Contributing

We welcome contributions! Please:

1. Read the [contribution guidelines](CONTRIBUTING.md)
2. Follow our [commit style guide](GIT_COMMIT_STYLE_GUIDE.md)
3. Write tests for new features
4. Update documentation

## 📄 License

CodeQuest is open source software licensed under the [MIT License](LICENSE).

## 🙏 Acknowledgments

Built with amazing tools from:
- [Charmbracelet](https://charm.sh) - Beautiful TUI tools for Go
- [Anthropic](https://anthropic.com) - Claude AI assistance
- The Go community

## 💬 Community

- **Issues**: [GitHub Issues](https://github.com/yourusername/codequest/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/codequest/discussions)
- **Discord**: Coming Soon!

---

**Ready to begin your quest?** 🎮⚔️

Start coding, earn XP, and level up your developer skills!