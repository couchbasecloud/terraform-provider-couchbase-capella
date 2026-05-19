package validator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExactlyOneOfNested_ValidSingleSpecified(t *testing.T) {
	v := ExactlyOneOfNested("aws_config", "gcp_config", "azure_config")

	// Create object with only aws_config specified (has non-null child)
	awsConfig, _ := types.ObjectValue(
		map[string]attr.Type{"account_id": types.StringType},
		map[string]attr.Value{"account_id": types.StringValue("123")},
	)
	gcpConfig := types.ObjectNull(map[string]attr.Type{"project_id": types.StringType})
	azureConfig := types.ObjectNull(map[string]attr.Type{"tenant_id": types.StringType})

	objValue, _ := types.ObjectValue(
		map[string]attr.Type{
			"aws_config":   awsConfig.Type(context.Background()),
			"gcp_config":   gcpConfig.Type(context.Background()),
			"azure_config": azureConfig.Type(context.Background()),
		},
		map[string]attr.Value{
			"aws_config":   awsConfig,
			"gcp_config":   gcpConfig,
			"azure_config": azureConfig,
		},
	)

	req := validator.ObjectRequest{
		ConfigValue: objValue,
	}
	resp := &validator.ObjectResponse{}

	v.ValidateObject(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no error for single specified, got: %v", resp.Diagnostics)
	}
}

func TestExactlyOneOfNested_InvalidMultipleSpecified(t *testing.T) {
	v := ExactlyOneOfNested("aws_config", "gcp_config", "azure_config")

	// Create object with both aws_config and gcp_config specified
	awsConfig, _ := types.ObjectValue(
		map[string]attr.Type{"account_id": types.StringType},
		map[string]attr.Value{"account_id": types.StringValue("123")},
	)
	gcpConfig, _ := types.ObjectValue(
		map[string]attr.Type{"project_id": types.StringType},
		map[string]attr.Value{"project_id": types.StringValue("my-project")},
	)
	azureConfig := types.ObjectNull(map[string]attr.Type{"tenant_id": types.StringType})

	objValue, _ := types.ObjectValue(
		map[string]attr.Type{
			"aws_config":   awsConfig.Type(context.Background()),
			"gcp_config":   gcpConfig.Type(context.Background()),
			"azure_config": azureConfig.Type(context.Background()),
		},
		map[string]attr.Value{
			"aws_config":   awsConfig,
			"gcp_config":   gcpConfig,
			"azure_config": azureConfig,
		},
	)

	req := validator.ObjectRequest{
		ConfigValue: objValue,
	}
	resp := &validator.ObjectResponse{}

	v.ValidateObject(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Error("Expected error for multiple specified, got none")
	}
}

func TestExactlyOneOfNested_InvalidNoneSpecified(t *testing.T) {
	v := ExactlyOneOfNested("aws_config", "gcp_config", "azure_config")

	// Create object with all configs null
	awsConfig := types.ObjectNull(map[string]attr.Type{"account_id": types.StringType})
	gcpConfig := types.ObjectNull(map[string]attr.Type{"project_id": types.StringType})
	azureConfig := types.ObjectNull(map[string]attr.Type{"tenant_id": types.StringType})

	objValue, _ := types.ObjectValue(
		map[string]attr.Type{
			"aws_config":   awsConfig.Type(context.Background()),
			"gcp_config":   gcpConfig.Type(context.Background()),
			"azure_config": azureConfig.Type(context.Background()),
		},
		map[string]attr.Value{
			"aws_config":   awsConfig,
			"gcp_config":   gcpConfig,
			"azure_config": azureConfig,
		},
	)

	req := validator.ObjectRequest{
		ConfigValue: objValue,
	}
	resp := &validator.ObjectResponse{}

	v.ValidateObject(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Error("Expected error for none specified, got none")
	}
}

func TestExactlyOneOfNested_EmptyObjectNotCounted(t *testing.T) {
	v := ExactlyOneOfNested("aws_config", "gcp_config")

	// Create object with aws_config having actual value, gcp_config being empty object
	awsConfig, _ := types.ObjectValue(
		map[string]attr.Type{"account_id": types.StringType},
		map[string]attr.Value{"account_id": types.StringValue("123")},
	)
	// Empty object (all null children) - should NOT be counted as specified
	gcpConfig, _ := types.ObjectValue(
		map[string]attr.Type{"project_id": types.StringType},
		map[string]attr.Value{"project_id": types.StringNull()},
	)

	objValue, _ := types.ObjectValue(
		map[string]attr.Type{
			"aws_config": awsConfig.Type(context.Background()),
			"gcp_config": gcpConfig.Type(context.Background()),
		},
		map[string]attr.Value{
			"aws_config": awsConfig,
			"gcp_config": gcpConfig,
		},
	)

	req := validator.ObjectRequest{
		ConfigValue: objValue,
	}
	resp := &validator.ObjectResponse{}

	v.ValidateObject(context.Background(), req, resp)

	// Should pass because empty object (all null children) is not counted as specified
	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no error (empty object not counted), got: %v", resp.Diagnostics)
	}
}

func TestAtLeastOneOfNested_ValidSingleSpecified(t *testing.T) {
	v := AtLeastOneOfNested("aws_config", "gcp_config")

	awsConfig, _ := types.ObjectValue(
		map[string]attr.Type{"account_id": types.StringType},
		map[string]attr.Value{"account_id": types.StringValue("123")},
	)
	gcpConfig := types.ObjectNull(map[string]attr.Type{"project_id": types.StringType})

	objValue, _ := types.ObjectValue(
		map[string]attr.Type{
			"aws_config": awsConfig.Type(context.Background()),
			"gcp_config": gcpConfig.Type(context.Background()),
		},
		map[string]attr.Value{
			"aws_config": awsConfig,
			"gcp_config": gcpConfig,
		},
	)

	req := validator.ObjectRequest{
		ConfigValue: objValue,
	}
	resp := &validator.ObjectResponse{}

	v.ValidateObject(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no error, got: %v", resp.Diagnostics)
	}
}

