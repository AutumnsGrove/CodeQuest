# CLAUDE.md - CodeQuest AI Development Guide

This document provides comprehensive guidance for Claude Code (or any AI assistant) when developing CodeQuest, a terminal-based gamified developer productivity RPG built with Go and the Charmbracelet ecosystem.

## Project Context

**CodeQuest** transforms coding work into an RPG adventure where:
- Every commit earns XP
- Bug fixes become quests
- Real development progress drives character growth
- Beautiful TUI powered by Bubble Tea, Lip Gloss, and Bubbles
- AI-powered mentorship through Crush, Mods, and Claude Code

**Current Status (October 2025):**
- ✅ All core internal packages implemented and tested (>80% coverage)
- ✅ Character system, quests, XP engine, storage, UI components, Git watcher, session tracking
- 🚧 Final integration in progress (wiring main.go to launch Bubble Tea app)
- 📍 See [DEVELOPMENT_STATUS.md](DEVELOPMENT_STATUS.md) for detailed progress

**Primary Goals:**
1. Educational - Learn Go through building a real application
2. Practical - Create a genuinely helpful developer tool
3. Beautiful - Showcase Charmbracelet's capabilities
4. Fun - Make coding feel like an adventure

<!-- BaseProject: Core Behavior -->
## Essential Instructions (Always Follow)

### Core Behavior
- Do what has been asked; nothing more, nothing less
- NEVER create files unless absolutely necessary for achieving your goal
- ALWAYS prefer editing existing files to creating new ones
- NEVER proactively create documentation files (*.md) or README files unless explicitly requested

### Naming Conventions
- **Directories**: Use CamelCase (e.g., `VideoProcessor`, `AudioTools`, `DataAnalysis`)
- **Date-based paths**: Use skewer-case with YYYY-MM-DD (e.g., `logs-2025-01-15`, `backup-2025-12-31`)
- **No spaces or underscores** in directory names (except date-based paths)

### Communication Style
- Be concise but thorough
- Explain reasoning for significant decisions
- Ask for clarification when requirements are ambiguous
- Proactively suggest improvements when appropriate

## Development Philosophy

### 1. Real Work First
- Game mechanics must map to actual productive activity
- XP comes from commits, not arbitrary actions
- Progress requires real code, not grinding

### 2. Learning Through Building
- Start simple with MVP, add complexity incrementally
- Each phase teaches new Go concepts
- Well-documented code with educational comments

### 3. Beautiful AND Functional
- Smooth animations and responsive layouts
- Rich colors and styling with Lip Gloss
- Delightful user experience

## Go Development Standards

### Code Organization
```go
// Package names match directory names
package game

// Group imports: standard library, third-party, local
import (
    "fmt"
    "time"

    "github.com/charmbracelet/bubbletea"

    "github.com/AutumnsGrove/codequest/internal/storage"
)

// Document all exported types and functions
// Character represents the player in the game world
type Character struct {
    // Use meaningful field names with json tags
    Name  string `json:"name"`
    Level int    `json:"level"`
}

// Methods use pointer receivers for mutation
func (c *Character) AddXP(amount int) bool {
    // Handle errors explicitly
    if amount < 0 {
        return false
    }
    // Implementation...
    return true
}
```

### Error Handling Pattern
```go
// Always wrap errors with context
result, err := DoSomething()
if err != nil {
    return fmt.Errorf("doing something: %w", err)
}

// Create custom error types for domain errors
type QuestError struct {
    QuestID string
    Reason  string
}

func (e QuestError) Error() string {
    return fmt.Sprintf("quest %s failed: %s", e.QuestID, e.Reason)
}
```

### Testing Requirements
- Every exported function needs tests
- Aim for >80% coverage on core packages
- Use table-driven tests
- Mock external dependencies

```go
func TestCharacter_AddXP(t *testing.T) {
    tests := []struct {
        name      string
        xp        int
        wantLevel int
    }{
        {"normal gain", 50, 1},
        {"level up", 100, 2},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Charmbracelet Ecosystem Usage

### Bubble Tea Patterns
```go
// Model represents application state
type Model struct {
    character *game.Character
    quests    []game.Quest
    screen    Screen
}

// Init performs initial setup
func (m Model) Init() tea.Cmd {
    return tea.Batch(
        loadCharacter,
        watchGitActivity,
    )
}

