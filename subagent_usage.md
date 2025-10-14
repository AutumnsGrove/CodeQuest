
# Subagent Usage Guidelines for Claude Code

## Core Principle: Subagent-First Development

**MANDATORY**: Use subagents extensively for ALL non-trivial tasks. Each subagent handles ONE focused responsibility within strict token limits, ensuring maximum quality and context efficiency.

## 🚨 CRITICAL PROCESS ORDERING 🚨

Tasks MUST proceed through these phases in strict order:

1. **RESEARCH PHASE** (Complete ALL research before ANY development)
2. **DEVELOPMENT PHASE** (Complete ALL development before ANY testing)
3. **TESTING PHASE** (Only after development is finalized)

**NEVER** skip ahead. **NEVER** mix phases. **NEVER** write code during research.

---

## Phase 1: Research (ALWAYS FIRST)

### Required Research Subagents

Before writing ANY code, spawn these research subagents:

1. **Requirements Analysis Subagent**
   - Extract and clarify all requirements
   - Identify ambiguities and assumptions
   - Define success criteria
   - Output: `requirements.md`
   - Commit Type: `docs`

2. **Technical Research Subagent**
   - Research best practices for the technology stack
   - Identify potential libraries/tools
   - Research common pitfalls and solutions
   - Output: `technical_research.md`
   - Commit Type: `docs`

3. **Architecture Planning Subagent**
   - Design high-level architecture
   - Define component boundaries
   - Plan data flow and interfaces
   - Output: `architecture_plan.md`
   - Commit Type: `docs`

4. **Edge Case Analysis Subagent**
   - Identify potential edge cases
   - Research error handling strategies
   - Define validation requirements
   - Output: `edge_cases.md`
   - Commit Type: `docs`

### Research Phase Completion Gate

✅ All research subagents have completed
✅ All research artifacts are created and committed
✅ No unresolved questions remain
✅ Architecture is fully planned
✅ All research commits are in git history

**Only then proceed to Development Phase**

---

## Phase 2: Development (AFTER Research Complete)

### Development Subagent Strategy

1. **Component Isolation**: Each major component gets its own subagent
2. **Sequential Building**: Build dependencies before dependents
3. **Incremental Integration**: Integrate components one at a time
4. **Atomic Commits**: Each subagent commits its work before handoff

### Required Development Subagents

1. **Core Infrastructure Subagent**
   - Set up project structure
   - Create base configurations
   - Implement shared utilities
   - Context: Research artifacts + requirements
   - Commit Type: `feat` or `chore`

2. **Component Development Subagents** (one per component)
   - Implement single component/module
   - Follow architecture plan exactly
   - Include inline documentation
   - Context: Component-specific requirements + interfaces
   - Commit Type: `feat`

3. **Integration Subagent**
   - Connect components
   - Implement data flow
   - Handle inter-component communication
   - Context: All component outputs + architecture
   - Commit Type: `feat` or `refactor`

4. **Documentation Subagent**
   - Create user documentation
   - Write API documentation
   - Create usage examples
   - Context: Completed codebase summary
   - Commit Type: `docs`

### Development Phase Completion Gate

✅ All components implemented
✅ All integrations complete
✅ Documentation created
✅ Code is functional (not tested)
✅ All development commits are in git history

**Only then proceed to Testing Phase**

---

## Phase 3: Testing (AFTER Development Complete)

### Testing Subagent Strategy

1. **Test Planning Subagent**
   - Create comprehensive test plan
   - Define test cases for each component
   - Plan integration test scenarios
   - Output: `test_plan.md`
   - Commit Type: `docs`

2. **Unit Test Development Subagent**
   - Write unit tests per test plan
   - One subagent per component's tests
   - Include edge case coverage
   - Context: Component code + test plan
   - Commit Type: `test`

3. **Integration Test Subagent**
   - Write integration tests
   - Test component interactions
   - Validate data flow
   - Context: Integration points + test plan
   - Commit Type: `test`

4. **Test Execution Subagent**
   - Run all tests
   - Document results
   - Identify any failures
   - Output: `test_results.md`
   - Commit Type: `docs` or `chore`

