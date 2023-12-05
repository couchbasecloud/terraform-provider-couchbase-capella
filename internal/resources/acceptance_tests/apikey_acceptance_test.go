package acceptance_tests_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"terraform-provider-capella/internal/api"
	providerschema "terraform-provider-capella/internal/schema"
	acctest "terraform-provider-capella/internal/testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccApiKeyResource(t *testing.T) {
	resourceName := "acc_apikey_" + acctest.GenerateRandomResourceName()
	resourceReference := "capella_apikey." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfig(acctest.Cfg, resourceName, projectResourceName, projectResourceReference),
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
	resourceReference := "capella_apikey." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfigWithMultipleResources(acctest.Cfg, resourceName, projectResourceName, projectResourceReference),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsApiKeyResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					resource.TestCheckResourceAttr(resourceReference, "expiry", "180"),
					resource.TestCheckResourceAttr(resourceReference, "allowed_cidrs.0", "10.1.42.0/23"),
					resource.TestCheckResourceAttr(resourceReference, "allowed_cidrs.1", "10.1.42.1/23"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "resources.1.roles.0", "projectDataReader"),
					resource.TestCheckResourceAttr(resourceReference, "resources.1.roles.1", "projectManager"),
					resource.TestCheckResourceAttr(resourceReference, "resources.1.type", "project"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.roles.0", "projectDataReader"),
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

func TestAccApiKeyResourceWithOnlyReqField(t *testing.T) {
	resourceName := "acc_apikey_" + acctest.GenerateRandomResourceName()
	resourceReference := "capella_apikey." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfigWithOnlyReqField(acctest.Cfg, resourceName),
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
				Config: testAccApiKeyResourceConfigWithOnlyReqFieldRotate(acctest.Cfg, resourceName),
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
	resourceReference := "capella_apikey." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfigForOrgOwner(acctest.Cfg, resourceName),
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
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccApiKeyResourceConfigRotateSet(acctest.Cfg, resourceName),
				ExpectError: regexp.MustCompile("rotate value should not be set"),
			},
		},
	})
}

func testAccApiKeyResourceConfig(cfg, resourceName, projectResourceName, projectResourceReference string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  description        = "description"
  expiry             = 150
  organization_roles = ["organizationMember"]
  allowed_cidrs      = ["10.1.42.1/23", "10.1.42.0/23"]
  resources = [
    {
      id    = var.project_id
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

resource "capella_project" "%[3]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  organization_roles = ["organizationMember"]
  allowed_cidrs      = ["10.1.42.1/23", "10.1.42.0/23"]
  resources = [
    {
      id    = %[4]s.id
      roles = ["projectManager", "projectDataReader"]
      type  = "project"
    },
	{
      id    = var.project_id
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

resource "capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  organization_roles = ["organizationOwner", "organizationMember"]
}
`, cfg, resourceName)
}

func testAccApiKeyResourceConfigForOrgOwner(cfg, resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  organization_roles = [ "organizationMember"]
  resources = [
	  {
		id = "1c50d827-cb90-49ca-a47e-dff850f53557"
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

resource "capella_apikey" "%[2]s" {
  organization_id    = var.organization_id
  name               = "%[2]s"
  organization_roles = [ "organizationMember"]
  resources = [
	  {
		id = "1c50d827-cb90-49ca-a47e-dff850f53557"
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

resource "capella_apikey" "%[2]s" {
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

resource "capella_project" "%[3]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "capella_apikey" "%[2]s" {
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
