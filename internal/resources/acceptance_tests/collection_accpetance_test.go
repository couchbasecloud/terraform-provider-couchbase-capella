package acceptance_tests

import (
	"context"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"net/http"
	"testing"
	"time"
)

func TestAccCollectionsTestCases(t *testing.T) {
	clusterResourceName := "new_cluster"
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	projectResourceName := "terraform_project"
	projectResourceReference := "couchbase-capella_project." + projectResourceName

	cidr, err := acctest.GetCIDR("aws")
	fmt.Println(cidr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	testCfg := acctest.ProjectCfg

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateCluster(&testCfg, clusterResourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(clusterResourceReference),
				),
			},
			//Create bucket
			{
				Config: testAccCreateBucketConfig(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.TestAccWait(time.Second * 10)),
				ExpectNonEmptyPlan: true,
			},
			//Create Scope
			{
				Config: testAccCreateMultiScopeConfigReference(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.TestAccWait(time.Second*10),
					resource.TestCheckResourceAttr("couchbase-capella_scope.new_scope1", "scope_name", "testScope1"),
					resource.TestCheckResourceAttrSet("couchbase-capella_scope.new_scope1", "cluster_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_scope.new_scope1", "project_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_scope.new_scope1", "organization_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_scope.new_scope1", "collections.#"),
					resource.TestCheckResourceAttr("couchbase-capella_scope.new_scope2", "scope_name", "testScope2"),
					resource.TestCheckResourceAttrSet("couchbase-capella_scope.new_scope2", "cluster_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_scope.new_scope2", "project_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_scope.new_scope2", "organization_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_scope.new_scope2", "collections.#"),
				),
				ExpectNonEmptyPlan: true,
			},

			//Create collections
			{
				Config: testAccCreateCollectionConfig(testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_collection.new_collection", "collection_name", "testCollection"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection", "bucket_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection", "cluster_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection", "project_id"),
					resource.TestCheckResourceAttr("couchbase-capella_collection.new_collection", "max_ttl", "200"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection", "organization_id"),
				),
				ExpectNonEmptyPlan: true,
			},

			////ImportState testing
			//{
			//	ResourceName:      "couchbase-capella_collection.new_collection",
			//	ImportStateIdFunc: generateCollectionDetailsForImport("couchbase-capella_collection.new_collection"),
			//	ImportState:       true,
			//},
			//Create multiple collections in multiple scopes
			{
				Config: testAccCreateMultiCollectionConfig(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_collection.new_collection1", "collection_name", "testCollection1"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection1", "bucket_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection1", "cluster_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection1", "project_id"),
					resource.TestCheckResourceAttr("couchbase-capella_collection.new_collection1", "max_ttl", "200"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection1", "organization_id"),
					resource.TestCheckResourceAttr("couchbase-capella_collection.new_collection2", "collection_name", "testCollection2"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection2", "bucket_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection2", "cluster_id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection2", "project_id"),
					resource.TestCheckResourceAttr("couchbase-capella_collection.new_collection2", "max_ttl", "200"),
					resource.TestCheckResourceAttrSet("couchbase-capella_collection.new_collection2", "organization_id"),
				),
				ExpectNonEmptyPlan: true,
			},

			//Delete collection
			{
				Config: testCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDeleteCollection("couchbase-capella_collection.new_collection1"),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
		},
	})
}

func testAccDeleteCollection(scopeResourceReference string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		var rawState map[string]string
		for _, module := range state.Modules {
			if len(module.Resources) > 0 {
				if val, ok := module.Resources[scopeResourceReference]; ok {
					rawState = val.Primary.Attributes
				}
			}
		}
		data, err := acctest.TestClient()
		if err != nil {
			return err
		}
		err = deleteCollection(data, rawState["organization_id"], rawState["project_id"], rawState["cluster_id"], rawState["bucket_id"], rawState["scope_name"], rawState["collection_name"])
		if err != nil {
			return err
		}
		fmt.Printf("collection deletion initiated")
		return nil
	}
}
func deleteCollection(data *providerschema.Data, orgId, projectId, clusterId, bucketId, scopeName, collectionName string) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s/collections/%s", data.HostURL, orgId, projectId, clusterId, bucketId, scopeName, collectionName)
	cfg := api.EndpointCfg{
		Url:           url,
		Method:        http.MethodDelete,
		SuccessStatus: http.StatusNoContent,
	}
	_, err := data.Client.ExecuteWithRetry(
		context.Background(),
		cfg,
		nil,
		data.Token,
		nil,
	)
	if err != nil {
		return err
	}
	//wait time to get the scope deleted
	time.Sleep(time.Second * 30)
	return nil
}

func generateCollectionDetailsForImport(collectoinResourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[collectoinResourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		fmt.Printf("raw state %s", rawState)
		str := fmt.Sprintf("collection_name=%s,scope_name=%s,bucket_id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["collection_name"], rawState["scope_name"], rawState["bucket_id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"])
		fmt.Printf("str %s", str)
		return str, nil
	}
}

func testAccCreateMultiCollectionConfig(cfg *string) string {
	*cfg = fmt.Sprintf(`%[1]s
	output "new_collection1" {
	  value = couchbase-capella_collection.new_collection1
	}
	
	resource "couchbase-capella_collection" "new_collection1" {
	  organization_id = var.organization_id
      project_id      = couchbase-capella_project.terraform_project.id
      cluster_id      = couchbase-capella_cluster.new_cluster.id
	  bucket_id       = couchbase-capella_bucket.new_bucket.id
	  scope_name      = couchbase-capella_scope.new_scope1.scope_name
	  collection_name = "testCollection1"
	  max_ttl         = "200"
	}
	output "new_collection2" {
	  value = couchbase-capella_collection.new_collection2
	}
	
	resource "couchbase-capella_collection" "new_collection2" {
	  organization_id = var.organization_id
	  project_id      = couchbase-capella_project.terraform_project.id
      cluster_id      = couchbase-capella_cluster.new_cluster.id
      bucket_id       = couchbase-capella_bucket.new_bucket.id
	  scope_name      = couchbase-capella_scope.new_scope2.scope_name
	  collection_name = "testCollection2"
	  max_ttl         = "200"
	}
	
`, *cfg)
	return *cfg

}

func testAccCreateCollectionConfig(cfg string) string {
	cfg = fmt.Sprintf(`
	%[1]s 
	output "new_collection" {
	  value = couchbase-capella_collection.new_collection
	}
	
	resource "couchbase-capella_collection" "new_collection" {
	  organization_id = var.organization_id
	  project_id      = couchbase-capella_project.terraform_project.id
      cluster_id      = couchbase-capella_cluster.new_cluster.id
	  bucket_id       = couchbase-capella_bucket.new_bucket.id
	  scope_name      = couchbase-capella_scope.new_scope1.scope_name
	  collection_name = "testCollection"
	  max_ttl         = "200"
	}
	
`, cfg)
	return cfg
}
