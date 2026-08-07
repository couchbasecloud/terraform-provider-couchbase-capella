package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)


var malformedLoadBalancerCIDRs = []string{
	"192.168.0.0/244",    // prefix length above 32
	"999.999.999.999/24", // octets out of range
	"not-a-cidr",         // not an address
	"2001:db8::/24",      // IPv6, not IPv4
	"192.168.0.0",        // missing prefix length
}

func TestAccAppServiceLoadBalancerCIDR(t *testing.T) {
	const (
		loadBalancerCIDR        = "192.168.0.0/24"
		loadBalancerCIDRChanged = "192.168.1.0/24"
	)

	resourceName := randomStringWithPrefix("tf_acc_app_svc_lb_")
	resourceReference := "couchbase-capella_app_service." + resourceName
	clusterName := randomStringWithPrefix("tf_acc_cluster_azure_")
	dataSourceName := randomStringWithPrefix("tf_acc_app_svcs_lb_ds_")
	dataSourceReference := "data.couchbase-capella_app_services." + dataSourceName
	appServiceName := randomStringWithPrefix("tf_acc_app_svc_lb_")
	clusterCIDR := generateRandomCIDR()

	steps := make([]resource.TestStep, 0, len(malformedLoadBalancerCIDRs)+4)


	for _, cidr := range malformedLoadBalancerCIDRs {
		steps = append(steps, resource.TestStep{
			Config: testAccAppServiceLoadBalancerCIDRMalformedConfig(resourceName, clusterName, appServiceName, clusterCIDR, cidr),
			ExpectError: regexp.MustCompile(`(?s)error during app service creation`),
		})
	}

	steps = append(steps,
		resource.TestStep{
			Config: testAccAppServiceLoadBalancerCIDRConfig(resourceName, clusterName, dataSourceName, appServiceName, clusterCIDR, loadBalancerCIDR, 2),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "load_balancer_cidr", loadBalancerCIDR),
				resource.TestCheckResourceAttr(resourceReference, "cloud_provider", "Azure"),
				resource.TestCheckResourceAttr(resourceReference, "name", appServiceName),
				resource.TestCheckResourceAttr(resourceReference, "nodes", "2"),
				resource.TestCheckResourceAttr(resourceReference, "compute.cpu", "2"),
				resource.TestCheckResourceAttr(resourceReference, "compute.ram", "4"),
				resource.TestCheckResourceAttrSet(resourceReference, "id"),
				resource.TestCheckResourceAttrSet(resourceReference, "cluster_id"),
				resource.TestCheckResourceAttrSet(resourceReference, "current_state"),
				resource.TestCheckResourceAttrSet(resourceReference, "version"),
				resource.TestCheckResourceAttrSet(resourceReference, "etag"),
				testAccCheckAppServicesDataSourceLoadBalancerCIDR(dataSourceReference, resourceReference, loadBalancerCIDR),
			),
		},
		resource.TestStep{
			ResourceName:      resourceReference,
			ImportStateIdFunc: generateAppServiceImportId(resourceReference),
			ImportState:       true,
			ImportStateVerify: true,
		},
		resource.TestStep{
			Config: testAccAppServiceLoadBalancerCIDRConfig(resourceName, clusterName, dataSourceName, appServiceName, clusterCIDR, loadBalancerCIDR, 3),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "nodes", "3"),
				resource.TestCheckResourceAttr(resourceReference, "load_balancer_cidr", loadBalancerCIDR),
				resource.TestCheckResourceAttr(resourceReference, "name", appServiceName),
				resource.TestCheckResourceAttr(resourceReference, "compute.cpu", "2"),
				resource.TestCheckResourceAttr(resourceReference, "compute.ram", "4"),
			),
		},
		resource.TestStep{
			Config: testAccAppServiceLoadBalancerCIDRConfig(resourceName, clusterName, dataSourceName, appServiceName, clusterCIDR, loadBalancerCIDRChanged, 3),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionReplace),
				},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "load_balancer_cidr", loadBalancerCIDRChanged),
			),
		},
	)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps:                    steps,
	})
}

func testAccAppServiceLoadBalancerCIDRConfig(resourceName, clusterName, dataSourceName, appServiceName, clusterCIDR, loadBalancerCIDR string, nodes int) string {
	return fmt.Sprintf(`
%[1]s

%[10]s

resource "couchbase-capella_app_service" "%[4]s" {
  organization_id    = "%[2]s"
  project_id         = "%[3]s"
  cluster_id         = couchbase-capella_cluster.%[5]s.id
  name               = "%[6]s"
  nodes              = %[9]d
  load_balancer_cidr = "%[8]s"
  compute = {
    cpu = 2
    ram = 4
  }
}

data "couchbase-capella_app_services" "%[7]s" {
  organization_id = "%[2]s"

  depends_on = [couchbase-capella_app_service.%[4]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, clusterName, appServiceName, dataSourceName, loadBalancerCIDR, nodes,
		testAccAzureClusterForAppServiceConfig(clusterName, clusterCIDR))
}

func testAccAppServiceLoadBalancerCIDRMalformedConfig(resourceName, clusterName, appServiceName, clusterCIDR, loadBalancerCIDR string) string {
	return fmt.Sprintf(`
%[1]s

%[7]s

resource "couchbase-capella_app_service" "%[4]s" {
  organization_id    = "%[2]s"
  project_id         = "%[3]s"
  cluster_id         = couchbase-capella_cluster.%[5]s.id
  name               = "%[6]s"
  load_balancer_cidr = "%[8]s"
  compute = {
    cpu = 2
    ram = 4
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, clusterName, appServiceName,
		testAccAzureClusterForAppServiceConfig(clusterName, clusterCIDR), loadBalancerCIDR)
}

func testAccAzureClusterForAppServiceConfig(clusterName, clusterCIDR string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_cluster" "%[3]s" {
  organization_id = "%[1]s"
  project_id      = "%[2]s"
  name            = "%[3]s"
  description     = "Terraform Acceptance Test app service load balancer CIDR"
  cloud_provider = {
    type   = "azure"
    region = "eastus"
    cidr   = "%[4]s"
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
      services     = ["data", "index", "query"]
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
`, globalOrgId, globalProjectId, clusterName, clusterCIDR)
}

func testAccCheckAppServicesDataSourceLoadBalancerCIDR(dataSourceReference, resourceReference, wantCIDR string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		appService, ok := state.RootModule().Resources[resourceReference]
		if !ok {
			return fmt.Errorf("resource %q not found in state", resourceReference)
		}
		appServiceID := appService.Primary.Attributes["id"]

		dataSource, ok := state.RootModule().Resources[dataSourceReference]
		if !ok {
			return fmt.Errorf("data source %q not found in state", dataSourceReference)
		}

		attrs := dataSource.Primary.Attributes
		count, err := strconv.Atoi(attrs["data.#"])
		if err != nil {
			return fmt.Errorf("invalid data.# on %q: %w", dataSourceReference, err)
		}

		for i := 0; i < count; i++ {
			if attrs[fmt.Sprintf("data.%d.id", i)] != appServiceID {
				continue
			}

			key := fmt.Sprintf("data.%d.load_balancer_cidr", i)
			if got := attrs[key]; got != wantCIDR {
				return fmt.Errorf("%s = %q, want %q", key, got, wantCIDR)
			}

			return nil
		}

		return fmt.Errorf("app service %q not found in %s.data across %d entries", appServiceID, dataSourceReference, count)
	}
}
