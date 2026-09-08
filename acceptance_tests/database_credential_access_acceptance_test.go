package acceptance_tests

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// This file covers the shape and semantics of the database credential `access`
// block: `access`, `access.resources.buckets`, `...buckets.scopes` and
// `...scopes.collections`, all of which are Set-typed attributes.
//
// They were briefly Lists. AV-139841 (#749) changed all four from Set to List,
// that shipped in v1.11.0, and AV-142331 (#772) reverted them to Set in v1.11.1.
// The tests here assert Set semantics, matching the schema as it stands.
//
// How a Set is guarded here, and why not by asserting element positions: cty
// stores set elements in a canonical order derived from a hash, which is not a
// documented contract and can shift between cty versions. Asserting "element 0
// is X" would therefore be pinning an implementation detail. Two properties are
// asserted instead, both of which follow from Set semantics directly:
//
//   - Values are checked by count and membership, never by position.
//   - Reordering elements in the configuration is not a change. Each ordering
//     test ends with a PlanOnly step that re-plans the same elements in a
//     different order and requires an empty plan. Under a List that step plans a
//     diff and fails, so it also catches an unintended return to List.
//
// Note that the pre-existing test in database_credential_acceptance_test.go
// cannot tell the two types apart: it configures a single access element with a
// single privilege, and terraform-plugin-testing flattens lists and sets alike
// by numeric index (internal/configs/hcl2shim/flatmap.go), so
// "access.0.privileges.0" holds either way.

// ─────────────────────────────────────────────────────────────────────────────
// Set semantics: unordered and deduplicated
// ─────────────────────────────────────────────────────────────────────────────

// TestAccDatabaseCredentialAccessBlocksAreUnordered guards the top-level
// `access` set: both grants survive, and swapping their order in the
// configuration is not a change.
func TestAccDatabaseCredentialAccessBlocksAreUnordered(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_dbc_order_")
	resourceReference := "couchbase-capella_database_credential." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialAccessBlocksConfig(resourceName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "password"),

					resource.TestCheckResourceAttr(resourceReference, "access.#", "2"),
					// Located by content rather than position: the bucket-scoped
					// writer grant, and the bare reader grant on all buckets.
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, "access.*", map[string]string{
						"privileges.#":             "1",
						"privileges.0":             "data_writer",
						"resources.buckets.#":      "1",
						"resources.buckets.0.name": globalBucketName,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, "access.*", map[string]string{
						"privileges.#": "1",
						"privileges.0": "data_reader",
					}),
				),
			},
			// The same two grants, written in the opposite order. A Set does not
			// distinguish them, so the plan must be empty; a List would diff.
			{
				Config:   testAccDatabaseCredentialAccessBlocksConfig(resourceName, false),
				PlanOnly: true,
			},
			// ImportState testing. ImportStateVerify is off because Read()
			// rebuilds access from prior state rather than the GET response, so
			// an imported credential carries no access at all - AV-142491, and
			// TestAccDatabaseCredentialImportPreservesAccess below.
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateDatabaseCredentialImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

// TestAccDatabaseCredentialBucketsAreUnordered guards
// `access.resources.buckets` using the two buckets TestMain provisions.
func TestAccDatabaseCredentialBucketsAreUnordered(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_dbc_buckets_")
	resourceReference := "couchbase-capella_database_credential." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialBucketsConfig(resourceName, globalBucketName, globalMetadataBucketName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_reader"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, "access.0.resources.buckets.*",
						map[string]string{"name": globalBucketName}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, "access.0.resources.buckets.*",
						map[string]string{"name": globalMetadataBucketName}),
				),
			},
			{
				Config:   testAccDatabaseCredentialBucketsConfig(resourceName, globalMetadataBucketName, globalBucketName),
				PlanOnly: true,
			},
		},
	})
}

