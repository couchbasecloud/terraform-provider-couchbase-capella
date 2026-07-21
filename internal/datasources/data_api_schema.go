package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var dataApiBuilder = capellaschema.NewSchemaBuilder("dataApi")

// DataApiSchema returns the schema for the DataApi data source.
func DataApiSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", dataApiBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", dataApiBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", dataApiBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "enable_data_api", dataApiBuilder, computedBool())
	capellaschema.AddAttr(attrs, "enable_network_peering", dataApiBuilder, computedBool())
	capellaschema.AddAttr(attrs, "state_for_data_api", dataApiBuilder, computedString())
	capellaschema.AddAttr(attrs, "state_for_network_peering", dataApiBuilder, computedString())
	capellaschema.AddAttr(attrs, "connection_string", dataApiBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the Data API and network peering status for an operational cluster.",
		Attributes:          attrs,
	}
}
