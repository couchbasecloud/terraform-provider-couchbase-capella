package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Negative and edge-case coverage for the database_role resource.
//
// The cases split into two groups. Schema validators and Terraform core reject the
// first group before any request is made, so those steps never reach the V4 API. The
// second group is only caught server-side, and the expected messages below were taken
// from live 4xx responses rather than inferred, so a change in API behaviour surfaces
// here rather than silently widening what the provider accepts.

// --- Rejected before the request is sent ---

// TestAccDatabaseRoleEmptyAccess covers the SizeAtLeast(1) validator on access.
// `access = []` satisfies Required but would send a role carrying no privileges.
func TestAccDatabaseRoleEmptyAccess(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_role_noacc_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccDatabaseRoleConfigAccess(resourceName, resourceName, "", `[]`),
				ExpectError: regexp.MustCompile(`(?s)access.*must contain at least 1`),
			},
		},
	})
}

// TestAccDatabaseRoleMissingAccess covers access being Required.
func TestAccDatabaseRoleMissingAccess(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_role_missacc_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_role" %[5]q {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  name            = %[5]q
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName),
				ExpectError: regexp.MustCompile(`(?s)(argument "access" is required|Missing required argument)`),
			},
		},
	})
}

// TestAccDatabaseRoleMissingRequiredIDs covers each of the three required IDs being
// absent. Terraform core rejects the config before the provider is consulted.
func TestAccDatabaseRoleMissingRequiredIDs(t *testing.T) {
	testCases := []struct {
		name    string
		omitted string
	}{
		{name: "organization_id", omitted: "organization_id"},
		{name: "project_id", omitted: "project_id"},
		{name: "cluster_id", omitted: "cluster_id"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceName := randomStringWithPrefix("tf_acc_role_missid_")
			ids := map[string]string{
				"organization_id": globalOrgId,
				"project_id":      globalProjectId,
				"cluster_id":      globalClusterId,
			}
			delete(ids, tc.omitted)

			idBlock := ""
			for _, key := range []string{"organization_id", "project_id", "cluster_id"} {
				if value, ok := ids[key]; ok {
					idBlock += fmt.Sprintf("  %s = %q\n", key, value)
				}
			}

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_role" %[3]q {
%[2]s  name = %[3]q
  access = %[4]s
}
`, globalProviderBlock, idBlock, resourceName, accessCollectionLevel("dataRead")),
						ExpectError: regexp.MustCompile(
							fmt.Sprintf(`(?s)(argument "%s" is required|Missing required argument)`, tc.omitted)),
					},
				},
			})
		})
	}
}

// TestAccDatabaseRoleInvalidUUID covers the UUID regex validator carried by the three
// ID attributes.
func TestAccDatabaseRoleInvalidUUID(t *testing.T) {
	testCases := []struct {
		name           string
		organizationId string
		projectId      string
		clusterId      string
	}{
		{name: "organization_id", organizationId: "not-a-uuid", projectId: globalProjectId, clusterId: globalClusterId},
		{name: "project_id", organizationId: globalOrgId, projectId: "not-a-uuid", clusterId: globalClusterId},
		{name: "cluster_id", organizationId: globalOrgId, projectId: globalProjectId, clusterId: "not-a-uuid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceName := randomStringWithPrefix("tf_acc_role_baduuid_")

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_role" %[5]q {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  name            = %[5]q
  access = %[6]s
}
`, globalProviderBlock, tc.organizationId, tc.projectId, tc.clusterId, resourceName,
							accessCollectionLevel("dataRead")),
						ExpectError: regexp.MustCompile(`(?s)must be a valid UUID`),
					},
				},
			})
		})
	}
}

// --- Rejected by the V4 API ---

// TestAccDatabaseRolePrivilegeLevelMismatch is the AV-141821 regression: queryManage is
// a scope-level privilege, so narrowing it to collections is rejected. The example
// shipped in examples/database_role carried exactly this shape.
func TestAccDatabaseRolePrivilegeLevelMismatch(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_role_lvlmis_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "",
					accessCollectionLevel("queryManage")),
				ExpectError: apiErrorPattern("Error creating database role",
					"collection-level access is not a allowed for the mentioned privilege", "422"),
			},
		},
	})
}

// TestAccDatabaseRoleUnknownPrivilege covers a privilege name absent from the cluster's
// RBAC template. Nothing client-side validates privilege names, so this reaches the API.
func TestAccDatabaseRoleUnknownPrivilege(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_role_badpriv_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "",
					accessCollectionLevel("notARealPrivilege")),
				ExpectError: apiErrorPattern("Error creating database role",
					"is not a valid privilege", "422"),
			},
		},
	})
}

