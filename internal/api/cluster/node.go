package cluster

import (
	"encoding/json"
)

// Node defines model for Node.
type Node struct {
	// Compute Following are the supported compute combinations for CPU
	// and RAM for different cloud providers. To learn more,
	// see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	Compute Compute         `json:"compute"`
	Disk    json.RawMessage `json:"disk"`
}

// Compute Following are the supported compute combinations for CPU
// and RAM for different cloud providers. To learn more,
// see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
type Compute struct {
	// Cpu depicts cpu units (cores).
	Cpu int `json:"cpu"`

	// Ram depicts ram units (GB).
	Ram int `json:"ram"`
}

// DiskAWS defines model for DiskAWS.
type DiskAWS struct {
	// Iops Please refer to documentation for supported IOPS.
	Iops int `json:"iops"`

	// Storage depicts storage in GB. See documentation for supported storage.
	Storage int `json:"storage"`

	// Type depicts type of disk. Please choose from the given list for
	// AWS cloud provider.
	Type DiskAWSType `json:"type"`
}

// DiskAWSType depicts type of disk. Please choose from the given list
// for AWS cloud provider.
type DiskAWSType string

// DiskAzure defines model for DiskAzure.
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

// AsDiskAWS returns the disk data as a DiskAWS
func (n *Node) AsDiskAWS() (DiskAWS, error) {
	var body DiskAWS
	err := json.Unmarshal(n.Disk, &body)
	return body, err
}

// FromDiskAWS overwrites any disk data inside as the provided DiskAWS
func (n *Node) FromDiskAWS(v DiskAWS) error {
	b, err := json.Marshal(v)
	n.Disk = b
	return err
}

// AsDiskAzure returns the disk data as a DiskAzure
func (n *Node) AsDiskAzure() (DiskAzure, error) {
	var body DiskAzure
	err := json.Unmarshal(n.Disk, &body)
	return body, err
}

// FromDiskAzure overwrites any disk data as the provided DiskAzure
func (n *Node) FromDiskAzure(v DiskAzure) error {
	b, err := json.Marshal(v)
	n.Disk = b
	return err
}

// AsDiskGCP returns the disk data as a DiskGCP
func (n *Node) AsDiskGCP() (DiskGCP, error) {
	var body DiskGCP
	err := json.Unmarshal(n.Disk, &body)
	return body, err
}

// FromDiskGCP overwrites any disk data as the provided DiskGCP
func (n *Node) FromDiskGCP(v DiskGCP) error {
	b, err := json.Marshal(v)
	n.Disk = b
	return err
}
