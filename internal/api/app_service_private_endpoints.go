package api

// GetAppServicePrivateEndpointResponse defines model for a CSP's private endpoint
// connected to an App Service's private endpoint service.
type GetAppServicePrivateEndpointResponse struct {
	// Id is the endpoint id.
	Id string `json:"id"`
	// Status is the endpoint status. Possible values are failed, linked, pending, pendingAcceptance, rejected and unrecognized.
	Status string `json:"status"`
	// ServiceName is the name of the endpoint service.
	ServiceName string `json:"serviceName"`
}

// GetAppServicePrivateEndpointsResponse is a list of private endpoints connected to an
// App Service's private endpoint service.
type GetAppServicePrivateEndpointsResponse struct {
	Endpoints []GetAppServicePrivateEndpointResponse `json:"endpoints"`
}
