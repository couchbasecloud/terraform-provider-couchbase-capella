package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

// Replication maps an XDCR replication between two Capella buckets onto Terraform state.
type Replication struct {
	// OrganizationId is the organization to which the source cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the project to which the source cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the source cluster of the replication.
	ClusterId types.String `tfsdk:"cluster_id"`

	// SourceBucket is the ID of the source bucket, not its name.
	SourceBucket types.String `tfsdk:"source_bucket"`

	// Target is the replication destination.
	Target *ReplicationTarget `tfsdk:"target"`

	// Direction is one of oneWay or twoWay.
	Direction types.String `tfsdk:"direction"`

	// Priority is one of low, medium or high.
	Priority types.String `tfsdk:"priority"`

	// NetworkUsageLimit is the bandwidth limit in MiB/s. 0 means unlimited.
	NetworkUsageLimit types.Int64 `tfsdk:"network_usage_limit"`

	// Mappings restricts replication to specific scopes and collections.
	// Empty means the whole bucket is replicated.
	Mappings []ReplicationMapping `tfsdk:"mappings"`

	// Filter excludes documents or restricts which mutations replicate.
	Filter *ReplicationFilter `tfsdk:"filter"`

	// Id is the ID of the replication.
	Id types.String `tfsdk:"id"`

	// ReverseReplicationId is set only when Direction is twoWay.
	ReverseReplicationId types.String `tfsdk:"reverse_replication_id"`

	// Status is one of pending, pausing, failed, paused or running.
	Status types.String `tfsdk:"status"`

	// ChangesLeft is the number of mutations still to replicate.
	ChangesLeft types.Int64 `tfsdk:"changes_left"`

	// Error is set when Status is failed.
	Error types.String `tfsdk:"error"`

	// SourceDetails is the resolved metadata for the replication source.
	SourceDetails *ReplicationEndpointDetails `tfsdk:"source_details"`

	// TargetDetails is the resolved metadata for the replication target.
	TargetDetails *ReplicationEndpointDetails `tfsdk:"target_details"`

	// Audit contains the replication's creation audit data.
	Audit ReplicationAuditData `tfsdk:"audit"`
}

// ReplicationTarget identifies the destination bucket of a replication.
type ReplicationTarget struct {
	// Cluster is the target cluster ID for capella targets, or the cluster reference name for external targets.
	Cluster types.String `tfsdk:"cluster"`

	// Bucket is the target bucket ID for capella targets, or the bucket name for external targets.
	Bucket types.String `tfsdk:"bucket"`

	// Type is one of capella or external.
	Type types.String `tfsdk:"type"`
}

func (r ReplicationTarget) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cluster": types.StringType,
		"bucket":  types.StringType,
		"type":    types.StringType,
	}
}

// ReplicationMapping maps one source scope onto one target scope.
type ReplicationMapping struct {
	SourceScope types.String `tfsdk:"source_scope"`
	TargetScope types.String `tfsdk:"target_scope"`

	// Collections restricts the mapping to specific collections.
	// Empty means every collection under the scope is replicated.
	Collections []ReplicationCollection `tfsdk:"collections"`
}

func (r ReplicationMapping) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source_scope": types.StringType,
		"target_scope": types.StringType,
		"collections":  types.ListType{ElemType: types.ObjectType{AttrTypes: ReplicationCollection{}.AttributeTypes()}},
	}
}

// ReplicationCollection maps one source collection onto one target collection.
type ReplicationCollection struct {
	SourceCollection types.String `tfsdk:"source_collection"`
	TargetCollection types.String `tfsdk:"target_collection"`
}

func (r ReplicationCollection) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source_collection": types.StringType,
		"target_collection": types.StringType,
	}
}

// ReplicationFilter holds the rules that exclude documents from a replication.
type ReplicationFilter struct {
	DocumentExcludeOptions *ReplicationDocumentExcludeOptions `tfsdk:"document_exclude_options"`
	Expressions            *ReplicationFilterExpressions      `tfsdk:"expressions"`
}

func (r ReplicationFilter) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"document_exclude_options": types.ObjectType{AttrTypes: ReplicationDocumentExcludeOptions{}.AttributeTypes()},
		"expressions":              types.ObjectType{AttrTypes: ReplicationFilterExpressions{}.AttributeTypes()},
	}
}

