package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var freeTierClusterOnOffBuilder = capellaschema.NewSchemaBuilder("freeTierClusterOnOff")

func FreeTierClusterOnOffSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", freeTierClusterOnOffBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", freeTierClusterOnOffBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "cluster_id", freeTierClusterOnOffBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "state", freeTierClusterOnOffBuilder, stringAttribute([]string{required}))

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the On/Off state of a free tier operational cluster.",
		Attributes:          attrs,
	}
}
