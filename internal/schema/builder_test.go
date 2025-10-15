package schema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestSchemaBuilder(t *testing.T) {
	// Create a builder for "project" resource
	builder := NewSchemaBuilder("project")

	// Verify resource name is stored correctly
	if builder.GetResourceName() != "project" {
		t.Errorf("Expected resource name 'project', got '%s'", builder.GetResourceName())
	}

	// Test WithOpenAPIDescription with string attribute
	attr := &schema.StringAttribute{
		Required: true,
	}

	// Use the generic function
	result := WithOpenAPIDescription(builder, attr, "name")

	// Verify description was set (should be non-empty for "name" field)
	if result.MarkdownDescription == "" {
		t.Error("Expected non-empty MarkdownDescription for 'name' field")
	}

	t.Logf("Description for project.name: %s", result.MarkdownDescription)
}

func TestSchemaBuilderMultipleResources(t *testing.T) {
	// Test that we can create builders for multiple resources and data sources
	items := []string{"project", "bucket", "cluster", "app_service", "allowlist", "users"}

	for _, itemName := range items {
		builder := NewSchemaBuilder(itemName)

		if builder.GetResourceName() != itemName {
			t.Errorf("Expected resource name '%s', got '%s'", itemName, builder.GetResourceName())
		}
	}
}

func TestSchemaBuilderWithDifferentAttributeTypes(t *testing.T) {
	builder := NewSchemaBuilder("test")

	t.Run("StringAttribute", func(t *testing.T) {
		attr := &schema.StringAttribute{}
		result := WithOpenAPIDescription(builder, attr, "test_field")
		// Type is preserved, no casting needed!
		if result != attr {
			t.Error("Result should be the same pointer as input")
		}
	})

	t.Run("Int64Attribute", func(t *testing.T) {
		attr := &schema.Int64Attribute{}
		result := WithOpenAPIDescription(builder, attr, "test_field")
		if result != attr {
			t.Error("Result should be the same pointer as input")
		}
	})

	t.Run("BoolAttribute", func(t *testing.T) {
		attr := &schema.BoolAttribute{}
		result := WithOpenAPIDescription(builder, attr, "test_field")
		if result != attr {
			t.Error("Result should be the same pointer as input")
		}
	})

	t.Run("ObjectAttribute", func(t *testing.T) {
		attr := &schema.ObjectAttribute{}
		result := WithOpenAPIDescription(builder, attr, "test_field")
		if result != attr {
			t.Error("Result should be the same pointer as input")
		}
	})

	t.Run("SingleNestedAttribute", func(t *testing.T) {
		attr := &schema.SingleNestedAttribute{}
		result := WithOpenAPIDescription(builder, attr, "test_field")
		if result != attr {
			t.Error("Result should be the same pointer as input")
		}
	})
}
