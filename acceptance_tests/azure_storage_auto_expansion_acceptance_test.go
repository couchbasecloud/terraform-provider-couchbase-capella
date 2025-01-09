package acceptance_tests

import (
	"context"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccClusterResourceAzureDiskAutoExpansion(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	cidr, err := getCIDR(context.Background(), api.NewClient(timeout), "azure")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfigAzureDiskAutoExpansion(resourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateClusterImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Disable autoexpansion
			{
				Config: testAccClusterConfigAzureDiskAutoExpansionOff(resourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "false"),
				),
			},
		},
	})
}

func testAccClusterConfigAzureDiskAutoExpansion(resourceName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Terraform Acceptance Test Azure auto expansion"
  cloud_provider = {
    type   = "azure"
    region = "eastus"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
  			type          = "P6"
			autoexpansion = true
		}
      }
      num_of_nodes = 3
      services     = ["data"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, cidr)
}

func testAccClusterConfigAzureDiskAutoExpansionOff(resourceName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Terraform Acceptance Test Azure auto expansion"
  cloud_provider = {
    type   = "azure"
    region = "eastus"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
  			type          = "P6"
			autoexpansion = false
		}
      }
      num_of_nodes = 3
      services     = ["data"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, cidr)
}
