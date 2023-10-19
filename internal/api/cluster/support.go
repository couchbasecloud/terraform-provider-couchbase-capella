package cluster

// Support defines the support plan and timezone for this particular cluster.
type Support struct {
	// Plan is plan type, either 'Basic', 'Developer Pro', or 'Enterprise'.
	Plan SupportPlan `json:"plan"`

	// Timezone is the standard timezone for the cluster. Should be the TZ identifier.
	Timezone SupportTimezone `json:"timezone"`
}

// SupportPlan is plan type, either 'Basic', 'Developer Pro', or 'Enterprise'.
type SupportPlan string

// SupportTimezone is the standard timezone for the cluster. Should be the TZ identifier.
type SupportTimezone string