// ReplicationDocumentExcludeOptions lists the document classes excluded from a replication.
type ReplicationDocumentExcludeOptions struct {
	Deletion   types.Bool `tfsdk:"deletion"`
	Expiration types.Bool `tfsdk:"expiration"`
	Ttl        types.Bool `tfsdk:"ttl"`
	Binary     types.Bool `tfsdk:"binary"`
}

func (r ReplicationDocumentExcludeOptions) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"deletion":   types.BoolType,
		"expiration": types.BoolType,
		"ttl":        types.BoolType,
		"binary":     types.BoolType,
	}
}

// ReplicationFilterExpressions holds the regular expression applied to replicated documents.
type ReplicationFilterExpressions struct {
	RegEx types.String `tfsdk:"reg_ex"`

	// SkipRestream controls whether an updated filter restarts the replication.
	// The API never returns it, so it is held from configuration only.
	SkipRestream types.Bool `tfsdk:"skip_restream"`
}

func (r ReplicationFilterExpressions) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"reg_ex":        types.StringType,
		"skip_restream": types.BoolType,
	}
}

// ReplicationEndpointDetails is the resolved metadata the API returns for either end of a replication.
type ReplicationEndpointDetails struct {
	Project *ReplicationProject `tfsdk:"project"`
	Cluster *ReplicationCluster `tfsdk:"cluster"`
	Bucket  *ReplicationBucket  `tfsdk:"bucket"`
	Scopes  []ReplicationScope  `tfsdk:"scopes"`
	Type    types.String        `tfsdk:"type"`
}

func (r ReplicationEndpointDetails) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project": types.ObjectType{AttrTypes: ReplicationProject{}.AttributeTypes()},
		"cluster": types.ObjectType{AttrTypes: ReplicationCluster{}.AttributeTypes()},
		"bucket":  types.ObjectType{AttrTypes: ReplicationBucket{}.AttributeTypes()},
		"scopes":  types.ListType{ElemType: types.ObjectType{AttrTypes: ReplicationScope{}.AttributeTypes()}},
		"type":    types.StringType,
	}
}

// ReplicationProject identifies the project owning one end of a replication.
type ReplicationProject struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r ReplicationProject) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

// ReplicationCluster identifies the cluster at one end of a replication.
type ReplicationCluster struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r ReplicationCluster) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

// ReplicationBucket identifies the bucket at one end of a replication.
type ReplicationBucket struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ConflictResolutionType types.String `tfsdk:"conflict_resolution_type"`
}

func (r ReplicationBucket) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                       types.StringType,
		"name":                     types.StringType,
		"conflict_resolution_type": types.StringType,
	}
}

// ReplicationScope is a scope, and optionally its collections, taking part in a replication.
type ReplicationScope struct {
	Name        types.String   `tfsdk:"name"`
	Collections []types.String `tfsdk:"collections"`
}

func (r ReplicationScope) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"collections": types.ListType{ElemType: types.StringType},
	}
}

// ReplicationAuditData is the audit data returned for a replication.
// Unlike CouchbaseAuditData it carries creation fields only; the replication API returns no
// modification or version fields.
type ReplicationAuditData struct {
	CreatedAt types.String `tfsdk:"created_at"`
	CreatedBy types.String `tfsdk:"created_by"`
}

func (r ReplicationAuditData) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_at": types.StringType,
		"created_by": types.StringType,
	}
}

func (r Replication) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id":        types.StringType,
		"project_id":             types.StringType,
		"cluster_id":             types.StringType,
		"source_bucket":          types.StringType,
		"target":                 types.ObjectType{AttrTypes: ReplicationTarget{}.AttributeTypes()},
		"direction":              types.StringType,
		"priority":               types.StringType,
		"network_usage_limit":    types.Int64Type,
		"mappings":               types.ListType{ElemType: types.ObjectType{AttrTypes: ReplicationMapping{}.AttributeTypes()}},
		"filter":                 types.ObjectType{AttrTypes: ReplicationFilter{}.AttributeTypes()},
		"id":                     types.StringType,
		"reverse_replication_id": types.StringType,
		"status":                 types.StringType,
		"changes_left":           types.Int64Type,
		"error":                  types.StringType,
		"source_details":         types.ObjectType{AttrTypes: ReplicationEndpointDetails{}.AttributeTypes()},
		"target_details":         types.ObjectType{AttrTypes: ReplicationEndpointDetails{}.AttributeTypes()},
		"audit":                  types.ObjectType{AttrTypes: ReplicationAuditData{}.AttributeTypes()},
	}
}

