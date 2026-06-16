package schema

import (
	"context"
	"reflect"
	"testing"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/enums"
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

// Composition validator tests

func TestAppendCompositionValidator_OneOf(t *testing.T) {
	def := &enums.CompositionDef{
		Kind:     "oneOf",
		Branches: []string{"AWSConfig", "GCPConfig", "AzureConfig"},
	}

	attr := &resourceschema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]resourceschema.Attribute{
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
		},
	}

	appendCompositionValidator(attr, def)

	if len(attr.Validators) != 1 {
		t.Fatalf("Expected 1 validator, got %d", len(attr.Validators))
	}
}

func TestAppendCompositionValidator_AnyOf(t *testing.T) {
	def := &enums.CompositionDef{
		Kind:     "anyOf",
		Branches: []string{"AWS", "GCP"},
	}

	attr := &resourceschema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]resourceschema.Attribute{
			"aws": &resourceschema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]resourceschema.Attribute{},
			},
			"gcp": &resourceschema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]resourceschema.Attribute{},
			},
		},
	}

	appendCompositionValidator(attr, def)

	if len(attr.Validators) != 1 {
		t.Fatalf("Expected 1 validator, got %d", len(attr.Validators))
	}
}

func TestAppendCompositionValidator_SkipsWithExistingValidators(t *testing.T) {
	def := &enums.CompositionDef{
		Kind:     "oneOf",
		Branches: []string{"A", "B"},
	}

	// Pre-populate with a validator (simulating call-site override)
	attr := &resourceschema.SingleNestedAttribute{
		Optional: true,
		Validators: []validator.Object{
			mockObjectValidator{},
		},
		Attributes: map[string]resourceschema.Attribute{
			"a": &resourceschema.SingleNestedAttribute{Optional: true},
			"b": &resourceschema.SingleNestedAttribute{Optional: true},
		},
	}

	appendCompositionValidator(attr, def)

	// Should still have only 1 validator (the original one)
	if len(attr.Validators) != 1 {
		t.Fatalf("Expected 1 validator (existing), got %d", len(attr.Validators))
	}
}

func TestAppendCompositionValidator_SkipsWithLessThan2Children(t *testing.T) {
	def := &enums.CompositionDef{
		Kind:     "oneOf",
		Branches: []string{"OnlyOne"},
	}

	attr := &resourceschema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]resourceschema.Attribute{
			"only_one": &resourceschema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]resourceschema.Attribute{},
			},
		},
	}

	appendCompositionValidator(attr, def)

	// Should not add validator when < 2 optional nested children
	if len(attr.Validators) != 0 {
		t.Fatalf("Expected 0 validators for single child, got %d", len(attr.Validators))
	}
}

func TestAppendCompositionValidator_IgnoresRequiredChildren(t *testing.T) {
	def := &enums.CompositionDef{
		Kind:     "oneOf",
		Branches: []string{"A", "B"},
	}

	attr := &resourceschema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]resourceschema.Attribute{
			"a": &resourceschema.SingleNestedAttribute{
				Required:   true, // Required, not a composition branch
				Attributes: map[string]resourceschema.Attribute{},
			},
			"b": &resourceschema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]resourceschema.Attribute{},
			},
		},
	}

	appendCompositionValidator(attr, def)

	// Should not add validator when only 1 optional nested child
	if len(attr.Validators) != 0 {
		t.Fatalf("Expected 0 validators (only 1 optional child), got %d", len(attr.Validators))
	}
}

func TestAppendCompositionValidator_IgnoresNonNestedAttributes(t *testing.T) {
	def := &enums.CompositionDef{
		Kind:     "oneOf",
		Branches: []string{"A", "B"},
	}

	attr := &resourceschema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]resourceschema.Attribute{
			"a": &resourceschema.StringAttribute{Optional: true}, // String, not nested
			"b": &resourceschema.StringAttribute{Optional: true}, // String, not nested
		},
	}

	appendCompositionValidator(attr, def)

	// Should not add validator when no optional nested children
	if len(attr.Validators) != 0 {
		t.Fatalf("Expected 0 validators (no nested children), got %d", len(attr.Validators))
	}
}

