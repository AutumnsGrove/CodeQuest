# CodeQuest MVP Implementation Plan (Subagent-Driven)

## 🎯 Strategy: Focused Subagent Architecture

NO RESEARCH NEEDED - Spec is comprehensive (CODEQUEST_SPEC.md, CLAUDE.md)

APPROACH: Use highly focused subagents - ONE component per subagent, clean handoffs, sequential building.

---

## 📅 Week 1: Foundation Layer (Days 1-7)

### Day 1: Project Infrastructure

**Subagent 1: Dependency Setup**
- Task: Install all Go dependencies, verify builds
- Output: Working go.mod, go.sum, successful make build
- Handoff: Dependency list + build verification

**Subagent 2: Configuration System**
- File: internal/config/config.go, defaults.go, validate.go
- Task: Config struct, TOML loading, defaults, validation
- Context: Spec section on Configuration System
- Output: Complete config package with tests

### Days 2-3: Character System

**Subagent 3: Character Model Core**
- File: internal/game/character.go
- Task: Complete Character struct, NewCharacter(), AddXP(), level-up logic
- Context: Data Models section from spec
- Output: Full character model implementation

**Subagent 4: XP & Leveling Engine**
- File: internal/game/engine.go
- Task: CalculateXPForLevel(), XP curves, stat calculations, streak logic
- Context: Character model interface from Subagent 3
- Output: Game engine functions

**Subagent 5: Character Tests**
- File: internal/game/character_test.go
- Task: Comprehensive tests - creation, XP, level-ups, edge cases
- Context: Character implementation from Subagent 3+4
- Output: >80% coverage test suite

### Days 3-4: Quest System

**Subagent 6: Quest Model**
- File: internal/game/quest.go
- Task: Quest struct, quest types, lifecycle methods
- Context: Quest model from spec
- Output: Complete quest model

**Subagent 7: Quest Logic**
- File: internal/game/quest.go (methods)
- Task: IsAvailable(), Start(), UpdateProgress(), CheckCompletion()
- Context: Quest model from Subagent 6
- Output: All quest methods implemented

**Subagent 8: Quest Tests**
- File: internal/game/quest_test.go
- Task: Test all quest types, lifecycle, progress tracking
- Context: Quest implementation from Subagent 6+7
- Output: Quest test suite

### Days 5-6: Storage & Events

**Subagent 9: Skate Storage Wrapper**
- File: internal/storage/skate.go
- Task: Skate CLI wrapper, save/load functions, error handling
- Context: Storage requirements from spec
- Output: Complete storage layer

**Subagent 10: Storage Tests**
- File: internal/storage/skate_test.go
- Task: Test save/load, mock Skate, error cases
- Context: Storage implementation from Subagent 9
- Output: Storage test suite

**Subagent 11: Event System**
- File: internal/game/events.go
- Task: EventBus, pub/sub pattern, event types
- Context: Event Model from spec
- Output: Working event system with tests

---

## 📅 Week 2: UI Layer (Days 8-14)

### Day 8: UI Foundation

**Subagent 12: Styles System**
- File: internal/ui/styles.go
- Task: Lip Gloss styles, color scheme, reusable components
- Context: UI requirements from spec
- Output: Complete style definitions

**Subagent 13: Key Bindings**
- File: internal/ui/keys.go
- Task: Key binding definitions for all screens
- Context: Keybind Conventions from spec
- Output: Key binding system

**Subagent 14: Main App Model**
- File: internal/ui/app.go
- Task: Bubble Tea model, Init/Update/View, screen management
- Context: Bubble Tea patterns from spec
- Output: Main TUI application shell

### Days 9-10: Dashboard Screen

**Subagent 15: Dashboard View**
- File: internal/ui/screens/dashboard.go
- Task: Dashboard screen with character stats, XP bar, active quest
- Context: MVP UI mockups from spec, styles from Subagent 12
- Output: Complete dashboard screen

**Subagent 16: UI Components - Header**
- File: internal/ui/components/header.go
- Task: Reusable header component
- Context: Component requirements
- Output: Header component

**Subagent 17: UI Components - Stat Bar**
- File: internal/ui/components/statbar.go
- Task: Progress bars, stat displays, formatted numbers
- Context: Dashboard requirements
- Output: Stat bar components

### Days 11-12: Quest Board Screen

**Subagent 18: Quest Board View**
- File: internal/ui/screens/questboard.go
- Task: Quest list, selection, accept/decline actions
- Context: Quest system from Week 1, UI styles
- Output: Complete quest board screen

**Subagent 19: UI Components - Modal**
- File: internal/ui/components/modal.go
- Task: Modal dialogs, confirmations, level-up animations
- Context: UI requirements
- Output: Modal component

