package acceptance_tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// TestAccBucketUpdateSendsPlanValues covers a defect in the bucket Update path: updating
// a bucket fails, or silently discards settings, whenever an Optional+Computed attribute
// is absent from configuration.
//
// Mechanism. On an update plan the framework marks every Computed attribute that is null
// in configuration as unknown (fwserver.MarkComputedNilsAsUnknown, which returns early
// only for attributes carrying a Default). Bucket's Update then builds the PUT body from
// the PLAN — note internal/resources/bucket.go reads req.Plan into a variable confusingly
// named "state" — using ValueString()/ValueInt64(), which return the zero value for an
// unknown. PutBucketRequest has no omitempty on any field, so the zero values are sent:
//
//	durability_level        -> ""  API rejects with 422, code 6005
//	replicas                -> 0   would silently drop replicas
//	time_to_live_in_seconds -> 0   would silently reset TTL
//
// flush is unaffected because it carries a Default, which exempts it from being marked
// unknown.
//
// Observed on 2026-07-30 against a live cluster: creating a bucket with only
// memory_allocation_in_mb set, then changing it, fails with
//
//	The durability level '' provided is not supported. The supported levels are
//	'none', 'majority', 'persistToMajority', and 'majorityAndPersistActive'.
//
// That is the common way to resize a bucket, so bucket updates are effectively broken
// unless durability_level is written explicitly. Present in released v1.9.1.
//
// The plan check below asserts the update is planned in place, and it passes today —
// it is kept to pin that behaviour, because an earlier reading of this schema wrongly
// predicted a destroy-and-recreate here. The failure is in the apply, not the plan.
//
// Gated so it does not red CI before the fix. Remove the gate with the fix.
func TestAccBucketUpdateSendsPlanValues(t *testing.T) {
	if os.Getenv("TF_ACC_BUCKET_UPDATE") == "" {
		t.Skip("known failure; set TF_ACC_BUCKET_UPDATE=1 to reproduce")
	}

	resourceName := randomStringWithPrefix("tf_acc_bucket_upd_")
	resourceReference := "couchbase-capella_bucket." + resourceName
	bucketName := randomStringWithPrefix("tfaccbupd")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create with only memory_allocation_in_mb — durability_level, replicas and
			// time_to_live_in_seconds are left to the API, which is the ordinary way to
			// declare a bucket.
			{
				Config: testAccBucketUpdateConfig(resourceName, bucketName, 100),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", bucketName),
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "100"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					// Populated by the API on create; these are what the update must
					// preserve rather than overwrite with zero values.
					resource.TestCheckResourceAttrSet(resourceReference, "durability_level"),
					resource.TestCheckResourceAttrSet(resourceReference, "replicas"),
				),
			},
			// Resize the bucket. The plan is an in-place update; the apply is what fails.
			{
				Config: testAccBucketUpdateConfig(resourceName, bucketName, 200),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "memory_allocation_in_mb", "200"),
					// The regression that matters: settings the practitioner never
					// specified must survive the update rather than being zeroed.
					resource.TestCheckResourceAttrSet(resourceReference, "durability_level"),
					resource.TestCheckResourceAttrSet(resourceReference, "replicas"),
				),
			},
		},
	})
}

// testAccBucketUpdateConfig sets only name and memory_allocation_in_mb. Setting
// durability_level, replicas or time_to_live_in_seconds would put non-null values in
// configuration and mask the defect.
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
