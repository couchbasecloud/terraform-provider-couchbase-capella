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

	capellaschema.AddAttr(attrs, "organization_id", privateEndpointServiceBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", privateEndpointServiceBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", privateEndpointServiceBuilder, requiredUUIDStringAttribute())

	capellaschema.AddAttr(attrs, "enabled", privateEndpointServiceBuilder, &schema.BoolAttribute{
		Required:      true,
		PlanModifiers: []planmodifier.Bool{custommodifier.BlockCreateWhenEnabledSetToFalse()},
	})

	capellaschema.AddAttr(attrs, "status", privateEndpointServiceBuilder, &schema.StringAttribute{
		Computed: true,
	})

	capellaschema.AddAttr(attrs, "service_name", privateEndpointServiceBuilder, &schema.StringAttribute{
		Computed: true,
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the private endpoint service for an operational cluster. " +
			"The private endpoint service must be enabled before you can create private endpoints to connect your Cloud Service Provider's private network (VPC/VNET) to your operational cluster. " +
			"This enables secure access to your cluster without exposing traffic to the public internet.\n\n" +
			"~> **Enablement failure handling:** If enablement terminally fails (`status` = `enableFailed`), the provider automatically issues a cleanup request to tear down the partially provisioned service and then removes the resource from Terraform state before the apply errors out. " +
			"This is intentional: the resource disappearing from state is expected, and the next `terraform apply` performs a clean re-create. " +
			"If the automatic cleanup cannot complete, the error will say so — contact Couchbase Capella Support to check for orphaned resources in your cloud account.",
		Attributes: attrs,
	}
}
