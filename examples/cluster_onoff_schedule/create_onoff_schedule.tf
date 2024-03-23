output "new_onoff_schedule" {
  value = couchbase-capella_onoff_schedule.new_onoff_schedule
}

resource "couchbase-capella_onoff_schedule" "new_onoff_schedule" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  timezone        = var.on_off_schedule.timezone
  days            = {
    day           = var.days.day
    state         = var.days.state
    from          = {
      hour        = var.from.hour
      minute      = var.from.minute
    }
    to          = {
      hour        = var.to.hour
      minute      = var.to.minute
    }
  }
}