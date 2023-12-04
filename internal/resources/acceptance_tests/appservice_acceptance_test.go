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

func TestAppServiceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAppServiceResourceConfig(acctest.Cfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "description", "description"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "compute.cpu", "2"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "compute.ram", "4"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "nodes", "2"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      "capella_app_service.new_app_service",
				ImportStateIdFunc: generateAppServiceImportId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAppServiceResourceConfigUpdate(acctest.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "description", "description"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "compute.cpu", "2"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "compute.ram", "4"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "nodes", "3"),
				),
			},
			{
				Config: testAccAppServiceResourceConfigUpdateWithIfMatch(acctest.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "description", "description"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "compute.cpu", "4"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "compute.ram", "8"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "nodes", "2"),
					resource.TestCheckResourceAttr("capella_app_service.new_app_service", "if_match", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAppServiceCreateWithReqFields(t *testing.T) {
	appServiceResourceName := "app_service_req_fields"
	appServiceResourceReference := "capella_app_service." + appServiceResourceName
	clusterResourceName := "new_cluster"
	clusterResourceReference := "capella_cluster." + clusterResourceName
	testCfg := acctest.ProjectCfg
	projectResourceName := "terraform_project"
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.1.68.0/23"
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
				Config: testAccAppServiceResourceReqConfig(acctest.ProjectCfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(appServiceResourceReference, "name", "test-terraform-app-service"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "description", ""),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.cpu", "2"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.ram", "4"),
					resource.TestCheckResourceAttr(appServiceResourceReference", "nodes", "2"),
				),
			},
		},
	},
	)
}
func TestAccAppServiceCreateWithOptFields(t *testing.T) {
	resourceName := "app_service_opt_fields"
	//cidr, _ := acctest.GetCIDR()
	//fmt.Println(cidr)
	appServiceResourceName := "app_service_opt_fields"
	appServiceResourceReference := "capella_app_service." + appServiceResourceName
	clusterResourceName := "new_cluster"
	clusterResourceReference := "capella_cluster." + clusterResourceName
	testCfg := acctest.ProjectCfg
	projectResourceName := "terraform_project"
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.1.68.0/23"
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
				Config: testAccAppServiceResourceOptConfig(acctest.ProjectCfg, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(appServiceResourceReference, "name", "app_service_opt_fields"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "description", "acceptance test app service"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.cpu", "2"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.ram", "4"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "nodes", "2"),
					//resource.TestCheckResourceAttr("capella_app_service.app_service_opt_fields", "version", "3.0"),
				),
			},

			//Invalid Update of fields
			{
				Config:      testAccAppServiceResourceUpdateInvalidClusterIdConfig(acctest.ProjectCfg, resourceName),
				ExpectError: regexp.MustCompile("wrong cluster id"),
			},
			{
				Config:      testAccAppServiceResourceUpdateInvalidProjectIdConfig(acctest.ProjectCfg, resourceName),
				ExpectError: regexp.MustCompile("wrong project id"),
			},
			{
				Config:      testAccAppServiceResourceUpdateInvalidOrgIdConfig(acctest.ProjectCfg, resourceName),
				ExpectError: regexp.MustCompile("wrong org id"),
			},
		},
	},
	)
}

func TestAccAppServiceDeleteAppService(t *testing.T) {
	appServiceResourceName := "app_service_opt_fields"
	appServiceResourceReference := "capella_app_service." + appServiceResourceName
	clusterResourceName := "new_cluster"
	clusterResourceReference := "capella_cluster." + clusterResourceName
	testCfg := acctest.ProjectCfg
	projectResourceName := "terraform_project"
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.1.68.0/23"
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
				Check: resource.ComposeTestCheckFunc(
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

func TestAccAppServiceDeleteCluster(t *testing.T) {
	appServiceResourceName := "app_service_opt_fields"
	appServiceResourceReference := "capella_app_service." + appServiceResourceName
	clusterResourceName := "new_cluster"
	clusterResourceReference := "capella_cluster." + clusterResourceName
	testCfg := acctest.ProjectCfg
	projectResourceName := "terraform_project"
	projectResourceReference := "capella_project." + projectResourceName
	cidr := "10.0.30.0/23"
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(appServiceResourceReference, "name", "app_service_opt_fields"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "description", "acceptance test app service"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.cpu", "2"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "compute.ram", "4"),
					resource.TestCheckResourceAttr(appServiceResourceReference, "nodes", "2"),
					testAccDeleteCluster(projectResourceReference, clusterResourceReference),
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
resource "capella_app_service" "%[2]s" {
  organization_id = var.organization_id
  project_id      = "90bafc4e-43fe-4577-9c6f-2893478bd392"
  cluster_id      = "c517165b-bd66-4f34-9bf5-31d89bae5e8c"
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
resource "capella_app_service" "app_service_req_fields" {
  organization_id = var.organization_id
  project_id      = "90bafc4e-43fe-4577-9c6f-2893478bd392"
  cluster_id      = "c517165b-bd66-4f34-9bf5-31d89bae5e8c"
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

resource "capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
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

resource "capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
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

resource "capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
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
	resourceName := "capella_app_service.new_app_service"
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
resource "capella_app_service" "%[2]s" {
  organization_id = var.organization_id
  project_id      = "90bafc4e-43fe-4577-9c6f-2893478bd392"
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
resource "capella_app_service" "%[2]s" {
  organization_id = var.organization_id
  project_id      = "55556666-4444-3333-2222-11111ffffff"
  cluster_id      = "c517165b-bd66-4f34-9bf5-31d89bae5e8c"
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
resource "capella_app_service" "%[2]s" {
  organization_id = "55556666-4444-3333-2222-11111ffffff"
  project_id      = "90bafc4e-43fe-4577-9c6f-2893478bd392"
  cluster_id      = "c517165b-bd66-4f34-9bf5-31d89bae5e8c"
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
		return nil
	}

}
