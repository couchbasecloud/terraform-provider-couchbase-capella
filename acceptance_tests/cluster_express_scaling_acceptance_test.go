package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

// Express scaling acceptance tests come in two groups.
//
// The default group never sends PUT .../expressScaling. It only manages the resource with
// express_scaling = false on the shared global cluster, which has never had express scaling
// enabled, so the provider finds the cluster already disabled and has nothing to send. These
// tests pass whether or not the FusionRebalances flag is on for the test tenant and never
// start a rebalance on the shared cluster.
//
// The second group sends the request and needs ACC_EXPRESS_SCALING=true. Each test there
// provisions its own cluster, because enabling express scaling rebalances the cluster for a
// long time and destroying the resource leaves express scaling as it is. The tenant must have
// the FusionRebalances flag on, or the endpoint answers 404.

// TestAccClusterExpressScalingResource covers the lifecycle on the global cluster without
// changing it: create, a follow-up plan with no changes, import with both ID forms and destroy.
func TestAccClusterExpressScalingResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_express_scaling_")
	resourceReference := "couchbase-capella_cluster_express_scaling." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { requireExpressScalingDisabled(t, globalClusterId) },
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterExpressScalingConfig(resourceName, globalClusterId, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "express_scaling", "false"),
					resource.TestCheckResourceAttr(resourceReference, "express_scaling_state", string(clusterapi.ExpressScalingDisabled)),
				),
			},
			{
				Config: testAccClusterExpressScalingConfig(resourceName, globalClusterId, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateExpressScalingImportId(resourceReference, "cluster_id"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "cluster_id",
				ImportStateVerifyIgnore:              []string{"timeouts"},
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateExpressScalingImportId(resourceReference, "id"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "cluster_id",
				ImportStateVerifyIgnore:              []string{"timeouts"},
			},
		},
		// Destroy only removes the resource from state, so the cluster must still be disabled.
		CheckDestroy: func(*terraform.State) error {
			return checkExpressScalingState(globalClusterId, clusterapi.ExpressScalingDisabled)
		},
	})
}

// TestAccClusterExpressScalingInvalidCluster verifies the error when the cluster does not exist.
// The provider reads the cluster before sending anything, so this needs no flag.
func TestAccClusterExpressScalingInvalidCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_express_scaling_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterExpressScalingConfig(resourceName, uuid.Nil.String(), true),
				ExpectError: apiErrorPattern("Error reading cluster express scaling", "Could not read cluster before changing express scaling", "404"),
			},
		},
	})
}

func TestAccClusterExpressScalingInvalidClusterId(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_express_scaling_invalid_id_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterExpressScalingConfig(resourceName, "not-a-uuid", true),
				ExpectError: regexp.MustCompile(`(?s)cluster_id.*must be a valid UUID`),
			},
		},
	})
}

