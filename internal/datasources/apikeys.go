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
			"organization_id": requiredStringAttribute,
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                 computedStringAttribute,
						"organization_id":    computedStringAttribute,
						"name":               computedStringAttribute,
						"description":        computedStringAttribute,
						"expiry":             computedFloat64Attribute,
						"allowed_cidrs":      computedListAttribute,
						"organization_roles": computedListAttribute,
						"resources": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id":    computedStringAttribute,
									"roles": computedListAttribute,
									"type":  computedStringAttribute,
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
			"Error Reading Api Keys in Capella",
			"Could not read Capella api keys in organization "+organizationId+": "+err.Error(),
		)
		return
	}

	response, err := d.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/apikeys", d.HostURL, organizationId),
		http.MethodGet,
		nil,
		d.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella ApiKeys",
			"Could not read api keys in organization "+organizationId+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella ApiKeys",
			"Could not read api keys in organization "+organizationId+": "+err.Error(),
		)
		return
	}

	apiKeyResp := api.GetApiKeysResponse{}
	err = json.Unmarshal(response.Body, &apiKeyResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing ApiKeys",
			"Could not list api keys, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to model
	for _, apiKey := range apiKeyResp.Data {
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
