package schema

import (
	"context"
	"fmt"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/freeTierCluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// FreeTierCluster is the struct for the free-tier cluster as read by the terraform provider schema for state file
type FreeTierCluster struct {
	//Id of the free-tier cluster
	Id types.String `tfsdk:"id"`
	//Availability of the free-tier cluster. It is single zone for free-tier clusters.
	Availability types.Object `tfsdk:"availability"`
	//CloudProvider of the free-tier cluster
	CloudProvider *CloudProvider `tfsdk:"cloud_provider"`
	//ProjectId of the free-tier cluster
	ProjectId types.String `tfsdk:"project_id"`
	//Audit data of the free-tier cluster
	Audit types.Object `tfsdk:"audit"`
	//Support plan used by the free-tier cluster
	Support types.Object `tfsdk:"support"`
	//OrganizationId of the free-tier cluster
	OrganizationId types.String `tfsdk:"organization_id"`
	// Name of the cluster (up to 256 characters).
	Name types.String `tfsdk:"name"`
	// CouchbaseServer is the version of the Couchbase Server to be installed in the cluster.
	CouchbaseServer types.Object `tfsdk:"couchbase_server"`
	// Description of the cluster (up to 1024 characters).
	Description types.String `tfsdk:"description"`
	// Id of the app service assosciated with the free-tier cluster
	AppServiceId types.String `tfsdk:"app_service_id"`
	// EnablePrivateDNSResolution signals that the cluster should have hostnames that are hosted in a public DNS zone that resolve to a private DNS address.
	// This exists to support the use case of customers connecting from their own data centers where it is not possible to make use of a cloud service provider DNS zone.
	EnablePrivateDNSResolution types.Bool `tfsdk:"enable_private_dns_resolution"`
	// CurrentState tells the status of the cluster - if it's healthy or degraded.
	CurrentState types.String `tfsdk:"current_state"`
	// ConnectionString specifies the Capella database endpoint for your client connection.
	ConnectionString types.String `tfsdk:"connection_string"`
	//ServiceGroups is the couchbase service groups to be run. At least one service group must contain the data service.
	ServiceGroups types.Set `tfsdk:"service_groups"`
	//cmekId is the customer managed encryption key id
	CmekId types.String `tfsdk:"cmek_id"`
}

func NewFreeTierCluster(ctx context.Context, getfreeClusterResponse *freeTierClusterapi.GetFreeTierClusterResponse, organizationId, projectId string, auditObject, availabilityObject, supportObject basetypes.ObjectValue, serviceGroupObj types.Set) (*FreeTierCluster, error) {
	newFreTierCluster := FreeTierCluster{
		Id:                         types.StringValue(getfreeClusterResponse.ID.String()),
		OrganizationId:             types.StringValue(organizationId),
		ProjectId:                  types.StringValue(projectId),
		Name:                       types.StringValue(getfreeClusterResponse.Name),
		Description:                types.StringValue(getfreeClusterResponse.Description),
		EnablePrivateDNSResolution: types.BoolValue(getfreeClusterResponse.EnablePrivateDNSResolution),
		Availability:               availabilityObject,
		CloudProvider: &CloudProvider{
			Cidr:   types.StringValue(getfreeClusterResponse.CloudProvider.Cidr),
			Region: types.StringValue(getfreeClusterResponse.CloudProvider.Region),
			Type:   types.StringValue(string(getfreeClusterResponse.CloudProvider.Type)),
		},
		Support:          supportObject,
		ConnectionString: types.StringValue(getfreeClusterResponse.ConnectionString),
		CurrentState:     types.StringValue(getfreeClusterResponse.CurrentState),
		Audit:            auditObject,
		ServiceGroups:    serviceGroupObj,
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

func NewTerraformServiceGroups(cluster *freeTierClusterapi.GetFreeTierClusterResponse) ([]ServiceGroup, error) {
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

// Validate validates the FreeTierCluster object
func (f *FreeTierCluster) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: f.OrganizationId,
		ProjectId:      f.ProjectId,
		Id:             f.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs, nil
}

// AttributeTypes returns a mapping of field names to their respective attribute types for the Availability struct.
func (a Availability) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type": types.StringType,
	}
}

// NewAvailability returns a new Availability object from the given API Availability object.
func NewAvailability(apiAvailability freeTierClusterapi.Availability) Availability {
	return Availability{
		Type: types.StringValue(apiAvailability.Type),
	}
}

// ServiceGroupAttributeTypes returns a mapping of field names to their respective attribute types for the ServiceGroup struct.
func ServiceGroupAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"node":         types.ObjectType{AttrTypes: Node{}.AttributeTypes()},
		"services":     types.SetType{ElemType: types.StringType},
		"num_of_nodes": types.Int64Type,
	}
}

// AttributeTypes returns a mapping of field names to their respective attribute types for the Node struct.
func (n Node) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"disk":    types.ObjectType{AttrTypes: Node_Disk{}.AttributeTypes()},
		"compute": types.ObjectType{AttrTypes: Compute{}.AttributeTypes()},
	}
}

// AttributeTypes returns a mapping of field names to their respective attribute types for the Node_Disk struct.
func (d Node_Disk) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":          types.StringType,
		"storage":       types.Int64Type,
		"iops":          types.Int64Type,
		"autoexpansion": types.BoolType,
	}
}

// AttributeTypes returns a mapping of field names to their respective attribute types for the Compute struct.
func (c Compute) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cpu": types.Int64Type,
		"ram": types.Int64Type,
	}
}

// AttributeTypes returns a mapping of field names to their respective attribute types for the ServiceGroup struct.
func (support Support) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"plan":     types.StringType,
		"timezone": types.StringType,
	}
}

// NewSupport returns a new Support object from the given API Support object.
func NewSupport(apiSupport freeTierClusterapi.Support) Support {
	return Support{
		Plan:     types.StringValue(apiSupport.Plan),
		Timezone: types.StringValue(apiSupport.Timezone),
	}
}

// NewServiceGroup returns a new ServiceGroup object from the given API ServiceGroup object.
func NewServiceGroups(ctx context.Context, serviceGroups []ServiceGroup) ([]types.Object, error, diag.Diagnostics) {
	serviceGroupObjList := make([]types.Object, 0)
	for _, serviceGroup := range serviceGroups {
		serviceGroupObj, diags := types.ObjectValueFrom(ctx, ServiceGroupAttributeTypes(), serviceGroup)
		if diags.HasError() {
			return nil, fmt.Errorf("service group object error"), diags
		}

		serviceGroupObjList = append(serviceGroupObjList, serviceGroupObj)

	}

	return serviceGroupObjList, nil, nil
}
