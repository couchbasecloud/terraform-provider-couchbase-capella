package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceLogStreamingBuilder = capellaschema.NewSchemaBuilder("appServiceLogStreaming", "PostLogStreamingRequest")

// AppServiceLogStreamingSchema returns the schema for the app_service_log_streaming resource.
func AppServiceLogStreamingSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	// Required hierarchical IDs
	capellaschema.AddAttr(attrs, "organization_id", appServiceLogStreamingBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", appServiceLogStreamingBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "cluster_id", appServiceLogStreamingBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "app_service_id", appServiceLogStreamingBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))

	// User configured attributes
	capellaschema.AddAttr(attrs, "output_type", appServiceLogStreamingBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "credentials", appServiceLogStreamingBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Sensitive:  true,
		Attributes: buildCredentialsAttributes(),
	})

	// Read-only attributes
	capellaschema.AddAttr(attrs, "config_state", appServiceLogStreamingBuilder, stringAttribute([]string{computed}), "GetLogStreamingResponse")
	capellaschema.AddAttr(attrs, "streaming_state", appServiceLogStreamingBuilder, stringAttribute([]string{computed}), "GetLogStreamingResponse")

	return schema.Schema{
		Attributes: attrs,
	}
}

// buildCredentialsAttributes builds the credentials attributes map with all provider-specific credential blocks.
func buildCredentialsAttributes() map[string]schema.Attribute {
	credentialsAttrs := make(map[string]schema.Attribute)

	// Datadog credentials
	datadogAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(datadogAttrs, "api_key", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "datadog")
	capellaschema.AddAttr(datadogAttrs, "url", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "datadog")

	capellaschema.AddAttr(credentialsAttrs, "datadog", appServiceLogStreamingBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Sensitive:  true,
		Attributes: datadogAttrs,
	})

	// Dynatrace credentials
	dynatraceAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dynatraceAttrs, "api_token", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "dynatrace")
	capellaschema.AddAttr(dynatraceAttrs, "url", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "dynatrace")

	capellaschema.AddAttr(credentialsAttrs, "dynatrace", appServiceLogStreamingBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Sensitive:  true,
		Attributes: dynatraceAttrs,
	})

	// Elastic credentials
	elasticAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(elasticAttrs, "user", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "elastic")
	capellaschema.AddAttr(elasticAttrs, "password", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "elastic")
	capellaschema.AddAttr(elasticAttrs, "url", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "elastic")

	capellaschema.AddAttr(credentialsAttrs, "elastic", appServiceLogStreamingBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Sensitive:  true,
		Attributes: elasticAttrs,
	})

	// Generic HTTP credentials
	genericHttpAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(genericHttpAttrs, "url", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "generic_http")
	capellaschema.AddAttr(genericHttpAttrs, "user", appServiceLogStreamingBuilder, stringAttribute([]string{optional, sensitive}), "generic_http")
	capellaschema.AddAttr(genericHttpAttrs, "password", appServiceLogStreamingBuilder, stringAttribute([]string{optional, sensitive}), "generic_http")

	capellaschema.AddAttr(credentialsAttrs, "generic_http", appServiceLogStreamingBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Sensitive:  true,
		Attributes: genericHttpAttrs,
	})

	// Loki credentials
	lokiAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(lokiAttrs, "user", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "loki")
	capellaschema.AddAttr(lokiAttrs, "password", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "loki")
	capellaschema.AddAttr(lokiAttrs, "url", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "loki")

	capellaschema.AddAttr(credentialsAttrs, "loki", appServiceLogStreamingBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Sensitive:  true,
		Attributes: lokiAttrs,
	})

	// Splunk credentials
	splunkAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(splunkAttrs, "splunk_token", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "splunk")
	capellaschema.AddAttr(splunkAttrs, "url", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "splunk")

	capellaschema.AddAttr(credentialsAttrs, "splunk", appServiceLogStreamingBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Sensitive:  true,
		Attributes: splunkAttrs,
	})

	// Sumologic credentials
	sumologicAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(sumologicAttrs, "url", appServiceLogStreamingBuilder, stringAttribute([]string{required, sensitive}), "sumologic")

	capellaschema.AddAttr(credentialsAttrs, "sumologic", appServiceLogStreamingBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Sensitive:  true,
		Attributes: sumologicAttrs,
	})

	return credentialsAttrs
}
