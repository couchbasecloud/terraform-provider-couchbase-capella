package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var replicationsBuilder = capellaschema.NewSchemaBuilder("replication", "ReplicationSummary")

// ReplicationsSchema returns the schema for the replications data source.
func ReplicationsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", replicationsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", replicationsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", replicationsBuilder, requiredString())

	// Data items
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "source_cluster", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "target_cluster", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "direction", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "audit", replicationsBuilder, computedReplicationAudit())

	capellaschema.AddAttr(attrs, "data", replicationsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "Retrieves a list of all replications for a specific cluster.",
		Attributes:          attrs,
	}
}
