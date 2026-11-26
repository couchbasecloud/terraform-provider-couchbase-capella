package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	outboundpe "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/outbound_private_endpoint"
)

// OutboundPrivateEndpointServiceVerification is the datasource model of a check of a private
// link service before an outbound private endpoint is created to it.
type OutboundPrivateEndpointServiceVerification struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`
	ServiceName    types.String `tfsdk:"service_name"`

	CloudProvider  *OutboundPrivateEndpointCloudProvider                     `tfsdk:"cloud_provider"`
	Ready          types.Bool                                                `tfsdk:"ready"`
	Checks         []OutboundPrivateEndpointServiceVerificationCheck         `tfsdk:"checks"`
	ProviderConfig *OutboundPrivateEndpointServiceVerificationProviderConfig `tfsdk:"provider_config"`
}

// OutboundPrivateEndpointServiceVerificationCheck is the result of one check.
type OutboundPrivateEndpointServiceVerificationCheck struct {
	Name    types.String `tfsdk:"name"`
	Result  types.String `tfsdk:"result"`
	Message types.String `tfsdk:"message"`
}

// OutboundPrivateEndpointServiceVerificationProviderConfig is what the cloud reports about the
// service. Only the cluster's cloud is set, and GCP and Azure have no fields yet.
type OutboundPrivateEndpointServiceVerificationProviderConfig struct {
	AWS *AWSOutboundPrivateEndpointServiceVerification `tfsdk:"aws"`
}

// AWSOutboundPrivateEndpointServiceVerification is what AWS reports about an endpoint service.
type AWSOutboundPrivateEndpointServiceVerification struct {
	AvailabilityZones          []types.String `tfsdk:"availability_zones"`
	AcceptanceRequired         types.Bool     `tfsdk:"acceptance_required"`
	ServiceOwner               types.String   `tfsdk:"service_owner"`
	PrivateDnsName             types.String   `tfsdk:"private_dns_name"`
	PrivateDnsNameVerification types.String   `tfsdk:"private_dns_name_verification"`
}

// SetResult fills the computed attributes from a verification response.
func (v *OutboundPrivateEndpointServiceVerification) SetResult(r *outboundpe.ServiceVerification) {
	v.CloudProvider = &OutboundPrivateEndpointCloudProvider{
		Type:   types.StringValue(r.CloudProvider.Type),
		Region: types.StringValue(r.CloudProvider.Region),
	}
	v.Ready = types.BoolValue(r.Ready)

	v.Checks = make([]OutboundPrivateEndpointServiceVerificationCheck, 0, len(r.Checks))
	for _, c := range r.Checks {
		v.Checks = append(v.Checks, OutboundPrivateEndpointServiceVerificationCheck{
			Name:    types.StringValue(c.Name),
			Result:  types.StringValue(c.Result),
			Message: types.StringPointerValue(c.Message),
		})
	}

	v.ProviderConfig = nil
	if r.ProviderConfig == nil {
		return
	}
	v.ProviderConfig = &OutboundPrivateEndpointServiceVerificationProviderConfig{}
	if aws := r.ProviderConfig.AWS; aws != nil {
		zones := make([]types.String, 0, len(aws.AvailabilityZones))
		for _, z := range aws.AvailabilityZones {
			zones = append(zones, types.StringValue(z))
		}
		v.ProviderConfig.AWS = &AWSOutboundPrivateEndpointServiceVerification{
			AvailabilityZones:          zones,
			AcceptanceRequired:         types.BoolValue(aws.AcceptanceRequired),
			ServiceOwner:               types.StringValue(aws.ServiceOwner),
			PrivateDnsName:             types.StringPointerValue(aws.PrivateDnsName),
			PrivateDnsNameVerification: types.StringValue(aws.PrivateDnsNameVerification),
		}
	}
}
