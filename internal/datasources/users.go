package datasources

import (
	"context"
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
	_ datasource.DataSource              = &Users{}
	_ datasource.DataSourceWithConfigure = &Users{}
)

// Users is the user data source implementation.
type Users struct {
	*providerschema.Data
}

// NewUsers is a helper function to simplify the provider implementation.
func NewUsers() datasource.DataSource {
	return &Users{}
}

// Metadata returns the user data source type name.
func (d *Users) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the User data source.
func (d *Users) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                   computedStringAttribute,
						"name":                 computedStringAttribute,
						"status":               computedStringAttribute,
						"inactive":             computedBoolAttribute,
						"email":                computedStringAttribute,
						"organization_id":      computedStringAttribute,
						"organization_roles":   computedListAttribute,
						"last_login":           computedStringAttribute,
						"region":               computedStringAttribute,
						"time_zone":            computedStringAttribute,
						"enable_notifications": computedBoolAttribute,
						"expires_at":           computedStringAttribute,
						"resources": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type":  computedStringAttribute,
									"id":    computedStringAttribute,
									"roles": computedListAttribute,
								},
							},
						},
						"audit": computedAuditAttribute,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of Users.
func (d *Users) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Users
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate state is not empty
	err := d.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Users",
			"Could not read users in organization "+state.OrganizationId.String()+": "+err.Error(),
		)
		return
	}

	organizationId := state.OrganizationId.ValueString()

	// Make request to list Users
	url := fmt.Sprintf("%s/v4/organizations/%s/users", d.HostURL, organizationId)
	response, err := api.GetPaginated[[]api.GetUserResponse](ctx, d.Client, d.Token, url, api.SortById)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Error Reading Capella Users",
				"Could not read users in organization "+state.OrganizationId.String()+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Users",
			"Could not read users in organization "+state.OrganizationId.String()+": "+api.ParseError(err),
		)
		return
	}

	state, err = d.mapResponseBody(ctx, response, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading User",
			"Could not create User, unexpected error: "+err.Error(),
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

// Configure adds the provider configured client to the User data source.
func (d *Users) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
// get Users to the Users schema that will be used by terraform.
func (d *Users) mapResponseBody(
	ctx context.Context,
	UsersResponse []api.GetUserResponse,
	state *providerschema.Users,
) (providerschema.Users, error) {
	for _, userResp := range UsersResponse {
		audit := providerschema.NewCouchbaseAuditData(userResp.Audit)

		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			return *state, fmt.Errorf("error occured while attempting to convert audit data")
		}

		// Set Optional Values
		var name string
		if userResp.Name != nil {
			name = *userResp.Name
		}

		UserState := providerschema.NewUser(
			types.StringValue(userResp.Id.String()),
			types.StringValue(name),
			types.StringValue(userResp.Email),
			types.StringValue(userResp.Status),
			types.BoolValue(userResp.Inactive),
			types.StringValue(userResp.OrganizationId.String()),
			providerschema.MorphRoles(userResp.OrganizationRoles),
			types.StringValue(userResp.LastLogin),
			types.StringValue(userResp.Region),
			types.StringValue(userResp.TimeZone),
			types.BoolValue(userResp.EnableNotifications),
			types.StringValue(userResp.ExpiresAt),
			providerschema.MorphResources(userResp.Resources),
			auditObj,
		)

		state.Data = append(state.Data, *UserState)
	}
	return *state, nil
}

// validate is used to verify that all the fields in the datasource
// have been populated.
func (d *Users) validate(state providerschema.Users) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	return nil
}
