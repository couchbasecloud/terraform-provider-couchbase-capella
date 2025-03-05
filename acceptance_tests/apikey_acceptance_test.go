package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccApiKeyResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resourceReference := "couchbase-capella_apikey." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "description"),
					resource.TestCheckResourceAttr(resourceReference, "expiry", "150"),
					resource.TestCheckResourceAttr(resourceReference, "allowed_cidrs.0", "10.1.42.0/23"),
					resource.TestCheckResourceAttr(resourceReference, "allowed_cidrs.1", "10.1.42.1/23"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.roles.0", "projectDataReader"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.roles.1", "projectManager"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.type", "project"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateApiKeyImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccApiKeyResourceWithOnlyReqField(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resourceReference := "couchbase-capella_apikey." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfigWithOnlyReqFieldConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					resource.TestCheckResourceAttr(resourceReference, "expiry", "180"),
					resource.TestCheckResourceAttr(resourceReference, "allowed_cidrs.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.1", "organizationOwner"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateApiKeyImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Rotate testing
			{
				Config: testAccApiKeyResourceConfigWithOnlyReqFieldRotateConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					resource.TestCheckResourceAttr(resourceReference, "expiry", "180"),
					resource.TestCheckResourceAttr(resourceReference, "allowed_cidrs.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.1", "organizationOwner"),
					resource.TestCheckResourceAttr(resourceReference, "rotate", "1"),
					resource.TestCheckResourceAttr(resourceReference, "secret", "abc"),
				),
			},
		},
	})
}

func TestAccApiKeyResourceForOrgOwner(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resourceReference := "couchbase-capella_apikey." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfigForOrgOwnerConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					resource.TestCheckResourceAttr(resourceReference, "expiry", "180"),
					resource.TestCheckResourceAttr(resourceReference, "allowed_cidrs.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.roles.0", "projectDataReader"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.roles.1", "projectManager"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.type", "project"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateApiKeyImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccApiKeyResourceInvalidScenarioRotateShouldNotPassedWhileCreate(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccApiKeyResourceConfigRotateSetConfig(resourceName),
				ExpectError: regexp.MustCompile("rotate value should not be set"),
			},
		},
	})
}

func testAccApiKeyResourceConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = "%[3]s"
  name               = "%[2]s"
  description        = "description"
  expiry             = 150
  organization_roles = ["organizationMember"]
  allowed_cidrs      = ["10.1.42.1/23", "10.1.42.0/23"]
  resources = [
    {
      id    = "%[4]s"
      roles = ["projectManager", "projectDataReader"]
      type  = "project"
    }
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId)
}

func testAccApiKeyResourceConfigWithOnlyReqFieldConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = "%[3]s"
  name               = "%[2]s"
  organization_roles = ["organizationOwner", "organizationMember"]
}
`, globalProviderBlock, resourceName, globalOrgId)
}

func testAccApiKeyResourceConfigForOrgOwnerConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = "%[3]s"
  name               = "%[2]s"
  organization_roles = [ "organizationMember"]
  resources = [
	  {
		id = "%[4]s"
		roles = [
		  "projectManager",
		  "projectDataReader"
		]
	  }
  ]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId)
}

func testAccApiKeyResourceConfigRotateSetConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = "%[3]s"
  name               = "%[2]s"
  organization_roles = [ "organizationMember"]
  resources = [
	  {
		id = "%[4]s"
		roles = [
		  "projectManager",
		  "projectDataReader"
		]
	  }
  ]
  rotate = 1
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId)
}

func testAccApiKeyResourceConfigWithOnlyReqFieldRotateConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = "%[3]s"
  name               = "%[2]s"
  organization_roles = ["organizationOwner", "organizationMember"]
  rotate             = 1
  secret             = "abc"
}
`, globalProviderBlock, resourceName, globalOrgId)
}

func generateApiKeyImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf("id=%s,organization_id=%s", rawState["id"], globalOrgId), nil
	}
}
