package schema

import (
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Bucket struct {
	DurabilityLevel          types.String `tfsdk:"durability_level"`
	Stats                    types.Object `tfsdk:"stats"`
	Type                     types.String `tfsdk:"type"`
	StorageBackend           types.String `tfsdk:"storage_backend"`
	ClusterId                types.String `tfsdk:"cluster_id"`
	BucketConflictResolution types.String `tfsdk:"bucket_conflict_resolution"`
	Name                     types.String `tfsdk:"name"`
	ProjectId                types.String `tfsdk:"project_id"`
	Id                       types.String `tfsdk:"id"`
	OrganizationId           types.String `tfsdk:"organization_id"`
	EvictionPolicy           types.String `tfsdk:"eviction_policy"`
	MemoryAllocationInMB     types.Int64  `tfsdk:"memory_allocation_in_mb"`
	TimeToLiveInSeconds      types.Int64  `tfsdk:"time_to_live_in_seconds"`
	Replicas                 types.Int64  `tfsdk:"replicas"`
	Flush                    types.Bool   `tfsdk:"flush"`
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
