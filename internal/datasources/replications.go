package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
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

// Metadata returns the data source type name.
func (d *Replications) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replications"
}

// Schema defines the schema for the data source.
func (d *Replications) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ReplicationsSchema()
}

// Read refreshes the Terraform state with the latest data of replications.
func (d *Replications) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ReplicationsData
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationId.ValueString()
	projectId := state.ProjectId.ValueString()
	clusterId := state.ClusterId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/replications", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]replicationapi.ReplicationSummary](ctx, d.ClientV1, d.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Replications",
			fmt.Sprintf(
				"Could not read replications in organization %s, project %s, cluster %s: %s",
				organizationId, projectId, clusterId, api.ParseError(err),
			),
		)
		return
	}

	for i := range response {
		replicationSummary := response[i]

		// Build audit object
		audit := providerschema.ReplicationAuditData{
			CreatedBy:  types.StringValue(replicationSummary.Audit.CreatedBy),
			CreatedAt:  types.StringValue(replicationSummary.Audit.CreatedAt),
			ModifiedBy: types.StringValue(replicationSummary.Audit.ModifiedBy),
			ModifiedAt: types.StringValue(replicationSummary.Audit.ModifiedAt),
			Version:    types.Int64Value(replicationSummary.Audit.Version),
		}
		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error Reading Replications",
				fmt.Sprintf(
					"Could not read replications in organization %s, project %s, cluster %s: %s",
					organizationId, projectId, clusterId, errors.ErrUnableToConvertAuditData,
				),
			)
			return
		}

		newReplicationData := providerschema.NewReplicationSummaryData(replicationSummary, auditObj)
		state.Data = append(state.Data, newReplicationData)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *Replications) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.Data = data
}