// TestAccDatabaseRoleEmptyPrivileges covers `privileges = []`. The privileges attribute
// carries no SizeAtLeast validator, so an empty set is only rejected server-side.
func TestAccDatabaseRoleEmptyPrivileges(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_role_nopriv_")

	accessHCL := fmt.Sprintf(`[
    {
      privileges = []
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
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "", accessHCL),
				ExpectError: apiErrorPattern("Error creating database role",
					"Provided privileges list must not be empty", "422"),
			},
		},
	})
}


func TestAccDatabaseRoleNonExistentKeyspace(t *testing.T) {
	testCases := []struct {
		name      string
		accessHCL func() string
		phrase    func() string
	}{
		{
			name: "bucket",
			phrase: func() string {
				return "references bucket 'tf_acc_no_such_bucket', which does not exist on the cluster"
			},
			accessHCL: func() string {
				return `[
    {
      privileges = ["dataRead"]
      resources = {
        buckets = [
          {
            name = "tf_acc_no_such_bucket"
          },
        ]
      }
    },
  ]`
			},
		},
		{
			name: "scope",
			phrase: func() string {
				return fmt.Sprintf(
					"references scope 'tf_acc_no_such_scope' in bucket '%s', which does not exist on the cluster",
					globalBucketName)
			},
			accessHCL: func() string {
				return fmt.Sprintf(`[
    {
      privileges = ["dataRead"]
      resources = {
        buckets = [
          {
            name = %q
            scopes = [
              {
                name = "tf_acc_no_such_scope"
              },
            ]
          },
        ]
      }
    },
  ]`, globalBucketName)
			},
		},
		{
			name: "collection",
			phrase: func() string {
				return fmt.Sprintf(
					"references collection 'tf_acc_no_such_collection' in scope '%s' of bucket '%s', which does not exist on the cluster",
					globalScopeName, globalBucketName)
			},
			accessHCL: func() string {
				return fmt.Sprintf(`[
    {
      privileges = ["dataRead"]
      resources = {
        buckets = [
          {
            name = %q
            scopes = [
              {
                name        = %q
                collections = ["tf_acc_no_such_collection"]
              },
            ]
          },
        ]
      }
    },
  ]`, globalBucketName, globalScopeName)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceName := randomStringWithPrefix("tf_acc_role_badks_")

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				CheckDestroy:             testAccCheckDatabaseRoleDestroy,
				Steps: []resource.TestStep{
					{
						Config:      testAccDatabaseRoleConfigAccess(resourceName, resourceName, "", tc.accessHCL()),
						ExpectError: apiErrorPattern("Error creating database role", tc.phrase(), "422"),
					},
				},
			})
		})
	}
}

// TestAccDatabaseRoleDuplicateName covers two roles sharing a name on one cluster.
// Role names are unique per cluster and the API rejects the second create.
//
// The two roles are applied in separate steps on purpose. Declaring both in one config
// lets Terraform create them concurrently, and the uniqueness check did not hold under
// that race: an earlier version of this test declared both together and the apply
// succeeded with no error. Sequencing the creates keeps the test about the uniqueness
// rule rather than about concurrency.
func TestAccDatabaseRoleDuplicateName(t *testing.T) {
	firstResource := randomStringWithPrefix("tf_acc_role_dup_a_")
	secondResource := randomStringWithPrefix("tf_acc_role_dup_b_")
	sharedRoleName := randomStringWithPrefix("tf_acc_role_shared_")

	firstConfig := testAccDatabaseRoleConfigAccess(firstResource, sharedRoleName, "",
		accessCollectionLevel("dataRead"))

	// Second resource, same role name. The provider block comes from firstConfig.
	secondRoleBlock := fmt.Sprintf(`
resource "couchbase-capella_database_role" %[1]q {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  name            = %[5]q
  access = %[6]s
}
`, secondResource, globalOrgId, globalProjectId, globalClusterId, sharedRoleName,
		accessBucketLevel("viewsReader"))

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: firstConfig,
				Check:  testAccExistsDatabaseRoleResource(t, "couchbase-capella_database_role."+firstResource),
			},
			{
				Config: firstConfig + secondRoleBlock,
				ExpectError: apiErrorPattern("Error creating database role",
					"is already in use", "422"),
			},
		},
	})
}

// TestAccDatabaseRoleInvalidCluster covers a well-formed but unknown cluster UUID.
func TestAccDatabaseRoleInvalidCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_role_badclus_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_role" %[4]q {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = "00000000-0000-0000-0000-000000000000"
  name            = %[4]q
  access = %[5]s
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName,
					accessCollectionLevel("dataRead")),
				ExpectError: apiErrorPattern("Error creating database role", "", "(403|404)"),
			},
		},
	})
}
