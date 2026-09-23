package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Access-shape coverage for the database_role resource.
//
// Every Capella privilege declares, in the RBAC template returned by the cluster's
// /privileges endpoint, the deepest resource level it may be granted at: collection,
// scope, bucket, or cluster-wide with no resources at all. Granting a privilege below
// its permitted level is rejected by the V4 API with a 422 (see
// TestAccDatabaseRolePrivilegeLevelMismatch in the negative tests).
//
// These tests pin one role per level plus a role combining all three scoped levels,
// which is what exercises reconcileAccess/reconcileBuckets/reconcileScopes/
// reconcileCollections end to end. Each apply is followed by the framework's built-in
// post-apply plan check, so a clean run also proves the round-trip through the API
// leaves no perpetual diff.

// TestAccDatabaseRoleCollectionLevelAccess grants a collection-level privilege at its
// deepest level: bucket -> scope -> collection.
func TestAccDatabaseRoleCollectionLevelAccess(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_coll_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "collection level", accessCollectionLevel("dataRead")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "dataRead"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.name", globalBucketName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.name", globalScopeName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.collections.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.collections.0", globalCollectionName),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateDatabaseRoleImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccDatabaseRoleScopeLevelAccess grants a scope-level privilege with the
// collections list omitted. queryManage cannot be narrowed to collections, so omitting
// it is the only valid shape; supplying collections is what AV-141821 hit.
func TestAccDatabaseRoleScopeLevelAccess(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_scope_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "scope level", accessScopeLevel("queryManage")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "queryManage"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.name", globalBucketName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.name", globalScopeName),
					// collections omitted stays null rather than becoming an empty list.
					resource.TestCheckNoResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.collections.#"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateDatabaseRoleImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccDatabaseRoleBucketLevelAccess grants a bucket-level privilege with no scopes.
func TestAccDatabaseRoleBucketLevelAccess(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_bucket_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "bucket level", accessBucketLevel("viewsReader")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "viewsReader"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.name", globalBucketName),
					// scopes omitted stays null.
					resource.TestCheckNoResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.#"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateDatabaseRoleImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccDatabaseRoleGlobalPrivilege grants a cluster-level privilege with the resources
// block omitted entirely, which is the only valid shape for privileges such as statsRead.
//
// This is the round-trip the provider works hardest to hide. createAccessFromSlice sends
// an empty buckets slice for an entry with no resources, and the V4 API answers the
// subsequent read with an implicit wildcard bucket, which mergeAccessEntry suppresses via
// isWildcardOnlyResourcesSchema. If that suppression regresses, the post-apply plan
// stops being empty and this test fails rather than the drift reaching a user.
//
// No import step: on import there is no prior state to reconcile against, so the wildcard
// the API supplies is kept and ImportStateVerify would compare it against the null the
// original apply stored. See the note in the summary for AV-141821/AV-141822 follow-up.
func TestAccDatabaseRoleGlobalPrivilege(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_global_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	accessHCL := `[
    {
      privileges = ["statsRead"]
    },
  ]`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "cluster level", accessHCL),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "statsRead"),
					// The wildcard bucket the API adds must not surface in state.
					resource.TestCheckNoResourceAttr(resourceReference, "access.0.resources.buckets.#"),
				),
			},
			// Re-applying the identical config must be a no-op. An empty plan here is the
			// assertion that the implicit wildcard is still being suppressed on read.
			{
				Config:   testAccDatabaseRoleConfigAccess(resourceName, resourceName, "cluster level", accessHCL),
				PlanOnly: true,
			},
		},
	})
}

// TestAccDatabaseRoleAllScopedLevels combines collection, scope and bucket level entries
// in a single role. The V4 API returns access entries keyed by their resources, so this
// is the case that exercises the whole reconcile chain: entries are paired by privilege
// set, buckets by name, scopes by name and collections by value. A clean post-apply plan
// means the ordering survived the round-trip.
func TestAccDatabaseRoleAllScopedLevels(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_multi_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "all scoped levels", accessAllScopedLevels()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "3"),
					// Entry order must match the configured order.
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "dataRead"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.1.privileges.*", "queryManage"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.2.privileges.*", "viewsReader"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.collections.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.1.resources.buckets.0.scopes.0.name", globalScopeName),
					resource.TestCheckResourceAttr(resourceReference, "access.2.resources.buckets.0.name", globalBucketName),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateDatabaseRoleImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccDatabaseRoleAccessAddAndRemove grows the access list from one entry to three
