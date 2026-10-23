package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = (*ClusterExpressScaling)(nil)
	_ resource.ResourceWithConfigure   = (*ClusterExpressScaling)(nil)
	_ resource.ResourceWithImportState = (*ClusterExpressScaling)(nil)
)

// expressScalingInvalidStateCode is the cp-open-api error code for ErrExpressScalingInvalidState,
// returned when the cluster is not healthy, which includes while it is rebalancing.
const expressScalingInvalidStateCode = 4050

const expressScalingDefaultTimeout = 60 * time.Minute

// Variables so tests can poll quickly.
var expressScalingPollInterval = time.Minute

var errExpressScalingTimeout = errors.New("timed out waiting for the express scaling change to finish")

// ClusterExpressScaling manages the express scaling setting of an existing cluster.
type ClusterExpressScaling struct {
	*providerschema.Data
}

func NewClusterExpressScaling() resource.Resource {
	return &ClusterExpressScaling{}
}

func (r *ClusterExpressScaling) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_express_scaling"
}

func (r *ClusterExpressScaling) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ClusterExpressScalingSchema()
}

func (r *ClusterExpressScaling) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.Data = data
}

func (r *ClusterExpressScaling) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ClusterExpressScaling
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, expressScalingDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshed, diags := r.apply(ctx, plan, timeout)
	resp.Diagnostics.Append(diags...)
	if refreshed != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, refreshed)...)
	}
}

func (r *ClusterExpressScaling) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.ClusterExpressScaling
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cluster express scaling",
			"Could not validate state: "+err.Error(),
		)
		return
	}

	current, err := r.getExpressScalingState(ctx, IDs[providerschema.OrganizationId], IDs[providerschema.ProjectId], IDs[providerschema.ClusterId])
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cluster express scaling",
			"Could not read cluster: "+errString,
		)
		return
	}

	state.OrganizationId = types.StringValue(IDs[providerschema.OrganizationId])
	state.ProjectId = types.StringValue(IDs[providerschema.ProjectId])
	state.ClusterId = types.StringValue(IDs[providerschema.ClusterId])
	state.ExpressScaling = types.BoolValue(current.IsEnabled())
	state.ExpressScalingState = types.StringValue(string(current))

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ClusterExpressScaling) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.ClusterExpressScaling
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, expressScalingDefaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshed, diags := r.apply(ctx, plan, timeout)
	resp.Diagnostics.Append(diags...)
	if refreshed != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, refreshed)...)
	}
}

// Delete only removes the resource from state. Disabling express scaling starts a long
// rebalance with cost implications, so destroying the resource leaves the cluster as it is.
func (r *ClusterExpressScaling) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *ClusterExpressScaling) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	req.ID = normalizeDeletionProtectionImportID(req.ID)
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// apply requests the planned setting and waits for the cluster to report it. It returns the
// state to save, which is nil only when nothing was changed in Capella.
//
// The PUT is skipped when the cluster already reports the planned setting, because each call
// queues a cluster deployment. If the cluster is already moving towards the planned setting,
// the API refuses the PUT because the cluster is rebalancing, and the provider waits instead.
// A change that failed part-way can still be retried, because the cluster is healthy again.
//
// When the timeout is reached the change carries on in Capella, so the planned value is saved
// with a warning rather than an error. Reads report a change in progress as the setting it is
// moving to, so this does not cause drift.
func (r *ClusterExpressScaling) apply(ctx context.Context, plan providerschema.ClusterExpressScaling, timeout time.Duration) (*providerschema.ClusterExpressScaling, diag.Diagnostics) {
	var diags diag.Diagnostics
	organizationId, projectId, clusterId := plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString()
	enabled := plan.ExpressScaling.ValueBool()

	current, err := r.getExpressScalingState(ctx, organizationId, projectId, clusterId)
	if err != nil {
		diags.AddError(
			"Error reading cluster express scaling",
			"Could not read cluster before changing express scaling: "+api.ParseError(err),
		)
		return nil, diags
	}

	if !current.IsSettled(enabled) {
		err := r.putExpressScaling(ctx, organizationId, projectId, clusterId, enabled)
		switch {
		case err == nil:
		case isExpressScalingInvalidState(err) && current.IsMovingTo(enabled):
			tflog.Info(ctx, "express scaling change already in progress, waiting for it to finish", map[string]interface{}{
				"express_scaling_state": current,
			})
		default:
			diags.AddError("Error changing cluster express scaling", expressScalingErrorDetail(err))
			return nil, diags
		}
	}

	final, err := r.waitForExpressScaling(ctx, organizationId, projectId, clusterId, enabled, timeout)
	if err != nil && !errors.Is(err, errExpressScalingTimeout) {
		diags.AddError(
			"Error waiting for cluster express scaling",
			"The express scaling change was accepted, but the cluster could not be read afterwards: "+api.ParseError(err),
		)
		return nil, diags
	}
	if err != nil {
		diags.AddWarning(
			"Express scaling change still in progress",
			fmt.Sprintf("The cluster still reports express scaling as %q after %s. The change carries on in Capella "+
				"while the cluster rebalances, which can take many hours on large clusters. The planned value has been saved; "+
				"increase the timeouts to wait longer.", final, timeout),
		)
	}

	return &providerschema.ClusterExpressScaling{
		OrganizationId:      plan.OrganizationId,
		ProjectId:           plan.ProjectId,
		ClusterId:           plan.ClusterId,
		ExpressScaling:      plan.ExpressScaling,
		ExpressScalingState: types.StringValue(string(final)),
		Timeouts:            plan.Timeouts,
	}, diags
}

