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

func TestAccScopeTestCases(t *testing.T) {
	clusterResourceName := "new_cluster"
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	projectResourceName := "terraform_project"
	projectResourceReference := "couchbase-capella_project." + projectResourceName

	scopeResourceName := "new_scope"
	scopeResourceReference := "couchbase-capella_scope." + scopeResourceName
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
				Config: testAccCreateScopeConfig(testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					acctest.TestAccWait(time.Second*10),
					resource.TestCheckResourceAttr(scopeResourceReference, "scope_name", "testScope"),
					resource.TestCheckResourceAttrSet(scopeResourceReference, "cluster_id"),
					resource.TestCheckResourceAttrSet(scopeResourceReference, "project_id"),
					resource.TestCheckResourceAttrSet(scopeResourceReference, "organization_id"),
					resource.TestCheckResourceAttrSet(scopeResourceReference, "collections.#"),
				),
				ExpectNonEmptyPlan: true,
			},
			// ImportState testing
			//{
			//	ResourceName:      scopeResourceReference,
			//	ImportStateIdFunc: generateScopeDetailsForImport(scopeResourceReference),
			//	ImportState:       true,
			//},

			//Create Multiple Scope
			{
				Config: testAccCreateMultiScopeConfigReference(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
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

			//List Scopes
			{
				Config: testAccListScopesConfig(testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.couchbase-capella_scopes.existing_scopes", "bucket_id"),
					resource.TestCheckResourceAttrSet("data.couchbase-capella_scopes.existing_scopes", "scopes.1.scope_name"),
					resource.TestCheckResourceAttrSet("data.couchbase-capella_scopes.existing_scopes", "scopes.0.scope_name"),
					resource.TestCheckResourceAttrSet("data.couchbase-capella_scopes.existing_scopes", "scopes.2.scope_name"),
				),
			},

			{
				Config: testCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDeleteScope("couchbase-capella_scope.new_scope1"),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
		},
	})
}

func testAccDeleteScope(scopeResourceReference string) resource.TestCheckFunc {
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
		err = deleteScope(data, rawState["organization_id"], rawState["project_id"], rawState["cluster_id"], rawState["bucket_id"], rawState["scope_name"])
		if err != nil {
			return err
		}
		fmt.Printf("scope deleted")
		return nil
	}
}

func deleteScope(data *providerschema.Data, orgId, projectId, clusterId, bucketId, scopeName string) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s", data.HostURL, orgId, projectId, clusterId, bucketId, scopeName)
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

func testAccListScopesConfig(cfg string) string {
	cfg = fmt.Sprintf(`
	%[1]s
	output "scopes_list" {
	  value = data.couchbase-capella_scopes.existing_scopes
	}
	
	data "couchbase-capella_scopes" "existing_scopes" {
	  organization_id = var.organization_id
	  project_id      = couchbase-capella_project.terraform_project.id
      cluster_id      = couchbase-capella_cluster.new_cluster.id
      bucket_id       = couchbase-capella_bucket.new_bucket.id
	}
`, cfg)
	return cfg
}
func testAccCreateScopeConfig(cfg string) string {
	cfg = fmt.Sprintf(`
%[1]s
output "new_scope" {
  value = couchbase-capella_scope.new_scope
}

resource "couchbase-capella_scope" "new_scope" {
  scope_name      = "testScope"
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  bucket_id       = couchbase-capella_bucket.new_bucket.id
}
`, cfg)
	return cfg

}

func testAccCreateMultiScopeConfigReference(cfg *string) string {
	*cfg = fmt.Sprintf(`
%[1]s
output "new_scope1" {
  value = couchbase-capella_scope.new_scope1
}

resource "couchbase-capella_scope" "new_scope1" {
  scope_name      = "testScope1"
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  bucket_id       = couchbase-capella_bucket.new_bucket.id
}
output "new_scope2" {
  value = couchbase-capella_scope.new_scope2
}

resource "couchbase-capella_scope" "new_scope2" {
  scope_name      = "testScope2"
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  bucket_id       = couchbase-capella_bucket.new_bucket.id
}
`, *cfg)
	return *cfg

}

func testAccListScope(cfg string) string {
	cfg = fmt.Sprintf(`
	%[1]s
	output "scopes_list" {
	  value = data.couchbase-capella_scopes.existing_scopes
	}
	
	data "couchbase-capella_scopes" "existing_scopes" {
	  organization_id = var.organization_id
	  project_id      = couchbase-capella_project.terraform_project.id
      cluster_id      = couchbase-capella_cluster.new_cluster.id
	  bucket_id       = couchbase-capella_bucket.new_bucket.id
	}
`, cfg)
	return cfg
}

func testAccCreateBucketConfig(cfg *string) string {
	*cfg = fmt.Sprintf(`
	%[1]s
	output "new_bucket" {
  		value = couchbase-capella_bucket.new_bucket
	}

	output "bucket_id" {
  		value = couchbase-capella_bucket.new_bucket.id
	}
	resource "couchbase-capella_bucket" "new_bucket" {
		name                       = "acceptanceTestBucket"
		organization_id            = var.organization_id
        project_id      = couchbase-capella_project.terraform_project.id
        cluster_id      = couchbase-capella_cluster.new_cluster.id
		type                       = "couchbase"
		storage_backend            = "couchstore"
		memory_allocation_in_mb    = 100
		bucket_conflict_resolution = "seqno"
		durability_level           = "none"
		replicas                   = 1
		flush                      = false
		time_to_live_in_seconds    = 0
		eviction_policy            = "fullEviction"
	}
`, *cfg)
	return *cfg
}

func generateScopeDetailsForImport(scopeResourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[scopeResourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		fmt.Printf("raw state %s", rawState)
		return fmt.Sprintf("scope_name=%s,bucket_id=%s,organization_id=%s,project_id=%s,cluster_id=%s", rawState["scope_name"], rawState["bucket_id"], rawState["organization_id"], rawState["project_id"], rawState["cluster_id"]), nil
	}
}
