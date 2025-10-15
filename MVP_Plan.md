CodeQuest MVP Implementation Plan

🎯 Goal

Build a minimum playable version in 3 weeks that demonstrates core RPG mechanics tied to real coding work.

📅 3-Week Development Timeline

Week 1: Foundation (Days 1-7) 🏗

Phase 1.1: Project Setup & Dependencies (Day 1)

1. Install all Go dependencies via go.mod:
  - Bubble Tea (TUI framework)
  - Lip Gloss (styling)
  - Bubbles (components)
  - Cobra (CLI commands)
  - go-git (Git operations)
  - fsnotify (file watching)
  - uuid (ID generation)
2. Run go mod download and verify all deps
3. Test basic build with make build

Phase 1.2: Character System (Days 2-3)

1. Complete internal/game/character.go:
  - Implement AddXP() with level-up logic
  - Add XP curve calculation (CalculateXPForLevel())
  - Implement streak tracking (UpdateStreak())
  - Add all fields from spec (RPG stats, progress tracking)
2. Create internal/game/character_test.go:
  - Test character creation
  - Test XP addition and level-ups
  - Test XP curve (level 1→2, 2→3, multi-level jumps)
  - Test streak calculations
  - Aim for >80% coverage

Phase 1.3: Quest System (Days 3-4)

1. Create internal/game/quest.go:
  - Define Quest struct with all fields
  - Implement quest types (Commit, Lines)
  - Add IsAvailable(), Start(), UpdateProgress(), CheckCompletion()
2. Create internal/game/quest_test.go:
  - Test quest lifecycle (available → active → completed)
  - Test progress tracking
  - Test completion detection
3. Create internal/game/engine.go:
  - XP calculation and reward logic
  - Quest management functions
  - Level-up handling

Phase 1.4: Storage Layer (Days 5-6)

1. Create internal/storage/skate.go:
  - Wrapper functions for Skate CLI
  - SaveCharacter(), LoadCharacter()
  - SaveQuest(), LoadQuest(), LoadAllQuests()
  - Error handling for missing data
2. Create internal/storage/skate_test.go:
  - Test save/load operations
  - Test data serialization
  - Mock Skate for testing
3. Create internal/config/config.go:
  - Config struct matching spec
  - Load from ~/.config/codequest/config.toml
  - Default config generation
  - Validation

Phase 1.5: Event System (Day 7)

1. Create internal/game/events.go:
  - Event bus for pub/sub pattern
  - Event types (commit, level_up, quest_start, quest_done)
  - Subscribe/Publish methods
2. Test event system

---
Week 2: User Interface (Days 8-14) 🎨

Phase 2.1: Bubble Tea Setup (Day 8)

1. Complete internal/ui/app.go:
  - Main Bubble Tea model
  - Screen state management
  - Init(), Update(), View() methods
  - Keyboard input routing
2. Create internal/ui/styles.go:
  - Lip Gloss style definitions
  - Color scheme (based on spec)
  - Reusable box/border styles
  - Text formatting helpers
3. Create internal/ui/keys.go:
  - Key binding definitions
  - Dashboard keys (single letter)
  - Other screen keys (with modifiers)

Phase 2.2: Dashboard Screen (Days 9-10)

1. Create internal/ui/screens/dashboard.go:
  - Welcome message with character name
  - Current level and XP display
  - XP progress bar
  - Active quest summary
  - Today's stats (commits, lines, time)
  - Key binding hints at bottom
2. Implement navigation:
  - Q → Quest Board
  - C → Character Sheet (placeholder)
  - M → Mentor dialog
  - T → Toggle timer
  - Esc → Exit

Phase 2.3: Quest Board Screen (Days 11-12)

1. Create internal/ui/screens/questboard.go:
  - List of available quests
  - Quest details (title, description, requirements, rewards)
  - Highlight selected quest
  - Accept/decline actions
  - Show active quest status
2. Create internal/ui/components/statbar.go:
  - Reusable stat display component
  - Progress bars
  - Formatted numbers

Phase 2.4: UI Components & Polish (Days 13-14)

1. Create internal/ui/components/header.go:
  - Screen title component
  - Character info header
2. Create internal/ui/components/modal.go:
  - Modal dialog for messages
  - Confirmation dialogs
  - Level-up animations (simple)
