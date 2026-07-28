package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	eventingapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/eventingfunction"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = (*EventingFunction)(nil)
	_ datasource.DataSourceWithConfigure = (*EventingFunction)(nil)
)

// EventingFunction is the eventing function data source implementation.
type EventingFunction struct {
	*providerschema.Data
}

// NewEventingFunction is a helper function to simplify the provider implementation.
func NewEventingFunction() datasource.DataSource {
	return &EventingFunction{}
}

// Metadata returns the eventing function data source type name.
func (d *EventingFunction) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eventing_function"
}

// Schema defines the schema for the eventing function data source.
func (d *EventingFunction) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = EventingFunctionSchema()
}

// Read refreshes the Terraform state with the latest eventing function data.
func (d *EventingFunction) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.EventingFunction
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
		functionName   = state.Name.ValueString()
	)

	requestUrl := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions/%s",
		d.HostURL, organizationId, projectId, clusterId, url.PathEscape(functionName),
	)
	// export omits the read-only status field so the response can be reused as a create payload.
	if state.Export.ValueBool() {
		requestUrl += "?export=true"
	}

	cfg := api.EndpointCfg{Url: requestUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Eventing Function",
			fmt.Sprintf("Could not read eventing function %s in cluster %s, unexpected error: %s", functionName, clusterId, api.ParseError(err)),
		)
		return
	}

	var function eventingapi.GetEventingFunctionResponse
	if err := json.Unmarshal(response.Body, &function); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Eventing Function",
			fmt.Sprintf("Could not read eventing function %s, unexpected error: %s", functionName, err.Error()),
		)
		return
	}

	model, diags := providerschema.NewEventingFunction(ctx, function, organizationId, projectId, clusterId, functionName, state.Export)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &model)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the eventing function data source.
func (d *EventingFunction) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.Data = data
}