// TestAccDatabaseCredentialScopesAndCollectionsAreUnordered guards the two
// deepest levels. It provisions two real scopes and three real collections,
// because the V4 API validates the keyspaces a grant names.
func TestAccDatabaseCredentialScopesAndCollectionsAreUnordered(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_dbc_deep_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	scopeA, scopeZ := orderedFixtureNames("tf_acc_dbc_scope_")
	colA, colM, colZ := orderedFixtureNamesTriple("tf_acc_dbc_col_")
	created := []string{colA, colM, colZ}

	scopesPath := "access.0.resources.buckets.0.scopes"
	collectionsPath := scopesPath + ".0.collections"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialDeepAccessConfig(
					resourceName, scopeZ, scopeA, created, []string{colZ, colA, colM},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.name", globalBucketName),

					resource.TestCheckResourceAttr(resourceReference, scopesPath+".#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, scopesPath+".*",
						map[string]string{"name": scopeZ, "collections.#": "3"}),
					// The second scope grants scope-level access with no collections.
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, scopesPath+".*",
						map[string]string{"name": scopeA}),

					resource.TestCheckTypeSetElemAttr(resourceReference, collectionsPath+".*", colZ),
					resource.TestCheckTypeSetElemAttr(resourceReference, collectionsPath+".*", colA),
					resource.TestCheckTypeSetElemAttr(resourceReference, collectionsPath+".*", colM),
				),
			},
			// Same scopes and collections, both reordered. Neither is a change.
			{
				Config: testAccDatabaseCredentialDeepAccessConfig(
					resourceName, scopeZ, scopeA, created, []string{colA, colM, colZ},
				),
				PlanOnly: true,
			},
		},
	})
}

// TestAccDatabaseCredentialDuplicateAccessBlocksAreDeduplicated pins the other
// half of Set semantics. Two identical grants collapse into one, silently -
// worth a test because it is a real footgun: a configuration generated by a
// loop can lose a grant without any diagnostic.
func TestAccDatabaseCredentialDuplicateAccessBlocksAreDeduplicated(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_dbc_dup_")
	resourceReference := "couchbase-capella_database_credential." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialDuplicateAccessConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_reader"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.name", globalBucketName),
				),
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Absent versus empty
// ─────────────────────────────────────────────────────────────────────────────

// TestAccDatabaseCredentialEmptyNestedSets exercises the empty-slice branches of
// createAccess and mapAccess at each nesting depth. Empty and absent are
// distinct values, and collapsing one into the other surfaces as "Provider
// produced inconsistent result after apply".
func TestAccDatabaseCredentialEmptyNestedSets(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_dbc_empty_")
	resourceReference := "couchbase-capella_database_credential." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// buckets = [] - the shape createAccess already synthesises as a
			// workaround for the PUT nil-pointer bug noted in AV-63388.
			{
				Config: testAccDatabaseCredentialEmptyBucketsConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.#", "0"),
				),
			},
			// scopes = [] - bucket-level grant, empty scope set.
			{
				Config: testAccDatabaseCredentialEmptyScopesConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.name", globalBucketName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.#", "0"),
				),
			},
			// collections = [] - scope-level grant, empty collection set.
			{
				Config: testAccDatabaseCredentialEmptyCollectionsConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.name", globalScopeName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.collections.#", "0"),
				),
			},
		},
	})
}

// TestAccDatabaseCredentialAccessResourcesWithoutBuckets covers a `resources`
// object supplied with no `buckets` key - the defect this branch fixes
// (AV-142487). mapAccess used to copy Resources only when Resources.Buckets was
// non-nil, so the applied state came back with resources = null while the plan
// held resources = {buckets = null}, and Terraform aborted the apply.
// createAccess had the mirror-image gap: its empty-buckets workaround sat in the
// Resources == nil branch only, so this configuration sent no resources to the
// V4 API at all.
//
// The two steps cross the fixed boundary in both directions: buckets absent,
// then buckets present. Nothing else about the credential changes between them.
func TestAccDatabaseCredentialAccessResourcesWithoutBuckets(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_dbc_nobuckets_")
	resourceReference := "couchbase-capella_database_credential." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialResourcesWithoutBucketsConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_reader"),
					// buckets must stay absent rather than being coerced to an
					// empty set, which would be a different value from the plan's.
					resource.TestCheckNoResourceAttr(resourceReference, "access.0.resources.buckets.#"),
				),
			},
			// nil -> non-nil buckets on the same resources object.
			{
				Config: testAccDatabaseCredentialEmptyScopesConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_reader"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.name", globalBucketName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.#", "0"),
				),
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Grant narrowing and privileges
// ─────────────────────────────────────────────────────────────────────────────

