package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
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

func TestAccDatasourceAppServices(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_svcs_ds_")
	dataSourceReference := "data.couchbase-capella_app_services." + dataSourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicesDataSourceConfig(dataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dataSourceReference, "data.#"),
					testAccCheckAppServicesDataSourceContainsGlobalAppService(dataSourceReference),
				),
			},
		},
	})
}

func TestAccDatasourceAppServicesMissingOrganization(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_svcs_ds_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_services" "%[2]s" {}
`, globalProviderBlock, dataSourceName),
				ExpectError: regexp.MustCompile(`The argument "organization_id" is required`),
			},
		},
	})
}

func testAccAppServiceResourceConfig(resourceName string) string {
	clusterName := randomStringWithPrefix("tf_acc_cluster_")
	cidr := generateRandomCIDR()
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[5]s"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[6]s"
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
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, clusterName, cidr)
}

func testAccAppServicesDataSourceConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_services" "%[3]s" {
  organization_id = "%[2]s"
}
`, globalProviderBlock, globalOrgId, dataSourceName)
}

func testAccCheckAppServicesDataSourceContainsGlobalAppService(dataSourceReference string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
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
			if attrs[fmt.Sprintf("data.%d.id", i)] != globalAppServiceId {
				continue
			}

			expectedAttrs := map[string]string{
				"organization_id": globalOrgId,
			}
			if globalAppServiceCreated {
				expectedAttrs["cluster_id"] = globalClusterId
				expectedAttrs["name"] = globalAppServiceName
				expectedAttrs["nodes"] = "2"
				expectedAttrs["compute.cpu"] = "2"
				expectedAttrs["compute.ram"] = "4"
			}

			for suffix, want := range expectedAttrs {
				if err := assertAppServicesDataSourceAttr(attrs, i, suffix, want); err != nil {
					return err
				}
			}

			for _, suffix := range []string{
				"cluster_id",
				"name",
				"nodes",
				"cloud_provider",
				"current_state",
				"version",
				"compute.cpu",
				"compute.ram",
				"audit.created_at",
				"audit.modified_at",
				"audit.version",
			} {
				key := fmt.Sprintf("data.%d.%s", i, suffix)
				if attrs[key] == "" {
					return fmt.Errorf("attribute %q expected to be set on matched app service %s", key, globalAppServiceId)
				}
			}

			return nil
		}

		return fmt.Errorf("expected app service %q in %s.data, not found across %d entries", globalAppServiceId, dataSourceReference, count)
	}
}

func assertAppServicesDataSourceAttr(attrs map[string]string, index int, suffix, want string) error {
	key := fmt.Sprintf("data.%d.%s", index, suffix)
	if got := attrs[key]; got != want {
		return fmt.Errorf("%s = %q, want %q", key, got, want)
	}
	return nil
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
