package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Replication{}
	_ datasource.DataSourceWithConfigure = &Replication{}
)

// Replication is the Replication data source implementation.
type Replication struct {
	*providerschema.Data
}

// NewReplication is a helper function to simplify the provider implementation.
func NewReplication() datasource.DataSource {
	return &Replication{}
}

// Metadata returns the data source type name.
func (d *Replication) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication"
}

// Schema defines the schema for the data source.
func (d *Replication) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ReplicationSchema()
}

// Read refreshes the Terraform state with the latest data from the API.
func (d *Replication) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ReplicationData
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationId.ValueString()
	projectId := state.ProjectId.ValueString()
	clusterId := state.ClusterId.ValueString()
	replicationId := state.ReplicationId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/replications/%s", d.HostURL, organizationId, projectId, clusterId, replicationId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	response, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		d.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Replication",
			fmt.Sprintf("Could not read replication %s: %s", replicationId, api.ParseError(err)),
		)
		return
	}

	var apiResponse replicationapi.GetReplicationResponse
	err = json.Unmarshal(response.Body, &apiResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Unmarshalling Replication",
			fmt.Sprintf("Could not unmarshal replication response: %s", err.Error()),
		)
		tflog.Debug(ctx, "error unmarshalling replication response", map[string]interface{}{
			"response_body": string(response.Body),
			"err":           err.Error(),
		})
		return
	}

	tflog.Info(ctx, "read replication", map[string]interface{}{
		"organization_id": organizationId,
		"project_id":      projectId,
		"cluster_id":      clusterId,
		"replication_id":  replicationId,
	})

	// Build source object
	sourceObj, diags := buildSourceObject(ctx, apiResponse.Source)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build target object
	targetObj, diags := buildTargetObject(ctx, apiResponse.Target)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build mappings list
	mappingsList, diags := buildMappingsList(ctx, apiResponse.Mappings)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build filter object
	filterObj, diags := buildFilterObject(ctx, apiResponse.Filter)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build audit object
	audit := providerschema.ReplicationAuditData{
		CreatedBy:  types.StringValue(apiResponse.Audit.CreatedBy),
		CreatedAt:  types.StringValue(apiResponse.Audit.CreatedAt),
		ModifiedBy: types.StringValue(apiResponse.Audit.ModifiedBy),
		ModifiedAt: types.StringValue(apiResponse.Audit.ModifiedAt),
		Version:    types.Int64Value(apiResponse.Audit.Version),
	}
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the final state
	newState := providerschema.NewReplicationData(
		organizationId,
		projectId,
		clusterId,
		replicationId,
		&apiResponse,
		sourceObj,
		targetObj,
		filterObj,
		mappingsList,
		auditObj,
	)

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the data source.
func (d *Replication) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.Data = data
}

// buildSourceObject builds the source object from API response.
func buildSourceObject(ctx context.Context, source replicationapi.Source) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build project object
	project := providerschema.ProjectReferenceData{
		Id:   types.StringValue(source.Project.Id),
		Name: types.StringValue(source.Project.Name),
	}
	projectObj, d := types.ObjectValueFrom(ctx, project.AttributeTypes(), project)
	if d.HasError() {
		diags.Append(d...)
	}

	// Build cluster object
	cluster := providerschema.ClusterReferenceData{
		Id:   types.StringValue(source.Cluster.Id),
		Name: types.StringValue(source.Cluster.Name),
	}
	clusterObj, d := types.ObjectValueFrom(ctx, cluster.AttributeTypes(), cluster)
	if d.HasError() {
		diags.Append(d...)
	}

	// Build bucket object
	bucket := providerschema.BucketReferenceData{
		Id:                     types.StringValue(source.Bucket.Id),
		Name:                   types.StringValue(source.Bucket.Name),
		ConflictResolutionType: types.StringValue(source.Bucket.ConflictResolutionType),
	}
	bucketObj, d := types.ObjectValueFrom(ctx, bucket.AttributeTypes(), bucket)
	if d.HasError() {
		diags.Append(d...)
	}

	// Build scopes list
	scopesList, d := buildScopesList(ctx, source.Scopes)
	if d.HasError() {
		diags.Append(d...)
	}

	sourceData := providerschema.ReplicationSourceData{
		Project: projectObj,
		Cluster: clusterObj,
		Bucket:  bucketObj,
		Scopes:  scopesList,
		Type:    types.StringValue(source.Type),
	}

	sourceObj, d := types.ObjectValueFrom(ctx, sourceData.AttributeTypes(), sourceData)
	if d.HasError() {
		diags.Append(d...)
	}

	return sourceObj, diags
}

