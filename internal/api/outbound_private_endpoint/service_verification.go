package outbound_private_endpoint

// VerifyServiceRequest is the body of a service verification.
type VerifyServiceRequest struct {
	// ServiceName is the customer's private link service to check.
	ServiceName string `json:"serviceName"`
}

// ServiceVerification is the result of a service verification. A failed check is reported
// here with a 200, not as an error.
type ServiceVerification struct {
	CloudProvider CloudProvider `json:"cloudProvider"`
	ServiceName   string        `json:"serviceName"`

	// Ready is true when no blocking check failed.
	Ready bool `json:"ready"`

	// Checks has one entry per check, in a fixed order.
	Checks []ServiceVerificationCheck `json:"checks"`

	// ProviderConfig is absent until the service is found.
	ProviderConfig *ServiceVerificationProviderConfig `json:"providerConfig,omitempty"`
}

// ServiceVerificationCheck is the result of one check.
type ServiceVerificationCheck struct {
	Name   string `json:"name"`
	Result string `json:"result"`

	// Message is present when the check failed or warned, and on GCP's notApplicable accessGranted.
	Message *string `json:"message,omitempty"`
}

// ServiceVerificationProviderConfig is what the cloud reports about the service. Only the
// cluster's cloud is present, and GCP and Azure have no fields yet.
type ServiceVerificationProviderConfig struct {
	AWS *AWSServiceVerification `json:"aws,omitempty"`
}

// AWSServiceVerification is what AWS reports about an endpoint service.
type AWSServiceVerification struct {
	AvailabilityZones          []string `json:"availabilityZones"`
	AcceptanceRequired         bool     `json:"acceptanceRequired"`
	ServiceOwner               string   `json:"serviceOwner"`
	PrivateDnsName             *string  `json:"privateDnsName,omitempty"`
	PrivateDnsNameVerification string   `json:"privateDnsNameVerification"`
}
