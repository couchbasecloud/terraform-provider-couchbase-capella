package cluster

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"

	"github.com/google/uuid"
)

// Availability defines if the cluster nodes will be deployed in multiple or single availability zones in the cloud.
type Availability struct {
	// Type is availability zone type, either 'single' or 'multi'.
	Type AvailabilityType `json:"type"`
}

// AvailabilityType is availability zone type, either 'single' or 'multi'.
type AvailabilityType string

// ConfigurationType defines model for ConfigurationType, either 'multiNode' or 'singleNode'.
type ConfigurationType string

// CreateClusterRequest is the request payload sent to the Capella V4 Public API in order to create a new cluster.
// A Couchbase cluster consists of one or more instances of Couchbase Capella, each running on an independent node.
// Data and services are shared across the cluster.
// A cluster may be referred to as a "database" in the documentation and in the Couchbase Capella user interface.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateClusterRequest struct {
	CouchbaseServer   *CouchbaseServer  `json:"couchbaseServer,omitempty"`
	Description       *string           `json:"description,omitempty"`
	CloudProvider     CloudProvider     `json:"cloudProvider"`
	Support           Support           `json:"support"`
	Availability      Availability      `json:"availability"`
	ConfigurationType ConfigurationType `json:"configurationType"`
	Name              string            `json:"name"`
	ServiceGroups     []ServiceGroup    `json:"serviceGroups"`
}

// CreateClusterResponse is the response received from the Capella V4 Public API when asked to create a new cluster.
type CreateClusterResponse struct {
	// Id is the UUID of the cluster created.
	Id uuid.UUID `json:"id"`
}

// GetClusterResponse is the response received from the Capella V4 Public API when asked to fetch details of an existing cluster.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// Project Viewer
// Database Data Reader/Writer
// Database Data Reader
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetClusterResponse struct {
	CouchbaseServer   CouchbaseServer   `json:"couchbaseServer"`
	AppServiceId      *uuid.UUID        `json:"appServiceId,omitempty"`
	CloudProvider     CloudProvider     `json:"cloudProvider"`
	Support           Support           `json:"support"`
	CurrentState      State             `json:"currentState"`
	ConfigurationType ConfigurationType `json:"configurationType"`
	Availability      Availability      `json:"availability"`
	Description       string            `json:"description"`
	Name              string            `json:"name"`
	Etag              string
	Audit             api.CouchbaseAuditData `json:"audit"`
	ServiceGroups     []ServiceGroup         `json:"serviceGroups"`
	Id                uuid.UUID              `json:"id"`
}

// UpdateClusterRequest is the payload sent to the Capella V4 Public API when asked to update an existing cluster.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type UpdateClusterRequest struct {
	Support       Support        `json:"support"`
	Description   string         `json:"description"`
	Name          string         `json:"name"`
	ServiceGroups []ServiceGroup `json:"serviceGroups"`
}
