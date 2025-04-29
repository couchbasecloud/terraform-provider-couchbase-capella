package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AllowLists{}
	_ datasource.DataSourceWithConfigure = &AllowLists{}
)

// AllowLists is the allow list data source implementation.
type AllowLists struct {
	*providerschema.Data
}

// NewAllowLists is a helper function to simplify the provider implementation.
func NewAllowLists() datasource.DataSource {
	return &AllowLists{}
}

// Metadata returns the allow list data source type name.
func (d *AllowLists) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_allowlists"
}

// Schema defines the schema for the allowlist data source.
func (d *AllowLists) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the allowlists details for a Capella cluster.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Capella organization.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Capella project.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Capella cluster.",
			},
			"data": schema.ListNestedAttribute{
				MarkdownDescription: "The list of allowlists.",
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
							MarkdownDescription: "An RFC3339 timestamp determining when the allowed CIDR will expire.",
						},
						"audit": computedAuditAttribute,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of allowlists.
func (d *AllowLists) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AllowLists
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate state is not empty
	err := d.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella AllowLists",
			"Could not read allow lists in cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	allowLists, err := d.listAllowLists(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella AllowLists",
			"Could not read allow lists in cluster "+state.ClusterId.String()+": "+api.ParseError(err),
		)
		return
	}

	state = d.mapResponseBody(allowLists, &state)
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

// listAllowLists executes calls to the list allowlist endpoint. It handles pagination and
// returns a slice of individual allowlists responses retrieved from multiple pages.
func (d *AllowLists) listAllowLists(ctx context.Context, organizationId, projectId, clusterId string) ([]api.GetAllowListResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/allowedcidrs",
		d.HostURL,
		organizationId,
		projectId,
		clusterId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	return api.GetPaginated[[]api.GetAllowListResponse](ctx, d.Client, d.Token, cfg, api.SortById)
}

// Configure adds the provider configured client to the allowlist data source.
func (d *AllowLists) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// mapResponseBody is used to map the response body from a call to
// listAllowlists to the allowlists schema that will be used by terraform.
func (d *AllowLists) mapResponseBody(
	allowLists []api.GetAllowListResponse,
	state *providerschema.AllowLists,
) providerschema.AllowLists {
	for _, allowList := range allowLists {
		allowListState := providerschema.OneAllowList{
			Id:             types.StringValue(allowList.Id.String()),
			OrganizationId: types.StringValue(state.OrganizationId.ValueString()),
			ProjectId:      types.StringValue(state.ProjectId.ValueString()),
			ClusterId:      types.StringValue(state.ClusterId.ValueString()),
			Cidr:           types.StringValue(allowList.Cidr),
			Audit: providerschema.CouchbaseAuditData{
				CreatedAt:  types.StringValue(allowList.Audit.CreatedAt.String()),
				CreatedBy:  types.StringValue(allowList.Audit.CreatedBy),
				ModifiedAt: types.StringValue(allowList.Audit.ModifiedAt.String()),
				ModifiedBy: types.StringValue(allowList.Audit.ModifiedBy),
				Version:    types.Int64Value(int64(allowList.Audit.Version)),
			},
		}
		if allowList.Comment != nil {
			allowListState.Comment = types.StringValue(*allowList.Comment)
		}
		if allowList.ExpiresAt != nil {
			allowListState.ExpiresAt = types.StringValue(*allowList.ExpiresAt)
		}
		state.Data = append(state.Data, allowListState)
	}
	return *state
}

// validate is used to verify that all the fields in the datasource
// have been populated.
func (d *AllowLists) validate(state providerschema.AllowLists) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if state.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if state.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	return nil
}
