#!/bin/bash
# Checks whether the environment is configured to run acceptance tests.
# Output is injected into the droid's session context at startup.

echo "=== Acceptance test environment ==="
echo ""

MISSING=0
for var in TF_VAR_host TF_VAR_auth_token TF_VAR_organization_id; do
  if [ -z "${!var}" ]; then
    echo "  $var: NOT SET"
    MISSING=1
  else
    echo "  $var: set"
  fi
done

echo ""

if [ "$MISSING" -eq 1 ]; then
  echo "Acceptance tests CANNOT be run this session — required credentials are missing."
  echo "Verification is limited to: go build ./acceptance_tests/"
  echo "Do not attempt to run tests or use make testacc."
else
  echo "Acceptance tests CAN be run this session."
  echo ""
  echo "Optional setup-skip variables (unset = TestMain creates the resource, ~15 min for cluster):"
  for var in TF_VAR_project_id TF_VAR_cluster_id TF_VAR_bucket_id TF_VAR_app_service_id; do
    if [ -z "${!var}" ]; then
      echo "  $var: NOT SET"
    else
      echo "  $var: set"
    fi
  done
fi
