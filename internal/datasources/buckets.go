package datasources

import (
	"context"
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/api/bucket"

	"terraform-provider-capella/internal/api"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"organization_id": schema.StringAttribute{
							Computed: true,
						},
						"project_id": schema.StringAttribute{
							Computed: true,
						},
						"cluster_id": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"storage_backend": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"memory_allocation_in_mb": schema.Int64Attribute{
							Optional: true,
							Computed: true,
						},
						"bucket_conflict_resolution": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"durability_level": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"replicas": schema.Int64Attribute{
							Optional: true,
							Computed: true,
						},
						"flush": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},
						"time_to_live_in_seconds": schema.Int64Attribute{
							Optional: true,
							Computed: true,
						},
						"eviction_policy": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"stats": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"item_count": schema.Int64Attribute{
									Computed: true,
								},
								"ops_per_second": schema.Int64Attribute{
									Computed: true,
								},
								"disk_used_in_mib": schema.Int64Attribute{
									Computed: true,
								},
								"memory_used_in_mib": schema.Int64Attribute{
									Computed: true,
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
	cfg := api.EndpointCfg{url, http.MethodGet, http.StatusOK}

	response, err := api.GetPaginated[[]bucket.GetBucketResponse](ctx, d.Client, d.Token, cfg, api.SortById)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != 404 {
			resp.Diagnostics.AddError(
				"Error Reading Capella Buckets",
				"Could not read buckets in cluster "+clusterId+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Buckets",
			"Could not read buckets in cluster "+clusterId+": "+err.Error(),
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
