package acceptance_tests

import (
	"fmt"
	"slices"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDatasourceEventingFunctions(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_ds_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	// allDsName does no status filtering (unfiltered)
	allDsName := randomStringWithPrefix("tf_acc_evt_fns_all_ds_")
	allDsReference := "data.couchbase-capella_eventing_functions." + allDsName

	// matchDsName filters by undeployed functions (which funcName is)
	matchDsName := randomStringWithPrefix("tf_acc_evt_fns_match_ds_")
	matchDsReference := "data.couchbase-capella_eventing_functions." + matchDsName

	// noMatchDsName filters by deploying which does not match the function made in this test
	noMatchDsName := randomStringWithPrefix("tf_acc_evt_fns_nomatch_ds_")
	noMatchDsReference := "data.couchbase-capella_eventing_functions." + noMatchDsName

	// expectedFunctionAttrs is the expected attribute set for the function in the all/match data sources.
	expectedFunctionAttrs := map[string]string{
		"name":                              funcName,
		"status":                            "undeployed",
		"event_source.bucket":               globalBucketName,
		"event_source.scope":                globalScopeName,
		"event_source.collection":           globalCollectionName,
		"event_metadata_storage.bucket":     globalMetadataBucketName,
		"event_metadata_storage.collection": globalCollectionName,
		"settings.worker_count":             "1",
		"settings.script_timeout":           "60",
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionsDataSourceConfig(funcName, allDsName, matchDsName, noMatchDsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Confirm the parent function created as expected.
					resource.TestCheckResourceAttr(funcReference, "name", funcName),

					resource.TestCheckResourceAttr(allDsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(allDsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(allDsReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrSet(allDsReference, "eventing_functions.#"),
					resource.TestCheckTypeSetElemNestedAttrs(allDsReference, "eventing_functions.*", expectedFunctionAttrs),

					// check status = ["undeployed"] filter includes the function and returns only undeployed entries
					resource.TestCheckTypeSetElemNestedAttrs(matchDsReference, "eventing_functions.*", expectedFunctionAttrs),
					testAccCheckEventingFunctionsAllHaveStatus(matchDsReference, "undeployed"),

					// check status = ["deploying"] filter excludes the undeployed function entirely.
					testAccCheckEventingFunctionsExcludes(noMatchDsReference, funcName),
				),
			},
		},
	})
}

func testAccEventingFunctionsDataSourceConfig(funcName, allDsName, matchDsName, noMatchDsName string) string {
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

data "couchbase-capella_eventing_functions" "%[11]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"

  depends_on = [couchbase-capella_eventing_function.%[5]s]
}

data "couchbase-capella_eventing_functions" "%[12]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  status          = ["undeployed"]

  depends_on = [couchbase-capella_eventing_function.%[5]s]
}

data "couchbase-capella_eventing_functions" "%[13]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  status          = ["deploying"]

  depends_on = [couchbase-capella_eventing_function.%[5]s]
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
		allDsName,
		matchDsName,
		noMatchDsName)
}

// TestAccDatasourceEventingFunctionsFullPayload (scenario 13): the list data source returns the full per-function payload (description, settings, bindings).
func TestAccDatasourceEventingFunctionsFullPayload(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_fns_full_fn_")
	dsName := randomStringWithPrefix("tf_acc_evt_fns_full_ds_")
	dsReference := "data.couchbase-capella_eventing_functions." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionsDataSourceConfigFull(funcName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsReference, "eventing_functions.#"),
					testAccCheckEventingFunctionInList(dsReference, funcName, map[string]string{
						"description":                             "Full payload in list.",
						"status":                                  "undeployed",
						"event_source.bucket":                     globalBucketName,
						"event_metadata_storage.bucket":           globalMetadataBucketName,
						"settings.worker_count":                   "2",
						"settings.feed_boundary":                  "everything",
						"settings.allow_sync_documents":           "false",
						"settings.cursor_aware":                   "true",
						"bindings.buckets.0.alias":                "dst_col",
						"bindings.buckets.0.permission":           "readWrite",
						"bindings.urls.0.alias":                   "apiBasic",
						"bindings.urls.0.authentication.type":     "basic",
						"bindings.urls.0.authentication.password": "*****",
						"bindings.constants.0.alias":              "maxRetries",
						"bindings.constants.0.value":              "3",
					}),
				),
			},
		},
	})
}

// testAccCheckEventingFunctionInList finds the list entry by name and asserts each expected attribute path.
func testAccCheckEventingFunctionInList(dsReference, funcName string, expected map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[dsReference]
		if !ok {
			return fmt.Errorf("data source %q not found in state", dsReference)
		}
		attrs := rs.Primary.Attributes

		count, err := strconv.Atoi(attrs["eventing_functions.#"])
		if err != nil {
			return fmt.Errorf("invalid eventing_functions.# on %q: %w", dsReference, err)
		}

		idx := -1
		for i := 0; i < count; i++ {
			if attrs[fmt.Sprintf("eventing_functions.%d.name", i)] == funcName {
				idx = i
				break
			}
		}
		if idx == -1 {
			return fmt.Errorf("eventing function %q not found in %s", funcName, dsReference)
		}

		for key, want := range expected {
			full := fmt.Sprintf("eventing_functions.%d.%s", idx, key)
			if got := attrs[full]; got != want {
				return fmt.Errorf("%s = %q, want %q", full, got, want)
			}
		}
		return nil
	}
}

