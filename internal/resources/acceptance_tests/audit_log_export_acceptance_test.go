package acceptance_tests

import (
	"fmt"
	"testing"
	"time"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAuditLogExportTestCases(t *testing.T) {
	auditLogExportResourceName := "acc_audit_log_export_" + acctest.GenerateRandomResourceName()
	auditLogExportResourceReference := "couchbase-capella_audit_log_export." + auditLogExportResourceName

	resource.ParallelTest(
		t, resource.TestCase{
			ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccAuditLogExportConfig(auditLogExportResourceName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet(auditLogExportResourceReference, "id"),
					),
					ExpectNonEmptyPlan: true,
				},
			},
		},
	)
}

func testAccAuditLogExportConfig(auditLogExportResourceName string) string {
	return fmt.Sprintf(
		`
%[1]s

resource "couchbase-capella_audit_log_export" "%[5]s" {
 organization_id = "%[2]s"
 project_id = "%[3]s"
 cluster_id = "%[4]s"
 start    = "%[6]s"
 end      = "%[7]s"
}

`, ProviderBlock, OrgId, ProjectId, ClusterId, auditLogExportResourceName,
		time.Now().Add(-2*time.Hour).Format("2006-01-02T15:04:05-07:00"),
		time.Now().Add(-1*time.Hour).Format("2006-01-02T15:04:05-07:00"),
	)
}
