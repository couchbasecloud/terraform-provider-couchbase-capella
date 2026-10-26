package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var alertConfigurationBuilder = capellaschema.NewSchemaBuilder("alertConfiguration")

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

// alertConfigurationValuesAttrs returns the computed attributes shared by the three levels.
func alertConfigurationValuesAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "configs", alertConfigurationBuilder, &schema.MapNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: alertConfigurationEntryAttrs()},
	})
	capellaschema.AddAttr(attrs, "resolved_configs", alertConfigurationBuilder, &schema.MapNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: resolvedAlertConfigurationEntryAttrs()},
	})
	capellaschema.AddAttr(attrs, "recipients", alertConfigurationBuilder, computedStringSet())
	capellaschema.AddAttr(attrs, "muted_alerts", alertConfigurationBuilder, &schema.MapNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: alertMuteStateAttrs()},
	})
	capellaschema.AddAttr(attrs, "resolved_muted_alerts", alertConfigurationBuilder, &schema.MapNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: alertMuteStateAttrs()},
	})
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