// TestAccDatabaseCredentialAccessNarrowingUpdate walks a grant from all buckets
// to one bucket, to scopes, to collections, and back out again. Each step
// asserts the whole access block rather than just the field that moved, so a
// change at one depth cannot silently reset another. This is the only coverage
// of an update crossing the nil-checking branches in createAccess and mapAccess.
func TestAccDatabaseCredentialAccessNarrowingUpdate(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_dbc_narrow_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	scopeA, scopeZ := orderedFixtureNames("tf_acc_dbc_nscope_")
	colA, colM, colZ := orderedFixtureNamesTriple("tf_acc_dbc_ncol_")
	created := []string{colA, colM, colZ}
	bucketsPath := "access.0.resources.buckets"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Widest: no resources at all, meaning all buckets.
			{
				Config: testAccDatabaseCredentialBareAccessConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_reader"),
					resource.TestCheckNoResourceAttr(resourceReference, bucketsPath+".#"),
				),
			},
			// Narrow to one bucket.
			{
				Config: testAccDatabaseCredentialEmptyScopesConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_reader"),
					resource.TestCheckResourceAttr(resourceReference, bucketsPath+".#", "1"),
					resource.TestCheckResourceAttr(resourceReference, bucketsPath+".0.name", globalBucketName),
					resource.TestCheckResourceAttr(resourceReference, bucketsPath+".0.scopes.#", "0"),
				),
			},
			// Narrow to two scopes, three collections on the first.
			{
				Config: testAccDatabaseCredentialDeepAccessConfig(
					resourceName, scopeZ, scopeA, created, []string{colZ, colA, colM},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_reader"),
					resource.TestCheckResourceAttr(resourceReference, bucketsPath+".#", "1"),
					resource.TestCheckResourceAttr(resourceReference, bucketsPath+".0.scopes.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, bucketsPath+".0.scopes.*",
						map[string]string{"name": scopeZ, "collections.#": "3"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, bucketsPath+".0.scopes.*",
						map[string]string{"name": scopeA}),
				),
			},
			// Widen all the way back out. This is the step that regresses if an
			// update ever stops clearing the deeper levels.
			{
				Config: testAccDatabaseCredentialBareAccessConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_reader"),
					resource.TestCheckNoResourceAttr(resourceReference, bucketsPath+".#"),
				),
			},
		},
	})
}

// TestAccDatabaseCredentialMultiplePrivileges covers `privileges` carrying both
// values. It is a Set, so membership is the only safe assertion - as it is for
// every attribute in this file.
//
// Worth noting while reading the schema: the database_credentials data source
// declares `access`, `buckets`, `scopes`, `collections` and `privileges` as
// Lists, where the resource declares all five as Sets. The two have disagreed
// since before AV-139841 and still disagree after AV-142331, and the generated
// docs contradict each other accordingly.
func TestAccDatabaseCredentialMultiplePrivileges(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_dbc_privs_")
	resourceReference := "couchbase-capella_database_credential." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialBothPrivilegesConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "data_reader"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "data_writer"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.name", globalBucketName),
				),
			},
			// Dropping one privilege must leave exactly the other.
			{
				Config: testAccDatabaseCredentialBucketsConfig(resourceName, globalBucketName, globalMetadataBucketName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "data_reader"),
				),
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Pre-existing defect. Skipped so the suite stays green; it asserts the
// behaviour that should hold once the referenced bug is fixed.
// ─────────────────────────────────────────────────────────────────────────────

