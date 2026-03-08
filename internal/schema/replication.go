package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	replication_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
)

// ReplicationData defines the Terraform state for the replication datasource.
type ReplicationData struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster from which the replication is established.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ReplicationId is the ID of the replication.
	ReplicationId types.String `tfsdk:"replication_id"`

	// Id is the ID of the specified replication.
	Id types.String `tfsdk:"id"`

	// Status is the status of the replication (pending, pausing, failed, paused, running).
	Status types.String `tfsdk:"status"`

	// ChangesLeft is the number of remaining mutations to the replication.
	ChangesLeft types.Int64 `tfsdk:"changes_left"`

	// Error is the error message if the replication has failed.
	Error types.String `tfsdk:"error"`

	// Source contains all the metadata about a replication source.
	Source types.Object `tfsdk:"source"`

	// Target contains all the metadata about a replication target.
	Target types.Object `tfsdk:"target"`

	// Mappings defines mappings from source to target scopes and collections.
	Mappings types.List `tfsdk:"mappings"`

	// Direction specifies the replication flow — whether it's oneWay or twoWay.
	Direction types.String `tfsdk:"direction"`

	// Priority represents the resource allocation to the replication (low, medium, high).
	Priority types.String `tfsdk:"priority"`

	// NetworkUsageLimit is the network usage limit in MiB per second. 0 means unlimited.
	NetworkUsageLimit types.Int64 `tfsdk:"network_usage_limit"`

	// Filter contains the replication filter settings.
	Filter types.Object `tfsdk:"filter"`

	// Audit contains the audit data for the replication.
	Audit types.Object `tfsdk:"audit"`
}

// Validate validates the replication datasource state and returns parsed IDs.
func (r *ReplicationData) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: r.OrganizationId,
		ProjectId:      r.ProjectId,
		ClusterId:      r.ClusterId,
		Id:             r.ReplicationId,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("failed to validate datasource state: %w", err)
	}

	return IDs, nil
}

// NewReplicationData creates a new ReplicationData from API response data.
func NewReplicationData(
	organizationId, projectId, clusterId, replicationId string,
	apiResponse *apigen.GetReplicationResponse,
) *ReplicationData {
	result := &ReplicationData{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		ReplicationId:  types.StringValue(replicationId),
		Id:             types.StringValue(apiResponse.Id),
		Status:         types.StringValue(string(apiResponse.Status)),
		ChangesLeft:    types.Int64Value(int64(apiResponse.ChangesLeft)),
		Direction:      types.StringValue(string(apiResponse.Direction)),
	}

	// Error is optional
	if apiResponse.Error != nil {
		result.Error = types.StringValue(*apiResponse.Error)
	} else {
		result.Error = types.StringNull()
	}

	// Priority is optional
	if apiResponse.Priority != nil {
		result.Priority = types.StringValue(string(*apiResponse.Priority))
	} else {
		result.Priority = types.StringNull()
	}

	// NetworkUsageLimit is optional
	if apiResponse.NetworkUsageLimit != nil {
		result.NetworkUsageLimit = types.Int64Value(int64(*apiResponse.NetworkUsageLimit))
	} else {
		result.NetworkUsageLimit = types.Int64Null()
	}

	return result
}

// SetNestedFields sets the nested fields (Source, Target, Mappings, Filter, Audit) from the API response.
func (r *ReplicationData) SetNestedFields(ctx context.Context, apiResponse *apigen.GetReplicationResponse) error {
	// Source
	source, err := buildSourceFromAPI(apiResponse.Source)
	if err != nil {
		return fmt.Errorf("failed to build source: %w", err)
	}
	sourceObj, diags := types.ObjectValue(sourceAttributeTypes(), source)
	if diags.HasError() {
		return fmt.Errorf("failed to create source object: %v", diags)
	}
	r.Source = sourceObj

	// Target
	target, err := buildTargetFromAPI(apiResponse.Target)
	if err != nil {
		return fmt.Errorf("failed to build target: %w", err)
	}
	targetObj, diags := types.ObjectValue(targetAttributeTypes(), target)
	if diags.HasError() {
		return fmt.Errorf("failed to create target object: %v", diags)
	}
	r.Target = targetObj

	// Mappings
	if apiResponse.Mappings != nil && len(*apiResponse.Mappings) > 0 {
		mappings, err := buildMappingsFromAPI(*apiResponse.Mappings)
		if err != nil {
			return fmt.Errorf("failed to build mappings: %w", err)
		}
		mappingsList, diags := types.ListValue(types.ObjectType{AttrTypes: mappingAttributeTypes()}, mappings)
		if diags.HasError() {
			return fmt.Errorf("failed to create mappings list: %v", diags)
		}
		r.Mappings = mappingsList
	} else {
		r.Mappings = types.ListNull(types.ObjectType{AttrTypes: mappingAttributeTypes()})
	}

	// Filter
	if apiResponse.Filter != nil {
		filter, err := buildFilterFromAPI(*apiResponse.Filter)
		if err != nil {
			return fmt.Errorf("failed to build filter: %w", err)
		}
		filterObj, diags := types.ObjectValue(filterAttributeTypes(), filter)
		if diags.HasError() {
			return fmt.Errorf("failed to create filter object: %v", diags)
		}
		r.Filter = filterObj
	} else {
		r.Filter = types.ObjectNull(filterAttributeTypes())
	}

	// Audit (replication uses a simpler audit structure with only created_at and created_by)
	audit := ReplicationAuditData{
		CreatedAt: types.StringValue(apiResponse.Audit.CreatedAt.String()),
		CreatedBy: types.StringValue(apiResponse.Audit.CreatedBy),
	}
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return fmt.Errorf("failed to create audit object: %v", diags)
	}
	r.Audit = auditObj

	return nil
}

