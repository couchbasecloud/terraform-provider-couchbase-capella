package acceptance_tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
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
					resource.TestCheckTypeSetElemAttr(resourceReference, "allowed_cidrs.*", "10.1.42.0/23"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "allowed_cidrs.*", "10.1.42.1/23"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "organization_roles.*", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "resources.0.roles.*", "projectDataReader"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "resources.0.roles.*", "projectManager"),
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
					resource.TestCheckTypeSetElemAttr(resourceReference, "allowed_cidrs.*", "0.0.0.0/0"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "organization_roles.*", "organizationMember"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "organization_roles.*", "organizationOwner"),
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
					resource.TestCheckTypeSetElemAttr(resourceReference, "allowed_cidrs.*", "0.0.0.0/0"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "organization_roles.*", "organizationMember"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "organization_roles.*", "organizationOwner"),
					resource.TestCheckResourceAttr(resourceReference, "rotate", "1"),
					resource.TestCheckResourceAttrSet(resourceReference, "secret"),
					resource.TestCheckResourceAttrSet(resourceReference, "token"),
				),
			},
		},
	})
}

func TestAccApiKeyResourceWithDefaults(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resourceReference := "couchbase-capella_apikey." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccApiKeyResourceConfigWithDefaults(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					resource.TestCheckResourceAttr(resourceReference, "expiry", "180"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "allowed_cidrs.*", "0.0.0.0/0"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "organization_roles.*", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.#", "0"),
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
					resource.TestCheckTypeSetElemAttr(resourceReference, "allowed_cidrs.*", "0.0.0.0/0"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "organization_roles.*", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "resources.0.roles.*", "projectDataReader"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "resources.0.roles.*", "projectManager"),
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
				ExpectError: regexp.MustCompile(`(?s)Invalid API Key Configuration.*rotate value should not be set during create`),
			},
		},
	})
}

func TestAccApiKeyResourceInvalidScenarioSecretShouldNotPassedWhileCreate(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccApiKeyResourceConfigSecretSetConfig(resourceName),
				ExpectError: regexp.MustCompile(`(?s)Invalid API Key Configuration.*secret can only be configured together with rotate`),
			},
		},
	})
}

func TestAccApiKeyResourceInvalidScenarioEmptyAllowedCIDRs(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccApiKeyResourceConfigWithEmptyAllowedCIDRs(resourceName),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value.*Attribute allowed_cidrs set must contain at least 1 elements, got: 0`),
			},
		},
	})
}

func TestAccApiKeyResourceInvalidScenarioEmptyOrganizationID(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccApiKeyResourceConfigWithOrganizationID(resourceName, ""),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value.*Attribute organization_id string length must be at least 1, got: 0`),
			},
		},
	})
}

func TestAccApiKeyResourceInvalidScenarioEmptyOrganizationRoles(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccApiKeyResourceConfigWithOrganizationRoles(resourceName, "[]"),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value.*Attribute organization_roles set must contain at least 1 elements, got: 0`),
			},
		},
	})
}

func TestAccApiKeyResourceInvalidScenarioInvalidOrganizationRole(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccApiKeyResourceConfigWithOrganizationRoles(resourceName, `["organizationWizard"]`),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.*organizationWizard.*organizationMember.*organizationOwner.*projectCreator`),
			},
		},
	})
}

func TestAccApiKeyResourceInvalidScenarioInvalidProjectRole(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccApiKeyResourceConfigWithProjectRole(resourceName, "projectWizard"),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.*projectWizard.*projectOwner.*projectManager.*projectViewer.*projectDataReaderWriter.*projectDataReader`),
			},
		},
	})
}

func TestAccApiKeyResourceInvalidScenarioResourceMissingRoles(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccApiKeyResourceConfigWithResourceMissingRoles(resourceName),
				ExpectError: regexp.MustCompile(`(?s)Incorrect attribute value type.*attribute "roles".*is required`),
			},
		},
	})
}

func TestAccApiKeyResourceInvalidScenarioResourceIdNotUUID(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccApiKeyResourceConfigWithResourceId(resourceName, "not-a-uuid"),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.*not-a-uuid.*resources.id must be a valid UUID`),
			},
		},
	})
}

