package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-capella/internal/api"
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

// Schema defines the schema for the user data source.
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
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Optional: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"inactive": schema.BoolAttribute{
							Computed: true,
						},
						"email": schema.StringAttribute{
							Required: true,
						},
						"organization_id": schema.StringAttribute{
							Required: true,
						},
						"last_login": schema.StringAttribute{
							Computed: true,
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"time_zone": schema.StringAttribute{
							Computed: true,
						},
						"enable_notifications": schema.BoolAttribute{
							Computed: true,
						},
						"expires_at": schema.StringAttribute{
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
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of users.
func (d *User) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Users
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Users in Capella",
			"Could not read Capella users in organization "+organizationId+": "+err.Error(),
		)
		return
	}

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
		if err.HttpStatusCode != 404 {
			resp.Diagnostics.AddError(
				"Error Reading Capella Users",
				"Could not read users in organization "+organizationId+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Users",
			"Could not read users in organization "+organizationId+": "+err.Error(),
		)
		return
	}

	userResp := api.GetUsersResponse{}
	err = json.Unmarshal(response.Body, &userResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing users",
			"Could not list users, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to model
	for _, user := range userResp.Data {
		userState := providerschema.OneUser{
			Audit: providerschema.CouchbaseAuditData{
				CreatedAt:  types.StringValue(user.Audit.CreatedAt.String()),
				CreatedBy:  types.StringValue(user.Audit.CreatedBy),
				ModifiedAt: types.StringValue(user.Audit.ModifiedAt.String()),
				ModifiedBy: types.StringValue(user.Audit.ModifiedBy),
				Version:    types.Int64Value(int64(user.Audit.Version)),
			},
			Id:                  types.StringValue(user.Id.String()),
			Name:                types.StringPointerValue(user.Name),
			Email:               types.StringValue(user.Email),
			Status:              types.StringValue(user.Status),
			Inactive:            types.BoolValue(user.Inactive),
			OrganizationId:      types.StringValue(user.OrganizationId.String()),
			LastLogin:           types.StringValue(user.LastLogin),
			Region:              types.StringValue(user.Region),
			TimeZone:            types.StringValue(user.TimeZone),
			EnableNotifications: types.BoolValue(user.EnableNotifications),
			ExpiresAt:           types.StringValue(user.ExpiresAt),
		}
		state.Data = append(state.Data, userState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the user data source.
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
