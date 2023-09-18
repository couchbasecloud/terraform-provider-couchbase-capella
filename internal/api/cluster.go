package api

import (
	"encoding/json"
	"github.com/apapsch/go-jsonmerge/v2"
	"github.com/google/uuid"
)

// Defines values for CloudProviderType.
const (
	Aws   CloudProviderType = "aws"
	Azure CloudProviderType = "azure"
	Gcp   CloudProviderType = "gcp"
)

// Defines values for CurrentState.
const (
	Degraded         CurrentState = "degraded"
	Deploying        CurrentState = "deploying"
	DeploymentFailed CurrentState = "deploymentFailed"
	DestroyFailed    CurrentState = "destroyFailed"
	Destroying       CurrentState = "destroying"
	Draft            CurrentState = "draft"
	Healthy          CurrentState = "healthy"
	Offline          CurrentState = "offline"
	Peering          CurrentState = "peering"
	PeeringFailed    CurrentState = "peeringFailed"
	RebalanceFailed  CurrentState = "rebalanceFailed"
	Rebalancing      CurrentState = "rebalancing"
	ScaleFailed      CurrentState = "scaleFailed"
	Scaling          CurrentState = "scaling"
	TurnedOff        CurrentState = "turnedOff"
	TurningOff       CurrentState = "turningOff"
	TurningOffFailed CurrentState = "turningOffFailed"
	TurningOn        CurrentState = "turningOn"
	TurningOnFailed  CurrentState = "turningOnFailed"
	UpgradeFailed    CurrentState = "upgradeFailed"
	Upgrading        CurrentState = "upgrading"
)

// Availability defines model for Availability.
type Availability struct {
	// Type is availability zone type, either 'single' or 'multi'.
	Type AvailabilityType `json:"type"`
}

// AvailabilityType is availability zone type, either 'single' or 'multi'.
type AvailabilityType string

// CloudProvider depicts where the cluster will be hosted. List of providers and the hosted regions -
// | Provider | Regions |
// | -------- | ------- |
// | AWS      | *Americas* - us-east-1, us-east-2, us-west-2, ca-central-1, sa-east-1, *Europe / Middle East* - eu-central-1, eu-west-1, eu-west-2, eu-west-3, eu-north-1, *AsiaPacific* - ap-southeast-1, ap-southeast-2, ap-northeast-1, ap-northeast-2, ap-south-1 |
// | GCP      | *Americas* - us-east1, us-east4, us-west1, us-west3, us-west4, us-central1, northamerica-northeast1, northamerica-northeast2, southamerica-east1, southamerica-west1, *Europe* - europe-west1, europe-west2, europe-west3, europe-west4, europe-west6, europe-west8, europe-central2, europe-north1, *Asia Pacific* - asia-east1, asia-east2, asia-northeast1, asia-northeast2, asia-northeast3, asia-south1, asia-south2, asia-southeast1, asia-southeast2, australia-southeast1, australia-southeast2 |
// | Azure    | *Americas* - eastus, canadacentral, westus3, brazilsouth, *Europe* - norwayeast, uksouth, westeurope, swedencentral, *Asia Pacific* - australiaeast, koreacentral, centralindia |
//
// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
type CloudProvider struct {
	// Cidr block for Cloud Provider.
	Cidr string `json:"cidr"`

	// Region is cloud provider region, e.g. 'us-west-2'. For information about supported regions, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	Region string `json:"region"`

	// Type is cloud provider type, either 'AWS', 'GCP', or 'Azure'.
	Type CloudProviderType `json:"type"`
}

// CloudProviderType is cloud provider type, either 'AWS', 'GCP', or 'Azure'.
type CloudProviderType string

// Compute Following are the supported compute combinations for CPU and RAM for different cloud providers.
// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
type Compute struct {
	// Cpu depicts cpu units (cores).
	Cpu int `json:"cpu"`

	// Ram depicts ram units (GB).
	Ram int `json:"ram"`
}

// CouchbaseServer defines model for CouchbaseServer.
type CouchbaseServer struct {
	// Version is version of the Couchbase Server to be installed in the cluster. Refer to documentation [here](https://docs.couchbase.com/cloud/clusters/upgrade-database.html#server-version-maintenance-support) for list of supported versions. The latest Couchbase Server version will be deployed by default.
	Version *string `json:"version,omitempty"`
}

