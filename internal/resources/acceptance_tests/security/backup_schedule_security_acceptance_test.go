package security

import (
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources/acceptance_tests"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCreateBackupScheduleNoAuth(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	os.Setenv("TF_VAR_auth_token", "")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupScheduleResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Missing Capella Authentication Token"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupScheduleOrgOwner(t *testing.T) {
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
				Config: testAccBackupScheduleResourceConfigRequired(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "organization_id", organizationId),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "project_id", projectId),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "cluster_id", clusterId),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "weekly_schedule.day_of_week", "sunday"),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "weekly_schedule.start_at", "10"),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "weekly_schedule.incremental_every", "4"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupScheduleOrgMember(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationMember")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupScheduleResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupScheduleProjCreator(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("projectCreator")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupScheduleResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupScheduleProjOwner(t *testing.T) {
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
				Config: testAccBackupScheduleResourceConfigRequired(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "organization_id", organizationId),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "project_id", projectId),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "cluster_id", clusterId),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "weekly_schedule.day_of_week", "sunday"),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "weekly_schedule.start_at", "10"),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "weekly_schedule.incremental_every", "4"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupScheduleProjManager(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	organizationId := os.Getenv("TF_VAR_organization_id")
	projectId := os.Getenv("TF_VAR_project_id")
	clusterId := os.Getenv("TF_VAR_cluster_id")
	testAccCreateProjAPI("projectCreator", projectId, "projectManager")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBackupScheduleResourceConfigRequired(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "organization_id", organizationId),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "project_id", projectId),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "cluster_id", clusterId),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "weekly_schedule.day_of_week", "sunday"),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "weekly_schedule.start_at", "10"),
					resource.TestCheckResourceAttr("capella_backup_schedule.new_backup_schedule", "weekly_schedule.incremental_every", "4"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupScheduleProjViewer(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectViewer")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupScheduleResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupScheduleDatabaseReaderWriter(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectDataReaderWriter")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupScheduleResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBackupScheduleDatabaseReader(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("projectCreator", projId, "projectDataReader")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBackupScheduleResourceConfigRequired(),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccBackupScheduleResourceConfigRequired() string {
	return fmt.Sprintf(`
%[1]s

output "new_bucket" {
	value = capella_bucket.new_bucket
}

output "bucket_id" {
	value = capella_bucket.new_bucket.id
}

resource "capella_bucket" "new_bucket" {
	name                       = "terraform-bucket"
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

output "new_backup_schedule" {
	value = capella_backup_schedule.new_backup_schedule
  }

  resource "capella_backup_schedule" "new_backup_schedule" {
	organization_id            = var.organization_id
	project_id                 = var.project_id
	cluster_id                 = var.cluster_id
	bucket_id                  = capella_bucket.new_bucket.id
	type = "weekly"
	weekly_schedule = {
	  day_of_week = "sunday"
	  start_at = 10
	  incremental_every = 4
	  retention_time = "90days"
	  cost_optimized_retention = false
	  }
}

`, acceptance_tests.ProviderBlock)
}
