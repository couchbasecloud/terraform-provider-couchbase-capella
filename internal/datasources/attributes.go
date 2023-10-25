package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	// computedStringAttribute returns a Terraform schema attribute
	// which is configured to be computed.
	computedStringAttribute = schema.StringAttribute{
		Computed: true,
	}

	// requiredStringAttribute returns a Terraform schema attribute
	// which is configured to be required.
	requiredStringAttribute = schema.StringAttribute{
		Required: true,
	}

	// computedBoolAttribute returns a Terraform schema attribute
	// which is configured to be computed.
	computedBoolAttribute = schema.BoolAttribute{
		Computed: true,
	}

	// computedBoolAttribute returns a Terraform schema attribute
	// which is configured to be computed.
	computedInt64Attribute = schema.Int64Attribute{
		Computed: true,
	}

	// computedBoolAttribute returns a Terraform schema attribute
	// which is configured to be computed.
	computedFloat64Attribute = schema.Float64Attribute{
		Computed: true,
	}

	// computedListAttribute returns a Terraform list schema attribute
	// which is configured to be computed and of type string.
	computedListAttribute = schema.ListAttribute{
		ElementType: types.StringType,
		Computed:    true,
	}

	// computedCursorAttribute returns a Terraform nested list
	// schema attribute which represents pagination metadata.
	computedCursorAttribute = schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"pages": schema.SingleNestedAttribute{
					Computed: true,
					Attributes: map[string]schema.Attribute{
						"page":       computedInt64Attribute,
						"next":       computedInt64Attribute,
						"previous":   computedInt64Attribute,
						"last":       computedInt64Attribute,
						"perPage":    computedInt64Attribute,
						"totalItems": computedInt64Attribute,
					},
				},
				"hrefs": schema.SingleNestedAttribute{
					Computed: true,
					Attributes: map[string]schema.Attribute{
						"first":    computedStringAttribute,
						"last":     computedStringAttribute,
						"previous": computedStringAttribute,
						"next":     computedStringAttribute,
					},
				},
			},
		},
	}

	// computedAuditAttribute retuns a SingleNestedAttribute to
	// represent couchbase audit data using terraform schema types.
	computedAuditAttribute = schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"created_at":  computedStringAttribute,
			"created_by":  computedStringAttribute,
			"modified_at": computedStringAttribute,
			"modified_by": computedStringAttribute,
			"version":     computedInt64Attribute,
		},
	}
)
