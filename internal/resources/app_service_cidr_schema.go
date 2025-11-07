package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceCIDRBuilder = capellaschema.NewSchemaBuilder("appServiceCIDR")

// AllowedCIDRsSchema returns the schema for the App Service allowed CIDRs resource.
func AllowedCIDRsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", appServiceCIDRBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "organization_id", appServiceCIDRBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", appServiceCIDRBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", appServiceCIDRBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "app_service_id", appServiceCIDRBuilder, stringAttribute([]string{required, requiresReplace}))

	attrs["cidr"] = schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The trusted CIDR block to allow the database connections from.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	attrs["comment"] = schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "A short description of the allowed CIDR block.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	attrs["expires_at"] = schema.StringAttribute{
		Optional: true,
		MarkdownDescription: "An RFC3339 timestamp determining when the allowed CIDR block will expire. " +
			"If this field is omitted then the allowed CIDR block is permanent. It will never automatically expire.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	capellaschema.AddAttr(attrs, "audit", appServiceCIDRBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "Manages the IP addresses allowed to connect to App Services in Couchbase Capella.",
		Attributes:          attrs,
	}
}
