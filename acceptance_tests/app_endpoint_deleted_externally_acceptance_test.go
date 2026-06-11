package acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAppEndpointDeletedExternally verifies that when an App Endpoint is
// deleted outside of Terraform (e.g. via the UI or management API), the
// provider gracefully removes the resource from state during the Read operation
// so that Terraform recreates it on the next apply.
//
// This tests the 403 → List check → remove-from-state flow implemented in
// checkAppEndpointDeletedOrForbidden.
func TestAccAppEndpointDeletedExternally(t *testing.T) {
	ensureFixtureBucketByName(t, globalDeletedExternallyEPBucketName)

	resourceName := randomStringWithPrefix("tf_acc_ep_del_ext_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_ep_del_ext_")

	cfg := testAccAppEndpointDeletedExternallyConfig(resourceName, epName, globalDeletedExternallyEPBucketName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: Create the App Endpoint via Terraform.
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", appEndpointClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", appEndpointAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "name", epName),
					resource.TestCheckResourceAttr(resourceReference, "bucket", globalDeletedExternallyEPBucketName),
				),
			},
			// Step 2: Delete the App Endpoint externally, then re-apply the same config.
			// The provider should detect the deletion via the 403 → List fallback,
			// remove the resource from state, and recreate it.
			{
				PreConfig: func() {
					ctx := context.Background()
					err := deleteAppEndpointFixtureEndpoint(
						ctx,
						globalClient,
						globalProjectId,
						appEndpointClusterId,
						appEndpointAppServiceId,
						epName,
					)
					if err != nil {
						t.Fatalf("failed to delete app endpoint externally: %v", err)
					}
				},
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", epName),
					resource.TestCheckResourceAttr(resourceReference, "bucket", globalDeletedExternallyEPBucketName),
				),
			},
		},
	})
}

// TestAccAppEndpointAccessControlFunctionDeletedExternally verifies that when
// an App Endpoint is deleted outside of Terraform, the access_control_function
// resource attached to it is gracefully removed from state during Read.
func TestAccAppEndpointAccessControlFunctionDeletedExternally(t *testing.T) {
	ensureFixtureBucketByName(t, globalACFDeletedExtEPBucketName)

	epResourceName := randomStringWithPrefix("tf_acc_ep_acf_del_")
	epReference := "couchbase-capella_app_endpoint." + epResourceName
	epName := randomStringWithPrefix("tf_acc_ep_acf_del_")

	acfResourceName := randomStringWithPrefix("tf_acc_acf_del_")
	acfReference := "couchbase-capella_app_endpoint_access_control_function." + acfResourceName

	acfBody := "function(doc, oldDoc, meta) { channel(doc.channels); }"

	cfgWithACF := testAccAppEndpointWithACFDeletedExternallyConfig(
		epResourceName, epName, globalACFDeletedExtEPBucketName,
		acfResourceName, acfBody,
	)
	cfgEndpointOnly := testAccAppEndpointDeletedExternallyConfig(epResourceName, epName, globalACFDeletedExtEPBucketName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: Create App Endpoint + ACF resource.
			{
				Config: cfgWithACF,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(epReference),
					resource.TestCheckResourceAttr(acfReference, "app_endpoint_name", epName),
					resource.TestCheckResourceAttr(acfReference, "access_control_function", acfBody),
				),
			},
			// Step 2: Delete the App Endpoint externally, then apply config without ACF.
			// The ACF resource should be removed from state automatically during Read
			// because the parent App Endpoint no longer exists.
			{
				PreConfig: func() {
					ctx := context.Background()
					err := deleteAppEndpointFixtureEndpoint(
						ctx,
						globalClient,
						globalProjectId,
						appEndpointClusterId,
						appEndpointAppServiceId,
						epName,
					)
					if err != nil {
						t.Fatalf("failed to delete app endpoint externally: %v", err)
					}
				},
				Config: cfgEndpointOnly,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(epReference),
					resource.TestCheckResourceAttr(epReference, "name", epName),
				),
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Config helpers
// ─────────────────────────────────────────────────────────────────────────────

func testAccAppEndpointDeletedExternallyConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		bucketName,
		endpointName,
	)
}

func testAccAppEndpointWithACFDeletedExternallyConfig(
	epResourceName, endpointName, bucketName string,
	acfResourceName, acfBody string,
) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
}

resource "couchbase-capella_app_endpoint_access_control_function" "%[9]s" {
	organization_id         = "%[3]s"
	project_id              = "%[4]s"
	cluster_id              = "%[5]s"
	app_service_id          = "%[6]s"
	app_endpoint_name       = couchbase-capella_app_endpoint.%[2]s.name
	access_control_function = "%[10]s"
}
`,
		globalProviderBlock,
		epResourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		bucketName,
		endpointName,
		acfResourceName,
		acfBody,
	)
}
