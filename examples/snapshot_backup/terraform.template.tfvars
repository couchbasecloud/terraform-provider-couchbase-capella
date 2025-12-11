auth_token = "<auth_token>"

organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

// Minimum retention is 24 hours and maximum is 720 hours.
// Retention must be an integer.
cloud_snapshot_backup = {
    retention = "<retention>"
    regions_to_copy = ["<region_1>", "<region_2>"]
    restore_times = "<restore_times>"
    cross_region_restore_preference = ["<region_1>", "<region_2>"]
}