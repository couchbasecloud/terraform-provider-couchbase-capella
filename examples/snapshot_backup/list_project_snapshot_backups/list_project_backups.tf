output "project_backups_list" {
  value = data.couchbase-capella_cloud_project_snapshot_backups.existing_cloud_project_snapshot_backups
}

data "couchbase-capella_cloud_project_snapshot_backups" "existing_cloud_project_snapshot_backups" {
  organization_id = var.organization_id
  project_id      = var.project_id
}

# Example of optional fields when fetching snapshot backups at project level 

#data "couchbase-capella_cloud_project_snapshot_backups" "existing_cloud_project_snapshot_backups" {
#  organization_id = var.organization_id
#  project_id      = var.project_id
#  page =  var.page
#  per_page = var.per_page
#  sort_by = var.sort_by
#  sort_direction = var.sort_direction
#}
