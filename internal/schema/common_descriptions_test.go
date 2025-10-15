package schema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestAddAttrWithCommonDescriptions(t *testing.T) {
	builder := NewSchemaBuilder("test")
	attrs := make(map[string]schema.Attribute)

	// Test common field that's not in OpenAPI
	AddAttr(attrs, "organization_id", builder, &schema.StringAttribute{
		Required: true,
	})

	// Verify the attribute was added and has the expected description
	attr := attrs["organization_id"]
	if attr == nil {
		t.Fatal("Expected organization_id attribute to be added")
	}

	// The attribute is stored as a pointer type
	strAttr, ok := attr.(*schema.StringAttribute)
	if !ok {
		t.Fatalf("Expected *schema.StringAttribute, got %T", attr)
	}

	if strAttr.MarkdownDescription != CommonDescriptions["organization_id"] {
		t.Errorf("Expected organization_id description from CommonDescriptions, got: %s", strAttr.MarkdownDescription)
	}
}

func TestAddAttrFallbackChain(t *testing.T) {
	builder := NewSchemaBuilder("project")
	attrs := make(map[string]schema.Attribute)

	// Test field that exists in OpenAPI (should prefer OpenAPI)
	AddAttr(attrs, "name", builder, &schema.StringAttribute{
		Required: true,
	})

	nameAttr, ok := attrs["name"].(*schema.StringAttribute)
	if !ok {
		t.Fatalf("Expected *schema.StringAttribute, got %T", attrs["name"])
	}

	// Should have OpenAPI description (not the common one)
	if nameAttr.MarkdownDescription == "" {
		t.Error("Expected name to have OpenAPI description")
	}
	// OpenAPI should have constraints section
	if len(nameAttr.MarkdownDescription) < 50 {
		t.Errorf("Expected rich OpenAPI description with constraints, got: %s", nameAttr.MarkdownDescription)
	}
}

func TestCommonDescriptionsRegistry(t *testing.T) {
	// Verify common fields are registered
	requiredFields := []string{
		"organization_id",
		"project_id",
		"cluster_id",
		"if_match",
		"etag",
		"audit",
	}

	for _, field := range requiredFields {
		if desc, ok := CommonDescriptions[field]; !ok || desc == "" {
			t.Errorf("Expected CommonDescriptions to have non-empty entry for %s", field)
		}
	}
}
