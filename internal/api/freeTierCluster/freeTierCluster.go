package freeTierClusterapi

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/google/uuid"
	"time"
)

type Type string

// CouchbaseServer is the version of the Couchbase Server installed in the free-tier cluster as received by the V4 Capella Public API.
type CouchbaseServer struct {
	Version *string `json:"version"`
}

// Availability is the availability type of the free-tier cluster as received by the V4 Capella Public API. It is single zone for free-tier clusters.
type Availability struct {
	Type string `json:"type"`
}

// Support is the support plan for the free-tier cluster as received by the V4 Capella Public API.
type Support struct {
	Plan     string `json:"plan"`
	Timezone string `json:"timezone"`
}

// Audit is the audit data for the free-tier cluster as received by the V4 Capella Public API.
type Audit struct {
	CreatedBy  string    `json:"createdBy"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedBy string    `json:"modifiedBy"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Version    int       `json:"version"`
}

// GetFreeTierCluster Response as received by the V4 Capella Public API.
type GetFreeTierClusterResponse struct {
	ID                         uuid.UUID              `json:"id"`
	AppServiceID               string                 `json:"appServiceId"`
	Name                       string                 `json:"name"`
	Description                string                 `json:"description"`
	ConnectionString           string                 `json:"connectionString"`
	CloudProvider              cluster.CloudProvider  `json:"cloudProvider"`
	CouchbaseServer            CouchbaseServer        `json:"couchbaseServer"`
	ServiceGroups              []cluster.ServiceGroup `json:"serviceGroups"`
	Availability               Availability           `json:"availability"`
	Support                    Support                `json:"support"`
	CurrentState               string                 `json:"currentState"`
	Audit                      api.CouchbaseAuditData `json:"audit"`
	CmekID                     string                 `json:"cmekId"`
	EnablePrivateDNSResolution bool                   `json:"enablePrivateDNSResolution"`
}

// CreateFreeTierClusterRequest is the request body for creating a free-tier cluster as sent to the V4 Capella Public API.
type CreateFreeTierClusterRequest struct {
	// Name is the name of the cluster (up to 256 characters).
	Name string `json:"name"`
	// Description depicts description of the cluster (up to 1024 characters).
	Description *string `json:"description,omitempty"`
	// CloudProvider is the cloud provider where the cluster will be hosted.
	// To learn more, see:
	// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
	// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
	// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
	CloudProvider cluster.CloudProvider `json:"cloudProvider"`
}

// UpdateFreeTierClusterRequest is the request body for updating a free-tier cluster as sent to the V4 Capella Public API.
type UpdateFreeTierClusterRequest struct {
	// Name is the name of the cluster (up to 256 characters).
	Name string `json:"name"`
	// Description depicts description of the cluster (up to 1024 characters).
	Description string `json:"description,omitempty"`
}
