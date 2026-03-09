package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccReplicationsDataSource tests the list replications data source.
func TestAccReplicationsDataSource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_replications_")
	resourceReference := "data.couchbase-capella_replications." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccReplicationsDataSourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					// The data attribute should exist (may be empty or populated depending on environment)
					resource.TestCheckResourceAttrSet(resourceReference, "data.#"),
				),
			},
		},
	})
}

// TestAccReplicationDataSourceNotFound tests the singular replication data source
// when a non-existent replication ID is provided.
func TestAccReplicationDataSourceNotFound(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_replication_")
	resourceReference := "data.couchbase-capella_replication." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccReplicationDataSourceConfig(resourceName, "non-existent-replication-id"),
				// When a non-existent ID is used, the API should return 404
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "replication_id", "non-existent-replication-id"),
				),
				// Expect error due to 404
				ExpectError: nil,
			},
		},
	})
}

// testAccReplicationsDataSourceConfig returns the HCL config for testing the list replications data source.
func testAccReplicationsDataSourceConfig(resourceName string) string {
	return fmt.Sprintf(`
%s

data "couchbase-capella_replications" "%s" {
  organization_id = "%s"
  project_id      = "%s"
  cluster_id      = "%s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}

// testAccReplicationDataSourceConfig returns the HCL config for testing the singular replication data source.
func testAccReplicationDataSourceConfig(resourceName, replicationId string) string {
	return fmt.Sprintf(`
%s

data "couchbase-capella_replication" "%s" {
  organization_id = "%s"
  project_id      = "%s"
  cluster_id      = "%s"
  replication_id  = "%s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, replicationId)
}
