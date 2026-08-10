package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)


func TestAccBucketUpdateSendsPlanValues(t *testing.T) {

	resourceName := randomStringWithPrefix("tf_acc_bucket_upd_")
	resourceReference := "couchbase-capella_bucket." + resourceName
	bucketName := randomStringWithPrefix("tfaccbupd")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketUpdateConfig(resourceName, bucketName, 100),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", bucketName),
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "100"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "durability_level"),
					resource.TestCheckResourceAttrSet(resourceReference, "replicas"),
				),
			},
			{
				Config: testAccBucketUpdateConfig(resourceName, bucketName, 200),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "200"),
					resource.TestCheckResourceAttrSet(resourceReference, "durability_level"),
					resource.TestCheckResourceAttrSet(resourceReference, "replicas"),
				),
			},
		},
	})
}

func testAccBucketUpdateConfig(resourceName, bucketName string, memoryMiB int) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[5]s" {
  organization_id         = "%[2]s"
  project_id              = "%[3]s"
  cluster_id              = "%[4]s"
  name                    = "%[6]s"
  memory_allocation_in_mb = %[7]d
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, bucketName, memoryMiB)
}
