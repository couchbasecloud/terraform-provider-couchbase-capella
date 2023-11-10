package datasources

import (
	"context"
	"fmt"
	"net/http"
	appservice "terraform-provider-capella/internal/api/appservice"

	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AppServices{}
	_ datasource.DataSourceWithConfigure = &AppServices{}
)

// AppServices is the AppServices data source implementation.
type AppServices struct {
	*providerschema.Data
}

// NewAppServices is a helper function to simplify the provider implementation.
func NewAppServices() datasource.DataSource {
	return &AppServices{}
}

// Metadata returns the app services data source type name.
func (d *AppServices) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_services"
}

// Schema defines the schema for the AppServices data source.
func (d *AppServices) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = AppServiceSchema()
}

// Read refreshes the Terraform state with the latest data of app services.
func (d *AppServices) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.AppServices
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error reading app services",
			"Could not read app services, unexpected error: organization ID cannot be empty.",
		)
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/appservices", d.HostURL, organizationId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]appservice.GetAppServiceResponse](ctx, d.Client, d.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella App Services",
			fmt.Sprintf("Could not read app services in organization %s, unexpected error: %s", organizationId, api.ParseError(err)),
		)
		return
	}

	for _, appService := range response {
		audit := providerschema.NewCouchbaseAuditData(appService.Audit)

		auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error Reading Capella App Services",
				fmt.Sprintf("Could not read app services in organization %s, unexpected error: %s", organizationId, errors.ErrUnableToConvertAuditData),
			)
		}

		newAppServiceData := providerschema.NewAppServiceData(&appService, organizationId, auditObj)
		state.Data = append(state.Data, *newAppServiceData)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the cluster data source.
func (d *AppServices) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.Data = data
}
