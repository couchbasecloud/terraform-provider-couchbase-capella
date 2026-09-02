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
