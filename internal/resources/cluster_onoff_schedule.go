package resources

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &ClusterOnOffSchedule{}
	_ resource.ResourceWithConfigure      = &ClusterOnOffSchedule{}
	_ resource.ResourceWithImportState    = &ClusterOnOffSchedule{}
	_ resource.ResourceWithValidateConfig = &ClusterOnOffSchedule{}
)

// weekdays is the order the V4 API requires for the on/off schedule days list:
// exactly one entry per day of the week, starting from Monday and ending with Sunday.
// The explicit length guards at compile time against entries being added or removed.
var weekdays = [7]string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

const errorMessageAfterOnOffScheduleCreation = "Cluster On/Off Schedule creation is successful, but encountered an error while checking the current" +
	" state of the cluster on/off schedule. Please run `terraform plan` after 1-2 minutes to know the" +
	" current on/off schedule state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileOnOffScheduleCreation = "There is an error during cluster on/off schedule creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

// ClusterOnOffSchedule is the cluster OnOffSchedule resource implementation.
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

// ValidateConfig enforces the constraints on the days list that the V4 API
// documents but the attribute-level validators cannot express: the schedule
// must contain exactly one entry per day of the week in Monday-to-Sunday order,
// the cluster cannot be scheduled to be off for every day of the week, custom
// days require a from time boundary, non-custom days cannot have time
// boundaries, and from must not be later than to.
func (c *ClusterOnOffSchedule) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var days types.List
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("days"), &days)...)
	if resp.Diagnostics.HasError() || days.IsNull() || days.IsUnknown() {
		return
	}

	var dayItems []providerschema.DayItem
	// Entries may still be wholly unknown (e.g. computed from other resources),
	// which cannot be decoded into DayItem; allow them to decode as zero values
	// so the per-item null/unknown checks below defer validation to apply time,
	// while genuine decoding errors still surface.
	resp.Diagnostics.Append(days.ElementsAs(ctx, &dayItems, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(dayItems) != len(weekdays) {
		resp.Diagnostics.AddAttributeError(
			path.Root("days"),
			"Invalid Cluster On/Off Schedule",
			fmt.Sprintf("The schedule requires exactly %d days, one for each day of the week starting"+
				" from Monday and ending with Sunday, got %d.", len(weekdays), len(dayItems)),
		)
		return
	}

	allOff := true
	for i, d := range dayItems {
		if d.Day.IsNull() || d.Day.IsUnknown() || d.State.IsNull() || d.State.IsUnknown() {
			return
		}
		if d.Day.ValueString() != weekdays[i] {
			resp.Diagnostics.AddAttributeError(
				path.Root("days").AtListIndex(i).AtName("day"),
				"Invalid Cluster On/Off Schedule",
				fmt.Sprintf("The days of the week must be in sequence starting from Monday and ending"+
					" with Sunday: expected %q at position %d, got %q.", weekdays[i], i, d.Day.ValueString()),
			)
			return
		}
		if d.State.ValueString() != "off" {
			allOff = false
		}
		validateDayTimeBoundaries(i, d, resp)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if allOff {
		resp.Diagnostics.AddAttributeError(
			path.Root("days"),
			"Invalid Cluster On/Off Schedule",
			"Clusters cannot be scheduled to be off for the entire day for every day of the week."+
				" At least one day must have state \"on\" or \"custom\".",
		)
	}
}

// validateDayTimeBoundaries enforces the V4 API rules tying a day's state to
// its from/to time boundaries: custom days must contain the from boundary,
// on/off days cannot contain any time boundary, and from must not be later
// than to.
func validateDayTimeBoundaries(i int, d providerschema.DayItem, resp *resource.ValidateConfigResponse) {
	day := d.Day.ValueString()

	switch d.State.ValueString() {
	case "custom":
		if d.From == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("days").AtListIndex(i).AtName("from"),
				"Invalid Cluster On/Off Schedule",
				fmt.Sprintf("The from time boundary is required when state is \"custom\" (day %q).", day),
			)
			return
		}
		if d.To == nil {
			return
		}
		from, fromKnown := boundaryMinutes(d.From)
		to, toKnown := boundaryMinutes(d.To)
		if fromKnown && toKnown && from > to {
			resp.Diagnostics.AddAttributeError(
				path.Root("days").AtListIndex(i).AtName("from"),
				"Invalid Cluster On/Off Schedule",
				fmt.Sprintf("The from time boundary must not be later than the to time boundary (day %q).", day),
			)
		}
	case "on", "off":
		if d.From != nil || d.To != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("days").AtListIndex(i),
				"Invalid Cluster On/Off Schedule",
				fmt.Sprintf("Days with state \"on\" or \"off\" cannot contain from/to time boundaries (day %q).", day),
			)
		}
	}
}

