package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/scope"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// Scope maps Scope resource schema data to the response received from V4 Capella Public API.
type Scope struct {
	// Collections is the array of Collections under a single scope
	Collections types.Set `tfsdk:"collections"`

	// Name is the name of the scope.
	Name types.String `tfsdk:"name"`

	// Uid is the UID of the scope.
	Uid types.String `tfsdk:"uid"`

	// BucketId is the id of the bucket for which the scope needs to be created.
	BucketId types.String `tfsdk:"bucket_id"`

	// ClusterId is the ID of the cluster for which the scope needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`
}

// Collection defines a Collection within the list of collections in a scope.
type Collection struct {
	// MaxTTL Max TTL of the collection.
	MaxTTL types.Int64 `tfsdk:"max_ttl"`

	// Name is the name of the collection.
	Name types.String `tfsdk:"name"`

	// Uid is the UID of the collection.
	Uid types.String `tfsdk:"uid"`
}

func CollectionAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"max_ttl": types.Int64Type,
		"name":    types.StringType,
		"uid":     types.StringType,
	}
}

// NewCollection creates a new collection object.
func NewCollection(collection scope.Collection) Collection {
	return Collection{
		//TODO: check nil too
		MaxTTL: types.Int64Value(*collection.MaxTTL),
		Name:   types.StringValue(*collection.Name),
		Uid:    types.StringValue(*collection.Uid),
	}
}

// Scopes defines structure based on the response received from V4 Capella Public API when asked to list scopes.
type Scopes struct {
	// Array of scopes. The server returns an array of scopes in the bucket under the single Uid.
	Scopes []ScopeData `tfsdk:"scopes"`

	// Uid is the UID of the whole scope containing all scopes.
	Uid types.String `tfsdk:"uid"`

	// BucketId is the id of the bucket for which the scope needs to be created.
	BucketId types.String `tfsdk:"bucket_id"`

	// ClusterId is the ID of the cluster for which the scope needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`
}

// ScopeData defines attributes for a single Scope when fetched from the V4 Capella Public API.
type ScopeData struct {
	Collections types.Set    `tfsdk:"collections"`
	Name        types.String `tfsdk:"name"`
	Uid         types.String `tfsdk:"uid"`
}

// NewScopeData creates new scope object.
func NewScopeData(scope *scope.GetScopeResponse,
	collectionSet basetypes.SetValue,
) *ScopeData {
	newScopeData := ScopeData{
		Collections: collectionSet,
		Name:        types.StringValue(*scope.Name),
		Uid:         types.StringValue(*scope.Uid),
	}
	return &newScopeData
}

// Validate will split the IDs by a delimiter i.e. comma , in case a terraform import CLI is invoked.
func (s Scope) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: s.OrganizationId,
		ProjectId:      s.ProjectId,
		ClusterId:      s.ClusterId,
		BucketId:       s.BucketId,
		ScopeName:      s.Name,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

// Validate is used to verify that all the fields in the datasource have been populated.
func (s Scopes) Validate() (bucketId, clusterId, projectId, organizationId string, err error) {
	if s.BucketId.IsNull() {
		return "", "", "", "", errors.ErrBucketIdMissing
	}
	if s.OrganizationId.IsNull() {
		return "", "", "", "", errors.ErrOrganizationIdMissing
	}
	if s.ProjectId.IsNull() {
		return "", "", "", "", errors.ErrProjectIdMissing
	}
	if s.ClusterId.IsNull() {
		return "", "", "", "", errors.ErrClusterIdMissing
	}

	return s.BucketId.ValueString(), s.ClusterId.ValueString(), s.ProjectId.ValueString(), s.OrganizationId.ValueString(), nil
}
