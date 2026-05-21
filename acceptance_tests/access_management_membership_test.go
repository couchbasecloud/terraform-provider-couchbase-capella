package acceptance_tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccDatasourceAttrAbsentOnAllElements(dsReference, listAttr, sensitiveField string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds := s.RootModule().Resources[dsReference]
		if ds == nil {
			return fmt.Errorf("datasource %s not found in state", dsReference)
		}
		count, _ := strconv.Atoi(ds.Primary.Attributes[listAttr+".#"])
		for i := 0; i < count; i++ {
			key := fmt.Sprintf("%s.%d.%s", listAttr, i, sensitiveField)
			if v, ok := ds.Primary.Attributes[key]; ok && v != "" {
				return fmt.Errorf("SECURITY: datasource %s exposed sensitive field %s on element %d (value redacted; len=%d)",
					dsReference, sensitiveField, i, len(v))
			}
		}
		return nil
	}
}

// Tests in this file cover pagination completeness, full field mapping, and
// sensitive field absence (token, password) across the list datasources.

func TestAccDatasourceAllowlistsMembership(t *testing.T) {
	a := randomStringWithPrefix("tf_acc_allowlist_mem_a_")
	b := randomStringWithPrefix("tf_acc_allowlist_mem_b_")
	c := randomStringWithPrefix("tf_acc_allowlist_mem_c_")
	dsName := randomStringWithPrefix("tf_acc_allowlists_mem_ds_")
	dsReference := "data.couchbase-capella_allowlists." + dsName

	cidrA, commentA := "198.51.100.11/32", "membership-a "+a
	cidrB, commentB := "198.51.100.12/32", "membership-b "+b
	cidrC, commentC := "198.51.100.13/32", "membership-c "+c

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAllowlistsMembershipConfig(a, b, c, dsName, cidrA, commentA, cidrB, commentB, cidrC, commentC),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"cidr":            cidrA,
						"comment":         commentA,
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
						"cluster_id":      globalClusterId,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"cidr":            cidrB,
						"comment":         commentB,
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
						"cluster_id":      globalClusterId,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"cidr":            cidrC,
						"comment":         commentC,
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
						"cluster_id":      globalClusterId,
					}),
				),
			},
		},
	})
}

func TestAccDatasourceApiKeysMembership(t *testing.T) {
	a := randomStringWithPrefix("tf_acc_apikey_mem_a_")
	b := randomStringWithPrefix("tf_acc_apikey_mem_b_")
	c := randomStringWithPrefix("tf_acc_apikey_mem_c_")
	dsName := randomStringWithPrefix("tf_acc_apikeys_mem_ds_")
	dsReference := "data.couchbase-capella_apikeys." + dsName

	descA := "membership-a " + a
	descB := "membership-b " + b
	descC := "membership-c " + c

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccApiKeysMembershipConfig(a, b, c, dsName, descA, descB, descC),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":            a,
						"description":     descA,
						"organization_id": globalOrgId,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":            b,
						"description":     descB,
						"organization_id": globalOrgId,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":            c,
						"description":     descC,
						"organization_id": globalOrgId,
					}),
					testAccDatasourceAttrAbsentOnAllElements(dsReference, "data", "token"),
				),
			},
		},
	})
}

func TestAccDatasourceDatabaseCredentialsMembership(t *testing.T) {
	a := randomStringWithPrefix("tf_acc_dbc_mem_a_")
	b := randomStringWithPrefix("tf_acc_dbc_mem_b_")
	c := randomStringWithPrefix("tf_acc_dbc_mem_c_")
	dsName := randomStringWithPrefix("tf_acc_dbc_mem_ds_")
	dsReference := "data.couchbase-capella_database_credentials." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialsMembershipConfig(a, b, c, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":            a,
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
						"cluster_id":      globalClusterId,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":            b,
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
						"cluster_id":      globalClusterId,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":            c,
						"organization_id": globalOrgId,
						"project_id":      globalProjectId,
						"cluster_id":      globalClusterId,
					}),
					testAccDatasourceAttrAbsentOnAllElements(dsReference, "data", "password"),
				),
			},
		},
	})
}

func TestAccDatasourceUsersFullFieldContent(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_users_fc_")
	dsName := randomStringWithPrefix("tf_acc_users_fc_ds_")
	resourceReference := "couchbase-capella_user." + resourceName
	dsReference := "data.couchbase-capella_users." + dsName

	username := resourceName
	email := username + "@couchbase.com"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				// perPage=0: walk all pages so the new user is found in data.*.
				Config: testAccUsersDataSourceConfig(resourceName, dsName, username, email, 0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", username),
					resource.TestCheckResourceAttr(resourceReference, "email", email),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),

					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":              username,
						"email":             email,
						"organization_id":   globalOrgId,
						"organization_roles.#": "1",
						"organization_roles.0": "organizationMember",
					}),
				),
			},
		},
	})
}

func testAccAllowlistsMembershipConfig(a, b, c, dsName, cidrA, commentA, cidrB, commentB, cidrC, commentC string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_allowlist" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  cidr            = "%[8]s"
  comment         = "%[9]s"
}

resource "couchbase-capella_allowlist" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  cidr            = "%[10]s"
  comment         = "%[11]s"
}

resource "couchbase-capella_allowlist" "%[7]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  cidr            = "%[12]s"
  comment         = "%[13]s"
}

data "couchbase-capella_allowlists" "%[14]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"

  depends_on = [
    couchbase-capella_allowlist.%[5]s,
    couchbase-capella_allowlist.%[6]s,
    couchbase-capella_allowlist.%[7]s,
  ]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, a, b, c, cidrA, commentA, cidrB, commentB, cidrC, commentC, dsName)
}

func testAccApiKeysMembershipConfig(a, b, c, dsName, descA, descB, descC string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_apikey" "%[4]s" {
  organization_id    = "%[2]s"
  name               = "%[4]s"
  description        = "%[8]s"
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

resource "couchbase-capella_apikey" "%[5]s" {
  organization_id    = "%[2]s"
  name               = "%[5]s"
  description        = "%[9]s"
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

resource "couchbase-capella_apikey" "%[6]s" {
  organization_id    = "%[2]s"
  name               = "%[6]s"
  description        = "%[10]s"
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

data "couchbase-capella_apikeys" "%[7]s" {
  organization_id = "%[2]s"

  depends_on = [
    couchbase-capella_apikey.%[4]s,
    couchbase-capella_apikey.%[5]s,
    couchbase-capella_apikey.%[6]s,
  ]
}
`, globalProviderBlock, globalOrgId, globalProjectId, a, b, c, dsName, descA, descB, descC)
}

func testAccDatabaseCredentialsMembershipConfig(a, b, c, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" "%[5]s" {
  name            = "%[5]s"
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  access = [{ privileges = ["data_writer"] }]
}

resource "couchbase-capella_database_credential" "%[6]s" {
  name            = "%[6]s"
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  access = [{ privileges = ["data_writer"] }]
}

resource "couchbase-capella_database_credential" "%[7]s" {
  name            = "%[7]s"
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  access = [{ privileges = ["data_writer"] }]
}

data "couchbase-capella_database_credentials" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"

  depends_on = [
    couchbase-capella_database_credential.%[5]s,
    couchbase-capella_database_credential.%[6]s,
    couchbase-capella_database_credential.%[7]s,
  ]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, a, b, c, dsName)
}
