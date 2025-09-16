package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/apigen"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// FreeTierCluster is the struct for the free-tier cluster as read by the terraform provider schema for state file.
type FreeTierCluster struct {
	//ID of the free-tier cluster.
	Id types.String `tfsdk:"id"`
	//Availability of the free-tier cluster. It is single zone for free-tier clusters.
	Availability types.Object `tfsdk:"availability"`
	//CloudProvider of the free-tier cluster.
	CloudProvider *CloudProvider `tfsdk:"cloud_provider"`
	//ProjectId of the free-tier cluster
	ProjectId types.String `tfsdk:"project_id"`
	//Audit data of the free-tier cluster.
	Audit types.Object `tfsdk:"audit"`
	//Support plan used by the free-tier cluster.
	Support types.Object `tfsdk:"support"`
	//OrganizationId of the free-tier cluster.
	OrganizationId types.String `tfsdk:"organization_id"`
	// Name of the cluster (up to 256 characters).
	Name types.String `tfsdk:"name"`
	// CouchbaseServer is the version of the Couchbase Server to be installed in the cluster.
	CouchbaseServer types.Object `tfsdk:"couchbase_server"`
	// Description of the cluster (up to 1024 characters).
	Description types.String `tfsdk:"description"`
	// ID of the app service assosciated with the free-tier cluster.
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
	//cmekId is the customer managed encryption key id.
	CmekId types.String `tfsdk:"cmek_id"`
	// Etag represents the version of the document.
	Etag types.String `tfsdk:"etag"`
}

func NewFreeTierCluster(ctx context.Context, getfreeTierClusterResponse *struct {
	AppServiceId               *openapi_types.UUID       `json:"appServiceId,omitempty"`
	Audit                      apigen.CouchbaseAuditData `json:"audit"`
	Availability               apigen.Availability       `json:"availability"`
	CloudProvider              apigen.CloudProvider      `json:"cloudProvider"`
	CmekId                     *string                   `json:"cmekId,omitempty"`
	ConfigurationType          apigen.ConfigurationType  `json:"configurationType"`
	ConnectionString           string                    `json:"connectionString"`
	CouchbaseServer            apigen.CouchbaseServer    `json:"couchbaseServer"`
	CurrentState               apigen.CurrentState       `json:"currentState"`
	Description                string                    `json:"description"`
	EnablePrivateDNSResolution *bool                     `json:"enablePrivateDNSResolution,omitempty"`
	Id                         openapi_types.UUID        `json:"id"`
	Name                       string                    `json:"name"`
	ServiceGroups              []apigen.ServiceGroup     `json:"serviceGroups"`
	Support                    struct {
		Plan     apigen.GetFreeTierCluster200SupportPlan `json:"plan"`
		Timezone apigen.SupportTimezone                  `json:"timezone"`
	} `json:"support"`
}, organizationId, projectId string, auditObject, availabilityObject, supportObject basetypes.ObjectValue, serviceGroupObj types.Set) (*FreeTierCluster, error) {
	newFreeTierCluster := FreeTierCluster{
		Id:                         types.StringValue(getfreeTierClusterResponse.Id.String()),
		OrganizationId:             types.StringValue(organizationId),
		ProjectId:                  types.StringValue(projectId),
		Name:                       types.StringValue(getfreeTierClusterResponse.Name),
		Description:                types.StringValue(getfreeTierClusterResponse.Description),
		EnablePrivateDNSResolution: types.BoolValue(getBoolOrFalse(getfreeTierClusterResponse.EnablePrivateDNSResolution)),
		Availability:               availabilityObject,
		CloudProvider: &CloudProvider{
			Cidr:   types.StringValue(valueOrEmpty(getfreeTierClusterResponse.CloudProvider.Cidr)),
			Region: types.StringValue(getfreeTierClusterResponse.CloudProvider.Region),
			Type:   types.StringValue(string(getfreeTierClusterResponse.CloudProvider.Type)),
		},
		Support:          supportObject,
		ConnectionString: types.StringValue(getfreeTierClusterResponse.ConnectionString),
		CurrentState:     types.StringValue(string(getfreeTierClusterResponse.CurrentState)),
		Audit:            auditObject,
		ServiceGroups:    serviceGroupObj,
		Etag:             types.StringNull(),
	}
	if clusterVersion := getfreeTierClusterResponse.CouchbaseServer.Version; clusterVersion != nil {
		version := *clusterVersion
		version = removePatch(version)
		couchbaseServer := CouchbaseServer{
			Version: types.StringValue(version),
		}
		couchbaseServerObject, diags := types.ObjectValueFrom(ctx, couchbaseServer.AttributeTypes(), couchbaseServer)
		if diags.HasError() {
			return nil, fmt.Errorf("error while converting couchbase server version")
		}
		newFreeTierCluster.CouchbaseServer = couchbaseServerObject
	}
	return &newFreeTierCluster, nil

}

