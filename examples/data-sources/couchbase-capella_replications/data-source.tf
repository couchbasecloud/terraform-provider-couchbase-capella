# Lists the replications that start on one cluster.
data "couchbase-capella_replications" "cluster_replications" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
  cluster_id      = "ffffffff-aaaa-1414-eeee-000000000000"
}

# Lists the replications of all the clusters in the project. A twoWay replication shows as two entries.
data "couchbase-capella_replications" "project_replications" {
  organization_id = "ffffffff-aaaa-1414-eeee-000000000000"
  project_id      = "ffffffff-aaaa-1414-eeee-000000000000"
}

output "cluster_replications" {
  value = data.couchbase-capella_replications.cluster_replications
}

output "project_replications" {
  value = data.couchbase-capella_replications.project_replications
}
