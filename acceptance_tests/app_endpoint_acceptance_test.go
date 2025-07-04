package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// This test is not parallel so it's done before all the other tests.
// The reason is buckets cannot be deleted while app service is being deleted.
func TestAppEndpointResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAppEndpointResourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "bucket", "test_bucket"),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					//resource.TestCheckResourceAttr(resourceReference, "app_service_id", "${couchbase-capella_app_service." + resourceName + ".id}"),
				),
			},
			// ImportState testing
			//{
			//	ResourceName:      resourceReference,
			//	ImportStateIdFunc: generateAppServiceImportId(resourceReference),
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAppEndpointResourceConfig(resourceName string) string {
	clusterName := randomStringWithPrefix("tf_acc_cluster_")
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[5]s"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "10.190.250.0/23"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 50
          type    = "gp3"
          iops    = 3000
        }
      }
      num_of_nodes = 1
      services     = ["data", "index", "query"]
    }
  ]
  availability = {
    "type" : "single"
  }
  support = {
    plan     = "developer pro"
    timezone = "PT"
  }
}

resource "couchbase-capella_app_service" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_cluster.%[5]s.id
  name            = "tf_acc_test_app_service"
  compute = {
    cpu = 2
    ram = 4
  }
}

resource "couchbase-capella_app_endpoint" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_cluster.%[5]s.id
  app_service_id  = couchbase-capella_app_service.%[4]s.id
  bucket          = "test_bucket"
  name            = "%[6]s"
  
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalAppServiceId, clusterName, globalAppEndpoint)
}

//func generateAppImportId(resourceReference string) resource.ImportStateIdFunc {
//	return func(state *terraform.State) (string, error) {
//		var rawState map[string]string
//		for _, m := range state.Modules {
//			if len(m.Resources) > 0 {
//				if v, ok := m.Resources[resourceReference]; ok {
//					rawState = v.Primary.Attributes
//				}
//			}
//		}
//		return fmt.Sprintf("id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
//	}
//}