// Update handles events
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case commitDetectedMsg:
        return m.handleCommit(msg)
    }
    return m, nil
}

// View renders the UI
func (m Model) View() string {
    return m.renderCurrentScreen()
}
```

### Lip Gloss Styling
```go
// Define reusable styles
var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205"))

    boxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("63"))
)

// Apply styles consistently
func renderTitle(text string) string {
    return titleStyle.Render(text)
}
```

<!-- BaseProject: House Agents -->
## House Agents - Quick Reference

### When to Use House Agents Proactively

**house-research**: Automatically invoke when searching across 20+ files for:
- Finding patterns across the codebase
- Searching for TODO/FIXME comments
- Locating API endpoints or function definitions
- Documentation searches
- Complex codebase analysis

**Example**: "Find all the authentication functions" → Automatically use house-research to find authentication functions

### Pattern Recognition
Main Claude should invoke house-research when:
- User mentions searching across many files
- Task involves finding patterns in the codebase
- Need to locate specific code patterns or TODOs
- Documentation or API searches required

## AI Integration Guidelines

### Provider Hierarchy
CodeQuest uses a fallback chain of AI providers for robust mentorship:

1. **Crush (OpenRouter)** - Online models via OpenRouter API
   - Primary for in-game mentor and quick help
   - Models: `openrouter/kimi/k2-0925`, `openrouter/deepseek/glm-4.5-air`
   - Requires API key stored in Skate: `codequest.openrouter_api_key`

2. **Mods (Local LLM)** - Offline local models
   - Fallback for code review and offline assistance
   - Models: `qwen3:30b` (complex), `qwen3:4b` (simple)
   - Requires Mods CLI installed via Homebrew

3. **Claude (Anthropic API)** - Claude models for complex tasks
   - Final fallback for advanced queries and quest generation
   - Models: `claude-sonnet-4-5-20250929`, `claude-haiku-4-5-20251001`
   - Requires API key stored in Skate: `codequest.anthropic_api_key`

### Implementation Pattern
```go
// Provider interface implemented by all AI backends
type AIProvider interface {
    Query(request Request) (Response, error)
    IsAvailable() bool
    Name() string
}

// AIManager orchestrates provider fallback chain
type AIManager struct {
    providers []AIProvider
    config    *config.AIConfig
}

// Query attempts each provider in order until one succeeds
func (ai *AIManager) Query(request Request) (Response, error) {
    for _, provider := range ai.providers {
        if !provider.IsAvailable() {
            continue
        }

        response, err := provider.Query(request)
        if err == nil {
            return response, nil
        }
        // Log error and try next provider
    }
    return Response{}, ErrNoProvidersAvailable
}
```

### Configuration Structure
```toml
# Config stored at ~/.config/codequest/config.toml
[ai.mentor]
provider = "crush"  # or "mods", "claude"
model_complex = "openrouter/kimi/k2-0925"
model_simple = "openrouter/deepseek/glm-4.5-air"
model_complex_offline = "qwen3:30b"
model_simple_offline = "qwen3:4b"
temperature = 0.7

[ai.review]
provider = "mods"
model_primary = "qwen3:30b"
model_fallback = "qwen3:4b"
auto_review = true
bonus_xp_enabled = true
```

### Rate Limiting
```go
type RateLimiter struct {
    requests int
    window   time.Duration
    mu       sync.Mutex
}

func (r *RateLimiter) Allow() bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    // Implementation
}
```

<!-- BaseProject: Git Workflow -->
## Git Workflow

### MANDATORY: Follow GIT_COMMIT_STYLE_GUIDE.md

**Every commit MUST follow this format:**
```
<type>: <description>

[optional body]
[optional footer]
```

**Valid types:**
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation only
- `style`: Formatting, no code change
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance
- `perf`: Performance improvements
- `build`: Build system changes
- `ci`: CI/CD changes

### After Completing Major Changes

**You MUST:**
1. Check git status: `git status`
2. Review recent commits for style: `git log --oneline -5`
3. Stage changes: `git add .`
4. Commit with proper message format (see below)
5. Verify commit succeeded: `git status && git log --oneline -1`

### Commit Message Template
```
[Action] [Brief description of what was changed]

- [Specific change 1 with technical detail]
- [Specific change 2 with technical detail]
- [Specific change 3 with technical detail]

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

