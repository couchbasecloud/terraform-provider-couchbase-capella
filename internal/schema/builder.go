package schema

import (
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/docs"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/enums"
)

// SchemaAttribute is a type constraint for supported attribute types across resources and datasources
type SchemaAttribute interface {
	*resourceschema.StringAttribute | *resourceschema.Int64Attribute | *resourceschema.BoolAttribute |
		*resourceschema.SetAttribute | *resourceschema.Float64Attribute | *resourceschema.NumberAttribute |
		*resourceschema.ListAttribute | *resourceschema.MapAttribute | *resourceschema.SingleNestedAttribute |
		*resourceschema.ObjectAttribute | *resourceschema.SetNestedAttribute | *resourceschema.ListNestedAttribute |
		*resourceschema.MapNestedAttribute | *datasourceschema.StringAttribute | *datasourceschema.Int64Attribute |
		*datasourceschema.BoolAttribute | *datasourceschema.SetAttribute | *datasourceschema.Float64Attribute |
		*datasourceschema.NumberAttribute | *datasourceschema.ListAttribute | *datasourceschema.MapAttribute |
		*datasourceschema.SingleNestedAttribute | *datasourceschema.ObjectAttribute | *datasourceschema.SetNestedAttribute |
		*datasourceschema.ListNestedAttribute | *datasourceschema.MapNestedAttribute
}

// SchemaAttributeMap is a type constraint for attribute maps
type SchemaAttributeMap interface {
	map[string]resourceschema.Attribute | map[string]datasourceschema.Attribute
}

// SchemaBuilder provides methods for building resource and data source schemas with OpenAPI integration.
type SchemaBuilder struct {
	resourceName      string
	openAPISchemaName string // Optional: OpenAPI schema name if different from resourceName
}

// NewSchemaBuilder creates a new SchemaBuilder for a specific resource or data source.
// The second parameter (openAPISchemaName) is optional and only needed when the OpenAPI schema name
// differs from the Terraform resource name.
//
// Examples:
//   - NewSchemaBuilder("project")                      // OpenAPI schema is also "project"
//   - NewSchemaBuilder("allowlist", "allowedCidr")     // OpenAPI schema is "allowedCidr"
func NewSchemaBuilder(resourceName string, openAPISchemaName ...string) *SchemaBuilder {
	schemaName := resourceName // Default to same name
	if len(openAPISchemaName) > 0 && openAPISchemaName[0] != "" {
		schemaName = openAPISchemaName[0]
	}

	return &SchemaBuilder{
		resourceName:      resourceName,
		openAPISchemaName: schemaName,
	}
}

// GetResourceName returns the Terraform resource name
func (b *SchemaBuilder) GetResourceName() string {
	return b.resourceName
}

// GetOpenAPISchemaName returns the OpenAPI schema name
func (b *SchemaBuilder) GetOpenAPISchemaName() string {
	return b.openAPISchemaName
}

// WithOpenAPIDescription sets the MarkdownDescription for the provided attribute
// by looking up the field description from the OpenAPI specification.
// If alternateSchemas are provided, it will try those schemas first before falling back to the builder's schema.
// Example: WithOpenAPIDescription(apiKeyBuilder, attr, "secret", "RotateAPIKeyRequest")
func WithOpenAPIDescription[T SchemaAttribute](b *SchemaBuilder, attr T, fieldName string, alternateSchemas ...string) T {
	var description string

	// Try alternate schemas first (e.g., RotateAPIKeyRequest, CreateAPIKeyResponse)
	for _, altSchema := range alternateSchemas {
		description = docs.GetOpenAPIDescription(altSchema, fieldName)
		if description != "" {
			break
		}
	}

	// Fall back to the builder's default schema
	if description == "" {
		description = docs.GetOpenAPIDescription(b.openAPISchemaName, fieldName)
	}

	setMarkdownDescription(attr, description)
	return attr
}

