package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var gsiMonitorBuilder = capellaschema.NewSchemaBuilder("gsiMonitor")

// GsiMonitorSchema returns the schema for the GsiMonitor data source.
func GsiMonitorSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", gsiMonitorBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", gsiMonitorBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", gsiMonitorBuilder, requiredString())
	capellaschema.AddAttr(attrs, "bucket_name", gsiMonitorBuilder, requiredString())
	capellaschema.AddAttr(attrs, "scope_name", gsiMonitorBuilder, requiredString())
	capellaschema.AddAttr(attrs, "collection_name", gsiMonitorBuilder, requiredString())
	capellaschema.AddAttr(attrs, "index_name", gsiMonitorBuilder, requiredString())
	capellaschema.AddAttr(attrs, "status", gsiMonitorBuilder, computedString())
	capellaschema.AddAttr(attrs, "progress", gsiMonitorBuilder, computedInt64())

	return schema.Schema{
		MarkdownDescription: "The data source to monitor the build progress of a GSI index.",
		Attributes:          attrs,
	}
}
