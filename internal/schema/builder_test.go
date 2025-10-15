package schema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestSchemaBuilder(t *testing.T) {
	// Create a builder for "project" resource
	builder := NewSchemaBuilder("project")

	// Verify it implements the interface
	var _ SchemaBuilder = builder

	// Verify resource name is stored correctly
	if builder.GetResourceName() != "project" {
		t.Errorf("Expected resource name 'project', got '%s'", builder.GetResourceName())
	}

	// Test WithOpenAPIDescription with string attribute
	attr := &schema.StringAttribute{
		Required: true,
	}

	result := builder.WithOpenAPIDescription(attr, "name")
	resultAttr, ok := result.(*schema.StringAttribute)
	if !ok {
		t.Fatalf("Expected *schema.StringAttribute, got %T", result)
	}

	// Verify description was set (should be non-empty for "name" field)
	if resultAttr.MarkdownDescription == "" {
		t.Error("Expected non-empty MarkdownDescription for 'name' field")
	}

	t.Logf("Description for project.name: %s", resultAttr.MarkdownDescription)
}

func TestSchemaBuilderInterface(t *testing.T) {
	// Test that we can create builders for multiple resources and data sources
	items := []string{"project", "bucket", "cluster", "app_service", "allowlist", "users"}

	for _, itemName := range items {
		builder := NewSchemaBuilder(itemName)

		if builder.GetResourceName() != itemName {
			t.Errorf("Expected resource name '%s', got '%s'", itemName, builder.GetResourceName())
		}

		// Verify it implements the interface
		var _ SchemaBuilder = builder
	}
}

func TestSchemaBuilderWithDifferentAttributeTypes(t *testing.T) {
	builder := NewSchemaBuilder("test")

	tests := []struct {
		name     string
		attr     any
		wantType string
	}{
		{
			name:     "StringAttribute",
			attr:     &schema.StringAttribute{},
			wantType: "*schema.StringAttribute",
		},
		{
			name:     "Int64Attribute",
			attr:     &schema.Int64Attribute{},
			wantType: "*schema.Int64Attribute",
		},
		{
			name:     "BoolAttribute",
			attr:     &schema.BoolAttribute{},
			wantType: "*schema.BoolAttribute",
		},
		{
			name:     "ObjectAttribute",
			attr:     &schema.ObjectAttribute{},
			wantType: "*schema.ObjectAttribute",
		},
		{
			name:     "SingleNestedAttribute",
			attr:     &schema.SingleNestedAttribute{},
			wantType: "*schema.SingleNestedAttribute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.WithOpenAPIDescription(tt.attr, "test_field")

			// Verify the type is preserved
			switch result.(type) {
			case *schema.StringAttribute:
				if tt.wantType != "*schema.StringAttribute" {
					t.Errorf("Expected %s, got *schema.StringAttribute", tt.wantType)
				}
			case *schema.Int64Attribute:
				if tt.wantType != "*schema.Int64Attribute" {
					t.Errorf("Expected %s, got *schema.Int64Attribute", tt.wantType)
				}
			case *schema.BoolAttribute:
				if tt.wantType != "*schema.BoolAttribute" {
					t.Errorf("Expected %s, got *schema.BoolAttribute", tt.wantType)
				}
			case *schema.ObjectAttribute:
				if tt.wantType != "*schema.ObjectAttribute" {
					t.Errorf("Expected %s, got *schema.ObjectAttribute", tt.wantType)
				}
			case *schema.SingleNestedAttribute:
				if tt.wantType != "*schema.SingleNestedAttribute" {
					t.Errorf("Expected %s, got *schema.SingleNestedAttribute", tt.wantType)
				}
			default:
				t.Errorf("Unexpected type: %T", result)
			}
		})
	}
}
