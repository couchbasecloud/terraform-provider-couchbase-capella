package acceptance_tests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// The tests in this file all cover the cluster `zones` attribute as AV-147905 reshaped it.
//
// Before that change `zones` was Optional with RequiresReplace, and Read echoed back
// whatever Create had stored (internal/resources/cluster.go, `if state.Zones != nil`),
// because NewCluster never assigned Zones at all. The attribute therefore could not drift,
// the `couchbase-capella_cluster` data source always reported it as null, and a cluster
// imported with zones in config was forced to replace.
//
// It is now Optional+Computed with RequiresReplaceIfConfigured, sourced from the API, which
// fixes import and the data source but creates three behaviours worth guarding:
//
//   - the value is now read back from Capella, so it could in principle disagree with the
//     configuration and force a destroy/recreate; the zone-name route to that was tested
//     live and is closed, leaving only plan-time validation to tighten (AV-149540);
//   - a computed attribute that stays null is re-marked unknown on every update plan, which
//     is permanent for GCP and Azure because Capella never reports zones for them
//     (AV-149539);
//   - the data source attribute changed from a list to a set, so configurations indexing it
//     stop validating (AV-149541).
//
// Cost note: Capella has no shared GCP or Azure cluster in this suite, so the non-AWS tests
// provision one each. Only the plan-time validation tests belong in sanity.list; the rest
// provision real clusters and belong in the full acceptance run.

// TestAccClusterZonesReadBackFromDatasource is the regression guard for the whole feature:
// it proves the provider maps zones out of the API response rather than only echoing config.
//
// It reads the shared cluster, which acceptance_tests/cluster.go creates as a multi AZ AWS
// cluster in us-east-1 with no zones in the request, so Capella chose them. Reverting the
// `Zones: NewZonesSet(cluster.Zones)` mapping in internal/schema/cluster.go leaves every
// other test in the suite green and fails this one.
//
// No cluster is provisioned - it reads whatever cluster the run is already using.
func TestAccClusterZonesReadBackFromDatasource(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_cluster_zones_ds_")
	dsRef := "data.couchbase-capella_cluster." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesDatasourceConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsRef, "id", globalClusterId),
					// Adaptive, because TF_VAR_cluster_id may point at a cluster this run
					// did not create and whose provider therefore is not known here.
					testAccCheckClusterZonesMatchPlatform(dsRef),
				),
			},
		},
	})
}

// TestAccClusterZonesMultiAZNotConfigured covers the same mapping as the data source test
// above, but through the resource, and additionally proves that leaving zones out does not
// produce a perpetual diff now that the attribute is computed: the step's own post-apply
// plan must come back empty.
//
// Provisions a cluster.
func TestAccClusterZonesMultiAZNotConfigured(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfig(resourceName, generateRandomCIDR(), "aws", "us-east-1", "multi", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					// Capella assigns one zone per node group placement; a multi AZ AWS
					// cluster is spread over exactly three.
					resource.TestCheckResourceAttr(resourceReference, "zones.#", "3"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateClusterImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccClusterZonesExplicitNull covers `zones = null` written out in the configuration,
// as distinct from leaving the attribute out altogether.
//
// This is the shape the repository's own example module produces:
// examples/cluster/create_cluster.tf passes `zones = var.cluster.zones` from an
// `optional(list(string))` variable, so an unset variable reaches the provider as an
// explicit null rather than as an absent attribute.
//
// HCL treats an argument set to null as equivalent to omitting it, so the two should be
// interchangeable - and the second step proves exactly that by switching from one to the
// other and requiring a no-op. If they ever diverge, a practitioner toggling an optional
// variable would see a spurious plan.
//
// Provisions a cluster.
func TestAccClusterZonesExplicitNull(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	// One CIDR across both steps: cloud_provider is RequiresReplace.
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				// zones = null - Capella chooses, and the chosen zones are read back.
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "multi", "null"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "zones.#", "3"),
				),
			},
			{
				// The attribute dropped entirely: the same value, so nothing to do.
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "multi", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.TestCheckResourceAttr(resourceReference, "zones.#", "3"),
			},
		},
	})
}