// Validate is used to verify that IDs have been properly imported.
func (r Replication) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: r.OrganizationId,
		ProjectId:      r.ProjectId,
		ClusterId:      r.ClusterId,
		Id:             r.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}
	return IDs, nil
}

// NewReplication converts a replication read from the API into Terraform state.
//
// skip_restream is deliberately left unset: the API's read model has no such field, so the value
// has to be carried over from configuration by the caller rather than being zeroed here.
func NewReplication(
	replication apigen.GetReplicationResponse,
	organizationId, projectId, clusterId string,
) *Replication {
	return &Replication{
		OrganizationId:       types.StringValue(organizationId),
		ProjectId:            types.StringValue(projectId),
		ClusterId:            types.StringValue(clusterId),
		SourceBucket:         types.StringValue(replication.Source.Bucket.Id),
		Target:               newReplicationTargetRef(replication.Target),
		Direction:            types.StringValue(string(replication.Direction)),
		Priority:             types.StringPointerValue((*string)(replication.Priority)),
		NetworkUsageLimit:    newOptionalInt64(replication.NetworkUsageLimit),
		Mappings:             newReplicationMappings(replication.Mappings),
		Filter:               newReplicationFilter(replication.Filter),
		Id:                   types.StringValue(replication.Id),
		ReverseReplicationId: types.StringNull(),
		Status:               types.StringValue(string(replication.Status)),
		ChangesLeft:          types.Int64Value(int64(replication.ChangesLeft)),
		Error:                types.StringPointerValue(replication.Error),
		SourceDetails:        newReplicationSourceDetails(replication.Source),
		TargetDetails:        newReplicationTargetDetails(replication.Target),
		Audit:                newReplicationAuditData(replication.Audit),
	}
}

// ToCreateRequest builds the create payload. The provider only ever creates replications
// asynchronously, so mode is always async.
func (r Replication) ToCreateRequest() apigen.CreateReplicationJSONRequestBody {
	// apigen.Async is the generator's unprefixed constant for CreateReplicationRequestMode "async".
	mode := apigen.Async
	request := apigen.CreateReplicationJSONRequestBody{
		SourceBucket:      r.SourceBucket.ValueString(),
		Mode:              &mode,
		NetworkUsageLimit: int64ToIntPointer(r.NetworkUsageLimit),
		Mappings:          toAPIMappings(r.Mappings),
		Filter:            toAPIFilter(r.Filter),
	}

	if r.Target != nil {
		request.Target.Cluster = r.Target.Cluster.ValueString()
		request.Target.Bucket = r.Target.Bucket.ValueString()
		if !r.Target.Type.IsNull() && !r.Target.Type.IsUnknown() {
			targetType := apigen.CreateReplicationRequestTargetType(r.Target.Type.ValueString())
			request.Target.Type = &targetType
		}
	}

	if !r.Direction.IsNull() && !r.Direction.IsUnknown() {
		direction := apigen.CreateReplicationRequestDirection(r.Direction.ValueString())
		request.Direction = &direction
	}

	if !r.Priority.IsNull() && !r.Priority.IsUnknown() {
		priority := apigen.CreateReplicationRequestPriority(r.Priority.ValueString())
		request.Priority = &priority
	}

	return request
}

// ToUpdateRequest builds the update payload. Fields omitted by the API are left unchanged, so
// allScopes is sent when the configuration has no mappings in order to switch a replication back
// to replicating the whole bucket.
func (r Replication) ToUpdateRequest() apigen.UpdateReplicationJSONRequestBody {
	request := apigen.UpdateReplicationJSONRequestBody{
		NetworkUsageLimit: int64ToIntPointer(r.NetworkUsageLimit),
		Mappings:          toAPIMappings(r.Mappings),
		Filter:            toAPIFilter(r.Filter),
	}

	if len(r.Mappings) == 0 {
		allScopes := true
		request.AllScopes = &allScopes
	}

	if !r.Priority.IsNull() && !r.Priority.IsUnknown() {
		priority := apigen.UpdateReplicationRequestPriority(r.Priority.ValueString())
		request.Priority = &priority
	}

	return request
}

