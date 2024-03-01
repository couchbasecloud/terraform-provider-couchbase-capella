output "new_samplebucket" {
  value = couchbase-capella_samplebucket.new_samplebucket
}

output "samplebucket_id" {
  value = couchbase-capella_samplebucket.new_samplebucket.id
}

resource "couchbase-capella_samplebucket" "new_samplebucket" {
  name             = var.samplebucket.name
  organization_id  = var.organization_id
  project_id       = var.project_id
  cluster_id       = var.cluster_id
}
