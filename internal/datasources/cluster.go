package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-capella/internal/api"
	clusterapi "terraform-provider-capella/internal/api/cluster"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Cluster{}
	_ datasource.DataSourceWithConfigure = &Cluster{}
)

// Cluster is the Cluster data source implementation.
type Cluster struct {
	*providerschema.Data
}

// NewCluster is a helper function to simplify the provider implementation.
func NewCluster() datasource.DataSource {
	return &Cluster{}
}

// Metadata returns the cluster data source type name.
func (d *Cluster) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clusters"
}

// Schema defines the schema for the Cluster data source.
func (d *Cluster) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ClusterSchema()
}

// Read refreshes the Terraform state with the latest data of clusters.
func (d *Cluster) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Clusters
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: organization ID cannot be empty.",
		)
		return
	}

	if state.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: project ID cannot be empty.",
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
	)

	response, err := d.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters", d.HostURL, organizationId, projectId),
		http.MethodGet,
		nil,
		d.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella Clusters",
			fmt.Sprintf("Could not read clusters in organization %s and project %s, unexpected error: %s", organizationId, projectId, err.CompleteError()),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Clusters",
			fmt.Sprintf("Could not read clusters in organization %s and project %s, unexpected error: %s", organizationId, projectId, err.Error()),
		)
		return
	}

	clusterResp := clusterapi.GetClustersResponse{}
	err = json.Unmarshal(response.Body, &clusterResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: "+err.Error(),
		)
		return
	}

	for _, cluster := range clusterResp.Data {
		audit := providerschema.NewCouchbaseAuditData(cluster.Audit)

		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error Reading Capella Clusters",
				fmt.Sprintf("Could not read clusters in organization %s and project %s, unexpected error: %s", organizationId, projectId, errors.ErrUnableToConvertAuditData),
			)
		}

		newClusterData, err := providerschema.NewClusterData(&cluster, organizationId, projectId, auditObj)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Capella Clusters",
				fmt.Sprintf("Could not read clusters in organization %s and project %s, unexpected error: %s", organizationId, projectId, err.Error()),
			)
		}
		state.Data = append(state.Data, *newClusterData)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the cluster data source.
func (d *Cluster) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
