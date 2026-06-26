package validator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCloudProviderRegion(t *testing.T) {
	allowed := map[string][]string{
		"aws":   {"us-west-2", "eu-west-1"},
		"azure": {"eastus", "swedencentral"},
	}

	objectType := map[string]attr.Type{
		"type":   types.StringType,
		"region": types.StringType,
	}

	tests := []struct {
		name      string
		provider  attr.Value
		region    attr.Value
		wantError bool
	}{
		{
			name:     "supported region",
			provider: types.StringValue("aws"),
			region:   types.StringValue("us-west-2"),
		},
		{
			name:      "unsupported region",
			provider:  types.StringValue("aws"),
			region:    types.StringValue("us-east-99"),
			wantError: true,
		},
		{
			name:     "unknown provider type is skipped (its enum validator reports it)",
			provider: types.StringValue("gcp"),
			region:   types.StringValue("anything"),
		},
		{
			name:     "unknown region value is skipped",
			provider: types.StringValue("aws"),
			region:   types.StringUnknown(),
		},
		{
			name:     "null region value is skipped",
			provider: types.StringValue("aws"),
			region:   types.StringNull(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			obj, diags := types.ObjectValue(objectType, map[string]attr.Value{
				"type":   tc.provider,
				"region": tc.region,
			})
			if diags.HasError() {
				t.Fatalf("failed to build object: %v", diags)
			}

			resp := &validator.ObjectResponse{}
			CloudProviderRegion(allowed).ValidateObject(
				context.Background(),
				validator.ObjectRequest{ConfigValue: obj},
				resp,
			)

			if got := resp.Diagnostics.HasError(); got != tc.wantError {
				t.Errorf("HasError() = %v, want %v (diags: %v)", got, tc.wantError, resp.Diagnostics)
			}
		})
	}
}

func TestCloudProviderRegion_NullObjectSkipped(t *testing.T) {
	obj := types.ObjectNull(map[string]attr.Type{
		"type":   types.StringType,
		"region": types.StringType,
	})

	resp := &validator.ObjectResponse{}
	CloudProviderRegion(map[string][]string{"aws": {"us-west-2"}}).ValidateObject(
		context.Background(),
		validator.ObjectRequest{ConfigValue: obj},
		resp,
	)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected null object to be skipped, got: %v", resp.Diagnostics)
	}
}
