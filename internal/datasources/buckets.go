package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Buckets{}
	_ datasource.DataSourceWithConfigure = &Buckets{}
)

// Buckets is the bucket data source implementation.
type Buckets struct {
	*providerschema.Data
}

// NewBuckets is a helper function to simplify the provider implementation.
func NewBuckets() datasource.DataSource {
	return &Buckets{}
}

// Metadata returns the bucket data source type name.
func (d *Buckets) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_buckets"
}

// Schema defines the schema for the bucket data source.
func (d *Buckets) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The bucket data source retrieves the bucket information for an operational cluster.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
			},
			"data": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Lists the bucket data in the cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the bucket.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the bucket.",
						},
						"organization_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the organization.",
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
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "The bucket type (Couchbase or Ephemeral).",
						},
						"storage_backend": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "The bucket storage engine type (Magma or Couchstore).",
						},
						"memory_allocation_in_mb": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Bucket size in MB.",
						},
						"bucket_conflict_resolution": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Conflict-resolution mechanism of the bucket.",
						},
						"durability_level": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Durability of the bucket.",
						},
						"replicas": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Number of replicas for the data.",
						},
						"flush": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Determines whether flushing is enabled in the bucket.",
						},
						"time_to_live_in_seconds": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Time-to-live (TTL) for items in the bucket, in seconds.",
						},
						"eviction_policy": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "Eviction policy for the bucket.",
						},
						"stats": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Bucket statistics.",
							Attributes: map[string]schema.Attribute{
								"item_count": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "Bucket item count.",
								},
								"ops_per_second": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "Bucket operations per second.",
								},
								"disk_used_in_mib": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "Disk used in MiB.",
								},
								"memory_used_in_mib": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "Memory used in MiB.",
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of buckets.
func (d *Buckets) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Buckets
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Buckets in Capella",
			"Could not read Capella buckets in cluster "+clusterId+": "+err.Error(),
		)
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", d.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := api.GetPaginated[[]bucket.GetBucketResponse](ctx, d.ClientV1, d.Token, cfg, api.SortById)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Buckets",
			"Could not read buckets in cluster "+clusterId+": "+api.ParseError(err),
		)
		return
	}

	// Map response body to model
	for _, bucket := range response {
		bucketState := providerschema.OneBucket{
			Id:                       types.StringValue(bucket.Id),
			Name:                     types.StringValue(bucket.Name),
			Type:                     types.StringValue(bucket.Type),
			OrganizationId:           types.StringValue(organizationId),
			ProjectId:                types.StringValue(projectId),
			ClusterId:                types.StringValue(clusterId),
			StorageBackend:           types.StringValue(bucket.StorageBackend),
			MemoryAllocationInMB:     types.Int64Value(bucket.MemoryAllocationInMb),
			BucketConflictResolution: types.StringValue(bucket.BucketConflictResolution),
			DurabilityLevel:          types.StringValue(bucket.DurabilityLevel),
			Replicas:                 types.Int64Value(bucket.Replicas),
			Flush:                    types.BoolValue(bucket.Flush),
			TimeToLiveInSeconds:      types.Int64Value(bucket.TimeToLiveInSeconds),
			EvictionPolicy:           types.StringValue(bucket.EvictionPolicy),
			Stats: &providerschema.Stats{
				ItemCount:       types.Int64Value(bucket.Stats.ItemCount),
				OpsPerSecond:    types.Int64Value(bucket.Stats.OpsPerSecond),
				DiskUsedInMiB:   types.Int64Value(bucket.Stats.DiskUsedInMib),
				MemoryUsedInMiB: types.Int64Value(bucket.Stats.MemoryUsedInMib),
			},
		}
		state.Data = append(state.Data, bucketState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the bucket data source.
func (d *Buckets) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
