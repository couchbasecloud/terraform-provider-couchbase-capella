package acceptance_tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	clusterapi "terraform-provider-capella/internal/api/cluster"
	cfg "terraform-provider-capella/internal/testing"

	//acctest "terraform-provider-capella/internal/resources/acceptance_tests"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
	"time"
)

func TestAccAllowListTestCases(t *testing.T) {
	resourceName := "new_cluster"
	resourceReference := "capella_cluster." + resourceName
	projectResourceName := "terraform_project"
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.0.2.0/23"

	testCfg := cfg.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//Creating cluster to check the allowlist configs
			{
				Config: testAccCreateCluster(&testCfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
				),
			},
			//IP with required fields
			{
				Config: testAccAddIpWithReqFields(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_allowlist.add_allowlist_req", "cidr", "10.1.1.1/32"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_allowlist_req", "id"),
				),
			},
			//IP with optional fields
			{
				Config: testAccAddIpWithOptionalFields(testCfg, "add_allowlist_opt", "10.4.5.6/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_allowlist.add_allowlist_opt", "cidr", "10.4.5.6/32"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_allowlist_opt", "id"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_allowlist_opt", "expires_at"),
					resource.TestCheckResourceAttr("capella_allowlist.add_allowlist_opt", "comment", "terraform allow list acceptance test"),
				),
			},
			//Unspecified IP address
			{
				Config: testAccAddIpWithOptionalFields(testCfg, "add_allowlist_quadzero", "0.0.0.0/0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_allowlist.add_allowlist_quadzero", "cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_allowlist_quadzero", "id"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_allowlist_quadzero", "expires_at"),
					resource.TestCheckResourceAttr("capella_allowlist.add_allowlist_quadzero", "comment", "terraform allow list acceptance test"),
				),
			},
			//expired IP
			{
				Config:      testAccAddIpWithExpiredIP(testCfg, "add_allowlist_expiredIP", "10.2.2.2/32"),
				ExpectError: regexp.MustCompile("Unable to create new allowlist\nfor database. The expiration time for the allowlist is not valid. Must be a\npoint in time greater than now."),
			},

			//Add SameIP (this ip is same as the one added with required fields teststep and the config of that test step is retained)
			{
				Config:      testAccAddIPSameIP(testCfg, "add_allowlist_sameIP", "10.1.1.1/32"),
				ExpectError: regexp.MustCompile("Could not execute request, unexpected error: Unable to add allowlist entry.\nThe CIDR provided already exists for the cluster. If you continue to have\nissues connecting to your cluster please contact support."),
			},
			//Delete expired IP
			{
				Config: testAccAddExpiringIP(testCfg, "add_expiring_ip", "10.1.2.3/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_allowlist.add_expiring_ip", "cidr", "10.1.2.3/32"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_expiring_ip", "id"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_expiring_ip", "expires_at"),
					resource.TestCheckResourceAttr("capella_allowlist.add_expiring_ip", "comment", "terraform allow list acceptance test"),
					testAccWait(time.Second*140)),
			},
		},
	})
}

// Delete the ip when the ip is deleted through api
func TestAccAllowedIPDeleteIP(t *testing.T) {
	clusterName := "new_cluster"
	clusterResourceReference := "capella_cluster." + clusterName
	projectResourceName := "terraform_project"
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.0.2.0/23"

	testCfg := cfg.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateCluster(&testCfg, clusterName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(clusterResourceReference),
				),
			},
			{
				Config: testAccAddIpWithOptionalFields(testCfg, "allowList_delete", "10.2.3.4/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_allowlist.allowList_delete", "cidr", "10.2.3.4/32"),
					resource.TestCheckResourceAttrSet("capella_allowlist.allowList_delete", "id"),
					resource.TestCheckResourceAttrSet("capella_allowlist.allowList_delete", "expires_at"),
					resource.TestCheckResourceAttr("capella_allowlist.allowList_delete", "comment", "terraform allow list acceptance test"),
					testAccDeleteAllowIP(clusterResourceReference, projectResourceReference, "capella_allowlist.allowList_delete"),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
		},
	})
}

