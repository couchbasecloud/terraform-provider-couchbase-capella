resource "couchbase-capella_database_credential" "new_database_credential" {
  name            = "ReadWriteOnSpecificCollections"
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  password        = "<password>"
  access = [
    {
      "privileges" : [
        "data_reader",
        "data_writer"
      ]
    }
  ]
}

# An advanced database credential is assigned capella user roles for
# fine-grained RBAC access instead of bucket-level access permissions.
# The user roles must already exist in the cluster.
resource "couchbase-capella_database_credential" "new_advanced_database_credential" {
  name            = "AdvancedCredential"
  organization_id = "<organization_id>"
  project_id      = "<project_id>"
  cluster_id      = "<cluster_id>"
  password        = "<password>"
  credential_type = "advanced"
  user_roles = [
    "developer",
    "bucket_admin"
  ]
}
