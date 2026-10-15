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
