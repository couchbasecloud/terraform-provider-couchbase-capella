package outbound_private_endpoint

// GetOutboundPrivateEndpointCapacityResponse is a cluster's outbound private endpoint capacity.
type GetOutboundPrivateEndpointCapacityResponse struct {
	ClusterEndpointLimit OutboundPrivateEndpointCapacityUsage `json:"clusterEndpointLimit"`
}

// OutboundPrivateEndpointCapacityUsage is the cluster's endpoints against its own limit.
type OutboundPrivateEndpointCapacityUsage struct {
	// Used counts every endpoint on the cluster, including pending, failed and deleting ones.
	Used int64 `json:"used"`

	// Limit is 3 unless Couchbase Support raises it.
	Limit int64 `json:"limit"`
}
