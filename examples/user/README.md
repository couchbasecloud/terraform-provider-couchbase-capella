# Capella User Example

This example shows how to create and manage users in Capella.

This creates a new user in the selected Capella project. It uses the organization ID and projectId to do so.

An invitation email is triggered and sent to the user. Upon receiving the invitation email, the user is required to click on a provided URL, which will redirect them to a page with a user interface (UI) where they can set their username and password.

The modification of any personal information related to a user can only be performed by the user through the UI. Similarly, the user can solely conduct password updates through the UI.

The "caller" possessing Organization Owner access rights retains the exclusive user creation capability. They hold the authority to assign roles at the organization and project levels.

At present, our support is limited to the capella resourceType of "project" exclusively.

To run, configure your Couchbase Capella provider as described in README in the root of this project.

# Example Walkthrough

// TODO (AV-62914) - Fill in examples after implementing logic. 

### View the plan for the resources that Terraform will create

### Apply the Plan, in order to create a new User in Capella
