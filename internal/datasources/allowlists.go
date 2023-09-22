package datasources

import (
	"context"
	"fmt"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	// TODO: Implement Read logic
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