// Helper functions for building nested objects

func buildSourceFromAPI(source apigen.ReplicationSource) (map[string]attr.Value, error) {
	result := map[string]attr.Value{
		"project_id":               types.StringValue(source.Project.Id),
		"project_name":             types.StringValue(source.Project.Name),
		"cluster_id":               types.StringValue(source.Cluster.Id),
		"cluster_name":             types.StringValue(source.Cluster.Name),
		"bucket_id":                types.StringValue(source.Bucket.Id),
		"bucket_name":              types.StringValue(source.Bucket.Name),
		"bucket_conflict_resolution": types.StringValue(source.Bucket.ConflictResolutionType),
	}

	if source.Type != nil {
		result["type"] = types.StringValue(string(*source.Type))
	} else {
		result["type"] = types.StringNull()
	}

	// Handle scopes
	if source.Scopes != nil && len(*source.Scopes) > 0 {
		scopes, err := buildScopesFromAPI(*source.Scopes)
		if err != nil {
			return nil, err
		}
		scopesList, diags := types.ListValue(types.ObjectType{AttrTypes: scopeAttributeTypes()}, scopes)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to create scopes list: %v", diags)
		}
		result["scopes"] = scopesList
	} else {
		result["scopes"] = types.ListNull(types.ObjectType{AttrTypes: scopeAttributeTypes()})
	}

	return result, nil
}

func buildTargetFromAPI(target apigen.ReplicationTarget) (map[string]attr.Value, error) {
	result := map[string]attr.Value{
		"cluster_id":   types.StringValue(target.Cluster.Id),
		"cluster_name": types.StringValue(target.Cluster.Name),
		"bucket_id":    types.StringValue(target.Bucket.Id),
		"bucket_name":  types.StringValue(target.Bucket.Name),
		"type":         types.StringValue(string(target.Type)),
	}

	// Bucket conflict resolution is optional for target
	if target.Bucket.ConflictResolutionType != nil {
		result["bucket_conflict_resolution"] = types.StringValue(*target.Bucket.ConflictResolutionType)
	} else {
		result["bucket_conflict_resolution"] = types.StringNull()
	}

	// Project is optional for target
	if target.Project != nil {
		if target.Project.Id != nil {
			result["project_id"] = types.StringValue(*target.Project.Id)
		} else {
			result["project_id"] = types.StringNull()
		}
		if target.Project.Name != nil {
			result["project_name"] = types.StringValue(*target.Project.Name)
		} else {
			result["project_name"] = types.StringNull()
		}
	} else {
		result["project_id"] = types.StringNull()
		result["project_name"] = types.StringNull()
	}

	// Handle scopes
	if target.Scopes != nil && len(*target.Scopes) > 0 {
		scopes, err := buildScopesFromAPI(*target.Scopes)
		if err != nil {
			return nil, err
		}
		scopesList, diags := types.ListValue(types.ObjectType{AttrTypes: scopeAttributeTypes()}, scopes)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to create scopes list: %v", diags)
		}
		result["scopes"] = scopesList
	} else {
		result["scopes"] = types.ListNull(types.ObjectType{AttrTypes: scopeAttributeTypes()})
	}

	return result, nil
}

