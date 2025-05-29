package datasources

import (
	"context"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource = (*AppServiceCidrs)(nil)
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
