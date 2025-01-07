package security

import (
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources/acceptance_tests"
	"os"

	"regexp"
	"testing"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNoAuth(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	os.Setenv("TF_VAR_auth_token", "")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccOrgAPIKeysConfig("organizationOwner"),
				ExpectError: regexp.MustCompile("empty value for the capella authentication token"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccAPIKeyRbacOrgOwner(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationOwner")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrgAPIKeysConfig("organizationOwner"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("capella_apikey.new_apikey", "id"),
					resource.TestCheckResourceAttrSet("capella_apikey.new_apikey", "token"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccAPIKeyRbacOrgMember(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationMember")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccOrgAPIKeysConfig("organizationOwner"),
				ExpectError: regexp.MustCompile("Could not create ApiKey"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccAPIKeyRbacProjCreator(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("projectCreator")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccOrgAPIKeysConfig("organizationOwner"),
				ExpectError: regexp.MustCompile("Could not create ApiKey"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccAPIKeyRbacProjOwner(t *testing.T) {
	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")

	testAccCreateProjAPI("organizationMember", projId, "projectOwner")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccOrgAPIKeysConfig("organizationOwner"),
				ExpectError: regexp.MustCompile("Could not create ApiKey"),
			},
			{
				Config: testAccProjAPIKeysConfig("organizationMember", projId, "projectViewer"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("capella_apikey.new_apikey", "id"),
					resource.TestCheckResourceAttrSet("capella_apikey.new_apikey", "token"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccAPIKeyRbacProjManager(t *testing.T) {
	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")

	testAccCreateProjAPI("organizationMember", projId, "projectManager")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccOrgAPIKeysConfig("organizationOwner"),
				ExpectError: regexp.MustCompile("Could not create ApiKey"),
			},
			{
				Config:      testAccProjAPIKeysConfig("organizationMember", projId, "projectViewer"),
				ExpectError: regexp.MustCompile("Could not create ApiKey"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccAPIKeyRbacProjViewer(t *testing.T) {
	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")

	testAccCreateProjAPI("organizationMember", projId, "projectManager")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccOrgAPIKeysConfig("organizationOwner"),
				ExpectError: regexp.MustCompile("Could not create ApiKey"),
			},
			{
				Config:      testAccProjAPIKeysConfig("organizationMember", projId, "projectViewer"),
				ExpectError: regexp.MustCompile("Could not create ApiKey"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccAPIKeyRbacDatabaseReaderWriter(t *testing.T) {
	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")

	testAccCreateProjAPI("organizationMember", projId, "projectDataReaderWriter")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccOrgAPIKeysConfig("organizationOwner"),
				ExpectError: regexp.MustCompile("Could not create ApiKey"),
			},
			{
				Config:      testAccProjAPIKeysConfig("organizationMember", projId, "projectViewer"),
				ExpectError: regexp.MustCompile("Could not create ApiKey"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccAPIKeyRbacDatabaseReader(t *testing.T) {
	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")

	testAccCreateProjAPI("organizationMember", projId, "projectDataReader")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccOrgAPIKeysConfig("organizationOwner"),
				ExpectError: regexp.MustCompile("Could not create ApiKey"),
			},
			{
				Config:      testAccProjAPIKeysConfig("organizationMember", projId, "projectViewer"),
				ExpectError: regexp.MustCompile("Could not create ApiKey"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccOrgAPIKeysConfig(organizationRole string) string {
	return fmt.Sprintf(`
%[1]s

output "new_apikey" {
	value     = capella_apikey.new_apikey
	sensitive = true
}

output "apikey_id" {
value = capella_apikey.new_apikey.id
}

resource "capella_apikey" "new_apikey" {
organization_id    = var.organization_id
name               = "Terraform Security"
description 	   = "APIKey to test Terraform Security Testing"
expiry 			   = 1
organization_roles = ["%[2]s"]
allowed_cidrs      = ["0.0.0.0/0"]
resources 		   = []
}

`, acceptance_tests.ProviderBlock, organizationRole)
}

func testAccProjAPIKeysConfig(organizationRole, projId, projectRole string) string {
	return fmt.Sprintf(`
%[1]s

output "new_apikey" {
	value     = capella_apikey.new_apikey
	sensitive = true
}

output "apikey_id" {
value = capella_apikey.new_apikey.id
}

resource "capella_apikey" "new_apikey" {
organization_id    = var.organization_id
name               = "Terraform Security"
description 	   = "APIKey to test Terraform Security Testing"
expiry 			   = 1
organization_roles = ["%[2]s"]
allowed_cidrs      = ["0.0.0.0/0"]
resources 		   = [
	{
		id = "%[3]s"
		roles = ["%[4]s"]
	}
]
}
`, acceptance_tests.ProviderBlock, organizationRole, projId, projectRole)
}
