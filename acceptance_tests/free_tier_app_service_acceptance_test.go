package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFreeTierAppServiceResource(t *testing.T) {
	clusterName := randomStringWithPrefix("tf_acc_free_tier_cluster_appsvc_")
	appServiceName := randomStringWithPrefix("tf_acc_free_tier_appsvc_")
	resourceReference := "couchbase-capella_free_tier_app_service." + appServiceName
	cidr := generateRandomCIDR()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeTierAppServiceResourceConfig(clusterName, appServiceName, cidr, "Free tier app service acceptance test."),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttr(resourceReference, "name", appServiceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "Free tier app service acceptance test."),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(resourceReference, "cluster_id"),
					resource.TestCheckResourceAttrSet(resourceReference, "current_state"),
					resource.TestCheckResourceAttrSet(resourceReference, "etag"),
				),
			},
			{
				Config: testAccFreeTierAppServiceResourceConfig(clusterName, appServiceName, cidr, "Updated free tier app service acceptance test."),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", appServiceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "Updated free tier app service acceptance test."),
				),
			},
		},
	})
}

func testAccFreeTierAppServiceResourceConfig(clusterName, appServiceName, cidr, description string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Free tier cluster for app service acceptance testing."

  cloud_provider = {
    type   = "aws"
    region = "us-east-2"
    cidr   = "%[6]s"
  }
}

resource "couchbase-capella_free_tier_app_service" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_free_tier_cluster.%[4]s.id
  name            = "%[5]s"
  description     = "%[7]s"
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterName, appServiceName, cidr, description)
}
