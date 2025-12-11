#!/bin/bash
# Generate changelog using github_changelog_generator
set -e

VERSION=$1
PREVIOUS_VERSION=$2

if [ -z "$VERSION" ] || [ -z "$PREVIOUS_VERSION" ]; then
    echo "Usage: $0 <new-version> <previous-version>"
    echo "Example: $0 v1.5.4 v1.5.3"
    exit 1
fi

echo "Generating changelog from ${PREVIOUS_VERSION} to ${VERSION}..."

# Extract GitHub repo info from git remote
REMOTE_URL=$(git remote get-url origin 2>/dev/null || echo "")
if [ -z "$REMOTE_URL" ]; then
    echo "ERROR: Could not detect git remote origin"
    exit 1
fi

# Parse owner and repo from URL (works for both HTTPS and SSH)
REPO_INFO=$(echo "$REMOTE_URL" | sed -E 's|.*[:/]([^/]+)/([^/]+)(\.git)?$|\1/\2|')
REPO_OWNER=$(echo "$REPO_INFO" | cut -d'/' -f1)
REPO_NAME=$(echo "$REPO_INFO" | cut -d'/' -f2 | sed 's/\.git$//')

echo "   Repository: $REPO_OWNER/$REPO_NAME"

# Backup existing changelog if it exists
if [ -f "CHANGELOG.md" ]; then
    echo "   Backing up existing CHANGELOG.md..."
    cp CHANGELOG.md CHANGELOG.md.backup
    
    # Remove existing section for this version to prevent duplicates
    # This finds "## [VERSION]" and removes everything until the next "## [" or end of file
    echo "   Removing existing ${VERSION} section (if any) to prevent duplicates..."
    sed -i.tmp "/## \[${VERSION#v}\]/,/## \[/{//!d;}" CHANGELOG.md || true
    sed -i.tmp "/## \[${VERSION#v}\]/d" CHANGELOG.md || true
    rm -f CHANGELOG.md.tmp
fi

# Update changelog generator config
cat > .github_changelog_generator << EOF
user=${REPO_OWNER}
project=${REPO_NAME}
future-release=${VERSION}
since-tag=${PREVIOUS_VERSION}
exclude-labels=no-changelog-needed,dependencies
enhancement-labels=enhancement,feature
bug-labels=bug,bugfix,fix
breaking-labels=breaking-change,breaking
deprecated-labels=deprecation,deprecated
date-format=%Y-%m-%d
base=CHANGELOG.md
exclude-tags-regex=^((?!v\d+\.\d+\.\d+$).)*$
EOF

echo "   Updated .github_changelog_generator config"

# Check if github_changelog_generator is installed
if ! command -v github_changelog_generator &> /dev/null; then
    echo ""
    echo "ERROR: github_changelog_generator not found!"
    echo "   Install it with: gem install github_changelog_generator"
    echo "   More info: https://github.com/github-changelog-generator/github-changelog-generator"
    exit 1
fi

# Check for GitHub token
if [ -z "$GITHUB_TOKEN" ]; then
    echo ""
    echo "WARNING: GITHUB_TOKEN not set - you may hit rate limits"
    echo "   Create token at: https://github.com/settings/tokens"
    echo "   Then: export GITHUB_TOKEN='your_token_here'"
    echo ""
fi

# Generate changelog
echo "   Running github_changelog_generator (this may take a minute)..."
if [ -n "$GITHUB_TOKEN" ]; then
    github_changelog_generator --token "$GITHUB_TOKEN"
else
    github_changelog_generator
fi

echo "SUCCESS: Changelog generated successfully!"
echo "   Updated: CHANGELOG.md"
echo ""
echo "   If you see duplicates in CHANGELOG.md:"
echo "   - Restore backup: cp CHANGELOG.md.backup CHANGELOG.md"
echo "   - Or manually remove duplicate entries"

