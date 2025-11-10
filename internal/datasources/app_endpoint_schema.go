package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointBuilder = capellaschema.NewSchemaBuilder("appEndpoint")

func AppEndpointSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointBuilder, &schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	})
	capellaschema.AddAttr(attrs, "project_id", appEndpointBuilder, &schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	})
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointBuilder, &schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	})
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointBuilder, &schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	})
	capellaschema.AddAttr(attrs, "name", appEndpointBuilder, &schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	})

	capellaschema.AddAttr(attrs, "bucket", appEndpointBuilder, computedString())
	capellaschema.AddAttr(attrs, "user_xattr_key", appEndpointBuilder, computedString())
	capellaschema.AddAttr(attrs, "delta_sync_enabled", appEndpointBuilder, computedBool())

	collectionAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(collectionAttrs, "access_control_function", appEndpointBuilder, computedString())
	capellaschema.AddAttr(collectionAttrs, "import_filter", appEndpointBuilder, computedString())

	scopeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(scopeAttrs, "collections", appEndpointBuilder, &schema.MapNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: collectionAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "scopes", appEndpointBuilder, &schema.MapNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: scopeAttrs,
		},
	})

	corsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(corsAttrs, "origin", appEndpointBuilder, &schema.SetAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(corsAttrs, "login_origin", appEndpointBuilder, &schema.SetAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(corsAttrs, "headers", appEndpointBuilder, &schema.SetAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(corsAttrs, "max_age", appEndpointBuilder, computedInt64())
	capellaschema.AddAttr(corsAttrs, "disabled", appEndpointBuilder, computedBool())

	capellaschema.AddAttr(attrs, "cors", appEndpointBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: corsAttrs,
	})

	oidcAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(oidcAttrs, "issuer", appEndpointBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "register", appEndpointBuilder, computedBool())
	capellaschema.AddAttr(oidcAttrs, "client_id", appEndpointBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "user_prefix", appEndpointBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "discovery_url", appEndpointBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "username_claim", appEndpointBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "roles_claim", appEndpointBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "provider_id", appEndpointBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "is_default", appEndpointBuilder, computedBool())

	capellaschema.AddAttr(attrs, "oidc", appEndpointBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: oidcAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "state", appEndpointBuilder, computedString())

	requireResyncAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(requireResyncAttrs, "items", appEndpointBuilder, &schema.SetAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})

	capellaschema.AddAttr(attrs, "require_resync", appEndpointBuilder, &schema.MapNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: requireResyncAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "admin_url", appEndpointBuilder, computedString())
	capellaschema.AddAttr(attrs, "metrics_url", appEndpointBuilder, computedString())
	capellaschema.AddAttr(attrs, "public_url", appEndpointBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source retrieves a single App Endpoint configuration for an App Service.",
		Attributes:          attrs,
	}
}
