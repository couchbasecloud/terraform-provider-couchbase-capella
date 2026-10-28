package resources

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var replicationBuilder = capellaschema.NewSchemaBuilder("replication")

// nonBlankPattern matches a value with at least one character that is not a space.
var nonBlankPattern = regexp.MustCompile(`\S`)

// ReplicationSchema defines the Terraform schema for the replication resource.
func ReplicationSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", replicationBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", replicationBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", replicationBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", replicationBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "source_bucket", replicationBuilder, requiredNonEmptyStringAttribute())

	targetAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(targetAttrs, "cluster", replicationBuilder, requiredStringAttributeNoReplace())
	capellaschema.AddAttr(targetAttrs, "bucket", replicationBuilder, requiredStringAttributeNoReplace())
	targetType := stringDefaultAttribute(capellaschema.ReplicationTargetTypeCapella, optional, computed)
	targetType.Validators = append(targetType.Validators,
		stringvalidator.OneOf(capellaschema.ReplicationTargetTypeCapella, capellaschema.ReplicationTargetTypeExternal))
	capellaschema.AddAttr(targetAttrs, "type", replicationBuilder, targetType)

	capellaschema.AddAttr(attrs, "target", replicationBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: targetAttrs,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
	})

	direction := stringDefaultAttributeReplaceOnChange(capellaschema.ReplicationDirectionOneWay, optional, computed)
	direction.Validators = append(direction.Validators,
		stringvalidator.OneOf(capellaschema.ReplicationDirectionOneWay, capellaschema.ReplicationDirectionTwoWay))
	capellaschema.AddAttr(attrs, "direction", replicationBuilder, direction)

	priority := stringDefaultAttribute(capellaschema.ReplicationPriorityHigh, optional, computed)
	priority.Validators = append(priority.Validators,
		stringvalidator.OneOf(capellaschema.ReplicationPriorityLow, capellaschema.ReplicationPriorityMedium, capellaschema.ReplicationPriorityHigh))
	capellaschema.AddAttr(attrs, "priority", replicationBuilder, priority)

	networkUsageLimit := int64DefaultAttribute(0, optional, computed)
	networkUsageLimit.Validators = append(networkUsageLimit.Validators,
		int64validator.Between(0, capellaschema.ReplicationNetworkUsageLimitMax))
	capellaschema.AddAttr(attrs, "network_usage_limit", replicationBuilder, networkUsageLimit)

	collectionAttrs := make(map[string]schema.Attribute)
	sourceCollection := requiredStringAttributeNoReplace()
	sourceCollection.MarkdownDescription = capellaschema.ReplicationSourceCollectionDescription
	capellaschema.AddAttr(collectionAttrs, "source_collection", replicationBuilder, sourceCollection)
	targetCollection := requiredStringAttributeNoReplace()
	targetCollection.MarkdownDescription = capellaschema.ReplicationTargetCollectionDescription
	capellaschema.AddAttr(collectionAttrs, "target_collection", replicationBuilder, targetCollection)

	mappingAttrs := make(map[string]schema.Attribute)
	sourceScope := requiredStringAttributeNoReplace()
	sourceScope.MarkdownDescription = capellaschema.ReplicationSourceScopeDescription
	capellaschema.AddAttr(mappingAttrs, "source_scope", replicationBuilder, sourceScope)
	targetScope := requiredStringAttributeNoReplace()
	targetScope.MarkdownDescription = capellaschema.ReplicationTargetScopeDescription
	capellaschema.AddAttr(mappingAttrs, "target_scope", replicationBuilder, targetScope)
	collections := &schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: collectionAttrs,
		},
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	}
	collections.MarkdownDescription = capellaschema.ReplicationCollectionsDescription
	capellaschema.AddAttr(mappingAttrs, "collections", replicationBuilder, collections)

	capellaschema.AddAttr(attrs, "mappings", replicationBuilder, &schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: mappingAttrs,
		},
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	})

	excludeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(excludeAttrs, "deletion", replicationBuilder, boolDefaultAttribute(false, optional, computed))
	capellaschema.AddAttr(excludeAttrs, "expiration", replicationBuilder, boolDefaultAttribute(false, optional, computed))
	capellaschema.AddAttr(excludeAttrs, "ttl", replicationBuilder, boolDefaultAttribute(false, optional, computed))
	capellaschema.AddAttr(excludeAttrs, "binary", replicationBuilder, boolDefaultAttribute(false, optional, computed))

	// reg_ex is required because GET drops the expressions object when the expression is empty.
	// Without it, a config like `expressions = { skip_restream = true }` reads back as null and
	// shows a diff on every plan. skip_restream also has no effect without an expression.
	expressionAttrs := make(map[string]schema.Attribute)
	regEx := requiredStringAttributeNoReplace()
	regEx.Validators = append(regEx.Validators, stringvalidator.RegexMatches(nonBlankPattern, "must not be only spaces"))
	capellaschema.AddAttr(expressionAttrs, "reg_ex", replicationBuilder, regEx)
	capellaschema.AddAttr(expressionAttrs, "skip_restream", replicationBuilder, boolDefaultAttribute(false, optional, computed))

	filterAttrs := make(map[string]schema.Attribute)
	documentExcludeOptions := &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: excludeAttrs,
	}
	documentExcludeOptions.MarkdownDescription = capellaschema.ReplicationDocumentExcludeOptionsDescription
	capellaschema.AddAttr(filterAttrs, "document_exclude_options", replicationBuilder, documentExcludeOptions)

	expressions := &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: expressionAttrs,
	}
	expressions.MarkdownDescription = capellaschema.ReplicationExpressionsDescription
	capellaschema.AddAttr(filterAttrs, "expressions", replicationBuilder, expressions)

	capellaschema.AddAttr(attrs, "filter", replicationBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: filterAttrs,
	})

	mobileSyncMode := stringDefaultAttribute(capellaschema.ReplicationMobileSyncModeOff, optional, computed)
	mobileSyncMode.Validators = append(mobileSyncMode.Validators,
		stringvalidator.OneOf(capellaschema.ReplicationMobileSyncModeOff, capellaschema.ReplicationMobileSyncModeActive))
	mobileSyncMode.MarkdownDescription = capellaschema.ReplicationMobileSyncModeDescription
	capellaschema.AddAttr(attrs, "mobile_sync_mode", replicationBuilder, mobileSyncMode)

	capellaschema.AddAttr(attrs, "conflict_logging", replicationBuilder, conflictLoggingAttribute())

	capellaschema.AddAttr(attrs, "reverse_replication_id", replicationBuilder,
		stringAttribute([]string{computed, useStateForUnknown}), "CreateReplicationJSONResponse")

	capellaschema.AddAttr(attrs, "status", replicationBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "changes_left", replicationBuilder, int64Attribute(computed))
	capellaschema.AddAttr(attrs, "error", replicationBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "source_details", replicationBuilder, computedReplicationDetailsAttribute(false))
	capellaschema.AddAttr(attrs, "target_details", replicationBuilder, computedReplicationDetailsAttribute(true))
	capellaschema.AddAttr(attrs, "audit", replicationBuilder, computedReplicationAuditAttribute())

	return schema.Schema{
		MarkdownDescription: capellaschema.ReplicationResourceDescription,
		Attributes:          attrs,
	}
}