3. Create internal/ui/screens/mentor.go:
  - Simple text input for questions
  - Response display
  - Loading indicator
4. Polish navigation and transitions
5. Add color and styling throughout

---
Week 3: Integration (Days 15-21) 🔌

Phase 3.1: Git Watcher (Days 15-16)

1. Create internal/watcher/git.go:
  - Watch .git directory for changes
  - Detect new commits via go-git
  - Parse commit data (files changed, lines added/removed)
  - Fire events on commit detection
2. Create internal/watcher/git_test.go:
  - Test commit detection
  - Test data parsing
  - Integration test with real git repo

Phase 3.2: Git Integration with Game (Day 17)

1. Connect Git watcher to event bus
2. Handle commit events in quest system:
  - Update quest progress
  - Award XP based on commit size
  - Check quest completion
  - Trigger level-up if threshold reached
3. Update UI to reflect changes in real-time

Phase 3.3: AI Integration - Crush/Mods (Days 18-19)

1. Create internal/ai/provider.go:
  - AI provider interface
  - Model selection logic
2. Create internal/ai/crush.go:
  - Execute Mods CLI with Crush model
  - Simple question/answer flow
  - Basic error handling
  - Timeout handling
3. Create internal/ai/mods.go:
  - Fallback to Mods directly
  - Local model support
4. Integrate mentor into UI:
  - /mentor <question> command
  - Display response in modal
  - Loading state

Phase 3.4: Session Tracking (Day 20)

1. Create internal/watcher/session.go:
  - Session timer tracking
  - Start/stop functionality
  - Save session time to character
2. Create internal/ui/components/timer.go:
  - Timer display component
  - Global hotkey (Ctrl+T) to toggle
  - Show/hide state
3. Integrate timer into all screens

Phase 3.5: Testing & Bug Fixes (Day 21)

1. Run full test suite:
  - All unit tests passing
  - Integration tests passing
  - Coverage >60% (MVP target)
2. Manual testing:
  - Create character flow
  - Accept and complete quest
  - Git commits award XP
  - Level-up works
  - Mentor responds
  - Data persists between sessions
3. Fix critical bugs
4. Update README with setup instructions

---
🎯 MVP Success Criteria

By end of Week 3, we must have:

- User can create a character
- Character saves/loads from Skate
- User can accept a quest from quest board
- Git commits are detected automatically
- Quest progress updates on commits
- Character earns XP from commits
- Character levels up at XP threshold
- Level-up shows visual feedback
- User can ask Crush/Mods for help via /mentor
- Session timer works with Ctrl+T toggle
- Dashboard shows all key stats
- Navigation works smoothly between screens
- All core models have tests (>60% coverage)
- TUI doesn't crash on edge cases
- README has clear setup instructions
- Config file generates with defaults

---
📦 Deliverables

1. Binary: codequest (single executable)
2. Config: ~/.config/codequest/config.toml (auto-generated)
3. Data: Stored in Skate KV store
4. Tests: Unit + integration tests
5. Docs: Updated README with quickstart

---
🚀 Post-MVP Phases

After MVP is complete and working, add features incrementally:

- Phase 2 (Week 4-5): Enhanced quests (PR-based, test coverage, chains)
- Phase 3 (Week 6-7): Skills & achievements
- Phase 4 (Week 8-9): Advanced UI (character sheet, inventory, animations)
- Phase 5 (Week 10): Advanced AI (code review, quest generation)
- Phase 6 (Week 11): External integrations (GitHub, WakaTime)
- Phase 7 (Week 12): Polish & documentation

---
🎓 Learning Path

This MVP will teach:
- Week 1: Go basics (structs, methods, interfaces, JSON, testing)
- Week 2: Bubble Tea patterns (Model/Update/View, components)
- Week 3: Go concurrency (goroutines, channels, file watching)

---
⚡ Implementation Strategy

1. Follow TDD - Write tests first when possible
2. Commit frequently - Following GIT_COMMIT_STYLE_GUIDE.md
3. Build incrementally - Test each component before moving on
4. Document as you go - Add educational comments
5. Use the spec - CODEQUEST_SPEC.md is the source of truth

Ready to start building! 🎮⚔️
