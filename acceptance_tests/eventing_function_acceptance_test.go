package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// eventingFunctionCode is a minimal valid eventing function handler reused across tests.
const eventingFunctionCode = `function OnUpdate(doc, meta) {\n  log(\"updated\", meta.id);\n}\n`

func TestAccEventingFunctionResourceRequiredOnly(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_req_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigRequiredOnly(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(funcReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(funcReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(funcReference, "name", funcName),
					resource.TestCheckResourceAttrSet(funcReference, "code"),
					resource.TestCheckResourceAttr(funcReference, "event_source.bucket", globalBucketName),
					resource.TestCheckResourceAttr(funcReference, "event_source.scope", globalScopeName),
					resource.TestCheckResourceAttr(funcReference, "event_source.collection", globalCollectionName),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.bucket", globalMetadataBucketName),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.collection", globalCollectionName),
					resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),
					// settings are optional+computed: server defaults are populated even when omitted.
					resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "1"),
					resource.TestCheckResourceAttr(funcReference, "settings.script_timeout", "60"),
					resource.TestCheckResourceAttr(funcReference, "settings.sql_consistency", "none"),
					resource.TestCheckResourceAttr(funcReference, "settings.language_compatibility", "7.2.0"),
					resource.TestCheckResourceAttr(funcReference, "settings.feed_boundary", "from_now"),
					resource.TestCheckResourceAttr(funcReference, "settings.max_timer_context_size", "1024"),
					resource.TestCheckResourceAttr(funcReference, "settings.allow_sync_documents", "true"),
					resource.TestCheckResourceAttr(funcReference, "settings.cursor_aware", "false"),
				),
			},
			{
				ResourceName:                         funcReference,
				ImportState:                          true,
				ImportStateIdFunc:                    generateEventingFunctionImportIdForResource(funcReference),
				ImportStateVerifyIdentifierAttribute: "name",
			},
			{
				// Updates are only applied while undeployed/paused, so the function stays undeployed here.
				Config: testAccEventingFunctionResourceConfigUpdatable(funcName, 3, 90),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "3"),
					resource.TestCheckResourceAttr(funcReference, "settings.script_timeout", "90"),
				),
			},
		},
	})
}

func testAccEventingFunctionResourceConfigRequiredOnly(funcName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, funcName, eventingFunctionCode, globalBucketName, globalMetadataBucketName, globalScopeName, globalCollectionName)
}

func testAccEventingFunctionResourceConfigUpdatable(funcName string, workerCount, scriptTimeout int) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  settings = {
    worker_count   = %[11]d
    script_timeout = %[12]d
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName,
		workerCount,
		scriptTimeout)
}

func generateEventingFunctionImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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
			"function_name=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["name"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

func TestAccDatasourceEventingFunction(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_ds_evt_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName
	dataSourceReference := "data.couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionDataSourceConfig(funcName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "name", funcName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.bucket", globalBucketName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.scope", globalScopeName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.collection", globalCollectionName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.bucket", globalMetadataBucketName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.scope", globalScopeName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.collection", globalCollectionName),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.worker_count", "1"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.script_timeout", "60"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.sql_consistency", "none"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.language_compatibility", "7.2.0"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.feed_boundary", "from_now"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.max_timer_context_size", "1024"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.allow_sync_documents", "true"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.cursor_aware", "false"),
					resource.TestCheckResourceAttrPair(dataSourceReference, "code", funcReference, "code"),
					// export defaults to false, so the read-only fields (currently just status) is returned
					resource.TestCheckResourceAttr(dataSourceReference, "status", "undeployed"),
				),
			},
		},
	})
}

// TestAccDatasourceEventingFunctionExport verifies that setting export = true on the
// couchbase-capella_eventing_function data source omits the read-only fields (just status currently).
func TestAccDatasourceEventingFunctionExport(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_ds_evt_exp_fn_")
	dataSourceReference := "data.couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionDataSourceConfig(funcName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "name", funcName),
					resource.TestCheckResourceAttr(dataSourceReference, "export", "true"),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.bucket", globalBucketName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.scope", globalScopeName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.collection", globalCollectionName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.bucket", globalMetadataBucketName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.scope", globalScopeName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.collection", globalCollectionName),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.worker_count", "1"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.script_timeout", "60"),
					resource.TestCheckResourceAttrSet(dataSourceReference, "code"),
					// export omits the read-only status field from the response, so check that is omitted
					resource.TestCheckNoResourceAttr(dataSourceReference, "status"),
				),
			},
		},
	})
}

// testAccEventingFunctionDataSourceConfig creates an eventing function and a single-function data
// source that reads it back.
func testAccEventingFunctionDataSourceConfig(funcName string, export bool) string {
	exportLine := ""
	if export {
		exportLine = "\n  export          = true"
	}
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }
}

data "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = couchbase-capella_eventing_function.%[5]s.name%[11]s
}
`,
		globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName,
		exportLine,
	)
}
