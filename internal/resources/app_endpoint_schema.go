package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointBuilder = capellaschema.NewSchemaBuilder("appEndpoint")

// AppEndpointSchema defines the schema for the app endpoint resource.
func AppEndpointSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", appEndpointBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "bucket", appEndpointBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "name", appEndpointBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "user_xattr_key", appEndpointBuilder, stringAttribute([]string{optional, computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "delta_sync_enabled", appEndpointBuilder, boolAttribute(optional, computed, useStateForUnknown))

	collectionAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(collectionAttrs, "access_control_function", appEndpointBuilder, stringAttribute([]string{optional, computed, useStateForUnknown}))
	capellaschema.AddAttr(collectionAttrs, "import_filter", appEndpointBuilder, stringAttribute([]string{optional, computed, useStateForUnknown}))

	scopeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(scopeAttrs, "collections", appEndpointBuilder, &schema.MapNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: collectionAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "scopes", appEndpointBuilder, &schema.MapNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: scopeAttrs,
		},
	})

	corsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(corsAttrs, "origin", appEndpointBuilder, &schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(corsAttrs, "login_origin", appEndpointBuilder, &schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(corsAttrs, "headers", appEndpointBuilder, &schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(corsAttrs, "max_age", appEndpointBuilder, int64Attribute(optional, computed, useStateForUnknown))
	capellaschema.AddAttr(corsAttrs, "disabled", appEndpointBuilder, &schema.BoolAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Bool{
			boolplanmodifier.UseStateForUnknown(),
		},
	})

	capellaschema.AddAttr(attrs, "cors", appEndpointBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: corsAttrs,
	})

	oidcAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(oidcAttrs, "issuer", appEndpointBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(oidcAttrs, "register", appEndpointBuilder, boolAttribute(optional, computed))
	capellaschema.AddAttr(oidcAttrs, "client_id", appEndpointBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(oidcAttrs, "user_prefix", appEndpointBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(oidcAttrs, "discovery_url", appEndpointBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(oidcAttrs, "username_claim", appEndpointBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(oidcAttrs, "roles_claim", appEndpointBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(oidcAttrs, "provider_id", appEndpointBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(oidcAttrs, "is_default", appEndpointBuilder, boolAttribute(computed))

	capellaschema.AddAttr(attrs, "oidc", appEndpointBuilder, &schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: oidcAttrs,
		},
	})

	requireResyncItemAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(requireResyncItemAttrs, "items", appEndpointBuilder, &schema.SetAttribute{
		Computed:    true,
		ElementType: types.StringType,
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
	})

	capellaschema.AddAttr(attrs, "require_resync", appEndpointBuilder, &schema.MapNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: requireResyncItemAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "state", appEndpointBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "admin_url", appEndpointBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "metrics_url", appEndpointBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "public_url", appEndpointBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage an App Endpoint configuration for a Couchbase Capella App Service.",
		Attributes:          attrs,
	}
}
