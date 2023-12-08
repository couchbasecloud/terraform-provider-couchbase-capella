package acceptance_tests

import (
	"context"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	appserviceapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"log"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAppServiceResource(t *testing.T) {

	clusterResourceName := "new_cluster"
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	testCfg := acctest.ProjectCfg
	projectResourceName := "terraform_project"
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr, err := acctest.GetCIDR("aws")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateCluster(&testCfg, clusterResourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(clusterResourceReference),
				),
			},
			// Create and Read testing
			{
				Config: testAccAppServiceResourceConfig(testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "description", "description"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.cpu", "2"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.ram", "4"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      "couchbase-capella_app_service.new_app_service",
				ImportStateIdFunc: generateAppServiceImportId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAppServiceResourceConfigUpdate(testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "description", "description"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.cpu", "2"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.ram", "4"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "nodes", "3"),
				),
			},
			{
				Config: testAccAppServiceResourceConfigUpdateWithIfMatch(testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "description", "description"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.cpu", "4"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "compute.ram", "8"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "nodes", "2"),
					resource.TestCheckResourceAttr("couchbase-capella_app_service.new_app_service", "if_match", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccAppServiceCreateWithReqFields(t *testing.T) {
	appServiceResourceName := "app_service_req_fields"
	appServiceResourceReference := "couchbase-capella_app_service." + appServiceResourceName
	clusterResourceName := "new_cluster"
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	testCfg := acctest.ProjectCfg
	projectResourceName := "terraform_project"
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr, err := acctest.GetCIDR("aws")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateCluster(&testCfg, clusterResourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(clusterResourceReference),
				),
			},
			{
				Config: testAccAppServiceResourceReqConfig(testCfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(appServiceResourceReference, "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "description", ""),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.cpu", "2"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.ram", "4"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "nodes", "2"),
				),
			},
		},
	},
	)
}
func TestAccAppServiceCreateWithOptFields(t *testing.T) {
	resourceName := "app_service_opt_fields"
	cidr, err := acctest.GetCIDR("aws")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	appServiceResourceName := "app_service_opt_fields"
	appServiceResourceReference := "couchbase-capella_app_service." + appServiceResourceName
	clusterResourceName := "new_cluster"
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	testCfg := acctest.ProjectCfg
	projectResourceName := "terraform_project"
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateCluster(&testCfg, clusterResourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(clusterResourceReference),
				),
			},
			{
				Config: testAccAppServiceResourceOptConfig(testCfg, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(appServiceResourceReference, "name", "app_service_opt_fields"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "description", "acceptance test app service"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.cpu", "2"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.ram", "4"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "nodes", "2"),
				),
			},

			//Invalid Update of fields
			//Expected error: Could not execute request, unexpected error: {"code":1000,"hint":"Check if all the required params are present in the request body.","httpStatusCode":400,"message":"The server cannot or will not process the
			//│ request due to something that is perceived to be a client error."}
			//we expect above error but due to formatting issues of expected error we are only checking unique word "perceived"
			{
				Config:      testAccAppServiceResourceUpdateInvalidProjectIdConfig(testCfg, resourceName),
				ExpectError: regexp.MustCompile("perceived"),
			},
			{
				Config:      testAccAppServiceResourceUpdateInvalidOrgIdConfig(testCfg, resourceName),
				ExpectError: regexp.MustCompile("perceived"),
			},
			{
				Config:      testAccAppServiceResourceUpdateInvalidClusterIdConfig(testCfg, resourceName),
				ExpectError: regexp.MustCompile("perceived"),
			},
		},
	},
	)
}

func TestAccAppServiceDeleteAppService(t *testing.T) {
	appServiceResourceName := "app_service_opt_fields"
	appServiceResourceReference := "couchbase-capella_app_service." + appServiceResourceName
	clusterResourceName := "new_cluster"
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	testCfg := acctest.ProjectCfg
	projectResourceName := "terraform_project"
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr, err := acctest.GetCIDR("aws")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateCluster(&testCfg, clusterResourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(clusterResourceReference),
				),
			},
			{
				Config: testAccAppServiceResourceOptConfig(testCfg, appServiceResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(appServiceResourceReference, "name", "app_service_opt_fields"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "description", "acceptance test app service"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.cpu", "2"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.ram", "4"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "nodes", "2"),
					testAccDeleteAppService(projectResourceReference, clusterResourceReference, appServiceResourceReference),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
		},
	})
}

func testAccAppServiceResourceOptConfig(cfg, resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_app_service" "%[2]s" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  description	  = "acceptance test app service"
  name            = "app_service_opt_fields"
  nodes			  = "2"
  compute = {
    cpu = 2
    ram = 4
}
}
`, cfg, resourceName)
}

func testAccAppServiceResourceReqConfig(cfg string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_app_service" "app_service_req_fields" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name            = "test-terraform-app-service"
  compute = {
    cpu = 2
    ram = 4
}
}
`, cfg)
}

