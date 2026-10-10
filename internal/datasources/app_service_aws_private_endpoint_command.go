package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = &AppServiceAWSPrivateEndpointCommand{}
	_ datasource.DataSourceWithConfigure = &AppServiceAWSPrivateEndpointCommand{}
)

// awsServiceNameFromCommandPattern matches the value following the AWS CLI
// `--service-name` flag in the generated `aws ec2 create-vpc-endpoint ...` command, e.g.
// "com.amazonaws.vpce.us-east-1.vpce-svc-0823b61a6d8cee231" out of
// "... --service-name com.amazonaws.vpce.us-east-1.vpce-svc-0823b61a6d8cee231 --vpc-endpoint-type ...".
var awsServiceNameFromCommandPattern = regexp.MustCompile(`--service-name\s+(\S+)`)

// extractAWSServiceNameFromCommand parses the AWS VPC endpoint service name out of the
// free-text AWS CLI command returned by the App Service private endpoint command API. There
// is no dedicated API field for this (the App Service private endpoint status response, unlike
// the operational cluster's, does not return a service name), so this is a best-effort,
// server-side convenience so the value can be fed directly into an `aws_vpc_endpoint` resource's
// `service_name` argument instead of requiring an HCL-level regex.
//
// This is deliberately best-effort: if Couchbase changes the command's phrasing and the pattern
// no longer matches, ok is false and the caller must leave the attribute null rather than fail
// the read — the command string itself is still returned untouched either way.
func extractAWSServiceNameFromCommand(command string) (serviceName string, ok bool) {
	match := awsServiceNameFromCommandPattern.FindStringSubmatch(command)
	if len(match) != 2 {
		return "", false
	}
	return match[1], true
}

// AppServiceAWSPrivateEndpointCommand is the data source implementation.
type AppServiceAWSPrivateEndpointCommand struct {
	*providerschema.Data
}

// NewAppServiceAWSPrivateEndpointCommand is a helper function to simplify the provider implementation.
func NewAppServiceAWSPrivateEndpointCommand() datasource.DataSource {
	return &AppServiceAWSPrivateEndpointCommand{}
}

// Metadata returns the data source type name.
func (a *AppServiceAWSPrivateEndpointCommand) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_aws_private_endpoint_command"
}

// Schema defines the schema for the App Service private endpoint command data source.
func (a *AppServiceAWSPrivateEndpointCommand) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceAWSPrivateEndpointCommandSchema()
}

// Read refreshes the Terraform state with the latest data of the App Service private endpoint command.
func (a *AppServiceAWSPrivateEndpointCommand) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServiceAWSCommandRequest
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating App Service AWS private endpoint command request",
			"Could not validate App Service AWS private endpoint command request: "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
		appServiceId   = state.AppServiceId.ValueString()
	)

	AWSCommandRequest := api.CreateVPCEndpointCommandRequest{
		VpcID:     state.VpcID.ValueString(),
		SubnetIDs: convertSubnetIDs(state.SubnetIDs),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/privateEndpointService/privateEndpointCommand", a.HostURL, organizationId, projectId, clusterId, appServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusOK}
	response, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		AWSCommandRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service AWS private endpoint command",
			"Could not read App Service AWS private endpoint command: "+api.ParseError(err),
		)
		return
	}

	var AWSCommandResponse api.CreatePrivateEndpointCommandResponse
	err = json.Unmarshal(response.Body, &AWSCommandResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling App Service AWS private endpoint command response",
			"Could not unmarshall App Service AWS private endpoint command response, unexpected error: "+err.Error(),
		)
		return
	}

	state.Command = types.StringValue(AWSCommandResponse.Command)
	if serviceName, ok := extractAWSServiceNameFromCommand(AWSCommandResponse.Command); ok {
		state.ServiceName = types.StringValue(serviceName)
	} else {
		state.ServiceName = types.StringNull()
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the App Service private endpoint command data source.
func (a *AppServiceAWSPrivateEndpointCommand) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// validate ensures organization id, project id, cluster id, app service id, and VPC id are valued.
func (a *AppServiceAWSPrivateEndpointCommand) validate(config providerschema.AppServiceAWSCommandRequest) error {
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
	if config.VpcID.IsNull() {
		return errors.ErrVPCIDMissing
	}
	return nil
}