// newReplicationTargetRef rebuilds the configured target block from the API's read model, which
// reports the target as resolved metadata rather than as the identifiers that were submitted.
func newReplicationTargetRef(target apigen.ReplicationTarget) *ReplicationTarget {
	return &ReplicationTarget{
		Cluster: types.StringValue(target.Cluster.Id),
		Bucket:  types.StringValue(target.Bucket.Id),
		Type:    types.StringValue(string(target.Type)),
	}
}

func newReplicationSourceDetails(source apigen.ReplicationSource) *ReplicationEndpointDetails {
	return &ReplicationEndpointDetails{
		Project: &ReplicationProject{
			Id:   types.StringValue(source.Project.Id),
			Name: types.StringValue(source.Project.Name),
		},
		Cluster: &ReplicationCluster{
			Id:   types.StringValue(source.Cluster.Id),
			Name: types.StringValue(source.Cluster.Name),
		},
		Bucket: &ReplicationBucket{
			Id:                     types.StringValue(source.Bucket.Id),
			Name:                   types.StringValue(source.Bucket.Name),
			ConflictResolutionType: types.StringValue(source.Bucket.ConflictResolutionType),
		},
		Scopes: newReplicationScopes(source.Scopes),
		Type:   types.StringPointerValue((*string)(source.Type)),
	}
}

func newReplicationTargetDetails(target apigen.ReplicationTarget) *ReplicationEndpointDetails {
	details := &ReplicationEndpointDetails{
		Cluster: &ReplicationCluster{
			Id:   types.StringValue(target.Cluster.Id),
			Name: types.StringValue(target.Cluster.Name),
		},
		Bucket: &ReplicationBucket{
			Id:                     types.StringValue(target.Bucket.Id),
			Name:                   types.StringValue(target.Bucket.Name),
			ConflictResolutionType: types.StringPointerValue(target.Bucket.ConflictResolutionType),
		},
		Scopes: newReplicationScopes(target.Scopes),
		Type:   types.StringValue(string(target.Type)),
	}

	// External targets have no Capella project.
	if target.Project != nil {
		details.Project = &ReplicationProject{
			Id:   types.StringPointerValue(target.Project.Id),
			Name: types.StringPointerValue(target.Project.Name),
		}
	}

	return details
}

func newReplicationScopes(scopes *apigen.Scopes) []ReplicationScope {
	if scopes == nil {
		return nil
	}

	converted := make([]ReplicationScope, 0, len(*scopes))
	for _, scope := range *scopes {
		replicationScope := ReplicationScope{Name: types.StringPointerValue(scope.Name)}
		if scope.Collections != nil {
			replicationScope.Collections = StringsToBaseStrings(*scope.Collections)
		}
		converted = append(converted, replicationScope)
	}
	return converted
}

func newReplicationMappings(mappings *apigen.Mappings) []ReplicationMapping {
	if mappings == nil {
		return nil
	}

	converted := make([]ReplicationMapping, 0, len(*mappings))
	for _, mapping := range *mappings {
		replicationMapping := ReplicationMapping{
			SourceScope: types.StringValue(mapping.SourceScope),
			TargetScope: types.StringValue(mapping.TargetScope),
		}
		if mapping.Collections != nil {
			replicationMapping.Collections = make([]ReplicationCollection, 0, len(*mapping.Collections))
			for _, collection := range *mapping.Collections {
				replicationMapping.Collections = append(replicationMapping.Collections, ReplicationCollection{
					SourceCollection: types.StringValue(collection.SourceCollection),
					TargetCollection: types.StringValue(collection.TargetCollection),
				})
			}
		}
		converted = append(converted, replicationMapping)
	}
	return converted
}

