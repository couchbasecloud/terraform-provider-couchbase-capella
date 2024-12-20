package security

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"regexp"
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
	cfg "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type OrganizationAPIPayload struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	AllowedCIDRs      []string `json:"allowedCIDRs"`
	OrganizationRoles []string `json:"organizationRoles"`
	Resources         []any    `json:"resources"`
	Expiry            int      `json:"expiry"`
}

type ProjectAPIPayload struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	AllowedCIDRs      []string `json:"allowedCIDRs"`
	OrganizationRoles []string `json:"organizationRoles"`
	Resources         []struct {
		ID    string   `json:"id"`
		Roles []string `json:"roles"`
	} `json:"resources"`
	Expiry int `json:"expiry"`
}

type ResponseStruct struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func TestAccOrganizationDataSourceNoAuth(t *testing.T) {
	organizationId := os.Getenv("TF_VAR_organization_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	os.Setenv("TF_VAR_auth_token", "")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccOrganizationResourceConfig(cfg.Cfg, organizationId),
				ExpectError: regexp.MustCompile("Missing Capella Authentication Token"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccOrganizationDataSourceRbacOrgOwner(t *testing.T) {
	organizationId := os.Getenv("TF_VAR_organization_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationOwner")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationResourceConfig(cfg.Cfg, organizationId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "name"),
					resource.TestCheckResourceAttr("data.capella_organization.get_organization", "organization_id", organizationId),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.created_at"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.modified_by"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.modified_at"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.version"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccOrganizationDataSourceRbacOrgMember(t *testing.T) {
	organizationId := os.Getenv("TF_VAR_organization_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationMember")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationResourceConfig(cfg.Cfg, organizationId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "name"),
					resource.TestCheckResourceAttr("data.capella_organization.get_organization", "organization_id", organizationId),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.created_at"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.modified_by"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.modified_at"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.version"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccOrganizationDataSourceRbacProjCreator(t *testing.T) {
	organizationId := os.Getenv("TF_VAR_organization_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("projectCreator")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationResourceConfig(cfg.Cfg, organizationId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "name"),
					resource.TestCheckResourceAttr("data.capella_organization.get_organization", "organization_id", organizationId),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.created_at"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.modified_by"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.modified_at"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.version"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccOrganizationResourceConfig(cfg string, organizationId string) string {
	return fmt.Sprintf(`
%[1]s

output "organizations_get" {
  value = data.capella_organization.get_organization
}

data "capella_organization" "get_organization" {
  organization_id = "%[2]s"
}

`, cfg, organizationId)
}

func testAccCreateOrgAPI(role string) string {
	log.Println("Creating api keys for Organization role")
	data, err := acctest.TestClient()
	if err != nil {
		fmt.Println("Error in Test Client")
	}

	host := os.Getenv("TF_VAR_host")
	orgid := os.Getenv("TF_VAR_organization_id")
	endPointCfg := api.EndpointCfg{
		Url:           fmt.Sprintf("%s/v4/organizations/%s/apikeys", host, orgid),
		Method:        http.MethodPost,
		SuccessStatus: 201,
	}
	payload := OrganizationAPIPayload{
		Name:              "Organization Owner API Key",
		Description:       "Creates an API key with a Organization Owner role.",
		Expiry:            720,
		AllowedCIDRs:      []string{"0.0.0.0/0"},
		OrganizationRoles: []string{role},
		Resources:         []any{},
	}
	resp, err := data.Client.Execute(
		endPointCfg,
		payload,
		data.Token,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	var responseStruct ResponseStruct
	err = json.Unmarshal(resp.Body, &responseStruct)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return ""
	}

	// Access the extracted value
	extractedValue := responseStruct.Token
	fmt.Println("Extracted Value:", extractedValue)
	os.Setenv("TF_VAR_auth_token", extractedValue)
	return extractedValue
}

func testAccCreateProjAPI(orgRole string, projId string, projRole string) string {
	log.Println("Creating api keys for Project role")
	data, err := acctest.TestClient()
	if err != nil {
		fmt.Println("Error in Test Client")
	}

	host := os.Getenv("TF_VAR_host")
	orgid := os.Getenv("TF_VAR_organization_id")
	authToken := os.Getenv("TF_VAR_auth_token")
	endPointCfg := api.EndpointCfg{
		Url:           fmt.Sprintf("%s/v4/organizations/%s/apikeys", host, orgid),
		Method:        http.MethodPost,
		SuccessStatus: 201,
	}
	payload := ProjectAPIPayload{
		Name:              "Organization Owner API Key",
		Description:       "Creates an API key with a Organization Owner role.",
		Expiry:            720,
		AllowedCIDRs:      []string{"0.0.0.0/0"},
		OrganizationRoles: []string{orgRole},
		Resources: []struct {
			ID    string   `json:"id"`
			Roles []string `json:"roles"`
		}{
			{ID: projId,
				Roles: []string{projRole}},
		},
	}
	resp, err := data.Client.Execute(
		endPointCfg,
		payload,
		authToken,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	var responseStruct ResponseStruct
	err = json.Unmarshal(resp.Body, &responseStruct)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return ""
	}

	// Access the extracted value
	extractedValue := responseStruct.Token
	fmt.Println("Extracted Value:", extractedValue)
	os.Setenv("TF_VAR_auth_token", extractedValue)
	return extractedValue
}
