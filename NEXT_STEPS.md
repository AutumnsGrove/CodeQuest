# CodeQuest: Next Steps & Development Roadmap

**Last Updated:** October 15, 2025
**Current Status:** MVP Complete (v0.1.0 Production Ready)
**Target:** v1.0.0 Stable Release (~10 weeks)

---

## 📊 Current Status

### ✅ MVP Complete (All 35 Subagents Done)

**Core Systems Implemented:**
- ✅ Character system with XP/leveling (`internal/game/character.go`)
- ✅ Quest system with commit & lines quests (`internal/game/quest.go`)
- ✅ XP engine with balanced progression (`internal/game/engine.go`)
- ✅ Event bus for pub/sub architecture (`internal/game/events.go`)
- ✅ Skate storage wrapper (`internal/storage/skate.go`)
- ✅ Git watcher and integration (`internal/watcher/git.go`, `integration.go`)
- ✅ Session tracking with Ctrl+T hotkey (`internal/watcher/session.go`)
- ✅ Complete UI with Bubble Tea (`internal/ui/`)
  - Dashboard, Quest Board, Character, Mentor, Settings screens
  - Header, StatBar, Modal, Timer components
- ✅ AI provider integration (`internal/ai/`)
  - Crush (OpenRouter), Mods (local), Claude (Anthropic) fallback chain
- ✅ Comprehensive test suite (>80% coverage on core packages)

**Test Coverage:**
- internal/config: 84.3%
- internal/game: 44.9%
- internal/storage: 80.3%
- internal/ui/components: 94.5%
- internal/ui/screens: 70.0%
- internal/watcher: 50.9%

**Build Status:** ✅ All tests passing, no race conditions, production-ready

---

## 🎯 Development Roadmap Overview

| Phase | Duration | Subagents | Key Features | Release |
|-------|----------|-----------|--------------|---------|
| **Phase 0** | 1-2 days | 3 | Release Validation & v0.1.0 Tag | v0.1.0 |
| **Phase 1** | 1 week | 3 | User Testing & Bug Fixes | v0.1.1 |
| **Phase 2** | 2 weeks | 8 | Enhanced Quest Types | v0.2.0 |
| **Phase 3** | 2 weeks | 7 | Skills & Achievements | v0.3.0 |
| **Phase 4** | 1.5 weeks | 6 | Advanced UI & Animations | v0.4.0 |
| **Phase 5** | 1 week | 5 | Advanced AI Features | v0.5.0 |
| **Phase 6** | 1 week | 4 | External Integrations | v0.6.0 |
| **Phase 7** | 1 week | 6 | Polish & v1.0.0 | v1.0.0 |
| **TOTAL** | **~10 weeks** | **42 subagents** | Complete Feature Set | v1.0.0 |

---

## Phase 0: Release Preparation & Validation (1-2 days)

**Goal:** Ensure MVP is truly production-ready and release v0.1.0

### Subagent 0.1: End-to-End Integration Validation
**Duration:** 4-6 hours
**Priority:** CRITICAL

**Description:** Verify main.go properly wires all components and the application launches successfully.

**Tasks:**
- [ ] Build binary: `make build`
- [ ] Test application launch: `./build/codequest`
- [ ] Verify character creation flow on first run
- [ ] Test quest acceptance from Quest Board
- [ ] Make test commits in a watched repository
- [ ] Verify XP gain and quest progress updates
- [ ] Test session timer with Ctrl+T hotkey
- [ ] Confirm data persistence (quit and relaunch)
- [ ] Test AI mentor commands (if API keys configured)
- [ ] Check all screen navigation (Dashboard, Quest Board, Character, Mentor, Settings)
- [ ] Test error handling (missing Skate, no git repo, network issues)

**Deliverable:**
- Integration test report with status of each feature
- Bug list if any issues found
- Confirmation that app is launch-ready OR list of blockers

**Dependencies:** None (starting point)

---

### Subagent 0.2: Pre-Release Checklist & Cleanup
**Duration:** 2-3 hours
**Priority:** HIGH

**Description:** Final polishing and preparation before tagging v0.1.0.

**Tasks:**
- [ ] Commit pending changes:
  - `.claude/settings.json`
  - `CLAUDE.md`
  - `README.md`
- [ ] Update version numbers in code (if any)
- [ ] Verify README.md accuracy:
  - Installation instructions current
  - Features list matches implementation
  - Status badges correct
- [ ] Consistency check across all docs:
  - README.md
  - DEVELOPMENT_STATUS.md
  - CODEQUEST_SPEC.md
  - CLAUDE.md
- [ ] Run full test suite: `make test`
- [ ] Run with race detector: `go test -race ./...`
- [ ] Format code: `go fmt ./...`
- [ ] Tidy dependencies: `go mod tidy`
- [ ] Create CHANGELOG.md for v0.1.0:
  - List all MVP features
  - Credit contributors
  - Known limitations

**Deliverable:**
- Clean git tree ready for tagging
- CHANGELOG.md with v0.1.0 notes
- All tests passing
- Documentation synchronized

**Dependencies:** Subagent 0.1 (must pass validation first)

---

### Subagent 0.3: v0.1.0 Release Execution
**Duration:** 2-3 hours
**Priority:** HIGH

**Description:** Tag the release, build binaries, and publish to GitHub.

**Tasks:**
- [ ] Create annotated git tag: `git tag -a v0.1.0 -m "MVP Release: Core features complete"`
- [ ] Push tag: `git push origin v0.1.0`
- [ ] Build release binaries:
  - macOS (arm64, amd64)
  - Linux (amd64, arm64)
  - Windows (amd64)
- [ ] Create GitHub release:
  - Title: "v0.1.0 - MVP Release"
  - Copy CHANGELOG.md content
  - Attach binaries
  - Mark as "pre-release" if desired
- [ ] Update README.md badges/status
- [ ] Announce release (if applicable)

**Deliverable:**
- Published v0.1.0 release on GitHub with binaries
- Public release notes
- Updated README with installation from release

**Dependencies:** Subagent 0.2 (must have clean tree)

---

## Phase 1: User Testing & Iteration (1 week)

**Goal:** Deploy to real workflow, gather feedback, and patch bugs

### Subagent 1.1: Personal Deployment & Daily Use
**Duration:** 5-7 days (ongoing)
**Priority:** HIGH

**Description:** Use CodeQuest in actual development workflow to identify issues.

**Tasks:**
- [ ] Install v0.1.0 on primary development machine
- [ ] Configure with real projects:
  - Set up git watch paths
  - Configure AI providers (Crush/Mods/Claude)
  - Customize settings
- [ ] Use for at least 5 coding sessions:
  - Accept quests
  - Make commits
  - Track progress
  - Use AI mentor
- [ ] Test edge cases:
  - Very large commits (1000+ lines)
  - Network loss during AI queries
  - Multiple active quests
  - Long coding sessions (>6 hours)
  - Repository switching
- [ ] Document issues:
  - Bugs (crashes, incorrect behavior)
  - UX problems (confusing UI, unclear messages)
  - Performance issues (lag, high CPU)
  - Missing features or QoL improvements

**Deliverable:**
- User experience report
- Detailed bug list with reproduction steps
- Feature request list prioritized by impact

**Dependencies:** Subagent 0.3 (need v0.1.0 installed)

---

### Subagent 1.2: Documentation Improvements
**Duration:** 1-2 days
**Priority:** MEDIUM

**Description:** Based on usage feedback, enhance user-facing documentation.

