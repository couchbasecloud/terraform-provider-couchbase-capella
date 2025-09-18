package datasources

import (
	"context"
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/apigen"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/google/uuid"
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
	FreeTierClusterFilter bool
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

	orgUUID, _ := uuid.Parse(state.OrganizationId.ValueString())
	projUUID, _ := uuid.Parse(state.ProjectId.ValueString())

	listResp, err := d.ClientV2.ListClustersWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Clusters",
			fmt.Sprintf(
				"Could not read clusters in organization %s and project %s, unexpected error: %s",
				state.OrganizationId.ValueString(), state.ProjectId.ValueString(), api.ParseError(err),
			),
		)
		return
	}
	if listResp.JSON200 == nil {
		resp.Diagnostics.AddError("Error Reading Capella Clusters", "unexpected response status: "+listResp.Status())
		return
	}

	for i := range listResp.JSON200.Data {
		cluster := listResp.JSON200.Data[i]
		if d.FreeTierClusterFilter {
			if cluster.Support.Plan != "free" {
				continue
			}
		}
		audit := providerschema.CouchbaseAuditData{
			CreatedAt:  types.StringValue(cluster.Audit.CreatedAt.String()),
			CreatedBy:  types.StringValue(cluster.Audit.CreatedBy),
			ModifiedAt: types.StringValue(cluster.Audit.ModifiedAt.String()),
			ModifiedBy: types.StringValue(cluster.Audit.ModifiedBy),
			Version:    types.Int64Value(int64(cluster.Audit.Version)),
		}

		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error Reading Capella Clusters",
				fmt.Sprintf("Could not read clusters in organization %s and project %s, unexpected error: %s", state.OrganizationId.ValueString(), state.ProjectId.ValueString(), errors.ErrUnableToConvertAuditData),
			)
		}

		cd := providerschema.ClusterData{
			Id:               types.StringValue(cluster.Id.String()),
			OrganizationId:   types.StringValue(state.OrganizationId.ValueString()),
			ProjectId:        types.StringValue(state.ProjectId.ValueString()),
			Name:             types.StringValue(cluster.Name),
			Description:      types.StringValue(cluster.Description),
			Audit:            auditObj,
			CurrentState:     types.StringValue(string(cluster.CurrentState)),
			ConnectionString: types.StringValue(cluster.ConnectionString),
		}
		state.Data = append(state.Data, cd)
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
