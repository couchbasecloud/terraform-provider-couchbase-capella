package acceptance_tests

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDatabaseCredentialTestCases(t *testing.T) {
	resourceName := "acc_database_credential_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_database_credential." + resourceName

	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName

	clusterResourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	cidr := "10.1.116.0/23"

	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Creating cluster to check the database credential configs
			{
				Config: testAccDatabaseCredentialCreateCluster(&testCfg, clusterResourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(clusterResourceReference),
				),
			},
			// Database Credential with required fields
			{
				Config: testAccAddDatabaseCredWithReqFields(testCfg, clusterResourceReference, resourceName, projectResourceReference),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_database_credential_name"),
					resource.TestCheckResourceAttr(resourceReference, "access.priviledges.0", "data_writer"),
					resource.TestCheckResourceAttr(resourceReference, "access.resources.buckets.0.name", "new_terraform_bucket"),
				),
			},
			// Database Credential with optional fields
			{
				Config: testAccAddDatabaseCredWithOptionalFields(testCfg, clusterResourceReference, resourceName, projectResourceReference),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_database_credential_name"),
					resource.TestCheckResourceAttr(resourceReference, "password", "acc_test_password"),
					resource.TestCheckResourceAttr(resourceReference, "access.priviledges.0", "data_writer"),
					resource.TestCheckResourceAttr(resourceReference, "access.resources.buckets.0.name", "new_terraform_bucket"),
				),
			},
			// Invalid name
			{
				Config:      testAccAddDatabaseCredWithInvalidName(&testCfg, clusterResourceReference, resourceName, projectResourceReference),
				ExpectError: regexp.MustCompile("Could not create database credential, unexpected error: The request was malformed or invalid."),
			},
		},
	})
}

// Attempt to delete the database credential when it has been already deleted through api
func TestAccAllowedDatabaseCredentialNotFound(t *testing.T) {
	resourceName := "acc_database_credential_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_database_credential." + resourceName

	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName

	clusterResourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	cidr := "10.1.116.0/23"

	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateCluster(&testCfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
				),
			},
			{
				Config: testAccAddIpWithOptionalFields(testCfg, "database_credential_delete", "10.2.3.4/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_database_credential_name"),
					resource.TestCheckResourceAttr(resourceReference, "password", "acc_test_password"),
					resource.TestCheckResourceAttr(resourceReference, "access.priviledges.0", "data_writer"),
					resource.TestCheckResourceAttr(resourceReference, "access.resources.buckets.0.name", "new_terraform_bucket"),
					testAccDeleteDatabaseCredential(clusterResourceReference, projectResourceReference, "couchbase-capella_database_credential.database_credential_delete"),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
		},
	})
}

func testAccAddDatabaseCredWithReqFields(cfg, clusterReference, resourceName, projectReference string) string {
	return fmt.Sprintf(
		`
		%[1]s
	
		output "add_database_credential_req"{
			value = couchbase-capella_database_credential.add_database_credential_req
			sensitive = true
		}
		
		resource "couchbase-capella_database_credential" "%[2]s" {
			name            = "acc_test_database_credential_name"
			organization_id = var.organization_id
			project_id      = %[3]s.id
			cluster_id      = %[2]s.id
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
		`, cfg, clusterReference, resourceName, projectReference)
}

func testAccAddDatabaseCredWithOptionalFields(cfg, clusterReference, resourceName, projectReference string) string {
	return fmt.Sprintf(
		`
		%[1]s
	
		output "add_database_credential_req"{
			value = couchbase-capella_database_credential.add_database_credential_req
			sensitive = true
		}
		
		resource "couchbase-capella_database_credential" "%[2]s" {
			name            = "acc_test_database_credential_name"
			organization_id = var.organization_id
			project_id      = %[3]s.id
			cluster_id      = %[2]s.id
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
		`, cfg, clusterReference, resourceName, projectReference)
}

func testAccAddDatabaseCredWithInvalidName(cfg *string, clusterReference, resourceName, projectReference string) string {
	*cfg = fmt.Sprintf(`
	%[1]s
	
	output "add_database_credential_req"{
		value = couchbase-capella_database_credential.add_database_credential_req
		sensitive = true
	}
	
	resource "couchbase-capella_database_credential" "add_database_credential_req" {
		name            = "acc_test_database_credential_invalid_name_="
		organization_id = var.organization_id
		project_id      = couchbase-capella_project.terraform_project.id
		cluster_id      = couchbase-capella_cluster.new_cluster.id
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

func testAccDatabaseCredentialCreateCluster(cfg *string, resourceName, projectResourceName, projectResourceReference, cidr string) string {
	log.Println("Creating cluster")
	*cfg = fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "couchbase-capella_cluster" "%[2]s" {
  organization_id = var.organization_id
  project_id      = %[4]s.id
  name            = "Terraform Acceptance Test Cluster"
  description     = "terraform acceptance test cluster"
  couchbase_server = {
    version = "7.1"
  }
  configuration_type = "multiNode"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 50
          type    = "gp3"
          iops    = 3000
        }
      }
      num_of_nodes = 2
      services     = ["index", "query"]
    },
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 50
          type    = "gp3"
          iops    = 3000
        }
      }
      num_of_nodes = 3
      services     = ["data"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "developer pro"
    timezone = "PT"
  }
}
`, *cfg, resourceName, projectResourceName, projectResourceReference, cidr)
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
