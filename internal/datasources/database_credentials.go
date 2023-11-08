package datasources

import (
	"context"
	"fmt"

	"terraform-provider-capella/internal/api"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DatabaseCredentials{}
	_ datasource.DataSourceWithConfigure = &DatabaseCredentials{}
)

// DatabaseCredentials is the database credential data source implementation.
type DatabaseCredentials struct {
	*providerschema.Data
}

// NewDatabaseCredentials is a helper function to simplify the provider implementation.
func NewDatabaseCredentials() datasource.DataSource {
	return &DatabaseCredentials{}
}

// Metadata returns the database credential data source type name.
func (d *DatabaseCredentials) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_credentials"
}

// Schema defines the schema for the database credential data source.
func (d *DatabaseCredentials) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"cluster_id":      requiredStringAttribute,
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              computedStringAttribute,
						"name":            computedStringAttribute,
						"organization_id": computedStringAttribute,
						"project_id":      computedStringAttribute,
						"cluster_id":      computedStringAttribute,
						"audit":           computedAuditAttribute,
						"access": schema.ListNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"privileges": schema.ListAttribute{
										Required:    true,
										ElementType: types.StringType,
									},
									"resources": schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{
											"buckets": schema.ListNestedAttribute{
												Optional: true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"name": requiredStringAttribute,
														"scopes": schema.ListNestedAttribute{
															Optional: true,
															NestedObject: schema.NestedAttributeObject{
																Attributes: map[string]schema.Attribute{
																	"name": requiredStringAttribute,
																	"collections": schema.ListAttribute{
																		Optional:    true,
																		ElementType: types.StringType,
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of database credentials.
func (d *DatabaseCredentials) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.DatabaseCredentials
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, projectId, organizationId, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Database Credentials in Capella",
			"Could not read Capella database credentials in cluster "+clusterId+": "+err.Error(),
		)
		return
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/users", d.HostURL, organizationId, projectId, clusterId)
	response, err := api.GetPaginated[[]api.GetDatabaseCredentialResponse](ctx, d.Client, d.Token, url, api.SortById)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != 404 {
			resp.Diagnostics.AddError(
				"Error Reading Capella Database Credentials",
				"Could not read database credentials in cluster "+clusterId+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella Database Credentials",
			"Could not read database credentials in cluster "+clusterId+": "+api.ParseError(err),
		)
		return
	}

	// Map response body to model
	for _, databaseCredential := range response {
		databaseCredentialState := providerschema.DatabaseCredentialItem{
			Id:             types.StringValue(databaseCredential.Id.String()),
			Name:           types.StringValue(databaseCredential.Name),
			OrganizationId: types.StringValue(organizationId),
			ProjectId:      types.StringValue(projectId),
			ClusterId:      types.StringValue(clusterId),
			Audit: providerschema.CouchbaseAuditData{
				CreatedAt:  types.StringValue(databaseCredential.Audit.CreatedAt.String()),
				CreatedBy:  types.StringValue(databaseCredential.Audit.CreatedBy),
				ModifiedAt: types.StringValue(databaseCredential.Audit.ModifiedAt.String()),
				ModifiedBy: types.StringValue(databaseCredential.Audit.ModifiedBy),
				Version:    types.Int64Value(int64(databaseCredential.Audit.Version)),
			},
		}
		databaseCredentialState.Access = mapAccess(databaseCredential)
		state.Data = append(state.Data, databaseCredentialState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the database credential data source.
func (d *DatabaseCredentials) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// mapAccess needs a 1:1 mapping when we store the output as the refreshed state.
// todo: add a unit test, tracking under: https://couchbasecloud.atlassian.net/browse/AV-63401
func mapAccess(plan api.GetDatabaseCredentialResponse) []providerschema.Access {
	var access = make([]providerschema.Access, len(plan.Access))

	for i, acc := range plan.Access {
		access[i] = providerschema.Access{Privileges: make([]types.String, len(acc.Privileges))}
		for j, permission := range acc.Privileges {
			access[i].Privileges[j] = types.StringValue(permission)
		}
		if acc.Resources != nil {
			if acc.Resources.Buckets != nil {
				access[i].Resources = &providerschema.Resources{Buckets: make([]providerschema.BucketResource, len(acc.Resources.Buckets))}
				for k, bucket := range acc.Resources.Buckets {
					access[i].Resources.Buckets[k].Name = types.StringValue(acc.Resources.Buckets[k].Name)
					if bucket.Scopes != nil {
						access[i].Resources.Buckets[k].Scopes = make([]providerschema.Scope, len(bucket.Scopes))
						for s, scope := range bucket.Scopes {
							access[i].Resources.Buckets[k].Scopes[s].Name = types.StringValue(scope.Name)
							if scope.Collections != nil {
								access[i].Resources.Buckets[k].Scopes[s].Collections = make([]types.String, len(scope.Collections))
								for c, coll := range scope.Collections {
									access[i].Resources.Buckets[k].Scopes[s].Collections[c] = types.StringValue(coll)
								}
							}
						}
					}
				}
			}
		}
	}

	return access
}