// buildTargetObject builds the target object from API response.
func buildTargetObject(ctx context.Context, target replicationapi.Target) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build project object (optional)
	var projectObj basetypes.ObjectValue
	if target.Project != nil {
		project := providerschema.ProjectReferenceData{
			Id:   types.StringValue(target.Project.Id),
			Name: types.StringValue(target.Project.Name),
		}
		projectObj, _ = types.ObjectValueFrom(ctx, project.AttributeTypes(), project)
	} else {
		projectObj = types.ObjectNull(providerschema.ProjectReferenceData{}.AttributeTypes())
	}

	// Build cluster object
	cluster := providerschema.ClusterReferenceData{
		Id:   types.StringValue(target.Cluster.Id),
		Name: types.StringValue(target.Cluster.Name),
	}
	clusterObj, d := types.ObjectValueFrom(ctx, cluster.AttributeTypes(), cluster)
	if d.HasError() {
		diags.Append(d...)
	}

	// Build bucket object
	bucket := providerschema.BucketReferenceData{
		Id:                     types.StringValue(target.Bucket.Id),
		Name:                   types.StringValue(target.Bucket.Name),
		ConflictResolutionType: types.StringValue(target.Bucket.ConflictResolutionType),
	}
	bucketObj, d := types.ObjectValueFrom(ctx, bucket.AttributeTypes(), bucket)
	if d.HasError() {
		diags.Append(d...)
	}

	// Build scopes list
	scopesList, d := buildScopesList(ctx, target.Scopes)
	if d.HasError() {
		diags.Append(d...)
	}

	targetData := providerschema.ReplicationTargetData{
		Project: projectObj,
		Cluster: clusterObj,
		Bucket:  bucketObj,
		Scopes:  scopesList,
		Type:    types.StringValue(target.Type),
	}

	targetObj, d := types.ObjectValueFrom(ctx, targetData.AttributeTypes(), targetData)
	if d.HasError() {
		diags.Append(d...)
	}

	return targetObj, diags
}

// buildScopesList builds a list of scope objects.
func buildScopesList(ctx context.Context, scopes []replicationapi.Scope) (basetypes.ListValue, diag.Diagnostics) {
	if len(scopes) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: providerschema.ReplicationScopeData{}.AttributeTypes()}), nil
	}

	var scopeObjs []attr.Value
	for _, s := range scopes {
		collections := make([]types.String, len(s.Collections))
		for i, c := range s.Collections {
			collections[i] = types.StringValue(c)
		}

		scopeData := providerschema.ReplicationScopeData{
			Name:        types.StringValue(s.Name),
			Collections: collections,
		}
		scopeObj, diags := types.ObjectValueFrom(ctx, scopeData.AttributeTypes(), scopeData)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: providerschema.ReplicationScopeData{}.AttributeTypes()}), diags
		}
		scopeObjs = append(scopeObjs, scopeObj)
	}

	listVal, diags := types.ListValue(types.ObjectType{AttrTypes: providerschema.ReplicationScopeData{}.AttributeTypes()}, scopeObjs)
	return listVal, diags
}

