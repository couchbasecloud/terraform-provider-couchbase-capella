package cluster

import (
	"terraform-provider-capella/internal/api"

	"github.com/google/uuid"
)

// Availability defines model for Availability.
type Availability struct {
	// Type is availability zone type, either 'single' or 'multi'.
	Type AvailabilityType `json:"type"`
}

// AvailabilityType is availability zone type, either 'single' or 'multi'.
type AvailabilityType string

// CreateClusterRequest defines model for CreateClusterRequest.
type CreateClusterRequest struct {
	Availability Availability `json:"availability"`

	// CloudProvider is the cloud provider where the cluster will be hosted.
	// To learn more, see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	CloudProvider   CloudProvider    `json:"cloudProvider"`
	CouchbaseServer *CouchbaseServer `json:"couchbaseServer,omitempty"`

	// Description depicts description of the cluster (up to 1024 characters).
	Description *string `json:"description,omitempty"`

	// Name is the name of the cluster (up to 256 characters).
	Name string `json:"name"`

	// ServiceGroups is the couchbase service groups to be run. At least one
	// service group must contain the data service.
	ServiceGroups []ServiceGroup `json:"serviceGroups"`
	Support       Support        `json:"support"`
}

// CreateClusterResponse defines model for CreateClusterResponse.
type CreateClusterResponse struct {
	// Id The ID of the cluster created.
	Id uuid.UUID `json:"id"`
}

// GetClusterResponse defines model for GetClusterResponse.
type GetClusterResponse struct {
	// AppServiceId is the ID of the linked app service.
	AppServiceId *uuid.UUID             `json:"appServiceId,omitempty"`
	Audit        api.CouchbaseAuditData `json:"audit"`
	Availability Availability           `json:"availability"`

	// CloudProvider is the cloud provider where the cluster will be hosted. To learn more,
	// see [Amazon Web Services](https://docs.couchbase.com/cloud/reference/aws.html).
	CloudProvider   CloudProvider   `json:"cloudProvider"`
	CouchbaseServer CouchbaseServer `json:"couchbaseServer"`
	CurrentState    State           `json:"currentState"`

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
