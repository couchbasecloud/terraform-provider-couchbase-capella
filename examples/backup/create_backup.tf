output "new_backup" {
  value = capella_backup.new_backup
}

resource "capella_backup" "new_backup" {
  organization_id            = var.organization_id
  project_id                 = var.project_id
  cluster_id                 = var.cluster_id
  bucket_id                  = var.bucket_id
}

// To trigger restore we need to change the resource like this

#resource "capella_backup" "new_backup" {
#  organization_id            = var.organization_id
#  project_id                 = var.project_id
#  cluster_id                 = var.cluster_id
#  bucket_id                  = var.bucket_id
#  restore = {
#    target_cluster_id = var.restore.target_cluster_id
#    source_cluster_id = var.restore.source_cluster_id
#    services = var.restore.services
#    force_updates = var.restore.force_updates
#    auto_remove_collections = var.restore.auto_remove_collections
#    filter_keys = var.restore.filter_keys
#    filter_values = var.restore.filter_values
#    include_data = var.restore.include_data
#    exclude_data = var.restore.exclude_data
#    map_data = var.restore.map_data
#    replace_ttl = var.restore.replace_ttl
#    replace_ttl_with = var.restore.replace_ttl
#  }
#  restore_times = var.restore.restore_times
#}