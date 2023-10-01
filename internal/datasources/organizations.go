package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/errors"

	providerschema "terraform-provider-capella/internal/schema"
)

var (
	_ datasource.DataSource = &Organization{}
)

type Organization struct {
	*providerschema.Data
}

func NewOrganization() datasource.DataSource {
	return &Organization{}
}

func (o Organization) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (o Organization) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"preferences": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"sessionDuration": schema.Int64Attribute{
									Computed: true,
								},
							},
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
					},
				},
			},
		},
	}
}

func (o Organization) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Organizations
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := o.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Organizations",
			"Could not read organizations in cluster"+state.OrganizationId.String()+": "+err.Error())
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
	)

	// Make request to list allowlists
	response, err := o.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s", o.HostURL, organizationId),
		http.MethodGet,
		nil,
		o.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Error Reading Capella Organizations",
				"Could not read organizations in cluster "+state.OrganizationId.String()+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Organizations",
			"Could not read organizations in cluster "+state.OrganizationId.String()+": "+err.Error(),
		)
		return
	}

	organizationsResponse := api.GetOrganizationsResponse{}
	err = json.Unmarshal(response.Body, &organizationsResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading organizations",
			"Could not create organizations, unexpected error: "+err.Error(),
		)
		return
	}

	for _, org := range organizationsResponse.Data {
		orgState := providerschema.OneOrganization{
			Id:          types.StringValue(org.Id.String()),
			Name:        types.StringValue(org.Name),
			Description: types.StringValue(*org.Description),
			Audit: providerschema.CouchbaseAuditData{
				CreatedAt:  types.StringValue(org.Audit.CreatedAt.String()),
				CreatedBy:  types.StringValue(org.Audit.CreatedBy),
				ModifiedAt: types.StringValue(org.Audit.ModifiedAt.String()),
				ModifiedBy: types.StringValue(org.Audit.ModifiedBy),
				Version:    types.Int64Value(int64(org.Audit.Version)),
			},
		}
		state.Data = append(state.Data, orgState)
	}

	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//state = o.mapResponseBody(organizationsResponse, &state)
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Error reading allowlist",
	//		"Could not create allowlist, unexpected error: "+err.Error(),
	//	)
	//	return
	//}
	//
	//// Set state
	//diags = resp.State.Set(ctx, &state)
	//
	//resp.Diagnostics.Append(diags...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}

}

// validate is used to verify that all the fields in the datasource
// have been populated.
func (o *Organization) validate(state providerschema.Organizations) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	return nil
}

// mapResponseBody is used to map the response body from a call to
// get allowlists to the allowlists schema that will be used by terraform.
//func (o *Organization) mapResponseBody(
//	organizationsResponse api.GetOrganizationResponse,
//	state *providerschema.Organizations,
//) providerschema.organizationsResponse {
//	for _, allowList := range organizationsResponse.Data {
//		allowListState := providerschema.OneAllowList{
//			Id:             types.StringValue(allowList.Id.String()),
//			OrganizationId: types.StringValue(state.OrganizationId.ValueString()),
//			ProjectId:      types.StringValue(state.ProjectId.ValueString()),
//			ClusterId:      types.StringValue(state.ClusterId.ValueString()),
//			Cidr:           types.StringValue(allowList.Cidr),
//			Comment:        types.StringValue(allowList.Comment),
//			ExpiresAt:      types.StringValue(allowList.ExpiresAt),
//			Audit: providerschema.CouchbaseAuditData{
//				CreatedAt:  types.StringValue(allowList.Audit.CreatedAt.String()),
//				CreatedBy:  types.StringValue(allowList.Audit.CreatedBy),
//				ModifiedAt: types.StringValue(allowList.Audit.ModifiedAt.String()),
//				ModifiedBy: types.StringValue(allowList.Audit.ModifiedBy),
//				Version:    types.Int64Value(int64(allowList.Audit.Version)),
//			},
//		}
//		state.Data = append(state.Data, allowListState)
//	}
//	return *state
//}
