package schema

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Events represents a structure for filtering and paginating events.
type Events struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectIds     types.Set    `tfsdk:"project_ids"`
	ClusterIds     types.Set    `tfsdk:"cluster_ids"`
	UserIds        types.Set    `tfsdk:"user_ids"`
	SeverityLevels types.Set    `tfsdk:"severity_levels"`
	Tags           types.Set    `tfsdk:"tags"`
	From           types.String `tfsdk:"from"`
	To             types.String `tfsdk:"to"`
	Page           types.Int64  `tfsdk:"page"`
	PerPage        types.Int64  `tfsdk:"per_page"`
	SortBy         types.String `tfsdk:"sort_by"`
	SortDirection  types.String `tfsdk:"sort_direction"`
	Data           []EventItem  `tfsdk:"data"`
	Cursor         *Cursor      `tfsdk:"cursor"`
}

// EventItem represents a single event item with detailed information.
type EventItem struct {
	AlertKey        types.String `tfsdk:"alert_key"`
	AppServiceId    types.String `tfsdk:"app_service_id"`
	AppServiceName  types.String `tfsdk:"app_service_name"`
	ClusterId       types.String `tfsdk:"cluster_id"`
	ClusterName     types.String `tfsdk:"cluster_name"`
	Id              types.String `tfsdk:"id"`
	ImageUrl        types.String `tfsdk:"image_url"`
	IncidentIds     types.Set    `tfsdk:"incident_ids"`
	Key             types.String `tfsdk:"key"`
	Kv              types.String `tfsdk:"kv"`
	OccurrenceCount types.Int64  `tfsdk:"occurrence_count"`
	ProjectId       types.String `tfsdk:"project_id"`
	ProjectName     types.String `tfsdk:"project_name"`
	RequestId       types.String `tfsdk:"request_id"`
	SessionId       types.String `tfsdk:"session_id"`
	Severity        types.String `tfsdk:"severity"`
	Source          types.String `tfsdk:"source"`
	Summary         types.String `tfsdk:"summary"`
	Timestamp       types.String `tfsdk:"timestamp"`
	UserEmail       types.String `tfsdk:"user_email"`
	UserId          types.String `tfsdk:"user_id"`
	UserName        types.String `tfsdk:"user_name"`
}

// NewEventItem creates a new EventItem instance from the provided API response.
func NewEventItem(event *api.GetEventResponse, incidentIdsSet types.Set, kvString types.String) (*EventItem, error) {
	newEventItem := EventItem{
		Id:          types.StringValue(event.Id.String()),
		Key:         types.StringValue(event.Key),
		Severity:    types.StringValue(event.Severity),
		Source:      types.StringValue(event.Source),
		IncidentIds: incidentIdsSet,
		Kv:          kvString,
		Timestamp:   types.StringValue(event.Timestamp.String()),
	}
	if event.AlertKey != nil {
		newEventItem.AlertKey = types.StringValue(*event.AlertKey)
	}
	if event.AppServiceId != nil {
		newEventItem.AppServiceId = types.StringValue(event.AppServiceId.String())
	}
	if event.AppServiceName != nil {
		newEventItem.AppServiceName = types.StringValue(*event.AppServiceName)
	}
	if event.ClusterId != nil {
		newEventItem.ClusterId = types.StringValue(event.ClusterId.String())
	}
	if event.ClusterName != nil {
		newEventItem.ClusterName = types.StringValue(*event.ClusterName)
	}
	if event.ImageURL != nil {
		newEventItem.ImageUrl = types.StringValue(*event.ImageURL)
	}
	if event.OccurrenceCount != nil {
		newEventItem.OccurrenceCount = types.Int64Value(int64(*event.OccurrenceCount))
	}
	if event.ProjectId != nil {
		newEventItem.ProjectId = types.StringValue(event.ProjectId.String())
	}
	if event.ProjectName != nil {
		newEventItem.ProjectName = types.StringValue(*event.ProjectName)
	}
	if event.RequestId != nil {
		newEventItem.RequestId = types.StringValue(event.RequestId.String())
	}
	if event.SessionId != nil {
		newEventItem.SessionId = types.StringValue(event.SessionId.String())
	}
	if event.Summary != nil {
		newEventItem.Summary = types.StringValue(*event.Summary)
	}
	if event.UserName != nil {
		newEventItem.UserName = types.StringValue(*event.UserName)
	}
	if event.UserId != nil {
		newEventItem.UserId = types.StringValue(event.UserId.String())
	}

	return &newEventItem, nil
}

