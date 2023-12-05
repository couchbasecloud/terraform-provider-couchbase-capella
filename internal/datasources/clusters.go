package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Clusters{}
	_ datasource.DataSourceWithConfigure = &Clusters{}
)

// Clusters is the Clusters data source implementation.
type Clusters struct {
	*providerschema.Data
}

// NewClusters is a helper function to simplify the provider implementation.
func NewClusters() datasource.DataSource {
	return &Clusters{}
}

// Metadata returns the cluster data source type name.
func (d *Clusters) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clusters"
}

// Schema defines the schema for the Clusters data source.
func (d *Clusters) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ClusterSchema()
}

// Read refreshes the Terraform state with the latest data of clusters.
func (d *Clusters) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Clusters
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			"Could not read cluster, unexpected error: organization ID cannot be empty.",
		)
		return
	}

	if state.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			"Could not read cluster, unexpected error: project ID cannot be empty.",
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters", d.HostURL, organizationId, projectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]clusterapi.GetClusterResponse](ctx, d.Client, d.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Clusters",
			fmt.Sprintf(
				"Could not read clusters in organization %s and project %s, unexpected error: %s",
				organizationId, projectId, api.ParseError(err),
			),
		)
		return
	}

	for _, cluster := range response {
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
func (d *Clusters) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
