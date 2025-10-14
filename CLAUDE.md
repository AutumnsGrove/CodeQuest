# CLAUDE.md - CodeQuest AI Development Guide

This document provides comprehensive guidance for Claude Code (or any AI assistant) when developing CodeQuest, a terminal-based gamified developer productivity RPG built with Go and the Charmbracelet ecosystem.

## Project Context

**CodeQuest** transforms coding work into an RPG adventure where:
- Every commit earns XP
- Bug fixes become quests
- Real development progress drives character growth
- Beautiful TUI powered by Bubble Tea, Lip Gloss, and Bubbles
- AI-powered mentorship through Crush, Mods, and Claude Code

**Primary Goals:**
1. Educational - Learn Go through building a real application
2. Practical - Create a genuinely helpful developer tool
3. Beautiful - Showcase Charmbracelet's capabilities
4. Fun - Make coding feel like an adventure

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

    "github.com/yourusername/codequest/internal/storage"
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

## AI Integration Guidelines

### Provider Hierarchy
1. **Crush (Primary)** - In-game mentor, quick help
2. **Mods (Local)** - Code review, offline fallback
3. **Claude Code** - Complex tasks, quest generation

### Implementation Pattern
```go
// Define provider interface
type AIProvider interface {
    Ask(question string, complexity string) (string, error)
    IsAvailable() bool
}

// Implement fallback chain
func (ai *AIManager) Query(prompt string) (string, error) {
    for _, provider := range ai.providers {
        if provider.IsAvailable() {
            response, err := provider.Ask(prompt)
            if err == nil {
                return response, nil
            }
        }
    }
    return "", fmt.Errorf("no AI providers available")
}
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

## Development Workflow

### Phase-Based Implementation

#### MVP (Weeks 1-3)
Focus on core functionality:
1. Character system with XP/levels
2. Basic quest types (commits, lines)
3. Simple TUI with dashboard
4. Git watcher for commits
5. Data persistence with Skate
6. Basic AI mentor

#### Enhancement Phases (Weeks 4-12)
Add features incrementally:
- Phase 2: Enhanced quests
- Phase 3: Skills & achievements
- Phase 4: Advanced UI
- Phase 5: AI integration
- Phase 6: External integrations
- Phase 7: Polish & documentation

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

## Security Best Practices

### API Key Management
- NEVER commit API keys
- Store in Skate's encrypted storage
- Use environment variables for CI/CD
- Implement key rotation

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

## Resource Links

### Essential Documentation
- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Lip Gloss Guide](https://github.com/charmbracelet/lipgloss)
- [Go Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Project Resources
- Specification: `CODEQUEST_SPEC.md`
- Commit Style: `GIT_COMMIT_STYLE_GUIDE.md`
- Architecture: `docs/ARCHITECTURE.md` (to be created)
- Roadmap: `docs/ROADMAP.md` (to be created)

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

---

*Remember: CodeQuest is not just a project, it's a learning adventure. Every line of code should teach something about Go, every feature should delight the user, and every commit should follow our standards.*