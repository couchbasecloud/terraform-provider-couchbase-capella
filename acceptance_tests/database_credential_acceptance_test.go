package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func TestAccDatabaseCredentialWithReqFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credential_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAddDatabaseCredWithReqFieldsConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_writer"),
					// credential_type defaults to basic when omitted.
					resource.TestCheckResourceAttr(resourceReference, "credential_type", "basic"),
					resource.TestCheckNoResourceAttr(resourceReference, "user_roles.#"),
				),
			},
		},
	})
}

func TestAccDatabaseCredentialWithOptionalFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credential_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAddDatabaseCredWithOptionalFieldsConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "password", "Secret12$#"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_writer"),
					resource.TestCheckResourceAttr(resourceReference, "credential_type", "basic"),
				),
			},
		},
	})
}

// TestAccDatabaseCredentialAdvanced tests the lifecycle of an advanced database credential:
// create with a capella user role, import, and update of the assigned user roles.
func TestAccDatabaseCredentialAdvanced(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credential_adv_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	role1Name := randomStringWithPrefix("tf_acc_db_role_")
	role2Name := randomStringWithPrefix("tf_acc_db_role_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccDatabaseCredentialAdvancedConfig(resourceName, role1Name, role2Name,
					"[couchbase-capella_database_role.role1.name]"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseCredentialResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "password", "Secret12$#"),
					resource.TestCheckResourceAttr(resourceReference, "credential_type", "advanced"),
					resource.TestCheckResourceAttr(resourceReference, "user_roles.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "user_roles.*", role1Name),
					// an advanced credential must not store any bucket-level access.
					resource.TestCheckNoResourceAttr(resourceReference, "access.#"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
				),
			},
			// ImportState
			{
				ResourceName:            resourceReference,
				ImportStateIdFunc:       generateDatabaseCredentialImportId(resourceReference),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update the assigned user roles
			{
				Config: testAccDatabaseCredentialAdvancedConfig(resourceName, role1Name, role2Name,
					"[couchbase-capella_database_role.role1.name, couchbase-capella_database_role.role2.name]"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseCredentialResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "password", "Secret12$#"),
					resource.TestCheckResourceAttr(resourceReference, "credential_type", "advanced"),
					resource.TestCheckResourceAttr(resourceReference, "user_roles.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "user_roles.*", role1Name),
					resource.TestCheckTypeSetElemAttr(resourceReference, "user_roles.*", role2Name),
					resource.TestCheckNoResourceAttr(resourceReference, "access.#"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
		},
	})
}

// TestAccDatabaseCredentialAccessAndUserRolesConflict verifies that configuring both
// access and user_roles is rejected at plan time by the exactly-one-of schema validator.
func TestAccDatabaseCredentialAccessAndUserRolesConflict(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credential_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				%[1]s

				resource "couchbase-capella_database_credential" "%[5]s" {
					name            = "%[5]s"
					organization_id = "%[2]s"
					project_id      = "%[3]s"
					cluster_id      = "%[4]s"
					access = [
						{
							privileges = ["data_writer"]
						},
					]
					user_roles = ["developer"]
				}
				`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

// TestAccDatabaseCredentialNeitherAccessNorUserRoles verifies that omitting both
// access and user_roles is rejected at plan time by the exactly-one-of schema validator.
func TestAccDatabaseCredentialNeitherAccessNorUserRoles(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credential_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				%[1]s

				resource "couchbase-capella_database_credential" "%[5]s" {
					name            = "%[5]s"
					organization_id = "%[2]s"
					project_id      = "%[3]s"
					cluster_id      = "%[4]s"
				}
				`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

// TestAccDatabaseCredentialAdvancedWithAccess verifies that an advanced credential
// configured with bucket-level access is rejected at plan time.
func TestAccDatabaseCredentialAdvancedWithAccess(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credential_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				%[1]s

				resource "couchbase-capella_database_credential" "%[5]s" {
					name            = "%[5]s"
					organization_id = "%[2]s"
					project_id      = "%[3]s"
					cluster_id      = "%[4]s"
					credential_type = "advanced"
					access = [
						{
							privileges = ["data_writer"]
						},
					]
				}
				`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName),
				ExpectError: regexp.MustCompile(`(?s)access can only be.*credential_type is.*"basic"`),
			},
		},
	})
}

// TestAccDatabaseCredentialBasicWithUserRoles verifies that a basic credential
// configured with user roles is rejected at plan time.
func TestAccDatabaseCredentialBasicWithUserRoles(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credential_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				%[1]s

				resource "couchbase-capella_database_credential" "%[5]s" {
					name            = "%[5]s"
					organization_id = "%[2]s"
					project_id      = "%[3]s"
					cluster_id      = "%[4]s"
					credential_type = "basic"
					user_roles      = ["developer"]
				}
				`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName),
				ExpectError: regexp.MustCompile(`(?s)user_roles can only be.*credential_type is.*"advanced"`),
			},
		},
	})
}

