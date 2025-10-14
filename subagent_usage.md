---
aliases: 
date created: Monday, October 6th 2025, 7:18:23 am
date modified: Monday, October 13th 2025, 12:26:44 am
tags: 
type:
---

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

2. **Technical Research Subagent**
   - Research best practices for the technology stack
   - Identify potential libraries/tools
   - Research common pitfalls and solutions
   - Output: `technical_research.md`

3. **Architecture Planning Subagent**
   - Design high-level architecture
   - Define component boundaries
   - Plan data flow and interfaces
   - Output: `architecture_plan.md`

4. **Edge Case Analysis Subagent**
   - Identify potential edge cases
   - Research error handling strategies
   - Define validation requirements
   - Output: `edge_cases.md`

### Research Phase Completion Gate

✅ All research subagents have completed
✅ All research artifacts are created
✅ No unresolved questions remain
✅ Architecture is fully planned

**Only then proceed to Development Phase**

---

## Phase 2: Development (AFTER Research Complete)

### Development Subagent Strategy

1. **Component Isolation**: Each major component gets its own subagent
2. **Sequential Building**: Build dependencies before dependents
3. **Incremental Integration**: Integrate components one at a time

### Required Development Subagents

1. **Core Infrastructure Subagent**
   - Set up project structure
   - Create base configurations
   - Implement shared utilities
   - Context: Research artifacts + requirements

2. **Component Development Subagents** (one per component)
   - Implement single component/module
   - Follow architecture plan exactly
   - Include inline documentation
   - Context: Component-specific requirements + interfaces

3. **Integration Subagent**
   - Connect components
   - Implement data flow
   - Handle inter-component communication
   - Context: All component outputs + architecture

4. **Documentation Subagent**
   - Create user documentation
   - Write API documentation
   - Create usage examples
   - Context: Completed codebase summary

### Development Phase Completion Gate

✅ All components implemented
✅ All integrations complete
✅ Documentation created
✅ Code is functional (not tested)

**Only then proceed to Testing Phase**

---

## Phase 3: Testing (AFTER Development Complete)

### Testing Subagent Strategy

1. **Test Planning Subagent**
   - Create comprehensive test plan
   - Define test cases for each component
   - Plan integration test scenarios
   - Output: `test_plan.md`

2. **Unit Test Development Subagent**
   - Write unit tests per test plan
   - One subagent per component's tests
   - Include edge case coverage
   - Context: Component code + test plan

3. **Integration Test Subagent**
   - Write integration tests
   - Test component interactions
   - Validate data flow
   - Context: Integration points + test plan

4. **Test Execution Subagent**
   - Run all tests
   - Document results
   - Identify any failures
   - Output: `test_results.md`

---

## Context Management Rules

### 1. Context Inheritance Hierarchy

```
Project Requirements (base context)
    ↓
Research Artifacts
    ↓
Architecture Plan
    ↓
Component Interfaces
    ↓
Implementation Details
```

### 2. Context Pruning Strategy

- **Pass Forward**: Only essential outputs, not full code
- **Summarize**: Previous phase outputs into bullet points
- **Reference**: Use file references instead of inline content
- **Focus**: Each subagent gets ONLY relevant context

### 3. Token Optimization

Each subagent should receive:
- Core requirements (< 500 tokens)
- Phase-specific context (< 2000 tokens)
- Relevant previous outputs (< 1000 tokens)
- Clear instructions (< 500 tokens)

**Total context per subagent: ~4000 tokens maximum**

---

## Handoff Protocol

### Subagent Output Format

Every subagent MUST produce:

```markdown
## Completion Artifact: [Subagent Name]

### Summary
[2-3 sentences of what was accomplished]

### Key Outputs
- [List of created files/artifacts]
- [Important decisions made]
- [Critical findings]

### Next Steps
- [Recommended next subagent]
- [Required inputs for next phase]

### Concerns/Blockers
- [Any issues encountered]
- [Risks identified]

### Confidence Score
[High/Medium/Low] - [Explanation]
```

### Handoff Rules

