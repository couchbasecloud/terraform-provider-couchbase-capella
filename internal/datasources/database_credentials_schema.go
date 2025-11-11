package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var databaseCredentialsBuilder = capellaschema.NewSchemaBuilder("databaseCredentials")

// DatabaseCredentialsSchema returns the schema for the DatabaseCredentials data source.
func DatabaseCredentialsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", databaseCredentialsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", databaseCredentialsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", databaseCredentialsBuilder, requiredString())

	// Build data attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", databaseCredentialsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", databaseCredentialsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", databaseCredentialsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", databaseCredentialsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", databaseCredentialsBuilder, computedString())

	// Build audit attributes
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_at", databaseCredentialsBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "created_by", databaseCredentialsBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_at", databaseCredentialsBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_by", databaseCredentialsBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "version", databaseCredentialsBuilder, computedInt64(), "CouchbaseAuditData")

	capellaschema.AddAttr(dataAttrs, "audit", databaseCredentialsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	})

	// Build bucket attributes for access.resources.buckets
	bucketAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(bucketAttrs, "name", databaseCredentialsBuilder, &schema.StringAttribute{
		Required: true,
	})
	capellaschema.AddAttr(bucketAttrs, "scopes", databaseCredentialsBuilder, &schema.ListNestedAttribute{
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
	capellaschema.AddAttr(resourcesAttrs, "buckets", databaseCredentialsBuilder, &schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: bucketAttrs,
		},
	})

	// Build access attributes
	accessAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(accessAttrs, "privileges", databaseCredentialsBuilder, &schema.ListAttribute{
		Required:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(accessAttrs, "resources", databaseCredentialsBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: resourcesAttrs,
	})

	capellaschema.AddAttr(dataAttrs, "access", databaseCredentialsBuilder, &schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: accessAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "data", databaseCredentialsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve database credentials for a cluster. Database credentials provide programmatic and application-level access to data on a database.",
		Attributes:          attrs,
	}
}
