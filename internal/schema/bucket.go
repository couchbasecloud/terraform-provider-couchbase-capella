package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

type Bucket struct {
	// Id is the id of the created bucket.
	ID types.String `tfsdk:"id"`

	// Name is the name of the bucket.
	Name types.String `tfsdk:"name"`

	// Type defines the type of the bucket.
	// Default: "couchbase"
	// Enum: "couchbase" "ephemeral"
	// If selected Ephemeral, it is not eligible for imports or App Endpoints creation. This field cannot be changed later.
	// The options may also be referred to as Memory and Disk (Couchbase), Memory Only (Ephemeral) in the Couchbase documentation.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/manage-buckets.html#add-bucket
	Type types.String `tfsdk:"type"`

	// StorageBackend defines the storage engine that is used by the bucket.
	// Default: "couchstore"
	// Enum: "couchstore" "magma"
	//
	// Ephemeral buckets do not support StorageBackend, hence not applicable for Ephemeral buckets and throws an error if this field is added.
	// This field is only applicable for a Couchbase bucket. The default value mentioned (Couchstore) is for Couchbase bucket.
	// This field cannot be changed later.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/data-service/storage-engines.html
	StorageBackend types.String `tfsdk:"storage_backend"`

	// MemoryAllocationInMB is the amount of memory to allocate for the bucket memory in MiB.
	// This is the maximum limit is dependent on the allocation of the KV service. For example, 80% of the allocation.
	// Default: 100
	//
	// The default value (100MiB) mentioned is for Couchbase type buckets with Couchstore as the Storage Backend.
	//
	// For Couchbase buckets, the default and minimum memory allocation changes according to the Storage Backend type as follows:
	// For Couchstore, the default and minimum memory allocation is 100 MiB.
	// For Magma, the default and minimum memory allocation is 1024 MiB.
	// For Ephemeral buckets, the default and minimum memory allocation is 100 MiB.
	MemoryAllocationInMB types.Int64 `tfsdk:"memory_allocation_in_mb"`

	// BucketConflictResolution is the means by which conflicts are resolved during replication.
	// Default: "seqno"
	// Enum: "seqno" "lww"
	// This field may be referred to as "conflict resolution" in the Couchbase documentation.
	// seqno and lww may be referred to as "sequence number" and "timestamp" respectively.
	// This field cannot be changed later.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/xdcr/xdcr.html#conflict-resolution
	BucketConflictResolution types.String `tfsdk:"bucket_conflict_resolution"`

	// DurabilityLevel is the minimum level at which all writes to the bucket must occur.
	// Default: "none"
	// Enum: "none" "majority" "majorityAndPersistActive" "persistToMajority"
	//
	// The options for Durability level are as follows, according to the bucket type.
	//
	// For a Couchbase bucket:
	// None
	// Replicate to Majority
	// Majority and Persist to Active
	// Persist to Majority
	//
	//For an Ephemeral bucket:
	// None
	// Replicate to Majority
	DurabilityLevel types.String `tfsdk:"durability_level"`

	// Replicas is the number of replicas for the bucket.
	// Default: 1
	// Enum: 1 2 3
	Replicas types.Int64 `tfsdk:"replicas"`

	// Flush determines whether flushing is enabled on the bucket.
	// Enable Flush to delete all items in this bucket at the earliest opportunity.
	// Disable Flush to avoid inadvertent data loss.
	// Default: false
	Flush types.Bool `tfsdk:"flush"`

	// TimeToLiveInSeconds specifies the time to live (TTL) value in seconds.
	// This is the maximum time to live for items in the bucket.
	// Default is 0, that means TTL is disabled. This is a non-negative value.
	TimeToLiveInSeconds types.Int64 `tfsdk:"time_to_live_in_seconds"`

	// EvictionPolicy is the policy which Capella adopts to prevent data loss due to memory exhaustion.
	// This may be also known as Ejection Policy in the Couchbase documentation.
	//
	// For Couchbase bucket, Eviction Policy is fullEviction by default.
	// For Ephemeral buckets, Eviction Policy is a required field, and should be one of the following:
	// noEviction
	// nruEviction
	// Default: "fullEviction"
	// Enum: "fullEviction" "noEviction" "nruEviction"
	// To learn more, see https://docs.couchbase.com/server/current/rest-api/rest-bucket-create.html#evictionpolicy
	EvictionPolicy types.String `tfsdk:"eviction_policy"`

	// Stats has the bucket stats that are related to memory and disk consumption.
	// itemCount: Number of documents in the bucket.
	// opsPerSecond: Number of operations per second.
	// diskUsedInMib: The amount of disk used (in MiB).
	// memoryUsedInMib: The amount of memory used (in MiB).
	Stats *Stats `tfsdk:"stats"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the database credential needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Audit contains all audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`
}

// Stats has the bucket stats that are related to memory and disk consumption.
type Stats struct {
	// ItemCount: Number of documents in the bucket.
	ItemCount types.Int64 `tfsdk:"item_count"`
	// OpsPerSecond: Number of operations per second.
	OpsPerSecond types.Int64 `tfsdk:"ops_per_second"`
	// DiskUsedInMib: The amount of disk used (in MiB).
	DiskUsedInMiB types.Int64 `tfsdk:"disk_used_in_mib"`
	// MemoryUsedInMib: The amount of memory used (in MiB).
	MemoryUsedInMiB types.Int64 `tfsdk:"memory_used_in_mib"`
}
