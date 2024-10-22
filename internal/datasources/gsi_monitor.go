package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

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
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"bucket_name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"scope_name":      optionalStringAttribute,
			"collection_name": optionalStringAttribute,
			"index_name":      requiredStringAttribute,
			"status":          computedStringAttribute,
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

	uri := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/queryService/indexBuildStatus/%s?bucket=%s&scope=%s&collection=%s",
		g.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		url.QueryEscape(config.IndexName.ValueString()),
		config.BucketName.ValueString(),
		scope,
		collection,
	)

	cfg := api.EndpointCfg{Url: uri, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	// 60 min is arbitrary.
	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	attempt := 0
	timer := time.NewTimer((1 << attempt) * time.Minute)

	for {
		select {
		case <-ctx.Done():
			resp.Diagnostics.AddWarning(
				"Index build did not complete",
				fmt.Sprintf(
					"Index build for %s in %s.%s.%s did not complete",
					config.IndexName.ValueString(),
					config.BucketName.ValueString(),
					scope,
					collection,
				),
			)
			return
		case <-timer.C:
			response, err := g.Client.ExecuteWithRetry(
				ctx,
				cfg,
				nil,
				g.Token,
				nil,
			)

			if err != nil {
				resp.Diagnostics.AddError(
					"Error getting index build status",
					fmt.Sprintf(
						"Could not get index build status for %s in %s.%s.%s.  Error: %s",
						config.IndexName.ValueString(),
						config.BucketName.ValueString(),
						scope,
						collection,
						err.Error(),
					),
				)
				return
			}

			status := api.IndexBuildStatusResponse{}
			if err = json.Unmarshal(response.Body, &status); err != nil {
				resp.Diagnostics.AddError(
					"Error unmarshaling index build status response",
					fmt.Sprintf(
						"Could not get index build status for %s in %s.%s.%s.  Error: %s",
						config.IndexName.ValueString(),
						config.BucketName.ValueString(),
						scope,
						collection,
						err.Error(),
					),
				)
				return
			}

			if status.Status == "Ready" {
				config.Status = types.StringValue(status.Status)
				resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
				if resp.Diagnostics.HasError() {
					return
				}
				return
			}

			attempt++
			// exponential backoff upto a max of 20 min.
			d := min(20, 1<<attempt)
			timer.Reset(time.Duration(d) * time.Minute)
		}

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
