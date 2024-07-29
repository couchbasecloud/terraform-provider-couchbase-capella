package api

import (
	"time"

	"github.com/google/uuid"
)

// GetEventResponse defines model for GetEventResponse.
type GetEventResponse struct {
	// AlertKey Populated on demand based on the Event.Key and select labels in KV.
	AlertKey *string `json:"alertKey,omitempty"`

	// AppServiceId SyncGatewayID this Event refers to.
	AppServiceId *uuid.UUID `json:"appServiceId,omitempty"`

	// AppServiceName Name of the sync gateway at the time of event emission.
	AppServiceName *string `json:"appServiceName,omitempty"`

	// ClusterId ClusterID this Event refers to.
	ClusterId *uuid.UUID `json:"clusterId,omitempty"`

	// ClusterName Name of the cluster at the time of event emission.
	ClusterName *string `json:"clusterName,omitempty"`

	// Id UUID for this instance of an Event.
	Id uuid.UUID `json:"id"`

	// ImageURL Rendered chart image to display for an Alert Event.
	ImageURL *string `json:"imageURL,omitempty"`

	// IncidentIds Group events related to an alert incident.
	IncidentIds *[]uuid.UUID `json:"incidentIds,omitempty"`

	// Key Defines the specific kind of Event.
	Key string `json:"key"`

	// Kv Key-value pairs for additional event data.
	Kv *map[string]interface{} `json:"kv,omitempty"`

	// OccurrenceCount Number of times the alert has fired within this "incident".
	OccurrenceCount *int `json:"occurrenceCount,omitempty"`

	// ProjectId ProjectID this Event refers to.
	ProjectId *uuid.UUID `json:"projectId,omitempty"`

	// ProjectName Name of the project at the time of event emission.
	ProjectName *string `json:"projectName,omitempty"`

	// RequestId RequestID for an Event.
	RequestId *uuid.UUID `json:"requestId,omitempty"`

	// SessionId User that initiated the request for this Event.
	SessionId *uuid.UUID `json:"sessionId,omitempty"`

	// Severity Severity of the event.
	Severity string `json:"severity"`

	// Source Identifies the originator of the event.
	Source string `json:"source"`

	// Summary Metadata.SummaryTemplate rendered for this event.
	Summary *string `json:"summary,omitempty"`

	// Timestamp Time when the event was emitted.
	Timestamp time.Time `json:"timestamp"`

	// UserEmail Email of the associated user at the time of event emission.
	UserEmail *string `json:"userEmail,omitempty"`

	// UserId User id that initiated the request for this Event.
	UserId *string `json:"userId,omitempty"`

	// UserName Name of the associated user at the time of event emission.
	UserName *string `json:"userName,omitempty"`
}

// GetEventsResponse defines model for GetEventsResponse.
type GetEventsResponse struct {
	Cursor Cursor             `json:"cursor"`
	Data   []GetEventResponse `json:"data"`
}
