package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Gaps in database_credential coverage introduced alongside fine-grained RBAC.
//
// The existing suite covers the two credential types and the four
// access/user_roles/credential_type conflict combinations. What it did not cover is the
// size validator on user_roles, the credential_type enum, the requires-replace
// behaviour when switching type, updating a basic credential's access in place, and
// the API's answer to a user role that does not exist.

// TestAccDatabaseCredentialEmptyUserRoles covers the SizeAtLeast(1) validator on
// user_roles. `user_roles = []` satisfies ExactlyOneOf but carries no permission, the
// mirror of the `access = []` case already covered for basic credentials.
func TestAccDatabaseCredentialEmptyUserRoles(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_cred_empty_roles_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" %[5]q {
  name            = %[5]q
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  credential_type = "advanced"
  user_roles      = []
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName),
				ExpectError: regexp.MustCompile(`(?s)user_roles.*(must contain at least 1|Invalid Attribute Combination)`),
			},
		},
	})
}

// TestAccDatabaseCredentialInvalidCredentialType covers the OneOf validator attached to
// credential_type. The validator is wired by hand because the published OpenAPI spec
// does not yet carry the enum, so it is worth pinning.
func TestAccDatabaseCredentialInvalidCredentialType(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_cred_bad_type_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" %[5]q {
  name            = %[5]q
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  credential_type = "superuser"
  access = [
    {
      privileges = ["data_writer"]
    },
  ]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName),
				ExpectError: regexp.MustCompile(`(?s)credential_type.*must be one of`),
			},
		},
	})
}

// TestAccDatabaseCredentialTypeRequiresReplace verifies that switching a credential from
// basic to advanced replaces it. credential_type carries a RequiresReplace plan modifier
// because the V4 API cannot convert an existing credential between the two forms.
func TestAccDatabaseCredentialTypeRequiresReplace(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_cred_retype_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	roleName := randomStringWithPrefix("tf_acc_db_role_")

	var originalID string

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRBACDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialBasicWithRoleConfig(resourceName, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseCredentialResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "credential_type", "basic"),
					testAccCheckCaptureResourceID(resourceReference, &originalID),
				),
			},
			{
				Config: testAccDatabaseCredentialAdvancedWithRoleConfig(resourceName, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseCredentialResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "credential_type", "advanced"),
					resource.TestCheckResourceAttr(resourceReference, "user_roles.#", "1"),
					resource.TestCheckNoResourceAttr(resourceReference, "access.#"),
					testAccCheckResourceIDChanged(resourceReference, &originalID),
				),
			},
		},
	})
}

// TestAccDatabaseCredentialBasicAccessUpdate updates a basic credential's access in
// place: first widening the privileges, then narrowing the resources to one bucket.
func TestAccDatabaseCredentialBasicAccessUpdate(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_cred_access_upd_")
	resourceReference := "couchbase-capella_database_credential." + resourceName

	clusterWideAccess := `[
    {
      privileges = ["data_reader"]
    },
  ]`
	widenedAccess := `[
    {
      privileges = ["data_reader", "data_writer"]
    },
  ]`
	bucketScopedAccess := fmt.Sprintf(`[
    {
      privileges = ["data_reader"]
      resources = {
        buckets = [
          {
            name = %q
          },
        ]
      }
    },
  ]`, globalBucketName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialAccessConfig(resourceName, clusterWideAccess),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseCredentialResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.#", "1"),
				),
			},
			{
				Config: testAccDatabaseCredentialAccessConfig(resourceName, widenedAccess),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseCredentialResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "data_reader"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "data_writer"),
				),
			},
			{
				Config: testAccDatabaseCredentialAccessConfig(resourceName, bucketScopedAccess),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseCredentialResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.name", globalBucketName),
				),
			},
		},
	})
}

// TestAccDatabaseCredentialUntrimmedName covers the trimmed-name check on create.
func TestAccDatabaseCredentialUntrimmedName(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_cred_untrimmed_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" %[5]q {
  name            = "  %[5]s  "
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  access = [
    {
      privileges = ["data_writer"]
    },
  ]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName),
				ExpectError: regexp.MustCompile(`(?s)leading.*trailing spaces`),
			},
		},
	})
}

// TestAccDatabaseCredentialNonExistentUserRole covers an advanced credential naming a
// role that was never created. Nothing client-side resolves role names, so the API is
// what rejects it.
func TestAccDatabaseCredentialNonExistentUserRole(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_cred_bad_role_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" %[5]q {
  name            = %[5]q
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  password        = "Secret12$#"
  credential_type = "advanced"
  user_roles      = ["tf_acc_no_such_role"]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName),
				ExpectError: apiErrorPattern("Error creating database credential",
					"application user roles do not exist on this cluster", "422"),
			},
		},
	})
}

// TestAccDatabaseCredentialInvalidCluster covers a well-formed but unknown cluster UUID.
func TestAccDatabaseCredentialInvalidCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_cred_bad_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" %[4]q {
  name            = %[4]q
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = "00000000-0000-0000-0000-000000000000"
  access = [
    {
      privileges = ["data_writer"]
    },
  ]
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName),
				ExpectError: apiErrorPattern("Error creating database credential", "", "(403|404)"),
			},
		},
	})
}

// --- Config builders ---

func testAccDatabaseCredentialAccessConfig(resourceName, accessHCL string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" %[5]q {
  name            = %[5]q
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  password        = "Secret12$#"
  access          = %[6]s
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, accessHCL)
}

// testAccDatabaseCredentialBasicWithRoleConfig renders a basic credential alongside the
// role the advanced variant will reference, so the two steps of the requires-replace
// test differ only in the credential itself.
func testAccDatabaseCredentialBasicWithRoleConfig(resourceName, roleName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_role" "retype_role" {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  name            = %[6]q
  access = %[7]s
}

resource "couchbase-capella_database_credential" %[5]q {
  name            = %[5]q
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  password        = "Secret12$#"
  credential_type = "basic"
  access = [
    {
      privileges = ["data_writer"]
    },
  ]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, roleName,
		accessCollectionLevel("dataRead"))
}

func testAccDatabaseCredentialAdvancedWithRoleConfig(resourceName, roleName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_role" "retype_role" {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  name            = %[6]q
  access = %[7]s
}

resource "couchbase-capella_database_credential" %[5]q {
  name            = %[5]q
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  password        = "Secret12$#"
  credential_type = "advanced"
  user_roles      = [couchbase-capella_database_role.retype_role.name]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName, roleName,
		accessCollectionLevel("dataRead"))
}
