package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceCIDRBuilder = capellaschema.NewSchemaBuilder("appServiceCIDR", "AppServiceAllowedCidr")

// AllowedCIDRsSchema returns the schema for the App Service allowed CIDRs resource.
func AllowedCIDRsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", appServiceCIDRBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "organization_id", appServiceCIDRBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServiceCIDRBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceCIDRBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceCIDRBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cidr", appServiceCIDRBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "comment", appServiceCIDRBuilder, stringAttribute([]string{optional, requiresReplace}))
	capellaschema.AddAttr(attrs, "expires_at", appServiceCIDRBuilder, stringAttribute([]string{optional, requiresReplace}))
	capellaschema.AddAttr(attrs, "audit", appServiceCIDRBuilder, computedAuditAttribute())
	// An omitted access_mode stays null in the state. The API uses admin for it, so null and admin do not replace.
	accessModeAttr := stringAttribute([]string{optional})
	accessModeAttr.PlanModifiers = append(accessModeAttr.PlanModifiers, requiresReplaceOnKnownChange(api.AllowedCIDRAccessModeAdmin))
	capellaschema.AddAttr(attrs, "access_mode", appServiceCIDRBuilder, accessModeAttr)

	return schema.Schema{
		MarkdownDescription: "Manages the IP addresses allowed to connect to App Services in Couchbase Capella.",
		Attributes:          attrs,
	}
}
