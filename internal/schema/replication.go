package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
)

// ReplicationData defines the Terraform state for the singular replication data source.
type ReplicationData struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which to fetch replications.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ReplicationId is the ID of the specific replication to fetch.
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

	// Direction specifies the replication flow — oneWay or twoWay.
	Direction types.String `tfsdk:"direction"`

	// Priority represents the resource allocation to the replication (low, medium, high).
	Priority types.String `tfsdk:"priority"`

	// NetworkUsageLimit is the network usage limit in MiB per second.
	NetworkUsageLimit types.Int64 `tfsdk:"network_usage_limit"`

	// Filter contains the replication settings for the Couchbase server API.
	Filter types.Object `tfsdk:"filter"`

	// Audit contains the audit data for the replication.
	Audit types.Object `tfsdk:"audit"`
}

// ReplicationSourceData defines the source metadata for a replication.
type ReplicationSourceData struct {
	Project types.Object `tfsdk:"project"`
	Cluster types.Object `tfsdk:"cluster"`
	Bucket  types.Object `tfsdk:"bucket"`
	Scopes  types.List   `tfsdk:"scopes"`
	Type    types.String `tfsdk:"type"`
}

// ReplicationTargetData defines the target metadata for a replication.
type ReplicationTargetData struct {
	Project types.Object `tfsdk:"project"`
	Cluster types.Object `tfsdk:"cluster"`
	Bucket  types.Object `tfsdk:"bucket"`
	Scopes  types.List   `tfsdk:"scopes"`
	Type    types.String `tfsdk:"type"`
}

// ProjectReferenceData contains project metadata.
type ProjectReferenceData struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// ClusterReferenceData contains cluster metadata.
type ClusterReferenceData struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// BucketReferenceData contains bucket metadata.
type BucketReferenceData struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ConflictResolutionType types.String `tfsdk:"conflict_resolution_type"`
}

// ReplicationScopeData contains scope and collection details.
type ReplicationScopeData struct {
	Name        types.String   `tfsdk:"name"`
	Collections []types.String `tfsdk:"collections"`
}

// MappingData defines a mapping from source to target scope/collection.
type MappingData struct {
	SourceScope types.String `tfsdk:"source_scope"`
	TargetScope types.String `tfsdk:"target_scope"`
	Collections types.List   `tfsdk:"collections"`
}

// CollectionMappingData defines a mapping between source and target collections.
type CollectionMappingData struct {
	SourceCollection types.String `tfsdk:"source_collection"`
	TargetCollection types.String `tfsdk:"target_collection"`
}

// FilterData contains the replication filter settings.
type FilterData struct {
	DocumentExcludeOptions types.Object `tfsdk:"document_exclude_options"`
	Expressions            types.Object `tfsdk:"expressions"`
}

// DocumentExcludeOptionsData specifies document types to filter out.
type DocumentExcludeOptionsData struct {
	Deletion   types.Bool `tfsdk:"deletion"`
	Expiration types.Bool `tfsdk:"expiration"`
	Ttl        types.Bool `tfsdk:"ttl"`
	Binary     types.Bool `tfsdk:"binary"`
}

// FilterExpressionsData contains filter expression settings.
type FilterExpressionsData struct {
	RegEx types.String `tfsdk:"reg_ex"`
}

// ReplicationAuditData defines the audit information for a replication.
type ReplicationAuditData struct {
	CreatedBy  types.String `tfsdk:"created_by"`
	CreatedAt  types.String `tfsdk:"created_at"`
	ModifiedBy types.String `tfsdk:"modified_by"`
	ModifiedAt types.String `tfsdk:"modified_at"`
	Version    types.Int64  `tfsdk:"version"`
}

// ReplicationsData defines the Terraform state for the list replications data source.
type ReplicationsData struct {
	// OrganizationId is the ID of the organization.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Data is the list of replications.
	Data []ReplicationSummaryData `tfsdk:"data"`
}

// ReplicationSummaryData defines the summary data for a replication in list response.
type ReplicationSummaryData struct {
	Id            types.String `tfsdk:"id"`
	SourceCluster types.String `tfsdk:"source_cluster"`
	TargetCluster types.String `tfsdk:"target_cluster"`
	Status        types.String `tfsdk:"status"`
	Direction     types.String `tfsdk:"direction"`
	Audit         types.Object `tfsdk:"audit"`
}

// AttributeTypes returns the attribute types for ProjectReferenceData.
func (p ProjectReferenceData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

// AttributeTypes returns the attribute types for ClusterReferenceData.
func (c ClusterReferenceData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

// AttributeTypes returns the attribute types for BucketReferenceData.
func (b BucketReferenceData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                       types.StringType,
		"name":                     types.StringType,
		"conflict_resolution_type": types.StringType,
	}
}

// AttributeTypes returns the attribute types for ReplicationScopeData.
func (s ReplicationScopeData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"collections": types.ListType{ElemType: types.StringType},
	}
}

