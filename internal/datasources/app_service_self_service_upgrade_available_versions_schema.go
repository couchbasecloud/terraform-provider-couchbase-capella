package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceSelfServiceUpgradeAvailableVersionsBuilder = capellaschema.NewSchemaBuilder("appServiceSelfServiceUpgradeAvailableVersions", "SelfServiceUpgradeAvailableVersion")

// AppServiceSelfServiceUpgradeAvailableVersionsSchema returns the schema for the AppServiceSelfServiceUpgradeAvailableVersions data source.
func AppServiceSelfServiceUpgradeAvailableVersionsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceSelfServiceUpgradeAvailableVersionsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appServiceSelfServiceUpgradeAvailableVersionsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceSelfServiceUpgradeAvailableVersionsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceSelfServiceUpgradeAvailableVersionsBuilder, requiredUUIDString())

	// The published OpenAPI spec does not carry the self-service upgrade schemas yet, so these attributes
	// supply a fallback description; AddAttr prefers the spec as soon as it has one.
	dataAttrs := make(map[string]schema.Attribute)
	for name, description := range map[string]string{
		"id":          capellaschema.AvailableVersionIdDescription,
		"title":       capellaschema.AvailableVersionTitleDescription,
		"description": capellaschema.AvailableVersionDescriptionDescription,
		"priority":    capellaschema.AvailableVersionPriorityDescription,
	} {
		attr := computedString()
		attr.MarkdownDescription = description
		capellaschema.AddAttr(dataAttrs, name, appServiceSelfServiceUpgradeAvailableVersionsBuilder, attr)
	}
	releaseNotes := computedStringList()
	releaseNotes.MarkdownDescription = capellaschema.AvailableVersionReleaseNotesDescription
	capellaschema.AddAttr(dataAttrs, "release_notes", appServiceSelfServiceUpgradeAvailableVersionsBuilder, releaseNotes)

	capellaschema.AddAttr(attrs, "data", appServiceSelfServiceUpgradeAvailableVersionsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source lists the versions an App Service can be upgraded to through a self-service upgrade. " +
			"Use a version's `id` as the `version_id` of the `couchbase-capella_app_service_self_service_upgrade` resource.",
		Attributes: attrs,
	}
}
