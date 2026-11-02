package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppEndpointAuditLogEventIDs)(nil)
	_ datasource.DataSourceWithConfigure = (*AppEndpointAuditLogEventIDs)(nil)
)

// AppEndpointAuditLogEventIDs is the App Endpoint audit log event IDs data source implementation.
type AppEndpointAuditLogEventIDs struct {
	*providerschema.Data
}

// NewAppEndpointAuditLogEventIDs is a helper function to simplify the provider implementation.
func NewAppEndpointAuditLogEventIDs() datasource.DataSource {
	return &AppEndpointAuditLogEventIDs{}
}

// Metadata returns the data source type name.
func (d *AppEndpointAuditLogEventIDs) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_audit_log_event_ids"
}

// Schema defines the schema for the data source.
func (d *AppEndpointAuditLogEventIDs) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppEndpointAuditLogEventIDsSchema()
}

// Configure adds the provider configured client to the data source.
func (d *AppEndpointAuditLogEventIDs) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// appEndpointAuditLogEventsResponse types the events map, which the generated AuditLogEventsResponse leaves untyped.
type appEndpointAuditLogEventsResponse struct {
	Events map[string]struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Enabled     bool   `json:"enabled"`
		Filterable  bool   `json:"filterable"`
	} `json:"events"`
}

// Read refreshes the Terraform state with the audit log events of the App Endpoint.
func (d *AppEndpointAuditLogEventIDs) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.AppEndpointAuditLogEventIDs
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/auditLogEvents",
		d.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.AppServiceId.ValueString(),
		config.AppEndpointName.ValueString(),
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Endpoint Audit Log Event IDs",
			fmt.Sprintf("Could not read audit log event IDs for App Endpoint %s: %s", config.AppEndpointName.ValueString(), api.ParseError(err)),
		)
		return
	}

	var eventsResp appEndpointAuditLogEventsResponse
	if err := json.Unmarshal(response.Body, &eventsResp); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Endpoint Audit Log Event IDs",
			"Could not unmarshal audit log events response: "+err.Error(),
		)
		return
	}

	config.Data = []providerschema.AppEndpointAuditLogEventID{}
	for id, event := range eventsResp.Events {
		eventID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading App Endpoint Audit Log Event IDs",
				fmt.Sprintf("Could not parse audit event ID %q: %s", id, err),
			)
			return
		}
		config.Data = append(config.Data, providerschema.AppEndpointAuditLogEventID{
			Id:          types.Int64Value(eventID),
			Name:        types.StringValue(event.Name),
			Description: types.StringValue(event.Description),
			Enabled:     types.BoolValue(event.Enabled),
			Filterable:  types.BoolValue(event.Filterable),
		})
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
