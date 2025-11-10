package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var networkPeerBuilder = capellaschema.NewSchemaBuilder("networkPeer")

func NetworkPeerSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", networkPeerBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "organization_id", networkPeerBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", networkPeerBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", networkPeerBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "name", networkPeerBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "provider_type", networkPeerBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "audit", networkPeerBuilder, computedAuditAttribute())
	capellaschema.AddAttr(attrs, "commands", networkPeerBuilder, &schema.SetAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})

	awsConfigAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(awsConfigAttrs, "account_id", networkPeerBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(awsConfigAttrs, "vpc_id", networkPeerBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(awsConfigAttrs, "region", networkPeerBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(awsConfigAttrs, "cidr", networkPeerBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(awsConfigAttrs, "provider_id", networkPeerBuilder, stringAttribute([]string{computed}))

	gcpConfigAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(gcpConfigAttrs, "cidr", networkPeerBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(gcpConfigAttrs, "network_name", networkPeerBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(gcpConfigAttrs, "project_id", networkPeerBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(gcpConfigAttrs, "service_account", networkPeerBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(gcpConfigAttrs, "provider_id", networkPeerBuilder, stringAttribute([]string{computed}))

	azureConfigAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(azureConfigAttrs, "tenant_id", networkPeerBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(azureConfigAttrs, "cidr", networkPeerBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(azureConfigAttrs, "resource_group", networkPeerBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(azureConfigAttrs, "subscription_id", networkPeerBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(azureConfigAttrs, "vnet_id", networkPeerBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(azureConfigAttrs, "provider_id", networkPeerBuilder, stringAttribute([]string{computed}))

	providerConfigAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(providerConfigAttrs, "aws_config", networkPeerBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: awsConfigAttrs,
	})
	capellaschema.AddAttr(providerConfigAttrs, "gcp_config", networkPeerBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: gcpConfigAttrs,
	})
	capellaschema.AddAttr(providerConfigAttrs, "azure_config", networkPeerBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: azureConfigAttrs,
	})

	capellaschema.AddAttr(attrs, "provider_config", networkPeerBuilder, &schema.SingleNestedAttribute{
		Required: true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
		Attributes: providerConfigAttrs,
	})

	statusAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(statusAttrs, "reasoning", networkPeerBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(statusAttrs, "state", networkPeerBuilder, stringAttribute([]string{computed}))

	capellaschema.AddAttr(attrs, "status", networkPeerBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: statusAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage network peering for an operational cluster.",
		Attributes:          attrs,
	}
}
