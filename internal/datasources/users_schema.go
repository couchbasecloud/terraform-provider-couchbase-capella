package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var usersBuilder = capellaschema.NewSchemaBuilder("users")

// UsersSchema returns the schema for the Users data source.
func UsersSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", usersBuilder, requiredString())

	// Build data attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", usersBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", usersBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", usersBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "inactive", usersBuilder, computedBool())
	capellaschema.AddAttr(dataAttrs, "email", usersBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", usersBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_roles", usersBuilder, &schema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
	})
	capellaschema.AddAttr(dataAttrs, "last_login", usersBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "region", usersBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "time_zone", usersBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "enable_notifications", usersBuilder, computedBool())
	capellaschema.AddAttr(dataAttrs, "expires_at", usersBuilder, computedString())

	// Build resources attributes
	resourcesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(resourcesAttrs, "type", usersBuilder, &schema.StringAttribute{
		Optional: true,
		Computed: true,
	}, "Resource")
	capellaschema.AddAttr(resourcesAttrs, "id", usersBuilder, computedString(), "Resource")
	capellaschema.AddAttr(resourcesAttrs, "roles", usersBuilder, &schema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
	}, "Resource")

	capellaschema.AddAttr(dataAttrs, "resources", usersBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: resourcesAttrs,
		},
	})

	// Build audit attributes
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_at", usersBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "created_by", usersBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_at", usersBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_by", usersBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "version", usersBuilder, computedInt64(), "CouchbaseAuditData")

	capellaschema.AddAttr(dataAttrs, "audit", usersBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	})

	capellaschema.AddAttr(attrs, "data", usersBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve users in a Capella organization.",
		Attributes:          attrs,
	}
}