// Delete the ip when the cluster is destroyed through api
func TestAccAllowedIPDeleteCluster(t *testing.T) {
	clusterName := "new_cluster"
	clusterResourceReference := "capella_cluster." + clusterName
	projectResourceName := "terraform_project"
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.0.2.0/23"
	testCfg := cfg.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateCluster(&testCfg, clusterName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(clusterResourceReference),
				),
			},
			{
				Config: testAccAddIpWithOptionalFields(testCfg, "allowList_delete", "10.2.3.4/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_allowlist.allowList_delete", "cidr", "10.2.3.4/32"),
					resource.TestCheckResourceAttrSet("capella_allowlist.allowList_delete", "id"),
					resource.TestCheckResourceAttrSet("capella_allowlist.allowList_delete", "expires_at"),
					resource.TestCheckResourceAttr("capella_allowlist.allowList_delete", "comment", "terraform allow list acceptance test"),
					testAccDeleteCluster(clusterResourceReference, projectResourceReference),
					testAccDeleteProject(projectResourceReference),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
		},
	})
}

func testAccCreateCluster(cfg *string, resourceName, projectResourceName, projectResourceReference, cidr string) string {
	log.Println("Creating cluster")
	*cfg = fmt.Sprintf(`
%[1]s

resource "capella_project" "%[3]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "capella_cluster" "%[2]s" {
  organization_id = var.organization_id
  project_id      = %[4]s.id
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
`, *cfg, resourceName, projectResourceName, projectResourceReference, cidr)
	return *cfg
}

func testAccAddIpWithReqFields(cfg *string) string {

	*cfg = fmt.Sprintf(`
%[1]s

output "add_allowlist_req"{
  value = capella_allowlist.add_allowlist_req
}

resource "capella_allowlist" "add_allowlist_req" {
  organization_id = var.organization_id
  project_id      = "capella_project.terraform_project.id"
  cluster_id      = "capella_cluster.new_cluster.id"
  cidr            = "10.1.1.1/32"
  comment		  = "Terraform acceptance tests"
}

`, *cfg)
	return *cfg
}

func testAccAddIpWithOptionalFields(cfg string, resourceName string, cidr string) string {
	timeNow := time.Now()
	timeNow = timeNow.AddDate(0, 0, 30).UTC()
	expiryTime := timeNow.Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

output "%[2]s"{
  value = capella_allowlist.%[2]s
}

resource "capella_allowlist" "%[2]s" {
  organization_id = var.organization_id
  project_id      = "capella_project.terraform_project.id"
  cluster_id      = "capella_cluster.new_cluster.id"
  cidr            = "%[3]s"
  comment		  = "terraform allow list acceptance test"
  expires_at      = "%[4]s"
}

`, cfg, resourceName, cidr, expiryTime)
}

func testAccAddIPSameIP(cfg string, resourceName string, cidr string) string {
	timeNow := time.Now()
	timeNow = timeNow.AddDate(0, 0, 30).UTC()
	expiryTime := timeNow.Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

output "%[2]s_1"{
  value = capella_allowlist.%[2]s_1
}

resource "capella_allowlist" "%[2]s_1" {
  organization_id = var.organization_id
  project_id      = "capella_project.terraform_project.id"
  cluster_id      = "capella_cluster.new_cluster.id"
  cidr            = "%[3]s"
  comment		  = "terraform allow list acceptance test"
  expires_at      = "%[4]s"
}

output "%[2]s_2"{
  value = capella_allowlist.%[2]s_2
}

resource "capella_allowlist" "%[2]s_2" {
  organization_id = var.organization_id
  project_id      = "capella_project.terraform_project.id"
  cluster_id      = "capella_cluster.new_cluster.id"
  cidr            = "%[3]s"
  comment		  = "terraform allow list acceptance test"
  expires_at      = "%[4]s"
}

`, cfg, resourceName, cidr, expiryTime)
}