// TestAccDatabaseCredentialImportPreservesAccess covers importing a credential.
// ImportState only passes the composite ID through, and Read() then rebuilds
// access from prior state via mapAccess rather than from the GET response, so an
// imported credential lands with access = [] and every nested value missing.
// The GET response does carry access - the database_credentials data source maps
// it straight from there - so Read should too.
func TestAccDatabaseCredentialImportPreservesAccess(t *testing.T) {
	t.Skip("AV-142491: Read() rebuilds access from prior state instead of the GET response, so an imported credential has access = [] and ImportStateVerify reports the whole block as missing; unskip once Read maps access from the API. Confirmed failing on a live cluster: all 18 access.* attributes are absent after import")

	resourceName := randomStringWithPrefix("tf_acc_dbc_import_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	scopeA, scopeZ := orderedFixtureNames("tf_acc_dbc_iscope_")
	colA, colM, colZ := orderedFixtureNamesTriple("tf_acc_dbc_icol_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialDeepAccessConfig(
					resourceName, scopeZ, scopeA, []string{colA, colM, colZ}, []string{colZ, colA, colM},
				),
				Check: resource.TestCheckResourceAttr(
					resourceReference, "access.0.resources.buckets.0.scopes.0.collections.#", "3",
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateDatabaseCredentialImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
				// password is write-only: the GET response never returns it.
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Data source
// ─────────────────────────────────────────────────────────────────────────────

// TestAccDatasourceDatabaseCredentialsNestedAccess covers the access block on
// the database_credentials data source, which the existing data source test does
// not touch at all - it only checks names and IDs.
//
// The data source maps access straight from the GET response
// (internal/datasources/database_credentials.go), unlike the resource, which
// rebuilds it from prior state. Assertions here are on cardinality and on
// single-element values, and a diagnostic logs the rest of what the API returned
// for AV-142491.
func TestAccDatasourceDatabaseCredentialsNestedAccess(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_dbc_ds_access_")
	dsName := randomStringWithPrefix("tf_acc_dbc_ds_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	dsReference := "data.couchbase-capella_database_credentials." + dsName
	scopeA, scopeZ := orderedFixtureNames("tf_acc_dbc_dsscope_")
	colA, colM, colZ := orderedFixtureNamesTriple("tf_acc_dbc_dscol_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseCredentialDeepAccessConfig(
					resourceName, scopeZ, scopeA, []string{colA, colM, colZ}, []string{colZ, colA, colM},
				) + testAccDatabaseCredentialsDataSourceOnlyConfig(dsName, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference,
						"access.0.resources.buckets.0.scopes.0.collections.#", "3"),

					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
					// Locate our credential in the list and assert the access
					// block came back with the right shape.
					resource.TestCheckTypeSetElemNestedAttrs(dsReference, "data.*", map[string]string{
						"name":                                  resourceName,
						"cluster_id":                            globalClusterId,
						"access.#":                              "1",
						"access.0.resources.buckets.#":          "1",
						"access.0.resources.buckets.0.name":     globalBucketName,
						"access.0.resources.buckets.0.scopes.#": "2",
					}),
					logDatasourceAccessShape(t, dsReference, resourceName, scopeZ, []string{colZ, colA, colM}),
				),
			},
		},
	})
}

// logDatasourceAccessShape records what the GET response returned for a
// credential, without asserting it.
//
// This is a diagnostic for AV-142491. Making Read() map access from the API is
// what would let an import carry its access block, and the obstacle named in the
// provider's own comment is that "the formats passed in terraform files and the
// GET response are different". Ordering is not the problem now that these
// attributes are Sets again, so what matters is whether the returned membership
// matches what was configured - in particular how an all-buckets grant comes
// back, which is the shape most likely to differ.
//
// It deliberately does not fail on a mismatch: a difference is a finding about
// the API, not a provider regression. Run with -v and read the log lines.
func logDatasourceAccessShape(
	t *testing.T,
	dsReference, credName, scopeWithCollections string,
	configuredCollections []string,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dsReference]
		if !ok {
			return fmt.Errorf("data source %s not found in state", dsReference)
		}
		attrs := rs.Primary.Attributes

		dataCount, err := strconv.Atoi(attrs["data.#"])
		if err != nil {
			return fmt.Errorf("reading %s data.#: %w", dsReference, err)
		}

		for i := range dataCount {
			if attrs[fmt.Sprintf("data.%d.name", i)] != credName {
				continue
			}

			scopesBase := fmt.Sprintf("data.%d.access.0.resources.buckets.0.scopes", i)
			scopeCount, err := strconv.Atoi(attrs[scopesBase+".#"])
			if err != nil {
				return fmt.Errorf("reading %s.#: %w", scopesBase, err)
			}

			gotScopes := make([]string, scopeCount)
			collectionsBase := ""
			for j := range scopeCount {
				gotScopes[j] = attrs[fmt.Sprintf("%s.%d.name", scopesBase, j)]
				if gotScopes[j] == scopeWithCollections {
					collectionsBase = fmt.Sprintf("%s.%d.collections", scopesBase, j)
				}
			}
			t.Logf("AV-142491 probe: scopes returned by the GET response: %v", gotScopes)

			if collectionsBase == "" {
				t.Logf("AV-142491 probe: scope %q absent from the GET response", scopeWithCollections)
				return nil
			}

			collectionCount, err := strconv.Atoi(attrs[collectionsBase+".#"])
			if err != nil {
				return fmt.Errorf("reading %s.#: %w", collectionsBase, err)
			}
			gotCollections := make([]string, collectionCount)
			for j := range collectionCount {
				gotCollections[j] = attrs[fmt.Sprintf("%s.%d", collectionsBase, j)]
			}

			want := slices.Sorted(slices.Values(configuredCollections))
			got := slices.Sorted(slices.Values(gotCollections))
			t.Logf("AV-142491 probe: collections configured %v, returned %v (same membership=%t)",
				configuredCollections, gotCollections, slices.Equal(want, got))

			return nil
		}

		return fmt.Errorf("credential %s not found among the %d entries of %s", credName, dataCount, dsReference)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Fixture name helpers
