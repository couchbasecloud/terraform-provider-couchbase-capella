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
	capellaschema.AddAttr(attrs, "organization_id", appServiceLogStreamingActivationStatusBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServiceLogStreamingActivationStatusBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceLogStreamingActivationStatusBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceLogStreamingActivationStatusBuilder, requiredUUIDStringAttribute())

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
		MarkdownDescription: `Manages the activation state of Log Streaming for an App Service.
This resource allows you to put Log Streaming in either a 'paused' state (meaning logs stop streaming but the log collector config is retained by Capella), or an 'enabled' state.
Log Streaming must already be set up on the App Service (i.e. it must not be in a 'disabled' state) in order to manage it's activation state.`,
		Attributes: attrs,
	}
}
