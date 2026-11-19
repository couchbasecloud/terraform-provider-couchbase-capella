package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceSelfServiceUpgradesBuilder = capellaschema.NewSchemaBuilder("appServiceSelfServiceUpgrades", "AppServiceMaintenance")

// AppServiceSelfServiceUpgradesSchema returns the schema for the AppServiceSelfServiceUpgrades data source.
func AppServiceSelfServiceUpgradesSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceSelfServiceUpgradesBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appServiceSelfServiceUpgradesBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceSelfServiceUpgradesBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceSelfServiceUpgradesBuilder, requiredUUIDString())

	// The API returns self-service upgrades in the App Service maintenance shape, so the items share
	// the maintenance fallback descriptions; AddAttr prefers the spec as soon as it has one.
	dataAttrs := make(map[string]schema.Attribute)
	for name, description := range map[string]string{
		"id":            capellaschema.SelfServiceUpgradeIdDescription,
		"kind":          capellaschema.MaintenanceKindDescription,
		"initiated_by":  capellaschema.MaintenanceInitiatedByDescription,
		"version_id":    capellaschema.MaintenanceVersionIdDescription,
		"current_state": capellaschema.MaintenanceCurrentStateDescription,
	} {
		attr := computedString()
		attr.MarkdownDescription = description
		capellaschema.AddAttr(dataAttrs, name, appServiceSelfServiceUpgradesBuilder, attr)
	}
	for name, description := range map[string]string{
		"queue_at":     capellaschema.MaintenanceQueueAtDescription,
		"completed_at": capellaschema.MaintenanceCompletedAtDescription,
	} {
		attr := computedRFC3339()
		attr.MarkdownDescription = description
		capellaschema.AddAttr(dataAttrs, name, appServiceSelfServiceUpgradesBuilder, attr)
	}
	capellaschema.AddAttr(dataAttrs, "audit", appServiceSelfServiceUpgradesBuilder, computedAudit())

	capellaschema.AddAttr(attrs, "data", appServiceSelfServiceUpgradesBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source lists the self-service upgrades of an App Service. " +
			"Maintenance initiated by Couchbase is listed by the `couchbase-capella_app_service_maintenances` data source instead.",
		Attributes: attrs,
	}
}
