package resources_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"net/http"
	"terraform-provider-capella/internal/api"
	clusterapi "terraform-provider-capella/internal/api/cluster"
	providerschema "terraform-provider-capella/internal/schema"
	"testing"
	"time"

	cfg "terraform-provider-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
//
//	var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
//		"capella": providerserver.NewProtocol6WithError(provider.New("test")()),
//	}
var (
// data *providerschema.Data
// organizationId, projectId, clusterId string
)

//func testAccPreCheck(t *testing.T) {
//	// You can add code here to run prior to any test case execution, for
//	// example assertions about the appropriate environment variables being set
//	// are common to see in a pre-check function.
//	//data, err := cfg.SharedClient("", "")
//	//if err != nil {
//	//	t.Fatalf(err.Error())
//	//}
//}

//func TestAccClusterResource(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { testAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			// Create and Read testing
//			{
//				Config: testAccClusterResourceConfig(cfg.Cfg),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("capella_cluster.new_cluster", "name", "Terraform Acceptance Test Cluster"),
//					//testAccDeleteClusterResource("capella_cluster.new_cluster"),
//					//resource.TestCheckResourceAttr("capella_project.acc_test", "description", "description"),
//					//resource.TestCheckResourceAttr("capella_project.acc_test", "etag", "Version: 1"),
//				),
//			},
//			////// ImportState testing
//			//{
//			//	ResourceName:      "capella_project.acc_test",
//			//	ImportStateIdFunc: generateProjectImportId,
//			//	ImportState:       true,
//			//	ImportStateVerify: true,
//			//},
//			//// Update and Read testing
//			//{
//			//	Config: testAccProjectResourceConfigUpdate(cfg.Cfg),
//			//	Check: resource.ComposeAggregateTestCheckFunc(
//			//		resource.TestCheckResourceAttr("capella_project.acc_test", "name", "acc_test_project_name_update"),
//			//		resource.TestCheckResourceAttr("capella_project.acc_test", "description", "description_update"),
//			//	),
//			//},
//			//{
//			//	Config: testAccProjectResourceConfigUpdateWithIfMatch(cfg.Cfg),
//			//	Check: resource.ComposeAggregateTestCheckFunc(
//			//		resource.TestCheckResourceAttr("capella_project.acc_test", "name", "acc_test_project_name_update_with_if_match"),
//			//		resource.TestCheckResourceAttr("capella_project.acc_test", "description", "description_update_with_match"),
//			//		resource.TestCheckResourceAttr("capella_project.acc_test", "etag", "Version: 3"),
//			//		resource.TestCheckResourceAttr("capella_project.acc_test", "if_match", "2"),
//			//	),
//			//},
//			// Delete testing automatically occurs in TestCase
//		},
//	})
//}

