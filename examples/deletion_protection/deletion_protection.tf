output "deletion_protection" {
  value = couchbase-capella_cluster_deletion_protection.cluster.deletion_protection
}

resource "couchbase-capella_cluster_deletion_protection" "cluster" {
  organization_id     = var.organization_id
  project_id          = var.project_id
  cluster_id          = var.cluster_id
  deletion_protection = var.deletion_protection
}