func testAccAddIpWithExpiredIP(cfg string, resourceName string, cidr string) string {
	timeNow := time.Now().UTC()
	time.Sleep(time.Second * 10)
	expiryTime := timeNow.Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

output "%[2]s"{
  value = capella_allowlist.%[2]s
}

resource "capella_allowlist" "%[2]s" {
  organization_id = var.organization_id
  project_id      = "capella_project.terraform_project.id"
  cluster_id      = "capella_cluster.new_cluster.id"
  cidr            = "%[3]s"
  comment		  = "terraform allow list acceptance test"
  expires_at      = "%[4]s"
}

`, cfg, resourceName, cidr, expiryTime)
}

func testAccAddExpiringIP(cfg string, resourceName string, cidr string) string {
	timeNow := time.Now().UTC()
	//Add two minutes for expiry so that IP can be expired
	timeNow = timeNow.Add(time.Minute * 2)
	expiryTime := timeNow.Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

output "%[2]s"{
  value = capella_allowlist.%[2]s
}

resource "capella_allowlist" "%[2]s" {
  organization_id = var.organization_id
  project_id      = "capella_project.terraform_project.id"
  cluster_id      = "capella_cluster.new_cluster.id"
  cidr            = "%[3]s"
  comment		  = "terraform allow list acceptance test"
  expires_at      = "%[4]s"
}

`, cfg, resourceName, cidr, expiryTime)
}

/*************** These functions to be moved to common util folder as they can be used with other tests as well **********************/

func testAccDeleteAllowIP(clusterResourceReference, projectResourceReference, allowIPResoureceReference string) resource.TestCheckFunc {
	log.Println("deleting the ip")
	return func(s *terraform.State) error {
		var clusterState, projectState, allowListState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[clusterResourceReference]; ok {
					clusterState = v.Primary.Attributes
				}
				if v, ok := m.Resources[projectResourceReference]; ok {
					projectState = v.Primary.Attributes
				}
				if v, ok := m.Resources[allowIPResoureceReference]; ok {
					allowListState = v.Primary.Attributes
				}
			}
		}
		data, err := TestClient()
		if err != nil {
			return err
		}
		host := os.Getenv("TF_VAR_host")
		orgid := os.Getenv("TF_VAR_organization_id")
		authToken := os.Getenv("TF_VAR_auth_token")
		_, err = data.Client.Execute(
			fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s//allowedcidrs/%s", host, orgid, projectState["id"], clusterState["id"], allowListState["id"]),
			http.MethodDelete,
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

func testAccDeleteProject(projectResourceReference string) resource.TestCheckFunc {
	log.Println("Deleting the project")
	return func(s *terraform.State) error {
		var projectState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[projectResourceReference]; ok {
					projectState = v.Primary.Attributes
				}
			}
		}
		data, err := TestClient()
		if err != nil {
			return err
		}
		host := os.Getenv("TF_VAR_host")
		orgid := os.Getenv("TF_VAR_organization_id")
		authToken := os.Getenv("TF_VAR_auth_token")
		_, err = data.Client.Execute(
			fmt.Sprintf("%s/v4/organizations/%s/projects/%s", host, orgid, projectState["id"]),
			http.MethodDelete,
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

func testAccDeleteCluster(clusterResourceReference, projectResourceReference string) resource.TestCheckFunc {
	log.Println("Deleting the cluster")
	return func(s *terraform.State) error {
		var clusterState, projectState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[clusterResourceReference]; ok {
					clusterState = v.Primary.Attributes
				}
				if v, ok := m.Resources[projectResourceReference]; ok {
					projectState = v.Primary.Attributes
				}
			}
		}
		data, err := TestClient()
		if err != nil {
			return err
		}
		host := os.Getenv("TF_VAR_host")
		orgid := os.Getenv("TF_VAR_organization_id")
		authToken := os.Getenv("TF_VAR_auth_token")
		_, err = data.Client.Execute(
			fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", host, orgid, projectState["id"], clusterState["id"]),
			http.MethodDelete,
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

func testAccWait(duration time.Duration) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		time.Sleep(duration)
		return nil
	}
}

// This function takes a resource reference string and returns a resource.TestCheckFunc. The returned function, when used
// in Terraform acceptance tests, ensures that the specified cluster resource exists in the Terraform state. It retrieves
// the resource by name from the Terraform state and checks its existence. If the resource exists, it returns nil; otherwise,
// it returns an error.
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
		data, err := TestClient()
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

// retrieveClusterFromServer checks cluster exists in server.
func retrieveClusterFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) (*clusterapi.GetClusterResponse, error) {
	response, err := data.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", data.HostURL, organizationId, projectId, clusterId),
		http.MethodGet,
		nil,
		data.Token,
		nil,
	)
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