// boundaryMinutes converts a time boundary to minutes since midnight. Null
// hour/minute values default to 0, matching the API default. Returns false
// when a value is unknown so validation is deferred to apply time.
func boundaryMinutes(b *providerschema.OnTimeBoundary) (int64, bool) {
	if b.Hour.IsUnknown() || b.Minute.IsUnknown() {
		return 0, false
	}
	var hour, minute int64
	if !b.Hour.IsNull() {
		hour = b.Hour.ValueInt64()
	}
	if !b.Minute.IsNull() {
		minute = b.Minute.ValueInt64()
	}
	return hour*60 + minute, true
}

// Create creates a new OnOffSchedule.
func (c *ClusterOnOffSchedule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ClusterOnOffSchedule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := c.validateCreateClusterOnOffScheduleRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create cluster on/off schedule request",
			"Could not create cluster on/off schedule "+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()

	var days = make([]api.DayItem, 0)
	for _, d := range plan.Days {
		var from, to *api.OnTimeBoundary
		if d.From != nil {
			from = &api.OnTimeBoundary{
				Hour:   d.From.Hour.ValueInt64(),
				Minute: d.From.Minute.ValueInt64(),
			}
		}
		if d.To != nil {
			to = &api.OnTimeBoundary{
				Hour:   d.To.Hour.ValueInt64(),
				Minute: d.To.Minute.ValueInt64(),
			}
		}
		days = append(days, api.DayItem{
			State: d.State.ValueString(),
			Day:   d.Day.ValueString(),
			From:  from,
			To:    to,
		})
	}

	scheduleRequest := api.CreateClusterOnOffScheduleRequest{
		Timezone: plan.Timezone.ValueString(),
		Days:     days,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/onOffSchedule", c.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusNoContent}
	err = c.scheduleRequestWithRetry(ctx, cfg, scheduleRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			errorMessageWhileOnOffScheduleCreation+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeScheduleWithPlan(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := c.retrieveClusterOnOffSchedule(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error Reading Capella Cluster On/Off Schedule",
			"Could not read Capella Cluster On/Off Schedule for the cluster: %s "+clusterId+"."+errorMessageAfterOnOffScheduleCreation+api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// initializeScheduleWithPlan initializes an instance of providerschema.ClusterOnOffSchedule
// with the specified plan. It marks all computed fields as null.
func initializeScheduleWithPlan(plan providerschema.ClusterOnOffSchedule) providerschema.ClusterOnOffSchedule {
	for _, d := range plan.Days {
		if d.From != nil {
			if d.From.Hour.IsUnknown() || d.From.Hour.IsNull() {
				d.From.Hour = types.Int64Null()
			}
			if d.From.Minute.IsUnknown() || d.From.Minute.IsNull() {
				d.From.Minute = types.Int64Null()
			}
		}
		if d.To != nil {
			if d.To.Hour.IsUnknown() || d.To.Hour.IsNull() {
				d.To.Hour = types.Int64Null()
			}
			if d.To.Minute.IsUnknown() || d.To.Minute.IsNull() {
				d.To.Minute = types.Int64Null()
			}
		}
	}

	return plan
}

func (c *ClusterOnOffSchedule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.ClusterOnOffSchedule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cluster on/off schedule",
			"Could not read on/off schedule for cluster with id "+state.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
	)

	// Get refreshed on/off schedule from Capella
	refreshedState, err := c.retrieveClusterOnOffSchedule(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cluster on/off schedule",
			"Could not read on/off schedule for cluster with id "+state.ClusterId.String()+": "+errString,
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *ClusterOnOffSchedule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.ClusterOnOffSchedule
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := plan.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating cluster on/off schedule",
			"Could not update on/off schedule for cluster with id"+plan.ClusterId.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
	)

	var days = make([]api.DayItem, 0)
	for _, d := range plan.Days {
		var from, to *api.OnTimeBoundary
		if d.From != nil {
			from = &api.OnTimeBoundary{
				Hour:   d.From.Hour.ValueInt64(),
				Minute: d.From.Minute.ValueInt64(),
			}
		}
		if d.To != nil {
			to = &api.OnTimeBoundary{
				Hour:   d.To.Hour.ValueInt64(),
				Minute: d.To.Minute.ValueInt64(),
			}
		}
		days = append(days, api.DayItem{
			State: d.State.ValueString(),
			Day:   d.Day.ValueString(),
			From:  from,
			To:    to,
		})
	}

	BackupScheduleRequest := api.UpdateClusterOnOffScheduleRequest{
		Timezone: plan.Timezone.ValueString(),
		Days:     days,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/onOffSchedule", c.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	if err = c.scheduleRequestWithRetry(ctx, cfg, BackupScheduleRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error updating cluster on/off schedule",
			"Could not update on/off schedule for cluster with id "+plan.ClusterId.String()+": "+api.ParseError(err),
		)
		return
	}

	currentState, err := c.retrieveClusterOnOffSchedule(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cluster on/off schedule",
			"Could not read on/off schedule for cluster with id "+plan.ClusterId.String()+": "+errString,
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (c *ClusterOnOffSchedule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state providerschema.ClusterOnOffSchedule
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting cluster on/off schedule",
			"Could not delete cluster on/off schedule for cluster with id "+state.ClusterId.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
	)

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/onOffSchedule", c.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}

	// Delete existing backup schedule
	_, err = c.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting cluster on/off schedule",
			"Could not delete cluster on/off schedule for cluster with id "+state.ClusterId.String()+" unexpected error: "+errString,
		)
		return
	}
}

// ImportState imports an already existing cluster on-off schedule that is not created by Terraform.
// Since Capella APIs may require multiple IDs, such as organizationId, projectId, clusterId,
// this function passes the root attribute which is a comma separated string of multiple IDs.
// example: "organization_id=<orgId>,project_id=<projId>,cluster_id=<clusterId>
func (c *ClusterOnOffSchedule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// Configure adds the provider configured client to the cluster on-off schedule resource.
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

func (c *ClusterOnOffSchedule) validateCreateClusterOnOffScheduleRequest(plan providerschema.ClusterOnOffSchedule) error {
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

// retrieveOnOffSchedule retrieves cluster on/off schedule information from the specified organization and project
// using the provided cluster ID by open-api call.
func (c *ClusterOnOffSchedule) retrieveClusterOnOffSchedule(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.ClusterOnOffSchedule, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/onOffSchedule", c.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := c.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	onOffScheduleResp := api.GetClusterOnOffScheduleResponse{}
	err = json.Unmarshal(response.Body, &onOffScheduleResp)
	if err != nil {
		return nil, err
	}

	refreshedState := providerschema.NewClusterOnOffSchedule(&onOffScheduleResp, organizationId, projectId, clusterId)
	return refreshedState, nil
}

func (c *ClusterOnOffSchedule) scheduleRequestWithRetry(ctx context.Context, cfg api.EndpointCfg, scheduleRequest any) error {
	const (
		maxRetryWindow = 5 * time.Minute
		retryInterval  = 15 * time.Second
	)
	deadline := time.Now().Add(maxRetryWindow)
	var lastErr error
	for {
		_, err := c.ClientV1.ExecuteWithRetry(ctx, cfg, scheduleRequest, c.Token, nil)
		if err == nil {
			return nil
		}
		if ctx.Err() != nil {
			cause := lastErr
			if cause == nil {
				cause = err
			}
			return fmt.Errorf("retry window (%v) exhausted for schedule create: %w", maxRetryWindow, cause)
		}
		// Only retry the specific Capella transient backend error (code 10000:
		// "Something went wrong on our end. We are actively investigating the issue.").
		// This 500 appears for ~minutes after a previous schedule write while the
		// backend propagates state. Other 500s likely indicate real failures and
		// should surface immediately rather than burn the retry window.
		var apiErr *api.Error
		if !stderrors.As(err, &apiErr) || apiErr.HttpStatusCode != http.StatusInternalServerError || apiErr.Code != 10000 {
			return err
		}
		lastErr = err
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return fmt.Errorf("retry window (%v) exhausted for schedule create: %w", maxRetryWindow, lastErr)
		}
		tflog.Debug(ctx, "schedule create returned 500; retrying", map[string]interface{}{"err": err})
		sleep := retryInterval
		if sleep > remaining {
			sleep = remaining
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry window (%v) exhausted for schedule create: %w", maxRetryWindow, lastErr)
		case <-time.After(sleep):
		}
	}
}
