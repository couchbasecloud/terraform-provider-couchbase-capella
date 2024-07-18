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

//// ConfigurationType defines model for ConfigurationType, either 'multiNode' or 'singleNode'.
//type ConfigurationType string

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
	// CouchbaseServer is the version of the Couchbase Server to be installed in the cluster.
	// Refer to documentation here (https://docs.couchbase.com/cloud/clusters/upgrade-database.html#server-version-maintenance-support)
	// for list of supported versions.
	// The latest Couchbase Server version will be deployed by default.
	CouchbaseServer *CouchbaseServer `json:"couchbaseServer,omitempty"`

	// Description depicts description of the cluster (up to 1024 characters).
	Description *string `json:"description,omitempty"`

	// CloudProvider is the cloud provider where the cluster will be hosted.
	// To learn more, see:
	// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
	// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
	// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
	CloudProvider CloudProvider `json:"cloudProvider"`

	// Support defines the support plan and timezone for this particular cluster.
	Support Support `json:"support"`

	// Availability zone type, either 'single' or 'multi'.
	// Enum: "single" "multi"
	Availability Availability `json:"availability"`

	//// ConfigurationType defines model for ConfigurationType, either 'multiNode' or 'singleNode'
	//ConfigurationType ConfigurationType `json:"configurationType"`

	// Name is the name of the cluster (up to 256 characters).
	Name string `json:"name"`

	// ServiceGroups is the couchbase service groups to be run. At least one
	// service group must contain the data service.
	ServiceGroups []ServiceGroup `json:"serviceGroups"`
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
	// CouchbaseServer is the version of the Couchbase Server to be installed in the cluster.
	// Refer to documentation here (https://docs.couchbase.com/cloud/clusters/upgrade-database.html#server-version-maintenance-support)
	// for list of supported versions.
	// The latest Couchbase Server version will be deployed by default.
	CouchbaseServer CouchbaseServer `json:"couchbaseServer"`

	// AppServiceId is the ID of the linked app service.
	AppServiceId *uuid.UUID `json:"appServiceId,omitempty"`

	// CloudProvider is the cloud provider where the cluster will be hosted.
	// To learn more, see:
	// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
	// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
	// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
	CloudProvider CloudProvider `json:"cloudProvider"`

	// Support defines the support plan and timezone for this particular cluster.
	Support Support `json:"support"`

	// CurrentState tells the status of the cluster - if it's healthy or degraded.
	CurrentState State `json:"currentState"`

	//// ConfigurationType defines model for ConfigurationType, either 'multiNode' or 'singleNode'
	//ConfigurationType ConfigurationType `json:"configurationType"`

	// Availability zone type, either 'single' or 'multi'.
	// Enum: "single" "multi"
	Availability Availability `json:"availability"`

	// Description depicts description of the cluster (up to 1024 characters).
	Description string `json:"description"`

	// Name is the name of the cluster (up to 256 characters).
	Name string `json:"name"`

	// Etag represents the version of the document
	Etag  string
	Audit api.CouchbaseAuditData `json:"audit"`

	// ServiceGroups is the couchbase service groups to be run. At least one
	// service group must contain the data service.
	ServiceGroups []ServiceGroup `json:"serviceGroups"`

	// Id is the ID of the cluster created.
	Id uuid.UUID `json:"id"`
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
	// Support defines the support plan and timezone for this particular cluster.
	Support Support `json:"support"`

	// Description is the new cluster description (up to 1024 characters).
	Description string `json:"description"`

	// Name is the new name of the cluster (up to 256 characters).
	Name string `json:"name"`

	// ServiceGroups is the couchbase service groups to be run. At least one
	// service group must contain the data service.
	ServiceGroups []ServiceGroup `json:"serviceGroups"`
}
