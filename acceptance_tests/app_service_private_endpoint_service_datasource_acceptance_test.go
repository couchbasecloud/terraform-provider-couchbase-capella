package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDatasourceAppServicePrivateEndpointService reads the private endpoint service state of the
// global App Service. Without the aws block, no command is requested, so command stays null.
func TestAccDatasourceAppServicePrivateEndpointService(t *testing.T) {
	dataSourceName := randomStringWithPrefix("tf_acc_app_service_pe_service_")
	dataSourceReference := "data.couchbase-capella_app_service_private_endpoint_service." + dataSourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicePrivateEndpointServiceDatasourceConfig(dataSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(dataSourceReference, "state"),
					resource.TestCheckResourceAttrSet(dataSourceReference, "target_state"),
					resource.TestCheckNoResourceAttr(dataSourceReference, "command"),
				),
			},
		},
	})
}

// TestAccDatasourceAppServicePrivateEndpointServiceInvalidVPCID verifies that a vpc_id shorter than the
// OpenAPI minimum length fails at plan time, before any API call, so dummy IDs are sufficient.
func TestAccDatasourceAppServicePrivateEndpointServiceInvalidVPCID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_service_private_endpoint_service" "test" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  app_service_id  = "33333333-3333-3333-3333-333333333333"

  aws = {
    vpc_id     = "vpc-short"
    subnet_ids = ["subnet-1234567890abcdef0"]
  }
}
`, globalProviderBlock),
				ExpectError: regexp.MustCompile(`(?s)vpc_id.*string length must be between 12 and 21`),
			},
		},
	})
}

func testAccAppServicePrivateEndpointServiceDatasourceConfig(dataSourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_app_service_private_endpoint_service" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
}
`, globalProviderBlock, dataSourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId)
}