func NewTerraformServiceGroups(cluster *apigen.GetClusterResponse) ([]ServiceGroup, error) {
	var newServiceGroups []ServiceGroup
	for _, serviceGroup := range cluster.ServiceGroups {
		newServiceGroup := ServiceGroup{
			Node: &Node{
				Compute: Compute{
					Ram: types.Int64Value(int64(serviceGroup.Node.Compute.Ram)),
					Cpu: types.Int64Value(int64(serviceGroup.Node.Compute.Cpu)),
				},
			},
		}

		switch cluster.CloudProvider.Type {
		case apigen.CloudProviderType("aws"):
			awsDisk, err := serviceGroup.Node.Disk.AsDiskAWS()
			if err != nil {
				return nil, fmt.Errorf("%s: %w", errors.ErrReadingAWSDisk, err)
			}
			newServiceGroup.Node.Disk = Node_Disk{
				Type:    types.StringValue(string(awsDisk.Type)),
				Storage: types.Int64Null(),
				IOPS:    types.Int64Null(),
			}
		case apigen.CloudProviderType("azure"):
			azureDisk, err := serviceGroup.Node.Disk.AsDiskAzure()
			if err != nil {
				return nil, fmt.Errorf("%s: %w", errors.ErrReadingAzureDisk, err)
			}
			nd := Node_Disk{Type: types.StringValue(string(azureDisk.Type))}
			if azureDisk.Storage != nil {
				nd.Storage = types.Int64Value(int64(*azureDisk.Storage))
			} else {
				nd.Storage = types.Int64Null()
			}
			if azureDisk.AutoExpansion != nil {
				nd.Autoexpansion = types.BoolValue(*azureDisk.AutoExpansion)
			} else {
				nd.Autoexpansion = types.BoolNull()
			}
			newServiceGroup.Node.Disk = nd
		case apigen.CloudProviderType("gcp"):
			gcpDisk, err := serviceGroup.Node.Disk.AsDiskGCP()
			if err != nil {
				return nil, fmt.Errorf("%s: %w", errors.ErrReadingGCPDisk, err)
			}
			newServiceGroup.Node.Disk = Node_Disk{
				Type:    types.StringValue(string(gcpDisk.Type)),
				Storage: types.Int64Value(int64(gcpDisk.Storage)),
				IOPS:    types.Int64Null(),
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

// Validate validates the FreeTierCluster object.
func (f FreeTierCluster) Validate() (map[Attr]string, error) {
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
func NewAvailability(apiAvailability apigen.Availability) Availability {
	return Availability{
		Type: types.StringValue(string(apiAvailability.Type)),
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
func NewSupport(apiSupport struct {
	Plan     apigen.GetFreeTierCluster200SupportPlan
	Timezone apigen.SupportTimezone
}) Support {
	return Support{
		Plan:     types.StringValue(string(apiSupport.Plan)),
		Timezone: types.StringValue(string(apiSupport.Timezone)),
	}
}

// NewServiceGroups returns a new ServiceGroup object from the given API ServiceGroup object.
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
