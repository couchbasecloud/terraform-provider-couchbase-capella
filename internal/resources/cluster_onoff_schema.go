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
	capellaschema.AddAttr(attrs, "state", clusterOnOffBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "turn_on_linked_app_service", clusterOnOffBuilder, boolAttribute(optional, computed))

	return schema.Schema{
		MarkdownDescription: "Manages the On/Off state of an operational cluster. This resource allows you to turn your cluster on or off.",
		Attributes:          attrs,
	}
}
