package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = (*GsiMonitor)(nil)
	_ datasource.DataSourceWithConfigure = (*GsiMonitor)(nil)
)

// GsiMonitor is the data source implementation.
type GsiMonitor struct {
	*providerschema.Data
}

func NewGsiMonitor() datasource.DataSource {
	return &GsiMonitor{}
}

func (g *GsiMonitor) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_query_index_monitor"
}

func (g *GsiMonitor) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The data source for monitoring Query Indexes in Couchbase Capella.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"bucket_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the bucket where the indexes exist. Specifies the bucket part of the key space.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"scope_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The name of the scope where the indexes exist. Specifies the scope portion of the keyspace. If unspecified, this will be the default scope.",
			},
			"collection_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The name of the collection where the indexes exist. Specifies the collection portion of the keyspace. If unspecified, this will be the default collection.",
			},
			"indexes": schema.SetAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Set of index names to monitor. These indexes must exist in the specified keyspace.",
			},
		},
	}
}

func (g *GsiMonitor) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.GsiBuildStatus
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// data sources don't support Default.
	// https://github.com/hashicorp/terraform-plugin-framework/issues/751 .
	var scope, collection string
	if config.ScopeName.IsNull() {
		scope = "_default"
	} else {
		scope = config.ScopeName.ValueString()
	}
	if config.CollectionName.IsNull() {
		collection = "_default"
	} else {
		collection = config.CollectionName.ValueString()
	}

	var indexes []string
	diags := config.Indexes.ElementsAs(ctx, &indexes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor := func(cfg api.EndpointCfg) (response *api.Response, err error) {
		if err := api.Limiter.Wait(ctx); err != nil {
			// do not block if rate limiter fails
			tflog.Error(ctx, "rate limiter error: "+err.Error())
		}

		return g.ClientV1.ExecuteWithRetry(
			ctx,
			cfg,
			nil,
			g.Token,
			nil,
		)
	}

	err := api.WatchIndexes(
		ctx,
		"Ready",
		indexes,
		monitor,
		api.Options{
			Host:       g.HostURL,
			OrgId:      config.OrganizationId.ValueString(),
			ProjectId:  config.ProjectId.ValueString(),
			ClusterId:  config.ClusterId.ValueString(),
			Bucket:     config.BucketName.ValueString(),
			Scope:      scope,
			Collection: collection,
		},
	)
	switch err {
	case nil:
		const msg = `All indexes are ready. Please run "terraform apply --refresh-only" to update state.`
		tflog.Info(ctx, msg)
	default:
		resp.Diagnostics.AddWarning(
			"All provided indexes are not ready",
			"All provided indexes have not completed building.  Please run \"terraform apply --refresh-only\" after some time.",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (g *GsiMonitor) Configure(
	_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	g.Data = data
}
