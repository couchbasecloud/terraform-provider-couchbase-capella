package schema

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
// for CPU and RAM for different cloud providers.
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

// AttributeTypes returns a mapping of field names to their respective attribute types for the CouchbaseServer struct.
// It is used during the conversion of a types.Object field to a CouchbaseServer type.
func (c CouchbaseServer) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"version": types.StringType,
	}
}

// Service is the couchbase service to run on the node.
type Service string

// ServiceGroup is the set of nodes that share the same disk, number of nodes and services.
type ServiceGroup struct {
	Node *Node `tfsdk:"node"`

	// Services is the couchbase service to run on the node.
	Services []types.String `tfsdk:"services"`

	// NumOfNodes is number of nodes. The minimum number of nodes for the cluster
	// can be 3 and maximum can be 27 nodes. Additional service groups can have
	// 2 nodes minimum and 24 nodes maximum.
	NumOfNodes types.Int64 `tfsdk:"num_of_nodes"`
}

// Node defines attributes of a cluster node.
type Node struct {
	// Disk is the type of disk that is supported per cloud provider during cluster creation.
	Disk Node_Disk `tfsdk:"disk"`
	// Compute Following are the supported compute combinations for CPU and RAM
	// for different cloud providers. To learn more, see
	// [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	Compute Compute `tfsdk:"compute"`
}

