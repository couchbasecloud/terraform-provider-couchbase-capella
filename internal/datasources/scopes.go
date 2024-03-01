package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	scope_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/scope"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Scopes{}
	_ datasource.DataSourceWithConfigure = &Scopes{}
)

// Scopes is the scope data source implementation.
type Scopes struct {
	*providerschema.Data
}

// NewScopes is a helper function to simplify the provider implementation.
func NewScopes() datasource.DataSource {
	return &Scopes{}
}

// Metadata returns the api key data source type name.
func (s *Scopes) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scopes"
}

// Schema defines the schema for the scope data source.
func (s *Scopes) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"cluster_id":      requiredStringAttribute,
			"bucket_id":       requiredStringAttribute,
			"uid":             computedStringAttribute,
			"scopes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uid":             computedStringAttribute,
						"organization_id": computedStringAttribute,
						"project_id":      computedStringAttribute,
						"cluster_id":      computedStringAttribute,
						"bucket_id":       computedStringAttribute,
						"name":            computedStringAttribute,
						"collections": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"max_ttl": computedInt64Attribute,
									"name":    computedStringAttribute,
									"uid":     computedStringAttribute,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of scopes.
func (s *Scopes) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.Scopes
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketId, clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Scopes in Capella",
			"Could not read Capella scopes in bucket "+bucketId+": "+err.Error(),
		)
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes", s.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := s.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		s.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Scopes",
			fmt.Sprintf("Could not read scopes in bucket %s, unexpected error: %s", bucketId, api.ParseError(err)),
		)
		return
	}

	scopesResp := scope_api.GetScopesResponse{}
	err = json.Unmarshal(response.Body, &scopesResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading scope",
			"Could not read scope in bucket, unexpected error: "+err.Error(),
		)
		return
	}

	for _, scope := range scopesResp.Scopes {
		//collections := make([]types.Object, 0)
		collections := make([]providerschema.Collection, 0)
		for _, apiCollection := range *scope.Collections {
			collection := providerschema.Collection{
				MaxTTL: types.Int64Value(*apiCollection.MaxTTL),
				Name:   types.StringValue(*apiCollection.Name),
				Uid:    types.StringValue(*apiCollection.Uid),
			}
			//collectionObj, diags := types.ObjectValueFrom(ctx, providerschema.CollectionAttributeTypes(), collection)
			//if diags.HasError() {
			//	fmt.Errorf("error Converting Collection to Object")
			//	return
			//}
			//collections = append(collections, collectionObj)
			collections = append(collections, collection)
		}
		//collectionSet, diags := types.SetValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(providerschema.CollectionAttributeTypes()), collections)
		//if diags.HasError() {
		//	fmt.Errorf("error Converting Collection object to set")
		//	return
		//}

		//newScope := providerschema.Scope{
		newScope := providerschema.ScopeData{
			Name:           types.StringValue(*scope.Name),
			Uid:            types.StringValue(*scope.Uid),
			Collections:    collections,
			OrganizationId: types.StringValue(organizationId),
			ProjectId:      types.StringValue(projectId),
			ClusterId:      types.StringValue(clusterId),
			BucketId:       types.StringValue(bucketId),
		}
		state.Scopes = append(state.Scopes, newScope)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (s *Scopes) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	s.Data = data
}