// TestAccClusterExpressScalingEnableDisable enables and disables express scaling on a dedicated
// Enterprise cluster. The first apply uses a short timeout, so it usually ends with a warning while
// the cluster is still enabling and saves the planned value. The second apply only raises the
// timeout and waits for the cluster to finish. Whether the cluster is still enabling at that point
// depends on Capella, so the in-progress path is pinned down by the resource's unit tests.
//
// The resource is then removed from the configuration, which must leave express scaling enabled.
// Recreating it with express_scaling = true sends nothing, since the cluster is already enabled,
// and changing it to false disables express scaling through Update.
func TestAccClusterExpressScalingEnableDisable(t *testing.T) {
	requireExpressScalingAcceptance(t)

	clusterResourceName := randomStringWithPrefix("tf_acc_cluster_express_scaling_")
	resourceName := randomStringWithPrefix("tf_acc_express_scaling_")
	clusterReference := "couchbase-capella_cluster." + clusterResourceName
	resourceReference := "couchbase-capella_cluster_express_scaling." + resourceName
	cidr := generateRandomCIDR()

	config := func(enabled bool, timeout string) string {
		return testAccClusterExpressScalingDedicatedClusterConfig(clusterResourceName, resourceName, cidr, "enterprise", enabled, timeout)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: config(true, "1m"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(resourceReference, "express_scaling", "true"),
					resource.TestMatchResourceAttr(resourceReference, "express_scaling_state", regexp.MustCompile(`^(enabling|enabled)$`)),
				),
			},
			{
				Config: config(true, "4h"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(resourceReference, "express_scaling", "true"),
					resource.TestCheckResourceAttr(resourceReference, "express_scaling_state", string(clusterapi.ExpressScalingEnabled)),
				),
			},
			{
				Config: config(true, "4h"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateExpressScalingImportId(resourceReference, "cluster_id"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "cluster_id",
				ImportStateVerifyIgnore:              []string{"timeouts"},
			},
			{
				Config: globalProviderBlock + testAccClusterExpressScalingClusterBlock(clusterResourceName, cidr, "enterprise", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionDestroy),
						plancheck.ExpectResourceAction(clusterReference, plancheck.ResourceActionNoop),
					},
				},
				Check: testAccCheckClusterExpressScalingState(clusterReference, clusterapi.ExpressScalingEnabled),
			},
			{
				Config: config(true, "4h"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "express_scaling", "true"),
					resource.TestCheckResourceAttr(resourceReference, "express_scaling_state", string(clusterapi.ExpressScalingEnabled)),
				),
			},
			{
				Config: config(false, "4h"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "express_scaling", "false"),
					resource.TestCheckResourceAttr(resourceReference, "express_scaling_state", string(clusterapi.ExpressScalingDisabled)),
				),
			},
		},
		CheckDestroy: testAccCheckClusterDestroyed(t, clusterReference),
	})
}

// TestAccClusterExpressScalingRejected covers the rejections that need no rebalance, on a
// dedicated Developer Pro cluster: the Enterprise plan check, then an API key that may only
// view the project. The Developer Pro plan also keeps a broken role check from starting a rebalance.
func TestAccClusterExpressScalingRejected(t *testing.T) {
	requireExpressScalingAcceptance(t)

	viewerToken := createProjectViewerApiKey(t)

	clusterResourceName := randomStringWithPrefix("tf_acc_cluster_express_scaling_dev_")
	resourceName := randomStringWithPrefix("tf_acc_express_scaling_dev_")
	clusterReference := "couchbase-capella_cluster." + clusterResourceName
	cidr := generateRandomCIDR()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterExpressScalingDedicatedClusterConfig(clusterResourceName, resourceName, cidr, "developer pro", true, "10m"),
				ExpectError: apiErrorPattern(
					"Error changing cluster express scaling",
					"Enterprise plan is required to enable express scaling",
					"422",
				),
			},
			{
				Config: testAccClusterExpressScalingViewerConfig(clusterResourceName, resourceName, cidr, viewerToken),
				ExpectError: apiErrorPattern(
					"Error changing cluster express scaling",
					"Access Denied",
					"403",
				),
			},
		},
		CheckDestroy: testAccCheckClusterDestroyed(t, clusterReference),
	})
}

// TestAccClusterResourceExpressScaling creates a cluster with express scaling on through the
// cluster resource attribute and checks that the next plan is empty. Enabling starts after the
// cluster is healthy, so create can save a state that still reads disabled. The test waits for
// Capella to report enabling or enabled before the follow-up plan refreshes the state.
func TestAccClusterResourceExpressScaling(t *testing.T) {
	requireExpressScalingAcceptance(t)

	clusterResourceName := randomStringWithPrefix("tf_acc_cluster_express_scaling_attr_")
	clusterReference := "couchbase-capella_cluster." + clusterResourceName
	config := testAccClusterResourceExpressScalingConfig(clusterResourceName, generateRandomCIDR())

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(t, clusterReference),
					resource.TestCheckResourceAttr(clusterReference, "express_scaling", "true"),
					testAccWaitForClusterExpressScalingEnabled(clusterReference, 15*time.Minute),
				),
			},
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterReference, "express_scaling", "true"),
					resource.TestMatchResourceAttr(clusterReference, "express_scaling_state", regexp.MustCompile(`^(enabling|enabled)$`)),
				),
			},
		},
		CheckDestroy: testAccCheckClusterDestroyed(t, clusterReference),
	})
}

