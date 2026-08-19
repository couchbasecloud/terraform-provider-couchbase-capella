package utils

import (
	"math"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestInt64ValueFromUint64Ptr(t *testing.T) {
	var (
		zero     uint64 = 0
		typical  uint64 = 1000
		maxInt64 uint64 = math.MaxInt64
		overflow uint64 = math.MaxUint64
	)

	tests := []struct {
		name     string
		value    *uint64
		expected types.Int64
	}{
		{
			name:     "nil becomes null so an unreported field is not mistaken for zero",
			value:    nil,
			expected: types.Int64Null(),
		},
		{
			name:     "zero is preserved as a real value",
			value:    &zero,
			expected: types.Int64Value(0),
		},
		{
			name:     "typical value is converted unchanged",
			value:    &typical,
			expected: types.Int64Value(1000),
		},
		{
			name:     "math.MaxInt64 is representable and passes through",
			value:    &maxInt64,
			expected: types.Int64Value(math.MaxInt64),
		},
		{
			name:     "value above math.MaxInt64 clamps instead of wrapping negative",
			value:    &overflow,
			expected: types.Int64Value(math.MaxInt64),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, Int64ValueFromUint64Ptr(test.value))
		})
	}
}
