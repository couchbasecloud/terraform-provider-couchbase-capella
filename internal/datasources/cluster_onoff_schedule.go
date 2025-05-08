package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ClusterOnOffSchedule{}
	_ datasource.DataSourceWithConfigure = &ClusterOnOffSchedule{}
)

const errorCodeForOnOffScheduleNotFound = 11040

// ClusterOnOffSchedule is the OnOffSchedule data source implementation.
type ClusterOnOffSchedule struct {
	*providerschema.Data
}

// NewClusterOnOffSchedule is a helper function to simplify the provider implementation.
func NewClusterOnOffSchedule() datasource.DataSource {
	return &ClusterOnOffSchedule{}
}

// Metadata returns the cluster on/off schedule data source type name.
func (c *ClusterOnOffSchedule) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_onoff_schedule"
}

// Schema defines the schema for the cluster on/off schedule data source.
func (c *ClusterOnOffSchedule) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The On/Off schedule data source allows you to retrieve the on/off schedule for a Capella cluster.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
			},
			"timezone": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The standard timezone for the cluster. Should be the TZ identifier. For example, 'ET'",
			},
			"days": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of days for the cluster on/off schedule. Each day should have a state, day, from, and to time.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The state of the cluster on/off schedule. Should be one of 'on' or 'off'.",
						},
						"day": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The day of the week for the cluster on/off schedule. Should be one of 'sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', or 'saturday'.",
						},
						"from": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "The time from which the cluster on/off schedule starts.",
							Attributes: map[string]schema.Attribute{
								"hour": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The hour of the day for the cluster on/off schedule. Should be between 0 and 23.",
								},
								"minute": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The minute of the hour for the cluster on/off schedule. Should be between 0 and 59.",
								},
							},
						},
						"to": schema.SingleNestedAttribute{
							Optional:            true,
							MarkdownDescription: "The time to which the cluster on/off schedule ends.",
							Attributes: map[string]schema.Attribute{
								"hour": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The hour of the day for the cluster on/off schedule. Should be between 0 and 23.",
								},
								"minute": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The minute of the hour for the cluster on/off schedule. Should be between 0 and 59.",
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data of cluster on/off schedules.
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
		apiError, ok := err.(*api.Error)
		if !ok {
			resp.Diagnostics.AddError(
				"Type assertion error",
				"Could not do type assertion for the error:"+apiError.Error(),
			)
		}
		if apiError.Code == errorCodeForOnOffScheduleNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella Cluster On/off schedule",
			"Could not read On/off schedule in cluster "+state.ClusterId.String()+": "+api.ParseError(err),
		)
		return
	}

	onOffScheduleResp := api.GetClusterOnOffScheduleResponse{}
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

// Configure adds the provider configured client to the cluster on/off schedule data source.
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
