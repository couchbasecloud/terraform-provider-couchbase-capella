package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAllowListTestCases(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//IP with required fields
			{
				Config: testAccAddIpWithReqFields(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_allowlist.add_allowlist_req", "cidr", "10.1.1.1/32"),
					resource.TestCheckResourceAttrSet("couchbase-capella_allowlist.add_allowlist_req", "id"),
				),
			},
			//IP with optional fields
			{
				Config: testAccAddIpWithOptionalFields("add_allowlist_opt", "10.4.5.6/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_allowlist.add_allowlist_opt", "cidr", "10.4.5.6/32"),
					resource.TestCheckResourceAttrSet("couchbase-capella_allowlist.add_allowlist_opt", "id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_allowlist.add_allowlist_opt", "expires_at"),
					resource.TestCheckResourceAttr("couchbase-capella_allowlist.add_allowlist_opt", "comment", "terraform allow list acceptance test"),
				),
			},
			//Unspecified IP address
			{
				Config: testAccAddIpWithOptionalFields("add_allowlist_quadzero", "0.0.0.0/0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("couchbase-capella_allowlist.add_allowlist_quadzero", "cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttrSet("couchbase-capella_allowlist.add_allowlist_quadzero", "id"),
					resource.TestCheckResourceAttrSet("couchbase-capella_allowlist.add_allowlist_quadzero", "expires_at"),
					resource.TestCheckResourceAttr("couchbase-capella_allowlist.add_allowlist_quadzero", "comment", "terraform allow list acceptance test"),
				),
			},
			//expired IP
			//expected error: "Unable to create new
			//        allowlist for database. The expiration time for the allowlist is not valid.
			//        Must be a point in time greater than now."
			{
				Config:      testAccAddIpWithExpiredIP("add_allowlist_expiredIP", "10.2.2.2/32"),
				ExpectError: regexp.MustCompile("The expiration time for the allowlist is not valid"),
			},
		},
	})
}

func testAccAddIpWithReqFields() string {
	cfg := fmt.Sprintf(`
%[1]s

output "add_allowlist_req"{
  value = couchbase-capella_allowlist.add_allowlist_req
}

resource "couchbase-capella_allowlist" "add_allowlist_req" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  cidr            = "10.1.1.1/32"
}

`, ProviderBlock, OrgId, ProjectId, ClusterId)
	return cfg
}

func testAccAddIpWithOptionalFields(resourceName, cidr string) string {
	timeNow := time.Now()
	timeNow = timeNow.AddDate(0, 0, 30).UTC()
	expiryTime := timeNow.Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

output "%[2]s"{
  value = couchbase-capella_allowlist.%[2]s
}

resource "couchbase-capella_allowlist" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  cidr            = "%[6]s"
  comment		  = "terraform allow list acceptance test"
  expires_at      = "%[7]s"
}

`, ProviderBlock, resourceName, OrgId, ProjectId, ClusterId, cidr, expiryTime)
}

func testAccAddIpWithExpiredIP(resourceName, cidr string) string {
	timeNow := time.Now().UTC()
	time.Sleep(time.Second * 10)
	expiryTime := timeNow.Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

output "%[5]s"{
  value = couchbase-capella_allowlist.%[5]s
}

resource "couchbase-capella_allowlist" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  cidr            = "%[6]s"
  comment		  = "terraform allow list acceptance test"
  expires_at      = "%[7]s"
}

`, ProviderBlock, OrgId, ProjectId, ClusterId, resourceName, cidr, expiryTime)
}
