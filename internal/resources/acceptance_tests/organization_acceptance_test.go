package acceptance_tests

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	cfg "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.

func AccPreCheck(t *testing.T) {
	if os.Getenv("TF_VAR_host") == "" {
		t.Fatalf("host not set")
	}
	if os.Getenv("TF_VAR_auth_token") == "" {
		t.Fatalf("auth token not set")
	}
	if os.Getenv("TF_VAR_organization_id") == "" {
		t.Fatalf("organization id not set")
	}
}

func TestAccOrganizationDataSource(t *testing.T) {

	organizationId := os.Getenv("TF_VAR_organization_id")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { AccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOrganizationResourceConfig(cfg.Cfg, organizationId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.couchbase-capella_organization.get_organization", "name"),
					resource.TestCheckResourceAttr("data.couchbase-capella_organization.get_organization", "organization_id", organizationId),
					resource.TestCheckResourceAttrSet("data.couchbase-capella_organization.get_organization", "audit.created_at"),
					resource.TestCheckResourceAttrSet("data.couchbase-capella_organization.get_organization", "audit.modified_by"),
					resource.TestCheckResourceAttrSet("data.couchbase-capella_organization.get_organization", "audit.modified_at"),
					resource.TestCheckResourceAttrSet("data.couchbase-capella_organization.get_organization", "audit.version"),
				),
			},

			{
				Config:      testAccOrganizationResourceConfig(cfg.Cfg, "123456-abcd-4567890"),
				ExpectError: regexp.MustCompile("server cannot or will not process.*"),
			},
		},
	})
}

func testAccOrganizationResourceConfig(cfg string, organizationId string) string {
	return fmt.Sprintf(`
%[1]s

output "organizations_get" {
  value = data.couchbase-capella_organization.get_organization
}

data "couchbase-capella_organization" "get_organization" {
  organization_id = "%[2]s"
}

`, cfg, organizationId)
}
