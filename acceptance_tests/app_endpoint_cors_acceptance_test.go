package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAppEndpointCorsResource exercises CRUD and import for the
// couchbase-capella_app_endpoint_cors resource.
//
// Uses a pre-created, dedicated endpoint (ensureCORSEndpoint) so it can run in
// parallel with other tests without competing for creation slots.
func TestAccAppEndpointCorsResource(t *testing.T) {
	ensureCORSEndpoint(t)

	resourceName := randomStringWithPrefix("tf_acc_cors_")
	resourceReference := "couchbase-capella_app_endpoint_cors." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCorsResourceOriginOnlyConfig(resourceName, globalCORSEndpointName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", appEndpointClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", appEndpointAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalCORSEndpointName),
					resource.TestCheckResourceAttr(resourceReference, "origin.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "origin.*", "*"),
				),
			},
			{
				Config: testAccCorsResourceConfig(resourceName, globalCORSEndpointName, `["*"]`, `["*"]`, `["Authorization", "Content-Type"]`, 3600, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "max_age", "3600"),
					resource.TestCheckResourceAttr(resourceReference, "origin.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "origin.*", "*"),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateCorsResourceImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "app_endpoint_name",
			},
		},
	})
}

// TestAccAppEndpointCorsResourceOriginOnly verifies that the CORS resource
// works correctly when only the required origin field is set.
//
// Uses its own pre-created endpoint (ensureCORSOriginOnlyEndpoint) so it can
// run in parallel with TestAccAppEndpointCorsResource without state conflicts.
func TestAccAppEndpointCorsResourceOriginOnly(t *testing.T) {
	ensureCORSOriginOnlyEndpoint(t)

	resourceName := randomStringWithPrefix("tf_acc_cors_")
	resourceReference := "couchbase-capella_app_endpoint_cors." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCorsResourceOriginOnlyConfig(resourceName, globalCORSOriginOnlyEndpointName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalCORSOriginOnlyEndpointName),
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

func testAccCorsResourceConfig(resourceName, epName, origin, loginOrigin, headers string, maxAge int, disabled bool) string {
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
		appEndpointClusterId,
		appEndpointAppServiceId,
		epName,
		origin,
		loginOrigin,
		headers,
		maxAge,
		disabled,
	)
}

func testAccCorsResourceOriginOnlyConfig(resourceName, epName string) string {
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
		appEndpointClusterId,
		appEndpointAppServiceId,
		epName,
	)
}

func generateCorsResourceImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if v, ok := m.Resources[resourceReference]; ok {
				rawState = v.Primary.Attributes
				break
			}
		}
		if rawState == nil {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
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
