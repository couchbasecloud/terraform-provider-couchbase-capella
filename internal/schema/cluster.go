package schema

import (
	"fmt"

	clusterapi "terraform-provider-capella/internal/api/cluster"
	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Availability defines the type of Availability Zone configuration for a cluster resource.
// single type means the nodes in the cluster will all be deployed in a single availability
// zone in the cloud region. multi type means the nodes in the cluster will all be deployed
// in separate multiple availability zones in the cloud region.
type Availability struct {
	// Type is the availability zone type, either 'single' or 'multi'.
	Type types.String `tfsdk:"type"`
}

// CloudProvider is the cloud provider where the cluster will be hosted.
// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
type CloudProvider struct {
	// Cidr is the cidr block for Cloud Provider.
	Cidr types.String `tfsdk:"cidr"`

	// Region is the cloud provider region, e.g. 'us-west-2'.
	// For information about supported regions, see
	// [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	// [Google Cloud Provider](https://docs.couchbase.com/cloud/reference/gcp.html).
	// [Azure Cloud](https://docs.couchbase.com/cloud/reference/azure.html).
	Region types.String `tfsdk:"region"`

	// Type is the cloud provider type, either 'AWS', 'GCP', or 'Azure'.
	Type types.String `tfsdk:"type"`
}

// Compute depicts the couchbase compute, following are the supported compute combinations
// for CPU and RAM for different cloud providers. To learn more,
// To learn more, see:
// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
type Compute struct {
	// Cpu depicts cpu units (cores).
	Cpu types.Int64 `tfsdk:"cpu"`

	// Ram depicts ram units (GB).
	Ram types.Int64 `tfsdk:"ram"`
}

// CouchbaseServer defines version for the Couchbase Server to be launched during the creation of the Capella cluster.
type CouchbaseServer struct {
	// Version is the version of the Couchbase Server to be installed in the cluster.
	// Refer to documentation [here](https://docs.couchbase.com/cloud/clusters/upgrade-database.html#server-version-maintenance-support)
	// for list of supported versions.
	// The latest Couchbase Server version will be deployed by default.
	Version types.String `tfsdk:"version"`
}

// Service is the couchbase service to run on the node.
type Service string

// ServiceGroup is the set of nodes that share the same disk, number of nodes and services.
type ServiceGroup struct {
	Node *Node `tfsdk:"node"`

	// NumOfNodes is number of nodes. The minimum number of nodes for the cluster
	// can be 3 and maximum can be 27 nodes. Additional service groups can have
	// 2 nodes minimum and 24 nodes maximum.
	NumOfNodes types.Int64 `tfsdk:"num_of_nodes"`

	// Services is the couchbase service to run on the node.
	Services []types.String `tfsdk:"services"`
}

// Node defines attributes of a cluster node.
type Node struct {
	// Compute Following are the supported compute combinations for CPU and RAM
	// for different cloud providers. To learn more, see
	// [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	Compute Compute `tfsdk:"compute"`
	// Disk is the type of disk that is supported per cloud provider during cluster creation.
	Disk Node_Disk `tfsdk:"disk"`
}

// Node_Disk is the type of disk on a particular node that is supported per cloud provider during cluster creation.
type Node_Disk struct {
	Type    types.String `tfsdk:"type"`
	Storage types.Int64  `tfsdk:"storage"`
	IOPS    types.Int64  `tfsdk:"iops"`
}

// Support defines the support plan and timezone for this particular cluster.
type Support struct {
	// Plan is the plan type, either 'Basic', 'Developer Pro', or 'Enterprise'.
	Plan types.String `tfsdk:"plan"`

	// Timezone is the standard timezone for the cluster.
	// Should be the TZ identifier.
	Timezone types.String `tfsdk:"timezone"`
}

// Cluster defines the response as received from V4 Capella Public API when asked to create a new cluster.
type Cluster struct {
	Id types.String `tfsdk:"id"`

	// AppServiceId is the ID of the linked app service.
	AppServiceId   types.String  `tfsdk:"app_service_id"`
	Audit          types.Object  `tfsdk:"audit"`
	OrganizationId types.String  `tfsdk:"organization_id"`
	ProjectId      types.String  `tfsdk:"project_id"`
	Availability   *Availability `tfsdk:"availability"`

	// CloudProvider The cloud provider where the cluster will be hosted.
	// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	CloudProvider *CloudProvider `tfsdk:"cloud_provider"`

	ConfigurationType types.String     `tfsdk:"configuration_type"`
	CouchbaseServer   *CouchbaseServer `tfsdk:"couchbase_server"`

	// Description of the cluster (up to 1024 characters).
	Description types.String `tfsdk:"description"`

	// Name of the cluster (up to 256 characters).
	Name types.String `tfsdk:"name"`

	// ServiceGroups is the couchbase service groups to be run. At least one service group must contain the data service.
	ServiceGroups []ServiceGroup `tfsdk:"service_groups"`
	Support       *Support       `tfsdk:"support"`
	CurrentState  types.String   `tfsdk:"current_state"`
	Etag          types.String   `tfsdk:"etag"`

	IfMatch types.String `tfsdk:"if_match"`
}

// NewCluster create new cluster object
func NewCluster(cluster *clusterapi.GetClusterResponse, organizationId, projectId string, auditObject basetypes.ObjectValue) (*Cluster, error) {
	newCluster := Cluster{
		Id:             types.StringValue(cluster.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		Name:           types.StringValue(cluster.Name),
		Description:    types.StringValue(cluster.Description),
		Availability: &Availability{
			Type: types.StringValue(string(cluster.Availability.Type)),
		},
		CloudProvider: &CloudProvider{
			Cidr:   types.StringValue(cluster.CloudProvider.Cidr),
			Region: types.StringValue(cluster.CloudProvider.Region),
			Type:   types.StringValue(string(cluster.CloudProvider.Type)),
		},
		ConfigurationType: types.StringValue(string(cluster.ConfigurationType)),
		Support: &Support{
			Plan:     types.StringValue(string(cluster.Support.Plan)),
			Timezone: types.StringValue(string(cluster.Support.Timezone)),
		},
		CurrentState: types.StringValue(string(cluster.CurrentState)),
		Audit:        auditObject,
		Etag:         types.StringValue(cluster.Etag),
	}

	if cluster.CouchbaseServer.Version != nil {
		version := *cluster.CouchbaseServer.Version
		newCluster.CouchbaseServer = &CouchbaseServer{
			Version: types.StringValue(version),
		}
	}

	newServiceGroups, err := morphToTerraformServiceGroups(cluster)
	if err != nil {
		return nil, err
	}
	newCluster.ServiceGroups = newServiceGroups
	return &newCluster, nil
}

func morphToTerraformServiceGroups(cluster *clusterapi.GetClusterResponse) ([]ServiceGroup, error) {
	var newServiceGroups []ServiceGroup
	for _, serviceGroup := range cluster.ServiceGroups {
		newServiceGroup := ServiceGroup{
			Node: &Node{
				Compute: Compute{
					Ram: types.Int64Value(int64(serviceGroup.Node.Compute.Ram)),
					Cpu: types.Int64Value(int64(serviceGroup.Node.Compute.Cpu)),
				},
			},
			NumOfNodes: types.Int64Value(int64(*serviceGroup.NumOfNodes)),
		}

		switch cluster.CloudProvider.Type {
		case clusterapi.Aws:
			awsDisk, err := serviceGroup.Node.AsDiskAWS()
			if err != nil {
				return nil, err
			}
			newServiceGroup.Node.Disk = Node_Disk{
				Type:    types.StringValue(string(awsDisk.Type)),
				Storage: types.Int64Value(int64(awsDisk.Storage)),
				IOPS:    types.Int64Value(int64(awsDisk.Iops)),
			}
		case clusterapi.Azure:
			azureDisk, err := serviceGroup.Node.AsDiskAzure()
			if err != nil {
				return nil, err
			}

			newServiceGroup.Node.Disk = Node_Disk{
				Type:    types.StringValue(string(azureDisk.Type)),
				Storage: types.Int64Value(int64(*azureDisk.Storage)),
				IOPS:    types.Int64Value(int64(*azureDisk.Iops)),
			}
		case clusterapi.Gcp:
			gcpDisk, err := serviceGroup.Node.AsDiskGCP()
			if err != nil {
				return nil, err
			}
			newServiceGroup.Node.Disk = Node_Disk{
				Type:    types.StringValue(string(gcpDisk.Type)),
				Storage: types.Int64Value(int64(gcpDisk.Storage)),
			}
		default:
			return nil, errors.ErrUnsupportedCloudProvider
		}

		if serviceGroup.NumOfNodes != nil {
			newServiceGroup.NumOfNodes = types.Int64Value(int64(*serviceGroup.NumOfNodes))
		}

		if serviceGroup.Services != nil {
			for _, service := range *serviceGroup.Services {
				tfService := types.StringValue(string(service))
				newServiceGroup.Services = append(newServiceGroup.Services, tfService)
			}
		}
		newServiceGroups = append(newServiceGroups, newServiceGroup)
	}
	return newServiceGroups, nil
}

func (c *Cluster) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: c.OrganizationId,
		ProjectId:      c.ProjectId,
		Id:             c.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs, nil
}

// Clusters defines structure based on the response received from V4 Capella Public API when asked to list clusters.
type Clusters struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the cluster
	ProjectId types.String `tfsdk:"project_id"`

	// Data It contains the list of resources.
	Data []ClusterData `tfsdk:"data"`
}

