package acceptance_tests

import (
	"os"
	"strconv"
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
	if bucketName := os.Getenv("TF_VAR_bucket_name"); bucketName != "" {
		globalBucketName = bucketName
	}
	dmClusterId = os.Getenv("TF_VAR_dm_cluster_id")

	// ACC_SKIP_APP_SERVICE skips the shared app service + app endpoint setup in
	// TestMain (see setup). Accepts standard bool forms (1/true/...); anything
	// unparseable, including unset, leaves it false.
	globalSkipAppService, _ = strconv.ParseBool(os.Getenv("ACC_SKIP_APP_SERVICE"))

	return nil
}
