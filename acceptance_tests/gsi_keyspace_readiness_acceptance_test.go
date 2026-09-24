package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccGSI_AV_133415 covers CBSE-22110: creating an index in the same apply as its
// bucket used to fail with `httpStatusCode: 424, message: "Keyspace not found in CB
// datastore"`. POST /v4/.../buckets returns 201 as soon as ns_server accepts the
// bucket, but the query nodes serve index DDL from an eventually-consistent metadata
// cache, so a CREATE INDEX issued immediately afterwards lands before the keyspace is
// visible and the apply fails.
//
// The config deliberately contains no time_sleep and no readiness poll: the indexes
// reach the bucket only through a bucket_name reference, so Terraform fires the index
// DDL the moment the bucket create returns — the exact window the customer hit. The
// bucket is ephemeral to match their automation.
//
// The secondary index is chained behind the primary rather than created alongside it:
// the indexer rejects concurrent builds on a keyspace with a warning, which would
// leave a partially-populated resource and make the assertions flap. Serializing keeps
// the assertion focused on the keyspace-propagation window, which only the first DDL
// after the bucket create can exercise.
func TestAccGSI_AV_133415(t *testing.T) {
	bucketName := randomStringWithPrefix("tf_acc_ks_bkt_")
	primaryIdxName := randomStringWithPrefix("tf_acc_ks_pidx_")
	secondaryIdxName := randomStringWithPrefix("tf_acc_ks_sidx_")

	bucketReference := "couchbase-capella_bucket." + bucketName
	primaryReference := "couchbase-capella_query_indexes." + primaryIdxName
	secondaryReference := "couchbase-capella_query_indexes." + secondaryIdxName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccGSIBucketAndIndexInOneApplyConfig(bucketName, primaryIdxName, secondaryIdxName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(bucketReference, "name", bucketName),
					resource.TestCheckResourceAttr(bucketReference, "type", "ephemeral"),
					resource.TestCheckResourceAttr(bucketReference, "cluster_id", globalClusterId),

					resource.TestCheckResourceAttr(primaryReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(primaryReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(primaryReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(primaryReference, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(primaryReference, "scope_name", "_default"),
					resource.TestCheckResourceAttr(primaryReference, "collection_name", "_default"),
					resource.TestCheckResourceAttr(primaryReference, "index_name", primaryIdxName),
					resource.TestCheckResourceAttr(primaryReference, "is_primary", "true"),
					resource.TestCheckResourceAttrSet(primaryReference, "status"),

					resource.TestCheckResourceAttr(secondaryReference, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(secondaryReference, "index_name", secondaryIdxName),
					resource.TestCheckResourceAttr(secondaryReference, "index_keys.0", "c1"),
					resource.TestCheckResourceAttrSet(secondaryReference, "status"),
				),
			},
		},
	})
}

// TestAccGSI_AV_133415_NewCollection is the same race one level down: a bucket, scope,
// collection and index all created in a single apply. Scope and collection creation is
// also accepted asynchronously, so the query service can be behind on any part of the
// keyspace, not just the bucket. This is the shape most real configurations take.
func TestAccGSI_AV_133415_NewCollection(t *testing.T) {
	bucketName := randomStringWithPrefix("tf_acc_kc_bkt_")
	scopeName := randomStringWithPrefix("tf_acc_kc_scp_")
	collectionName := randomStringWithPrefix("tf_acc_kc_col_")
	idxName := randomStringWithPrefix("tf_acc_kc_idx_")

	idxReference := "couchbase-capella_query_indexes." + idxName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccGSINewCollectionAndIndexInOneApplyConfig(bucketName, scopeName, collectionName, idxName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(idxReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(idxReference, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(idxReference, "scope_name", scopeName),
					resource.TestCheckResourceAttr(idxReference, "collection_name", collectionName),
					resource.TestCheckResourceAttr(idxReference, "index_name", idxName),
					resource.TestCheckResourceAttr(idxReference, "index_keys.0", "c1"),
					resource.TestCheckResourceAttrSet(idxReference, "status"),
				),
			},
		},
	})
}

func testAccGSIBucketAndIndexInOneApplyConfig(bucketName, primaryIdxName, secondaryIdxName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
  type            = "ephemeral"
  eviction_policy = "nruEviction"
}

resource "couchbase-capella_query_indexes" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = couchbase-capella_bucket.%[2]s.name
  scope_name      = "_default"
  collection_name = "_default"
  index_name      = "%[6]s"
  is_primary      = true
}

resource "couchbase-capella_query_indexes" "%[7]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = couchbase-capella_bucket.%[2]s.name
  scope_name      = "_default"
  collection_name = "_default"
  index_name      = "%[7]s"
  index_keys      = ["c1"]
  with = {
    defer_build = false
  }

  depends_on = [couchbase-capella_query_indexes.%[6]s]
}
`, globalProviderBlock, bucketName, globalOrgId, globalProjectId, globalClusterId,
		primaryIdxName, secondaryIdxName)
}

func testAccGSINewCollectionAndIndexInOneApplyConfig(bucketName, scopeName, collectionName, idxName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  name            = "%[2]s"
}

resource "couchbase-capella_scope" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = couchbase-capella_bucket.%[2]s.id
  scope_name      = "%[6]s"
}

resource "couchbase-capella_collection" "%[7]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = couchbase-capella_bucket.%[2]s.id
  scope_name      = "%[6]s"
  collection_name = "%[7]s"

  depends_on = [couchbase-capella_scope.%[6]s]
}

resource "couchbase-capella_query_indexes" "%[8]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_name     = couchbase-capella_bucket.%[2]s.name
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "%[8]s"
  index_keys      = ["c1"]
  with = {
    defer_build = false
  }

  depends_on = [couchbase-capella_collection.%[7]s]
}
`, globalProviderBlock, bucketName, globalOrgId, globalProjectId, globalClusterId,
		scopeName, collectionName, idxName)
}
