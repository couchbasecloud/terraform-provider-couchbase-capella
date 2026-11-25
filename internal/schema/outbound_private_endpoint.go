package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	outboundpe "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/outbound_private_endpoint"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// OutboundPrivateEndpoint is the Terraform state of an outbound private endpoint: an endpoint
// in the cluster's network that connects to a private link service the customer owns.
type OutboundPrivateEndpoint struct {
	Id             types.String `tfsdk:"id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`

	Name         types.String `tfsdk:"name"`
	EndpointType types.String `tfsdk:"endpoint_type"`
	ServiceName  types.String `tfsdk:"service_name"`
	Hostname     types.String `tfsdk:"hostname"`

	CloudProvider      types.Object `tfsdk:"cloud_provider"`
	ProviderEndpointId types.String `tfsdk:"provider_endpoint_id"`
	Status             types.String `tfsdk:"status"`
	FailureReason      types.String `tfsdk:"failure_reason"`
	ProviderStatus     types.String `tfsdk:"provider_status"`
	IsCustomHostname   types.Bool   `tfsdk:"is_custom_hostname"`
	PrivateDnsEnabled  types.Bool   `tfsdk:"private_dns_enabled"`
	IpAddresses        types.List   `tfsdk:"ip_addresses"`
	ProviderConfig     types.Object `tfsdk:"provider_config"`
	Audit              types.Object `tfsdk:"audit"`
}

// OutboundPrivateEndpointCloudProvider is the cluster's cloud provider and region.
type OutboundPrivateEndpointCloudProvider struct {
	Type   types.String `tfsdk:"type"`
	Region types.String `tfsdk:"region"`
}

func (OutboundPrivateEndpointCloudProvider) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":   types.StringType,
		"region": types.StringType,
	}
}

// OutboundPrivateEndpointProviderConfig holds the cloud-specific details. Only the cluster's
// cloud is set, and GCP and Azure have no fields yet.
type OutboundPrivateEndpointProviderConfig struct {
	AWS types.Object `tfsdk:"aws"`
}

func (OutboundPrivateEndpointProviderConfig) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"aws": types.ObjectType{AttrTypes: OutboundPrivateEndpointAWSConfig{}.AttributeTypes()},
	}
}

// OutboundPrivateEndpointAWSConfig is the AWS-specific detail of an endpoint.
type OutboundPrivateEndpointAWSConfig struct {
	RegionalDnsName   types.String `tfsdk:"regional_dns_name"`
	IpAddresses       types.List   `tfsdk:"ip_addresses"`
	AvailabilityZones types.List   `tfsdk:"availability_zones"`
}

func (OutboundPrivateEndpointAWSConfig) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"regional_dns_name":  types.StringType,
		"ip_addresses":       types.ListType{ElemType: types.ObjectType{AttrTypes: OutboundPrivateEndpointAWSIpAddress{}.AttributeTypes()}},
		"availability_zones": types.ListType{ElemType: types.StringType},
	}
}

// OutboundPrivateEndpointAWSIpAddress is one of the endpoint's addresses with the zone it is in.
type OutboundPrivateEndpointAWSIpAddress struct {
	IpAddress        types.String `tfsdk:"ip_address"`
	AvailabilityZone types.String `tfsdk:"availability_zone"`
}

func (OutboundPrivateEndpointAWSIpAddress) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ip_address":        types.StringType,
		"availability_zone": types.StringType,
	}
}

// Validate returns the IDs the endpoint's routes need. After a terraform import only id is
// set, in the form id=<uuid>,cluster_id=<uuid>,project_id=<uuid>,organization_id=<uuid>.
func (o *OutboundPrivateEndpoint) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: o.OrganizationId,
		ProjectId:      o.ProjectId,
		ClusterId:      o.ClusterId,
		Id:             o.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %w", err)
	}

	return IDs, nil
}

