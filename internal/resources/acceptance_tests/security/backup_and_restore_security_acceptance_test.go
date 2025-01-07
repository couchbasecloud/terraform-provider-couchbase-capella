package security

import (
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources/acceptance_tests"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCreateBackupRestoreNoAuth(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	os.Setenv("TF_VAR_auth_token", "")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupRestoreResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Missing Capella Authentication Token"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupRestoreOrgOwner(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	organizationId := os.Getenv("TF_VAR_organization_id")
	projectId := os.Getenv("TF_VAR_project_id")
	clusterId := os.Getenv("TF_VAR_cluster_id")
	testAccCreateOrgAPI("organizationOwner")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBackupRestoreResourceConfigRequired(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_backup.new_backup", "organization_id", organizationId),
					resource.TestCheckResourceAttr("capella_backup.new_backup", "project_id", projectId),
					resource.TestCheckResourceAttr("capella_backup.new_backup", "cluster_id", clusterId),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupRestoreOrgMember(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationMember")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupRestoreResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupRestoreProjCreator(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("projectCreator")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupRestoreResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupRestoreProjOwner(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	organizationId := os.Getenv("TF_VAR_organization_id")
	projectId := os.Getenv("TF_VAR_project_id")
	clusterId := os.Getenv("TF_VAR_cluster_id")
	testAccCreateProjAPI("projectCreator", projectId, "projectOwner")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBackupRestoreResourceConfigRequired(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_backup.new_backup", "organization_id", organizationId),
					resource.TestCheckResourceAttr("capella_backup.new_backup", "project_id", projectId),
					resource.TestCheckResourceAttr("capella_backup.new_backup", "cluster_id", clusterId),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupRestoreProjManager(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectManager")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupRestoreResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Could not get the latest bucket backup"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupRestoreProjViewer(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectViewer")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupRestoreResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupRestoreDatabaseReaderWriter(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectDataReaderWriter")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupRestoreResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupRestoreDatabaseReader(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectDataReader")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupRestoreResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccBackupRestoreResourceConfigRequired() string {
	return fmt.Sprintf(`
%[1]s

output "new_bucket" {
	value = capella_bucket.new_bucket
}

output "bucket_id" {
	value = capella_bucket.new_bucket.id
}

resource "capella_bucket" "new_bucket" {
	name                       = "terraform-security-bucket"
	organization_id            = var.organization_id
	project_id                 = var.project_id
	cluster_id                 = var.cluster_id
	type                       = "couchbase"
	storage_backend            = "couchstore"
	memory_allocation_in_mb    = 105
	bucket_conflict_resolution = "seqno"
	durability_level           = "majorityAndPersistActive"
	replicas                   = 2
	flush                      = true
	time_to_live_in_seconds    = 100
}

output "new_backup" {
	value = capella_backup.new_backup
}

resource "capella_backup" "new_backup" {
	organization_id            = var.organization_id
	project_id                 = var.project_id
	cluster_id                 = var.cluster_id
	bucket_id                  = capella_bucket.new_bucket.id
}

`, acceptance_tests.ProviderBlock)

}
