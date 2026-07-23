package datasources

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// Helper functions for common datasource attribute patterns

func requiredString() *schema.StringAttribute {
	return &schema.StringAttribute{
		Required: true,
	}
}

func requiredStringWithValidator() *schema.StringAttribute {
	return &schema.StringAttribute{
		Required:   true,
		Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
	}
}

func requiredUUIDString() *schema.StringAttribute {
	return &schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
			stringvalidator.RegexMatches(uuidRegex, "must be a valid UUID"),
		},
	}
}

func computedString() *schema.StringAttribute {
	return &schema.StringAttribute{
		Computed: true,
	}
}

func computedRFC3339() *schema.StringAttribute {
	return &schema.StringAttribute{
		Computed:   true,
		CustomType: timetypes.RFC3339Type{},
	}
}

func optionalString() *schema.StringAttribute {
	return &schema.StringAttribute{
		Optional: true,
	}
}

func requiredInt64() *schema.Int64Attribute {
	return &schema.Int64Attribute{
		Required: true,
	}
}

func computedInt64() *schema.Int64Attribute {
	return &schema.Int64Attribute{
		Computed: true,
	}
}

func optionalInt64() *schema.Int64Attribute {
	return &schema.Int64Attribute{
		Optional: true,
	}
}

func computedBool() *schema.BoolAttribute {
	return &schema.BoolAttribute{
		Computed: true,
	}
}

func optionalStringSet() *schema.SetAttribute {
	return &schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
	}
}

func computedStringSet() *schema.SetAttribute {
	return &schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
	}
}

// computedResourcesAttribute builds the resources schema block (buckets.scopes.collections)
// using Set types. This is reused by both access and privilege schemas.
func computedResourcesAttribute(builder *capellaschema.SchemaBuilder) *schema.SingleNestedAttribute {
	scopeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(scopeAttrs, "name", builder, computedString())
	capellaschema.AddAttr(scopeAttrs, "collections", builder, computedStringSet())

	bucketAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(bucketAttrs, "name", builder, computedString())
	capellaschema.AddAttr(bucketAttrs, "scopes", builder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: scopeAttrs,
		},
	})

	resourcesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(resourcesAttrs, "buckets", builder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: bucketAttrs,
		},
	})

	return &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: resourcesAttrs,
	}
}

// computedAccessAttribute builds the access schema block (privileges, resources.buckets.scopes.collections)
// using Set types consistent with the corresponding resource schemas.
func computedAccessAttribute(builder *capellaschema.SchemaBuilder) *schema.SetNestedAttribute {
	accessAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(accessAttrs, "privileges", builder, computedStringSet())
	capellaschema.AddAttr(accessAttrs, "resources", builder, computedResourcesAttribute(builder))

	return &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: accessAttrs,
		},
	}
}

func computedAudit() *schema.SingleNestedAttribute {
	tempBuilder := capellaschema.NewSchemaBuilder("audit")
	auditAttrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(auditAttrs, "created_at", tempBuilder, &schema.StringAttribute{
		Computed: true,
	}, "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "created_by", tempBuilder, &schema.StringAttribute{
		Computed: true,
	}, "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_at", tempBuilder, &schema.StringAttribute{
		Computed: true,
	}, "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_by", tempBuilder, &schema.StringAttribute{
		Computed: true,
	}, "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "version", tempBuilder, &schema.Int64Attribute{
		Computed: true,
	}, "CouchbaseAuditData")

	return &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	}
}
