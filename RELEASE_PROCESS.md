# Release Process

This document describes the automated release process for the Terraform Provider Couchbase Capella.

## Overview

Releases are semi-automated using the `make release-prep` command, which:
1. Generates the changelog from GitHub PRs
2. Creates an upgrade guide scaffold
3. Builds the documentation

The actual release is triggered automatically via GitHub Actions when a tag is pushed.

## Prerequisites

### Required Tools

```bash
# Ruby gem for changelog generation
gem install github_changelog_generator
```

### Recommended (Enhanced Automation)

For the best results with automatic content extraction:

```bash
# Python with PyGithub for intelligent PR analysis
pip install PyGithub

# Set GitHub token for API access (avoids rate limiting)
export GITHUB_TOKEN="your_github_token_here"
```

**With PyGithub installed, you get:**
- Automatic extraction of PR descriptions
- Code example detection and inclusion
- New resource/data source detection
- Deprecation detection
- PR quality validation
- 80%+ complete upgrade guides (vs 30% without)

**Without PyGithub:**
- Falls back to basic bash scripts
- Creates simple scaffolds with TODOs
- Still works, just requires more manual editing

To create a GitHub token:
1. Go to https://github.com/settings/tokens
2. Generate a new token (classic)
3. Select scope: `repo` (for private repos) or `public_repo`
4. Copy the token and set it as an environment variable

## Important: Templates vs Docs Directory

**Always edit files in `templates/guides/`, never in `docs/guides/`**

- `templates/guides/` - SOURCE files (you edit these)
- `docs/guides/` - GENERATED files (auto-created by `make build-docs`)

The workflow:
1. Scripts create upgrade guides in `templates/guides/X.Y.Z-upgrade-guide.md`
2. You edit the file in `templates/guides/`
3. Run `make build-docs` to copy `templates/` → `docs/`
4. Commit both `templates/` and `docs/` directories

Why? The `tfplugindocs` tool generates documentation from schema definitions and copies custom content from `templates/` to `docs/`. If you edit `docs/` directly, your changes will be overwritten next time docs are built.

## Quick Start

### Step 1: Prepare the Release

```bash
# Prepare release artifacts for version 1.5.4
make release-prep VERSION=1.5.4

# This will:
# 1. Update .github_changelog_generator config with enhanced categorization
# 2. Generate CHANGELOG.md from GitHub PRs
# 3. Validate PR quality (checks labels, descriptions, etc.)
# 4. Generate upgrade guide with auto-extracted content:
#    - PR descriptions
#    - Code examples from PR bodies
#    - New resource detection
#    - Deprecation detection
# 5. Build documentation (copies to docs/)
```

### Step 2: Review and Enhance

1. **Review the changelog:**
   ```bash
   # Check CHANGELOG.md for accuracy
   git diff CHANGELOG.md
   ```

2. **(Optional) Edit the upgrade guide:**
   ```bash
   # Edit the scaffold with real content (replace <VERSION> with your version, e.g., 1.7.0)
   vim templates/guides/<VERSION>-upgrade-guide.md
   ```
   
   Add:
   - Detailed feature descriptions
   - Code examples for new features
   - Migration steps for breaking changes
   - Bug fix details

3. **(Optional) Test examples:**
   ```bash
   # Test any code examples in the upgrade guide
   cd examples/feature_name
   terraform init && terraform plan
   ```

4. **(Optional) Rebuild docs:**
   ```bash
   # After editing templates, rebuild docs
   make build-docs
   ```

### Step 3: Create Pull Request

> **Note:** Replace `<VERSION>` with your actual version number (e.g., `1.7.0`)

```bash
# Create a new branch for the release
git checkout -b release/v<VERSION>

# Commit all changes
git add .
git commit -m "Prepare release v<VERSION>"

# Push the branch
git push origin release/v<VERSION>
```

**Example for version 1.7.0:**
```bash
git checkout -b release/v1.7.0
git add .
git commit -m "Prepare release v1.7.0"
git push origin release/v1.7.0
```

Open a Pull Request for review on GitHub.

### Step 4: Merge and Tag

After the PR is reviewed and approved:

> **Note:** Replace `<VERSION>` with your actual version number (e.g., `1.7.0`)

```bash
# Merge the PR to main (via GitHub UI or CLI)
gh pr merge release/v<VERSION> --merge

# Switch to main and pull the latest changes
git checkout main
git pull origin main

# Create and push the tag from main branch
git tag v<VERSION>
git push origin v<VERSION>

# GitHub Actions will automatically:
# - Build binaries
# - Create GitHub release
# - Upload to Terraform Registry
```

**Example for version 1.7.0:**
```bash
gh pr merge release/v1.7.0 --merge
git checkout main
git pull origin main
git tag v1.7.0
git push origin v1.7.0
```

### Step 5: Verify

1. Check GitHub Actions workflow: https://github.com/couchbasecloud/terraform-provider-couchbase-capella/actions
2. Verify release: https://github.com/couchbasecloud/terraform-provider-couchbase-capella/releases
3. Check Terraform Registry: https://registry.terraform.io/providers/couchbasecloud/couchbase-capella/latest

## Release Checklist

Use this command to display the full checklist:

```bash
make release-checklist
```

## Manual Process (Without Automation)

If you prefer not to use the automation or it's unavailable:

### 1. Update Changelog Config

Edit `.github_changelog_generator`:
```
future-release=v1.5.4
since-tag=v1.5.3
exclude-labels=no-changelog-needed
date-format=%Y-%m-%d
base=CHANGELOG.md
```

