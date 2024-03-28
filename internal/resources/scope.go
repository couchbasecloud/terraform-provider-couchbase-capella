package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	scope_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/scope"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &Scope{}
	_ resource.ResourceWithConfigure   = &Scope{}
	_ resource.ResourceWithImportState = &Scope{}
)

const errorMessageWhileScopeCreation = "There is an error during scope creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

const errorMessageAfterScopeCreation = "Scope creation is successful, but encountered an error while checking the current" +
	" state of the scope. Please run `terraform plan` after 1-2 minutes to know the" +
	" current scope state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

// Scope is the scope resource implementation.
type Scope struct {
	*providerschema.Data
}

// NewScope is a helper function to simplify the provider implementation.
func NewScope() resource.Resource {
	return &Scope{}
}

// ImportState imports a remote scope that is not created by Terraform.
func (s *Scope) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import name and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("scope_name"), req, resp)
}

// Metadata returns the Scope resource type name.
func (s *Scope) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scope"
}

// Schema defines the schema for the Scope resource.
func (s *Scope) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ScopeSchema()
}

// Configure It adds the provider configured api to the scope resource.
func (s *Scope) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	s.Data = data
}

// Create creates a new scope.
func (s *Scope) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.Scope
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	scopeRequest := scope_api.CreateScopeRequest{
		Name: plan.Name.ValueString(),
	}

	if err := s.validateCreateScopeRequest(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing create scope request",
			"Could not create scope, unexpected error: "+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()
	var projectId = plan.ProjectId.ValueString()
	var clusterId = plan.ClusterId.ValueString()
	var bucketId = plan.BucketId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes", s.HostURL, organizationId, projectId, clusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	_, err := s.Client.ExecuteWithRetry(
		ctx,
		cfg,
		scopeRequest,
		s.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating scope",
			errorMessageWhileScopeCreation+api.ParseError(err),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeScopeWithPlan(plan))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err, diag := s.retrieveScope(ctx, organizationId, projectId, clusterId, bucketId, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Scope",
			"Could not read Capella Scope for the bucket: %s "+bucketId+"."+errorMessageAfterScopeCreation+api.ParseError(err),
		)
		return
	}
	if diag.HasError() {
		diags.Append(diag...)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// initializeScopeWithPlan initializes an instance of providerschema.Scope
// with the specified plan. It marks all computed fields as null.
func initializeScopeWithPlan(plan providerschema.Scope) providerschema.Scope {
	if plan.Collections.IsNull() || plan.Collections.IsUnknown() {
		plan.Collections = types.SetNull(types.ObjectType{}.WithAttributeTypes(providerschema.CollectionAttributeTypes()))
	}
	types.SetNull(types.SetType{})
	return plan
}

func (s *Scope) validateCreateScopeRequest(plan providerschema.Scope) error {
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
	return s.validateScopeAttributesTrimmed(plan)
}

func (s *Scope) validateScopeAttributesTrimmed(plan providerschema.Scope) error {
	if (!plan.Name.IsNull() && !plan.Name.IsUnknown()) && !providerschema.IsTrimmed(plan.Name.ValueString()) {
		return fmt.Errorf("name %s", errors.ErrNotTrimmed)
	}
	return nil
}

// retrieveScope retrieves scope information from the specified organization and project using the provided bucket ID by open-api call.
func (s *Scope) retrieveScope(ctx context.Context, organizationId, projectId, clusterId, bucketId, scopeName string) (*providerschema.Scope, error, diag.Diagnostics) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s", s.HostURL, organizationId, projectId, clusterId, bucketId, scopeName)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := s.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		s.Token,
		nil,
	)
	if err != nil {
		return nil, err, nil
	}

	scopeResp := scope_api.GetScopeResponse{}
	err = json.Unmarshal(response.Body, &scopeResp)
	if err != nil {
		return nil, err, nil
	}

	if validateScopeNameIsSameInPlanAndState(scopeName, *scopeResp.Name) {
		scopeResp.Name = &scopeName
	}

	refreshedState := providerschema.Scope{
		Name:           types.StringValue(*scopeResp.Name),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		BucketId:       types.StringValue(bucketId),
	}

	objectList := make([]types.Object, 0)
	//traverse collections
	for _, apiCollection := range *scopeResp.Collections {
		//create object
		providerschemaCollection := providerschema.NewCollection(apiCollection)
		collectionObj, diag := types.ObjectValueFrom(ctx, providerschema.CollectionAttributeTypes(), providerschemaCollection)
		if diag.HasError() {
			return nil, fmt.Errorf("collection object error"), diag
		}
		objectList = append(objectList, collectionObj)
	}

	//create collection set
	//using a set instead of list as we need unique list of collections, no duplicates
	collectionSet, diag := types.SetValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(providerschema.CollectionAttributeTypes()), objectList)
	if diag.HasError() {
		return nil, fmt.Errorf("collection set error"), diag
	}

	refreshedState.Collections = collectionSet

	return &refreshedState, nil, nil
}

func validateScopeNameIsSameInPlanAndState(planScopeName, stateScopeName string) bool {
	return strings.EqualFold(planScopeName, stateScopeName)
}

// Read reads the scope information.
func (s *Scope) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.Scope
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Scope in Capella",
			"Could not read Capella Scope with name "+state.Name.String()+": "+err.Error(),
		)
		return
	}
	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		bucketId       = IDs[providerschema.BucketId]
		scopeName      = IDs[providerschema.ScopeName]
	)

	refreshedState, err, diag := s.retrieveScope(ctx, organizationId, projectId, clusterId, bucketId, scopeName)
	if diag.HasError() {

		diags.Append(diag...)
		resp.Diagnostics.Append(diags...)
		return
	}
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading scope",
			"Could not read scope name "+state.Name.String()+": "+errString,
		)
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the scope.
func (s *Scope) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// Couchbase Capella's v4 does not support a PUT endpoint for scopes.
	// Scopes can only be created, read and deleted.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/buckets-scopes-and-collections
	//
	// Note: In this situation, terraform apply will default to deleting and executing a new create.
	// The update implementation should simply be left empty.
	// https://developer.hashicorp.com/terraform/plugin/framework/resources/update
}

// Delete deletes the scope.
func (s *Scope) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.Scope
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting scope",
			"Could not delete scope name "+state.Name.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
		bucketId       = resourceIDs[providerschema.BucketId]
		scopeName      = resourceIDs[providerschema.ScopeName]
	)

	// Delete existing scope
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s", s.HostURL, organizationId, projectId, clusterId, bucketId, scopeName)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err = s.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		s.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting scope",
			"Could not delete scope name "+state.Name.String()+": "+errString,
		)
		return
	}

}
