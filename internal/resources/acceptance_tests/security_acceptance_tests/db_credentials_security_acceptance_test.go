package security_acceptance_tests

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"terraform-provider-capella/internal/api"
	acctest "terraform-provider-capella/internal/testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOrgOwnerDatabaseCredentialTest(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationOwner")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				Config: testAccAddDatabaseCredWithReqFields(&testCfg),
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
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateOrgAPI("organizationMember") },
				Config:      testAccAddDatabaseCredWithReqFields(&testCfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccProjCreatorDatabaseCredentialTest(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateOrgAPI("projectCreator") },
				Config:      testAccAddDatabaseCredWithReqFields(&testCfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccProjOwnerDatabaseCredentialTest(t *testing.T) {

	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig: func() { testAccCreateProjAPI("organizationMember", projId, "projectOwner") },
				Config:    testAccAddDatabaseCredWithReqFields(&testCfg),
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
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateProjAPI("organizationMember", projId, "projectManager") },
				Config:      testAccAddDatabaseCredWithReqFields(&testCfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccProjViewerDatabaseCredentialTest(t *testing.T) {

	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateProjAPI("organizationMember", projId, "projectViewer") },
				Config:      testAccAddDatabaseCredWithReqFields(&testCfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccDatabaseReaderWriterDatabaseCredentialTest(t *testing.T) {

	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateProjAPI("organizationMember", projId, "projectDataReaderWriter") },
				Config:      testAccAddDatabaseCredWithReqFields(&testCfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccDatabaseReaderDatabaseCredentialTest(t *testing.T) {

	projId := os.Getenv("TF_VAR_project_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Database Credential with required fields
			{
				PreConfig:   func() { testAccCreateProjAPI("organizationMember", projId, "projectDataReader") },
				Config:      testAccAddDatabaseCredWithReqFields(&testCfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccAddDatabaseCredWithReqFields(cfg *string) string {
	*cfg = fmt.Sprintf(`
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

	`, *cfg)
	return *cfg
}

func testAccDeleteDatabaseCredential(clusterResourceReference, projectResourceReference, databaseCredentialResourceReference string) resource.TestCheckFunc {
	log.Println("deleting the database credential")
	return func(s *terraform.State) error {
		var clusterState, projectState, databaseCredentialState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[clusterResourceReference]; ok {
					clusterState = v.Primary.Attributes
				}
				if v, ok := m.Resources[projectResourceReference]; ok {
					projectState = v.Primary.Attributes
				}
				if v, ok := m.Resources[databaseCredentialResourceReference]; ok {
					databaseCredentialState = v.Primary.Attributes
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
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/users/%s", host, orgid, projectState["id"], clusterState["id"], databaseCredentialState["id"])
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
