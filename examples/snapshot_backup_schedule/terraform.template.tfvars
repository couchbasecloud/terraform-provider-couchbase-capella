auth_token = "<auth_token>"

organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"


// Interval must be 1, 2, 4, 6, 8, 12, or 24 hours.
// Minimum retention is 24 hours and maximum is 720 hours.
// Retention must be an integer.
// Start time must be a valid RFC3339 timestamp, with the minutes set to either 00 or 30, and the seconds set to 00.

snapshot_backup_schedule = {
    retention = "<retention>"
    interval = "<interval>"
    start_time = "<start_time>"
}