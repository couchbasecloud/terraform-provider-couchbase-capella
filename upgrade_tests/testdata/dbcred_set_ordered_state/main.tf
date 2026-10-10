# Both dbcred_* fixtures use this configuration verbatim. They differ only in the
# element order of the prior state, which is the variable under test.
terraform {
  required_providers {
    couchbase-capella = {
      source  = "couchbasecloud/couchbase-capella"
      version = "99.0.0"
    }
  }
}

# host and authentication_token come from CAPELLA_HOST and
# CAPELLA_AUTHENTICATION_TOKEN, set by the test to placeholders.
provider "couchbase-capella" {}

resource "couchbase-capella_database_credential" "upgrade_fixture" {
  name            = "tf_upgrade_fixture"
  organization_id = "6af08c0a-8cab-4c1c-b257-b521575c16d0"
  project_id      = "f14134f2-7943-4dd9-b0dc-6f4a4d4b1b2c"
  cluster_id      = "ffffffff-aaaa-1111-2222-333333333333"

  # Every list below is deliberately in non-canonical order, so a provider that
  # stores these as Sets disagrees with this configuration and plans a diff.
  access = [
    {
      privileges = ["data_writer"]
      resources = {
        buckets = [
          {
            name = "zeta"
            scopes = [
              {
                name        = "z_scope"
                collections = ["z_col", "a_col", "m_col"]
              },
            ]
          },
          {
            name = "alpha"
          },
        ]
      }
    },
    {
      privileges = ["data_reader"]
    },
  ]
}
