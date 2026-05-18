package validator

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ validator.Object = (*exactlyOneOfNestedValidator)(nil)
	_ validator.Object = (*atLeastOneOfNestedValidator)(nil)
)

// ExactlyOneOfNested returns a validator that ensures exactly one of the specified
// nested attributes is non-null. Unlike the standard ExactlyOneOf validator, this
// correctly handles SingleNestedAttribute by checking if the object value is actually
// null, not just empty or unknown.
func ExactlyOneOfNested(attributeNames ...string) validator.Object {
	return &exactlyOneOfNestedValidator{
		attributeNames: attributeNames,
	}
}

// AtLeastOneOfNested returns a validator that ensures at least one of the specified
// nested attributes is non-null. Unlike the standard AtLeastOneOf validator, this
// correctly handles SingleNestedAttribute by checking if the object value is actually
// null, not just empty or unknown.
func AtLeastOneOfNested(attributeNames ...string) validator.Object {
	return &atLeastOneOfNestedValidator{
		attributeNames: attributeNames,
	}
}

type exactlyOneOfNestedValidator struct {
	attributeNames []string
}

func (v *exactlyOneOfNestedValidator) Description(_ context.Context) string {
	return v.MarkdownDescription(context.Background())
}

func (v *exactlyOneOfNestedValidator) MarkdownDescription(_ context.Context) string {
	sorted := make([]string, len(v.attributeNames))
	copy(sorted, v.attributeNames)
	sort.Strings(sorted)
	return fmt.Sprintf("Exactly one of %s must be specified", formatAttributeList(sorted))
}

func (v *exactlyOneOfNestedValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	specifiedAttrs := v.countSpecifiedAttributes(ctx, req.ConfigValue)

	if len(specifiedAttrs) == 0 {
		sorted := make([]string, len(v.attributeNames))
		copy(sorted, v.attributeNames)
		sort.Strings(sorted)
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing Required Attribute",
			fmt.Sprintf("Exactly one of %s must be specified, but none were provided.", formatAttributeList(sorted)),
		)
		return
	}

	if len(specifiedAttrs) > 1 {
		sort.Strings(specifiedAttrs)
		sorted := make([]string, len(v.attributeNames))
		copy(sorted, v.attributeNames)
		sort.Strings(sorted)
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Combination",
			fmt.Sprintf("Exactly one of %s must be specified, but %s were all specified.", formatAttributeList(sorted), formatAttributeList(specifiedAttrs)),
		)
	}
}

func (v *exactlyOneOfNestedValidator) countSpecifiedAttributes(_ context.Context, obj types.Object) []string {
	var specified []string

	attrs := obj.Attributes()
	for _, name := range v.attributeNames {
		attrValue, ok := attrs[name]
		if !ok {
			continue
		}

		if isAttributeSpecified(attrValue) {
			specified = append(specified, name)
		}
	}

	return specified
}

type atLeastOneOfNestedValidator struct {
	attributeNames []string
}

func (v *atLeastOneOfNestedValidator) Description(_ context.Context) string {
	return v.MarkdownDescription(context.Background())
}

func (v *atLeastOneOfNestedValidator) MarkdownDescription(_ context.Context) string {
	sorted := make([]string, len(v.attributeNames))
	copy(sorted, v.attributeNames)
	sort.Strings(sorted)
	return fmt.Sprintf("At least one of %s must be specified", formatAttributeList(sorted))
}

func (v *atLeastOneOfNestedValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	specifiedAttrs := v.countSpecifiedAttributes(ctx, req.ConfigValue)

	if len(specifiedAttrs) == 0 {
		sorted := make([]string, len(v.attributeNames))
		copy(sorted, v.attributeNames)
		sort.Strings(sorted)
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing Required Attribute",
			fmt.Sprintf("At least one of %s must be specified, but none were provided.", formatAttributeList(sorted)),
		)
	}
}

func (v *atLeastOneOfNestedValidator) countSpecifiedAttributes(_ context.Context, obj types.Object) []string {
	var specified []string

	attrs := obj.Attributes()
	for _, name := range v.attributeNames {
		attrValue, ok := attrs[name]
		if !ok {
			continue
		}

		if isAttributeSpecified(attrValue) {
			specified = append(specified, name)
		}
	}

	return specified
}

// isAttributeSpecified checks if an attribute value is actually specified by the user.
// For nested objects, this checks if the object is non-null AND has at least one
// configured (non-null or unknown) attribute inside it.
// Unknown values are treated as "specified" because they represent user-configured
// values that reference computed attributes (e.g., account_id = other_resource.id).
func isAttributeSpecified(attrValue attr.Value) bool {
	if attrValue.IsNull() {
		return false
	}

	// If the entire object is unknown, treat it as specified
	// (user configured it, but values aren't resolved yet)
	if attrValue.IsUnknown() {
		return true
	}

	// For object types, check if any child attribute is configured (non-null or unknown)
	if objVal, ok := attrValue.(types.Object); ok {
		return hasConfiguredAttribute(objVal)
	}

	// For other types (string, int, etc.), non-null means specified
	return true
}

// hasConfiguredAttribute checks if an object has at least one configured attribute.
// A configured attribute is one that is either:
// - non-null and known (user provided a literal value), OR
// - unknown (user provided a reference to a computed value)
// This distinguishes between empty objects initialized by Terraform and objects
// where the user actually specified values (even if those values are computed).
func hasConfiguredAttribute(obj types.Object) bool {
	for _, v := range obj.Attributes() {
		// Null means not configured
		if v.IsNull() {
			continue
		}

		// Unknown means configured (e.g., references a computed value)
		if v.IsUnknown() {
			return true
		}

		// For nested objects, recurse
		if nestedObj, ok := v.(types.Object); ok {
			if hasConfiguredAttribute(nestedObj) {
				return true
			}
			continue
		}

		// Non-null, known, non-object value found
		return true
	}
	return false
}

func formatAttributeList(names []string) string {
	if len(names) == 0 {
		return "[]"
	}
	return "[" + strings.Join(names, ", ") + "]"
}

// PathExpressionsFromNames converts attribute names to path expressions relative to the parent.
func PathExpressionsFromNames(names []string) []path.Expression {
	exprs := make([]path.Expression, len(names))
	for i, name := range names {
		exprs[i] = path.MatchRelative().AtName(name)
	}
	return exprs
}
