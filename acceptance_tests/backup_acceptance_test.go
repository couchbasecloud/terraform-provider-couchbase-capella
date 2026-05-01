package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccBackupEnumValidators_AV_129336 verifies that the schema-level enum validator
// rejects out-of-range values at plan time for the backup resource's restore.replace_ttl field.
func TestAccBackupEnumValidators_AV_129336(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_validators_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Invalid replace_ttl — rejected by schema validator before reaching the API
			{
				Config:      testAccBackupConfigWithInvalidReplaceTTL(resourceName),
				ExpectError: regexp.MustCompile(`value must be one of:`),
			},
			// Valid replace_ttl = "none" — validator passes at plan time; PlanOnly avoids
			// checking post-apply state since Create does not persist the restore block.
			{
				Config:   testAccBackupConfigWithValidReplaceTTL(resourceName, "none"),
				PlanOnly: true,
			},
		},
	})
}

func testAccBackupConfigWithInvalidReplaceTTL(resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"

  restore = {
    target_cluster_id = "%[5]s"
    source_cluster_id = "%[5]s"
    services          = ["data"]
    replace_ttl       = "invalid_ttl"
  }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalBucketId)
}

func testAccBackupConfigWithValidReplaceTTL(resourceName, replaceTTL string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_backup" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"

  restore = {
    target_cluster_id = "%[5]s"
    source_cluster_id = "%[5]s"
    services          = ["data"]
    replace_ttl       = "%[7]s"
  }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalBucketId, replaceTTL)
}
