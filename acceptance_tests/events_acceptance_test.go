package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

type listEventsResponse struct {
	Data []struct {
		Id string `json:"id"`
	} `json:"data"`
}

func fetchFirstEventId(t *testing.T) string {
	t.Helper()
	url := fmt.Sprintf("%s/v4/organizations/%s/events?perPage=1", globalHost, globalOrgId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	resp, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
	if err != nil {
		t.Logf("could not list org events: %v", err)
		return ""
	}
	var body listEventsResponse
	if err := json.Unmarshal(resp.Body, &body); err != nil || len(body.Data) == 0 {
		return ""
	}
	return body.Data[0].Id
}

func fetchFirstProjectEventId(t *testing.T) string {
	t.Helper()
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/events?perPage=1", globalHost, globalOrgId, globalProjectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	resp, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
	if err != nil {
		t.Logf("could not list project events: %v", err)
		return ""
	}
	var body listEventsResponse
	if err := json.Unmarshal(resp.Body, &body); err != nil || len(body.Data) == 0 {
		return ""
	}
	return body.Data[0].Id
}

func TestAccDatasourceEvents(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_events_")
	dsReference := "data.couchbase-capella_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsListConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
					resource.TestCheckResourceAttrSet(dsReference, "cursor.pages.total_items"),
					resource.TestCheckResourceAttrSet(dsReference, "cursor.pages.page"),
					resource.TestCheckResourceAttrSet(dsReference, "cursor.pages.per_page"),
				),
			},
		},
	})
}

func TestAccDatasourceEventsFilteredByProject(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_events_proj_")
	dsReference := "data.couchbase-capella_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsListFilteredByProjectConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceEventsFilteredByCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_events_cluster_")
	dsReference := "data.couchbase-capella_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsListFilteredByClusterConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceEventsWithPagination(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_events_page_")
	dsReference := "data.couchbase-capella_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsListPaginatedConfig(dsName, 1, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "page", "1"),
					resource.TestCheckResourceAttr(dsReference, "per_page", "5"),
					resource.TestCheckResourceAttrSet(dsReference, "cursor.pages.per_page"),
				),
			},
		},
	})
}

func TestAccDatasourceEventsWithSortAsc(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_events_sort_asc_")
	dsReference := "data.couchbase-capella_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsListSortedConfig(dsName, "timestamp", "asc"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "sort_by", "timestamp"),
					resource.TestCheckResourceAttr(dsReference, "sort_direction", "asc"),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceEventsWithSortDesc(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_events_sort_desc_")
	dsReference := "data.couchbase-capella_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsListSortedConfig(dsName, "timestamp", "desc"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "sort_by", "timestamp"),
					resource.TestCheckResourceAttr(dsReference, "sort_direction", "desc"),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceEventsWithTimeRange(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_events_time_")
	dsReference := "data.couchbase-capella_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsListTimeRangeConfig(dsName, "2020-01-01T00:00:00.000Z", "2099-12-31T23:59:59.999Z"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceEventsFilterBySeverity(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_events_sev_")
	dsReference := "data.couchbase-capella_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventsListBySeverityConfig(dsName, "info"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceEventsInvalidOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_events_bad_org_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_events" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Events|organization.*not found|access to the requested resource is denied|Not Found|Forbidden`),
			},
		},
	})
}

func TestAccDatasourceEventsMissingOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_events_no_org_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_events" "%[2]s" {
}
`, globalProviderBlock, dsName),
				ExpectError: regexp.MustCompile(`(?s)organization_id|argument.*required|Missing required argument`),
			},
		},
	})
}

func TestAccDatasourceEventsInvalidSortDirection(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_events_bad_sort_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventsListSortedConfig(dsName, "timestamp", "invalid"),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Events|sort_direction|invalid|asc.*desc`),
			},
		},
	})
}

func TestAccDatasourceEvent(t *testing.T) {
	eventId := fetchFirstEventId(t)
	if eventId == "" {
		t.Skip("no events found in organization, skipping single event read test")
	}

	dsName := randomStringWithPrefix("tf_acc_event_")
	dsReference := "data.couchbase-capella_event." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventConfig(dsName, eventId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "id", eventId),
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dsReference, "key"),
					resource.TestCheckResourceAttrSet(dsReference, "severity"),
					resource.TestCheckResourceAttrSet(dsReference, "source"),
					resource.TestCheckResourceAttrSet(dsReference, "timestamp"),
				),
			},
		},
	})
}

func TestAccDatasourceEventNotFound(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_event_notfound_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventConfig(dsName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Event|event.*not found|Not Found|access to the requested resource is denied`),
			},
		},
	})
}

func TestAccDatasourceEventInvalidOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_event_bad_org_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_event" "%[2]s" {
  id              = "00000000-0000-0000-0000-000000000001"
  organization_id = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Event|organization.*not found|access to the requested resource is denied|Not Found|Forbidden`),
			},
		},
	})
}

func TestAccDatasourceEventMissingId(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_event_no_id_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_event" "%[2]s" {
  organization_id = "%[3]s"
}
`, globalProviderBlock, dsName, globalOrgId),
				ExpectError: regexp.MustCompile(`(?s)id|argument.*required|Missing required argument`),
			},
		},
	})
}

func TestAccDatasourceProjectEvents(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_events_")
	dsReference := "data.couchbase-capella_project_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectEventsListConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
					resource.TestCheckResourceAttrSet(dsReference, "cursor.pages.total_items"),
				),
			},
		},
	})
}

func TestAccDatasourceProjectEventsWithoutProjectId(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_events_noproj_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_events" "%[2]s" {
  organization_id = "%[3]s"
}
`, globalProviderBlock, dsName, globalOrgId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Project Events|400|client error`),
			},
		},
	})
}

func TestAccDatasourceProjectEventsWithPagination(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_events_page_")
	dsReference := "data.couchbase-capella_project_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectEventsListPaginatedConfig(dsName, 1, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "page", "1"),
					resource.TestCheckResourceAttr(dsReference, "per_page", "5"),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceProjectEventsWithSort(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_events_sort_")
	dsReference := "data.couchbase-capella_project_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectEventsSortedConfig(dsName, "timestamp", "desc"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "sort_by", "timestamp"),
					resource.TestCheckResourceAttr(dsReference, "sort_direction", "desc"),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceProjectEventsFilterByCluster(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_events_cls_")
	dsReference := "data.couchbase-capella_project_events." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectEventsFilterByClusterConfig(dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttrSet(dsReference, "data.#"),
				),
			},
		},
	})
}

func TestAccDatasourceProjectEventsInvalidOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_events_bad_org_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_events" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "%[3]s"
}
`, globalProviderBlock, dsName, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Project Events|organization.*not found|access to the requested resource is denied|Not Found|Forbidden`),
			},
		},
	})
}

func TestAccDatasourceProjectEventsInvalidSortDirection(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_events_bad_sort_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccProjectEventsSortedConfig(dsName, "timestamp", "badDirection"),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Project Events|sort_direction|invalid|asc.*desc`),
			},
		},
	})
}

func TestAccDatasourceProjectEventsMissingOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_events_no_org_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_events" "%[2]s" {
  project_id = "%[3]s"
}
`, globalProviderBlock, dsName, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)organization_id|argument.*required|Missing required argument`),
			},
		},
	})
}

func TestAccDatasourceProjectEvent(t *testing.T) {
	eventId := fetchFirstProjectEventId(t)
	if eventId == "" {
		t.Skip("no project events found, skipping single project event read test")
	}

	dsName := randomStringWithPrefix("tf_acc_prj_event_")
	dsReference := "data.couchbase-capella_project_event." + dsName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectEventConfig(dsName, eventId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsReference, "id", eventId),
					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dsReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttrSet(dsReference, "key"),
					resource.TestCheckResourceAttrSet(dsReference, "severity"),
					resource.TestCheckResourceAttrSet(dsReference, "source"),
					resource.TestCheckResourceAttrSet(dsReference, "timestamp"),
				),
			},
		},
	})
}

func TestAccDatasourceProjectEventNotFound(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_event_notfound_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccProjectEventConfig(dsName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Project Event|event.*not found|Not Found|access to the requested resource is denied`),
			},
		},
	})
}

