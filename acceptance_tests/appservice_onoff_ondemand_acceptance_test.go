package acceptance_tests

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
)

func TestAccAppServiceOnOffOnDemandResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_service_onoff_ondemand_")
	resourceReference := "couchbase-capella_app_service_onoff_ondemand." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceOnOffOnDemandResourceConfig(resourceName, globalAppServiceId, "on"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "state", "on"),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateAppServiceOnOffOnDemandImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "app_service_id",
			},
		},
	})
}

func TestAccAppServiceOnOffOnDemandAlreadyInRequestedState(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_service_onoff_already_on_")
	resourceReference := "couchbase-capella_app_service_onoff_ondemand." + resourceName
	cfg := testAccAppServiceOnOffOnDemandResourceConfig(resourceName, globalAppServiceId, "on")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		// The fixture app service has to be on already, otherwise the first apply performs a
		// real activation and never reaches the redundant-request path under test. TestMain
		// waits for Healthy only when it creates the app service itself; one supplied through
		// TF_VAR_app_service_id is adopted as-is, so its state is asserted here rather than
		// assumed. Failing loudly beats passing without covering the fix.
		PreCheck: func() {
			currentState, err := appServiceCurrentState(context.Background(), globalClient)
			if err != nil {
				t.Fatalf("could not read state of app service %s: %s", globalAppServiceId, err)
			}
			if currentState != appservice.Healthy {
				t.Fatalf(
					"app service %s must be %q so that requesting state \"on\" is redundant, got %q",
					globalAppServiceId, appservice.Healthy, currentState,
				)
			}
		},
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "state", "on"),
				),
			},
			// Tainting forces a second Create against the app service, which is still on
			// because Delete is a no-op for this resource. That exercises the redundant
			// activation request independently of what the first apply happened to do.
			{
				Config: cfg,
				Taint:  []string{resourceReference},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "state", "on"),
				),
			},
			// Re-apply the same config; expect no changes (no perpetual drift).
			{
				Config:             cfg,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAppServiceOnOffOnDemandResourceInvalidState(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_service_onoff_invalid_state_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServiceOnOffOnDemandResourceConfig(resourceName, globalAppServiceId, "frozen"),
				ExpectError: regexp.MustCompile(`(?s)state must be either 'on' or 'off'|App Service activation failed`),
			},
		},
	})
}

func TestAccAppServiceOnOffOnDemandResourceInvalidAppService(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_service_onoff_invalid_appsvc_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServiceOnOffOnDemandResourceConfig(resourceName, "33333333-3333-3333-3333-333333333333", "on"),
				ExpectError: regexp.MustCompile(`(?s)App Service activation failed|switching on/off|not found|access to the requested resource is denied`),
			},
		},
	})
}

func TestAccAppServiceOnOffOnDemandResourceInvalidUUIDs(t *testing.T) {
	tests := []struct {
		name           string
		organizationID string
		projectID      string
		clusterID      string
		appServiceID   string
	}{
		{
			name:           "organization_id",
			organizationID: "not-a-uuid",
			projectID:      "11111111-1111-1111-1111-111111111111",
			clusterID:      "22222222-2222-2222-2222-222222222222",
			appServiceID:   "33333333-3333-3333-3333-333333333333",
		},
		{
			name:           "project_id",
			organizationID: "00000000-0000-0000-0000-000000000000",
			projectID:      "not-a-uuid",
			clusterID:      "22222222-2222-2222-2222-222222222222",
			appServiceID:   "33333333-3333-3333-3333-333333333333",
		},
		{
			name:           "cluster_id",
			organizationID: "00000000-0000-0000-0000-000000000000",
			projectID:      "11111111-1111-1111-1111-111111111111",
			clusterID:      "not-a-uuid",
			appServiceID:   "33333333-3333-3333-3333-333333333333",
		},
		{
			name:           "app_service_id",
			organizationID: "00000000-0000-0000-0000-000000000000",
			projectID:      "11111111-1111-1111-1111-111111111111",
			clusterID:      "22222222-2222-2222-2222-222222222222",
			appServiceID:   "not-a-uuid",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			resourceName := randomStringWithPrefix("tf_acc_app_service_onoff_invalid_uuid_")
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: testAccAppServiceOnOffOnDemandResourceConfigWithIDs(
							resourceName,
							test.organizationID,
							test.projectID,
							test.clusterID,
							test.appServiceID,
							"on",
						),
						ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.*` + test.name + `.*must be a valid UUID`),
					},
				},
			})
		})
	}
}

func TestAccAppServiceOnOffOnDemandResourceMissingRequiredFields(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "app_service_id",
			body: fmt.Sprintf(`
  organization_id = "%s"
  project_id      = "%s"
  cluster_id      = "%s"
  state           = "on"
`, globalOrgId, globalProjectId, globalClusterId),
		},
		{
			name: "state",
			body: fmt.Sprintf(`
  organization_id = "%s"
  project_id      = "%s"
  cluster_id      = "%s"
  app_service_id  = "%s"
`, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId),
		},
		{
			name: "cluster_id",
			body: fmt.Sprintf(`
  organization_id = "%s"
  project_id      = "%s"
  app_service_id  = "%s"
  state           = "on"
`, globalOrgId, globalProjectId, globalAppServiceId),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			resourceName := randomStringWithPrefix("tf_acc_app_service_onoff_missing_")
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_onoff_ondemand" "%[2]s" {%[3]s}
`, globalProviderBlock, resourceName, test.body),
						ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument "%s" is required`, test.name)),
					},
				},
			})
		})
	}
}

func testAccAppServiceOnOffOnDemandResourceConfig(resourceName, appServiceID, state string) string {
	return testAccAppServiceOnOffOnDemandResourceConfigWithIDs(
		resourceName, globalOrgId, globalProjectId, globalClusterId, appServiceID, state,
	)
}

func testAccAppServiceOnOffOnDemandResourceConfigWithIDs(resourceName, organizationID, projectID, clusterID, appServiceID, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_onoff_ondemand" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  state           = "%[7]s"
}
`, globalProviderBlock, resourceName, organizationID, projectID, clusterID, appServiceID, state)
}

func generateAppServiceOnOffOnDemandImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		res, ok := state.RootModule().Resources[resourceReference]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		attrs := res.Primary.Attributes
		return fmt.Sprintf(
			"app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			attrs["app_service_id"],
			attrs["cluster_id"],
			attrs["project_id"],
			attrs["organization_id"],
		), nil
	}
}
