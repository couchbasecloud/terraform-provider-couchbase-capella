package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceLogStreamingActivationStatusBuilder = capellaschema.NewSchemaBuilder("appServiceLogStreamingActivationStatus")

// AppServiceLogStreamingActivationStatusSchema returns the schema for the
// app_service_log_streaming_activation_status resource.
func AppServiceLogStreamingActivationStatusSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	// Required hierarchical IDs
	capellaschema.AddAttr(attrs, "organization_id", appServiceLogStreamingActivationStatusBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", appServiceLogStreamingActivationStatusBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "cluster_id", appServiceLogStreamingActivationStatusBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "app_service_id", appServiceLogStreamingActivationStatusBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))

	// User-specified desired Log Streaming state
	capellaschema.AddAttr(attrs, "state", appServiceLogStreamingActivationStatusBuilder, stringAttribute(
		[]string{required},
		validator.String(
			stringvalidator.OneOf(
				string(apigen.GetLogStreamingResponseConfigStatePaused),
				string(apigen.GetLogStreamingResponseConfigStateEnabled),
			),
		),
	))

	return schema.Schema{
		MarkdownDescription: "Manages the activation state (paused/enabled) of log streaming on an App Service.",
		Attributes:          attrs,
	}
}