// requireExpressScalingAcceptance skips tests that send the express scaling request. They take
// hours and need the FusionRebalances flag on the tenant, so they run only when asked for.
func requireExpressScalingAcceptance(t *testing.T) {
	t.Helper()
	if enabled, _ := strconv.ParseBool(os.Getenv("ACC_EXPRESS_SCALING")); !enabled {
		t.Skip("skipping: set ACC_EXPRESS_SCALING=true to run express scaling tests that provision their own cluster. " +
			"The tenant needs the FusionRebalances flag on.")
	}
}

// testAccCheckClusterExpressScalingState checks the express scaling state Capella reports for the
// cluster in state, independently of any express scaling resource.
func testAccCheckClusterExpressScalingState(clusterReference string, want clusterapi.ExpressScalingState) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		clusterResource, ok := state.RootModule().Resources[clusterReference]
		if !ok {
			return fmt.Errorf("resource %s not found in state", clusterReference)
		}
		return checkExpressScalingState(clusterResource.Primary.Attributes["id"], want)
	}
}

// requireExpressScalingDisabled skips the test unless the cluster reports express scaling as
// disabled. A cluster with express scaling on or changing would be disabled by the test, and one
// that reports no value would fail the assertions, which expect the value Capella returns as is.
func requireExpressScalingDisabled(t *testing.T, clusterID string) {
	t.Helper()
	if err := checkExpressScalingState(clusterID, clusterapi.ExpressScalingDisabled); err != nil {
		t.Skipf("skipping: %v", err)
	}
}

func checkExpressScalingState(clusterID string, want clusterapi.ExpressScalingState) error {
	got, err := getClusterExpressScalingState(clusterID)
	if err != nil {
		return err
	}
	if got != want {
		return fmt.Errorf("cluster %s reports express scaling %q, want %q", clusterID, got, want)
	}
	return nil
}

func getClusterExpressScalingState(clusterID string) (clusterapi.ExpressScalingState, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", globalHost, globalOrgId, globalProjectId, clusterID)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
	if err != nil {
		return "", fmt.Errorf("could not read cluster %s: %w", clusterID, err)
	}

	var cluster clusterapi.GetClusterResponse
	if err := json.Unmarshal(response.Body, &cluster); err != nil {
		return "", fmt.Errorf("could not decode cluster %s: %w", clusterID, err)
	}
	return cluster.ExpressScaling, nil
}

// testAccWaitForClusterExpressScalingEnabled waits until Capella reports express scaling as enabling
// or enabled for the cluster in state, failing after timeout.
func testAccWaitForClusterExpressScalingEnabled(clusterReference string, timeout time.Duration) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		clusterResource, ok := state.RootModule().Resources[clusterReference]
		if !ok {
			return fmt.Errorf("resource %s not found in state", clusterReference)
		}
		clusterID := clusterResource.Primary.Attributes["id"]

		deadline := time.Now().Add(timeout)
		for {
			got, err := getClusterExpressScalingState(clusterID)
			if err == nil && got.IsEnabled() {
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("cluster %s did not report express scaling as enabling or enabled within %s (last state %q, error %v)", clusterID, timeout, got, err)
			}
			time.Sleep(30 * time.Second)
		}
	}
}

// createProjectViewerApiKey creates an API key that can view the global project but not manage
// its clusters, and deletes it when the test ends.
func createProjectViewerApiKey(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	projectType := "project"
	expiry := float32(1)
	request := api.CreateApiKeyRequest{
		Name:              randomStringWithPrefix("tf_acc_express_scaling_viewer_"),
		OrganizationRoles: []string{"organizationMember"},
		Expiry:            &expiry,
		Resources: &api.Resources{
			{Type: &projectType, Id: uuid.MustParse(globalProjectId), Roles: []string{"projectViewer"}},
		},
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/apikeys", globalHost, globalOrgId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := globalClient.ExecuteWithRetry(ctx, cfg, request, globalToken, nil)
	if err != nil {
		t.Fatalf("could not create project viewer API key: %v", err)
	}

	var apiKey api.CreateApiKeyResponse
	if err := json.Unmarshal(response.Body, &apiKey); err != nil {
		t.Fatalf("could not decode project viewer API key: %v", err)
	}

	t.Cleanup(func() {
		url := fmt.Sprintf("%s/v4/organizations/%s/apikeys/%s", globalHost, globalOrgId, apiKey.Id)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
		if _, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil); err != nil {
			t.Logf("WARNING: failed to delete project viewer API key %s: %v", apiKey.Id, err)
		}
	})

	return apiKey.Token
}

