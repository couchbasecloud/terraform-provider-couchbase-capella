package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDatasourceDatabasePrivileges(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_db_privs_ds_")
	dsReference := "data.couchbase-capella_database_privileges." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabasePrivilegesDataSourceConfig(dsName, globalClusterId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					testAccCheckDatabasePrivilegesNonEmpty(dsReference),
				),
			},
		},
	})
}

func TestAccDatasourceDatabasePrivilegesInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_db_privs_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccDatabasePrivilegesDataSourceConfig(dsName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Database Privileges.*"httpStatusCode":(403|404)`),
			},
		},
	})
}

func TestAccDatasourceDatabasePrivilegesMissingCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_db_privs_missing_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_database_privileges" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)cluster_id|argument.*required`),
			},
		},
	})
}

func TestAccDatasourceDatabasePrivilegesMissingOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_db_privs_missing_org_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_database_privileges" "%[2]s" {
  project_id = "%[3]s"
  cluster_id = "%[4]s"
}
`, globalProviderBlock, dsName, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`(?s)organization_id|argument.*required`),
			},
		},
	})
}

func TestAccDatasourceDatabasePrivilegesMissingProject(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_db_privs_missing_project_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_database_privileges" "%[2]s" {
  organization_id = "%[3]s"
  cluster_id      = "%[4]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalClusterId),
				ExpectError: regexp.MustCompile(`(?s)project_id|argument.*required`),
			},
		},
	})
}

// testAccCheckDatabasePrivilegesNonEmpty verifies the datasource returned at
// least one privilege item with name and group populated.
func testAccCheckDatabasePrivilegesNonEmpty(dsReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds := s.RootModule().Resources[dsReference]
		if ds == nil {
			return fmt.Errorf("datasource %s not found in state", dsReference)
		}
		count, _ := strconv.Atoi(ds.Primary.Attributes["data.#"])
		if count == 0 {
			return fmt.Errorf("datasource %s returned no database privileges", dsReference)
		}
		for i := 0; i < count; i++ {
			name := ds.Primary.Attributes[fmt.Sprintf("data.%d.name", i)]
			group := ds.Primary.Attributes[fmt.Sprintf("data.%d.group", i)]
			if name != "" && group != "" {
				return nil
			}
		}
		return fmt.Errorf("datasource %s has %d elements but none had name and group both set", dsReference, count)
	}
}

func testAccDatabasePrivilegesDataSourceConfig(dsName, clusterID string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_database_privileges" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, clusterID)
}

// privilegeLevel is the deepest resource level a privilege may be granted at, as
// declared by the RBAC template in the data source's resources attribute.
type privilegeLevel int

const (
	// privilegeLevelCluster has no resources block at all.
	privilegeLevelCluster privilegeLevel = iota
	// privilegeLevelBucket has a bucket with no scopes.
	privilegeLevelBucket
	// privilegeLevelScope has a bucket and a scope but no collections.
	privilegeLevelScope
	// privilegeLevelCollection nests all the way down to collections.
	privilegeLevelCollection
)

func (l privilegeLevel) String() string {
	switch l {
	case privilegeLevelCluster:
		return "cluster"
	case privilegeLevelBucket:
		return "bucket"
	case privilegeLevelScope:
		return "scope"
	case privilegeLevelCollection:
		return "collection"
	}
	return "unknown"
}

// TestAccDatasourceDatabasePrivilegesRBACTemplate asserts the data source reproduces the
// RBAC template faithfully, not merely that it returned something.
//
// This is the regression guard for AV-141822, where the privileges endpoint's
// {"data": [...]} envelope was unmarshalled into a bare slice and every read failed. It
// also pins the level each privilege may be granted at, which is what a config author
// consults to avoid the 422 covered by TestAccDatabaseRolePrivilegeLevelMismatch.
func TestAccDatasourceDatabasePrivilegesRBACTemplate(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_db_privs_tmpl_")
	dsReference := "data.couchbase-capella_database_privileges." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabasePrivilegesDataSourceConfig(dsName, globalClusterId),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Every privilege in the catalogue must carry a group.
					testAccCheckDatabasePrivilegesAllPopulated(dsReference),
					// One representative privilege per level.
					testAccCheckDatabasePrivilegeLevel(dsReference, "dataRead", "data", privilegeLevelCollection),
					testAccCheckDatabasePrivilegeLevel(dsReference, "queryManage", "query", privilegeLevelScope),
					testAccCheckDatabasePrivilegeLevel(dsReference, "viewsReader", "data", privilegeLevelBucket),
					testAccCheckDatabasePrivilegeLevel(dsReference, "statsRead", "global", privilegeLevelCluster),
				),
			},
		},
	})
}

// testAccCheckDatabasePrivilegesAllPopulated requires every returned privilege to have
// both a name and a group. The pre-existing non-empty check passes as soon as one entry
// is well formed, which would hide a partially decoded envelope.
func testAccCheckDatabasePrivilegesAllPopulated(dsReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, count, err := databasePrivilegeAttrs(s, dsReference)
		if err != nil {
			return err
		}
		for i := 0; i < count; i++ {
			if attrs[fmt.Sprintf("data.%d.name", i)] == "" {
				return fmt.Errorf("privilege at index %d has no name", i)
			}
			if attrs[fmt.Sprintf("data.%d.group", i)] == "" {
				return fmt.Errorf("privilege %q at index %d has no group",
					attrs[fmt.Sprintf("data.%d.name", i)], i)
			}
		}
		return nil
	}
}

// testAccCheckDatabasePrivilegeLevel finds privilegeName in the data source output and
// asserts its group and the depth of its RBAC template.
func testAccCheckDatabasePrivilegeLevel(dsReference, privilegeName, wantGroup string, wantLevel privilegeLevel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, count, err := databasePrivilegeAttrs(s, dsReference)
		if err != nil {
			return err
		}

		for i := 0; i < count; i++ {
			prefix := fmt.Sprintf("data.%d.", i)
			if attrs[prefix+"name"] != privilegeName {
				continue
			}

			if got := attrs[prefix+"group"]; got != wantGroup {
				return fmt.Errorf("privilege %q group = %q, want %q", privilegeName, got, wantGroup)
			}

			gotLevel := privilegeLevelCluster
			switch {
			case attrs[prefix+"resources.buckets.0.scopes.0.collections.#"] != "":
				gotLevel = privilegeLevelCollection
			case attrs[prefix+"resources.buckets.0.scopes.#"] != "":
				gotLevel = privilegeLevelScope
			case attrs[prefix+"resources.buckets.#"] != "":
				gotLevel = privilegeLevelBucket
			}

			if gotLevel != wantLevel {
				return fmt.Errorf("privilege %q level = %s, want %s", privilegeName, gotLevel, wantLevel)
			}
			return nil
		}

		return fmt.Errorf("privilege %q not found among the %d privileges returned by %s",
			privilegeName, count, dsReference)
	}
}

func databasePrivilegeAttrs(s *terraform.State, dsReference string) (map[string]string, int, error) {
	ds, ok := s.RootModule().Resources[dsReference]
	if !ok {
		return nil, 0, fmt.Errorf("datasource %s not found in state", dsReference)
	}
	count, err := strconv.Atoi(ds.Primary.Attributes["data.#"])
	if err != nil {
		return nil, 0, fmt.Errorf("datasource %s has no data count: %w", dsReference, err)
	}
	if count == 0 {
		return nil, 0, fmt.Errorf("datasource %s returned no database privileges", dsReference)
	}
	return ds.Primary.Attributes, count, nil
}