// NewOutboundPrivateEndpoint builds the Terraform state from the API's endpoint. The resource and
// the outbound_private_endpoints data source both use it, so they map every field the same way.
func NewOutboundPrivateEndpoint(
	ctx context.Context,
	endpoint *outboundpe.GetOutboundPrivateEndpointResponse,
	organizationId, projectId, clusterId string,
) (*OutboundPrivateEndpoint, error) {
	audit := NewCouchbaseAuditData(endpoint.Audit)
	auditObject, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("%w: %v", errors.ErrUnableToConvertAuditData, diags.Errors())
	}

	cloudProvider, diags := types.ObjectValueFrom(ctx, OutboundPrivateEndpointCloudProvider{}.AttributeTypes(), OutboundPrivateEndpointCloudProvider{
		Type:   types.StringValue(endpoint.CloudProvider.Type),
		Region: types.StringValue(endpoint.CloudProvider.Region),
	})
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert cloud provider: %v", diags.Errors())
	}

	ipAddresses, diags := types.ListValueFrom(ctx, types.StringType, nonNilStrings(endpoint.IpAddresses))
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert ip addresses: %v", diags.Errors())
	}

	providerConfig, err := newOutboundPrivateEndpointProviderConfig(ctx, endpoint.ProviderConfig)
	if err != nil {
		return nil, err
	}

	return &OutboundPrivateEndpoint{
		Id:                 types.StringValue(endpoint.Id),
		OrganizationId:     types.StringValue(organizationId),
		ProjectId:          types.StringValue(projectId),
		ClusterId:          types.StringValue(clusterId),
		Name:               types.StringValue(endpoint.Name),
		EndpointType:       types.StringValue(endpoint.EndpointType),
		ServiceName:        types.StringValue(endpoint.ServiceName),
		Hostname:           types.StringPointerValue(endpoint.Hostname),
		CloudProvider:      cloudProvider,
		ProviderEndpointId: types.StringPointerValue(endpoint.ProviderEndpointId),
		Status:             types.StringValue(endpoint.Status),
		FailureReason:      types.StringPointerValue(endpoint.FailureReason),
		ProviderStatus:     types.StringPointerValue(endpoint.ProviderStatus),
		IsCustomHostname:   types.BoolValue(endpoint.IsCustomHostname),
		PrivateDnsEnabled:  types.BoolValue(endpoint.PrivateDnsEnabled),
		IpAddresses:        ipAddresses,
		ProviderConfig:     providerConfig,
		Audit:              auditObject,
	}, nil
}

// newOutboundPrivateEndpointProviderConfig is null when the API sends no providerConfig, and
// its aws is null on GCP and Azure.
func newOutboundPrivateEndpointProviderConfig(ctx context.Context, config *outboundpe.ProviderConfig) (types.Object, error) {
	attrTypes := OutboundPrivateEndpointProviderConfig{}.AttributeTypes()
	if config == nil {
		return types.ObjectNull(attrTypes), nil
	}

	awsConfig := types.ObjectNull(OutboundPrivateEndpointAWSConfig{}.AttributeTypes())
	if config.AWS != nil {
		var err error
		awsConfig, err = newOutboundPrivateEndpointAWSConfig(ctx, config.AWS)
		if err != nil {
			return types.ObjectNull(attrTypes), err
		}
	}

	providerConfig, diags := types.ObjectValueFrom(ctx, attrTypes, OutboundPrivateEndpointProviderConfig{AWS: awsConfig})
	if diags.HasError() {
		return types.ObjectNull(attrTypes), fmt.Errorf("%w: %v", errors.ErrConvertingProviderConfig, diags.Errors())
	}

	return providerConfig, nil
}

func newOutboundPrivateEndpointAWSConfig(ctx context.Context, config *outboundpe.AWSConfig) (types.Object, error) {
	attrTypes := OutboundPrivateEndpointAWSConfig{}.AttributeTypes()

	addresses := make([]OutboundPrivateEndpointAWSIpAddress, 0, len(config.IpAddresses))
	for _, address := range config.IpAddresses {
		addresses = append(addresses, OutboundPrivateEndpointAWSIpAddress{
			IpAddress:        types.StringValue(address.IpAddress),
			AvailabilityZone: types.StringValue(address.AvailabilityZone),
		})
	}
	ipAddresses, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: OutboundPrivateEndpointAWSIpAddress{}.AttributeTypes()}, addresses)
	if diags.HasError() {
		return types.ObjectNull(attrTypes), fmt.Errorf("%w: %v", errors.ErrConvertingProviderConfig, diags.Errors())
	}

	zones, diags := types.ListValueFrom(ctx, types.StringType, nonNilStrings(config.AvailabilityZones))
	if diags.HasError() {
		return types.ObjectNull(attrTypes), fmt.Errorf("%w: %v", errors.ErrConvertingProviderConfig, diags.Errors())
	}

	awsConfig, diags := types.ObjectValueFrom(ctx, attrTypes, OutboundPrivateEndpointAWSConfig{
		RegionalDnsName:   types.StringPointerValue(config.RegionalDnsName),
		IpAddresses:       ipAddresses,
		AvailabilityZones: zones,
	})
	if diags.HasError() {
		return types.ObjectNull(attrTypes), fmt.Errorf("%w: %v", errors.ErrConvertingProviderConfig, diags.Errors())
	}

	return awsConfig, nil
}

// nonNilStrings keeps a missing list an empty one, so state doesn't flip between null and [].
func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
