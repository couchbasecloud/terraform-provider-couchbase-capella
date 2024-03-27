package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	scheduleapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster_onoff_schedule"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ClusterOnOffSchedule{}
	_ datasource.DataSourceWithConfigure = &ClusterOnOffSchedule{}
)

// ClusterOnOffSchedule is the OnOffSchedule data source implementation.
type ClusterOnOffSchedule struct {
	*providerschema.Data
}

// NewClusterOnOffSchedule is a helper function to simplify the provider implementation.
func NewClusterOnOffSchedule() datasource.DataSource {
	return &ClusterOnOffSchedule{}
}

// Metadata returns the certificates data source type name.
func (c *ClusterOnOffSchedule) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_onoff_schedule"
}

// Schema defines the schema for the allowlist data source.
func (c *ClusterOnOffSchedule) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": requiredStringAttribute,
			"project_id":      requiredStringAttribute,
			"cluster_id":      requiredStringAttribute,
			"timezone":        computedStringAttribute,
			"days": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"state": computedStringAttribute,
						"day":   computedStringAttribute,
						"from": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"hour":   computedInt64Attribute,
								"minute": computedInt64Attribute,
							},
						},
						"to": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"hour":   computedInt64Attribute,
								"minute": computedInt64Attribute,
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of projects.
func (c *ClusterOnOffSchedule) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state providerschema.ClusterOnOffSchedule
	diags := req.Config.Get(ctx, &state)
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

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/onOffSchedule", c.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)

	if err != nil {
		// Adding the condition to check if the error is 404-cluster on/off schedule not found, if yes, then skip throwing the error.
		// This is because it always throws 404-schedule not found as initially no schedule exists.
		apiError := err.(*api.Error)
		if apiError.Code == 11040 {
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella Cluster On/off schedule",
			"Could not read On/off schedule in cluster "+state.ClusterId.String()+": "+api.ParseError(err),
		)
		return
	}

	onOffScheduleResp := scheduleapi.GetClusterOnOffScheduleResponse{}
	err = json.Unmarshal(response.Body, &onOffScheduleResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Cluster On/off schedule",
			"Could not read Cluster On/off schedule, unexpected error: "+err.Error(),
		)
		return
	}

	dayItems := make([]providerschema.DayItem, len(onOffScheduleResp.Days))
	for i, d := range onOffScheduleResp.Days {
		day := providerschema.DayItem{}
		if d.State != "" {
			day.State = types.StringValue(d.State)
		}
		if d.Day != "" {
			day.Day = types.StringValue(d.Day)
		}

		if d.From != nil {
			day.From = &providerschema.OnTimeBoundary{
				Hour:   types.Int64Value(d.From.Hour),
				Minute: types.Int64Value(d.From.Minute),
			}
		}

		if d.To != nil {
			day.To = &providerschema.OnTimeBoundary{
				Hour:   types.Int64Value(d.To.Hour),
				Minute: types.Int64Value(d.To.Minute),
			}
		}

		dayItems[i] = day
	}

	state.Timezone = types.StringValue(onOffScheduleResp.Timezone)
	state.Days = dayItems

	// Set state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the project data source.
func (c *ClusterOnOffSchedule) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	c.Data = data
}
