package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Replication{}
	_ datasource.DataSourceWithConfigure = &Replication{}
)

// Replication is the replication data source implementation.
type Replication struct {
	*providerschema.Data
}

// NewReplication is a helper function to simplify the provider implementation.
func NewReplication() datasource.DataSource {
	return &Replication{}
}

// Metadata returns the replication data source type name.
func (r *Replication) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replication"
}

// Schema defines the schema for the replication data source.
func (r *Replication) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ReplicationSchema()
}

// Read reads the replication information.
func (r *Replication) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ReplicationData
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizationId := state.OrganizationId.ValueString()
	projectId := state.ProjectId.ValueString()
	clusterId := state.ClusterId.ValueString()
	replicationId := state.ReplicationId.ValueString()

	// Parse string IDs to UUIDs for the API client
	orgUUID, projUUID, clusterUUID, err := r.parseUUIDs(organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing IDs",
			"Could not parse resource IDs: "+err.Error(),
		)
		return
	}

	response, err := r.ClientV2.GetReplicationWithResponse(
		ctx,
		orgUUID,
		projUUID,
		clusterUUID,
		replicationId,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Replication",
			fmt.Sprintf("Could not read replication %s: %s: %s", replicationId, errors.ErrExecutingRequest, err.Error()),
		)
		return
	}

	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error Reading Replication",
			fmt.Sprintf("Unexpected response while reading replication: %s", string(response.Body)),
		)
		return
	}

	if response.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Error Reading Replication",
			"API returned an empty response body.",
		)
		return
	}

	tflog.Info(ctx, "read replication", map[string]interface{}{
		"organization_id": organizationId,
		"project_id":      projectId,
		"cluster_id":      clusterId,
		"replication_id":  replicationId,
	})

	// Map the API response to the datasource state.
	state = *providerschema.NewReplicationData(
		organizationId,
		projectId,
		clusterId,
		replicationId,
		response.JSON200,
	)

	// Set nested fields
	if err := state.SetNestedFields(ctx, response.JSON200); err != nil {
		resp.Diagnostics.AddError(
			"Error Building Replication State",
			fmt.Sprintf("Could not build replication state: %s", err.Error()),
		)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Configure adds the provider configured client to the data source.
func (r *Replication) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.Data = data
}

// parseUUIDs parses the string IDs into UUID types for the generated API client.
func (r *Replication) parseUUIDs(organizationId, projectId, clusterId string) (uuid.UUID, uuid.UUID, uuid.UUID, error) {
	orgUUID, err := uuid.Parse(organizationId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, fmt.Errorf("invalid organization_id: %w", err)
	}

	projUUID, err := uuid.Parse(projectId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, fmt.Errorf("invalid project_id: %w", err)
	}

	clusterUUID, err := uuid.Parse(clusterId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, fmt.Errorf("invalid cluster_id: %w", err)
	}

	return orgUUID, projUUID, clusterUUID, nil
}
