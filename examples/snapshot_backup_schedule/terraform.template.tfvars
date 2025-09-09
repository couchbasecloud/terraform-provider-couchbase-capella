auth_token = "<auth_token>"

organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"


// Minimum interval is 1 hour and maximum is 24 hours.
// Minimum retention is 24 hours and maximum is 720 hours.

snapshot_backup_schedule = {
    retention = "<retention>"
    interval = "<interval>"
    start_time = "<start_time>"
}