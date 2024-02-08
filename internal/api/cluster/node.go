package cluster

import (
	"encoding/json"
)

// Node defines attributes of a cluster node.
type Node struct {
	// Disk is the type of disk that is supported per cloud provider during cluster creation.
	Disk json.RawMessage `json:"disk"`

	// Compute is the family of instances in cloud that are supported during cluster creation.
	// Following are the supported compute combinations for CPU
	// and RAM for different cloud providers. To learn more,
	// see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	Compute Compute `json:"compute"`
}

// Compute depicts the couchbase compute, following are the supported compute combinations for CPU
// and RAM for different cloud providers.
// To learn more, see:
// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
type Compute struct {
	// Cpu depicts cpu units (cores).
	Cpu int `json:"cpu"`

	// Ram depicts ram units (GB).
	Ram int `json:"ram"`
}

// DiskAWS defines the disk metadata as supported by AWS.
type DiskAWS struct {
	// Type depicts type of disk. Please choose from the given list for
	// AWS cloud provider.
	Type DiskAWSType `json:"type"`

	// Iops Please refer to documentation for supported IOPS.
	Iops int `json:"iops"`

	// Storage depicts storage in GB. See documentation for supported storage.
	Storage int `json:"storage"`
}

// DiskAWSType depicts type of disk. Please choose from the given list
// for AWS cloud provider.
type DiskAWSType string

// DiskAzure defines attributes for disks metadata supported in Azure.
type DiskAzure struct {
	// Iops is required for ultra disk types. Please refer to documentation
	// for supported IOPS.
	Iops *int `json:"iops,omitempty"`

	// Storage depicts storage in GB. Required for ultra disk types.
	// Please refer to documentation for supported storage.
	Storage *int `json:"storage,omitempty"`

	// Type depicts type of disk. Please choose from the given list
	// for Azure cloud provider.
	Type DiskAzureType `json:"type"`

	// Autoexpansion determines if disk auto expansion is enabled
	Autoexpansion *bool `json:"autoexpansion,omitempty"`
}

// DiskAzureType depicts type of disk. Please choose from the given list for Azure cloud provider.
type DiskAzureType string

// DiskGCP defines the disk metadata as supported by GCP.
type DiskGCP struct {
	// Type is type of disk. Please choose from the given list for GCP cloud provider.
	Type DiskGCPType `json:"type"`

	// Storage is storage in GB. Please refer to documentation for supported storage.
	Storage int `json:"storage"`
}

// DiskGCPType is type of disk. Please choose from the given list for GCP cloud provider.
type DiskGCPType string

// AsDiskAWS returns the disk data as a DiskAWS.
func (n *Node) AsDiskAWS() (DiskAWS, error) {
	var body DiskAWS
	err := json.Unmarshal(n.Disk, &body)
	return body, err
}

// FromDiskAWS overwrites any disk data inside as the provided DiskAWS.
func (n *Node) FromDiskAWS(v DiskAWS) error {
	b, err := json.Marshal(v)
	n.Disk = b
	return err
}

// AsDiskAzure returns the disk data as a DiskAzure.
func (n *Node) AsDiskAzure() (DiskAzure, error) {
	var body DiskAzure
	err := json.Unmarshal(n.Disk, &body)
	return body, err
}

// FromDiskAzure overwrites any disk data as the provided DiskAzure.
func (n *Node) FromDiskAzure(v DiskAzure) error {
	b, err := json.Marshal(v)
	n.Disk = b
	return err
}

// AsDiskGCP returns the disk data as a DiskGCP.
func (n *Node) AsDiskGCP() (DiskGCP, error) {
	var body DiskGCP
	err := json.Unmarshal(n.Disk, &body)
	return body, err
}

// FromDiskGCP overwrites any disk data as the provided DiskGCP.
func (n *Node) FromDiskGCP(v DiskGCP) error {
	b, err := json.Marshal(v)
	n.Disk = b
	return err
}
