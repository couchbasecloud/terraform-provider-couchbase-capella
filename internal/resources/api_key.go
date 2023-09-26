package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"net/http"
	"terraform-provider-capella/internal/api/api_key"

	"terraform-provider-capella/internal/api"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	ApiKeyRequest := api_key.CreateApiKeyRequest{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.Name.ValueString(),
	}
	var newOrganizationRoles []api_key.ApiKeyOrganizationRole
	for _, organizationRole := range plan.OrganizationRoles {
		newOrganizationRoles = append(newOrganizationRoles, api_key.ApiKeyOrganizationRole(organizationRole.ValueString()))
	}
	ApiKeyRequest.OrganizationRoles = newOrganizationRoles

	elements := make([]types.String, 0, len(plan.AllowedCIDRs.Elements()))
	_ = plan.AllowedCIDRs.ElementsAs(ctx, &elements, false)

	var newAllowedCidrs []string
	for _, allowedCidr := range elements {
		newAllowedCidrs = append(newAllowedCidrs, allowedCidr.ValueString())
	}

	if !plan.AllowedCIDRs.IsNull() {
		ApiKeyRequest.AllowedCIDRs = &newAllowedCidrs
	}
	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/apikeys", r.HostURL, organizationId),
		http.MethodPost,
		ApiKeyRequest,
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

	ApiKeyResponse := api_key.CreateAPIKeyResponse{}
	err = json.Unmarshal(response.Body, &ApiKeyResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ApiKey",
			"Could not create ApiKey, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := r.retrieveApiKey(ctx, organizationId, ApiKeyResponse.Id)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error Reading Capella ApiKeys",
			"Could not read Capella ApiKey ID "+ApiKeyResponse.Id+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella ApiKeys",
			"Could not read Capella ApiKey ID "+ApiKeyResponse.Id+": "+err.Error(),
		)
		return
	}

	refreshedState.Token = types.StringValue(ApiKeyResponse.Token)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads ApiKey information.
func (r *ApiKey) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//TODO
}

// Update updates the ApiKey.
func (r *ApiKey) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//TODO
}

// Delete deletes the ApiKey.
func (r *ApiKey) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//TODO
}

func (r *ApiKey) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//TODO
}

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

	apiKeyResp := api_key.GetApiKeyResponse{}
	err = json.Unmarshal(response.Body, &apiKeyResp)
	if err != nil {
		return nil, err
	}

	audit := providerschema.NewCouchbaseAuditData(apiKeyResp.Audit)

	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("error while audit conversion")
	}

	refreshedState := providerschema.ApiKey{
		Id:             types.StringValue(apiKeyResp.Id),
		OrganizationId: types.StringValue(organizationId),
		Name:           types.StringValue(apiKeyResp.Name),
		Description:    types.StringValue(apiKeyResp.Description),
		Expiry:         types.Float64Value(float64(apiKeyResp.Expiry)),
		Audit:          auditObj,
	}

	var newAllowedCidr []attr.Value
	for _, allowedCidr := range apiKeyResp.AllowedCIDRs {
		newAllowedCidr = append(newAllowedCidr, types.StringValue(allowedCidr))
	}
	//refreshedState.AllowedCIDRs

	x, _ := types.ListValue(types.StringType, newAllowedCidr)

	refreshedState.AllowedCIDRs = x

	var newOrganizationRoles []types.String
	for _, organizationRole := range apiKeyResp.OrganizationRoles {
		newOrganizationRoles = append(newOrganizationRoles, types.StringValue(string(organizationRole)))
	}
	refreshedState.OrganizationRoles = newOrganizationRoles

	var newApiKeyResourcesItems []providerschema.APIKeyResourcesItems
	for _, resource := range apiKeyResp.Resources {
		newResourceItem := providerschema.APIKeyResourcesItems{
			Id: types.StringValue(resource.Id.String()),
		}
		if resource.Type != nil {
			newResourceItem.Type = types.StringValue(*resource.Type)
		}
		var newRoles []types.String
		for _, role := range resource.Roles {
			newRoles = append(newRoles, types.StringValue(string(role)))
		}
		newResourceItem.Roles = newRoles
		newApiKeyResourcesItems = append(newApiKeyResourcesItems, newResourceItem)
	}
	refreshedState.Resources = newApiKeyResourcesItems

	return &refreshedState, nil
}
