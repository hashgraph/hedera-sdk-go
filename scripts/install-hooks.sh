#!/bin/bash
# source
HOOKS_DIR="$(dirname "$0")/git-hooks"

# target
GIT_HOOKS_DIR="$(git rev-parse --show-toplevel)/.git/hooks"

# Copy hooks to the .git/hooks directory
for hook in "$HOOKS_DIR"/*; do
    echo "Installing $(basename "$hook") hook..."
    cp "$hook" "$GIT_HOOKS_DIR/$(basename "$hook")"
    chmod +x "$GIT_HOOKS_DIR/$(basename "$hook")"
done

echo "Git hooks have been set up."
