package datasources

import (
	"context"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"net/http"
)

var (
	_ datasource.DataSource              = (*AppServiceCidrs)(nil)
	_ datasource.DataSourceWithConfigure = (*AppServiceCidrs)(nil)
)

type AppServiceCidrs struct {
	*providerschema.Data
}

func NewAppServiceCidrs() datasource.DataSource {
	return &AppServiceCidrs{}
}

func (a *AppServiceCidrs) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_app_service_cidrs"
}

func (a *AppServiceCidrs) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the allowlists details for a Capella App Service.",
		Attributes: map[string]schema.Attribute{
			"data": schema.ListNestedAttribute{
				MarkdownDescription: "The list of allowed CIDRs on an App Service. ",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the allowed CIDR.",
						},
						"organization_id": computedStringAttribute,
						"project_id":      computedStringAttribute,
						"cluster_id":      computedStringAttribute,
						"cidr": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The trusted CIDR to allow the database connections from.",
						},
						"comment": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The trusted CIDR to allow the database connections from.",
						},
						"expires_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "An RFC3339 timestamp determining when the allowed CIDR will expire. If this field is empty/omitted then the allowed CIDR is permanent and will never automatically expire.",
						},
						"audit": computedAuditAttribute,
					},
				},
			},
		},
	}
}

// listAllowLists executes calls to the list allowlist endpoint. It handles pagination and
// returns a slice of individual allowlists responses retrieved from multiple pages.
func (d *AllowLists) listAllowedCIDRs(ctx context.Context, organizationId, projectId, clusterId, appServiceId string) ([]api.GetAllowListResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/allowedcidrs",
		d.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	return api.GetPaginated[[]api.GetAllowListResponse](ctx, d.Client, d.Token, cfg, api.SortById)
}

func (a *AppServiceCidrs) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// TODO
}

func (a *AppServiceCidrs) Configure(
	_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	// TODO
}
