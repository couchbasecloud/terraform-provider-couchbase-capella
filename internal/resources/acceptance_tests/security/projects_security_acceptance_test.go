package security

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
)

func TestAccCreateProjectNoAuth(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	tempId := os.Getenv("TF_VAR_auth_token")
	os.Setenv("TF_VAR_auth_token", "")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(acctest.ProjectCfg, rnd),
				ExpectError: regexp.MustCompile("Missing Capella Authentication Token"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectOrgOwner(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	resourceName := "capella_project." + rnd
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationOwner")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfigRequired(acctest.ProjectCfg, rnd),
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
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	resourceName := "capella_project." + rnd
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("projectCreator")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfigRequired(acctest.ProjectCfg, rnd),
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
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationMember")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(acctest.ProjectCfg, rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectProjOwner(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectOwner")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(acctest.ProjectCfg, rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectProjManager(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectManager")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(acctest.ProjectCfg, rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectProjViewer(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectViewer")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(acctest.ProjectCfg, rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectDatabaseDataReaderWriter(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectDataReaderWriter")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(acctest.ProjectCfg, rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateProjectDatabaseDataReader(t *testing.T) {
	rnd := "acc_project_" + acctest.GenerateRandomResourceName()
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectDataReader")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccProjectResourceConfigRequired(acctest.ProjectCfg, rnd),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccProjectResourceConfigRequired(cfg string, rnd string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "%[2]s" {
    organization_id = var.organization_id
	name            = "%[2]s"
	description     = "terraform acceptance test project"
}
`, cfg, rnd)
}
