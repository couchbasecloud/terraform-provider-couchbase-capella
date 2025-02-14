package schema

import (
	"context"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/freeTierCluster"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

//type FreeTierCluster struct {
//	Name          types.String   `tfsdk:"name"`
//	Description   types.String   `tfsdk:"description"`
//	CloudProvider *CloudProvider `tfsdk:"cloud_provider"`
//}

//type CloudProvider struct {
//	Type   types.String `tfsdk:"type"`
//	Region types.String `tfsdk:"region"`
//	Cidr   types.String `tfsdk:"cidr"`
//}

type FreeTierCluster struct {
	Id            types.String   `tfsdk:"id"`
	Availability  *Availability  `tfsdk:"availability"`
	CloudProvider *CloudProvider `tfsdk:"cloud_provider"`
	ProjectId     types.String   `tfsdk:"project_id"`
	Audit         types.Object   `tfsdk:"audit"`
	Support       *Support       `tfsdk:"support"`

	OrganizationId types.String `tfsdk:"organization_id"`
	// ConfigurationType represents whether a cluster is configured as a single-node or multi-node cluster.
	ConfigurationType types.String `tfsdk:"configuration_type"`
	// Name of the cluster (up to 256 characters).
	Name types.String `tfsdk:"name"`
	// CouchbaseServer is the version of the Couchbase Server to be installed in the cluster.
	CouchbaseServer types.Object `tfsdk:"couchbase_server"`

	// Description of the cluster (up to 1024 characters).
	Description types.String `tfsdk:"description"`

	// EnablePrivateDNSResolution signals that the cluster should have hostnames that are hosted in a public DNS zone that resolve to a private DNS address.
	// This exists to support the use case of customers connecting from their own data centers where it is not possible to make use of a cloud service provider DNS zone.
	EnablePrivateDNSResolution types.Bool `tfsdk:"enable_private_dns_resolution"`
	// CurrentState tells the status of the cluster - if it's healthy or degraded.
	CurrentState types.String `tfsdk:"current_state"`
	// ConnectionString specifies the Capella database endpoint for your client connection.
	ConnectionString types.String `tfsdk:"connection_string"`
	// ServiceGroups is the couchbase service groups to be run. At least one service group must contain the data service.
	ServiceGroups []ServiceGroup `tfsdk:"service_groups"`
}

func NewFreeTierCluster(ctx context.Context, getfreeClusterResponse *freeTierClusterapi.GetFreeTierClusterResponse, organizationId, projectId string, auditObject basetypes.ObjectValue) (*FreeTierCluster, error) {
	newFreTierCluster := FreeTierCluster{
		Id:                         types.StringValue(getfreeClusterResponse.ID.String()),
		OrganizationId:             types.StringValue(organizationId),
		ProjectId:                  types.StringValue(projectId),
		Name:                       types.StringValue(getfreeClusterResponse.Name),
		Description:                types.StringValue(getfreeClusterResponse.Description),
		EnablePrivateDNSResolution: types.BoolValue(getfreeClusterResponse.EnablePrivateDNSResolution),
		Availability: &Availability{
			Type: types.StringValue(getfreeClusterResponse.Availability.Type),
		},
		CloudProvider: &CloudProvider{
			Cidr:   types.StringValue(getfreeClusterResponse.CloudProvider.Cidr),
			Region: types.StringValue(getfreeClusterResponse.CloudProvider.Region),
			Type:   types.StringValue(string(getfreeClusterResponse.CloudProvider.Type)),
		},
		Support: &Support{
			Plan:     types.StringValue(getfreeClusterResponse.Support.Plan),
			Timezone: types.StringValue(getfreeClusterResponse.Support.Timezone),
		},
		ConnectionString: types.StringValue(getfreeClusterResponse.ConnectionString),
		CurrentState:     types.StringValue(getfreeClusterResponse.CurrentState),
		Audit:            auditObject,
		//Etag:             types.StringValue(getClusterResponse.Etag),
	}
	if getfreeClusterResponse.CouchbaseServer.Version != nil {
		version := *getfreeClusterResponse.CouchbaseServer.Version
		version = removePatch(version)
		couchbaseServer := CouchbaseServer{
			Version: types.StringValue(version),
		}
		couchbaseServerObject, diags := types.ObjectValueFrom(ctx, couchbaseServer.AttributeTypes(), couchbaseServer)
		if diags.HasError() {
			return nil, fmt.Errorf("error while converting couchbase server version")
		}
		newFreTierCluster.CouchbaseServer = couchbaseServerObject
	}
	return &newFreTierCluster, nil

}