### Commit Examples
```bash
# Feature addition
git commit -m "feat: Add character leveling system"

# Bug fix
git commit -m "fix: Correct XP calculation overflow"

# Documentation
git commit -m "docs: Add API documentation for game package"

# Refactoring
git commit -m "refactor: Extract quest validation logic"
```

### Branch Strategy
- `main` - Stable releases only
- `develop` - Active development
- `feature/*` - New features
- `fix/*` - Bug fixes

### When to Commit
Commit changes immediately after:
- ✅ Completing a significant feature or bug fix
- ✅ Adding new functionality that works correctly
- ✅ Making configuration or structural improvements
- ✅ Implementing user-requested features
- ✅ Fixing critical errors or security issues

<!-- BaseProject: TODO Management -->
## TODO Management

### MANDATORY: Maintain TODOS.md

You MUST actively maintain the `TODOS.md` file in the project root. This is a critical part of the workflow.

**Always check TODOS.md first** when starting a new task or session

**Update TODOS.md immediately** when:
- A task is completed (mark with ✅ or remove)
- A new task is identified (add it)
- A task's priority or status changes
- You discover subtasks or dependencies

**Format for TODOS.md:**
```markdown
# Project TODOs

## High Priority
- [ ] Task description here
- [x] Completed task (keep for reference)

## Medium Priority
- [ ] Another task

## Low Priority / Future Ideas
- [ ] Nice to have feature

## Blocked
- [ ] Task blocked by X (waiting on...)
```

**Use clear task descriptions** that include:
- What needs to be done
- Why it's important (if not obvious)
- Any dependencies or blockers

**Keep it current**: Remove or archive completed tasks regularly to keep the list manageable

## Development Workflow

### Current Phase: MVP Integration (Week 4)

**Completed (Weeks 1-3):**
- ✅ Character system with XP/levels (`internal/game/character.go`)
- ✅ Quest system with lifecycle (`internal/game/quest.go`)
- ✅ XP engine with balanced progression (`internal/game/engine.go`)
- ✅ Event bus for pub/sub (`internal/game/events.go`)
- ✅ Skate storage wrapper (`internal/storage/skate.go`)
- ✅ Git watcher and integration (`internal/watcher/git.go`, `integration.go`)
- ✅ Session tracking (`internal/watcher/session.go`)
- ✅ UI components and screens (`internal/ui/`)
- ✅ AI provider interfaces (`internal/ai/`)
- ✅ Comprehensive test suite (>80% coverage on core packages)

**Current Focus:**
- 🚧 Wire `cmd/codequest/main.go` to initialize and launch Bubble Tea app
- 🚧 Connect all components: Config → Storage → EventBus → UI → Watchers
- 🚧 End-to-end integration testing
- 🚧 User documentation and setup guides

**Next: Post-MVP Enhancement Phases**
- Phase 2: Advanced quest types (tests, PR, refactoring)
- Phase 3: Skill tree and achievement systems
- Phase 4: Enhanced UI with animations
- Phase 5: AI-generated quests and code review
- Phase 6: External integrations (GitHub, WakaTime)
- Phase 7: Polish, optimization, and extended documentation

### Testing During Development
```bash
# Run tests continuously
go test ./... -v

# Check coverage
go test ./... -cover

# Run specific package tests
go test ./internal/game -v

# Benchmark performance
go test -bench=. ./internal/game
```

### Building and Running
```bash
# Build the application
make build

# Run in development mode
make dev

# Install globally
make install

# Clean build artifacts
make clean
```

## Application Architecture

### Current Implementation Status

**Implemented Packages:**
```
internal/
├── config/      ✅ Configuration management (TOML + Skate)
├── game/        ✅ Core game logic (character, quests, XP, events)
├── storage/     ✅ Data persistence (Skate wrapper)
├── watcher/     ✅ Git monitoring & session tracking
├── ai/          ✅ AI provider interfaces (Crush, Mods, Claude)
└── ui/          ✅ Bubble Tea TUI (screens + components)
```

