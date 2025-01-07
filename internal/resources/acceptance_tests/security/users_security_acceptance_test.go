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

func TestAccCreateUserNoAuth(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	os.Setenv("TF_VAR_auth_token", "")
	name := "terraform_security"
	email := "koushal.sharma+10@couchbase.com"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccUsersResourceConfigRequired(name, email),
				ExpectError: regexp.MustCompile("Missing Capella Authentication Token"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateUserOrgOwner(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationOwner")
	name := "terraform_security"
	email := "koushal.sharma+10@couchbase.com"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUsersResourceConfigRequired(name, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_user.new_user", "name", name),
					resource.TestCheckResourceAttr("capella_user.new_user", "email", email),
					resource.TestCheckResourceAttrSet("capella_user.new_user", "id"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateUserOrgMember(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationMember")
	name := "terraform_security"
	email := "koushal.sharma+10@couchbase.com"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccUsersResourceConfigRequired(name, email),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateUserProjCreator(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("projectCreator")
	name := "terraform_security"
	email := "koushal.sharma+10@couchbase.com"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccUsersResourceConfigRequired(name, email),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateUserProjOwner(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectOwner")
	name := "terraform_security"
	email := "koushal.sharma+10@couchbase.com"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccUsersResourceConfigRequired(name, email),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateUserProjManager(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectManager")
	name := "terraform_security"
	email := "koushal.sharma+10@couchbase.com"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccUsersResourceConfigRequired(name, email),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateUserProjViewer(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectViewer")
	name := "terraform_security"
	email := "koushal.sharma+10@couchbase.com"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccUsersResourceConfigRequired(name, email),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateUserDatabaseReaderWriter(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectDataReaderWriter")
	name := "terraform_security"
	email := "koushal.sharma+10@couchbase.com"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccUsersResourceConfigRequired(name, email),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateUserDatabaseReader(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectDataReader")
	name := "terraform_security"
	email := "koushal.sharma+10@couchbase.com"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccUsersResourceConfigRequired(name, email),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccUsersResourceConfigRequired(name, email string) string {
	return fmt.Sprintf(`
%[1]s

output "new_user" {
	value = capella_user.new_user
}

output "user_id" {
	value = capella_user.new_user.id
}

resource "capella_user" "new_user" {
	organization_id = var.organization_id
	name = "%[2]s"
	email = "%[3]s"
	organization_roles = ["organizationMember"]
	resources = [
		{
			type = "project"
			id = var.project_id
			roles = ["projectOwner"]
		}
	]
}

`, acceptance_tests.ProviderBlock, name, email)
}
