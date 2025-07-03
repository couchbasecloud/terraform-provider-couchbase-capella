package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// BaseStringsToStrings converts a slice of basetypes.StringValue to a slice of strings.
// It extracts the string values from each basetypes.StringValue element.
// If a string value is nil or unknown, it will be represented as an empty string in the resulting slice.
func BaseStringsToStrings(stringValues []basetypes.StringValue) []string {
	result := make([]string, len(stringValues))
	for i, sv := range stringValues {
		result[i] = sv.ValueString()
	}
	return result
}

// StringsToBaseStrings converts a slice of strings to a slice of basetypes.StringValue.
// It wraps each string in a basetypes.StringValue.
func StringsToBaseStrings(genericStrings []string) []basetypes.StringValue {
	result := make([]basetypes.StringValue, len(genericStrings))
	for i, s := range genericStrings {
		result[i] = basetypes.NewStringValue(s)
	}
	return result
}
