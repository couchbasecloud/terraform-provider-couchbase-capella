package api

// CreateClusterOnOffScheduleRequest is the request payload sent to the Capella V4 Public API
// in order to create a new cluster on/off schedule.
// In order to access this endpoint, the provided API key must have at least one of the roles referenced below:
//
// Organization Owner
// Project Owner.
type CreateClusterOnOffScheduleRequest struct {
	// Timezone for the schedule
	// Enum: "Pacific/Midway" "US/Hawaii" "US/Alaska" "US/Pacific" "US/Mountain" "US/Central" "US/Eastern"
	// "America/Puerto_Rico" "Canada/Newfoundland" "America/Argentina/Buenos_Aires" "Atlantic/Cape_Verde"
	// "Europe/London" "Europe/Amsterdam" "Europe/Athens" "Africa/Nairobi" "Asia/Tehran" "Indian/Mauritius"
	// "Asia/Karachi" "Asia/Calcutta" "Asia/Dhaka" "Asia/Bangkok" "Asia/Hong_Kong" "Asia/Tokyo" "Australia/North"
	// "Australia/Sydney" "Pacific/Ponape" "Antarctica/South_Pole"
	Timezone string `json:"timezone"`

	// Days is an array of day-wise schedule to manage the cluster on/off state.
	Days []DayItem `json:"days"`
}

// DayItem is an array of Days.
type DayItem struct {
	// State to be set for the cluster (on, off, or custom).
	//
	// On state turns the cluster on (healthy state) for the whole day.
	// Off state turns the cluster off for the whole day.
	// Custom state should be used for the days when a cluster needs to be in the turned on (healthy) state
	// during the specified time window instead of all day.
	State string `json:"state"`

	// Day of the week for scheduling on/off.
	//
	// The days of the week should be in proper sequence starting from Monday and ending with Sunday.
	// The On/Off schedule requires 7 days for the schedule, one for each day of the week.
	// There cannot be more or less than 7 days in the schedule.
	// Clusters cannot be scheduled to be off for the entire day for every day of the week.
	// Enum: "monday" "tuesday" "wednesday" "thursday" "friday" "saturday" "sunday"
	Day string `json:"day"`

	// From is the starting time boundary for the cluster on/off schedule.
	// The time boundary should be according to the following rules:
	//
	// If the schedule is a non-custom day - with state on or off, it cannot contain a time boundary.
	// If the schedule is a custom day -
	// It should contain the from time boundary. If the to time boundary is not specified then
	// the default value of 0 hour 0 minute is set and the cluster will be turned on for the entire day
	// from the time set in from time boundary.
	//
	// Time boundary should have a valid hour value. The valid hour values are from 0 to 23 inclusive.
	// Time boundary should have a valid minute value. The valid minute values are 0 and 30.
	// The from time boundary should not be later than the to time boundary.
	//
	// If the hour and minute values are not provided for the time boundaries,
	// it is set to a default value of 0 for both. (0 hour 0 minute)
	From *OnTimeBoundary `json:"from"`

	// To is the ending time boundary for the cluster on/off schedule.
	To *OnTimeBoundary `json:"to"`
}

// OnTimeBoundary corresponds to "from" and "to" time boundaries for when the cluster needs to be in the turned on
// (healthy) state on a day with "custom" scheduling timings.
type OnTimeBoundary struct {
	// Hour of the time boundary.
	// Default: 0
	Hour int64 `json:"hour"`

	// Minute of the time boundary.
	// Default: 0
	Minute int64 `json:"minute"`
}

// GetClusterOnOffScheduleResponse is the API response received from the Capella V4 Public API
// while attempting to fetch the on/off schedule for a cluster.
// In order to access this endpoint, the provided API key must have at least one of the roles referenced below:
//
// Organization Owner
// Project Owner.
type GetClusterOnOffScheduleResponse struct {
	Timezone string    `json:"timezone"`
	Days     []DayItem `json:"days"`
}

// UpdateClusterOnOffScheduleRequest is the request payload sent to the Capella V4 Public API
// in order to create a update a cluster's on/off schedule.
// In order to access this endpoint, the provided API key must have at least one of the roles referenced below:
//
// Organization Owner
// Project Owner.
type UpdateClusterOnOffScheduleRequest struct {
	Timezone string    `json:"timezone"`
	Days     []DayItem `json:"days"`
}
