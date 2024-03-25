output "existing_onoff_schedule" {
  value = data.couchbase-capella_onoff_schedule.existing_onoff_schedule
}

data "couchbase-capella_onoff_schedule" "existing_onoff_schedule" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
