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
					resource.TestCheckResourceAttrSet(dsReference, "data.0.id"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.name"),
					resource.TestCheckResourceAttr(dsReference, "data.0.organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.audit.created_at"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.audit.modified_at"),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.audit.version"),
					testAccProjectsDataSourceContains(dsReference, globalProjectId),
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
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Projects|access to the requested resource is denied|organization`),
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
				ExpectError: regexp.MustCompile(`(?s)organization_id|argument.*required`),
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

func testAccProjectsDataSourceContains(dsReference, projectId string) resource.TestCheckFunc {
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
			if attrs[fmt.Sprintf("data.%d.id", i)] == projectId {
				return nil
			}
		}
		return fmt.Errorf("expected project %q in %s.data, not found across %d entries", projectId, dsReference, count)
	}
}