**Tasks:**
- [ ] Expand README.md:
  - Add troubleshooting section
  - Include screenshots/GIFs
  - Add "Common Workflows" examples
- [ ] Create docs/ directory with:
  - `QUICKSTART.md` - 5-minute getting started
  - `CONFIGURATION.md` - Detailed config guide
  - `FAQ.md` - Frequently asked questions
  - `TROUBLESHOOTING.md` - Common issues
- [ ] Record demo GIFs with VHS:
  - Installation process
  - First run and character creation
  - Quest acceptance and completion
  - Using AI mentor
  - Session timer usage
- [ ] Update CLAUDE.md with any new patterns learned

**Deliverable:**
- Enhanced documentation suite in docs/
- Demo GIFs in assets/ or embedded in README
- Improved onboarding experience

**Dependencies:** Subagent 1.1 (need real usage experience)

---

### Subagent 1.3: Bug Fixes & Polish (v0.1.1)
**Duration:** 2-3 days
**Priority:** HIGH

**Description:** Address critical bugs and UX issues found during testing.

**Tasks:**
- [ ] Triage bug list from Subagent 1.1:
  - P0: Crashes, data loss (fix immediately)
  - P1: Major functionality broken (fix for v0.1.1)
  - P2: Minor bugs (consider for v0.1.1 or v0.2.0)
  - P3: Polish/nice-to-have (backlog)
- [ ] Fix P0 and P1 bugs:
  - Write regression tests first
  - Implement fixes
  - Verify with end-to-end testing
- [ ] UX improvements:
  - Better error messages
  - Loading indicators
  - Confirmation dialogs
  - Keyboard shortcut hints
- [ ] Performance optimizations if needed:
  - Profile hot paths
  - Reduce allocations
  - Optimize git watcher
- [ ] Update tests and docs
- [ ] Create CHANGELOG entry for v0.1.1
- [ ] Tag and release v0.1.1

**Deliverable:**
- v0.1.1 patch release
- Bug fixes validated
- Improved stability and UX

**Dependencies:** Subagent 1.1 (need bug list)

---

## Phase 2: Enhanced Quest System (2 weeks)

**Goal:** Implement advanced quest types from CODEQUEST_SPEC.md Phase 2

### Subagent 2.1: Test Coverage Quest Type
**Duration:** 2-3 days
**Priority:** HIGH

**Description:** Add quests for writing tests.

**Tasks:**
- [ ] Define model:
  - `QuestTypeTests` constant
  - Target: number of tests or coverage percentage
- [ ] Implement test detection:
  - Parse test files (`*_test.go`)
  - Count test functions (`func Test...`)
  - Optional: Calculate coverage with `go test -cover`
- [ ] Add quest creation:
  - "Write 10 Tests" quest
  - "Achieve 80% Coverage" quest
- [ ] Update progress tracking:
  - Hook into git watcher
  - Detect new test commits
  - Update quest progress
- [ ] UI updates:
  - Display test count in quest description
  - Show coverage percentage if applicable
- [ ] Tests:
  - Test quest creation
  - Test progress updates
  - Test completion detection

**Deliverable:**
- `internal/game/quest_tests.go` (or add to `quest.go`)
- Test detection logic in watcher
- Tests with >80% coverage

**Files to Modify:**
- `internal/game/quest.go` (add QuestTypeTests)
- `internal/watcher/git.go` (add test file detection)
- `internal/game/quest_test.go` (add test quest tests)

**Dependencies:** None (builds on existing quest system)

---

### Subagent 2.2: Pull Request Quest Type
**Duration:** 3-4 days
**Priority:** MEDIUM

**Description:** Add quests for creating and merging pull requests.

**Tasks:**
- [ ] Define model:
  - `QuestTypePR` constant
  - Criteria: create PR, get reviews, merge PR
- [ ] GitHub API integration:
  - Create `internal/integrations/github.go`
  - Authenticate with token from Skate
  - Implement PR detection
  - Monitor PR status (open, reviewed, merged)
- [ ] Add quest types:
  - "Create a Pull Request"
  - "Get PR Approved"
  - "Merge a PR"
- [ ] Update event system:
  - New event: `EventPRCreated`, `EventPRMerged`
  - Publish from GitHub watcher
- [ ] UI updates:
  - Display PR URL in quest
  - Show review status
- [ ] Tests:
  - Mock GitHub API
  - Test PR detection
  - Test quest completion

**Deliverable:**
- `internal/game/quest_pr.go`
- `internal/integrations/github.go`
- GitHub API integration tests

**Files to Create:**
- `internal/integrations/github.go`
- `internal/integrations/github_test.go`

**Files to Modify:**
- `internal/game/quest.go`
- `internal/game/events.go` (new event types)
- `internal/config/config.go` (GitHub config section)

**Dependencies:** None (new subsystem)

---

### Subagent 2.3: Refactoring Quest Type
**Duration:** 3-4 days
**Priority:** MEDIUM

**Description:** Add quests for code refactoring with AI validation.

**Tasks:**
- [ ] Define model:
  - `QuestTypeRefactor` constant
  - Criteria: refactor N files, improve code quality
- [ ] Refactor detection:
  - Analyze git diffs for refactor patterns:
    - Method extraction
    - Variable renaming
    - Code simplification
  - Use heuristics (lines deleted > added, complexity reduced)
- [ ] AI validation:
  - Send diff to AI for review
  - Ask: "Is this a refactoring? Rate quality 0-100"
  - Bonus XP for high-quality refactors
- [ ] Add quest types:
  - "Refactor 5 Functions"
  - "Simplify Complex Code"
  - "Extract Common Logic"
- [ ] UI updates:
  - Show refactor quality score
  - Display AI feedback
- [ ] Tests:
  - Test refactor detection
  - Mock AI responses
  - Test quality scoring

**Deliverable:**
- `internal/game/quest_refactor.go`
- Refactor detection logic
- AI-powered quality assessment

**Files to Modify:**
- `internal/game/quest.go`
- `internal/watcher/git.go` (refactor analysis)
- `internal/ai/review.go` (refactor validation)

**Dependencies:** Existing AI integration

---

### Subagent 2.4: Quest Chains & Prerequisites
**Duration:** 2-3 days
**Priority:** HIGH

**Description:** Implement quest dependencies and unlock system.

**Tasks:**
- [ ] Enhance Quest model:
  - `Prerequisites []string` (already exists, implement logic)
  - `UnlocksQuests []string` (already exists, implement logic)
- [ ] Prerequisite validation:
  - Check if required quests completed before showing
  - Mark quests as "locked" in UI
- [ ] Quest unlock system:
  - When quest completes, unlock dependent quests
  - Publish `EventQuestUnlocked` event
  - Show notification to user
- [ ] Create quest chains:
  - Define 2-3 starter chains in data/
  - Example: "First Commit" → "10 Commits" → "100 Commits"
- [ ] UI updates:
  - Show locked quests with icon
  - Display prerequisite requirements
  - Visual quest chain connections (optional)
- [ ] Tests:
  - Test prerequisite checking
  - Test quest unlocking
  - Test chain progression

**Deliverable:**
- Quest chain system functional
- Data files with starter chains
- UI showing locked/unlocked states

**Files to Modify:**
- `internal/game/quest.go` (implement prerequisite logic)
- `internal/ui/screens/questboard.go` (locked quest UI)
- `data/quests/` (quest chain definitions)

**Dependencies:** None (enhances existing quests)

---

