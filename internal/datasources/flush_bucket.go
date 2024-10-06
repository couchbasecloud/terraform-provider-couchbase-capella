package datasources

import (
	"context"
	"fmt"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &FlushBucket{}
	_ datasource.DataSourceWithConfigure = &FlushBucket{}
)

// FlushBucket is the bucket data source implementation.
type FlushBucket struct {
	*providerschema.Data
}

// NewFlushBucket is a helper function to simplify the provider implementation.
func NewFlushBucket() datasource.DataSource {
	return &FlushBucket{}
}

// Metadata returns the bucket data source type name.
func (d *FlushBucket) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flush"
}

// Schema defines the schema for the bucket data source.
func (s *FlushBucket) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"cluster_id":      requiredStringAttribute,
			"bucket_id":       requiredStringAttribute,
		},
	}
}

// Read refreshes the Terraform state with the latest data of buckets.
func (d *FlushBucket) Read(_ context.Context, _ datasource.ReadRequest, _ *datasource.ReadResponse) {
	// Do we even need this? Flush bucket information is all purely in terraform and not Capella.
}

// Configure adds the provider configured client to the bucket data source.
func (d *FlushBucket) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