// ClusterData defines attributes for a single cluster when fetched from the V4 Capella Public API.
type ClusterData struct {
	Id types.String `tfsdk:"id"`

	// AppServiceId is the ID of the linked app service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Audit contains all audit-related fields.
	Audit types.Object `tfsdk:"audit"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// Availability defines if the cluster nodes will be deployed in multiple or single availability zones in the cloud.
	Availability *Availability `tfsdk:"availability"`

	// CloudProvider The cloud provider where the cluster will be hosted.
	// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	CloudProvider *CloudProvider `tfsdk:"cloud_provider"`

	// CouchbaseServer defines version for the Couchbase Server to be launched during the creation of the Capella cluster.
	CouchbaseServer *CouchbaseServer `tfsdk:"couchbase_server"`

	// Description of the cluster (up to 1024 characters).
	Description types.String `tfsdk:"description"`

	// Name of the cluster (up to 256 characters).
	Name types.String `tfsdk:"name"`

	// ServiceGroups is the couchbase service groups to be run. At least one service group must contain the data service.
	ServiceGroups []ServiceGroup `tfsdk:"service_groups"`

	// Support defines the support plan and timezone for this particular cluster.
	Support *Support `tfsdk:"support"`

	// State defines the current state of cluster
	CurrentState types.String `tfsdk:"current_state"`
}

// NewClusterData creates a new cluster data object
func NewClusterData(cluster *clusterapi.GetClusterResponse, organizationId, projectId string, auditObject basetypes.ObjectValue) (*ClusterData, error) {
	newClusterData := ClusterData{
		Id:             types.StringValue(cluster.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		Name:           types.StringValue(cluster.Name),
		Description:    types.StringValue(cluster.Description),
		Availability: &Availability{
			Type: types.StringValue(string(cluster.Availability.Type)),
		},
		CloudProvider: &CloudProvider{
			Cidr:   types.StringValue(cluster.CloudProvider.Cidr),
			Region: types.StringValue(cluster.CloudProvider.Region),
			Type:   types.StringValue(string(cluster.CloudProvider.Type)),
		},
		Support: &Support{
			Plan:     types.StringValue(string(cluster.Support.Plan)),
			Timezone: types.StringValue(string(cluster.Support.Timezone)),
		},
		CurrentState: types.StringValue(string(cluster.CurrentState)),
		Audit:        auditObject,
	}

	if cluster.CouchbaseServer.Version != nil {
		version := *cluster.CouchbaseServer.Version
		newClusterData.CouchbaseServer = &CouchbaseServer{
			Version: types.StringValue(version),
		}
	}

	newServiceGroups, err := morphToTerraformServiceGroups(cluster)
	if err != nil {
		return nil, err
	}
	newClusterData.ServiceGroups = newServiceGroups
	return &newClusterData, nil
}
