package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var dataAPIBuilder = capellaschema.NewSchemaBuilder("dataAPI")

// DataAPISchema returns the schema for the DataAPI data source.
func DataAPISchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", dataAPIBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", dataAPIBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", dataAPIBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "enabled", dataAPIBuilder, computedBool())
	capellaschema.AddAttr(attrs, "state", dataAPIBuilder, computedString())
	capellaschema.AddAttr(attrs, "enabled_for_network_peering", dataAPIBuilder, computedBool())
	capellaschema.AddAttr(attrs, "state_for_network_peering", dataAPIBuilder, computedString())
	capellaschema.AddAttr(attrs, "connection_string", dataAPIBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The Data API data source allows you to retrieve the Data API status for a cluster.",
		Attributes:          attrs,
	}
}