// AddAttr adds an attribute with automatic description to the attributes map.
// Works for both resource and datasource schemas.
//
// Description is automatically loaded from the OpenAPI spec:
// 1. Path parameters (organization_id, project_id, etc.) from components.parameters
// 2. Header parameters (if_match) from components.parameters
// 3. Response headers (etag) from components.headers
// 4. Schema references (audit) from components.schemas
// 5. Schema properties (name, description, etc.) from request/response schemas
//
// If alternateSchemas are provided, it will try those schemas first before falling back to the builder's schema.
//
// Example:
//
//	attrs := make(map[string]schema.Attribute)
//	capellaschema.AddAttr(attrs, "name", projectBuilder, stringAttribute([]string{required}))
//	capellaschema.AddAttr(attrs, "secret", apiKeyBuilder, stringAttribute([]string{optional}), "RotateAPIKeyRequest")
func AddAttr[M SchemaAttributeMap, T SchemaAttribute](
	attrs M,
	fieldName string,
	builder *SchemaBuilder,
	attr T,
	alternateSchemas ...string,
) {
	var description string

	// Try alternate schemas first (e.g., RotateAPIKeyRequest, CreateAPIKeyResponse)
	for _, altSchema := range alternateSchemas {
		description = docs.GetOpenAPIDescription(altSchema, fieldName)
		if description != "" {
			break
		}
	}

	// Fall back to the builder's default schema
	if description == "" {
		description = docs.GetOpenAPIDescription(builder.openAPISchemaName, fieldName)
	}

	// If still not found and resourceName differs from openAPISchemaName, try resourceName
	// This allows pattern-based lookups like Create{Resource}Request to work
	if description == "" && builder.resourceName != builder.openAPISchemaName {
		description = docs.GetOpenAPIDescription(builder.resourceName, fieldName)
	}

	setMarkdownDescription(attr, description)

	if def := enums.Lookup(builder, alternateSchemas, fieldName); def != nil {
		appendOneOfValidator(attr, def)
	}

	if compDef := enums.CompositionLookup(builder, alternateSchemas, fieldName); compDef != nil {
		appendCompositionValidator(attr, compDef)
	}

	// Add to map based on map type
	switch m := any(&attrs).(type) {
	case *map[string]resourceschema.Attribute:
		result, ok := any(attr).(resourceschema.Attribute)
		if !ok {
			panic("failed to convert attribute to resourceschema.Attribute")
		}
		(*m)[fieldName] = result
	case *map[string]datasourceschema.Attribute:
		result, ok := any(attr).(datasourceschema.Attribute)
		if !ok {
			panic("failed to convert attribute to datasourceschema.Attribute")
		}
		(*m)[fieldName] = result
	default:
		panic("unsupported attribute map type")
	}
}

// setMarkdownDescription uses reflection to set the MarkdownDescription field on any attribute type
func setMarkdownDescription(attr any, description string) {
	v := reflect.ValueOf(attr)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}
	field := v.FieldByName("MarkdownDescription")
	if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString(description)
	}
}

// appendOneOfValidator attaches a OneOf validator derived from def to the
// attribute. Skips when the call site has already attached any validator
// — that's the override discipline: a hand-coded OneOf (or any other
// validator) at the call site wins. Only the four shapes the spec
// produces are handled: scalar string, scalar int64, list/set of
// strings, list/set of ints. Nested attributes get validators on their
// inner fields via separate AddAttr calls.
func appendOneOfValidator(a any, def *enums.EnumDef) {
	switch x := a.(type) {

	case *resourceschema.StringAttribute:
		if def.IsArray || def.Type != "string" || len(x.Validators) > 0 {
			return
		}
		x.Validators = append(x.Validators, stringvalidator.OneOf(def.Values...))
	case *datasourceschema.StringAttribute:
		if def.IsArray || def.Type != "string" || len(x.Validators) > 0 {
			return
		}
		x.Validators = append(x.Validators, stringvalidator.OneOf(def.Values...))

	case *resourceschema.Int64Attribute:
		ints, ok := parseEnumInt64s(def)
		if !ok || len(x.Validators) > 0 {
			return
		}
		x.Validators = append(x.Validators, int64validator.OneOf(ints...))
	case *datasourceschema.Int64Attribute:
		ints, ok := parseEnumInt64s(def)
		if !ok || len(x.Validators) > 0 {
			return
		}
		x.Validators = append(x.Validators, int64validator.OneOf(ints...))

	case *resourceschema.ListAttribute:
		if !def.IsArray || len(x.Validators) > 0 {
			return
		}
		if v := elementOneOfList(x.ElementType, def); v != nil {
			x.Validators = append(x.Validators, v)
		}
	case *datasourceschema.ListAttribute:
		if !def.IsArray || len(x.Validators) > 0 {
			return
		}
		if v := elementOneOfList(x.ElementType, def); v != nil {
			x.Validators = append(x.Validators, v)
		}

	case *resourceschema.SetAttribute:
		if !def.IsArray || len(x.Validators) > 0 {
			return
		}
		if v := elementOneOfSet(x.ElementType, def); v != nil {
			x.Validators = append(x.Validators, v)
		}
	case *datasourceschema.SetAttribute:
		if !def.IsArray || len(x.Validators) > 0 {
			return
		}
		if v := elementOneOfSet(x.ElementType, def); v != nil {
			x.Validators = append(x.Validators, v)
		}
	}
}

