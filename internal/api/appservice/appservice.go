package appservice

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"

	"github.com/google/uuid"
)

// CreateAppServiceRequest is the request payload sent to the Capella V4 Public API in order to create a new app service.
// An App Service synchronizes data between the Couchbase Capella database and any mobile applications.
// App Service is a fully managed application backend designed to provide data synchronization for mobile/IoT applications and the Capella Cloud Service.
//
// To learn more about App Services, see https://docs.couchbase.com/cloud/app-services/index.html
//
// In order to access this endpoint, the provided API key must have at least one of the roles referenced below:
//
// Organization Owner
// Project Owner
// Project Manager
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateAppServiceRequest struct {
	// Description is the description for the app service (up to 256 characters).
	Description *string `json:"description,omitempty"`

	// Nodes is the number of nodes configured for the App Service.
	// The number of nodes can range from 2 to 12
	Nodes *int64 `json:"nodes,omitempty"`

	// Version is version of the App Service Server to be installed.
	// The latest Server version will be deployed by default.
	Version *string `json:"version,omitempty"`

	// Name is the name of the app service, the name of the app service should follow this naming criteria:
	// An app service name should have at least 2 characters and up to 256 characters.
	Name string `json:"name"`

	// Compute is the CPU and RAM configuration of the app service.
	Compute AppServiceCompute `json:"compute"`
}

// CreateAppServiceResponse is the response received from the Capella V4 Public API when asked to create a new app service.
type CreateAppServiceResponse struct {
	// ID is the UUID of the app service
	Id uuid.UUID `json:"id"`
}

// GetAppServiceResponse is the response received from the Capella V4 Public API when asked to fetch details of an existing app service.
//
// In order to access this endpoint, the provided API key must have at least one of the roles referenced below:
// Organization Owner
// Project Owner
// Project Manager
// Project Viewer
// Database Data Reader/Writer
// Database Data Reader
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetAppServiceResponse struct {
	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId string `json:"organizationId"`

	// Name is the name of the app service, the name of the app service should follow this naming criteria:
	// An app service name should have at least 2 characters and up to 256 characters.
	Name string `json:"name"`

	// Description is the description for the app service (up to 256 characters).
	Description string `json:"description"`

	// CloudProvider is the cloud provider where the app service will be hosted.
	// To learn more, see:
	// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
	// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
	// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
	CloudProvider string `json:"cloudProvider"`

	// ProjectId is the projectId of the cluster.
	ProjectId string `json:"projectId"`

	// ClusterId is the clusterId of the cluster.
	ClusterId string `json:"clusterId"`

	// CurrentState defines the current state of app service.
	CurrentState State `json:"currentState"`

	// Version defines the version of the app service server.
	Version string `json:"version"`

	// Etag represents the version of the document
	Etag string

	// Audit contains all audit-related fields.
	Audit api.CouchbaseAuditData `json:"audit"`

	// Compute is the CPU and RAM configuration of the app service.
	Compute AppServiceCompute `json:"compute"`

	// Nodes is the number of nodes configured for the app service.
	Nodes int `json:"nodes"`

	// Id is the ID of the app service created.
	Id uuid.UUID `json:"id"`
}

type UpdateAppServiceRequest struct {
	// Nodes is the number of nodes configured for the App Service.
	// The number of nodes can range from 2 to 12
	Nodes int64 `json:"nodes"`

	// Compute is the CPU and RAM configuration of the app service.
	Compute AppServiceCompute `json:"compute"`
}
