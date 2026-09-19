package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	custommodifier "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources/custom_plan_modifiers"
	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServicePrivateEndpointServiceBuilder = capellaschema.NewSchemaBuilder("appServicePrivateEndpointService")

func AppServicePrivateEndpointServiceSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServicePrivateEndpointServiceBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServicePrivateEndpointServiceBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServicePrivateEndpointServiceBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServicePrivateEndpointServiceBuilder, requiredUUIDStringAttribute())

	capellaschema.AddAttr(attrs, "enabled", appServicePrivateEndpointServiceBuilder, &schema.BoolAttribute{
		Required:      true,
		PlanModifiers: []planmodifier.Bool{custommodifier.BlockCreateWhenEnabledSetToFalse()},
	})

	capellaschema.AddAttr(attrs, "state", appServicePrivateEndpointServiceBuilder, &schema.StringAttribute{
		Computed: true,
	}, "GetAppServicePrivateEndpointStateResponse")

	capellaschema.AddAttr(attrs, "target_state", appServicePrivateEndpointServiceBuilder, &schema.StringAttribute{
		Computed: true,
	}, "GetAppServicePrivateEndpointStateResponse")

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the private endpoint service for an App Service (Sync Gateway). " +
			"The private endpoint service must be enabled before you can create private endpoints to connect your Cloud Service Provider's private network (VPC/VNET) to your App Service. " +
			"This enables secure access to your App Service without exposing traffic to the public internet.\n\n" +
			"~> **Enablement failure handling:** If enablement terminally fails (`state` = `failed`, `target_state` = `enabled`), the provider automatically issues a cleanup request to tear down the partially provisioned service and then removes the resource from Terraform state before the apply errors out. " +
			"This is intentional: the resource disappearing from state is expected, and the next `terraform apply` performs a clean re-create. " +
			"If the automatic cleanup cannot complete, the error will say so — contact Couchbase Capella Support to check for orphaned resources in your cloud account.\n\n" +
			"~> **AWS-only as of this writing:** Couchbase Capella's private endpoints for App Services are AWS-only; Azure/GCP support is not confirmed GA. This resource's schema does not restrict cloud provider, but it has only been validated against an AWS-backed App Service.",
		Attributes: attrs,
	}
}
