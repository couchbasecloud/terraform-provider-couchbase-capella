package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var outboundPrivateEndpointServiceVerificationBuilder = capellaschema.NewSchemaBuilder(
	"outboundPrivateEndpointServiceVerification", "OutboundPrivateEndpointServiceVerification",
)

// OutboundPrivateEndpointServiceVerificationSchema returns the schema for the
// outbound_private_endpoint_service_verification data source.
//
// Attribute descriptions come from the published OpenAPI spec through AddAttr. The route is
// x-internal until the feature ships, so the descriptions appear once it is public.
func OutboundPrivateEndpointServiceVerificationSchema() schema.Schema {
	b := outboundPrivateEndpointServiceVerificationBuilder
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", b, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", b, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", b, requiredUUIDString())
	capellaschema.AddAttr(attrs, "service_name", b, &schema.StringAttribute{
		Required:   true,
		Validators: []validator.String{stringvalidator.LengthBetween(1, 512)},
	}, "VerifyOutboundPrivateEndpointServiceRequest")

	cloudProviderAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cloudProviderAttrs, "type", b, computedString(), "OutboundPrivateEndpointCloudProvider")
	capellaschema.AddAttr(cloudProviderAttrs, "region", b, computedString(), "OutboundPrivateEndpointCloudProvider")
	capellaschema.AddAttr(attrs, "cloud_provider", b, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: cloudProviderAttrs,
	})

	capellaschema.AddAttr(attrs, "ready", b, computedBool())

	checkAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(checkAttrs, "name", b, computedString(), "OutboundPrivateEndpointServiceVerificationCheck")
	capellaschema.AddAttr(checkAttrs, "result", b, computedString(), "OutboundPrivateEndpointServiceVerificationCheck")
	capellaschema.AddAttr(checkAttrs, "message", b, computedString(), "OutboundPrivateEndpointServiceVerificationCheck")
	capellaschema.AddAttr(attrs, "checks", b, &schema.ListNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: checkAttrs},
	})

	awsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(awsAttrs, "availability_zones", b, computedStringList(), "AWSOutboundPrivateEndpointServiceVerification")
	capellaschema.AddAttr(awsAttrs, "acceptance_required", b, computedBool(), "AWSOutboundPrivateEndpointServiceVerification")
	capellaschema.AddAttr(awsAttrs, "service_owner", b, computedString(), "AWSOutboundPrivateEndpointServiceVerification")
	capellaschema.AddAttr(awsAttrs, "private_dns_name", b, computedString(), "AWSOutboundPrivateEndpointServiceVerification")
	capellaschema.AddAttr(awsAttrs, "private_dns_name_verification", b, computedString(), "AWSOutboundPrivateEndpointServiceVerification")

	providerConfigAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(providerConfigAttrs, "aws", b, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: awsAttrs,
	}, "OutboundPrivateEndpointServiceVerificationProviderConfig")
	capellaschema.AddAttr(attrs, "provider_config", b, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: providerConfigAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "Checks a customer's private link service from the cluster's cloud account and region before an " +
			"outbound private endpoint is created to it. Nothing is created. A failed check is reported in `checks` with " +
			"`ready` false rather than as an error, so a configuration can gate on `ready`. `provider_config` is set only " +
			"once the service is found, and only for the cluster's cloud. GCP and Azure have no fields yet.",
		Attributes: attrs,
	}
}
