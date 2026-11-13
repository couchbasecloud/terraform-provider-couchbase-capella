package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceMaintenanceScheduleBuilder = capellaschema.NewSchemaBuilder("appServiceMaintenanceSchedule", "AppServiceMaintenance")

// AppServiceMaintenanceScheduleSchema returns the schema for the AppServiceMaintenanceSchedule resource.
func AppServiceMaintenanceScheduleSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceMaintenanceScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServiceMaintenanceScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceMaintenanceScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceMaintenanceScheduleBuilder, requiredUUIDStringAttribute())

	// The published OpenAPI spec does not carry the maintenance schemas yet, so the attributes below
	// supply a fallback description; AddAttr prefers the spec as soon as it has one.
	maintenanceId := requiredNonEmptyStringAttribute()
	maintenanceId.MarkdownDescription = capellaschema.MaintenanceIdDescription
	capellaschema.AddAttr(attrs, "maintenance_id", appServiceMaintenanceScheduleBuilder, maintenanceId)

	queueAt := rfc3339Attribute(required)
	queueAt.MarkdownDescription = capellaschema.MaintenanceQueueAtDescription
	capellaschema.AddAttr(attrs, "queue_at", appServiceMaintenanceScheduleBuilder, queueAt)

	for name, description := range map[string]string{
		"kind":          capellaschema.MaintenanceKindDescription,
		"initiated_by":  capellaschema.MaintenanceInitiatedByDescription,
		"version_id":    capellaschema.MaintenanceVersionIdDescription,
		"current_state": capellaschema.MaintenanceCurrentStateDescription,
	} {
		attr := stringAttribute([]string{computed})
		attr.MarkdownDescription = description
		capellaschema.AddAttr(attrs, name, appServiceMaintenanceScheduleBuilder, attr)
	}

	completedAt := rfc3339Attribute(computed)
	completedAt.MarkdownDescription = capellaschema.MaintenanceCompletedAtDescription
	capellaschema.AddAttr(attrs, "completed_at", appServiceMaintenanceScheduleBuilder, completedAt)

	rescheduleBy := rfc3339Attribute(computed)
	rescheduleBy.MarkdownDescription = capellaschema.MaintenanceRescheduleByDescription
	capellaschema.AddAttr(attrs, "reschedule_by", appServiceMaintenanceScheduleBuilder, rescheduleBy)

	capellaschema.AddAttr(attrs, "audit", appServiceMaintenanceScheduleBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "Manages the queue time of an existing App Service maintenance. " +
			"This resource does not create a maintenance: it adopts one listed by the `couchbase-capella_app_service_maintenances` data source " +
			"and reschedules it to `queue_at`. Only `pending` maintenance can be rescheduled, to a time in the future and no later than `reschedule_by`.\n\n" +
			"Destroying this resource only removes it from the Terraform state; the maintenance keeps its last queue time.",
		Attributes: attrs,
	}
}
