package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Helper functions for common datasource attribute patterns

func requiredString() *schema.StringAttribute {
	return &schema.StringAttribute{
		Required: true,
	}
}

func computedString() *schema.StringAttribute {
	return &schema.StringAttribute{
		Computed: true,
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
