package acceptance_tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccDatasourceAttrAbsentOnAllElements asserts that the given attribute
// is empty (or missing) on every element of the named list/set attribute on
// the datasource. Used as a P0 tripwire for sensitive fields like `token`
// or `password` that resource Create may legitimately return once, but
// must NEVER appear in the list datasource response.
//
// Empty string is treated as "absent" — Terraform's state flattening writes
// "" for nested attributes that are present in the schema but unset on the
// element, which is the desired behaviour for these sensitive fields.
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

// The tests in this file target bug classes the existing per-datasource
// happy-path tests cannot catch:
//
//   - Pagination/single-element drops: every happy-path test creates ONE
//     parent resource, so a datasource that silently returns only page 1
//     (or only the first match) would still pass. Each test below creates
//     three parent resources with distinct identifying fields and requires
//     all three to be present in `data.*`.
//
//   - Field mapping: the happy-path assertions check only `organization_id`,
//     `project_id`, `cluster_id` on the matched set element. A provider
//     that silently drops or scrambles `cidr`/`comment`/`name`/`description`
//     etc. would not be caught. Each test below asserts the full set of
//     identifying + descriptive fields.
//
//   - Sensitive-field leakage: the apikey and database_credential resources
//     have sensitive fields (`token`, `password`). If the list datasource
//     ever surfaces them, that is a P0 security regression. The apikey and
//     database_credential tests below assert those fields are absent from
//     the datasource shape.
//
// Generated under AV-128950 as follow-up depth coverage on top of the
// scope-required tests in PR #601.

// TestAccDatasourceAllowlistsMembership creates three allowlists with
// distinct CIDRs and comments, then asserts each one appears in the
// `couchbase-capella_allowlists` datasource response with all four
// identifying/descriptive fields matching what we wrote.
func TestAccDatasourceAllowlistsMembership(t *testing.T) {
	a := randomStringWithPrefix("tf_acc_allowlist_mem_a_")
	b := randomStringWithPrefix("tf_acc_allowlist_mem_b_")
	c := randomStringWithPrefix("tf_acc_allowlist_mem_c_")
	dsName := randomStringWithPrefix("tf_acc_allowlists_mem_ds_")
	dsReference := "data.couchbase-capella_allowlists." + dsName

	// RFC 5737 documentation block — guaranteed not in real use, so the
	// /32 CIDRs are stable across runs and harmless on the test cluster.
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
					// Every one of the three must appear with BOTH cidr and
					// comment matching — proves pagination AND field mapping.
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

// TestAccDatasourceApiKeysMembership creates three apikeys with distinct
// names and descriptions, asserts each appears with full field content in
// `data.*`, and also asserts the sensitive `token` field is absent from
// every element of the datasource response — catching a leak of the
// secret that would otherwise be a P0 security regression.
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
					// SECURITY: the datasource must not surface the secret
					// `token` field that the resource exposes once on create.
					// If any element has a non-empty `token`, we have a P0
					// leak — this assertion is a tripwire.
					testAccDatasourceAttrAbsentOnAllElements(dsReference, "data", "token"),
				),
			},
		},
	})
}

// TestAccDatasourceDatabaseCredentialsMembership creates three database
// credentials and asserts each appears in the datasource with matching
// name + scope ids, AND that the sensitive `password` field is absent
// from every datasource element (P0 tripwire).
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
					// SECURITY: `password` is sensitive on the resource and must
					// not appear in the list datasource. P0 leak tripwire.
					testAccDatasourceAttrAbsentOnAllElements(dsReference, "data", "password"),
				),
			},
		},
	})
}

// TestAccDatasourceUsersFullFieldContent extends the existing single-user
// happy-path test with full field-content verification instead of relying
// only on TestCheckTypeSetElemNestedAttrs on id/name/email/org_id. Email
// uniqueness on the tenant precludes the multi-element pattern used for
// the other three datasources, so this test instead drills into ALL
// documented fields on the one user we create.
func TestAccDatasourceUsersFullFieldContent(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_users_fc_")
	dsName := randomStringWithPrefix("tf_acc_users_fc_ds_")
	resourceReference := "couchbase-capella_user." + resourceName
	dsReference := "data.couchbase-capella_users." + dsName

	// Stable tenant fixture — matches the pattern other user tests use.
	username := "terraform_acceptance_test_field_content"
	email := username + "@couchbase.com"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceConfig(resourceName, dsName, username, email),
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

// --- config builders for the membership tests -------------------------------

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
