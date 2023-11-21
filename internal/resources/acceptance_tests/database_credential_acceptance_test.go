package acceptance_tests

import (
	"fmt"
	acctest "terraform-provider-capella/internal/testing"
	cfg "terraform-provider-capella/internal/testing"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccDatabaseCredentialResourceWithOnlyReqFields is an acceptance test which tests
// creating and deleting a database credential which has only the
// required fields populated.
func TestAccDatabaseCredentialWithOnlyReqFields(t *testing.T) {
	resourceName := "acc_database_credential" + acctest.GenerateRandomResourceName()
	resourceReference := "capella_database_credential." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				PreConfig: func() {
					time.Sleep(1 * time.Second)
				},
				Config: generateDatabaseCredentialConfig(cfg.Cfg, map[string]string{
					"name":            "var.database_credential_name",
					"organization_id": "var.organization_id",
					"project_id":      "var.project_id",
					"cluster_id":      "var.cluster_id",
					"access":          "access",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_database_credential_name"),
					resource.TestCheckResourceAttr(resourceReference, "password", "password"),
					resource.TestCheckResourceAttr(resourceReference, "access", "access"),
				),
			},
			// NOTE: No delete case is provided - this occurs automatically
		},
	})
}

// TestAccDatabaseCredentialResourceWithOptionalField is an acceptance test which tests
// creating, reading, updating and deleting a database credential which both the
// required and optional fields populated. Importing a database credential created externally is
// also tested.
func TestAccDatabaseCredentialResourceWithOptionalField(t *testing.T) {
	resourceName := "acc_database_credential" + acctest.GenerateRandomResourceName()
	resourceReference := "capella_database_credential." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				PreConfig: func() {
					time.Sleep(1 * time.Second)
				},
				Config: generateDatabaseCredentialConfig(cfg.Cfg, map[string]string{
					"name":            "var.database_credential_name",
					"organization_id": "var.organization_id",
					"project_id":      "var.project_id",
					"cluster_id":      "var.cluster_id",
					"password":        "password",
					"access":          "access",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_database_credential_name"),
					resource.TestCheckResourceAttr(resourceReference, "password", "password"),
					resource.TestCheckResourceAttr(resourceReference, "access", "access"),
				),
			},
			// Import state
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateDatabaseCredentialImportId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: generateDatabaseCredentialConfig(cfg.Cfg, map[string]string{
					"name":            "var.database_credential_name",
					"organization_id": "var.organization_id",
					"project_id":      "var.project_id",
					"cluster_id":      "var.cluster_id",
					"password":        "updated_password",
					"access":          "access",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_project_name_update_with_if_match"),
					resource.TestCheckResourceAttr(resourceReference, "password", "updated_password"),
					resource.TestCheckResourceAttr(resourceReference, "access", "access"),
				),
			},
			// NOTE: No delete case is provided - this occurs automatically
		},
	})
}

// TestAccDatabaseCredentialInvalidScenario is a Terraform acceptance test that that simulates the
// scenario where a database credential is created with all possible fields, but with an invalid name.
func TestAccDatabaseCredentialInvalidScenario(t *testing.T) {
	// TODO: Implement test
}

// TestAccDatabaseCredentialResourceNotFound is a Terraform acceptance test that that simulates the
// scenario where a database credential is created from Terraform, but it is deleted by a REST API
// call and the deletion is successful.
//
// This test ensures that Terraform can handle the scenario where the original database credential
// no longer exists and can create a database credential with the specified configuration when updating.
func TestAccDatabaseCredentialResourceNotFound(t *testing.T) {
	// TODO: Implement test
}

// generateDatabaseCredentialConfig is used to build configs with varying fields and
// values to be stored within the fields. It constructs a config with the following format.
// Any omitted fields will not be included.
//
//	return fmt.Sprintf(`
//	%[1]s
//
//	resource "capella_database_credential" "new_database_credential" {
//		name            = <database_credential_name>
//		organization_id = <organization_id>
//		project_id      = <project_id>
//		cluster_id      = <cluster_id>
//		password        = <password>
//		access          = <access>
//	  }
//	`, cfg)
func generateDatabaseCredentialConfig(cfg string, configFields map[string]string) string {
	databaseCredentialCfg := fmt.Sprintf(`
	%[1]s

	resource "capella_database_credential" "new_database_credential" {
	`, cfg)

	// add specific fields
	for k, v := range configFields {
		databaseCredentialCfg += fmt.Sprintf("	%s= %s\n ", k, v)
	}

	// close the config
	databaseCredentialCfg += "}"
	return databaseCredentialCfg
}

func generateDatabaseCredentialImportId(state *terraform.State) (string, error) {
	resourceName := "capella_database_credential.acc_test"
	var rawState map[string]string
	for _, m := range state.Modules {
		if len(m.Resources) > 0 {
			if v, ok := m.Resources[resourceName]; ok {
				rawState = v.Primary.Attributes
			}
		}
	}
	fmt.Printf("raw state %s", rawState)

	return fmt.Sprintf(
			"id=%s,organization_id=%s,project_id=%s,cluster_id=%s",
			rawState["id"],
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"]),
		nil
}
