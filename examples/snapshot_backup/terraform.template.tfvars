auth_token = "<auth_token>"

organization_id = "<organization_id>"
project_id      = "<project_id>"
cluster_id      = "<cluster_id>"

// Minimum retention is 24 hours and maximum is 720 hours.
// Retention must be an integer.
// Uncomment restore_times to trigger a restore, and increment restore_times to trigger susbequent restores.
// Uncomment cross_region_restore_preference to specify the regions to restore from.
cloud_snapshot_backup = {
    retention = "<retention>"
    regions_to_copy = ["<region_1>", "<region_2>"]
    //restore_times = 1
    //cross_region_restore_preference = ["<region_1>", "<region_2>"]
}