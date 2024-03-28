package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	collection "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// Collection maps Collection resource schema data to the response received from V4 Capella Public API.
type Collection struct {

	// Name is the name of the collection.
	Name types.String `tfsdk:"collection_name"`

	// MaxTTL Max TTL of the collection.
	MaxTTL types.Int64 `tfsdk:"max_ttl"`

	// ScopeName is the name of the scope for which the collection needs to be created.
	ScopeName types.String `tfsdk:"scope_name"`

	// BucketId is the id of the bucket for which the collection needs to be created.
	BucketId types.String `tfsdk:"bucket_id"`

	// ClusterId is the ID of the cluster for which the collection needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`
}

// Collections defines structure based on the response received from V4 Capella Public API when asked to list collections.
type Collections struct {
	// Array of scopes. The server returns an array of scopes in the bucket under the single Uid.
	Data []CollectionData `tfsdk:"data"`

	// ScopeName is the name of the scope for which the collection needs to be created.
	ScopeName types.String `tfsdk:"scope_name"`

	// BucketId is the id of the bucket for which the collection needs to be created.
	BucketId types.String `tfsdk:"bucket_id"`

	// ClusterId is the ID of the cluster for which the collection needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`
}

// CollectionData defines attributes for a single Collection when fetched from the V4 Capella Public API.
type CollectionData struct {
	Name   types.String `tfsdk:"collection_name"`
	MaxTTL types.Int64  `tfsdk:"max_ttl"`
}

// Validate will split the IDs by a delimiter i.e. comma , in case a terraform import CLI is invoked.
func (c *Collection) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: c.OrganizationId,
		ProjectId:      c.ProjectId,
		ClusterId:      c.ClusterId,
		BucketId:       c.BucketId,
		ScopeName:      c.ScopeName,
		CollectionName: c.Name,
	}

	IDs, err := validateSchemaState(state, CollectionName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

// Validate is used to verify that all the fields in the datasource have been populated.
func (c Collections) Validate() (bucketId, clusterId, projectId, organizationId, scopeName string, err error) {
	if c.BucketId.IsNull() {
		return "", "", "", "", "", errors.ErrBucketIdMissing
	}
	if c.OrganizationId.IsNull() {
		return "", "", "", "", "", errors.ErrOrganizationIdMissing
	}
	if c.ProjectId.IsNull() {
		return "", "", "", "", "", errors.ErrProjectIdMissing
	}
	if c.ClusterId.IsNull() {
		return "", "", "", "", "", errors.ErrClusterIdMissing
	}
	if c.ScopeName.IsNull() {
		return "", "", "", "", "", errors.ErrScopeNameMissing
	}

	return c.BucketId.ValueString(), c.ClusterId.ValueString(), c.ProjectId.ValueString(), c.OrganizationId.ValueString(), c.ScopeName.ValueString(), nil
}

// NewCollectionData creates a new collectionData object.
func NewCollectionData(collection *collection.GetCollectionResponse) (*CollectionData, error) {
	newCollectionData := CollectionData{}
	if collection.MaxTTL != nil {
		newCollectionData.MaxTTL = types.Int64Value(*collection.MaxTTL)
	}

	if collection.Name != nil {
		newCollectionData.Name = types.StringValue(*collection.Name)
	}

	return &newCollectionData, nil
}
