package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources"
)

func TestAccUserResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	username := resourceName
	resourceReference := "couchbase-capella_user." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccUserResourceConfig(resourceName, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", username),
					resource.TestCheckResourceAttr(resourceReference, "email", username+"@couchbase.com"),
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
				Config: testAccUserResourceConfigUpdate(resourceName, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", username),
					resource.TestCheckResourceAttr(resourceReference, "email", username+"@couchbase.com"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.type", "project"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.roles.0", "projectViewer"),
				),
			},
			// NOTE: No delete case is provided - this occurs automatically
		},
	})
}

func TestAccUserResourceWithProjectRoles(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	username := resourceName
	resourceReference := "couchbase-capella_user." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfigWithProjectRoles(resourceName, username),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", username),
					resource.TestCheckResourceAttr(resourceReference, "email", username+"@couchbase.com"),
					resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", "organizationMember"),
					resource.TestCheckResourceAttr(resourceReference, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "resources.0.type", "project"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "resources.0.roles.*", "projectViewer"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "resources.0.roles.*", "projectDataReaderWriter"),
				),
			},
		},
	})
}

func TestAccUserResourceOrganizationRoles(t *testing.T) {
	for _, role := range resources.ValidOrganizationRoles() {
		t.Run(role, func(t *testing.T) {
			resourceName := randomStringWithPrefix("tf_acc_user_")
			username := resourceName
			resourceReference := "couchbase-capella_user." + resourceName

			// Not parallel: concurrent user deletes return HTTP 500 from the Capella API.
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: testAccUserResourceConfigWithOrganizationRoles(resourceName, username, fmt.Sprintf(`["%s"]`, role)),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(resourceReference, "email", username+"@couchbase.com"),
							resource.TestCheckResourceAttr(resourceReference, "organization_roles.#", "1"),
							resource.TestCheckResourceAttr(resourceReference, "organization_roles.0", role),
						),
					},
					{
						ResourceName:      resourceReference,
						ImportStateIdFunc: generateUserImportIdForResource(resourceReference),
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			})
		})
	}
}

func TestAccUserResourceMultipleOrganizationRoles(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	username := resourceName
	resourceReference := "couchbase-capella_user." + resourceName
	roles := resources.ValidOrganizationRoles()

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceReference, "organization_roles.#", strconv.Itoa(len(roles))),
	}
	for i, role := range roles {
		checks = append(checks, resource.TestCheckResourceAttr(resourceReference, "organization_roles."+strconv.Itoa(i), role))
	}

	// Not parallel: concurrent user deletes return HTTP 500 from the Capella API.
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfigWithOrganizationRoles(resourceName, username, fmt.Sprintf(`["%s"]`, strings.Join(roles, `", "`))),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateUserImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccUserResourceUpdateOrganizationRoles(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	username := resourceName
	resourceReference := "couchbase-capella_user." + resourceName
	roles := resources.ValidOrganizationRoles()

	// Each step appends one role. The API appends roles added by update to the
	// end of the list, so inserting anywhere else fails with an inconsistent result.
	var steps []resource.TestStep
	for i := range roles {
		assigned := roles[:i+1]

		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceReference, "organization_roles.#", strconv.Itoa(len(assigned))),
		}
		for j, role := range assigned {
			checks = append(checks, resource.TestCheckResourceAttr(resourceReference, "organization_roles."+strconv.Itoa(j), role))
		}

		step := resource.TestStep{
			Config: testAccUserResourceConfigWithOrganizationRoles(resourceName, username, fmt.Sprintf(`["%s"]`, strings.Join(assigned, `", "`))),
			Check:  resource.ComposeAggregateTestCheckFunc(checks...),
		}
		if i > 0 {
			step.ConfigPlanChecks = resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
				},
			}
		}
		steps = append(steps, step)
	}

	// Not parallel: concurrent user deletes return HTTP 500 from the Capella API.
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps:                    steps,
	})
}

