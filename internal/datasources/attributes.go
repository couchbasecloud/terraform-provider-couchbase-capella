package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// computedStringAttribute returns a Terraform schema attribute
// which is configured to be computed.
func computedStringAttribute() *schema.StringAttribute {
	return &schema.StringAttribute{
		Computed: true,
	}
}

// requiredStringAttribute returns a Terraform schema attribute
// which is configured to be required.
func requiredStringAttribute() *schema.StringAttribute {
	return &schema.StringAttribute{
		Required: true,
	}
}

// computedBoolAttribute returns a Terraform schema attribute
// which is configured to be computed.
func computedBoolAttribute() *schema.BoolAttribute {
	return &schema.BoolAttribute{
		Computed: true,
	}
}

// computedBoolAttribute returns a Terraform schema attribute
// which is configured to be computed.
func computedInt64Attribute() *schema.Int64Attribute {
	return &schema.Int64Attribute{
		Computed: true,
	}
}

// computedListAttribute returns a Terraform list schema attribute
// which is configured to be computed and of type string.
func computedListAttribute() *schema.ListAttribute {
	return &schema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
	}
}

// computedAuditAttribute retuns a SingleNestedAttribute to
// represent couchbase audit data using terraform schema types.
func computedAuditAttribute() *schema.SingleNestedAttribute {
	return &schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"created_at":  computedStringAttribute(),
			"created_by":  computedStringAttribute(),
			"modified_at": computedStringAttribute(),
			"modified_by": computedStringAttribute(),
			"version":     computedInt64Attribute(),
		},
	}
}
