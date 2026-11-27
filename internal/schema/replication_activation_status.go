package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// ReplicationActivationStatusStateDescription is the fallback description of state.
const ReplicationActivationStatusStateDescription = "\n - The requested status: `running` or `paused`. A replication " +
	"that is starting shows as `running`, and one that is pausing shows as `paused`. A `failed` replication shows as a " +
	"change, and the next apply sends the request again."

// ReplicationActivationStatus is the Terraform state of the replication activation status resource.
type ReplicationActivationStatus struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`
	ReplicationId  types.String `tfsdk:"replication_id"`
	State          types.String `tfsdk:"state"`
}

// NewReplicationActivationStatus builds the state with all IDs set, so an import does not keep the
// import string in replication_id.
func NewReplicationActivationStatus(organizationId, projectId, clusterId, replicationId, state string) *ReplicationActivationStatus {
	return &ReplicationActivationStatus{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		ReplicationId:  types.StringValue(replicationId),
		State:          types.StringValue(state),
	}
}

// Validate checks the IDs and returns them. On import, it reads them from the import string in replication_id.
func (r ReplicationActivationStatus) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: r.OrganizationId,
		ProjectId:      r.ProjectId,
		ClusterId:      r.ClusterId,
		ReplicationId:  r.ReplicationId,
	}

	IDs, err := validateSchemaState(state, ReplicationId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}
	return IDs, nil
}
