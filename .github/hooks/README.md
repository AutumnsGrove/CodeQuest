# Git Pre-Commit Hooks for CodeQuest

This directory contains Git hooks to maintain code quality and enforce standards in the CodeQuest project.

## Quick Start

### Install Hooks

Copy the hook files to your local `.git/hooks/` directory and make them executable:

```bash
# From project root
cp .github/hooks/pre-commit .git/hooks/
cp .github/hooks/commit-msg .git/hooks/
cp .github/hooks/pre-push .git/hooks/

# Make executable
bash -c 'chmod +x .git/hooks/{pre-commit,commit-msg,pre-push}'

# Verify installation
ls -la .git/hooks/pre-* .git/hooks/commit-msg
```

### Test Installation

```bash
# Test commit-msg hook with a sample message
echo "feat: test message" | .git/hooks/commit-msg /tmp/test_msg.txt

# Test pre-commit hook
.git/hooks/pre-commit

# Test pre-push hook (requires go to be installed)
.git/hooks/pre-push
```

## Hook Descriptions

### 1. pre-commit Hook

**When it runs:** Before each commit is created

**What it checks:**
- **Secrets Scanning**: Detects API keys and sensitive data before committing
  - Anthropic API keys (`sk-ant-api...`)
  - OpenAI/OpenRouter keys (`sk-...`)
  - Google API keys
  - AWS access keys
  - `secrets.json` files
  - `credentials.json` files
  - `.env` files

- **Code Formatting**: Ensures all Go code follows project standards
  - Runs `gofmt` on staged Go files
  - Identifies files that need formatting

- **Static Analysis**: Runs Go's built-in code analyzer
  - Executes `go vet ./...`
  - Catches common programming errors

- **Test Execution**: Runs tests for any modified packages
  - Runs tests only for packages with changes
  - Ensures no regressions are introduced

**Exit behavior:**
- Exits with code 0 (success) if all checks pass
- Exits with code 1 (failure) if any check fails - commit is blocked
- Shows clear error messages and fixes for each issue

**Files checked:**
- Only files staged with `git add` are checked
- Untracked or unstaged files are ignored

### 2. commit-msg Hook

**When it runs:** After writing commit message but before finalizing commit

**What it validates:**
- **Message Format**: Ensures commit messages follow project standards

Supports two formats:

**Format 1: Custom Action-Based (Recommended)**
```
[Action] Brief description
Valid actions: Add, Update, Fix, Refactor, Remove, Enhance, Improve, Create, Delete, Rename
Example: Add character leveling system
```

**Format 2: Conventional Commits**
```
type: description  OR  type(scope): description
Valid types: feat, fix, docs, style, refactor, test, chore, perf, ci, build, revert
Example: feat: add character leveling system
```

- **Message Length**: Warns if first line exceeds 72 characters
- **Co-Authored Attribution**: Recommends adding `Co-Authored-By: Claude <noreply@anthropic.com>`

**Exit behavior:**
- Exits with code 0 (success) if message format is valid
- Exits with code 1 (failure) if format is invalid - commit is blocked
- Shows examples of valid formats

**Why this matters:**
Good commit messages are essential for:
- Maintaining readable project history
- Enabling proper git operations (bisect, blame, log)
- Documenting design decisions
- Creating clear release notes

### 3. pre-push Hook

**When it runs:** Before pushing to remote repository