// and back down to one, so both the append and the drop paths of the update are covered.
func TestAccDatabaseRoleAccessAddAndRemove(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_addrm_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "", accessCollectionLevel("dataRead")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
				),
			},
			// Grow to three entries.
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "", accessAllScopedLevels()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.2.privileges.*", "viewsReader"),
				),
			},
			// Shrink back to one.
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "", accessCollectionLevel("dataRead")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "dataRead"),
				),
			},
		},
	})
}

// TestAccDatabaseRoleMultiplePrivilegesInOneEntry grants several privileges that share
// one resource block. The API returns them as a single entry, so the privilege set must
// round-trip without reordering into a diff.
func TestAccDatabaseRoleMultiplePrivilegesInOneEntry(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_role_mpriv_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	accessHCL := fmt.Sprintf(`[
    {
      privileges = ["dataRead", "dataManage", "dataMonitor"]
      resources = {
        buckets = [
          {
            name = %q
            scopes = [
              {
                name        = %q
                collections = [%q]
              },
            ]
          },
        ]
      }
    },
  ]`, globalBucketName, globalScopeName, globalCollectionName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "", accessHCL),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "dataRead"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "dataManage"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "dataMonitor"),
				),
			},
		},
	})
}

// TestAccDatabaseRoleMultipleCollections grants one privilege across several collections
// in the same scope, covering reconcileCollections with more than one element.
//
// The collections are created in the same config rather than assumed: the V4 API
// answers a role referencing a keyspace that does not exist with a 400, so the role
// must depend on collections that really are present on the cluster.
func TestAccDatabaseRoleMultipleCollections(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_role_mcoll_")
	resourceReference := "couchbase-capella_database_role." + resourceName
	collectionA := randomStringWithPrefix("tfacccolla")
	collectionB := randomStringWithPrefix("tfacccollb")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleMultiCollectionConfig(resourceName, collectionA, collectionB),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.collections.#", "2"),
					// Configured order must be preserved: collections is a List, not a Set.
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.collections.0", collectionA),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.collections.1", collectionB),
				),
			},
		},
	})
}

// testAccDatabaseRoleMultiCollectionConfig provisions two collections in the shared
// bucket's default scope and grants a role access to both.
func testAccDatabaseRoleMultiCollectionConfig(resourceName, collectionA, collectionB string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_collection" %[6]q {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  bucket_id       = %[5]q
  scope_name      = %[8]q
  collection_name = %[6]q
}

resource "couchbase-capella_collection" %[7]q {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  bucket_id       = %[5]q
  scope_name      = %[8]q
  collection_name = %[7]q
}

resource "couchbase-capella_database_role" %[9]q {
  organization_id = %[2]q
  project_id      = %[3]q
  cluster_id      = %[4]q
  name            = %[9]q
  access = [
    {
      privileges = ["dataRead"]
      resources = {
        buckets = [
          {
            name = %[10]q
            scopes = [
              {
                name        = %[8]q
                collections = [%[6]q, %[7]q]
              },
            ]
          },
        ]
      }
    },
  ]

  depends_on = [
    couchbase-capella_collection.%[6]s,
    couchbase-capella_collection.%[7]s,
  ]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, globalBucketId,
		collectionA, collectionB, globalScopeName, resourceName, globalBucketName)
}

// TestAccDatabaseRoleNameRequiresReplace verifies that renaming a role replaces it
// rather than updating in place: name carries a RequiresReplace plan modifier because
// the V4 API has no rename operation.
func TestAccDatabaseRoleNameRequiresReplace(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_role_replace_")
	resourceReference := "couchbase-capella_database_role." + resourceName
	originalRoleName := randomStringWithPrefix("tf_acc_role_original_")
	renamedRoleName := randomStringWithPrefix("tf_acc_role_renamed_")

	var originalID string

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, originalRoleName, "", accessCollectionLevel("dataRead")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", originalRoleName),
					testAccCheckCaptureResourceID(resourceReference, &originalID),
				),
			},
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, renamedRoleName, "", accessCollectionLevel("dataRead")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", renamedRoleName),
					testAccCheckResourceIDChanged(resourceReference, &originalID),
				),
			},
		},
	})
}

