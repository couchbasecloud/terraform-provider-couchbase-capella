package acceptance_tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	cfg "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccUserResource(t *testing.T) {
	resourceName := "acc_user_" + cfg.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_user." + resourceName
	projectResourceName := "acc_project_" + cfg.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { cfg.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: cfg.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccUserResourceConfig(cfg.Cfg, resourceName, projectResourceName, projectResourceReference),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_user_name"),
					resource.TestCheckResourceAttr(resourceReference, "email", "terraformacceptancetest@couchbase.com"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationOwner"),
				),
			},
			// Import state
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateUserImportId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccUserResourceConfigUpdate(cfg.Cfg, resourceName, projectResourceName, projectResourceReference),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_user_name"),
					resource.TestCheckResourceAttr(resourceReference, "email", "terraformacceptancetest@couchbase.com"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.type", "project"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.roles.0", "projectViewer"),
				),
			},
			// NOTE: No delete case is provided - this occurs automatically
		},
	})
}

// This function takes a resource reference string and returns a resource.TestCheckFunc. The returned function, when used
// in Terraform acceptance tests, ensures the successful deletion of the specified cluster resource. It retrieves
// the resource by name from the Terraform state, initiates the deletion, checks the status of the deletion, and
// confirms that the resource no longer exists. If the resource is successfully deleted, it returns nil; otherwise,
// it returns an error.
func testAccDeleteUserResource(resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}

		data, err := cfg.TestClient()
		if err != nil {
			return err
		}
		err = deleteUserFromServer(data, rawState["organization_id"], rawState["id"])
		if err != nil {
			return err
		}
		fmt.Printf("delete initiated")
		err = readUserFromServer(data, rawState["organization_id"], rawState["id"])
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if !resourceNotFound {
			return fmt.Errorf(errString)
		}
		fmt.Printf("successfully deleted")
		return nil
	}
}

// deleteUserFromServer deletes user from server
func deleteUserFromServer(data *providerschema.Data, organizationId, clusterId string) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/users/%s", data.HostURL, organizationId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err := data.Client.Execute(
		cfg,
		nil,
		data.Token,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// readUserFromServer reads user from server
func readUserFromServer(data *providerschema.Data, organizationId, clusterId string) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/users/%s", data.HostURL, organizationId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	_, err := data.Client.Execute(
		cfg,
		nil,
		data.Token,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func TestAccUserResourceResourceNotFound(t *testing.T) {
	resourceName := "acc_user_" + cfg.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_user." + resourceName
	projectResourceName := "acc_project_" + cfg.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { cfg.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: cfg.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig(cfg.Cfg, resourceName, projectResourceName, projectResourceReference),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_user_name"),
					resource.TestCheckResourceAttr(resourceReference, "email", "acc_test_email"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationOwner"),
					// Delete the user from the server and wait until deletion is successful
					testAccDeleteUserResource(resourceReference),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},

			// Attempt to update - since the orginal has been deleted, a new user will be created.
			{
				Config: testAccUserResourceConfig(cfg.Cfg, resourceName, projectResourceName, projectResourceReference),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "acc_test_user_name"),
					resource.TestCheckResourceAttr(resourceReference, "email", "acc_test_email"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.type", "project"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.roles.0", "projectViewer"),
				),
			},
		},
	})
}

func testAccUserResourceConfig(cfg, resourceReference, projectResourceName, projectResourceReference string) string {
	return fmt.Sprintf(`
	%[1]s
	  
	resource "couchbase-capella_project" "%[3]s" {
		organization_id = var.organization_id
		name            = "acc_test_project_name"
		description     = "description"
	}
	
	resource "couchbase-capella_user" "%[2]s" {
		organization_id = var.organization_id
	  
		name  = "acc_test_user_name"
		email = "terraformacceptancetest@couchbase.com"
	  
		organization_roles = [
			"organizationOwner"
		]
	  }
	`, cfg, resourceReference, projectResourceName, projectResourceReference)
}

func testAccUserResourceConfigUpdate(cfg, resourceReference, projectResourceName, projectResourceReference string) string {
	return fmt.Sprintf(`
	%[1]s
	  
	resource "couchbase-capellacapella_project" "%[3]s" {
		organization_id = var.organization_id
		name            = "acc_test_project_name"
		description     = "description"
	}
	
	resource "couchbase-capella_user" "%[2]s" {
		organization_id = var.organization_id
	  
		name  = "acc_test_user_name"
		email = "terraformacceptancetest@couchbase.com"
	  
		organization_roles = [
			"organizationMember"
		]
	  
		resources = [
		  {
			type = "project"
			id   = %[4]s.id
			roles = [
			  "projectViewer",
			]
		  }
		]
	  }
	`, cfg, resourceReference, projectResourceName, projectResourceReference)
}

func generateUserImportId(state *terraform.State) (string, error) {
	resourceName := "couchbase-capella_user.acc_test"
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
			"id=%s,organization_id=%s",
			rawState["id"],
			rawState["organization_id"]),
		nil
}
