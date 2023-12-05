package acceptance_tests

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"github.com/couchbasecloud/couchbase-capella/internal/api"
	acctest "github.com/couchbasecloud/couchbase-capella/internal/testing"

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
	cidr := "10.250.250.0/23"

	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
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
				ExpectError: regexp.MustCompile("Unable\nto create new allowlist for database. The expiration time for the allowlist\nis not valid. Must be a point in time greater than now."),
			},

			//Add SameIP (this ip is same as the one added with required fields teststep and the config of that test step is retained)
			{
				Config:      testAccAddIPSameIP(testCfg, "add_allowlist_sameIP", "10.1.1.1/32"),
				ExpectError: regexp.MustCompile("CIDR provided already exists for the cluster"),
			},
			//Delete expired IP
			{
				Config: testAccAddExpiringIP(testCfg, "add_expiring_ip", "10.1.2.3/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_allowlist.add_expiring_ip", "cidr", "10.1.2.3/32"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_expiring_ip", "id"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_expiring_ip", "expires_at"),
					resource.TestCheckResourceAttr("capella_allowlist.add_expiring_ip", "comment", "terraform allow list acceptance test"),
					acctest.TestAccWait(time.Second*250)),
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
	cidr := "10.4.2.0/23"
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
				Config: testAccAddIpWithOptionalFields(testCfg, "allowList_delete", "10.4.3.4/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_allowlist.allowList_delete", "cidr", "10.4.3.4/32"),
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

func testAccAddIpWithReqFields(cfg *string) string {

	*cfg = fmt.Sprintf(`
%[1]s

output "add_allowlist_req"{
  value = capella_allowlist.add_allowlist_req
}

resource "capella_allowlist" "add_allowlist_req" {
  organization_id = var.organization_id
  project_id      = capella_project.terraform_project.id
  cluster_id      = capella_cluster.new_cluster.id
  cidr            = "10.1.1.1/32"
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
  project_id      = capella_project.terraform_project.id
  cluster_id      = capella_cluster.new_cluster.id
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
  project_id      = capella_project.terraform_project.id
  cluster_id      = capella_cluster.new_cluster.id
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
  project_id      = capella_project.terraform_project.id
  cluster_id      = capella_cluster.new_cluster.id
  cidr            = "%[3]s"
  comment		  = "terraform allow list acceptance test"
  expires_at      = "%[4]s"
}

`, cfg, resourceName, cidr, expiryTime)
}

func testAccAddExpiringIP(cfg string, resourceName string, cidr string) string {
	timeNow := time.Now().UTC()
	//Add two minutes for expiry so that IP can be expired
	timeNow = timeNow.Add(time.Minute * 4)
	expiryTime := timeNow.Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

output "%[2]s"{
  value = capella_allowlist.%[2]s
}

resource "capella_allowlist" "%[2]s" {
  organization_id = var.organization_id
  project_id      = capella_project.terraform_project.id
  cluster_id      = capella_cluster.new_cluster.id
  cidr            = "%[3]s"
  comment		  = "terraform allow list acceptance test"
  expires_at      = "%[4]s"
}

`, cfg, resourceName, cidr, expiryTime)
}

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
		data, err := acctest.TestClient()
		if err != nil {
			return err
		}
		host := os.Getenv("TF_VAR_host")
		orgid := os.Getenv("TF_VAR_organization_id")
		authToken := os.Getenv("TF_VAR_auth_token")
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/allowedcidrs/%s", host, orgid, projectState["id"], clusterState["id"], allowListState["id"])
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
