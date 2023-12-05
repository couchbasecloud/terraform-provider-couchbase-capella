package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	optional           = "optional"
	computed           = "computed"
	required           = "required"
	sensitive          = "sensitive"
	requiresReplace    = "requiresReplace"
	useStateForUnknown = "useStateForUnknown"
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
			attribute.PlanModifiers = append(attribute.PlanModifiers, planModifiers...)
		case useStateForUnknown:
			var planModifiers = []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}
			attribute.PlanModifiers = append(attribute.PlanModifiers, planModifiers...)
		}
	}
	return &attribute
}

// stringDefaultAttribute sets the default values for a string field and returns the string attribute
func stringDefaultAttribute(defaultValue string, fields ...string) *schema.StringAttribute {
	attribute := stringAttribute(fields...)
	attribute.Default = stringdefault.StaticString(defaultValue)
	return attribute
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

// boolDefaultAttribute sets the default values for a boolean field and returns the bool attribute
func boolDefaultAttribute(defaultValue bool, fields ...string) *schema.BoolAttribute {
	attribute := boolAttribute(fields...)
	attribute.Default = booldefault.StaticBool(defaultValue)
	return attribute
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
		case requiresReplace:
			var planModifiers = []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			}
			attribute.PlanModifiers = append(attribute.PlanModifiers, planModifiers...)
		case useStateForUnknown:
			var planModifiers = []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			}
			attribute.PlanModifiers = append(attribute.PlanModifiers, planModifiers...)
		}
	}
	return &attribute
}

// int64DefaultAttribute sets the default values for an int field and returns the int64 attribute
func int64DefaultAttribute(defaultValue int64, fields ...string) *schema.Int64Attribute {
	attribute := int64Attribute(fields...)
	attribute.Default = int64default.StaticInt64(defaultValue)
	return attribute
}

// numberAttribute is a variadic function which sets the requested fields
// in an number attribute to true and then returns the string attribute.
func numberAttribute(fields ...string) *schema.NumberAttribute {
	attribute := schema.NumberAttribute{}

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
		case requiresReplace:
			var planModifiers = []planmodifier.Float64{
				float64planmodifier.RequiresReplace(),
			}
			attribute.PlanModifiers = append(attribute.PlanModifiers, planModifiers...)
		case useStateForUnknown:
			var planModifiers = []planmodifier.Float64{
				float64planmodifier.UseStateForUnknown(),
			}
			attribute.PlanModifiers = append(attribute.PlanModifiers, planModifiers...)
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

// stringSetAttribute returns a Terraform string set schema attribute
// which is configured to be of type string.
func stringSetAttribute(fields ...string) *schema.SetAttribute {
	attribute := schema.SetAttribute{
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
			var planModifiers = []planmodifier.Set{
				setplanmodifier.RequiresReplace(),
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
