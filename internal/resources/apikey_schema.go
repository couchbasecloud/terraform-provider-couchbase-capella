package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var apiKeyBuilder = capellaschema.NewSchemaBuilder("apiKey")

func ApiKeySchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", apiKeyBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "organization_id", apiKeyBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "name", apiKeyBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "description", apiKeyBuilder, stringDefaultAttribute("", optional, computed, requiresReplace, useStateForUnknown))
	capellaschema.AddAttr(attrs, "expiry", apiKeyBuilder, float64DefaultAttribute(180, optional, computed, requiresReplace, useStateForUnknown))
	capellaschema.AddAttr(attrs, "secret", apiKeyBuilder, stringAttribute([]string{optional, computed, sensitive}))
	capellaschema.AddAttr(attrs, "token", apiKeyBuilder, stringAttribute([]string{computed, sensitive}))
	capellaschema.AddAttr(attrs, "audit", apiKeyBuilder, computedAuditAttribute())

	attrs["allowed_cidrs"] = &schema.SetAttribute{
		Optional:    true,
		Computed:    true,
		ElementType: types.StringType,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
			setplanmodifier.RequiresReplace(),
		},
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
		Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{types.StringValue("0.0.0.0/0")})),
	}

	attrs["organization_roles"] = stringSetAttribute(required, requiresReplace)

	resourceAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(resourceAttrs, "id", apiKeyBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(resourceAttrs, "roles", apiKeyBuilder, stringSetAttribute(required))
	capellaschema.AddAttr(resourceAttrs, "type", apiKeyBuilder, stringDefaultAttribute("project", optional, computed))

	attrs["resources"] = &schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: resourceAttrs,
		},
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.RequiresReplace(),
		},
	}

	attrs["rotate"] = &schema.NumberAttribute{
		Optional: true,
		Computed: true,
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to create and manage API keys in Capella. API keys are used to authenticate and authorize access to Capella resources and services.",
		Attributes:          attrs,
	}
}
