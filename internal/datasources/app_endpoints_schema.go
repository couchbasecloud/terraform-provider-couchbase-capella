package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointsBuilder = capellaschema.NewSchemaBuilder("appEndpoints")

func AppEndpointsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointsBuilder, &schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	})
	capellaschema.AddAttr(attrs, "project_id", appEndpointsBuilder, &schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	})
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointsBuilder, &schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	})
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointsBuilder, &schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	})

	collectionAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(collectionAttrs, "access_control_function", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(collectionAttrs, "import_filter", appEndpointsBuilder, computedString())

	scopeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(scopeAttrs, "collections", appEndpointsBuilder, &schema.MapNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: collectionAttrs,
		},
	})

	corsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(corsAttrs, "origin", appEndpointsBuilder, &schema.SetAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(corsAttrs, "login_origin", appEndpointsBuilder, &schema.SetAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(corsAttrs, "headers", appEndpointsBuilder, &schema.SetAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})
	capellaschema.AddAttr(corsAttrs, "max_age", appEndpointsBuilder, computedInt64())
	capellaschema.AddAttr(corsAttrs, "disabled", appEndpointsBuilder, computedBool())

	oidcAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(oidcAttrs, "issuer", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "register", appEndpointsBuilder, computedBool())
	capellaschema.AddAttr(oidcAttrs, "client_id", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "user_prefix", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "discovery_url", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "username_claim", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "roles_claim", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "provider_id", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(oidcAttrs, "is_default", appEndpointsBuilder, computedBool())

	requireResyncAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(requireResyncAttrs, "items", appEndpointsBuilder, &schema.SetAttribute{
		Computed:    true,
		ElementType: types.StringType,
	})

	appEndpointAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(appEndpointAttrs, "bucket", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(appEndpointAttrs, "name", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(appEndpointAttrs, "user_xattr_key", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(appEndpointAttrs, "delta_sync_enabled", appEndpointsBuilder, computedBool())
	capellaschema.AddAttr(appEndpointAttrs, "scopes", appEndpointsBuilder, &schema.MapNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: scopeAttrs,
		},
	})
	capellaschema.AddAttr(appEndpointAttrs, "cors", appEndpointsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: corsAttrs,
	})
	capellaschema.AddAttr(appEndpointAttrs, "oidc", appEndpointsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: oidcAttrs,
		},
	})
	capellaschema.AddAttr(appEndpointAttrs, "state", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(appEndpointAttrs, "require_resync", appEndpointsBuilder, &schema.MapNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: requireResyncAttrs,
		},
	})
	capellaschema.AddAttr(appEndpointAttrs, "admin_url", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(appEndpointAttrs, "metrics_url", appEndpointsBuilder, computedString())
	capellaschema.AddAttr(appEndpointAttrs, "public_url", appEndpointsBuilder, computedString())

	capellaschema.AddAttr(attrs, "app_endpoints", appEndpointsBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: appEndpointAttrs,
		},
	})

	filterAttrs := make(map[string]schema.Attribute)
	filterAttrs["name"] = schema.StringAttribute{
		MarkdownDescription: "The name of the attribute to filter.",
		Optional:            true,
		Validators: []validator.String{
			stringvalidator.OneOf("name"),
		},
	}
	filterAttrs["values"] = schema.SetAttribute{
		MarkdownDescription: "List of values to match against.",
		Optional:            true,
		ElementType:         types.StringType,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	}

	return schema.Schema{
		MarkdownDescription: "The data source retrieves App Endpoint configurations for an App Service.",
		Attributes:          attrs,
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				MarkdownDescription: "Filter criteria for App Endpoints.  Only filtering by App Endpoint name is supported.",
				Attributes:          filterAttrs,
			},
		},
	}
}
