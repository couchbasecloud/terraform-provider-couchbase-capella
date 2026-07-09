package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// PrivateEndpointService represents the status of private endpoint service on a cluster.
type PrivateEndpointService struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster associated with the private endpoint service.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Enabled indicates if private endpoint service is enabled/disabled on cluster.
	Enabled types.Bool `tfsdk:"enabled"`

	// Status is the lifecycle state of the private endpoint service derived from
	// the most recent enable/disable/update operation. Terminal states are
	// enableFailed and disableFailed; transient states are enabling, disabling,
	// and unknown; idle means no operation has run. It may be empty when the
	// control plane does not report a status.
	Status types.String `tfsdk:"status"`

	// ServiceName is the CSP endpoint service name that customer VPC/VNET
	// endpoints connect to (for AWS, the VPC endpoint service name, e.g.
	// com.amazonaws.vpce.us-east-1.vpce-svc-1234). It is populated once the
	// service is enabled and can be fed directly into an aws_vpc_endpoint
	// resource. It may be empty on older control planes or when the name cannot
	// be determined.
	ServiceName types.String `tfsdk:"service_name"`
}

// Validate is used to verify that IDs have been properly imported.
func (p *PrivateEndpointService) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: p.OrganizationId,
		ProjectId:      p.ProjectId,
		ClusterId:      p.ClusterId,
	}

	IDs, err := validateSchemaState(state, ClusterId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
