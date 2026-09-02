package utils

import (
	"math"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

// Int64ValueFromUint64Ptr converts an optional uint64 API field to a Terraform Int64.
// A nil pointer becomes null so callers can distinguish "not reported by the API" from a real zero.
// Terraform has no unsigned integer type, so values above math.MaxInt64 are clamped rather than
// wrapping to a negative number; the counters this is used for never approach that bound.
func Int64ValueFromUint64Ptr(v *uint64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	if *v > math.MaxInt64 {
		return types.Int64Value(math.MaxInt64)
	}
	return types.Int64Value(int64(*v))
}
