package acceptance_tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
)

func TestAccGSITestCases(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateGSINonDeferredIndexConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "index_name", "index1"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "index_keys.0", "c1"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "with.num_replica", "1"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "where", "geo.alt > 1000"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "bucket_name", "default"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "scope_name", "_default"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "collection_name", "_default"),
					resource.TestCheckResourceAttrSet("couchbase-capella_query_indexes.non_deferred_index", "organization_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_query_indexes.non_deferred_index", "project_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_query_indexes.non_deferred_index", "cluster_id"),
				),
			},
			{
				Config:                               testAccCreateGSINonDeferredIndexConfig(),
				ResourceName:                         "couchbase-capella_query_indexes.non_deferred_index",
				ImportStateIdFunc:                    generateGSIImportIdForResource("couchbase-capella_query_indexes.non_deferred_index"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "index_name",
				ImportStateVerifyIgnore: []string{
					"where",
					"index_keys",
					"partition_by",
					"with",
					"is_primary",
				},
			},
		},
	})
}

// project_id      = couchbase-capella_project.terraform_project.id
// cluster_id      = couchbase-capella_cluster.new_cluster.id
// organization_id = var.organization_id
func testAccCreateGSINonDeferredIndexConfig() string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "non_deferred_index" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "default"
  scope_name      = "_default"
  collection_name = "_default"
  index_name      = "index1"
  index_keys      = ["c1"]
  where = "geo.alt > 1000"
  with ={
        num_replica = 1
  }
}
`, ProviderBlock, OrgId, ProjectId, ClusterId)
}

func generateGSIImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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
