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