func testAccAppServiceResourceConfig(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name            = "test-terraform-app-service"
  description     = "description"
  compute = {
    cpu = 2
    ram = 4
}
}
`, cfg)
}

func testAccAppServiceResourceConfigUpdate(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name            = "test-terraform-app-service"
  description     = "description"
  compute = {
    cpu = 2
    ram = 4
  }
  nodes = 3
}
`, cfg)
}

func testAccAppServiceResourceConfigUpdateWithIfMatch(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name            = "test-terraform-app-service"
  description     = "description"
  if_match        =  2
  compute = {
    cpu = 4
    ram = 8
  }
  nodes = 2
}
`, cfg)
}

func generateAppServiceImportId(state *terraform.State) (string, error) {
	resourceName := "couchbase-capella_app_service.new_app_service"
	var rawState map[string]string
	for _, m := range state.Modules {
		if len(m.Resources) > 0 {
			if v, ok := m.Resources[resourceName]; ok {
				rawState = v.Primary.Attributes
			}
		}
	}
	fmt.Printf("raw state %s", rawState)
	return fmt.Sprintf("id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
}

func testAccAppServiceResourceUpdateInvalidClusterIdConfig(cfg, resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_app_service" "%[2]s" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = "55556666-4444-3333-2222-11111ffffff"
  description	  = "acceptance test app service"
  name            = "app_service_opt_fields"
  nodes			  = "2"
  compute = {
    cpu = 2
    ram = 4
}
}
`, cfg, resourceName)
}

func testAccAppServiceResourceUpdateInvalidProjectIdConfig(cfg, resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_app_service" "%[2]s" {
  organization_id = var.organization_id
  project_id      = "55556666-4444-3333-2222-11111ffffff"
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  description	  = "acceptance test app service"
  name            = "app_service_opt_fields"
  nodes			  = "2"
  compute = {
    cpu = 2
    ram = 4
}
}
`, cfg, resourceName)
}

func testAccAppServiceResourceUpdateInvalidOrgIdConfig(cfg, resourceName string) string {
	return fmt.Sprintf(`
%[1]s
resource "couchbase-capella_app_service" "%[2]s" {
  organization_id = "55556666-4444-3333-2222-11111ffffff"
  project_id      = couchbase-capella_project.terraform_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  description	  = "acceptance test app service"
  name            = "app_service_opt_fields"
  nodes			  = "2"
  compute = {
    cpu = 2
    ram = 4
}
}
`, cfg, resourceName)
}

func testAccDeleteAppService(projectResourceReference, clusterResourceReference, appServiceResourceReference string) resource.TestCheckFunc {
	log.Println("deleting the appService")
	return func(s *terraform.State) error {
		var clusterState, projectState, appServiceState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[clusterResourceReference]; ok {
					clusterState = v.Primary.Attributes
				}
				if v, ok := m.Resources[projectResourceReference]; ok {
					projectState = v.Primary.Attributes
				}
				if v, ok := m.Resources[appServiceResourceReference]; ok {
					appServiceState = v.Primary.Attributes
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
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s", host, orgid, projectState["id"], clusterState["id"], appServiceState["id"])
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
		_, err = data.Client.Execute(
			cfg,
			nil,
			authToken,
			nil,
		)
		if err != nil {
			return err
		}
		err = checkAppServiceStatus(data, context.Background(), orgid, projectState["id"], clusterState["id"], appServiceState["id"])
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if !resourceNotFound {
			return fmt.Errorf(errString)
		}
		fmt.Printf("successfully deleted")
		return nil
	}

}

func checkAppServiceStatus(data *providerschema.Data, ctx context.Context, orgId, projectId, clusterId, appServiceId string) error {

	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	const sleep = time.Second * 3
	timer := time.NewTimer(2 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("wrong cluster state")

		case <-timer.C:
			res, err := retrieveAppService(data, orgId, projectId, clusterId, appServiceId)
			switch err {
			case nil:
				if appserviceapi.IsFinalState(res.CurrentState) {
					return nil
				}
				const msg = "waiting for app service to get deleted"
				tflog.Info(ctx, msg)
			default:
				return err
			}
			timer.Reset(sleep)

		}
	}
}

func retrieveAppService(data *providerschema.Data, orgId, projectId, clusterId, appServiceId string) (*appserviceapi.GetAppServiceResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s", data.HostURL, orgId, projectId, clusterId, appServiceId)
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
	appserviceresponse := appserviceapi.GetAppServiceResponse{}
	if err != nil {
		return nil, err
	}
	appserviceresponse.Etag = response.Response.Header.Get("ETag")
	return &appserviceresponse, nil
}
