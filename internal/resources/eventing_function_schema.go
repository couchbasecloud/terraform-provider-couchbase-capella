package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

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

// URL binding authentication types accepted by the eventing API. Each type selects which credential
// fields of a URL binding authentication block must be set and which must be absent.
const (
	eventingURLAuthNone   = "none"
	eventingURLAuthBasic  = "basic"
	eventingURLAuthBearer = "bearer"
	eventingURLAuthDigest = "digest"
)

var eventingFunctionBuilder = capellaschema.NewSchemaBuilder("eventingFunction")

// EventingFunctionSchema defines the schema for the eventing function resource.
func EventingFunctionSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", eventingFunctionBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", eventingFunctionBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", eventingFunctionBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "name", eventingFunctionBuilder, requiredNonEmptyStringAttribute())
	capellaschema.AddAttr(attrs, "description", eventingFunctionBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(attrs, "code", eventingFunctionBuilder, stringAttribute([]string{required, sensitive}))
	capellaschema.AddAttr(
		attrs,
		"state",
		eventingFunctionBuilder,
		&schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  stringdefault.StaticString(eventingStateUndeployed),
			Validators: []validator.String{
				stringvalidator.OneOf(
					eventingStateDeployed,
					eventingStateUndeployed,
					eventingStatePaused,
				),
			},
		},
	)

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
	capellaschema.AddAttr(attrs, "worker_count", eventingFunctionBuilder, int64Attribute(optional, computed, useStateForUnknown), "EventingFunctionSettings")
	capellaschema.AddAttr(attrs, "script_timeout", eventingFunctionBuilder, int64Attribute(optional, computed, useStateForUnknown), "EventingFunctionSettings")

	capellaschema.AddAttr(attrs, "sql_consistency", eventingFunctionBuilder, stringAttribute(
		[]string{optional, computed, useStateForUnknown}, validator.String(stringvalidator.OneOf("none", "request"))), "EventingFunctionSettings")

	capellaschema.AddAttr(attrs, "language_compatibility", eventingFunctionBuilder, stringAttribute(
		[]string{optional, computed, useStateForUnknown}, validator.String(stringvalidator.OneOf("6.0.0", "6.5.0", "6.6.2", "7.2.0"))), "EventingFunctionSettings")

	capellaschema.AddAttr(attrs, "feed_boundary", eventingFunctionBuilder, stringAttribute(
		[]string{optional, computed, useStateForUnknown}, validator.String(stringvalidator.OneOf("everything", "from_now"))), "EventingFunctionSettings")

	capellaschema.AddAttr(attrs, "max_timer_context_size", eventingFunctionBuilder, int64Attribute(optional, computed, useStateForUnknown), "EventingFunctionSettings")
	capellaschema.AddAttr(attrs, "allow_sync_documents", eventingFunctionBuilder, boolAttribute(optional, computed, useStateForUnknown), "EventingFunctionSettings")
	capellaschema.AddAttr(attrs, "cursor_aware", eventingFunctionBuilder, boolAttribute(optional, computed, useStateForUnknown), "EventingFunctionSettings")
	return attrs
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
		[]string{required},
		validator.String(
			stringvalidator.OneOf(
				eventingURLAuthNone,
				eventingURLAuthBasic,
				eventingURLAuthBearer,
				eventingURLAuthDigest,
			),
		),
	))
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