// Node_Disk is the type of disk on a particular node that is supported per cloud provider during cluster creation.
type Node_Disk struct {
	Type          types.String `tfsdk:"type"`
	Storage       types.Int64  `tfsdk:"storage"`
	IOPS          types.Int64  `tfsdk:"iops"`
	Autoexpansion types.Bool   `tfsdk:"autoexpansion"`
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
	// Availability zone type, either 'single' or 'multi'.
	Availability *Availability `tfsdk:"availability"`

	// Support defines the support plan and timezone for this particular cluster.
	Support *Support `tfsdk:"support"`

	// CloudProvider The cloud provider where the cluster will be hosted.
	// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	CloudProvider *CloudProvider `tfsdk:"cloud_provider"`

	ProjectId      types.String `tfsdk:"project_id"`
	Id             types.String `tfsdk:"id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	Audit          types.Object `tfsdk:"audit"`

	// ConfigurationType represents whether a cluster is configured as a single-node or multi-node cluster.
	ConfigurationType types.String `tfsdk:"configuration_type"`

	// CouchbaseServer is the version of the Couchbase Server to be installed in the cluster.
	CouchbaseServer types.Object `tfsdk:"couchbase_server"`

	// Description of the cluster (up to 1024 characters).
	Description types.String `tfsdk:"description"`

	// EnablePrivateDNSResolution signals that the cluster should have hostnames that are hosted in a public DNS zone that resolve to a private DNS address.
	// This exists to support the use case of customers connecting from their own data centers where it is not possible to make use of a cloud service provider DNS zone.
	EnablePrivateDNSResolution types.Bool `tfsdk:"enable_private_dns_resolution"`

	// Zones is the cloud services provider availability zones for the cluster. Currently Supported only for single AZ clusters so only 1 zone is allowed in list.
	Zones types.Set `tfsdk:"zones"`

	// Name of the cluster (up to 256 characters).
	Name types.String `tfsdk:"name"`

	// AppServiceId is the ID of the linked app service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// ConnectionString specifies the Capella database endpoint for your client connection.
	ConnectionString types.String `tfsdk:"connection_string"`

	// CurrentState tells the status of the cluster - if it's healthy or degraded.
	CurrentState types.String `tfsdk:"current_state"`

	// Etag represents the version of the document
	Etag    types.String `tfsdk:"etag"`
	IfMatch types.String `tfsdk:"if_match"`

	// ServiceGroups is the couchbase service groups to be run. At least one service group must contain the data service.
	ServiceGroups []ServiceGroup `tfsdk:"service_groups"`
}

// removePatch removes the patch version from the provided cluster server version.
func removePatch(version string) string {
	// Split the version string by '.'
	parts := strings.Split(version, ".")

	if len(parts) >= 2 {
		// Remove the last part (patch) if it's a digit
		if _, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
			parts = parts[:len(parts)-1]
		}

		// Join the parts back together
		result := strings.Join(parts, ".")
		return result
	}

	// If the version is in an invalid format (e.g., '7'), return the same version
	return version
}

// NewCluster create new cluster object.
func NewCluster(ctx context.Context, cluster *clusterapi.GetClusterResponse, organizationId, projectId string, auditObject basetypes.ObjectValue) (*Cluster, error) {
	newCluster := Cluster{
		Id:                         types.StringValue(cluster.Id.String()),
		OrganizationId:             types.StringValue(organizationId),
		ProjectId:                  types.StringValue(projectId),
		Name:                       types.StringValue(cluster.Name),
		Description:                types.StringValue(cluster.Description),
		EnablePrivateDNSResolution: types.BoolValue(cluster.EnablePrivateDNSResolution),
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
		ConnectionString: types.StringValue(cluster.ConnectionString),
		CurrentState:     types.StringValue(string(cluster.CurrentState)),
		Audit:            auditObject,
		Etag:             types.StringValue(cluster.Etag),
	}

	if cluster.CouchbaseServer.Version != nil {
		version := *cluster.CouchbaseServer.Version
		version = removePatch(version)
		couchbaseServer := CouchbaseServer{
			Version: types.StringValue(version),
		}
		couchbaseServerObject, diags := types.ObjectValueFrom(ctx, couchbaseServer.AttributeTypes(), couchbaseServer)
		if diags.HasError() {
			return nil, fmt.Errorf("error while converting couchbase server version")
		}
		newCluster.CouchbaseServer = couchbaseServerObject
	}

	newServiceGroups, err := morphToTerraformServiceGroups(cluster)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrConvertingServiceGroups, err)
	}
	newCluster.ServiceGroups = newServiceGroups

	newZones, err := MorphZones(cluster.Zones)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrConvertingZone, err)
	}

	newCluster.Zones = newZones

	return &newCluster, nil
}

func MorphZones(zones []string) (basetypes.SetValue, error) {
	var newZone []attr.Value
	for _, zone := range zones {
		newZone = append(newZone, types.StringValue(zone))
	}

	newZones, diags := types.SetValue(types.StringType, newZone)
	if diags.HasError() {
		return types.SetUnknown(types.StringType), fmt.Errorf("error while converting zones")
	}

	return newZones, nil
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
				return nil, fmt.Errorf("%s: %w", errors.ErrReadingAWSDisk, err)
			}
			newServiceGroup.Node.Disk = Node_Disk{
				Type:    types.StringValue(string(awsDisk.Type)),
				Storage: types.Int64Value(int64(awsDisk.Storage)),
				IOPS:    types.Int64Value(int64(awsDisk.Iops)),
			}
		case clusterapi.Azure:
			azureDisk, err := serviceGroup.Node.AsDiskAzure()
			if err != nil {
				return nil, fmt.Errorf("%s: %w", errors.ErrReadingAzureDisk, err)
			}

			newServiceGroup.Node.Disk = Node_Disk{
				Type:          types.StringValue(string(azureDisk.Type)),
				Storage:       types.Int64Value(int64(*azureDisk.Storage)),
				IOPS:          types.Int64Value(int64(*azureDisk.Iops)),
				Autoexpansion: types.BoolValue(*azureDisk.Autoexpansion),
			}
		case clusterapi.Gcp:
			gcpDisk, err := serviceGroup.Node.AsDiskGCP()
			if err != nil {
				return nil, fmt.Errorf("%s: %w", errors.ErrReadingGCPDisk, err)
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
	Availability               *Availability    `tfsdk:"availability"`
	Support                    *Support         `tfsdk:"support"`
	CouchbaseServer            *CouchbaseServer `tfsdk:"couchbase_server"`
	CloudProvider              *CloudProvider   `tfsdk:"cloud_provider"`
	OrganizationId             types.String     `tfsdk:"organization_id"`
	ProjectId                  types.String     `tfsdk:"project_id"`
	Id                         types.String     `tfsdk:"id"`
	Audit                      types.Object     `tfsdk:"audit"`
	Description                types.String     `tfsdk:"description"`
	EnablePrivateDNSResolution types.Bool       `tfsdk:"enable_private_dns_resolution"`
	Zones                      types.List       `tfsdk:"zones"`
	Name                       types.String     `tfsdk:"name"`
	AppServiceId               types.String     `tfsdk:"app_service_id"`
	ConnectionString           types.String     `tfsdk:"connection_string"`
	CurrentState               types.String     `tfsdk:"current_state"`
	ServiceGroups              []ServiceGroup   `tfsdk:"service_groups"`
}

// NewClusterData creates a new cluster data object.
func NewClusterData(cluster *clusterapi.GetClusterResponse, organizationId, projectId string, auditObject basetypes.ObjectValue) (*ClusterData, error) {
	newClusterData := ClusterData{
		Id:                         types.StringValue(cluster.Id.String()),
		OrganizationId:             types.StringValue(organizationId),
		ProjectId:                  types.StringValue(projectId),
		Name:                       types.StringValue(cluster.Name),
		Description:                types.StringValue(cluster.Description),
		EnablePrivateDNSResolution: types.BoolValue(cluster.EnablePrivateDNSResolution),
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
		ConnectionString: types.StringValue(cluster.ConnectionString),
		CurrentState:     types.StringValue(string(cluster.CurrentState)),
		Audit:            auditObject,
	}

	if cluster.CouchbaseServer.Version != nil {
		version := *cluster.CouchbaseServer.Version
		newClusterData.CouchbaseServer = &CouchbaseServer{
			Version: types.StringValue(version),
		}
	}

	newServiceGroups, err := morphToTerraformServiceGroups(cluster)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrConvertingServiceGroups, err)
	}
	newClusterData.ServiceGroups = newServiceGroups

	var newZone []attr.Value
	for _, zone := range cluster.Zones {
		newZone = append(newZone, types.StringValue(zone))
	}

	zones, diags := types.ListValue(types.StringType, newZone)
	if diags.HasError() {
		return &ClusterData{}, fmt.Errorf("error while converting zones")
	}

	newClusterData.Zones = zones

	return &newClusterData, nil
}
