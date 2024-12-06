package acceptance_tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
)

func TestAccGSITestCases(t *testing.T) {
	resourceName := "new_cluster"
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "terraform_project"
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr := "10.0.9.0/24"
	//if err != nil {
	//	t.Error(err)
	//	t.FailNow()
	//}

	testCfg := acctest.ProjectCfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			//Creating cluster to check the gsi index creation
			{
				Config: testAccCreateCluster(&testCfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					// importSampleBucket(resourceReference),
				),
			},
			{
				Config: SampleBucketWithTravelSample(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.TestAccWait(time.Minute * 2)),
			},
			{
				Config: testAccCreateGSINonDeferredIndexConfig(testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "index_name", "index1"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "index_keys.0", "c1"),
					//resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "with.defer_build", "false"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "with.num_replica", "1"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "with.num_partition", "24"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "where", "geo.alt > 1000"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "partition_by.0", "meta().id"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "bucket_name", "travel-sample"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "scope_name", "_default"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "collection_name", "_default"),
					resource.TestCheckResourceAttrSet("couchbase-capella_query_indexes.non_deferred_index", "organization_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_query_indexes.non_deferred_index", "project_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_query_indexes.non_deferred_index", "cluster_id"),
				),
			},
			{
				Config:                               testAccCreateGSINonDeferredIndexConfig(testCfg),
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
			{
				Config: testAccCreateGSIDeferredIndexConfig(testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.deferred_index", "index_name", "index2"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.deferred_index", "index_keys.0", "c2"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.deferred_index", "with.defer_build", "true"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.deferred_index", "with.num_replica", "1"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.deferred_index", "with.num_partition", "24"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.deferred_index", "where", "geo.alt > 1000"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.deferred_index", "partition_by.0", "meta().id"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.deferred_index", "bucket_name", "travel-sample"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.deferred_index", "scope_name", "_default"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.deferred_index", "collection_name", "_default"),
					resource.TestCheckResourceAttrSet("couchbase-capella_query_indexes.deferred_index", "organization_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_query_indexes.deferred_index", "project_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_query_indexes.deferred_index", "cluster_id"),
				),
			},

			{

				Config: testCfg + "\n" + GsiMultiNonDeferredIndexfile,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.multi_non_deferred_index.instances.0", "index_name", "idx_non_deferred1"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.multi_non_deferred_index.instances.0", "index_keys.0", "field1"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.multi_non_deferred_index.instances.99", "index_name", "idx_non_deferred99"),
				),
			},
		},
	})
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
		fmt.Printf("raw state %s", rawState)
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

// project_id      = var.project_id
// cluster_id      = var.cluster_id
func multiNonDeferredIndexConfig(cfg string) string {
	cfg = fmt.Sprintf(`
%[1]s

locals {
  index_template = templatefile("/Users/saicharanrachamadugu/repos2/terraform-provider-couchbase-capella/indexes_with_non_deferred_build.json",{
  organization_id = ""
  project_id      = ""
  cluster_id      = ""
})
}


resource "couchbase-capella_query_indexes" "multi_non_deferred_index" {
  for_each        = jsondecode(local.index_template).resource["couchbase-capella_indexes"]
  organization_id = ""
  project_id      = ""
  cluster_id      = ""
  bucket_name     = each.value.bucket_name
  scope_name      = each.value.scope_name
  collection_name = each.value.collection_name
  index_name      = each.value.index_name
  index_keys      = each.value.index_keys
}

`, cfg)
	return cfg
}

// project_id      = couchbase-capella_project.terraform_project.id
// cluster_id      = couchbase-capella_cluster.new_cluster.id
// organization_id = var.organization_id
func testAccCreateGSINonDeferredIndexConfig(cfg string) string {
	cfg = fmt.Sprintf(`
%[1]s

output "create_gsi_index" {
  value = couchbase-capella_query_indexes.non_deferred_index
}

resource "couchbase-capella_query_indexes" "non_deferred_index" {
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  organization_id = var.organization_id
  bucket_name     = "travel-sample"
  scope_name      = "_default"
  collection_name = "_default"
  index_name      = "index1"
  index_keys      = ["c1"]
  partition_by   = ["meta().id"]
  where = "geo.alt > 1000"
  with ={
        num_replica = 1
        num_partition = 24
  }

}

`, cfg)
	fmt.Println(cfg)
	return cfg
}

func testAccCreateGSIDeferredIndexConfig(cfg string) string {
	cfg = fmt.Sprintf(`
%[1]s

output "create_gsi_index" {
  value = couchbase-capella_query_indexes.deferred_index
}

resource "couchbase-capella_query_indexes" "deferred_index" {
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  organization_id = var.organization_id
  bucket_name     = "travel-sample"
  scope_name      = "_default"
  collection_name = "_default"
  index_name      = "index2"
  index_keys      = ["c2"]
  partition_by   = ["meta().id"]
  where = "geo.alt > 1000"
  with ={
        num_replica = 1
        num_partition = 24
		defer_build = true
  }

}

`, cfg)
	return cfg
}

func SampleBucketWithTravelSample(cfg *string) string {
	*cfg = fmt.Sprintf(`
%[1]s
output "sample_bucket_travel"{
  value = couchbase-capella_sample_bucket.add_sample_bucket_travel
}
resource "couchbase-capella_sample_bucket" "add_sample_bucket_travel" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name			  = "travel-sample"
}
`, *cfg)
	return *cfg
}