### Subagent 2.5: Daily Quests System
**Duration:** 2-3 days
**Priority:** MEDIUM

**Description:** Add rotating daily challenges with bonus XP.

**Tasks:**
- [ ] Define daily quest model:
  - `QuestTypeDaily` constant
  - Daily generation logic
  - Reset mechanism
- [ ] Daily generation:
  - Generate 3 daily quests on login
  - Random selection from daily pool
  - Appropriate difficulty for player level
- [ ] Reset mechanism:
  - Check for date rollover
  - Clear yesterday's dailies
  - Generate new dailies
  - Notify player
- [ ] Streak tracking:
  - Track consecutive daily completions
  - Bonus XP for streaks (7-day, 30-day)
- [ ] UI updates:
  - "Daily Quests" section on Quest Board
  - Timer showing reset time
  - Streak counter display
- [ ] Tests:
  - Test daily generation
  - Test reset logic
  - Test streak tracking

**Deliverable:**
- Daily quest system functional
- Streak tracking with bonuses
- UI section for dailies

**Files to Create:**
- `internal/game/quest_daily.go`

**Files to Modify:**
- `internal/game/quest.go`
- `internal/ui/screens/questboard.go`
- `data/quests/daily/` (daily quest templates)

**Dependencies:** None (new subsystem)

---

### Subagent 2.6: Quest Difficulty Levels
**Duration:** 2 days
**Priority:** LOW

**Description:** Add easy/normal/hard difficulty variants for quests.

**Tasks:**
- [ ] Add difficulty field:
  - `Difficulty string` to Quest struct ("easy", "normal", "hard")
- [ ] Scale rewards by difficulty:
  - Easy: 0.5x XP
  - Normal: 1.0x XP
  - Hard: 2.0x XP
- [ ] Adjust targets:
  - Easy: Lower goals (e.g., 5 commits)
  - Normal: Standard (e.g., 10 commits)
  - Hard: Higher goals (e.g., 20 commits)
- [ ] UI updates:
  - Difficulty badges (🟢 Easy, 🟡 Normal, 🔴 Hard)
  - Display in quest list
  - Color-code quest cards
- [ ] Quest Board filtering:
  - Filter by difficulty
  - Show recommended difficulty based on level
- [ ] Tests:
  - Test reward scaling
  - Test difficulty display

**Deliverable:**
- Quest difficulty system
- UI badges and filtering

**Files to Modify:**
- `internal/game/quest.go` (add Difficulty field)
- `internal/game/engine.go` (difficulty XP multipliers)
- `internal/ui/screens/questboard.go` (difficulty UI)
- `data/quests/` (add difficulty to quest definitions)

**Dependencies:** None (enhances existing quests)

---

### Subagent 2.7: Quest Abandonment
**Duration:** 1-2 days
**Priority:** LOW

**Description:** Allow players to drop quests they don't want.

**Tasks:**
- [ ] Add abandon method:
  - `Abandon()` method to Quest
  - Set status to `QuestAbandoned` (new status)
  - Clear progress
- [ ] Cooldown system:
  - Track abandon time
  - Prevent re-accepting for N hours (e.g., 24h)
- [ ] UI implementation:
  - "Abandon Quest" button on quest detail
  - Confirmation modal: "Are you sure? Progress will be lost."
  - Display cooldown time if re-attempting
- [ ] Partial progress handling:
  - Option: Save partial XP (e.g., 25% of current progress)
  - Or: No reward for abandoned quests
- [ ] Tests:
  - Test abandon flow
  - Test cooldown
  - Test re-acceptance blocking

**Deliverable:**
- Quest abandonment feature
- Confirmation modal and cooldown

**Files to Modify:**
- `internal/game/quest.go` (Abandon method)
- `internal/ui/screens/questboard.go` (abandon UI)
- `internal/ui/components/modal.go` (confirmation modal)

**Dependencies:** None (enhances existing quests)

---

### Subagent 2.8: Phase 2 Testing & Integration
**Duration:** 2-3 days
**Priority:** HIGH

**Description:** Comprehensive testing of all new quest types and integration.

**Tasks:**
- [ ] Integration tests:
  - Test all 5 new quest types working together
  - Test quest chains with new types
  - Test daily quests with other active quests
  - Test abandoning different quest types
- [ ] UI testing:
  - Navigate all screens with new quests
  - Test difficulty filtering
  - Test locked quest display
  - Test daily quest section
- [ ] Performance testing:
  - Test with 10+ active quests
  - Test with 50+ completed quests
  - Profile quest progress updates
- [ ] Bug fixes:
  - Fix any issues found
  - Polish UX
  - Improve error messages
- [ ] Documentation:
  - Update README with new quest types
  - Add examples to docs/
  - Update CHANGELOG for v0.2.0
- [ ] Release preparation:
  - Tag v0.2.0
  - Build binaries
  - Create GitHub release

**Deliverable:**
- v0.2.0 release with Enhanced Quests
- All tests passing
- Updated documentation

**Files to Create:**
- `test/integration/quest_types_test.go`

**Files to Modify:**
- README.md (new features)
- CHANGELOG.md (v0.2.0 notes)

**Dependencies:** Subagents 2.1-2.7 (all new quest features)

---

## Phase 3: Skills & Achievements (2 weeks)

**Goal:** Add RPG progression depth with skill trees and achievements

### Subagent 3.1: Skill Tree Model
**Duration:** 3-4 days
**Priority:** HIGH

**Description:** Design and implement the skill tree system.

**Tasks:**
- [ ] Define Skill struct:
  ```go
  type Skill struct {
      ID            string
      Name          string
      Description   string
      Icon          string
      MaxLevel      int
      CurrentLevel  int
      Prerequisites []string  // Skill IDs
      Cost          int       // Skill points to unlock
      Effects       []SkillEffect
  }
  ```
- [ ] Define SkillEffect:
  ```go
  type SkillEffect struct {
      Type   string  // "xp_multiplier", "quest_bonus", etc.
      Value  float64
  }
  ```
- [ ] Skill tree structure:
  - Create data structure (tree, graph, or branching paths)
  - Define 3-5 skill trees (e.g., "Coder", "Debugger", "Collaborator")
- [ ] Skill point system:
  - Award skill points on level-up (e.g., 1 per level)
  - Track available vs. spent points
- [ ] Prerequisite validation:
  - Check if required skills unlocked
  - Validate skill point availability
- [ ] Skill application:
  - Apply effects when skill learned
  - Update character stats
- [ ] Save/load:
  - Persist skill tree state to Skate
  - Load on character initialization
- [ ] Tests:
  - Test skill unlocking
  - Test prerequisite checking
  - Test effect application
  - Test save/load

**Deliverable:**
- `internal/game/skills.go` with complete skill system
- Skill data definitions
- Tests with >80% coverage

**Files to Create:**
- `internal/game/skills.go`
- `internal/game/skills_test.go`
- `data/skills.json` (skill definitions)

**Files to Modify:**
- `internal/game/character.go` (add SkillPoints, Skills fields)
- `internal/storage/skate.go` (save/load skills)

**Dependencies:** None (new subsystem)

---

### Subagent 3.2: Skill Tree UI
**Duration:** 3-4 days
**Priority:** HIGH

**Description:** Build interactive skill tree visualization.

**Tasks:**
- [ ] Design layout:
  - Tree structure with Lip Gloss
  - Nodes for each skill
  - Lines connecting prerequisites
- [ ] Skill node rendering:
  - Show icon, name, level
  - Color-code: locked (gray), available (yellow), learned (green)
  - Display cost and effects on hover/selection
