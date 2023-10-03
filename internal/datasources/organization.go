package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/api/organization"
	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	providerschema "terraform-provider-capella/internal/schema"
)

var (
	_ datasource.DataSource              = &Organization{}
	_ datasource.DataSourceWithConfigure = &Organization{}
)

type Organization struct {
	*providerschema.Data
}

func NewOrganization() datasource.DataSource {
	return &Organization{}
}

func (o *Organization) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (o *Organization) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
								"session_duration": schema.Int64Attribute{
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

func (o *Organization) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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

	var organizationId = state.OrganizationId.ValueString()

	// Make request to get organization
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

	organizationsResponse := organization.GetOrganizationResponse{}
	err = json.Unmarshal(response.Body, &organizationsResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading organizations",
			"Could not create organizations, unexpected error: "+err.Error(),
		)
		return
	}

	orgState := providerschema.OneOrganization{
		Id:          types.StringValue(organizationsResponse.Id.String()),
		Name:        types.StringValue(organizationsResponse.Name),
		Description: types.StringValue(*organizationsResponse.Description),
		Audit: providerschema.CouchbaseAuditData{
			CreatedAt:  types.StringValue(organizationsResponse.Audit.CreatedAt.String()),
			CreatedBy:  types.StringValue(organizationsResponse.Audit.CreatedBy),
			ModifiedAt: types.StringValue(organizationsResponse.Audit.ModifiedAt.String()),
			ModifiedBy: types.StringValue(organizationsResponse.Audit.ModifiedBy),
			Version:    types.Int64Value(int64(organizationsResponse.Audit.Version)),
		},
	}
	state.Data = append(state.Data, orgState)

	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// validate is used to verify that all the fields in the datasource
// have been populated.
func (o *Organization) validate(state providerschema.Organizations) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	return nil
}

// Configure adds the provider configured client to the organization data source.
func (o *Organization) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	o.Data = data
}