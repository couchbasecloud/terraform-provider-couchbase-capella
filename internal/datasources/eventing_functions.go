package datasources

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/eventing_function"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = (*EventingFunctions)(nil)
	_ datasource.DataSourceWithConfigure = (*EventingFunctions)(nil)
)

// EventingFunctions is the list eventing functions data source implementation.
type EventingFunctions struct {
	*providerschema.Data
}

// NewEventingFunctions is a helper function to simplify the provider implementation.
func NewEventingFunctions() datasource.DataSource {
	return &EventingFunctions{}
}

// Metadata returns the eventing functions data source type name.
func (d *EventingFunctions) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eventing_functions"
}

// Schema defines the schema for the eventing functions data source.
func (d *EventingFunctions) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = EventingFunctionsSchema()
}

// Read refreshes the Terraform state with the latest list of eventing functions.
func (d *EventingFunctions) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.EventingFunctions
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	requestUrl := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions",
		d.HostURL, organizationId, projectId, clusterId,
	)

	// The list endpoint filters by state server-side via a comma-separated status query parameter.
	if !state.Status.IsNull() && !state.Status.IsUnknown() {
		var statuses []string
		diags := state.Status.ElementsAs(ctx, &statuses, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(statuses) > 0 {
			requestUrl += "?status=" + strings.Join(statuses, ",")
		}
	}

	cfg := api.EndpointCfg{Url: requestUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	functions, err := api.GetPaginated[[]eventing_function.EventingFunction](ctx, d.ClientV1, d.Token, cfg, api.SortByName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Eventing Functions",
			fmt.Sprintf("Could not read eventing functions in cluster %s, unexpected error: %s", clusterId, api.ParseError(err)),
		)
		return
	}

	state.EventingFunctions = make([]providerschema.OneEventingFunction, 0, len(functions))
	for _, function := range functions {
		model, diags := providerschema.NewOneEventingFunction(ctx, function)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.EventingFunctions = append(state.EventingFunctions, model)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the eventing functions data source.
func (d *EventingFunctions) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
