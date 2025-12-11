#!/usr/bin/env python3
"""
Helper module for extracting meaningful content from GitHub PRs.
Used by generate-upgrade-guide.py to create better documentation.
"""

import re
from typing import Dict, List, Optional, Tuple


def extract_description(pr_body: str) -> Optional[str]:
    """
    Extract meaningful description from PR body.
    
    Looks for common section headers like "## Description", "## What", "## Summary"
    and extracts the content following them.
    """
    if not pr_body:
        return None
    
    # Common description section patterns
    patterns = [
        r'##\s*Description\s*\n(.*?)(?=\n##|\n---|\Z)',
        r'##\s*What\s*\n(.*?)(?=\n##|\n---|\Z)',
        r'##\s*Summary\s*\n(.*?)(?=\n##|\n---|\Z)',
        r'##\s*Overview\s*\n(.*?)(?=\n##|\n---|\Z)',
    ]
    
    for pattern in patterns:
        match = re.search(pattern, pr_body, re.DOTALL | re.IGNORECASE)
        if match:
            description = match.group(1).strip()
            # Clean up the description
            description = clean_description(description)
            if description and len(description) > 20:  # Meaningful length
                return description
    
    # Fallback: Extract first paragraph if no section found
    paragraphs = pr_body.strip().split('\n\n')
    for para in paragraphs:
        cleaned = clean_description(para)
        if cleaned and len(cleaned) > 20 and not para.startswith('#'):
            return cleaned
    
    return None


def clean_description(text: str) -> str:
    """Clean up extracted description text."""
    # Remove markdown checkboxes
    text = re.sub(r'- \[[ x]\]\s*', '- ', text)
    
    # Remove excessive newlines
    text = re.sub(r'\n{3,}', '\n\n', text)
    
    # Remove HTML comments
    text = re.sub(r'<!--.*?-->', '', text, flags=re.DOTALL)
    
    # Trim
    text = text.strip()
    
    return text


def extract_terraform_examples(pr_body: str) -> List[str]:
    """
    Extract Terraform/HCL code blocks from PR description.
    
    Returns a list of code blocks found.
    """
    if not pr_body:
        return []
    
    examples = []
    
    # Pattern for markdown code blocks with terraform/hcl/tf language hint
    patterns = [
        r'```(?:terraform|hcl|tf)\n(.*?)```',
        r'```\n(resource "couchbase-capella_.*?)```',  # Detect by content
    ]
    
    for pattern in patterns:
        matches = re.finditer(pattern, pr_body, re.DOTALL)
        for match in matches:
            code = match.group(1).strip()
            if code and 'couchbase-capella' in code:
                examples.append(code)
    
    return examples


def detect_new_resources(pr_files: List[Dict]) -> List[str]:
    """
    Detect new Terraform resources from PR file changes.
    
    Args:
        pr_files: List of file objects from GitHub API (each has 'filename', 'status')
    
    Returns:
        List of resource names (e.g., ['free_tier_cluster', 'app_service'])
    """
    new_resources = []
    
    for file_obj in pr_files:
        filename = file_obj.get('filename', '')
        status = file_obj.get('status', '')
        
        # Look for new resource files
        if status == 'added' and 'internal/resources/' in filename:
            # Extract resource name from filename
            # e.g., internal/resources/free_tier_cluster.go -> free_tier_cluster
            match = re.search(r'internal/resources/([a-z_]+)\.go$', filename)
            if match:
                resource_name = match.group(1)
                # Skip utility files
                if not resource_name.endswith('_schema') and resource_name not in [
                    'schema', 'state', 'attributes', 'utils'
                ]:
                    new_resources.append(resource_name)
    
    return new_resources


def detect_new_datasources(pr_files: List[Dict]) -> List[str]:
    """
    Detect new Terraform data sources from PR file changes.
    
    Args:
        pr_files: List of file objects from GitHub API
    
    Returns:
        List of data source names
    """
    new_datasources = []
    
    for file_obj in pr_files:
        filename = file_obj.get('filename', '')
        status = file_obj.get('status', '')
        
        # Look for new data source files
        if status == 'added' and 'internal/datasources/' in filename:
            match = re.search(r'internal/datasources/([a-z_]+)\.go$', filename)
            if match:
                datasource_name = match.group(1)
                # Skip utility files
                if not datasource_name.endswith('_schema') and datasource_name not in [
                    'schema', 'state', 'attributes', 'utils'
                ]:
                    new_datasources.append(datasource_name)
    
    return new_datasources