// TestAccClusterZonesSingleAZImportRoundTrip closes the gap left by
// TestAccClusterResourceWithConfigurationTypeFieldAdded, whose import step sets
// ImportStateVerify: false and so cannot prove zones survive an import.
//
// On main this round trip was broken: import left zones null in state, so the next plan
// against a config that set zones forced a replacement.
//
// Provisions a cluster.
func TestAccClusterZonesSingleAZImportRoundTrip(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	// One CIDR for every step: cloud_provider is RequiresReplace, so re-generating it per
	// step would plan a replacement and the no-op assertion below would fail for the wrong
	// reason.
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "single", `["use1-az1"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "single"),
					resource.TestCheckResourceAttr(resourceReference, "zones.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "zones.*", "use1-az1"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateClusterImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
				// Without this the imported state is discarded at the end of this step and
				// the plan check below would run against the state left by the create,
				// proving nothing about import. ImportStateVerify only compares attributes;
				// it does not persist them.
				ImportStatePersist: true,
			},
			{
				// The plan against the IMPORTED state must be clean rather than a replacement.
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "single", `["use1-az1"]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// TestAccClusterZonesReorderIsNotAChange is the set-versus-list guard.
//
// The zones are configured in an order that is deliberately NOT sorted, then the second step
// reorders them into sorted order. Under a Set that is the same value and nothing happens;
// under a List it is a different value and, with RequiresReplaceIfConfigured, the cluster is
// destroyed and recreated. Do not "tidy" these into alphabetical order - the test proves
// nothing if the configured order already matches the canonical one.
//
// Note this cannot be a PlanOnly step: PlanOnly skips ConfigPlanChecks.PreApply, so the
// assertion would silently guard nothing.
//
// Provisions a cluster.
func TestAccClusterZonesReorderIsNotAChange(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "multi", `["use1-az4", "use1-az2", "use1-az1"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "zones.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "zones.*", "use1-az1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "zones.*", "use1-az2"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "zones.*", "use1-az4"),
				),
			},
			{
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "multi", `["use1-az1", "use1-az2", "use1-az4"]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// TestAccClusterZonesRemovedFromConfig covers a transition the schema change quietly altered.
//
// Under the old Optional+RequiresReplace schema, deleting zones from a configuration that
// had them planned a replacement. Now the attribute is Optional+Computed, so an absent
// configuration keeps whatever Capella reports and nothing happens. Clearing a field the
// practitioner no longer wants is a shape people reach naturally, and it must not destroy a
// cluster.
//
// Provisions a cluster.
func TestAccClusterZonesRemovedFromConfig(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "multi", `["use1-az1", "use1-az2", "use1-az4"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "zones.#", "3"),
				),
			},
			{
				// zones omitted entirely: no replacement, and the assigned zones stay in state.
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "multi", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "zones.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "zones.*", "use1-az1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "zones.*", "use1-az2"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "zones.*", "use1-az4"),
				),
			},
		},
	})
}

// TestAccClusterZonesUnchangedByUnrelatedUpdate guards the behaviour AV-149539's fix must
// preserve, and is the reason that ticket says not to simply delete useStateForUnknown.
//
// On any plan where something else changed, the framework marks every computed attribute
// whose config is null as unknown. For a cluster whose zones Capella chose, useStateForUnknown
// copies the stored value back; without it, an unrelated edit renders as
// `zones = [...] -> (known after apply)` on every single update.
//
// Provisions a cluster.
func TestAccClusterZonesUnchangedByUnrelatedUpdate(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfigWithDescription(resourceName, cidr, "aws", "us-east-1", "multi", "", "before"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "zones.#", "3"),
				),
			},
			{
				// Only the description changes. zones must stay known in the plan.
				Config: testAccClusterZonesConfigWithDescription(resourceName, cidr, "aws", "us-east-1", "multi", "", "after"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(resourceReference, tfjsonpath.New("zones"), knownvalue.SetSizeExact(3)),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "description", "after"),
					resource.TestCheckResourceAttr(resourceReference, "zones.#", "3"),
				),
			},
		},
	})
}

// TestAccClusterZonesAbsentForGCP proves zones stays null for a provider Capella does not
// report zone IDs for, rather than becoming an empty set. Null and empty are different
// values to Terraform, and NewZonesSet deliberately returns a null set for both an absent
// and an empty API field.
//
// Provisions a GCP cluster.
func TestAccClusterZonesAbsentForGCP(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfig(resourceName, generateRandomCIDR(), "gcp", "us-east1", "multi", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "gcp"),
					// Null, not empty: a null set has no `zones.#` entry in the flatmap at
					// all, whereas an empty one would be "0". Asserting absence is correct
					// here precisely because null is the intended value.
					resource.TestCheckNoResourceAttr(resourceReference, "zones.#"),
				),
			},
		},
	})
}

