package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/data_api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// TestAccDataApiResource walks enable, import, peering and disable on the shared cluster.
func TestAccDataApiResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_data_api_")
	resourceReference := "couchbase-capella_data_api." + resourceName
	dsName := randomStringWithPrefix("tf_acc_data_api_ds_")
	dsReference := "data.couchbase-capella_data_api." + dsName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApiResourceAndDatasourceConfig(resourceName, dsName, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDataApiResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_data_api", "enabled"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_network_peering", "disabled"),
					resource.TestCheckResourceAttrSet(resourceReference, "connection_string"),
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dsReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(dsReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttr(dsReference, "state_for_data_api", "enabled"),
					resource.TestCheckResourceAttr(dsReference, "state_for_network_peering", "disabled"),
					resource.TestCheckResourceAttrSet(dsReference, "connection_string"),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    generateDataApiImportId(resourceReference),
				ImportStateVerifyIdentifierAttribute: "cluster_id",
			},
			// Malformed imports. Patterns use \s+ because Terraform wraps diagnostics at ~76 cols.
			// Too few pairs.
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateIdFunc: malformedDataApiImportId(resourceReference, "cluster_id=%[1]s,project_id=%[2]s"),
				ExpectError:       regexp.MustCompile(`(?s)error\s+parsing\s+terraform\s+import`),
			},
			// Too many pairs.
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateIdFunc: malformedDataApiImportId(resourceReference, "cluster_id=%[1]s,project_id=%[2]s,organization_id=%[3]s,id=extra"),
				ExpectError:       regexp.MustCompile(`(?s)error\s+parsing\s+terraform\s+import`),
			},
			// A key that is absent from the provider's importIds map entirely.
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateIdFunc: malformedDataApiImportId(resourceReference, "cluster_id=%[1]s,nonsense_id=%[2]s,organization_id=%[3]s"),
				ExpectError:       regexp.MustCompile(`(?s)error\s+parsing\s+terraform\s+import`),
			},
			// In importIds but not a key of this resource: parses, then fails the presence check.
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateIdFunc: malformedDataApiImportId(resourceReference, "cluster_id=%[1]s,bucket_id=%[2]s,organization_id=%[3]s"),
				ExpectError:       regexp.MustCompile(`(?s)resource\s+was\s+(missing|empty)`),
			},
			// Empty value for a required key.
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateIdFunc: malformedDataApiImportId(resourceReference, "cluster_id=,project_id=%[2]s,organization_id=%[3]s"),
				ExpectError:       regexp.MustCompile(`(?s)resource\s+was\s+(missing|empty)`),
			},
			// Bare UUID: must fail parsing rather than reaching the API as a confusing 404.
			{
				ResourceName:      resourceReference,
				ImportState:       true,
				ImportStateIdFunc: malformedDataApiImportId(resourceReference, "%[1]s"),
				ExpectError:       regexp.MustCompile(`(?s)error\s+parsing\s+terraform\s+import`),
			},
			{
				Config: testAccDataApiResourceAndDatasourceConfig(resourceName, dsName, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDataApiResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "true"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_data_api", "enabled"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_network_peering", "enabled"),
					resource.TestCheckResourceAttrSet(resourceReference, "connection_string"),
					resource.TestCheckResourceAttr(dsReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(dsReference, "enable_network_peering", "true"),
					resource.TestCheckResourceAttr(dsReference, "state_for_data_api", "enabled"),
					resource.TestCheckResourceAttr(dsReference, "state_for_network_peering", "enabled"),
					resource.TestCheckResourceAttrSet(dsReference, "connection_string"),
				),
			},
			// Peering off, Data API on: the update must not cascade or drop connection_string.
			{
				Config: testAccDataApiResourceAndDatasourceConfig(resourceName, dsName, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDataApiResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_data_api", "enabled"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_network_peering", "disabled"),
					resource.TestCheckResourceAttrSet(resourceReference, "connection_string"),
					resource.TestCheckResourceAttr(dsReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(dsReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttr(dsReference, "state_for_network_peering", "disabled"),
				),
			},
			{
				Config: testAccDataApiResourceAndDatasourceConfig(resourceName, dsName, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDataApiResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "enable_data_api", "false"),
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_data_api", "disabled"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_network_peering", "disabled"),
					// Exact "": AttrSet would pass on a stale hostname from the enabled steps.
					resource.TestCheckResourceAttr(resourceReference, "connection_string", ""),
					resource.TestCheckResourceAttr(dsReference, "enable_data_api", "false"),
					resource.TestCheckResourceAttr(dsReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttr(dsReference, "state_for_data_api", "disabled"),
					resource.TestCheckResourceAttr(dsReference, "state_for_network_peering", "disabled"),
					resource.TestCheckResourceAttr(dsReference, "connection_string", ""),
				),
			},
		},
	})
}

