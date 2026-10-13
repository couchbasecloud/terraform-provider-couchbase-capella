package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// ReplicationTargetTypeCapella and ReplicationTargetTypeExternal are the values of target.type.
const (
	ReplicationTargetTypeCapella  = "capella"
	ReplicationTargetTypeExternal = "external"
)

// ReplicationDirectionOneWay and ReplicationDirectionTwoWay are the values of direction.
const (
	ReplicationDirectionOneWay = "oneWay"
	ReplicationDirectionTwoWay = "twoWay"
)

// ReplicationPriorityLow, ReplicationPriorityMedium and ReplicationPriorityHigh are the values of priority.
const (
	ReplicationPriorityLow    = "low"
	ReplicationPriorityMedium = "medium"
	ReplicationPriorityHigh   = "high"
)

// Fallback descriptions for attributes that the published OpenAPI spec does not describe, or
// describes for a different object.
const (
	ReplicationSourceScopeDescription      = "\n - The name of the scope in the source bucket."
	ReplicationTargetScopeDescription      = "\n - The name of the scope in the target bucket."
	ReplicationCollectionsDescription      = "\n - The collections to replicate from the source scope. When omitted, all collections in the scope are replicated."
	ReplicationSourceCollectionDescription = "\n - The name of the collection in the source scope."
	ReplicationTargetCollectionDescription = "\n - The name of the collection in the target scope."

	ReplicationDocumentExcludeOptionsDescription = "\n - The kinds of document change that are not replicated."
	ReplicationExpressionsDescription            = "\n - The filter expression. Only documents that match it are replicated."

	ReplicationSourceDetailsDescription = "\n - Details of the source of the replication, as returned by the API."
	ReplicationTargetDetailsDescription = "\n - Details of the target of the replication, as returned by the API."

	ReplicationDetailsIdDescription                     = "\n - The ID."
	ReplicationDetailsNameDescription                   = "\n - The name."
	ReplicationDetailsConflictResolutionTypeDescription = "\n - The conflict resolution type of the bucket."
)

// Replication maps the schema for the XDCR replication resource.
type Replication struct {
	Id             types.String `tfsdk:"id"`
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`

	// SourceBucket is the source bucket ID.
	SourceBucket types.String       `tfsdk:"source_bucket"`
	Target       *ReplicationTarget `tfsdk:"target"`

	Direction         types.String `tfsdk:"direction"`
	Priority          types.String `tfsdk:"priority"`
	NetworkUsageLimit types.Int64  `tfsdk:"network_usage_limit"`

	// Mappings is a set of ReplicationScopeMapping. GET returns mappings in a random order.
	Mappings types.Set          `tfsdk:"mappings"`
	Filter   *ReplicationFilter `tfsdk:"filter"`

	// ReverseReplicationId is only returned by the create job, so it is kept from prior state.
	ReverseReplicationId types.String `tfsdk:"reverse_replication_id"`

	Status        types.String `tfsdk:"status"`
	ChangesLeft   types.Int64  `tfsdk:"changes_left"`
	Error         types.String `tfsdk:"error"`
	SourceDetails types.Object `tfsdk:"source_details"`
	TargetDetails types.Object `tfsdk:"target_details"`
	Audit         types.Object `tfsdk:"audit"`
}

// ReplicationTarget is the configured target of a replication.
type ReplicationTarget struct {
	Cluster types.String `tfsdk:"cluster"`
	Bucket  types.String `tfsdk:"bucket"`
	Type    types.String `tfsdk:"type"`
}

// ReplicationScopeMapping is one element of Replication.Mappings.
type ReplicationScopeMapping struct {
	SourceScope types.String `tfsdk:"source_scope"`
	TargetScope types.String `tfsdk:"target_scope"`

	// Collections is a set of ReplicationCollectionMapping.
	Collections types.Set `tfsdk:"collections"`
}

func (m ReplicationScopeMapping) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source_scope": types.StringType,
		"target_scope": types.StringType,
		"collections":  types.SetType{ElemType: types.ObjectType{AttrTypes: ReplicationCollectionMapping{}.AttributeTypes()}},
	}
}

// ReplicationCollectionMapping is one element of ReplicationScopeMapping.Collections.
type ReplicationCollectionMapping struct {
	SourceCollection types.String `tfsdk:"source_collection"`
	TargetCollection types.String `tfsdk:"target_collection"`
}

func (c ReplicationCollectionMapping) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source_collection": types.StringType,
		"target_collection": types.StringType,
	}
}

// ReplicationFilter is the configured document filter of a replication.
type ReplicationFilter struct {
	DocumentExcludeOptions *ReplicationDocumentExcludeOptions `tfsdk:"document_exclude_options"`
	Expressions            *ReplicationFilterExpressions      `tfsdk:"expressions"`
}

// ReplicationDocumentExcludeOptions is the configured filter.document_exclude_options.
type ReplicationDocumentExcludeOptions struct {
	Deletion   types.Bool `tfsdk:"deletion"`
	Expiration types.Bool `tfsdk:"expiration"`
	TTL        types.Bool `tfsdk:"ttl"`
	Binary     types.Bool `tfsdk:"binary"`
}

// ReplicationFilterExpressions is the configured filter.expressions.
type ReplicationFilterExpressions struct {
	RegEx        types.String `tfsdk:"reg_ex"`
	SkipRestream types.Bool   `tfsdk:"skip_restream"`
}

// ReplicationAudit is the audit data of a replication. The API returns only these two fields.
type ReplicationAudit struct {
	CreatedBy types.String `tfsdk:"created_by"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func (a ReplicationAudit) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_by": types.StringType,
		"created_at": types.StringType,
	}
}

func replicationReferenceAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

func replicationBucketReferenceAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                       types.StringType,
		"name":                     types.StringType,
		"conflict_resolution_type": types.StringType,
	}
}

// ReplicationSourceDetailsAttributeTypes is the object type of Replication.SourceDetails.
func ReplicationSourceDetailsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project": types.ObjectType{AttrTypes: replicationReferenceAttributeTypes()},
		"cluster": types.ObjectType{AttrTypes: replicationReferenceAttributeTypes()},
		"bucket":  types.ObjectType{AttrTypes: replicationBucketReferenceAttributeTypes()},
	}
}

// ReplicationTargetDetailsAttributeTypes is the object type of Replication.TargetDetails.
func ReplicationTargetDetailsAttributeTypes() map[string]attr.Type {
	attrTypes := ReplicationSourceDetailsAttributeTypes()
	attrTypes["type"] = types.StringType
	return attrTypes
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

// NewReplication builds a fresh replication state from a GET response.
//
// prior is the plan (after create or update) or the prior state (on refresh). It is nil-safe.
// Only the values the API does not return, or returns in a different shape from the config, are
// taken from prior. The prior state is never patched in place.
func NewReplication(
	ctx context.Context,
	resp *replicationapi.GetReplicationResponse,
	organizationId, projectId, clusterId string,
	prior *Replication,
) (*Replication, error) {
	if prior == nil {
		prior = &Replication{}
	}

	mappings, err := newMappingsSet(ctx, resp.Mappings)
	if err != nil {
		return nil, err
	}

	sourceDetails, err := newSourceDetails(resp.Source)
	if err != nil {
		return nil, err
	}

	targetDetails, err := newTargetDetails(resp.Target)
	if err != nil {
		return nil, err
	}

	audit, diags := types.ObjectValueFrom(ctx, ReplicationAudit{}.AttributeTypes(), ReplicationAudit{
		CreatedBy: types.StringValue(resp.Audit.CreatedBy),
		CreatedAt: types.StringValue(resp.Audit.CreatedAt.String()),
	})
	if diags.HasError() {
		return nil, fmt.Errorf("%w: %s", errors.ErrUnableToConvertAuditData, diags.Errors())
	}

	sourceBucket := types.StringValue(resp.Source.Bucket.Id)
	if resp.Source.Bucket.Id == "" {
		sourceBucket = prior.SourceBucket
	}

	reverseReplicationId := prior.ReverseReplicationId
	if reverseReplicationId.IsUnknown() {
		reverseReplicationId = types.StringNull()
	}

	replicationError := types.StringNull()
	if resp.Error != "" {
		replicationError = types.StringValue(resp.Error)
	}

	return &Replication{
		Id:                   types.StringValue(resp.Id),
		OrganizationId:       types.StringValue(organizationId),
		ProjectId:            types.StringValue(projectId),
		ClusterId:            types.StringValue(clusterId),
		SourceBucket:         sourceBucket,
		Target:               newReplicationTarget(resp.Target, prior.Target),
		Direction:            types.StringValue(resp.Direction),
		Priority:             types.StringValue(resp.Priority),
		NetworkUsageLimit:    types.Int64Value(resp.NetworkUsageLimit),
		Mappings:             mappings,
		Filter:               newReplicationFilter(resp.Filter, prior.Filter),
		ReverseReplicationId: reverseReplicationId,
		Status:               types.StringValue(string(resp.Status)),
		ChangesLeft:          types.Int64Value(resp.ChangesLeft),
		Error:                replicationError,
		SourceDetails:        sourceDetails,
		TargetDetails:        targetDetails,
		Audit:                audit,
	}, nil
}

// newReplicationTarget maps the GET target back to the configured shape: IDs for a Capella
// target, the reference name and bucket name for an external target.
func newReplicationTarget(target replicationapi.TargetDetails, prior *ReplicationTarget) *ReplicationTarget {
	switch target.Type {
	case ReplicationTargetTypeCapella:
		return &ReplicationTarget{
			Cluster: types.StringValue(target.Cluster.Id),
			Bucket:  types.StringValue(target.Bucket.Id),
			Type:    types.StringValue(ReplicationTargetTypeCapella),
		}
	case ReplicationTargetTypeExternal:
		return &ReplicationTarget{
			Cluster: types.StringValue(target.Cluster.Name),
			Bucket:  types.StringValue(target.Bucket.Name),
			Type:    types.StringValue(ReplicationTargetTypeExternal),
		}
	default:
		return prior
	}
}

// newReplicationFilter builds the filter state. GET drops each part of the filter that has only
// default values, drops the whole filter when both parts are gone, and never returns skipRestream.
func newReplicationFilter(filter *replicationapi.Filter, prior *ReplicationFilter) *ReplicationFilter {
	var (
		apiExcludeOptions   *replicationapi.DocumentExcludeOptions
		apiExpressions      *replicationapi.FilterExpressions
		priorExcludeOptions *ReplicationDocumentExcludeOptions
		priorExpressions    *ReplicationFilterExpressions
	)
	if filter != nil {
		apiExcludeOptions, apiExpressions = filter.DocumentExcludeOptions, filter.Expressions
	}
	if prior != nil {
		priorExcludeOptions, priorExpressions = prior.DocumentExcludeOptions, prior.Expressions
	}

	excludeOptions := newDocumentExcludeOptions(apiExcludeOptions, priorExcludeOptions)
	expressions := newFilterExpressions(apiExpressions, priorExpressions)

	if excludeOptions == nil && expressions == nil && prior == nil {
		return nil
	}

	return &ReplicationFilter{
		DocumentExcludeOptions: excludeOptions,
		Expressions:            expressions,
	}
}

func newDocumentExcludeOptions(
	options *replicationapi.DocumentExcludeOptions,
	prior *ReplicationDocumentExcludeOptions,
) *ReplicationDocumentExcludeOptions {
	if options != nil {
		return &ReplicationDocumentExcludeOptions{
			Deletion:   types.BoolValue(options.Deletion),
			Expiration: types.BoolValue(options.Expiration),
			TTL:        types.BoolValue(options.TTL),
			Binary:     types.BoolValue(options.Binary),
		}
	}
	if prior != nil && prior.allFalse() {
		return &ReplicationDocumentExcludeOptions{
			Deletion:   types.BoolValue(false),
			Expiration: types.BoolValue(false),
			TTL:        types.BoolValue(false),
			Binary:     types.BoolValue(false),
		}
	}
	return nil
}

func newFilterExpressions(
	expressions *replicationapi.FilterExpressions,
	prior *ReplicationFilterExpressions,
) *ReplicationFilterExpressions {
	if expressions == nil {
		return nil
	}

	skipRestream := types.BoolValue(false)
	if prior != nil && !prior.SkipRestream.IsNull() && !prior.SkipRestream.IsUnknown() {
		skipRestream = prior.SkipRestream
	}

	return &ReplicationFilterExpressions{
		RegEx:        types.StringValue(expressions.RegEx),
		SkipRestream: skipRestream,
	}
}

func (o *ReplicationDocumentExcludeOptions) allFalse() bool {
	return !o.Deletion.ValueBool() && !o.Expiration.ValueBool() && !o.TTL.ValueBool() && !o.Binary.ValueBool()
}

// newMappingsSet converts the GET mappings to a set. No mappings means all scopes, which is null.
func newMappingsSet(ctx context.Context, mappings []replicationapi.ScopeMapping) (types.Set, error) {
	scopeType := types.ObjectType{AttrTypes: ReplicationScopeMapping{}.AttributeTypes()}
	collectionType := types.ObjectType{AttrTypes: ReplicationCollectionMapping{}.AttributeTypes()}

	if len(mappings) == 0 {
		return types.SetNull(scopeType), nil
	}

	scopes := make([]ReplicationScopeMapping, 0, len(mappings))
	for _, mapping := range mappings {
		collections := types.SetNull(collectionType)
		if len(mapping.Collections) > 0 {
			elements := make([]ReplicationCollectionMapping, 0, len(mapping.Collections))
			for _, collection := range mapping.Collections {
				elements = append(elements, ReplicationCollectionMapping{
					SourceCollection: types.StringValue(collection.SourceCollection),
					TargetCollection: types.StringValue(collection.TargetCollection),
				})
			}
			var diags diag.Diagnostics
			collections, diags = types.SetValueFrom(ctx, collectionType, elements)
			if diags.HasError() {
				return types.SetNull(scopeType), fmt.Errorf("%w: %s", errors.ErrConvertingReplicationMappings, diags.Errors())
			}
		}
		scopes = append(scopes, ReplicationScopeMapping{
			SourceScope: types.StringValue(mapping.SourceScope),
			TargetScope: types.StringValue(mapping.TargetScope),
			Collections: collections,
		})
	}

	set, diags := types.SetValueFrom(ctx, scopeType, scopes)
	if diags.HasError() {
		return types.SetNull(scopeType), fmt.Errorf("%w: %s", errors.ErrConvertingReplicationMappings, diags.Errors())
	}
	return set, nil
}

func newReferenceObject(reference replicationapi.Reference) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(replicationReferenceAttributeTypes(), map[string]attr.Value{
		"id":   types.StringValue(reference.Id),
		"name": types.StringValue(reference.Name),
	})
}

