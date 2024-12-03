package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	IndexKeys types.List `tfsdk:"index_keys"`

	// PartitionBy is the partition by clause.
	PartitionBy types.List `tfsdk:"partition_by"`

	// Where is the where clause.
	Where types.String `tfsdk:"where"`

	// Status is the index state such as Ready, etc
	Status types.String `tfsdk:"status"`

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

	// NumPartition is the number of partitions for a partitioned index.
	NumPartition types.Int64 `tfsdk:"num_partition"`
}

type GsiDefinitions struct {
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

	Data []GsiData `tfsdk:"data"`
}

type GsiData struct {
	// IndexName is the name of the index.
	IndexName types.String `tfsdk:"index_name"`

	// Definition is the index definition.
	Definition types.String `tfsdk:"definition"`
}

type GsiBuildStatus struct {
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

	// Indexes are the list of indexes.
	Indexes types.Set `tfsdk:"indexes"`
}

func (g *GsiDefinition) GetAttributeValues() (map[Attr]string, error) {
	// handle terraform import.
	if g.OrganizationId.IsNull() {
		attrs, err := splitImportString(
			g.IndexName.ValueString(),
			[]Attr{OrganizationId,
				ProjectId,
				ClusterId,
				BucketName,
				ScopeName,
				CollectionName,
				IndexName},
		)
		if err != nil {
			return nil, err
		}
		return attrs, nil
	}

	// handle Read().
	attrs := map[Attr]string{
		OrganizationId: g.OrganizationId.ValueString(),
		ProjectId:      g.ProjectId.ValueString(),
		ClusterId:      g.ClusterId.ValueString(),
		BucketName:     g.BucketName.ValueString(),
		ScopeName:      g.ScopeName.ValueString(),
		CollectionName: g.CollectionName.ValueString(),
		IndexName:      "",
	}
	// if a primary index was created without a name,
	// indexer uses name #primary.
	if !g.IsPrimary.IsNull() && g.IndexName.IsNull() {
		attrs[IndexName] = "#primary"
	} else {
		attrs[IndexName] = g.IndexName.ValueString()
	}

	return attrs, nil
}
