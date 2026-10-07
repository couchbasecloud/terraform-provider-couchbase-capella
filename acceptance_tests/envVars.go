package acceptance_tests

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// serverVersionPattern matches the MAJOR.MINOR form the cluster resource keeps in state.
// The provider strips the patch component from the version the API reports, so a
// patch-level value would never match state. See the check in getEnvVars.
var serverVersionPattern = regexp.MustCompile(`^\d+\.\d+$`)

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

	// TF_VAR_server_version pins the Couchbase Server version for every cluster the
	// suite creates. It is also picked up by Terraform itself as the value for
	// var.server_version, which the cluster HCL in the tests references.
	globalServerVersion = os.Getenv("TF_VAR_server_version")
	if globalServerVersion == "" {
		globalServerVersion = defaultServerVersion
	}
	// couchbase_server is RequiresReplace and the provider stores only MAJOR.MINOR, so a
	// patch-level value here would never match state: every plan after create would be
	// non-empty and would propose replacing the cluster. Fail loudly instead.
	if !serverVersionPattern.MatchString(globalServerVersion) {
		return fmt.Errorf("%w  Got %q", ErrInvalidServerVersion, globalServerVersion)
	}

	// ACC_SKIP_APP_SERVICE skips the shared app service + app endpoint setup in
	// TestMain (see setup). Accepts standard bool forms (1/true/...); anything
	// unparseable, including unset, leaves it false.
	globalSkipAppService, _ = strconv.ParseBool(os.Getenv("ACC_SKIP_APP_SERVICE"))

	return nil
}
