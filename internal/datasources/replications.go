package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	replication_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Replications{}
	_ datasource.DataSourceWithConfigure = &Replications{}
)

// Replications is the Replications data source implementation.
type Replications struct {
	*providerschema.Data
}

// NewReplications is a helper function to simplify the provider implementation.
func NewReplications() datasource.DataSource {
	return &Replications{}
}

// Metadata returns the replications data source type name.
func (r *Replications) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replications"
}

// Schema defines the schema for the Replications data source.
func (r *Replications) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ReplicationsSchema()
}

// Read refreshes the Terraform state with the latest data of replications.
func (r *Replications) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ReplicationsData
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Replications in Capella",
			"Could not read Capella Replications in cluster "+clusterId+": "+err.Error(),
		)
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/replications", r.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]replication_api.GetReplicationSummaryResponse](ctx, r.ClientV1, r.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Replications",
			fmt.Sprintf(
				"Could not read replications in organization %s, project %s, and cluster %s, unexpected error: %s",
				organizationId, projectId, clusterId, api.ParseError(err),
			),
		)
		return
	}

	for i := range response {
		replication := response[i]
		audit := providerschema.ReplicationAuditData{
			CreatedAt: types.StringValue(replication.Audit.CreatedAt.String()),
			CreatedBy: types.StringValue(replication.Audit.CreatedBy),
		}

		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error Reading Replication",
				fmt.Sprintf("Could not read replication list in organization %s, project %s, and cluster %s, unexpected error: %s", organizationId, projectId, clusterId, "unable to convert audit data"),
			)
			continue
		}

		newReplicationData := providerschema.NewReplicationSummaryData(replication, auditObj)
		state.Data = append(state.Data, *newReplicationData)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (r *Replications) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	r.Data = data
}
