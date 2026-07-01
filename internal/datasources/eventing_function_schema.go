package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Builders for the eventing function schema and its nested objects. The OpenAPI schema names are
// passed so that field descriptions resolve from the published Capella API specification.
var (
	eventingFunctionBuilder                = capellaschema.NewSchemaBuilder("eventing_function", "EventingFunction")
	eventingFunctionKeyspaceBuilder        = capellaschema.NewSchemaBuilder("eventing_function_keyspace", "EventingFunctionKeyspace")
	eventingFunctionSettingsBuilder        = capellaschema.NewSchemaBuilder("eventing_function_settings", "EventingFunctionSettings")
	eventingFunctionBindingsBuilder        = capellaschema.NewSchemaBuilder("eventing_function_bindings", "EventingFunctionBindings")
	eventingFunctionBucketBindingBuilder   = capellaschema.NewSchemaBuilder("eventing_function_bucket_binding", "EventingFunctionBucketBinding")
	eventingFunctionUrlBindingBuilder      = capellaschema.NewSchemaBuilder("eventing_function_url_binding", "EventingFunctionUrlBinding")
	eventingFunctionAuthenticationBuilder  = capellaschema.NewSchemaBuilder("eventing_function_url_binding_authentication", "URLBindingAuthentication")
	eventingFunctionConstantBindingBuilder = capellaschema.NewSchemaBuilder("eventing_function_constant_binding", "EventingFunctionConstantBinding")
)

func getEventingFunctionKeyspaceAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "bucket", eventingFunctionKeyspaceBuilder, computedString())
	capellaschema.AddAttr(attrs, "scope", eventingFunctionKeyspaceBuilder, computedString())
	capellaschema.AddAttr(attrs, "collection", eventingFunctionKeyspaceBuilder, computedString())
	return attrs
}

func getEventingFunctionSettingsAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "worker_count", eventingFunctionSettingsBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "script_timeout", eventingFunctionSettingsBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "sql_consistency", eventingFunctionSettingsBuilder, computedString())
	capellaschema.AddAttr(attrs, "language_compatibility", eventingFunctionSettingsBuilder, computedString())
	capellaschema.AddAttr(attrs, "feed_boundary", eventingFunctionSettingsBuilder, computedString())
	capellaschema.AddAttr(attrs, "max_timer_context_size", eventingFunctionSettingsBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "allow_sync_documents", eventingFunctionSettingsBuilder, computedBool())
	capellaschema.AddAttr(attrs, "cursor_aware", eventingFunctionSettingsBuilder, computedBool())
	return attrs
}

func getEventingFunctionBucketBindingAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "alias", eventingFunctionBucketBindingBuilder, computedString())
	capellaschema.AddAttr(attrs, "bucket", eventingFunctionBucketBindingBuilder, computedString())
	capellaschema.AddAttr(attrs, "scope", eventingFunctionBucketBindingBuilder, computedString())
	capellaschema.AddAttr(attrs, "collection", eventingFunctionBucketBindingBuilder, computedString())
	capellaschema.AddAttr(attrs, "permission", eventingFunctionBucketBindingBuilder, computedString())
	return attrs
}

func getEventingFunctionAuthenticationAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "type", eventingFunctionAuthenticationBuilder, computedString())
	capellaschema.AddAttr(attrs, "username", eventingFunctionAuthenticationBuilder, computedString())
	// password and bearer_token are credential fields, so they are marked sensitive. The eventing
	// service already redacts them, returning five asterisks rather than the stored secret.
	capellaschema.AddAttr(attrs, "password", eventingFunctionAuthenticationBuilder, &schema.StringAttribute{
		Computed:  true,
		Sensitive: true,
	})
	capellaschema.AddAttr(attrs, "bearer_token", eventingFunctionAuthenticationBuilder, &schema.StringAttribute{
		Computed:  true,
		Sensitive: true,
	})
	return attrs
}

func getEventingFunctionUrlBindingAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "alias", eventingFunctionUrlBindingBuilder, computedString())
	capellaschema.AddAttr(attrs, "url", eventingFunctionUrlBindingBuilder, computedString())
	capellaschema.AddAttr(attrs, "allow_cookies", eventingFunctionUrlBindingBuilder, computedBool())
	capellaschema.AddAttr(attrs, "validate_tls_certificate", eventingFunctionUrlBindingBuilder, computedBool())
	capellaschema.AddAttr(attrs, "authentication", eventingFunctionUrlBindingBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getEventingFunctionAuthenticationAttrs(),
	})
	return attrs
}

func getEventingFunctionConstantBindingAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "alias", eventingFunctionConstantBindingBuilder, computedString())
	capellaschema.AddAttr(attrs, "value", eventingFunctionConstantBindingBuilder, computedString())
	return attrs
}

func getEventingFunctionBindingsAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "buckets", eventingFunctionBindingsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getEventingFunctionBucketBindingAttrs(),
		},
	})
	capellaschema.AddAttr(attrs, "urls", eventingFunctionBindingsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getEventingFunctionUrlBindingAttrs(),
		},
	})
	capellaschema.AddAttr(attrs, "constants", eventingFunctionBindingsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getEventingFunctionConstantBindingAttrs(),
		},
	})
	return attrs
}

// EventingFunctionSchema returns the schema for the single eventing function data source.
func EventingFunctionSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", eventingFunctionBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", eventingFunctionBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", eventingFunctionBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "name", eventingFunctionBuilder, requiredString())

	// export omits the read-only status field from the response so the result can be reused as a
	// create payload. It maps to the export query parameter on the get eventing function endpoint.
	capellaschema.AddAttr(attrs, "export", eventingFunctionBuilder, &schema.BoolAttribute{
		Optional: true,
	})

	capellaschema.AddAttr(attrs, "description", eventingFunctionBuilder, computedString())
	capellaschema.AddAttr(attrs, "status", eventingFunctionBuilder, computedString())
	capellaschema.AddAttr(attrs, "code", eventingFunctionBuilder, &schema.StringAttribute{
		Computed:  true,
		Sensitive: true,
	})

	capellaschema.AddAttr(attrs, "event_source", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getEventingFunctionKeyspaceAttrs(),
	})
	capellaschema.AddAttr(attrs, "event_metadata_storage", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getEventingFunctionKeyspaceAttrs(),
	})
	capellaschema.AddAttr(attrs, "settings", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getEventingFunctionSettingsAttrs(),
	})
	capellaschema.AddAttr(attrs, "bindings", eventingFunctionBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getEventingFunctionBindingsAttrs(),
	})

	return schema.Schema{
		MarkdownDescription: "The eventing function data source retrieves a single eventing function in a cluster.",
		Attributes:          attrs,
	}
}
