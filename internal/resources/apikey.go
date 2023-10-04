package resources

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-capella/internal/api"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ApiKey{}
	_ resource.ResourceWithConfigure   = &ApiKey{}
	_ resource.ResourceWithImportState = &ApiKey{}
)

// ApiKey is the ApiKey resource implementation.
type ApiKey struct {
	*providerschema.Data
}

func NewApiKey() resource.Resource {
	return &ApiKey{}
}

// Metadata returns the apiKey resource type name.
func (r *ApiKey) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apikey"
}

// Schema defines the schema for the apiKey resource.
func (r *ApiKey) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ApiKeySchema()
}

// Configure adds the provider configured client to the apiKey resource.
func (r *ApiKey) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.Data = data
}

// Create creates a new apiKey.
func (a *ApiKey) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ApiKey
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.validateCreateApiKeyRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ApiKey",
			"Could not create ApiKey, unexpected error:"+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()

	apiKeyRequest := api.CreateApiKeyRequest{
		Name:              plan.Name.ValueString(),
		OrganizationRoles: a.convertOrganizationRoles(plan.OrganizationRoles),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiKeyRequest.Description = plan.Description.ValueStringPointer()
	}

	if !plan.Expiry.IsNull() && !plan.Expiry.IsUnknown() {
		expiry := float32(plan.Expiry.ValueFloat64())
		apiKeyRequest.Expiry = &expiry
	}

	convertedResources, err := a.convertResources(plan.Resources)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ApiKey",
			"Could not create ApiKey, unexpected error:"+err.Error(),
		)
		return
	}
	apiKeyRequest.Resources = &convertedResources

	if !plan.AllowedCIDRs.IsNull() && !plan.AllowedCIDRs.IsUnknown() {
		convertedAllowedCidr, err := a.convertAllowedCidrs(ctx, plan.AllowedCIDRs)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating ApiKey",
				"Could not create ApiKey, unexpected error:"+err.Error(),
			)
			return
		}
		apiKeyRequest.AllowedCIDRs = &convertedAllowedCidr
	}

	response, err := a.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/apikeys", a.HostURL, organizationId),
		http.MethodPost,
		apiKeyRequest,
		a.Token,
		nil,
	)
	_, err = handleApiKeyError(err)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ApiKey",
			"Could not create ApiKey, unexpected error: "+err.Error(),
		)
		return
	}

	apiKeyResponse := api.CreateApiKeyResponse{}
	err = json.Unmarshal(response.Body, &apiKeyResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ApiKey",
			"Could not create ApiKey, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := a.retrieveApiKey(ctx, organizationId, apiKeyResponse.Id)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella ApiKeys",
			"Could not read Capella ApiKey ID "+apiKeyResponse.Id+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella ApiKeys",
			"Could not read Capella ApiKey ID "+apiKeyResponse.Id+": "+err.Error(),
		)
		return
	}

	resources, err := providerschema.OrderList2(plan.Resources, refreshedState.Resources)
	switch err {
	case nil:
		refreshedState.Resources = resources
	default:
		tflog.Error(ctx, err.Error())
	}

	for i, resource := range refreshedState.Resources {
		if providerschema.AreEqual(resource.Roles, plan.Resources[i].Roles) {
			refreshedState.Resources[i].Roles = plan.Resources[i].Roles
		}
	}

	if providerschema.AreEqual(refreshedState.OrganizationRoles, plan.OrganizationRoles) {
		refreshedState.OrganizationRoles = plan.OrganizationRoles
	}

	refreshedState.Token = types.StringValue(apiKeyResponse.Token)
	refreshedState.Rotate = plan.Rotate

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads ApiKey information.
func (a *ApiKey) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// TODO
}

// Update rotates the ApiKey.
func (a *ApiKey) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, state providerschema.ApiKey
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rotate api key",
			"Could not rotate api key id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		apiKeyId       = resourceIDs[providerschema.ApiKeyId]
	)

	if plan.Rotate.IsNull() || plan.Rotate.IsUnknown() {
		resp.Diagnostics.AddError(
			"Error rotating api key",
			"Could not rotate api key id "+state.Id.String()+": rotate value is not set",
		)
		return
	}

	if !state.Rotate.IsNull() && !state.Rotate.IsUnknown() {
		planRotate := *plan.Rotate.ValueBigFloat()
		stateRotate := *state.Rotate.ValueBigFloat()
		if planRotate.Cmp(&stateRotate) != 1 {
			resp.Diagnostics.AddError(
				"Error rotating api key",
				"Could not rotate api key id "+state.Id.String()+": plan rotate value is not greater than state rotate value",
			)
			return
		}
	}

	var rotateApiRequest api.RotateAPIKeyRequest
	if !plan.Secret.IsNull() || !plan.Secret.IsUnknown() {
		rotateApiRequest = api.RotateAPIKeyRequest{
			Secret: plan.Secret.ValueStringPointer(),
		}
	}

	response, err := a.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/apikeys/%s/rotate", a.HostURL, organizationId, apiKeyId),
		http.MethodPost,
		rotateApiRequest,
		a.Token,
		nil,
	)
	_, err = handleApiKeyError(err)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rotating api key",
			"Could not rotate api key id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	rotateApiKeyResponse := api.RotateAPIKeyResponse{}
	err = json.Unmarshal(response.Body, &rotateApiKeyResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rotating api key",
			"Could not rotate api key id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	currentState, err := a.retrieveApiKey(ctx, organizationId, apiKeyId)
	_, err = handleApiKeyError(err)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rotating api key",
			"Could not rotate api key id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	resources, err := providerschema.OrderList2(state.Resources, currentState.Resources)
	switch err {
	case nil:
		currentState.Resources = resources
	default:
		tflog.Error(ctx, err.Error())
	}

	for i, resource := range currentState.Resources {
		if providerschema.AreEqual(resource.Roles, state.Resources[i].Roles) {
			currentState.Resources[i].Roles = state.Resources[i].Roles
		}
	}

	if providerschema.AreEqual(currentState.OrganizationRoles, state.OrganizationRoles) {
		currentState.OrganizationRoles = state.OrganizationRoles
	}

	currentState.Secret = types.StringValue(rotateApiKeyResponse.SecretKey)
	if !state.Id.IsNull() && !state.Id.IsUnknown() && !state.Secret.IsNull() && !state.Secret.IsUnknown() {
		currentState.Token = types.StringValue(base64.StdEncoding.EncodeToString([]byte(state.Id.ValueString() + ":" + state.Secret.ValueString())))
	}
	currentState.Rotate = plan.Rotate

	// Set state to fully populated data
	diags = resp.State.Set(ctx, currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the ApiKey.
func (a *ApiKey) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state providerschema.ApiKey
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting api key",
			"Could not delete api key id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		apiKeyId       = resourceIDs[providerschema.ApiKeyId]
	)

	// Delete existing api key
	_, err = a.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/apikeys/%s", a.HostURL, organizationId, apiKeyId),
		http.MethodDelete,
		nil,
		a.Token,
		nil,
	)
	resourceNotFound, err := handleApiKeyError(err)
	if resourceNotFound {
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting api key",
			"Could not delete api key id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}
}

