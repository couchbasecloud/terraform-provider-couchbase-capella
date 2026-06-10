package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceCertificate(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_certificate_ds_")
	dsReference := "data.couchbase-capella_certificate." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateDatasourceConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrSet(dsReference, "data.0.certificate"),
				),
			},
		},
	})
}

func TestAccDatasourceCertificateInvalidCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_certificate_ds_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_certificate" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`Error Reading Capella Certificate`),
			},
		},
	})
}

func testAccCertificateDatasourceConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_certificate" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId)
}