func TestAccDatasourceProjectEventInvalidProject(t *testing.T) {
	eventId := fetchFirstProjectEventId(t)
	if eventId == "" {
		eventId = "00000000-0000-0000-0000-000000000001"
	}
	dsName := randomStringWithPrefix("tf_acc_prj_event_bad_proj_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_event" "%[2]s" {
  id              = "%[3]s"
  organization_id = "%[4]s"
  project_id      = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName, eventId, globalOrgId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Project Event|project.*not found|access to the requested resource is denied|Not Found`),
			},
		},
	})
}

func TestAccDatasourceProjectEventMissingId(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_event_no_id_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_event" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)id|argument.*required|Missing required argument`),
			},
		},
	})
}

func TestAccDatasourceProjectEventMissingProjectId(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_event_no_proj_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_event" "%[2]s" {
  id              = "00000000-0000-0000-0000-000000000001"
  organization_id = "%[3]s"
}
`, globalProviderBlock, dsName, globalOrgId),
				ExpectError: regexp.MustCompile(`(?s)project_id|argument.*required|Missing required argument`),
			},
		},
	})
}

func TestAccDatasourceProjectEventMissingOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_event_no_org_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_event" "%[2]s" {
  id         = "00000000-0000-0000-0000-000000000001"
  project_id = "%[3]s"
}
`, globalProviderBlock, dsName, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)organization_id|argument.*required|Missing required argument`),
			},
		},
	})
}

func TestAccDatasourceProjectEventInvalidOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_prj_event_bad_org_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_event" "%[2]s" {
  id              = "00000000-0000-0000-0000-000000000001"
  organization_id = "00000000-0000-0000-0000-000000000000"
  project_id      = "%[3]s"
}
`, globalProviderBlock, dsName, globalProjectId),
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Project Event|organization.*not found|access to the requested resource is denied|Not Found|Forbidden`),
			},
		},
	})
}

func testAccEventsListConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_events" "%[2]s" {
  organization_id = "%[3]s"
}
`, globalProviderBlock, dsName, globalOrgId)
}

func testAccEventsListFilteredByProjectConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_events" "%[2]s" {
  organization_id = "%[3]s"
  project_ids     = ["%[4]s"]
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId)
}

func testAccEventsListFilteredByClusterConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_events" "%[2]s" {
  organization_id = "%[3]s"
  cluster_ids     = ["%[4]s"]
}
`, globalProviderBlock, dsName, globalOrgId, globalClusterId)
}

func testAccEventsListPaginatedConfig(dsName string, page, perPage int) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_events" "%[2]s" {
  organization_id = "%[3]s"
  page            = %[4]d
  per_page        = %[5]d
}
`, globalProviderBlock, dsName, globalOrgId, page, perPage)
}

func testAccEventsListSortedConfig(dsName, sortBy, sortDirection string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_events" "%[2]s" {
  organization_id = "%[3]s"
  sort_by         = "%[4]s"
  sort_direction  = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, sortBy, sortDirection)
}

func testAccEventsListTimeRangeConfig(dsName, from, to string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_events" "%[2]s" {
  organization_id = "%[3]s"
  from            = "%[4]s"
  to              = "%[5]s"
}
`, globalProviderBlock, dsName, globalOrgId, from, to)
}

func testAccEventsListBySeverityConfig(dsName, severity string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_events" "%[2]s" {
  organization_id = "%[3]s"
  severity_levels = ["%[4]s"]
}
`, globalProviderBlock, dsName, globalOrgId, severity)
}

func testAccEventConfig(dsName, eventId string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_event" "%[2]s" {
  id              = "%[3]s"
  organization_id = "%[4]s"
}
`, globalProviderBlock, dsName, eventId, globalOrgId)
}

func testAccProjectEventsListConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_events" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId)
}

func testAccProjectEventsListPaginatedConfig(dsName string, page, perPage int) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_events" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  page            = %[5]d
  per_page        = %[6]d
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, page, perPage)
}

func testAccProjectEventsSortedConfig(dsName, sortBy, sortDirection string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_events" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  sort_by         = "%[5]s"
  sort_direction  = "%[6]s"
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, sortBy, sortDirection)
}

func testAccProjectEventsFilterByClusterConfig(dsName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_events" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_ids     = ["%[5]s"]
}
`, globalProviderBlock, dsName, globalOrgId, globalProjectId, globalClusterId)
}

func testAccProjectEventConfig(dsName, eventId string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_project_event" "%[2]s" {
  id              = "%[3]s"
  organization_id = "%[4]s"
  project_id      = "%[5]s"
}
`, globalProviderBlock, dsName, eventId, globalOrgId, globalProjectId)
}
