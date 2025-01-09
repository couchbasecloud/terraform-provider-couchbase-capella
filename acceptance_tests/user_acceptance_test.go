package acceptance_tests

import (
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"net/http"
	"testing"
)

func TestAccUserResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	resourceReference := "couchbase-capella_user." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccUserResourceConfig(resourceName, "terraform_acceptance_test1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "terraform_acceptance_test1"),
					resource.TestCheckResourceAttr(resourceReference, "email", "terraform_acceptance_test1@couchbase.com"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationOwner"),
				),
			},
			// Import state
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateUserImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read
			{
				Config: testAccUserResourceConfigUpdate(resourceName, "terraform_acceptance_test1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "terraform_acceptance_test1"),
					resource.TestCheckResourceAttr(resourceReference, "email", "terraform_acceptance_test1@couchbase.com"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.type", "project"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.roles.0", "projectViewer"),
				),
			},
			// NOTE: No delete case is provided - this occurs automatically
		},
	})
}

func TestAccUserResourceResourceNotFound(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	resourceReference := "couchbase-capella_user." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserResourceConfig(resourceName, "terraform_acceptance_test2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "terraform_acceptance_test2"),
					resource.TestCheckResourceAttr(resourceReference, "email", "terraform_acceptance_test2@couchbase.com"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationOwner"),
					// Delete the user from the server and wait until deletion is successful
					testAccDeleteUserResource(resourceReference),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},

			// Attempt to update - since the orginal has been deleted, a new user will be created.
			{
				Config: testAccUserResourceConfigUpdate(resourceName, "terraform_acceptance_test2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", "terraform_acceptance_test2"),
					resource.TestCheckResourceAttr(resourceReference, "email", "terraform_acceptance_test2@couchbase.com"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.type", "project"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.roles.0", "projectViewer"),
				),
			},
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

		data := newTestClient()
		err := deleteUserFromServer(data, rawState["organization_id"], rawState["id"])
		if err != nil {
			return err
		}
		err = readUserFromServer(data, rawState["organization_id"], rawState["id"])
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if !resourceNotFound {
			return fmt.Errorf(errString)
		}
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

func testAccUserResourceConfig(resourceName, username string) string {
	return fmt.Sprintf(`
	%[1]s
	
	resource "couchbase-capella_user" "%[2]s" {
		organization_id = "%[3]s"
	  
		name  = "%[4]s"
		email = "%[5]s"
	  
		organization_roles = [
			"organizationOwner"
		]
	  }
	`, globalProviderBlock, resourceName, globalOrgId, username, username+"@couchbase.com")
}

func testAccUserResourceConfigUpdate(resourceName, username string) string {
	return fmt.Sprintf(`
	%[1]s
	resource "couchbase-capella_user" "%[2]s" {
		organization_id = "%[3]s"
	  
		name  = "%[5]s"
		email = "%[6]s"
	  
		organization_roles = [
			"organizationMember"
		]
	  
		resources = [
		  {
			type = "project"
			id   = "%[4]s"
			roles = [
			  "projectViewer",
			]
		  }
		]
	  }
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, username, username+"@couchbase.com")
}

func generateUserImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf("id=%s,organization_id=%s", rawState["id"], globalOrgId), nil
	}
}
