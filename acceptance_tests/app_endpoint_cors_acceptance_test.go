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
// Creates its own bucket and app endpoint so it can run in parallel with other
// tests without competing for the shared common endpoint's CORS state.
func TestAccAppEndpointCorsResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cors_")
	resourceReference := "couchbase-capella_app_endpoint_cors." + resourceName
	bucketName := randomStringWithPrefix("tf_acc_cors_bkt_")
	epName := randomStringWithPrefix("tf_acc_cors_ep_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCorsResourceOriginOnlyConfig(resourceName, bucketName, epName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", epName),
					resource.TestCheckResourceAttr(resourceReference, "origin.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "origin.*", "*"),
				),
			},
			{
				Config: testAccCorsResourceConfig(resourceName, bucketName, epName, `["*"]`, `["*"]`, `["Authorization", "Content-Type"]`, 3600, false),
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
func TestAccAppEndpointCorsResourceOriginOnly(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cors_")
	resourceReference := "couchbase-capella_app_endpoint_cors." + resourceName
	bucketName := randomStringWithPrefix("tf_acc_cors_bkt_")
	epName := randomStringWithPrefix("tf_acc_cors_ep_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCorsResourceOriginOnlyConfig(resourceName, bucketName, epName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", epName),
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

func testAccCorsResourceConfig(resourceName, bucketName, epName, origin, loginOrigin, headers string, maxAge int, disabled bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[6]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s_ep" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[7]s"
	bucket          = "%[6]s"
	name            = "%[8]s"
	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}

resource "couchbase-capella_app_endpoint_cors" "%[2]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[7]s"
	app_endpoint_name = "%[8]s"
	origin            = %[9]s
	login_origin      = %[10]s
	headers           = %[11]s
	max_age           = %[12]d
	disabled          = %[13]t
	depends_on        = [couchbase-capella_app_endpoint.%[2]s_ep]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		bucketName,
		globalAppServiceId,
		epName,
		origin,
		loginOrigin,
		headers,
		maxAge,
		disabled,
	)
}

func testAccCorsResourceOriginOnlyConfig(resourceName, bucketName, epName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[6]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s_ep" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[7]s"
	bucket          = "%[6]s"
	name            = "%[8]s"
	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}

resource "couchbase-capella_app_endpoint_cors" "%[2]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[7]s"
	app_endpoint_name = "%[8]s"
	origin            = ["*"]
	depends_on        = [couchbase-capella_app_endpoint.%[2]s_ep]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		bucketName,
		globalAppServiceId,
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
