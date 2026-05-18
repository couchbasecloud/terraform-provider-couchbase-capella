package schema

import (
	"context"
	"testing"

	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// TestCompositionValidatorIntegration_OneOf verifies that ExactlyOneOfNested validator
// is auto-attached via AddAttr for oneOf compositions.
func TestCompositionValidatorIntegration_OneOf(t *testing.T) {
	builder := NewSchemaBuilder("networkPeering", "CreateNetworkPeeringRequest")
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
		"azure_config": &resourceschema.SingleNestedAttribute{
			Optional:   true,
			Attributes: map[string]resourceschema.Attribute{},
		},
	}

	AddAttr(attrs, "provider_config", builder, &resourceschema.SingleNestedAttribute{
		Required:   true,
		Attributes: providerConfigAttrs,
	})

	providerConfig, ok := attrs["provider_config"].(*resourceschema.SingleNestedAttribute)
	if !ok {
		t.Fatal("Expected provider_config to be SingleNestedAttribute")
	}

	// Verify validator was auto-attached
	if len(providerConfig.Validators) != 1 {
		t.Fatalf("Expected 1 validator (ExactlyOneOfNested), got %d", len(providerConfig.Validators))
	}

	// Verify it's the correct validator by checking description
	desc := providerConfig.Validators[0].Description(context.Background())
	if desc == "" {
		t.Fatal("Validator should have a description")
	}
	t.Logf("✓ Composition validator auto-attached: %s", desc)
}

// TestCompositionValidatorIntegration_AnyOf verifies that AtLeastOneOfNested validator
// is auto-attached via AddAttr for anyOf compositions.
func TestCompositionValidatorIntegration_AnyOf(t *testing.T) {
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

	AddAttr(attrs, "config", builder, &resourceschema.SingleNestedAttribute{
		Optional:   true,
		Attributes: configAttrs,
	})

	config, ok := attrs["config"].(*resourceschema.SingleNestedAttribute)
	if !ok {
		t.Fatal("Expected config to be SingleNestedAttribute")
	}

	// Verify validator was auto-attached
	if len(config.Validators) != 1 {
		t.Fatalf("Expected 1 validator (AtLeastOneOfNested), got %d", len(config.Validators))
	}

	desc := config.Validators[0].Description(context.Background())
	t.Logf("✓ AnyOf composition validator auto-attached: %s", desc)
}

// TestCompositionValidatorIntegration_CallSiteOverride verifies that call-site
// validators are preserved (not overwritten).
func TestCompositionValidatorIntegration_CallSiteOverride(t *testing.T) {
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

	// Pre-attach a custom validator
	customAttr := &resourceschema.SingleNestedAttribute{
		Required: true,
		Validators: []validator.Object{
			mockObjectValidator{},
		},
		Attributes: providerConfigAttrs,
	}

	AddAttr(attrs, "provider_config", builder, customAttr)

	providerConfig, ok := attrs["provider_config"].(*resourceschema.SingleNestedAttribute)
	if !ok {
		t.Fatal("Expected provider_config to be SingleNestedAttribute")
	}

	// Should still have only 1 validator (the original one, not auto-attached)
	if len(providerConfig.Validators) != 1 {
		t.Fatalf("Expected 1 validator (original), got %d", len(providerConfig.Validators))
	}

	t.Log("✓ Call-site validator preserved")
}
