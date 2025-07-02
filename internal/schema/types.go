package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// BaseStringsToStrings converts a slice of basetypes.StringValue to a slice of strings.
// It extracts the string values from each basetypes.StringValue element.
func BaseStringsToStrings(stringValues []basetypes.StringValue) []string {
	result := make([]string, len(stringValues))
	for i, sv := range stringValues {
		result[i] = sv.ValueString()
	}
	return result
}

// StringsToBaseStrings converts a slice of strings to a slice of basetypes.StringValue.
// It wraps each string in a basetypes.StringValue.
func StringsToBaseStrings(strings []string) []basetypes.StringValue {
	result := make([]basetypes.StringValue, len(strings))
	for i, s := range strings {
		result[i] = basetypes.NewStringValue(s)
	}
	return result
}
