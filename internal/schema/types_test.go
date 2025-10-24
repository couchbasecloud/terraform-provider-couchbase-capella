package schema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestStringValuesToStrings(t *testing.T) {
	tests := []struct {
		name           string
		input          []basetypes.StringValue
		expectedOutput []string
	}{
		{
			name: "[POSITIVE] Convert non-empty slice of StringValues to strings",
			input: []basetypes.StringValue{
				basetypes.NewStringValue("hello"),
				basetypes.NewStringValue("world"),
				basetypes.NewStringValue("test"),
			},
			expectedOutput: []string{"hello", "world", "test"},
		},
		{
			name:           "[POSITIVE] Convert empty slice",
			input:          []basetypes.StringValue{},
			expectedOutput: []string{},
		},
		{
			name: "[POSITIVE] Convert slice with empty strings",
			input: []basetypes.StringValue{
				basetypes.NewStringValue(""),
				basetypes.NewStringValue("non-empty"),
				basetypes.NewStringValue(""),
			},
			expectedOutput: []string{"", "non-empty", ""},
		},
		{
			name: "[POSITIVE] Convert slice with special characters",
			input: []basetypes.StringValue{
				basetypes.NewStringValue("hello@world.com"),
				basetypes.NewStringValue("test-123"),
				basetypes.NewStringValue("path/to/file"),
			},
			expectedOutput: []string{"hello@world.com", "test-123", "path/to/file"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := BaseStringsToStrings(test.input)
			assert.Equal(t, test.expectedOutput, result)
		})
	}
}

func TestStringsToStringValues(t *testing.T) {
	tests := []struct {
		name           string
		input          []string
		expectedOutput []basetypes.StringValue
	}{
		{
			name:  "[POSITIVE] Convert non-empty slice of strings to StringValues",
			input: []string{"hello", "world", "test"},
			expectedOutput: []basetypes.StringValue{
				basetypes.NewStringValue("hello"),
				basetypes.NewStringValue("world"),
				basetypes.NewStringValue("test"),
			},
		},
		{
			name:           "[POSITIVE] Convert empty slice",
			input:          []string{},
			expectedOutput: []basetypes.StringValue{},
		},
		{
			name:  "[POSITIVE] Convert slice with empty strings",
			input: []string{"", "non-empty", ""},
			expectedOutput: []basetypes.StringValue{
				basetypes.NewStringValue(""),
				basetypes.NewStringValue("non-empty"),
				basetypes.NewStringValue(""),
			},
		},
		{
			name:  "[POSITIVE] Convert slice with special characters",
			input: []string{"hello@world.com", "test-123", "path/to/file"},
			expectedOutput: []basetypes.StringValue{
				basetypes.NewStringValue("hello@world.com"),
				basetypes.NewStringValue("test-123"),
				basetypes.NewStringValue("path/to/file"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := StringsToBaseStrings(test.input)
			assert.Equal(t, len(test.expectedOutput), len(result))

			for i, expected := range test.expectedOutput {
				assert.Equal(t, expected.ValueString(), result[i].ValueString())
			}
		})
	}
}
