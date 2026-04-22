#!/bin/bash
# Checks whether the environment is configured to run acceptance tests.
# Output is injected into the droid's session context at startup.

lines=()
lines+=("=== Acceptance test environment ===")
lines+=("")

MISSING=0
for var in TF_VAR_host TF_VAR_auth_token TF_VAR_organization_id; do
  if [ -z "${!var}" ]; then
    lines+=("  $var: NOT SET")
    MISSING=1
  else
    lines+=("  $var: set")
  fi
done

lines+=("")

if [ "$MISSING" -eq 1 ]; then
  lines+=("Acceptance tests CANNOT be run this session — required credentials are missing.")
  lines+=("Verification is limited to: go build ./acceptance_tests/")
  lines+=("Do not attempt to run tests or use make testacc.")
else
  lines+=("Acceptance tests CAN be run this session.")
  lines+=("")
  lines+=("Optional setup-skip variables (unset = TestMain creates the resource, ~15 min for cluster):")
  for var in TF_VAR_project_id TF_VAR_cluster_id TF_VAR_bucket_id TF_VAR_app_service_id; do
    if [ -z "${!var}" ]; then
      lines+=("  $var: NOT SET")
    else
      lines+=("  $var: set")
    fi
  done
fi

# Join lines into a single string with newlines, then emit as Factory additionalContext JSON
output=$(printf '%s\n' "${lines[@]}")
printf '{"additionalContext": %s}' "$(printf '%s' "$output" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))')"
