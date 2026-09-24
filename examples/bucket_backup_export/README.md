# Capella Bucket Backup Export Example

This example shows how to export an existing bucket backup in Capella into a downloadable zip
archive, and how to read back a pre-signed URL to download it.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

In this example, we are going to do the following.

1. Create a new backup export job for an existing backup with the specified configuration.
2. Refresh it until the export completes.
3. Read the export job back through the data source to obtain the download URL.

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
See [Export lifecycle](#export-lifecycle) for how to do that.

## Read the download URL

Command: `terraform output download_url`

The URL is generated fresh on every read and is valid for one hour. It is marked sensitive because
it carries short lived cloud storage credentials. The data source exists so a configuration other
than the one that created the export can fetch a fresh URL for it.

## Import an existing backup export

Command: `terraform import couchbase-capella_bucket_backup_export.new_bucket_backup_export id=<export_id>,backup_id=<backup_id>,bucket_id=<bucket_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

## Destroy the backup export

Command: `terraform destroy`

There is no API to cancel or delete an export, so destroy only removes the export from Terraform
state. The archive expires from cloud storage on its own and Capella drops the export record about
7 days after the export completes.

# Export lifecycle

An export is a server side job that Capella cannot cancel or delete, so it behaves differently from
most resources over time.

| Stage                                     | `status`                  | `terraform plan` | Replacing the resource |
|-------------------------------------------|---------------------------|------------------|------------------------|
| Running                                   | `pending` or `processing` | No changes       | Fails with error 14061 |
| Up to 12 hours after completing           | `complete`                | No changes       | Fails with error 14062 |
| 12 hours to about 7 days after completing | `expired`                 | No changes       | Starts a new export    |
| About 7 days after completing             | Removed from state        | `+ create`       | Not applicable         |

- **The archive expires 12 hours after the export completes.** `status` becomes `expired` and
  `backup_download_url` becomes null. Both are computed attributes, so `terraform plan` proposes no
  change. Check `status` or `expiration` directly rather than relying on the plan. To get a new
  archive, run
  `terraform apply -replace=couchbase-capella_bucket_backup_export.new_bucket_backup_export`.
- **Replacing fails until the archive expires.** Only one active export can exist per backup
  cycle, so `-replace` and `taint` fail with error 14061 ("already pending or being processed") or
  14062 ("already completed"). Wait for `expiration` to pass and try again.
- **The export is re-created about 7 days after it completes.** Capella then drops the export
  record, so the resource is removed from state and the next `terraform apply` starts a new export
  of the same backup, which can be up to 5 TB. Configurations applied on a schedule or from CI will
  keep re-exporting without anything in the configuration changing, so remove the resource from the
  configuration once you have the download.
- **A failed export stays in state** until you replace it, since Capella keeps failed export
  records indefinitely. A failed export does not block the cycle, so `-replace` works straight away.
