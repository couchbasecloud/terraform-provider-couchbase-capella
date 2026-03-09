package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// replicationBuilder is the schema builder for the replication data source.
var replicationBuilder = capellaschema.NewSchemaBuilder("replication", "GetReplicationResponse")

// ReplicationSchema returns the schema for the singular replication data source.
func ReplicationSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	// Required hierarchical IDs
	capellaschema.AddAttr(attrs, "organization_id", replicationBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", replicationBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", replicationBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "replication_id", replicationBuilder, requiredStringWithValidator())

	// Project reference attributes
	projectAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(projectAttrs, "id", replicationBuilder, computedString())
	capellaschema.AddAttr(projectAttrs, "name", replicationBuilder, computedString())

	// Cluster reference attributes
	clusterAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(clusterAttrs, "id", replicationBuilder, computedString())
	capellaschema.AddAttr(clusterAttrs, "name", replicationBuilder, computedString())

	// Bucket reference attributes
	bucketAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(bucketAttrs, "id", replicationBuilder, computedString())
	capellaschema.AddAttr(bucketAttrs, "name", replicationBuilder, computedString())
	capellaschema.AddAttr(bucketAttrs, "conflict_resolution_type", replicationBuilder, computedString())

	// Scope attributes for replication
	scopeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(scopeAttrs, "name", replicationBuilder, computedString())
	capellaschema.AddAttr(scopeAttrs, "collections", replicationBuilder, &schema.ListAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})

	// Source attributes
	sourceAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(sourceAttrs, "project", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: projectAttrs,
	})
	capellaschema.AddAttr(sourceAttrs, "cluster", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: clusterAttrs,
	})
	capellaschema.AddAttr(sourceAttrs, "bucket", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: bucketAttrs,
	})
	capellaschema.AddAttr(sourceAttrs, "scopes", replicationBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: scopeAttrs,
		},
	})
	capellaschema.AddAttr(sourceAttrs, "type", replicationBuilder, computedString())

	// Target attributes
	targetAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(targetAttrs, "project", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: projectAttrs,
	})
	capellaschema.AddAttr(targetAttrs, "cluster", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: clusterAttrs,
	})
	capellaschema.AddAttr(targetAttrs, "bucket", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: bucketAttrs,
	})
	capellaschema.AddAttr(targetAttrs, "scopes", replicationBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: scopeAttrs,
		},
	})
	capellaschema.AddAttr(targetAttrs, "type", replicationBuilder, computedString())

	// Collection mapping attributes
	collectionMappingAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(collectionMappingAttrs, "source_collection", replicationBuilder, computedString())
	capellaschema.AddAttr(collectionMappingAttrs, "target_collection", replicationBuilder, computedString())

	// Mapping attributes
	mappingAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(mappingAttrs, "source_scope", replicationBuilder, computedString())
	capellaschema.AddAttr(mappingAttrs, "target_scope", replicationBuilder, computedString())
	capellaschema.AddAttr(mappingAttrs, "collections", replicationBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: collectionMappingAttrs,
		},
	})

	// Document exclude options attributes
	docExcludeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(docExcludeAttrs, "deletion", replicationBuilder, computedBool())
	capellaschema.AddAttr(docExcludeAttrs, "expiration", replicationBuilder, computedBool())
	capellaschema.AddAttr(docExcludeAttrs, "ttl", replicationBuilder, computedBool())
	capellaschema.AddAttr(docExcludeAttrs, "binary", replicationBuilder, computedBool())

	// Filter expressions attributes
	expressionsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(expressionsAttrs, "reg_ex", replicationBuilder, computedString())

	// Filter attributes
	filterAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(filterAttrs, "document_exclude_options", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: docExcludeAttrs,
	})
	capellaschema.AddAttr(filterAttrs, "expressions", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: expressionsAttrs,
	})

	// Audit attributes
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_by", replicationBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "created_at", replicationBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "modified_by", replicationBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "modified_at", replicationBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "version", replicationBuilder, computedInt64())

	// Computed attributes at top level
	capellaschema.AddAttr(attrs, "id", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "status", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "changes_left", replicationBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "error", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "direction", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "priority", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "network_usage_limit", replicationBuilder, computedInt64())

	// Nested computed attributes
	capellaschema.AddAttr(attrs, "source", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: sourceAttrs,
	})
	capellaschema.AddAttr(attrs, "target", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: targetAttrs,
	})
	capellaschema.AddAttr(attrs, "mappings", replicationBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: mappingAttrs,
		},
	})
	capellaschema.AddAttr(attrs, "filter", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: filterAttrs,
	})
	capellaschema.AddAttr(attrs, "audit", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "The data source retrieves details of a specific Couchbase Capella replication.",
		Attributes:          attrs,
	}
}
