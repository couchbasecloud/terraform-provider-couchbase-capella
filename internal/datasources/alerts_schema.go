package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var alertsBuilder = capellaschema.NewSchemaBuilder("alerts", "Alert")

// AlertsSchema returns the schema for the alerts data source.
func AlertsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", alertsBuilder, requiredUUIDString())

	capellaschema.AddAttr(attrs, "status", alertsBuilder, optionalStringSetOf(stringvalidator.OneOf("firing", "resolved")))
	capellaschema.AddAttr(attrs, "cluster_ids", alertsBuilder, optionalStringSetOf(stringvalidator.RegexMatches(uuidRegex, "must be a valid UUID")))
	capellaschema.AddAttr(attrs, "project_ids", alertsBuilder, optionalStringSetOf(stringvalidator.RegexMatches(uuidRegex, "must be a valid UUID")))
	capellaschema.AddAttr(attrs, "severity_levels", alertsBuilder, optionalStringSetOf(stringvalidator.OneOf("debug", "info", "warning", "critical")))
	capellaschema.AddAttr(attrs, "from", alertsBuilder, optionalRFC3339())
	capellaschema.AddAttr(attrs, "to", alertsBuilder, optionalRFC3339())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "alert_id", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "event_id", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "kind", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "resource_id", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "resource_type", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "resource_name", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_name", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "bucket", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "node_hostname", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "severity", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", alertsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "starts_at", alertsBuilder, computedRFC3339())
	capellaschema.AddAttr(dataAttrs, "ends_at", alertsBuilder, computedRFC3339())
	capellaschema.AddAttr(dataAttrs, "event_time", alertsBuilder, computedRFC3339())

	capellaschema.AddAttr(attrs, "data", alertsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	}, "ListAlertsResponse")

	return schema.Schema{
		MarkdownDescription: "The alerts data source lists the firing and resolved alerts in an organization, with one entry for each alert on each resource. Muted alerts are not included. By default it covers the last 48 hours; use from and to to read up to 31 days back.",
		Attributes:          attrs,
	}
}

func optionalStringSetOf(validators ...validator.String) *schema.SetAttribute {
	return &schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
			setvalidator.ValueStringsAre(validators...),
		},
	}
}

func optionalRFC3339() *schema.StringAttribute {
	return &schema.StringAttribute{
		Optional:   true,
		CustomType: timetypes.RFC3339Type{},
	}
}