func TestAccClusterResourceNotFound(t *testing.T) {
	var (
	//organizationId, projectId, clusterId string
	//data                                 *providerschema.Data
	//err                                  error
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: cfg.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					//data, err = cfg.SharedClient("", "")
					//if err != nil {
					//	t.Fatalf(err.Error())
					//}
				},
				Config: testAccClusterResourceConfig(cfg.Cfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("capella_cluster.new_cluster", "name", "Terraform Acceptance Test Cluster"),
					//preConfigValue("capella_cluster.new_cluster"),
					testAccDeleteClusterResource("capella_cluster.new_cluster"),
					//resource.TestCheckResourceAttr("capella_project.acc_test", "description", "description"),
					//resource.TestCheckResourceAttr("capella_project.acc_test", "etag", "Version: 1"),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
			////// ImportState testing
			//{
			//	ResourceName:      "capella_project.acc_test",
			//	ImportStateIdFunc: generateProjectImportId,
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
			//// Update and Read testing
			{
				//PreConfig: func() {
				//	fmt.Printf("temporary dir")
				//	fmt.Printf(t.TempDir())
				//	fmt.Printf("organization Id:")
				//	fmt.Printf(organizationId)
				//	fmt.Printf("project Id: ")
				//	fmt.Printf(projectId)
				//	fmt.Printf("cluster Id: ")
				//	fmt.Printf(clusterId)
				//	if organizationId == "" {
				//		t.Fatalf(errors.ErrOrganizationIdCannotBeEmpty.Error())
				//	}
				//	if projectId == "" {
				//		t.Fatalf(errors.ErrProjectIdCannotBeEmpty.Error())
				//	}
				//	if clusterId == "" {
				//		t.Fatalf(errors.ErrClusterIdCannotBeEmpty.Error())
				//	}
				//	fmt.Printf("organization Id: " + organizationId)
				//	fmt.Printf("project Id: " + projectId)
				//	fmt.Printf("cluster Id: " + clusterId)
				//	deleteClusterFromServer(data, organizationId, projectId, clusterId)
				//	fmt.Printf("started Deletion")
				//},
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
				Config:             testAccClusterResourceConfigUpdate(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_cluster.new_cluster", "name", "Terraform Acceptance Test Cluster Update"),
				),
			},
			//{
			//	Config: testAccProjectResourceConfigUpdateWithIfMatch(cfg.Cfg),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("capella_project.acc_test", "name", "acc_test_project_name_update_with_if_match"),
			//		resource.TestCheckResourceAttr("capella_project.acc_test", "description", "description_update_with_match"),
			//		resource.TestCheckResourceAttr("capella_project.acc_test", "etag", "Version: 3"),
			//		resource.TestCheckResourceAttr("capella_project.acc_test", "if_match", "2"),
			//	),
			//},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccClusterResourceConfig(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "new_project" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "capella_cluster" "new_cluster" {
  organization_id = var.organization_id
  project_id      = capella_project.new_project.id
  name            = "Terraform Acceptance Test Cluster"
  description     = "My first test cluster for multiple services."
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "10.250.250.0/23"
  }
  configuration_type = "multiNode"
  couchbase_server = {
    version = "7.1"
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
          type    = "io2"
          iops    = 5000
        }
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
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
`, cfg)
}

func testAccClusterResourceConfigUpdate(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "new_project" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "capella_cluster" "new_cluster" {
  organization_id = var.organization_id
  project_id      = capella_project.new_project.id
  name            = "Terraform Acceptance Test Cluster Update"
  description     = "My first test cluster for multiple services."
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "10.250.250.0/23"
  }
  configuration_type = "multiNode"
  couchbase_server = {
    version = "7.1"
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
          type    = "io2"
          iops    = 5000
        }
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
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
`, cfg)
}

//func testAccProjectResourceConfigUpdate(cfg string) string {
//	return fmt.Sprintf(`
//%[1]s
//
//resource "capella_project" "acc_test" {
//   organization_id = var.organization_id
//	name            = "acc_test_project_name_update"
//	description     = "description_update"
//}
//`, cfg)
//}
//
//func testAccProjectResourceConfigUpdateWithIfMatch(cfg string) string {
//	return fmt.Sprintf(`
//%[1]s
//
//resource "capella_project" "acc_test" {
//    organization_id = var.organization_id
//	name            = "acc_test_project_name_update_with_if_match"
//	description     = "description_update_with_match"
//	if_match        =  2
//}
//`, cfg)
//}
//
//func generateProjectImportId(state *terraform.State) (string, error) {
//	resourceName := "capella_project.acc_test"
//	var rawState map[string]string
//	for _, m := range state.Modules {
//		if len(m.Resources) > 0 {
//			if v, ok := m.Resources[resourceName]; ok {
//				rawState = v.Primary.Attributes
//			}
//		}
//	}
//	fmt.Printf("raw state %s", rawState)
//	return fmt.Sprintf("id=%s,organization_id=%s", rawState["id"], rawState["organization_id"]), nil
//}

// testAccCheckExampleWidgetExists uses the Example SDK directly to retrieve
// the Widget description, and stores it in the provided
// *example.WidgetDescription

