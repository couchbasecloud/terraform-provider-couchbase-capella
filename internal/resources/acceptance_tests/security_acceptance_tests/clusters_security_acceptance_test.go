package security_acceptance_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCreateClusterNoAuth(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.252.250.0/23"
	tempId := os.Getenv("TF_VAR_auth_token")
	os.Setenv("TF_VAR_auth_token", "")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccClusterResourceConfig(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				ExpectError: regexp.MustCompile("Missing Capella Authentication Token"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateClusterOrgOwner(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.252.250.0/23"
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationOwner")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccClusterResourceConfig(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster"),
					resource.TestCheckResourceAttr(resourceReference, "description", "My first test cluster for multiple services."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.1"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "PT"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateClusterOrgMember(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.254.250.0/23"
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationMember")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccClusterResourceConfig(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateClusterProjCreator(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.254.250.0/23"
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("projectCreator")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccClusterResourceConfig(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateClusterProjOwner(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.252.250.0/23"
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectOwner")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccClusterResourceConfig(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster"),
					resource.TestCheckResourceAttr(resourceReference, "description", "My first test cluster for multiple services."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.1"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "PT"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateClusterProjManager(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.242.250.0/23"
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectManager")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccClusterResourceConfig(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster"),
					resource.TestCheckResourceAttr(resourceReference, "description", "My first test cluster for multiple services."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.1"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "PT"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateClusterProjViewer(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.256.250.0/23"
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectViewer")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccClusterResourceConfig(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateClusterDatabaseReaderWriter(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.256.250.0/23"
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectDataReaderWriter")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccClusterResourceConfig(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateClusterDatabaseReader(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.256.250.0/23"
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectDataReader")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccClusterResourceConfig(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccClusterResourceConfig(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_cluster" "%[2]s" {
  organization_id = var.organization_id
  project_id      = var.project_id
  name            = "Terraform Acceptance Test Cluster"
  description     = "My first test cluster for multiple services."
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
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}

// retrieveClusterFromServer checks cluster exists in server.
func retrieveClusterFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) (*clusterapi.GetClusterResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", data.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.Client.Execute(
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

func testAccExistsClusterResource(resourceReference string) resource.TestCheckFunc {
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
		_, err = retrieveClusterFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["id"])
		if err != nil {
			return err
		}
		return nil
	}
}
