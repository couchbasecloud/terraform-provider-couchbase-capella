package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &PrivateEndpointService{}
	_ datasource.DataSourceWithConfigure = &PrivateEndpointService{}
)

// PrivateEndpointService is the data source implementation.
type PrivateEndpointService struct {
	*providerschema.Data
}

// NewPrivateEndpointService is a helper function to simplify the provider implementation.
func NewPrivateEndpointService() datasource.DataSource {
	return &PrivateEndpointService{}
}

// Metadata returns the data source type name.
func (p *PrivateEndpointService) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_endpoint_service"
}

func (p *PrivateEndpointService) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"cluster_id":      requiredStringAttribute,
			"enabled":         computedBoolAttribute,
		},
	}
}

func (p *PrivateEndpointService) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.PrivateEndpointService
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := state.Validate()
	if err != nil {
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Capella Private Endpoint Service",
				"Could not read private endpoint service in cluster "+state.ClusterId.String()+": "+err.Error(),
			)
			return
		}
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService", p.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := p.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Private Endpoint Service",
			"Could not read private endpoint service in cluster "+state.OrganizationId.String()+": "+api.ParseError(err),
		)
		return
	}

	privateEndpointServiceStatus := api.GetPrivateEndpointServiceStatusResponse{}
	err = json.Unmarshal(response.Body, &privateEndpointServiceStatus)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling private endpoint service status",
			"Could not unmarshall private endpoint service status, unexpected error: "+err.Error(),
		)
		return
	}

	state.Enabled = types.BoolValue(privateEndpointServiceStatus.Enabled)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *PrivateEndpointService) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	p.Data = data
}