func testAccDeleteClusterResource(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state

		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceName]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}

		//rs, ok := s.RootModule().Resources[resourceName]
		//if !ok {
		//	return fmt.Errorf("Not found: %s", resourceName)
		//}
		//
		//if rs.Primary.ID == "" {
		//	return fmt.Errorf("Widget ID is not set")
		//}
		//
		//if err != nil {
		//	return nil, err
		//}

		//response, err := sharedClient..DescribeWidgets(&example.DescribeWidgetsInput{
		//	WidgetIDs: []string{rs.Primary.ID},
		//})

		//if err != nil {
		//	return err
		//}

		// we expect only a single widget by this ID. If we find zero, or many,
		// then we consider this an error
		//if len(response.WidgetDescriptions) != 1 ||
		//	*response.WidgetDescriptions[0].WidgetID != rs.Primary.ID {
		//	return fmt.Errorf("Widget not found")
		//}
		//
		//// store the resulting widget in the *example.WidgetDescription pointer
		//*widget = *response.WidgetDescriptions[0]

		data, err := cfg.SharedClient("", "")
		if err != nil {
			return err
		}
		err = deleteClusterFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["id"])
		if err != nil {
			return err
		}
		fmt.Printf("delete initiated")
		err = checkClusterStatus(data, context.Background(), rawState["organization_id"], rawState["project_id"], rawState["id"])
		resourceNotFound, err := handleClusterError(err)
		if !resourceNotFound {
			return err
		}
		fmt.Printf("successfully deleted")
		return nil
	}
}

//func preConfigValue(resourceName string) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		// retrieve the resource by name from state
//
//		var rawState map[string]string
//		for _, m := range s.Modules {
//			if len(m.Resources) > 0 {
//				if v, ok := m.Resources[resourceName]; ok {
//					rawState = v.Primary.Attributes
//				}
//			}
//		}
//		stateOrganizationId := rawState["organization_id"]
//		fmt.Printf("******************************************** organizationId***************************")
//		fmt.Printf(stateOrganizationId)
//		organizationId = stateOrganizationId
//
//		stateProjectId := rawState["project_id"]
//		projectId = stateProjectId
//		fmt.Printf("******************************************** project_id ***************************")
//		fmt.Printf(stateProjectId)
//
//		stateClusterId := rawState["id"]
//		clusterId = stateClusterId
//		fmt.Printf("******************************************** cluster_id ***************************")
//		fmt.Printf(stateClusterId)
//
//		return nil
//	}
//}

func retriveClusterFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) (*clusterapi.GetClusterResponse, error) {
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

func deleteClusterFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) error {
	fmt.Println(fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", data.HostURL, organizationId, projectId, clusterId))
	_, err := data.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", data.HostURL, organizationId, projectId, clusterId),
		http.MethodDelete,
		nil,
		data.Token,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func checkClusterStatus(data *providerschema.Data, ctx context.Context, organizationId, projectId, ClusterId string) error {
	var (
		clusterResp *clusterapi.GetClusterResponse
		err         error
	)

	// Assuming 60 minutes is the max time deployment takes, can change after discussion
	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	const sleep = time.Second * 3

	timer := time.NewTimer(2 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			const msg = "cluster creation status transition timed out after initiation"
			return fmt.Errorf(msg)

		case <-timer.C:
			clusterResp, err = retriveClusterFromServer(data, organizationId, projectId, ClusterId)
			switch err {
			case nil:
				if clusterapi.IsFinalState(clusterResp.CurrentState) {
					return nil
				}
				const msg = "waiting for cluster to complete the execution"
				tflog.Info(ctx, msg)
			default:
				return err
			}
			timer.Reset(sleep)
		}
	}
}

func handleClusterError(err error) (bool, error) {
	switch err := err.(type) {
	case nil:
		return false, nil
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			return false, fmt.Errorf(err.CompleteError())
		}
		return true, fmt.Errorf(err.CompleteError())
	default:
		return false, err
	}
}
