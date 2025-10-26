#!/bin/bash
# Installation script for CodeQuest Git hooks
# Makes it easy to set up all pre-commit hooks in one command

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  CodeQuest Git Hooks Installation     ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Find project root
PROJECT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)

if [ -z "$PROJECT_ROOT" ]; then
    echo -e "${RED}✗ Not in a Git repository${NC}"
    echo "Please run this script from within your CodeQuest Git repository."
    exit 1
fi

HOOKS_SOURCE="$PROJECT_ROOT/.github/hooks"
HOOKS_DEST="$PROJECT_ROOT/.git/hooks"

# Verify source directory exists
if [ ! -d "$HOOKS_SOURCE" ]; then
    echo -e "${RED}✗ Hooks source directory not found: $HOOKS_SOURCE${NC}"
    exit 1
fi

echo -e "${YELLOW}Project root: $PROJECT_ROOT${NC}"
echo ""

# List of hooks to install
HOOKS=("pre-commit" "commit-msg" "pre-push")

echo -e "${YELLOW}Installing hooks:${NC}"
for hook in "${HOOKS[@]}"; do
    if [ ! -f "$HOOKS_SOURCE/$hook" ]; then
        echo -e "${RED}✗ Hook file not found: $HOOKS_SOURCE/$hook${NC}"
        exit 1
    fi

    echo -e "  ${YELLOW}→${NC} Installing $hook..."

    # Copy hook
    cp "$HOOKS_SOURCE/$hook" "$HOOKS_DEST/$hook"

    # Make executable
    bash -c "chmod +x '$HOOKS_DEST/$hook'" || {
        echo -e "${RED}✗ Failed to make $hook executable${NC}"
        exit 1
    }

    echo -e "    ${GREEN}✓${NC} $hook installed"
done

echo ""
echo -e "${YELLOW}Verifying installation...${NC}"
echo ""

all_installed=true
for hook in "${HOOKS[@]}"; do
    hook_path="$HOOKS_DEST/$hook"

    if [ ! -f "$hook_path" ]; then
        echo -e "  ${RED}✗${NC} $hook: Not found"
        all_installed=false
        continue
    fi

    if [ ! -x "$hook_path" ]; then
        echo -e "  ${RED}✗${NC} $hook: Not executable"
        all_installed=false
        continue
    fi

    echo -e "  ${GREEN}✓${NC} $hook: Installed and executable"
done

echo ""

if [ "$all_installed" = true ]; then
    echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║  Installation Successful! ✓             ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${YELLOW}What's installed:${NC}"
    echo "  • pre-commit  - Checks code formatting, runs tests, scans for secrets"
    echo "  • commit-msg  - Validates commit message format"
    echo "  • pre-push    - Runs full test suite before pushing"
    echo ""
    echo -e "${YELLOW}Next steps:${NC}"
    echo "  1. Try making a test commit to verify hooks work:"
    echo "     git add . && git commit -m 'test: verify hooks'"
    echo ""
    echo "  2. For more information, see:"
    echo "     .github/hooks/README.md"
    echo ""
    echo -e "${YELLOW}To temporarily skip hooks (use sparingly):${NC}"
    echo "  git commit --no-verify -m 'Your message'"
    echo ""
    exit 0
else
    echo -e "${RED}╔════════════════════════════════════════╗${NC}"
    echo -e "${RED}║  Installation Failed ✗                 ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${YELLOW}Troubleshooting:${NC}"
    echo "  1. Check Git repository exists: git status"
    echo "  2. Check .github/hooks/ directory exists"
    echo "  3. Check file permissions: ls -la .github/hooks/"
    echo "  4. Try manual installation: cp .github/hooks/* .git/hooks/"
    echo ""
    exit 1
fi
