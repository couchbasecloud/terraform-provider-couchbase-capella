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
