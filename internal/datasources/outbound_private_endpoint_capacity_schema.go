package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var outboundPrivateEndpointCapacityBuilder = capellaschema.NewSchemaBuilder("outboundPrivateEndpointCapacity", "OutboundPrivateEndpointCapacity")

// OutboundPrivateEndpointCapacitySchema returns the schema for the outbound_private_endpoint_capacity data source.
//
// The capacity route is x-internal until the feature ships, so AddAttr finds no descriptions
// for these attributes until the spec is public.
func OutboundPrivateEndpointCapacitySchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", outboundPrivateEndpointCapacityBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", outboundPrivateEndpointCapacityBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", outboundPrivateEndpointCapacityBuilder, requiredUUIDString())

	usageAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(usageAttrs, "used", outboundPrivateEndpointCapacityBuilder, computedInt64(), "OutboundPrivateEndpointCapacityUsage")
	capellaschema.AddAttr(usageAttrs, "limit", outboundPrivateEndpointCapacityBuilder, computedInt64(), "OutboundPrivateEndpointCapacityUsage")

	capellaschema.AddAttr(attrs, "cluster_endpoint_limit", outboundPrivateEndpointCapacityBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: usageAttrs,
	}, "OutboundPrivateEndpointCapacity")

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve how many outbound private endpoints a cluster has against its limit. " +
			"Every endpoint counts, including pending, failed and deleting ones. The cluster must be on.",
		Attributes: attrs,
	}
}
