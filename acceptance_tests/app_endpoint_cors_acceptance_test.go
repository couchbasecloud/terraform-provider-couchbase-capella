package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAppEndpointCorsResource exercises CRUD and import for the
// couchbase-capella_app_endpoint_cors resource against the common pre-created
// endpoint.
//
// Runs sequentially (resource.Test, not ParallelTest) to avoid a race
// condition with TestAccAppEndpointCorsResourceOriginOnly: both resources
// write CORS configuration to the same endpoint, and concurrent updates
// cause the post-apply refresh plan to become non-empty.
func TestAccAppEndpointCorsResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cors_")
	resourceReference := "couchbase-capella_app_endpoint_cors." + resourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCorsResourceOriginOnlyConfig(resourceName, globalAppEndpointName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalAppEndpointName),
					resource.TestCheckResourceAttr(resourceReference, "origin.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "origin.*", "*"),
				),
			},
			{
				Config: testAccCorsResourceConfig(resourceName, globalAppEndpointName, `["*"]`, `["*"]`, `["Authorization", "Content-Type"]`, 3600, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "max_age", "3600"),
					resource.TestCheckResourceAttr(resourceReference, "origin.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "origin.*", "*"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateCorsResourceImportId(resourceReference),
				ImportState:       true,
				// ImportStateVerify omitted: resource has no "id" attribute.
			},
		},
	})
}

// TestAccAppEndpointCorsResourceOriginOnly verifies that the CORS resource
// works correctly when only the required origin field is set.
func TestAccAppEndpointCorsResourceOriginOnly(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cors_")
	resourceReference := "couchbase-capella_app_endpoint_cors." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCorsResourceOriginOnlyConfig(resourceName, globalAppEndpointName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalAppEndpointName),
					resource.TestCheckResourceAttr(resourceReference, "origin.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "origin.*", "*"),
				),
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Config helpers
// ─────────────────────────────────────────────────────────────────────────────

func testAccCorsResourceConfig(resourceName, endpointName, origin, loginOrigin, headers string, maxAge int, disabled bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint_cors" "%[2]s" {
  organization_id   = "%[3]s"
  project_id        = "%[4]s"
  cluster_id        = "%[5]s"
  app_service_id    = "%[6]s"
  app_endpoint_name = "%[7]s"
  origin            = %[8]s
  login_origin      = %[9]s
  headers           = %[10]s
  max_age           = %[11]d
  disabled          = %[12]t
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		endpointName,
		origin,
		loginOrigin,
		headers,
		maxAge,
		disabled,
	)
}

func testAccCorsResourceOriginOnlyConfig(resourceName, endpointName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint_cors" "%[2]s" {
  organization_id   = "%[3]s"
  project_id        = "%[4]s"
  cluster_id        = "%[5]s"
  app_service_id    = "%[6]s"
  app_endpoint_name = "%[7]s"
  origin            = ["*"]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		endpointName,
	)
}

func generateCorsResourceImportId(resourceReference string) resource.ImportStateIdFunc {
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
			"organization_id=%s,project_id=%s,cluster_id=%s,app_service_id=%s,app_endpoint_name=%s",
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"],
			rawState["app_service_id"],
			rawState["app_endpoint_name"],
		), nil
	}
}
