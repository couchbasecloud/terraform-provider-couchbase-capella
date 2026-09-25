package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// The tests in this file all cover one behaviour: what the cluster resource puts on the
// wire for a disk attribute that the practitioner left out of config.
//
// `storage`, `iops` and `autoexpansion` are all Optional+Computed
// (internal/resources/cluster_schema.go:63-65), so an omitted attribute arrives at the
// provider as *unknown*, not null. A guard that tests only IsNull() therefore falls
// through, and ValueInt64()/ValueBool() unwrap the unknown to its zero value, which then
// goes into the request body as an explicit 0 or false. Capella cannot tell that apart
// from a deliberate choice, so its own default never applies.
//
// AV-143040 and AV-143042 fixed the two Azure occurrences of this in
// morphToApiServiceGroups (internal/resources/cluster.go:702-726). Two gaps remain:
//
//   - the AWS branch (cluster.go:681-687) still has it, and unlike Azure cannot be fixed by
//     adding an IsUnknown() guard - see AV-145015 and the skipped tests below, which assert
//     the intended behaviour rather than the current one;
//   - on Azure updates the provider now omits `autoexpansion` whenever set-element
//     correlation breaks, and relies on the server treating an absent field as "leave
//     as-is" (AV-143041, Capella 2.2.288) rather than upholding the contract itself.
//
// The Azure assertions below were run live on 2026-09-23; the per-test comments record what
// each one established, including one expectation that turned out to be wrong. The AWS tests
// are skipped: they state the contract AV-145015 will deliver, and fail against today's
// behaviour by design.
//
// The request-shaping itself is pinned deterministically by the unit test in
// internal/resources/cluster_azure_disk_test.go. These tests cover what the practitioner
// actually experiences end to end, which the unit test cannot show.

// TestAccClusterResourceAzureUltraDiskIopsOmitted_AV_143042 covers AV-143042: the Azure
// Ultra-disk `iops` guard checked Storage.IsUnknown() instead of IOPS.IsUnknown(), so with
// `storage` set and `iops` omitted the provider sent "iops": 0.
//
// Capella requires IOPS for Ultra disks either way, so both the broken and the fixed
// provider fail this apply - but they fail differently, and that difference is the assertion:
//
//	before the fix: "The storage IOPS value is invalid for the service group '...',
//	                 should be between <min> and <max> inclusive but is 0."
//	after the fix:  "The storage IOPS value is missing for the service group '...'."
//
// Matching on "missing" fails against the unfixed provider, which is what makes this a
// regression guard rather than a test that merely confirms Ultra disks need IOPS.
//
// This fails during request validation, so no cluster is provisioned and the test is cheap.
func TestAccClusterResourceAzureUltraDiskIopsOmitted_AV_143042(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfigAzureUltraDisk(resourceName, cidr, "storage = 1024"),
				ExpectError: apiErrorPattern("Error initiating cluster create request",
					"storage IOPS value is missing", "422"),
			},
		},
	})
}

// TestAccClusterResourceAzureUltraDiskIopsInvalid_AV_143042 is the control for the test
// above. An out-of-range `iops` must come back as "invalid ... but is 1", which proves the
// value the practitioner wrote is genuinely reaching Capella and being range-checked. Without
// this, a "missing" error in the test above could just as easily mean the provider dropped
// `iops` unconditionally, and the fix would look correct while being wrong in a new way.
func TestAccClusterResourceAzureUltraDiskIopsInvalid_AV_143042(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfigAzureUltraDisk(resourceName, cidr,
					"storage = 1024\n          iops    = 1"),
				ExpectError: apiErrorPattern("storage IOPS value is invalid",
					"inclusive but is 1.", "422"),
			},
		},
	})
}

// TestAccClusterResourceAwsDiskIopsOmitted_AV_145015 and its `storage` sibling below state
// the behaviour the provider should have and does not yet, so both are skipped.
//
// They assert a provider-side rejection naming the missing attribute. Today the config is
// accepted, `iops: 0` goes on the wire, and Capella answers with a range error about a value
// the practitioner never wrote:
//
//	"The storage IOPS value is invalid for the service group '...', should be between
//	 3000 and 16000 inclusive but is 0."
//
// The patterns below do not match that, so these tests genuinely fail today - which is why
// they are skipped rather than left red, and what makes them one line from becoming a
// regression guard.
//
// Rejection rather than omission is the right contract: `storage` and `iops` are plain ints
// with no `omitempty` in the API contract on both sides, so leaving them out of the JSON
// still unmarshals to 0 server-side and produces the identical 422. There is no AWS default
// for Capella to apply. See AV-145015 for the full analysis.
//
// The patterns accept either plausible fix shape - a validateCreateCluster check following
// the existing AWS `autoexpansion` and GCP `iops` rules (internal/resources/cluster.go:833-849),
// or making the attributes Required in the schema. Tighten to whichever ships.
//
// When unskipping, add both tests to acceptance_tests/sanity.list: they are rejected before
// anything is provisioned and cost under two seconds each.
func TestAccClusterResourceAwsDiskIopsOmitted_AV_145015(t *testing.T) {
	t.Skip("AV-145015: an omitted AWS iops is unknown, serialised as 0 and rejected by " +
		"Capella with a range error; unskip once the provider rejects the config itself")

	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterConfigAwsDisk(resourceName, cidr, "type    = \"gp3\"\n          storage = 50"),
				ExpectError: regexp.MustCompile(`(?is)invalid configuration.*iops|attribute "iops" is required`),
			},
		},
	})
}

