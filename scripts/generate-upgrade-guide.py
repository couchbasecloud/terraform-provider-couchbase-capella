#!/usr/bin/env python3
"""
Generate a comprehensive upgrade guide with data from GitHub PRs.

This enhanced version automatically extracts descriptions, code examples,
and resource information from PRs to create near-complete upgrade guides.

Requirements:
    pip install PyGithub

Usage:
    export GITHUB_TOKEN="your_token_here"
    python3 generate-upgrade-guide.py 1.5.4 v1.5.3
"""

import os
import sys
from datetime import datetime
from typing import Dict, List

# Add scripts directory to path so we can import local modules
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
if SCRIPT_DIR not in sys.path:
    sys.path.insert(0, SCRIPT_DIR)

# Import helper functions
import extract_pr_content as extractor

try:
    from github import Github, Auth
except ImportError:
    print("ERROR: PyGithub not installed")
    print("   Install it with: pip install PyGithub")
    sys.exit(1)


def get_prs_since_tag(repo, since_tag: str):
    """Get all merged PRs since a specific tag."""
    print(f"   Looking up tag {since_tag}...")
    
    try:
        tag = repo.get_git_ref(f"tags/{since_tag}")
        tag_sha = tag.object.sha
        tag_commit = repo.get_commit(tag_sha)
        tag_date = tag_commit.commit.author.date
        print(f"   Tag date: {tag_date.strftime('%Y-%m-%d')}")
    except Exception as e:
        print(f"   WARNING: Could not find tag {since_tag}: {e}")
        print("   Using last 30 days instead...")
        from datetime import timedelta
        tag_date = datetime.now() - timedelta(days=30)
    
    print(f"   Fetching closed PRs (this may take 30-60 seconds)...")
    pulls = repo.get_pulls(state='closed', sort='updated', direction='desc')
    recent_prs = []
    
    pr_count = 0
    for pr in pulls:
        pr_count += 1
        if pr_count % 20 == 0:
            print(f"   Scanned {pr_count} PRs, found {len(recent_prs)} relevant...")
        if pr.merged and pr.merged_at and pr.merged_at > tag_date:
            recent_prs.append(pr)
    
    return recent_prs


def enrich_pr_data(pr, repo) -> Dict:
    """
    Enrich PR with extracted content (descriptions, examples, resources).
    """
    labels = [label.name for label in pr.labels]
    
    # Extract ticket ID
    ticket_id = extractor.extract_ticket_id(pr.title)
    title = pr.title
    if ticket_id:
        title = title.replace(f'[{ticket_id}]', '').strip()
    
    # Categorize
    category = extractor.categorize_pr_by_content(pr.title, pr.body or '', labels)
    
    # Extract description
    description = extractor.extract_description(pr.body or '')
    
    # Extract code examples
    examples = extractor.extract_terraform_examples(pr.body or '')
    
    # Get file changes to detect new resources
    try:
        files = list(pr.get_files())
        new_resources = extractor.detect_new_resources(files)
        new_datasources = extractor.detect_new_datasources(files)
    except Exception:
        files = []
        new_resources = []
        new_datasources = []
    
    # Detect deprecations
    deprecations = extractor.detect_deprecations(pr.body or '', files)
    
    return {
        'pr': pr,
        'title': title,
        'ticket_id': ticket_id,
        'category': category,
        'description': description,
        'examples': examples,
        'new_resources': new_resources,
        'new_datasources': new_datasources,
        'deprecations': deprecations,
        'labels': labels
    }


def categorize_enriched_prs(enriched_prs: List[Dict]) -> Dict:
    """Categorize enriched PR data."""
    categories = {
        'features': [],
        'enhancements': [],
        'bug_fixes': [],
        'breaking': [],
        'docs': [],
        'other': []
    }
    
    for pr_data in enriched_prs:
        # Skip no-changelog-needed
        if 'no-changelog-needed' in [l.lower() for l in pr_data['labels']]:
            continue
        
        category = pr_data['category']
        if category == 'breaking':
            categories['breaking'].append(pr_data)
        elif category == 'feature':
            categories['features'].append(pr_data)
        elif category == 'enhancement':
            categories['enhancements'].append(pr_data)
        elif category == 'bug':
            categories['bug_fixes'].append(pr_data)
        elif category == 'docs':
            categories['docs'].append(pr_data)
        else:
            categories['other'].append(pr_data)
    
    return categories


