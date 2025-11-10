#!/bin/bash
# Generate upgrade guide scaffold
set -e

VERSION=$1
PREVIOUS_VERSION=$2

if [ -z "$VERSION" ] || [ -z "$PREVIOUS_VERSION" ]; then
    echo "Usage: $0 <new-version> <previous-version>"
    echo "Example: $0 1.5.4 v1.5.3"
    exit 1
fi

# Remove 'v' prefix if present in PREVIOUS_VERSION
PREVIOUS_VERSION=${PREVIOUS_VERSION#v}
# Remove 'v' prefix from VERSION if present
VERSION=${VERSION#v}

GUIDE_FILE="templates/guides/${VERSION}-upgrade-guide.md"

echo "Generating upgrade guide scaffold for v${VERSION}..."

# Check if file already exists
if [ -f "$GUIDE_FILE" ]; then
    echo "WARNING: File already exists: $GUIDE_FILE"
    read -p "   Overwrite? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted"
        exit 1
    fi
fi

# Get list of PRs merged since last release (using git log)
echo "   Analyzing commits since v${PREVIOUS_VERSION}..."

# Extract PR numbers from commit messages
PR_NUMBERS=$(git log --oneline v${PREVIOUS_VERSION}..HEAD 2>/dev/null | grep -oE '#[0-9]+' | sort -u | sed 's/#//' || echo "")

if [ -z "$PR_NUMBERS" ]; then
    echo "   WARNING: No PRs found in git history since v${PREVIOUS_VERSION}"
    echo "   Creating basic template..."
fi

# Create version without dots for sidebar ID
VERSION_CLEAN=$(echo "$VERSION" | tr -d '.')

# Create the guide template matching the existing format
cat > "$GUIDE_FILE" << 'EOF'
---
layout: "couchbase-capella"
page_title: "Couchbase Capella Provider VERSION: Upgrade and Information Guide"
sidebar_current: "docs-couchbase-capella-guides-VERSION_CLEAN-upgrade-guide"
description: |-
Couchbase Capella Provider VERSION: Upgrade and Information Guide
---

# Couchbase Capella Provider VERSION: Upgrade and Information Guide

Here is a list of what's new in VERSION

## New Features

<!-- TODO: Add feature bullets with inline resource links
Example:
* Create and manage X with [`couchbase-capella_resource_name`](https://registry.terraform.io/providers/couchbasecloud/couchbase-capella/latest/docs/resources/resource_name)
-->

## Bug Fixes

<!-- TODO: List bug fixes (simple bullets)
Example:
* Fixed issue where X was not working correctly
-->

## Changes

There are no deprecations as part of this release.

VERSION includes new features and general improvements. See the [CHANGELOG](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/master/CHANGELOG.md) for more specific information.

<!-- TODO: List new resources/datasources
Example:
* Manage Feature Name [`couchbase-capella_resource`](registry_url)
-->

<!-- TODO: Add detailed sections for major features (optional)
## Feature Name

Description of the feature.

To use the resource:

```
resource "couchbase-capella_resource" "example" {
  # configuration
}
```

For more information, see the [examples](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples/resource_name)
-->

### Helpful Links

- [Getting Started with the Terraform Provider](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/master/examples/getting_started)
- [Capella Management API v4.0](https://docs.couchbase.com/cloud/management-api-reference/index.html)
- [See Specific Examples](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/master/examples)
EOF

# Replace VERSION placeholders
sed -i.bak "s/VERSION/${VERSION}/g" "$GUIDE_FILE"
sed -i.bak "s/VERSION_CLEAN/${VERSION_CLEAN}/g" "$GUIDE_FILE"
rm "${GUIDE_FILE}.bak" 2>/dev/null || true

echo "SUCCESS: Upgrade guide scaffold created at: $GUIDE_FILE"
echo ""
echo "PRs since v${PREVIOUS_VERSION}:"
if [ -n "$PR_NUMBERS" ]; then
    echo "$PR_NUMBERS" | while read -r pr; do
        echo "   - https://github.com/couchbasecloud/terraform-provider-couchbase-capella/pull/$pr"
    done
else
    echo "   (Check CHANGELOG.md for full list)"
fi
echo ""
echo "ACTION REQUIRED:"
echo "   1. Review PRs listed above and in CHANGELOG.md"
echo "   2. Edit $GUIDE_FILE:"
echo "      - Add feature descriptions with code examples"
echo "      - Document bug fixes"
echo "      - List any breaking changes with migration steps"
echo "      - Add deprecation notices"
echo "   3. Test all code examples"
echo "   4. Run 'make build-docs' to publish to docs/"