// ─────────────────────────────────────────────────────────────────────────────

// orderedFixtureNames returns two names sharing a random suffix, distinguished
// by a fixed a/z infix. randomStringWithPrefix alone gives no way to tell which
// generated name is which, which matters when an assertion names one of them.
func orderedFixtureNames(prefix string) (a, z string) {
	suffix := randomString()
	return prefix + "a_" + suffix, prefix + "z_" + suffix
}

// orderedFixtureNamesTriple is orderedFixtureNames for three names.
func orderedFixtureNamesTriple(prefix string) (a, m, z string) {
	suffix := randomString()
	return prefix + "a_" + suffix, prefix + "m_" + suffix, prefix + "z_" + suffix
}

func generateDatabaseCredentialImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf(
			"id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Config helpers
// ─────────────────────────────────────────────────────────────────────────────

// dbCredScopeBlock and dbCredCollectionBlock provision real keyspaces on the
// shared fixture bucket. The credential configs reference the created names as
// literals and depend on the resources explicitly, so every name in an access
// block is known at plan time and the ordering assertions stay unambiguous.
func dbCredScopeBlock(scopeName string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_scope" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_id       = "%[5]s"
  scope_name      = "%[1]s"
}
`, scopeName, globalOrgId, globalProjectId, globalClusterId, globalBucketId)
}

func dbCredCollectionBlock(collectionName, scopeName string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_collection" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_id       = "%[5]s"
  scope_name      = couchbase-capella_scope.%[6]s.scope_name
  collection_name = "%[1]s"
}
`, collectionName, globalOrgId, globalProjectId, globalClusterId, globalBucketId, scopeName)
}

// testAccDatabaseCredentialDeepAccessConfig builds a credential whose single
// access block reaches collection depth. scopeWithCollections is listed first
// and scopeBare second; createdCollections are all provisioned in
// scopeWithCollections, and orderedCollections is the order they are granted in.
func testAccDatabaseCredentialDeepAccessConfig(
	credName, scopeWithCollections, scopeBare string,
	createdCollections, orderedCollections []string,
) string {
	var b strings.Builder
	b.WriteString(globalProviderBlock)
	b.WriteString(dbCredScopeBlock(scopeWithCollections))
	b.WriteString(dbCredScopeBlock(scopeBare))
	for _, c := range createdCollections {
		b.WriteString(dbCredCollectionBlock(c, scopeWithCollections))
	}

	quoted := make([]string, len(orderedCollections))
	for i, c := range orderedCollections {
		quoted[i] = fmt.Sprintf("%q", c)
	}
	// Both scopes are named as literals inside the access block, so nothing in
	// the configuration orders the credential after them. scopeWithCollections is
	// covered transitively - the collections reference it - but scopeBare has no
	// other reference at all, and the V4 API validates the keyspaces a grant
	// names, so without this the credential can be created too early and fail
	// sporadically. Listing both scopes keeps the ordering explicit rather than
	// relying on that transitive edge.
	deps := []string{
		"couchbase-capella_scope." + scopeWithCollections,
		"couchbase-capella_scope." + scopeBare,
	}
	for _, c := range createdCollections {
		deps = append(deps, "couchbase-capella_collection."+c)
	}

	fmt.Fprintf(&b, `
resource "couchbase-capella_database_credential" "%[1]s" {
  name            = "%[1]s"
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"

  access = [
    {
      privileges = ["data_reader"]
      resources = {
        buckets = [
          {
            name = "%[5]s"
            scopes = [
              {
                name        = "%[6]s"
                collections = [%[7]s]
              },
              {
                name = "%[8]s"
              },
            ]
          },
        ]
      }
    },
  ]

  depends_on = [%[9]s]
}
`, credName, globalOrgId, globalProjectId, globalClusterId, globalBucketName,
		scopeWithCollections, strings.Join(quoted, ", "), scopeBare, strings.Join(deps, ", "))

	return b.String()
}

