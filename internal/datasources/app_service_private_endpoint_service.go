package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AppServicePrivateEndpointService{}
	_ datasource.DataSourceWithConfigure = &AppServicePrivateEndpointService{}
)

// AppServicePrivateEndpointService is the data source implementation.
type AppServicePrivateEndpointService struct {
	*providerschema.Data
}

// NewAppServicePrivateEndpointService is a helper function to simplify the provider implementation.
func NewAppServicePrivateEndpointService() datasource.DataSource {
	return &AppServicePrivateEndpointService{}
}

// Metadata returns the data source type name.
func (p *AppServicePrivateEndpointService) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_private_endpoint_service"
}

// Schema defines schema for the App Service private endpoint service data source.
func (p *AppServicePrivateEndpointService) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServicePrivateEndpointServiceSchema()
}

// Read refreshes the Terraform state with the latest data of the App Service private endpoint service.
func (p *AppServicePrivateEndpointService) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServicePrivateEndpointService
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := p.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating Capella App Service Private Endpoint Service",
			"Could not validate App Service private endpoint service on app service "+state.AppServiceId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
		appServiceId   = state.AppServiceId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService", p.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := p.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		p.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Service Private Endpoint Service",
			"Could not read App Service private endpoint service on app service "+state.AppServiceId.String()+": "+api.ParseError(err),
		)
		return
	}

	privateEndpointServiceStatus := api.GetAppServicePrivateEndpointServiceStatusResponse{}
	err = json.Unmarshal(response.Body, &privateEndpointServiceStatus)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling App Service private endpoint service status",
			"Could not unmarshall App Service private endpoint service status, unexpected error: "+err.Error(),
		)
		return
	}

	state.Enabled = types.BoolValue(privateEndpointServiceStatus.State != nil && *privateEndpointServiceStatus.State == "enabled")
	state.State = types.StringNull()
	if privateEndpointServiceStatus.State != nil {
		state.State = types.StringValue(*privateEndpointServiceStatus.State)
	}
	state.TargetState = types.StringNull()
	if privateEndpointServiceStatus.TargetState != nil {
		state.TargetState = types.StringValue(*privateEndpointServiceStatus.TargetState)
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the App Service private endpoint service data source.
func (p *AppServicePrivateEndpointService) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// validate ensures organization id, project id, cluster id, and app service id are valued.
func (p *AppServicePrivateEndpointService) validate(state providerschema.AppServicePrivateEndpointService) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if state.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if state.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if state.AppServiceId.IsNull() {
		return errors.ErrAppServiceIdMissing
	}
	return nil
}
