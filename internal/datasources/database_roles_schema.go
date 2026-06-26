package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var databaseRolesBuilder = capellaschema.NewSchemaBuilder("databaseRoles")

// DatabaseRolesSchema returns the schema for the DatabaseRoles data source.
func DatabaseRolesSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", databaseRolesBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", databaseRolesBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", databaseRolesBuilder, requiredString())

	// Build data item attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", databaseRolesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", databaseRolesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "description", databaseRolesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", databaseRolesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", databaseRolesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", databaseRolesBuilder, computedString())

	// Build audit attributes
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_at", databaseRolesBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "created_by", databaseRolesBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_at", databaseRolesBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_by", databaseRolesBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "version", databaseRolesBuilder, computedInt64(), "CouchbaseAuditData")

	capellaschema.AddAttr(dataAttrs, "audit", databaseRolesBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	})

	// Build bucket attributes for access.resources.buckets
	bucketAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(bucketAttrs, "name", databaseRolesBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(bucketAttrs, "scopes", databaseRolesBuilder, &schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Required: true,
				},
				"collections": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},
			},
		},
	})

	// Build resources attributes for access
	resourcesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(resourcesAttrs, "buckets", databaseRolesBuilder, &schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: bucketAttrs,
		},
	})

	// Build access attributes
	accessAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(accessAttrs, "privileges", databaseRolesBuilder, &schema.ListAttribute{
		Required:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(accessAttrs, "resources", databaseRolesBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: resourcesAttrs,
	})

	capellaschema.AddAttr(dataAttrs, "access", databaseRolesBuilder, &schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: accessAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "data", databaseRolesBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to list database user roles for a cluster. Database roles define reusable sets of privileges that can be assigned to database credentials.",
		Attributes:          attrs,
	}
}
