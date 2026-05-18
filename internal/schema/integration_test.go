package schema

import (
	"testing"

	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// TestCompositionValidatorIntegration tests the end-to-end flow of
// composition validator auto-attach using a real schema structure that
// matches the OpenAPI spec pattern.
func TestCompositionValidatorIntegration(t *testing.T) {
	// Simulate a schema builder for CreateNetworkPeeringRequest which has:
	// providerConfig: oneOf [AWSConfigData, GCPConfigData, AzureConfigData]
	builder := NewSchemaBuilder("networkPeering", "CreateNetworkPeeringRequest")
	attrs := make(map[string]resourceschema.Attribute)

	// Create the nested attributes for the composition branches
	awsConfigAttrs := map[string]resourceschema.Attribute{
		"account_id": &resourceschema.StringAttribute{Required: true},
		"vpc_id":     &resourceschema.StringAttribute{Required: true},
	}

	gcpConfigAttrs := map[string]resourceschema.Attribute{
		"project_id":   &resourceschema.StringAttribute{Required: true},
		"network_name": &resourceschema.StringAttribute{Required: true},
	}

	azureConfigAttrs := map[string]resourceschema.Attribute{
		"subscription_id": &resourceschema.StringAttribute{Required: true},
		"vnet_name":       &resourceschema.StringAttribute{Required: true},
	}

	// Create the provider_config nested attribute with child branches
	providerConfigAttrs := map[string]resourceschema.Attribute{
		"aws_config": &resourceschema.SingleNestedAttribute{
			Optional:   true,
			Attributes: awsConfigAttrs,
		},
		"gcp_config": &resourceschema.SingleNestedAttribute{
			Optional:   true,
			Attributes: gcpConfigAttrs,
		},
		"azure_config": &resourceschema.SingleNestedAttribute{
			Optional:   true,
			Attributes: azureConfigAttrs,
		},
	}

	// AddAttr should auto-attach ExactlyOneOf validator because:
	// 1. CompositionLookup finds CreateNetworkPeeringRequest.providerConfig (oneOf)
	// 2. appendCompositionValidator introspects providerConfigAttrs
	// 3. Finds 3 optional SingleNestedAttribute children
	// 4. Attaches ExactlyOneOf validator
	AddAttr(attrs, "provider_config", builder, &resourceschema.SingleNestedAttribute{
		Required:   true,
		Attributes: providerConfigAttrs,
	})

	// Verify the attribute was added
	if attrs["provider_config"] == nil {
		t.Fatal("Expected provider_config attribute to be added")
	}

	// Verify it's a SingleNestedAttribute
	providerConfig, ok := attrs["provider_config"].(*resourceschema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("Expected *resourceschema.SingleNestedAttribute, got %T", attrs["provider_config"])
	}

	// Verify the composition validator was auto-attached
	if len(providerConfig.Validators) != 1 {
		t.Fatalf("Expected 1 validator (ExactlyOneOf), got %d", len(providerConfig.Validators))
	}

	t.Logf("✓ Composition validator successfully auto-attached to provider_config")
}

// TestCompositionValidatorIntegration_AnyOf tests anyOf composition auto-attach.
func TestCompositionValidatorIntegration_AnyOf(t *testing.T) {
	// Simulate GetCMEKMetadata which has: config: anyOf [AWSConfig, GCPConfig, AzureConfig]
	builder := NewSchemaBuilder("cmek", "GetCMEKMetadata")
	attrs := make(map[string]resourceschema.Attribute)

	configAttrs := map[string]resourceschema.Attribute{
		"aws_config": &resourceschema.SingleNestedAttribute{
			Optional:   true,
			Attributes: map[string]resourceschema.Attribute{},
		},
		"gcp_config": &resourceschema.SingleNestedAttribute{
			Optional:   true,
			Attributes: map[string]resourceschema.Attribute{},
		},
		"azure_config": &resourceschema.SingleNestedAttribute{
			Optional:   true,
			Attributes: map[string]resourceschema.Attribute{},
		},
	}

	// AddAttr should auto-attach AtLeastOneOf validator for anyOf
	AddAttr(attrs, "config", builder, &resourceschema.SingleNestedAttribute{
		Optional:   true,
		Attributes: configAttrs,
	})

	config, ok := attrs["config"].(*resourceschema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("Expected *resourceschema.SingleNestedAttribute, got %T", attrs["config"])
	}

	// Verify AtLeastOneOf validator was attached
	if len(config.Validators) != 1 {
		t.Fatalf("Expected 1 validator (AtLeastOneOf for anyOf), got %d", len(config.Validators))
	}

	t.Logf("✓ AtLeastOneOf validator successfully auto-attached for anyOf composition")
}

// TestCompositionValidatorIntegration_Override tests that call-site validators
// take precedence over auto-attached validators.
func TestCompositionValidatorIntegration_Override(t *testing.T) {
	builder := NewSchemaBuilder("test", "CreateNetworkPeeringRequest")
	attrs := make(map[string]resourceschema.Attribute)

	providerConfigAttrs := map[string]resourceschema.Attribute{
		"aws_config": &resourceschema.SingleNestedAttribute{
			Optional:   true,
			Attributes: map[string]resourceschema.Attribute{},
		},
		"gcp_config": &resourceschema.SingleNestedAttribute{
			Optional:   true,
			Attributes: map[string]resourceschema.Attribute{},
		},
	}

	// Pre-attach a custom validator (simulating manual override)
	customAttr := &resourceschema.SingleNestedAttribute{
		Required: true,
		Validators: []validator.Object{
			mockObjectValidator{},
		},
		Attributes: providerConfigAttrs,
	}

	// AddAttr should NOT add another validator (call-site wins)
	AddAttr(attrs, "provider_config", builder, customAttr)

	providerConfig, ok := attrs["provider_config"].(*resourceschema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("Expected *resourceschema.SingleNestedAttribute, got %T", attrs["provider_config"])
	}

	// Should still have only 1 validator (the original one, not auto-attached)
	if len(providerConfig.Validators) != 1 {
		t.Fatalf("Expected 1 validator (original), got %d", len(providerConfig.Validators))
	}

	t.Logf("✓ Call-site validator override respected (no auto-attach)")
}
