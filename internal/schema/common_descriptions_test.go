package schema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestAddAttrWithCommonDescriptions(t *testing.T) {
	builder := NewSchemaBuilder("test")
	attrs := make(map[string]schema.Attribute)

	// Test path parameter field (should come from OpenAPI, not CommonDescriptions)
	AddAttr(attrs, "organization_id", builder, &schema.StringAttribute{
		Required: true,
	})

	// Verify the attribute was added and has a description
	attr := attrs["organization_id"]
	if attr == nil {
		t.Fatal("Expected organization_id attribute to be added")
	}

	// The attribute is stored as a pointer type
	strAttr, ok := attr.(*schema.StringAttribute)
	if !ok {
		t.Fatalf("Expected *schema.StringAttribute, got %T", attr)
	}

	// Should get description from OpenAPI path parameters now
	if strAttr.MarkdownDescription == "" {
		t.Error("Expected organization_id to have description from OpenAPI path parameters")
	}
	if strAttr.MarkdownDescription != "The GUID4 ID of the organization." {
		t.Errorf("Expected OpenAPI path parameter description, got: %s", strAttr.MarkdownDescription)
	}

	// Test HTTP header field (should come from CommonDescriptions)
	AddAttr(attrs, "if_match", builder, &schema.StringAttribute{
		Optional: true,
	})

	ifMatchAttr, ok := attrs["if_match"].(*schema.StringAttribute)
	if !ok {
		t.Fatalf("Expected *schema.StringAttribute for if_match, got %T", attrs["if_match"])
	}

	if ifMatchAttr.MarkdownDescription != CommonDescriptions["if_match"] {
		t.Errorf("Expected if_match description from CommonDescriptions, got: %s", ifMatchAttr.MarkdownDescription)
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
	// Verify common fields are registered (only non-OpenAPI fields)
	requiredFields := []string{
		"if_match",
		"etag",
		"audit",
	}

	for _, field := range requiredFields {
		if desc, ok := CommonDescriptions[field]; !ok || desc == "" {
			t.Errorf("Expected CommonDescriptions to have non-empty entry for %s", field)
		}
	}

	// Verify ID fields are NOT in CommonDescriptions (they come from OpenAPI now)
	idFieldsShouldNotExist := []string{
		"organization_id",
		"project_id",
		"cluster_id",
	}

	for _, field := range idFieldsShouldNotExist {
		if _, ok := CommonDescriptions[field]; ok {
			t.Errorf("Expected %s to NOT be in CommonDescriptions (should come from OpenAPI parameters)", field)
		}
	}
}
