package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceFreeTierClusters(t *testing.T) {
	clusterName := randomStringWithPrefix("tf_acc_free_tier_clusters_ds_")
	dsName := randomStringWithPrefix("tf_acc_free_tier_clusters_ds_")
	clusterReference := "couchbase-capella_free_tier_cluster." + clusterName
	dsReference := "data.couchbase-capella_free_tier_clusters." + dsName
	cidr := generateRandomCIDR()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeTierClustersDatasourceClusterConfig(clusterName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterReference, "name", clusterName),
					resource.TestCheckResourceAttrSet(clusterReference, "id"),
				),
			},
			{
				Config: testAccFreeTierClustersDatasourceConfig(clusterName, dsName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterReference, "name", clusterName),
					resource.TestCheckResourceAttrSet(clusterReference, "id"),
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
					testAccCheckListElemNestedAttrs(dsReference, "data", map[string]string{
						"name":            clusterName,
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
						"support.plan":    "free",
					}),
				),
			},
		},
	})
}

func testAccFreeTierClustersDatasourceClusterConfig(clusterName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster" "%[4]s" {
	organization_id = "%[2]s"
	project_id      = "%[3]s"
	name            = "%[4]s"
	description     = "Free tier cluster for clusters data source acceptance testing."

	cloud_provider = {
		type   = "aws"
		region = "us-east-2"
		cidr   = "%[5]s"
	}
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterName, cidr)
}

func testAccFreeTierClustersDatasourceConfig(clusterName, dsName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Free tier cluster for clusters data source acceptance testing."

  cloud_provider = {
    type   = "aws"
    region = "us-east-2"
    cidr   = "%[6]s"
  }
}

data "couchbase-capella_free_tier_clusters" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"

  depends_on = [couchbase-capella_free_tier_cluster.%[4]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterName, dsName, cidr)
}
