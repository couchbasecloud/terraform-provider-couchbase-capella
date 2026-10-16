package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccReplicationDataSourceReference = "data.couchbase-capella_replication.replication"

// TestAccReplicationDataSource reads a replication that the resource created, and checks that the data
// source has the same attributes as the resource.
func TestAccReplicationDataSource(t *testing.T) {
	requireReplicationAnchor(t)
	buckets := newReplicationTestBuckets(true)
	settings := replicationMappings + replicationFullFilter

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckReplicationClustersHealthy(),
		Steps: []resource.TestStep{
			{
				PreConfig: testAccPreConfigWaitForReplicationClusters(t),
				Config:    testAccReplicationBucketsConfig(buckets),
				Check:     testAccCheckReplicationClustersHealthy(),
			},
			{
				Config: testAccReplicationConfig(buckets, settings),
				Check:  testAccCheckReplicationClustersHealthy(),
			},
			{
				Config: testAccReplicationConfig(buckets, settings) + testAccReplicationDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckReplicationDataSourceMatchesResource(
						"id",
						"organization_id",
						"project_id",
						"cluster_id",
						"source_bucket",
						"target.cluster",
						"target.bucket",
						"target.type",
						"direction",
						"priority",
						"network_usage_limit",
						"mappings.#",
						"filter.document_exclude_options.deletion",
						"filter.document_exclude_options.expiration",
						"filter.document_exclude_options.ttl",
						"filter.document_exclude_options.binary",
						"filter.expressions.reg_ex",
						"mobile_sync_mode",
						"source_details.project.id",
						"source_details.project.name",
						"source_details.cluster.id",
						"source_details.cluster.name",
						"source_details.bucket.id",
						"source_details.bucket.name",
						"source_details.bucket.conflict_resolution_type",
						"target_details.project.id",
						"target_details.cluster.id",
						"target_details.cluster.name",
						"target_details.bucket.id",
						"target_details.bucket.name",
						"target_details.bucket.conflict_resolution_type",
						"target_details.type",
						"audit.created_by",
						"audit.created_at",
					),
					resource.TestCheckTypeSetElemNestedAttrs(testAccReplicationDataSourceReference, "mappings.*", map[string]string{
						"source_scope":                    "sales",
						"target_scope":                    "sales",
						"collections.#":                   "1",
						"collections.0.source_collection": "orders",
						"collections.0.target_collection": "orders",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(testAccReplicationDataSourceReference, "mappings.*", map[string]string{
						"source_scope": "inventory",
						"target_scope": "inventory",
					}),
					resource.TestCheckResourceAttrSet(testAccReplicationDataSourceReference, "status"),
					resource.TestCheckResourceAttrSet(testAccReplicationDataSourceReference, "changes_left"),
					resource.TestCheckNoResourceAttr(testAccReplicationDataSourceReference, "conflict_logging.%"),
					resource.TestCheckNoResourceAttr(testAccReplicationDataSourceReference, "filter.expressions.skip_restream"),
					resource.TestCheckNoResourceAttr(testAccReplicationDataSourceReference, "reverse_replication_id"),
				),
			},
		},
	})
}

// testAccCheckReplicationDataSourceMatchesResource checks that each attribute of the data source has the
// same value as in the replication resource.
func testAccCheckReplicationDataSourceMatchesResource(attributes ...string) resource.TestCheckFunc {
	checks := make([]resource.TestCheckFunc, 0, len(attributes))
	for _, attribute := range attributes {
		checks = append(checks, resource.TestCheckResourceAttrPair(
			testAccReplicationDataSourceReference, attribute, testAccReplicationReference, attribute,
		))
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func testAccReplicationDataSourceConfig() string {
	return fmt.Sprintf(`
data "couchbase-capella_replication" "replication" {
  organization_id = %[1]s.organization_id
  project_id      = %[1]s.project_id
  cluster_id      = %[1]s.cluster_id
  id              = %[1]s.id
}
`, testAccReplicationReference)
}