// TestAccClusterZonesNoUnknownOnUpdateGCP_AV_149539 asserts the behaviour AV-149539 asks for,
// not the behaviour today, so it fails until that ticket is fixed.
//
// Capella never reports zone IDs for GCP or Azure, so the attribute is null in state forever
// and the framework re-marks it unknown on every update plan, giving a permanent
// `+ zones = (known after apply)` line that no apply ever settles.
//
// Provisions a GCP cluster when unskipped.
func TestAccClusterZonesNoUnknownOnUpdateGCP_AV_149539(t *testing.T) {
	t.Skip("AV-149539: zones is re-marked unknown on every update plan for GCP and Azure " +
		"clusters because UseStateForUnknown bails on a null prior state; unskip once the " +
		"plan modifier carries a null state value across")

	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfigWithDescription(resourceName, cidr, "gcp", "us-east1", "multi", "", "before"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckNoResourceAttr(resourceReference, "zones.#"),
				),
			},
			{
				Config: testAccClusterZonesConfigWithDescription(resourceName, cidr, "gcp", "us-east1", "multi", "", "after"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						// Today this is unknown. It should be a known null: the provider
						// knows Capella will not report zones for a GCP cluster.
						plancheck.ExpectKnownValue(resourceReference, tfjsonpath.New("zones"), knownvalue.Null()),
					},
				},
			},
		},
	})
}

// TestAccClusterZonesNameRejected pins what Capella does with an AWS availability zone NAME
// (us-east-1a) where a zone ID (use1-az1) is required.
//
// This started life as the experiment for AV-149540, which asked whether Capella silently
// normalises a name to an ID - if it did, the refreshed state would disagree with the
// configuration and RequiresReplaceIfConfigured would destroy and recreate a healthy
// cluster. Run live on 2026-10-27, it does not: the create is rejected outright with
//
//	{"hint":"Please ensure availability zones specified is valid.",
//	 "message":"Unable to process request. The availability zone 'us-east-1a' is invalid.",
//	 "code":4042,"httpStatusCode":422}
//
// so that drift vector is closed and this test guards it staying closed.
//
// What the run did expose is that the rejection happens at apply rather than at plan:
// validateZones checks the cloud provider and the zone count but never the identifier
// format, so a plan over a config with zone names looks clean and the practitioner only
// finds out once the create is attempted. See TestAccClusterZonesNameRejectedAtPlanTime
// below for the behaviour that would be better.
//
// Provisions nothing - the API rejects the request before any cluster is built.
func TestAccClusterZonesNameRejected(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfig(resourceName, generateRandomCIDR(), "aws", "us-east-1", "single", `["us-east-1a"]`),
				// Terraform hard-wraps diagnostic text, and in the observed output the wrap
				// falls between "is" and "invalid" - so the whitespace here must be \s+ and
				// not a literal space, or the regex silently never matches.
				ExpectError: regexp.MustCompile(
					`(?s)Error initiating cluster create request.*availability\s+zone\s+'us-east-1a'\s+is\s+invalid`),
			},
		},
	})
}