// TestAccDatabaseRoleDescriptionLifecycle covers description being set, changed and
// dropped. description is Optional+Computed, so removing it from config leaves the last
// known value in state rather than clearing it.
func TestAccDatabaseRoleDescriptionLifecycle(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_desc_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "first description", accessCollectionLevel("dataRead")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "description", "first description"),
				),
			},
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "second description", accessCollectionLevel("dataRead")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "description", "second description"),
				),
			},
			{
				Config: testAccDatabaseRoleConfigAccess(resourceName, resourceName, "", accessCollectionLevel("dataRead")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDatabaseRoleResource(t, resourceReference),
					resource.TestCheckResourceAttrSet(resourceReference, "description"),
				),
			},
		},
	})
}

// TestAccDatabaseRoleDeletedOutOfBand deletes the role behind Terraform's back and
// verifies the next refresh drops it from state instead of erroring, which is the
// not-found branch of Read.
func TestAccDatabaseRoleDeletedOutOfBand(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_db_role_oob_")
	resourceReference := "couchbase-capella_database_role." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testAccDatabaseRoleConfigAccess(resourceName, resourceName, "", accessCollectionLevel("dataRead")),
				Check:              testAccDatabaseRoleDeleteOOB(resourceReference),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// --- Config builders ---

// testAccDatabaseRoleConfigAccess renders a database role whose access list is the
// caller-supplied HCL, so each level of the Capella RBAC template can be exercised
// independently. roleName is separate from resourceName so rename tests can vary the
// role name without moving the resource address.
func testAccDatabaseRoleConfigAccess(resourceName, roleName, description, accessHCL string) string {
	descBlock := ""
	if description != "" {
		descBlock = fmt.Sprintf("description     = %q", description)
	}
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_role" %[2]q {
  organization_id = %[3]q
  project_id      = %[4]q
  cluster_id      = %[5]q
  name            = %[6]q
  %[7]s
  access = %[8]s
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, roleName, descBlock, accessHCL)
}

// accessCollectionLevel grants privilege at bucket -> scope -> collection, the deepest
// level the RBAC template allows for privileges such as dataRead.
func accessCollectionLevel(privilege string) string {
	return fmt.Sprintf(`[
    {
      privileges = [%q]
      resources = {
        buckets = [
          {
            name = %q
            scopes = [
              {
                name        = %q
                collections = [%q]
              },
            ]
          },
        ]
      }
    },
  ]`, privilege, globalBucketName, globalScopeName, globalCollectionName)
}

// accessScopeLevel grants privilege at bucket -> scope with collections omitted, the
// only valid shape for scope-level privileges such as queryManage.
func accessScopeLevel(privilege string) string {
	return fmt.Sprintf(`[
    {
      privileges = [%q]
      resources = {
        buckets = [
          {
            name = %q
            scopes = [
              {
                name = %q
              },
            ]
          },
        ]
      }
    },
  ]`, privilege, globalBucketName, globalScopeName)
}

// accessBucketLevel grants privilege at bucket level with scopes omitted, the only
// valid shape for bucket-level privileges such as viewsReader.
func accessBucketLevel(privilege string) string {
	return fmt.Sprintf(`[
    {
      privileges = [%q]
      resources = {
        buckets = [
          {
            name = %q
          },
        ]
      }
    },
  ]`, privilege, globalBucketName)
}

// accessAllScopedLevels combines the three scoped levels in one access list. The
// entries carry distinct resources, so the API keeps them as three separate entries.
func accessAllScopedLevels() string {
	return fmt.Sprintf(`[
    {
      privileges = ["dataRead"]
      resources = {
        buckets = [
          {
            name = %[1]q
            scopes = [
              {
                name        = %[2]q
                collections = [%[3]q]
              },
            ]
          },
        ]
      }
    },
    {
      privileges = ["queryManage"]
      resources = {
        buckets = [
          {
            name = %[1]q
            scopes = [
              {
                name = %[2]q
              },
            ]
          },
        ]
      }
    },
    {
      privileges = ["viewsReader"]
      resources = {
        buckets = [
          {
            name = %[1]q
          },
        ]
      }
    },
  ]`, globalBucketName, globalScopeName, globalCollectionName)
}
