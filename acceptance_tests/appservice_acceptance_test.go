package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

func TestAccAppServiceResourceOptionalFieldsAndScale(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_")
	resourceReference := "couchbase-capella_app_service." + resourceName
	clusterName := randomStringWithPrefix("tf_acc_cluster_")
	cidr := generateRandomCIDR()
	appServiceName := randomStringWithPrefix("tf_acc_app_svc_")
	description := "terraform app service optional fields acceptance test"

	const version = "4.0"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceResourceOptionalFieldsConfig(resourceName, clusterName, cidr, appServiceName, description, version, 2, 2, 4),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "name", appServiceName),
					resource.TestCheckResourceAttr(resourceReference, "description", description),
					resource.TestCheckResourceAttr(resourceReference, "nodes", "2"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr(resourceReference, "compute.cpu", "2"),
					resource.TestCheckResourceAttr(resourceReference, "compute.ram", "4"),
					resource.TestCheckResourceAttr(resourceReference, "version", version),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "cluster_id"),
					resource.TestCheckResourceAttrSet(resourceReference, "current_state"),
					resource.TestCheckResourceAttrSet(resourceReference, "etag"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppServiceImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAppServiceResourceOptionalFieldsConfig(resourceName, clusterName, cidr, appServiceName, description, version, 3, 4, 8),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "name", appServiceName),
					resource.TestCheckResourceAttr(resourceReference, "description", description),
					resource.TestCheckResourceAttr(resourceReference, "nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider", "AWS"),
					resource.TestCheckResourceAttr(resourceReference, "compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "compute.ram", "8"),
					resource.TestCheckResourceAttr(resourceReference, "version", version),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "cluster_id"),
					resource.TestCheckResourceAttrSet(resourceReference, "current_state"),
					resource.TestCheckResourceAttrSet(resourceReference, "etag"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
				),
			},
			{
				// Changing the major.minor version is not supported and must be rejected at apply time.
				Config:      testAccAppServiceResourceOptionalFieldsConfig(resourceName, clusterName, cidr, appServiceName, description, "1.0", 3, 4, 8),
				ExpectError: regexp.MustCompile(`(?s)version as this is not supported`),
			},
		},
	})
}

func TestAccAppServiceResourceNameTooLong(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServiceResourceNameTooLongConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)Attribute name string length must be at most 256, got: 257`),
			},
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

func testAccAppServiceResourceOptionalFieldsConfig(resourceName, clusterName, cidr, appServiceName, description, version string, nodes, cpu, ram int) string {
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
	organization_id = "%[2]s"
	project_id      = "%[3]s"
	cluster_id      = couchbase-capella_cluster.%[5]s.id
	name            = "%[7]s"
	description     = "%[8]s"
	version         = "%[12]s"
	nodes           = %[9]d
	compute = {
		cpu = %[10]d
		ram = %[11]d
	}
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, clusterName, cidr, appServiceName, description, nodes, cpu, ram, version)
}

func testAccAppServiceResourceNameTooLongConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service" "%[2]s" {
	organization_id = "00000000-0000-0000-0000-000000000000"
	project_id      = "11111111-1111-1111-1111-111111111111"
	cluster_id      = "22222222-2222-2222-2222-222222222222"
	name            = "%[3]s"
	description     = "Invalid App Service name length."
	nodes           = 2

	compute = {
		cpu = 2
		ram = 4
	}
}
`, globalProviderBlock, resourceName, strings.Repeat("a", 257))
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
