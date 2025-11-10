package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var apiKeyBuilder = capellaschema.NewSchemaBuilder("apiKey", "APIKey")

func ApiKeySchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", apiKeyBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", apiKeyBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "name", apiKeyBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "description", apiKeyBuilder, stringDefaultAttribute("", optional, computed, requiresReplace, useStateForUnknown))
	capellaschema.AddAttr(attrs, "expiry", apiKeyBuilder, float64DefaultAttribute(180, optional, computed, requiresReplace, useStateForUnknown))

	// Secret field - description from RotateAPIKeyRequest schema
	capellaschema.AddAttr(attrs, "secret", apiKeyBuilder, stringAttribute([]string{optional, computed, sensitive}), "RotateAPIKeyRequest")

	// Token field - description from CreateAPIKeyResponse and Token schema
	capellaschema.AddAttr(attrs, "token", apiKeyBuilder, stringAttribute([]string{computed, sensitive}), "CreateAPIKeyResponse", "Token")

	capellaschema.AddAttr(attrs, "audit", apiKeyBuilder, computedAuditAttribute())
	capellaschema.AddAttr(attrs, "allowed_cidrs", apiKeyBuilder, &schema.SetAttribute{
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
	})
	capellaschema.AddAttr(attrs, "organization_roles", apiKeyBuilder, stringSetAttribute(required, requiresReplace))

	resourceAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(resourceAttrs, "id", apiKeyBuilder, stringAttribute([]string{required}), "Resource")
	capellaschema.AddAttr(resourceAttrs, "roles", apiKeyBuilder, stringSetAttribute(required), "Resource")
	capellaschema.AddAttr(resourceAttrs, "type", apiKeyBuilder, stringDefaultAttribute("project", optional, computed), "Resource")

	capellaschema.AddAttr(attrs, "resources", apiKeyBuilder, &schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: resourceAttrs,
		},
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.RequiresReplace(),
		},
	})

	// Rotate field - Terraform-specific field to trigger key rotation
	// NOTE: This field does not exist in the OpenAPI spec - it's a Terraform provider
	// implementation detail for managing the key rotation lifecycle. Therefore, we
	// cannot use AddAttr and must manually define the description.
	attrs["rotate"] = &schema.NumberAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "\n - Set this value in incremental order from the previously set rotate value (starting from 1) to trigger rotation of the API key. When rotated, a new token and secret are generated.",
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to create and manage API keys in Capella. API keys are used to authenticate and authorize access to Capella resources and services.",
		Attributes:          attrs,
	}
}
