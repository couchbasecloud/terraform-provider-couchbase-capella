package security

import (
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources/acceptance_tests"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCreateProjectNoAuth(t *testing.T) {
	rnd := acceptance_tests.RandomStringWithPrefix("tf_acc_project_")
	tempId := os.Getenv("TF_VAR_auth_token")
	os.Setenv("TF_VAR_auth_token", "")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(rnd),
				ExpectError: regexp.MustCompile("Missing Capella Authentication Token"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectOrgOwner(t *testing.T) {
	rnd := acceptance_tests.RandomStringWithPrefix("tf_acc_project_")
	resourceName := "capella_project." + rnd
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationOwner")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfigRequired(rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test project"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 1"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectProjCreator(t *testing.T) {
	rnd := acceptance_tests.RandomStringWithPrefix("tf_acc_project_")
	resourceName := "capella_project." + rnd
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("projectCreator")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfigRequired(rnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform acceptance test project"),
					resource.TestCheckResourceAttr(resourceName, "etag", "Version: 1"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectOrgMember(t *testing.T) {
	rnd := acceptance_tests.RandomStringWithPrefix("tf_acc_project_")
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationMember")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectProjOwner(t *testing.T) {
	rnd := acceptance_tests.RandomStringWithPrefix("tf_acc_project_")
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectOwner")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectProjManager(t *testing.T) {
	rnd := acceptance_tests.RandomStringWithPrefix("tf_acc_project_")
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectManager")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectProjViewer(t *testing.T) {
	rnd := acceptance_tests.RandomStringWithPrefix("tf_acc_project_")
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectViewer")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectDatabaseDataReaderWriter(t *testing.T) {
	rnd := acceptance_tests.RandomStringWithPrefix("tf_acc_project_")
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectDataReaderWriter")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectDatabaseDataReader(t *testing.T) {
	rnd := acceptance_tests.RandomStringWithPrefix("tf_acc_project_")
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectDataReader")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccProjectResourceConfigRequired(rnd string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "%[2]s" {
    organization_id = var.organization_id
	name            = "%[2]s"
	description     = "terraform acceptance test project"
}
`, acceptance_tests.ProviderBlock, rnd)
}
