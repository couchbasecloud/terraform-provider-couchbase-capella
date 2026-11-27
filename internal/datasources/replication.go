package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ datasource.DataSource              = (*Replication)(nil)
	_ datasource.DataSourceWithConfigure = (*Replication)(nil)
)

// Replication is the replication data source implementation.
type Replication struct {
	*providerschema.Data
}

// NewReplication is a helper function to simplify the provider implementation.
func NewReplication() datasource.DataSource {
	return &Replication{}
}

// Metadata returns the replication data source type name.
func (d *Replication) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication"
}

// Schema defines the schema for the replication data source.
func (d *Replication) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ReplicationSchema()
}

// Configure adds the provider configured client to the replication data source.
func (d *Replication) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read gets the replication and writes it to the state.
func (d *Replication) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.ReplicationData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = config.OrganizationId.ValueString()
		projectId      = config.ProjectId.ValueString()
		clusterId      = config.ClusterId.ValueString()
		replicationId  = config.Id.ValueString()
	)

	cfg := api.EndpointCfg{
		Url:           replicationapi.ReplicationURL(d.HostURL, organizationId, projectId, clusterId, replicationId),
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}
	response, err := d.ClientV1.ExecuteWithRetry(ctx, cfg, nil, d.Token, nil)
	if err != nil {
		detail := "Could not read replication " + replicationId + ": " + api.ParseError(err)
		if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
			detail += "\n\nCheck that cluster_id is the source cluster of this replication. For the reverse direction " +
				"of a twoWay replication, use the reverse replication ID and the target cluster ID."
		}
		resp.Diagnostics.AddError("Error reading replication", detail)
		return
	}

	replicationResp := replicationapi.GetReplicationResponse{}
	if err = json.Unmarshal(response.Body, &replicationResp); err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication",
			"Could not read replication "+replicationId+": "+errors.ErrUnmarshallingResponse.Error()+": "+err.Error(),
		)
		return
	}

	state, err := providerschema.NewReplicationData(ctx, &replicationResp, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading replication",
			"Could not read replication "+replicationId+": "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
