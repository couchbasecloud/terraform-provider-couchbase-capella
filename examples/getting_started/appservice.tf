# Stores the app service details in an output variable.
# Can be viewed using `terraform output app_service` command
output "app_service" {
  value = couchbase-capella_app_service.new_app_service
}

resource "couchbase-capella_app_service" "new_app_service" {
  organization_id = var.organization_id
  project_id      = couchbase-capella_project.new_project.id
  cluster_id      = couchbase-capella_cluster.new_cluster.id
  name            = var.app_service.name
  description     = var.app_service.description
  nodes           = var.app_service.nodes
  compute = {
    cpu = var.app_service.compute.cpu
    ram = var.app_service.compute.ram
  }
  depends_on = [couchbase-capella_cluster.new_cluster, couchbase-capella_bucket.new_bucket, couchbase-capella_sample_bucket.new_sample_bucket]
}