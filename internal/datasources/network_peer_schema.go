package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var networkPeerBuilder = capellaschema.NewSchemaBuilder("networkPeer")

func NetworkPeerSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", networkPeerBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", networkPeerBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", networkPeerBuilder, requiredString())

	awsConfigAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(awsConfigAttrs, "account_id", networkPeerBuilder, computedString())
	capellaschema.AddAttr(awsConfigAttrs, "vpc_id", networkPeerBuilder, computedString())
	capellaschema.AddAttr(awsConfigAttrs, "region", networkPeerBuilder, computedString())
	capellaschema.AddAttr(awsConfigAttrs, "cidr", networkPeerBuilder, computedString())
	capellaschema.AddAttr(awsConfigAttrs, "provider_id", networkPeerBuilder, computedString())

	gcpConfigAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(gcpConfigAttrs, "cidr", networkPeerBuilder, computedString())
	capellaschema.AddAttr(gcpConfigAttrs, "network_name", networkPeerBuilder, computedString())
	capellaschema.AddAttr(gcpConfigAttrs, "project_id", networkPeerBuilder, computedString())
	capellaschema.AddAttr(gcpConfigAttrs, "service_account", networkPeerBuilder, computedString())
	capellaschema.AddAttr(gcpConfigAttrs, "provider_id", networkPeerBuilder, computedString())

	azureConfigAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(azureConfigAttrs, "tenant_id", networkPeerBuilder, computedString())
	capellaschema.AddAttr(azureConfigAttrs, "cidr", networkPeerBuilder, computedString())
	capellaschema.AddAttr(azureConfigAttrs, "resource_group", networkPeerBuilder, computedString())
	capellaschema.AddAttr(azureConfigAttrs, "subscription_id", networkPeerBuilder, computedString())
	capellaschema.AddAttr(azureConfigAttrs, "vnet_id", networkPeerBuilder, computedString())
	capellaschema.AddAttr(azureConfigAttrs, "provider_id", networkPeerBuilder, computedString())

	providerConfigAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(providerConfigAttrs, "aws_config", networkPeerBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: awsConfigAttrs,
	})
	capellaschema.AddAttr(providerConfigAttrs, "gcp_config", networkPeerBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: gcpConfigAttrs,
	})
	capellaschema.AddAttr(providerConfigAttrs, "azure_config", networkPeerBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: azureConfigAttrs,
	})

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", networkPeerBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", networkPeerBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "provider_config", networkPeerBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: providerConfigAttrs,
	})
	capellaschema.AddAttr(dataAttrs, "audit", networkPeerBuilder, computedAudit())
	capellaschema.AddAttr(dataAttrs, "status", networkPeerBuilder, &schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"state":     schema.StringAttribute{Computed: true},
			"reasoning": schema.StringAttribute{Computed: true},
		},
	})

	capellaschema.AddAttr(attrs, "data", networkPeerBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The network peers data source retrieves all network peers that have been created for a specific cluster.",
		Attributes:          attrs,
	}
}
