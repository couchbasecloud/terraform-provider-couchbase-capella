package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var dataAPIBuilder = capellaschema.NewSchemaBuilder("dataAPI", "UpdateDataAPINetworkPeering")

func DataAPISchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", dataAPIBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", dataAPIBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", dataAPIBuilder, requiredUUIDStringAttribute())

	capellaschema.AddAttr(attrs, "enable_data_api", dataAPIBuilder, boolAttribute(required), "UpdateDataAPINetworkPeering")
	capellaschema.AddAttr(attrs, "enable_network_peering", dataAPIBuilder, boolAttribute(required), "UpdateDataAPINetworkPeering")

	capellaschema.AddAttr(attrs, "state", dataAPIBuilder, stringAttribute([]string{computed}), "GetDataAPIStatusResponse")
	capellaschema.AddAttr(attrs, "state_for_network_peering", dataAPIBuilder, stringAttribute([]string{computed}), "GetDataAPIStatusResponse")
	capellaschema.AddAttr(attrs, "connection_string", dataAPIBuilder, stringAttribute([]string{computed}), "GetDataAPIStatusResponse")

	return schema.Schema{
		MarkdownDescription: "Manages Data API on a Couchbase Capella cluster. Enables or disables Data API and network peering for the cluster. Enabling Data API is an asynchronous operation.",
		Attributes:          attrs,
	}
}
