package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDatasourceProjects(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_projects_ds_")
	dsReference := "data.couchbase-capella_projects." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectsDataSourceConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
					// Locate globalProjectId in the list and assert required
					// fields on that specific entry — asserting data.0.* alone
					// would silently pass when data.0 is a different project
					// and globalProjectId sits at a later index.
					testAccCheckProjectsDataSourceContainsWithFields(dsReference, globalProjectId, globalOrgId),
				),
			},
		},
	})
}

func TestAccDatasourceProjectsInvalidOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_projects_ds_invalid_org_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_projects" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName),
				// Require the provider's specific summary AND a 403/404 from the
				// API. A bare "|organization" matched many unrelated diagnostics
				// (auth/transport/etc.) and could pass for the wrong reason.
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Projects.*"httpStatusCode":(403|404)`),
			},
		},
	})
}

func TestAccDatasourceProjectsMissingOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_projects_ds_missing_org_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_projects" "%[2]s" {}
`, globalProviderBlock, dsName),
				// Match Terraform's exact "argument X is required" diagnostic for
				// the missing organization_id, instead of any error mentioning
				// the field name.
				ExpectError: regexp.MustCompile(`(?s)The argument "organization_id" is required`),
			},
		},
	})
}

func testAccProjectsDataSourceConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_projects" "%[3]s" {
  organization_id = "%[2]s"
}
`, globalProviderBlock, globalOrgId, dsName)
}

// testAccCheckProjectsDataSourceContainsWithFields locates the entry with
// id == projectId in the projects datasource list and asserts that this
// specific entry has the required computed fields populated and the expected
// organization_id. Asserting at data.0.* alone is unsafe because list ordering
// is not guaranteed and other projects in the tenant may sit ahead of the one
// under test.
func testAccCheckProjectsDataSourceContainsWithFields(dsReference, projectId, orgId string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[dsReference]
		if !ok {
			return fmt.Errorf("data source %q not found in state", dsReference)
		}
		attrs := rs.Primary.Attributes
		count, err := strconv.Atoi(attrs["data.#"])
		if err != nil {
			return fmt.Errorf("invalid data.# on %q: %w", dsReference, err)
		}
		for i := 0; i < count; i++ {
			if attrs[fmt.Sprintf("data.%d.id", i)] != projectId {
				continue
			}
			if got := attrs[fmt.Sprintf("data.%d.organization_id", i)]; got != orgId {
				return fmt.Errorf("data.%d.organization_id = %q, want %q", i, got, orgId)
			}
			for _, suffix := range []string{
				"name",
				"audit.created_at",
				"audit.modified_at",
				"audit.version",
			} {
				key := fmt.Sprintf("data.%d.%s", i, suffix)
				if attrs[key] == "" {
					return fmt.Errorf("attribute %q expected to be set on matched project %s", key, projectId)
				}
			}
			return nil
		}
		return fmt.Errorf("expected project %q in %s.data, not found across %d entries", projectId, dsReference, count)
	}
}
