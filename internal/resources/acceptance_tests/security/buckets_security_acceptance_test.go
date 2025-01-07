package security

import (
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources/acceptance_tests"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCreateBucketNoAuth(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	os.Setenv("TF_VAR_auth_token", "")
	name := "terraform_security_bucket"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBucketsResourceConfigRequired(name),
				ExpectError: regexp.MustCompile("Missing Capella Authentication Token"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBucketOrgOwner(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationOwner")
	name := "terraform_security_bucket"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBucketsResourceConfigRequired(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "name", name),
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "type", "couchbase"),
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "storage_backend", "couchstore"),
					resource.TestCheckResourceAttrSet("capella_bucket.new_bucket", "id"),
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "memory_allocation_in_mb", "105"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBucketOrgMember(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("organizationMember")
	name := "terraform_security_bucket"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBucketsResourceConfigRequired(name),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBucketProjCreator(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testAccCreateOrgAPI("projectCreator")
	name := "terraform_security_bucket"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBucketsResourceConfigRequired(name),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBucketProjOwner(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectOwner")
	name := "terraform_security_bucket"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBucketsResourceConfigRequired(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "name", name),
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "type", "couchbase"),
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "storage_backend", "couchstore"),
					resource.TestCheckResourceAttrSet("capella_bucket.new_bucket", "id"),
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "memory_allocation_in_mb", "105"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBucketProjManager(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectManager")
	name := "terraform_security_bucket"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBucketsResourceConfigRequired(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "name", name),
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "type", "couchbase"),
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "storage_backend", "couchstore"),
					resource.TestCheckResourceAttrSet("capella_bucket.new_bucket", "id"),
					resource.TestCheckResourceAttr("capella_bucket.new_bucket", "memory_allocation_in_mb", "105"),
				),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBucketProjViewer(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectViewer")
	name := "terraform_security_bucket"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBucketsResourceConfigRequired(name),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBucketDatabaseReaderWriter(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectDataReaderWriter")
	name := "terraform_security_bucket"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBucketsResourceConfigRequired(name),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func TestAccCreateBucketDatabaseReader(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testAccCreateProjAPI("organizationMember", projId, "projectDataReader")
	name := "terraform_security_bucket"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acceptance_tests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config:      testAccBucketsResourceConfigRequired(name),
				ExpectError: regexp.MustCompile("Could not create bucket"),
			},
		},
	})
	os.Setenv("TF_VAR_auth_token", tempId)
}

func testAccBucketsResourceConfigRequired(name string) string {
	return fmt.Sprintf(`
%[1]s

output "new_bucket" {
	value = capella_bucket.new_bucket
}

output "bucket_id" {
	value = capella_bucket.new_bucket.id
}

resource "capella_bucket" "new_bucket" {
	name                       = "%[2]s"
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

`, acceptance_tests.ProviderBlock, name)

}
