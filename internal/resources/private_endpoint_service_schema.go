package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	custommodifier "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources/custom_plan_modifiers"
	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var privateEndpointServiceBuilder = capellaschema.NewSchemaBuilder("privateEndpointService")

func PrivateEndpointServiceSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", privateEndpointServiceBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", privateEndpointServiceBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", privateEndpointServiceBuilder, stringAttribute([]string{required, requiresReplace}))

	capellaschema.AddAttr(attrs, "enabled", privateEndpointServiceBuilder, &schema.BoolAttribute{
		Required:      true,
		PlanModifiers: []planmodifier.Bool{custommodifier.BlockCreateWhenEnabledSetToFalse()},
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the private endpoint service for an operational cluster. The private endpoint service must be enabled before you can create private endpoints to connect your Cloud Service Provider's private network (VPC/VNET) to your operational cluster. This enables secure access to your cluster without exposing traffic to the public internet.",
		Attributes:          attrs,
	}
}
