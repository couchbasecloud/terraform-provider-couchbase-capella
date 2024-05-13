package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AuditLogEventIDs{}
	_ datasource.DataSourceWithConfigure = &AuditLogEventIDs{}
)

// AuditLogEventIDs is a list of audit log event ids.
type AuditLogEventIDs struct {
	*providerschema.Data
}

// NewAuditLogEventIDs is a helper function to simplify the provider implementation.
func NewAuditLogEventIDs() datasource.DataSource {
	return &AuditLogEventIDs{}
}

// Metadata returns the certificates data source type name.
func (a *AuditLogEventIDs) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audit_log_event_ids"
}

// Schema defines the schema for the audit log event ids data source.
func (a *AuditLogEventIDs) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"cluster_id":      requiredStringAttribute,
			"data": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"module": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of projects.
func (a *AuditLogEventIDs) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AuditLogEventIDs
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate state is not empty
	err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading audit log event ids",
			"Could not read audit log event ids: "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/auditLogEvents", a.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading audit log event ids",
			"Could not read audit log event ids: "+api.ParseError(err),
		)
		return
	}

	auditLogEventIDsResponse := api.GetAuditLogEventIDsResponse{}
	err = json.Unmarshal(response.Body, &auditLogEventIDsResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading audit log event ids",
			"Could not read audit log event ids: "+api.ParseError(err),
		)
		return
	}

	for _, event := range auditLogEventIDsResponse.Events {
		eventState := providerschema.AuditLogEventID{
			Description: types.StringValue(event.Description),
			Id:          types.Int64Value(int64(event.Id)),
			Module:      types.StringValue(event.Module),
			Name:        types.StringValue(event.Name),
		}

		state.Data = append(state.Data, eventState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the project data source.
func (a *AuditLogEventIDs) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	a.Data = data
}
