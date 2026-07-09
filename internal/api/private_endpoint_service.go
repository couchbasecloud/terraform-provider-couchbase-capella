package api

// GetPrivateEndpointServiceStatusResponse is the response received from the Capella V4 Public API
// when getting private endpoint service status.
type GetPrivateEndpointServiceStatusResponse struct {
	Enabled bool `json:"enabled"`

	// Status is the lifecycle state of the private endpoint service derived from
	// the most recent enable/disable/update operation (for example "enableFailed"
	// or "enabling"). It is optional and best-effort: it is omitted on GCP, when
	// the status feature flag is disabled, and on older control planes, in which
	// case callers fall back to the Enabled boolean.
	Status *string `json:"status,omitempty"`

	// ServiceName is the CSP endpoint service name that customer VPC/VNET
	// endpoints connect to (for AWS, the VPC endpoint service name). It is
	// available once the private endpoint service is enabled so it can be fed
	// directly into the customer-side endpoint resource. It is optional and
	// best-effort: it is omitted when the service is not enabled, when the name
	// cannot be determined, and on older control planes.
	ServiceName *string `json:"serviceName,omitempty"`

	PrivateDns string `json:"privateDns"`
}
