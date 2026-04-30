package api

// GetPrivateEndpointResponse defines model for a CSP's private endpoint.
type GetPrivateEndpointResponse struct {
	// Id is the endpoint id.
	Id string `json:"id"`
	// Status is the endpoint status.  Possible values are failed, linked, pending, pendingAcceptance, rejected and unrecognized.
	Status string `json:"status"`
	// ServiceName is the name of the endpoint service.
	ServiceName string `json:"serviceName"`
}

// GetPrivateEndpointsResponse is a list of private endpoints.
type GetPrivateEndpointsResponse struct {
	Endpoints []GetPrivateEndpointResponse `json:"endpoints"`
}
