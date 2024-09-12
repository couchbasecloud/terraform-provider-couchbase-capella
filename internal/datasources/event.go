package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Event{}
	_ datasource.DataSourceWithConfigure = &Event{}
)

// Event is the Event data source implementation.
type Event struct {
	*providerschema.Data
}

// NewEvent is a helper function to simplify the provider implementation.
func NewEvent() datasource.DataSource {
	return &Event{}
}

// Metadata returns the backup data source type name.
func (d *Event) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event"
}

// Schema defines the schema for the Event data source.
func (d *Event) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = EventSchema()
}

// Read refreshes the Terraform state with the latest data of Event.
func (d *Event) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Event
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		eventId        = state.Id.ValueString()
		err            error
	)

	event, err := d.getEvent(ctx, organizationId, eventId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Event",
			"Could not read event : "+err.Error(),
		)
		return
	}

	newEventState, err := MapEventResponseBody(ctx, event, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading event",
			"Could not read event, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, newEventState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the event data source.
func (d *Event) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *Event) getEvent(ctx context.Context, organizationId, eventId string) (api.GetEventResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/events/%s",
		d.HostURL,
		organizationId,
		eventId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := d.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		d.Token,
		nil,
	)
	if err != nil {
		return api.GetEventResponse{}, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	var event api.GetEventResponse
	err = json.Unmarshal(response.Body, &event)
	if err != nil {
		return api.GetEventResponse{}, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}
	return event, nil
}
