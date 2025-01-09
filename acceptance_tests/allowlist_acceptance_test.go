package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAllowListWithRequiredFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_allowlist_")
	resourceReference := "couchbase-capella_allowlist." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAddIpWithReqFields(resourceName, "10.1.1.1/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.1.1.1/32"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
		},
	})
}

func TestAccAllowListWithOptionalFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_allowlist_")
	resourceReference := "couchbase-capella_allowlist." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAddIpWithOptionalFields(resourceName, "10.4.5.6/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.4.5.6/32"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "expires_at"),
					resource.TestCheckResourceAttr(resourceReference, "comment", "terraform allow list acceptance test"),
				),
			},
		},
	})
}

func TestAccAllowListAllowAllIP(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_allowlist_")
	resourceReference := "couchbase-capella_allowlist." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAddIpWithOptionalFields(resourceName, "0.0.0.0/0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "expires_at"),
					resource.TestCheckResourceAttr(resourceReference, "comment", "terraform allow list acceptance test"),
				),
			},
		},
	})
}

func TestAccAllowListWithExpiredIP(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_allowlist_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAddIpWithExpiredIP(resourceName, "10.2.2.2/32"),
				ExpectError: regexp.MustCompile("The expiration time for the allowlist is not valid"),
			},
		},
	})
}

func testAccAddIpWithReqFields(resourceName, cidr string) string {
	cfg := fmt.Sprintf(`
%[1]s

resource "couchbase-capella_allowlist" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  cidr            = "%[6]s"
}

`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, cidr)
	return cfg
}

func testAccAddIpWithOptionalFields(resourceName, cidr string) string {
	timeNow := time.Now()
	timeNow = timeNow.AddDate(0, 0, 30).UTC()
	expiryTime := timeNow.Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_allowlist" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  cidr            = "%[6]s"
  comment		  = "terraform allow list acceptance test"
  expires_at      = "%[7]s"
}

`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, cidr, expiryTime)
}

func testAccAddIpWithExpiredIP(resourceName, cidr string) string {
	timeNow := time.Now().UTC()
	time.Sleep(time.Second * 10)
	expiryTime := timeNow.Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_allowlist" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  cidr            = "%[6]s"
  comment		  = "terraform allow list acceptance test"
  expires_at      = "%[7]s"
}

`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, cidr, expiryTime)
}
