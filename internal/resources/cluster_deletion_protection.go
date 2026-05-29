package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = &ClusterDeletionProtection{}
	_ resource.ResourceWithConfigure   = &ClusterDeletionProtection{}
	_ resource.ResourceWithImportState = &ClusterDeletionProtection{}
)

// ClusterDeletionProtection manages the deletion protection state of a cluster.
type ClusterDeletionProtection struct {
	*providerschema.Data
}

func NewClusterDeletionProtection() resource.Resource {
	return &ClusterDeletionProtection{}
}

func (r *ClusterDeletionProtection) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_deletion_protection"
}

func (r *ClusterDeletionProtection) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ClusterDeletionProtectionSchema()
}

func (r *ClusterDeletionProtection) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ClusterDeletionProtection) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ClusterDeletionProtection
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.putDeletionProtection(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString(), plan.DeletionProtection.ValueBool()); err != nil {
		resp.Diagnostics.AddError(
			"Error setting cluster deletion protection",
			"Could not set deletion protection: "+api.ParseError(err),
		)
		return
	}

	refreshed, err := r.retrieveState(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cluster deletion protection",
			"Could not read cluster after setting deletion protection: "+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, refreshed)
	resp.Diagnostics.Append(diags...)
}

func (r *ClusterDeletionProtection) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.ClusterDeletionProtection
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cluster deletion protection",
			"Could not validate state: "+err.Error(),
		)
		return
	}

	refreshed, err := r.retrieveState(ctx, IDs[providerschema.OrganizationId], IDs[providerschema.ProjectId], IDs[providerschema.ClusterId])
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cluster deletion protection",
			"Could not read cluster: "+errString,
		)
		return
	}

	diags = resp.State.Set(ctx, refreshed)
	resp.Diagnostics.Append(diags...)
}

func (r *ClusterDeletionProtection) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.ClusterDeletionProtection
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.putDeletionProtection(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString(), plan.DeletionProtection.ValueBool()); err != nil {
		resp.Diagnostics.AddError(
			"Error updating cluster deletion protection",
			"Could not update deletion protection: "+api.ParseError(err),
		)
		return
	}

	refreshed, err := r.retrieveState(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cluster deletion protection",
			"Could not read cluster after update: "+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, refreshed)
	resp.Diagnostics.Append(diags...)
}

// Delete is a no-op. Removing this resource from state does not alter the cluster.
// Capella's v4 does not support a DELETE for deletion protection — it is toggled via PUT.
// https://docs.couchbase.com/cloud/management-api-reference/index.html
func (r *ClusterDeletionProtection) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *ClusterDeletionProtection) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	req.ID = normalizeDeletionProtectionImportID(req.ID)
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

// normalizeDeletionProtectionImportID replaces an "id=<value>" key-value pair in the
// import string with "cluster_id=<value>", so both forms are accepted on import.
// Other keys that happen to contain "id" as a substring (e.g. organization_id) are unaffected.
func normalizeDeletionProtectionImportID(importID string) string {
	pairs := strings.Split(importID, ",")
	for i, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 && kv[0] == "id" {
			pairs[i] = "cluster_id=" + kv[1]
		}
	}
	return strings.Join(pairs, ",")
}

func (r *ClusterDeletionProtection) putDeletionProtection(ctx context.Context, organizationId, projectId, clusterId string, enabled bool) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/deletionProtection", r.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := r.ClientV1.ExecuteWithRetry(
		ctx,
		cfg,
		clusterapi.UpdateDeletionProtectionRequest{DeletionProtection: enabled},
		r.Token,
		nil,
	)
	return err
}

func (r *ClusterDeletionProtection) retrieveState(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.ClusterDeletionProtection, error) {
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
		return nil, err
	}

	var clusterResp clusterapi.GetClusterResponse
	if err := json.Unmarshal(response.Body, &clusterResp); err != nil {
		return nil, err
	}

	return &providerschema.ClusterDeletionProtection{
		OrganizationId:     types.StringValue(organizationId),
		ProjectId:          types.StringValue(projectId),
		ClusterId:          types.StringValue(clusterId),
		DeletionProtection: types.BoolValue(clusterResp.DeletionProtection),
	}, nil
}
