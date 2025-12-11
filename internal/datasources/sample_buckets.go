package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	samplebucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/sample_bucket"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &SampleBuckets{}
	_ datasource.DataSourceWithConfigure = &SampleBuckets{}
)

// Sample buckets is the sample bucket data source implementation.
type SampleBuckets struct {
	*providerschema.Data
}

// NewSampleBuckets is a helper function to simplify the provider implementation.
func NewSampleBuckets() datasource.DataSource {
	return &SampleBuckets{}
}

// Metadata returns the sample bucket data source type name.
func (d *SampleBuckets) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sample_buckets"
}

// Schema defines the schema for the sample bucket data source.
func (s *SampleBuckets) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SampleBucketsSchema()
}

// Read refreshes the Terraform state with the latest data of sample buckets.
func (d *SampleBuckets) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.SampleBuckets
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Sample Buckets in Capella",
			"Could not read sample buckets in cluster "+clusterId+": "+err.Error(),
		)
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/sampleBuckets", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := api.GetPaginated[[]samplebucketapi.GetSampleBucketResponse](ctx, d.ClientV1, d.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Sample Buckets in Capella",
			"Could not read sample buckets in cluster "+clusterId+": "+api.ParseError(err),
		)
		return
	}

	// Map response body to model
	for _, sampleBucket := range response {
		var sampleStats providerschema.Stats

		if sampleBucket.Stats != nil {
			sampleStats = providerschema.NewStats(*sampleBucket.Stats)
		}

		sampleBucketStatsObj, diags := types.ObjectValueFrom(ctx, sampleStats.AttributeTypes(), sampleStats)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error Error Reading Sample Bucket Info",
				fmt.Sprintf("Could not read sample bucket info from record, unexpected error: %s", fmt.Errorf("error while sample bucket info conversion")),
			)
			return
		}

		sampleBucketState := providerschema.SampleBucket{
			Id:                       types.StringValue(sampleBucket.Id),
			Name:                     types.StringValue(sampleBucket.Name),
			Type:                     types.StringValue(sampleBucket.Type),
			OrganizationId:           types.StringValue(organizationId),
			ProjectId:                types.StringValue(projectId),
			ClusterId:                types.StringValue(clusterId),
			StorageBackend:           types.StringValue(sampleBucket.StorageBackend),
			MemoryAllocationInMB:     types.Int64Value(sampleBucket.MemoryAllocationInMb),
			BucketConflictResolution: types.StringValue(sampleBucket.BucketConflictResolution),
			DurabilityLevel:          types.StringValue(sampleBucket.DurabilityLevel),
			Replicas:                 types.Int64Value(sampleBucket.Replicas),
			Flush:                    types.BoolValue(sampleBucket.Flush),
			TimeToLiveInSeconds:      types.Int64Value(sampleBucket.TimeToLiveInSeconds),
			EvictionPolicy:           types.StringValue(sampleBucket.EvictionPolicy),
			Stats:                    sampleBucketStatsObj,
		}
		state.Data = append(state.Data, sampleBucketState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the bucket data source.
func (d *SampleBuckets) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
