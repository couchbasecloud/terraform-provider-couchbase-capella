#!/bin/bash
# Script to run tests from a list file
# Usage: ./scripts/run-test-list.sh acceptance_tests/sanity.list

set -e
set -o pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ACC_DIR="$REPO_ROOT/acceptance_tests"

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

# Resolve the list to an absolute path before moving to the repo root, so a
# relative argument keeps working from any directory.
TEST_LIST_FILE="$(cd "$(dirname "$TEST_LIST_FILE")" && pwd)/$(basename "$TEST_LIST_FILE")"
cd "$REPO_ROOT"

# Check required environment variables
if [ -z "$TF_VAR_auth_token" ] || [ -z "$TF_VAR_host" ] || [ -z "$TF_VAR_organization_id" ]; then
    echo "ERROR: Required environment variables not set"
    echo "Please export: TF_VAR_auth_token, TF_VAR_host, TF_VAR_organization_id"
    exit 1
fi

# Read test names from file (skip comments and empty lines, tolerate CRLF and
# stray indentation, drop duplicates while preserving file order).
# The `|| true` keeps a comments-only list from tripping `set -e` via pipefail,
# so the empty-list check below reports it instead of exiting silently.
TEST_NAMES="$(sed -e 's/\r$//' -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' "$TEST_LIST_FILE" \
    | grep -v '^#' | grep -v '^$' | awk '!seen[$0]++' || true)"

if [ -z "$TEST_NAMES" ]; then
    echo "Error: No tests found in $TEST_LIST_FILE"
    exit 1
fi

EXPECTED_COUNT="$(printf '%s\n' "$TEST_NAMES" | grep -c . || true)"
TEST_PATTERN="$(printf '%s\n' "$TEST_NAMES" | tr '\n' '|' | sed 's/|$//')"

# Fail on names that match no test function.
#
# The pattern below is anchored, so a typo or a renamed/deleted test makes
# `go test -run` match nothing, print "no tests to run" and still exit 0 - a
# green gate that exercised nothing. Validate up front instead.
#
# `go test -list` is deliberately not used for this: TestMain provisions real
# Capella resources before the testing framework handles -list, so listing would
# cost a full setup cycle. Reading the declarations costs nothing.
DECLARED_TESTS="$(grep -rhoE '^func Test[A-Za-z0-9_]+\(' "$ACC_DIR"/*.go \
    | sed -E 's/^func (Test[A-Za-z0-9_]+)\(.*/\1/' | sort -u)"

UNKNOWN_NAMES=""
while IFS= read -r test_name; do
    if ! printf '%s\n' "$DECLARED_TESTS" | grep -qxF "$test_name"; then
        UNKNOWN_NAMES="${UNKNOWN_NAMES}  ${test_name}"$'\n'
    fi
done <<< "$TEST_NAMES"

if [ -n "$UNKNOWN_NAMES" ]; then
    echo "ERROR: $TEST_LIST_FILE names tests that do not exist in acceptance_tests/:"
    echo ""
    printf '%s' "$UNKNOWN_NAMES"
    echo ""
    echo "These were renamed, deleted or misspelled. Left in place they would"
    echo "match nothing and the run would pass without executing them."
    exit 1
fi

# Matches the timeout used by `make testacc`. The sanity list provisions real
# clusters and app services, so a short timeout panics mid-run rather than
# failing a test.
TEST_TIMEOUT="${TEST_TIMEOUT:-180m}"

echo "=========================================="
echo "Running Sanity Tests from: $TEST_LIST_FILE"
echo "=========================================="
echo "Tests requested: $EXPECTED_COUNT"
echo "Test pattern: $TEST_PATTERN"
echo "Timeout: $TEST_TIMEOUT"
echo ""

TEST_LOG="$(mktemp)"
trap 'rm -f "$TEST_LOG"' EXIT

# Run the tests
set +e
CAPELLA_OPENAPI_SPEC_PATH="${CAPELLA_OPENAPI_SPEC_PATH:-$REPO_ROOT/openapi.generated.yaml}" \
TF_ACC=1 \
go test -timeout="${TEST_TIMEOUT}" -v ./acceptance_tests/ -run "^(${TEST_PATTERN})$" 2>&1 | tee "$TEST_LOG"
TEST_STATUS="${PIPESTATUS[0]}"
set -e

# Assert the run actually executed what was asked for. `go test` reports success
# for an empty selection, so exit status alone cannot be trusted here.
RAN_COUNT="$(grep -cE '^=== RUN[[:space:]]+Test[A-Za-z0-9_]+$' "$TEST_LOG" || true)"

echo ""
echo "=========================================="
echo "Tests requested: $EXPECTED_COUNT   executed: $RAN_COUNT"
echo "=========================================="

# A non-zero status is a genuine test or setup failure; the log already explains
# it, and the counts below would only add a misleading second diagnosis.
if [ "$TEST_STATUS" -ne 0 ]; then
    exit "$TEST_STATUS"
fi

if grep -q 'no tests to run' "$TEST_LOG"; then
    echo "ERROR: go test reported \"no tests to run\" - the selection matched nothing."
    exit 1
fi

if [ "$RAN_COUNT" -eq 0 ]; then
    echo "ERROR: no tests executed. Refusing to report success for an empty run."
    exit 1
fi

if [ "$RAN_COUNT" -ne "$EXPECTED_COUNT" ]; then
    echo "ERROR: go test succeeded but executed $RAN_COUNT of $EXPECTED_COUNT requested tests."
    echo "Some listed tests were never run (build tags, or excluded from the package)."
    exit 1
fi

exit 0
