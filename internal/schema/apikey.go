package schema

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/errors"
)

// ApiKeyResourcesItems defines model for APIKeyResourcesItems.
type ApiKeyResourcesItems struct {
	// Id is the id of the project.
	Id types.String `tfsdk:"id"`

	// Roles is the Project Roles associated with the API key.
	// To learn more about Project Roles,
	// see [Project Roles](https://docs.couchbase.com/cloud/projects/project-roles.html).
	Roles []types.String `tfsdk:"roles"`

	// Type is the type of the resource.
	Type types.String `tfsdk:"type"`
}

// ApiKey maps ApiKey resource schema data
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
	Name types.String `tfsdk:"name"`

	// OrganizationRoles are the organization level roles granted to the API key.
	OrganizationRoles []types.String `tfsdk:"organization_roles"`

	// Resources  is the resources are the resource level permissions associated
	// with the API key. To learn more about Organization Roles, see
	// [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
	Resources []ApiKeyResourcesItems `tfsdk:"resources"`

	// Rotate is set only when updating(rotating) the API key,
	// and it should be set be set in incremental order from
	// the previously set rotate value, ideally we should start.
	// it from 1 when we are rotating for first time.
	Rotate types.Number `tfsdk:"rotate"`

	// Secret associated with API key. One has to follow the secret key policy,
	// such as allowed characters and a length of 64 characters. If this field
	// is left empty, a secret will be auto-generated.
	Secret types.String `tfsdk:"secret"`

	// Token is a confidential piece of information that is used to authorize
	// requests made to v4 endpoints.
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

	newAllowedCidrs, err := MorphAllowedCidrs(apiKey.AllowedCIDRs)
	if err != nil {
		return nil, err
	}

	newApiKey.AllowedCIDRs = newAllowedCidrs

	newApiKey.OrganizationRoles = MorphApiKeyOrganizationRoles(apiKey.OrganizationRoles)

	newApiKey.Resources = MorphApiKeyResources(apiKey.Resources)

	return &newApiKey, nil
}

// MorphAllowedCidrs is used to convert string list to basetypes.ListValue
// TODO : add unit testing
func MorphAllowedCidrs(allowedCIDRs []string) (basetypes.ListValue, error) {
	var newAllowedCidr []attr.Value
	for _, allowedCidr := range allowedCIDRs {
		newAllowedCidr = append(newAllowedCidr, types.StringValue(allowedCidr))
	}

	newAllowedCidrs, diags := types.ListValue(types.StringType, newAllowedCidr)
	if diags.HasError() {
		return types.ListUnknown(types.StringType), fmt.Errorf("error while converting allowedcidrs")
	}

	return newAllowedCidrs, nil
}

// MorphApiKeyOrganizationRoles is used to convert nested organizationRoles from
// strings to terraform type.String.
// TODO : add unit testing
func MorphApiKeyOrganizationRoles(organizationRoles []string) []basetypes.StringValue {
	var newOrganizationRoles []types.String
	for _, organizationRole := range organizationRoles {
		newOrganizationRoles = append(newOrganizationRoles, types.StringValue(string(organizationRole)))
	}
	return newOrganizationRoles
}

// MorphApiKeyResources is used to covert nested resources from strings
// to terraform types.String
// TODO : add unit testing
func MorphApiKeyResources(resources api.Resources) []ApiKeyResourcesItems {
	var newApiKeyResourcesItems []ApiKeyResourcesItems
	for _, resource := range resources {
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
	return newApiKeyResourcesItems
}

// Validate checks the validity of an API key and extracts associated IDs.
// TODO : add unit testing
func (a *ApiKey) Validate() (map[string]string, error) {
	const idDelimiter = ","
	var found bool

	organizationId := a.OrganizationId.ValueString()
	apiKeyId := a.Id.ValueString()

	// check if the id is a comma separated string of multiple IDs, usually passed during the terraform import CLI
	if a.OrganizationId.IsNull() {
		strs := strings.Split(a.Id.ValueString(), idDelimiter)
		if len(strs) != 2 {
			return nil, errors.ErrIdMissing
		}

		_, apiKeyId, found = strings.Cut(strs[0], "id=")
		if !found {
			return nil, errors.ErrApiKeyIdMissing
		}

		_, organizationId, found = strings.Cut(strs[1], "organization_id=")
		if !found {
			return nil, errors.ErrOrganizationIdMissing
		}

	}

	resourceIDs := a.generateResourceIdMap(organizationId, apiKeyId)

	err := a.checkEmpty(resourceIDs)
	if err != nil {
		return nil, fmt.Errorf("resource import unsuccessful: %s", err)
	}

	return resourceIDs, nil
}

// generateResourceIdMap is used to populate a map with selected IDs
func (a *ApiKey) generateResourceIdMap(organizationId, apiKeyId string) map[string]string {
	return map[string]string{
		OrganizationId: organizationId,
		ApiKeyId:       apiKeyId,
	}
}

// checkEmpty is used to verify that a supplied resourceId map has been populated
func (a *ApiKey) checkEmpty(resourceIdMap map[string]string) error {
	if resourceIdMap[ApiKeyId] == "" {
		return errors.ErrApiKeyIdCannotBeEmpty
	}

	if resourceIdMap[OrganizationId] == "" {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	return nil
}

// ApiKeys defines model for GetApiKeysResponse.
type ApiKeys struct {
	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// Data It contains the list of resources.
	Data []ApiKeyData `tfsdk:"data"`
}

// ApiKeyData maps api key resource schema data
type ApiKeyData struct {
	// OrganizationId is the organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// AllowedCIDRs represents the list of inbound CIDRs for the API key.
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
	Name types.String `tfsdk:"name"`

	// OrganizationRoles are the organization level roles granted to the API key.
	OrganizationRoles []types.String `tfsdk:"organization_roles"`

	// Resources are the resource level permissions associated
	// with the API key. To learn more about Organization Roles, see
	// [Organization Roles](https://docs.couchbase.com/cloud/organizations/organization-user-roles.html).
	Resources []ApiKeyResourcesItems `tfsdk:"resources"`
}

// NewApiKeyData creates a new apiKeyData object
func NewApiKeyData(apiKey *api.GetApiKeyResponse, organizationId string, auditObject basetypes.ObjectValue) (ApiKeyData, error) {
	newApiKeyData := ApiKeyData{
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
		return ApiKeyData{}, fmt.Errorf("error while converting allowedcidrs")
	}

	newApiKeyData.AllowedCIDRs = allowedCidrs

	var newOrganizationRoles []types.String
	for _, organizationRole := range apiKey.OrganizationRoles {
		newOrganizationRoles = append(newOrganizationRoles, types.StringValue(organizationRole))
	}
	newApiKeyData.OrganizationRoles = newOrganizationRoles

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
			newRoles = append(newRoles, types.StringValue(role))
		}
		newResourceItem.Roles = newRoles
		newApiKeyResourcesItems = append(newApiKeyResourcesItems, newResourceItem)
	}
	newApiKeyData.Resources = newApiKeyResourcesItems

	return newApiKeyData, nil
}

// Validate is used to verify that all the fields in the datasource
// have been populated.
func (a ApiKeys) Validate() (organizationId string, err error) {
	if a.OrganizationId.IsNull() {
		return "", errors.ErrOrganizationIdMissing
	}

	return a.OrganizationId.ValueString(), nil
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

// AreEqual returns true if the two arrays contain the same elements,
// without any extra values, False otherwise.
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
