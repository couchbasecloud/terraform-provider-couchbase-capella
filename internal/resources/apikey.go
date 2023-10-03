package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-capella/internal/api"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
func (r *ApiKey) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ApiKey
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating ApiKey",
			"Could not create ApiKey, unexpected error: organization ID cannot be empty.",
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()

	if !plan.Rotate.IsNull() && !plan.Rotate.IsUnknown() {
		if plan.Rotate.ValueBool() == true {
			resp.Diagnostics.AddError(
				"Error creating api key",
				"Could not create api key id: rotate flag should not be set or set to false",
			)
			return
		}
	}

	apiKeyRequest := api.CreateApiKeyRequest{
		Name: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiKeyRequest.Description = plan.Description.ValueStringPointer()
	}

	if !plan.Expiry.IsNull() && !plan.Expiry.IsUnknown() {
		expiry := float32(plan.Expiry.ValueFloat64())
		apiKeyRequest.Expiry = &expiry
	}

	var newOrganizationRoles []string
	for _, organizationRole := range plan.OrganizationRoles {
		newOrganizationRoles = append(newOrganizationRoles, organizationRole.ValueString())
	}
	apiKeyRequest.OrganizationRoles = newOrganizationRoles

	var newResources []api.ResourcesItems
	for _, resource := range plan.Resources {
		id, err := uuid.Parse(resource.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating ApiKey",
				"Could not create ApiKey, unexpected error: resource id is not valid uuid: err: "+err.Error(),
			)
			return
		}
		newResource := api.ResourcesItems{
			Id: id,
		}

		var newRoles []string
		for _, role := range resource.Roles {
			newRoles = append(newRoles, role.ValueString())
		}
		newResource.Roles = newRoles

		if !resource.Type.IsNull() && !resource.Type.IsUnknown() {
			newResource.Type = resource.Type.ValueStringPointer()
		}
		newResources = append(newResources, newResource)
	}
	apiKeyRequest.Resources = &newResources

	elements := make([]types.String, 0, len(plan.AllowedCIDRs.Elements()))
	_ = plan.AllowedCIDRs.ElementsAs(ctx, &elements, false)

	var newAllowedCidrs []string
	for _, allowedCidr := range elements {
		newAllowedCidrs = append(newAllowedCidrs, allowedCidr.ValueString())
	}

	if !plan.AllowedCIDRs.IsNull() {
		apiKeyRequest.AllowedCIDRs = &newAllowedCidrs
	}

	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/apikeys", r.HostURL, organizationId),
		http.MethodPost,
		apiKeyRequest,
		r.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error creating ApiKey",
			"Could not create ApiKey, unexpected error: "+err.CompleteError(),
		)
		return
	default:
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

	refreshedState, err := r.retrieveApiKey(ctx, organizationId, apiKeyResponse.Id)
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

// Update updates the ApiKey.
func (a *ApiKey) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//TODO
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

func (r *ApiKey) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO
}

// retrieveApiKey retrieves apikey information for a specified organization and apiKeyId.
func (r *ApiKey) retrieveApiKey(ctx context.Context, organizationId, apiKeyId string) (*providerschema.ApiKey, error) {
	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/apikeys/%s", r.HostURL, organizationId, apiKeyId),
		http.MethodGet,
		nil,
		r.Token,
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
