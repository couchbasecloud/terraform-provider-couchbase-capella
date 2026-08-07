package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccFreeTierClusters_AV_134714 verifies that the free-tier clusters
// datasource lists free-tier clusters via the dedicated /clusters/freeTier
// endpoint instead of filtering the regular clusters list, which always
// produced an empty (and previously state-breaking) result. See AV-134714.
func TestAccFreeTierClusters_AV_134714(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_free_tier_clusters_ds_")
	resourceReference := "data.couchbase-capella_free_tier_clusters." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeTierClustersDatasourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(resourceReference, "data.#"),
				),
			},
		},
	})
}

func testAccFreeTierClustersDatasourceConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_free_tier_clusters" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName)
}
