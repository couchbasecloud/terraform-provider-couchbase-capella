package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	collection_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Collections{}
	_ datasource.DataSourceWithConfigure = &Collections{}
)

// Collections is the collections data source implementation.
type Collections struct {
	*providerschema.Data
}

// NewCollections is a helper function to simplify the provider implementation.
func NewCollections() datasource.DataSource {
	return &Collections{}
}

// Metadata returns the collection data source type name.
func (c *Collections) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collections"
}

// Schema defines the schema for the collection data source.
func (c *Collections) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The data source for retrieving collections within a scope in a bucket.",
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
			"bucket_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the bucket. It is the URL-compatible base64 encoding of the bucket name",
			},
			"scope_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the scope.",
			},
			"data": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of collections in the scope.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"collection_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the collection.",
						},
						"max_ttl": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The maximum Time To Live (TTL) for documents in the collection.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of collections.
func (c *Collections) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Collections
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketId, clusterId, projectId, organizationId, scopeName, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Collections in Capella",
			"Could not read Capella collections in scope "+scopeName+": "+err.Error(),
		)
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s/collections", c.HostURL, organizationId, projectId, clusterId, bucketId, scopeName)
	cfg := collection_api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Collections",
			fmt.Sprintf("Could not read collections in scope %s, unexpected error: %s", scopeName, collection_api.ParseError(err)),
		)
		return
	}

	collectionsResp := collection_api.GetCollectionsResponse{}
	err = json.Unmarshal(response.Body, &collectionsResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading collections",
			"Could not read collections in scope, unexpected error: "+err.Error(),
		)
		return
	}
	for i := range collectionsResp.Data {
		collection := collectionsResp.Data[i]
		newCollectionData, err := providerschema.NewCollectionData(&collection)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing Collections",
				fmt.Sprintf("Could not list collections, unexpected error: %s", err.Error()),
			)
			return
		}
		state.Data = append(state.Data, *newCollectionData)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (c *Collections) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	c.Data = data
}
