package acceptance_tests

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/provider"
	acctest "terraform-provider-capella/internal/testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"capella": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// TestAccProjectResource is a Terraform acceptance test that covers the lifecycle of a Capella project resource.
//
// The test includes the following steps:
//  1. PreCheck: Ensure the prerequisites for acceptance testing.
//  2. Create and Read Testing: Configure the test environment and create a Capella project resource.
//     Verify that the project has the expected attributes such as name, description, and etag.
//  3. ImportState Testing: Import the state of the created project and verify the imported state matches the expected state.
//  4. Update and Read Testing: Modify the project's attributes and ensure the changes are applied successfully.
//  5. Delete Testing: Automatically occurs in the TestCase as part of cleanup.
func TestAccProjectResource(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	resourceName := "capella_project." + rnd
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig(acctest.Cfg, rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 1"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: generateProjectImportIdForResource(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfigUpdate(acctest.Cfg, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project_name_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "description_update"),
				),
			},
			{
				Config: testAccProjectResourceConfigUpdateWithIfMatch(acctest.Cfg, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project_name_update_with_if_match"),
					resource.TestCheckResourceAttr(resourceName, "description", "description_update_with_match"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 3"),
					resource.TestCheckResourceAttr(resourceName, "if_match", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCreateProjectWithReqFields(t *testing.T) {

	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	resourceName := "capella_project." + rnd
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfigRequired(acctest.Cfg, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 1"),
				),
			},
		},
	})
}

func TestAccCreateProjectOptFields(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	resourceName := "capella_project." + rnd
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(acctest.Cfg, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test project"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 1"),
				),
			},
			{
				ResourceName:      "capella_project.acc_test",
				ImportStateIdFunc: generateProjectImportIdForResource(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfigUpdate(acctest.Cfg, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "description update"),
				),
			},
			{
				Config: testAccProjectResourceConfigUpdateWithIfMatch(acctest.Cfg, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project_name_update_with_if_match"),
					resource.TestCheckResourceAttr(resourceName, "description", "description_update_with_match"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 3"),
					resource.TestCheckResourceAttr(resourceName, "if_match", "2"),
				),
			},
		},
	})
}

// testAccProjectResourceConfig generates a Terraform configuration string for creating a Capella project resource.
func testAccProjectResourceConfig(cfg, rnd string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "%[2]s" {
    organization_id = var.organization_id
	name            = "%[2]s"
	description     = "description"
}
`, cfg, rnd)
}

// testAccProjectResourceConfigUpdate generates a Terraform configuration string for updating a Capella project resource.
func testAccProjectResourceConfigUpdate(cfg, rnd string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "%[2]s" {
   organization_id = var.organization_id
	name            = "acc_test_project_name_update"
	description     = "description_update"
}
`, cfg, rnd)
}

// testAccProjectResourceConfigUpdateWithIfMatch generates a Terraform configuration string for updating a Capella project resource
// with an "if_match" attribute.
func testAccProjectResourceConfigUpdateWithIfMatch(cfg, rnd string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "%[2]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name_update_with_if_match"
	description     = "description_update_with_match"
	if_match        =  2
}
`, cfg, rnd)
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
		fmt.Printf("raw state %s", rawState)
		return fmt.Sprintf("id=%s,organization_id=%s", rawState["id"], rawState["organization_id"]), nil
	}
}

func testAccDeleteProject(projectResourceReference string) resource.TestCheckFunc {
	log.Println("Deleting the project")
	return func(s *terraform.State) error {
		var projectState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[projectResourceReference]; ok {
					projectState = v.Primary.Attributes
				}
			}
		}
		data, err := acctest.TestClient()
		if err != nil {
			return err
		}
		host := os.Getenv("TF_VAR_host")
		orgid := os.Getenv("TF_VAR_organization_id")
		authToken := os.Getenv("TF_VAR_auth_token")
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s", host, orgid, projectState["id"])
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
		_, err = data.Client.Execute(
			cfg,
			nil,
			authToken,
			nil,
		)
		if err != nil {
			return err
		}
		return nil
	}
}

// Update

func TestAccValidProjcetUpdate(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	resourceName := "capella_project." + rnd
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(acctest.Cfg, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test project"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 1"),
				),
			},
			//update the project name and description
			{
				Config: testAccProjectResourceConfigUpdate(acctest.Cfg, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "description update"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 2"),
				),
			},
		},
	})
}

func TestAccInvalidProjectResource(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	resourceName := "capella_project." + rnd
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Invalid field in create testing
			{
				Config:      testAccProjectResourceConfigInvalid(acctest.Cfg, rnd),
				ExpectError: regexp.MustCompile("An argument named \"unwantedfiled\" is not expected here"),
			},
			{
				Config: testAccProjectResourceConfig(acctest.Cfg, rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test project"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 1"),
				),
			},
			//Update the organisation id
			{
				Config:      testAccProjectResourceConfigUpdateInvalid(acctest.Cfg, rnd),
				ExpectError: regexp.MustCompile("The server cannot or will not\nprocess the request due to something that is perceived to be a client\nerror"),
			},
		},
	})
}

func TestAccDeleteProjectBeforeDestroy(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	resourceName := "capella_project." + rnd
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(acctest.Cfg, rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "acc_test_project"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test project"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 1"),
					testAccDeleteProject(resourceName),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
		},
	})
}

func testAccProjectResourceConfigRequired(cfg string, rnd string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "%[2]s" {
    organization_id = var.organization_id
	name            = "acc_test_project"
}
`, cfg, rnd)

}

func testAccProjectResourceConfigUpdateInvalid(cfg, rnd string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "%[2]s" {
    organization_id = "abc-def"
	name            = "acc_test_project"
}
`, cfg, rnd)

}
func testAccProjectResourceConfigInvalid(cfg, rnd string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "%[2]s" {
    organization_id = var.organization_id
	name            = "acc_test_project"
	unwantedfiled   = "unwanted value"
}
`, cfg, rnd)

}
