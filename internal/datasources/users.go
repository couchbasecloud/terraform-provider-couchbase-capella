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
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of then Capella organization.",
			},
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the user.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the user.",
						},
						"status": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "Status depicts user status whether they are verified or not." +
								"It can be one of the following values: verified, not-verified, pending-primary.",
						},
						"inactive": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Inactive depicts whether the user has accepted the invite for the organization.",
						},
						"email": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The email of the user.",
						},
						"organization_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the Capella organization.",
						},
						"organization_roles": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "The organization roles associated to the user. They determines the privileges user possesses in the organization.",
						},
						"last_login": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Time(UTC) at which user last logged in.",
						},
						"region": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The region of the user",
						},
						"time_zone": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Time zone of the user.",
						},
						"enable_notifications": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "After enabling email notifications for your account, you will start receiving email notification alerts from all databases in projects you are a part of.",
						},
						"expires_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Time at which the user expires.",
						},
						"resources": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Optional:            true,
										Computed:            true,
										MarkdownDescription: "Type of the resource.",
									},
									"id": schema.StringAttribute{Computed: true, MarkdownDescription: "The ID of the project."},
									"roles": schema.ListAttribute{
										ElementType:         types.StringType,
										Computed:            true,
										MarkdownDescription: "Project Roles associated with the User.",
									},
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
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]api.GetUserResponse](ctx, d.Client, d.Token, cfg, api.SortById)
	if err != nil {
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
			return *state, fmt.Errorf("error occurred while attempting to convert audit data")
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
