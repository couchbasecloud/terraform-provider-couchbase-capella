package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceAllowlists(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_allowlists_")
	dsName := randomStringWithPrefix("tf_acc_allowlists_ds_")
	dsReference := "data.couchbase-capella_allowlists." + dsName

	cidr := "10.7.8.9/32"
	comment := "terraform allowlists datasource acceptance test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAllowlistsDataSourceConfig(resourceName, dsName, cidr, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"cidr":    cidr,
						"comment": comment,
					}),
				),
			},
		},
	})
}

func TestAccDatasourceAllowlistsInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_allowlists_ds_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_allowlists" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				// Read() wraps the API error with "Error Reading Capella AllowLists";
				// the bogus cluster UUID hits the v4 allowedcidrs endpoint and the
				// backend returns 404 — match both so this only passes for the right
				// reason.
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella AllowLists.*"httpStatusCode":(403|404)`),
			},
		},
	})
}

func testAccAllowlistsDataSourceConfig(resourceName, dsName, cidr, comment string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_allowlist" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  cidr            = "%[7]s"
  comment         = "%[8]s"
}

data "couchbase-capella_allowlists" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"

  depends_on = [couchbase-capella_allowlist.%[5]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, dsName, cidr, comment)
}
