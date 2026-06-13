package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Eventing function terminal activation states exposed via the `state` attribute. They mirror the
// terminal runtime status reported by the API so a refresh maps a status directly back to a state.
const (
	eventingStateDeployed   = "deployed"
	eventingStateUndeployed = "undeployed"
	eventingStatePaused     = "paused"
	eventingStateResumed    = "resumed"
)

var eventingFunctionBuilder = capellaschema.NewSchemaBuilder("eventingFunction")

// EventingFunctionSchema defines the schema for the eventing function resource.
func EventingFunctionSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", eventingFunctionBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", eventingFunctionBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", eventingFunctionBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "name", eventingFunctionBuilder, requiredNonEmptyStringAttribute())
	capellaschema.AddAttr(attrs, "description", eventingFunctionBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "code", eventingFunctionBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "state", eventingFunctionBuilder, stringAttribute(
		[]string{optional},
		validator.String(stringvalidator.OneOf(
			eventingStateDeployed, eventingStateUndeployed, eventingStatePaused, eventingStateResumed,
		)),
	))

	capellaschema.AddAttr(attrs, "event_source", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: keyspaceAttributes(),
	})
	capellaschema.AddAttr(attrs, "event_metadata_storage", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: keyspaceAttributes(),
	})

	capellaschema.AddAttr(attrs, "settings", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Computed:   true,
		Default:    settingsDefault(),
		Attributes: settingsAttributes(),
	})

	capellaschema.AddAttr(attrs, "bindings", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: bindingsAttributes(),
	})

	return schema.Schema{
		MarkdownDescription: "Manages an eventing function on a Capella cluster, including its JavaScript code, " +
			"source and metadata keyspaces, runtime settings, bindings, and deployment state.",
		Attributes: attrs,
	}
}

func keyspaceAttributes() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "bucket", eventingFunctionBuilder, requiredStringAttributeNoReplace(), "EventingFunctionKeyspace")
	capellaschema.AddAttr(attrs, "scope", eventingFunctionBuilder, stringAttribute([]string{optional, computed}), "EventingFunctionKeyspace")
	capellaschema.AddAttr(attrs, "collection", eventingFunctionBuilder, stringAttribute([]string{optional, computed}), "EventingFunctionKeyspace")
	return attrs
}

func settingsAttributes() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "worker_count", eventingFunctionBuilder, int64DefaultAttribute(1, optional, computed), "EventingFunctionSettings")
	capellaschema.AddAttr(attrs, "script_timeout", eventingFunctionBuilder, int64DefaultAttribute(60, optional, computed), "EventingFunctionSettings")

	sqlConsistency := stringDefaultAttribute("none", optional, computed)
	sqlConsistency.Validators = append(sqlConsistency.Validators, validator.String(stringvalidator.OneOf("none", "request")))
	capellaschema.AddAttr(attrs, "sql_consistency", eventingFunctionBuilder, sqlConsistency, "EventingFunctionSettings")

	languageCompatibility := stringDefaultAttribute("7.2.0", optional, computed)
	languageCompatibility.Validators = append(languageCompatibility.Validators, validator.String(stringvalidator.OneOf("6.0.0", "6.5.0", "6.6.2", "7.2.0")))
	capellaschema.AddAttr(attrs, "language_compatibility", eventingFunctionBuilder, languageCompatibility, "EventingFunctionSettings")

	feedBoundary := stringDefaultAttribute("from_now", optional, computed)
	feedBoundary.Validators = append(feedBoundary.Validators, validator.String(stringvalidator.OneOf("everything", "from_now")))
	capellaschema.AddAttr(attrs, "feed_boundary", eventingFunctionBuilder, feedBoundary, "EventingFunctionSettings")

	capellaschema.AddAttr(attrs, "max_timer_context_size", eventingFunctionBuilder, int64DefaultAttribute(1024, optional, computed), "EventingFunctionSettings")
	capellaschema.AddAttr(attrs, "allow_sync_documents", eventingFunctionBuilder, boolDefaultAttribute(false, optional, computed), "EventingFunctionSettings")
	capellaschema.AddAttr(attrs, "cursor_aware", eventingFunctionBuilder, boolDefaultAttribute(false, optional, computed), "EventingFunctionSettings")
	return attrs
}

