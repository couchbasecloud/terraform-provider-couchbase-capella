package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = &OrganizationAlertConfiguration{}
	_ resource.ResourceWithConfigure   = &OrganizationAlertConfiguration{}
	_ resource.ResourceWithImportState = &OrganizationAlertConfiguration{}
)

const organizationAlertConfigurationLevel = "organization"

// OrganizationAlertConfiguration is the organization alert configuration resource implementation.
type OrganizationAlertConfiguration struct {
	*providerschema.Data
}

// NewOrganizationAlertConfiguration is a helper function to simplify the provider implementation.
func NewOrganizationAlertConfiguration() resource.Resource {
	return &OrganizationAlertConfiguration{}
}

// Metadata returns the organization alert configuration resource type name.
func (r *OrganizationAlertConfiguration) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_alert_configuration"
}

// Schema defines the schema for the organization alert configuration resource.
func (r *OrganizationAlertConfiguration) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = OrganizationAlertConfigurationSchema()
}

// Configure adds the provider configured client to the organization alert configuration resource.
func (r *OrganizationAlertConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create sets the organization's alert configuration. There is no create endpoint, since an
// organization always has exactly one alert configuration, so this updates it.
func (r *OrganizationAlertConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.OrganizationAlertConfigurationResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	levelPath := plan.OrganizationId.ValueString()
	settings, err := applyAlertConfiguration(ctx, r.Data, levelPath, plan.AlertConfigurationSettings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Capella Organization Alert Configuration",
			"Could not set the organization alert configuration, unexpected error: "+api.ParseError(err),
		)
		return
	}
	plan.AlertConfigurationSettings = settings

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the organization's alert configuration.
func (r *OrganizationAlertConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.OrganizationAlertConfigurationResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Organization Alert Configuration",
			"Could not validate the organization alert configuration state: "+err.Error(),
		)
		return
	}

	levelPath := IDs[providerschema.OrganizationId]
	settings, err := readAlertConfiguration(ctx, r.Data, levelPath, state.AlertConfigurationSettings)
	if err != nil {
		notFound, errString := api.CheckResourceNotFoundError(err)
		if notFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella Organization Alert Configuration",
			"Could not read the organization alert configuration, unexpected error: "+errString,
		)
		return
	}

	state.OrganizationId = types.StringValue(IDs[providerschema.OrganizationId])
	state.AlertConfigurationSettings = settings

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update sends only what differs between the plan and the organization's current alert configuration.
func (r *OrganizationAlertConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state providerschema.OrganizationAlertConfigurationResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	levelPath := plan.OrganizationId.ValueString()
	settings, err := applyAlertConfiguration(ctx, r.Data, levelPath, plan.AlertConfigurationSettings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Capella Organization Alert Configuration",
			"Could not update the organization alert configuration, unexpected error: "+api.ParseError(err),
		)
		return
	}

	if warning := alertConfigurationRemovalWarning(organizationAlertConfigurationLevel, state.AlertConfigurationSettings, plan.AlertConfigurationSettings); warning != "" {
		resp.Diagnostics.AddWarning("Alert Configuration Left In Capella", warning)
	}
	plan.AlertConfigurationSettings = settings

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete clears the mutes and recipients this resource manages. There is no delete endpoint, and
// thresholds and frequencies cannot be removed, so they stay with the values last set.
func (r *OrganizationAlertConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.OrganizationAlertConfigurationResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	levelPath := state.OrganizationId.ValueString()
	if err := deleteAlertConfiguration(ctx, r.Data, levelPath, state.AlertConfigurationSettings); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Capella Organization Alert Configuration",
			"Could not clear the organization's mutes and recipients, unexpected error: "+api.ParseError(err),
		)
		return
	}

	if warning := alertConfigurationLeftBehindWarning(organizationAlertConfigurationLevel, state.AlertConfigurationSettings); warning != "" {
		resp.Diagnostics.AddWarning("Alert Configuration Left In Capella", warning)
	}
}

// ImportState imports the organization's alert configuration from "organization_id=<organization_id>",
// the key=value form the other levels use. Only the ID is imported, as at the other levels. The
// shared import parser cannot be used: it only parses an import string stored in another
// attribute while organization_id is still null.
func (r *OrganizationAlertConfiguration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organizationId, err := parseOrganizationAlertConfigurationImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Capella Organization Alert Configuration",
			"Could not import the organization alert configuration: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), organizationId)...)
}

func parseOrganizationAlertConfigurationImportID(importID string) (string, error) {
	organizationId, ok := strings.CutPrefix(importID, "organization_id=")
	if !ok || organizationId == "" || strings.Contains(organizationId, ",") {
		return "", fmt.Errorf("%w: expected organization_id=<organization_id>, got %q", errors.ErrInvalidImport, importID)
	}
	return organizationId, nil
}
