package appservice

import (
	"github.com/google/uuid"
	"terraform-provider-capella/internal/api"
)

// CreateAppServiceRequest defines the model for CreateAppServiceRequest.
type CreateAppServiceRequest struct {
	// Name is the name of the app service, the name of the app service should follow this naming criteria:
	// An app service name should have at least 2 characters and up to 256 characters.
	Name string `json:"name"`

	// Description is the description for the app service.
	Description string `json:"description"`

	// Nodes is the number of nodes configured for the app service.
	Nodes int64 `json:"nodes"`

	// Compute is the CPU and RAM configuration of the app service.
	Compute Compute `tfsdk:"compute"`

	// Version is version of the App Service Server to be installed.
	// The latest Server version will be deployed by default.
	Version *string `json:"version,omitempty"`
}

// CreateAppServiceResponse defines the model for CreateAppServiceResponse.
type CreateAppServiceResponse struct {
	// ID is the ID of the app service
	Id uuid.UUID `json:"id"`
}

// GetAppServiceResponse defines the model for GetAppServiceResponse.
type GetAppServiceResponse struct {
	// Id is the ID of the app service.
	Id uuid.UUID `json:"id"`

	// Name is the name of the app service, the name of the app service should follow this naming criteria:
	// An app service name should have at least 2 characters and up to 256 characters.
	Name string `json:"name"`

	// Description is the description for the app service.
	Description string `json:"description"`

	// CloudProvider is the cloud provider where the cluster will be hosted.
	CloudProvider string `json:"cloudProvider"`

	// Nodes is the number of nodes configured for the app service.
	Nodes int `json:"nodes"`

	// Compute is the CPU and RAM configuration of the app service.
	Compute Compute `tfsdk:"compute"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId string `json:"organizationId"`

	// ProjectId is the projectId of the cluster.
	ProjectId string `tfsdk:"projectId"`

	// ClusterId is the clusterId of the cluster.
	ClusterId string `tfsdk:"clusterId"`

	// CurrentState defines the current state of app service.
	CurrentState string `json:"currentState"`

	// Version defines the version of the app service server
	Version string `tfsdk:"version"`

	// Audit contains all audit-related fields.
	Audit api.CouchbaseAuditData `json:"audit"`
}

// GetAppServicesResponse defines the model for GetAppServicesResponse
type GetAppServicesResponse struct {
	Data []GetAppServiceResponse `json:"data"`
}