### Days 13-14: Mentor & Timer

**Subagent 20: Mentor Screen**
- File: internal/ui/screens/mentor.go
- Task: Text input, response display, loading indicator
- Context: AI Integration requirements
- Output: Mentor interface screen

**Subagent 21: Timer Component**
- File: internal/ui/components/timer.go
- Task: Session timer display, toggle functionality
- Context: Session tracking requirements
- Output: Timer component with global hotkey

**Subagent 22: UI Integration & Polish**
- Files: All UI files
- Task: Wire all screens together, smooth navigation, final styling
- Context: All UI components from Subagents 12-21
- Output: Fully integrated TUI

---

## 📅 Week 3: Integration Layer (Days 15-21)

### Days 15-16: Git Watcher

**Subagent 23: Git Repository Watcher**
- File: internal/watcher/git.go
- Task: Watch .git, detect commits, parse commit data
- Context: Git integration requirements from spec
- Output: Git watcher implementation

**Subagent 24: Git Event Integration**
- Files: internal/watcher/git.go + event system
- Task: Connect git watcher to event bus, fire commit events
- Context: Event system from Subagent 11, git watcher from Subagent 23
- Output: Git events flowing to game engine

**Subagent 25: Git Watcher Tests**
- File: internal/watcher/git_test.go
- Task: Test commit detection, parsing, event firing
- Context: Git watcher implementation
- Output: Git watcher test suite

### Day 17: Game-Git Integration

**Subagent 26: Commit Handler**
- Files: Game engine + event handlers
- Task: Handle commit events, update quests, award XP, check level-ups
- Context: Character/quest systems, git events
- Output: Complete commit → XP → level-up flow

**Subagent 27: Real-time UI Updates**
- Files: UI components + event handlers
- Task: Update UI in real-time when commits happen
- Context: UI system, event system
- Output: Live UI updates on game events

### Days 18-19: AI Integration

**Subagent 28: AI Provider Interface**
- File: internal/ai/provider.go
- Task: AI provider interface definition
- Context: AI Integration Strategy from spec
- Output: Provider interface

**Subagent 29: Crush/Mods Client**
- File: internal/ai/crush.go
- Task: Execute Mods/Crush CLI, handle responses, timeouts
- Context: Crush Integration section from spec
- Output: Working Crush client

**Subagent 30: AI-Mentor Integration**
- Files: Mentor screen + AI client
- Task: Connect mentor UI to AI client, handle loading/errors
- Context: Mentor screen + Crush client
- Output: Working mentor feature

### Day 20: Session Tracking

**Subagent 31: Session Timer**
- File: internal/watcher/session.go
- Task: Session time tracking, start/stop, persistence
- Context: Session tracking requirements
- Output: Session timer implementation

**Subagent 32: Timer-UI Integration**
- Files: Timer component + session tracker
- Task: Wire timer to session tracking, global hotkey
- Context: Timer component + session tracker
- Output: Working timer with Ctrl+T

### Day 21: Testing & Integration

**Subagent 33: Integration Tests**
- File: test/integration/mvp_test.go
- Task: End-to-end integration tests for MVP user flows
- Context: All systems from Subagents 1-32
- Output: Integration test suite

**Subagent 34: Bug Fixing & Polish**
- Files: All
- Task: Run full test suite, fix bugs, polish rough edges
- Context: Test results from all previous subagents
- Output: Stable MVP build

**Subagent 35: Documentation**
- File: README.md
- Task: Update README with setup instructions, quickstart guide
- Context: Complete MVP implementation
- Output: User-ready documentation

---

## 🎯 MVP Success Criteria (Verified by Subagent 34)

- User can create character → Character persists
- Accept quest → Make commits → Quest updates → Earn XP → Level up
- Mentor responds → Timer works → UI smooth → Tests pass (>60%)

---

## 📦 Subagent Handoff Protocol

Each subagent produces:

### Completion: [Name]
- **Summary:** [What was built]
- **Files Created/Modified:** [List]
- **Key Decisions:** [Important choices]
- **Next Subagent Needs:** [Context for next]
- **Confidence:** [High/Medium/Low] + reason

---

## 🚀 Execution Strategy

1. One subagent at a time - Complete before next
2. Sequential dependencies - Build foundation first
3. Clean handoffs - Structured artifacts between agents
4. Context pruning - Only essential context (<4000 tokens)
5. Test after build - Development before testing
6. Commit frequently - Follow GIT_COMMIT_STYLE_GUIDE.md

**35 focused subagents × focused tasks = clean, testable MVP** 🎮⚔️