def extract_ticket_id(pr_title: str) -> Optional[str]:
    """
    Extract ticket ID from PR title.
    
    Examples:
        "[AV-12345] Add feature" -> "AV-12345"
        "[JIRA-123] Fix bug" -> "JIRA-123"
    """
    match = re.search(r'\[([A-Z]+-\d+)\]', pr_title)
    if match:
        return match.group(1)
    return None


def categorize_pr_by_content(pr_title: str, pr_body: str, pr_labels: List[str]) -> str:
    """
    Categorize PR based on title keywords, body content, and labels.
    
    Returns: 'feature', 'enhancement', 'bug', 'breaking', 'docs', 'other'
    """
    title_lower = pr_title.lower()
    body_lower = (pr_body or '').lower()
    labels_lower = [label.lower() for label in pr_labels]
    
    # Check labels first (most reliable)
    if any(label in labels_lower for label in ['breaking-change', 'breaking']):
        return 'breaking'
    if any(label in labels_lower for label in ['feature']):
        return 'feature'
    if any(label in labels_lower for label in ['enhancement', 'improvement']):
        return 'enhancement'
    if any(label in labels_lower for label in ['bug', 'bugfix', 'fix']):
        return 'bug'
    if any(label in labels_lower for label in ['documentation', 'docs']):
        return 'docs'
    
    # Check title keywords
    if any(keyword in title_lower for keyword in ['breaking', 'break:']):
        return 'breaking'
    if any(keyword in title_lower for keyword in ['add support', 'implement', 'new resource', 'new feature']):
        return 'feature'
    if any(keyword in title_lower for keyword in ['enhance', 'improve', 'update']):
        return 'enhancement'
    if any(keyword in title_lower for keyword in ['fix', 'bug', 'resolve']):
        return 'bug'
    if any(keyword in title_lower for keyword in ['docs', 'documentation']):
        return 'docs'
    
    # Check body keywords
    if 'breaking change' in body_lower:
        return 'breaking'
    
    return 'other'


def detect_deprecations(pr_body: str, pr_files: List[Dict] = None) -> List[Dict[str, str]]:
    """
    Detect deprecated fields or resources mentioned in PR.
    
    Returns:
        List of dicts with 'field' and 'reason' keys
    """
    deprecations = []
    
    if not pr_body:
        return deprecations
    
    # Look for deprecation mentions in body
    deprecation_patterns = [
        r'deprecat(?:e|ing|ed)\s+`([^`]+)`',
        r'`([^`]+)`\s+is\s+(?:now\s+)?deprecated',
        r'removed?\s+`([^`]+)`',
    ]
    
    for pattern in deprecation_patterns:
        matches = re.finditer(pattern, pr_body, re.IGNORECASE)
        for match in matches:
            field = match.group(1)
            deprecations.append({
                'field': field,
                'context': match.group(0)
            })
    
    return deprecations


def validate_terraform_code(code: str) -> Tuple[bool, Optional[str]]:
    """
    Basic validation of Terraform code.
    
    Returns:
        (is_valid, error_message)
    """
    # Check for basic structure
    if not code.strip():
        return False, "Empty code block"
    
    # Check for resource/data block
    if not re.search(r'(resource|data)\s+"[^"]+"\s+"[^"]+"', code):
        return False, "No resource or data source declaration found"
    
    # Check for basic syntax issues
    if code.count('{') != code.count('}'):
        return False, "Mismatched braces"
    
    # Check if it references couchbase-capella
    if 'couchbase-capella' not in code:
        return False, "Does not reference couchbase-capella provider"
    
    return True, None


def format_resource_name(internal_name: str) -> str:
    """
    Convert internal resource name to Terraform resource name.
    
    Examples:
        "free_tier_cluster" -> "couchbase-capella_free_tier_cluster"
        "app_service" -> "couchbase-capella_app_service"
    """
    return f"couchbase-capella_{internal_name}"


def generate_registry_url(resource_name: str, is_datasource: bool = False) -> str:
    """
    Generate Terraform Registry URL for a resource or data source.
    
    Args:
        resource_name: Internal name (e.g., "free_tier_cluster")
        is_datasource: True if this is a data source, False for resource
    """
    base_url = "https://registry.terraform.io/providers/couchbasecloud/couchbase-capella/latest/docs"
    type_path = "data-sources" if is_datasource else "resources"
    return f"{base_url}/{type_path}/{resource_name}"


def generate_examples_url(resource_name: str) -> str:
    """
    Generate GitHub URL to examples folder for a resource.
    """
    base_url = "https://github.com/couchbasecloud/terraform-provider-couchbase-capella/tree/main/examples"
    return f"{base_url}/{resource_name}"