**Integration Needed:**
```go
// cmd/codequest/main.go needs to:
// 1. Load configuration
cfg, err := config.Load()

// 2. Initialize storage client
storage, err := storage.NewSkateClient()

// 3. Create event bus
eventBus := game.NewEventBus()

// 4. Initialize UI model
model := ui.NewModel(storage, cfg)

// 5. Start Git watcher
watcher := watcher.NewWatcherManager(eventBus)
watcher.AddRepository(cfg.Git.WatchPaths...)

// 6. Register game event handlers
engine := game.NewGameEngine(character, quests, storage)
engine.RegisterHandlers(eventBus)

// 7. Launch Bubble Tea program
program := tea.NewProgram(model, tea.WithAltScreen())
program.Run()
```

### Key Integration Points

1. **Config → Everything**: All components need config for initialization
2. **Storage → Character/Quests**: Load/save game state
3. **EventBus → UI**: Real-time updates on commits, level-ups, quest completion
4. **GitWatcher → EventBus**: Publishes commit events
5. **GameEngine → EventBus**: Subscribes to events, updates character/quests
6. **UI → Storage**: Periodic saves and manual save commands

## Common Tasks

### Adding a New Feature
1. Create feature branch: `git checkout -b feature/feature-name`
2. Write tests first (TDD approach)
3. Implement feature
4. Update documentation
5. Run tests: `make test`
6. Commit with proper type: `git commit -m "feat: Add feature description"`

### Creating a Quest Type
```go
// 1. Define the quest type
const QuestTypeCustom QuestType = "custom"

// 2. Implement validation
func (q *Quest) validateCustom() error {
    // Validation logic
}

// 3. Implement progress tracking
func (q *Quest) updateCustomProgress(event Event) {
    // Progress logic
}

// 4. Add tests
func TestCustomQuest(t *testing.T) {
    // Test implementation
}
```

### Adding UI Screen
```go
// 1. Create screen component
type CustomScreen struct {
    // Screen state
}

// 2. Implement Update method
func (s CustomScreen) Update(msg tea.Msg) (CustomScreen, tea.Cmd) {
    // Handle events
}

// 3. Implement View method
func (s CustomScreen) View() string {
    // Render UI
}

// 4. Register in main app
func (m Model) switchToCustom() (Model, tea.Cmd) {
    m.screen = ScreenCustom
    return m, nil
}
```

## Performance Considerations

### Optimization Guidelines
- Profile before optimizing
- Use buffered channels for events
- Pool frequently allocated objects
- Cache expensive computations
- Minimize allocations in hot paths

### Profiling
```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

<!-- BaseProject: Security Best Practices -->
## Security Best Practices

### API Key Management

**CodeQuest-Specific:**
- Primary storage: Skate's encrypted storage (`codequest.openrouter_api_key`, `codequest.anthropic_api_key`)
- Access via: `skate get codequest.openrouter_api_key`
- NEVER commit API keys to version control
- Implement key rotation for compromised keys

**General Best Practices (for other projects):**
- Store API keys in `secrets.json` files (add to `.gitignore` immediately)
- Provide `secrets_template.json` with empty/example values for setup
- Use environment variables as fallbacks
- Show clear status messages about key loading source

**Example Loading Pattern (Python projects):**
```python
def load_secrets():
    """Load API keys from secrets.json file."""
    secrets_path = os.path.join(os.path.dirname(__file__), "secrets.json")
    try:
        with open(secrets_path, 'r') as f:
            secrets = json.load(f)
        return secrets
    except FileNotFoundError:
        print(f"Warning: secrets.json not found. Using environment variables as fallback.")
        return {}
    except json.JSONDecodeError as e:
        print(f"Error parsing secrets.json: {e}. Using environment variables as fallback.")
        return {}

# Load secrets at startup
SECRETS = load_secrets()
API_KEY = SECRETS.get("anthropic_api_key", os.getenv("ANTHROPIC_API_KEY", ""))
```

**DO:**
- Store all API keys securely (Skate for CodeQuest, secrets.json for others)
- Add secrets files to `.gitignore` immediately
- Provide template files with empty values for setup
- Use environment variables as fallbacks
- Include error handling for missing/malformed secrets

**DON'T:**
- Hardcode API keys directly in source code
- Commit actual API keys to version control
- Store keys in configuration files that might be shared
- Log or print actual API key values

### Input Validation
```go
func validateInput(input string) error {
    if len(input) > MaxInputLength {
        return ErrInputTooLong
    }
    if !isValidUTF8(input) {
        return ErrInvalidEncoding
    }
    return nil
}
```

## Documentation Standards

### Code Documentation
- Document all exported items
- Explain "why" not just "what"
- Include examples for complex functions
- Keep comments up-to-date

### User Documentation
- Maintain README.md with quickstart
- Create user guides in docs/
- Document configuration options
- Provide troubleshooting guides

## Debugging Tips

### Common Issues
1. **TUI not rendering**: Check terminal capabilities
2. **Git watcher not working**: Verify git repository
3. **AI timeout**: Check internet/API keys
4. **Data not persisting**: Verify Skate installation

### Debug Mode
```go
// Enable debug logging
if os.Getenv("CODEQUEST_DEBUG") == "1" {
    log.SetLevel(log.DebugLevel)
}