// CreateClusterRequest defines model for CreateClusterRequest.
type CreateClusterRequest struct {
	Availability Availability `json:"availability"`

	// CloudProvider is the cloud provider where the cluster will be hosted. List of providers and the hosted regions -
	// | Provider | Regions |
	// | -------- | ------- |
	// | AWS      | *Americas* - us-east-1, us-east-2, us-west-2, ca-central-1, sa-east-1, *Europe / Middle East* - eu-central-1, eu-west-1, eu-west-2, eu-west-3, eu-north-1, *AsiaPacific* - ap-southeast-1, ap-southeast-2, ap-northeast-1, ap-northeast-2, ap-south-1 |
	// | GCP      | *Americas* - us-east1, us-east4, us-west1, us-west3, us-west4, us-central1, northamerica-northeast1, northamerica-northeast2, southamerica-east1, southamerica-west1, *Europe* - europe-west1, europe-west2, europe-west3, europe-west4, europe-west6, europe-west8, europe-central2, europe-north1, *Asia Pacific* - asia-east1, asia-east2, asia-northeast1, asia-northeast2, asia-northeast3, asia-south1, asia-south2, asia-southeast1, asia-southeast2, australia-southeast1, australia-southeast2 |
	// | Azure    | *Americas* - eastus, canadacentral, westus3, brazilsouth, *Europe* - norwayeast, uksouth, westeurope, swedencentral, *Asia Pacific* - australiaeast, koreacentral, centralindia |
	//
	// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	CloudProvider   CloudProvider    `json:"cloudProvider"`
	CouchbaseServer *CouchbaseServer `json:"couchbaseServer,omitempty"`

	// Description depicts description of the cluster (up to 1024 characters).
	Description *string `json:"description,omitempty"`

	// Name is the name of the cluster (up to 256 characters).
	Name string `json:"name"`

	// ServiceGroups is the couchbase service groups to be run. At least one service group must contain the data service.
	ServiceGroups []ServiceGroup `json:"serviceGroups"`
	Support       Support        `json:"support"`
}

// CreateClusterResponse defines model for CreateClusterResponse.
type CreateClusterResponse struct {
	// Id The ID of the cluster created.
	Id uuid.UUID `json:"id"`
}

// CurrentState defines model for CurrentState.
type CurrentState string

// DiskAWS defines model for DiskAWS.
type DiskAWS struct {
	// Iops Please refer to documentation for supported IOPS.
	Iops int `json:"iops"`

	// Storage depicts storage in GB. See documentation for supported storage.
	Storage int `json:"storage"`

	// Type depicts type of disk. Please choose from the given list for AWS cloud provider.
	Type DiskAWSType `json:"type"`
}

// DiskAWSType depicts type of disk. Please choose from the given list for AWS cloud provider.
type DiskAWSType string

// DiskAzure defines model for DiskAzure.
type DiskAzure struct {
	// Iops is required for ultra disk types. Please refer to documentation for supported IOPS.
	Iops *int `json:"iops,omitempty"`

	// Storage depicts storage in GB. Required for ultra disk types. Please refer to documentation for supported storage.
	Storage *int `json:"storage,omitempty"`

	// Type depicts type of disk. Please choose from the given list for Azure cloud provider.
	Type DiskAzureType `json:"type"`
}

// DiskAzureType depicts type of disk. Please choose from the given list for Azure cloud provider.
type DiskAzureType string

// DiskGCP defines model for DiskGCP.
type DiskGCP struct {
	// Storage is storage in GB. Please refer to documentation for supported storage.
	Storage int `json:"storage"`

	// Type is type of disk. Please choose from the given list for GCP cloud provider.
	Type DiskGCPType `json:"type"`
}

// DiskGCPType is type of disk. Please choose from the given list for GCP cloud provider.
type DiskGCPType string

// GetClusterResponse defines model for GetClusterResponse.
type GetClusterResponse struct {
	// AppServiceId is the ID of the linked app service.
	AppServiceId *uuid.UUID         `json:"appServiceId,omitempty"`
	Audit        CouchbaseAuditData `json:"audit"`
	Availability Availability       `json:"availability"`

	// CloudProvider is the cloud provider where the cluster will be hosted. List of providers and the hosted regions -
	// | Provider | Regions |
	// | -------- | ------- |
	// | AWS      | *Americas* - us-east-1, us-east-2, us-west-2, ca-central-1, sa-east-1, *Europe / Middle East* - eu-central-1, eu-west-1, eu-west-2, eu-west-3, eu-north-1, *AsiaPacific* - ap-southeast-1, ap-southeast-2, ap-northeast-1, ap-northeast-2, ap-south-1 |
	// | GCP      | *Americas* - us-east1, us-east4, us-west1, us-west3, us-west4, us-central1, northamerica-northeast1, northamerica-northeast2, southamerica-east1, southamerica-west1, *Europe* - europe-west1, europe-west2, europe-west3, europe-west4, europe-west6, europe-west8, europe-central2, europe-north1, *Asia Pacific* - asia-east1, asia-east2, asia-northeast1, asia-northeast2, asia-northeast3, asia-south1, asia-south2, asia-southeast1, asia-southeast2, australia-southeast1, australia-southeast2 |
	// | Azure    | *Americas* - eastus, canadacentral, westus3, brazilsouth, *Europe* - norwayeast, uksouth, westeurope, swedencentral, *Asia Pacific* - australiaeast, koreacentral, centralindia |
	//
	// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	CloudProvider   CloudProvider   `json:"cloudProvider"`
	CouchbaseServer CouchbaseServer `json:"couchbaseServer"`
	CurrentState    CurrentState    `json:"currentState"`

	// Description  depicts description of the cluster (up to 1024 characters).
	Description string `json:"description"`

	// Id is the ID of the cluster created.
	Id uuid.UUID `json:"id"`

	// Name Name of the cluster (up to 256 characters).
	Name          string         `json:"name"`
	ServiceGroups []ServiceGroup `json:"serviceGroups"`
	Support       Support        `json:"support"`

	Etag string
}

