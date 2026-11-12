package alerts

// AlertConfigurationEntry is the configuration for a single alert at one level of the hierarchy.
// A field that is not set at that level is omitted.
type AlertConfigurationEntry struct {
	FrequencyHours    *int64   `json:"frequencyHours,omitempty"`
	Threshold         *float64 `json:"threshold,omitempty"`
	WarningThreshold  *float64 `json:"warningThreshold,omitempty"`
	CriticalThreshold *float64 `json:"criticalThreshold,omitempty"`
	UpdatedAt         *string  `json:"updatedAt,omitempty"`
	UpdatedBy         *string  `json:"updatedBy,omitempty"`
}

// AlertMuteState is the mute override set for one alert at one level.
type AlertMuteState struct {
	Action     string  `json:"action"`
	MutedFrom  *string `json:"mutedFrom,omitempty"`
	MutedUntil *string `json:"mutedUntil,omitempty"`
	UpdatedAt  *string `json:"updatedAt,omitempty"`
	UpdatedBy  *string `json:"updatedBy,omitempty"`
}

// GetAlertConfigurationResponse is the alert configuration of one level, as returned by the get
// alert configuration endpoints. The resolved fields are what is in effect after inheritance.
type GetAlertConfigurationResponse struct {
	Configs             map[string]AlertConfigurationEntry `json:"configs"`
	ResolvedConfigs     map[string]AlertConfigurationEntry `json:"resolvedConfigs"`
	Recipients          []string                           `json:"recipients"`
	MutedAlerts         map[string]AlertMuteState          `json:"mutedAlerts"`
	ResolvedMutedAlerts map[string]AlertMuteState          `json:"resolvedMutedAlerts"`
}

// UpdateAlertConfigurationRequest is the body of the update alert configuration endpoints. Every
// field is optional: an omitted field leaves that part of the level's configuration unchanged.
type UpdateAlertConfigurationRequest struct {
	Configs    map[string]AlertConfigurationUpdate `json:"configs,omitempty"`
	Recipients *[]string                           `json:"recipients,omitempty"`
	Alerts     []AlertMuteAction                   `json:"alerts,omitempty"`
}

// AlertConfigurationUpdate sets the configuration of one alert at one level. A field that is not
// set keeps its current value.
type AlertConfigurationUpdate struct {
	FrequencyHours    *int64 `json:"frequencyHours,omitempty"`
	Threshold         *int64 `json:"threshold,omitempty"`
	WarningThreshold  *int64 `json:"warningThreshold,omitempty"`
	CriticalThreshold *int64 `json:"criticalThreshold,omitempty"`
}

// AlertMuteAction mutes, unmutes or clears the mute override of one alert at one level.
type AlertMuteAction struct {
	AlertName  string  `json:"alertName"`
	Action     string  `json:"action"`
	MutedFrom  *string `json:"mutedFrom,omitempty"`
	MutedUntil *string `json:"mutedUntil,omitempty"`
}

const (
	AlertMuteActionMute   = "mute"
	AlertMuteActionUnmute = "unmute"
	AlertMuteActionClear  = "clear"
)
