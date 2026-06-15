package validator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var testRegions = map[string][]string{
	"aws":   {"us-west-2", "eu-west-1"},
	"azure": {"eastus", "swedencentral"},
}

func cloudProviderObject(t *testing.T, provider, region attr.Value) types.Object {
	t.Helper()
	obj, diags := types.ObjectValue(
		map[string]attr.Type{
			"type":   types.StringType,
			"region": types.StringType,
		},
		map[string]attr.Value{
			"type":   provider,
			"region": region,
		},
	)
	if diags.HasError() {
		t.Fatalf("failed to build object: %v", diags)
	}
	return obj
}

func runRegionValidator(t *testing.T, obj types.Object) *validator.ObjectResponse {
	t.Helper()
	v := CloudProviderRegion(testRegions)
	resp := &validator.ObjectResponse{}
	v.ValidateObject(context.Background(), validator.ObjectRequest{ConfigValue: obj}, resp)
	return resp
}

func TestCloudProviderRegion_SupportedRegion(t *testing.T) {
	obj := cloudProviderObject(t, types.StringValue("aws"), types.StringValue("us-west-2"))
	if resp := runRegionValidator(t, obj); resp.Diagnostics.HasError() {
		t.Errorf("expected no error for supported region, got: %v", resp.Diagnostics)
	}
}

func TestCloudProviderRegion_UnsupportedRegion(t *testing.T) {
	obj := cloudProviderObject(t, types.StringValue("aws"), types.StringValue("us-east-99"))
	resp := runRegionValidator(t, obj)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected an error for unsupported region, got none")
	}
}

func TestCloudProviderRegion_CaseInsensitiveProviderType(t *testing.T) {
	obj := cloudProviderObject(t, types.StringValue("AWS"), types.StringValue("eu-west-1"))
	if resp := runRegionValidator(t, obj); resp.Diagnostics.HasError() {
		t.Errorf("expected provider type match to be case-insensitive, got: %v", resp.Diagnostics)
	}
}

func TestCloudProviderRegion_UnknownProviderSkipped(t *testing.T) {
	// gcp is absent from testRegions; the type enum validator owns that error.
	obj := cloudProviderObject(t, types.StringValue("gcp"), types.StringValue("anything"))
	if resp := runRegionValidator(t, obj); resp.Diagnostics.HasError() {
		t.Errorf("expected unknown provider type to be skipped, got: %v", resp.Diagnostics)
	}
}

func TestCloudProviderRegion_UnknownValuesSkipped(t *testing.T) {
	obj := cloudProviderObject(t, types.StringValue("aws"), types.StringUnknown())
	if resp := runRegionValidator(t, obj); resp.Diagnostics.HasError() {
		t.Errorf("expected unknown region to be skipped, got: %v", resp.Diagnostics)
	}
}

func TestCloudProviderRegion_NullObjectSkipped(t *testing.T) {
	obj := types.ObjectNull(map[string]attr.Type{
		"type":   types.StringType,
		"region": types.StringType,
	})
	if resp := runRegionValidator(t, obj); resp.Diagnostics.HasError() {
		t.Errorf("expected null object to be skipped, got: %v", resp.Diagnostics)
	}
}