// Node defines model for Node.
type Node struct {
	// Compute Following are the supported compute combinations for CPU and RAM for different cloud providers.
	// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	Compute Compute         `json:"compute"`
	Disk    json.RawMessage `json:"disk"`
}

// Service defines model for Service.
type Service string

// ServiceGroup The set of nodes that share the same disk, number of nodes and services.
type ServiceGroup struct {
	Node *Node `json:"node,omitempty"`

	// NumOfNodes is the number of nodes. The minimum number of nodes for the cluster can be 3 and maximum can be 27 nodes. Additional service groups can have 2 nodes minimum and 24 nodes maximum.
	NumOfNodes *int `json:"numOfNodes,omitempty"`

	// Services is the couchbase service to run on the node.
	Services *[]Service `json:"services,omitempty"`
}

// Support defines model for Support.
type Support struct {
	// Plan is plan type, either 'Basic', 'Developer Pro', or 'Enterprise'.
	Plan SupportPlan `json:"plan"`

	// Timezone is the standard timezone for the cluster. Should be the TZ identifier.
	Timezone SupportTimezone `json:"timezone"`
}

// SupportPlan is plan type, either 'Basic', 'Developer Pro', or 'Enterprise'.
type SupportPlan string

// SupportTimezone is the standard timezone for the cluster. Should be the TZ identifier.
type SupportTimezone string

// UpdateClusterRequest defines model for UpdateClusterRequest.
type UpdateClusterRequest struct {
	// Description is the new cluster description (up to 1024 characters).
	Description string `json:"description"`

	// Name is the new name of the cluster (up to 256 characters).
	Name          string         `json:"name"`
	ServiceGroups []ServiceGroup `json:"serviceGroups"`
	Support       Support        `json:"support"`
}

// AsDiskAWS returns the union data inside the Node_Disk as a DiskAWS
func (n *Node) AsDiskAWS() (DiskAWS, error) {
	var body DiskAWS
	err := json.Unmarshal(n.Disk, &body)
	return body, err
}

// FromDiskAWS overwrites any union data inside the Node_Disk as the provided DiskAWS
func (n *Node) FromDiskAWS(v DiskAWS) error {
	b, err := json.Marshal(v)
	n.Disk = b
	return err
}

// MergeDiskAWS performs a merge with any union data inside the Node_Disk, using the provided DiskAWS
func (n *Node) MergeDiskAWS(v DiskAWS) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := JsonMerge(n.Disk, b)
	n.Disk = merged
	return err
}

// AsDiskAzure returns the union data inside the Node_Disk as a DiskAzure
func (n *Node) AsDiskAzure() (DiskAzure, error) {
	var body DiskAzure
	err := json.Unmarshal(n.Disk, &body)
	return body, err
}

// FromDiskAzure overwrites any union data inside the Node_Disk as the provided DiskAzure
func (n *Node) FromDiskAzure(v DiskAzure) error {
	b, err := json.Marshal(v)
	n.Disk = b
	return err
}

// MergeDiskAzure performs a merge with any union data inside the Node_Disk, using the provided DiskAzure
func (n *Node) MergeDiskAzure(v DiskAzure) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := JsonMerge(n.Disk, b)
	n.Disk = merged
	return err
}

// AsDiskGCP returns the union data inside the Node_Disk as a DiskGCP
func (n *Node) AsDiskGCP() (DiskGCP, error) {
	var body DiskGCP
	err := json.Unmarshal(n.Disk, &body)
	return body, err
}

// FromDiskGCP overwrites any union data inside the Node_Disk as the provided DiskGCP
func (n *Node) FromDiskGCP(v DiskGCP) error {
	b, err := json.Marshal(v)
	n.Disk = b
	return err
}

// MergeDiskGCP performs a merge with any union data inside the Node_Disk, using the provided DiskGCP
func (n *Node) MergeDiskGCP(v DiskGCP) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := JsonMerge(n.Disk, b)
	n.Disk = merged
	return err
}

func JsonMerge(data, patch json.RawMessage) (json.RawMessage, error) {
	merger := jsonmerge.Merger{
		CopyNonexistent: true,
	}
	if data == nil {
		data = []byte(`{}`)
	}
	if patch == nil {
		patch = []byte(`{}`)
	}
	merged, err := merger.MergeBytes(data, patch)
	if err != nil {
		return nil, err
	}
	return merged, nil
}
