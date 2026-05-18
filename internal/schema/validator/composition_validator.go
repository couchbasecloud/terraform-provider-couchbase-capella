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
// non-null attribute inside it (to handle Terraform's empty object initialization).
func isAttributeSpecified(attrValue attr.Value) bool {
	if attrValue.IsNull() {
		return false
	}

	// For unknown values, we can't determine if specified - treat as not specified
	// to avoid false positives during planning
	if attrValue.IsUnknown() {
		return false
	}

	// For object types, check if any child attribute is non-null
	if objVal, ok := attrValue.(types.Object); ok {
		return hasNonNullAttribute(objVal)
	}

	// For other types (string, int, etc.), non-null means specified
	return true
}

// hasNonNullAttribute checks if an object has at least one non-null, non-unknown attribute.
func hasNonNullAttribute(obj types.Object) bool {
	for _, v := range obj.Attributes() {
		if v.IsNull() || v.IsUnknown() {
			continue
		}

		// For nested objects, recurse
		if nestedObj, ok := v.(types.Object); ok {
			if hasNonNullAttribute(nestedObj) {
				return true
			}
			continue
		}

		// Non-null, non-unknown, non-object value found
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
