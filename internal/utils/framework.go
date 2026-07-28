package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

// BoolPointerIfKnown returns a pointer to the bool value when it is known and non-null, and nil otherwise.
func BoolPointerIfKnown(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return v.ValueBoolPointer()
}

// StringPointerIfKnown returns a pointer to the string value when it is known and non-null, and nil otherwise.
func StringPointerIfKnown(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return v.ValueStringPointer()
}

// Int64PointerIfKnown returns a pointer to the int64 value when it is known and non-null, and nil otherwise.
func Int64PointerIfKnown(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return v.ValueInt64Pointer()
}
