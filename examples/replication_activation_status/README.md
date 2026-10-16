# Capella XDCR Replication Activation Status Example

This example shows how to pause and resume an XDCR replication in Capella.

It uses the organization ID, project ID, cluster ID and replication ID. The replication must already exist, for example from the `examples/replication` example.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Before you start

- Pause and resume change one direction only. For a `oneWay` replication, use the replication ID and the source cluster ID.
- For the reverse direction of a `twoWay` replication, use the `reverse_replication_id` attribute of the replication and the target cluster ID. To pause both directions, use one resource for each direction.
- If the replication already has the requested status, Terraform does not change it.
- A replication that is starting shows as `running`, and one that is pausing shows as `paused`. A `failed` replication shows as a change in `terraform plan`, and the next `terraform apply` sends the pause or resume again.
- Destroying this resource does not change the replication. It stays paused or running.
- Both clusters must be healthy. If a cluster is not healthy, for example soon after the replication is created, Terraform tries again for up to 15 minutes.

# Example Walkthrough

In this example, we are going to do the following.

1. CREATE: Pause the replication as stated in the `create_replication_activation_status.tf` file.
2. IMPORT: Import the activation status of a replication that is not in the terraform state file.
3. UPDATE: Resume the replication.
4. DESTROY: Remove the activation status from the terraform state file. The replication does not change.

If you check the `terraform.template.tfvars` file - Make sure you copy the file to `terraform.tfvars` and update the values of the variables as per the correct organization access.

## CREATE
### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
TODO
```

### Apply the Plan, in order to pause the replication

Command: `terraform apply`

Sample Output:
```
TODO
```

### List the resources that are present in the Terraform State file.

Command: `terraform state list`

Sample Output:
```
TODO
```

## IMPORT
### Remove the resource `new_replication_activation_status` from the Terraform State file

Command: `terraform state rm couchbase-capella_replication_activation_status.new_replication_activation_status`

Sample Output:
```
TODO
```

Please note, this command only removes the resource from the Terraform State file. The replication is still paused in Capella.

### Now, let's import the resource in Terraform

Command: `terraform import couchbase-capella_replication_activation_status.new_replication_activation_status replication_id=<replication_id>,cluster_id=<cluster_id>,project_id=<project_id>,organization_id=<organization_id>`

Sample Output:
```
TODO
```

Here, we pass the IDs as a single comma-separated string.
The first ID in the string is the replication ID.
The second ID is the cluster ID, the cluster where this direction of the replication starts.
The third ID is the project ID, the project to which the cluster belongs.
The fourth ID is the organization ID, the organization to which the project belongs.

### Let's run a terraform plan to confirm that the import was successful and no resource states were impacted

Command: `terraform plan`

Sample Output:
```
TODO
```

## UPDATE
### Resume the replication by overriding the `state` variable.

Command: `terraform apply -var 'state=running'`

Sample Output:
```
TODO
```

## DESTROY
### Finally, remove the activation status from the Terraform State file

The replication does not change. It keeps the status from the last apply.

Command: `terraform destroy -var 'state=running'`

Sample Output:
```
TODO
```
