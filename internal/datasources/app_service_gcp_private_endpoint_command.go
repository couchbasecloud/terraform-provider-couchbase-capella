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

var (
	_ datasource.DataSource              = &AppServiceGCPPrivateEndpointCommand{}
	_ datasource.DataSourceWithConfigure = &AppServiceGCPPrivateEndpointCommand{}
)

// AppServiceGCPPrivateEndpointCommand is the data source implementation.
//
// NOTE: as of this writing, Couchbase Capella's private endpoints for App Services are
// AWS-only; Azure/GCP support is not confirmed GA (the generated OpenAPI client models the
// connection-command request body as an unresolved oneOf/anyOf union, consistent with the
// schema being written generically/forward-compatible even if only AWS is wired up in the
// backend today). This data source is included for API/schema parity with the cluster-level
// gcp_private_endpoint_command resource, but has not been validated against a live GCP-backed
// App Service.
type AppServiceGCPPrivateEndpointCommand struct {
	*providerschema.Data
}

// NewAppServiceGCPPrivateEndpointCommand is a helper function to simplify the provider implementation.
func NewAppServiceGCPPrivateEndpointCommand() datasource.DataSource {
	return &AppServiceGCPPrivateEndpointCommand{}
}

// Metadata returns the data source type name.
func (a *AppServiceGCPPrivateEndpointCommand) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_gcp_private_endpoint_command"
}

// Schema defines the schema for the App Service private endpoint command data source.
func (a *AppServiceGCPPrivateEndpointCommand) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceGcpPrivateEndpointCommandSchema()
}

// Read refreshes the Terraform state with the latest data of the App Service private endpoint command.
func (a *AppServiceGCPPrivateEndpointCommand) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServiceGCPCommandRequest
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating App Service GCP private endpoint command request",
			"Could not validate App Service GCP private endpoint command request: "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
		appServiceId   = state.AppServiceId.ValueString()
	)

	GCPCommandRequest := api.CreateGCPEndpointCommandRequest{
		VpcNetworkID: state.VpcNetworkID.ValueString(),
		SubnetIDs:    convertSubnetIDs(state.SubnetIDs),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService/privateEndpointCommand", a.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusOK}
	response, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		GCPCommandRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service GCP private endpoint command",
			"Could not read App Service GCP private endpoint command: "+api.ParseError(err),
		)
		return
	}

	var GCPCommandResponse api.CreatePrivateEndpointCommandResponse
	err = json.Unmarshal(response.Body, &GCPCommandResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling App Service GCP private endpoint command response",
			"Could not unmarshall App Service GCP private endpoint command response, unexpected error: "+err.Error(),
		)
		return
	}

	state.Command = types.StringValue(GCPCommandResponse.Command)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the App Service private endpoint command data source.
func (a *AppServiceGCPPrivateEndpointCommand) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	a.Data = data
}

// validate ensures organization id, project id, cluster id, app service id, and VPC Network id are valued.
func (a *AppServiceGCPPrivateEndpointCommand) validate(config providerschema.AppServiceGCPCommandRequest) error {
	if config.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if config.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if config.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if config.AppServiceId.IsNull() {
		return errors.ErrAppServiceIdMissing
	}
	if config.VpcNetworkID.IsNull() {
		return errors.ErrVPCIDMissing
	}
	return nil
}
