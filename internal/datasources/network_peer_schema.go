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

	commandAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(commandAttrs, "command", networkPeerBuilder, computedString())

	commandsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(commandsAttrs, "aws", networkPeerBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: commandAttrs,
	})
	capellaschema.AddAttr(commandsAttrs, "gcp", networkPeerBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: commandAttrs,
	})
	capellaschema.AddAttr(commandsAttrs, "azure", networkPeerBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: commandAttrs,
	})

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", networkPeerBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", networkPeerBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "provider_type", networkPeerBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "provider_config", networkPeerBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "audit", networkPeerBuilder, computedAudit())
	capellaschema.AddAttr(dataAttrs, "status", networkPeerBuilder, &schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"state":     schema.StringAttribute{Computed: true},
			"reasoning": schema.StringAttribute{Computed: true},
		},
	})
	capellaschema.AddAttr(dataAttrs, "commands", networkPeerBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: commandsAttrs,
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
