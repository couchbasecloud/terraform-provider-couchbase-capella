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
	capellaschema.AddAttr(attrs, "cidr", allowlistBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "comment", allowlistBuilder, stringAttribute([]string{optional, computed, requiresReplace}))
	capellaschema.AddAttr(attrs, "expires_at", allowlistBuilder, stringAttribute([]string{optional, requiresReplace}))
	capellaschema.AddAttr(attrs, "audit", allowlistBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "Manages the Allowed IP addresses to connect to databases for Couchbase Capella.",
		Attributes:          attrs,
	}
}
