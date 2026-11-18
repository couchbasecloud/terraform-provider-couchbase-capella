package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
)

// ReplicationsDataDescription is the fallback description of the data attribute of the replications data source.
const ReplicationsDataDescription = "\n - The replications, sorted by ID. Each direction of a replication is one entry, " +
	"and only the cluster where it starts lists it. Thus, without cluster_id, a `twoWay` replication shows as two entries."

// Replications is the Terraform state of the replications data source.
type Replications struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`

	Data []ReplicationSummaryData `tfsdk:"data"`
}

// ReplicationSummaryData is one entry of the replications data source.
type ReplicationSummaryData struct {
	Id            types.String `tfsdk:"id"`
	SourceCluster types.String `tfsdk:"source_cluster"`
	TargetCluster types.String `tfsdk:"target_cluster"`
	Status        types.String `tfsdk:"status"`
	Direction     types.String `tfsdk:"direction"`
	Audit         types.Object `tfsdk:"audit"`
}

// NewReplicationSummaryData builds one entry of the replications data source from a list entry.
func NewReplicationSummaryData(ctx context.Context, summary replicationapi.ReplicationSummary) (ReplicationSummaryData, error) {
	audit, err := newReplicationAudit(ctx, summary.Audit)
	if err != nil {
		return ReplicationSummaryData{}, err
	}

	return ReplicationSummaryData{
		Id:            types.StringValue(summary.Id),
		SourceCluster: types.StringValue(summary.SourceCluster),
		TargetCluster: types.StringValue(summary.TargetCluster),
		Status:        types.StringValue(string(summary.Status)),
		Direction:     types.StringValue(summary.Direction),
		Audit:         audit,
	}, nil
}
