package acceptance_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccApiKeyResource(t *testing.T) {
	resourceName := "acc_apikey_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_apikey." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfig(acctest.ProjectCfg, resourceName, projectResourceName, projectResourceReference),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsApiKeyResource(resourceReference),
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
			//// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateApiKeyImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccApiKeyResourceWithMultipleResources(t *testing.T) {
	resourceName := "acc_apikey_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_apikey." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfigWithMultipleResources(acctest.ProjectCfg, resourceName, projectResourceName, projectResourceReference),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsApiKeyResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					resource.TestCheckResourceAttr(resourceReference, "expiry", "180"),
					resource.TestCheckResourceAttr(resourceReference, "allowed_cidrs.0", "10.1.42.0/23"),
					resource.TestCheckResourceAttr(resourceReference, "allowed_cidrs.1", "10.1.42.1/23"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.#", "2"),
				),
			},
			//// ImportState testing
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
	resourceName := "acc_apikey_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_apikey." + resourceName
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfigWithOnlyReqField(acctest.ProjectCfg, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsApiKeyResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					resource.TestCheckResourceAttr(resourceReference, "expiry", "180"),
					resource.TestCheckResourceAttr(resourceReference, "allowed_cidrs.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.1", "organizationOwner"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateApiKeyImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Rotate testing
			{
				Config: testAccApiKeyResourceConfigWithOnlyReqFieldRotate(acctest.ProjectCfg, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsApiKeyResource(resourceReference),
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
	resourceName := "acc_apikey_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_apikey." + resourceName
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfigForOrgOwner(acctest.ProjectCfg, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsApiKeyResource(resourceReference),
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
			//// ImportState testing
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
	resourceName := "acc_apikey_" + acctest.GenerateRandomResourceName()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccApiKeyResourceConfigRotateSet(acctest.ProjectCfg, resourceName),
				ExpectError: regexp.MustCompile("rotate value should not be set"),
			},
		},
	})
}

func testAccApiKeyResourceConfig(cfg, resourceName, projectResourceName, projectResourceReference string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_project" "terraform_api_test_project" {
    organization_id = var.organization_id
	name            = "terraform_api_test_project"
}
resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  description        = "description"
  expiry             = 150
  organization_roles = ["organizationMember"]
  allowed_cidrs      = ["10.1.42.1/23", "10.1.42.0/23"]
  resources = [
    {
      id    = couchbase-capella_project.terraform_api_test_project.id
      roles = ["projectManager", "projectDataReader"]
      type  = "project"
    }
  ]
}
`, cfg, resourceName)
}

func testAccApiKeyResourceConfigWithMultipleResources(cfg, resourceName, projectResourceName, projectResourceReference string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s_1" {
    organization_id = var.organization_id
	name            = "terraform_api_test_project"
}
resource "couchbase-capella_project" "%[3]s_2" {
    organization_id = var.organization_id
	name            = "terraform_api_test_project"
  depends_on = [couchbase-capella_project.%[3]s_1]

}


resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  organization_roles = ["organizationMember"]
  allowed_cidrs      = ["10.1.42.1/23", "10.1.42.0/23"]
  resources = [
    {
      id    = %[4]s_1.id
      roles = ["projectManager", "projectDataReader"]
      type  = "project"
    },
	{
      id    = %[4]s_2.id
      roles = ["projectDataReader"]
      type  = "project"
    }
  ]
}
`, cfg, resourceName, projectResourceName, projectResourceReference)
}

func testAccApiKeyResourceConfigWithOnlyReqField(cfg, resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  organization_roles = ["organizationOwner", "organizationMember"]
}
`, cfg, resourceName)
}

func testAccApiKeyResourceConfigForOrgOwner(cfg, resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_project" "terraform_api_test_project" {
    organization_id = var.organization_id
	name            = "terraform_api_test_project"
}
resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  organization_roles = [ "organizationMember"]
  resources = [
	  {
		id = couchbase-capella_project.terraform_api_test_project.id
		roles = [
		  "projectManager",
		  "projectDataReader"
		]
	  }
  ]
}
`, cfg, resourceName)
}

func testAccApiKeyResourceConfigRotateSet(cfg, resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_project" "terraform_api_test_project" {
    organization_id = var.organization_id
	name            = "terraform_api_test_project"
}
resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  organization_roles = [ "organizationMember"]
  resources = [
	  {
		id = couchbase-capella_project.terraform_api_test_project.id
		roles = [
		  "projectManager",
		  "projectDataReader"
		]
	  }
  ]
  rotate = 1
}
`, cfg, resourceName)
}

func testAccApiKeyResourceConfigWithOnlyReqFieldRotate(cfg, resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  organization_roles = ["organizationOwner", "organizationMember"]
  rotate             = 1
  secret             = "abc"
}
`, cfg, resourceName)
}

func testAccApiKeyResourceConfigWithoutResource(cfg, resourceName, projectResourceName, projectResourceReference string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  organization_roles = ["organizationOwner", "organizationMember"]
  allowed_cidrs      = ["10.1.42.0/23", "10.1.42.0/23"]
}
`, cfg, resourceName, projectResourceName, projectResourceReference)
}

func testAccExistsApiKeyResource(resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state

		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		fmt.Printf("raw state %s", rawState)
		data, err := acctest.TestClient()
		if err != nil {
			return err
		}
		_, err = retrieveApiKeyFromServer(data, rawState["organization_id"], rawState["id"])
		if err != nil {
			return err
		}
		return nil
	}
}

func retrieveApiKeyFromServer(data *providerschema.Data, organizationId, apiKeyId string) (*api.GetApiKeyResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/apikeys/%s", data.HostURL, organizationId, apiKeyId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.Client.Execute(
		cfg,
		nil,
		data.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}
	apiKeyResp := api.GetApiKeyResponse{}
	err = json.Unmarshal(response.Body, &apiKeyResp)
	if err != nil {
		return nil, err
	}
	return &apiKeyResp, nil
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
		fmt.Printf("raw state %s", rawState)
		return fmt.Sprintf("id=%s,organization_id=%s", rawState["id"], rawState["organization_id"]), nil
	}
}
