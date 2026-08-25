package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Coverage for the database_roles data source. The lifecycle test in
// database_role_acceptance_test.go only asserts the list is populated; these tests
// assert the created role is actually in it and that its nested access and audit
// blocks survive the trip through the list endpoint.

// TestAccDatasourceDatabaseRolesContent verifies the role created in the same config is
// returned by the data source with its identifiers, description and audit populated.
func TestAccDatasourceDatabaseRolesContent(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_role_dsrole_")
	dsName := randomStringWithPrefix("tf_acc_db_roles_ds_")
	dsReference := "data.couchbase-capella_database_roles." + dsName
	description := "listed by the roles datasource"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRolesDatasourceContentConfig(resourceName, dsName, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					// The role created above must be one of the listed entries.
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":            resourceName,
						"description":     description,
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
						"cluster_id":      globalClusterId,
					}),
					testAccCheckDatabaseRolesEntryComplete(dsReference, resourceName),
				),
			},
		},
	})
}

// TestAccDatasourceDatabaseRolesInvalidCluster covers a well-formed but unknown cluster.
func TestAccDatasourceDatabaseRolesInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_db_roles_ds_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_database_roles" %[2]q {
  organization_id = %[3]q
  project_id      = %[4]q
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Database Roles.*"httpStatusCode":(403|404)`),
			},
		},
	})
}

// TestAccDatasourceDatabaseRolesMissingIDs covers each required ID being absent.
func TestAccDatasourceDatabaseRolesMissingIDs(t *testing.T) {
	testCases := []string{"organization_id", "project_id", "cluster_id"}

	for _, omitted := range testCases {
		t.Run(omitted, func(t *testing.T) {
			dsName := randomStringWithPrefix("tf_acc_db_roles_ds_missing_")
			ids := map[string]string{
				"organization_id": globalOrgId,
				"project_id":      globalProjectId,
				"cluster_id":      globalClusterId,
			}
			delete(ids, omitted)

			idBlock := ""
			for _, key := range []string{"organization_id", "project_id", "cluster_id"} {
				if value, ok := ids[key]; ok {
					idBlock += fmt.Sprintf("  %s = %q\n", key, value)
				}
			}

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_database_roles" %[2]q {
%[3]s}
`, globalProviderBlock, dsName, idBlock),
						ExpectError: regexp.MustCompile(
							fmt.Sprintf(`(?s)(%s|argument.*required)`, omitted)),
					},
				},
			})
		})
	}
}

func testAccDatabaseRolesDatasourceContentConfig(resourceName, dsName, description string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_role" %[5]q {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  name            = %[5]q
  description     = %[7]q
  access = %[8]s
}

data "couchbase-capella_database_roles" %[6]q {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q

  depends_on = [couchbase-capella_database_role.%[5]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, dsName,
		description, accessCollectionLevel("dataRead"))
}

// testAccCheckDatabaseRolesEntryComplete locates roleName in the data source output and
// asserts the nested attributes the list endpoint is responsible for populating: the
// role id, the audit block, and the access entry down to its collection.
//
// TestCheckTypeSetElemNestedAttrs cannot express "id is set" or reach two levels of
// nested list, so the entry is found by index and checked directly.
func testAccCheckDatabaseRolesEntryComplete(dsReference, roleName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dsReference]
		if !ok {
			return fmt.Errorf("datasource %s not found in state", dsReference)
		}
		attrs := ds.Primary.Attributes

		count, err := strconv.Atoi(attrs["data.#"])
		if err != nil {
			return fmt.Errorf("datasource %s has no data count: %w", dsReference, err)
		}
		if count == 0 {
			return fmt.Errorf("datasource %s returned no database roles", dsReference)
		}

		for i := 0; i < count; i++ {
			prefix := fmt.Sprintf("data.%d.", i)
			if attrs[prefix+"name"] != roleName {
				continue
			}

			required := []string{
				"id",
				"audit.created_at",
				"audit.created_by",
				"audit.modified_at",
				"audit.modified_by",
				"audit.version",
				"access.0.privileges.0",
				"access.0.resources.buckets.0.name",
				"access.0.resources.buckets.0.scopes.0.name",
				"access.0.resources.buckets.0.scopes.0.collections.0",
			}
			for _, suffix := range required {
				if attrs[prefix+suffix] == "" {
					return fmt.Errorf("role %q in %s is missing %s", roleName, dsReference, suffix)
				}
			}

			if got := attrs[prefix+"access.0.privileges.0"]; got != "dataRead" {
				return fmt.Errorf("role %q privilege = %q, want dataRead", roleName, got)
			}
			if got := attrs[prefix+"access.0.resources.buckets.0.name"]; got != globalBucketName {
				return fmt.Errorf("role %q bucket = %q, want %q", roleName, got, globalBucketName)
			}
			return nil
		}

		return fmt.Errorf("role %q not found among the %d roles returned by %s", roleName, count, dsReference)
	}
}