// TestAccClusterZonesNameRejectedAtPlanTime asserts the behaviour AV-149540 now asks for
// rather than the behaviour today, so it fails until that ticket is fixed.
//
// A zone identifier the provider can see is malformed should be rejected while planning,
// like every other zones rule, instead of being sent to Capella and bounced back as a 422
// part-way through an apply.
func TestAccClusterZonesNameRejectedAtPlanTime(t *testing.T) {
	t.Skip("AV-149540: a zone name passes validateZones, which checks only the cloud " +
		"provider and the zone count, and is rejected by Capella at apply instead; unskip " +
		"once the provider validates the zone ID format at plan time")

	resourceName := randomStringWithPrefix("tf_acc_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfig(resourceName, generateRandomCIDR(), "aws", "us-east-1", "single", `["us-east-1a"]`),
				// The provider's own diagnostic, not Capella's 422.
				ExpectError: regexp.MustCompile(`(?s)Invalid zones.*us-east-1a`),
			},
		},
	})
}

// TestAccClusterZonesChangedForcesReplacement is the positive control for
// RequiresReplaceIfConfigured: a deliberately changed zone must replace the cluster. Without
// it, a test suite that only ever asserts "no replacement" would also pass if the modifier
// were dropped entirely.
//
// Provisions a cluster and then replaces it, so it is the most expensive test in this file.
func TestAccClusterZonesChangedForcesReplacement(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "multi", `["use1-az1", "use1-az2", "use1-az4"]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckTypeSetElemAttr(resourceReference, "zones.*", "use1-az4"),
				),
			},
			{
				// use1-az3 is avoided elsewhere in the suite because some compute types are
				// not offered there; use1-az5 is the spare.
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "multi", `["use1-az1", "use1-az2", "use1-az5"]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "zones.#", "3"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "zones.*", "use1-az5"),
				),
			},
		},
	})
}

// The remaining tests are rejected during plan, so they provision nothing and run in about a
// second each. They belong in sanity.list.

// TestAccClusterZonesSingleAZTooMany covers the single AZ branch of validateZones, which the
// multi AZ count test added with AV-147905 does not reach.
func TestAccClusterZonesSingleAZTooMany(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterZonesConfig(resourceName, generateRandomCIDR(), "aws", "us-east-1", "single", `["use1-az1", "use1-az2"]`),
				ExpectError: regexp.MustCompile(`(?s)Invalid zones.*got 2 zones for a single AZ cluster`),
			},
		},
	})
}

