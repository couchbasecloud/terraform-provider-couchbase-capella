package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

type AppServiceCIDR struct {
	OrganizationId types.String `tfsdk:"organization_id"`

	ProjectId types.String `tfsdk:"project_id"`

	ClusterId types.String `tfsdk:"cluster_id"`

	AppServiceId types.String `tfsdk:"app_service_id"`

	Cidr types.String `tfsdk:"cidr"`

	Comment types.String `tfsdk:"comment"`

	ExpiresAt types.String `tfsdk:"expires_at"`

	Audit types.Object `tfsdk:"audit"`
}
