package schema

import (
	"reflect"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/docs"
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

	setMarkdownDescription(attr, description)

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
