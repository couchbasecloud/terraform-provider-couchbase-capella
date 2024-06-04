package network_peer

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// CreateNetworkPeeringRequest Creates a network peering record for Capella. Capella does not support peering of networks between different cloud providers. For example, you cannot peer GCP VPC that hosts Capella cluster with an AWS VPC hosting an application.
//
// - Create configures a Couchbase Capella private networking with the cloud provider. Setting up a private network enables your application to interact with Couchbase Capella over a private connection by co-locating them through network peering.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
// - Organization Owner
// - Project Owner
// - Project Manager
//
// To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type CreateNetworkPeeringRequest struct {
	// Name is the name of the peering relationship. -  The name of the peering relationship must be at least 2 characters long. -  The name can not exceed 128 characters.
	Name string `json:"name"`

	// ProviderConfig The config data for a peering relationship for a cluster on AWS, GCP.
	ProviderConfig json.RawMessage `json:"providerConfig"`

	// ProviderType Type of the cloud provider for which the peering connection is created. Which are- 1. aws 2. gcp
	ProviderType string `json:"providerType"`
}

// CreateNetworkPeeringResponse is the response received from the Capella V4 Public API when asked to create a new network peering connection.
type CreateNetworkPeeringResponse struct {
	// Id is the ID is the unique UUID generated when a VPC record is created.
	Id uuid.UUID `json:"id"`
}

// GetNetworkPeeringRecordResponse Fetches the details of the network peering meta data based on the peerID provided.
//
//	In order to access this endpoint, the provided API key must have at least one of the following roles:
//	 - Organization Owner
//	 - Project Owner
//	 - Project Manager
//
//	To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type GetNetworkPeeringRecordResponse struct {
	Audit api.CouchbaseAuditData `json:"audit"`

	// Commands contains the list of commands that the user must execute in order to complete the association of the network.
	Commands []string `json:"commands"`

	// Id The ID is the unique UUID generated when a VPC record is created.
	Id string `json:"id"`

	// Name is the name of the peering relationship.
	Name string `json:"name"`

	// ProviderConfig This provides details about the configuration and the ID of the VPC peer on AWS, GCP.
	ProviderConfig json.RawMessage `json:"providerConfig"`
	Status         PeeringStatus   `json:"status"`
}