def generate_feature_bullet(pr_data: Dict) -> str:
    """Generate a simple bullet point for the New Features list."""
    title = pr_data['title']
    new_resources = pr_data['new_resources']
    new_datasources = pr_data['new_datasources']
    description = pr_data['description']
    
    # Extract first sentence of description if available
    desc_text = ""
    if description:
        first_sentence = description.split('.')[0].strip()
        if first_sentence and len(first_sentence) < 150:
            desc_text = first_sentence
        else:
            desc_text = title
    else:
        desc_text = title
    
    bullet = f"* {desc_text}"
    
    # Add resource links inline
    if new_resources:
        for resource in new_resources:
            resource_name = extractor.format_resource_name(resource)
            registry_url = extractor.generate_registry_url(resource, is_datasource=False)
            bullet += f" [`{resource_name}`]({registry_url})"
    
    if new_datasources:
        for datasource in new_datasources:
            datasource_name = extractor.format_resource_name(datasource)
            registry_url = extractor.generate_registry_url(datasource, is_datasource=True)
            bullet += f" [`{datasource_name}`]({registry_url})"
    
    bullet += "\n"
    return bullet


def generate_detailed_feature_section(pr_data: Dict, version: str) -> str:
    """Generate a detailed section for a major feature (goes after Changes section)."""
    pr = pr_data['pr']
    title = pr_data['title']
    description = pr_data['description']
    examples = pr_data['examples']
    new_resources = pr_data['new_resources']
    
    # Only generate detailed section if we have substantial content
    if not description or not examples:
        return ""
    
    section = f"## {title}\n\n"
    
    # Add full description
    section += f"{description}\n\n"
    
    # Add usage instructions
    if new_resources:
        resource = new_resources[0]
        resource_name = extractor.format_resource_name(resource)
        section += f"To use the `{resource}` resource:\n\n"
    
    # Add code example
    if examples:
        for example in examples:
            is_valid, error = extractor.validate_terraform_code(example)
            if is_valid:
                section += "```\n"  # No language hint, plain markdown code block
                section += example
                if not example.endswith('\n'):
                    section += '\n'
                section += "```\n\n"
                break
    
    # Add link to examples
    if new_resources:
        for resource in new_resources:
            examples_url = extractor.generate_examples_url(resource)
            section += f"For more information, see the [examples for {resource}]({examples_url})\n\n"
    
    return section


