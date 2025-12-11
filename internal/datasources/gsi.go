package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
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

func (g *GsiDefinitions) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GsiSchema()
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

	response, err := g.ClientV1.ExecuteWithRetry(
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
