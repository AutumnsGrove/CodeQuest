# CodeQuest Installation Guide

Complete installation instructions for CodeQuest, the terminal-based gamified developer productivity RPG.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Install](#quick-install)
3. [Step-by-Step Installation](#step-by-step-installation)
4. [Configuration Setup](#configuration-setup)
5. [Verification](#verification)
6. [Troubleshooting](#troubleshooting)
7. [Uninstallation](#uninstallation)
8. [Getting Help](#getting-help)

---

## Prerequisites

Before installing CodeQuest, ensure you have the following installed:

### Required

- **Go 1.21 or higher**: [Install Go](https://go.dev/doc/install)
  - Verify installation:
    ```bash
    go version
    ```

- **Git**: [Install Git](https://git-scm.com/downloads)
  - Verify installation:
    ```bash
    git --version
    ```

- **Skate**: Data persistence for character and quest data
  - Install via Homebrew (macOS/Linux):
    ```bash
    brew install charmbracelet/tap/skate
    ```
  - Verify installation:
    ```bash
    skate --version
    ```

### Optional (for AI Mentor Features)

- **Mods** (local AI models): Install via Homebrew
  ```bash
  brew install charmbracelet/tap/mods
  ```
  - Requires Ollama or similar for running local models

- **API Keys** for online AI:
  - OpenRouter (Crush): [openrouter.ai/keys](https://openrouter.ai/keys)
  - Anthropic (Claude): [console.anthropic.com](https://console.anthropic.com)

---

## Quick Install

For most users, the fastest way to install CodeQuest:

```bash
# Clone the repository
git clone https://github.com/AutumnsGrove/codequest.git
cd codequest

# Build the application
make build

# Run CodeQuest
./build/codequest
```

That's it! CodeQuest will create its configuration on first run.

**Alternative:** Install via `go install` (latest version):
```bash
go install github.com/AutumnsGrove/codequest/cmd/codequest@latest
```

---

## Step-by-Step Installation

### 1. Check Prerequisites

Before starting, verify all required tools are installed:

```bash
# Check Go
go version
# Output should be: go version go1.21+ ...

# Check Git
git --version
# Output should be: git version 2.x.x ...

# Check Skate
skate --version
# Output should show skate version
```

If any of these fail, install the missing tools using the links in the [Prerequisites](#prerequisites) section.

### 2. Clone the Repository

```bash
# Navigate to where you want to store the project
cd ~/projects  # or your preferred directory

# Clone the CodeQuest repository
git clone https://github.com/AutumnsGrove/codequest.git
cd codequest

# Verify you're in the correct directory
ls -la
# Should show: Makefile, README.md, cmd/, internal/, etc.
```

### 3. Install Dependencies

Download and prepare the Go dependencies:

```bash
# Install Go module dependencies
make deps

# Verify dependencies are installed
go mod verify
# Output: Should show "all modules verified" or similar
```

### 4. Build the Application

Compile CodeQuest into an executable:

```bash
# Build for your system (macOS/Linux)
make build

# Verify the build succeeded
ls -la ./build/
# Should show: codequest (executable)

# Test the binary
./build/codequest --version
# Should show version information
```

### 5. Make It Available Globally (Optional)

Install the binary to your system PATH so you can run `codequest` from anywhere:

```bash
# Install globally
make install

# Verify installation
which codequest
# Should show: /usr/local/bin/codequest (or similar path)

# Test global installation
codequest --version
```

Or manually:
```bash
# Copy the binary to a directory in your PATH
cp ./build/codequest /usr/local/bin/
chmod +x /usr/local/bin/codequest

# Verify
codequest --version
```

### 6. Launch CodeQuest

```bash
# If installed globally
codequest

# Or if running from the build directory
./build/codequest
```

On first run, CodeQuest will:
1. Create `~/.config/codequest/` directory
2. Generate a default `config.toml`
3. Prompt you to create a character
4. Display the main dashboard

---

## Configuration Setup

### First Run Configuration

When you first launch CodeQuest, you'll be guided through:

1. **Character Creation**: Enter your character name
2. **Game Settings**: Choose difficulty (easy/normal/hard)
3. **Repository Paths**: Specify where your Git projects are located

### Configuration File Location

CodeQuest stores its configuration at:
```
~/.config/codequest/config.toml
```

Edit this file directly to customize settings:

```bash
# Open configuration in your preferred editor
nano ~/.config/codequest/config.toml
# or
vim ~/.config/codequest/config.toml
# or
open ~/.config/codequest/config.toml  # macOS
```

### Basic Configuration Example

```toml
[character]
name = "YourCharacterName"
starting_level = 1

[game]
difficulty = "normal"  # easy, normal, or hard
auto_detect_repos = true

[git]
watch_paths = [
  "~/projects",
  "~/work"
]

[ai.mentor]
provider = "crush"  # or "mods", "claude"
model_complex = "openrouter/deepseek/glm-4.5-air"
model_simple = "openrouter/deepseek/glm-4.5-air"
temperature = 0.7

[ai.review]
provider = "mods"
auto_review = true
bonus_xp_enabled = true
```

### Configuring AI Providers

CodeQuest uses a fallback system: **Crush (OpenRouter) → Mods (Local) → Claude (Anthropic)**

#### Option A: Use Crush (OpenRouter) - Recommended for Beginners

Crush provides online AI models with no local setup required.

1. **Get an API key**:
   - Visit [openrouter.ai/keys](https://openrouter.ai/keys)
   - Create a free account
   - Copy your API key

2. **Store the key securely**:
   ```bash
   # Save the key to Skate (encrypted)
   skate set codequest.openrouter_api_key "sk-or-v1-..."

   # Verify it was stored
   skate get codequest.openrouter_api_key
   # Output: sk-or-v1-...
   ```

3. **Update your config**:
   ```toml
   [ai.mentor]
   provider = "crush"
   model_complex = "openrouter/deepseek/glm-4.5-air"
   model_simple = "openrouter/deepseek/glm-4.5-air"
   ```

#### Option B: Use Mods (Local Models) - Privacy-Focused

Run AI models locally without external API calls.

1. **Install Mods**:
   ```bash
   brew install charmbracelet/tap/mods
   ```

2. **Configure local models**:
   ```bash
   # Open Mods settings
   mods --settings

   # Or manually check Mods documentation for your preferred model
   ```

3. **Update your config**:
   ```toml
   [ai.mentor]
   provider = "mods"
   model_complex_offline = "qwen3:30b"
   model_simple_offline = "qwen3:4b"
   ```

#### Option C: Use Claude (Anthropic) - Advanced

For access to Claude models.

1. **Get an API key**:
   - Visit [console.anthropic.com](https://console.anthropic.com)
   - Create account or sign in
   - Go to API keys section
   - Create a new key

2. **Store the key securely**:
   ```bash
   skate set codequest.anthropic_api_key "sk-ant-..."
   ```

3. **Update your config**:
   ```toml
   [ai.mentor]
   provider = "claude"
   model_complex = "claude-sonnet-4-5-20250929"
   model_simple = "claude-haiku-4-5-20251001"
   ```

### Configuring Git Repository Paths

Add the directories where your Git projects are located:

```bash
# Find your project directories
cd ~
find . -name ".git" -type d 2>/dev/null | head -5
```

Update your config:

```toml
[git]
auto_detect_repos = true
watch_paths = [
  "~/projects",
  "~/work/client-projects",
  "~/open-source"
]
```

CodeQuest will monitor these directories for Git commits and award XP automatically.

---

## Verification

### Step 1: Verify Installation

```bash
# Check the binary exists and runs
codequest --version

# Expected output:
# CodeQuest v0.1.0-beta
```

### Step 2: Verify Configuration

```bash
# Check configuration file exists
ls -la ~/.config/codequest/config.toml

# Expected output:
# -rw-r--r--  1 user  group  2048 Oct 25 12:00 ~/.config/codequest/config.toml
```

### Step 3: Verify Data Storage

```bash
# Check Skate is working
skate get codequest.test_key

# Skate might return an empty value (expected) or error if not configured
```

### Step 4: Run a Test

1. **Launch CodeQuest**:
   ```bash
   codequest
   ```

2. **In another terminal, test Git activity**:
   ```bash
   cd ~/projects  # or any of your watch_paths
   git status
   echo "# Test" >> test.txt
   git add test.txt
   git commit -m "test: Verify CodeQuest installation"
   ```

3. **Return to CodeQuest and verify**:
   - You should see XP awarded for the commit
   - Character stats should update
   - Commits appear in quest progress

### Step 5: Verify AI (Optional)

Test AI mentor functionality:

```bash
# In CodeQuest, press 'm' to open the Mentor screen
# Ask a question like "How do I use git branches?"
# You should receive a response from the AI provider
```

---

## Troubleshooting

### Common Issues and Solutions

#### "command not found: codequest"

**Problem**: The codequest command is not in your PATH.

**Solution**:
```bash
# Reinstall globally
make install

# Or add the build directory to PATH
export PATH="$PATH:~/projects/codequest/build"

# Make permanent (add to ~/.bashrc, ~/.zshrc, etc.)
echo 'export PATH="$PATH:~/projects/codequest/build"' >> ~/.zshrc
```

#### "skate: command not found"

**Problem**: Skate is not installed.

**Solution**:
```bash
# Install Skate
brew install charmbracelet/tap/skate

# Verify installation
skate --version

# If Homebrew doesn't find it
brew tap charmbracelet/tap
brew install skate
```

#### "Failed to load configuration"

**Problem**: Config file is missing or corrupted.

**Solution**:
```bash
# Create the config directory
mkdir -p ~/.config/codequest

# Delete the bad config (backup first!)
cp ~/.config/codequest/config.toml ~/.config/codequest/config.toml.backup
rm ~/.config/codequest/config.toml

# Launch CodeQuest to regenerate it
codequest
```

#### "No AI providers available"

**Problem**: No AI provider is properly configured.

**Solution - Option 1: Disable AI (CodeQuest works fine without it)**:
```toml
[ai.mentor]
provider = "none"
```

**Solution - Option 2: Set up OpenRouter (easiest)**:
```bash
# Get key from openrouter.ai/keys
skate set codequest.openrouter_api_key "sk-or-v1-..."

# Update config
# provider = "crush"
```

**Solution - Option 3: Set up local models**:
```bash
brew install charmbracelet/tap/mods
# Then configure provider = "mods" in config
```

#### "Build fails with missing go.mod"

**Problem**: Go modules aren't initialized.

**Solution**:
```bash
cd ~/projects/codequest
make deps
make build
```

#### "Git commits not detected"

**Problem**: CodeQuest isn't watching your repository.

**Solution**:
```bash
# Check your config
cat ~/.config/codequest/config.toml

# Verify your repository path is listed
# Make sure it's an absolute path or use ~/

# Add the path if missing:
# watch_paths = ["~/projects"]

# Verify the repository has commits
cd ~/your/project
git log --oneline -1
```

#### "API key not found in Skate"

**Problem**: API key wasn't stored properly.

**Solution**:
```bash
# Store the key again
skate set codequest.openrouter_api_key "YOUR_KEY_HERE"

# Verify it was stored
skate get codequest.openrouter_api_key

# Check the config has the correct provider
cat ~/.config/codequest/config.toml
```

#### "Character/progress data lost"

**Problem**: Data wasn't persisted properly.

**Solution**:
```bash
# Check Skate storage
skate list | grep codequest

# Verify Skate is working
skate set test.value "test"
skate get test.value

# If Skate is broken, reinstall it
brew reinstall charmbracelet/tap/skate
```

#### Build Issues on Different Platforms

**For macOS (Intel)**:
```bash
make build
```

**For macOS (Apple Silicon)**:
```bash
make build
# Or specify architecture
GOARCH=arm64 make build
```

**For Linux**:
```bash
make build
# Adjust build command if needed
```

**For Windows**:
```bash
# If using WSL (Windows Subsystem for Linux)
make build

# Or use Go directly
go build -o build/codequest.exe ./cmd/codequest
```

---

## Uninstallation

### Complete Uninstallation

If you want to completely remove CodeQuest:

```bash
# Remove the binary
rm /usr/local/bin/codequest
# Or wherever make install placed it

# Remove the source code
rm -rf ~/projects/codequest
# Or wherever you cloned it

# Remove configuration (OPTIONAL - keeps your character/data)
rm -rf ~/.config/codequest

# Remove stored secrets (OPTIONAL - removes API keys)
skate delete codequest.openrouter_api_key
skate delete codequest.anthropic_api_key
```

### Keep Data, Remove Application

To keep your character and progress:

```bash
# Remove only the binary and source
rm /usr/local/bin/codequest
rm -rf ~/projects/codequest

# Your data remains in:
# - ~/.config/codequest/config.toml
# - Skate's encrypted storage
```

Reinstalling later will restore your character!

---

## Getting Help

### Documentation

- **[README.md](README.md)** - Overview, features, and quick start
- **[CLAUDE.md](CLAUDE.md)** - Developer documentation and Go standards
- **[CODEQUEST_SPEC.md](CODEQUEST_SPEC.md)** - Complete technical specification
- **[DEVELOPMENT_STATUS.md](DEVELOPMENT_STATUS.md)** - Current development progress

### Support Resources

- **Issues**: [GitHub Issues](https://github.com/AutumnsGrove/codequest/issues)
- **Discussions**: [GitHub Discussions](https://github.com/AutumnsGrove/codequest/discussions)
- **Questions**: Check existing issues before opening a new one

### Reporting Bugs

When reporting an issue, include:

```bash
# Your system information
go version
git --version
skate --version

# Relevant config (without API keys)
cat ~/.config/codequest/config.toml

# Error messages (copy/paste exactly)

# Steps to reproduce
```

### Getting Support for Specific Providers

**OpenRouter Issues**:
- Visit [openrouter.ai/status](https://openrouter.ai/status)
- Check your account balance
- Review API key permissions

**Mods Issues**:
- Run `mods --settings` to verify configuration
- Check if Ollama/models are installed
- Review Mods documentation

**Skate Issues**:
- Run `skate --version` to verify installation
- Check file permissions: `ls -la ~/.local/share/skate/` (or appropriate path)
- Review Skate documentation

---

## Next Steps

Once installed and configured:

1. **Create your character** - Launch `codequest` and create your hero
2. **Set up Git watching** - Configure `watch_paths` for your repositories
3. **Accept a quest** - Browse the Quest Board (press `q`)
4. **Start coding** - Make commits and earn XP
5. **Level up** - Watch your character grow with your productivity

Happy questing! 🎮

---

**Questions or issues?** Check the [Getting Help](#getting-help) section or visit the [GitHub repository](https://github.com/AutumnsGrove/codequest).

For developers, see [CLAUDE.md](CLAUDE.md) for contribution guidelines and development setup.
