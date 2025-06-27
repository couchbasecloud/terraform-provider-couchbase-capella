package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*AppServiceCidrs)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServiceCidrs)(nil)
)

// AppServiceCidrs is the data source implementation for retrieving allowed CIDRs for an App Service.
type AppServiceCidrs struct {
	*providerschema.Data
}

// NewAppServiceCidrs is used in (p *capellaProvider) DataSources for building the provider.
func NewAppServiceCidrs() datasource.DataSource {
	return &AppServiceCidrs{}
}

// Metadata returns the App Service CIDRs data source type name.
func (a *AppServiceCidrs) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_services_cidr"
}

// Schema defines the schema for the App Service CIDRs data source.
func (a *AppServiceCidrs) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the allowed CIDR blocks for a Capella App Service.",
		Attributes: map[string]schema.Attribute{
			"data": schema.ListNestedAttribute{
				MarkdownDescription: "The list of allowed CIDR blocks on an App Service. ",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the allowed CIDR block.",
						},
						"organization_id": computedStringAttribute,
						"project_id":      computedStringAttribute,
						"cluster_id":      computedStringAttribute,
						"app_service_id":  computedStringAttribute,
						"cidr": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The trusted CIDR block to allow the database connections from.",
						},
						"comment": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A short description of the allowed CIDR block.",
						},
						"expires_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "An RFC3339 timestamp determining when the allowed CIDR block will expire. If this field is empty/omitted then the allowed CIDR block is permanent. It will never automatically expire.",
						},
						"audit": computedAuditAttribute,
					},
				},
			},
		},
	}
}

// listAllowedCIDRs executes calls to the list app service allowed cidrs endpoint. It handles pagination and
// returns a slice of individual allowed cidr responses retrieved from multiple pages.
func (a *AppServiceCidrs) listAllowedCIDRs(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) ([]api.AppServiceAllowedCIDRResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/allowedcidrs",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	return api.GetPaginated[[]api.AppServiceAllowedCIDRResponse](ctx, a.Client, a.Token, cfg, api.SortById)
}

// Read refreshes the Terraform state with the allowed CIDRs in the Terraform state file.
func (a *AppServiceCidrs) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state *providerschema.AppServiceCIDRs
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate state is not empty
	organizationId, projectId, clusterId, appServiceId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Allowed CIDRs",
			"Could not validate App Service allowed CIDRS in state file "+state.AppServiceId.String()+": "+err.Error(),
		)
		return
	}

	response, err := a.listAllowedCIDRs(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Listing App Service allowed CIDRs",
			fmt.Sprintf("Could not list App Service allowed CIDRs in organization %s, unexpected error: %s", organizationId, api.ParseError(err)),
		)
		return
	}

	state = a.mapResponseBody(ctx, response, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error morphing allowed CIDRs",
			"Could not morph list allowed CIDRs response, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// mapResponseBody is used to map the response body from a call to
// listAllowedCidrs to the allowed cidrs schema that will be used by terraform.
func (a *AppServiceCidrs) mapResponseBody(
	ctx context.Context,
	allowLists []api.AppServiceAllowedCIDRResponse,
	state *providerschema.AppServiceCIDRs,
) *providerschema.AppServiceCIDRs {
	state = &providerschema.AppServiceCIDRs{
		OrganizationId: types.StringValue(state.OrganizationId.ValueString()),
		ProjectId:      types.StringValue(state.ProjectId.ValueString()),
		ClusterId:      types.StringValue(state.ClusterId.ValueString()),
		AppServiceId:   types.StringValue(state.AppServiceId.ValueString()),
	}
	for _, allowList := range allowLists {
		// Create audit data object
		audit := providerschema.NewCouchbaseAuditData(allowList.Audit)
		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			// If we have an error, we should not continue processing the rest of the allowlists.
			return nil
		}

		allowListState := providerschema.AppServiceCIDRData{
			Id:    types.StringValue(allowList.Id),
			Cidr:  types.StringValue(allowList.Cidr),
			Audit: auditObj,
		}
		if allowList.Comment != "" {
			allowListState.Comment = types.StringValue(allowList.Comment)
		}
		if allowList.ExpiresAt != "" {
			allowListState.ExpiresAt = types.StringValue(allowList.ExpiresAt)
		}
		state.Data = append(state.Data, allowListState)
	}
	return state
}

// Configure is used to configure the AppServiceCidrs data source with the provider data.
func (a *AppServiceCidrs) Configure(
	_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
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
