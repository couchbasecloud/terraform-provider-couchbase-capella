package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

type Certificate struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	Certificate types.String `tfsdk:"certificate"`
}

type Certificates struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	Data []OneCertificate `tfsdk:"data"`
}

type OneCertificate struct {
	Certificate types.String `tfsdk:"certificate"`
}
