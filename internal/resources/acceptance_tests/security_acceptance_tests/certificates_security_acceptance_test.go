package security_acceptance_tests

import (
	"fmt"
	"os"

	"regexp"
	"testing"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
	cfg "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCertificateDataSourceNoAuth(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	os.Setenv("TF_VAR_auth_token", "")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCertificateResourceConfig(cfg.Cfg),
				ExpectError: regexp.MustCompile("Missing Capella Authentication Token"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCertificateDataSourceRbacOrgOwner(t *testing.T) {
	organizationId := os.Getenv("TF_VAR_organization_id")
	projectId := os.Getenv("TF_VAR_project_id")
	clusterId := os.Getenv("TF_VAR_cluster_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationOwner")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateResourceConfig(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.capella_certificate.existing_certificate", "certificate"),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "organization_id", organizationId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "project_id", projectId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "cluster_id", clusterId),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCertificateDataSourceRbacOrgMember(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationMember")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCertificateResourceConfig(cfg.Cfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCertificateDataSourceRbacProjCreator(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("projectCreator")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCertificateResourceConfig(cfg.Cfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCertificateDataSourceRbacProjOwner(t *testing.T) {
	organizationId := os.Getenv("TF_VAR_organization_id")
	projectId := os.Getenv("TF_VAR_project_id")
	clusterId := os.Getenv("TF_VAR_cluster_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectOwner")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateResourceConfig(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.capella_certificate.existing_certificate", "certificate"),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "organization_id", organizationId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "project_id", projectId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "cluster_id", clusterId),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCertificateDataSourceRbacProjManager(t *testing.T) {
	organizationId := os.Getenv("TF_VAR_organization_id")
	projectId := os.Getenv("TF_VAR_project_id")
	clusterId := os.Getenv("TF_VAR_cluster_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectManager")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateResourceConfig(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.capella_certificate.existing_certificate", "certificate"),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "organization_id", organizationId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "project_id", projectId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "cluster_id", clusterId),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCertificateDataSourceRbacProjViewer(t *testing.T) {
	organizationId := os.Getenv("TF_VAR_organization_id")
	projectId := os.Getenv("TF_VAR_project_id")
	clusterId := os.Getenv("TF_VAR_cluster_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectViewer")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateResourceConfig(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.capella_certificate.existing_certificate", "certificate"),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "organization_id", organizationId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "project_id", projectId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "cluster_id", clusterId),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCertificateDataSourceRbacProjDataReaderWriter(t *testing.T) {
	organizationId := os.Getenv("TF_VAR_organization_id")
	projectId := os.Getenv("TF_VAR_project_id")
	clusterId := os.Getenv("TF_VAR_cluster_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectDataReaderWriter")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateResourceConfig(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.capella_certificate.existing_certificate", "certificate"),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "organization_id", organizationId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "project_id", projectId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "cluster_id", clusterId),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCertificateDataSourceRbacProjDataReader(t *testing.T) {
	organizationId := os.Getenv("TF_VAR_organization_id")
	projectId := os.Getenv("TF_VAR_project_id")
	clusterId := os.Getenv("TF_VAR_cluster_id")
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectDataReader")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateResourceConfig(cfg.Cfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.capella_certificate.existing_certificate", "certificate"),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "organization_id", organizationId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "project_id", projectId),
					resource.TestCheckResourceAttr("data.capella_certificate.existing_certificate", "cluster_id", clusterId),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccCertificateResourceConfig(cfg string) string {
	return fmt.Sprintf(`
%[1]s

output "certificate" {
	value = data.capella_certificate.existing_certificate
}

data "capella_certificate" "existing_certificate" {
	organization_id = var.organization_id
	project_id      = var.project_id
	cluster_id      = var.cluster_id
}

`, cfg)
}