### 2. Generate Changelog

```bash
github-changelog-generator --token $GITHUB_TOKEN
```

### 3. Create Upgrade Guide

```bash
# Create the file
touch templates/guides/1.5.4-upgrade-guide.md

# Add content based on template structure
# See existing guides for examples
```

### 4. Build Documentation

```bash
make build-docs
```

### 5. Continue with Step 3 above (Create Pull Request, Merge and Tag)

## Troubleshooting

### Issue: `github_changelog_generator` not found

```bash
# Install the gem
gem install github_changelog_generator

# Or on macOS with homebrew
brew install github_changelog_generator
```

### Issue: GitHub API rate limiting

```bash
# Set a GitHub token to increase rate limits
export GITHUB_TOKEN="your_token_here"

# Then re-run the release prep
make release-prep VERSION=1.5.4
```

### Issue: PyGithub not installed

```bash
# For enhanced upgrade guide generation
pip install PyGithub

# Or use the basic bash version (it will fall back automatically)
```

### Issue: Previous version not detected correctly

```bash
# Manually specify previous version
bash scripts/generate-changelog.sh v1.5.4 v1.5.3
bash scripts/generate-upgrade-guide.sh 1.5.4 v1.5.3
```

## Upgrade Guide Template

When creating an upgrade guide, use this structure:

```markdown
---
layout: "couchbase-capella"
page_title: "Couchbase Capella Provider X.Y.Z: Upgrade and Information Guide"
sidebar_current: "docs-couchbase-capella-guides-XYZ-upgrade-guide"
description: |-
Couchbase Capella Provider X.Y.Z: Upgrade and Information Guide
---

# Couchbase Capella Provider X.Y.Z: Upgrade and Information Guide

## New Features

* Feature name [`resource_name`](link)
  - Brief description
  - Code example

## Bug Fixes

* Brief description [#PR](link)

## Breaking Changes

If none:
There are no breaking changes as part of this release.

If any:
WARNING: **ACTION REQUIRED**
* Change description
* Migration steps

## Changes

List deprecations and other changes.

### Helpful Links

- [Getting Started](link)
- [API Reference](link)
- [Examples](link)
```

## Best Practices

### PR Labeling

For best results with automated changelog generation, ensure all PRs have proper labels:

- `enhancement` or `feature` - New features
- `bug` or `bugfix` - Bug fixes
- `breaking-change` - Breaking changes
- `documentation` - Documentation updates
- `no-changelog-needed` - Skip in changelog (internal changes)

### Commit Messages

Use conventional commit messages for clarity:

```
[AV-12345] Add support for feature X
[AV-12346] Fix bug in resource Y
[AV-12347] BREAKING: Remove deprecated field Z
```

### Testing

Before releasing:
-  Run all tests: `make test`
-  Run acceptance tests: `make testacc`
-  Run linters: `make check`
-  Test examples in upgrade guide

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Breaking changes
- **MINOR** (0.X.0): New features (backwards compatible)
- **PATCH** (0.0.X): Bug fixes (backwards compatible)

## What Gets Automated

### Changelog Generation

The `generate-changelog.sh` script creates a comprehensive changelog with:
- Automatic PR categorization (features, bugs, breaking changes, etc.)
- Links to PRs and issues
- Author attribution
- Filtered by labels (excludes dependencies, no-changelog-needed)

### Upgrade Guide Generation (Enhanced)

The enhanced `generate-upgrade-guide.py` script automatically:

1. **Extracts Descriptions** - Parses PR bodies for feature descriptions
2. **Finds Code Examples** - Detects Terraform code blocks in PR descriptions
3. **Detects Resources** - Scans file changes for new resources/data sources
4. **Auto-links** - Creates links to Terraform Registry and examples
5. **Finds Deprecations** - Detects deprecated fields mentioned in PRs
6. **Validates Quality** - Checks PR labels and descriptions

**Example Output (matching existing template):**

For feature releases, you get:
```markdown
Here is a list of what's new in 1.5.4

## New Features

* Create and manage free tier clusters with [`couchbase-capella_free_tier_cluster`](registry_link)
* Turn clusters on/off on demand with [`couchbase-capella_cluster_onoff_ondemand`](registry_link)

## Changes

There are no deprecations as part of this release.

1.5.4 includes new features and general improvements. See the [CHANGELOG](...) for more specific information.

* Manage Free Tier Cluster [`couchbase-capella_free_tier_cluster`](registry_link)

## Managing Free Tier Clusters

Use the new free_tier_cluster resource to create and manage free tier operational clusters...

```
resource "couchbase-capella_free_tier_cluster" "example" {
  organization_id = var.organization_id
  ...
}
```

For more information, see the [examples for free_tier_cluster](examples_link)
```

The format matches your existing upgrade guides in `docs/guides/`.

### PR Validation

The `validate-prs.py` script checks:
- All PRs have appropriate type labels
- PRs have meaningful descriptions
- Feature PRs include code examples
- Breaking changes are properly labeled
- Ticket references are present

## Additional Resources

- [GoReleaser Documentation](https://goreleaser.com/)
- [GitHub Changelog Generator](https://github.com/github-changelog-generator/github-changelog-generator)
- [Terraform Registry Publishing](https://www.terraform.io/docs/registry/providers/publishing.html)
- [Terraform Provider Development](https://www.terraform.io/docs/extend/index.html)

