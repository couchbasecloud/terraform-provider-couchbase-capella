package enums

import "testing"

type stubBuilder struct {
	resource string
	schema   string
}

func (s stubBuilder) GetResourceName() string      { return s.resource }
func (s stubBuilder) GetOpenAPISchemaName() string { return s.schema }

func TestLookup_LowerCamelSchemaName(t *testing.T) {
	// NewSchemaBuilder("allowlist", "allowedCidr"): enumTable key is
	// PascalCase "AllowedCidr", so Lookup must capitalize the schema name.
	b := stubBuilder{resource: "allowlist", schema: "allowedCidr"}

	def := Lookup(b, nil, "type")
	if def == nil {
		t.Fatal("expected enum for allowlist.type via AllowedCidr, got nil")
	}
	if def.Type != "string" {
		t.Errorf("expected string enum, got %q", def.Type)
	}
}

func TestLookup_PascalCaseSchemaName(t *testing.T) {
	// NewSchemaBuilder("sampleBucket", "PostSampleBucket") — schema already
	// PascalCase; capitalize is a no-op and Lookup still resolves it.
	b := stubBuilder{resource: "sampleBucket", schema: "PostSampleBucket"}

	def := Lookup(b, nil, "name")
	if def == nil {
		t.Fatal("expected enum for sampleBucket.name, got nil")
	}
}

func TestLookup_ResourceNameOnly(t *testing.T) {
	// NewSchemaBuilder("project") — resource == schema. Patterns derived
	// from the resource name (Project, CreateProjectRequest, ...) cover it.
	b := stubBuilder{resource: "project", schema: "project"}

	if def := Lookup(b, nil, "non_existent_field"); def != nil {
		t.Errorf("expected nil for unknown field, got %+v", def)
	}
}

func TestLookup_AlternateSchemaWins(t *testing.T) {
	b := stubBuilder{resource: "unknown", schema: "unknown"}

	def := Lookup(b, []string{"AllowedCidr"}, "status")
	if def == nil {
		t.Fatal("expected enum via alternate schema AllowedCidr, got nil")
	}
}
