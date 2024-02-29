package schema

import (
	"fmt"

	scheduleapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster_onoff_schedule"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ClusterOnOffSchedule defines the response as received from V4 Capella Public API when asked to create a new cluster activation schedule.
// Couchbase supports the feature where you can schedule when your provisioned database is on and off to save costs.
// Turning off your database turns off the compute for your cluster but the storage remains.
// All of your data, schema (buckets, scopes, and collections), and indexes remain, as well as your cluster configuration,
// including users and allow lists. When you turn your provisioned database off,
// you will be charged the OFF amount for the database.
type ClusterOnOffSchedule struct {
	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the projectId of the capella tenant.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the capella tenant.
	ClusterId types.String `tfsdk:"cluster_id"`

	// Timezone for the schedule
	// Enum: "Pacific/Midway" "US/Hawaii" "US/Alaska" "US/Pacific" "US/Mountain" "US/Central" "US/Eastern"
	// "America/Puerto_Rico" "Canada/Newfoundland" "America/Argentina/Buenos_Aires" "Atlantic/Cape_Verde"
	// "Europe/London" "Europe/Amsterdam" "Europe/Athens" "Africa/Nairobi" "Asia/Tehran" "Indian/Mauritius"
	// "Asia/Karachi" "Asia/Calcutta" "Asia/Dhaka" "Asia/Bangkok" "Asia/Hong_Kong" "Asia/Tokyo" "Australia/North"
	// "Australia/Sydney" "Pacific/Ponape" "Antarctica/South_Pole"
	Timezone types.String `tfsdk:"timezone"`

	// Days is an array of day-wise schedule to manage the cluster on/off state.
	Days []DayItem `tfsdk:"days"`
}

// DayItem is an array of Days.
type DayItem struct {
	// State to be set for the cluster (on, off, or custom).
	//
	// On state turns the cluster on (healthy state) for the whole day.
	// Off state turns the cluster off for the whole day.
	// Custom state should be used for the days when a cluster needs to be in the turned on (healthy) state
	// during the specified time window instead of all day.
	State types.String `tfsdk:"state"`

	// Day of the week for scheduling on/off.
	//
	// The days of the week should be in proper sequence starting from Monday and ending with Sunday.
	// The On/Off schedule requires 7 days for the schedule, one for each day of the week.
	// There cannot be more or less than 7 days in the schedule.
	// Clusters cannot be scheduled to be off for the entire day for every day of the week.
	// Enum: "monday" "tuesday" "wednesday" "thursday" "friday" "saturday" "sunday"
	Day types.String `tfsdk:"day"`

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
	From *OnTimeBoundary `tfsdk:"from"`
	// To is the ending time boundary for the cluster on/off schedule.
	To *OnTimeBoundary `tfsdk:"to"`
}

// OnTimeBoundary corresponds to "from" and "to" time boundaries for when the cluster needs to be in the turned on
// (healthy) state on a day with "custom" scheduling timings.
type OnTimeBoundary struct {
	// Hour of the time boundary.
	// Default: 0
	Hour types.Int64 `tfsdk:"hour"`

	// Minute of the time boundary.
	// Default: 0
	Minute types.Int64 `tfsdk:"minute"`
}

// Validate checks the validity of an API key and extracts associated IDs.
func (a *ClusterOnOffSchedule) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
	}

	IDs, err := validateSchemaState(state, ClusterId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

// NewClusterOnOffSchedule creates new cluster on/off schedule object.
func NewClusterOnOffSchedule(onOffSchedule *scheduleapi.GetClusterOnOffScheduleResponse,
	organizationId, projectId, clusterId string,
) *ClusterOnOffSchedule {
	var days = make([]DayItem, 0)

	for _, d := range onOffSchedule.Days {
		days = append(days, DayItem{
			Day:   types.StringValue(d.Day),
			State: types.StringValue(d.State),
			From: &OnTimeBoundary{
				Hour:   types.Int64Value(d.From.Hour),
				Minute: types.Int64Value(d.From.Minute),
			},
			To: &OnTimeBoundary{
				Hour:   types.Int64Value(d.To.Hour),
				Minute: types.Int64Value(d.To.Minute),
			},
		})
	}

	newObj := ClusterOnOffSchedule{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		Timezone:       types.StringValue(onOffSchedule.Timezone),
		Days:           days,
	}
	return &newObj
}
