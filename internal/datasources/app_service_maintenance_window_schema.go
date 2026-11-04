package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceMaintenanceWindowBuilder = capellaschema.NewSchemaBuilder("appServiceMaintenanceWindow")

// AppServiceMaintenanceWindowSchema returns the schema for the AppServiceMaintenanceWindow data source.
func AppServiceMaintenanceWindowSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceMaintenanceWindowBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appServiceMaintenanceWindowBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceMaintenanceWindowBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceMaintenanceWindowBuilder, requiredUUIDString())

	preferredDay := computedString()
	preferredDay.MarkdownDescription = "The day of the week the maintenance window starts on, in lowercase. Null when no window is set."
	capellaschema.AddAttr(attrs, "preferred_day", appServiceMaintenanceWindowBuilder, preferredDay)

	preferredHour := computedInt64()
	preferredHour.MarkdownDescription = "The hour of the day (0-23) the maintenance window starts at. Null when no window is set."
	capellaschema.AddAttr(attrs, "preferred_hour", appServiceMaintenanceWindowBuilder, preferredHour)

	preferredMinute := computedInt64()
	preferredMinute.MarkdownDescription = "The minute of the hour (0-59) the maintenance window starts at. Null when no window is set."
	capellaschema.AddAttr(attrs, "preferred_minute", appServiceMaintenanceWindowBuilder, preferredMinute)

	return schema.Schema{
		MarkdownDescription: "The App Service maintenance window data source retrieves the preferred weekly maintenance window for an operational App Service. " +
			"The maintenance window endpoints are only available to organizations enabled for self-service upgrades.",
		Attributes: attrs,
	}
}
