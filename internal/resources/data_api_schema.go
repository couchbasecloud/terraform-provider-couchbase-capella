package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var dataApiBuilder = capellaschema.NewSchemaBuilder("dataApi")

func DataApiSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", dataApiBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", dataApiBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", dataApiBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "enable_data_api", dataApiBuilder, boolAttribute(required))
	capellaschema.AddAttr(attrs, "enable_network_peering", dataApiBuilder, boolAttribute(optional, computed, useStateForUnknown))
	capellaschema.AddAttr(attrs, "state_for_data_api", dataApiBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "state_for_network_peering", dataApiBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "connection_string", dataApiBuilder, stringAttribute([]string{computed, useStateForUnknown}))

	return schema.Schema{
		MarkdownDescription: "This resource allows you to enable or disable the Data API and network peering for an operational cluster.",
		Attributes:          attrs,
	}
}
