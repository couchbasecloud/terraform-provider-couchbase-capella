package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var outboundPrivateEndpointsBuilder = capellaschema.NewSchemaBuilder("outboundPrivateEndpoints", "OutboundPrivateEndpoint")

// OutboundPrivateEndpointsSchema returns the schema for the outbound_private_endpoints data source.
// Each item matches the couchbase-capella_outbound_private_endpoint resource's attributes.
func OutboundPrivateEndpointsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", outboundPrivateEndpointsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", outboundPrivateEndpointsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", outboundPrivateEndpointsBuilder, requiredUUIDString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", outboundPrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", outboundPrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", outboundPrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", outboundPrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", outboundPrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "endpoint_type", outboundPrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "service_name", outboundPrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "hostname", outboundPrivateEndpointsBuilder, computedString())

	cloudProviderAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cloudProviderAttrs, "type", outboundPrivateEndpointsBuilder, computedString(), "OutboundPrivateEndpointCloudProvider")
	capellaschema.AddAttr(cloudProviderAttrs, "region", outboundPrivateEndpointsBuilder, computedString(), "OutboundPrivateEndpointCloudProvider")
	capellaschema.AddAttr(dataAttrs, "cloud_provider", outboundPrivateEndpointsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: cloudProviderAttrs,
	})

	capellaschema.AddAttr(dataAttrs, "provider_endpoint_id", outboundPrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", outboundPrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "failure_reason", outboundPrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "provider_status", outboundPrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "is_custom_hostname", outboundPrivateEndpointsBuilder, computedBool())
	capellaschema.AddAttr(dataAttrs, "private_dns_enabled", outboundPrivateEndpointsBuilder, computedBool())
	capellaschema.AddAttr(dataAttrs, "ip_addresses", outboundPrivateEndpointsBuilder, computedStringList())

	awsIpAddressAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(awsIpAddressAttrs, "ip_address", outboundPrivateEndpointsBuilder, computedString(), "AWSOutboundPrivateEndpointIpAddress")
	capellaschema.AddAttr(awsIpAddressAttrs, "availability_zone", outboundPrivateEndpointsBuilder, computedString(), "AWSOutboundPrivateEndpointIpAddress")

	awsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(awsAttrs, "regional_dns_name", outboundPrivateEndpointsBuilder, computedString(), "AWSOutboundPrivateEndpoint")
	capellaschema.AddAttr(awsAttrs, "ip_addresses", outboundPrivateEndpointsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: awsIpAddressAttrs,
		},
	}, "AWSOutboundPrivateEndpoint")
	capellaschema.AddAttr(awsAttrs, "availability_zones", outboundPrivateEndpointsBuilder, computedStringList(), "AWSOutboundPrivateEndpoint")

	providerConfigAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(providerConfigAttrs, "aws", outboundPrivateEndpointsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: awsAttrs,
	})
	capellaschema.AddAttr(dataAttrs, "provider_config", outboundPrivateEndpointsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: providerConfigAttrs,
	})

	capellaschema.AddAttr(dataAttrs, "audit", outboundPrivateEndpointsBuilder, computedAudit())

	capellaschema.AddAttr(attrs, "data", outboundPrivateEndpointsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	}, "ListOutboundPrivateEndpointsResponse")

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve every outbound private endpoint on a cluster, newest first. " +
			"Pending, failed and deleting endpoints are included.",
		Attributes: attrs,
	}
}