// buildMappingsList builds a list of mapping objects.
func buildMappingsList(ctx context.Context, mappings []replicationapi.Mapping) (basetypes.ListValue, diag.Diagnostics) {
	if len(mappings) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: providerschema.MappingData{}.AttributeTypes()}), nil
	}

	var mappingObjs []attr.Value
	for _, m := range mappings {
		// Build collections list for this mapping
		var collectionObjs []attr.Value
		for _, c := range m.Collections {
			collectionData := providerschema.CollectionMappingData{
				SourceCollection: types.StringValue(c.SourceCollection),
				TargetCollection: types.StringValue(c.TargetCollection),
			}
			collectionObj, diags := types.ObjectValueFrom(ctx, collectionData.AttributeTypes(), collectionData)
			if diags.HasError() {
				return types.ListNull(types.ObjectType{AttrTypes: providerschema.MappingData{}.AttributeTypes()}), diags
			}
			collectionObjs = append(collectionObjs, collectionObj)
		}

		var collectionsList basetypes.ListValue
		if len(collectionObjs) > 0 {
			var diags diag.Diagnostics
			collectionsList, diags = types.ListValue(types.ObjectType{AttrTypes: providerschema.CollectionMappingData{}.AttributeTypes()}, collectionObjs)
			if diags.HasError() {
				return types.ListNull(types.ObjectType{AttrTypes: providerschema.MappingData{}.AttributeTypes()}), diags
			}
		} else {
			collectionsList = types.ListNull(types.ObjectType{AttrTypes: providerschema.CollectionMappingData{}.AttributeTypes()})
		}

		mappingData := providerschema.MappingData{
			SourceScope: types.StringValue(m.SourceScope),
			TargetScope: types.StringValue(m.TargetScope),
			Collections: collectionsList,
		}
		mappingObj, diags := types.ObjectValueFrom(ctx, mappingData.AttributeTypes(), mappingData)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: providerschema.MappingData{}.AttributeTypes()}), diags
		}
		mappingObjs = append(mappingObjs, mappingObj)
	}

	listVal, diags := types.ListValue(types.ObjectType{AttrTypes: providerschema.MappingData{}.AttributeTypes()}, mappingObjs)
	return listVal, diags
}

// buildFilterObject builds the filter object from API response.
func buildFilterObject(ctx context.Context, filter *replicationapi.Filter) (basetypes.ObjectValue, diag.Diagnostics) {
	if filter == nil {
		return types.ObjectNull(providerschema.FilterData{}.AttributeTypes()), nil
	}

	var diags diag.Diagnostics

	// Build document exclude options
	var excludeOptionsObj basetypes.ObjectValue
	if filter.DocumentExcludeOptions != nil {
		excludeData := providerschema.DocumentExcludeOptionsData{
			Deletion:   types.BoolValue(filter.DocumentExcludeOptions.Deletion),
			Expiration: types.BoolValue(filter.DocumentExcludeOptions.Expiration),
			Ttl:        types.BoolValue(filter.DocumentExcludeOptions.Ttl),
			Binary:     types.BoolValue(filter.DocumentExcludeOptions.Binary),
		}
		excludeOptionsObj, diags = types.ObjectValueFrom(ctx, excludeData.AttributeTypes(), excludeData)
		if diags.HasError() {
			return types.ObjectNull(providerschema.FilterData{}.AttributeTypes()), diags
		}
	} else {
		excludeOptionsObj = types.ObjectNull(providerschema.DocumentExcludeOptionsData{}.AttributeTypes())
	}

	// Build expressions
	var expressionsObj basetypes.ObjectValue
	if filter.Expressions != nil {
		expressionsData := providerschema.FilterExpressionsData{
			RegEx: types.StringValue(filter.Expressions.RegEx),
		}
		expressionsObj, diags = types.ObjectValueFrom(ctx, expressionsData.AttributeTypes(), expressionsData)
		if diags.HasError() {
			return types.ObjectNull(providerschema.FilterData{}.AttributeTypes()), diags
		}
	} else {
		expressionsObj = types.ObjectNull(providerschema.FilterExpressionsData{}.AttributeTypes())
	}

	filterData := providerschema.FilterData{
		DocumentExcludeOptions: excludeOptionsObj,
		Expressions:            expressionsObj,
	}

	filterObj, diags := types.ObjectValueFrom(ctx, filterData.AttributeTypes(), filterData)
	if diags.HasError() {
		return types.ObjectNull(providerschema.FilterData{}.AttributeTypes()), diags
	}

	return filterObj, nil
}
