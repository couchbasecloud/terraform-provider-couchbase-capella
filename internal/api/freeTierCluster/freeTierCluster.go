package freeTierClusterapi

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

type Type string

// CreateFreeTierClusterRequest is the request body for creating a free-tier cluster as sent to the V4 Capella Public API.
type CreateFreeTierClusterRequest struct {
	// Name is the name of the cluster (up to 256 characters).
	Name string `json:"name"`
	// Description depicts description of the cluster (up to 1024 characters).
	Description *string `json:"description,omitempty"`
	// CloudProvider is the cloud provider where the cluster will be hosted.
	// To learn more, see:
	// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
	// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
	// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
	CloudProvider cluster.CloudProvider `json:"cloudProvider"`
}

// UpdateFreeTierClusterRequest is the request body for updating a free-tier cluster as sent to the V4 Capella Public API.
type UpdateFreeTierClusterRequest struct {
	// Name is the name of the cluster (up to 256 characters).
	Name string `json:"name"`
	// Description depicts description of the cluster (up to 1024 characters).
	Description string `json:"description,omitempty"`
}
