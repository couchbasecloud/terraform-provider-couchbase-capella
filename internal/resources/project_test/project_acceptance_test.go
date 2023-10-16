package project_test

import (
	"fmt"
	"testing"

	"terraform-provider-capella/internal/provider"
	cfg "terraform-provider-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"capella": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for
	// example assertions about the appropriate environment variables being set
	// are common to see in a pre-check function.
}

func TestAccProjectResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig(cfg.Cfg),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("capella_project.acc_test", "description", "description"),
					//resource.TestCheckResourceAttr("cml2_group.test", "labs.#", "2"),
					//resource.TestCheckResourceAttr("cml2_group.test", "members.#", "1"),
				),
			},
			//// ImportState testing
			//{
			//	ResourceName:      "cml2_group.test",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
			//// Update and Read testing
			//{
			//	Config: testAccGroupResourceConfigUpdate(cfg.Cfg),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("cml2_group.test", "name", "new name"),
			//		resource.TestCheckResourceAttr("cml2_group.test", "description", "new description"),
			//		resource.TestCheckResourceAttr("cml2_group.test", "labs.#", "1"),
			//		resource.TestCheckResourceAttr("cml2_group.test", "members.#", "0"),
			//	),
			//},
			//{
			//	Config: testAccGroupResourceConfigUpdate2(cfg.Cfg, "read_write"),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("cml2_group.test", "labs.#", "1"),
			//		resource.TestCheckResourceAttr("cml2_group.test", "members.#", "2"),
			//		resource.TestCheckTypeSetElemNestedAttrs("cml2_group.test", "labs.*", map[string]string{
			//			"permission": "read_write",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccGroupResourceConfigUpdate2(cfg.Cfg, "read_only"),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("cml2_group.test", "labs.#", "1"),
			//		resource.TestCheckResourceAttr("cml2_group.test", "members.#", "2"),
			//		resource.TestCheckTypeSetElemNestedAttrs("cml2_group.test", "labs.*", map[string]string{
			//			"permission": "read_only",
			//		}),
			//	),
			//},
			// Delete testing automatically occurs in TestCase
		},
	})
}

//func TestAccGroupResourceNoLists(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { testAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			// Create and Read testing
//			{
//				Config: testAccGroupResourceConfigNoLists(cfg.Cfg),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("cml2_group.test", "description", "description"),
//					resource.TestCheckResourceAttr("cml2_group.test", "labs.#", "0"),
//					resource.TestCheckResourceAttr("cml2_group.test", "members.#", "0"),
//				),
//			},
//			// Delete testing automatically occurs in TestCase
//		},
//	})
//}

func testAccProjectResourceConfig(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "capella_project" "acc_test" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}
`, cfg)
}

func testAccGroupResourceConfigUpdate(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "cml2_user" "acc_test" {
	username      = "acc_test_group_user"
	password      = "süpersücret"
	fullname      = "firstname, lastname"
	email         = "bla@cml.lab"
	description   = "acc test user description"
	is_admin      = false
}

resource "cml2_lab" "lab1" {
	title       = "group_acc_test_lab1"
}

resource "cml2_lab" "lab2" {
	title       = "group_acc_test_lab2"
}

resource "cml2_group" "test" {
	description = "new description"
	name = "new name"
	members = []
	labs = [
		{
			id = cml2_lab.lab1.id
			permission = "read_only"
		}
	]
}
`, cfg)
}

func testAccGroupResourceConfigUpdate2(cfg, permission string) string {
	return fmt.Sprintf(`
%[1]s

resource "cml2_user" "acc_test" {
	username      = "acc_test_group_user"
	password      = "süpersücret"
	fullname      = "firstname, lastname"
	email         = "bla@cml.lab"
	description   = "acc test user description"
	is_admin      = false
}

resource "cml2_user" "acc_test_2" {
	username      = "acc_test_group_user_2"
	password      = "süpersücret"
	fullname      = "firstname, lastname"
	email         = "bla@cml.lab"
	description   = "acc test user description"
	is_admin      = false
}

resource "cml2_lab" "lab1" {
	title       = "group_acc_test_lab1"
}

resource "cml2_group" "test" {
	description = "new description"
	name = "new name"
	members = [ cml2_user.acc_test.id, cml2_user.acc_test_2.id ]
	labs = [
		{
			id = cml2_lab.lab1.id
			permission = %[2]q
		},
	]
}
`, cfg, permission)
}

func testAccGroupResourceConfigNoLists(cfg string) string {
	return fmt.Sprintf(`
%[1]s

resource "cml2_group" "test" {
	description = "description"
	name = "new name"
}
`, cfg)
}
