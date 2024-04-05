# Stores the cluster onoff schedule details in an output variable.
# Can be viewed using `terraform output cluster_onoff_schedule` command
output "cluster_onoff_schedule" {
  value = couchbase-capella_cluster_onoff_schedule.new_cluster_onoff_schedule
}

resource "couchbase-capella_cluster_onoff_schedule" "new_cluster_onoff_schedule" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
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