// settingsDefault supplies the full default settings object so that the defaults are applied even
// when the entire settings block is omitted from configuration (the per-attribute defaults only
// apply when the block is present but a child is omitted). The attribute types are derived from
// settingsAttributes so the two stay in sync.
func settingsDefault() defaults.Object {
	attrTypes := make(map[string]attr.Type)
	for name, attribute := range settingsAttributes() {
		attrTypes[name] = attribute.GetType()
	}

	return objectdefault.StaticValue(types.ObjectValueMust(attrTypes, map[string]attr.Value{
		"worker_count":           types.Int64Value(1),
		"script_timeout":         types.Int64Value(60),
		"sql_consistency":        types.StringValue("none"),
		"language_compatibility": types.StringValue("7.2.0"),
		"feed_boundary":          types.StringValue("from_now"),
		"max_timer_context_size": types.Int64Value(1024),
		"allow_sync_documents":   types.BoolValue(false),
		"cursor_aware":           types.BoolValue(false),
	}))
}

func bindingsAttributes() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)

	bucketAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(bucketAttrs, "alias", eventingFunctionBuilder, requiredStringAttributeNoReplace(), "EventingFunctionBucketBinding")
	capellaschema.AddAttr(bucketAttrs, "bucket", eventingFunctionBuilder, requiredStringAttributeNoReplace(), "EventingFunctionBucketBinding")
	capellaschema.AddAttr(bucketAttrs, "scope", eventingFunctionBuilder, stringAttribute([]string{optional, computed}), "EventingFunctionBucketBinding")
	capellaschema.AddAttr(bucketAttrs, "collection", eventingFunctionBuilder, stringAttribute([]string{optional, computed}), "EventingFunctionBucketBinding")
	capellaschema.AddAttr(bucketAttrs, "permission", eventingFunctionBuilder, stringAttribute(
		[]string{optional, computed}, validator.String(stringvalidator.OneOf("read", "readWrite"))), "EventingFunctionBucketBinding")

	capellaschema.AddAttr(attrs, "buckets", eventingFunctionBuilder, &schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: bucketAttrs,
		},
	})

	authAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(authAttrs, "type", eventingFunctionBuilder, stringAttribute(
		[]string{required}, validator.String(stringvalidator.OneOf("none", "basic", "bearer", "digest"))))
	capellaschema.AddAttr(authAttrs, "username", eventingFunctionBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(authAttrs, "password", eventingFunctionBuilder, stringAttribute([]string{optional, sensitive}))
	capellaschema.AddAttr(authAttrs, "bearer_token", eventingFunctionBuilder, stringAttribute([]string{optional, sensitive}))

	urlAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(urlAttrs, "alias", eventingFunctionBuilder, requiredStringAttributeNoReplace(), "EventingFunctionUrlBinding")
	capellaschema.AddAttr(urlAttrs, "url", eventingFunctionBuilder, requiredStringAttributeNoReplace(), "EventingFunctionUrlBinding")
	capellaschema.AddAttr(urlAttrs, "allow_cookies", eventingFunctionBuilder, boolAttribute(optional, computed), "EventingFunctionUrlBinding")
	capellaschema.AddAttr(urlAttrs, "validate_tls_certificate", eventingFunctionBuilder, boolAttribute(optional, computed), "EventingFunctionUrlBinding")
	capellaschema.AddAttr(urlAttrs, "authentication", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: authAttrs,
	})

	capellaschema.AddAttr(attrs, "urls", eventingFunctionBuilder, &schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: urlAttrs,
		},
	})

	constantAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(constantAttrs, "alias", eventingFunctionBuilder, requiredStringAttributeNoReplace(), "EventingFunctionConstantBinding")
	capellaschema.AddAttr(constantAttrs, "value", eventingFunctionBuilder, requiredStringAttributeNoReplace(), "EventingFunctionConstantBinding")

	capellaschema.AddAttr(attrs, "constants", eventingFunctionBuilder, &schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: constantAttrs,
		},
	})

	return attrs
}
