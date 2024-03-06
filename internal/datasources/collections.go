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
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"cluster_id":      requiredStringAttribute,
			"bucket_id":       requiredStringAttribute,
			"scope_name":      requiredStringAttribute,
			"collection_name": requiredStringAttribute,
			"max_ttl":         requiredInt64Attribute,
			"uid":             computedStringAttribute,
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":    computedStringAttribute,
						"max_ttl": computedInt64Attribute,
						"uid":     computedStringAttribute,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of scopes.
func (c *Collections) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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

	for i := range scopesResp.Scopes {
		scope := scopesResp.Scopes[i]
		//for _, scope := range scopesResp.Scopes {
		objectList := make([]types.Object, 0)
		for _, apiCollection := range *scope.Collections {
			providerschemaCollection := providerschema.NewCollection(apiCollection)
			collectionObj, diag := types.ObjectValueFrom(ctx, providerschema.CollectionAttributeTypes(), providerschemaCollection)
			if diag.HasError() {
				resp.Diagnostics.AddError(
					"Collection obj error",
					"Could not read collection object",
				)
			}

			objectList = append(objectList, collectionObj)

		}

		collectionSet, diag := types.SetValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(providerschema.CollectionAttributeTypes()), objectList)
		if diag.HasError() {
			resp.Diagnostics.AddError(
				"collectionSet error",
				"Could not read collection set",
			)
		}

		//the array of scopes is listed together under one uid
		stateScopes := state.Scopes
		newScope := providerschema.NewScopeData(&scope, collectionSet)

		stateScopes = append(stateScopes, *newScope)

		scopeState := providerschema.Scopes{
			Scopes:         stateScopes,
			Uid:            types.StringValue(*scopesResp.Uid),
			OrganizationId: types.StringValue(organizationId),
			ProjectId:      types.StringValue(projectId),
			ClusterId:      types.StringValue(clusterId),
			BucketId:       types.StringValue(bucketId),
		}

		state = scopeState
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (c *Collections) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
