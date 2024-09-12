package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

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

	// IndexName is the name of the index.
	IndexName types.String `tfsdk:"index_name"`

	// IsPrimary indicates if it's a primary index.
	IsPrimary types.Bool `tfsdk:"is_primary"`

	// IndexKeys is a list of index keys.
	IndexKeys types.Set `tfsdk:"index_keys"`

	// PartitionBy is the partition by clause.
	PartitionBy types.String `tfsdk:"partition_by"`

	// Where is the where clause.
	Where types.String `tfsdk:"where"`

	// With represents the WITH clause of an index.
	With *WithOptions `tfsdk:"with"`

	BuildIndexes types.Set `tfsdk:"build_indexes"`
}

// WithOptions represents the attributes of the WITH clause.
type WithOptions struct {
	// DeferBuild indicates if this is a deferred index build.
	DeferBuild types.Bool `tfsdk:"defer_build"`

	// NumReplica is the number of replicas for the index.
	NumReplica types.Int64 `tfsdk:"num_replica"`

	// NumPartitions is the number of partitions for a partitioned index.
	NumPartitions types.Int64 `tfsdk:"num_partitions"`
}

// Validate will split the IDs by comma, if a terraform import CLI is invoked.
func (g *GsiDefinition) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: g.OrganizationId,
		ProjectId:      g.ProjectId,
		ClusterId:      g.ClusterId,
		BucketName:     g.BucketName,
		ScopeName:      g.ScopeName,
		CollectionName: g.CollectionName,
		IndexName:      g.IndexName,
	}

	IDs, err := validateSchemaState(state, IndexName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