func buildScopesFromAPI(scopes apigen.Scopes) ([]attr.Value, error) {
	var result []attr.Value

	for _, scope := range scopes {
		scopeMap := map[string]attr.Value{
			"name": types.StringNull(),
		}
		if scope.Name != nil {
			scopeMap["name"] = types.StringValue(*scope.Name)
		}

		// Handle collections
		if scope.Collections != nil && len(*scope.Collections) > 0 {
			collections, err := buildCollectionsFromAPI(*scope.Collections)
			if err != nil {
				return nil, err
			}
			collectionsList, diags := types.ListValue(types.StringType, collections)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to create collections list: %v", diags)
			}
			scopeMap["collections"] = collectionsList
		} else {
			scopeMap["collections"] = types.ListNull(types.StringType)
		}

		scopeObj, diags := types.ObjectValue(scopeAttributeTypes(), scopeMap)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to create scope object: %v", diags)
		}
		result = append(result, scopeObj)
	}

	return result, nil
}

func buildCollectionsFromAPI(collections apigen.Collections) ([]attr.Value, error) {
	var result []attr.Value
	for _, collection := range collections {
		result = append(result, types.StringValue(collection))
	}
	return result, nil
}

func buildMappingsFromAPI(mappings apigen.Mappings) ([]attr.Value, error) {
	var result []attr.Value

	for _, mapping := range mappings {
		mappingMap := map[string]attr.Value{
			"source_scope": types.StringValue(mapping.SourceScope),
			"target_scope": types.StringValue(mapping.TargetScope),
		}

		// Handle collections
		if mapping.Collections != nil && len(*mapping.Collections) > 0 {
			var collections []attr.Value
			for _, collection := range *mapping.Collections {
				collectionMap := map[string]attr.Value{
					"source_collection": types.StringValue(collection.SourceCollection),
					"target_collection": types.StringValue(collection.TargetCollection),
				}
				collectionObj, diags := types.ObjectValue(collectionMappingAttributeTypes(), collectionMap)
				if diags.HasError() {
					return nil, fmt.Errorf("failed to create collection object: %v", diags)
				}
				collections = append(collections, collectionObj)
			}
			collectionsList, diags := types.ListValue(types.ObjectType{AttrTypes: collectionMappingAttributeTypes()}, collections)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to create collections list: %v", diags)
			}
			mappingMap["collections"] = collectionsList
		} else {
			mappingMap["collections"] = types.ListNull(types.ObjectType{AttrTypes: collectionMappingAttributeTypes()})
		}

		mappingObj, diags := types.ObjectValue(mappingAttributeTypes(), mappingMap)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to create mapping object: %v", diags)
		}
		result = append(result, mappingObj)
	}

	return result, nil
}

func buildFilterFromAPI(filter apigen.GetFilter) (map[string]attr.Value, error) {
	result := map[string]attr.Value{}

	// Document exclude options
	if filter.DocumentExcludeOptions != nil {
		docOpts := map[string]attr.Value{
			"binary":     types.BoolValue(filter.DocumentExcludeOptions.Binary != nil && *filter.DocumentExcludeOptions.Binary),
			"deletion":   types.BoolValue(filter.DocumentExcludeOptions.Deletion != nil && *filter.DocumentExcludeOptions.Deletion),
			"expiration": types.BoolValue(filter.DocumentExcludeOptions.Expiration != nil && *filter.DocumentExcludeOptions.Expiration),
			"ttl":        types.BoolValue(filter.DocumentExcludeOptions.Ttl != nil && *filter.DocumentExcludeOptions.Ttl),
		}
		docOptsObj, diags := types.ObjectValue(documentExcludeOptionsAttributeTypes(), docOpts)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to create document_exclude_options object: %v", diags)
		}
		result["document_exclude_options"] = docOptsObj
	} else {
		result["document_exclude_options"] = types.ObjectNull(documentExcludeOptionsAttributeTypes())
	}

	// Expressions
	if filter.Expressions != nil {
		exprs := map[string]attr.Value{}
		if filter.Expressions.RegEx != nil {
			exprs["regex"] = types.StringValue(*filter.Expressions.RegEx)
		} else {
			exprs["regex"] = types.StringNull()
		}
		exprsObj, diags := types.ObjectValue(expressionsAttributeTypes(), exprs)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to create expressions object: %v", diags)
		}
		result["expressions"] = exprsObj
	} else {
		result["expressions"] = types.ObjectNull(expressionsAttributeTypes())
	}

	return result, nil
}

