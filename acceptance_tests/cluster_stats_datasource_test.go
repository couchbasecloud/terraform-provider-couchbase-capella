package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceClusterStats(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cluster_stats_ds_")
	dsReference := "data.couchbase-capella_cluster_stats." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterStatsDataSourceConfig(dsName, globalClusterId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrSet(dsReference, "free_memory_in_mb"),
					resource.TestCheckResourceAttrSet(dsReference, "max_replicas"),
					resource.TestCheckResourceAttrSet(dsReference, "total_memory_in_mb"),
				),
			},
		},
	})
}

func TestAccDatasourceClusterStatsInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cluster_stats_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterStatsDataSourceConfig(dsName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile(`(?s)Error.*[Cc]luster|cluster.*not found|access to the requested resource is denied`),
			},
		},
	})
}

func TestAccDatasourceClusterStatsMissingCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cluster_stats_missing_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cluster_stats" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)cluster_id|argument.*required`),
			},
		},
	})
}

func testAccClusterStatsDataSourceConfig(dsName, clusterID string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_cluster_stats" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, clusterID)
}