def generate_guide(version: str, previous_version: str, prs_by_category: Dict) -> str:
    """Generate the complete upgrade guide markdown matching the existing template."""
    
    version_clean = version.replace('.', '')
    
    guide = f"""---
layout: "couchbase-capella"
page_title: "Couchbase Capella Provider {version}: Upgrade and Information Guide"
sidebar_current: "docs-couchbase-capella-guides-{version_clean}-upgrade-guide"
description: |-
Couchbase Capella Provider {version}: Upgrade and Information Guide
---

# Couchbase Capella Provider {version}: Upgrade and Information Guide

"""
    
    # Determine release type
    has_features = prs_by_category['features'] or prs_by_category['enhancements']
    has_breaking = bool(prs_by_category['breaking'])
    is_bugfix_only = not has_features and not has_breaking
    
    # Add intro line for feature releases
    if has_features:
        guide += f"Here is a list of what's new in {version}\n\n"
    
    # New Features section - simple bullet list
    if has_features:
        guide += "## New Features\n\n"
        
        all_features = prs_by_category['features'] + prs_by_category['enhancements']
        for pr_data in all_features:
            guide += generate_feature_bullet(pr_data)
        guide += "\n"
    
    # Bug Fixes section - simple list
    if prs_by_category['bug_fixes']:
        guide += "## Bug Fixes\n\n"
        for pr_data in prs_by_category['bug_fixes']:
            title = pr_data['title']
            description = pr_data['description']
            
            # Use description if available, otherwise title
            if description:
                # Use first sentence for bug fixes
                first_sentence = description.split('.')[0].strip()
                if first_sentence and len(first_sentence) < 200:
                    guide += f"* {first_sentence}\n"
                else:
                    guide += f"* {title}\n"
            else:
                guide += f"* {title}\n"
        guide += "\n"
    
    # Changes section
    guide += "## Changes\n\n"
    
    if not has_breaking:
        guide += "There are no deprecations as part of this release.\n\n"
    
    # For bug fix releases, mention general improvements
    if is_bugfix_only:
        guide += f"{version} also includes general improvements and bug fixes. "
        guide += "See the [CHANGELOG](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/master/CHANGELOG.md) "
        guide += "for more specific information.\n\n"
    elif has_features:
        # For feature releases
        guide += f"{version} includes new features and general improvements. "
        guide += "See the [CHANGELOG](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/master/CHANGELOG.md) "
        guide += "for more specific information.\n\n"
    
    # List new resources/datasources in Changes section (like 1.5.3 does)
    new_items = []
    for pr_data in (prs_by_category['features'] + prs_by_category['enhancements']):
        for resource in pr_data['new_resources']:
            resource_name = extractor.format_resource_name(resource)
            registry_url = extractor.generate_registry_url(resource, is_datasource=False)
            new_items.append((resource_name, registry_url, 'resource'))
        for datasource in pr_data['new_datasources']:
            datasource_name = extractor.format_resource_name(datasource)
            registry_url = extractor.generate_registry_url(datasource, is_datasource=True)
            new_items.append((datasource_name, registry_url, 'datasource'))
    
    if new_items:
        for name, url, item_type in new_items:
            desc = "Manage" if item_type == 'resource' else "Retrieve"
            # Make a friendly name from the resource
            friendly_name = name.replace('couchbase-capella_', '').replace('_', ' ').title()
            guide += f"* {desc} {friendly_name} [`{name}`]({url})\n"
        guide += "\n"
    
    # Breaking Changes section (if any)
    if prs_by_category['breaking']:
        guide += "## Breaking Changes\n\n"
        guide += "WARNING: **ACTION REQUIRED** - The following changes may require updates to your Terraform configurations:\n\n"
        
        for pr_data in prs_by_category['breaking']:
            title = pr_data['title']
            description = pr_data['description']
            deprecations = pr_data['deprecations']
            
            guide += f"### {title}\n\n"
            
            if description:
                guide += f"{description}\n\n"
            
            if deprecations:
                guide += "**Deprecated:**\n"
                for dep in deprecations:
                    guide += f"- `{dep['field']}`\n"
                guide += "\n"
            
            guide += "<!-- TODO: Add migration steps if needed -->\n\n"
    
    # Detailed feature sections (only for features with good descriptions and examples)
    detailed_sections = []
    all_features = prs_by_category['features'] + prs_by_category['enhancements']
    for pr_data in all_features:
        section = generate_detailed_feature_section(pr_data, version)
        if section:
            detailed_sections.append(section)
    
    if detailed_sections:
        for section in detailed_sections:
            guide += section
    
    # Footer
    guide += """### Helpful Links

- [Getting Started with the Terraform Provider](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/master/examples/getting_started)
- [Capella Management API v4.0](https://docs.couchbase.com/cloud/management-api-reference/index.html)
- [See Specific Examples](https://github.com/couchbasecloud/terraform-provider-couchbase-capella/blob/master/examples)
"""
    
    return guide