func TestAtLeastOneOfNested_ValidMultipleSpecified(t *testing.T) {
	v := AtLeastOneOfNested("aws_config", "gcp_config")

	awsConfig, _ := types.ObjectValue(
		map[string]attr.Type{"account_id": types.StringType},
		map[string]attr.Value{"account_id": types.StringValue("123")},
	)
	gcpConfig, _ := types.ObjectValue(
		map[string]attr.Type{"project_id": types.StringType},
		map[string]attr.Value{"project_id": types.StringValue("my-project")},
	)

	objValue, _ := types.ObjectValue(
		map[string]attr.Type{
			"aws_config": awsConfig.Type(context.Background()),
			"gcp_config": gcpConfig.Type(context.Background()),
		},
		map[string]attr.Value{
			"aws_config": awsConfig,
			"gcp_config": gcpConfig,
		},
	)

	req := validator.ObjectRequest{
		ConfigValue: objValue,
	}
	resp := &validator.ObjectResponse{}

	v.ValidateObject(context.Background(), req, resp)

	// Multiple specified is OK for AtLeastOneOf
	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no error for multiple specified, got: %v", resp.Diagnostics)
	}
}

func TestAtLeastOneOfNested_InvalidNoneSpecified(t *testing.T) {
	v := AtLeastOneOfNested("aws_config", "gcp_config")

	awsConfig := types.ObjectNull(map[string]attr.Type{"account_id": types.StringType})
	gcpConfig := types.ObjectNull(map[string]attr.Type{"project_id": types.StringType})

	objValue, _ := types.ObjectValue(
		map[string]attr.Type{
			"aws_config": awsConfig.Type(context.Background()),
			"gcp_config": gcpConfig.Type(context.Background()),
		},
		map[string]attr.Value{
			"aws_config": awsConfig,
			"gcp_config": gcpConfig,
		},
	)

	req := validator.ObjectRequest{
		ConfigValue: objValue,
	}
	resp := &validator.ObjectResponse{}

	v.ValidateObject(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Error("Expected error for none specified, got none")
	}
}

func TestDescription(t *testing.T) {
	exactlyOne := ExactlyOneOfNested("a", "b", "c")
	atLeastOne := AtLeastOneOfNested("x", "y")

	exactlyOneDesc := exactlyOne.Description(context.Background())
	atLeastOneDesc := atLeastOne.Description(context.Background())

	if exactlyOneDesc == "" {
		t.Error("ExactlyOneOfNested should have a description")
	}
	if atLeastOneDesc == "" {
		t.Error("AtLeastOneOfNested should have a description")
	}

	t.Logf("ExactlyOneOfNested description: %s", exactlyOneDesc)
	t.Logf("AtLeastOneOfNested description: %s", atLeastOneDesc)
}

func TestExactlyOneOfNested_UnknownValuesAreSpecified(t *testing.T) {
	// Test that objects with unknown values (e.g., from computed attributes)
	// are correctly treated as "specified" during plan phase
	v := ExactlyOneOfNested("a", "b", "c")

	// Object "a" has an unknown value (e.g., account_id = other_resource.id)
	objA, _ := types.ObjectValue(
		map[string]attr.Type{"field": types.StringType},
		map[string]attr.Value{"field": types.StringUnknown()},
	)

	// Object "b" is an empty object (not specified by user)
	objB, _ := types.ObjectValue(
		map[string]attr.Type{"field": types.StringType},
		map[string]attr.Value{"field": types.StringNull()},
	)

	// Object "c" is null
	objC := types.ObjectNull(map[string]attr.Type{"field": types.StringType})

	objValue, _ := types.ObjectValue(
		map[string]attr.Type{
			"a": objA.Type(context.Background()),
			"b": objB.Type(context.Background()),
			"c": objC.Type(context.Background()),
		},
		map[string]attr.Value{
			"a": objA,
			"b": objB,
			"c": objC,
		},
	)

	req := validator.ObjectRequest{
		ConfigValue: objValue,
	}
	resp := &validator.ObjectResponse{}

	v.ValidateObject(context.Background(), req, resp)

	// Should not have errors - "a" is specified (even with unknown values)
	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no errors, got: %v", resp.Diagnostics)
	}
}

func TestExactlyOneOfNested_EntireObjectUnknown(t *testing.T) {
	// Test that an entirely unknown object is treated as "specified"
	v := ExactlyOneOfNested("a", "b")

	// Object "a" is entirely unknown
	objA := types.ObjectUnknown(map[string]attr.Type{"field": types.StringType})

	// Object "b" is null
	objB := types.ObjectNull(map[string]attr.Type{"field": types.StringType})

	objValue, _ := types.ObjectValue(
		map[string]attr.Type{
			"a": objA.Type(context.Background()),
			"b": objB.Type(context.Background()),
		},
		map[string]attr.Value{
			"a": objA,
			"b": objB,
		},
	)

	req := validator.ObjectRequest{
		ConfigValue: objValue,
	}
	resp := &validator.ObjectResponse{}

	v.ValidateObject(context.Background(), req, resp)

	// Should not have errors - "a" is specified (entire object is unknown)
	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no errors, got: %v", resp.Diagnostics)
	}
}
