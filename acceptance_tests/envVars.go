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

	return nil
}
