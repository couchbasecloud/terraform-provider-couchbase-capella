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