// Add debug output
log.Debug("Character state: %+v", character)
```

## Quality Checklist

Before committing:
- [ ] Tests pass: `make test`
- [ ] Code formatted: `go fmt ./...`
- [ ] Linter clean: `golangci-lint run`
- [ ] Documentation updated
- [ ] Commit message follows style guide
- [ ] No API keys or secrets

## Technical Requirements

### Go Version
- **Development:** Go 1.25.1 (current as of October 2025)
- **Minimum:** Go 1.21+ for users
- **Module:** `github.com/AutumnsGrove/codequest`

### Key Dependencies
```go
// Core Charmbracelet TUI stack
github.com/charmbracelet/bubbletea  v1.3.10   // TUI framework
github.com/charmbracelet/lipgloss   v1.1.0    // Styling
github.com/charmbracelet/bubbles    v0.21.0   // Components

// Configuration and storage
github.com/BurntSushi/toml          v1.5.0    // Config parsing
// Skate CLI (installed separately via Homebrew)

// Git operations and file watching
github.com/go-git/go-git/v5         v5.16.3   // Git operations
github.com/fsnotify/fsnotify        v1.9.0    // File watching

// Utilities
github.com/google/uuid              v1.6.0    // UUID generation
```

### External Tools Required
- **Skate** (data persistence): `brew install charmbracelet/tap/skate`
- **Mods** (optional, local AI): `brew install charmbracelet/tap/mods`

## Resource Links

### Essential Documentation
- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Lip Gloss Guide](https://github.com/charmbracelet/lipgloss)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Go Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Project Documentation
- **README.md** - User-facing documentation and setup guide
- **CODEQUEST_SPEC.md** - Complete technical specification and design
- **GIT_COMMIT_STYLE_GUIDE.md** - Mandatory commit message format
- **DEVELOPMENT_STATUS.md** - Current progress and subagent completion status
- **CLAUDE.md** - This file - AI development guide

## AI Assistant Reminders

When working on CodeQuest:
1. **Always follow the commit style guide** - Every single commit
2. **Think educationally** - Add comments explaining Go concepts
3. **Test everything** - TDD approach preferred
4. **Build incrementally** - MVP first, features later
5. **Document thoroughly** - Future developers will thank you
6. **Consider the user** - Make it fun and engaging
7. **Respect the spec** - But suggest improvements
8. **Handle errors gracefully** - Never panic in production
9. **Optimize wisely** - Profile first, optimize second
10. **Enjoy the journey** - This is about learning!

## Next Development Phase: Main Application Wiring

**Priority Task:** Implement `cmd/codequest/main.go` to launch the full application.

**Key Steps:**
1. Initialize configuration (create default if missing)
2. Set up error handling and graceful shutdown
3. Initialize Skate storage client (handle missing Skate gracefully)
4. Load or create character
5. Create event bus and register handlers
6. Start Git watcher with configured paths
7. Initialize session tracker
8. Create and start Bubble Tea UI
9. Handle cleanup on exit (save state, stop watchers)

**Important Considerations:**
- Graceful degradation if Skate not installed (show helpful error)
- Handle missing config file (auto-create with defaults)
- Proper cleanup on Ctrl+C (save character, stop watchers)
- Connect EventBus to both UI and game logic
- Ensure all goroutines are properly managed and stopped

**Testing Strategy:**
- Manual testing of full application flow
- Integration tests for component initialization
- Test graceful shutdown behavior
- Verify data persistence across restarts

---

*Remember: CodeQuest is not just a project, it's a learning adventure. Every line of code should teach something about Go, every feature should delight the user, and every commit should follow our standards.*

*Current Status: All internal packages complete. Final integration to make it all come together!* 🚀