// --- Config builders ---

func testAccClusterExpressScalingConfig(resourceName, clusterID string, enabled bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_express_scaling" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  express_scaling = %[6]t
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, clusterID, enabled)
}

// testAccClusterExpressScalingClusterBlock returns a cluster resource block. extra is added inside
// the block as is.
func testAccClusterExpressScalingClusterBlock(clusterResourceName, cidr, plan, extra string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_cluster" "%[3]s" {
	organization_id = "%[1]s"
	project_id      = "%[2]s"
	name            = "%[3]s"
	couchbase_server = {
		version = var.server_version
	}
	cloud_provider = {
		type   = "aws"
		region = "us-east-1"
		cidr   = "%[4]s"
	}
	service_groups = [
		{
			node = {
				compute = {
					cpu = 4
					ram = 16
				}
				disk = {
					storage = 50
					type    = "gp3"
					iops    = 3000
				}
			}
			num_of_nodes = 3
			services     = ["data", "index", "query"]
		}
	]
	availability = {
		"type" : "multi"
	}
	support = {
		plan     = "%[5]s"
		timezone = "PT"
	}
	%[6]s
}
`, globalOrgId, globalProjectId, clusterResourceName, cidr, plan, extra)
}

func testAccClusterExpressScalingDedicatedClusterConfig(clusterResourceName, resourceName, cidr, plan string, enabled bool, timeout string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s
resource "couchbase-capella_cluster_express_scaling" "%[5]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = couchbase-capella_cluster.%[6]s.id
	express_scaling = %[7]t

	timeouts {
		create = "%[8]s"
		update = "%[8]s"
	}
}
`, globalProviderBlock, testAccClusterExpressScalingClusterBlock(clusterResourceName, cidr, plan, ""),
		globalOrgId, globalProjectId, resourceName, clusterResourceName, enabled, timeout)
}

func testAccClusterResourceExpressScalingConfig(clusterResourceName, cidr string) string {
	return globalProviderBlock + testAccClusterExpressScalingClusterBlock(clusterResourceName, cidr, "enterprise", "express_scaling = true")
}

// testAccClusterExpressScalingViewerConfig keeps the cluster on the default provider and sends the
// express scaling request through a second provider authenticated with the project viewer key.
func testAccClusterExpressScalingViewerConfig(clusterResourceName, resourceName, cidr, viewerToken string) string {
	return fmt.Sprintf(`
%[1]s

provider "couchbase-capella" {
  alias                = "viewer"
  host                 = var.host
  authentication_token = %[2]q
}
%[3]s
resource "couchbase-capella_cluster_express_scaling" "%[6]s" {
	provider        = couchbase-capella.viewer
	organization_id = "%[4]s"
	project_id      = "%[5]s"
	cluster_id      = couchbase-capella_cluster.%[7]s.id
	express_scaling = true
}
`, globalProviderBlock, viewerToken, testAccClusterExpressScalingClusterBlock(clusterResourceName, cidr, "developer pro", ""),
		globalOrgId, globalProjectId, resourceName, clusterResourceName)
}

// --- Import ID ---

// generateExpressScalingImportId builds the import ID with the cluster under idKey, which is
// "cluster_id" or its "id" alias.
func generateExpressScalingImportId(resourceReference, idKey string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		if rawState == nil {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		return fmt.Sprintf(
			"%s=%s,project_id=%s,organization_id=%s",
			idKey, rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}