// TestAccClusterResourceAwsDiskStorageOmitted_AV_145015 is the `storage` half of the same
// defect. See the comment above for why this asserts rejection and is skipped.
func TestAccClusterResourceAwsDiskStorageOmitted_AV_145015(t *testing.T) {
	t.Skip("AV-145015: an omitted AWS storage is unknown, serialised as 0 and rejected by " +
		"Capella with a range error; unskip once the provider rejects the config itself")

	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterConfigAwsDisk(resourceName, cidr, "type = \"gp3\"\n          iops = 3000"),
				ExpectError: regexp.MustCompile(`(?is)invalid configuration.*storage|attribute "storage" is required`),
			},
		},
	})
}

// TestAccClusterResourceAzureScaleOmitsAutoexpansion_AV_143041 covers the Terraform side of
// AV-143041, which the ticket's "Terraform impact: None" note concluded was unreachable.
//
// That conclusion holds only for an update that touches nothing inside the service group.
// `service_groups` is a SetNestedAttribute (internal/resources/cluster_schema.go:90), and
// Terraform correlates set elements by value rather than by position, so changing any
// configured attribute inside an element leaves it with no prior element to correlate
// against and its Optional+Computed attributes plan as unknown.
//
// Verified with `terraform plan -refresh=false` against a handwritten prior state holding
// autoexpansion=false, using Terraform 1.14.1 and a dev_overrides build of this provider
// (no Capella calls involved, since Configure opens no connection and Cluster has no
// ModifyPlan). Reading after_unknown out of the JSON plan:
//
//	description changed (outside the set)   autoexpansion = false      <- correlates
//	num_of_nodes  3 -> 4                    autoexpansion = UNKNOWN
//	compute cpu/ram 4/16 -> 8/32            autoexpansion = UNKNOWN
//	services ["data"] -> ["data","index"]   autoexpansion = UNKNOWN
//	second service group added              UNKNOWN on the new group only
//
// So the provider omits `autoexpansion` on most Azure cluster updates where the practitioner
// has not pinned it in config - horizontal scaling, vertical scaling and service changes, not
// just adding a service group. `storage` and `iops` go unknown on the same elements.
//
// What saves this today is the server, not the provider. AV-143041 (fixed in Capella 2.2.288)
// changed the update path so an absent autoExpansion means "leave as-is" rather than adopting
// the create-time default of true. Ran live 2026-09-23 against a non-production Capella environment: scaling 3 -> 4 nodes
// with `autoexpansion` omitted left it false, and no 422. Before that server fix, the same
// request would have been read as "enable" - an unrequested Azure swap rebalance (node
// replacement and data movement) on a cluster deliberately configured with auto-expansion
// off, or a flat 422 "Disk auto expansion is disabled on the tenant" for a tenant without
// the EnableAzureStorageAutoExpansion flag, blocking a scale-up that has nothing to do with
// disk settings.
//
// So this test guards a contract the provider does not currently uphold on its own. Keep it:
// it is what will catch the server-side behaviour regressing, and it fails loudly if the
// provider ever stops being shielded.
//
// The assertion is deliberately on the outcome and not on the plan: whether Terraform
// correlates the set element is an implementation detail that may change, but
// auto-expansion staying off across a scale-up is the contract either way. Note a pass does
// not distinguish "the server preserved it" from "the provider sent it explicitly" - it
// pins the result the practitioner sees, which is the thing worth pinning.
func TestAccClusterResourceAzureScaleOmitsAutoexpansion_AV_143041(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Auto-expansion explicitly off, so there is no ambiguity about what the
			// operator asked for.
			{
				Config: testAccClusterConfigAzureScale(resourceName, cidr, 3, "\n          autoexpansion = false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "false"),
				),
			},
			// Scale out, and drop autoexpansion from config. It must stay off: the operator
			// asked to add a node, not to turn on auto-expansion.
			{
				Config: testAccClusterConfigAzureScale(resourceName, cidr, 4, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "service_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "false"),
					// The scale-up must have actually happened rather than the provider
					// returning stale state, otherwise the assertion above proves nothing.
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "P6"),
				),
			},
		},
	})
}

