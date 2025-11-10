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

# Update changelog generator config with better categorization
cat > .github_changelog_generator << EOF
future-release=${VERSION}
since-tag=${PREVIOUS_VERSION}
exclude-labels=no-changelog-needed,dependencies
enhancement-labels=enhancement,feature
bug-labels=bug,bugfix,fix
breaking-labels=breaking-change,breaking
deprecated-labels=deprecation,deprecated
documentation-labels=documentation,docs
date-format=%Y-%m-%d
base=CHANGELOG.md
EOF

echo "   Updated .github_changelog_generator config with enhanced categorization"

# Check if github_changelog_generator is installed
if ! command -v github_changelog_generator &> /dev/null; then
    echo ""
    echo "WARNING: github_changelog_generator not found!"
    echo "   Install it with: gem install github_changelog_generator"
    echo "   More info: https://github.com/github-changelog-generator/github-changelog-generator"
    exit 1
fi

# Generate changelog
echo "   Running github_changelog_generator (this may take a minute)..."
github_changelog_generator --token "$GITHUB_TOKEN" 2>/dev/null || {
    echo ""
    echo "WARNING: GitHub token not set or invalid"
    echo "   Set GITHUB_TOKEN environment variable to avoid rate limiting"
    echo "   Create token at: https://github.com/settings/tokens"
    echo ""
    echo "   Trying without token (may hit rate limits)..."
    github_changelog_generator
}

echo "SUCCESS: Changelog generated successfully!"
echo "   Updated: CHANGELOG.md"

