package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGSI(t *testing.T) {
	resourceType := "couchbase-capella_query_indexes"
	primaryIndexResourceName := randomStringWithPrefix("tf_acc_gsi_")
	secondaryIndexResourceName := randomStringWithPrefix("tf_acc_gsi_")
	primaryIndexResourceReference := fmt.Sprintf("%s.%s", resourceType, primaryIndexResourceName)
	secondaryIndexResourceReference := fmt.Sprintf("%s.%s", resourceType, secondaryIndexResourceName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateGSINonDeferredIndexConfig(secondaryIndexResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(secondaryIndexResourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(secondaryIndexResourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(secondaryIndexResourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(secondaryIndexResourceReference, "bucket_name", globalBucketName),
					resource.TestCheckResourceAttr(secondaryIndexResourceReference, "scope_name", globalScopeName),
					resource.TestCheckResourceAttr(secondaryIndexResourceReference, "collection_name", globalCollectionName),
					resource.TestCheckResourceAttr(secondaryIndexResourceReference, "index_name", "index1"),
					resource.TestCheckResourceAttr(secondaryIndexResourceReference, "index_keys.0", "c1"),
					resource.TestCheckResourceAttr(secondaryIndexResourceReference, "where", "geo.alt > 1000"),
					resource.TestCheckResourceAttr(secondaryIndexResourceReference, "with.num_replica", "1"),
				),
			},
			{
				ResourceName:      secondaryIndexResourceReference,
				ImportStateIdFunc: generateGsiImportIdForResource(secondaryIndexResourceReference),
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

func testAccCreateGSINonDeferredIndexConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "%[6]s"
  collection_name = "%[7]s"
  index_name      = "index1"
  index_keys      = ["c1"]
  where = "geo.alt > 1000"
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