// TestAccDataApiResourceNetworkPeeringWithoutDataApi verifies that enabling network peering while the Data API is
// disabled is rejected by local config validation instead of being sent to the API.
func TestAccDataApiResourceNetworkPeeringWithoutDataApi(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_data_api_peering_only_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataApiResourceConfig(resourceName, false, true),
				ExpectError: regexp.MustCompile(`Network peering cannot be enabled`),
			},
		},
	})
}

// TestAccDataApiResourceInvalidUUIDs verifies that each ID attribute rejects
// non-UUID values via local schema validation, one attribute per subtest.
func TestAccDataApiResourceInvalidUUIDs(t *testing.T) {
	tests := []struct {
		name           string
		organizationID string
		projectID      string
		clusterID      string
	}{
		{
			name:           "organization_id",
			organizationID: "not-a-uuid",
			projectID:      "11111111-1111-1111-1111-111111111111",
			clusterID:      "22222222-2222-2222-2222-222222222222",
		},
		{
			name:           "project_id",
			organizationID: "00000000-0000-0000-0000-000000000000",
			projectID:      "not-a-uuid",
			clusterID:      "22222222-2222-2222-2222-222222222222",
		},
		{
			name:           "cluster_id",
			organizationID: "00000000-0000-0000-0000-000000000000",
			projectID:      "11111111-1111-1111-1111-111111111111",
			clusterID:      "not-a-uuid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resourceName := randomStringWithPrefix("tf_acc_data_api_non_uuid_")

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: testAccDataApiResourceIDsConfig(
							resourceName, test.organizationID, test.projectID, test.clusterID),
						ExpectError: regexp.MustCompile(
							`(?s)Invalid Attribute Value Match.*` + test.name + `.*must be a valid UUID`),
					},
				},
			})
		})
	}
}

// TestAccDatasourceDataApiInvalidUUIDs verifies that each ID attribute rejects
// non-UUID values via local schema validation, one attribute per subtest.
func TestAccDatasourceDataApiInvalidUUIDs(t *testing.T) {
	tests := []struct {
		name           string
		organizationID string
		projectID      string
		clusterID      string
	}{
		{
			name:           "organization_id",
			organizationID: "not-a-uuid",
			projectID:      "11111111-1111-1111-1111-111111111111",
			clusterID:      "22222222-2222-2222-2222-222222222222",
		},
		{
			name:           "project_id",
			organizationID: "00000000-0000-0000-0000-000000000000",
			projectID:      "not-a-uuid",
			clusterID:      "22222222-2222-2222-2222-222222222222",
		},
		{
			name:           "cluster_id",
			organizationID: "00000000-0000-0000-0000-000000000000",
			projectID:      "11111111-1111-1111-1111-111111111111",
			clusterID:      "not-a-uuid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dsName := randomStringWithPrefix("tf_acc_data_api_ds_non_uuid_")

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: testAccDataApiDatasourceIDsConfig(
							dsName, test.organizationID, test.projectID, test.clusterID),
						ExpectError: regexp.MustCompile(
							`(?s)Invalid Attribute Value Match.*` + test.name + `.*must be a valid UUID`),
					},
				},
			})
		})
	}
}

func testAccDataApiResourceConfig(resourceName string, enableDataApi, enableNetworkPeering bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_data_api" "%[2]s" {
  organization_id        = "%[3]s"
  project_id             = "%[4]s"
  cluster_id             = "%[5]s"
  enable_data_api        = %[6]t
  enable_network_peering = %[7]t
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, enableDataApi, enableNetworkPeering)
}

