package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccGSI(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_gsi_")
	resourceReference := "couchbase-capella_query_indexes." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateGSINonDeferredIndexConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "index_name", "index1"),
					resource.TestCheckResourceAttr(resourceReference, "index_keys.0", "c1"),
					resource.TestCheckResourceAttr(resourceReference, "with.num_replica", "1"),
					resource.TestCheckResourceAttr(resourceReference, "where", "geo.alt > 1000"),
					resource.TestCheckResourceAttr(resourceReference, "bucket_name", "default"),
					resource.TestCheckResourceAttr(resourceReference, "scope_name", "_default"),
					resource.TestCheckResourceAttr(resourceReference, "collection_name", "_default"),
					resource.TestCheckResourceAttrSet(resourceReference, "organization_id"),
					resource.TestCheckResourceAttrSet(resourceReference, "project_id"),
					resource.TestCheckResourceAttrSet(resourceReference, "cluster_id"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateGSIImportIdForResource(resourceReference),
				ImportState:       true,
			},
		},
	})
}

func testAccCreateGSINonDeferredIndexConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_query_indexes" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_name     = "%[5]s"
  scope_name      = "_default"
  collection_name = "_default"
  index_name      = "index1"
  index_keys      = ["c1"]
  where = "geo.alt > 1000"
  with ={
        num_replica = 1
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, globalBucketName, resourceName)
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
