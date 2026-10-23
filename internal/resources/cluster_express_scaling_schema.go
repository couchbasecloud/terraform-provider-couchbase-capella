package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var clusterExpressScalingBuilder = capellaschema.NewSchemaBuilder("clusterExpressScaling")

func ClusterExpressScalingSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", clusterExpressScalingBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", clusterExpressScalingBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", clusterExpressScalingBuilder, requiredUUIDStringAttribute())

	expressScaling := boolAttribute(required)
	expressScaling.MarkdownDescription = capellaschema.ClusterExpressScalingDescription
	capellaschema.AddAttr(attrs, "express_scaling", clusterExpressScalingBuilder, expressScaling)

	expressScalingState := stringAttribute([]string{computed})
	expressScalingState.MarkdownDescription = capellaschema.ExpressScalingStateDescription
	capellaschema.AddAttr(attrs, "express_scaling_state", clusterExpressScalingBuilder, expressScalingState)

	return schema.Schema{
		MarkdownDescription: capellaschema.ClusterExpressScalingResourceDescription,
		Attributes:          attrs,
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(context.Background(), timeouts.Opts{
				Create:            true,
				Update:            true,
				CreateDescription: "How long to wait for the change to finish when the resource is created. Defaults to `60m`.",
				UpdateDescription: "How long to wait for the change to finish when the resource is updated. Defaults to `60m`.",
			}),
		},
	}
}