// AttributeTypes returns the attribute types for CollectionMappingData.
func (c CollectionMappingData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source_collection": types.StringType,
		"target_collection": types.StringType,
	}
}

// AttributeTypes returns the attribute types for MappingData.
func (m MappingData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source_scope": types.StringType,
		"target_scope": types.StringType,
		"collections":  types.ListType{ElemType: types.ObjectType{AttrTypes: CollectionMappingData{}.AttributeTypes()}},
	}
}

// AttributeTypes returns the attribute types for DocumentExcludeOptionsData.
func (d DocumentExcludeOptionsData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"deletion":   types.BoolType,
		"expiration": types.BoolType,
		"ttl":        types.BoolType,
		"binary":     types.BoolType,
	}
}

// AttributeTypes returns the attribute types for FilterExpressionsData.
func (f FilterExpressionsData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"reg_ex": types.StringType,
	}
}

// AttributeTypes returns the attribute types for FilterData.
func (f FilterData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"document_exclude_options": types.ObjectType{AttrTypes: DocumentExcludeOptionsData{}.AttributeTypes()},
		"expressions":              types.ObjectType{AttrTypes: FilterExpressionsData{}.AttributeTypes()},
	}
}

// AttributeTypes returns the attribute types for ReplicationAuditData.
func (r ReplicationAuditData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_by":  types.StringType,
		"created_at":  types.StringType,
		"modified_by": types.StringType,
		"modified_at": types.StringType,
		"version":     types.Int64Type,
	}
}

// AttributeTypes returns the attribute types for ReplicationSourceData.
func (r ReplicationSourceData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project": types.ObjectType{AttrTypes: ProjectReferenceData{}.AttributeTypes()},
		"cluster": types.ObjectType{AttrTypes: ClusterReferenceData{}.AttributeTypes()},
		"bucket":  types.ObjectType{AttrTypes: BucketReferenceData{}.AttributeTypes()},
		"scopes":  types.ListType{ElemType: types.ObjectType{AttrTypes: ReplicationScopeData{}.AttributeTypes()}},
		"type":    types.StringType,
	}
}

// AttributeTypes returns the attribute types for ReplicationTargetData.
func (r ReplicationTargetData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project": types.ObjectType{AttrTypes: ProjectReferenceData{}.AttributeTypes()},
		"cluster": types.ObjectType{AttrTypes: ClusterReferenceData{}.AttributeTypes()},
		"bucket":  types.ObjectType{AttrTypes: BucketReferenceData{}.AttributeTypes()},
		"scopes":  types.ListType{ElemType: types.ObjectType{AttrTypes: ReplicationScopeData{}.AttributeTypes()}},
		"type":    types.StringType,
	}
}

// AttributeTypes returns the attribute types for ReplicationSummaryData.
func (r ReplicationSummaryData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             types.StringType,
		"source_cluster": types.StringType,
		"target_cluster": types.StringType,
		"status":         types.StringType,
		"direction":      types.StringType,
		"audit":          types.ObjectType{AttrTypes: ReplicationAuditData{}.AttributeTypes()},
	}
}

// NewReplicationData creates a new ReplicationData from API response.
func NewReplicationData(
	organizationId, projectId, clusterId, replicationId string,
	apiResponse *replicationapi.GetReplicationResponse,
	sourceObj, targetObj, filterObj types.Object,
	mappingsList types.List,
	auditObj types.Object,
) *ReplicationData {
	return &ReplicationData{
		OrganizationId:    types.StringValue(organizationId),
		ProjectId:         types.StringValue(projectId),
		ClusterId:         types.StringValue(clusterId),
		ReplicationId:     types.StringValue(replicationId),
		Id:                types.StringValue(apiResponse.Id),
		Status:            types.StringValue(apiResponse.Status),
		ChangesLeft:       types.Int64Value(apiResponse.ChangesLeft),
		Error:             types.StringValue(apiResponse.Error),
		Direction:         types.StringValue(apiResponse.Direction),
		Priority:          types.StringValue(apiResponse.Priority),
		NetworkUsageLimit: types.Int64Value(apiResponse.NetworkUsageLimit),
		Source:            sourceObj,
		Target:            targetObj,
		Mappings:          mappingsList,
		Filter:            filterObj,
		Audit:             auditObj,
	}
}

// NewReplicationSummaryData creates a new ReplicationSummaryData from API response.
func NewReplicationSummaryData(
	apiResponse replicationapi.ReplicationSummary,
	auditObj types.Object,
) ReplicationSummaryData {
	return ReplicationSummaryData{
		Id:            types.StringValue(apiResponse.Id),
		SourceCluster: types.StringValue(apiResponse.SourceCluster),
		TargetCluster: types.StringValue(apiResponse.TargetCluster),
		Status:        types.StringValue(apiResponse.Status),
		Direction:     types.StringValue(apiResponse.Direction),
		Audit:         auditObj,
	}
}
