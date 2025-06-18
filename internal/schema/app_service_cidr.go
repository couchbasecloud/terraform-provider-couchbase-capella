package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// AppServiceCIDRs defines the attributes for an individual App service allowed CIDR.
type AppServiceCIDR struct {
	// OrganizationId is the Capella tenant id associated with the App Service.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the id of the Capella project associated with the App Service.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the id of the Capella cluster associated with the App Service.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the id of the Capella App Service associated with the App Service CIDR.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Cidr is the CIDR to allow connections from.
	Cidr types.String `tfsdk:"cidr"`
	// Id is a GUID4 identifier of the App service CIDR.
	Comment types.String `tfsdk:"comment"`
	// ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt types.String `tfsdk:"expires_at"`
	// Audit contains the audit information for the App service CIDR. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`
	// Id is the ID is the unique UUID generated when an allowed cidr is created.
	Id types.String `tfsdk:"id"`
}

// AppServiceCIDRs defines the attributes for an individual App service allowed CIDR.
type AppServiceCIDRData struct {
	// Id is the ID is the unique UUID generated when an allowed cidr is created.
	Id types.String `tfsdk:"id"`
	// Cidr is the CIDR or ip address range to allow connections from.
	Cidr types.String `tfsdk:"cidr"`
	// Id is a GUID4 identifier of the App service CIDR.
	Comment types.String `tfsdk:"comment"`
	// ExpiresAt is an RFC3339 timestamp determining when the allowed CIDR should expire.
	ExpiresAt types.String `tfsdk:"expires_at"`
	// Audit contains the audit information for the App service CIDR. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`
}

// AppServiceCIDRs defines the attributes as received from the
// V4 Capella Public API when asked to list App service allowed CIDRs.
type AppServiceCIDRs struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the GUID4 identifier for the project.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the GUID4 identifier for the Cluster.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the GUID4 identifier for the App Service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Data contains the list of resources.
	Data []AppServiceCIDRData `tfsdk:"data"`
}

// Validate is used to verify that all the fields in the datasource
// have been populated.
func (a *AppServiceCIDRs) Validate() (clusterId, projectId, organizationId, appserviceId string, err error) {
	if a.OrganizationId.IsNull() {
		return "", "", "", "", errors.ErrOrganizationIdMissing
	}
	if a.ProjectId.IsNull() {
		return "", "", "", "", errors.ErrProjectIdMissing
	}
	if a.ClusterId.IsNull() {
		return "", "", "", "", errors.ErrClusterIdMissing
	}
	if a.AppServiceId.IsNull() {
		return "", "", "", "", errors.ErrAppServiceIdMissing
	}
	return a.ClusterId.ValueString(), a.ProjectId.ValueString(), a.OrganizationId.ValueString(), a.AppServiceId.ValueString(), nil
}

// Validate is used to verify that all the fields in the datasource
// have been populated.
func (a *AppServiceCIDR) Validate() (organizationId, projectId, clusterId, appserviceId, allowedCIDRId string, err error) {
	if a.OrganizationId.IsNull() {
		return "", "", "", "", "", errors.ErrOrganizationIdMissing
	}
	if a.ProjectId.IsNull() {
		return "", "", "", "", "", errors.ErrProjectIdMissing
	}
	if a.ClusterId.IsNull() {
		return "", "", "", "", "", errors.ErrClusterIdMissing
	}
	if a.AppServiceId.IsNull() {
		return "", "", "", "", "", errors.ErrAppServiceIdMissing
	}
	if a.Id.IsNull() {
		return "", "", "", "", "", errors.ErrAllowListIdMissing
	}
	return a.OrganizationId.ValueString(), a.ProjectId.ValueString(), a.ClusterId.ValueString(), a.AppServiceId.ValueString(), a.Id.ValueString(), nil
}
