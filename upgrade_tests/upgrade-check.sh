#!/usr/bin/env bash
#
# End-to-end provider upgrade check against a live Capella organisation.
#
# Installs a released provider version from the registry, applies a
# configuration, swaps in the provider built from the working tree, and re-plans.
# A non-empty plan at that point is an upgrade regression: practitioners who
# upgrade without touching their configuration would see the same diff.
#
# This is the slow, high-fidelity counterpart to the Go tests alongside it, which
# cover the state-decode and plan-merge path offline on every PR. Use this before a
# release, or when a change touches how state is read.
#
# Usage:
#   FROM_VERSION=1.9.1 \
#   TF_VAR_host=https://cloudapi.cloud.couchbase.com \
#   TF_VAR_auth_token=... \
#   TF_VAR_organization_id=... \
#   TF_VAR_project_id=... \
#   TF_VAR_cluster_id=... \
#   TF_VAR_bucket_name=default \
#   upgrade_tests/upgrade-check.sh
#
# Exit codes: 0 clean upgrade, 1 setup or apply failure, 2 upgrade regression.

set -euo pipefail

require() {
  local name=$1
  if [[ -z "${!name:-}" ]]; then
    echo "error: $name must be set" >&2
    exit 1
  fi
}

require FROM_VERSION
require TF_VAR_host
require TF_VAR_auth_token
require TF_VAR_organization_id
require TF_VAR_project_id
require TF_VAR_cluster_id

BUCKET_NAME=${TF_VAR_bucket_name:-default}
REPO_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
LOCAL_VERSION=99.0.0
PROVIDER_ADDR=registry.terraform.io/couchbasecloud/couchbase-capella
CRED_NAME="tf_upgrade_check_$RANDOM"

WORK=$(mktemp -d)
DESTROY_NEEDED=0

cleanup() {
  local rc=$?
  if [[ $DESTROY_NEEDED -eq 1 ]]; then
    echo
    echo "==> Destroying test resources"
    # Always destroy with the local provider: it is the one whose state is on disk.
    terraform -chdir="$WORK" destroy -auto-approve -no-color >/dev/null 2>&1 ||
      echo "warning: destroy failed; check for leftover credential $CRED_NAME" >&2
  fi
  rm -rf "$WORK"
  exit $rc
}
trap cleanup EXIT

# The provider reads these; TF_VAR_* are only wired up for the acceptance tests.
export CAPELLA_HOST="$TF_VAR_host"
export CAPELLA_AUTHENTICATION_TOKEN="$TF_VAR_auth_token"
export TF_IN_AUTOMATION=1 TF_INPUT=0 CHECKPOINT_DISABLE=1

write_config() {
  local version=$1
  cat > "$WORK/main.tf" <<EOF
terraform {
  required_providers {
    couchbase-capella = {
      source  = "couchbasecloud/couchbase-capella"
      version = "$version"
    }
  }
}

provider "couchbase-capella" {}

# Nested access at every depth, each list in non-canonical order so that any
# change in how these are stored or compared shows up as a diff.
resource "couchbase-capella_database_credential" "upgrade_check" {
  name            = "$CRED_NAME"
  organization_id = "$TF_VAR_organization_id"
  project_id      = "$TF_VAR_project_id"
  cluster_id      = "$TF_VAR_cluster_id"

  access = [
    {
      privileges = ["data_writer"]
      resources = {
        buckets = [
          {
            name = "$BUCKET_NAME"
            scopes = [
              { name = "_default" },
            ]
          },
        ]
      }
    },
    {
      privileges = ["data_reader"]
    },
  ]
}
EOF
}

echo "==> Phase 1: install couchbasecloud/couchbase-capella $FROM_VERSION from the registry"

# Phase 1 needs its own CLI config just as much as phase 2 does. Left unset,
# Terraform reads the user's ~/.terraformrc, and a dev_overrides block for this
# provider there - the normal way to develop against a local build - silently
# substitutes that build for the released one. Terraform downloads
# $FROM_VERSION, reports installing it, then runs the local binary anyway. The
# script would compare the working tree against itself and pass for the wrong
# reason, which is worse than failing.
cat > "$WORK/cli-registry.tfrc" <<EOF
provider_installation {
  direct {}
}
EOF
export TF_CLI_CONFIG_FILE="$WORK/cli-registry.tfrc"

