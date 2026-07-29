package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccGSIIncludeMissingImport reproduces CBSE-23077 / AV-134985: importing a
// query index whose key uses INCLUDE MISSING must preserve the modifier in
// index_keys, otherwise the next plan forces a destroy-and-recreate. Skipped
// until AV-134985 (import reads the CREATE INDEX DDL) lands on main; un-skip to
// reproduce the bug or to verify the fix.
func TestAccGSIIncludeMissingImport(t *testing.T) {
	t.Skip("CBSE-23077/AV-134985: import drops INCLUDE MISSING from index_keys; un-skip once the import fix lands")

	const includeMissingKey = "`name` INCLUDE MISSING"
	name := randomStringWithPrefix("tf_acc_gsi_incmiss_")
	ref := "couchbase-capella_query_indexes." + name

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccGSIIncludeMissingConfig(name, includeMissingKey),
				Check:  resource.TestCheckResourceAttr(ref, "index_keys.0", includeMissingKey),
			},
			{
				ResourceName:      ref,
				ImportStateIdFunc: generateGsiImportIdForResource(ref),
				ImportState:       true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported state, got %d", len(states))
					}
					if got := states[0].Attributes["index_keys.0"]; got != includeMissingKey {
						return fmt.Errorf(
							"AV-134985: import dropped INCLUDE MISSING; index_keys.0 = %q, want %q",
							got, includeMissingKey,
						)
					}
					return nil
				},
			},
		},
	})
}

func testAccGSIIncludeMissingConfig(resourceName, indexKey string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "%[8]s"
  index_keys      = [%[9]q]
  with = {
    defer_build = false
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, globalBucketName, globalScopeName, globalCollectionName, resourceName, indexKey)
}