func parseEnumInt64s(def *enums.EnumDef) ([]int64, bool) {
	if def.IsArray || def.Type != "integer" {
		return nil, false
	}
	return parseInt64Slice(def.Values)
}

func parseInt64Slice(values []string) ([]int64, bool) {
	out := make([]int64, 0, len(values))
	for _, v := range values {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, false
		}
		out = append(out, n)
	}
	return out, true
}

func elementOneOfList(elem attr.Type, def *enums.EnumDef) validator.List {
	switch {
	case elem == types.StringType && def.Type == "string":
		return listvalidator.ValueStringsAre(stringvalidator.OneOf(def.Values...))
	case elem == types.Int64Type && def.Type == "integer":
		ints, ok := parseInt64Slice(def.Values)
		if !ok {
			return nil
		}
		return listvalidator.ValueInt64sAre(int64validator.OneOf(ints...))
	}
	return nil
}

func elementOneOfSet(elem attr.Type, def *enums.EnumDef) validator.Set {
	switch {
	case elem == types.StringType && def.Type == "string":
		return setvalidator.ValueStringsAre(stringvalidator.OneOf(def.Values...))
	case elem == types.Int64Type && def.Type == "integer":
		ints, ok := parseInt64Slice(def.Values)
		if !ok {
			return nil
		}
		return setvalidator.ValueInt64sAre(int64validator.OneOf(ints...))
	}
	return nil
}

// appendCompositionValidator attaches ExactlyOneOf or AtLeastOneOf validators
// derived from the composition definition to a SingleNestedAttribute.
// The composition branches are converted to sibling attribute paths.
// Skips when the call site has already attached any validator — that's the
// override discipline: a hand-coded validator at the call site wins.
// Only handles SingleNestedAttribute since oneOf/anyOf in the OpenAPI spec
// typically map to nested objects with mutually exclusive or at-least-one
// child attributes.
func appendCompositionValidator(a any, def *enums.CompositionDef) {
	// Only allOf doesn't need a validator (it's just merged constraints)
	if def.Kind == "allOf" {
		return
	}

	// Build paths from branch names (convert to snake_case for Terraform attribute names)
	paths := make([]path.Expression, 0, len(def.Branches))
	for _, branch := range def.Branches {
		// Convert branch schema name to likely Terraform attribute name
		// e.g., "AWSConfig" -> "aws_config", "GCPConfigData" -> "gcp_config_data"
		attrName := camelToSnake(branch)
		paths = append(paths, path.MatchRelative().AtParent().AtName(attrName))
	}

	if len(paths) < 2 {
		return
	}

	switch x := a.(type) {
	case *resourceschema.SingleNestedAttribute:
		if len(x.Validators) > 0 {
			return
		}
		switch def.Kind {
		case "oneOf":
			x.Validators = append(x.Validators, objectvalidator.ExactlyOneOf(paths...))
		case "anyOf":
			x.Validators = append(x.Validators, objectvalidator.AtLeastOneOf(paths...))
		}
	case *datasourceschema.SingleNestedAttribute:
		if len(x.Validators) > 0 {
			return
		}
		switch def.Kind {
		case "oneOf":
			x.Validators = append(x.Validators, objectvalidator.ExactlyOneOf(paths...))
		case "anyOf":
			x.Validators = append(x.Validators, objectvalidator.AtLeastOneOf(paths...))
		}
	}
}

func camelToSnake(s string) string {
	if s == "" {
		return s
	}
	var result []rune
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, r-'A'+'a')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
