package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var replicationBuilder = capellaschema.NewSchemaBuilder("replication")

// ReplicationSchema returns the schema for the replication data source.
func ReplicationSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", replicationBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", replicationBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", replicationBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "id", replicationBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "source_bucket", replicationBuilder, computedString())

	targetAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(targetAttrs, "cluster", replicationBuilder, computedString())
	capellaschema.AddAttr(targetAttrs, "bucket", replicationBuilder, computedString())
	capellaschema.AddAttr(targetAttrs, "type", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "target", replicationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: targetAttrs,
	})

	capellaschema.AddAttr(attrs, "direction", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "priority", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "network_usage_limit", replicationBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "mappings", replicationBuilder, computedReplicationMappingsAttribute())
	capellaschema.AddAttr(attrs, "filter", replicationBuilder, computedReplicationFilterAttribute())

	mobileSyncMode := computedString()
	mobileSyncMode.MarkdownDescription = capellaschema.ReplicationDataMobileSyncModeDescription
	capellaschema.AddAttr(attrs, "mobile_sync_mode", replicationBuilder, mobileSyncMode)

	capellaschema.AddAttr(attrs, "conflict_logging", replicationBuilder, computedConflictLoggingAttribute())

	capellaschema.AddAttr(attrs, "status", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "changes_left", replicationBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "error", replicationBuilder, computedString())
	capellaschema.AddAttr(attrs, "source_details", replicationBuilder, computedReplicationDetailsAttribute(false))
	capellaschema.AddAttr(attrs, "target_details", replicationBuilder, computedReplicationDetailsAttribute(true))
	capellaschema.AddAttr(attrs, "audit", replicationBuilder, computedReplicationAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve an XDCR replication.\n\nThe read fails while the source " +
			"cluster is not healthy, for example during a rebalance or soon after a replication is created. Run the " +
			"command again when the cluster is healthy.",
		Attributes: attrs,
	}
}

func computedReplicationMappingsAttribute() *schema.SetNestedAttribute {
	collectionAttrs := make(map[string]schema.Attribute)
	sourceCollection := computedString()
	sourceCollection.MarkdownDescription = capellaschema.ReplicationSourceCollectionDescription
	capellaschema.AddAttr(collectionAttrs, "source_collection", replicationBuilder, sourceCollection)
	targetCollection := computedString()
	targetCollection.MarkdownDescription = capellaschema.ReplicationTargetCollectionDescription
	capellaschema.AddAttr(collectionAttrs, "target_collection", replicationBuilder, targetCollection)

	mappingAttrs := make(map[string]schema.Attribute)
	sourceScope := computedString()
	sourceScope.MarkdownDescription = capellaschema.ReplicationSourceScopeDescription
	capellaschema.AddAttr(mappingAttrs, "source_scope", replicationBuilder, sourceScope)
	targetScope := computedString()
	targetScope.MarkdownDescription = capellaschema.ReplicationTargetScopeDescription
	capellaschema.AddAttr(mappingAttrs, "target_scope", replicationBuilder, targetScope)
	collections := &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: collectionAttrs,
		},
	}
	collections.MarkdownDescription = capellaschema.ReplicationCollectionsDescription
	capellaschema.AddAttr(mappingAttrs, "collections", replicationBuilder, collections)

	return &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: mappingAttrs,
		},
	}
}

func computedReplicationFilterAttribute() *schema.SingleNestedAttribute {
	excludeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(excludeAttrs, "deletion", replicationBuilder, computedBool())
	capellaschema.AddAttr(excludeAttrs, "expiration", replicationBuilder, computedBool())
	capellaschema.AddAttr(excludeAttrs, "ttl", replicationBuilder, computedBool())
	capellaschema.AddAttr(excludeAttrs, "binary", replicationBuilder, computedBool())

	expressionAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(expressionAttrs, "reg_ex", replicationBuilder, computedString())

	filterAttrs := make(map[string]schema.Attribute)
	documentExcludeOptions := &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: excludeAttrs,
	}
	documentExcludeOptions.MarkdownDescription = capellaschema.ReplicationDocumentExcludeOptionsDescription
	capellaschema.AddAttr(filterAttrs, "document_exclude_options", replicationBuilder, documentExcludeOptions)

	expressions := &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: expressionAttrs,
	}
	expressions.MarkdownDescription = capellaschema.ReplicationExpressionsDescription
	capellaschema.AddAttr(filterAttrs, "expressions", replicationBuilder, expressions)

	return &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: filterAttrs,
	}
}

