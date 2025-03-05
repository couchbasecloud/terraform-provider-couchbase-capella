package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccProjectResource is a Terraform acceptance test that covers the lifecycle of a Capella project resource.
//
// The test includes the following steps:
//  1. PreCheck: Ensure the prerequisites for acceptance testing.
//  2. Create and Read Testing: Configure the test environment and create a Capella project resource.
//     Verify that the project has the expected attributes such as name, description, and etag.
//  3. ImportState Testing: Import the state of the created project and verify the imported state matches the expected state.
//  4. Update and Read Testing: Modify the project's attributes and ensure the changes are applied successfully.
//  5. Delete Testing: Automatically occurs in the TestCase as part of cleanup.

// TODO:  AV-96937: make project acceptance tests concurrent
func TestAccProjectResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_project_")
	resourceReference := "couchbase-capella_project." + resourceName
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "terraform acceptance test project"),
					resource.TestCheckResourceAttr(resourceReference, "etag", "Version: 1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateProjectImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfigUpdate(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "description_update"),
				),
			},
			{
				Config: testAccProjectResourceConfigUpdateWithIfMatch(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "description_update_with_match"),
					resource.TestCheckResourceAttr(resourceReference, "etag", "Version: 3"),
					resource.TestCheckResourceAttr(resourceReference, "if_match", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccProjectCreateWithReqFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_project_")
	resourceReference := "couchbase-capella_project." + resourceName
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfigRequired(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					resource.TestCheckResourceAttr(resourceReference, "etag", "Version: 1"),
				),
			},
		},
	})
}

func TestAccProjectValidUpdate(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_project_")
	resourceReference := "couchbase-capella_project." + resourceName
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "terraform acceptance test project"),
					resource.TestCheckResourceAttr(resourceReference, "etag", "Version: 1"),
				),
			},
			//update the project name and description
			{
				Config: testAccProjectResourceConfigUpdate(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "description_update"),
					resource.TestCheckResourceAttr(resourceReference, "etag", "Version: 2"),
				),
			},
		},
	})
}

func TestAccProjectInvalidResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_project_")
	resourceReference := "couchbase-capella_project." + resourceName
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Invalid field in create testing
			{
				Config:      testAccProjectResourceConfigInvalid(resourceName),
				ExpectError: regexp.MustCompile("An argument named \"unwantedfiled\" is not expected here"),
			},
			{
				Config: testAccProjectResourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "terraform acceptance test project"),
					resource.TestCheckResourceAttr(resourceReference, "etag", "Version: 1"),
				),
			},
			{
				Config:      testAccProjectResourceConfigUpdateInvalid(resourceName),
				ExpectError: regexp.MustCompile("server cannot or will not process the request.*"),
			},
		},
	})
}

// testAccProjectResourceConfig generates a Terraform configuration string for creating a Capella project resource.
func testAccProjectResourceConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
    organization_id = "%[2]s"
	name            = "%[3]s"
	description     = "terraform acceptance test project"
}
`, globalProviderBlock, globalOrgId, resourceName)
}

// testAccProjectResourceConfigUpdate generates a Terraform configuration string for updating a Capella project resource.
func testAccProjectResourceConfigUpdate(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
   organization_id = "%[2]s"
	name            = "%[3]s"
	description     = "description_update"
}
`, globalProviderBlock, globalOrgId, resourceName)
}

// testAccProjectResourceConfigUpdateWithIfMatch generates a Terraform configuration string for updating a Capella project resource
// with an "if_match" attribute.
func testAccProjectResourceConfigUpdateWithIfMatch(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
    organization_id = "%[2]s"
	name            = "%[3]s"
	description     = "description_update_with_match"
	if_match        =  2
}
`, globalProviderBlock, globalOrgId, resourceName)
}

// generateProjectImportIdForResource generates a project import ID based on the provided resource name
// and the attributes in the Terraform state.
//
// This function takes a resource name as input and returns a function of type `resource.ImportStateIdFunc`.
// The generated import ID is in the format "id=<value>,organization_id=<value>".
func generateProjectImportIdForResource(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceName]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf("id=%s,organization_id=%s", rawState["id"], rawState["organization_id"]), nil
	}
}

func testAccProjectResourceConfigRequired(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[2]s" {
    organization_id = "%[3]s"
	name            = "%[2]s"
}
`, globalProviderBlock, resourceName, globalOrgId)

}

func testAccProjectResourceConfigUpdateInvalid(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[2]s" {
    organization_id = "abc-def"
	name            = "%[2]s"
}
`, globalProviderBlock, resourceName)

}

func testAccProjectResourceConfigInvalid(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
    organization_id = "%[2]s"
	name            = "%[3]s"
	unwantedfiled   = "unwanted value"
}
`, globalProviderBlock, globalOrgId, resourceName)
}