// waitForExpressScaling polls the cluster until express scaling settles at the requested
// setting. Read errors are retried, except for not found. On timeout it returns the last
// state seen together with errExpressScalingTimeout.
func (r *ClusterExpressScaling) waitForExpressScaling(
	ctx context.Context,
	organizationId, projectId, clusterId string,
	enabled bool,
	timeout time.Duration,
) (clusterapi.ExpressScalingState, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(expressScalingPollInterval)
	defer ticker.Stop()

	var last clusterapi.ExpressScalingState
	for {
		current, err := r.getExpressScalingState(ctx, organizationId, projectId, clusterId)
		switch {
		case err != nil:
			if resourceNotFound, _ := api.CheckResourceNotFoundError(err); resourceNotFound {
				return last, err
			}
			tflog.Info(ctx, "retrying after error polling cluster express scaling", map[string]interface{}{"error": err.Error()})
		case current.IsSettled(enabled):
			return current, nil
		default:
			last = current
			tflog.Info(ctx, "waiting for cluster express scaling change to finish", map[string]interface{}{
				"express_scaling_state": current,
			})
		}

		select {
		case <-ctx.Done():
			return last, fmt.Errorf("%w: %w", errExpressScalingTimeout, ctx.Err())
		case <-ticker.C:
		}
	}
}

func (r *ClusterExpressScaling) putExpressScaling(ctx context.Context, organizationId, projectId, clusterId string, enabled bool) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/expressScaling", r.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := r.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		clusterapi.UpdateExpressScalingRequest{ExpressScaling: enabled},
		r.Token,
		nil,
	)
	return err
}

func (r *ClusterExpressScaling) getExpressScalingState(ctx context.Context, organizationId, projectId, clusterId string) (clusterapi.ExpressScalingState, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", r.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := r.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return "", err
	}

	var clusterResp clusterapi.GetClusterResponse
	if err := json.Unmarshal(response.Body, &clusterResp); err != nil {
		return "", err
	}

	return clusterResp.ExpressScaling, nil
}

func isExpressScalingInvalidState(err error) bool {
	var apiError *api.Error
	return errors.As(err, &apiError) && apiError.Code == expressScalingInvalidStateCode
}

// expressScalingErrorDetail explains the express scaling errors the API can return. The
// remaining 422 responses (Enterprise plan required, express scaling not available for the
// organization, rate limit reached) carry a message that already says what to do.
func expressScalingErrorDetail(err error) string {
	var apiError *api.Error
	if !errors.As(err, &apiError) {
		return "Could not change express scaling: " + err.Error()
	}

	var lead string
	switch {
	case apiError.Code == expressScalingInvalidStateCode:
		lead = "The cluster must be healthy to change express scaling. Wait for any scaling, upgrade or rebalance on the cluster to finish, then apply again."
	case apiError.HttpStatusCode == http.StatusNotFound:
		lead = "The cluster was not found, or express scaling is not available for this organization."
	case apiError.HttpStatusCode == http.StatusBadRequest:
		lead = "Capella rejected the express scaling request. Please report this issue to the provider developers."
	default:
		lead = "Could not change express scaling."
	}

	return lead + " " + apiError.CompleteError()
}