// TestAccClusterZonesEmptyList covers `zones = []`, which is a third distinct shape
// alongside an omitted attribute and an explicit null.
//
// It is accepted and means "let Capella choose", the same as omitting - but unlike omitting,
// the empty set is kept in state rather than being replaced by the zones Capella picked.
// That is deliberate (see zonesForState in internal/resources/cluster.go): the attribute is
// configured, so storing the real zones would plan a replacement against the configured [].
// The consequence is that state reports no zones for a cluster that really is in three, and
// the data source will disagree with the resource for this config shape.
//
// An earlier revision of AV-147905 rejected [] outright; that was reverted to keep accepting
// what releases before it accepted for single AZ clusters.
//
// Provisions a cluster.
func TestAccClusterZonesEmptyList(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "multi", `[]`),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					// Empty, not absent: the two are different values to Terraform, and the
					// empty set is what the configuration asked for.
					resource.TestCheckResourceAttr(resourceReference, "zones.#", "0"),
				),
			},
			{
				// The step's own post-apply plan already asserts emptiness is stable; this
				// one proves the empty set does not drift into the zones Capella assigned,
				// which would plan a replacement.
				Config: testAccClusterZonesConfig(resourceName, cidr, "aws", "us-east-1", "multi", `[]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.TestCheckResourceAttr(resourceReference, "zones.#", "0"),
			},
		},
	})
}

// TestAccClusterZonesAzureRejected covers the Azure arm of the non-AWS rejection; the test
// added with AV-147905 covers only GCP.
func TestAccClusterZonesAzureRejected(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterZonesConfig(resourceName, generateRandomCIDR(), "azure", "eastus", "multi", `["1", "2", "3"]`),
				ExpectError: regexp.MustCompile(`(?s)Invalid zones.*got zones for a azure cluster`),
			},
		},
	})
}

// TestAccClusterZonesEmptyStringsCollapse covers the empty string, which is distinct from
// both an empty list and an omitted attribute.
//
// Three empty strings are three identical elements, so the Set collapses them to one before
// validateZones counts - the config is caught, but as "1 zone" rather than as three blank
// ones. Pinned so a future change to that message is a deliberate one.
func TestAccClusterZonesEmptyStringsCollapse(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterZonesConfig(resourceName, generateRandomCIDR(), "aws", "us-east-1", "multi", `["", "", ""]`),
				ExpectError: regexp.MustCompile(`(?s)Invalid zones.*got 1 zones for a multi AZ cluster`),
			},
		},
	})
}

// TestAccClusterZonesEmptyStringRejected covers a single empty string on a single AZ cluster,
// where the count rule cannot catch it: one element is a legal count, so the blank string
// goes to Capella and is rejected there.
//
// This is AV-149540 again - validateZones checks the cloud provider and the count but never
// the identifier itself - and it is the cheapest demonstration of it, since an empty string
// is something the provider could obviously reject on its own.
//
// Provisions nothing.
func TestAccClusterZonesEmptyStringRejected(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfig(resourceName, generateRandomCIDR(), "aws", "us-east-1", "single", `[""]`),
				// Tied to the zone rejection specifically, by both the phrasing and the 4042
				// code Capella returns for an invalid availability zone. Matching only the
				// generic "Error initiating cluster create request" prefix would let this
				// test pass on any create failure - a quota, auth or transient API error -
				// without the blank zone having been rejected at all.
				//
				// Whitespace is \s+ throughout because Terraform hard-wraps diagnostics and
				// the wrap position depends on the message length.
				//
				// Capella's exact rendering of a blank zone is not yet confirmed by a live
				// run; if this fails, tighten it to whatever the message actually says
				// rather than loosening it back to the generic prefix.
				ExpectError: regexp.MustCompile(
					`(?s)Error\s+initiating\s+cluster\s+create\s+request.*availability\s+zone.*4042`),
			},
		},
	})
}

// TestAccClusterZonesNonexistentIdRejected covers a well-formed zone ID that does not exist.
//
// Unlike a zone name, no format check could catch this, and nor could it catch a valid ID
// from the wrong region (usw2-az1 in us-east-1, which also plans cleanly). Both are
// necessarily the API's job, which bounds how much AV-149540 can reasonably fix: a format
// rule would catch names and blanks at plan time, but existence and region always need
// Capella.
//
// Provisions nothing.
func TestAccClusterZonesNonexistentIdRejected(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterZonesConfig(resourceName, generateRandomCIDR(), "aws", "us-east-1", "single", `["use1-az99"]`),
				ExpectError: regexp.MustCompile(
					`(?s)Error\s+initiating\s+cluster\s+create\s+request.*availability\s+zone\s+'use1-az99'\s+is\s+invalid`),
			},
		},
	})
}

// TestAccClusterZonesDuplicatesCollapse documents how a duplicated zone is reported. The Set
// deduplicates before validateZones counts, so three entries with a repeat are rejected as
// two zones rather than as a duplicate. The message is confusing but the config is caught,
// which is what matters; this test pins the behaviour so a future change to the message is a
// deliberate one.
func TestAccClusterZonesDuplicatesCollapse(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterZonesConfig(resourceName, generateRandomCIDR(), "aws", "us-east-1", "multi", `["use1-az1", "use1-az1", "use1-az2"]`),
				ExpectError: regexp.MustCompile(`(?s)Invalid zones.*got 2 zones for a multi AZ cluster`),
			},
		},
	})
}

// testAccCheckClusterZonesMatchPlatform asserts the zones reported for a cluster are
// consistent with its cloud provider and availability type, whichever cluster the run is
// pointed at. AWS clusters must report zones - that is the feature - while GCP and Azure
// clusters must report none.
func testAccCheckClusterZonesMatchPlatform(reference string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		var attrs map[string]string
		for _, m := range state.Modules {
			if v, ok := m.Resources[reference]; ok {
				attrs = v.Primary.Attributes
			}
		}
		if attrs == nil {
			return fmt.Errorf("%s not found in state", reference)
		}

		// Presence matters as well as the count: a null collection has no zones.# entry at
		// all, while an empty one has "0". Terraform treats those as different values, so
		// counting alone would accept either.
		raw, present := attrs["zones.#"]
		count := 0
		if present {
			parsed, err := strconv.Atoi(raw)
			if err != nil {
				return fmt.Errorf("%s: zones.# is %q, not a number", reference, raw)
			}
			count = parsed
		}

		provider := attrs["cloud_provider.type"]
		availability := attrs["availability.type"]

		// Without this the check has a vacuous path: an unreadable provider type would fall
		// into the non-AWS branch, which passes whenever zones is absent - so the test would
		// report success for a cluster it never actually inspected.
		if provider == "" {
			return fmt.Errorf("%s: cloud_provider.type is not set in state, so the zone "+
				"expectation cannot be determined", reference)
		}
		if availability == "" {
			return fmt.Errorf("%s: availability.type is not set in state, so the expected "+
				"zone count cannot be determined", reference)
		}

		if provider != "aws" {
			if present {
				return fmt.Errorf(
					"%s: a %s cluster reported a zones collection of %d, expected the attribute "+
						"to be null - Capella does not report zone IDs for %s, and null is not "+
						"the same value as an empty collection",
					reference, provider, count, provider)
			}
			return nil
		}

		if !present {
			return fmt.Errorf(
				"%s: an aws cluster reported no zones attribute at all, expected the zone IDs "+
					"Capella assigned - the provider is not reading zones back from the API",
				reference)
		}

		want := 3
		if availability == "single" {
			want = 1
		}
		if count != want {
			return fmt.Errorf(
				"%s: a %s AZ aws cluster reported %d zones, expected %d - the provider is not reading zones back from the API",
				reference, availability, count, want,
			)
		}
		return nil
	}
}

func testAccClusterZonesDatasourceConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_cluster" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  id              = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId)
}

func testAccClusterZonesConfig(resourceName, cidr, cloudProvider, region, availability, zones string) string {
	return testAccClusterZonesConfigWithDescription(resourceName, cidr, cloudProvider, region, availability, zones, "Terraform acceptance test for cluster zones")
}

// testAccClusterZonesConfigWithDescription builds a cluster configuration. An empty zones
// argument omits the attribute entirely, which is a different configuration from `[]`.
func testAccClusterZonesConfigWithDescription(resourceName, cidr, cloudProvider, region, availability, zones, description string) string {
	// iops is set for AWS because an omitted iops is sent as 0 (AV-145015); GCP and Azure
	// reject iops, and Azure disks are sized by tier rather than by storage.
	disk := `storage = 50
          type    = "gp3"
          iops    = 3000`
	switch cloudProvider {
	case "gcp":
		disk = `storage = 50
          type    = "pd-ssd"`
	case "azure":
		disk = `type = "P6"`
	}

	zonesAttr := ""
	if zones != "" {
		zonesAttr = "zones = " + zones
	}

	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  name            = "%[2]s"
  description     = "%[9]s"
  couchbase_server = {
    version = var.server_version
  }
  cloud_provider = {
    type   = "%[6]s"
    region = "%[7]s"
    cidr   = "%[5]s"
  }
  %[10]s
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          %[11]s
        }
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
    }
  ]
  availability = {
    "type" : "%[8]s"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, cidr, cloudProvider, region, availability, description, zonesAttr, disk)
}
