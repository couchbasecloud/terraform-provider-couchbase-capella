#!/bin/bash
# Script to run tests from a list file
# Usage: ./scripts/run-test-list.sh acceptance_tests/sanity.list

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <test-list-file>"
    echo "Example: $0 acceptance_tests/sanity.list"
    exit 1
fi

TEST_LIST_FILE="$1"

if [ ! -f "$TEST_LIST_FILE" ]; then
    echo "Error: Test list file not found: $TEST_LIST_FILE"
    exit 1
fi

# Check required environment variables
if [ -z "$TF_VAR_auth_token" ] || [ -z "$TF_VAR_host" ] || [ -z "$TF_VAR_organization_id" ]; then
    echo "ERROR: Required environment variables not set"
    echo "Please export: TF_VAR_auth_token, TF_VAR_host, TF_VAR_organization_id"
    exit 1
fi

# Read test names from file (skip comments and empty lines)
TEST_PATTERN=$(grep -v '^#' "$TEST_LIST_FILE" | grep -v '^[[:space:]]*$' | tr '\n' '|' | sed 's/|$//')

if [ -z "$TEST_PATTERN" ]; then
    echo "Error: No tests found in $TEST_LIST_FILE"
    exit 1
fi

echo "=========================================="
echo "Running Sanity Tests from: $TEST_LIST_FILE"
echo "=========================================="
echo "Test pattern: $TEST_PATTERN"
echo ""

# Run the tests
CAPELLA_OPENAPI_SPEC_PATH="$(pwd)/openapi.generated.yaml" \
TF_ACC=1 \
go test -timeout=30m -v ./acceptance_tests/ -run "^(${TEST_PATTERN})$"
