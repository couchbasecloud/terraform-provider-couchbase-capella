package acceptance_tests

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAppServicesCIDRDataSource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_services_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName

	dataSourceName := randomStringWithPrefix("tf_acc_app_services_cidr_ds_")
	dataSourceReference := "data.couchbase-capella_app_services_cidr." + dataSourceName

	// Generate a random CIDR block for testing. This ensures that the test is not dependent on any pre-existing data and can be run multiple times without conflicts.
	// This is used as an allowed IP on the App Service for purposes of testing the data source.
	allowedCIDR := fmt.Sprintf("172.16.%d.%d/32", rand.Intn(256), rand.Intn(256)) // #nosec G404

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicesCIDRDataSourceConfig(resourceName, dataSourceName, allowedCIDR),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", allowedCIDR),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),

					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(dataSourceReference, "data.#"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceReference, "data.*", map[string]string{
						"cidr": allowedCIDR,
					}),
				),
			},
		},
	})
}

func TestAccAppServicesCIDRDataSourceMissingRequiredFields(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_services_cidr_ds_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_services_cidr" "%[2]s" {}
`, globalProviderBlock, dataSourceName),
				ExpectError: regexp.MustCompile(`The argument "(organization_id|project_id|cluster_id|app_service_id)" is required`),
			},
		},
	})
}

func testAccAppServicesCIDRDataSourceConfig(resourceName, dataSourceName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_services_cidr" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  cidr            = "%[8]s"
  comment         = "terraform app services cidr datasource acceptance test"
}

data "couchbase-capella_app_services_cidr" "%[7]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"

  depends_on = [couchbase-capella_app_services_cidr.%[6]s]
}
`,
		globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		resourceName,
		dataSourceName,
		cidr,
	)
}