// ReplicationAuditData contains the simplified audit fields for replication.
type ReplicationAuditData struct {
	// CreatedAt The RFC3339 timestamp associated with when the replication was initially created.
	CreatedAt types.String `tfsdk:"created_at"`

	// CreatedBy The user who created the replication.
	CreatedBy types.String `tfsdk:"created_by"`
}

// AttributeTypes returns the attribute types for ReplicationAuditData.
func (r ReplicationAuditData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_at": types.StringType,
		"created_by": types.StringType,
	}
}

// Attribute type functions

func sourceAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project_id":                 types.StringType,
		"project_name":               types.StringType,
		"cluster_id":                 types.StringType,
		"cluster_name":               types.StringType,
		"bucket_id":                  types.StringType,
		"bucket_name":                types.StringType,
		"bucket_conflict_resolution": types.StringType,
		"scopes":                     types.ListType{ElemType: types.ObjectType{AttrTypes: scopeAttributeTypes()}},
		"type":                       types.StringType,
	}
}

func targetAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project_id":                 types.StringType,
		"project_name":             types.StringType,
		"cluster_id":                 types.StringType,
		"cluster_name":               types.StringType,
		"bucket_id":                  types.StringType,
		"bucket_name":                types.StringType,
		"bucket_conflict_resolution": types.StringType,
		"scopes":                     types.ListType{ElemType: types.ObjectType{AttrTypes: scopeAttributeTypes()}},
		"type":                       types.StringType,
	}
}

func scopeAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"collections": types.ListType{ElemType: types.StringType},
	}
}

func mappingAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source_scope": types.StringType,
		"target_scope": types.StringType,
		"collections":  types.ListType{ElemType: types.ObjectType{AttrTypes: collectionMappingAttributeTypes()}},
	}
}

func collectionMappingAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source_collection": types.StringType,
		"target_collection": types.StringType,
	}
}

func filterAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"document_exclude_options": types.ObjectType{AttrTypes: documentExcludeOptionsAttributeTypes()},
		"expressions":            types.ObjectType{AttrTypes: expressionsAttributeTypes()},
	}
}

func documentExcludeOptionsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"binary":     types.BoolType,
		"deletion":   types.BoolType,
		"expiration": types.BoolType,
		"ttl":        types.BoolType,
	}
}

func expressionsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"regex": types.StringType,
	}
}

// ReplicationSummaryData defines attributes for a single replication summary when fetched from the list API.
type ReplicationSummaryData struct {
	// Id is the ID of the specified replication.
	Id types.String `tfsdk:"id"`

	// SourceCluster is the name of the source cluster.
	SourceCluster types.String `tfsdk:"source_cluster"`

	// TargetCluster is the name of the target cluster.
	TargetCluster types.String `tfsdk:"target_cluster"`

	// Status is the status of the replication.
	Status types.String `tfsdk:"status"`

	// Direction specifies the replication flow.
	Direction types.String `tfsdk:"direction"`

	// Audit contains the audit data for the replication.
	Audit types.Object `tfsdk:"audit"`
}

// NewReplicationSummaryData creates a new ReplicationSummaryData from API response data.
func NewReplicationSummaryData(
	apiResponse replication_api.GetReplicationSummaryResponse,
	auditObj basetypes.ObjectValue,
) *ReplicationSummaryData {
	result := &ReplicationSummaryData{
		Id:            types.StringValue(apiResponse.Id),
		SourceCluster: types.StringValue(apiResponse.SourceCluster),
		TargetCluster: types.StringValue(apiResponse.TargetCluster),
		Status:        types.StringValue(apiResponse.Status),
		Audit:         auditObj,
	}

	if apiResponse.Direction != nil {
		result.Direction = types.StringValue(*apiResponse.Direction)
	} else {
		result.Direction = types.StringNull()
	}

	return result
}

// ReplicationsData defines structure based on the response received from list replications API.
type ReplicationsData struct {
	// OrganizationId is the ID of the organization.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Data contains the list of replication summaries.
	Data []ReplicationSummaryData `tfsdk:"data"`
}

// Validate is used to verify that all the fields in the datasource have been populated.
func (r ReplicationsData) Validate() (clusterId, projectId, organizationId string, err error) {
	if r.OrganizationId.IsNull() {
		return "", "", "", errors.ErrOrganizationIdMissing
	}
	if r.ProjectId.IsNull() {
		return "", "", "", errors.ErrProjectIdMissing
	}
	if r.ClusterId.IsNull() {
		return "", "", "", errors.ErrClusterIdMissing
	}
	return r.ClusterId.ValueString(), r.ProjectId.ValueString(), r.OrganizationId.ValueString(), nil
}