---

## Git Commit Protocol

### 🔴 MANDATORY: Commit Before Handoff

Every subagent MUST commit its work before producing the completion artifact.

**Process Flow:**
1. Complete your assigned task
2. Review all files created/modified
3. **Read `GIT_COMMIT_STYLE_GUIDE.md`**
4. Determine appropriate commit type
5. Write commit message following the guide
6. Stage changes with `git add`
7. Execute `git commit`
8. Capture commit hash
9. Include commit information in completion artifact
10. Produce handoff documentation

### Commit Type Mapping by Subagent

| Subagent Type | Commit Type | Example Message |
|---------------|-------------|-----------------|
| **Research Phase** |
| Requirements Analysis | `docs` | `docs: Add requirements analysis and success criteria` |
| Technical Research | `docs` | `docs: Add technical research for authentication system` |
| Architecture Planning | `docs` | `docs: Add system architecture plan` |
| Edge Case Analysis | `docs` | `docs: Document edge cases and validation requirements` |
| **Development Phase** |
| Core Infrastructure | `feat` or `chore` | `feat: Set up project structure with configuration` |
| Component Development | `feat` | `feat: Implement JWT authentication module` |
| Bug Fix | `fix` | `fix: Correct token expiration validation logic` |
| Integration | `feat` or `refactor` | `feat: Integrate auth module with API endpoints` |
| Code Refactoring | `refactor` | `refactor: Extract validation logic into helpers` |
| Documentation | `docs` | `docs: Add API documentation with usage examples` |
| **Testing Phase** |
| Test Planning | `docs` | `docs: Add comprehensive test plan for auth system` |
| Unit Test Development | `test` | `test: Add unit tests for authentication module` |
| Integration Tests | `test` | `test: Add integration tests for login flow` |
| Test Execution | `docs` or `chore` | `docs: Add test execution results` |
| **Special Cases** |
| Performance Improvement | `perf` | `perf: Optimize database query for user lookup` |
| Styling/Formatting | `style` | `style: Format code with Black` |
| Build Configuration | `build` | `build: Add Docker configuration` |
| CI/CD Updates | `ci` | `ci: Add GitHub Actions workflow` |

### Commit Message Template for Subagents

```bash
# Good commit message structure:
<type>: <short description (≤50 chars)>

<optional body explaining what and why>

Subagent: [Subagent Name/Role]
Phase: [Research/Development/Testing]
```

### Examples of Good Subagent Commits

```bash
# Research Phase
git commit -m "docs: Add requirements analysis

Extracted and clarified all requirements from initial brief.
Identified 3 key ambiguities requiring user clarification.
Defined success criteria for MVP.

Subagent: Requirements Analysis
Phase: Research"

# Development Phase
git commit -m "feat: Implement user authentication module

Added JWT-based authentication with:
- Login/logout endpoints
- Token generation and validation
- Password hashing with bcrypt
- Middleware for protected routes

Subagent: Component Development (Auth)
Phase: Development"

# Testing Phase
git commit -m "test: Add unit tests for authentication

Covers:
- Token generation and validation
- Password hashing verification
- Edge cases for expired tokens
- Invalid credentials handling

Test coverage: 95%

Subagent: Unit Test Development (Auth)
Phase: Testing"

# Integration
git commit -m "feat: Integrate authentication with API

Connected auth module to existing API endpoints.
Added middleware to protect routes requiring authentication.
Updated API documentation with auth requirements.

Subagent: Integration
Phase: Development"

# Documentation
git commit -m "docs: Add API documentation with examples

Created comprehensive API docs including:
- Authentication flow diagrams
- Endpoint specifications
- Code examples for common use cases
- Error handling guide

Subagent: Documentation
Phase: Development"
```

### Commit Guidelines for Subagents

✅ **DO:**
- Read `GIT_COMMIT_STYLE_GUIDE.md` before committing
- Use imperative mood: "Add feature" not "Added feature"
- Be specific about what changed
- Keep subject line under 50 characters
- Add body for complex changes explaining "what" and "why"
- Include subagent name and phase in body
- Commit only successful, complete work
- Review all changes before committing
- Use appropriate commit type from the guide