write_config "$FROM_VERSION"
terraform -chdir="$WORK" init -no-color 2>&1 | tee "$WORK/phase1-init.log"

# Belt and braces: TF_CLI_CONFIG_FILE replaces the default config rather than
# merging with it, so the block above is enough. Assert it anyway, because a
# false pass here is invisible in the final output.
if grep -qi "development overrides" "$WORK/phase1-init.log"; then
  echo >&2
  echo "error: provider development overrides are in effect, so phase 1 would not" >&2
  echo "       actually run $FROM_VERSION. Refusing to report a meaningless result." >&2
  exit 1
fi

# Armed before the apply, not after: apply can create the credential and then
# fail, and with `set -e` a post-apply assignment would never run, so the EXIT
# trap would skip destroy and leak a live credential. Destroying an empty state
# is a no-op, so arming early costs nothing.
DESTROY_NEEDED=1
terraform -chdir="$WORK" apply -auto-approve -no-color

echo
# Informational only. `terraform version` reports the version selected in the
# lock file, and it still reports that version when dev_overrides has swapped in
# a different binary - so this line cannot prove which code ran. The grep above
# is the actual safeguard.
echo "==> Phase 1 provider selection (from the lock file):"
terraform -chdir="$WORK" version -no-color | grep -i couchbase || true

echo
echo "==> Phase 2: build the working tree and serve it as $LOCAL_VERSION"
PLATFORM="$(go env GOOS)_$(go env GOARCH)"
MIRROR="$WORK/mirror/$PROVIDER_ADDR/$LOCAL_VERSION/$PLATFORM"
mkdir -p "$MIRROR"
(cd "$REPO_ROOT" && go build -o "$MIRROR/terraform-provider-couchbase-capella_v$LOCAL_VERSION" .)

cat > "$WORK/cli.tfrc" <<EOF
provider_installation {
  filesystem_mirror {
    path    = "$WORK/mirror"
    include = ["$PROVIDER_ADDR"]
  }
  direct {
    exclude = ["$PROVIDER_ADDR"]
  }
}
EOF
export TF_CLI_CONFIG_FILE="$WORK/cli.tfrc"

# The lock file pins a checksum of the released binary, and the local build has a
# different one, so it has to go before re-initialising.
write_config "$LOCAL_VERSION"
rm -f "$WORK/.terraform.lock.hcl"
terraform -chdir="$WORK" init -upgrade -no-color

echo
echo "==> Phase 2 provider selection (from the lock file, expect v$LOCAL_VERSION):"
terraform -chdir="$WORK" version -no-color | grep -i couchbase || true

echo
echo "==> Phase 3: plan with the upgraded provider against unchanged configuration"
set +e
terraform -chdir="$WORK" plan -detailed-exitcode -no-color
PLAN_CODE=$?
set -e

case $PLAN_CODE in
  0)
    echo
    echo "PASS: upgrading from $FROM_VERSION to the working tree produces an empty plan."
    exit 0
    ;;
  2)
    echo
    echo "UPGRADE REGRESSION: the plan above is what a practitioner sees after upgrading"
    echo "from $FROM_VERSION without changing their configuration."
    echo
    echo "==> Phase 4: checking whether one apply converges"
    terraform -chdir="$WORK" apply -auto-approve -no-color
    set +e
    terraform -chdir="$WORK" plan -detailed-exitcode -no-color
    SECOND=$?
    set -e
    if [[ $SECOND -eq 0 ]]; then
      echo
      echo "The diff is self-healing: it clears after a single apply."
    else
      echo
      echo "The diff PERSISTS after apply - this is a permanent diff, not a one-off upgrade cost."
    fi
    exit 2
    ;;
  *)
    echo "error: terraform plan failed (exit $PLAN_CODE)" >&2
    exit 1
    ;;
esac
