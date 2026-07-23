package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/data_api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// DataApi maps the couchbase-capella_data_api resource and data source schema data.
type DataApi struct {
	OrganizationId         types.String `tfsdk:"organization_id"`
	ProjectId              types.String `tfsdk:"project_id"`
	ClusterId              types.String `tfsdk:"cluster_id"`
	EnableDataApi          types.Bool   `tfsdk:"enable_data_api"`
	EnableNetworkPeering   types.Bool   `tfsdk:"enable_network_peering"`
	StateForDataApi        types.String `tfsdk:"state_for_data_api"`
	StateForNetworkPeering types.String `tfsdk:"state_for_network_peering"`
	ConnectionString       types.String `tfsdk:"connection_string"`
}

// NewDataApi morphs the Data API status response into the schema struct.
func NewDataApi(organizationId, projectId, clusterId string, status *data_api.GetDataApiStatusResponse) *DataApi {
	return &DataApi{
		OrganizationId:         types.StringValue(organizationId),
		ProjectId:              types.StringValue(projectId),
		ClusterId:              types.StringValue(clusterId),
		EnableDataApi:          types.BoolValue(status.Enabled),
		EnableNetworkPeering:   types.BoolValue(status.EnabledForNetworkPeering),
		StateForDataApi:        types.StringValue(string(status.State)),
		StateForNetworkPeering: types.StringValue(string(status.StateForNetworkPeering)),
		ConnectionString:       types.StringValue(status.ConnectionString),
	}
}

// Validate is used to verify that IDs have been properly imported.
func (d *DataApi) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: d.OrganizationId,
		ProjectId:      d.ProjectId,
		ClusterId:      d.ClusterId,
	}

	IDs, err := validateSchemaState(state, ClusterId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