def main():
    if len(sys.argv) != 3:
        print("Usage: generate-upgrade-guide.py <new-version> <previous-version>")
        print("Example: generate-upgrade-guide.py 1.5.4 v1.5.3")
        sys.exit(1)
    
    version = sys.argv[1].lstrip('v')
    previous_version = sys.argv[2]
    if not previous_version.startswith('v'):
        previous_version = f'v{previous_version}'
    
    # Get GitHub token from environment
    token = os.environ.get('GITHUB_TOKEN')
    if not token:
        print("ERROR: GITHUB_TOKEN environment variable not set")
        print("   Create a token at: https://github.com/settings/tokens")
        print("   Required scopes: repo (for private repos) or public_repo")
        print("")
        print("   Set it with: export GITHUB_TOKEN='your_token_here'")
        sys.exit(1)
    
    print(f" Generating upgrade guide for v{version}...")
    print(f"   Analyzing changes since {previous_version}...")
    
    # Connect to GitHub
    try:
        auth = Auth.Token(token)
        g = Github(auth=auth)
        repo = g.get_repo("couchbasecloud/terraform-provider-couchbase-capella")
    except Exception as e:
        print(f"ERROR: Failed to connect to GitHub: {e}")
        sys.exit(1)
    
    # Get PRs since last release
    try:
        prs = get_prs_since_tag(repo, previous_version)
        print(f"   Found {len(prs)} merged PRs")
    except Exception as e:
        print(f"ERROR: Failed to fetch PRs: {e}")
        sys.exit(1)
    
    # Enrich PR data with extracted content
    print(f"   Extracting content from {len(prs)} PRs (this may take 1-2 minutes)...")
    enriched_prs = []
    for i, pr in enumerate(prs, 1):
        if i % 3 == 0 or i == 1:
            print(f"   Processing PR {i}/{len(prs)}: #{pr.number} - {pr.title[:50]}...")
        enriched = enrich_pr_data(pr, repo)
        enriched_prs.append(enriched)
    
    # Categorize PRs
    print(f"   Categorizing PRs by type...")
    prs_by_category = categorize_enriched_prs(enriched_prs)
    
    # Print summary
    print(f"")
    print(f"   PR Summary:")
    print(f"   - Features: {len(prs_by_category['features'])}")
    print(f"   - Enhancements: {len(prs_by_category['enhancements'])}")
    print(f"   - Bug Fixes: {len(prs_by_category['bug_fixes'])}")
    print(f"   - Breaking Changes: {len(prs_by_category['breaking'])}")
    print(f"   - Documentation: {len(prs_by_category['docs'])}")
    print(f"   - Other: {len(prs_by_category['other'])}")
    
    # Count how many have descriptions/examples
    total_features = len(prs_by_category['features']) + len(prs_by_category['enhancements'])
    with_description = sum(1 for pd in (prs_by_category['features'] + prs_by_category['enhancements']) 
                          if pd['description'])
    with_examples = sum(1 for pd in (prs_by_category['features'] + prs_by_category['enhancements']) 
                       if pd['examples'])
    with_resources = sum(1 for pd in (prs_by_category['features'] + prs_by_category['enhancements']) 
                        if pd['new_resources'] or pd['new_datasources'])
    
    if total_features > 0:
        print(f"   - Features with descriptions: {with_description}/{total_features}")
        print(f"   - Features with code examples: {with_examples}/{total_features}")
        print(f"   - Features with detected resources: {with_resources}/{total_features}")
    
    # Generate guide
    print(f"")
    print(f"   Generating upgrade guide document...")
    guide_content = generate_guide(version, previous_version, prs_by_category)
    
    # Write to file
    print(f"   Writing guide to file...")
    guide_file = f"templates/guides/{version}-upgrade-guide.md"
    try:
        with open(guide_file, 'w') as f:
            f.write(guide_content)
    except Exception as e:
        print(f"ERROR: Failed to write file: {e}")
        sys.exit(1)
    
    print(f"SUCCESS: Upgrade guide created at: {guide_file}")
    print("")
    
    # Calculate completeness
    todos_count = guide_content.count('<!-- TODO')
    if todos_count == 0:
        print("SUCCESS: Guide is complete! No TODOs remaining.")
    elif todos_count <= 2:
        print(f"INFO: Guide is mostly complete! Only {todos_count} TODO(s) remaining.")
    else:
        print(f"WARNING:  {todos_count} TODO(s) remaining - review and enhance:")
    
    print("")
    print(" Next steps:")
    print(f"   1. Review {guide_file}")
    print("   2. Fill in any remaining TODOs")
    print("   3. Enhance descriptions for clarity")
    print("   4. Test any code examples")
    print("   5. Run 'make build-docs' to publish to docs/")


if __name__ == "__main__":
    main()
