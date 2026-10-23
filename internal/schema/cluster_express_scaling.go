package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// The express scaling endpoint is not in the published API reference yet, so these
// descriptions are used until the spec lookup finds one.
const (
	ClusterExpressScalingResourceDescription = "Manages express scaling on an existing Couchbase Capella cluster. Express scaling requires the Enterprise support plan. " +
		"Enabling or disabling it starts a rebalance that can run for many hours on large clusters, so create and update wait for the change to finish, " +
		"up to the configured timeouts. If a timeout is reached, the change carries on in Capella and the resource is saved with a warning. " +
		"Destroying this resource only removes it from the Terraform state and leaves express scaling on the cluster unchanged."
	ClusterExpressScalingDescription = "\n - Whether express scaling is enabled on the cluster. Set to `true` to enable it and `false` to disable it. " +
		"Reads as `true` while express scaling is being enabled, and `false` while it is being disabled."
)

// ClusterExpressScaling manages the express scaling setting of a cluster.
type ClusterExpressScaling struct {
	OrganizationId      types.String   `tfsdk:"organization_id"`
	ProjectId           types.String   `tfsdk:"project_id"`
	ClusterId           types.String   `tfsdk:"cluster_id"`
	ExpressScaling      types.Bool     `tfsdk:"express_scaling"`
	ExpressScalingState types.String   `tfsdk:"express_scaling_state"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

func (c *ClusterExpressScaling) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: c.OrganizationId,
		ProjectId:      c.ProjectId,
		ClusterId:      c.ClusterId,
	}

	IDs, err := validateSchemaState(state, ClusterId)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %w", err)
	}

	return IDs, nil
}
