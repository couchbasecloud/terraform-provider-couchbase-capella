package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProjectEventSchema returns the schema for the ProjectEvent data source.
func ProjectEventSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Data source to retrieve an event for a project. Events represent a trail of actions that users performs within Capella at project level.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The event ID of the event.",
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"alert_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Populated on demand based on the Event.Key and select labels in KV.",
			},
			"app_service_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SyncGatewayID this Event refers to.",
			},
			"app_service_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the sync gateway at the time of event emission.",
			},
			"cluster_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ClusterID this Event refers to.",
			},
			"cluster_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the cluster at the time of event emission.",
			},
			"image_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Rendered chart image to display for an Alert Event.",
			},
			"incident_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Group events related to an alert incident.",
			},
			"key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Defines the specific kind of Event.",
			},
			"kv": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Key-value pairs for additional event data.",
			},
			"occurrence_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of times the alert has fired within this \"incident\".",
			},
			"project_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the project at the time of event emission.",
			},
			"request_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RequestID for an Event.",
			},
			"session_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "User that initiated the request for this Event.",
			},
			"severity": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Severity of the event.",
			},
			"source": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifies the originator of the event.",
			},
			"summary": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Metadata.SummaryTemplate rendered for this event.",
			},
			"timestamp": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Time when the event was emitted.",
			},
			"user_email": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Email of the associated user at the time of event emission.",
			},
			"user_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "User id that initiated the request for this Event.",
			},
			"user_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the associated user at the time of event emission.",
			},
		},
	}
}