func testAccDataApiResourceIDsConfig(resourceName, organizationID, projectID, clusterID string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_data_api" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  enable_data_api = true
}
`, globalProviderBlock, resourceName, organizationID, projectID, clusterID)
}

func testAccDataApiDatasourceIDsConfig(dsName, organizationID, projectID, clusterID string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_data_api" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, dsName, organizationID, projectID, clusterID)
}

func testAccDataApiResourceAndDatasourceConfig(resourceName, dsName string, enableDataApi, enableNetworkPeering bool) string {
	// enable_network_peering is omitted when false so those steps also exercise the attribute's default.
	enableNetworkPeeringAttr := ""
	if enableNetworkPeering {
		enableNetworkPeeringAttr = "\n  enable_network_peering = true"
	}

	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_data_api" "%[2]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"
  enable_data_api = %[7]t%[8]s
}

data "couchbase-capella_data_api" "%[3]s" {
  organization_id = "%[4]s"
  project_id      = "%[5]s"
  cluster_id      = "%[6]s"

  depends_on = [couchbase-capella_data_api.%[2]s]
}
`, globalProviderBlock, resourceName, dsName, globalOrgId, globalProjectId, globalClusterId, enableDataApi, enableNetworkPeeringAttr)
}

func generateDataApiImportId(resourceReference string) resource.ImportStateIdFunc {
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
			"cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["cluster_id"],
			rawState["project_id"],
			rawState["organization_id"],
		), nil
	}
}

func retrieveDataApiStatusFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) (*data_api.GetDataApiStatusResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI",
		data.HostURL, organizationId, projectId, clusterId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
	if err != nil {
		return nil, err
	}
	status := &data_api.GetDataApiStatusResponse{}
	if err := json.Unmarshal(response.Body, status); err != nil {
		return nil, err
	}
	return status, nil
}

// testAccExistsDataApiResource asserts server status matches state; a bare GET would not.
func testAccExistsDataApiResource(t *testing.T, resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		if rawState == nil {
			return fmt.Errorf("resource %s not found in state", resourceReference)
		}

		data := newTestClient(t)
		status, err := retrieveDataApiStatusFromServer(
			data, rawState["organization_id"], rawState["project_id"], rawState["cluster_id"],
		)
		if err != nil {
			return err
		}

		for _, check := range []struct {
			attr   string
			server string
		}{
			{"enable_data_api", strconv.FormatBool(status.Enabled)},
			{"enable_network_peering", strconv.FormatBool(status.EnabledForNetworkPeering)},
			{"state_for_data_api", string(status.State)},
			{"state_for_network_peering", string(status.StateForNetworkPeering)},
			{"connection_string", status.ConnectionString},
		} {
			if rawState[check.attr] != check.server {
				return fmt.Errorf(
					"state disagrees with the server for %s: terraform has %q, API reports %q",
					check.attr, rawState[check.attr], check.server,
				)
			}
		}

		return nil
	}
}

// malformedDataApiImportId formats an invalid import string from cluster/project/org IDs.
func malformedDataApiImportId(resourceReference, format string) resource.ImportStateIdFunc {
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
			format,
			rawState["cluster_id"],
			rawState["project_id"],
			rawState["organization_id"],
		), nil
	}
}

// setDataApiOnServer updates the Data API outside Terraform, waiting either side of the PUT
// because the API 422s any write mid-transition.
func setDataApiOnServer(data *providerschema.Data, organizationId, projectId, clusterId string, enableDataApi, enableNetworkPeering bool) error {
	if err := waitForDataApiTerminalState(data, organizationId, projectId, clusterId); err != nil {
		return fmt.Errorf("cluster was not settled before the update: %w", err)
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI",
		data.HostURL, organizationId, projectId, clusterId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusAccepted}
	if _, err := data.ClientV1.ExecuteWithRetry(
		context.Background(),
		cfg,
		data_api.UpdateDataApiRequest{
			EnableDataApi:        enableDataApi,
			EnableNetworkPeering: enableNetworkPeering,
		},
		data.Token,
		nil,
	); err != nil {
		return err
	}

	return waitForDataApiTerminalState(data, organizationId, projectId, clusterId)
}

// TestAccDataApiResourceDestroyIsNoOp pins the no-op destroy, ending enabled so it proves something.
func TestAccDataApiResourceDestroyIsNoOp(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_data_api_noop_destroy_")
	resourceReference := "couchbase-capella_data_api." + resourceName

	// Not ParallelTest: mutates the Data API on the shared cluster.
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckDataApiDestroyIsNoOp(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataApiResourceConfig(resourceName, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDataApiResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "enable_data_api", "true"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_data_api", "enabled"),
				),
			},
		},
	})
}

