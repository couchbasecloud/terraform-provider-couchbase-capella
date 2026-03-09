package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// replicationsBuilder is the schema builder for the replications data source.
var replicationsBuilder = capellaschema.NewSchemaBuilder("replications", "ReplicationSummary")

// ReplicationsSchema returns the schema for the list replications data source.
func ReplicationsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	// Required hierarchical IDs
	capellaschema.AddAttr(attrs, "organization_id", replicationsBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", replicationsBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", replicationsBuilder, requiredStringWithValidator())

	// Audit attributes for data items
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_by", replicationsBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "created_at", replicationsBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "modified_by", replicationsBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "modified_at", replicationsBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "version", replicationsBuilder, computedInt64())

	// Data item attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "source_cluster", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "target_cluster", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "direction", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "audit", replicationsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	})

	// Data list attribute
	capellaschema.AddAttr(attrs, "data", replicationsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source retrieves a list of Couchbase Capella replications for a cluster.",
		Attributes:          attrs,
	}
}
