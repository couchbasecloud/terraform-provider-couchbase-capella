package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
	_ datasource.DataSource              = (*GsiDefinitions)(nil)
	_ datasource.DataSourceWithConfigure = (*GsiDefinitions)(nil)
)

// GsiDefinitions is the data source implementation.
type GsiDefinitions struct {
	*providerschema.Data
}

func NewGsiDefinitions() datasource.DataSource {
	return &GsiDefinitions{}
}

func (g *GsiDefinitions) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_query_indexes"
}

func (g *GsiDefinitions) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for retrieving Query Indexes in Couchbase Capella",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Capella organization.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Capella project.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Capella cluster where the indexes exist.",
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
				MarkdownDescription: "The name of the scope where the indexes exist. Specifies the scope part of the key space. If unspecified, this will be the default scope.",
			},
			"collection_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies the collection part of the key space. If unspecified, this will be the default collection.",
			},
			"data": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of indexes in the specified keyspace.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the index.",
						},
						"is_primary": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether this is a primary index.",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The current state of the index. For example 'Created', 'Ready', etc.",
						},
						"keyspace_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The full keyspace identifier for the index (bucket.scope.collection).",
						},
						"index_key": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "List of document fields being indexed.",
						},
						"condition": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The WHERE clause condition for the index.",
						},
						"partition": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "List of fields the index is partitioned by.",
						},
						"replica_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Number of index replicas.",
						},
						"partition_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Number of partitions for the index.",
						},
					},
				},
			},
		},
	}
}

func (g *GsiDefinitions) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config providerschema.GsiDefinitions
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

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/queryService/indexes?bucket=%s&scope=%s&collection=%s",
		g.HostURL,
		config.OrganizationId.ValueString(),
		config.ProjectId.ValueString(),
		config.ClusterId.ValueString(),
		config.BucketName.ValueString(),
		scope,
		collection,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := g.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		g.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Listing Query Indexes",
			fmt.Sprintf(
				"Could not list query indexes in %s.%s.%s.  Error: %s",
				config.BucketName.ValueString(),
				scope,
				collection,
				api.ParseError(err),
			),
		)
		return
	}

	gsiList := api.ListIndexDefinitionsResponse{}
	err = json.Unmarshal(response.Body, &gsiList)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshaling index definitions",
			"Could not unmarshal index definitions.  Error: "+err.Error(),
		)
		return
	}

	config.Data = make([]providerschema.GsiData, len(gsiList.Definitions))
	for i, index := range gsiList.Definitions {
		config.Data[i].IndexName = types.StringValue(index.IndexName)
		config.Data[i].Definition = types.StringValue(index.Definition)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (g *GsiDefinitions) Configure(
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
