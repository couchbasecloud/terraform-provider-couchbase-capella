package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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

	attrs["scopes"] = schema.MapNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"collections": schema.MapNestedAttribute{
					Optional: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"access_control_function": stringAttribute([]string{optional, computed, useStateForUnknown}),
							"import_filter":           stringAttribute([]string{optional, computed, useStateForUnknown}),
						},
					},
				},
			},
		},
	}
	attrs["cors"] = schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"origin": &schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"login_origin": &schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"headers": &schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"max_age": &schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"disabled": &schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	attrs["oidc"] = schema.ListNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"issuer":         stringAttribute([]string{required}),
				"register":       boolAttribute(optional, computed),
				"client_id":      stringAttribute([]string{required}),
				"user_prefix":    stringAttribute([]string{optional, computed}),
				"discovery_url":  stringAttribute([]string{optional, computed}),
				"username_claim": stringAttribute([]string{optional, computed}),
				"roles_claim":    stringAttribute([]string{optional, computed}),
				"provider_id":    stringAttribute([]string{computed}),
				"is_default":     boolAttribute(computed),
			},
		},
	}
	attrs["require_resync"] = schema.MapNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"items": &schema.SetAttribute{
					Computed:    true,
					ElementType: types.StringType,
					PlanModifiers: []planmodifier.Set{
						setplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
	}
	attrs["state"] = &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	attrs["admin_url"] = &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	attrs["metrics_url"] = &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	attrs["public_url"] = &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage an App Endpoint configuration for a Couchbase Capella App Service.",
		Attributes:          attrs,
	}
}