func testAccEventingFunctionsDataSourceConfigFull(funcName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  description     = "Full payload in list."
  code            = "%[6]s"
  state           = "undeployed"

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
    worker_count         = 2
    feed_boundary        = "everything"
    allow_sync_documents = false
    cursor_aware         = true
  }

  bindings = {
    buckets = [
      {
        alias      = "dst_col"
        bucket     = "%[7]s"
        scope      = "%[9]s"
        collection = "%[10]s"
        permission = "readWrite"
      }
    ]

    urls = [
      {
        alias                    = "apiBasic"
        url                      = "https://example.com/basic"
        allow_cookies            = true
        validate_tls_certificate = false
        authentication = {
          type     = "basic"
          username = "svc_user"
          password = "svc_pass"
        }
      }
    ]

    constants = [
      {
        alias = "maxRetries"
        value = "3"
      }
    ]
  }
}

data "couchbase-capella_eventing_functions" "%[11]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"

  depends_on = [couchbase-capella_eventing_function.%[5]s]
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
		dsName)
}

// TestAccDatasourceEventingFunctionsStateFilter (scenario 14): status filter across deployed/undeployed/paused plus a multi-value filter (paused created via deploy then pause).
func TestAccDatasourceEventingFunctionsStateFilter(t *testing.T) {
	deployedName := randomStringWithPrefix("tf_acc_evt_flt_dep_fn_")
	undeployedName := randomStringWithPrefix("tf_acc_evt_flt_undep_fn_")
	pausedName := randomStringWithPrefix("tf_acc_evt_flt_paused_fn_")

	deployedDs := randomStringWithPrefix("tf_acc_evt_flt_dep_ds_")
	pausedDs := randomStringWithPrefix("tf_acc_evt_flt_paused_ds_")
	undeployedDs := randomStringWithPrefix("tf_acc_evt_flt_undep_ds_")
	multiDs := randomStringWithPrefix("tf_acc_evt_flt_multi_ds_")

	deployedDsRef := "data.couchbase-capella_eventing_functions." + deployedDs
	pausedDsRef := "data.couchbase-capella_eventing_functions." + pausedDs
	undeployedDsRef := "data.couchbase-capella_eventing_functions." + undeployedDs
	multiDsRef := "data.couchbase-capella_eventing_functions." + multiDs

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Step 1: create all three functions; the paused one starts deployed so it can later be paused.
			{
				Config: testAccEventingFunctionsStateFilterConfig(
					deployedName, undeployedName, pausedName, "deployed",
					deployedDs, pausedDs, undeployedDs, multiDs),
			},
			// Step 2: move the paused function to paused, then assert each status filter.
			{
				Config: testAccEventingFunctionsStateFilterConfig(
					deployedName, undeployedName, pausedName, "paused",
					deployedDs, pausedDs, undeployedDs, multiDs),
				Check: resource.ComposeAggregateTestCheckFunc(
					// status=["deployed"]: includes the deployed function, all entries deployed, undeployed/paused absent.
					testAccCheckEventingFunctionInList(deployedDsRef, deployedName, map[string]string{"status": "deployed"}),
					testAccCheckEventingFunctionsAllHaveStatus(deployedDsRef, "deployed"),
					testAccCheckEventingFunctionsExcludes(deployedDsRef, undeployedName),
					testAccCheckEventingFunctionsExcludes(deployedDsRef, pausedName),

					// status = ["paused"]: includes the paused function, every entry is paused.
					testAccCheckEventingFunctionInList(pausedDsRef, pausedName, map[string]string{"status": "paused"}),
					testAccCheckEventingFunctionsAllHaveStatus(pausedDsRef, "paused"),
					testAccCheckEventingFunctionsExcludes(pausedDsRef, deployedName),
					testAccCheckEventingFunctionsExcludes(pausedDsRef, undeployedName),

					// status = ["undeployed"]: includes the undeployed function, every entry is undeployed.
					testAccCheckEventingFunctionInList(undeployedDsRef, undeployedName, map[string]string{"status": "undeployed"}),
					testAccCheckEventingFunctionsAllHaveStatus(undeployedDsRef, "undeployed"),
					testAccCheckEventingFunctionsExcludes(undeployedDsRef, deployedName),
					testAccCheckEventingFunctionsExcludes(undeployedDsRef, pausedName),

					// status = ["deployed","paused"]: includes both, excludes the undeployed function.
					testAccCheckEventingFunctionInList(multiDsRef, deployedName, map[string]string{"status": "deployed"}),
					testAccCheckEventingFunctionInList(multiDsRef, pausedName, map[string]string{"status": "paused"}),
					testAccCheckEventingFunctionsExcludes(multiDsRef, undeployedName),
				),
			},
		},
	})
}

