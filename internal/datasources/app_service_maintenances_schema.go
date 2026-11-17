package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceMaintenancesBuilder = capellaschema.NewSchemaBuilder("appServiceMaintenances", "AppServiceMaintenance")

// AppServiceMaintenancesSchema returns the schema for the AppServiceMaintenances data source.
func AppServiceMaintenancesSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceMaintenancesBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appServiceMaintenancesBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceMaintenancesBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceMaintenancesBuilder, requiredUUIDString())

	// The published OpenAPI spec does not carry the maintenance schemas yet, so these attributes
	// supply a fallback description; AddAttr prefers the spec as soon as it has one.
	dataAttrs := make(map[string]schema.Attribute)
	for name, description := range map[string]string{
		"id":            capellaschema.MaintenanceIdDescription,
		"kind":          capellaschema.MaintenanceKindDescription,
		"initiated_by":  capellaschema.MaintenanceInitiatedByDescription,
		"version_id":    capellaschema.MaintenanceVersionIdDescription,
		"current_state": capellaschema.MaintenanceCurrentStateDescription,
	} {
		attr := computedString()
		attr.MarkdownDescription = description
		capellaschema.AddAttr(dataAttrs, name, appServiceMaintenancesBuilder, attr)
	}
	for name, description := range map[string]string{
		"queue_at":     capellaschema.MaintenanceQueueAtDescription,
		"completed_at": capellaschema.MaintenanceCompletedAtDescription,
	} {
		attr := computedRFC3339()
		attr.MarkdownDescription = description
		capellaschema.AddAttr(dataAttrs, name, appServiceMaintenancesBuilder, attr)
	}
	capellaschema.AddAttr(dataAttrs, "audit", appServiceMaintenancesBuilder, computedAudit())

	capellaschema.AddAttr(attrs, "data", appServiceMaintenancesBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source lists the customer-visible maintenance records of an App Service. " +
			"Use a record's `id` as the `maintenance_id` of the `couchbase-capella_app_service_maintenance_schedule` resource to reschedule it.",
		Attributes: attrs,
	}
}
