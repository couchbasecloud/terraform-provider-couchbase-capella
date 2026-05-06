package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAppEndpointAccessControlFunction exercises CRUD and import for the
// couchbase-capella_app_endpoint_access_control_function resource.
//
// Uses a pre-created, dedicated endpoint (ensureACFEndpoint) so it can run in
// parallel with other tests without competing for creation slots.
//
// Step 1 also verifies that scope and collection default to "_default" when
// omitted from configuration.
func TestAccAppEndpointAccessControlFunction(t *testing.T) {
	ensureACFEndpoint(t)

	resourceName := randomStringWithPrefix("tf_acc_acf_")
	resourceReference := "couchbase-capella_app_endpoint_access_control_function." + resourceName

	initialFn := "function(doc, oldDoc, meta) { channel(doc.channels); }"
	updatedFn := "function(doc, oldDoc, meta) { channel(doc.type); }"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				// scope and collection omitted — verify they default to "_default"
				Config: testAccACFConfigNoScope(resourceName, globalACFEndpointName, initialFn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalACFEndpointName),
					resource.TestCheckResourceAttr(resourceReference, "scope", "_default"),
					resource.TestCheckResourceAttr(resourceReference, "collection", "_default"),
					resource.TestCheckResourceAttr(resourceReference, "access_control_function", initialFn),
				),
			},
			{
				Config: testAccACFConfig(resourceName, globalACFEndpointName, "_default", "_default", updatedFn),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access_control_function", updatedFn),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateACFImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "app_endpoint_name",
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Config helpers
// ─────────────────────────────────────────────────────────────────────────────

func testAccACFConfig(resourceName, epName, scope, collection, acfBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint_access_control_function" "%[2]s" {
	organization_id         = "%[3]s"
	project_id              = "%[4]s"
	cluster_id              = "%[5]s"
	app_service_id          = "%[6]s"
	app_endpoint_name       = "%[7]s"
	scope                   = "%[8]s"
	collection              = "%[9]s"
	access_control_function = "%[10]s"
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		epName,
		scope,
		collection,
		acfBody,
	)
}

func testAccACFConfigNoScope(resourceName, epName, acfBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint_access_control_function" "%[2]s" {
	organization_id         = "%[3]s"
	project_id              = "%[4]s"
	cluster_id              = "%[5]s"
	app_service_id          = "%[6]s"
	app_endpoint_name       = "%[7]s"
	access_control_function = "%[8]s"
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		epName,
		acfBody,
	)
}

func generateACFImportId(resourceReference string) resource.ImportStateIdFunc {
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
			"organization_id=%s,project_id=%s,cluster_id=%s,app_service_id=%s,app_endpoint_name=%s,scope_name=%s,collection_name=%s",
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"],
			rawState["app_service_id"],
			rawState["app_endpoint_name"],
			rawState["scope"],
			rawState["collection"],
		), nil
	}
}
