package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	samplebucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/sample_bucket"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"cluster_id": schema.StringAttribute{
				Required: true,
			},
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the bucket. This is the base64 encoding of the bucket name.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the sample dataset to be loaded. The name has to be one of the following sample datasets: : \"travel-sample\", \"gamesim-sample\" or \"beer-sample\".",
						},
						"organization_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the Capella organization.",
						},
						"project_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the project.",
						},
						"cluster_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the cluster.",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Type of the bucket. If selected Ephemeral, it is not eligible for imports or App Endpoints creation. The options may also be referred to as Memory and Disk (Couchbase), Memory Only (Ephemeral) in the Couchbase documentation.",
						},
						"storage_backend": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Type of the bucket. If selected Ephemeral, it is not eligible for imports or App Endpoints creation. The options may also be referred to as Memory and Disk (Couchbase), Memory Only (Ephemeral) in the Couchbase documentation.",
						},
						"memory_allocation_in_mb": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The amount of memory to allocate for the bucket memory in MiB. The maximum limit is dependent on the allocation of the KV service. For example, 80% of the allocation.",
						},
						"bucket_conflict_resolution": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The means in which conflicts are resolved during replication. This field may be referred to as conflictResolution in the Couchbase documentation, and seqno and lww may be referred to as sequence Number and Timestamp respectively.",
						},
						"durability_level": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "This is the minimum level at which all writes to the Couchbase bucket must occur. The options for Durability level are as follows, according to the bucket type. For a Couchbase bucket: None, Replicate to Majority, Majority and Persist to Active, Persist to Majority. For an Ephemeral bucket: None, Replicate to Majority",
						},
						"replicas": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of replicas for the bucket.",
						},
						"flush": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Replaced by flushEnabled. Determines whether bucket flush is enabled. Set property to true to be able to delete all items in this bucket using the /flush endpoint. Disable property to avoid inadvertent data loss by calling the the /flush endpoint.",
						},
						"time_to_live_in_seconds": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Specifies the time to live (TTL) value in seconds. This is the maximum time to live for items in the bucket. If specified as 0, TTL is disabled. This is a non-negative value.",
						},
						"eviction_policy": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The policy which Capella adopts to prevent data loss due to memory exhaustion. This may be also known as Ejection Policy in the Couchbase documentation. For Couchbase bucket, Eviction Policy is fullEviction by default. For Ephemeral buckets, Eviction Policy is a required field, and should be one of the following: noEviction, nruEviction",
						},
						"stats": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"item_count": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "Number of documents in the bucket.",
								},
								"ops_per_second": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "Number of operations per second.",
								},
								"disk_used_in_mib": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The amount of disk used (in MiB).",
								},
								"memory_used_in_mib": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The amount of memory used (in MiB).",
								},
							},
						},
					},
				},
			},
		},
	}
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
	response, err := api.GetPaginated[[]samplebucketapi.GetSampleBucketResponse](ctx, d.Client, d.Token, cfg, api.SortById)
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
