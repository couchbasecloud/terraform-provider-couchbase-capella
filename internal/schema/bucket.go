package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

//const (
//	OrganizationId = "organizationId"
//	ProjectId      = "projectId"
//	ClusterId      = "clusterId"
//)

type Type string

type StorageBackend string

type ConflictResolution string

type DurabilityLevel string

type BucketEviction string

//type Bucket struct {
//	Name types.String `tfsdk:"name"`
//
//	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
//	// The database credential will be created for the cluster.
//	OrganizationId types.String `tfsdk:"organization_id"`
//
//	// ProjectId is the ID of the project to which the Capella cluster belongs.
//	// The database credential will be created for the cluster.
//	ProjectId types.String `tfsdk:"project_id"`
//
//	// ClusterId is the ID of the cluster for which the database credential needs to be created.
//	ClusterId types.String `tfsdk:"cluster_id"`
//
//	Type Type `tfsdk:"type"`
//
//	StorageBackend StorageBackend `tfsdk:"storage_backend"`
//
//	MemoryAllocationInMb int `tfsdk:"memory_allocationinmb"`
//
//	BucketConflictResolution ConflictResolution `tfsdk:"conflict_resolution"`
//
//	DurabilityLevel DurabilityLevel `tfsdk:"durability_level"`
//
//	Replicas int `tfsdk:"replicas"`
//
//	Flush bool `tfsdk:"flush"`
//
//	TimeToLiveInSeconds int `tfsdk:"ttl"`
//
//	EvictionPolicy BucketEviction `tfsdk:"eviction_policy"`
//}

type Bucket struct {
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
}
