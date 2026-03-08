package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var replicationBuilder = capellaschema.NewSchemaBuilder("replication", "GetReplicationResponse")

// ReplicationSchema returns the schema for the replication data source.
func ReplicationSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	// Required hierarchical IDs
	capellaschema.AddAttr(attrs, "organization_id", replicationBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", replicationBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", replicationBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "replication_id", replicationBuilder, requiredStringWithValidator())

	// Computed attributes
	capellaschema.AddAttr(attrs, "id", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "status", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "changes_left", replicationBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "error", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "direction", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "priority", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "network_usage_limit", replicationBuilder, computedInt64())

	// Nested computed attributes
	// Source
	sourceAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(sourceAttrs, "project_id", replicationBuilder, computedString(), "ReplicationSource", "Source")
	capellaschema.AddAttr(sourceAttrs, "project_name", replicationBuilder, computedString(), "ReplicationSource", "Source")
	capellaschema.AddAttr(sourceAttrs, "cluster_id", replicationBuilder, computedString(), "ReplicationSource", "Source")
	capellaschema.AddAttr(sourceAttrs, "cluster_name", replicationBuilder, computedString(), "ReplicationSource", "Source")
	capellaschema.AddAttr(sourceAttrs, "bucket_id", replicationBuilder, computedString(), "ReplicationSource", "Source")
	capellaschema.AddAttr(sourceAttrs, "bucket_name", replicationBuilder, computedString(), "ReplicationSource", "Source")
	capellaschema.AddAttr(sourceAttrs, "bucket_conflict_resolution", replicationBuilder, computedString(), "ReplicationSource", "Source")
	capellaschema.AddAttr(sourceAttrs, "type", replicationBuilder, computedString(), "ReplicationSourceType")

	// Source scopes
	scopeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(scopeAttrs, "name", replicationBuilder, computedString(), "Scopes")
	capellaschema.AddAttr(scopeAttrs, "collections", replicationBuilder, computedStringSet(), "Collections")
	capellaschema.AddAttr(sourceAttrs, "scopes", replicationBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: scopeAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "source", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: sourceAttrs,
	})

	// Target
	targetAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(targetAttrs, "project_id", replicationBuilder, computedString(), "ReplicationTarget", "Target")
	capellaschema.AddAttr(targetAttrs, "project_name", replicationBuilder, computedString(), "ReplicationTarget", "Target")
	capellaschema.AddAttr(targetAttrs, "cluster_id", replicationBuilder, computedString(), "ReplicationTarget", "Target")
	capellaschema.AddAttr(targetAttrs, "cluster_name", replicationBuilder, computedString(), "ReplicationTarget", "Target")
	capellaschema.AddAttr(targetAttrs, "bucket_id", replicationBuilder, computedString(), "ReplicationTarget", "Target")
	capellaschema.AddAttr(targetAttrs, "bucket_name", replicationBuilder, computedString(), "ReplicationTarget", "Target")
	capellaschema.AddAttr(targetAttrs, "bucket_conflict_resolution", replicationBuilder, computedString(), "ReplicationTarget", "Target")
	capellaschema.AddAttr(targetAttrs, "type", replicationBuilder, computedString(), "ReplicationTargetType")

	// Target scopes (same structure as source scopes)
	capellaschema.AddAttr(targetAttrs, "scopes", replicationBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: scopeAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "target", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: targetAttrs,
	})

	// Mappings
	collectionMappingAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(collectionMappingAttrs, "source_collection", replicationBuilder, computedString(), "Mappings")
	capellaschema.AddAttr(collectionMappingAttrs, "target_collection", replicationBuilder, computedString(), "Mappings")

	mappingAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(mappingAttrs, "source_scope", replicationBuilder, computedString(), "Mappings")
	capellaschema.AddAttr(mappingAttrs, "target_scope", replicationBuilder, computedString(), "Mappings")
	capellaschema.AddAttr(mappingAttrs, "collections", replicationBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: collectionMappingAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "mappings", replicationBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: mappingAttrs,
		},
	})

	// Filter
	docExcludeOptsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(docExcludeOptsAttrs, "binary", replicationBuilder, computedBool(), "GetFilter")
	capellaschema.AddAttr(docExcludeOptsAttrs, "deletion", replicationBuilder, computedBool(), "GetFilter")
	capellaschema.AddAttr(docExcludeOptsAttrs, "expiration", replicationBuilder, computedBool(), "GetFilter")
	capellaschema.AddAttr(docExcludeOptsAttrs, "ttl", replicationBuilder, computedBool(), "GetFilter")

	expressionsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(expressionsAttrs, "regex", replicationBuilder, computedString(), "GetFilter")

	filterAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(filterAttrs, "document_exclude_options", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: docExcludeOptsAttrs,
	}, "GetFilter")
	capellaschema.AddAttr(filterAttrs, "expressions", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: expressionsAttrs,
	}, "GetFilter")

	capellaschema.AddAttr(attrs, "filter", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: filterAttrs,
	})

	// Audit (replication uses simpler audit with only created_at and created_by)
	capellaschema.AddAttr(attrs, "audit", replicationBuilder, computedReplicationAudit())

	return schema.Schema{
		MarkdownDescription: "Retrieves the details of a specific replication.",
		Attributes:          attrs,
	}
}
