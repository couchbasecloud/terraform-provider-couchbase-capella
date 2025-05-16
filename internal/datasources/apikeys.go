package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ApiKeys{}
	_ datasource.DataSourceWithConfigure = &ApiKeys{}
)

// ApiKeys is the api key data source implementation.
type ApiKeys struct {
	*providerschema.Data
}

// NewApiKeys is a helper function to simplify the provider implementation.
func NewApiKeys() datasource.DataSource {
	return &ApiKeys{}
}

// Metadata returns the api key data source type name.
func (d *ApiKeys) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apikeys"
}

// Schema defines the schema for the api key data source.
func (d *ApiKeys) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The data source to retrieve API keys in an organization. API keys are used to authenticate and authorize access to Capella resources and services.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                 schema.StringAttribute{Computed: true, MarkdownDescription: "The unique identifier for the API key."},
						"organization_id":    schema.StringAttribute{Computed: true, MarkdownDescription: "The GUID4 ID of the organization."},
						"name":               schema.StringAttribute{Computed: true, MarkdownDescription: "Name of the API key."},
						"description":        schema.StringAttribute{Computed: true, MarkdownDescription: "A description for the API key."},
						"expiry":             schema.Float64Attribute{Computed: true, MarkdownDescription: "Expiry of the API key in number of days. If set to -1, the token will not expire."},
						"allowed_cidrs":      schema.ListAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "The list of inbound CIDRs for the API key. The system making a request must come from one of the allowed CIDRs.."},
						"organization_roles": schema.ListAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "List of the organization roles associated with the API key."},
						"resources": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Resource level permissions associated with the API key.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id":    schema.StringAttribute{Computed: true, MarkdownDescription: "The GUID4 ID of the project."},
									"roles": schema.ListAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Project Roles associated with the API key."},
									"type":  schema.StringAttribute{Computed: true, MarkdownDescription: "Resource type."},
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

// Read refreshes the Terraform state with the latest data of api keys.
func (d *ApiKeys) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ApiKeys
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Api Keys in Capella",
			"Could not read Capella api keys in organization "+organizationId+": "+err.Error(),
		)
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/apikeys", d.HostURL, organizationId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]api.GetApiKeyResponse](ctx, d.Client, d.Token, cfg, api.SortByName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella ApiKeys",
			"Could not read api keys in organization "+organizationId+": "+api.ParseError(err),
		)
		return
	}

	for i := range response {
		apiKey := response[i]
		audit := providerschema.NewCouchbaseAuditData(apiKey.Audit)

		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error listing ApiKeys",
				fmt.Sprintf("Could not list api keys, unexpected error: %s", fmt.Errorf("error while audit conversion")),
			)
			return
		}
		newApiKeyData, err := providerschema.NewApiKeyData(&apiKey, organizationId, auditObj)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing ApiKeys",
				fmt.Sprintf("Could not list api keys, unexpected error: %s", err.Error()),
			)
			return
		}
		state.Data = append(state.Data, newApiKeyData)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the api key data source.
func (d *ApiKeys) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