// conflictLoggingAttribute returns the conflict_logging block.
func conflictLoggingAttribute() *schema.SingleNestedAttribute {
	ruleBucket := stringAttribute([]string{optional}, stringvalidator.LengthAtLeast(1))
	ruleBucket.MarkdownDescription = capellaschema.ReplicationLoggingRuleBucketDescription
	ruleCollection := stringAttribute([]string{optional}, stringvalidator.LengthAtLeast(1))
	ruleCollection.MarkdownDescription = capellaschema.ReplicationLoggingRuleCollectionDescription
	ruleAttrs := map[string]schema.Attribute{
		"bucket":     ruleBucket,
		"collection": ruleCollection,
	}

	bucket := stringAttribute([]string{required})
	bucket.MarkdownDescription = capellaschema.ReplicationConflictLoggingBucketDescription
	collection := stringAttribute([]string{required})
	collection.MarkdownDescription = capellaschema.ReplicationConflictLoggingCollectionDescription
	conflictLoggingAttrs := map[string]schema.Attribute{
		"bucket":     bucket,
		"collection": collection,
	}
	disabled := boolAttribute(required)
	disabled.MarkdownDescription = capellaschema.ReplicationConflictLoggingDisabledDescription
	capellaschema.AddAttr(conflictLoggingAttrs, "disabled", replicationBuilder, disabled)

	// GET drops an empty loggingRules map, so an empty map in the config would show a diff on every plan.
	loggingRules := &schema.MapNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: ruleAttrs,
		},
		Validators: []validator.Map{
			mapvalidator.SizeAtLeast(1),
		},
	}
	loggingRules.MarkdownDescription = capellaschema.ReplicationLoggingRulesDescription
	capellaschema.AddAttr(conflictLoggingAttrs, "logging_rules", replicationBuilder, loggingRules)

	conflictLogging := &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: conflictLoggingAttrs,
	}
	conflictLogging.MarkdownDescription = capellaschema.ReplicationConflictLoggingDescription
	return conflictLogging
}

