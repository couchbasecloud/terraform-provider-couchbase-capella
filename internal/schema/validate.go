package schema

import (
	"fmt"
	"strings"
	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	importIds = map[string]Attr{
		"organization_id": OrganizationId,
		"project_id":      ProjectId,
		"cluster_id":      ClusterId,
		"id":              Id,
	}
)

// validateSchemaState validates that the IDs passed in as variadic
// parameters were successfully imported.
func validateSchemaState(state map[Attr]basetypes.StringValue) (map[Attr]string, error) {
	IDs, keyParams := morphState(state)

	// If the state was passed in via terraform import we need to
	// retrieve the individual IDs from the ID string.
	if checkForImportString(state[OrganizationId]) {
		var err error
		IDs, err = splitImportString(state[Id].ValueString(), keyParams)
		if err != nil {
			return nil, fmt.Errorf("failed to validate imported state: %s", err)
		}
	}

	err := checkKeysAndValues(IDs, keyParams)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs, nil
}

// morphState is used to convert a map representing the IDs associated
// with a state from type basetypes.StringValue to type String. It also returns
// a sorted slice of the keys in the map.
func morphState(state map[Attr]basetypes.StringValue) (map[Attr]string, []Attr) {
	var keyParams []Attr
	IDs := map[Attr]string{}
	for key, value := range state {
		keyParams = append(keyParams, key)
		IDs[key] = value.ValueString()
	}

	return IDs, keyParams
}

// checkForImportString is used to determine whether the state is
// preexisting or has been passed in during a terraform import
// via the CLI.
func checkForImportString(organizationId basetypes.StringValue) bool {
	return organizationId.IsNull()
}

// splitImportString is used to populate a map with each ID name
// and its corresponding value retrieved from the terraform import string.
//
// Note: The import string is passed in the following format:
// "id=100,cluster_id=200,project_id=300,organization_id=400"
func splitImportString(importString string, keyParams []Attr) (map[Attr]string, error) {
	const (
		idDelimiter     = ","
		equalsDelimiter = "="
	)

	pairs := strings.Split(importString, idDelimiter)
	if len(pairs) != len(keyParams) {
		return nil, fmt.Errorf("error parsing terraform import: %s", errors.ErrIdMissing)
	}

	// Use the equals delimiter to further split each pair and
	// retrieve each key and value.
	IDs := make(map[Attr]string)
	for _, pair := range pairs {
		keyValue := strings.SplitN(pair, equalsDelimiter, 2)
		IDs[importIds[keyValue[0]]] = keyValue[1]
	}

	return IDs, nil
}

// checkKeysAndValues is used to validate that an ID map
// has been populated with the expected ID keys and that the
// associated values are not empty.
func checkKeysAndValues(IDs map[Attr]string, keyParams []Attr) error {
	for _, key := range keyParams {
		value, ok := IDs[key]
		if !ok {
			return fmt.Errorf("terraform resource was missing: %w: %s", errors.ErrIdMissing, key)

		}
		if value == "" {
			return fmt.Errorf("terraform resource was empty: %w: %s", errors.ErrIdMissing, key)

		}
	}
	return nil
}
