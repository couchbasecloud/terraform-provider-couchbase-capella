package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAppEndpointImportFilter exercises CRUD and import for the
// couchbase-capella_app_endpoint_import_filter resource.
//
// Uses a pre-created, dedicated endpoint (ensureIFEndpoint) so it can run in
// parallel with other tests without competing for creation slots.
//
// Step 1 also verifies that scope and collection default to "_default" when
// omitted from configuration.
func TestAccAppEndpointImportFilter(t *testing.T) {
	ensureIFEndpoint(t)

	resourceName := randomStringWithPrefix("tf_acc_if_")
	resourceReference := "couchbase-capella_app_endpoint_import_filter." + resourceName

	initialFilter := "function(doc) { return true; }"
	updatedFilter := "function(doc) { return doc.type === 'user'; }"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				// scope and collection omitted — verify they default to "_default"
				Config: testAccImportFilterConfigNoScope(resourceName, globalIFEndpointName, initialFilter),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalIFEndpointName),
					resource.TestCheckResourceAttr(resourceReference, "scope", "_default"),
					resource.TestCheckResourceAttr(resourceReference, "collection", "_default"),
					resource.TestCheckResourceAttr(resourceReference, "import_filter", initialFilter),
				),
			},
			{
				Config: testAccImportFilterConfig(resourceName, globalIFEndpointName, "_default", "_default", updatedFilter),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "import_filter", updatedFilter),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateImportFilterImportId(resourceReference),
				ImportState:       true,
				// ImportStateVerify cannot be used here: ImportStatePassthroughID
				// stores the full composite ID in app_endpoint_name, but Read
				// normalises it to just the endpoint name. The verifier then
				// fails to locate the resource by the original composite value.
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Config helpers
// ─────────────────────────────────────────────────────────────────────────────

func testAccImportFilterConfig(resourceName, epName, scope, collection, filterBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint_import_filter" "%[2]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[6]s"
	app_endpoint_name = "%[7]s"
	scope             = "%[8]s"
	collection        = "%[9]s"
	import_filter     = "%[10]s"
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
		filterBody,
	)
}

func testAccImportFilterConfigNoScope(resourceName, epName, filterBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint_import_filter" "%[2]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[6]s"
	app_endpoint_name = "%[7]s"
	import_filter     = "%[8]s"
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		epName,
		filterBody,
	)
}

func generateImportFilterImportId(resourceReference string) resource.ImportStateIdFunc {
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
