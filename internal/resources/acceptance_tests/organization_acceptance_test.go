package acceptance_tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"regexp"
	"terraform-provider-capella/internal/provider"
	cfg "terraform-provider-capella/internal/testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"capella": providerserver.NewProtocol6WithError(provider.New("test")()),
}

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
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "name"),
					resource.TestCheckResourceAttr("data.capella_organization.get_organization", "organization_id", organizationId),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.created_at"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.modified_by"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.modified_at"),
					resource.TestCheckResourceAttrSet("data.capella_organization.get_organization", "audit.version"),
				),
			},

			{
				Config:      testAccOrganizationResourceConfig(cfg.Cfg, "123456-abcd-4567890"),
				ExpectError: regexp.MustCompile("The server cannot or will not\nprocess the request due to something that is perceived to be a client\nerror"),
			},
		},
	})
}

func testAccOrganizationResourceConfig(cfg string, orgnizationId string) string {
	return fmt.Sprintf(`
%[1]s

output "organizations_get" {
  value = data.capella_organization.get_organization
}

data "capella_organization" "get_organization" {
  organization_id = "%[2]s"
}

`, cfg, orgnizationId)
}