func newReplicationFilter(filter *apigen.GetFilter) *ReplicationFilter {
	if filter == nil {
		return nil
	}

	converted := &ReplicationFilter{}
	if filter.DocumentExcludeOptions != nil {
		converted.DocumentExcludeOptions = &ReplicationDocumentExcludeOptions{
			Deletion:   types.BoolPointerValue(filter.DocumentExcludeOptions.Deletion),
			Expiration: types.BoolPointerValue(filter.DocumentExcludeOptions.Expiration),
			Ttl:        types.BoolPointerValue(filter.DocumentExcludeOptions.Ttl),
			Binary:     types.BoolPointerValue(filter.DocumentExcludeOptions.Binary),
		}
	}
	if filter.Expressions != nil {
		converted.Expressions = &ReplicationFilterExpressions{
			RegEx: types.StringPointerValue(filter.Expressions.RegEx),
		}
	}
	return converted
}

func newReplicationAuditData(audit apigen.ReplicationAuditData) ReplicationAuditData {
	return ReplicationAuditData{
		CreatedAt: types.StringValue(audit.CreatedAt.String()),
		CreatedBy: types.StringValue(audit.CreatedBy),
	}
}

// toAPIMappings assembles the generated client's anonymous mapping structs field by field.
// Go type identity forbids converting named structs into them, so the nesting is spelled out here
// rather than in each caller.
func toAPIMappings(mappings []ReplicationMapping) *apigen.Mappings {
	if len(mappings) == 0 {
		return nil
	}

	converted := make(apigen.Mappings, 0, len(mappings))
	for _, mapping := range mappings {
		var apiMapping struct {
			Collections *[]struct {
				SourceCollection string `json:"sourceCollection"`
				TargetCollection string `json:"targetCollection"`
			} `json:"collections,omitempty"`
			SourceScope string `json:"sourceScope"`
			TargetScope string `json:"targetScope"`
		}
		apiMapping.SourceScope = mapping.SourceScope.ValueString()
		apiMapping.TargetScope = mapping.TargetScope.ValueString()

		if len(mapping.Collections) > 0 {
			collections := make([]struct {
				SourceCollection string `json:"sourceCollection"`
				TargetCollection string `json:"targetCollection"`
			}, 0, len(mapping.Collections))
			for _, collection := range mapping.Collections {
				collections = append(collections, struct {
					SourceCollection string `json:"sourceCollection"`
					TargetCollection string `json:"targetCollection"`
				}{
					SourceCollection: collection.SourceCollection.ValueString(),
					TargetCollection: collection.TargetCollection.ValueString(),
				})
			}
			apiMapping.Collections = &collections
		}

		converted = append(converted, apiMapping)
	}
	return &converted
}

// toAPIFilter assembles the generated client's anonymous filter structs field by field, for the
// same reason as toAPIMappings.
func toAPIFilter(filter *ReplicationFilter) *apigen.Filter {
	if filter == nil {
		return nil
	}

	converted := &apigen.Filter{}

	if options := filter.DocumentExcludeOptions; options != nil {
		converted.DocumentExcludeOptions = &struct {
			Binary     *bool `json:"binary,omitempty"`
			Deletion   *bool `json:"deletion,omitempty"`
			Expiration *bool `json:"expiration,omitempty"`
			Ttl        *bool `json:"ttl,omitempty"`
		}{
			Binary:     options.Binary.ValueBoolPointer(),
			Deletion:   options.Deletion.ValueBoolPointer(),
			Expiration: options.Expiration.ValueBoolPointer(),
			Ttl:        options.Ttl.ValueBoolPointer(),
		}
	}

	if expressions := filter.Expressions; expressions != nil {
		converted.Expressions = &struct {
			RegEx        *string `json:"regEx,omitempty"`
			SkipRestream *bool   `json:"skipRestream,omitempty"`
		}{
			RegEx:        expressions.RegEx.ValueStringPointer(),
			SkipRestream: expressions.SkipRestream.ValueBoolPointer(),
		}
	}

	return converted
}

// newOptionalInt64 widens an optional API int, mapping absent to a Terraform null rather than to 0.
func newOptionalInt64(value *int) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*value))
}

// int64ToIntPointer narrows a Terraform Int64 to the API's int, mapping null and unknown to nil so
// that an unset attribute is omitted from the payload rather than sent as 0.
func int64ToIntPointer(value types.Int64) *int {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	narrowed := int(value.ValueInt64())
	return &narrowed
}
