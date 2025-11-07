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
	capellaschema.AddAttr(attrs, "organization_id", allowlistBuilder, &schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "project_id", allowlistBuilder, &schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "cluster_id", allowlistBuilder, &schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "cidr", allowlistBuilder, &schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "comment", allowlistBuilder, &schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "expires_at", allowlistBuilder, &schema.StringAttribute{
		Optional: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "audit", allowlistBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "Manages the Allowed IP addresses to connect to databases for Couchbase Capella.",
		Attributes:          attrs,
	}
}
