output "new_onoff_schedule" {
  value = couchbase-capella_onoff_schedule.new_onoff_schedule
}

resource "couchbase-capella_onoff_schedule" "new_onoff_schedule" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  timezone        = var.onoff_schedule.timezone
  days = [
      for day in var.days : {
        state = day.state
        day   = day.day
        from  = {
          hour   = day.from.hour
          minute = day.from.minute
        }
        to    = {
          hour   = day.to.hour
          minute = day.to.minute
        }
      }
    ]
}