❌ **DON'T:**
- Commit with vague messages like "Update files" or "WIP"
- Commit failed or incomplete work
- Use past tense: "Added feature"
- Skip reading the style guide
- Exceed 50 characters in subject line
- End subject line with a period
- Commit during error recovery (wait for successful retry)
- Mix multiple concerns in one commit

### When NOT to Commit

- ❌ Subagent failed to complete task
- ❌ No files were created or modified
- ❌ During error recovery/analysis phase
- ❌ Coordinator subagent analyzing failures (no actual changes)
- ❌ Work is incomplete or doesn't meet success criteria

### Git Commands for Subagents

```bash
# 1. Check status and review changes
git status
git diff

# 2. Stage changes
git add path/to/file1.py path/to/file2.md
# Or stage all changes
git add .

# 3. Commit with message
git commit -m "type: short description

Longer explanation if needed.

Subagent: [Name]
Phase: [Phase]"

# 4. Capture commit hash
git log -1 --format="%H"
```

---

## Context Management Rules

### 1. Context Inheritance Hierarchy

```
Project Requirements (base context)
    ↓
Research Artifacts (+ commit hashes)
    ↓
Architecture Plan (+ commit hash)
    ↓
Component Interfaces (+ commit hashes)
    ↓
Implementation Details (+ commit hashes)
```

### 2. Context Pruning Strategy

- **Pass Forward**: Only essential outputs, not full code
- **Summarize**: Previous phase outputs into bullet points
- **Reference**: Use file references and commit hashes instead of inline content
- **Focus**: Each subagent gets ONLY relevant context
- **Git History**: Reference previous subagent commits for context

### 3. Token Optimization

Each subagent should receive:
- Core requirements (< 500 tokens)
- Phase-specific context (< 2000 tokens)
- Relevant previous outputs (< 1000 tokens)
- Clear instructions (< 500 tokens)
- Commit references (< 100 tokens)

**Total context per subagent: ~4000 tokens maximum**

---

## Handoff Protocol

### Subagent Output Format

Every subagent MUST produce:

```markdown
## Completion Artifact: [Subagent Name]

### Git Commit Information
**Commit Hash**: `abc123def456789...` (full hash)
**Commit Type**: `[type]`
**Commit Message**: 
```
[type]: [description]

[optional body]
```

**Files Changed**:
- `path/to/file1.py` (created, +150 lines)
- `path/to/file2.md` (modified, +45/-12 lines)
- `path/to/file3.ts` (created, +200 lines)

### Summary
[2-3 sentences of what was accomplished]

### Key Outputs
- [List of created files/artifacts]
- [Important decisions made]
- [Critical findings]

### Key Decisions Made
- [Decision 1 and rationale]
- [Decision 2 and rationale]

### Next Steps
- [Recommended next subagent]
- [Required inputs for next phase]
- [Context to pass forward]

### Concerns/Blockers
- [Any issues encountered]
- [Risks identified]
- [Technical debt incurred]

### Confidence Score
[High/Medium/Low] - [Explanation]

### Success Criteria Met
- [x] [Criterion 1]
- [x] [Criterion 2]
- [ ] [Criterion 3 - not met, reason]
```

### Handoff Rules

1. **Explicit Handoff**: Previous subagent explicitly states completion
2. **Artifact Validation**: Verify all expected outputs exist
3. **Commit Verification**: Confirm commit was successful
4. **Context Transfer**: Next subagent receives completion artifact + relevant context
5. **No Backtracking**: Cannot return to previous phase without explicit rollback
6. **Commit Reference**: Next subagent can reference previous commits

---

## Common Task Decompositions

### Web Application

1. **Research Phase** (5-6 subagents, ~6 commits):
   - Requirements analysis → `docs: Add requirements analysis`
   - Frontend framework research → `docs: Add frontend research`
   - Backend architecture research → `docs: Add backend architecture research`
   - Database design research → `docs: Add database design research`
   - Security requirements research → `docs: Add security requirements`
   - Deployment strategy research → `docs: Add deployment strategy`

