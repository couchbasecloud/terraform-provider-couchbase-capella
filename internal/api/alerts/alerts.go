package alerts

// Alert is one alert on one resource, firing or resolved, as returned by the list alerts endpoint.
type Alert struct {
	AlertId      string  `json:"alertId"`
	EventId      string  `json:"eventId"`
	Name         string  `json:"name"`
	Kind         string  `json:"kind"`
	ResourceId   string  `json:"resourceId"`
	ResourceType string  `json:"resourceType"`
	ResourceName *string `json:"resourceName,omitempty"`
	ProjectId    *string `json:"projectId,omitempty"`
	ProjectName  *string `json:"projectName,omitempty"`
	Bucket       *string `json:"bucket,omitempty"`
	NodeHostname *string `json:"nodeHostname,omitempty"`
	Severity     string  `json:"severity"`
	Status       string  `json:"status"`
	StartsAt     string  `json:"startsAt"`
	EndsAt       *string `json:"endsAt,omitempty"`
	EventTime    string  `json:"eventTime"`
}
