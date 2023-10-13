package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-capella/internal/errors"
)

type Bucket struct {
	Id types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the database credential needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	Type types.String `tfsdk:"type"`

	StorageBackend types.String `tfsdk:"storage_backend"`

	MemoryAllocationInMb int `tfsdk:"memory_allocationinmb"`

	BucketConflictResolution types.String `tfsdk:"conflict_resolution"`

	DurabilityLevel types.String `tfsdk:"durability_level"`

	Replicas int `tfsdk:"replicas"`

	Flush bool `tfsdk:"flush"`

	TimeToLiveInSeconds int `tfsdk:"ttl"`

	EvictionPolicy types.String `tfsdk:"eviction_policy"`

	Stats types.Object `tfsdk:"stats"`
}

type Stats struct {
	ItemCount       types.Int64 `tfsdk:"item_count"`
	OpsPerSecond    types.Int64 `tfsdk:"ops_per_second"`
	DiskUsedInMib   types.Int64 `tfsdk:"disk_used_in_mib"`
	MemoryUsedInMib types.Int64 `tfsdk:"memory_used_in_mib"`
}

type OneBucket struct {
	Id types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the database credential needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	Type types.String `tfsdk:"type"`

	StorageBackend types.String `tfsdk:"storage_backend"`

	MemoryAllocationInMb int `tfsdk:"memory_allocationinmb"`

	BucketConflictResolution types.String `tfsdk:"conflict_resolution"`

	DurabilityLevel types.String `tfsdk:"durability_level"`

	Replicas int `tfsdk:"replicas"`

	Flush bool `tfsdk:"flush"`

	TimeToLiveInSeconds int `tfsdk:"ttl"`

	EvictionPolicy types.String `tfsdk:"eviction_policy"`

	Stats *Stats `tfsdk:"stats"`
}

// Validate will split the IDs by a delimiter i.e. comma , in case a terraform import CLI is invoked.
// The format of the terraform import CLI would include the IDs as follows -
// `terraform import capella_bucket.new_bucket id=<uuid>,cluster_id=<uuid>,project_id=<uuid>,organization_id=<uuid>`
func (c Bucket) Validate() (bucketId, clusterId, projectId, organizationId string, err error) {

	const (
		idDelimiter       = ","
		organizationIdSep = "organization_id="
		projectIdSep      = "project_id="
		clusterIdSep      = "cluster_id="
		bucketIdSep       = "id="
	)

	organizationId = c.OrganizationId.ValueString()
	projectId = c.ProjectId.ValueString()
	clusterId = c.ClusterId.ValueString()
	bucketId = c.Id.ValueString()
	var found bool

	// check if the id is a comma separated string of multiple IDs, usually passed during the terraform import CLI
	if c.OrganizationId.IsNull() {
		strs := strings.Split(c.Id.ValueString(), idDelimiter)
		if len(strs) != 4 {
			err = errors.ErrIdMissing
			return
		}
		_, bucketId, found = strings.Cut(strs[0], bucketIdSep)
		if !found {
			err = errors.ErrDatabaseCredentialIdMissing
			return
		}

		_, clusterId, found = strings.Cut(strs[1], clusterIdSep)
		if !found {
			err = errors.ErrClusterIdMissing
			return
		}

		_, projectId, found = strings.Cut(strs[2], projectIdSep)
		if !found {
			err = errors.ErrProjectIdMissing
			return
		}

		_, organizationId, found = strings.Cut(strs[3], organizationIdSep)
		if !found {
			err = errors.ErrOrganizationIdMissing
			return
		}
	}

	if bucketId == "" {
		err = errors.ErrBucketIdCannotBeEmpty
		return
	}

	if clusterId == "" {
		err = errors.ErrClusterIdCannotBeEmpty
		return
	}

	if projectId == "" {
		err = errors.ErrProjectIdCannotBeEmpty
		return
	}

	if organizationId == "" {
		err = errors.ErrOrganizationIdCannotBeEmpty
		return
	}

	return bucketId, clusterId, projectId, organizationId, nil
}