func testAccAddDatabaseCredWithReqFieldsConfig(resourceName string) string {
	return fmt.Sprintf(
		`
		%[1]s

		resource "couchbase-capella_database_credential" "%[5]s" {
			name            = "%[5]s"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			access = [
				{
					privileges = ["data_writer"]
				},
			]
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}

func testAccAddDatabaseCredWithOptionalFieldsConfig(resourceName string) string {
	return fmt.Sprintf(
		`
		%[1]s
		resource "couchbase-capella_database_credential" "%[5]s" {
			name            = "%[5]s"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			password        = "Secret12$#"
			access = [
				{
					privileges = ["data_writer"]
				},
			]
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}

// testAccDatabaseCredentialAdvancedConfig renders an advanced database credential along
// with the two database roles it can be assigned. userRolesExpr is the HCL expression
// for the user_roles attribute so tests can vary the assigned roles between steps.
func testAccDatabaseCredentialAdvancedConfig(resourceName, role1Name, role2Name, userRolesExpr string) string {
	return fmt.Sprintf(
		`
		%[1]s

		resource "couchbase-capella_database_role" "role1" {
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			name            = "%[6]s"
			%[8]s
		}

		resource "couchbase-capella_database_role" "role2" {
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			name            = "%[7]s"
			%[8]s
		}

		resource "couchbase-capella_database_credential" "%[5]s" {
			name            = "%[5]s"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			password        = "Secret12$#"
			credential_type = "advanced"
			user_roles      = %[9]s
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName,
		role1Name, role2Name, databaseRoleAccessBlock(`"dataRead"`), userRolesExpr)
}

// --- Import ID Generator ---

func generateDatabaseCredentialImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		attrs := databaseRoleAttrsFromState(state, resourceReference)
		return fmt.Sprintf(
			"id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			attrs["id"], attrs["cluster_id"], attrs["project_id"], attrs["organization_id"],
		), nil
	}
}

// --- Existence Check ---

func testAccExistsDatabaseCredentialResource(t *testing.T, resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := databaseRoleAttrsFromState(s, resourceReference)
		data := newTestClient(t)
		return retrieveDatabaseCredentialFromServer(
			data,
			attrs["organization_id"],
			attrs["project_id"],
			attrs["cluster_id"],
			attrs["id"],
		)
	}
}

func retrieveDatabaseCredentialFromServer(data *providerschema.Data, organizationId, projectId, clusterId, id string) error {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/users/%s",
		data.HostURL, organizationId, projectId, clusterId, id,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
	if err != nil {
		return err
	}

	var dbResp api.GetDatabaseCredentialResponse
	err = json.Unmarshal(response.Body, &dbResp)
	if err != nil {
		return err
	}
	if dbResp.Id.String() != id {
		return errors.ErrNotFound
	}
	return nil
}
