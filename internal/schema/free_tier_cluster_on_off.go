package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ClusterOnOffOnDemand provides the means to turn the given cluster to on or off state.
type FreeTierClusterOnOff struct {
	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// State is the state (on/off) the cluster needs to be turned on/off of the cluster
	State types.String `tfsdk:"state"`
}

func (f *FreeTierClusterOnOff) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: f.OrganizationId,
		ProjectId:      f.ProjectId,
		ClusterId:      f.ClusterId,
	}

	IDs, err := validateSchemaState(state, ClusterId)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs, nil
}
