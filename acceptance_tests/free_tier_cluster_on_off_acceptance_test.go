package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFreeTierClusterOnOffResource(t *testing.T) {
	clusterName := randomStringWithPrefix("tf_acc_free_tier_cluster_on_off_")
	resourceName := randomStringWithPrefix("tf_acc_free_tier_cluster_on_off_")
	resourceReference := "couchbase-capella_free_tier_cluster_on_off." + resourceName
	cidr := generateRandomCIDR()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccFreeTierClusterOnOffResourceLiveConfig(clusterName, resourceName, cidr, "off"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(resourceReference, "cluster_id"),
					resource.TestCheckResourceAttr(resourceReference, "state", "off"),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateFreeTierClusterOnOffImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "cluster_id",
			},
			{
				Config: testAccFreeTierClusterOnOffResourceLiveConfig(clusterName, resourceName, cidr, "on"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "cluster_id"),
					resource.TestCheckResourceAttr(resourceReference, "state", "on"),
				),
			},
			{
				Config: testAccFreeTierClusterOnOffResourceLiveConfig(clusterName, resourceName, cidr, "off"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "cluster_id"),
					resource.TestCheckResourceAttr(resourceReference, "state", "off"),
				),
			},
		},
	})
}

func TestAccFreeTierClusterOnOffResourceInvalidState(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_free_tier_cluster_on_off_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccFreeTierClusterOnOffResourceConfig(resourceName, "frozen"),
				ExpectError: regexp.MustCompile(`(?s)invalid state value|invalid state`),
			},
		},
	})
}

func testAccFreeTierClusterOnOffResourceLiveConfig(clusterName, resourceName, cidr, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Free tier cluster for on/off acceptance testing."

  cloud_provider = {
    type   = "aws"
    region = "us-east-2"
    cidr   = "%[6]s"
  }
}

resource "couchbase-capella_free_tier_cluster_on_off" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_free_tier_cluster.%[4]s.id
  state           = "%[7]s"
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterName, resourceName, cidr, state)
}

func generateFreeTierClusterOnOffImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, module := range state.Modules {
			if len(module.Resources) > 0 {
				if v, ok := module.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		if rawState == nil {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		return fmt.Sprintf(
			"organization_id=%s,project_id=%s,cluster_id=%s",
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"],
		), nil
	}
}

func testAccFreeTierClusterOnOffResourceConfig(resourceName, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster_on_off" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "11111111-1111-1111-1111-111111111111"
  cluster_id      = "22222222-2222-2222-2222-222222222222"
  state           = "%[3]s"
}
`, globalProviderBlock, resourceName, state)
}
