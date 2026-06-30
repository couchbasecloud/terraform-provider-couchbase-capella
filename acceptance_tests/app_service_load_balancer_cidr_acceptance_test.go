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

// TestAccAppServiceLoadBalancerCIDR tests the load_balancer_cidr attribute on the
// app service resource and data source against a dedicated Azure cluster.
func TestAccAppServiceLoadBalancerCIDR(t *testing.T) {
	// load_balancer_cidr is Azure only and must be a valid IPv4 /24 network address.
	// The global cluster is AWS, so this test stands up its own Azure cluster. The
	// 192.168.x.0/24 ranges cannot overlap the random 10.x.y.0/23 cluster CIDR, and
	// the control plane adds the /24 as an extra VNet address space, so the app
	// service deploys cleanly.
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

	// This test uses its own cluster and app service (no shared global resources),
	// so it can run in parallel.
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing, plus data source coverage.
			{
				Config: testAccAppServiceLoadBalancerCIDRConfig(resourceName, clusterName, dataSourceName, appServiceName, clusterCIDR, loadBalancerCIDR),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "load_balancer_cidr", loadBalancerCIDR),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider", "Azure"),
					resource.TestCheckResourceAttr(resourceReference, "name", appServiceName),
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
			// ImportState testing confirms the API round-trips load_balancer_cidr:
			// on import there is no prior state, so the value comes from the API.
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppServiceImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// load_balancer_cidr is create-only, so changing it must force a
			// replace. Plan checks require an apply step, so this recreates the
			// app service with the new CIDR.
			{
				Config: testAccAppServiceLoadBalancerCIDRConfig(resourceName, clusterName, dataSourceName, appServiceName, clusterCIDR, loadBalancerCIDRChanged),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "load_balancer_cidr", loadBalancerCIDRChanged),
				),
			},
		},
	})
}

// TestAccAppServiceLoadBalancerCIDRAzureOnly verifies the API rejects a requested
// load balancer CIDR on a non-Azure cluster and that the provider surfaces the error.
// The global cluster is AWS, and the Azure-only check runs before the "app service
// already exists" check, so the global cluster's existing app service is irrelevant.
func TestAccAppServiceLoadBalancerCIDRAzureOnly(t *testing.T) {
	const loadBalancerCIDR = "192.168.0.0/24"

	resourceName := randomStringWithPrefix("tf_acc_app_svc_lb_aws_")
	appServiceName := randomStringWithPrefix("tf_acc_app_svc_lb_aws_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServiceLoadBalancerCIDRAzureOnlyConfig(resourceName, appServiceName, loadBalancerCIDR),
				ExpectError: regexp.MustCompile(`(?s)only supported for Azure App Services`),
			},
		},
	})
}

func testAccAppServiceLoadBalancerCIDRConfig(resourceName, clusterName, dataSourceName, appServiceName, clusterCIDR, loadBalancerCIDR string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[5]s"
  description     = "Terraform Acceptance Test app service load balancer CIDR"
  cloud_provider = {
    type   = "azure"
    region = "eastus"
    cidr   = "%[7]s"
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

data "couchbase-capella_app_services" "%[9]s" {
  organization_id = "%[2]s"

  depends_on = [couchbase-capella_app_service.%[4]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, clusterName, appServiceName, clusterCIDR, loadBalancerCIDR, dataSourceName)
}

func testAccAppServiceLoadBalancerCIDRAzureOnlyConfig(resourceName, appServiceName, loadBalancerCIDR string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service" "%[5]s" {
  organization_id    = "%[2]s"
  project_id         = "%[3]s"
  cluster_id         = "%[4]s"
  name               = "%[6]s"
  load_balancer_cidr = "%[7]s"
  compute = {
    cpu = 2
    ram = 4
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, appServiceName, loadBalancerCIDR)
}

// testAccCheckAppServicesDataSourceLoadBalancerCIDR finds the app service by id in
// the app_services data source list and asserts its load_balancer_cidr.
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
