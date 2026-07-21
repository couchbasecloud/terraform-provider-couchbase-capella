package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/data_api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = (*DataApi)(nil)
	_ resource.ResourceWithConfigure   = (*DataApi)(nil)
	_ resource.ResourceWithImportState = (*DataApi)(nil)
)

const errorMessageWhileDataApiUpdate = "There is an error during the Data API configuration update. Please check in Capella to see if any" +
	" hanging resources have been created, unexpected error: "

const errorMessageAfterDataApiUpdate = "Data API configuration update has been processed, but encountered an error while checking the current" +
	" state of the Data API. Please run `terraform plan` after 1-2 minutes to know the" +
	" current state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

// DataApi is the Data API configuration resource implementation.
type DataApi struct {
	*providerschema.Data
}

func NewDataApi() resource.Resource {
	return &DataApi{}
}

// Metadata returns the Data API configuration resource type name.
func (d *DataApi) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_api"
}

// Schema defines the schema for the Data API configuration resource.
func (d *DataApi) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DataApiSchema()
}

// Configure adds the provider configured client to the Data API configuration resource.
func (d *DataApi) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.Data = data
}

// Create updates a clusters data API configuration and adds it to the state.
func (d *DataApi) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.DataApi
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	err := d.updateDataApi(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Data API configuration",
			errorMessageWhileDataApiUpdate+api.ParseError(err),
		)
		return
	}

	refreshedState, err := d.checkDataApiStatus(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Data API configuration",
			errorMessageAfterDataApiUpdate+api.ParseError(err),
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

// Read retrieves the Data API and network peering status.
func (d *DataApi) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.DataApi
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Data API Configuration",
			"Could not read Capella Data API configuration: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
	)

	statusResp, err := d.getDataApiStatus(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "cluster not found, removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella Data API Configuration",
			"Could not read Capella Data API configuration: "+errString,
		)
		return
	}

	refreshedState := providerschema.NewDataApi(organizationId, projectId, clusterId, statusResp)

	// Set refreshed state
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the Data API configuration.
func (d *DataApi) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.DataApi
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	err := d.updateDataApi(ctx, plan)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "cluster not found, removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error updating Data API configuration",
			errorMessageWhileDataApiUpdate+errString,
		)
		return
	}

	refreshedState, err := d.checkDataApiStatus(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Data API configuration",
			errorMessageAfterDataApiUpdate+api.ParseError(err),
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

// Delete is a no-op operation because the Data API configuration cannot be deleted. The Data API and network peering are left in their current state on the cluster.
func (d *DataApi) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// No-op operation.
}

func (d *DataApi) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// updateDataApi calls the update endpoint, which asynchronously enables or disables the Data API and network peering on the cluster.
func (d *DataApi) updateDataApi(ctx context.Context, plan providerschema.DataApi) error {
	updateRequest := data_api.UpdateDataApiRequest{
		EnableDataApi:        plan.EnableDataApi.ValueBool(),
		EnableNetworkPeering: plan.EnableNetworkPeering.ValueBool(),
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI",
		d.HostURL,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusAccepted}
	_, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		updateRequest,
		d.Token,
		nil,
	)

	return err
}

// checkDataApiStatus polls the Data API status until both the Data API and its network peering reach a final state,
// then returns the refreshed resource state.
// Polls up to 30 minutes, with a 15-second interval between polls to allow for the asynchronous operation to complete.
func (d *DataApi) checkDataApiStatus(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.DataApi, error) {
	const timeout = time.Minute * 30

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		statusResp, err := d.getDataApiStatus(ctx, organizationId, projectId, clusterId)
		switch {
		case err != nil:
			tflog.Info(ctx, "retrying after error polling Data API status", map[string]interface{}{"error": err.Error()})
		case data_api.IsFinalState(statusResp.State) && data_api.IsFinalState(statusResp.StateForNetworkPeering):
			tflog.Debug(ctx, "data api has reached a final state", map[string]interface{}{
				"state_for_data_api":        statusResp.State,
				"state_for_network_peering": statusResp.StateForNetworkPeering,
			})
			return providerschema.NewDataApi(organizationId, projectId, clusterId, statusResp), nil
		default:
			tflog.Info(ctx, "waiting for Data API to finish transitioning to a final state", map[string]interface{}{
				"state_for_data_api":        statusResp.State,
				"state_for_network_peering": statusResp.StateForNetworkPeering,
			})
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out while waiting for data API state to finish transitioning to a final state: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

func (d *DataApi) getDataApiStatus(ctx context.Context, organizationId, projectId, clusterId string) (*data_api.GetDataApiStatusResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI",
		d.HostURL,
		organizationId,
		projectId,
		clusterId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := d.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		d.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
	}

	statusResp := data_api.GetDataApiStatusResponse{}
	err = json.Unmarshal(response.Body, &statusResp)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}

	return &statusResp, nil
}
