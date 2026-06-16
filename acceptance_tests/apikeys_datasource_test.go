package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceApiKeys(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_apikeys_")
	dsName := randomStringWithPrefix("tf_acc_apikeys_ds_")
	resourceReference := "couchbase-capella_apikey." + resourceName
	dsReference := "data.couchbase-capella_apikeys." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccApiKeysDataSourceConfig(resourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Confirm the parent apikey was created so we can rely on it
					// being in the datasource response.
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),

					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":            resourceName,
						"description":     "terraform apikeys datasource acceptance test",
						"organization_id": globalOrgId,
					}),
				),
			},
		},
	})
}

func TestAccDatasourceApiKeysInvalidOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_apikeys_ds_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_apikeys" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName),
				// Read() emits "Error Reading Capella ApiKeys" and the backend rejects
				// the bogus org with a 403/404 — both must appear so we know the
				// failure came from the apikeys list call, not transport/auth.
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella ApiKeys.*"httpStatusCode":(403|404)`),
			},
		},
	})
}

func TestAccDatasourceApiKeysEmptyOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_apikeys_ds_empty_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_apikeys" "%[2]s" {
  organization_id = ""
}
`, globalProviderBlock, dsName),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value.*Attribute organization_id string length must be at least 1, got: 0`),
			},
		},
	})
}

func testAccApiKeysDataSourceConfig(resourceName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[4]s" {
  organization_id    = "%[2]s"
  name               = "%[4]s"
  description        = "terraform apikeys datasource acceptance test"
  expiry             = 180
  organization_roles = ["organizationMember"]
  allowed_cidrs      = ["10.1.42.0/23"]
  resources = [
    {
      id    = "%[3]s"
      roles = ["projectManager", "projectDataReader"]
      type  = "project"
    }
  ]
}

data "couchbase-capella_apikeys" "%[5]s" {
  organization_id = "%[2]s"

  depends_on = [couchbase-capella_apikey.%[4]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, dsName)
}
