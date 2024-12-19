package schema

import (
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
