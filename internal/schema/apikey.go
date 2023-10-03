package schema

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-capella/internal/api"
)

// ApiKeyResourcesItems defines model for APIKeyResourcesItems.
type ApiKeyResourcesItems struct {
	// Id is the id of the project.
	Id types.String `tfsdk:"id"`

	// Roles is the Project Roles associated with the API key.
	// To learn more about Project Roles, see [Project Roles](https://docs.couchbase.com/cloud/projects/project-roles.html).
	Roles []types.String `tfsdk:"roles"`

	// Type is the type of the resource.
	Type types.String `tfsdk:"type"`
}

type ApiKey struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// AllowedCIDRs is the list of inbound CIDRs for the API key.
	// The system making a request must come from one of the allowed CIDRs.
	AllowedCIDRs types.List   `tfsdk:"allowed_cidrs"`
	Audit        types.Object `tfsdk:"audit"`

	// Description is the description for the API key.
	Description types.String `tfsdk:"description"`

	// Expiry is the expiry of the API key in number of days.
	// If set to -1, the token will not expire.
	Expiry types.Float64 `tfsdk:"expiry"`

	// Id is the id is a unique identifier for an apiKey.
	Id types.String `tfsdk:"id"`

	// Name is the name of the API key.
	Name              types.String   `tfsdk:"name"`
	OrganizationRoles []types.String `tfsdk:"organization_roles"`

	// Resources  is the resources are the resource level permissions associated with the API key.
	// To learn more about Organization Roles, see [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
	Resources []ApiKeyResourcesItems `tfsdk:"resources"`

	Rotate types.Bool `tfsdk:"rotate"`

	Secret types.String `tfsdk:"secret"`

	Token types.String `tfsdk:"token"`
}

// NewApiKey creates new apikey object
func NewApiKey(apiKey *api.GetApiKeyResponse, organizationId string, auditObject basetypes.ObjectValue) (*ApiKey, error) {
	newApiKey := ApiKey{
		Id:             types.StringValue(apiKey.Id),
		OrganizationId: types.StringValue(organizationId),
		Name:           types.StringValue(apiKey.Name),
		Description:    types.StringValue(apiKey.Description),
		Expiry:         types.Float64Value(float64(apiKey.Expiry)),
		Audit:          auditObject,
	}

	var newAllowedCidr []attr.Value
	for _, allowedCidr := range apiKey.AllowedCIDRs {
		newAllowedCidr = append(newAllowedCidr, types.StringValue(allowedCidr))
	}

	allowedCidrs, diags := types.ListValue(types.StringType, newAllowedCidr)
	if diags.HasError() {
		return nil, fmt.Errorf("error while converting allowedcidrs")
	}

	newApiKey.AllowedCIDRs = allowedCidrs

	var newOrganizationRoles []types.String
	for _, organizationRole := range apiKey.OrganizationRoles {
		newOrganizationRoles = append(newOrganizationRoles, types.StringValue(string(organizationRole)))
	}
	newApiKey.OrganizationRoles = newOrganizationRoles

	var newApiKeyResourcesItems []ApiKeyResourcesItems
	for _, resource := range apiKey.Resources {
		newResourceItem := ApiKeyResourcesItems{
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
	newApiKey.Resources = newApiKeyResourcesItems

	return &newApiKey, nil
}

// OrderList2 function to order list2 based on list1's Ids
func OrderList2(list1, list2 []ApiKeyResourcesItems) ([]ApiKeyResourcesItems, error) {
	if len(list1) != len(list2) {
		return nil, fmt.Errorf("returned resources is not same as in plan")
	}
	// Create a map from Id to APIKeyResourcesItems for list2
	idToItem := make(map[string]ApiKeyResourcesItems)
	for _, item := range list2 {
		idToItem[item.Id.ValueString()] = item
	}

	// Create a new ordered list2 based on the order of list1's Ids
	orderedList2 := make([]ApiKeyResourcesItems, len(list1))
	for i, item1 := range list1 {
		orderedList2[i] = idToItem[item1.Id.ValueString()]
	}

	if len(orderedList2) != len(list2) {
		return nil, fmt.Errorf("returned resources is not same as in plan")
	}

	return orderedList2, nil
}

// AreEqual returns true if the two arrays contain the same elements, without any extra values, False otherwise.
func AreEqual[T comparable](array1 []T, array2 []T) bool {
	if len(array1) != len(array2) {
		return false
	}
	set1 := make(map[T]bool)
	for _, element := range array1 {
		set1[element] = true
	}

	for _, element := range array2 {
		if !set1[element] {
			return false
		}
	}

	return len(set1) == len(array1)
}
