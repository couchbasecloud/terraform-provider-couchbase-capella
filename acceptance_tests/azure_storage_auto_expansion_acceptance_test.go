package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterResourceAzureDiskAutoExpansion(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName

	// Generate a random CIDR to avoid conflicts with existing clusters
	cidr := generateRandomCIDR()

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

// TestAccClusterResourceAzureDiskAutoExpansionDefault validates AV-139797: Azure clusters
// must default autoexpansion to true when it is omitted from config, and toggling it off/on
// afterwards must not cause a regression (perpetual diff or apply failure).
func TestAccClusterResourceAzureDiskAutoExpansionDefault(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName

	// Generate a random CIDR to avoid conflicts with existing clusters
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// autoexpansion omitted from config: Azure must default it to true.
			{
				Config: testAccClusterConfigAzureDiskAutoExpansionUnset(resourceName, cidr),
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
			// Re-applying the same config must not produce a perpetual diff.
			{
				Config:   testAccClusterConfigAzureDiskAutoExpansionUnset(resourceName, cidr),
				PlanOnly: true,
			},
			// Turn autoexpansion off.
			{
				Config: testAccClusterConfigAzureDiskAutoExpansionOff(resourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "false"),
				),
			},
			// Turn autoexpansion back on.
			{
				Config: testAccClusterConfigAzureDiskAutoExpansion(resourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "true"),
				),
			},
		},
	})
}

// TestAccClusterResourceAzureDiskAutoExpansionOffOmittedOnUpdate covers the consequential
// direction of AV-143041's "Symptom 2": if an omitted autoExpansion were ever read as "enable"
// by Terraform the way it is by a raw API PUT, a cluster with auto-expansion explicitly off
// would get silently switched on - and on Azure that triggers an unrequested swap rebalance
// (node replacement and data movement). Starts from explicit autoexpansion=false, then forces
// an update that omits the field, and asserts it stays false.
func TestAccClusterResourceAzureDiskAutoExpansionOffOmittedOnUpdate(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName

	// Generate a random CIDR to avoid conflicts with existing clusters
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfigAzureDiskAutoExpansionOff(resourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "false"),
				),
			},
			// Update: change description (forces a real PUT) while dropping autoexpansion from
			// config. Must stay false - never silently switch on. Also assert the description
			// actually changed, proving the PUT/read cycle ran rather than returning stale state.
			{
				Config: testAccClusterConfigAzureDiskAutoExpansionOffOmittedUpdate(resourceName, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "false"),
					resource.TestCheckResourceAttr(resourceReference, "description", "Terraform Acceptance Test Azure auto expansion off update omitted"),
				),
			},
		},
	})
}

func testAccClusterConfigAzureDiskAutoExpansionOffOmittedUpdate(resourceName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Terraform Acceptance Test Azure auto expansion off update omitted"
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
          type = "P6"
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

func testAccClusterConfigAzureDiskAutoExpansionUnset(resourceName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Terraform Acceptance Test Azure auto expansion default"
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
          type = "P6"
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