func TestAccApiKeyResourceInvalidScenarioResourceIdEmpty(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikey_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccApiKeyResourceConfigWithResourceId(resourceName, ""),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Length.*Attribute\s+resources.*\.id\s+string length must be at least 1, got: 0`),
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

func testAccApiKeyResourceConfigWithDefaults(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
	organization_id    = "%[3]s"
	name               = "%[2]s"
	organization_roles = ["organizationMember"]
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

func testAccApiKeyResourceConfigSecretSetConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
	organization_id    = "%[3]s"
	name               = "%[2]s"
	organization_roles = ["organizationMember"]
	secret             = "abc"
}
`, globalProviderBlock, resourceName, globalOrgId)
}

func testAccApiKeyResourceConfigWithEmptyAllowedCIDRs(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
	organization_id    = "%[3]s"
	name               = "%[2]s"
	organization_roles = ["organizationMember"]
	allowed_cidrs      = []
}
`, globalProviderBlock, resourceName, globalOrgId)
}

func testAccApiKeyResourceConfigWithOrganizationID(resourceName, organizationID string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
	organization_id    = "%[3]s"
	name               = "%[2]s"
	organization_roles = ["organizationMember"]
}
`, globalProviderBlock, resourceName, organizationID)
}

func testAccApiKeyResourceConfigWithOrganizationRoles(resourceName, organizationRoles string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
	organization_id    = "%[3]s"
	name               = "%[2]s"
	organization_roles = %[4]s
}
`, globalProviderBlock, resourceName, globalOrgId, organizationRoles)
}

func testAccApiKeyResourceConfigWithProjectRole(resourceName, projectRole string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
	organization_id    = "%[3]s"
	name               = "%[2]s"
	organization_roles = ["organizationMember"]
	resources = [
		{
			id    = "%[4]s"
			roles = ["%[5]s"]
			type  = "project"
		}
	]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, projectRole)
}

func testAccApiKeyResourceConfigWithResourceMissingRoles(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
	organization_id    = "%[3]s"
	name               = "%[2]s"
	organization_roles = ["organizationMember"]
	resources = [
		{
			id   = "%[4]s"
			type = "project"
		}
	]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId)
}

func testAccApiKeyResourceConfigWithResourceId(resourceName, resourceId string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
	organization_id    = "%[3]s"
	name               = "%[2]s"
	organization_roles = ["organizationMember"]
	resources = [
		{
			id    = "%[4]s"
			roles = ["projectViewer"]
			type  = "project"
		}
	]
}
`, globalProviderBlock, resourceName, globalOrgId, resourceId)
}

func testAccApiKeyResourceConfigWithOnlyReqFieldRotateConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[2]s" {
  organization_id    = "%[3]s"
  name               = "%[2]s"
  organization_roles = ["organizationOwner", "organizationMember"]
  rotate             = 1
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

// TestAccApiKeyResourceInvalidExpiryDirectAPI hits the API directly with an
// invalid expiry (0) to verify CBSE-23222 / AV-137404: the API must return a
// proper validation error (HTTP 422) with an expiry-specific message, not a
// generic 500/code 10000.
func TestAccApiKeyResourceInvalidExpiryDirectAPI(t *testing.T) {
	data := newTestClient(t)

	name := randomStringWithPrefix("tf_acc_apikey_expiry_direct_")
	expiry := float32(0)
	req := api.CreateApiKeyRequest{
		Name:              name,
		OrganizationRoles: []string{"organizationMember"},
		Expiry:            &expiry,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/apikeys", data.HostURL, globalOrgId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}

	_, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, req, data.Token, nil)
	if err == nil {
		t.Fatal("unexpected success: API accepted expiry=0, but it should be rejected")
	}

	var apiErr *api.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected api.Error, got %T: %v", err, err)
	}

	if apiErr.HttpStatusCode == http.StatusInternalServerError {
		t.Fatalf(
			"CBSE-23222 regression: API returns 500 for invalid expiry (0). "+
				"Expected HTTP 422 with message about expiry being too small. Got: %s",
			apiErr.Error(),
		)
	}
	if apiErr.HttpStatusCode != http.StatusUnprocessableEntity {
		t.Fatalf(
			"unexpected HTTP status for invalid expiry: got %d, want %d. Body: %s",
			apiErr.HttpStatusCode, http.StatusUnprocessableEntity, apiErr.Error(),
		)
	}
	if apiErr.Code != 422 {
		t.Fatalf("unexpected error code: got %d, want 422. Body: %s", apiErr.Code, apiErr.Error())
	}
	if !strings.Contains(apiErr.Message, "expiry") {
		t.Fatalf("expected expiry validation message, got: %s", apiErr.Message)
	}
	if apiErr.Hint == "" {
		t.Fatal("expected a non-empty error hint about required parameters")
	}
}