2. **Development Phase** (8-10 subagents, ~10-12 commits):
   - Project setup → `chore: Initialize project structure`
   - Database schema → `feat: Add database schema`
   - Backend API (one per major endpoint group)
     - Auth endpoints → `feat: Implement authentication endpoints`
     - User endpoints → `feat: Implement user management endpoints`
     - Data endpoints → `feat: Implement data processing endpoints`
   - Frontend components (one per major feature)
     - Login component → `feat: Add login component`
     - Dashboard → `feat: Add dashboard component`
     - Data visualization → `feat: Add data visualization components`
   - Styling/UI → `feat: Add UI styling and theme`
   - Integration → `feat: Integrate frontend with API`
   - Configuration → `chore: Add environment configuration`

3. **Testing Phase** (4-5 subagents, ~5-6 commits):
   - Test planning → `docs: Add comprehensive test plan`
   - Backend unit tests → `test: Add backend unit tests`
   - Frontend unit tests → `test: Add frontend component tests`
   - Integration tests → `test: Add API integration tests`
   - E2E test scenarios → `test: Add end-to-end test suite`

**Total commits for web app: ~21-24 atomic, well-documented commits**

### CLI Tool

1. **Research Phase** (3-4 subagents, ~4 commits):
   - Requirements analysis → `docs: Add requirements analysis`
   - CLI framework research → `docs: Add CLI framework research`
   - Architecture planning → `docs: Add architecture plan`
   - Error handling strategy → `docs: Add error handling strategy`

2. **Development Phase** (5-6 subagents, ~6-7 commits):
   - Project structure → `chore: Initialize CLI project structure`
   - Core logic implementation → `feat: Implement core processing logic`
   - CLI interface → `feat: Add command-line interface`
   - Configuration handling → `feat: Add configuration file support`
   - Help/documentation → `docs: Add CLI help and usage docs`
   - Integration → `feat: Integrate CLI with core logic`

3. **Testing Phase** (3-4 subagents, ~4 commits):
   - Test planning → `docs: Add test plan`
   - Unit tests → `test: Add unit tests for core logic`
   - CLI integration tests → `test: Add CLI integration tests`
   - Documentation validation → `docs: Validate and update documentation`

**Total commits for CLI tool: ~14-15 atomic, well-documented commits**

---

## Error Recovery

### When a Subagent Fails

1. **Analyze Failure**: Spawn diagnostic subagent to understand issue (no commit)
2. **Minimal Rollback**: Only rollback failed subagent's work
   ```bash
   git revert [commit-hash]  # If commit was made
   ```
3. **Context Adjustment**: Refine context and retry
4. **Successful Retry**: Commit with reference to retry
5. **Escalation Path**: After 2 failures, spawn coordinator subagent

### Commit Handling for Failed Subagents

- **Failed Subagent**: Do NOT commit anything
- **Diagnostic Subagent**: Do NOT commit (no code changes, analysis only)
- **Successful Retry**: Commit with reference to previous attempt:
  ```bash
  git commit -m "feat: Implement validation logic
  
  Second attempt after refining error handling requirements.
  Original attempt encountered issues with edge case handling.
  
  Subagent: Component Development (Validation) - Retry
  Phase: Development"
  ```
- **Rollback Required**: Use git to revert failed commit:
  ```bash
  # Revert specific commit
  git revert [commit-hash]
  
  # Or reset if nothing depends on it yet
  git reset --hard HEAD~1
  ```

### Coordinator Subagent

Use when:
- Multiple subagent failures occur
- Requirements seem contradictory
- Architecture needs major revision
- User clarification needed

Responsibilities:
- Analyze all previous attempts and commits
- Identify root cause
- Propose solution strategy
- Coordinate recovery plan
- **Does NOT commit** (analysis only)

After coordinator analysis, spawn new subagent to implement fix, which WILL commit.

---

## Best Practices

### DO ✅

**Process & Workflow:**
- **Always** complete research before coding
- **Always** use subagents for components > 100 lines
- **Always** provide structured handoff artifacts
- **Always** validate phase completion before proceeding
- **Always** separate concerns between subagents
- **Always** prefer multiple focused subagents over one large subagent