func (a *ApiKey) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//TODO
}

// retrieveApiKey retrieves apikey information for a specified organization and apiKeyId.
func (a *ApiKey) retrieveApiKey(ctx context.Context, organizationId, apiKeyId string) (*providerschema.ApiKey, error) {
	response, err := a.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/apikeys/%s", a.HostURL, organizationId, apiKeyId),
		http.MethodGet,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	apiKeyResp := api.GetApiKeyResponse{}
	err = json.Unmarshal(response.Body, &apiKeyResp)
	if err != nil {
		return nil, err
	}

	audit := providerschema.NewCouchbaseAuditData(apiKeyResp.Audit)

	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("error while audit conversion")
	}

	refreshedState, err := providerschema.NewApiKey(&apiKeyResp, organizationId, auditObj)
	if err != nil {
		return nil, err
	}
	return refreshedState, nil
}

// this func extract error message if error is api.Error and also checks whether error is
// resource not found
func handleApiKeyError(err error) (bool, error) {
	switch err := err.(type) {
	case nil:
		return false, nil
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			return false, fmt.Errorf(err.CompleteError())
		}
		return true, fmt.Errorf(err.CompleteError())
	default:
		return false, err
	}
}

// validateCreateApiKeyRequest validates the required fields in the create request.
func (a *ApiKey) validateCreateApiKeyRequest(plan providerschema.ApiKey) error {
	if plan.OrganizationId.IsNull() {
		return fmt.Errorf("organizationId cannot be empty")
	}
	if plan.Name.IsNull() {
		return fmt.Errorf("name cannot be empty")
	}
	if plan.OrganizationRoles == nil {
		return fmt.Errorf("organizationRoles cannot be empty")
	}
	if plan.Resources == nil {
		return fmt.Errorf("resource cannot be nil")
	}
	if !plan.Rotate.IsNull() && !plan.Rotate.IsUnknown() {
		return fmt.Errorf("rotate value should not be set")
	}
	if !plan.Secret.IsNull() && !plan.Secret.IsUnknown() {
		return fmt.Errorf("secret should not be set while create operation")
	}
	return nil
}

// convertOrganizationRoles is used to convert all roles
// in an array of basetypes.StringValue to strings.
func (a *ApiKey) convertOrganizationRoles(organizationRoles []basetypes.StringValue) []string {
	var convertedRoles []string
	for _, role := range organizationRoles {
		convertedRoles = append(convertedRoles, role.ValueString())
	}
	return convertedRoles
}

// convertResource is used to convert a resource object containing nested fields
// of type basetypes.StringValue to a resource object containing nested fields of go defined type.
func (a *ApiKey) convertResources(resources []providerschema.ApiKeyResourcesItems) ([]api.ResourcesItems, error) {
	var convertedResources []api.ResourcesItems
	for _, resource := range resources {
		id, err := uuid.Parse(resource.Id.ValueString())
		if err != nil {
			return nil, fmt.Errorf("resource id is not valid uuid")
		}
		convertedResource := api.ResourcesItems{
			Id: id,
		}

		var convertedRoles []string
		for _, role := range resource.Roles {
			convertedRoles = append(convertedRoles, role.ValueString())
		}
		convertedResource.Roles = convertedRoles

		if !resource.Type.IsNull() && !resource.Type.IsUnknown() {
			convertedResource.Type = resource.Type.ValueStringPointer()
		}
		convertedResources = append(convertedResources, convertedResource)
	}
	return convertedResources, nil
}

// convertAllowedCidrs is used to convert allowed cidrs in types.List to array of string.
func (a *ApiKey) convertAllowedCidrs(ctx context.Context, allowedCidrs types.List) ([]string, error) {
	elements := make([]types.String, 0, len(allowedCidrs.Elements()))
	diags := allowedCidrs.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, fmt.Errorf("error while extracting allowedCidrs elements")
	}

	var convertedAllowedCidrs []string
	for _, allowedCidr := range elements {
		convertedAllowedCidrs = append(convertedAllowedCidrs, allowedCidr.ValueString())
	}
	return convertedAllowedCidrs, nil
}
