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
	Description *string           `json:"description,omitempty"`
	Nodes       *int64            `json:"nodes,omitempty"`
	Version     *string           `json:"version,omitempty"`
	Name        string            `json:"name"`
	Compute     AppServiceCompute `json:"compute"`
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
	OrganizationId string `json:"organizationId"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	CloudProvider  string `json:"cloudProvider"`
	ProjectId      string `json:"projectId"`
	ClusterId      string `json:"clusterId"`
	CurrentState   State  `json:"currentState"`
	Version        string `json:"version"`
	Etag           string
	Audit          api.CouchbaseAuditData `json:"audit"`
	Compute        AppServiceCompute      `json:"compute"`
	Nodes          int                    `json:"nodes"`
	Id             uuid.UUID              `json:"id"`
}

type UpdateAppServiceRequest struct {
	// Nodes is the number of nodes configured for the App Service.
	// The number of nodes can range from 2 to 12
	Nodes int64 `json:"nodes"`

	// Compute is the CPU and RAM configuration of the app service.
	Compute AppServiceCompute `json:"compute"`
}