// TestAccClusterResourceAzureAddServiceGroupOmitsAutoexpansion_AV_143041 covers the other
// route to the same AV-143041 exposure: a service group added to an existing Azure cluster
// has no prior state at all, so its `autoexpansion` is necessarily unknown and the provider
// necessarily omits it.
//
// The added group inherits the cluster's existing auto-expansion setting - here false - it
// does NOT get Capella's create-time default of true. AV-139797's default lives on the
// create path (setV2Disk); the update path (setDisk) was changed by AV-143041 to preserve
// whatever the cluster already had whenever the field is absent. An added service group goes
// through the update path, so it follows the cluster, not the create default.
//
// This was established the hard way: the test originally asserted true on the added group
// and failed on 2026-09-23 with both groups reading false. That is the
// correct behaviour and the assertion was wrong, so it now pins inheritance.
//
// The property that actually matters for AV-143041 is the one this asserts: adding a service
// group must not flip auto-expansion ON for a cluster deliberately running with it off,
// because enabling it on Azure means a swap rebalance. If the server-side fix regresses, the
// added group comes back true and this fails.
//
// Assertions use TestCheckTypeSetElemNestedAttrs because service_groups is a set and the
// order of its two elements in state is not guaranteed - in the live run the newly added
// index/query group landed at index 0 and the original data group at index 1.
func TestAccClusterResourceAzureAddServiceGroupOmitsAutoexpansion_AV_143041(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cluster_")
	resourceReference := "couchbase-capella_cluster." + resourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfigAzureScale(resourceName, cidr, 3, "\n          autoexpansion = false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "false"),
				),
			},
			{
				Config: testAccClusterConfigAzureTwoServiceGroups(resourceName, cidr),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "service_groups.#", "2"),
					// The original data group keeps auto-expansion off.
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, "service_groups.*", map[string]string{
						"num_of_nodes":            "3",
						"node.disk.autoexpansion": "false",
						"node.disk.type":          "P6",
					}),
					// The added index/query group inherits the cluster's setting rather than
					// the create-time default. True here would mean an unrequested swap
					// rebalance on a cluster the operator runs with auto-expansion off.
					resource.TestCheckTypeSetElemNestedAttrs(resourceReference, "service_groups.*", map[string]string{
						"num_of_nodes":            "2",
						"node.disk.autoexpansion": "false",
						"node.disk.type":          "P6",
					}),
				),
			},
		},
	})
}

// testAccClusterConfigAzureUltraDisk builds an Azure cluster with an Ultra disk. diskFields
// is injected verbatim into the disk block so a caller can omit an attribute entirely,
// which is the whole point of these tests - an attribute left out of config is unknown,
// and unknown is what the guards under test mishandle.
func testAccClusterConfigAzureUltraDisk(resourceName, cidr, diskFields string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Terraform Acceptance Test Azure Ultra disk"
  cloud_provider = {
    type   = "azure"
    region = "eastus"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          type    = "Ultra"
          %[6]s
        }
      }
      num_of_nodes = 3
      services     = ["data"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, cidr, diskFields)
}

// testAccClusterConfigAwsDisk builds an AWS cluster whose entire disk block is supplied by
// the caller, so a test can leave out `storage` or `iops` and observe what the provider
// sends in their place.
func testAccClusterConfigAwsDisk(resourceName, cidr, diskBlock string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Terraform Acceptance Test AWS disk defaults"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          %[6]s
        }
      }
      num_of_nodes = 3
      services     = ["data"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, cidr, diskBlock)
}

// testAccClusterConfigAzureScale builds a single-service-group Azure cluster with a
// caller-controlled node count, and autoexpansionLine either set or left out entirely.
func testAccClusterConfigAzureScale(resourceName, cidr string, numOfNodes int, autoexpansionLine string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Terraform Acceptance Test Azure autoexpansion across scaling"
  cloud_provider = {
    type   = "azure"
    region = "eastus"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          type = "P6"%[7]s
        }
      }
      num_of_nodes = %[6]d
      services     = ["data"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, cidr, numOfNodes, autoexpansionLine)
}

// testAccClusterConfigAzureTwoServiceGroups adds an index/query service group to the cluster
// built by testAccClusterConfigAzureScale. The original data group keeps its explicit
// autoexpansion = false; the added group omits the attribute, which is the case under test.
func testAccClusterConfigAzureTwoServiceGroups(resourceName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  description     = "Terraform Acceptance Test Azure autoexpansion across scaling"
  cloud_provider = {
    type   = "azure"
    region = "eastus"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          type          = "P6"
          autoexpansion = false
        }
      }
      num_of_nodes = 3
      services     = ["data"]
    },
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          type = "P6"
        }
      }
      num_of_nodes = 2
      services     = ["index", "query"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, resourceName, cidr)
}
