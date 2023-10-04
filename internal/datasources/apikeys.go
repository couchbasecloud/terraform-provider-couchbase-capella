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
	_ datasource.DataSource              = &ApiKey{}
	_ datasource.DataSourceWithConfigure = &ApiKey{}
)

// ApiKey is the api key data source implementation.
type ApiKey struct {
	*providerschema.Data
}

// NewApiKey is a helper function to simplify the provider implementation.
func NewApiKey() datasource.DataSource {
	return &ApiKey{}
}

// Metadata returns the api key data source type name.
func (d *ApiKey) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apikeys"
}

// Schema defines the schema for the api key data source.
func (d *ApiKey) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"expiry": schema.Float64Attribute{
							Computed: true,
						},
						"allowed_cidrs": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"organization_roles": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"resources": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"roles": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
									"type": schema.StringAttribute{
										Computed: true,
									},
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

// Read refreshes the Terraform state with the latest data of api keys.
func (d *ApiKey) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ApiKeys
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading ApiKeys in Capella",
			"Could not read Capella ApiKeys in organization "+organizationId+": "+err.Error(),
		)
		return
	}

	response, err := d.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/ApiKeys", d.HostURL, organizationId),
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
				"Error Reading Capella ApiKeys",
				"Could not read ApiKeys in organization "+organizationId+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella ApiKeys",
			"Could not read ApiKeys in organization "+organizationId+": "+err.Error(),
		)
		return
	}

	ApiKeyResp := api.GetApiKeysResponse{}
	err = json.Unmarshal(response.Body, &ApiKeyResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing ApiKeys",
			"Could not list ApiKeys, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to model
	for _, ApiKey := range ApiKeyResp.Data {
		ApiKeyState := providerschema.OneApiKey{
			Audit: providerschema.CouchbaseAuditData{
				CreatedAt:  types.StringValue(ApiKey.Audit.CreatedAt.String()),
				CreatedBy:  types.StringValue(ApiKey.Audit.CreatedBy),
				ModifiedAt: types.StringValue(ApiKey.Audit.ModifiedAt.String()),
				ModifiedBy: types.StringValue(ApiKey.Audit.ModifiedBy),
				Version:    types.Int64Value(int64(ApiKey.Audit.Version)),
			},
			Id:                  types.StringValue(ApiKey.Id.String()),
			Name:                types.StringPointerValue(ApiKey.Name),
			Email:               types.StringValue(ApiKey.Email),
			Status:              types.StringValue(ApiKey.Status),
			Inactive:            types.BoolValue(ApiKey.Inactive),
			OrganizationId:      types.StringValue(ApiKey.OrganizationId.String()),
			LastLogin:           types.StringValue(ApiKey.LastLogin),
			Region:              types.StringValue(ApiKey.Region),
			TimeZone:            types.StringValue(ApiKey.TimeZone),
			EnableNotifications: types.BoolValue(ApiKey.EnableNotifications),
			ExpiresAt:           types.StringValue(ApiKey.ExpiresAt),
		}
		state.Data = append(state.Data, ApiKeyState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the api key data source.
func (d *ApiKey) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
