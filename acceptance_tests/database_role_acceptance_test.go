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

// TestAccDatabaseRoleResource tests the full lifecycle: create, read, import, update, delete.
func TestAccDatabaseRoleResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_")
	resourceReference := "couchbase-capella_database_role." + resourceName
	description := "initial description"
	updatedDescription := "updated description"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccDatabaseRoleResourceConfig(resourceName, "", description, `"dataRead"`),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", description),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.*.privileges.*", "dataRead"),
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
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateDatabaseRoleImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update description and access
			{
				Config: testAccDatabaseRoleResourceConfig(resourceName, "", updatedDescription, `"dataManage"`),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", updatedDescription),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.*.privileges.*", "dataManage"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
				),
			},
		},
	})
}

// TestAccDatabaseRoleResourceWithRequiredFieldsOnly tests creation with no optional fields.
func TestAccDatabaseRoleResourceWithRequiredFieldsOnly(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleResourceConfig(resourceName, "", "", `"dataRead"`),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.*.privileges.*", "dataRead"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
				),
			},
		},
	})
}

// TestAccDatabaseRoleResourceWithScopedAccess tests creation with bucket and scope level access.
func TestAccDatabaseRoleResourceWithScopedAccess(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleResourceConfigWithScopedAccess(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.*.privileges.*", "dataRead"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.*.privileges.*", "dataManage"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
		},
	})
}

// TestAccDatabaseRoleResourceInvalidName tests that a name with leading/trailing spaces is rejected.
func TestAccDatabaseRoleResourceInvalidName(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccDatabaseRoleResourceConfig(resourceName, "  "+resourceName+"  ", "", `"dataRead"`),
				ExpectError: regexp.MustCompile(`(?s)leading.*trailing spaces`),
			},
		},
	})
}

// TestAccDatasourceDatabaseRoles tests the list database roles datasource.
func TestAccDatasourceDatabaseRoles(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_")
	datasourceReference := "data.couchbase-capella_database_roles." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRolesDatasourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(datasourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(datasourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrSet(datasourceReference, "data.#"),
					resource.TestCheckResourceAttrSet(datasourceReference, "data.0.id"),
					resource.TestCheckResourceAttrSet(datasourceReference, "data.0.name"),
				),
			},
		},
	})
}

// --- Config Builders ---

// databaseRoleAccessBlock returns a Terraform HCL access block with a wildcard bucket
// to match the API's default response when no resources are explicitly scoped.
func databaseRoleAccessBlock(privileges string) string {
	return fmt.Sprintf(`
		access = [
			{
				privileges = [%s]
				resources = {
					buckets = [
						{
							name = "*"
						},
					]
				}
			},
		]`, privileges)
}

func testAccDatabaseRoleResourceConfig(resourceName, name, description, privileges string) string {
	if name == "" {
		name = resourceName
	}
	descBlock := ""
	if description != "" {
		descBlock = fmt.Sprintf(`description     = "%s"`, description)
	}
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_database_role" "%[2]s" {
		organization_id = "%[3]s"
		project_id      = "%[4]s"
		cluster_id      = "%[5]s"
		name            = "%[6]s"
		%[7]s
		%[8]s
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId,
		name, descBlock, databaseRoleAccessBlock(privileges))
}

func testAccDatabaseRoleResourceConfigWithScopedAccess(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_database_role" "%[2]s" {
		organization_id = "%[3]s"
		project_id      = "%[4]s"
		cluster_id      = "%[5]s"
		name            = "%[2]s"
		description     = "role with scoped access"
		access = [
			{
				privileges = ["dataRead", "dataManage"]
				resources = {
					buckets = [
						{
							name = "%[6]s"
							scopes = [
								{
									name        = "%[7]s"
									collections = ["%[8]s"]
								},
							]
						},
					]
				}
			},
		]
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId,
		globalBucketName, globalScopeName, globalCollectionName)
}

func testAccDatabaseRolesDatasourceConfig(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

	data "couchbase-capella_database_roles" "%[2]s" {
		organization_id = "%[3]s"
		project_id      = "%[4]s"
		cluster_id      = "%[5]s"
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}

// --- Import ID Generator ---

func generateDatabaseRoleImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf(
			"id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

// --- Existence Check ---

func testAccExistsDatabaseRoleResource(t *testing.T, resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		data := newTestClient(t)
		return retrieveDatabaseRoleFromServer(
			data,
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"],
			rawState["id"],
		)
	}
}

func retrieveDatabaseRoleFromServer(data *providerschema.Data, organizationId, projectId, clusterId, id string) error {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/roles/%s",
		data.HostURL, organizationId, projectId, clusterId, id,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
	if err != nil {
		return err
	}

	var roleResp api.GetDatabaseRoleResponse
	err = json.Unmarshal(response.Body, &roleResp)
	if err != nil {
		return err
	}
	if roleResp.Id.String() != id {
		return errors.ErrNotFound
	}
	return nil
}
