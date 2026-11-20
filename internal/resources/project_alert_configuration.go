package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = &ProjectAlertConfiguration{}
	_ resource.ResourceWithConfigure   = &ProjectAlertConfiguration{}
	_ resource.ResourceWithImportState = &ProjectAlertConfiguration{}
)

const projectAlertConfigurationLevel = "project"

// ProjectAlertConfiguration is the project alert configuration resource implementation.
type ProjectAlertConfiguration struct {
	*providerschema.Data
}

// NewProjectAlertConfiguration is a helper function to simplify the provider implementation.
func NewProjectAlertConfiguration() resource.Resource {
	return &ProjectAlertConfiguration{}
}

// Metadata returns the project alert configuration resource type name.
func (r *ProjectAlertConfiguration) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_alert_configuration"
}

// Schema defines the schema for the project alert configuration resource.
func (r *ProjectAlertConfiguration) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ProjectAlertConfigurationSchema()
}

// Configure adds the provider configured client to the project alert configuration resource.
func (r *ProjectAlertConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.Data = data
}

// Create sets the project's alert configuration. There is no create endpoint, since a project
// always has exactly one alert configuration, so this updates it.
func (r *ProjectAlertConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ProjectAlertConfigurationResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	levelPath := projectAlertConfigurationPath(plan.OrganizationId.ValueString(), plan.ProjectId.ValueString())
	settings, err := applyAlertConfiguration(ctx, r.Data, levelPath, plan.AlertConfigurationSettings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Capella Project Alert Configuration",
			"Could not set the project alert configuration, unexpected error: "+api.ParseError(err),
		)
		return
	}
	plan.AlertConfigurationSettings = settings

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the project's alert configuration.
func (r *ProjectAlertConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.ProjectAlertConfigurationResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Project Alert Configuration",
			"Could not validate the project alert configuration state: "+err.Error(),
		)
		return
	}

	levelPath := projectAlertConfigurationPath(IDs[providerschema.OrganizationId], IDs[providerschema.ProjectId])
	settings, err := readAlertConfiguration(ctx, r.Data, levelPath, state.AlertConfigurationSettings)
	if err != nil {
		notFound, errString := api.CheckResourceNotFoundError(err)
		if notFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella Project Alert Configuration",
			"Could not read the project alert configuration, unexpected error: "+errString,
		)
		return
	}

	state.OrganizationId = types.StringValue(IDs[providerschema.OrganizationId])
	state.ProjectId = types.StringValue(IDs[providerschema.ProjectId])
	state.AlertConfigurationSettings = settings

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update sends only what differs between the plan and the project's current alert configuration.
func (r *ProjectAlertConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state providerschema.ProjectAlertConfigurationResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	levelPath := projectAlertConfigurationPath(plan.OrganizationId.ValueString(), plan.ProjectId.ValueString())
	settings, err := applyAlertConfiguration(ctx, r.Data, levelPath, plan.AlertConfigurationSettings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Capella Project Alert Configuration",
			"Could not update the project alert configuration, unexpected error: "+api.ParseError(err),
		)
		return
	}

	if warning := alertConfigurationRemovalWarning(projectAlertConfigurationLevel, state.AlertConfigurationSettings, plan.AlertConfigurationSettings); warning != "" {
		resp.Diagnostics.AddWarning("Alert Configuration Left In Capella", warning)
	}
	plan.AlertConfigurationSettings = settings

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete clears the mutes and recipients this resource manages. There is no delete endpoint, and
// thresholds and frequencies cannot be removed, so they stay with the values last set.
func (r *ProjectAlertConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.ProjectAlertConfigurationResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	levelPath := projectAlertConfigurationPath(state.OrganizationId.ValueString(), state.ProjectId.ValueString())
	if err := deleteAlertConfiguration(ctx, r.Data, levelPath, state.AlertConfigurationSettings); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Capella Project Alert Configuration",
			"Could not clear the project's mutes and recipients, unexpected error: "+api.ParseError(err),
		)
		return
	}

	if warning := alertConfigurationLeftBehindWarning(projectAlertConfigurationLevel, state.AlertConfigurationSettings); warning != "" {
		resp.Diagnostics.AddWarning("Alert Configuration Left In Capella", warning)
	}
}

// ImportState imports the project's alert configuration by its IDs. Only the IDs are imported:
// configs, recipients and mutes are not, for the same reason Read ignores alerts the config does
// not list. The first apply after import sends the values in the user's config.
func (r *ProjectAlertConfiguration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("project_id"), req, resp)
}

func projectAlertConfigurationPath(organizationId, projectId string) string {
	return fmt.Sprintf("%s/projects/%s", organizationId, projectId)
}
