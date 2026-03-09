package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceOnOffBuilder = capellaschema.NewSchemaBuilder("appServiceOnOff")

func AppServiceOnOffOnDemandSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceOnOffBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServiceOnOffBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceOnOffBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceOnOffBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "state", appServiceOnOffBuilder, stringAttribute([]string{required}))

	return schema.Schema{
		MarkdownDescription: "Manages the on-demand state of an App Service. This resource is used to turn the App Service on or off on-demand.",
		Attributes:          attrs,
	}
}
