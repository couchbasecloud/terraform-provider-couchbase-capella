package schema

import (
	"fmt"

	samplebucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/sample_bucket"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type SampleBucket struct {
	// DurabilityLevel is the minimum level at which all writes to the sample bucket must occur.
	// Default: "none"
	// Enum: "none" "majority" "majorityAndPersistActive" "persistToMajority"
	//
	// The options for Durability level are as follows, according to the sample bucket type.
	//
	// For a Couchbase sample bucket:
	// None
	// Replicate to Majority
	// Majority and Persist to Active
	// Persist to Majority
	DurabilityLevel types.String `tfsdk:"durability_level"`

	// Stats has the sample bucket stats that are related to memory and disk consumption.
	// itemCount: Number of documents in the sample bucket.
	// opsPerSecond: Number of operations per second.
	// diskUsedInMib: The amount of disk used (in MiB).
	// memoryUsedInMib: The amount of memory used (in MiB).
	Stats types.Object `tfsdk:"stats"`

	// Type defines the type of the sample bucket.
	// Default: "couchbase"
	//
	// This field for sample buckets is always the default and cannot be changed.
	// The options may also be referred to as Memory and Disk (Couchbase), Memory Only (Ephemeral) in the Couchbase documentation.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	Type types.String `tfsdk:"type"`

	// StorageBackend defines the storage engine that is used by the sample bucket.
	// Default: "couchstore"
	//
	// This field for sample buckets is always the default and cannot be changed.
	// This field cannot be changed later.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/storage-engines.html
	StorageBackend types.String `tfsdk:"storage_backend"`

	// ClusterId is the ID of the cluster for which the database credential needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// BucketConflictResolution is the means by which conflicts are resolved during replication.
	// Default: "seqno"
	//
	// This field for sample buckets is always the default and cannot be changed.
	// This field may be referred to as "conflict resolution" in the Couchbase documentation.
	// seqno may be referred to as "sequence number".
	// This field cannot be changed later.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/xdcr/xdcr.html#conflict-resolution
	BucketConflictResolution types.String `tfsdk:"bucket_conflict_resolution"`

	// Name is the name of the sample bucket.
	// Enum: "travel-sample", "beer-sample", "gamesim-sample"
	Name types.String `tfsdk:"name"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	ProjectId types.String `tfsdk:"project_id"`

	// Id is the id of the created sample bucket.
	Id types.String `tfsdk:"id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	OrganizationId types.String `tfsdk:"organization_id"`

	// EvictionPolicy is the policy which Capella adopts to prevent data loss due to memory exhaustion.
	// This may be also known as Ejection Policy in the Couchbase documentation.
	//
	// For Couchbase sample bucket, Eviction Policy is fullEviction by default and cannot be changed
	// To learn more, see https://docs.couchbase.com/server/current/rest-api/rest-bucket-create.html#evictionpolicy
	EvictionPolicy types.String `tfsdk:"eviction_policy"`

	// MemoryAllocationInMB is the amount of memory to allocate for the sample bucket memory in MiB.
	// This is the maximum limit is dependent on the allocation of the KV service. For example, 80% of the allocation.
	// Default: 200
	//
	// For Couchbase sample buckets, the default and minimum memory allocation is different. The minimum allocation is 100MiB
	MemoryAllocationInMB types.Int64 `tfsdk:"memory_allocation_in_mb"`

	// TimeToLiveInSeconds specifies the time to live (TTL) value in seconds.
	// This is the maximum time to live for items in the sample bucket.
	// Default is 0, that means TTL is disabled. This is a non-negative value.
	TimeToLiveInSeconds types.Int64 `tfsdk:"time_to_live_in_seconds"`

	// Replicas is the number of replicas for the sample bucket.
	// Default: 1
	// Enum: 1 2 3
	Replicas types.Int64 `tfsdk:"replicas"`

	// Flush determines whether flushing is enabled on the sample bucket.
	// Enable Flush to delete all items in this sample bucket at the earliest opportunity.
	// Disable Flush to avoid inadvertent data loss.
	// Default: false
	Flush types.Bool `tfsdk:"flush"`
}

// SampleBuckets defines attributes for the LIST buckets response received from V4 Capella Public API.
type SampleBuckets struct {
	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Data It contains the list of resources.
	Data []SampleBucket `tfsdk:"data"`
}

// Validate will split the IDs by a delimiter i.e. comma , in case a terraform import CLI is invoked.
// The format of the terraform import CLI would include the IDs as follows -
// `terraform import couchbase-capella_sample_bucket.new_sample_bucket id=<uuid>,cluster_id=<uuid>,project_id=<uuid>,organization_id=<uuid>`.
func (b SampleBucket) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: b.OrganizationId,
		ProjectId:      b.ProjectId,
		ClusterId:      b.ClusterId,
		Id:             b.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

// Validate is used to verify that all the fields in the datasource
// have been populated.
func (b SampleBuckets) Validate() (clusterId, projectId, organizationId string, err error) {
	if b.OrganizationId.IsNull() {
		return "", "", "", errors.ErrOrganizationIdMissing
	}
	if b.ProjectId.IsNull() {
		return "", "", "", errors.ErrProjectIdMissing
	}
	if b.ClusterId.IsNull() {
		return "", "", "", errors.ErrClusterIdMissing
	}
	return b.ClusterId.ValueString(), b.ProjectId.ValueString(), b.OrganizationId.ValueString(), nil
}

// NewStats creates a new Stats data object.
func NewStats(stats samplebucketapi.Stats) Stats {
	return Stats{
		ItemCount:       types.Int64Value(stats.ItemCount),
		OpsPerSecond:    types.Int64Value(stats.OpsPerSecond),
		DiskUsedInMiB:   types.Int64Value(stats.DiskUsedInMib),
		MemoryUsedInMiB: types.Int64Value(stats.MemoryUsedInMib),
	}
}
