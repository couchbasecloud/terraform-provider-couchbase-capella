package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

// AppServiceAdminUserData is the data source model for a single App Service admin user.
type AppServiceAdminUserData struct {
	// OrganizationId is the ID of the organization to which the App Service belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the App Service belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster to which the App Service belongs.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service to which the admin user belongs.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Id is the ID of the admin user.
	Id types.String `tfsdk:"id"`

	// Name is the name of the admin user.
	Name types.String `tfsdk:"name"`

	// AccessAllEndpoints is true when the admin user has access to all App Endpoints.
	AccessAllEndpoints types.Bool `tfsdk:"access_all_endpoints"`

	// Endpoints is the set of App Endpoints the admin user has access to.
	Endpoints types.Set `tfsdk:"endpoints"`

	// Audit contains the audit information for the admin user.
	Audit types.Object `tfsdk:"audit"`
}

// AppServiceAdminUsers is the data source model for listing App Service admin users.
type AppServiceAdminUsers struct {
	// OrganizationId is the ID of the organization to which the App Service belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the App Service belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster to which the App Service belongs.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service to which the admin users belong.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Filters is an optional filter to apply to the admin users list.
	Filters *Filter `tfsdk:"filter"`

	// Data is the list of admin users.
	Data []OneAppServiceAdminUser `tfsdk:"data"`
}

// OneAppServiceAdminUser is a single admin user in the list data source.
type OneAppServiceAdminUser struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	AccessAllEndpoints types.Bool   `tfsdk:"access_all_endpoints"`
	Endpoints          types.Set    `tfsdk:"endpoints"`
	Audit              types.Object `tfsdk:"audit"`
}

// NewOneAppServiceAdminUser maps an API admin user to its Terraform model.
func NewOneAppServiceAdminUser(ctx context.Context, user apigen.AppServiceAdminUser) (OneAppServiceAdminUser, diag.Diagnostics) {
	var endpoints []string
	if user.Endpoints != nil {
		endpoints = *user.Endpoints
	}
	endpointsSet, diags := types.SetValueFrom(ctx, types.StringType, endpoints)

	auditObj := types.ObjectNull(CouchbaseAuditData{}.AttributeTypes())
	if user.Audit != nil {
		audit := NewCouchbaseAuditData(api.CouchbaseAuditData(*user.Audit))
		var d diag.Diagnostics
		auditObj, d = types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		diags.Append(d...)
	}

	return OneAppServiceAdminUser{
		Id:                 types.StringPointerValue(user.Id),
		Name:               types.StringPointerValue(user.Name),
		AccessAllEndpoints: types.BoolPointerValue(user.AccessAllEndpoints),
		Endpoints:          endpointsSet,
		Audit:              auditObj,
	}, diags
}
