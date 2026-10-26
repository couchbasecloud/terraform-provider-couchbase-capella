package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAdminUsersBuilder = capellaschema.NewSchemaBuilder("appServiceAdminUsers", "appServiceAdminUser")

// AppServiceAdminUsersSchema returns the schema for the App Service admin users data source.
func AppServiceAdminUsersSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceAdminUsersBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", appServiceAdminUsersBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAdminUsersBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAdminUsersBuilder, requiredStringWithValidator())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", appServiceAdminUsersBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", appServiceAdminUsersBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "access_all_endpoints", appServiceAdminUsersBuilder, computedBool())
	capellaschema.AddAttr(dataAttrs, "endpoints", appServiceAdminUsersBuilder, computedStringSet())
	capellaschema.AddAttr(dataAttrs, "audit", appServiceAdminUsersBuilder, computedAudit())

	capellaschema.AddAttr(attrs, "data", appServiceAdminUsersBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	filterAttrs := make(map[string]schema.Attribute)
	filterAttrs["name"] = schema.StringAttribute{
		MarkdownDescription: "The name of the attribute to filter.",
		Optional:            true,
		Validators: []validator.String{
			stringvalidator.OneOf("name"),
		},
	}
	filterAttrs["values"] = schema.SetAttribute{
		MarkdownDescription: "List of values to match against.",
		Optional:            true,
		ElementType:         types.StringType,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	}

	return schema.Schema{
		MarkdownDescription: "Retrieves the admin users of a Capella App Service.",
		Attributes:          attrs,
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				MarkdownDescription: "Filter criteria for admin users. Only filtering by admin user name is supported.",
				Attributes:          filterAttrs,
			},
		},
	}
}
