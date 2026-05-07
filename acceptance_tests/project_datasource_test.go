package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReadProject(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_project_")
	resourceReference := "data.couchbase-capella_project." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDataSourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "name"),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "id", globalProjectId),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
				),
			},
		},
	})
}

func TestAccReadProject_InvalidID(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_project_invalid_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccProjectDataSourceConfigWithID(resourceName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile("Error Reading Capella Project|Could not read project"),
			},
		},
	})
}

func testAccProjectDataSourceConfig(resourceName string) string {
	return testAccProjectDataSourceConfigWithID(resourceName, globalProjectId)
}

func testAccProjectDataSourceConfigWithID(resourceName, projectID string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_project" "%[4]s" {
  organization_id = "%[2]s"
  id              = "%[3]s"
}

`, globalProviderBlock, globalOrgId, projectID, resourceName)
}
