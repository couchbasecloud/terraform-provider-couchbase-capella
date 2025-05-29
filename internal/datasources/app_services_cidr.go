package datasources

import (
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
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
