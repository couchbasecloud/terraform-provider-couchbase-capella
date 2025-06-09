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

type AppServiceCidrs struct {
	*providerschema.Data
}

func NewAppServiceCidrs() datasource.DataSource {
	return &AppServiceCidrs{}
}

func (a *AppServiceCidrs) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_services_cidr"
}

func (a *AppServiceCidrs) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the allowed CIDRs for a Capella App Service.",
		Attributes: map[string]schema.Attribute{
			"data": schema.ListNestedAttribute{
				MarkdownDescription: "The list of allowed CIDRs on an App Service. ",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the allowed CIDR.",
						},
						"organization_id": computedStringAttribute,
						"project_id":      computedStringAttribute,
						"cluster_id":      computedStringAttribute,
						"cidr": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The trusted CIDR to allow the database connections from.",
						},
						"comment": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The trusted CIDR to allow the database connections from.",
						},
						"expires_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "An RFC3339 timestamp determining when the allowed CIDR will expire. If this field is empty/omitted then the allowed CIDR is permanent and will never automatically expire.",
						},
						"audit": computedAuditAttribute,
					},
				},
			},
		},
	}
}

// listAllowLists executes calls to the list allowlist endpoint. It handles pagination and
// returns a slice of individual allowlists responses retrieved from multiple pages.
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

func (a *AppServiceCidrs) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServiceCIDRs
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate state is not empty
	organizationId, projectId, clusterId, appServiceId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella AllowLists",
			"Could not read allow lists in cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	response, err := a.listAllowedCIDRs(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Services",
			fmt.Sprintf("Could not read app services in organization %s, unexpected error: %s", organizationId, api.ParseError(err)),
		)
		return
	}

	state = a.mapResponseBody(ctx, response, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading allowlist",
			"Could not read allowlist, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// mapResponseBody is used to map the response body from a call to
// listAllowlists to the allowlists schema that will be used by terraform.
func (a *AppServiceCidrs) mapResponseBody(
	ctx context.Context,
	allowLists []api.AppServiceAllowedCIDRResponse,
	state *providerschema.AppServiceCIDRs,
) providerschema.AppServiceCIDRs {
	state = &providerschema.AppServiceCIDRs{
		OrganizationId: types.StringValue(state.OrganizationId.ValueString()),
		ProjectId:      types.StringValue(state.ProjectId.ValueString()),
		ClusterId:      types.StringValue(state.ClusterId.ValueString()),
	}
	for _, allowList := range allowLists {
		// Create audit data object
		audit := providerschema.NewCouchbaseAuditData(allowList.Audit)
		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			// Handle error - for now we'll set audit to null
			auditObj = types.ObjectNull(audit.AttributeTypes())
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
	return *state
}

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
