#!/bin/bash
set -e

# Print timestamped status messages
echo "[$(date '+%Y-%m-%d %H:%M:%S')] Starting acceptance test script..."

# Load variables from terraform.tfvars
if [ ! -f terraform.tfvars ]; then
  echo "terraform.tfvars not found! Please create it in the project root."
  exit 1
fi

# Export variables from terraform.tfvars (robust to spaces and quotes)
while IFS='=' read -r key value; do
  # Remove whitespace and comments
  key=$(echo "$key" | xargs)
  value=$(echo "$value" | sed 's/["\r\n]*$//;s/^[ \t]*//;s/[ \t]*$//')
  # Remove inline comments
  value=$(echo "$value" | cut -d'#' -f1 | xargs)
  # Remove surrounding quotes
  value=$(echo "$value" | sed 's/^"//;s/"$//')
  if [[ -n "$key" && -n "$value" ]]; then
    export TF_VAR_${key}="${value}"
  fi
done < <(grep -v '^#' terraform.tfvars | grep '=')

echo "[$(date '+%Y-%m-%d %H:%M:%S')] Exported variables:"
env | grep TF_VAR_

# Accept test name(s) as arguments
if [ $# -gt 0 ]; then
  TEST_ARGS="-run $*"
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] Running acceptance tests matching: $*"
else
  TEST_ARGS=""
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] Running all acceptance tests..."
fi

# Run tests with progress output
echo "[$(date '+%Y-%m-%d %H:%M:%S')] make testacc GO_TEST_ARGS=\"$TEST_ARGS\""
TF_ACC=1 go test -timeout=120m -v ./acceptance_tests/ $TEST_ARGS | tee acceptance-test.log

EXIT_CODE=${PIPESTATUS[0]}
echo "[$(date '+%Y-%m-%d %H:%M:%S')] Acceptance tests finished with exit code $EXIT_CODE. Log saved to acceptance-test.log."
exit $EXIT_CODE
