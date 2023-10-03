output "new_bucket" {
  value = capella_bucket.new_bucket
}

resource "capella_bucket" "new_bucket" {
  name            = var.bucket_name
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
}
