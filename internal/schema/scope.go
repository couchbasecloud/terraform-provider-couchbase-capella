package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

type Scope struct {
	// Collections is the array of Collections under a single scope
	Collections []Collection `tfsdk:"collections"`

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

// Collection defines model for Collection.
type Collection struct {
	// MaxTTL Max TTL of the collection.
	MaxTTL types.Int64 `tfsdk:"max_ttl"`

	// Name is the name of the collection.
	Name types.String `tfsdk:"name"`

	// Uid is the UID of the collection.
	Uid types.String `tfsdk:"uid"`
}

// Scopes defines attributes for the LIST scopes response received from V4 Capella Public API.
type Scopes struct {
	// Array of scopes. The server returns an array of scopes in the bucket under the single Uid.
	Scopes []OneScope `tfsdk:"scopes"`

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

type OneScope struct {
	Collections    []Collection `tfsdk:"collections"`
	Name           types.String `tfsdk:"name"`
	Uid            types.String `tfsdk:"uid"`
	BucketId       types.String `tfsdk:"bucket_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	OrganizationId types.String `tfsdk:"organization_id"`
}
