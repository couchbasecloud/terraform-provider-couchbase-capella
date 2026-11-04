package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceMaintenanceWindowBuilder = capellaschema.NewSchemaBuilder("appServiceMaintenanceWindow")

func AppServiceMaintenanceWindowSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceMaintenanceWindowBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServiceMaintenanceWindowBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceMaintenanceWindowBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceMaintenanceWindowBuilder, requiredUUIDStringAttribute())

	preferredDay := stringAttribute([]string{required},
		validator.String(stringvalidator.OneOf("monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday")))
	preferredDay.MarkdownDescription = "The day of the week the maintenance window starts on, in lowercase: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday` or `sunday`."
	capellaschema.AddAttr(attrs, "preferred_day", appServiceMaintenanceWindowBuilder, preferredDay)

	preferredHour := int64Attribute(required)
	preferredHour.Validators = []validator.Int64{int64validator.Between(0, 23)}
	preferredHour.MarkdownDescription = "The hour of the day the maintenance window starts at, from 0 to 23."
	capellaschema.AddAttr(attrs, "preferred_hour", appServiceMaintenanceWindowBuilder, preferredHour)

	preferredMinute := int64Attribute(required)
	preferredMinute.Validators = []validator.Int64{int64validator.Between(0, 59)}
	preferredMinute.MarkdownDescription = "The minute of the hour the maintenance window starts at, from 0 to 59."
	capellaschema.AddAttr(attrs, "preferred_minute", appServiceMaintenanceWindowBuilder, preferredMinute)

	return schema.Schema{
		MarkdownDescription: "Manages the preferred weekly maintenance window for an operational App Service.\n\n" +
			"An App Service has at most one maintenance window. If a window is already set (for example through the Capella UI), " +
			"creating this resource fails instead of overwriting it; import the existing window with " +
			"`terraform import couchbase-capella_app_service_maintenance_window.<name> organization_id=<organization_id>,project_id=<project_id>,cluster_id=<cluster_id>,app_service_id=<app_service_id>` " +
			"and manage it from Terraform.\n\n" +
			"The maintenance window endpoints are only available to organizations enabled for self-service upgrades.",
		Attributes: attrs,
	}
}
