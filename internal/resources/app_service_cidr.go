package resources

import (
	"context"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                = (*AppServiceCidr)(nil)
	_ resource.ResourceWithConfigure   = (*AppServiceCidr)(nil)
	_ resource.ResourceWithImportState = (*AppServiceCidr)(nil)
)

type AppServiceCidr struct {
	*providerschema.Data
}

func NewAppServiceCidr() resource.Resource {
	return &AppServiceCidr{}
}

func (a *AppServiceCidr) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (a *AppServiceCidr) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

	// TODO
}
