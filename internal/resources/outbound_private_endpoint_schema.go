package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var outboundPrivateEndpointBuilder = capellaschema.NewSchemaBuilder("outboundPrivateEndpoint")

// Attribute descriptions come from the published OpenAPI spec through AddAttr. The outbound private
// endpoint routes are x-internal until the feature ships, so the descriptions appear once they are public.
func OutboundPrivateEndpointSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", outboundPrivateEndpointBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", outboundPrivateEndpointBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", outboundPrivateEndpointBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", outboundPrivateEndpointBuilder, requiredUUIDStringAttribute())

	capellaschema.AddAttr(attrs, "name", outboundPrivateEndpointBuilder, stringAttribute([]string{required}, validator.String(stringvalidator.LengthBetween(2, 128))))
	capellaschema.AddAttr(attrs, "endpoint_type", outboundPrivateEndpointBuilder, stringDefaultAttributeReplaceOnChange("interface", optional, computed))
	capellaschema.AddAttr(attrs, "service_name", outboundPrivateEndpointBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthBetween(1, 512))))
	capellaschema.AddAttr(attrs, "hostname", outboundPrivateEndpointBuilder, &schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.LengthBetween(1, 254),
		},
		PlanModifiers: []planmodifier.String{
			outboundPrivateEndpointHostnameModifier{},
		},
	})

	cloudProviderAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cloudProviderAttrs, "type", outboundPrivateEndpointBuilder, stringAttribute([]string{computed}), "OutboundPrivateEndpointCloudProvider")
	capellaschema.AddAttr(cloudProviderAttrs, "region", outboundPrivateEndpointBuilder, stringAttribute([]string{computed}), "OutboundPrivateEndpointCloudProvider")
	capellaschema.AddAttr(attrs, "cloud_provider", outboundPrivateEndpointBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: cloudProviderAttrs,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
	})

	capellaschema.AddAttr(attrs, "provider_endpoint_id", outboundPrivateEndpointBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "status", outboundPrivateEndpointBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "failure_reason", outboundPrivateEndpointBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "provider_status", outboundPrivateEndpointBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "is_custom_hostname", outboundPrivateEndpointBuilder, boolAttribute(computed, useStateForUnknown))
	capellaschema.AddAttr(attrs, "private_dns_enabled", outboundPrivateEndpointBuilder, boolAttribute(computed))
	capellaschema.AddAttr(attrs, "ip_addresses", outboundPrivateEndpointBuilder, stringListAttribute(computed))

	awsIpAddressAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(awsIpAddressAttrs, "ip_address", outboundPrivateEndpointBuilder, stringAttribute([]string{computed}), "AWSOutboundPrivateEndpointIpAddress")
	capellaschema.AddAttr(awsIpAddressAttrs, "availability_zone", outboundPrivateEndpointBuilder, stringAttribute([]string{computed}), "AWSOutboundPrivateEndpointIpAddress")

	awsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(awsAttrs, "regional_dns_name", outboundPrivateEndpointBuilder, stringAttribute([]string{computed}), "AWSOutboundPrivateEndpoint")
	capellaschema.AddAttr(awsAttrs, "ip_addresses", outboundPrivateEndpointBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: awsIpAddressAttrs,
		},
	}, "AWSOutboundPrivateEndpoint")
	capellaschema.AddAttr(awsAttrs, "availability_zones", outboundPrivateEndpointBuilder, stringListAttribute(computed), "AWSOutboundPrivateEndpoint")

	providerConfigAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(providerConfigAttrs, "aws", outboundPrivateEndpointBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: awsAttrs,
	})
	capellaschema.AddAttr(attrs, "provider_config", outboundPrivateEndpointBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: providerConfigAttrs,
	})

	capellaschema.AddAttr(attrs, "audit", outboundPrivateEndpointBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "Manages an outbound private endpoint: an endpoint in a Capella cluster's network that connects to " +
			"a private link service the customer owns, such as a Cloud Native Gateway in front of a self-managed cluster. " +
			"Use its `hostname` as the target of an XDCR replication. Create returns once the connection request reaches " +
			"the service. The customer then accepts it in their own cloud, and `status` moves from `pendingAcceptance` to `linked`.",
		Attributes: attrs,
	}
}

// outboundPrivateEndpointHostnameModifier plans hostname, which the customer may set or leave to
// the cloud or Capella. Changing a set hostname replaces the endpoint, and so does removing it from
// the config, since the endpoint then needs an assigned name. Left unset on an endpoint whose name
// was assigned, it keeps the name from state rather than reading as unknown on every plan.
type outboundPrivateEndpointHostnameModifier struct{}

func (outboundPrivateEndpointHostnameModifier) Description(context.Context) string {
	return "Replaces the endpoint when a set hostname changes or is removed, and keeps an assigned hostname from state."
}

func (m outboundPrivateEndpointHostnameModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (outboundPrivateEndpointHostnameModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Nothing to compare on create or destroy.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	if !req.ConfigValue.IsNull() {
		if !req.ConfigValue.IsUnknown() && !req.ConfigValue.Equal(req.StateValue) {
			resp.RequiresReplace = true
		}
		return
	}

	var isCustom types.Bool
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("is_custom_hostname"), &isCustom)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Terraform only replaces for an attribute whose planned value differs from state, and the
	// replacement's hostname is assigned after apply.
	if isCustom.ValueBool() {
		resp.PlanValue = types.StringUnknown()
		resp.RequiresReplace = true
		return
	}

	resp.PlanValue = req.StateValue
}
