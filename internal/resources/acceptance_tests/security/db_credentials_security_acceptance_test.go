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

func TestAccOrgOwnerDatabaseCredentialTest(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationOwner")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				Config: testAccAddDatabaseCredWithReqFields(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("capella_database_credential.add_database_credential_req", "id"),
					resource.TestCheckResourceAttr("capella_database_credential.add_database_credential_req", "name", "acc_test_database_credential_name"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccOrgMemberDatabaseCredentialTest(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateOrgAPI("organizationMember") },
				Config:      testAccAddDatabaseCredWithReqFields(),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccProjCreatorDatabaseCredentialTest(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateOrgAPI("projectCreator") },
				Config:      testAccAddDatabaseCredWithReqFields(),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccProjOwnerDatabaseCredentialTest(t *testing.T) {
	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig: func() { testAccCreateProjAPI("organizationMember", projId, "projectOwner") },
				Config:    testAccAddDatabaseCredWithReqFields(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("capella_database_credential.add_database_credential_req", "id"),
					resource.TestCheckResourceAttr("capella_database_credential.add_database_credential_req", "name", "acc_test_database_credential_name"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccProjManagerDatabaseCredentialTest(t *testing.T) {
	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateProjAPI("organizationMember", projId, "projectManager") },
				Config:      testAccAddDatabaseCredWithReqFields(),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccProjViewerDatabaseCredentialTest(t *testing.T) {
	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateProjAPI("organizationMember", projId, "projectViewer") },
				Config:      testAccAddDatabaseCredWithReqFields(),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccDatabaseReaderWriterDatabaseCredentialTest(t *testing.T) {
	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateProjAPI("organizationMember", projId, "projectDataReaderWriter") },
				Config:      testAccAddDatabaseCredWithReqFields(),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccDatabaseReaderDatabaseCredentialTest(t *testing.T) {
	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateProjAPI("organizationMember", projId, "projectDataReader") },
				Config:      testAccAddDatabaseCredWithReqFields(),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccAddDatabaseCredWithReqFields() string {
	return fmt.Sprintf(`
	%[1]s

	output "add_database_credential_req"{
		value = capella_database_credential.add_database_credential_req
		sensitive = true
	}

	resource "capella_database_credential" "add_database_credential_req" {
		name            = "acc_test_database_credential_name"
		organization_id = var.organization_id
		project_id      = var.project_id
		cluster_id      = var.cluster_id
		access = [
			{
				privileges = ["data_writer"]
			}
		]
	}

	`, acceptance_tests.ProviderBlock)
}
