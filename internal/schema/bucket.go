package schema

import (
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Bucket struct {
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

	// Stats has the bucket stats that are related to memory and disk consumption.
	// itemCount: Number of documents in the bucket.
	// opsPerSecond: Number of operations per second.
	// diskUsedInMib: The amount of disk used (in MiB).
	// memoryUsedInMib: The amount of memory used (in MiB).
	Stats types.Object `tfsdk:"stats"`

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

	// ClusterId is the ID of the cluster for which the database credential needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// BucketConflictResolution is the means by which conflicts are resolved during replication.
	// Default: "seqno"
	// Enum: "seqno" "lww"
	// This field may be referred to as "conflict resolution" in the Couchbase documentation.
	// seqno and lww may be referred to as "sequence number" and "timestamp" respectively.
	// This field cannot be changed later.
	// To learn more, see https://docs.couchbase.com/cloud/clusters/xdcr/xdcr.html#conflict-resolution
	BucketConflictResolution types.String `tfsdk:"bucket_conflict_resolution"`

	// Name is the name of the bucket.
	Name types.String `tfsdk:"name"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	ProjectId types.String `tfsdk:"project_id"`

	// Id is the id of the created bucket.
	Id types.String `tfsdk:"id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	// The database credential will be created for the cluster.
	OrganizationId types.String `tfsdk:"organization_id"`

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

	// TimeToLiveInSeconds specifies the time to live (TTL) value in seconds.
	// This is the maximum time to live for items in the bucket.
	// Default is 0, that means TTL is disabled. This is a non-negative value.
	TimeToLiveInSeconds types.Int64 `tfsdk:"time_to_live_in_seconds"`

	// Replicas is the number of replicas for the bucket.
	// Default: 1
	// Enum: 1 2 3
	Replicas types.Int64 `tfsdk:"replicas"`

	// Flush determines whether flushing is enabled on the bucket.
	// Enable Flush to delete all items in this bucket at the earliest opportunity.
	// Disable Flush to avoid inadvertent data loss.
	// Default: false
	Flush types.Bool `tfsdk:"flush"`
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

func (s Stats) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"item_count":         types.Int64Type,
		"ops_per_second":     types.Int64Type,
		"disk_used_in_mib":   types.Int64Type,
		"memory_used_in_mib": types.Int64Type,
	}
}

// Buckets defines attributes for the LIST buckets response received from V4 Capella Public API.
type Buckets struct {
	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Data It contains the list of resources.
	Data []OneBucket `tfsdk:"data"`
}

type OneBucket struct {
	Stats                    *Stats       `tfsdk:"stats"`
	DurabilityLevel          types.String `tfsdk:"durability_level"`
	Name                     types.String `tfsdk:"name"`
	StorageBackend           types.String `tfsdk:"storage_backend"`
	ClusterId                types.String `tfsdk:"cluster_id"`
	BucketConflictResolution types.String `tfsdk:"bucket_conflict_resolution"`
	Id                       types.String `tfsdk:"id"`
	ProjectId                types.String `tfsdk:"project_id"`
	OrganizationId           types.String `tfsdk:"organization_id"`
	Type                     types.String `tfsdk:"type"`
	EvictionPolicy           types.String `tfsdk:"eviction_policy"`
	TimeToLiveInSeconds      types.Int64  `tfsdk:"time_to_live_in_seconds"`
	Replicas                 types.Int64  `tfsdk:"replicas"`
	MemoryAllocationInMB     types.Int64  `tfsdk:"memory_allocation_in_mb"`
	Flush                    types.Bool   `tfsdk:"flush"`
}

// Validate will split the IDs by a delimiter i.e. comma , in case a terraform import CLI is invoked.
// The format of the terraform import CLI would include the IDs as follows -
// `terraform import capella_bucket.new_bucket id=<uuid>,cluster_id=<uuid>,project_id=<uuid>,organization_id=<uuid>`.
func (b Bucket) Validate() (map[Attr]string, error) {
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
func (b Buckets) Validate() (clusterId, projectId, organizationId string, err error) {
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
