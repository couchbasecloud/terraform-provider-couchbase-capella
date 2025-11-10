package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func optionalBool() *schema.BoolAttribute {
	return &schema.BoolAttribute{
		Optional: true,
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
	return &schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"modified_at": schema.StringAttribute{
				Computed: true,
			},
			"modified_by": schema.StringAttribute{
				Computed: true,
			},
			"version": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}
