package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"testing"

	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDatabaseCredentialTestCases(t *testing.T) {
	resourceName := "new_cluster"
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "terraform_project"
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr := "10.1.42.0/23"

	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Creating cluster to check the database_credential configs
			{
				Config: testAccDatabaseCredentialCreateCluster(&testCfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccDatabaseCredentialExistsClusterResource(resourceReference),
				),
			},
			//database_credential with required fields
			{
				Config: testAccAddDatabaseCredWithReqFields(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_database_credential.add_database_credential_req", "name", "acc_test_database_credential_name"),
					resource.TestCheckResourceAttr("couchbase-capella_database_credential.add_database_credential_req", "access.0.privileges.0", "data_writer"),
				),
			},
			//database_credential with optional fields
			{
				Config: testAccAddDatabaseCredWithOptionalFields(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_database_credential.add_database_credential_opt", "name", "acc_test_database_credential_name2"),
					resource.TestCheckResourceAttr("couchbase-capella_database_credential.add_database_credential_opt", "password", "Secret12$#"),
					resource.TestCheckResourceAttr("couchbase-capella_database_credential.add_database_credential_opt", "access.0.privileges.0", "data_writer"),
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

// Delete the database_credential when the cluster is destroyed through api
func testAccAddDatabaseCredWithReqFields(cfg *string) string {
	return fmt.Sprintf(
		`
		%[1]s
	
		output "add_database_credential_req"{
			value = couchbase-capella_database_credential.add_database_credential_req
			sensitive = true
		}
		
		resource "couchbase-capella_database_credential" "add_database_credential_req" {
			name            = "acc_test_database_credential_name"
			organization_id = var.organization_id
			project_id      = couchbase-capella_project.terraform_project.id
			cluster_id      = couchbase-capella_cluster.new_cluster.id
			access = [
				{
					privileges = ["data_writer"]
				},
			]
		}
		`, *cfg)
}

func testAccAddDatabaseCredWithOptionalFields(cfg *string) string {
	return fmt.Sprintf(
		`
		%[1]s
	
		output "add_database_credential_opt"{
			value = couchbase-capella_database_credential.add_database_credential_opt
			sensitive = true
		}
		
		resource "couchbase-capella_database_credential" "add_database_credential_opt" {
			name            = "acc_test_database_credential_name2"
			organization_id = var.organization_id
			project_id      = couchbase-capella_project.terraform_project.id
			cluster_id      = couchbase-capella_cluster.new_cluster.id
			password        = "Secret12$#"
			access = [
				{
					privileges = ["data_writer"]
				},
			]
		}
		`, *cfg)
}

func testAccAddDatabaseCredWithInvalidName(cfg *string) string {
	*cfg = fmt.Sprintf(`
	%[1]s
	
	output "add_database_credential_invalid_name"{
		value = couchbase-capella_database_credential.add_database_credential_invalid_name
		sensitive = true
	}
	
	resource "couchbase-capella_database_credential" "add_database_credential_invalid_name" {
		name            = "acc_test_database_credential_invalid_name3"
		organization_id = var.organization_id
		project_id      = couchbase-capella_project.terraform_project.id
		cluster_id      = couchbase-capella_cluster.new_cluster.id
		password        = "Secret12$#"
		access = [
			{
				privileges = ["data_writer"]
			},
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
  name            = "terraform database credential acceptance test cluster"
  description     = "terraform database credential acceptance test cluster"
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
      num_of_nodes = 3
      services     = ["data"]
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

func testAccDatabaseCredentialExistsClusterResource(resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state

		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		fmt.Printf("raw state %s", rawState)
		data, err := acctest.TestClient()
		if err != nil {
			return err
		}
		_, err = retrieveDatabaseCredentialClusterFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["id"])
		if err != nil {
			return err
		}
		return nil
	}
}

func retrieveDatabaseCredentialClusterFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) (*clusterapi.GetClusterResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", data.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.Client.ExecuteWithRetry(
		context.Background(),
		cfg,
		nil,
		data.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	clusterResp := clusterapi.GetClusterResponse{}
	err = json.Unmarshal(response.Body, &clusterResp)
	if err != nil {
		return nil, err
	}
	clusterResp.Etag = response.Response.Header.Get("ETag")
	return &clusterResp, nil
}
