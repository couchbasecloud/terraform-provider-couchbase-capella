package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
	capellaschema.AddAttr(attrs, "indexes", gsiMonitorBuilder, &schema.SetAttribute{
		Required:    true,
		ElementType: types.StringType,
	})

	return schema.Schema{
		MarkdownDescription: "The data source to monitor the build progress of a GSI index.",
		Attributes:          attrs,
	}
}
