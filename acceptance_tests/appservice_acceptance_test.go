package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// This test is not parallel so it's done before all the other tests.
// The reason is buckets cannot be deleted while app service is being deleted.
func TestAppServiceResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_")
	resourceReference := "couchbase-capella_app_service." + resourceName
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAppServiceResourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "description"),
					resource.TestCheckResourceAttr(resourceReference, "compute.cpu", "2"),
					resource.TestCheckResourceAttr(resourceReference, "compute.ram", "4"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppServiceImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAppServiceResourceConfig(resourceName string) string {
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
      services     = ["data"]
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
  cluster_id      = "couchbase-capella_cluster.%[5]s.id"
  name            = "tf_acc_test_app_service"
  compute = {
    cpu = 2
    ram = 4
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, clusterName)
}

func generateAppServiceImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf("id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
	}
}
