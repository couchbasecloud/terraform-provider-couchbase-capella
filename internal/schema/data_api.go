package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// DataAPI defines the response as received from V4 Capella Public API when asked to fetch the Data API status.
type DataAPI struct {
	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Enabled indicates whether Data API is enabled or disabled on the cluster.
	Enabled types.Bool `tfsdk:"enabled"`

	// State is the current status of the Data API.
	State types.String `tfsdk:"state"`

	// EnabledForNetworkPeering indicates whether network peering was enabled or disabled for the Data API.
	EnabledForNetworkPeering types.Bool `tfsdk:"enabled_for_network_peering"`

	// StateForNetworkPeering is the current status for vpc peering for Data API.
	StateForNetworkPeering types.String `tfsdk:"state_for_network_peering"`

	// ConnectionString is the connection string for the Data API service.
	ConnectionString types.String `tfsdk:"connection_string"`
}

// Validate checks the validity of the DataAPI state and extracts associated IDs.
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

// NewDataAPI creates a new DataAPI object from the API response.
func NewDataAPI(dataAPIResp *api.GetDataAPIStatusResponse,
	organizationId, projectId, clusterId string,
) *DataAPI {
	newObj := DataAPI{
		OrganizationId:           types.StringValue(organizationId),
		ProjectId:                types.StringValue(projectId),
		ClusterId:                types.StringValue(clusterId),
		Enabled:                  types.BoolValue(dataAPIResp.Enabled),
		State:                    types.StringValue(dataAPIResp.State),
		EnabledForNetworkPeering: types.BoolValue(dataAPIResp.EnabledForNetworkPeering),
		StateForNetworkPeering:   types.StringValue(dataAPIResp.StateForNetworkPeering),
		ConnectionString:         types.StringValue(dataAPIResp.ConnectionString),
	}
	return &newObj
}
