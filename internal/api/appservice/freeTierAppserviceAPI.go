package appservice

// CreateFreeTierAppServiceRequest is the request payload sent to the Capella V4 Public API in order to create a new free-tier app service.
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

type CreateFreeTierAppServiceRequest struct {
	// Description is the description for the free-tier app service (up to 256 characters).
	Description *string `json:"description,omitempty"`
	// Name of the free-tier app service.
	Name string `json:"name"`
}

// CreateFreeTierAppServiceResponse is the request payload for the Capella V4 Public API free-tier app service update request.
type UpdateFreeTierAppServiceRequest struct {
	// Description is the description for the free-tier app service (up to 256 characters).
	Description *string `json:"description,omitempty"`
	// Name of the free-tier app service.
	Name string `json:"name"`
}