// testAccDatabaseCredentialAccessBlocksConfig mirrors the two-block shape
// shipped in examples/database_credential/terraform.template.tfvars. writerFirst
// selects which of the two blocks is written first, so a step can re-plan the
// same two grants in the opposite order without changing anything else.
func testAccDatabaseCredentialAccessBlocksConfig(credName string, writerFirst bool) string {
	writer := `    {
      privileges = ["data_writer"]
      resources = {
        buckets = [
          {
            name = "` + globalBucketName + `"
          },
        ]
      }
    },
`
	reader := `    {
      privileges = ["data_reader"]
    },
`
	blocks := writer + reader
	if !writerFirst {
		blocks = reader + writer
	}

	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" "%[2]s" {
  name            = "%[2]s"
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"

  access = [
%[6]s  ]
}
`, globalProviderBlock, credName, globalOrgId, globalProjectId, globalClusterId, blocks)
}

func testAccDatabaseCredentialBucketsConfig(credName, firstBucket, secondBucket string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" "%[2]s" {
  name            = "%[2]s"
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"

  access = [
    {
      privileges = ["data_reader"]
      resources = {
        buckets = [
          {
            name = "%[6]s"
          },
          {
            name = "%[7]s"
          },
        ]
      }
    },
  ]
}
`, globalProviderBlock, credName, globalOrgId, globalProjectId, globalClusterId, firstBucket, secondBucket)
}

// testAccDatabaseCredentialBareAccessConfig grants on all buckets by omitting
// the resources object entirely.
func testAccDatabaseCredentialBareAccessConfig(credName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" "%[2]s" {
  name            = "%[2]s"
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"

  access = [
    {
      privileges = ["data_reader"]
    },
  ]
}
`, globalProviderBlock, credName, globalOrgId, globalProjectId, globalClusterId)
}

func testAccDatabaseCredentialEmptyBucketsConfig(credName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" "%[2]s" {
  name            = "%[2]s"
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"

  access = [
    {
      privileges = ["data_reader"]
      resources = {
        buckets = []
      }
    },
  ]
}
`, globalProviderBlock, credName, globalOrgId, globalProjectId, globalClusterId)
}

func testAccDatabaseCredentialEmptyScopesConfig(credName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" "%[2]s" {
  name            = "%[2]s"
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"

  access = [
    {
      privileges = ["data_reader"]
      resources = {
        buckets = [
          {
            name   = "%[6]s"
            scopes = []
          },
        ]
      }
    },
  ]
}
`, globalProviderBlock, credName, globalOrgId, globalProjectId, globalClusterId, globalBucketName)
}

func testAccDatabaseCredentialEmptyCollectionsConfig(credName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" "%[2]s" {
  name            = "%[2]s"
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"

  access = [
    {
      privileges = ["data_reader"]
      resources = {
        buckets = [
          {
            name = "%[6]s"
            scopes = [
              {
                name        = "%[7]s"
                collections = []
              },
            ]
          },
        ]
      }
    },
  ]
}
`, globalProviderBlock, credName, globalOrgId, globalProjectId, globalClusterId, globalBucketName, globalScopeName)
}

func testAccDatabaseCredentialDuplicateAccessConfig(credName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" "%[2]s" {
  name            = "%[2]s"
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"

  access = [
    {
      privileges = ["data_reader"]
      resources = {
        buckets = [
          {
            name = "%[6]s"
          },
        ]
      }
    },
    {
      privileges = ["data_reader"]
      resources = {
        buckets = [
          {
            name = "%[6]s"
          },
        ]
      }
    },
  ]
}
`, globalProviderBlock, credName, globalOrgId, globalProjectId, globalClusterId, globalBucketName)
}

func testAccDatabaseCredentialBothPrivilegesConfig(credName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" "%[2]s" {
  name            = "%[2]s"
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"

  access = [
    {
      privileges = ["data_writer", "data_reader"]
      resources = {
        buckets = [
          {
            name = "%[6]s"
          },
        ]
      }
    },
  ]
}
`, globalProviderBlock, credName, globalOrgId, globalProjectId, globalClusterId, globalBucketName)
}

func testAccDatabaseCredentialResourcesWithoutBucketsConfig(credName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_database_credential" "%[2]s" {
  name            = "%[2]s"
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"

  access = [
    {
      privileges = ["data_reader"]
      resources  = {}
    },
  ]
}
`, globalProviderBlock, credName, globalOrgId, globalProjectId, globalClusterId)
}

// testAccDatabaseCredentialsDataSourceOnlyConfig emits just the data source, to
// be appended to a config that already carries a provider block and the
// credential it depends on.
func testAccDatabaseCredentialsDataSourceOnlyConfig(dsName, credName string) string {
	return fmt.Sprintf(`
data "couchbase-capella_database_credentials" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"

  depends_on = [couchbase-capella_database_credential.%[5]s]
}
`, dsName, globalOrgId, globalProjectId, globalClusterId, credName)
}
