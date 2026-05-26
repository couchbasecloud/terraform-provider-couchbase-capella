package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccClusterDeletionProtectionResource tests the full lifecycle:
// create with protection disabled → update to enabled → import state.
func TestAccClusterDeletionProtectionResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_del_prot_")
	resourceReference := "couchbase-capella_cluster_deletion_protection." + resourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: Create with deletion_protection = false
			{
				Config: testAccClusterDeletionProtectionConfig(resourceName, globalClusterId, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "deletion_protection", "false"),
				),
			},
			// Step 2: Update from false to true
			{
				Config: testAccClusterDeletionProtectionConfig(resourceName, globalClusterId, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "deletion_protection", "true"),
				),
			},
			// Step 3: ImportState
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateDeletionProtectionImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "cluster_id",
			},
			// Step 4: Disable to leave cluster in clean state
			{
				Config: testAccClusterDeletionProtectionConfig(resourceName, globalClusterId, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "deletion_protection", "false"),
				),
			},
		},
	})
}

// TestAccClusterDeletionProtectionToggle verifies that deletion protection
// can be toggled back and forth and the state is correctly reflected.
func TestAccClusterDeletionProtectionToggle(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_del_prot_toggle_")
	resourceReference := "couchbase-capella_cluster_deletion_protection." + resourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Enable protection
			{
				Config: testAccClusterDeletionProtectionConfig(resourceName, globalClusterId, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "deletion_protection", "true"),
				),
			},
			// Disable protection
			{
				Config: testAccClusterDeletionProtectionConfig(resourceName, globalClusterId, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "deletion_protection", "false"),
				),
			},
			// Re-enable protection
			{
				Config: testAccClusterDeletionProtectionConfig(resourceName, globalClusterId, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "deletion_protection", "true"),
				),
			},
			// Disable to leave cluster in clean state
			{
				Config: testAccClusterDeletionProtectionConfig(resourceName, globalClusterId, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "deletion_protection", "false"),
				),
			},
		},
	})
}

// TestAccClusterDeletionProtectionInvalidCluster verifies correct error for nonexistent cluster.
func TestAccClusterDeletionProtectionInvalidCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_del_prot_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterDeletionProtectionConfig(resourceName, "00000000-0000-0000-0000-000000000000", true),
				ExpectError: regexp.MustCompile(`(?s)Error.*cluster|not found|access.*denied`),
			},
		},
	})
}

// --- Config builders ---

func testAccClusterDeletionProtectionConfig(resourceName, clusterID string, enabled bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_deletion_protection" "%[2]s" {
  organization_id     = "%[3]s"
  project_id          = "%[4]s"
  cluster_id          = "%[5]s"
  deletion_protection = %[6]t
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, clusterID, enabled)
}

// --- Import ID ---

func generateDeletionProtectionImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		if rawState == nil {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		return fmt.Sprintf(
			"cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}
