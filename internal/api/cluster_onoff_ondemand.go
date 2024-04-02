package api

// CreateClusterOnRequest is the payload passed to V4 Capella Public API to turn on the cluster.
// Turn cluster on.
//
// In order to access this endpoint, the provided API key must have at least one of the roles referenced below:
// - Organization Owner
// - Project Owner
//
// To learn more, see [Organization, Project, and Database Access Overview](https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html).
type CreateClusterOnRequest struct {
	// TurnOnLinkedAppService Set this value to true if you want to turn on the app service linked with the cluster, false if not.
	// If set to true, the app service, if present, will turn on with the cluster.
	// Default value for this is false, which means the linked app service will be kept off.
	TurnOnLinkedAppService bool `json:"turnOnLinkedAppService"`
}