func TestAppendCompositionValidator_Datasource(t *testing.T) {
	def := &enums.CompositionDef{
		Kind:     "oneOf",
		Branches: []string{"AWS", "GCP"},
	}

	attr := &datasourceschema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]datasourceschema.Attribute{
			"aws": &datasourceschema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]datasourceschema.Attribute{},
			},
			"gcp": &datasourceschema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]datasourceschema.Attribute{},
			},
		},
	}

	appendCompositionValidator(attr, def)

	if len(attr.Validators) != 1 {
		t.Fatalf("Expected 1 validator for datasource, got %d", len(attr.Validators))
	}
}

func TestAppendCompositionValidator_IgnoresComputedOnlyChildren(t *testing.T) {
	def := &enums.CompositionDef{
		Kind:     "oneOf",
		Branches: []string{"A", "B", "C"},
	}

	attr := &resourceschema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]resourceschema.Attribute{
			"a": &resourceschema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]resourceschema.Attribute{},
			},
			"b": &resourceschema.SingleNestedAttribute{
				Computed:   true, // Computed-only, not a valid composition branch
				Attributes: map[string]resourceschema.Attribute{},
			},
			"c": &resourceschema.SingleNestedAttribute{
				Computed:   true, // Computed-only, not a valid composition branch
				Attributes: map[string]resourceschema.Attribute{},
			},
		},
	}

	appendCompositionValidator(attr, def)

	// Should not add validator when only 1 optional nested child (computed-only excluded)
	if len(attr.Validators) != 0 {
		t.Fatalf("Expected 0 validators (only 1 optional child after excluding computed-only), got %d", len(attr.Validators))
	}
}

func TestAppendCompositionValidator_AllOfSkipped(t *testing.T) {
	// allOf is for schema composition/inheritance, not mutual exclusion
	// We don't attach validators for allOf
	def := &enums.CompositionDef{
		Kind:     "allOf",
		Branches: []string{"Base", "Extended"},
	}

	attr := &resourceschema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]resourceschema.Attribute{
			"base": &resourceschema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]resourceschema.Attribute{},
			},
			"extended": &resourceschema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]resourceschema.Attribute{},
			},
		},
	}

	appendCompositionValidator(attr, def)

	// allOf should not add any validator
	if len(attr.Validators) != 0 {
		t.Fatalf("Expected 0 validators for allOf, got %d", len(attr.Validators))
	}
}

func TestAppendCompositionValidator_WorksWithComputedInnerFields(t *testing.T) {
	// Custom ExactlyOneOfNested/AtLeastOneOfNested validators correctly handle
	// nested attributes with computed fields by checking for actual user-provided values.
	def := &enums.CompositionDef{
		Kind:     "oneOf",
		Branches: []string{"AWS", "GCP", "Azure"},
	}

	attr := &resourceschema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]resourceschema.Attribute{
			"aws_config": &resourceschema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]resourceschema.Attribute{
					"account_id":  &resourceschema.StringAttribute{Optional: true},
					"provider_id": &resourceschema.StringAttribute{Computed: true}, // Computed field inside
				},
			},
			"gcp_config": &resourceschema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]resourceschema.Attribute{
					"project_id":  &resourceschema.StringAttribute{Optional: true},
					"provider_id": &resourceschema.StringAttribute{Computed: true}, // Computed field inside
				},
			},
			"azure_config": &resourceschema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]resourceschema.Attribute{
					"tenant_id":   &resourceschema.StringAttribute{Optional: true},
					"provider_id": &resourceschema.StringAttribute{Computed: true}, // Computed field inside
				},
			},
		},
	}

	appendCompositionValidator(attr, def)

	// Custom validator SHOULD be added - it handles computed fields correctly
	if len(attr.Validators) != 1 {
		t.Fatalf("Expected 1 validator (custom validator handles computed fields), got %d", len(attr.Validators))
	}
}

func TestExtractChildNames_IncludesAllOptionalNested(t *testing.T) {
	attrs := map[string]resourceschema.Attribute{
		"pure_optional": &resourceschema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]resourceschema.Attribute{
				"name": &resourceschema.StringAttribute{Optional: true},
			},
		},
		"has_computed_inside": &resourceschema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]resourceschema.Attribute{
				"name":        &resourceschema.StringAttribute{Optional: true},
				"provider_id": &resourceschema.StringAttribute{Computed: true},
			},
		},
		"required_field": &resourceschema.SingleNestedAttribute{
			Required:   true,
			Attributes: map[string]resourceschema.Attribute{},
		},
	}

	names := extractChildNames(attrs)

	// Should have 2 names: pure_optional and has_computed_inside
	// (required_field excluded, but computed inner fields are OK now)
	if len(names) != 2 {
		t.Fatalf("Expected 2 names (all optional nested), got %d", len(names))
	}
}

