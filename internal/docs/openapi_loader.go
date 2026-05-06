package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/net/html"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schemawalk"
)

var openAPIDoc *openapi3.T

func init() {
	loadOpenAPISpec()
}

func loadOpenAPISpec() {
	openAPISource := os.Getenv("OPENAPI_SPEC_URL")
	if openAPISource == "" {
		return
	}

	data, err := fetchFromURL(openAPISource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not fetch OpenAPI spec from %s: %v\n", openAPISource, err)
		fmt.Fprintf(os.Stderr, "Field descriptions will not be enhanced with OpenAPI metadata.\n")
		return
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromData(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to parse OpenAPI spec: %v\n", err)
		fmt.Fprintf(os.Stderr, "Field descriptions will not be enhanced with OpenAPI metadata.\n")
		return
	}

	openAPIDoc = doc
}

func fetchFromURL(url string) ([]byte, error) {
	client := &http.Client{Timeout: 60 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if isHTML(data) {
		return extractEmbeddedSpec(data)
	}
	return data, nil
}

func isHTML(data []byte) bool {
	return strings.HasPrefix(http.DetectContentType(data), "text/html")
}

func extractEmbeddedSpec(htmlData []byte) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var findSpec func(*html.Node) string
	findSpec = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "script" && n.FirstChild != nil {
			if idx := strings.Index(n.FirstChild.Data, `{"openapi":"3.0`); idx != -1 {
				spec, _ := extractJSONObject(n.FirstChild.Data[idx:])
				return spec
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if spec := findSpec(c); spec != "" {
				return spec
			}
		}
		return ""
	}

	if spec := findSpec(doc); spec != "" {
		return []byte(spec), nil
	}
	return nil, fmt.Errorf("could not find embedded OpenAPI spec in HTML page")
}

func extractJSONObject(s string) (string, error) {
	dec := json.NewDecoder(strings.NewReader(s))
	var raw json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return "", fmt.Errorf("failed to decode JSON: %w", err)
	}
	return string(raw), nil
}

// GetOpenAPIDescription retrieves an enhanced description for a Terraform field
// from the loaded OpenAPI spec. Returns "" when the spec is not loaded or the
// field is not found.
func GetOpenAPIDescription(resourceName, tfFieldName string) string {
	if openAPIDoc == nil || openAPIDoc.Components == nil {
		return ""
	}

	switch tfFieldName {
	case "if_match":
		if openAPIDoc.Components.Parameters != nil {
			if p, ok := openAPIDoc.Components.Parameters["If-Match"]; ok && p.Value != nil && p.Value.Description != "" {
				return strings.TrimSpace(p.Value.Description)
			}
		}

	case "etag":
		if openAPIDoc.Components.Headers != nil {
			if h, ok := openAPIDoc.Components.Headers["ETag"]; ok && h.Value != nil && h.Value.Description != "" {
				return strings.TrimSpace(h.Value.Description)
			}
		}

	case "audit":
		if openAPIDoc.Components.Schemas != nil {
			if s, ok := openAPIDoc.Components.Schemas["CouchbaseAuditData"]; ok && s.Value != nil {
				return "Couchbase audit data."
			}
		}

	default:
		stringSchemaMap := map[string]string{
			"access_control_function": "AccessFunction",
			"import_filter":           "ImportFilter",
		}
		if schemaName, ok := stringSchemaMap[tfFieldName]; ok {
			if s, ok := openAPIDoc.Components.Schemas[schemaName]; ok && s.Value != nil && s.Value.Description != "" {
				return buildEnhancedDescription(s.Value)
			}
		}

		if openAPIDoc.Components.Parameters != nil {
			for _, paramName := range paramCandidates(tfFieldName) {
				if p, ok := openAPIDoc.Components.Parameters[paramName]; ok && p.Value != nil && p.Value.Description != "" {
					return strings.TrimSpace(p.Value.Description)
				}
			}
		}
	}

	if openAPIDoc.Components.Schemas == nil {
		return ""
	}

	camelField := schemawalk.SnakeToCamel(tfFieldName)
	for _, schemaName := range schemawalk.SchemaNameCandidates(resourceName) {
		ref := openAPIDoc.Components.Schemas[schemaName]
		if ref == nil || ref.Value == nil {
			continue
		}
		if prop := schemawalk.FindProperty(openAPIDoc.Components.Schemas, ref.Value, camelField); prop != nil {
			return buildEnhancedDescription(prop)
		}
	}

	return ""
}

// paramCandidates returns the OpenAPI parameter name candidates for a Terraform field.
func paramCandidates(tfFieldName string) []string {
	candidates := []string{schemawalk.SnakeToPascal(tfFieldName)}
	switch tfFieldName {
	case "app_endpoint_name":
		candidates = append(candidates, "appEndpointId", "appEndpointKeyspace")
	case "scope":
		candidates = append(candidates, "scopeName", "appEndpointKeyspace")
	case "collection":
		candidates = append(candidates, "collectionName", "appEndpointKeyspace")
	}
	return candidates
}

