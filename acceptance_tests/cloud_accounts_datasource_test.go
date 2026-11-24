package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccCloudAccountsDataSource reads the test organization's cloud accounts. The shared test
// cluster is on AWS, so the organization always has a Capella AWS account. It may have no GCP
// or Azure account, so those aren't checked.
func TestAccCloudAccountsDataSource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cloud_accounts_")
	resourceReference := "data.couchbase-capella_cloud_accounts." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudAccountsDataSourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestMatchResourceAttr(resourceReference, "aws_capella_account", regexp.MustCompile(`^\d{12}$`)),
				),
			},
		},
	})
}

func testAccCloudAccountsDataSourceConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_cloud_accounts" "%[3]s" {
  organization_id = "%[2]s"
}
`, globalProviderBlock, globalOrgId, resourceName)
}
