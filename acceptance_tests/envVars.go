package acceptance_tests

import (
	"os"
)

func getEnvVars() error {
	globalHost = os.Getenv("TF_VAR_host")
	if globalHost == "" {
		return ErrHostMissing
	}
	globalToken = os.Getenv("TF_VAR_auth_token")
	if globalToken == "" {
		return ErrTokenMissing
	}
	globalOrgId = os.Getenv("TF_VAR_organization_id")
	if globalOrgId == "" {
		return ErrOrgIdMissing
	}

	// Optional: Use existing resources instead of creating new ones
	globalProjectId = os.Getenv("TF_VAR_project_id")
	globalClusterId = os.Getenv("TF_VAR_cluster_id")
	globalAppServiceId = os.Getenv("TF_VAR_app_service_id")
	globalBucketId = os.Getenv("TF_VAR_bucket_id")

	// Optional dedicated cluster for snapshot/restore tests so that long
	// restore windows don't interfere with the shared cluster used by other
	// tests. Falls back to the primary cluster/bucket when unset.
	globalSnapshotClusterId = os.Getenv("TF_VAR_snapshot_cluster_id")
	globalSnapshotBucketId = os.Getenv("TF_VAR_snapshot_bucket_id")

	return nil
}

func snapshotClusterId() string {
	if globalSnapshotClusterId != "" {
		return globalSnapshotClusterId
	}
	return globalClusterId
}

func snapshotBucketId() string {
	if globalSnapshotBucketId != "" {
		return globalSnapshotBucketId
	}
	return globalBucketId
}