func testAccCheckDataApiDestroyIsNoOp(t *testing.T) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		data := newTestClient(t)

		status, err := retrieveDataApiStatusFromServer(data, globalOrgId, globalProjectId, globalClusterId)
		if err != nil {
			return fmt.Errorf("failed to read Data API status after destroy: %w", err)
		}

		if !status.Enabled {
			return fmt.Errorf(
				"destroy changed the cluster: Data API is %q after destroy, expected it to remain enabled",
				status.State,
			)
		}

		// Delete is a no-op, so nothing else restores the cluster for later tests.
		if err := setDataApiOnServer(data, globalOrgId, globalProjectId, globalClusterId, false, false); err != nil {
			return fmt.Errorf("failed to disable Data API during cleanup: %w", err)
		}

		return nil
	}
}

// waitForDataApiTerminalState polls until the Data API and its peering both settle.
func waitForDataApiTerminalState(data *providerschema.Data, organizationId, projectId, clusterId string) error {
	const (
		timeout  = 20 * time.Minute
		interval = 15 * time.Second
	)

	deadline := time.Now().Add(timeout)
	for {
		status, err := retrieveDataApiStatusFromServer(data, organizationId, projectId, clusterId)
		if err != nil {
			return fmt.Errorf("failed to poll Data API status: %w", err)
		}
		if data_api.IsFinalState(status.State) && data_api.IsFinalState(status.StateForNetworkPeering) {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf(
				"Data API did not settle within %s: state=%q state_for_network_peering=%q",
				timeout, status.State, status.StateForNetworkPeering,
			)
		}
		time.Sleep(interval)
	}
}

// TestAccDataApiDriftDetection covers Read against a change Terraform did not make.
func TestAccDataApiDriftDetection(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_data_api_drift_")
	resourceReference := "couchbase-capella_data_api." + resourceName

	// Not ParallelTest: mutates the Data API on the shared cluster.
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCleanupDataApi(t),
		Steps: []resource.TestStep{
			// Baseline. The wait stops a preceding test's in-flight transition 422ing this apply.
			{
				PreConfig: func() {
					if err := waitForDataApiTerminalState(newTestClient(t), globalOrgId, globalProjectId, globalClusterId); err != nil {
						t.Fatalf("cluster was not settled before the test started: %v", err)
					}
				},
				Config: testAccDataApiResourceConfig(resourceName, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDataApiResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "false"),
				),
			},
			// Enable peering behind Terraform's back; refresh must then plan a return to config.
			{
				PreConfig: func() {
					if err := setDataApiOnServer(newTestClient(t), globalOrgId, globalProjectId, globalClusterId, true, true); err != nil {
						t.Fatalf("failed to enable network peering out of band: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Refresh records what the cluster reports, not what the config asks for.
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "true"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_network_peering", "enabled"),
				),
			},
			// Re-applying the unchanged config must reconcile the drift.
			{
				Config: testAccDataApiResourceConfig(resourceName, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsDataApiResource(t, resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "enable_network_peering", "false"),
					resource.TestCheckResourceAttr(resourceReference, "state_for_network_peering", "disabled"),
					resource.TestCheckResourceAttr(resourceReference, "enable_data_api", "true"),
				),
			},
		},
	})
}

// testAccCleanupDataApi disables the Data API; the framework cannot, as Delete is a no-op.
func testAccCleanupDataApi(t *testing.T) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		data := newTestClient(t)
		if err := setDataApiOnServer(data, globalOrgId, globalProjectId, globalClusterId, false, false); err != nil {
			return fmt.Errorf("failed to disable Data API during cleanup: %w", err)
		}
		return nil
	}
}

