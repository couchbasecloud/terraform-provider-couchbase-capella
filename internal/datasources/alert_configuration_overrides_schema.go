package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	alertConfigurationOverridesBuilder = capellaschema.NewSchemaBuilder("alertConfigurationOverrides", "AlertConfigurationOverride")
	alertConfigurationEntryBuilder     = capellaschema.NewSchemaBuilder("alertConfigurationEntry", "AlertConfigurationEntry")
	alertMuteStateBuilder              = capellaschema.NewSchemaBuilder("alertMuteState", "AlertMuteState")
)

// AlertConfigurationOverridesSchema returns the schema for the alert configuration overrides data source.
func AlertConfigurationOverridesSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", alertConfigurationOverridesBuilder, requiredUUIDString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "scope", alertConfigurationOverridesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", alertConfigurationOverridesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", alertConfigurationOverridesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "resource_id", alertConfigurationOverridesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "alert_kind", alertConfigurationOverridesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "config", alertConfigurationOverridesBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: alertConfigurationEntryAttrs(),
	})
	capellaschema.AddAttr(dataAttrs, "previous_config", alertConfigurationOverridesBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: alertConfigurationEntryAttrs(),
	})
	capellaschema.AddAttr(dataAttrs, "muted_alert", alertConfigurationOverridesBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: alertMuteStateAttrs(),
	})
	capellaschema.AddAttr(dataAttrs, "recipients", alertConfigurationOverridesBuilder, computedStringSet())
	capellaschema.AddAttr(dataAttrs, "previous_recipients", alertConfigurationOverridesBuilder, computedStringSet())

	capellaschema.AddAttr(attrs, "data", alertConfigurationOverridesBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	}, "ListAlertConfigurationOverridesResponse")

	return schema.Schema{
		MarkdownDescription: "The alert configuration overrides data source lists every alert that the organization, a project or a cluster configures or disables itself instead of inheriting from the level above. Each entry is one alert at one level.",
		Attributes:          attrs,
	}
}

func alertConfigurationEntryAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "frequency_hours", alertConfigurationEntryBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "threshold", alertConfigurationEntryBuilder, &schema.Float64Attribute{Computed: true})
	capellaschema.AddAttr(attrs, "warning_threshold", alertConfigurationEntryBuilder, &schema.Float64Attribute{Computed: true})
	capellaschema.AddAttr(attrs, "critical_threshold", alertConfigurationEntryBuilder, &schema.Float64Attribute{Computed: true})
	capellaschema.AddAttr(attrs, "updated_at", alertConfigurationEntryBuilder, computedRFC3339())
	capellaschema.AddAttr(attrs, "updated_by", alertConfigurationEntryBuilder, computedString())
	return attrs
}

func alertMuteStateAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "action", alertMuteStateBuilder, computedString())
	capellaschema.AddAttr(attrs, "muted_from", alertMuteStateBuilder, computedRFC3339())
	capellaschema.AddAttr(attrs, "muted_until", alertMuteStateBuilder, computedRFC3339())
	capellaschema.AddAttr(attrs, "updated_at", alertMuteStateBuilder, computedRFC3339())
	capellaschema.AddAttr(attrs, "updated_by", alertMuteStateBuilder, computedString())
	return attrs
}
