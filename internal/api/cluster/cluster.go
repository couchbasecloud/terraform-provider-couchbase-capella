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
	// Availability zone type, either 'single' or 'multi'.
	// Enum: "single" "multi"
	Availability Availability `json:"availability"`

	// CloudProvider is the cloud provider where the cluster will be hosted.
	// To learn more, see:
	// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
	// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
	// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
	CloudProvider CloudProvider `json:"cloudProvider"`

	// CouchbaseServer is the version of the Couchbase Server to be installed in the cluster.
	// Refer to documentation here (https://docs.couchbase.com/cloud/clusters/upgrade-database.html#server-version-maintenance-support)
	// for list of supported versions.
	// The latest Couchbase Server version will be deployed by default.
	CouchbaseServer *CouchbaseServer `json:"couchbaseServer,omitempty"`

	// Description depicts description of the cluster (up to 1024 characters).
	Description *string `json:"description,omitempty"`

	// Name is the name of the cluster (up to 256 characters).
	Name string `json:"name"`

	// ServiceGroups is the couchbase service groups to be run. At least one
	// service group must contain the data service.
	ServiceGroups []ServiceGroup `json:"serviceGroups"`

	// Support defines the support plan and timezone for this particular cluster.
	Support Support `json:"support"`
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
	// AppServiceId is the ID of the linked app service.
	AppServiceId *uuid.UUID             `json:"appServiceId,omitempty"`
	Audit        api.CouchbaseAuditData `json:"audit"`
	// Availability zone type, either 'single' or 'multi'.
	// Enum: "single" "multi"
	Availability Availability `json:"availability"`

	// CloudProvider is the cloud provider where the cluster will be hosted.
	// To learn more, see:
	// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
	// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
	// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
	CloudProvider CloudProvider `json:"cloudProvider"`

	// CouchbaseServer is the version of the Couchbase Server to be installed in the cluster.
	// Refer to documentation here (https://docs.couchbase.com/cloud/clusters/upgrade-database.html#server-version-maintenance-support)
	// for list of supported versions.
	// The latest Couchbase Server version will be deployed by default.
	CouchbaseServer CouchbaseServer `json:"couchbaseServer"`

	// CurrentState tells the status of the cluster - if it's healthy or degraded.
	CurrentState State `json:"currentState"`

	// Description  depicts description of the cluster (up to 1024 characters).
	Description string `json:"description"`

	// Id is the ID of the cluster created.
	Id uuid.UUID `json:"id"`

	// Name Name of the cluster (up to 256 characters).
	Name string `json:"name"`

	// ServiceGroups is the couchbase service groups to be run. At least one
	// service group must contain the data service.
	ServiceGroups []ServiceGroup `json:"serviceGroups"`

	// Support defines the support plan and timezone for this particular cluster.
	Support Support `json:"support"`

	Etag string
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
	// Description is the new cluster description (up to 1024 characters).
	Description string `json:"description"`

	// Name is the new name of the cluster (up to 256 characters).
	Name string `json:"name"`

	// ServiceGroups is the couchbase service groups to be run. At least one
	// service group must contain the data service.
	ServiceGroups []ServiceGroup `json:"serviceGroups"`

	// Support defines the support plan and timezone for this particular cluster.
	Support Support `json:"support"`
}

// GetClustersResponse is the response received from the Capella V4 Public API when asked to list all clusters.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// Project Manager
// Project Viewer
// Database Data Reader/Writer
// Database Data Reader
// Returned set of clusters is reduced to what the caller has access to view.
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetClustersResponse struct {
	Data []GetClusterResponse `json:"data"`
}
