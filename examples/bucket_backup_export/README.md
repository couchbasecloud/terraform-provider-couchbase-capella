# Capella Bucket Backup Export Example

This example shows how to export an existing bucket backup in Capella into a downloadable zip
archive.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. Create a new backup export job for an existing backup with the specified configuration.
2. Refresh it until the export completes and exposes a download URL.

`bucket_id` is the ID Capella returns for the bucket, which is the URL-compatible base64 encoding of
the bucket name, not the bucket name itself. `backup_id` is the ID of the backup to export, which
you can get from the `couchbase-capella_backups` data source.

The export runs asynchronously on the backup infrastructure, so `terraform apply` returns while the
job is still `pending`. Refresh until `status` is `complete` to get `backup_download_url`.

A backup can only be exported if it is in a `ready` state, is at most 5 TB, and belongs to a
finished backup cycle. Only one active export can exist per backup cycle.

## Create the backup export

Command: `terraform plan`

Sample Output:
```
$ terraform plan

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # couchbase-capella_bucket_backup_export.new_bucket_backup_export will be created
  + resource "couchbase-capella_bucket_backup_export" "new_bucket_backup_export" {
      + backup_download_url = (sensitive value)
      + backup_id           = "ffffffff-aaaa-1414-eeee-000000000000"
      + bucket_id           = "dHJhdmVsLXNhbXBsZQ=="
      + bucket_name         = (known after apply)
      + cluster_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + created_at          = (known after apply)
      + cycle_id            = (known after apply)
      + expiration          = (known after apply)
      + id                  = (known after apply)
      + organization_id     = "ffffffff-aaaa-1414-eeee-000000000000"
      + project_id          = "ffffffff-aaaa-1414-eeee-000000000000"
      + sha256_checksum     = (known after apply)
      + size_in_bytes       = (known after apply)
      + status              = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

Command: `terraform apply`

## Poll until the export completes

Command: `terraform apply --refresh-only`

Repeat until `status` is `complete`. The archive is deleted from cloud storage at `expiration`,
roughly 12 hours after the export completes, after which the backup must be exported again.

## Read the download URL

Command: `terraform output download_url`

The URL is generated fresh on every refresh and is valid for one hour. It is marked sensitive
because it carries short lived cloud storage credentials.

## Import an existing backup export

Command: `terraform import couchbase-capella_bucket_backup_export.new_bucket_backup_export id=<export_id>,backup_id=<backup_id>,bucket_id=<bucket_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

## Destroy the backup export

Command: `terraform destroy`

There is no API to cancel or delete an export, so destroy only removes the export from Terraform
state. The archive expires from cloud storage on its own and Capella drops the export record after
about 7 days.
