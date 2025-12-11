package schema

import (
	"testing"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestSchemaBuilder(t *testing.T) {
	// Create a builder for "project" resource
	builder := NewSchemaBuilder("project")

	// Verify resource name is stored correctly
	if builder.GetResourceName() != "project" {
		t.Errorf("Expected resource name 'project', got '%s'", builder.GetResourceName())
	}

	// Test WithOpenAPIDescription with string attribute
	attr := &resourceschema.StringAttribute{
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

	t.Run("ResourceStringAttribute", func(t *testing.T) {
		attr := &resourceschema.StringAttribute{}
		result := WithOpenAPIDescription(builder, attr, "test_field")
		if result != attr {
			t.Error("Result should be the same pointer as input")
		}
	})

	t.Run("ResourceInt64Attribute", func(t *testing.T) {
		attr := &resourceschema.Int64Attribute{}
		result := WithOpenAPIDescription(builder, attr, "test_field")
		if result != attr {
			t.Error("Result should be the same pointer as input")
		}
	})

	t.Run("ResourceBoolAttribute", func(t *testing.T) {
		attr := &resourceschema.BoolAttribute{}
		result := WithOpenAPIDescription(builder, attr, "test_field")
		if result != attr {
			t.Error("Result should be the same pointer as input")
		}
	})

	t.Run("DatasourceStringAttribute", func(t *testing.T) {
		attr := &datasourceschema.StringAttribute{}
		result := WithOpenAPIDescription(builder, attr, "test_field")
		if result != attr {
			t.Error("Result should be the same pointer as input")
		}
	})

	t.Run("DatasourceInt64Attribute", func(t *testing.T) {
		attr := &datasourceschema.Int64Attribute{}
		result := WithOpenAPIDescription(builder, attr, "test_field")
		if result != attr {
			t.Error("Result should be the same pointer as input")
		}
	})
}

func TestAddAttrWithCommonDescriptions(t *testing.T) {
	builder := NewSchemaBuilder("test")
	attrs := make(map[string]resourceschema.Attribute)

	// Test path parameter field (should come from OpenAPI, not CommonDescriptions)
	AddAttr(attrs, "organization_id", builder, &resourceschema.StringAttribute{
		Required: true,
	})

	// Verify the attribute was added and has a description
	attr := attrs["organization_id"]
	if attr == nil {
		t.Fatal("Expected organization_id attribute to be added")
	}

	// The attribute is stored as a pointer type
	strAttr, ok := attr.(*resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("Expected *resourceschema.StringAttribute, got %T", attr)
	}

	// Should get description from OpenAPI path parameters now
	if strAttr.MarkdownDescription == "" {
		t.Error("Expected organization_id to have description from OpenAPI path parameters")
	}
	if strAttr.MarkdownDescription != "The GUID4 ID of the organization." {
		t.Errorf("Expected OpenAPI path parameter description, got: %s", strAttr.MarkdownDescription)
	}

	// Test HTTP header field (should come from OpenAPI parameters)
	AddAttr(attrs, "if_match", builder, &resourceschema.StringAttribute{
		Optional: true,
	})

	ifMatchAttr, ok := attrs["if_match"].(*resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("Expected *resourceschema.StringAttribute for if_match, got %T", attrs["if_match"])
	}

	if ifMatchAttr.MarkdownDescription == "" {
		t.Error("Expected if_match to have description from OpenAPI parameters")
	}
	if ifMatchAttr.MarkdownDescription != "A precondition header that specifies the entity tag of a resource." {
		t.Errorf("Expected OpenAPI parameter description, got: %s", ifMatchAttr.MarkdownDescription)
	}
}

func TestAddAttrFallbackChain(t *testing.T) {
	builder := NewSchemaBuilder("project")
	attrs := make(map[string]resourceschema.Attribute)

	// Test field that exists in OpenAPI (should prefer OpenAPI)
	AddAttr(attrs, "name", builder, &resourceschema.StringAttribute{
		Required: true,
	})

	nameAttr, ok := attrs["name"].(*resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("Expected *resourceschema.StringAttribute, got %T", attrs["name"])
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

func TestAllDescriptionsFromOpenAPI(t *testing.T) {
	// Test that all common fields now come from OpenAPI spec (no hardcoded descriptions)
	builder := NewSchemaBuilder("test")
	attrs := make(map[string]resourceschema.Attribute)

	testFields := []struct {
		name     string
		expected string
	}{
		{"organization_id", "The GUID4 ID of the organization."},
		{"project_id", "The GUID4 ID of the project."},
		{"cluster_id", "The GUID4 ID of the cluster."},
		{"if_match", "A precondition header that specifies the entity tag of a resource."},
		{"audit", "Couchbase audit data."},
	}

	for _, tf := range testFields {
		AddAttr(attrs, tf.name, builder, &resourceschema.StringAttribute{})

		attr, ok := attrs[tf.name].(*resourceschema.StringAttribute)
		if !ok {
			t.Fatalf("Expected *resourceschema.StringAttribute for %s, got %T", tf.name, attrs[tf.name])
		}

		if attr.MarkdownDescription == "" {
			t.Errorf("Expected %s to have description from OpenAPI spec", tf.name)
		}
		if attr.MarkdownDescription != tf.expected {
			t.Errorf("Field %s: expected %q, got %q", tf.name, tf.expected, attr.MarkdownDescription)
		}
	}
}

func TestAddAttrWithDatasourceSchema(t *testing.T) {
	builder := NewSchemaBuilder("event")
	attrs := make(map[string]datasourceschema.Attribute)

	// Test datasource attributes work with AddAttr
	AddAttr(attrs, "id", builder, &datasourceschema.StringAttribute{
		Required: true,
	})
	AddAttr(attrs, "organization_id", builder, &datasourceschema.StringAttribute{
		Required: true,
	})

	// Verify attributes were added
	if attrs["id"] == nil {
		t.Fatal("Expected id attribute to be added")
	}
	if attrs["organization_id"] == nil {
		t.Fatal("Expected organization_id attribute to be added")
	}

	// Verify descriptions were set
	idAttr, ok := attrs["id"].(*datasourceschema.StringAttribute)
	if !ok {
		t.Fatalf("Expected *datasourceschema.StringAttribute for id, got %T", attrs["id"])
	}
	if idAttr.MarkdownDescription == "" {
		t.Error("Expected id to have description")
	}
}

func TestSetMarkdownDescriptionReflection(t *testing.T) {
	t.Run("ResourceStringAttribute", func(t *testing.T) {
		attr := &resourceschema.StringAttribute{}
		setMarkdownDescription(attr, "test description")
		if attr.MarkdownDescription != "test description" {
			t.Errorf("Expected 'test description', got '%s'", attr.MarkdownDescription)
		}
	})

	t.Run("DatasourceInt64Attribute", func(t *testing.T) {
		attr := &datasourceschema.Int64Attribute{}
		setMarkdownDescription(attr, "test int64 description")
		if attr.MarkdownDescription != "test int64 description" {
			t.Errorf("Expected 'test int64 description', got '%s'", attr.MarkdownDescription)
		}
	})

	t.Run("ResourceNestedAttribute", func(t *testing.T) {
		attr := &resourceschema.SingleNestedAttribute{}
		setMarkdownDescription(attr, "nested description")
		if attr.MarkdownDescription != "nested description" {
			t.Errorf("Expected 'nested description', got '%s'", attr.MarkdownDescription)
		}
	})

	t.Run("NilPointer", func(t *testing.T) {
		// Should not panic on nil
		setMarkdownDescription(nil, "test")
	})

	t.Run("NonStructType", func(t *testing.T) {
		// Should not panic on non-struct
		str := "test"
		setMarkdownDescription(&str, "test")
	})
}
