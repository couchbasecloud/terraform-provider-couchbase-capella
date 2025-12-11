package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
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
	resp.Schema = ScopesSchema()
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

	response, err := s.ClientV1.ExecuteWithRetry(
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
			"Error reading scopes",
			"Could not read scopes in bucket, unexpected error: "+err.Error(),
		)
		return
	}

	for i := range scopesResp.Scopes {
		scope := scopesResp.Scopes[i]
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

func (s *Scopes) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
