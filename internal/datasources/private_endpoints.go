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
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &PrivateEndpoints{}
	_ datasource.DataSourceWithConfigure = &PrivateEndpoints{}
)

// PrivateEndpoints is the data source implementation.
type PrivateEndpoints struct {
	*providerschema.Data
}

// NewPrivateEndpoints is a helper function to simplify the provider implementation.
func NewPrivateEndpoints() datasource.DataSource {
	return &PrivateEndpoints{}
}

// Metadata returns the data source type name.
func (p *PrivateEndpoints) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_endpoints"
}

// Schema defines the schema for the private endpoint data source.
func (p *PrivateEndpoints) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The data source to retrieve private endpoints for an operational cluster. Access your operational cluster from your Cloud Service Provider's private network. Returns a list of private endpoints associated with the endpoint service for your operational cluster, along with the endpoint state. Each private endpoint connects a private network to the operational cluster.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster. Private endpoints enable secure access to this cluster through your Cloud Service Provider's private network.",
			},
			"data": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Lists the private endpoints associated with the cluster. Each entry represents a connection point between your private network and the Capella cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the private endpoint.",
						},
						"status": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "The current status of the private endpoint. Possible values are:\n" +
								"* `pending` - The endpoint creation is in progress\n" +
								"* `pendingAcceptance` - The endpoint is waiting for acceptance from Capella\n" +
								"* `linked` - The endpoint is successfully connected and active\n" +
								"* `rejected` - The endpoint connection request was rejected\n" +
								"* `unrecognized` - The endpoint state cannot be determined\n" +
								"* `failed` - The endpoint creation or connection attempt failed",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of private endpoint.
func (p *PrivateEndpoints) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.PrivateEndpoints
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := p.validate(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Private Endpoints",
			"Could not read private endpoints in cluster "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/privateEndpointService/endpoints", p.HostURL, organizationId, projectId, clusterId)
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
			"Error Reading Capella Private Endpoints",
			"Could not read private endpoints in cluster "+state.ClusterId.String()+": "+api.ParseError(err),
		)
		return
	}

	privateEndpointsResp := api.GetPrivateEndpointsResponse{}
	err = json.Unmarshal(response.Body, &privateEndpointsResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshalling private endpoints",
			"Could not read private endpoints, unexpected error: "+err.Error(),
		)
		return
	}

	for _, e := range privateEndpointsResp.Endpoints {
		endpointData := providerschema.PrivateEndpointData{}
		endpointData.Id = types.StringValue(e.Id)
		endpointData.Status = types.StringValue(e.Status)
		state.Data = append(state.Data, endpointData)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the private endpoint data source.
func (p *PrivateEndpoints) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// validate ensures organization id, project id and cluster id are valued.
func (p *PrivateEndpoints) validate(state providerschema.PrivateEndpoints) error {
	if state.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if state.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if state.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	return nil
}
