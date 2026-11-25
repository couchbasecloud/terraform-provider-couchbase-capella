package outbound_private_endpoint

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// Statuses an outbound private endpoint can report.
const (
	StatusPending           = "pending"
	StatusPendingAcceptance = "pendingAcceptance"
	StatusLinked            = "linked"
	StatusRejected          = "rejected"
	StatusFailed            = "failed"
	StatusDeleting          = "deleting"
	StatusUnrecognized      = "unrecognized"
)

// FailureReasonDeleteFailed is the failure reason of an endpoint whose teardown failed.
const FailureReasonDeleteFailed = "deleteFailed"

// Error codes the outbound private endpoint routes answer with when the endpoint
// or the cluster is busy. Both are worth retrying.
const (
	// ErrCodeBusy is returned when updating an endpoint that is pending or deleting.
	ErrCodeBusy = 16016

	// ErrCodeOperationInProgress is returned when another change to the cluster's
	// outbound private endpoints holds the cluster.
	ErrCodeOperationInProgress = 16017
)

// CreateOutboundPrivateEndpointRequest is the body of a create.
type CreateOutboundPrivateEndpointRequest struct {
	// Name is Capella's label for the endpoint.
	Name string `json:"name"`

	// EndpointType is the AWS VPC endpoint type. interface is the only type.
	EndpointType string `json:"endpointType,omitempty"`

	// ServiceName is the customer's private link service.
	ServiceName string `json:"serviceName"`

	// Hostname is the customer's own hostname. Left out to have the cloud or Capella assign one.
	Hostname string `json:"hostname,omitempty"`
}

// CreateOutboundPrivateEndpointResponse is the body of a create's 202.
type CreateOutboundPrivateEndpointResponse struct {
	Id string `json:"id"`
}

// UpdateOutboundPrivateEndpointRequest is the body of an update. The name is the only updatable field.
type UpdateOutboundPrivateEndpointRequest struct {
	Name string `json:"name"`
}

// GetOutboundPrivateEndpointResponse is an outbound private endpoint as the API returns it.
type GetOutboundPrivateEndpointResponse struct {
	Id                 string                 `json:"id"`
	Name               string                 `json:"name"`
	CloudProvider      CloudProvider          `json:"cloudProvider"`
	EndpointType       string                 `json:"endpointType"`
	ServiceName        string                 `json:"serviceName"`
	ProviderEndpointId *string                `json:"providerEndpointId,omitempty"`
	Status             string                 `json:"status"`
	FailureReason      *string                `json:"failureReason,omitempty"`
	ProviderStatus     *string                `json:"providerStatus,omitempty"`
	Hostname           *string                `json:"hostname,omitempty"`
	IsCustomHostname   bool                   `json:"isCustomHostname"`
	PrivateDnsEnabled  bool                   `json:"privateDnsEnabled"`
	IpAddresses        []string               `json:"ipAddresses"`
	ProviderConfig     *ProviderConfig        `json:"providerConfig,omitempty"`
	Audit              api.CouchbaseAuditData `json:"audit"`
}

// CloudProvider is the cluster's cloud provider and region.
type CloudProvider struct {
	Type   string `json:"type"`
	Region string `json:"region"`
}

// ProviderConfig holds the cloud-specific details. Only the cluster's cloud is present,
// and GCP and Azure have no fields yet.
type ProviderConfig struct {
	AWS *AWSConfig `json:"aws,omitempty"`
}

// AWSConfig is the AWS-specific detail of an endpoint.
type AWSConfig struct {
	RegionalDnsName   *string        `json:"regionalDnsName,omitempty"`
	IpAddresses       []AWSIpAddress `json:"ipAddresses"`
	AvailabilityZones []string       `json:"availabilityZones"`
}

// AWSIpAddress is one of the endpoint's addresses with the zone it is in.
type AWSIpAddress struct {
	IpAddress        string `json:"ipAddress"`
	AvailabilityZone string `json:"availabilityZone"`
}
