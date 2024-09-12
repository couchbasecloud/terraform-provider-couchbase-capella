package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// GsiDefinition represents the primary or secondary index.
type GsiDefinition struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the scope needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// BucketName is the bucket for the index.
	BucketName types.String `tfsdk:"bucket_name"`

	// ScopeName is the scope for the index.
	ScopeName types.String `tfsdk:"scope_name"`

	// CollectionName is the collection for the index.
	CollectionName types.String `tfsdk:"collection_name"`

	// DefId is the index definition id.
	DefId types.String `tfsdk:"def_id"`

	// IsPrimary indicates if it's a primary index.
	IsPrimary types.Bool `tfsdk:"is_primary"`

	// IndexKeys is a list of index keys.
	IndexKeys types.Set `tfsdk:"index_keys"`

	// Where is the where clause.
	Where types.String `tfsdk:"where"`

	// With represents the WITH clause of an index.
	With IndexOptions
}

// IndexOptions represents the attributes of the WITH clause.
type IndexOptions struct {
	// DeferBuild indicates if this is a deferred index build.
	DeferBuild types.Bool `tfsdk:"defer_build"`

	// NumReplica is the number of replicas for the index.
	NumReplica types.Int64 `tfsdk:"num_replica"`

	NumPartitions types.Int64 `tfsdk:"num_partitions"`
}
