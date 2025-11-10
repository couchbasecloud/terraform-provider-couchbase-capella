#!/usr/bin/env python3
"""
Validate PRs for release documentation quality.

Pre-flight check before generating upgrade guides to catch poorly-documented PRs.

Usage:
    export GITHUB_TOKEN="your_token"
    python3 validate-prs.py v1.5.3 v1.5.2
"""

import os
import sys
from typing import List, Dict

try:
    from github import Github
except ImportError:
    print("ERROR: PyGithub not installed")
    print("   Install it with: pip install PyGithub")
    sys.exit(1)


def get_prs_since_tag(repo, since_tag: str):
    """Get all merged PRs since a specific tag."""
    try:
        tag = repo.get_git_ref(f"tags/{since_tag}")
        tag_sha = tag.object.sha
        tag_commit = repo.get_commit(tag_sha)
        tag_date = tag_commit.commit.author.date
    except Exception as e:
        print(f"WARNING:  Could not find tag {since_tag}: {e}")
        from datetime import datetime, timedelta
        tag_date = datetime.now() - timedelta(days=30)
    
    pulls = repo.get_pulls(state='closed', sort='updated', direction='desc')
    recent_prs = []
    
    for pr in pulls:
        if pr.merged and pr.merged_at and pr.merged_at > tag_date:
            recent_prs.append(pr)
    
    return recent_prs


def validate_pr(pr) -> Dict[str, any]:
    """
    Validate a single PR for documentation quality.
    
    Returns a dict with validation results.
    """
    issues = []
    warnings = []
    
    # Check if PR is skipped from changelog
    labels = [label.name.lower() for label in pr.labels]
    if 'no-changelog-needed' in labels:
        return {'pr': pr, 'skipped': True, 'issues': [], 'warnings': []}
    
    # Check for labels
    has_type_label = any(label in labels for label in [
        'feature', 'enhancement', 'bug', 'bugfix', 'fix', 
        'documentation', 'docs', 'breaking-change', 'breaking'
    ])
    
    if not has_type_label:
        warnings.append("No type label (feature/bug/enhancement/docs)")
    
    # Check for description
    if not pr.body or len(pr.body.strip()) < 20:
        issues.append("PR description is empty or too short (< 20 chars)")
    
    # Check for ticket reference in title
    import re
    if not re.search(r'\[([A-Z]+-\d+)\]', pr.title):
        warnings.append("No ticket ID in title (e.g., [AV-12345])")
    
    # Check if title is meaningful
    title_lower = pr.title.lower()
    vague_titles = ['update', 'fix', 'change', 'modify', 'refactor']
    if any(title_lower.strip().startswith(word) for word in vague_titles):
        if len(pr.title.split()) < 4:
            warnings.append("Title is vague - add more context")
    
    # For feature PRs, check for code examples
    if 'feature' in labels or 'enhancement' in labels:
        if not pr.body or '```' not in pr.body:
            warnings.append("Feature PR should include code examples")
    
    # Check for breaking changes mentioned but not labeled
    if pr.body and 'breaking' in pr.body.lower():
        if 'breaking-change' not in labels and 'breaking' not in labels:
            issues.append("Mentions breaking changes but missing 'breaking-change' label")
    
    return {
        'pr': pr,
        'skipped': False,
        'issues': issues,
        'warnings': warnings
    }


def print_validation_results(results: List[Dict]):
    """Print validation results in a readable format."""
    
    total_prs = len(results)
    skipped_prs = len([r for r in results if r['skipped']])
    validated_prs = total_prs - skipped_prs
    
    prs_with_issues = [r for r in results if not r['skipped'] and r['issues']]
    prs_with_warnings = [r for r in results if not r['skipped'] and r['warnings']]
    
    print(f"\n📊 PR Validation Results")
    print(f"=" * 60)
    print(f"Total PRs: {total_prs}")
    print(f"Validated: {validated_prs}")
    print(f"Skipped (no-changelog-needed): {skipped_prs}")
    print(f"PRs with issues: {len(prs_with_issues)}")
    print(f"PRs with warnings: {len(prs_with_warnings)}")
    
    if prs_with_issues:
        print(f"\nERROR: Issues Found ({len(prs_with_issues)} PRs)")
        print("=" * 60)
        for result in prs_with_issues:
            pr = result['pr']
            print(f"\nPR #{pr.number}: {pr.title}")
            print(f"   URL: {pr.html_url}")
            for issue in result['issues']:
                print(f"   ERROR: {issue}")
    
    if prs_with_warnings:
        print(f"\nWARNING:  Warnings ({len(prs_with_warnings)} PRs)")
        print("=" * 60)
        for result in prs_with_warnings:
            pr = result['pr']
            print(f"\nPR #{pr.number}: {pr.title}")
            print(f"   URL: {pr.html_url}")
            for warning in result['warnings']:
                print(f"   WARNING:  {warning}")
    
    if not prs_with_issues and not prs_with_warnings:
        print(f"\nSUCCESS: All PRs look good!")
    
    # Return exit code
    return 0 if not prs_with_issues else 1


def main():
    if len(sys.argv) != 3:
        print("Usage: validate-prs.py <new-version> <previous-version>")
        print("Example: validate-prs.py v1.5.4 v1.5.3")
        sys.exit(1)
    
    new_version = sys.argv[1]
    previous_version = sys.argv[2]
    
    if not previous_version.startswith('v'):
        previous_version = f'v{previous_version}'
    
    # Get GitHub token
    token = os.environ.get('GITHUB_TOKEN')
    if not token:
        print("ERROR: GITHUB_TOKEN environment variable not set")
        print("   Set it with: export GITHUB_TOKEN='your_token_here'")
        sys.exit(1)
    
    print(f"🔍 Validating PRs for release {new_version}")
    print(f"   Checking PRs since {previous_version}...")
    
    # Connect to GitHub
    try:
        g = Github(token)
        repo = g.get_repo("couchbasecloud/terraform-provider-couchbase-capella")
    except Exception as e:
        print(f"ERROR: Failed to connect to GitHub: {e}")
        sys.exit(1)
    
    # Get PRs
    try:
        prs = get_prs_since_tag(repo, previous_version)
        print(f"   Found {len(prs)} merged PRs")
    except Exception as e:
        print(f"ERROR: Failed to fetch PRs: {e}")
        sys.exit(1)
    
    # Validate each PR
    results = []
    for pr in prs:
        result = validate_pr(pr)
        results.append(result)
    
    # Print results
    exit_code = print_validation_results(results)
    
    if exit_code != 0:
        print(f"\nTIP: Recommendation:")
        print(f"   Fix the issues above to improve upgrade guide quality")
        print(f"   Or proceed anyway - the generator will work with what's available")
    
    sys.exit(exit_code)


if __name__ == "__main__":
    main()

