package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EventSchema returns the schema for the Event data source.
func EventSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID for this instance of an Event.",
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of then Capella organization.",
			},
			"alert_key":        schema.StringAttribute{Computed: true, MarkdownDescription: "Populated on demand based on the Event.Key and select labels in KV."},
			"app_service_id":   schema.StringAttribute{Computed: true, MarkdownDescription: "SyncGatewayID this Event refers to."},
			"app_service_name": schema.StringAttribute{Computed: true, MarkdownDescription: "Name of the sync gateway at the time of event emission."},
			"cluster_id":       schema.StringAttribute{Computed: true, MarkdownDescription: "ClusterID this Event refers to."},
			"cluster_name":     schema.StringAttribute{Computed: true, MarkdownDescription: "Name of the cluster at the time of event emission."},
			"image_url":        schema.StringAttribute{Computed: true, MarkdownDescription: "Rendered chart image to display for an Alert Event."},
			"incident_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Group events related to an alert incident.",
			},
			"key":              schema.StringAttribute{Computed: true, MarkdownDescription: "Defines the specific kind of Event."},
			"kv":               schema.StringAttribute{Computed: true, MarkdownDescription: "Key-value pairs for additional event data."},
			"occurrence_count": schema.Int64Attribute{Computed: true, MarkdownDescription: "Number of times the alert has fired within this \"incident\"."},
			"project_id":       schema.StringAttribute{Computed: true, MarkdownDescription: "ProjectID this Event refers to."},
			"project_name":     schema.StringAttribute{Computed: true, MarkdownDescription: "Name of the project at the time of event emission."},
			"request_id":       schema.StringAttribute{Computed: true, MarkdownDescription: "RequestID for an Event."},
			"session_id":       schema.StringAttribute{Computed: true, MarkdownDescription: "User that initiated the request for this Event."},
			"severity":         schema.StringAttribute{Computed: true, MarkdownDescription: "Severity of the event."},
			"source":           schema.StringAttribute{Computed: true, MarkdownDescription: "Identifies the originator of the event."},
			"summary":          schema.StringAttribute{Computed: true, MarkdownDescription: "Metadata.SummaryTemplate rendered for this event."},
			"timestamp":        schema.StringAttribute{Computed: true, MarkdownDescription: "Time when the event was emitted."},
			"user_email":       schema.StringAttribute{Computed: true, MarkdownDescription: "Email of the associated user at the time of event emission."},
			"user_id":          schema.StringAttribute{Computed: true, MarkdownDescription: "User id that initiated the request for this Event."},
			"user_name":        schema.StringAttribute{Computed: true, MarkdownDescription: "Name of the associated user at the time of event emission."},
		},
	}
}
