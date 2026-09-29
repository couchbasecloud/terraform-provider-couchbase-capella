package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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

// TestAccClusterResourceAzureDiskAutoExpansionOffOmittedOnUpdate locks down the case where
// Terraform shields the provider from AV-143041's "Symptom 2" - an omitted autoExpansion
// being read as "enable" by a raw API PUT, silently switching on a cluster the operator
// configured with it off and triggering an unrequested Azure swap rebalance.
//
// Be clear about what this does and does not guard. Only `description` changes here, which
// lives outside the service_groups set, so the set element still correlates with prior state
// and Terraform supplies the previously stored `false` for the omitted Optional+Computed
// attribute. The provider therefore sends an explicit value and never reaches the server's
// omitted-field path - and it did so before AV-143040 was fixed too. This test passes against
// both the fixed and the unfixed provider; it is a lock on Terraform's Optional+Computed
// behaviour, not a regression guard for that fix. The PreApply check below states that
// explicitly by asserting the planned value is known-false rather than unknown.
//
// The case where correlation does NOT hold, and the provider really does omit the field, is
// covered in cluster_disk_unknown_acceptance_test.go.
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
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						// The point of the test: omitted from config, but planned as a known
						// false carried over from prior state - not unknown. That is why the
						// provider still sends it explicitly here.
						plancheck.ExpectKnownValue(
							resourceReference,
							tfjsonpath.New("service_groups").
								AtSliceIndex(0).
								AtMapKey("node").
								AtMapKey("disk").
								AtMapKey("autoexpansion"),
							knownvalue.Bool(false),
						),
					},
				},
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