func newBucketReferenceObject(bucket replicationapi.BucketReference) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(replicationBucketReferenceAttributeTypes(), map[string]attr.Value{
		"id":                       types.StringValue(bucket.Id),
		"name":                     types.StringValue(bucket.Name),
		"conflict_resolution_type": types.StringValue(bucket.ConflictResolutionType),
	})
}

// newDetailsAttributes builds the project, cluster and bucket objects shared by the source and target details.
func newDetailsAttributes(project, cluster replicationapi.Reference, bucket replicationapi.BucketReference) (map[string]attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	projectObj, d := newReferenceObject(project)
	diags.Append(d...)
	clusterObj, d := newReferenceObject(cluster)
	diags.Append(d...)
	bucketObj, d := newBucketReferenceObject(bucket)
	diags.Append(d...)

	return map[string]attr.Value{
		"project": projectObj,
		"cluster": clusterObj,
		"bucket":  bucketObj,
	}, diags
}

func newSourceDetails(source replicationapi.SourceDetails) (types.Object, error) {
	attrs, diags := newDetailsAttributes(source.Project, source.Cluster, source.Bucket)
	if diags.HasError() {
		return types.ObjectNull(ReplicationSourceDetailsAttributeTypes()), fmt.Errorf("%w: %s", errors.ErrConvertingReplicationDetails, diags.Errors())
	}

	obj, diags := types.ObjectValue(ReplicationSourceDetailsAttributeTypes(), attrs)
	if diags.HasError() {
		return types.ObjectNull(ReplicationSourceDetailsAttributeTypes()), fmt.Errorf("%w: %s", errors.ErrConvertingReplicationDetails, diags.Errors())
	}
	return obj, nil
}

func newTargetDetails(target replicationapi.TargetDetails) (types.Object, error) {
	attrs, diags := newDetailsAttributes(target.Project, target.Cluster, target.Bucket)
	if diags.HasError() {
		return types.ObjectNull(ReplicationTargetDetailsAttributeTypes()), fmt.Errorf("%w: %s", errors.ErrConvertingReplicationDetails, diags.Errors())
	}
	attrs["type"] = types.StringValue(target.Type)

	obj, diags := types.ObjectValue(ReplicationTargetDetailsAttributeTypes(), attrs)
	if diags.HasError() {
		return types.ObjectNull(ReplicationTargetDetailsAttributeTypes()), fmt.Errorf("%w: %s", errors.ErrConvertingReplicationDetails, diags.Errors())
	}
	return obj, nil
}
