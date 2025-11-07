package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var clusterOnOffBuilder = capellaschema.NewSchemaBuilder("clusterOnOff")

func ClusterOnOffOnDemandSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", clusterOnOffBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", clusterOnOffBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", clusterOnOffBuilder, stringAttribute([]string{required, requiresReplace}))

	attrs["state"] = WithDescription(stringAttribute([]string{required}), "State of the cluster. It can be `on` and `off`.")
	attrs["turn_on_linked_app_service"] = WithDescription(boolAttribute(optional, computed), " Whether to turn on the linked App Service when the cluster is turned on.")

	return schema.Schema{
		MarkdownDescription: "Manages the On/Off state of an operational cluster. This resource allows you to turn your cluster on or off.",
		Attributes:          attrs,
	}
}
