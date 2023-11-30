package acceptance_tests

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

func TestAccDatabaseCredentialTestCases(t *testing.T) {
	resourceName := "new_cluster"
	resourceReference := "capella_cluster." + resourceName
	projectResourceName := "terraform_project"
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.250.250.0/23"

	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Creating cluster to check the database credential configs
			{
				Config: testAccCreateCluster(&testCfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
				),
			},
			// Database Credential with required fields
			{
				Config: testAccAddDatabaseCredWithReqFields(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("capella_allowlist.add_allowlist_req", "id"),
					resource.TestCheckResourceAttr("capella_database_credential.add_database_credential_req", "name", "acc_test_database_credential_name"),
					resource.TestCheckResourceAttr("capella_database_credential.add_database_credential_req", "access.priviledges.0", "data_writer"),
					resource.TestCheckResourceAttr("capella_database_credential.add_database_credential_req", "access.resources.buckets.0.name", "new_terraform_bucket"),
				),
			},
			// Database Credential with optional fields
			{
				Config: testAccAddDatabaseCredWithOptionalFields(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("capella_allowlist.add_allowlist_req", "id"),
					resource.TestCheckResourceAttr("capella_database_credential.add_database_credential_opt", "name", "acc_test_database_credential_name"),
					resource.TestCheckResourceAttr("capella_database_credential.add_database_credential_opt", "password", "acc_test_password"),
					resource.TestCheckResourceAttr("capella_database_credential.add_database_credential_req", "access.priviledges.0", "data_writer"),
					resource.TestCheckResourceAttr("capella_database_credential.add_database_credential_req", "access.resources.buckets.0.name", "new_terraform_bucket"),
				),
			},
			// Invalid name
			{
				Config:      testAccAddDatabaseCredWithInvalidName(&testCfg),
				ExpectError: regexp.MustCompile("Could not create database credential, unexpected error: The request was malformed or invalid."),
			},
		},
	})
}

// Attempt to delete the database credential when it has been already deleted through api
func TestAccAllowedDatabaseCredentialNotFound(t *testing.T) {
	clusterName := "new_cluster"
	clusterResourceReference := "capella_cluster." + clusterName
	projectResourceName := "terraform_project"
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.250.250.0/23"

	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateCluster(&testCfg, clusterName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(clusterResourceReference),
				),
			},
			{
				Config:             testAccAddIpWithOptionalFields(testCfg, "databaseCredential_delete", "10.2.3.4/32"),
				Check:              resource.ComposeAggregateTestCheckFunc(),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
		},
	})
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
		project_id      = capella_project.terraform_project.id
		cluster_id      = capella_cluster.new_cluster.id
		access = [
			{
				privileges = ["data_writer"]
				resources = {
				buckets = [{
					name = "new_terraform_bucket"
					scopes = [
					{
						name        = "_default"
						collections = ["_default"]
					}
					]
				}]
				}
			},
			{
				privileges = ["data_reader"]
			}
		]
	}
	
	`, *cfg)
	return *cfg
}

func testAccAddDatabaseCredWithOptionalFields(cfg *string) string {
	*cfg = fmt.Sprintf(`
	%[1]s
	
	output "add_database_credential_req"{
		value = capella_database_credential.add_database_credential_req
		sensitive = true
	}
	
	resource "capella_database_credential" "add_database_credential_req" {
		name            = "acc_test_database_credential_name"
		organization_id = var.organization_id
		project_id      = capella_project.terraform_project.id
		cluster_id      = capella_cluster.new_cluster.id
		password        = "acc_test_password"
		access = [
			{
				privileges = ["data_writer"]
				resources = {
				buckets = [{
					name = "new_terraform_bucket"
					scopes = [
					{
						name        = "_default"
						collections = ["_default"]
					}
					]
				}]
				}
			},
			{
				privileges = ["data_reader"]
			}
		]
	}
	
	`, *cfg)
	return *cfg
}

func testAccAddDatabaseCredWithInvalidName(cfg *string) string {
	*cfg = fmt.Sprintf(`
	%[1]s
	
	output "add_database_credential_req"{
		value = capella_database_credential.add_database_credential_req
		sensitive = true
	}
	
	resource "capella_database_credential" "add_database_credential_req" {
		name            = "acc_test_database_credential_invalid_name_="
		organization_id = var.organization_id
		project_id      = capella_project.terraform_project.id
		cluster_id      = capella_cluster.new_cluster.id
		password        = "acc_test_password"
		access          = "acc_test_access"
		access = [
			{
				privileges = ["data_writer"]
				resources = {
				buckets = [{
					name = "new_terraform_bucket"
					scopes = [
					{
						name        = "_default"
						collections = ["_default"]
					}
					]
				}]
				}
			},
			{
				privileges = ["data_reader"]
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
