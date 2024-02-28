package resources

import (
	"context"
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ClusterOnOffSchedule{}
	_ resource.ResourceWithConfigure   = &ClusterOnOffSchedule{}
	_ resource.ResourceWithImportState = &ClusterOnOffSchedule{}
)

// ClusterOnOffSchedule is the OnOffSchedule resource implementation.
type ClusterOnOffSchedule struct {
	*providerschema.Data
}

// NewClusterOnOffSchedule is a helper function to simplify the provider implementation.
func NewClusterOnOffSchedule() resource.Resource {
	return &ClusterOnOffSchedule{}
}

// Metadata returns the OnOffSchedule resource type name.
func (c *ClusterOnOffSchedule) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_onoff_schedule"

}

// Schema defines the schema for the OnOffSchedule resource.
func (c *ClusterOnOffSchedule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = OnOffScheduleSchema()
}

// Create creates a new OnOffSchedule.
func (c *ClusterOnOffSchedule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Read reads OnOffSchedule information.
func (c *ClusterOnOffSchedule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the OnOffSchedule.
func (c *ClusterOnOffSchedule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the OnOffSchedule.
func (c *ClusterOnOffSchedule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// ImportState imports an already existing cluster on-off schedule that is not created by Terraform.
// Since Capella APIs may require multiple IDs, such as organizationId, projectId, clusterId,
// this function passes the root attribute which is a comma separated string of multiple IDs.
// example: "organization_id=<orgId>,project_id=<projId>,cluster_id=<clusterId>
func (c *ClusterOnOffSchedule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

func (c *ClusterOnOffSchedule) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	c.Data = data
}

func (c *ClusterOnOffSchedule) validateCreateOnOffScheduleRequest(plan providerschema.ClusterOnOffSchedule) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdCannotBeEmpty
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdCannotBeEmpty
	}

	return nil
}
