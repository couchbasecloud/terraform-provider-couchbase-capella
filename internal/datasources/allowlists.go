package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AllowList{}
	_ datasource.DataSourceWithConfigure = &AllowList{}
)

// AllowList is the allow list data source implementation.
type AllowList struct {
	*providerschema.Data
}

// NewAllowList is a helper function to simplify the provider implementation.
func NewAllowList() datasource.DataSource {
	return &AllowList{}
}

// Metadata returns the allow list data source type name.
func (d *AllowList) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_allowlists"
}

// Schema defines the schema for the allowlist data source.
func (d *AllowList) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
						},
						"project_id": schema.StringAttribute{
							Computed: true,
						},
						"cluster_id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"audit": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"created_at": schema.StringAttribute{
									Computed: true,
								},
								"created_by": schema.StringAttribute{
									Computed: true,
								},
								"modified_at": schema.StringAttribute{
									Computed: true,
								},
								"modified_by": schema.StringAttribute{
									Computed: true,
								},
								"version": schema.Int64Attribute{
									Computed: true,
								},
							},
						},
						"if_match": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of allowlists.
func (d *AllowList) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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

	// Make request to list allowlists
	response, err := d.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/allowedcidrs", d.HostURL, organizationId, projectId, clusterId),
		http.MethodGet,
		nil,
		d.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Error Reading Capella AllowLists",
				"Could not read allow lists in cluster "+state.ClusterId.String()+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading AllowLists",
			"Could not read allow lists in cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	allowListsResponse := api.GetAllowListsResponse{}
	err = json.Unmarshal(response.Body, &allowListsResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading allowlist",
			"Could not create allowlist, unexpected error: "+err.Error(),
		)
		return
	}

	state = d.mapResponseBody(allowListsResponse, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading allowlist",
			"Could not create allowlist, unexpected error: "+err.Error(),
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

// Configure adds the provider configured client to the allowlist data source.
func (d *AllowList) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
// get allowlists to the allowlists schema that will be used by terraform.
func (d *AllowList) mapResponseBody(
	allowListsResponse api.GetAllowListsResponse,
	state *providerschema.AllowLists,
) providerschema.AllowLists {
	for _, allowList := range allowListsResponse.Data {
		allowListState := providerschema.OneAllowList{
			Id:             types.StringValue(allowList.Id.String()),
			OrganizationId: types.StringValue(state.OrganizationId.ValueString()),
			ProjectId:      types.StringValue(state.ProjectId.ValueString()),
			ClusterId:      types.StringValue(state.ClusterId.ValueString()),
			Cidr:           types.StringValue(allowList.Cidr),
			Comment:        types.StringValue(allowList.Comment),
			ExpiresAt:      types.StringValue(allowList.ExpiresAt),
			Audit: providerschema.CouchbaseAuditData{
				CreatedAt:  types.StringValue(allowList.Audit.CreatedAt.String()),
				CreatedBy:  types.StringValue(allowList.Audit.CreatedBy),
				ModifiedAt: types.StringValue(allowList.Audit.ModifiedAt.String()),
				ModifiedBy: types.StringValue(allowList.Audit.ModifiedBy),
				Version:    types.Int64Value(int64(allowList.Audit.Version)),
			},
		}
		state.Data = append(state.Data, allowListState)
	}
	return *state
}

// validate is used to verify that all the fields in the datasource
// have been populated.
func (d *AllowList) validate(state providerschema.AllowLists) error {
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
