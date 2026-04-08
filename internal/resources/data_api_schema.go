package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var dataAPIBuilder = capellaschema.NewSchemaBuilder("dataAPI")

func DataAPISchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", dataAPIBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", dataAPIBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", dataAPIBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "enable_data_api", dataAPIBuilder, boolAttribute(required))
	capellaschema.AddAttr(attrs, "enable_network_peering", dataAPIBuilder, boolAttribute(required))
	capellaschema.AddAttr(attrs, "enabled", dataAPIBuilder, boolAttribute(computed))
	capellaschema.AddAttr(attrs, "state", dataAPIBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "enabled_for_network_peering", dataAPIBuilder, boolAttribute(computed))
	capellaschema.AddAttr(attrs, "state_for_network_peering", dataAPIBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "connection_string", dataAPIBuilder, stringAttribute([]string{computed}))

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the Data API configuration for a cluster.",
		Attributes:          attrs,
	}
}
