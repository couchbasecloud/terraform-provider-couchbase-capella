package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var replicationsBuilder = capellaschema.NewSchemaBuilder("replications", "replicationSummary")

// ReplicationsSchema returns the schema for the replications data source.
func ReplicationsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", replicationsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", replicationsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", replicationsBuilder, optionalUUIDString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "source_cluster", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "target_cluster", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "direction", replicationsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "audit", replicationsBuilder, computedReplicationAuditAttribute())

	data := &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	}
	data.MarkdownDescription = capellaschema.ReplicationsDataDescription
	capellaschema.AddAttr(attrs, "data", replicationsBuilder, data)

	return schema.Schema{
		MarkdownDescription: "The data source to list the XDCR replications of a cluster, or of all the clusters in a " +
			"project when cluster_id is not set.",
		Attributes: attrs,
	}
}
