package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var freeTierAppServiceBuilder = capellaschema.NewSchemaBuilder("freeTierAppService")

func FreeTierAppServiceSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", freeTierAppServiceBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", freeTierAppServiceBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", freeTierAppServiceBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "cluster_id", freeTierAppServiceBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "name", freeTierAppServiceBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "nodes", freeTierAppServiceBuilder, int64Attribute(computed, useStateForUnknown))
	capellaschema.AddAttr(attrs, "cloud_provider", freeTierAppServiceBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "current_state", freeTierAppServiceBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "version", freeTierAppServiceBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "audit", freeTierAppServiceBuilder, computedAuditAttribute())
	capellaschema.AddAttr(attrs, "plan", freeTierAppServiceBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "etag", freeTierAppServiceBuilder, stringAttribute([]string{computed}))

	attrs["description"] = &schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
	}

	attrs["compute"] = schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"cpu": int64Attribute(computed, useStateForUnknown),
			"ram": int64Attribute(computed, useStateForUnknown),
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage free tier App Services associated with a free tier operational cluster.",
		Attributes:          attrs,
	}
}
