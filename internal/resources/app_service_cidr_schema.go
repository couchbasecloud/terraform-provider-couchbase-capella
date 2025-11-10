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
	capellaschema.AddAttr(attrs, "organization_id", appServiceCIDRBuilder, &schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "project_id", appServiceCIDRBuilder, &schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "cluster_id", appServiceCIDRBuilder, &schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "app_service_id", appServiceCIDRBuilder, &schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "cidr", appServiceCIDRBuilder, &schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "comment", appServiceCIDRBuilder, &schema.StringAttribute{
		Optional: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "expires_at", appServiceCIDRBuilder, &schema.StringAttribute{
		Optional: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	})
	capellaschema.AddAttr(attrs, "audit", appServiceCIDRBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "Manages the IP addresses allowed to connect to App Services in Couchbase Capella.",
		Attributes:          attrs,
	}
}
