package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// DataAPI represents the Terraform state for the Data API resource.
type DataAPI struct {
	OrganizationId         types.String `tfsdk:"organization_id"`
	ProjectId              types.String `tfsdk:"project_id"`
	ClusterId              types.String `tfsdk:"cluster_id"`
	EnableDataApi          types.Bool   `tfsdk:"enable_data_api"`
	EnableNetworkPeering   types.Bool   `tfsdk:"enable_network_peering"`
	State                  types.String `tfsdk:"state"`
	StateForNetworkPeering types.String `tfsdk:"state_for_network_peering"`
	ConnectionString       types.String `tfsdk:"connection_string"`
}

// Validate is used to verify that IDs have been properly imported.
func (d *DataAPI) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: d.OrganizationId,
		ProjectId:      d.ProjectId,
		ClusterId:      d.ClusterId,
	}

	IDs, err := validateSchemaState(state, ClusterId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