- [ ] Navigation:
  - Arrow keys to move between skills
  - Enter to view details
  - Space to unlock (if available)
- [ ] Skill detail modal:
  - Full description
  - Prerequisites list
  - Effects breakdown
  - "Unlock" or "Upgrade" button
- [ ] Skill point display:
  - Show available points in header
  - Update in real-time
- [ ] Multiple skill trees:
  - Tab navigation between trees
  - Or: Single scrollable view with all trees
- [ ] Animations:
  - Skill unlock animation
  - Visual feedback on point spend
- [ ] Tests:
  - Test UI rendering
  - Test navigation
  - Test unlock flow

**Deliverable:**
- `internal/ui/screens/skills.go` with interactive skill tree
- Skill tree accessible from Dashboard (I key or menu)

**Files to Create:**
- `internal/ui/screens/skills.go`

**Files to Modify:**
- `internal/ui/app.go` (add ScreenSkills)
- `internal/ui/keys.go` (add skill tree hotkey)
- `internal/ui/screens/dashboard.go` (add skill tree menu item)

**Dependencies:** Subagent 3.1 (need skill model)

---

### Subagent 3.3: Skill Effects System
**Duration:** 2-3 days
**Priority:** HIGH

**Description:** Implement skill bonuses and effects.

**Tasks:**
- [ ] Define effect types:
  - `xp_multiplier`: Increase XP gain (e.g., +10%)
  - `quest_reward_bonus`: Increase quest rewards (e.g., +20%)
  - `commit_quality_bonus`: Extra XP for good commits
  - `streak_protection`: Skip 1 day without breaking streak
  - `unlock_feature`: Enable special features (e.g., custom quests)
- [ ] Effect application:
  - Apply multipliers in XP engine
  - Apply bonuses in quest completion
  - Enable features in UI
- [ ] Passive bonuses:
  - Always active once learned
  - Compound with other bonuses
- [ ] Active abilities (future):
  - Skills with cooldowns
  - Manual activation
- [ ] UI display:
  - Show active effects on character screen
  - Tooltip explaining each effect
- [ ] Tests:
  - Test XP multipliers
  - Test quest bonuses
  - Test effect stacking

**Deliverable:**
- Skill effect system integrated with game engine
- Visible impact on gameplay

**Files to Modify:**
- `internal/game/engine.go` (apply skill effects)
- `internal/game/quest.go` (apply skill bonuses)
- `internal/ui/screens/character.go` (display active effects)

**Dependencies:** Subagent 3.1 (need skill model)

---

### Subagent 3.4: Achievement Model
**Duration:** 2 days
**Priority:** MEDIUM

**Description:** Create achievement tracking system.

**Tasks:**
- [ ] Define Achievement struct:
  ```go
  type Achievement struct {
      ID          string
      Name        string
      Description string
      Icon        string
      Criteria    AchievementCriteria
      Reward      AchievementReward
      Secret      bool  // Hidden until unlocked
      UnlockedAt  *time.Time
  }
  ```
- [ ] Define criteria types:
  - Reach level N
  - Complete N quests
  - Earn N total XP
  - Maintain N-day streak
  - Unlock specific skill
  - Make N commits
- [ ] Achievement checking:
  - Subscribe to relevant events
  - Check criteria on each event
  - Unlock if met
- [ ] Rewards:
  - Bonus XP
  - Skill points
  - Unlock cosmetics (future)
- [ ] Secret achievements:
  - Hidden from list until unlocked
  - Surprise bonus
- [ ] Save/load:
  - Persist achievement state
  - Track unlock timestamps
- [ ] Tests:
  - Test criteria checking
  - Test unlock logic
  - Test rewards

**Deliverable:**
- `internal/game/achievements.go` with tracking system
- Achievement definitions

**Files to Create:**
- `internal/game/achievements.go`
- `internal/game/achievements_test.go`
- `data/achievements.json`

**Files to Modify:**
- `internal/game/character.go` (add Achievements field)
- `internal/game/events.go` (add EventAchievementUnlock)
- `internal/storage/skate.go` (save/load achievements)

**Dependencies:** None (new subsystem)

---

### Subagent 3.5: Achievement UI
**Duration:** 2-3 days
**Priority:** MEDIUM

**Description:** Display achievements in game.

**Tasks:**
- [ ] Achievement list screen:
  - Grid or list layout
  - Show unlocked achievements (with date)
  - Show locked achievements (grayed out)
  - Hide secret achievements until unlocked
- [ ] Achievement card:
  - Icon, name, description
  - Unlock date
  - Reward received
- [ ] Progress tracking:
  - Show progress bars for in-progress achievements
  - E.g., "Complete 100 Quests: 45/100"
- [ ] Unlock animation:
  - Toast notification on unlock
  - Confetti or visual effect
  - Sound effect (optional)
- [ ] Recent achievements widget:
  - Show last 3 unlocked on Dashboard
  - Quick access to achievement screen
- [ ] Filters:
  - View by category
  - Unlocked/locked filter
- [ ] Tests:
  - Test UI rendering
  - Test unlock notifications

**Deliverable:**
- `internal/ui/screens/achievements.go` with achievement display
- Unlock notifications

**Files to Create:**
- `internal/ui/screens/achievements.go`

**Files to Modify:**
- `internal/ui/app.go` (add ScreenAchievements)
- `internal/ui/screens/dashboard.go` (add recent achievements widget)
- `internal/ui/components/modal.go` (achievement unlock modal)

**Dependencies:** Subagent 3.4 (need achievement model)

---

### Subagent 3.6: Achievement Data & Definitions
**Duration:** 1-2 days
**Priority:** LOW

**Description:** Create initial set of achievements.

**Tasks:**
- [ ] Design achievement categories:
  - Progression (level-based)
  - Questing (quest-based)
  - Coding (commit-based)
  - Dedication (streak-based)
  - Mastery (skill-based)
  - Secret (hidden goals)
- [ ] Create 20-30 starter achievements:
  - "First Steps" - Complete first quest
  - "Level 10" - Reach level 10
  - "Commit Streak" - 7-day streak
  - "Century" - 100 commits
  - "Test Master" - Complete 10 test quests
  - "Code Reviewer" - Get 5 perfect reviews
  - "Skill Collector" - Unlock 10 skills
  - "Quest Champion" - Complete 50 quests
  - Secret: "Night Owl" - Commit at 2AM
  - Secret: "Weekend Warrior" - 20 commits on weekend
- [ ] Define in JSON:
  - All fields populated
  - Criteria clearly specified
  - Rewards balanced
- [ ] Validation:
  - Ensure criteria are achievable
  - Balance rewards
  - Test loading system

**Deliverable:**
- `data/achievements.json` with 20-30 achievements
- Achievement loading validated

**Files to Create:**
- `data/achievements.json`

**Dependencies:** Subagent 3.4 (need achievement model defined)

---

### Subagent 3.7: Phase 3 Testing & Release
**Duration:** 2-3 days
**Priority:** HIGH

**Description:** Comprehensive testing of skills and achievements.

**Tasks:**
- [ ] Integration tests:
  - Test skill unlocking and effects
  - Test achievement unlocking
  - Test skill + achievement interaction
  - Test save/load of both systems
- [ ] Gameplay testing:
  - Unlock several skills
  - Verify effects apply correctly
  - Unlock achievements naturally
  - Test secret achievements
- [ ] Balance testing:
  - Verify skill costs are fair
  - Check effect magnitudes
  - Ensure achievements are achievable
