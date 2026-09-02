package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// Regression guards for defects found reviewing the fine-grained RBAC work that are
// documented but not yet fixed. Each test asserts the behaviour we want rather than
// the behaviour we have, and is skipped with its ticket so it becomes a guard the
// moment the fix lands - unskip, do not rewrite.

// TestAccDatabaseRole_AV_142505_DescriptionCanBeCleared asserts that a description set
// on a database role can be cleared by assigning an empty string.
//
// AV-142505: `Description` is `string` with `json:"description,omitempty"` on both
// database role request bodies, so an empty string is dropped from the PUT entirely,
// the API keeps the previous description, and the GET reads that stale value back into
// state. Because `description` is Optional+Computed the plan value is a known "", so
// the apply fails with "Provider produced inconsistent result after apply".
//
// This is deliberately distinct from TestAccDatabaseRoleDescriptionLifecycle, whose
// final step passes "" to testAccDatabaseRoleConfigAccess. That helper omits the
// attribute entirely when the description is empty, so that step exercises attribute
// removal - which Optional+Computed answers by keeping the prior value - and asserts
// only TestCheckResourceAttrSet, which holds either way. The config builder below
// emits an explicit `description = ""` instead, which is the case that breaks.
//
// Same defect class as AV-136448 on the eventing function resource, fixed there.
func TestAccDatabaseRole_AV_142505_DescriptionCanBeCleared(t *testing.T) {
	t.Skip("AV-142505: clearing a database role description fails with an inconsistent result after apply; unskip once the description field is transmitted when empty")

	resourceName := randomStringWithPrefix("tf_acc_role_clrdesc_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			// Create with a description.
			{
				Config: testAccDatabaseRoleConfigExplicitDescription(resourceName, resourceName, "first description", accessCollectionLevel("dataRead")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", "first description"),
				),
			},
			// Clear it. The apply must succeed and the description must be empty,
			// not the value the API still holds.
			{
				Config: testAccDatabaseRoleConfigExplicitDescription(resourceName, resourceName, "", accessCollectionLevel("dataRead")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					// The rest of the role must survive the clear untouched.
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
				),
			},
			// A refresh must not resurrect the old description.
			{
				Config:   testAccDatabaseRoleConfigExplicitDescription(resourceName, resourceName, "", accessCollectionLevel("dataRead")),
				PlanOnly: true,
			},
		},
	})
}

// testAccDatabaseRoleConfigExplicitDescription always emits a description attribute,
// including when the value is empty. testAccDatabaseRoleConfigAccess omits the
// attribute in that case, which is a different scenario - see AV-142505 above.
func testAccDatabaseRoleConfigExplicitDescription(resourceName, roleName, description, accessHCL string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_role" %[2]q {
  organization_id = %[3]q
  project_id      = %[4]q
  cluster_id      = %[5]q
  name            = %[6]q
  description     = %[7]q
  access          = %[8]s
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, roleName, description, accessHCL)
}

// TestAccDatasourceDatabasePrivileges_AV_142506_ReturnsEveryPage asserts the database
// privileges data source returns the complete privilege set, not just whatever the
// first response carried.
//
// AV-142506: DatabasePrivileges.listPrivileges issues a single unpaged
// ClientV1.ExecuteWithRetry and returns the envelope's Data slice, while the sibling
// database_roles data source walks every page through api.GetPaginated. If the
// privileges endpoint paginates the way other V4 list endpoints do, the data source
// silently truncates and a practitioner sees valid privileges as unavailable.
//
// The assertion compares the data source against api.GetPaginated over the same
// endpoint, which is exactly what the proposed fix uses, so the two agree once the fix
// lands. TestAccDatasourceDatabasePrivileges cannot catch this: it asserts only that
// the result is non-empty, which a truncated list satisfies.
//
// Note the outcome if the endpoint turns out not to paginate at all: both fetches
// return the same set, this test passes, and AV-142506 should be closed as Not a Bug.
// That is a legitimate result of unskipping, not a reason to change the assertion.
func TestAccDatasourceDatabasePrivileges_AV_142506_ReturnsEveryPage(t *testing.T) {
	t.Skip("AV-142506: the database_privileges data source does not paginate; unskip once listPrivileges walks every page, or to confirm the endpoint returns everything unpaged")

	dsName := randomStringWithPrefix("tf_acc_db_privs_pages_")
	dsReference := "data.couchbase-capella_database_privileges." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabasePrivilegesDataSourceConfig(dsName, globalClusterId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDatabasePrivilegesNonEmpty(dsReference),
					testAccCheckDatabasePrivilegesMatchAllPages(dsReference),
				),
			},
		},
	})
}

// testAccCheckDatabasePrivilegesMatchAllPages fetches every page of the privileges
// endpoint directly and asserts the data source returned the same privilege names.
// A data source that stops after one page yields a strict subset, and the diff names
// exactly which privileges went missing.
func testAccCheckDatabasePrivilegesMatchAllPages(dsReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds := s.RootModule().Resources[dsReference]
		if ds == nil {
			return fmt.Errorf("datasource %s not found in state", dsReference)
		}

		count, err := strconv.Atoi(ds.Primary.Attributes["data.#"])
		if err != nil {
			return fmt.Errorf("datasource %s has no readable data count: %w", dsReference, err)
		}

		fromDatasource := make([]string, 0, count)
		for i := 0; i < count; i++ {
			fromDatasource = append(fromDatasource, ds.Primary.Attributes[fmt.Sprintf("data.%d.name", i)])
		}

		fromAPI, err := allDatabasePrivilegeNames(globalOrgId, globalProjectId, globalClusterId)
		if err != nil {
			return err
		}

		sort.Strings(fromDatasource)
		sort.Strings(fromAPI)

		if missing := missingFrom(fromDatasource, fromAPI); len(missing) > 0 {
			return fmt.Errorf(
				"datasource %s returned %d of %d privileges; missing %v",
				dsReference, len(fromDatasource), len(fromAPI), missing,
			)
		}
		if extra := missingFrom(fromAPI, fromDatasource); len(extra) > 0 {
			return fmt.Errorf(
				"datasource %s returned privileges the API did not: %v",
				dsReference, extra,
			)
		}
		return nil
	}
}

// allDatabasePrivilegeNames reads every page of the privileges endpoint. An empty
// sortParameter leaves sortBy off the query string, which matters because privileges
// carry no id to sort by.
func allDatabasePrivilegeNames(organizationId, projectId, clusterId string) ([]string, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/privileges",
		globalHost, organizationId, projectId, clusterId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	privileges, err := api.GetPaginated[[]api.GetDatabasePrivilegeResponse](
		context.Background(), globalClient, globalToken, cfg, "",
	)
	if err != nil {
		return nil, fmt.Errorf("listing every page of database privileges: %w", err)
	}

	names := make([]string, 0, len(privileges))
	for _, privilege := range privileges {
		names = append(names, privilege.Name)
	}
	return names, nil
}

// missingFrom returns the entries of want that have absent from got.
func missingFrom(got, want []string) []string {
	present := make(map[string]struct{}, len(got))
	for _, name := range got {
		present[name] = struct{}{}
	}

	var missing []string
	for _, name := range want {
		if _, ok := present[name]; !ok {
			missing = append(missing, name)
		}
	}
	return missing
}
