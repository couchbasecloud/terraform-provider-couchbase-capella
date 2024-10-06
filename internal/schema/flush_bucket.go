package schema

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FlushBucket struct {
	// BucketId is the id of the bucket for which the scope needs to be created.
	BucketId types.String `tfsdk:"bucket_id"`

	// ClusterId is the ID of the cluster for which the scope needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`
}

// Validate is used to verify that all the fields in the datasource have been populated.
func (b FlushBucket) Validate() (bucketId, clusterId, projectId, organizationId string, err error) {
	if b.BucketId.IsNull() {
		return "", "", "", "", errors.ErrBucketIdMissing
	}
	if b.OrganizationId.IsNull() {
		return "", "", "", "", errors.ErrOrganizationIdMissing
	}
	if b.ProjectId.IsNull() {
		return "", "", "", "", errors.ErrProjectIdMissing
	}
	if b.ClusterId.IsNull() {
		return "", "", "", "", errors.ErrClusterIdMissing
	}

	return b.BucketId.ValueString(), b.ClusterId.ValueString(), b.ProjectId.ValueString(), b.OrganizationId.ValueString(), nil
}
