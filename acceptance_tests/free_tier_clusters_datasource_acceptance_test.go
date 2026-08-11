package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccFreeTierClustersDatasourceEmptyResult verifies that the free-tier
// clusters data source always sets the data attribute, even when no free-tier
// cluster matches. It previously left data as a nil slice, which lands in state
// as a null list and fails "Attribute 'data.#' expected to be set". See
// AV-134714.
//
// Listing free-tier clusters themselves is not asserted here: the Capella v4
// API has no free-tier cluster list endpoint (only the free-tier buckets have
// one), so the data source can only filter the general /clusters list.
func TestAccFreeTierClustersDatasourceEmptyResult(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_free_tier_clusters_ds_")
	dsRef := "data.couchbase-capella_free_tier_clusters." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeTierClustersDatasourceConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsRef, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsRef, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(dsRef, "data.#"),
				),
			},
		},
	})
}

func testAccFreeTierClustersDatasourceConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_free_tier_clusters" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
}
`, globalProviderBlock, globalOrgId, globalProjectId, dsName)
}