// mockObjectValidator implements validator.Object for testing
type mockObjectValidator struct{}

func (m mockObjectValidator) Description(_ context.Context) string         { return "mock" }
func (m mockObjectValidator) MarkdownDescription(_ context.Context) string { return "mock" }
func (m mockObjectValidator) ValidateObject(_ context.Context, _ validator.ObjectRequest, _ *validator.ObjectResponse) {
}

func TestSetRequiredIfUnset_SetsRequired(t *testing.T) {
	// Attribute with no flags set - should set Required: true
	attr := &resourceschema.StringAttribute{}

	setRequiredIfUnset(attr)

	if !attr.Required {
		t.Error("Expected Required to be set to true")
	}
}

func TestSetRequiredIfUnset_RespectsExistingRequired(t *testing.T) {
	// Attribute already has Required: true - should not change
	attr := &resourceschema.StringAttribute{Required: true}

	setRequiredIfUnset(attr)

	if !attr.Required {
		t.Error("Expected Required to remain true")
	}
}

func TestSetRequiredIfUnset_RespectsExistingOptional(t *testing.T) {
	// Attribute has Optional: true - should not set Required
	attr := &resourceschema.StringAttribute{Optional: true}

	setRequiredIfUnset(attr)

	if attr.Required {
		t.Error("Expected Required to remain false when Optional is set")
	}
	if !attr.Optional {
		t.Error("Expected Optional to remain true")
	}
}

func TestSetRequiredIfUnset_RespectsExistingComputed(t *testing.T) {
	// Attribute has Computed: true - should not set Required
	attr := &resourceschema.StringAttribute{Computed: true}

	setRequiredIfUnset(attr)

	if attr.Required {
		t.Error("Expected Required to remain false when Computed is set")
	}
	if !attr.Computed {
		t.Error("Expected Computed to remain true")
	}
}

