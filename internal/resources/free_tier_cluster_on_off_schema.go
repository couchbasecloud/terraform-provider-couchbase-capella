package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var freeTierClusterOnOffBuilder = capellaschema.NewSchemaBuilder("freeTierClusterOnOff")

func FreeTierClusterOnOffSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", freeTierClusterOnOffBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", freeTierClusterOnOffBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", freeTierClusterOnOffBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "state", freeTierClusterOnOffBuilder, stringAttribute([]string{required}))

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the On/Off state of a free tier operational cluster.",
		Attributes:          attrs,
	}
}