**Git & Commits:**
- **Always** read `GIT_COMMIT_STYLE_GUIDE.md` before committing
- **Always** commit your work before producing completion artifact
- **Always** use appropriate commit type for your subagent role
- **Always** include commit hash and details in completion artifact
- **Always** write descriptive commit messages with context
- **Always** review changes before committing (`git diff`)
- **Always** stage files explicitly (`git add`)
- **Always** include "Subagent:" and "Phase:" in commit body

**Context Management:**
- **Always** pass forward only essential information
- **Always** reference commits instead of duplicating code
- **Always** keep context under 4000 tokens per subagent

### DON'T ❌

**Process & Workflow:**
- **Never** mix research and implementation
- **Never** pass entire codebases between subagents
- **Never** skip the research phase
- **Never** test before development is complete
- **Never** use a single subagent for complex tasks
- **Never** exceed 4000 tokens of context per subagent

**Git & Commits:**
- **Never** commit with "WIP", "temp", or vague messages
- **Never** commit failed or incomplete work
- **Never** skip reading GIT_COMMIT_STYLE_GUIDE.md
- **Never** use past tense in commit messages ("Added" → "Add")
- **Never** commit without reviewing changes first
- **Never** commit during error recovery (wait for successful completion)
- **Never** mix multiple unrelated changes in one commit
- **Never** exceed 50 characters in commit subject line
- **Never** end subject line with a period

**Context Management:**
- **Never** duplicate code in handoff artifacts
- **Never** pass irrelevant context
- **Never** ignore token limits

---

## Subagent Invocation Template

```markdown
You are a specialized subagent responsible for: [SPECIFIC TASK]

## Your Role
[One sentence description]

## Pre-Task Requirements (MANDATORY)
Before starting your work, you MUST read these files:
1. **`GIT_COMMIT_STYLE_GUIDE.md`** - Learn commit message conventions
2. [Any other relevant guides or documents]

## Context Provided
- **Requirements**: [Reference or summary]
- **Previous Phase**: [Completion artifact from previous subagent]
- **Previous Commits**: [Relevant commit hashes and summaries]
- **Your Focus**: [Specific area of responsibility]

## Your Tasks
1. [Specific task 1]
2. [Specific task 2]
3. [Specific task 3]

## Expected Outputs
- **Files to Create/Modify**:
  - `path/to/file1.py` - [description]
  - `path/to/file2.md` - [description]
- **Artifacts**: [specific artifacts]

## Constraints
- Do NOT [specific prohibition]
- Focus ONLY on [specific scope]
- Assume [specific assumptions from research]
- Stay within [specific boundaries]

## Success Criteria
- [ ] [Measurable criterion 1]
- [ ] [Measurable criterion 2]
- [ ] [Measurable criterion 3]

## Completion Checklist
Before producing your completion artifact, ensure you have:
- [ ] Completed all assigned tasks
- [ ] Created/modified all expected output files
- [ ] Reviewed all changes (`git diff`)
- [ ] Read `GIT_COMMIT_STYLE_GUIDE.md`
- [ ] Determined appropriate commit type
- [ ] Written clear commit message (≤50 char subject)
- [ ] Staged changes (`git add`)
- [ ] Executed git commit
- [ ] Captured commit hash
- [ ] Verified all success criteria are met
- [ ] Prepared completion artifact with commit info

## Commit Type Guidance
Based on your role, your commit should likely be:
- **[Expected commit type]** - [Reasoning]

Example commit message:
```
[type]: [Short description of your specific work]

[What you did and why]

Subagent: [Your Role Name]
Phase: [Research/Development/Testing]
```

Please proceed with your focused task. When complete, commit your work and provide a structured completion artifact.
```

---

## Metrics for Success

Track these metrics to ensure effective subagent usage:

### Process Metrics
- **Phase Separation**: 100% of research complete before development starts
- **Subagent Focus**: Each subagent handles ≤1 major component
- **Context Efficiency**: No subagent receives >4000 tokens
- **Handoff Success**: 100% of handoffs include completion artifacts
- **Token Usage**: Total tokens ≤ sum of individual subagent tokens
- **Rework Rate**: <10% of subagents need to be re-run

