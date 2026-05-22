package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ClusterDeletionProtection manages the deletion protection state of a cluster.
type ClusterDeletionProtection struct {
	OrganizationId     types.String `tfsdk:"organization_id"`
	ProjectId          types.String `tfsdk:"project_id"`
	ClusterId          types.String `tfsdk:"cluster_id"`
	DeletionProtection types.Bool   `tfsdk:"deletion_protection"`
}

func (c *ClusterDeletionProtection) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: c.OrganizationId,
		ProjectId:      c.ProjectId,
		ClusterId:      c.ClusterId,
	}

	IDs, err := validateSchemaState(state, ClusterId)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs, nil
}
