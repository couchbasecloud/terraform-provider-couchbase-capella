package schema

import (
	"strings"
	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Availability defines the type of Availability Zone configuration for a cluster resource.
// single type means the nodes in the cluster will all be deployed in a single availability zone in the cloud region.
// multi type means the nodes in the cluster will all be deployed in separate multiple availability zones in the cloud region.
type Availability struct {
	// Type Availability zone type, either 'single' or 'multi'.
	Type types.String `tfsdk:"type"`
}

// CloudProvider The cloud provider where the cluster will be hosted. List of providers and the hosted regions -
// | Provider | Regions |
// | -------- | ------- |
// | AWS      | *Americas* - us-east-1, us-east-2, us-west-2, ca-central-1, sa-east-1
// |          | *Europe / Middle East* - eu-central-1, eu-west-1, eu-west-2, eu-west-3, eu-north-1,
// |          | *AsiaPacific* - ap-southeast-1, ap-southeast-2, ap-northeast-1, ap-northeast-2, ap-south-1 |
// | GCP      | *Americas* - us-east1, us-east4, us-west1, us-west3, us-west4, us-central1, northamerica-northeast1, northamerica-northeast2, southamerica-east1, southamerica-west1,
// |          | *Europe* - europe-west1, europe-west2, europe-west3, europe-west4, europe-west6, europe-west8, europe-central2, europe-north1,
// |          | *Asia Pacific* - asia-east1, asia-east2, asia-northeast1, asia-northeast2, asia-northeast3, asia-south1, asia-south2, asia-southeast1, asia-southeast2, australia-southeast1, australia-southeast2 |
// | Azure    | *Americas* - eastus, canadacentral, westus3, brazilsouth,
// |          | *Europe* - norwayeast, uksouth, westeurope, swedencentral,
// |          | *Asia Pacific* - australiaeast, koreacentral, centralindia |
//
// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
type CloudProvider struct {
	// Cidr CIDR block for Cloud Provider.
	Cidr types.String `tfsdk:"cidr"`

	// Region Cloud provider region, e.g. 'us-west-2'. For information about supported regions,
	// see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	Region types.String `tfsdk:"region"`

	// Type Cloud provider type, either 'AWS', 'GCP', or 'Azure'.
	Type types.String `tfsdk:"type"`
}

// Compute Following are the supported compute combinations for CPU and RAM for different cloud providers.
// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
type Compute struct {
	// Cpu CPU units (cores).
	Cpu types.Int64 `tfsdk:"cpu"`

	// Ram RAM units (GB).
	Ram types.Int64 `tfsdk:"ram"`
}

// CouchbaseServer defines model for CouchbaseServer.
type CouchbaseServer struct {
	// Version Version of the Couchbase Server to be installed in the cluster.
	// Refer to documentation [here](https://docs.couchbase.com/cloud/clusters/upgrade-database.html#server-version-maintenance-support)
	// for list of supported versions.
	// The latest Couchbase Server version will be deployed by default.
	Version types.String `tfsdk:"version"`
}

type Service string

// ServiceGroup The set of nodes that share the same disk, number of nodes and services.
type ServiceGroup struct {
	Node *Node `tfsdk:"node"`

	// NumOfNodes Number of nodes. The minimum number of nodes for the cluster can be 3 and maximum can be 27 nodes.
	// Additional service groups can have 2 nodes minimum and 24 nodes maximum.
	NumOfNodes types.Int64 `tfsdk:"num_of_nodes"`

	// Services The couchbase service to run on the node.
	Services []types.String `tfsdk:"services"`
}

type Node struct {
	// Compute Following are the supported compute combinations for CPU and RAM for different cloud providers.
	// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	Compute Compute   `tfsdk:"compute"`
	Disk    Node_Disk `tfsdk:"disk"`
}

// Node_Disk defines model for Node.Disk.
type Node_Disk struct {
	Type    types.String `tfsdk:"type"`
	Storage types.Int64  `tfsdk:"storage"`
	IOPS    types.Int64  `tfsdk:"iops"`
}

// Support defines model for Support.
type Support struct {
	// Plan Plan type, either 'Basic', 'Developer Pro', or 'Enterprise'.
	Plan types.String `tfsdk:"plan"`

	// Timezone The standard timezone for the cluster. Should be the TZ identifier.
	Timezone types.String `tfsdk:"timezone"`
}

// ClusterResourceModel defines model for CreateClusterRequest.
type ClusterResourceModel struct {
	Id types.String `tfsdk:"id"`

	// AppServiceId The ID of the linked app service.
	AppServiceId   types.String  `tfsdk:"app_service_id"`
	Audit          types.Object  `tfsdk:"audit"`
	OrganizationId types.String  `tfsdk:"organization_id""`
	ProjectId      types.String  `tfsdk:"project_id"`
	Availability   *Availability `tfsdk:"availability"`

	// CloudProvider The cloud provider where the cluster will be hosted. List of providers and the hosted regions -
	// | Provider | Regions |
	// | -------- | ------- |
	// | AWS      | *Americas* - us-east-1, us-east-2, us-west-2, ca-central-1, sa-east-1,
	// |		  |	*Europe / Middle East* - eu-central-1, eu-west-1, eu-west-2, eu-west-3, eu-north-1,
	// |          | *AsiaPacific* - ap-southeast-1, ap-southeast-2, ap-northeast-1, ap-northeast-2, ap-south-1 |
	// | GCP      | *Americas* - us-east1, us-east4, us-west1, us-west3, us-west4, us-central1, northamerica-northeast1, northamerica-northeast2, southamerica-east1, southamerica-west1,
	// |          | *Europe* - europe-west1, europe-west2, europe-west3, europe-west4, europe-west6, europe-west8, europe-central2, europe-north1,
	// |          | *Asia Pacific* - asia-east1, asia-east2, asia-northeast1, asia-northeast2, asia-northeast3, asia-south1, asia-south2, asia-southeast1, asia-southeast2, australia-southeast1, australia-southeast2 |
	// | Azure    | *Americas* - eastus, canadacentral, westus3, brazilsouth,
	// |          | *Europe* - norwayeast, uksouth, westeurope, swedencentral,
	// |          | *Asia Pacific* - australiaeast, koreacentral, centralindia |
	//
	// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	CloudProvider   *CloudProvider   `tfsdk:"cloud_provider"`
	CouchbaseServer *CouchbaseServer `tfsdk:"couchbase_server"`

	// Description of the cluster (up to 1024 characters).
	Description types.String `tfsdk:"description"`

	// Name of the cluster (up to 256 characters).
	Name types.String `tfsdk:"name"`

	// ServiceGroups The couchbase service groups to be run. At least one service group must contain the data service.
	ServiceGroups []ServiceGroup `tfsdk:"service_groups"`
	Support       *Support       `tfsdk:"support"`
	CurrentState  types.String   `tfsdk:"current_state"`
	Etag          types.String   `tfsdk:"etag"`

	IfMatch types.String `tfsdk:"if_match"`
}

func (c *ClusterResourceModel) Validate() error {
	if c.Id.IsNull() {
		return errors.ErrClusterIdCannotBeEmpty
	}

	if c.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}

	if c.ProjectId.IsNull() {
		return errors.ErrProjectIdCannotBeEmpty
	}
	return nil
}

func (c *ClusterResourceModel) PopulateParamsForImport() error {
	combinedIDs := c.Id.ValueString()
	splitIDs := strings.Split(combinedIDs, "#")

	if c.OrganizationId.IsNull() && len(splitIDs) > 2 {
		c.OrganizationId = types.StringValue(splitIDs[0])
	}

	if c.ProjectId.IsNull() && len(splitIDs) > 2 {
		c.ProjectId = types.StringValue(splitIDs[1])
	}

	if len(c.Id.ValueString()) > 36 && len(splitIDs) > 2 {
		c.Id = types.StringValue(splitIDs[2])
	}

	return c.Validate()
}
