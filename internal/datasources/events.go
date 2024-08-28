package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Events{}
	_ datasource.DataSourceWithConfigure = &Events{}
)

// Events is the Events data source implementation.
type Events struct {
	*providerschema.Data
}

// NewEvents is a helper function to simplify the provider implementation.
func NewEvents() datasource.DataSource {
	return &Events{}
}

// Metadata returns the backup data source type name.
func (d *Events) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_events"
}

// Schema defines the schema for the Events data source.
func (d *Events) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = EventsSchema()
}

// Read refreshes the Terraform state with the latest data of Events.
func (d *Events) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Events
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		err            error
	)
	queryParam, err := d.buildQueryParams(ctx, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Events",
			"Could not read events : "+err.Error(),
		)
		return
	}

	finalUrl := fmt.Sprintf("%s/v4/organizations/%s/events", d.HostURL, organizationId)
	if len(queryParam) > 0 {
		finalUrl = finalUrl + BuildQueryParams(queryParam)
	}

	events, err := d.listEvents(ctx, finalUrl)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Events",
			"Could not read events : "+err.Error(),
		)
		return
	}

	eventItems, err := MapResponseEventsBody(ctx, events.Data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading events",
			"Could not read events, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state
	state.Data = eventItems
	state.Cursor = &providerschema.Cursor{
		Hrefs: &providerschema.Hrefs{
			First:    types.StringValue(events.Cursor.Hrefs.First),
			Last:     types.StringValue(events.Cursor.Hrefs.Last),
			Next:     types.StringValue(events.Cursor.Hrefs.Next),
			Previous: types.StringValue(events.Cursor.Hrefs.Previous),
		},
		Pages: &providerschema.Pages{
			Last:       types.Int64Value(int64(events.Cursor.Pages.Last)),
			Next:       types.Int64Value(int64(events.Cursor.Pages.Next)),
			Page:       types.Int64Value(int64(events.Cursor.Pages.Page)),
			PerPage:    types.Int64Value(int64(events.Cursor.Pages.PerPage)),
			Previous:   types.Int64Value(int64(events.Cursor.Pages.Previous)),
			TotalItems: types.Int64Value(int64(events.Cursor.Pages.TotalItems)),
		},
	}
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the events data source.
func (d *Events) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.Data = data
}

func (d *Events) listEvents(ctx context.Context, url string) (api.GetEventsResponse, error) {
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := d.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		d.Token,
		nil,
	)
	if err != nil {
		return api.GetEventsResponse{}, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	var events api.GetEventsResponse
	err = json.Unmarshal(response.Body, &events)
	if err != nil {
		return api.GetEventsResponse{}, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}
	return events, nil
}

func (d *Events) buildQueryParams(ctx context.Context, state *providerschema.Events) (map[string][]string, error) {
	queryParam := make(map[string][]string)
	if !state.ProjectIds.IsNull() && !state.ProjectIds.IsUnknown() {
		projectIds, err := ConvertToList(ctx, state.ProjectIds)
		if err != nil {
			return nil, err
		}
		queryParam["projectIds"] = projectIds
	}
	if !state.ClusterIds.IsNull() && !state.ClusterIds.IsUnknown() {
		clusterIds, err := ConvertToList(ctx, state.ClusterIds)
		if err != nil {
			return nil, err
		}
		queryParam["clusterIds"] = clusterIds
	}
	if !state.UserIds.IsNull() && !state.UserIds.IsUnknown() {
		userIds, err := ConvertToList(ctx, state.UserIds)
		if err != nil {
			return nil, err
		}
		queryParam["userIds"] = userIds
	}
	if !state.SeverityLevels.IsNull() && !state.SeverityLevels.IsUnknown() {
		severityLevels, err := ConvertToList(ctx, state.SeverityLevels)
		if err != nil {
			return nil, err
		}
		queryParam["severityLevels"] = severityLevels
	}
	if !state.Tags.IsNull() && !state.Tags.IsUnknown() {
		tags, err := ConvertToList(ctx, state.Tags)
		if err != nil {
			return nil, err
		}
		queryParam["tags"] = tags
	}
	if !state.From.IsNull() && !state.From.IsUnknown() {
		from := state.From.ValueString()
		queryParam["from"] = []string{from}
	}
	if !state.To.IsNull() && !state.To.IsUnknown() {
		to := state.To.ValueString()
		queryParam["to"] = []string{to}
	}
	if !state.Page.IsNull() && !state.Page.IsUnknown() {
		page := int(state.Page.ValueInt64())
		queryParam["page"] = []string{strconv.Itoa(page)}
	}
	if !state.PerPage.IsNull() && !state.PerPage.IsUnknown() {
		perPage := int(state.PerPage.ValueInt64())
		queryParam["perPage"] = []string{strconv.Itoa(perPage)}
	}
	if !state.SortBy.IsNull() && !state.SortBy.IsUnknown() {
		sortBy := state.SortBy.ValueString()
		queryParam["sortBy"] = []string{sortBy}
	}
	if !state.SortDirection.IsNull() && !state.SortDirection.IsUnknown() {
		sortDir := state.SortDirection.ValueString()
		queryParam["sortDirection"] = []string{sortDir}
	}
	return queryParam, nil
}
