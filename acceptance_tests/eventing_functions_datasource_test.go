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

	// expectedFunctionAttrs is used to verify that allDsReference and matchDsReference data sources contain an eventing function
	// with the expected attributes.
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

// eventingFunctionAttrValues returns the value of suffix (e.g. "name" or "status") for every entry
// in the eventing_functions list exposed by the datasource in state.
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

// testAccCheckEventingFunctionsAllHaveStatus verifies every entry returned by a status-filtered
// eventing functions datasource carries the expected status.
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

// testAccCheckEventingFunctionsExcludes asserts a named function is absent from a filtered
// eventing functions datasource response.
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
