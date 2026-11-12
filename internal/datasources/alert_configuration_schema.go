package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var alertConfigurationBuilder = capellaschema.NewSchemaBuilder("alertConfiguration", "GetAlertConfigurationResponse")

// OrganizationAlertConfigurationSchema returns the schema for the organization alert configuration data source.
func OrganizationAlertConfigurationSchema() schema.Schema {
	attrs := alertConfigurationValuesAttrs()
	capellaschema.AddAttr(attrs, "organization_id", alertConfigurationBuilder, requiredUUIDString())

	return schema.Schema{
		MarkdownDescription: "The organization alert configuration data source retrieves the alert configuration that an organization sets for itself, and the configuration in effect for it after the Capella defaults are applied.",
		Attributes:          attrs,
	}
}

// ProjectAlertConfigurationSchema returns the schema for the project alert configuration data source.
func ProjectAlertConfigurationSchema() schema.Schema {
	attrs := alertConfigurationValuesAttrs()
	capellaschema.AddAttr(attrs, "organization_id", alertConfigurationBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", alertConfigurationBuilder, requiredUUIDString())

	return schema.Schema{
		MarkdownDescription: "The project alert configuration data source retrieves the alert configuration that a project sets for itself, and the configuration in effect for it after inheritance from the organization and the Capella defaults.",
		Attributes:          attrs,
	}
}

// ClusterAlertConfigurationSchema returns the schema for the cluster alert configuration data source.
func ClusterAlertConfigurationSchema() schema.Schema {
	attrs := alertConfigurationValuesAttrs()
	capellaschema.AddAttr(attrs, "organization_id", alertConfigurationBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", alertConfigurationBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", alertConfigurationBuilder, requiredUUIDString())

	return schema.Schema{
		MarkdownDescription: "The cluster alert configuration data source retrieves the alert configuration that a cluster sets for itself, and the configuration in effect for it after inheritance from the project, the organization and the Capella defaults.",
		Attributes:          attrs,
	}
}

// The OpenAPI loader does not read a description written next to a $ref inside allOf, which is how
// the spec documents these four maps, so they carry a fallback in the loader's bullet format.
// AddAttr prefers the spec as soon as it can read one.
const (
	configsDescription             = "\n - The alert configs set at this level only. Empty when this level has never been configured."
	resolvedConfigsDescription     = "\n - The alert configs actually in effect at this level, after inheritance. Each field is taken from the closest level that sets it, so an alert configured at the organization and overridden at the cluster resolves to the cluster's fields on top of the organization's."
	mutedAlertsDescription         = "\n - The mute overrides set at this level only. Empty when this level sets none."
	resolvedMutedAlertsDescription = "\n - The mute overrides in effect at this level, after inheritance and after expiry has been applied. An alert listed here with `action` `mute` is not notifying now. Empty when no level sets any."
)

// alertConfigurationValuesAttrs returns the computed attributes shared by the three levels.
func alertConfigurationValuesAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)

	configs := &schema.MapNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: alertConfigurationEntryAttrs()},
	}
	configs.MarkdownDescription = configsDescription
	capellaschema.AddAttr(attrs, "configs", alertConfigurationBuilder, configs)

	resolvedConfigs := &schema.MapNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: resolvedAlertConfigurationEntryAttrs()},
	}
	resolvedConfigs.MarkdownDescription = resolvedConfigsDescription
	capellaschema.AddAttr(attrs, "resolved_configs", alertConfigurationBuilder, resolvedConfigs)

	capellaschema.AddAttr(attrs, "recipients", alertConfigurationBuilder, computedStringSet())

	mutedAlerts := &schema.MapNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: alertMuteStateAttrs()},
	}
	mutedAlerts.MarkdownDescription = mutedAlertsDescription
	capellaschema.AddAttr(attrs, "muted_alerts", alertConfigurationBuilder, mutedAlerts)

	resolvedMutedAlerts := &schema.MapNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: alertMuteStateAttrs()},
	}
	resolvedMutedAlerts.MarkdownDescription = resolvedMutedAlertsDescription
	capellaschema.AddAttr(attrs, "resolved_muted_alerts", alertConfigurationBuilder, resolvedMutedAlerts)

	return attrs
}

// resolvedAlertConfigurationEntryAttrs returns the attributes of a resolved config entry. Unlike
// alertConfigurationEntryAttrs it has no updated_at or updated_by, which the API never returns
// for resolved entries.
func resolvedAlertConfigurationEntryAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "frequency_hours", alertConfigurationEntryBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "threshold", alertConfigurationEntryBuilder, &schema.Float64Attribute{Computed: true})
	capellaschema.AddAttr(attrs, "warning_threshold", alertConfigurationEntryBuilder, &schema.Float64Attribute{Computed: true})
	capellaschema.AddAttr(attrs, "critical_threshold", alertConfigurationEntryBuilder, &schema.Float64Attribute{Computed: true})
	return attrs
}
