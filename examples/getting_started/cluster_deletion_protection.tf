output "cluster_deletion_protection" {
  value = couchbase-capella_cluster_deletion_protection.new_cluster.deletion_protection
}

resource "couchbase-capella_cluster_deletion_protection" "new_cluster" {
  organization_id     = var.organization_id
  project_id          = couchbase-capella_project.new_project.id
  cluster_id          = couchbase-capella_cluster.new_cluster.id
  deletion_protection = var.deletion_protection
}

