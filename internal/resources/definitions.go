package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	optional        = "optional"
	computed        = "computed"
	required        = "required"
	sensitive       = "sensitive"
	requiresReplace = "requiresReplace"
)

// stringAttribute is a variadic function which sets the requested fields
// in a string attribute to true and then returns the string attribute.
func stringAttribute(fields ...string) *schema.StringAttribute {
	attribute := schema.StringAttribute{}

	for _, field := range fields {
		switch field {
		case required:
			attribute.Required = true
		case optional:
			attribute.Optional = true
		case computed:
			attribute.Computed = true
		case sensitive:
			attribute.Sensitive = true
		case requiresReplace:
			var planModifiers = []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}
			attribute.PlanModifiers = planModifiers
		}
	}
	return &attribute
}

// boolAttribute is a variadic function which sets the requested fields
// in a bool attribute to true and then returns the string attribute.
func boolAttribute(fields ...string) *schema.BoolAttribute {
	attribute := schema.BoolAttribute{}

	for _, field := range fields {
		switch field {
		case required:
			attribute.Required = true
		case optional:
			attribute.Optional = true
		case computed:
			attribute.Computed = true
		case sensitive:
			attribute.Sensitive = true
		case requiresReplace:
			var planModifiers = []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			}
			attribute.PlanModifiers = planModifiers
		}
	}
	return &attribute
}

// int64Attribute is a variadic function which sets the requested fields
// in an Int64 attribute to true and then returns the string attribute.
func int64Attribute(fields ...string) *schema.Int64Attribute {
	attribute := schema.Int64Attribute{}

	for _, field := range fields {
		switch field {
		case required:
			attribute.Required = true
		case optional:
			attribute.Optional = true
		case computed:
			attribute.Computed = true
		case sensitive:
			attribute.Sensitive = true
		}
	}
	return &attribute
}

// float64Attribute is a variadic function which sets the requested fields
// in a float64 attribute to true and then returns the string attribute.
func float64Attribute(fields ...string) *schema.Float64Attribute {
	attribute := schema.Float64Attribute{}

	for _, field := range fields {
		switch field {
		case required:
			attribute.Required = true
		case optional:
			attribute.Optional = true
		case computed:
			attribute.Computed = true
		case sensitive:
			attribute.Sensitive = true
		}
	}
	return &attribute
}

// stringListAttribute returns a Terraform string list schema attribute
// which is configured to be of type string.
func stringListAttribute(fields ...string) *schema.ListAttribute {
	attribute := schema.ListAttribute{
		ElementType: types.StringType,
	}

	for _, field := range fields {
		switch field {
		case required:
			attribute.Required = true
		case optional:
			attribute.Optional = true
		case computed:
			attribute.Computed = true
		case sensitive:
			attribute.Sensitive = true
		case requiresReplace:
			var planModifiers = []planmodifier.List{
				listplanmodifier.RequiresReplace(),
			}
			attribute.PlanModifiers = planModifiers
		}
	}
	return &attribute
}

// computedAuditAttribute retuns a SingleNestedAttribute to
// represent couchbase audit data using terraform schema types.
func computedAuditAttribute() *schema.SingleNestedAttribute {
	return &schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"created_at":  stringAttribute(computed),
			"created_by":  stringAttribute(computed),
			"modified_at": stringAttribute(computed),
			"modified_by": stringAttribute(computed),
			"version":     int64Attribute(computed),
		},
	}
}
