package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var privateEndpointsBuilder = capellaschema.NewSchemaBuilder("privateEndpoints")

func PrivateEndpointsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", privateEndpointsBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", privateEndpointsBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", privateEndpointsBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "endpoint_id", privateEndpointsBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "status", privateEndpointsBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "service_name", privateEndpointsBuilder, stringAttribute([]string{computed}))

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage private endpoints for an operational cluster. Private endpoints allow you to securely connect your Cloud Service Provider's private network (VPC/VNET) to your operational cluster without exposing traffic to the public internet.",
		Attributes:          attrs,
	}
}