func TestSetRequiredIfUnset_WorksWithDifferentAttributeTypes(t *testing.T) {
	tests := []struct {
		name string
		attr any
	}{
		{"Int64Attribute", &resourceschema.Int64Attribute{}},
		{"BoolAttribute", &resourceschema.BoolAttribute{}},
		{"SingleNestedAttribute", &resourceschema.SingleNestedAttribute{Attributes: map[string]resourceschema.Attribute{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setRequiredIfUnset(tt.attr)

			// Use reflection to check Required field
			v := reflect.ValueOf(tt.attr).Elem()
			required := v.FieldByName("Required")
			if !required.Bool() {
				t.Errorf("Expected Required to be set to true for %s", tt.name)
			}
		})
	}
}

func f64(v float64) *float64 { return &v }
func i64(v int64) *int64     { return &v }

func TestAppendConstraintValidator_StringLength(t *testing.T) {
	tests := []struct {
		name string
		def  *enums.ConstraintDef
	}{
		{"min only", &enums.ConstraintDef{MinLength: i64(2)}},
		{"max only", &enums.ConstraintDef{MaxLength: i64(64)}},
		{"min and max", &enums.ConstraintDef{MinLength: i64(2), MaxLength: i64(64)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := &resourceschema.StringAttribute{}
			appendConstraintValidator(attr, tt.def)
			if len(attr.Validators) != 1 {
				t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
			}
		})
	}
}

func TestAppendConstraintValidator_Int64Range(t *testing.T) {
	tests := []struct {
		name string
		def  *enums.ConstraintDef
	}{
		{"min only", &enums.ConstraintDef{Minimum: f64(1)}},
		{"max only", &enums.ConstraintDef{Maximum: f64(32)}},
		{"min and max", &enums.ConstraintDef{Minimum: f64(1), Maximum: f64(32)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := &resourceschema.Int64Attribute{}
			appendConstraintValidator(attr, tt.def)
			if len(attr.Validators) != 1 {
				t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
			}
		})
	}
}

func TestAppendConstraintValidator_Float64Range(t *testing.T) {
	attr := &resourceschema.Float64Attribute{}
	appendConstraintValidator(attr, &enums.ConstraintDef{Minimum: f64(0.5), Maximum: f64(99.5)})
	if len(attr.Validators) != 1 {
		t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
	}
}

func TestAppendConstraintValidator_ListSize(t *testing.T) {
	attr := &resourceschema.ListAttribute{}
	appendConstraintValidator(attr, &enums.ConstraintDef{MinItems: i64(1), MaxItems: i64(10)})
	if len(attr.Validators) != 1 {
		t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
	}
}

func TestAppendConstraintValidator_SetSize(t *testing.T) {
	attr := &resourceschema.SetAttribute{}
	appendConstraintValidator(attr, &enums.ConstraintDef{MinItems: i64(1)})
	if len(attr.Validators) != 1 {
		t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
	}
}

func TestAppendConstraintValidator_MapSize(t *testing.T) {
	attr := &resourceschema.MapAttribute{}
	appendConstraintValidator(attr, &enums.ConstraintDef{MaxItems: i64(5)})
	if len(attr.Validators) != 1 {
		t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
	}
}

func TestAppendConstraintValidator_ListNestedSize(t *testing.T) {
	attr := &resourceschema.ListNestedAttribute{}
	appendConstraintValidator(attr, &enums.ConstraintDef{MinItems: i64(1), MaxItems: i64(3)})
	if len(attr.Validators) != 1 {
		t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
	}
}

func TestAppendConstraintValidator_SetNestedSize(t *testing.T) {
	attr := &resourceschema.SetNestedAttribute{}
	appendConstraintValidator(attr, &enums.ConstraintDef{MaxItems: i64(8)})
	if len(attr.Validators) != 1 {
		t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
	}
}

func TestAppendConstraintValidator_SkipsWithExistingValidators(t *testing.T) {
	attr := &resourceschema.StringAttribute{
		Validators: []validator.String{mockStringValidator{}},
	}
	appendConstraintValidator(attr, &enums.ConstraintDef{MinLength: i64(2), MaxLength: i64(64)})
	if len(attr.Validators) != 1 {
		t.Fatalf("expected validator count to remain 1 (call-site wins), got %d", len(attr.Validators))
	}
	if _, ok := attr.Validators[0].(mockStringValidator); !ok {
		t.Error("expected pre-existing mock validator to be retained")
	}
}

func TestAppendConstraintValidator_NoApplicableConstraints(t *testing.T) {
	// String attribute given numeric constraints — nothing applicable, nothing attached.
	attr := &resourceschema.StringAttribute{}
	appendConstraintValidator(attr, &enums.ConstraintDef{Minimum: f64(1), Maximum: f64(10)})
	if len(attr.Validators) != 0 {
		t.Fatalf("expected 0 validators when constraints don't apply, got %d", len(attr.Validators))
	}
}

func TestAppendConstraintValidator_Datasource(t *testing.T) {
	// Cover datasource variants for the four core attribute kinds.
	t.Run("StringAttribute", func(t *testing.T) {
		attr := &datasourceschema.StringAttribute{}
		appendConstraintValidator(attr, &enums.ConstraintDef{MaxLength: i64(128)})
		if len(attr.Validators) != 1 {
			t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
		}
	})
	t.Run("Int64Attribute", func(t *testing.T) {
		attr := &datasourceschema.Int64Attribute{}
		appendConstraintValidator(attr, &enums.ConstraintDef{Minimum: f64(1)})
		if len(attr.Validators) != 1 {
			t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
		}
	})
	t.Run("ListAttribute", func(t *testing.T) {
		attr := &datasourceschema.ListAttribute{}
		appendConstraintValidator(attr, &enums.ConstraintDef{MinItems: i64(1)})
		if len(attr.Validators) != 1 {
			t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
		}
	})
	t.Run("SetAttribute", func(t *testing.T) {
		attr := &datasourceschema.SetAttribute{}
		appendConstraintValidator(attr, &enums.ConstraintDef{MaxItems: i64(10)})
		if len(attr.Validators) != 1 {
			t.Fatalf("expected 1 validator, got %d", len(attr.Validators))
		}
	})
}

// mockStringValidator implements validator.String for override-discipline tests.
type mockStringValidator struct{}

func (m mockStringValidator) Description(_ context.Context) string         { return "mock" }
func (m mockStringValidator) MarkdownDescription(_ context.Context) string { return "mock" }
func (m mockStringValidator) ValidateString(_ context.Context, _ validator.StringRequest, _ *validator.StringResponse) {
}
