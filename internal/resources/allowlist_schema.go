package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var allowlistBuilder = capellaschema.NewSchemaBuilder("allowlist")

func AllowlistsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", allowlistBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "organization_id", allowlistBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", allowlistBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", allowlistBuilder, stringAttribute([]string{required, requiresReplace}))

	attrs["cidr"] = schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The trusted CIDR to allow the database connections from.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	attrs["comment"] = schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "A short description of the allowed CIDR.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	attrs["expires_at"] = schema.StringAttribute{
		Optional: true,
		MarkdownDescription: "An RFC3339 timestamp determining when the allowed CIDR will expire. " +
			"If this field is omitted then the allowed CIDR is permanent and will never automatically expire.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	capellaschema.AddAttr(attrs, "audit", allowlistBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "Manages the Allowed IP addresses to connect to databases for Couchbase Capella.",
		Attributes:          attrs,
	}
}
