package acceptance_tests

import (
	"fmt"
	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccGSITestCases(t *testing.T) {
	resourceName := "new_cluster"
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "terraform_project"
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr, err := acctest.GetCIDR("aws")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	testCfg := acctest.ProjectCfg
	//cfg := ""
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Creating cluster to check the gsi index creation
			{
				Config: testAccCreateCluster(&testCfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
				),
			},
			{
				Config: testAccCreateGSINonDeferredIndexConfig(testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "index_name", "index1"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.non_deferred_index", "index_keys.0", "c1"),
					resource.TestCheckResourceAttr("couchbase-capella_query_indexes.deferred_index", "with.defer_build", "false"),
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
				Config: testAccCreateGSIDeferredIndexConfig(&testCfg),
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
		},
	})
}

// project_id      = couchbase-capella_project.terraform_project.id
// cluster_id      = couchbase-capella_cluster.new_cluster.id
func testAccCreateGSINonDeferredIndexConfig(cfg string) string {
	cfg = fmt.Sprintf(`
%[1]s

output "create_gsi_index" {
  value = couchbase-capella_query_indexes.non_deferred_index
}

resource "couchbase-capella_query_indexes" "non_deferred_index" {
  organization_id = var.organization_id
  project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  cluster_id      = "a1b11035-d1d3-4e2f-a856-7b9c63c62276"
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
	return cfg
}

func testAccCreateGSIDeferredIndexConfig(cfg *string) string {
	*cfg = fmt.Sprintf(`
%[1]s

output "create_gsi_index" {
  value = couchbase-capella_query_indexes.deferred_index
}

resource "couchbase-capella_query_indexes" "deferred_index" {
  organization_id = var.organization_id
  project_id      = "c1fade1a-9f27-4a3c-af73-d1b2301890e3"
  cluster_id      = "a1b11035-d1d3-4e2f-a856-7b9c63c62276"
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

`, *cfg)
	return *cfg
}
