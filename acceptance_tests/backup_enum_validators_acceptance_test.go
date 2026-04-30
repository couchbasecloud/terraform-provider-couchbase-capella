package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccBackupRestoreReplaceTTLEnumValidator verifies that the restore.replace_ttl
// field rejects out-of-enum values at plan time and accepts valid ones.
func TestAccBackupRestoreReplaceTTLEnumValidator(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_backup_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Rejected: invalid replace_ttl value
			{
				Config:      testAccBackupWithReplaceTTL(resourceName, "invalid_ttl"),
				ExpectError: regexp.MustCompile("value must be one of:"),
			},
			// Accepted: valid replace_ttl value — plan only, no backup deployed
			{
				Config:             testAccBackupWithReplaceTTL(resourceName, "none"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccBackupWithReplaceTTL(resourceName, replaceTTL string) string {
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
