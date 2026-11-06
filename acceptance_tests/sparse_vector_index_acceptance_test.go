package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccSparseVectorIndexDefaultTunables covers SPV-01/SPV-25: a default-tunable
// SPARSE VECTOR index, created via the existing couchbase-capella_query_indexes
// resource with zero provider code changes, builds and comes online.
func TestAccSparseVectorIndexDefaultTunables(t *testing.T) {
	requireSparseVectorCluster(t)

	const resourceType = "couchbase-capella_query_indexes"
	resourceReference := fmt.Sprintf("%s.%s", resourceType, "sv_idx")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccSparseVectorIndexConfig("sv_idx", "tf_acc_sv_idx"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", sparseVectorClusterId),
					resource.TestCheckResourceAttr(resourceReference, "bucket_name", sparseVectorBucketName),
					resource.TestCheckResourceAttr(resourceReference, "scope_name", sparseVectorScopeName),
					resource.TestCheckResourceAttr(resourceReference, "collection_name", sparseVectorCollectionName),
					resource.TestCheckResourceAttr(resourceReference, "index_name", "tf_acc_sv_idx"),
					resource.TestCheckResourceAttr(resourceReference, "index_keys.0", "`sv` SPARSE VECTOR"),
					resource.TestCheckResourceAttr(resourceReference, "status", "Ready"),
				),
			},
		},
	})
}

// TestAccSparseVectorIndexMultiple covers SPV-03: two sparse vector indexes on
// the same collection coexist without conflict.
func TestAccSparseVectorIndexMultiple(t *testing.T) {
	requireSparseVectorCluster(t)

	const resourceType = "couchbase-capella_query_indexes"
	idx1Reference := fmt.Sprintf("%s.%s", resourceType, "sv_idx_multi1")
	idx2Reference := fmt.Sprintf("%s.%s", resourceType, "sv_idx_multi2")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccSparseVectorIndexMultipleConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(idx1Reference, "status", "Ready"),
					resource.TestCheckResourceAttr(idx2Reference, "status", "Ready"),
				),
			},
		},
	})
}

func testAccSparseVectorIndexConfig(resourceName, indexName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  is_primary      = false
  index_name      = "%[9]s"
  index_keys      = ["`+"`"+`sv`+"`"+` SPARSE VECTOR"]
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		sparseVectorClusterId,
		sparseVectorBucketName,
		sparseVectorScopeName,
		sparseVectorCollectionName,
		resourceName,
		indexName,
	)
}

func testAccSparseVectorIndexMultipleConfig() string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "sv_idx_multi1" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  is_primary      = false
  index_name      = "tf_acc_sv_idx_multi1"
  index_keys      = ["`+"`"+`sv`+"`"+` SPARSE VECTOR"]
}

resource "couchbase-capella_query_indexes" "sv_idx_multi2" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  is_primary      = false
  index_name      = "tf_acc_sv_idx_multi2"
  index_keys      = ["`+"`"+`sv`+"`"+` SPARSE VECTOR"]
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		sparseVectorClusterId,
		sparseVectorBucketName,
		sparseVectorScopeName,
		sparseVectorCollectionName,
	)
}
