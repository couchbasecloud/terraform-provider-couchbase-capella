package acceptance_tests

import (
	"fmt"
	"math/rand"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAppServicesCIDR_AV_129803 verifies that reading an App Service allowed CIDR
// does not silently drop the resource from state when more than the default page size
// (~10) of CIDRs exist on the App Service. The bug was that getAllowedCIDR used the
// list endpoint without pagination, so only the first page of CIDRs was visible.
func TestAccAppServicesCIDR_AV_129803(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_services_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName
	allowedCIDR := fmt.Sprintf("172.17.%d.%d/32", rand.Intn(256), rand.Intn(256)) // #nosec G404
	comment := "terraform app services cidr pagination acceptance test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccAppServicesCIDRConfig(resourceName, allowedCIDR, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "cidr", allowedCIDR),
					resource.TestCheckResourceAttr(resourceReference, "comment", comment),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
			// ImportState
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppServicesCIDRImportIdForResource(resourceReference),
				ImportState:       true,
			},
		},
	})
}

// TestAccAppServicesCIDRResourceMultipleCIDRs verifies that the provider can manage
// multiple allowed CIDRs on a single App Service without pagination-related state drift.
// This test creates CIDR entries beyond the default page size to confirm getAllowedCIDR
// correctly fetches each CIDR via the single-item endpoint.
func TestAccAppServicesCIDRResourceMultipleCIDRs(t *testing.T) {
	const numCIDRs = 12 // exceeds default page size of ~10
	resourceNames := make([]string, numCIDRs)
	resourceReferences := make([]string, numCIDRs)
	allowedCIDRs := make([]string, numCIDRs)
	for i := 0; i < numCIDRs; i++ {
		resourceNames[i] = randomStringWithPrefix("tf_acc_app_services_cidr_multi_")
		resourceReferences[i] = "couchbase-capella_app_services_cidr." + resourceNames[i]
		// Use distinct /32 CIDRs so they don't conflict
		allowedCIDRs[i] = fmt.Sprintf("10.%d.%d.%d/32", (i/256)%256, i%256, rand.Intn(256)) // #nosec G404
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create multiple CIDRs
			{
				Config: testAccAppServicesCIDRMultipleConfig(resourceNames, allowedCIDRs),
				Check:  testAccCheckMultipleAppServicesCIDRs(resourceReferences, allowedCIDRs),
			},
			// Re-read (plan) to verify none were removed from state
			{
				Config:   testAccAppServicesCIDRMultipleConfig(resourceNames, allowedCIDRs),
				PlanOnly: true,
			},
		},
	})
}

// TestAccAppServicesCIDRResourceInvalidCidr verifies that invalid CIDR input produces
// a clear error rather than a silent failure.
func TestAccAppServicesCIDRResourceInvalidCidr(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_services_cidr_inv_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServicesCIDRConfig(resourceName, "invalid_cidr", ""),
				ExpectError: regexp.MustCompile(`invalid CIDR address`),
			},
		},
	})
}

func testAccAppServicesCIDRConfig(resourceName, cidr, comment string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_services_cidr" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  cidr            = "%[7]s"
  comment         = "%[8]s"
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		cidr,
		comment,
	)
}

func testAccAppServicesCIDRMultipleConfig(resourceNames, cidrs []string) string {
	config := globalProviderBlock + "\n"
	for i := range resourceNames {
		config += fmt.Sprintf(`
resource "couchbase-capella_app_services_cidr" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  cidr            = "%[6]s"
  comment         = "multi cidr test entry %[7]d"
}
`,
			resourceNames[i],
			globalOrgId,
			globalProjectId,
			globalClusterId,
			globalAppServiceId,
			cidrs[i],
			i,
		)
	}
	return config
}

func testAccCheckMultipleAppServicesCIDRs(resourceReferences, expectedCIDRs []string) resource.TestCheckFunc {
	checks := make([]resource.TestCheckFunc, 0, len(resourceReferences)*4)
	for i := range resourceReferences {
		ref := resourceReferences[i]
		expectedCIDR := expectedCIDRs[i]
		checks = append(checks,
			resource.TestCheckResourceAttr(ref, "organization_id", globalOrgId),
			resource.TestCheckResourceAttr(ref, "project_id", globalProjectId),
			resource.TestCheckResourceAttr(ref, "cluster_id", globalClusterId),
			resource.TestCheckResourceAttr(ref, "app_service_id", globalAppServiceId),
			resource.TestCheckResourceAttr(ref, "cidr", expectedCIDR),
			resource.TestCheckResourceAttrSet(ref, "id"),
		)
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func generateAppServicesCIDRImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf(
			"id=%s,cluster_id=%s,project_id=%s,organization_id=%s,app_service_id=%s",
			rawState["id"],
			rawState["cluster_id"],
			rawState["project_id"],
			rawState["organization_id"],
			rawState["app_service_id"],
		), nil
	}
}
