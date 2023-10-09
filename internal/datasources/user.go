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
	_ datasource.DataSource              = &User{}
	_ datasource.DataSourceWithConfigure = &User{}
)

// User is the user data source implementation.
type User struct {
	*providerschema.Data
}

// NewUser is a helper function to simplify the provider implementation.
func NewUser() datasource.DataSource {
	return &User{}
}

// Metadata returns the user data source type name.
func (d *User) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the User data source.
func (d *User) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              computedStringAttribute(),
						"name":            computedStringAttribute(),
						"status":          computedStringAttribute(),
						"inactive":        computedBoolAttribute(),
						"email":           computedStringAttribute(),
						"organization_id": computedStringAttribute(),
						"organization_roles": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"last_login":           computedStringAttribute(),
						"region":               computedStringAttribute(),
						"time_zone":            computedStringAttribute(),
						"enable_notifications": computedBoolAttribute(),
						"expires_at":           computedStringAttribute(),
						"resources": schema.ListNestedAttribute{
							Required: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": computedStringAttribute(),
									"id":   computedStringAttribute(),
									"roles": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
								},
							},
						},
						"audit": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"created_at":  computedStringAttribute(),
								"created_by":  computedStringAttribute(),
								"modified_at": computedStringAttribute(),
								"modified_by": computedStringAttribute(),
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

// Read refreshes the Terraform state with the latest data of Users.
func (d *User) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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

	var (
		organizationId = state.OrganizationId.ValueString()
	)

	// Make request to list Users
	response, err := d.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/users", d.HostURL, organizationId),
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
			"Error Reading Users",
			"Could not read users in organization "+state.OrganizationId.String()+": "+err.Error(),
		)
		return
	}

	UsersResponse := api.GetUsersResponse{}
	err = json.Unmarshal(response.Body, &UsersResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading User",
			"Could not create User, unexpected error: "+err.Error(),
		)
		return
	}

	state, err = d.mapResponseBody(ctx, UsersResponse, &state)
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
func (d *User) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *User) mapResponseBody(
	ctx context.Context,
	UsersResponse api.GetUsersResponse,
	state *providerschema.Users,
) (providerschema.Users, error) {
	for _, userResp := range UsersResponse.Data {
		audit := providerschema.NewCouchbaseAuditData(userResp.Audit)

		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			return *state, fmt.Errorf("error occured while attempting to convert audit data")
		}

		UserState := providerschema.NewUser(
			types.StringValue(userResp.Id.String()),
			types.StringValue(*userResp.Name),
			types.StringValue(userResp.Email),
			types.StringValue(userResp.Status),
			types.BoolValue(userResp.Inactive),
			types.StringValue(userResp.OrganizationId.String()),
			providerschema.MorphOrganizationRoles(*userResp.OrganizationRoles),
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

// computedStringAttribute returns a Terraform schema attribute
// which is configured to be computed.
func computedStringAttribute() *schema.StringAttribute {
	return &schema.StringAttribute{
		Computed: true,
	}
}

// computedBoolAttribute returns a Terraform schema attribute
// which is configured to be computed.
func computedBoolAttribute() *schema.BoolAttribute {
	return &schema.BoolAttribute{
		Computed: true,
	}
}

// validate is used to verify that all the fields in the datasource
// have been populated.
func (d *User) validate(state providerschema.Users) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	return nil
}
