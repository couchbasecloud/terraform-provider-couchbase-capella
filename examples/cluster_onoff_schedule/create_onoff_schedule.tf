output "new_cluster_onoff_schedule" {
  value = couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule
}

resource "couchbase-capella_cluster_onoff_schedule" "new_cluster_onoff_schedule" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  timezone        = var.cluster_onoff_schedule.timezone
  days = [
    for day in var.days : {
      state = day.state
      day   = day.day
      from  = day.from
      to    = day.to
    }
  ]
}
