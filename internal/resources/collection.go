package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	collection_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &Collection{}
	_ resource.ResourceWithConfigure   = &Collection{}
	_ resource.ResourceWithImportState = &Collection{}
)

const errorMessageWhileCollectionCreation = "There is an error during collection creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

const errorMessageAfterCollectionCreation = "Collection creation is successful, but encountered an error while checking the current" +
	" state of the collection. Please run `terraform plan` after 1-2 minutes to know the" +
	" current collection state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

// Collection is the collection resource implementation.
type Collection struct {
	*providerschema.Data
}

// NewCollection is a helper function to simplify the provider implementation.
func NewCollection() resource.Resource {
	return &Collection{}
}

// ImportState imports a remote collection that is not created by Terraform.
func (c *Collection) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import name and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("collection_name"), req, resp)
}

// Metadata returns the Collection resource type name.
func (c *Collection) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

// Schema defines the schema for the Collection resource.
func (c *Collection) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CollectionSchema()
}

// Configure It adds the provider configured api to the collection resource.
func (c *Collection) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	c.Data = data
}

// Create creates a new collection.
func (c *Collection) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.Collection
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	collectionRequest := collection_api.CreateCollectionRequest{
		Name: plan.Name.ValueString(),
	}

	// Check for optional fields
	if !plan.MaxTTL.IsNull() && !plan.MaxTTL.IsUnknown() {
		collectionRequest.MaxTTL = plan.MaxTTL.ValueInt64Pointer()
	}

	if err := c.validateCreateCollectionRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create collection request",
			"Could not create collection, unexpected error: "+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var bucketId = plan.BucketId.ValueString()
	var scopeName = plan.ScopeName.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s/collections", c.HostURL, organizationId, projectId, clusterId, bucketId, scopeName)
	cfg := collection_api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	_, err := c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		collectionRequest,
		c.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating collection",
			errorMessageWhileCollectionCreation+collection_api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeCollectionWithPlan(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := c.retrieveCollection(ctx, organizationId, projectId, clusterId, bucketId, scopeName, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Collection",
			"Could not read Capella Collection for the scope: %s "+scopeName+"."+errorMessageAfterCollectionCreation+collection_api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// initializeCollectionWithPlan initializes an instance of providerschema.Collection
// with the specified plan. It marks all computed fields as null.
func initializeCollectionWithPlan(plan providerschema.Collection) providerschema.Collection {
	if plan.MaxTTL.IsNull() || plan.MaxTTL.IsUnknown() {
		plan.MaxTTL = types.Int64Null()
	}
	return plan
}

func (c *Collection) validateCreateCollectionRequest(plan providerschema.Collection) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	if plan.BucketId.IsNull() {
		return errors.ErrBucketIdMissing
	}
	if plan.ScopeName.IsNull() {
		return errors.ErrScopeNameMissing
	}

	return c.validateCollectionAttributesTrimmed(plan)
}

func (c *Collection) validateCollectionAttributesTrimmed(plan providerschema.Collection) error {
	if (!plan.Name.IsNull() && !plan.Name.IsUnknown()) && !providerschema.IsTrimmed(plan.Name.ValueString()) {
		return fmt.Errorf("collection_name %s", errors.ErrNotTrimmed)
	}
	return nil
}

// retrieveCollection retrieves collection information from the specified organization and project using the provided bucket ID and scope name by open-api call.
func (c *Collection) retrieveCollection(ctx context.Context, organizationId, projectId, clusterId, bucketId, scopeName, collectionName string) (*providerschema.Collection, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s/collections/%s", c.HostURL, organizationId, projectId, clusterId, bucketId, scopeName, collectionName)
	cfg := collection_api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	collectionResp := collection_api.GetCollectionResponse{}
	err = json.Unmarshal(response.Body, &collectionResp)
	if err != nil {
		return nil, err
	}

	if collectionResp.Name != nil && validateCollectionNameIsSameInPlanAndState(collectionName, *collectionResp.Name) {
		collectionResp.Name = &collectionName
	}

	refreshedState := providerschema.Collection{
		ScopeName:      types.StringValue(scopeName),
		BucketId:       types.StringValue(bucketId),
		ClusterId:      types.StringValue(clusterId),
		ProjectId:      types.StringValue(projectId),
		OrganizationId: types.StringValue(organizationId),
	}

	//check nil pointers
	if collectionResp.Name != nil {
		refreshedState.Name = types.StringValue(*collectionResp.Name)
	}

	if collectionResp.MaxTTL != nil {
		refreshedState.MaxTTL = types.Int64Value(*collectionResp.MaxTTL)
	}
	return &refreshedState, nil
}

func validateCollectionNameIsSameInPlanAndState(planCollectionName, stateCollectionName string) bool {
	return strings.EqualFold(planCollectionName, stateCollectionName)
}

// Read reads the collection information.
func (c *Collection) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.Collection
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Collection in Capella",
			"Could not read Capella Collection with name "+state.Name.String()+": "+err.Error(),
		)
		return
	}
	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		bucketId       = IDs[providerschema.BucketId]
		scopeName      = IDs[providerschema.ScopeName]
		collectionName = IDs[providerschema.CollectionName]
	)

	refreshedState, err := c.retrieveCollection(ctx, organizationId, projectId, clusterId, bucketId, scopeName, collectionName)
	if err != nil {
		resourceNotFound, errString := collection_api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading collection",
			"Could not read collection name "+state.Name.String()+": "+errString,
		)
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the collection. Only maxTTL value for the collection can be updated. This endpoint only applies for clusters with server version >= 7.6.0.
// Collections cannot be updated for clusters with server version < 7.6.0.
func (c *Collection) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan providerschema.Collection
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := plan.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating collection",
			"Could not update collection for scope"+plan.ScopeName.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		bucketId       = resourceIDs[providerschema.BucketId]
		scopeName      = resourceIDs[providerschema.ScopeName]
		collectionName = resourceIDs[providerschema.CollectionName]
	)

	updateCollectionRequest := collection_api.UpdateCollectionRequest{
		MaxTTL: plan.MaxTTL.ValueInt64(),
	}

	// Update existing collection
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s/collections/%s", c.HostURL, organizationId, projectId, clusterId, bucketId, scopeName, collectionName)
	cfg := collection_api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err = c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		updateCollectionRequest,
		c.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating collection",
			"Could not update collection for scope"+plan.ScopeName.String()+": "+collection_api.ParseError(err),
		)
		return
	}

	currentState, err := c.retrieveCollection(ctx, organizationId, projectId, clusterId, bucketId, scopeName, collectionName)
	if err != nil {
		resourceNotFound, errString := collection_api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading collection",
			"Could not read collection name "+plan.Name.String()+": "+errString,
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the collection.
func (c *Collection) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.Collection
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting collection",
			"Could not delete collection name "+state.Name.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		bucketId       = resourceIDs[providerschema.BucketId]
		scopeName      = resourceIDs[providerschema.ScopeName]
		collectionName = resourceIDs[providerschema.CollectionName]
	)

	// Delete existing collection
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s/collections/%s", c.HostURL, organizationId, projectId, clusterId, bucketId, scopeName, collectionName)
	cfg := collection_api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err = c.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := collection_api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting collection",
			"Could not delete collection name "+state.Name.String()+": "+errString,
		)
		return
	}

}