// Event represents a detailed event with all associated information.
type Event struct {
	Id              types.String `tfsdk:"id"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	AlertKey        types.String `tfsdk:"alert_key"`
	AppServiceId    types.String `tfsdk:"app_service_id"`
	AppServiceName  types.String `tfsdk:"app_service_name"`
	ClusterId       types.String `tfsdk:"cluster_id"`
	ClusterName     types.String `tfsdk:"cluster_name"`
	ImageUrl        types.String `tfsdk:"image_url"`
	IncidentIds     types.Set    `tfsdk:"incident_ids"`
	Key             types.String `tfsdk:"key"`
	Kv              types.String `tfsdk:"kv"`
	OccurrenceCount types.Int64  `tfsdk:"occurrence_count"`
	ProjectId       types.String `tfsdk:"project_id"`
	ProjectName     types.String `tfsdk:"project_name"`
	RequestId       types.String `tfsdk:"request_id"`
	SessionId       types.String `tfsdk:"session_id"`
	Severity        types.String `tfsdk:"severity"`
	Source          types.String `tfsdk:"source"`
	Summary         types.String `tfsdk:"summary"`
	Timestamp       types.String `tfsdk:"timestamp"`
	UserEmail       types.String `tfsdk:"user_email"`
	UserId          types.String `tfsdk:"user_id"`
	UserName        types.String `tfsdk:"user_name"`
}

// NewEvent creates a new Event instance from the provided API response.
func NewEvent(event *api.GetEventResponse, organizationId types.String, incidentIdsSet types.Set, kvString types.String) (*Event, error) {
	newEvent := Event{
		OrganizationId: organizationId,
		Id:             types.StringValue(event.Id.String()),
		Key:            types.StringValue(event.Key),
		Severity:       types.StringValue(event.Severity),
		Source:         types.StringValue(event.Source),
		IncidentIds:    incidentIdsSet,
		Kv:             kvString,
		Timestamp:      types.StringValue(event.Timestamp.String()),
	}
	if event.AlertKey != nil {
		newEvent.AlertKey = types.StringValue(*event.AlertKey)
	}
	if event.AppServiceId != nil {
		newEvent.AppServiceId = types.StringValue(event.AppServiceId.String())
	}
	if event.AppServiceName != nil {
		newEvent.AppServiceName = types.StringValue(*event.AppServiceName)
	}
	if event.ClusterId != nil {
		newEvent.ClusterId = types.StringValue(event.ClusterId.String())
	}
	if event.ClusterName != nil {
		newEvent.ClusterName = types.StringValue(*event.ClusterName)
	}
	if event.ImageURL != nil {
		newEvent.ImageUrl = types.StringValue(*event.ImageURL)
	}
	if event.OccurrenceCount != nil {
		newEvent.OccurrenceCount = types.Int64Value(int64(*event.OccurrenceCount))
	}
	if event.ProjectId != nil {
		newEvent.ProjectId = types.StringValue(event.ProjectId.String())
	}
	if event.ProjectName != nil {
		newEvent.ProjectName = types.StringValue(*event.ProjectName)
	}
	if event.RequestId != nil {
		newEvent.RequestId = types.StringValue(event.RequestId.String())
	}
	if event.SessionId != nil {
		newEvent.SessionId = types.StringValue(event.SessionId.String())
	}
	if event.Summary != nil {
		newEvent.Summary = types.StringValue(*event.Summary)
	}
	if event.UserName != nil {
		newEvent.UserName = types.StringValue(*event.UserName)
	}
	if event.UserId != nil {
		newEvent.UserId = types.StringValue(event.UserId.String())
	}

	return &newEvent, nil
}