- [ ] UI testing:
  - Navigate skill tree
  - Browse achievements
  - Test unlock animations
- [ ] Bug fixes:
  - Fix any issues found
  - Polish animations
  - Improve tooltips
- [ ] Documentation:
  - Document skill trees in README
  - Add achievement guide
  - Update CHANGELOG for v0.3.0
- [ ] Release:
  - Tag v0.3.0
  - Build binaries
  - Publish release

**Deliverable:**
- v0.3.0 release with Skills & Achievements
- Balanced and tested systems

**Files to Create:**
- `test/integration/skills_achievements_test.go`

**Files to Modify:**
- README.md (new features)
- CHANGELOG.md (v0.3.0 notes)

**Dependencies:** Subagents 3.1-3.6 (all Phase 3 features)

---

## Phase 4: Advanced UI (1.5 weeks)

**Goal:** Enhanced UI screens, animations, and help system

### Subagent 4.1: Character Sheet Screen Enhancement
**Duration:** 2-3 days
**Priority:** MEDIUM

**Description:** Create detailed character view with statistics.

**Tasks:**
- [ ] Enhance character screen:
  - Current: Basic stats display
  - Add: Detailed stat breakdown
  - Add: Progress history graph
  - Add: All-time statistics
- [ ] Statistics sections:
  - Core stats (Level, XP, CodePower, Wisdom, Agility)
  - Progress (Total commits, lines, quests)
  - Streaks (Current, longest)
  - Session stats (Today's commits, time)
  - Achievements (count, recent)
  - Skills (learned, points available)
- [ ] Visual enhancements:
  - Stat comparison bars
  - Progress graphs (commits over time)
  - Color-coded stats
- [ ] Equipment section (future):
  - Placeholder for items
  - "Coming soon" message
- [ ] Tests:
  - Test stat display
  - Test graph rendering

**Deliverable:**
- Enhanced `internal/ui/screens/character.go`
- Comprehensive character statistics

**Files to Modify:**
- `internal/ui/screens/character.go`
- `internal/ui/styles.go` (stat display styles)

**Dependencies:** Phase 3 complete (need skills/achievements)

---

### Subagent 4.2: Inventory/Skills Combined Screen
**Duration:** 2 days
**Priority:** LOW

**Description:** Create unified inventory view.

**Tasks:**
- [ ] Tab-based screen:
  - Tab 1: Skills (from Subagent 3.2)
  - Tab 2: Achievements (from Subagent 3.5)
  - Tab 3: Inventory (placeholder for items)
- [ ] Navigation:
  - Tab key to switch tabs
  - Arrow keys within each tab
- [ ] Unified header:
  - Show relevant stats for current tab
  - E.g., "Skill Points: 5" on Skills tab
- [ ] Quick stats:
  - Summary bar at bottom
  - Show counts: "12 Skills | 18 Achievements | 0 Items"
- [ ] Tests:
  - Test tab switching
  - Test navigation

**Deliverable:**
- Unified inventory screen
- Easy navigation between progression systems

**Files to Create:**
- `internal/ui/screens/inventory.go`

**Files to Modify:**
- `internal/ui/app.go` (ScreenInventory replaces separate screens)
- `internal/ui/keys.go` (inventory hotkey)

**Dependencies:** Subagents 3.2, 3.5 (need skills and achievements)

---

### Subagent 4.3: Settings Screen Enhancement
**Duration:** 2-3 days
**Priority:** MEDIUM

**Description:** Comprehensive settings UI.

**Tasks:**
- [ ] Settings categories:
  - AI Providers
  - Git Configuration
  - UI Preferences
  - Keybindings
  - Notifications
  - Advanced
- [ ] AI Provider settings:
  - Select provider (Crush/Mods/Claude)
  - Configure models
  - Test connection button
  - Manage API keys (via Skate, not shown)
- [ ] Git settings:
  - Add/remove watch paths
  - Auto-detect toggle
  - Ignore patterns
- [ ] UI preferences:
  - Theme (dark/light)
  - Animations (on/off)
  - Compact mode
  - Show keybind hints
- [ ] Keybinding customization:
  - List all keybinds
  - Click to rebind
  - Reset to defaults
- [ ] Save/apply:
  - Live preview changes
  - "Apply" and "Cancel" buttons
  - Write to config.toml
- [ ] Tests:
  - Test config updates
  - Test save/load

**Deliverable:**
- Comprehensive settings screen
- User-friendly configuration

**Files to Modify:**
- `internal/ui/screens/settings.go` (enhance existing)
- `internal/config/config.go` (validation)

**Dependencies:** None (enhances existing)

---

### Subagent 4.4: Animation System
**Duration:** 3-4 days
**Priority:** MEDIUM

**Description:** Add visual feedback animations.

**Tasks:**
- [ ] Animation framework:
  - Define animation types
  - Tick-based system with Bubble Tea
  - Easing functions
- [ ] Level-up animation:
  - Flash effect
  - "LEVEL UP!" banner
  - Stat increase display
  - Confetti particles (ASCII art)
- [ ] Achievement unlock animation:
  - Slide-in notification
  - Icon display
  - Achievement name + description
  - Auto-dismiss after 5 seconds
- [ ] XP gain animation:
  - "+50 XP" floating text
  - Progress bar fill animation
  - Smooth transitions
- [ ] Quest completion animation:
  - "Quest Complete!" banner
  - Reward display
  - Fanfare effect
- [ ] Particle system (optional):
  - ASCII particles for effects
  - Stars, dots, sparkles
- [ ] Performance:
  - Animations don't block UI
  - Can be disabled in settings
- [ ] Tests:
  - Test animation triggering
  - Test disable setting

**Deliverable:**
- Animation framework
- Delightful visual feedback

**Files to Create:**
- `internal/ui/animations.go`
- `internal/ui/particles.go` (optional)

**Files to Modify:**
- `internal/ui/app.go` (animation system integration)
- `internal/ui/screens/dashboard.go` (trigger animations)

**Dependencies:** None (new subsystem)

---

### Subagent 4.5: Help Overlay
**Duration:** 2 days
**Priority:** MEDIUM

**Description:** In-game help system.

**Tasks:**
- [ ] Help modal:
  - Trigger with `?` key (already in keys.go)
  - Overlay on current screen
  - Dim background
- [ ] Contextual help:
  - Show help relevant to current screen
  - E.g., Quest Board help shows quest actions
- [ ] Help sections:
  - Keybindings (auto-generated from keys.go)
  - Navigation guide
  - Feature explanations
  - Tips and tricks
- [ ] Search functionality:
  - Type to search help topics
  - Fuzzy matching
- [ ] Quick reference:
  - One-page cheat sheet
  - Printable format
- [ ] Tutorial mode (optional):
  - First-time user walkthrough
  - Highlight features
  - Tooltips
- [ ] Tests:
  - Test help display
  - Test search

**Deliverable:**
- Help overlay system
- Comprehensive user guidance

**Files to Create:**
- `internal/ui/components/help.go`
- `data/help/` (help content files)

**Files to Modify:**
- `internal/ui/app.go` (help modal integration)
- `internal/ui/keys.go` (ensure ? key handled)

**Dependencies:** None (new feature)

---

### Subagent 4.6: Phase 4 Polish & Release
**Duration:** 2 days
**Priority:** HIGH

**Description:** UI refinement and v0.4.0 release.

**Tasks:**
- [ ] UI testing:
  - Test all new screens
  - Test animations
  - Test help system
  - Test settings changes
- [ ] Performance testing:
  - Benchmark animation rendering
  - Profile UI updates
  - Optimize if needed
- [ ] UX improvements:
  - Smooth transitions
  - Consistent spacing
  - Color harmony
  - Accessibility features
- [ ] Bug fixes:
  - Fix any UI glitches
  - Fix animation issues
  - Improve responsiveness
- [ ] Documentation:
  - Screenshot new screens
  - Update README with UI features
  - Add UI guide to docs/
  - Update CHANGELOG
- [ ] Release:
  - Tag v0.4.0
  - Build binaries
  - Publish release

**Deliverable:**
- v0.4.0 release with Advanced UI
- Polished user experience

**Files to Modify:**
- README.md (UI screenshots)
- CHANGELOG.md (v0.4.0 notes)
- docs/UI_GUIDE.md (new)

**Dependencies:** Subagents 4.1-4.5 (all Phase 4 features)

---

## Phase 5: Advanced AI Features (1 week)

**Goal:** Enhanced AI integration for code review, quest generation, and learning

### Subagent 5.1: Enhanced Code Review System
**Duration:** 2-3 days
**Priority:** HIGH

**Description:** Deeper AI code analysis and feedback.

**Tasks:**
- [ ] Enhance review triggers:
  - Automatic after N commits
  - Manual trigger command
  - Review on quest completion
- [ ] Detailed analysis:
  - Code quality score (0-100)
  - Strengths list
  - Improvement areas
  - Go idiom suggestions
  - Performance recommendations
  - Security concerns
- [ ] Review storage:
  - Save review history
  - Track quality over time
  - Show trends
- [ ] Bonus XP system:
  - Award bonus XP for high-quality code (>80 score)
  - Scale bonus by review score
  - Show bonus in notification
- [ ] Review UI:
  - Display review in modal
  - Syntax highlighting for suggestions
  - "Apply suggestion" button (future)
- [ ] Multiple AI models:
  - Use complex model for deep reviews
  - Use simple model for quick checks
- [ ] Tests:
  - Mock AI responses
  - Test bonus XP calculation
  - Test review storage

**Deliverable:**
- Enhanced `internal/ai/review.go`
- Detailed code review system

**Files to Modify:**
- `internal/ai/review.go`
- `internal/watcher/integration.go` (trigger reviews)
- `internal/game/engine.go` (bonus XP)
- `internal/ui/components/modal.go` (review display)

**Dependencies:** Existing AI integration

---

### Subagent 5.2: AI Quest Generation
**Duration:** 2-3 days
**Priority:** MEDIUM

**Description:** Generate custom quests using AI.

**Tasks:**
- [ ] Quest generation prompt:
  - Input: topic, difficulty, player level
  - Output: JSON quest definition
  - Example: "Generate a quest about testing"
- [ ] Quest validation:
  - Verify JSON format
  - Check quest is achievable
  - Validate rewards
  - Ensure description is clear
- [ ] Generation UI:
  - "Generate Quest" button on Quest Board
  - Input field for topic/keywords
  - Difficulty slider
  - "Generate" button
  - Preview before accepting
- [ ] Quest library:
  - Save generated quests to library
  - Share quests (export/import JSON)
  - Community quest packs (future)
- [ ] AI provider selection:
  - Use Claude for complex quests
  - Use Crush for simple quests
  - Fallback chain
- [ ] Tests:
  - Mock AI responses
  - Test quest validation
  - Test generation flow

**Deliverable:**
- `internal/ai/questgen.go` with generation system
- UI for quest creation

**Files to Create:**
- `internal/ai/questgen.go`

**Files to Modify:**
- `internal/ui/screens/questboard.go` (add generation button)
- `internal/game/quest.go` (import generated quests)

**Dependencies:** Existing AI integration

---

### Subagent 5.3: Learning Path Suggestions
**Duration:** 2 days
**Priority:** LOW

**Description:** AI-powered learning recommendations.

**Tasks:**
- [ ] Commit pattern analysis:
  - Analyze last N commits
  - Identify patterns:
    - Languages used
    - File types
    - Commit frequency
    - Code complexity
- [ ] Improvement areas:
  - Identify weaknesses
  - E.g., "Low test coverage", "No documentation"
- [ ] Personalized suggestions:
  - Suggest relevant quests
  - Suggest skills to learn
  - Suggest tutorials
- [ ] Learning path generation:
  - Create multi-quest learning path
  - E.g., "Master Testing" path:
    1. Write your first test
    2. Achieve 50% coverage
    3. Use table-driven tests
    4. Add integration tests
- [ ] UI display:
  - "Suggested for You" section on Dashboard
  - Learning path detail view
  - "Start Path" button
- [ ] Tests:
  - Test pattern analysis
  - Mock AI suggestions

**Deliverable:**
- Learning path system
- Personalized recommendations

**Files to Create:**
- `internal/ai/learning.go`

**Files to Modify:**
- `internal/ui/screens/dashboard.go` (suggestions widget)
- `internal/watcher/git.go` (commit analysis)

**Dependencies:** Subagent 5.2 (quest generation)

---

### Subagent 5.4: Tutorial Quest Creation
**Duration:** 2 days
**Priority:** LOW

**Description:** AI-generated tutorials for learning.

**Tasks:**
- [ ] Tutorial quest type:
  - Interactive learning quests
  - Step-by-step instructions
  - Check understanding questions
- [ ] Tutorial generation:
  - Input: topic (e.g., "Go interfaces")
  - Generate multi-step tutorial
  - Include code examples
  - Add practice exercises
- [ ] Built-in tutorials:
  - Go basics
  - TUI development
  - Git workflows
  - Testing practices
- [ ] Tutorial UI:
  - Step-by-step display
  - Code snippets
  - "Next Step" button
  - Progress tracking
- [ ] Completion verification:
  - Optional: AI checks if user completed step
  - Self-reported completion
- [ ] Tests:
  - Test tutorial loading
  - Test step progression

**Deliverable:**
- Tutorial quest system
- AI-generated learning content

**Files to Create:**
- `internal/game/quest_tutorial.go`

**Files to Modify:**
- `internal/ai/questgen.go` (tutorial generation)
- `internal/ui/screens/questboard.go` (tutorial display)
- `data/tutorials/` (built-in tutorials)

**Dependencies:** Subagent 5.2 (quest generation)

---

### Subagent 5.5: Phase 5 Testing & Release
**Duration:** 1-2 days
**Priority:** HIGH

**Description:** Test AI features and release v0.5.0.

**Tasks:**
- [ ] AI testing:
  - Test code review quality
  - Test quest generation
  - Test learning suggestions
  - Test tutorial quests
- [ ] Provider testing:
  - Test Crush (online)
  - Test Mods (offline)
  - Test Claude fallback
  - Test rate limiting
- [ ] Integration testing:
  - Test AI + game logic
  - Test AI + UI
  - Test error handling
- [ ] Bug fixes:
  - Fix AI response parsing
  - Handle API errors gracefully
  - Improve prompts
- [ ] Documentation:
  - Document AI features
  - Add AI setup guide
  - Update CHANGELOG
- [ ] Release:
  - Tag v0.5.0
  - Build binaries
  - Publish release

**Deliverable:**
- v0.5.0 release with Advanced AI
- Robust AI integration

**Files to Modify:**
- README.md (AI features)
- docs/AI_SETUP.md (new)
- CHANGELOG.md (v0.5.0 notes)

**Dependencies:** Subagents 5.1-5.4 (all Phase 5 features)

---

## Phase 6: External Integrations (1 week)

**Goal:** Integrate with GitHub, WakaTime, and add data export

### Subagent 6.1: GitHub API Integration
**Duration:** 3-4 days
**Priority:** HIGH

**Description:** Full GitHub integration for issues and PRs.

**Tasks:**
- [ ] Enhance GitHub client (from Subagent 2.2):
  - Issue tracking
  - PR monitoring
  - Commit verification
  - Repository stats
- [ ] Issue sync:
  - Fetch repository issues
  - Convert issues to quests
  - Update quest when issue closed
  - Bi-directional sync
- [ ] PR integration:
  - Track PR creation
  - Monitor review status
  - Detect PR merge
  - Award XP on merge
- [ ] Repository stats:
  - Commit count
  - PR count
  - Issue count
  - Contributor stats
- [ ] Settings:
  - Enable/disable GitHub integration
  - Select repositories to watch
  - Configure sync frequency
- [ ] API key management:
  - Store token in Skate
  - Validate token on setup
  - Handle rate limiting
- [ ] Tests:
  - Mock GitHub API
  - Test issue sync
  - Test PR tracking

**Deliverable:**
- Full GitHub integration
- Issue-to-quest sync

**Files to Modify:**
- `internal/integrations/github.go` (enhance from 2.2)
- `internal/integrations/github_test.go`
- `internal/config/config.go` (GitHub settings)
- `internal/ui/screens/settings.go` (GitHub config UI)

**Dependencies:** Subagent 2.2 (basic GitHub integration)

---

### Subagent 6.2: WakaTime Integration
**Duration:** 2-3 days
**Priority:** LOW

**Description:** Optional time tracking with WakaTime.

**Tasks:**
- [ ] WakaTime API client:
  - Authenticate with API key (from Skate)
  - Fetch time stats
  - Get project breakdown
  - Get language stats
- [ ] Session sync:
  - Compare CodeQuest session time vs. WakaTime
  - Reconcile differences
  - Display both in Dashboard
- [ ] Detailed stats:
  - Time by project
  - Time by language
  - Time by day/week/month
  - Coding activity graph
- [ ] UI integration:
  - WakaTime widget on Dashboard
  - Detailed stats in Character screen
  - "Powered by WakaTime" badge
- [ ] Settings:
  - Enable/disable WakaTime
  - Configure API key
  - Select projects to track
- [ ] Tests:
  - Mock WakaTime API
  - Test stat fetching
  - Test sync logic

**Deliverable:**
- `internal/integrations/wakatime.go`
- Optional WakaTime integration

**Files to Create:**
- `internal/integrations/wakatime.go`
- `internal/integrations/wakatime_test.go`

**Files to Modify:**
- `internal/config/config.go` (WakaTime settings)
- `internal/ui/screens/dashboard.go` (WakaTime widget)
- `internal/ui/screens/settings.go` (WakaTime config UI)

**Dependencies:** None (new subsystem)

---

### Subagent 6.3: Data Export System
**Duration:** 2 days
**Priority:** MEDIUM

**Description:** Export stats to various formats.

**Tasks:**
- [ ] Export formats:
  - JSON (full data dump)
  - CSV (spreadsheet-friendly)
  - Markdown (human-readable reports)
  - HTML (rich reports with graphs)
- [ ] Export data:
  - Character stats
  - Quest history
  - Achievements
  - Skills
  - Commit history
  - Session stats
  - AI review history
- [ ] Report generation:
  - Daily summary
  - Weekly summary
  - Monthly summary
  - All-time stats
- [ ] Chart generation:
  - Commits over time (ASCII art)
  - XP progression graph
  - Quest completion rate
- [ ] Export UI:
  - "Export Data" button in Settings
  - Select format
  - Select date range
  - Select data types
  - Save to file
- [ ] Tests:
  - Test each export format
  - Verify data completeness

**Deliverable:**
- Data export functionality
- Multiple format support

**Files to Create:**
- `internal/export/exporter.go`
- `internal/export/formats.go`
- `internal/export/exporter_test.go`

**Files to Modify:**
- `internal/ui/screens/settings.go` (export UI)

**Dependencies:** None (reads existing data)

---

### Subagent 6.4: Phase 6 Testing & Release
**Duration:** 1-2 days
**Priority:** HIGH

**Description:** Integration testing and v0.6.0 release.

**Tasks:**
- [ ] Integration testing:
  - Test GitHub sync
  - Test WakaTime sync
  - Test data export
  - Test with rate limiting
- [ ] Error handling:
  - Test API failures
  - Test network issues
  - Test invalid credentials
- [ ] Performance:
  - Test with large data sets
  - Optimize API calls
  - Cache where appropriate
- [ ] Bug fixes:
  - Fix integration issues
  - Improve error messages
  - Handle edge cases
- [ ] Documentation:
  - Document GitHub setup
  - Document WakaTime setup
  - Document data export
  - Update CHANGELOG
- [ ] Release:
  - Tag v0.6.0
  - Build binaries
  - Publish release

**Deliverable:**
- v0.6.0 release with Integrations
- External API support

**Files to Modify:**
- README.md (integration features)
- docs/INTEGRATIONS.md (new)
- CHANGELOG.md (v0.6.0 notes)

**Dependencies:** Subagents 6.1-6.3 (all Phase 6 features)

---

## Phase 7: Final Polish (1 week)

**Goal:** Production-ready v1.0.0 with comprehensive docs and demos

### Subagent 7.1: VHS Demo Recordings
**Duration:** 1-2 days
**Priority:** MEDIUM

**Description:** Create demo videos and GIFs.

**Tasks:**
- [ ] Install VHS: `brew install vhs`
- [ ] Write tape scripts:
  - `demo.tape` - Full application demo
  - `installation.tape` - Installation process
  - `quest-flow.tape` - Quest acceptance to completion
  - `level-up.tape` - Level-up animation
  - `ai-mentor.tape` - Using AI mentor
  - `skills.tape` - Skill tree navigation
- [ ] Record demos:
  - Run each tape script
  - Generate GIFs
  - Optimize file sizes
- [ ] Create assets:
  - Store in `assets/demos/`
  - Embed in README
  - Add to docs
- [ ] YouTube video (optional):
  - Longer form demo
  - Narrated walkthrough
- [ ] Tests:
  - Verify all demos run correctly

**Deliverable:**
- Demo GIFs and videos
- VHS tape scripts

**Files to Create:**
- `scripts/recordings/demo.tape`
- `scripts/recordings/installation.tape`
- `scripts/recordings/quest-flow.tape`
- `scripts/recordings/level-up.tape`
- `scripts/recordings/ai-mentor.tape`
- `scripts/recordings/skills.tape`
- `assets/demos/` (GIF outputs)

**Dependencies:** Complete application (Phases 0-6)

---

### Subagent 7.2: Comprehensive Documentation
**Duration:** 2-3 days
**Priority:** HIGH

**Description:** Complete documentation suite.

**Tasks:**
- [ ] API documentation:
  - Generate godoc
  - Host on pkg.go.dev
  - Add code examples
- [ ] Architecture guide:
  - System overview
  - Component interaction
  - Data flow diagrams
  - Design decisions
- [ ] Contributing guide:
  - How to contribute
  - Code style guide
  - PR process
  - Development setup
- [ ] Developer docs:
  - Adding new quest types
  - Creating custom UI components
  - Extending AI integration
  - Plugin development (future)
- [ ] User guides:
  - Getting started (enhanced)
  - Configuration guide
  - Troubleshooting
  - FAQ
  - Tips and tricks
- [ ] Update all docs:
  - README.md (final polish)
  - CLAUDE.md (development guide)
  - CODEQUEST_SPEC.md (mark complete)
  - DEVELOPMENT_STATUS.md (v1.0.0 status)
- [ ] Doc site (optional):
  - Create with Hugo/MkDocs
  - Deploy to GitHub Pages

**Deliverable:**
- Complete docs/ directory
- Comprehensive documentation

**Files to Create:**
- `docs/ARCHITECTURE.md`
- `docs/CONTRIBUTING.md`
- `docs/DEVELOPMENT.md`
- `docs/API.md`
- `docs/PLUGIN_GUIDE.md` (future)

**Files to Modify:**
- README.md (final polish)
- CLAUDE.md (v1.0.0 notes)
- CODEQUEST_SPEC.md (mark complete)
- DEVELOPMENT_STATUS.md (v1.0.0 status)

**Dependencies:** Complete application (Phases 0-6)

---

### Subagent 7.3: CLI Helpers with Gum
**Duration:** 2 days
**Priority:** LOW

**Description:** Enhanced CLI commands with Gum.

**Tasks:**
- [ ] Install Gum support: `brew install gum`
- [ ] Interactive setup:
  - `codequest setup` - Guided configuration
  - Use Gum input/select for choices
  - Create config.toml
  - Set up API keys in Skate
  - Test connections
- [ ] Quest creation wizard:
  - `codequest quest create` - Interactive quest creation
  - Gum forms for quest details
  - Save to custom quest file
- [ ] Config wizard:
  - `codequest config` - Interactive config editing
  - Gum choose for options
  - Update config.toml
- [ ] Status commands:
  - `codequest status` - Quick stats display
  - `codequest stats` - Detailed statistics
  - `codequest quest list` - List quests
  - `codequest skills` - Show skills
  - Gum formatted output
- [ ] Tests:
  - Test CLI commands
  - Test Gum integration

**Deliverable:**
- Enhanced CLI with Gum
- User-friendly commands

**Files to Create:**
- `internal/cli/setup.go`
- `internal/cli/wizard.go`
- `internal/cli/status.go`

**Files to Modify:**
- `cmd/codequest/main.go` (add CLI commands)

**Dependencies:** None (enhances existing CLI)

---

### Subagent 7.4: Error Handling & Recovery
**Duration:** 2 days
**Priority:** HIGH

**Description:** Robust error handling throughout application.

**Tasks:**
- [ ] Graceful degradation:
  - App works without Skate (memory storage)
  - App works without AI (disable features)
  - App works without git (manual mode)
- [ ] Error recovery:
  - Recover from panics
  - Log errors to file
  - Display user-friendly messages
  - Offer recovery actions
- [ ] Better error messages:
  - Clear, actionable errors
  - Include fix suggestions
  - Provide documentation links
- [ ] Crash reporting:
  - Optional anonymous crash reports
  - Stack traces
  - System info
  - Opt-in only
- [ ] Data recovery:
  - Backup character data
  - Restore from backup
  - Corruption detection
- [ ] Tests:
  - Test error scenarios
  - Test recovery mechanisms

**Deliverable:**
- Robust error handling
- Enhanced reliability

**Files to Modify:**
- `internal/storage/skate.go` (graceful fallback)
- `internal/ai/manager.go` (AI fallback)
- `internal/watcher/git.go` (git error handling)
- `cmd/codequest/main.go` (panic recovery)

**Dependencies:** None (improves existing code)

---

### Subagent 7.5: Performance Optimization
**Duration:** 2 days
**Priority:** MEDIUM

**Description:** Optimize hot paths and reduce resource usage.

**Tasks:**
- [ ] Profiling:
  - CPU profiling
  - Memory profiling
  - Goroutine profiling
  - Identify bottlenecks
- [ ] Optimization targets:
  - UI rendering (reduce redraws)
  - Git watcher (efficient polling)
  - Event bus (minimize allocations)
  - Quest checking (cache results)
- [ ] Memory optimization:
  - Pool frequently allocated objects
  - Reduce string allocations
  - Optimize data structures
- [ ] Performance benchmarks:
  - Benchmark critical paths
  - Set performance targets
  - Regression testing
- [ ] Startup time:
  - Lazy loading
  - Parallel initialization
  - Reduce dependencies
- [ ] Tests:
  - Benchmark tests
  - Performance regression tests

**Deliverable:**
- Optimized application
- Performance benchmarks

**Files to Modify:**
- All `*_test.go` (add benchmarks)
- `internal/ui/app.go` (optimize rendering)
- `internal/watcher/git.go` (optimize polling)
- `internal/game/events.go` (optimize pub/sub)

**Dependencies:** None (improves existing code)

---

### Subagent 7.6: v1.0.0 Release
**Duration:** 1 day
**Priority:** CRITICAL

**Description:** Major stable release preparation.

**Tasks:**
- [ ] Final testing:
  - Full application walkthrough
  - All features tested
  - All integrations verified
  - Documentation reviewed
- [ ] Version updates:
  - Update version to 1.0.0 in code
  - Update all documentation
  - Finalize CHANGELOG
- [ ] Release notes:
  - Comprehensive feature list
  - Migration guide (if needed)
  - Known issues
  - Roadmap preview
- [ ] Build artifacts:
  - Build for all platforms
  - Create installers
  - Sign binaries (if applicable)
- [ ] Release preparation:
  - Tag v1.0.0
  - Create GitHub release
  - Attach binaries
  - Publish release notes
- [ ] Announcement:
  - Blog post (if applicable)
  - Social media
  - Reddit/HN (if appropriate)
  - Email list (if exists)
- [ ] Post-release:
  - Monitor for issues
  - Respond to feedback
  - Plan v1.1.0

**Deliverable:**
- v1.0.0 stable release
- Public launch

**Files to Modify:**
- All version numbers
- README.md (v1.0.0 announcement)
- CHANGELOG.md (v1.0.0 complete notes)
- DEVELOPMENT_STATUS.md (v1.0.0 status)

**Dependencies:** Subagents 7.1-7.5 (all Phase 7 complete)

---

## Immediate Next Steps

**Start Here:** Once you begin a fresh session:

1. **Phase 0, Subagent 0.1** - End-to-End Integration Validation
   - Build and test the application
   - Verify all features work
   - Document any issues

2. **Phase 0, Subagent 0.2** - Pre-Release Checklist
   - Clean up pending commits
   - Verify documentation
   - Run full test suite

3. **Phase 0, Subagent 0.3** - Release v0.1.0
   - Tag the release
   - Build binaries
   - Publish on GitHub

**Then:** Move to Phase 1 for user testing and iteration.

---

## Summary

**Total Development Time:** ~10 weeks (50-60 working days)

**Total Subagents:** 42

**Major Milestones:**
- ✅ v0.1.0 - MVP Complete (Week 0)
- v0.2.0 - Enhanced Quests (Week 3)
- v0.3.0 - Skills & Achievements (Week 5)
- v0.4.0 - Advanced UI (Week 6.5)
- v0.5.0 - Advanced AI (Week 7.5)
- v0.6.0 - Integrations (Week 8.5)
- v1.0.0 - Stable Release (Week 10)

**This roadmap takes CodeQuest from a complete MVP to a feature-rich, production-ready developer productivity RPG!** 🎮⚔️🚀
