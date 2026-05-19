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

// Composition lookup tests

func TestCompositionLookup_OneOf(t *testing.T) {
	b := stubBuilder{resource: "networkPeering", schema: "CreateNetworkPeeringRequest"}

	def := CompositionLookup(b, nil, "provider_config")
	if def == nil {
		t.Fatal("expected composition for providerConfig, got nil")
	}
	if def.Kind != "oneOf" {
		t.Errorf("expected Kind=oneOf, got %q", def.Kind)
	}
	if len(def.Branches) != 3 {
		t.Errorf("expected 3 branches, got %d: %v", len(def.Branches), def.Branches)
	}
}

func TestCompositionLookup_AnyOf(t *testing.T) {
	b := stubBuilder{resource: "cmek", schema: "GetCMEKMetadata"}

	def := CompositionLookup(b, nil, "config")
	if def == nil {
		t.Fatal("expected composition for config, got nil")
	}
	if def.Kind != "anyOf" {
		t.Errorf("expected Kind=anyOf, got %q", def.Kind)
	}
}

func TestCompositionLookup_PatternMatching(t *testing.T) {
	// Should find via Create{Resource}Request pattern
	b := stubBuilder{resource: "networkPeering", schema: "networkPeering"}

	def := CompositionLookup(b, nil, "provider_config")
	if def == nil {
		t.Fatal("expected composition via CreateNetworkPeeringRequest pattern, got nil")
	}
}

func TestCompositionLookup_AlternateSchemaWins(t *testing.T) {
	b := stubBuilder{resource: "unknown", schema: "unknown"}

	def := CompositionLookup(b, []string{"CreateNetworkPeeringRequest"}, "provider_config")
	if def == nil {
		t.Fatal("expected composition via alternate schema, got nil")
	}
}

func TestCompositionLookup_NotFound(t *testing.T) {
	b := stubBuilder{resource: "unknown", schema: "unknown"}

	def := CompositionLookup(b, nil, "non_existent_field")
	if def != nil {
		t.Errorf("expected nil for unknown field, got %+v", def)
	}
}
