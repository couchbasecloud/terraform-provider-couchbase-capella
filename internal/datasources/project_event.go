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
	_ datasource.DataSource              = &ProjectEvent{}
	_ datasource.DataSourceWithConfigure = &ProjectEvent{}
)

// ProjectEvent is the ProjectEvent data source implementation.
type ProjectEvent struct {
	*providerschema.Data
}

// NewProjectEvent is a helper function to simplify the provider implementation.
func NewProjectEvent() datasource.DataSource {
	return &ProjectEvent{}
}

// Metadata returns the backup data source type name.
func (d *ProjectEvent) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_event"
}

// Schema defines the schema for the ProjectEvent data source.
func (d *ProjectEvent) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ProjectEventSchema()
}

// Read refreshes the Terraform state with the latest data of ProjectEvent.
func (d *ProjectEvent) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Event
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		eventId        = state.Id.ValueString()
		err            error
	)

	event, err := d.getProjectEvent(ctx, organizationId, projectId, eventId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Project Event",
			"Could not read project event : "+err.Error(),
		)
		return
	}

	newEventState, err := MapEventResponseBody(ctx, event, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Project Event",
			"Could not read project event, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, newEventState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the ProjectEvent data source.
func (d *ProjectEvent) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProjectEvent) getProjectEvent(ctx context.Context, organizationId, projectId, eventId string) (api.GetEventResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/events/%s",
		d.HostURL,
		organizationId,
		projectId,
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
