package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &FreeTierBuckets{}
	_ datasource.DataSourceWithConfigure = &FreeTierBuckets{}
)

// FreeTierBuckets is the free tier bucket data source implementation.
type FreeTierBuckets struct {
	*Buckets
}

// Metadata is a override method of buckets Metadata.
func (f *FreeTierBuckets) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_free_tier_buckets"
}

// NewFreeTierBuckets is a helper function to simplify the provider implementation.
func NewFreeTierBuckets() datasource.DataSource {
	return &FreeTierBuckets{
		Buckets: &Buckets{},
	}
}
