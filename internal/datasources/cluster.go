package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
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
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

// Schema defines the schema for the Cluster data source.
func (d *Cluster) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = OneClusterSchema()
}

// Read refreshes the Terraform state with the latest data of a single cluster.
func (d *Cluster) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Cluster
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Cluster",
			"Could not read cluster: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.Id]
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		d.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Cluster",
			fmt.Sprintf("Could not read cluster %s: %s", clusterId, api.ParseError(err)),
		)
		return
	}

	var clusterResponse clusterapi.GetClusterResponse
	err = json.Unmarshal(response.Body, &clusterResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Cluster",
			fmt.Sprintf("Could not unmarshal cluster %s, unexpected error: %s", clusterId, err.Error()),
		)
		return
	}
	clusterResponse.Etag = response.Response.Header.Get("ETag")

	audit := providerschema.NewCouchbaseAuditData(clusterResponse.Audit)
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error Reading Capella Cluster",
			fmt.Sprintf("Could not convert audit data for cluster %s", clusterId),
		)
		return
	}

	newClusterState, err := providerschema.NewCluster(ctx, &clusterResponse, organizationId, projectId, auditObj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Cluster",
			fmt.Sprintf("Could not read cluster %s, unexpected error: %s", clusterId, err.Error()),
		)
		return
	}

	diags = resp.State.Set(ctx, newClusterState)
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