func TestAccUserResourceInvalidScenarioEmptyOrganizationRoles(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	username := resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccUserResourceConfigWithOrganizationRoles(resourceName, username, "[]"),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value.*Attribute organization_roles list must contain at least 1 elements, got: 0`),
			},
		},
	})
}

func TestAccUserResourceInvalidScenarioEmptyOrganizationID(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	username := resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccUserResourceConfigWithOrganizationID(resourceName, username, ""),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value.*Attribute organization_id string length must be at least 1, got: 0`),
			},
		},
	})
}

func TestAccUserResourceInvalidScenarioMissingEmail(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccUserResourceConfigMissingEmail(resourceName),
				ExpectError: regexp.MustCompile(`(?s)Missing required argument.*email`),
			},
		},
	})
}

func TestAccUserResourceInvalidScenarioInvalidOrganizationRole(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	username := resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccUserResourceConfigWithOrganizationRoles(resourceName, username, `["organizationWizard"]`),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.*organizationWizard.*organizationMember.*organizationOwner.*projectCreator.*organizationBillingAdmin.*organizationBillingViewer.*organizationReadOnly`),
			},
		},
	})
}

func TestAccUserResourceInvalidScenarioInvalidProjectRole(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	username := resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccUserResourceConfigWithProjectRole(resourceName, username, "projectWizard"),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.*projectWizard.*projectOwner.*projectManager.*projectViewer.*projectDataReaderWriter.*projectDataReader`),
			},
		},
	})
}

func TestAccUserResourceInvalidScenarioResourceIdEmpty(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_user_")
	username := resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccUserResourceConfigWithResourceId(resourceName, username, ""),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Length.*Attribute\s+resources.*\.id\s+string length must be at least 1, got: 0`),
			},
		},
	})
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

func testAccUserResourceConfigWithProjectRoles(resourceName, username string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_user" "%[2]s" {
	organization_id = "%[3]s"
	name            = "%[4]s"
	email           = "%[5]s"

	organization_roles = ["organizationMember"]
	resources = [
		{
			type  = "project"
			id    = "%[6]s"
			roles = ["projectViewer", "projectDataReaderWriter"]
		}
	]
}
`, globalProviderBlock, resourceName, globalOrgId, username, username+"@couchbase.com", globalProjectId)
}

func testAccUserResourceConfigWithOrganizationRoles(resourceName, username, organizationRoles string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_user" "%[2]s" {
	organization_id = "%[3]s"
	name            = "%[4]s"
	email           = "%[5]s"

	organization_roles = %[6]s
}
`, globalProviderBlock, resourceName, globalOrgId, username, username+"@couchbase.com", organizationRoles)
}

func testAccUserResourceConfigWithOrganizationID(resourceName, username, organizationID string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_user" "%[2]s" {
	organization_id = "%[3]s"
	name            = "%[4]s"
	email           = "%[5]s"

	organization_roles = ["organizationMember"]
}
`, globalProviderBlock, resourceName, organizationID, username, username+"@couchbase.com")
}

func testAccUserResourceConfigMissingEmail(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_user" "%[2]s" {
	organization_id    = "%[3]s"
	organization_roles = ["organizationMember"]
}
`, globalProviderBlock, resourceName, globalOrgId)
}

func testAccUserResourceConfigWithProjectRole(resourceName, username, projectRole string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_user" "%[2]s" {
	organization_id = "%[3]s"
	name            = "%[4]s"
	email           = "%[5]s"

	organization_roles = ["organizationMember"]
	resources = [
		{
			type  = "project"
			id    = "%[6]s"
			roles = ["%[7]s"]
		}
	]
}
`, globalProviderBlock, resourceName, globalOrgId, username, username+"@couchbase.com", globalProjectId, projectRole)
}

func testAccUserResourceConfigWithResourceId(resourceName, username, resourceId string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_user" "%[2]s" {
	organization_id = "%[3]s"
	name            = "%[4]s"
	email           = "%[5]s"

	organization_roles = ["organizationMember"]
	resources = [
		{
			type  = "project"
			id    = "%[6]s"
			roles = ["projectViewer"]
		}
	]
}
`, globalProviderBlock, resourceName, globalOrgId, username, username+"@couchbase.com", resourceId)
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
