package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var apiKeysBuilder = capellaschema.NewSchemaBuilder("apiKeys")

// ApiKeysSchema returns the schema for the ApiKeys data source.
func ApiKeysSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", apiKeysBuilder, requiredString())

	// Build data attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", apiKeysBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", apiKeysBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", apiKeysBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "description", apiKeysBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "expiry", apiKeysBuilder, &schema.Float64Attribute{
		Computed: true,
	})
	capellaschema.AddAttr(dataAttrs, "allowed_cidrs", apiKeysBuilder, &schema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
	})
	capellaschema.AddAttr(dataAttrs, "organization_roles", apiKeysBuilder, &schema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
	})

	// Build resources attributes
	resourcesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(resourcesAttrs, "id", apiKeysBuilder, computedString())
	capellaschema.AddAttr(resourcesAttrs, "roles", apiKeysBuilder, &schema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
	})
	capellaschema.AddAttr(resourcesAttrs, "type", apiKeysBuilder, computedString())

	capellaschema.AddAttr(dataAttrs, "resources", apiKeysBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: resourcesAttrs,
		},
	})

	// Build audit attributes
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_at", apiKeysBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "created_by", apiKeysBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "modified_at", apiKeysBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "modified_by", apiKeysBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "version", apiKeysBuilder, computedInt64())

	capellaschema.AddAttr(dataAttrs, "audit", apiKeysBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	})

	capellaschema.AddAttr(attrs, "data", apiKeysBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve API keys in an organization. API keys are used to authenticate and authorize access to Capella resources and services.",
		Attributes:          attrs,
	}
}
