package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGSI(t *testing.T) {
	const resourceType = "couchbase-capella_query_indexes"
	primaryIndexResourceName := randomStringWithPrefix("tf_acc_gsi_")
	primaryIndexResourceReference := fmt.Sprintf("%s.%s", resourceType, primaryIndexResourceName)
	gsiResourceIdx1 := fmt.Sprintf("%s.%s", resourceType, "idx1")
	gsiResourceIdx2 := fmt.Sprintf("%s.%s", resourceType, "idx2")
	gsiResourceIdx3 := fmt.Sprintf("%s.%s", resourceType, "idx3")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateGSINonDeferredIndexConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(gsiResourceIdx1, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "index_name", "idx1"),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "index_keys.0", "c1"),
					resource.TestCheckResourceAttr(gsiResourceIdx1, "where", "geo.alt > 1000"),

					resource.TestCheckResourceAttr(gsiResourceIdx2, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "index_name", "idx2"),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "index_keys.0", "c2"),
					resource.TestCheckResourceAttr(gsiResourceIdx2, "where", "geo.alt > 2000"),

					resource.TestCheckResourceAttr(gsiResourceIdx3, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "index_name", "idx3"),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "index_keys.0", "c3"),
					resource.TestCheckResourceAttr(gsiResourceIdx3, "where", "geo.alt > 3000"),
				),
			},
			{
				ResourceName:      gsiResourceIdx1,
				ImportStateIdFunc: generateGsiImportIdForResource(gsiResourceIdx1),
				ImportState:       true,
			},
			{
				Config: testAccCreatePrimaryIndexConfig(primaryIndexResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "index_name", "primary_index"),
					resource.TestCheckResourceAttr(primaryIndexResourceReference, "with.num_replica", "1"),
				),
			},
		},
	})
}

func testAccCreateGSINonDeferredIndexConfig() string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "idx1" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "idx1"
  index_keys      = ["c1"]
  where = "geo.alt > 1000"
  with = {
        defer_build = false
  }
}

resource "couchbase-capella_query_indexes" "idx2" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "idx2"
  index_keys      = ["c2"]
  where = "geo.alt > 2000"
  with = {
        defer_build = false
  }
}

resource "couchbase-capella_query_indexes" "idx3" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "idx3"
  index_keys      = ["c3"]
  where = "geo.alt > 3000"
  with = {
        defer_build = false
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalBucketName,
		globalScopeName,
		globalCollectionName,
	)
}

func testAccCreatePrimaryIndexConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "primary_index"
  is_primary      = true
  with = {
        num_replica = 1
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalBucketName,
		globalScopeName,
		globalCollectionName,
		resourceName)
}

// TestAccGSIImportWithoutIsPrimary verifies that importing an existing index
// without specifying is_primary in the config does not trigger a spurious
// in-place update. The is_primary attribute is computed, so Terraform uses the
// imported value and produces no diff.
func TestAccGSIImportWithoutIsPrimary(t *testing.T) {
	const resourceType = "couchbase-capella_query_indexes"
	indexResourceName := "idx"
	indexResourceReference := fmt.Sprintf("%s.%s", resourceType, indexResourceName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "import_test_idx"
  index_keys      = ["import_test_key"]
}
`, globalProviderBlock,
					globalOrgId,
					globalProjectId,
					globalClusterId,
					globalBucketName,
					globalScopeName,
					globalCollectionName,
					indexResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(indexResourceReference, "index_name", "import_test_idx"),
					resource.TestCheckResourceAttr(indexResourceReference, "index_keys.0", "import_test_key"),
				),
			},
			{
				ResourceName:            indexResourceReference,
				ImportStateIdFunc:       generateGsiImportIdForResource(indexResourceReference),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"with"},
			},
		},
	})
}

// TestAccGSINumReplicaUnchangedNoop verifies that re-applying the same config
// after an index is created does not trigger a failing ALTER INDEX call when
// num_replica has not changed. Previously the provider issued an unconditional
// ALTER INDEX on every in-place update, which Couchbase Server rejects when
// num_replica is unchanged ("Index already has X number of replica").
func TestAccGSINumReplicaUnchangedNoop(t *testing.T) {
	const resourceType = "couchbase-capella_query_indexes"
	indexResourceName := "idx"
	indexResourceReference := fmt.Sprintf("%s.%s", resourceType, indexResourceName)

	configWithNumReplica0 := fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "noop_test_idx"
  index_keys      = ["noop_test_key"]
  with = {
    num_replica = 0
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalBucketName,
		globalScopeName,
		globalCollectionName,
		indexResourceName)

	configWithNumReplica1 := fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "noop_test_idx"
  index_keys      = ["noop_test_key"]
  with = {
    num_replica = 1
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalBucketName,
		globalScopeName,
		globalCollectionName,
		indexResourceName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: Create index with num_replica = 0
			{
				Config: configWithNumReplica0,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(indexResourceReference, "index_name", "noop_test_idx"),
					resource.TestCheckResourceAttr(indexResourceReference, "with.num_replica", "0"),
				),
			},
			// Step 2: Re-apply the same config — should be a no-op, not fail
			{
				Config: configWithNumReplica0,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(indexResourceReference, "index_name", "noop_test_idx"),
					resource.TestCheckResourceAttr(indexResourceReference, "with.num_replica", "0"),
				),
			},
			// Step 3: Update num_replica from 0 to 1
			{
				Config: configWithNumReplica1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(indexResourceReference, "index_name", "noop_test_idx"),
					resource.TestCheckResourceAttr(indexResourceReference, "with.num_replica", "1"),
				),
			},
			// Step 4: Re-apply with num_replica = 1 — should be a no-op, not fail
			{
				Config: configWithNumReplica1,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(indexResourceReference, "index_name", "noop_test_idx"),
					resource.TestCheckResourceAttr(indexResourceReference, "with.num_replica", "1"),
				),
			},
		},
	})
}

func generateGsiImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf(
			"index_name=%s,collection_name=%s,scope_name=%s,bucket_name=%s,cluster_id=%s,organization_id=%s,project_id=%s",
			rawState["index_name"],
			rawState["collection_name"],
			rawState["scope_name"],
			rawState["bucket_name"],
			rawState["cluster_id"],
			rawState["organization_id"],
			rawState["project_id"],
		), nil
	}
}