### Git Metrics
- **Commit Rate**: 100% of successful subagents produce commits
- **Commit Quality**: 100% of commits follow style guide
- **Atomic Commits**: Each commit represents one logical change
- **Commit Coverage**: All code changes are committed before handoff
- **Message Quality**: Commit messages are clear and descriptive
- **Revert Rate**: <5% of commits need to be reverted

### Quality Metrics
- **Test Coverage**: >80% code coverage from testing phase
- **Documentation**: All components have documentation
- **Architecture Adherence**: 100% of code follows architecture plan
- **Success Rate**: >90% of subagents complete without retry

---

## Git History Benefits

With proper subagent commits, your git history becomes:

### 📊 A Clear Timeline
```bash
git log --oneline --graph

* abc1234 test: Add integration tests for auth flow
* def5678 test: Add unit tests for authentication module
* ghi9012 docs: Add comprehensive test plan
* jkl3456 docs: Add API documentation with examples
* mno7890 feat: Integrate auth module with API endpoints
* pqr1234 feat: Implement JWT authentication module
* stu5678 feat: Set up project structure with configuration
* vwx9012 docs: Add system architecture plan
* yza3456 docs: Add technical research for auth system
* bcd7890 docs: Add requirements analysis
```

### 🔍 Easy Debugging
```bash
# Find when a specific component was added
git log --grep="authentication"

# See what a specific subagent did
git log --grep="Subagent: Component Development (Auth)"

# View changes from development phase
git log --grep="Phase: Development"
```

### ⏮️ Safe Rollback
```bash
# Revert a specific subagent's work
git revert [commit-hash]

# Reset to before a problematic subagent
git reset --hard [commit-before-issue]

# Create a branch from a specific phase
git checkout -b feature-branch [last-research-commit]
```

### 📝 Automatic Documentation
```bash
# Generate changelog from commits
git log --pretty=format:"- %s" --grep="feat:"

# See all testing work
git log --pretty=format:"- %s (%h)" --grep="test:"

# Review all research
git log --pretty=format:"- %s" --grep="Phase: Research"
```

---

## Remember: The Power of Subagents

The effectiveness of subagents lies in:

1. **Focused Expertise**: Each subagent excels at ONE thing
2. **Context Efficiency**: Small, relevant contexts produce better outputs
3. **Process Discipline**: Strict ordering prevents wasted effort
4. **Clear Handoffs**: Structured artifacts ensure continuity
5. **Parallel Thinking**: Multiple specialized agents > one generalist
6. **Atomic Commits**: Each change is isolated, documented, and revertible
7. **Git History**: Timeline of work becomes a valuable project artifact

**USE SUBAGENTS EXTENSIVELY. COMMIT ATOMIC CHANGES. YOUR EFFECTIVENESS DEPENDS ON IT.**

---

## Quick Reference Card

### Subagent Workflow
```
1. Read GIT_COMMIT_STYLE_GUIDE.md
2. Receive context (≤4000 tokens)
3. Execute focused task
4. Review changes (git diff)
5. Determine commit type
6. Write commit message
7. Commit (git add + git commit)
8. Capture commit hash
9. Produce completion artifact
10. Hand off to next subagent
```

### Phase → Commit Type Mapping
- **Research Phase** → `docs:`
- **Development Phase** → `feat:`, `fix:`, `refactor:`
- **Testing Phase** → `test:`, `docs:`
- **Infrastructure** → `chore:`, `build:`, `ci:`
- **Performance** → `perf:`
- **Styling** → `style:`

### Commit Message Format
```
<type>: <description (≤50 chars)>

<body explaining what and why>

Subagent: [Name]
Phase: [Phase]
```

### Emergency Commands
```bash
# Revert last commit
git revert HEAD

# Reset (dangerous - use carefully)
git reset --hard HEAD~1

# See what changed
git diff HEAD~1

# View commit details
git show [commit-hash]
```

---

*This guide ensures that subagent-driven development produces not only high-quality code, but also a high-quality, navigable git history that serves as living documentation of the development process.*