// computedReplicationDetailsAttribute returns the computed source details, or the target details
// when isTarget is set. Only the target has a type.
func computedReplicationDetailsAttribute(isTarget bool) *schema.SingleNestedAttribute {
	detailsBuilder := capellaschema.NewSchemaBuilder("source_details", "replicationSource")
	detailsDescription := capellaschema.ReplicationSourceDetailsDescription
	if isTarget {
		detailsBuilder = capellaschema.NewSchemaBuilder("target_details", "replicationTarget")
		detailsDescription = capellaschema.ReplicationTargetDetailsDescription
	}

	referenceAttrs := func() map[string]schema.Attribute {
		attrs := make(map[string]schema.Attribute)
		id := stringAttribute([]string{computed})
		id.MarkdownDescription = capellaschema.ReplicationDetailsIdDescription
		capellaschema.AddAttr(attrs, "id", detailsBuilder, id)
		name := stringAttribute([]string{computed})
		name.MarkdownDescription = capellaschema.ReplicationDetailsNameDescription
		capellaschema.AddAttr(attrs, "name", detailsBuilder, name)
		return attrs
	}

	bucketAttrs := referenceAttrs()
	conflictResolutionType := stringAttribute([]string{computed})
	conflictResolutionType.MarkdownDescription = capellaschema.ReplicationDetailsConflictResolutionTypeDescription
	capellaschema.AddAttr(bucketAttrs, "conflict_resolution_type", detailsBuilder, conflictResolutionType)

	detailsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(detailsAttrs, "project", detailsBuilder, &schema.SingleNestedAttribute{Computed: true, Attributes: referenceAttrs()})
	capellaschema.AddAttr(detailsAttrs, "cluster", detailsBuilder, &schema.SingleNestedAttribute{Computed: true, Attributes: referenceAttrs()})
	capellaschema.AddAttr(detailsAttrs, "bucket", detailsBuilder, &schema.SingleNestedAttribute{Computed: true, Attributes: bucketAttrs})
	if isTarget {
		capellaschema.AddAttr(detailsAttrs, "type", detailsBuilder, stringAttribute([]string{computed}))
	}

	details := &schema.SingleNestedAttribute{
		Computed:      true,
		Attributes:    detailsAttrs,
		PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
	}
	details.MarkdownDescription = detailsDescription
	return details
}

// computedReplicationAuditAttribute returns the replication audit. Unlike computedAuditAttribute,
// it has only the two fields that the replication API returns.
func computedReplicationAuditAttribute() *schema.SingleNestedAttribute {
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_by", replicationBuilder, stringAttribute([]string{computed}), "ReplicationAuditData")
	capellaschema.AddAttr(auditAttrs, "created_at", replicationBuilder, stringAttribute([]string{computed}), "ReplicationAuditData")

	return &schema.SingleNestedAttribute{
		Computed:      true,
		Attributes:    auditAttrs,
		PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
	}
}
