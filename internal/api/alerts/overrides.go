package alerts

// AlertConfigurationOverride is one alert's override at one level, as returned by the list
// alert configuration overrides endpoint.
type AlertConfigurationOverride struct {
	Scope              string                   `json:"scope"`
	OrganizationId     string                   `json:"organizationId"`
	ProjectId          *string                  `json:"projectId,omitempty"`
	ResourceId         *string                  `json:"resourceId,omitempty"`
	AlertKind          string                   `json:"alertKind"`
	Config             *AlertConfigurationEntry `json:"config,omitempty"`
	PreviousConfig     *AlertConfigurationEntry `json:"previousConfig,omitempty"`
	MutedAlert         *AlertMuteState          `json:"mutedAlert,omitempty"`
	Recipients         []string                 `json:"recipients,omitempty"`
	PreviousRecipients []string                 `json:"previousRecipients,omitempty"`
}