**What it checks:**
- **Full Test Suite**: Runs `go test ./...` with verbose output
- **Code Analysis**: Runs `go vet ./...` on entire project
- **Code Formatting**: Checks `gofmt` compliance on all files
- **Test Coverage**: Reports average coverage percentage
  - Warns if coverage drops below 70%
  - Informational only (doesn't block push)

**Exit behavior:**
- Exits with code 0 (success) if all checks pass - push proceeds
- Exits with code 1 (failure) if any critical check fails - push is blocked
- Test coverage below 70% is a warning only

**Why this matters:**
- Ensures remote repository only receives stable, tested code
- Maintains project health by catching issues before they're shared
- Prevents failed CI/CD builds from incomplete work
- Verifies code quality standards across entire project

## Hook Execution Order

When you run `git commit`:

```
1. pre-commit runs
   ├─ Scans for secrets
   ├─ Checks Go formatting
   ├─ Runs go vet analysis
   └─ Runs tests on modified packages

2. (if pre-commit passes) commit-msg runs
   └─ Validates message format

3. (if both pass) Commit is created
```

When you run `git push`:

```
1. pre-push runs
   ├─ Runs full test suite
   ├─ Runs code analysis
   ├─ Checks formatting
   └─ Reports test coverage

2. (if pre-push passes) Push proceeds
```

## Troubleshooting

### Hook Not Running

**Problem:** Hook script not executing during commit/push

**Solution:**
```bash
# Verify hook exists
ls -la .git/hooks/pre-commit

# Verify it's executable (should show -rwxr-xr-x)
ls -la .git/hooks/pre-commit

# Make executable if needed
bash -c 'chmod +x .git/hooks/pre-commit'

# Test running manually
.git/hooks/pre-commit
```

### "Permission Denied" Error

**Problem:** Hooks exist but fail with permission error

**Solution:**
```bash
# Fix permissions for all hooks
bash -c 'chmod +x .git/hooks/{pre-commit,commit-msg,pre-push}'

# Verify
ls -la .git/hooks/pre-* .git/hooks/commit-msg
```

### Pre-commit Hook Fails on Formatting

**Problem:** Hook fails because code isn't formatted

**Solution:**
```bash
# Run gofmt to automatically format code
gofmt -w ./...

# Or format specific file
gofmt -w path/to/file.go

# Stage the formatted files
git add path/to/file.go

# Try commit again
git commit -m "Your message"
```

### Pre-commit Hook Fails go vet

**Problem:** `go vet` finds issues like unused variables, logic errors

**Solution:**
```bash
# See detailed go vet output
go vet ./...

# Fix issues manually in your code, then try again
git add your/file.go
git commit -m "Your message"
```

### Pre-commit Hook Fails Tests

**Problem:** Tests fail before commit is allowed

**Solution:**
```bash
# Run tests to see failures
go test ./... -v

# Fix failing tests in your code

# Re-run tests to verify
go test ./... -v

# Stage your fixes and try commit again
git add your/file.go
git commit -m "Your message"
```

### Pre-push Hook Fails on Coverage

**Problem:** Coverage below 70% warning appears

**Note:** This is informational and doesn't block the push. It's a reminder to add tests.

**Solution:**
```bash
# Check which packages have low coverage
go test ./... -cover

# Add tests to low-coverage packages
# Re-run to verify coverage improved
go test ./... -cover
```

### Pre-commit Detects Secrets

**Problem:** Hook blocks commit due to suspected API key

**Solution:**
```bash
# Check what was detected
git diff --cached | grep -E "sk-ant-api|sk-[a-zA-Z0-9]{32,}"

# Remove or move the sensitive data
# Example if you accidentally added a secrets file:
git reset HEAD secrets.json
# Add secrets.json to .gitignore if not already there
echo "secrets.json" >> .gitignore

# Stage the fixed files
git add .gitignore

# Try commit again
git commit -m "Your message"
```

## Bypassing Hooks (Use Sparingly)

If you need to bypass hooks for legitimate reasons:

```bash
# Skip pre-commit and commit-msg checks
git commit --no-verify -m "Your message"

# Skip pre-push check
git push --no-verify

# Push to specific branch without pre-push checks
git push origin main --no-verify
```

**When to use `--no-verify`:**
- Emergency production hotfixes
- Hook is temporarily broken/misconfigured
- Intentional policy override approved by team

**When NOT to use `--no-verify`:**
- Regular development work
- Avoiding test failures
- Ignoring code quality issues

## Go Environment Requirements

The hooks require Go to be installed and properly configured:

```bash
# Verify Go is installed
go version

# Verify GOPATH is set (should return non-empty path)
echo $GOPATH

# Verify your project has a go.mod file
cat go.mod
```

If Go tools are not found:
- Install Go: https://go.dev/doc/install
- Add Go bin to PATH (usually automatic with macOS installation)
- Restart your terminal and verify with `go version`

## Customizing Hooks

### Disabling a Hook Temporarily

Rename the hook to disable it:

```bash
# Disable pre-commit hook
mv .git/hooks/pre-commit .git/hooks/pre-commit.disabled

# Re-enable it
mv .git/hooks/pre-commit.disabled .git/hooks/pre-commit
```

### Modifying Hook Behavior

Edit `.github/hooks/*` files and copy to `.git/hooks/`:

```bash
# Edit the hook
nano .github/hooks/pre-commit

# Copy to active hooks directory
cp .github/hooks/pre-commit .git/hooks/

# Make executable
bash -c 'chmod +x .git/hooks/pre-commit'
```

## Common Hook Patterns

### Running Only on Specific Files

```bash
# Modify pre-commit hook to filter files
go_files=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' | grep -E 'internal/game|cmd' || true)
```

### Adding Additional Checks

```bash
# Example: Check for println statements (debug logging)
if git diff --cached | grep -E "fmt\.Println|log\.Println"; then
    echo "Debug print statements found. Remove before committing."
    exit 1
fi
```

## Related Documentation

- **CLAUDE.md** - Development guide with Go standards
- **GIT_COMMIT_STYLE_GUIDE.md** - Detailed commit message standards
- **DEVELOPMENT_STATUS.md** - Current project status and roadmap
- [Effective Go](https://go.dev/doc/effective_go) - Go code style guide
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Support and Questions

For issues with hooks:

1. Check the Troubleshooting section above
2. Verify Go installation with `go version`
3. Run hooks manually to see detailed output
4. Check hook shell syntax: `bash -n .git/hooks/pre-commit`
5. See CLAUDE.md for project-specific development guidelines

## Installation Script

If manual installation is tedious, you can create a setup script:

```bash
#!/bin/bash
# setup_hooks.sh

echo "Installing CodeQuest Git hooks..."

# Copy hooks
cp .github/hooks/pre-commit .git/hooks/
cp .github/hooks/commit-msg .git/hooks/
cp .github/hooks/pre-push .git/hooks/

# Make executable
bash -c 'chmod +x .git/hooks/{pre-commit,commit-msg,pre-push}'

# Verify
if [ -x .git/hooks/pre-commit ]; then
    echo "✓ Hooks installed successfully"
    ls -la .git/hooks/pre-* .git/hooks/commit-msg
else
    echo "✗ Hook installation failed"
    exit 1
fi
```

Save as `setup_hooks.sh`, make executable, and run:

```bash
bash setup_hooks.sh
```

---

**Last Updated:** October 25, 2025

For CodeQuest project standards, see CLAUDE.md
