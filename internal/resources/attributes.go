package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	optional           = "optional"
	computed           = "computed"
	required           = "required"
	sensitive          = "sensitive"
	requiresReplace    = "requiresReplace"
	useStateForUnknown = "useStateForUnknown"
	deprecated         = "deprecated"
	deprecationMessage = "Remove this attribute's configuration as it no longer in use and the attribute will be removed in the next major version of the provider."
)

// SchemaAttribute is a type alias that encapsulates the allowed attribute types.
// It is used to define a set of types that can be used as attributes in the schema.
// This alias includes various attribute types such as StringAttribute, Int64Attribute, BoolAttribute, SetAttribute,
// Float64Attribute, NumberAttribute, and ListAttribute.
type SchemaAttribute interface {
	*schema.StringAttribute | *schema.Int64Attribute | *schema.BoolAttribute | *schema.SetAttribute |
		*schema.Float64Attribute | *schema.NumberAttribute | *schema.ListAttribute
}

// WithDescription sets the MarkdownDescription for the provided attribute.
// It accepts an attribute of type SchemaAttribute, and a description string.
// The function returns the modified attribute.
func WithDescription[T SchemaAttribute](attr T, description string) T {
	switch v := any(attr).(type) {
	case *schema.StringAttribute:
		v.MarkdownDescription = description
	case *schema.Int64Attribute:
		v.MarkdownDescription = description
	case *schema.BoolAttribute:
		v.MarkdownDescription = description
	case *schema.SetAttribute:
		v.MarkdownDescription = description
	case *schema.Float64Attribute:
		v.MarkdownDescription = description
	case *schema.NumberAttribute:
		v.MarkdownDescription = description
	case *schema.ListAttribute:
		v.MarkdownDescription = description
	}
	return attr
}

// stringAttribute is a variadic function which sets the requested fields
// in a string attribute to true and then returns the string attribute.
func stringAttribute(fields []string, validators ...validator.String) *schema.StringAttribute {
	attribute := schema.StringAttribute{}
	attribute.Validators = make([]validator.String, 0)

	attribute.Validators = append(attribute.Validators, validators...)

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
		case deprecated:
			attribute.DeprecationMessage = deprecationMessage
		}
	}
	return &attribute
}

// stringDefaultAttribute sets the default values for a string field and returns the string attribute.
func stringDefaultAttribute(defaultValue string, fields ...string) *schema.StringAttribute {
	attribute := stringAttribute(fields)
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

		case useStateForUnknown:
			var planModifiers = []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			}
			attribute.PlanModifiers = planModifiers
		}
	}

	return &attribute
}

// boolDefaultAttribute sets the default values for a boolean field and returns the bool attribute.
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

// int64DefaultAttribute sets the default values for an int field and returns the int64 attribute.
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

// float64DefaultAttribute sets the default values for an float field and returns the float64 attribute.
func float64DefaultAttribute(defaultValue float64, fields ...string) *schema.Float64Attribute {
	attribute := float64Attribute(fields...)
	attribute.Default = float64default.StaticFloat64(defaultValue)
	return attribute
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

// computedAuditAttribute returns a SingleNestedAttribute to
// represent couchbase audit data using terraform schema types.
func computedAuditAttribute() *schema.SingleNestedAttribute {
	return &schema.SingleNestedAttribute{
		Description: "Couchbase audit data.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the resource.",
			},
			"modified_at": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp when the resource was last modified.",
			},
			"modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who last modified the resource.",
			},
			"version": schema.Int64Attribute{
				Computed: true,
				Description: "The version of the document. " +
					"This value is incremented each time the resource is modified.",
			},
		},
	}
}