// buildEnhancedDescription creates a rich markdown description from an OpenAPI property schema.
func buildEnhancedDescription(prop *openapi3.Schema) string {
	referencedSchema := schemawalk.ResolveArrayRef(openAPIDoc.Components.Schemas, prop)
	description := extractDescription(prop, referencedSchema)
	enumValues := collectEnumValues(prop, referencedSchema)

	var parts []string

	if description != "" {
		parts = append(parts, formatDescriptionBullet(description))
	} else if len(enumValues) > 0 {
		parts = append(parts, formatValidValues(enumValues))
	}

	if constraints := buildConstraints(prop); len(constraints) > 0 {
		parts = append(parts, formatConstraintsBullet(constraints))
	}

	if description != "" && len(enumValues) > 0 {
		parts = append(parts, formatValidValues(enumValues))
	}

	if prop.Default != nil {
		parts = append(parts, formatDefaultBullet(prop.Default))
	}

	if prop.Format != "" {
		if fd := formatDescription(prop.Format); fd != "" {
			parts = append(parts, formatFormatBullet(fd))
		}
	}

	if prop.Deprecated {
		parts = append(parts, formatDeprecationBullet())
	}

	return strings.Join(parts, "")
}

func formatDescriptionBullet(desc string) string {
	return fmt.Sprintf("\n - %s", desc)
}

func formatConstraintsBullet(constraints []string) string {
	return fmt.Sprintf("\n - **Constraints**: %s", strings.Join(constraints, ", "))
}

func formatDefaultBullet(def interface{}) string {
	return fmt.Sprintf("\n - **Default**: `%v`", def)
}

func formatFormatBullet(fd string) string {
	return fmt.Sprintf("\n - **Format**: %s", fd)
}

func formatDeprecationBullet() string {
	return "\n - **Deprecated**: This field is deprecated and will be removed in a future release."
}

func extractDescription(prop, referenced *openapi3.Schema) string {
	if prop.Description != "" {
		return cleanDescription(prop.Description)
	}
	if referenced != nil && referenced.Description != "" {
		return cleanDescription(referenced.Description)
	}
	return ""
}

func collectEnumValues(prop, referenced *openapi3.Schema) []interface{} {
	if len(prop.Enum) > 0 {
		return prop.Enum
	}
	if referenced != nil {
		return referenced.Enum
	}
	return nil
}

func formatEnumValues(enumValues []interface{}) []string {
	quoted := make([]string, len(enumValues))
	for i, v := range enumValues {
		quoted[i] = fmt.Sprintf("`%v`", v)
	}
	return quoted
}

func formatValidValues(enumValues []interface{}) string {
	return fmt.Sprintf("\n - **Valid Values**: %s", strings.Join(formatEnumValues(enumValues), ", "))
}

func buildConstraints(prop *openapi3.Schema) []string {
	var c []string
	if prop.MinLength > 0 {
		c = append(c, fmt.Sprintf("Minimum length: %d characters", prop.MinLength))
	}
	if prop.MaxLength != nil && *prop.MaxLength > 0 {
		c = append(c, fmt.Sprintf("Maximum length: %d characters", *prop.MaxLength))
	}
	if prop.Min != nil {
		c = append(c, fmt.Sprintf("Minimum: %v", *prop.Min))
	}
	if prop.Max != nil {
		c = append(c, fmt.Sprintf("Maximum: %v", *prop.Max))
	}
	if prop.Pattern != "" {
		c = append(c, fmt.Sprintf("Pattern: `%s`", prop.Pattern))
	}
	return c
}

func formatDescription(format string) string {
	formats := map[string]string{
		"uuid":      "UUID (GUID4)",
		"date":      "Date in RFC3339 format",
		"date-time": "Date-time in RFC3339 format",
		"email":     "Email address",
		"uri":       "URI",
		"hostname":  "Hostname",
		"ipv4":      "IPv4 address",
		"ipv6":      "IPv6 address",
	}
	return formats[format]
}

func isTableRow(line string) bool {
	return strings.Count(line, "|") >= 2
}

func isListMarker(b byte) bool {
	return b == '-' || b == '*' || b == '+'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func appendTextBuffer(result, buffer []string) []string {
	if len(buffer) == 0 {
		return result
	}
	return append(result, strings.Join(buffer, " "))
}

func cleanDescription(desc string) string {
	if desc = strings.TrimSpace(desc); desc == "" {
		return ""
	}

	lines := strings.Split(desc, "\n")
	cleaned := make([]string, 0, len(lines))
	textBuf := make([]string, 0)
	inTable := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		switch {
		case isTableRow(line) && !inTable:
			cleaned = appendTextBuffer(cleaned, textBuf)
			textBuf = textBuf[:0]
			cleaned = append(cleaned, "", line)
			inTable = true
		case !isTableRow(line) && inTable:
			inTable = false
			if line != "" {
				cleaned = append(cleaned, "")
			}
			fallthrough
		case !inTable && line != "":
			if s := removeListMarkers(line); s != "" {
				textBuf = append(textBuf, s)
			}
		default:
			if inTable {
				cleaned = append(cleaned, line)
			}
		}
	}
	cleaned = appendTextBuffer(cleaned, textBuf)
	return strings.Join(cleaned, "\n")
}

func removeListMarkers(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}
	if isListMarker(line[0]) {
		return strings.TrimSpace(line[1:])
	}
	if isDigit(line[0]) {
		if dot := strings.Index(line, "."); dot > 0 && dot < len(line)-1 {
			allDigits := true
			for i := 0; i < dot; i++ {
				if !isDigit(line[i]) {
					allDigits = false
					break
				}
			}
			if allDigits {
				return strings.TrimSpace(line[dot+1:])
			}
		}
	}
	return line
}
