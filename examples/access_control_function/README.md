# Capella Access Control Function Example

This example shows how to create and manage access control and validation functions for App Endpoints in Capella.

This creates a new access control function for a specific collection within an App Endpoint. Access control functions are JavaScript functions that specify access control policies applied to documents in collections. Every document update is processed by this function.

The default access control function is `function(doc){channel(doc.channels);}` for the default collection and `function(doc){channel(collectionName);}` for named collections.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

## Example Walkthrough

Command: `terraform apply`

Sample Output:
```
# (Output will be added after implementation)
```

### View the plan for the resources that Terraform will create

Command: `terraform plan`

Sample Output:
```
# (Output will be added after implementation)
```

### Apply the Plan, in order to create a new Access Control Function in Capella

Command: `terraform apply`

Sample Output:
```
# (Output will be added after implementation)
```

### View the current list of resources that are present in Terraform State

Command: `terraform show`

Sample Output:
```
# (Output will be added after implementation)
```

### Update the access control function

Make changes to the `access_control_function` attribute in `create_access_control_function.tf` and run:

Command: `terraform apply`

Sample Output:
```
# (Output will be added after implementation)
```

### Delete the resources that Terraform manages

Command: `terraform destroy`

Sample Output:
```
# (Output will be added after implementation)
```

## Prerequisites

- Couchbase Capella organization with appropriate permissions
- Existing App Service and App Endpoint
- Valid scope and collection within the App Endpoint

## Important Notes

- Access control functions are applied at the collection level
- Functions must be valid JavaScript
- Changes to access control functions affect all documents in the collection
- The resource requires `app_endpoint_name` rather than ID for identification 