// TestAccDataApiResourceMissingRequiredAttributes covers required attributes; no API calls.
func TestAccDataApiResourceMissingRequiredAttributes(t *testing.T) {
	validUUID := "ffffffff-aaaa-1414-eeee-000000000000"

	tests := []struct {
		name        string
		config      string
		expectError string
	}{
		{
			name: "resource missing enable_data_api",
			config: fmt.Sprintf(`
resource "couchbase-capella_data_api" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[2]s"
  cluster_id      = "%[2]s"
}
`, randomStringWithPrefix("tf_acc_data_api_missing_flag_"), validUUID),
			expectError: `(?s)The\s+argument\s+"enable_data_api"\s+is\s+required`,
		},
		{
			name: "resource missing cluster_id",
			config: fmt.Sprintf(`
resource "couchbase-capella_data_api" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[2]s"
  enable_data_api = true
}
`, randomStringWithPrefix("tf_acc_data_api_missing_cluster_"), validUUID),
			expectError: `(?s)The\s+argument\s+"cluster_id"\s+is\s+required`,
		},
		{
			name: "data source missing cluster_id",
			config: fmt.Sprintf(`
data "couchbase-capella_data_api" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[2]s"
}
`, randomStringWithPrefix("tf_acc_data_api_ds_missing_cluster_"), validUUID),
			expectError: `(?s)The\s+argument\s+"cluster_id"\s+is\s+required`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config:      globalProviderBlock + test.config,
						ExpectError: regexp.MustCompile(test.expectError),
					},
				},
			})
		})
	}
}

// TestAccDataApiResourceEmptyAndWhitespaceIDs covers both ID validators; "   " proves no trimming.
func TestAccDataApiResourceEmptyAndWhitespaceIDs(t *testing.T) {
	validUUID := "ffffffff-aaaa-1414-eeee-000000000000"

	tests := []struct {
		name        string
		clusterID   string
		expectError string
	}{
		{
			name:        "empty cluster_id fails the length validator",
			clusterID:   "",
			expectError: `(?s)Invalid\s+Attribute\s+Value\s+Length.*cluster_id.*at\s+least\s+1`,
		},
		{
			name:        "whitespace cluster_id fails only the UUID validator",
			clusterID:   "   ",
			expectError: `(?s)Invalid\s+Attribute\s+Value\s+Match.*cluster_id.*must\s+be\s+a\s+valid\s+UUID`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resourceName := randomStringWithPrefix("tf_acc_data_api_bad_id_")

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: testAccDataApiResourceIDsConfig(
							resourceName, validUUID, validUUID, test.clusterID),
						ExpectError: regexp.MustCompile(test.expectError),
					},
				},
			})
		})
	}
}

// TestAccDataApiResourceNonExistentCluster checks an unknown UUID fails fast, not after 30 min.
func TestAccDataApiResourceNonExistentCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_data_api_ghost_")
	const nonExistentClusterID = "deadbeef-dead-beef-dead-beefdeadbeef"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApiResourceIDsConfig(
					resourceName, globalOrgId, globalProjectId, nonExistentClusterID),
				ExpectError: regexp.MustCompile(`(?s)Error\s+creating\s+Data\s+API\s+configuration`),
			},
		},
	})
}

// TestAccDatasourceDataApiNonExistentCluster matches the summary only, so a better 404 still passes.
func TestAccDatasourceDataApiNonExistentCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_data_api_ds_ghost_")
	const nonExistentClusterID = "deadbeef-dead-beef-dead-beefdeadbeef"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApiDatasourceIDsConfig(
					dsName, globalOrgId, globalProjectId, nonExistentClusterID),
				ExpectError: regexp.MustCompile(`(?s)Error\s+Reading\s+Capella\s+Data\s+API\s+Status`),
			},
		},
	})
}

// TestAccDataApiResourceComputedFlagsRejectInvalidCombination covers flags unknown at plan time;
// the provider's own message, not an API error, proves nothing reached Capella.
func TestAccDataApiResourceComputedFlagsRejectInvalidCombination(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_data_api_computed_flags_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

resource "terraform_data" "flags_%[2]s" {
  input = {
    enable_data_api        = false
    enable_network_peering = true
  }
}

resource "couchbase-capella_data_api" "%[2]s" {
  organization_id        = "%[3]s"
  project_id             = "%[4]s"
  cluster_id             = "%[5]s"
  enable_data_api        = terraform_data.flags_%[2]s.output.enable_data_api
  enable_network_peering = terraform_data.flags_%[2]s.output.enable_network_peering
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId),
				ExpectError: regexp.MustCompile(`(?s)Network\s+peering\s+cannot\s+be\s+enabled`),
			},
		},
	})
}
