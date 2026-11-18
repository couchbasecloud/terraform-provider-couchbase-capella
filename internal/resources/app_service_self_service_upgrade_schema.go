package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceSelfServiceUpgradeBuilder = capellaschema.NewSchemaBuilder("appServiceSelfServiceUpgrade", "AppServiceMaintenance")

// AppServiceSelfServiceUpgradeSchema returns the schema for the AppServiceSelfServiceUpgrade resource.
func AppServiceSelfServiceUpgradeSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceSelfServiceUpgradeBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServiceSelfServiceUpgradeBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceSelfServiceUpgradeBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceSelfServiceUpgradeBuilder, requiredUUIDStringAttribute())

	// The published OpenAPI spec does not carry the self-service upgrade schemas yet, so the attributes
	// below supply a fallback description; AddAttr prefers the spec as soon as it has one.
	id := stringAttribute([]string{computed, useStateForUnknown})
	id.MarkdownDescription = capellaschema.SelfServiceUpgradeIdDescription
	capellaschema.AddAttr(attrs, "id", appServiceSelfServiceUpgradeBuilder, id)

	// PUT replaces the version in place, so a new version does not require replace.
	versionId := requiredStringAttributeNoReplace()
	versionId.MarkdownDescription = capellaschema.SelfServiceUpgradeVersionIdDescription
	capellaschema.AddAttr(attrs, "version_id", appServiceSelfServiceUpgradeBuilder, versionId, "CreateSelfServiceUpgradeRequest")

	queueAt := rfc3339Attribute(optional, computed, useStateForUnknown)
	queueAt.MarkdownDescription = capellaschema.SelfServiceUpgradeQueueAtDescription
	capellaschema.AddAttr(attrs, "queue_at", appServiceSelfServiceUpgradeBuilder, queueAt, "CreateSelfServiceUpgradeRequest")

	currentState := stringAttribute([]string{computed})
	currentState.MarkdownDescription = capellaschema.MaintenanceCurrentStateDescription
	capellaschema.AddAttr(attrs, "current_state", appServiceSelfServiceUpgradeBuilder, currentState)

	completedAt := rfc3339Attribute(computed)
	completedAt.MarkdownDescription = capellaschema.MaintenanceCompletedAtDescription
	capellaschema.AddAttr(attrs, "completed_at", appServiceSelfServiceUpgradeBuilder, completedAt)

	capellaschema.AddAttr(attrs, "audit", appServiceSelfServiceUpgradeBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "Manages a self-service upgrade of an App Service, either now or at `queue_at`. " +
			"Pick `version_id` from the `couchbase-capella_app_service_self_service_upgrade_available_versions` data source.\n\n" +
			"`version_id` and `queue_at` can be changed in place while the upgrade is `pending`. " +
			"Destroying this resource cancels a `pending` upgrade; an upgrade that has finished (`completed`, `failed` or `cancelled`) is only removed from the Terraform state. " +
			"An upgrade that has already started cannot be cancelled.\n\n" +
			"The self-service upgrade endpoints are only available to organizations enabled for self-service upgrades.",
		Attributes: attrs,
	}
}