1. **Explicit Handoff**: Previous subagent explicitly states completion
2. **Artifact Validation**: Verify all expected outputs exist
3. **Context Transfer**: Next subagent receives completion artifact + relevant context
4. **No Backtracking**: Cannot return to previous phase

---

## Common Task Decompositions

### Web Application

1. **Research Phase** (5-6 subagents):
   - Requirements analysis
   - Frontend framework research
   - Backend architecture research
   - Database design research
   - Security requirements research
   - Deployment strategy research

2. **Development Phase** (8-10 subagents):
   - Project setup
   - Database schema
   - Backend API (one per major endpoint group)
   - Frontend components (one per major feature)
   - Styling/UI
   - Integration
   - Configuration

3. **Testing Phase** (4-5 subagents):
   - Test planning
   - Backend unit tests
   - Frontend unit tests
   - Integration tests
   - E2E test scenarios

### CLI Tool

1. **Research Phase** (3-4 subagents):
   - Requirements analysis
   - CLI framework research
   - Architecture planning
   - Error handling strategy

2. **Development Phase** (5-6 subagents):
   - Project structure
   - Core logic implementation
   - CLI interface
   - Configuration handling
   - Help/documentation
   - Integration

3. **Testing Phase** (3-4 subagents):
   - Test planning
   - Unit tests
   - CLI integration tests
   - Documentation validation

---

## Error Recovery

### When a Subagent Fails

1. **Analyze Failure**: Spawn diagnostic subagent to understand issue
2. **Minimal Rollback**: Only rollback failed subagent's work
3. **Context Adjustment**: Refine context and retry
4. **Escalation Path**: After 2 failures, spawn coordinator subagent

### Coordinator Subagent

Use when:
- Multiple subagent failures occur
- Requirements seem contradictory
- Architecture needs major revision
- User clarification needed

Responsibilities:
- Analyze all previous attempts
- Identify root cause
- Propose solution strategy
- Coordinate recovery plan

---

## Best Practices

### DO ✅

- **Always** complete research before coding
- **Always** use subagents for components > 100 lines
- **Always** provide structured handoff artifacts
- **Always** validate phase completion before proceeding
- **Always** separate concerns between subagents
- **Always** prefer multiple focused subagents over one large subagent

### DON'T ❌

- **Never** mix research and implementation
- **Never** pass entire codebases between subagents
- **Never** skip the research phase
- **Never** test before development is complete
- **Never** use a single subagent for complex tasks
- **Never** exceed 4000 tokens of context per subagent

---

## Subagent Invocation Template

```markdown
You are a specialized subagent responsible for: [SPECIFIC TASK]

## Your Role
[One sentence description]

## Context Provided
- Requirements: [Reference or summary]
- Previous Phase: [Completion artifact from previous subagent]
- Your Focus: [Specific area of responsibility]

## Your Tasks
1. [Specific task 1]
2. [Specific task 2]
3. [Specific task 3]

## Expected Outputs
- [Specific file/artifact 1]
- [Specific file/artifact 2]

## Constraints
- Do NOT [specific prohibition]
- Focus ONLY on [specific scope]
- Assume [specific assumptions from research]

## Success Criteria
- [Measurable criterion 1]
- [Measurable criterion 2]

Please proceed with your focused task and provide a completion artifact when done.
```

---

## Metrics for Success

Track these metrics to ensure effective subagent usage:

- **Phase Separation**: 100% of research complete before development starts
- **Subagent Focus**: Each subagent handles ≤1 major component
- **Context Efficiency**: No subagent receives >4000 tokens
- **Handoff Success**: 100% of handoffs include completion artifacts
- **Token Usage**: Total tokens ≤ sum of individual subagent tokens
- **Rework Rate**: <10% of subagents need to be re-run

---

## Remember

The power of subagents lies in:
1. **Focused Expertise**: Each subagent excels at ONE thing
2. **Context Efficiency**: Small, relevant contexts produce better outputs
3. **Process Discipline**: Strict ordering prevents wasted effort
4. **Clear Handoffs**: Structured artifacts ensure continuity
5. **Parallel Thinking**: Multiple specialized agents > one generalist

**USE SUBAGENTS EXTENSIVELY. YOUR EFFECTIVENESS DEPENDS ON IT.**