// computedConflictLoggingAttribute returns the conflict_logging block.
func computedConflictLoggingAttribute() *schema.SingleNestedAttribute {
	ruleBucket := computedString()
	ruleBucket.MarkdownDescription = capellaschema.ReplicationDataRuleBucketDescription
	ruleCollection := computedString()
	ruleCollection.MarkdownDescription = capellaschema.ReplicationDataRuleCollectionDescription
	ruleAttrs := map[string]schema.Attribute{
		"bucket":     ruleBucket,
		"collection": ruleCollection,
	}

	bucket := computedString()
	bucket.MarkdownDescription = capellaschema.ReplicationConflictLoggingBucketDescription
	collection := computedString()
	collection.MarkdownDescription = capellaschema.ReplicationConflictLoggingCollectionDescription
	conflictLoggingAttrs := map[string]schema.Attribute{
		"bucket":     bucket,
		"collection": collection,
	}
	disabled := computedBool()
	disabled.MarkdownDescription = capellaschema.ReplicationDataDisabledDescription
	capellaschema.AddAttr(conflictLoggingAttrs, "disabled", replicationBuilder, disabled)

	loggingRules := &schema.MapNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: ruleAttrs,
		},
	}
	loggingRules.MarkdownDescription = capellaschema.ReplicationLoggingRulesDescription
	capellaschema.AddAttr(conflictLoggingAttrs, "logging_rules", replicationBuilder, loggingRules)

	conflictLogging := &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: conflictLoggingAttrs,
	}
	conflictLogging.MarkdownDescription = capellaschema.ReplicationDataConflictLoggingDescription
	return conflictLogging
}

// computedReplicationDetailsAttribute returns the source details, or the target details when isTarget
// is set. Only the target has a type.
func computedReplicationDetailsAttribute(isTarget bool) *schema.SingleNestedAttribute {
	detailsBuilder := capellaschema.NewSchemaBuilder("source_details", "replicationSource")
	detailsDescription := capellaschema.ReplicationSourceDetailsDescription
	if isTarget {
		detailsBuilder = capellaschema.NewSchemaBuilder("target_details", "replicationTarget")
		detailsDescription = capellaschema.ReplicationTargetDetailsDescription
	}

	referenceAttrs := func() map[string]schema.Attribute {
		attrs := make(map[string]schema.Attribute)
		id := computedString()
		id.MarkdownDescription = capellaschema.ReplicationDetailsIdDescription
		capellaschema.AddAttr(attrs, "id", detailsBuilder, id)
		name := computedString()
		name.MarkdownDescription = capellaschema.ReplicationDetailsNameDescription
		capellaschema.AddAttr(attrs, "name", detailsBuilder, name)
		return attrs
	}

	bucketAttrs := referenceAttrs()
	conflictResolutionType := computedString()
	conflictResolutionType.MarkdownDescription = capellaschema.ReplicationDetailsConflictResolutionTypeDescription
	capellaschema.AddAttr(bucketAttrs, "conflict_resolution_type", detailsBuilder, conflictResolutionType)

	detailsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(detailsAttrs, "project", detailsBuilder, &schema.SingleNestedAttribute{Computed: true, Attributes: referenceAttrs()})
	capellaschema.AddAttr(detailsAttrs, "cluster", detailsBuilder, &schema.SingleNestedAttribute{Computed: true, Attributes: referenceAttrs()})
	capellaschema.AddAttr(detailsAttrs, "bucket", detailsBuilder, &schema.SingleNestedAttribute{Computed: true, Attributes: bucketAttrs})
	if isTarget {
		capellaschema.AddAttr(detailsAttrs, "type", detailsBuilder, computedString())
	}

	details := &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: detailsAttrs,
	}
	details.MarkdownDescription = detailsDescription
	return details
}

// computedReplicationAuditAttribute returns the replication audit.
func computedReplicationAuditAttribute() *schema.SingleNestedAttribute {
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_by", replicationBuilder, computedString(), "ReplicationAuditData")
	capellaschema.AddAttr(auditAttrs, "created_at", replicationBuilder, computedString(), "ReplicationAuditData")

	return &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	}
}