func testAccEventingFunctionsStateFilterConfig(
	deployedName, undeployedName, pausedName, pausedState,
	deployedDs, pausedDs, undeployedDs, multiDs string,
) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id        = "%[2]s"
  project_id             = "%[3]s"
  cluster_id             = "%[4]s"
  name                   = "%[5]s"
  code                   = "%[9]s"
  state                  = "deployed"
  event_source           = { bucket = "%[7]s" }
  event_metadata_storage = { bucket = "%[8]s" }
}

resource "couchbase-capella_eventing_function" "%[6]s" {
  organization_id        = "%[2]s"
  project_id             = "%[3]s"
  cluster_id             = "%[4]s"
  name                   = "%[6]s"
  code                   = "%[9]s"
  state                  = "undeployed"
  event_source           = { bucket = "%[7]s" }
  event_metadata_storage = { bucket = "%[8]s" }
}

resource "couchbase-capella_eventing_function" "%[10]s" {
  organization_id        = "%[2]s"
  project_id             = "%[3]s"
  cluster_id             = "%[4]s"
  name                   = "%[10]s"
  code                   = "%[9]s"
  state                  = "%[11]s"
  event_source           = { bucket = "%[7]s" }
  event_metadata_storage = { bucket = "%[8]s" }
}

data "couchbase-capella_eventing_functions" "%[12]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  status          = ["deployed"]

  depends_on = [
    couchbase-capella_eventing_function.%[5]s,
    couchbase-capella_eventing_function.%[6]s,
    couchbase-capella_eventing_function.%[10]s,
  ]
}

data "couchbase-capella_eventing_functions" "%[13]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  status          = ["paused"]

  depends_on = [
    couchbase-capella_eventing_function.%[5]s,
    couchbase-capella_eventing_function.%[6]s,
    couchbase-capella_eventing_function.%[10]s,
  ]
}

data "couchbase-capella_eventing_functions" "%[14]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  status          = ["undeployed"]

  depends_on = [
    couchbase-capella_eventing_function.%[5]s,
    couchbase-capella_eventing_function.%[6]s,
    couchbase-capella_eventing_function.%[10]s,
  ]
}

data "couchbase-capella_eventing_functions" "%[15]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  status          = ["deployed", "paused"]

  depends_on = [
    couchbase-capella_eventing_function.%[5]s,
    couchbase-capella_eventing_function.%[6]s,
    couchbase-capella_eventing_function.%[10]s,
  ]
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		deployedName,
		undeployedName,
		globalBucketName,
		globalMetadataBucketName,
		eventingFunctionCode,
		pausedName,
		pausedState,
		deployedDs,
		pausedDs,
		undeployedDs,
		multiDs)
}

// eventingFunctionAttrValues returns the value of suffix for every entry in the list data source's eventing_functions.
func eventingFunctionAttrValues(state *terraform.State, dsReference, suffix string) ([]string, error) {
	resourceState, ok := state.RootModule().Resources[dsReference]
	if !ok {
		return nil, fmt.Errorf("data source %q not found in state", dsReference)
	}

	attrs := resourceState.Primary.Attributes

	count, err := strconv.Atoi(attrs["eventing_functions.#"])
	if err != nil {
		return nil, fmt.Errorf("invalid eventing_functions.# on %q: %w", dsReference, err)
	}

	values := make([]string, count)
	for i := range values {
		values[i] = attrs[fmt.Sprintf("eventing_functions.%d.%s", i, suffix)]
	}
	return values, nil
}

// testAccCheckEventingFunctionsAllHaveStatus verifies every entry in a status-filtered data source has the expected status.
func testAccCheckEventingFunctionsAllHaveStatus(dsReference, status string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		statuses, err := eventingFunctionAttrValues(state, dsReference, "status")
		if err != nil {
			return err
		}

		if len(statuses) == 0 {
			return fmt.Errorf("expected at least one eventing function in %s, got none", dsReference)
		}

		for i, got := range statuses {
			if got != status {
				return fmt.Errorf("eventing_functions.%d.status = %q, want %q", i, got, status)
			}
		}
		return nil
	}
}

// testAccCheckEventingFunctionsExcludes asserts a named function is absent from a filtered data source response.
func testAccCheckEventingFunctionsExcludes(dsReference, funcName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		names, err := eventingFunctionAttrValues(state, dsReference, "name")
		if err != nil {
			return err
		}

		if slices.Contains(names, funcName) {
			return fmt.Errorf("eventing function %q unexpectedly present in %s", funcName, dsReference)
		}
		return nil
	}
}
