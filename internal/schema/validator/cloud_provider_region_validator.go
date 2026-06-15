package validator

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.Object = (*cloudProviderRegionValidator)(nil)

// CloudProviderRegion returns an object validator for a cloud_provider block that
// rejects a region not supported for the configured provider type. allowed maps a
// lowercased provider type (e.g. "aws") to its supported region codes. The check
// is skipped when type or region is null/unknown, or when the provider type is
// absent from the map (its own enum validator reports that).
func CloudProviderRegion(allowed map[string][]string) validator.Object {
	return &cloudProviderRegionValidator{allowed: allowed}
}

type cloudProviderRegionValidator struct {
	allowed map[string][]string
}

func (v *cloudProviderRegionValidator) Description(_ context.Context) string {
	return v.MarkdownDescription(context.Background())
}

func (v *cloudProviderRegionValidator) MarkdownDescription(_ context.Context) string {
	return "region must be one of the supported regions for the configured cloud provider type"
}

func (v *cloudProviderRegionValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	providerType, ok := stringValue(attrs["type"])
	if !ok {
		return
	}
	region, ok := stringValue(attrs["region"])
	if !ok {
		return
	}

	regions, ok := v.allowed[strings.ToLower(providerType)]
	if !ok {
		// Unknown provider type — the type attribute's own enum validator reports it.
		return
	}

	if slices.Contains(regions, region) {
		return
	}

	sorted := append([]string(nil), regions...)
	sort.Strings(sorted)
	resp.Diagnostics.AddAttributeError(
		req.Path.AtName("region"),
		"Unsupported Region",
		fmt.Sprintf("region %q is not supported for cloud provider %q. Supported regions are: %s.",
			region, providerType, strings.Join(sorted, ", ")),
	)
}

// stringValue returns the string and true only when value is a known, non-null
// String; otherwise it returns false so the caller can skip validation.
func stringValue(value any) (string, bool) {
	s, ok := value.(types.String)
	if !ok || s.IsNull() || s.IsUnknown() {
		return "", false
	}
	return s.ValueString(), true
}
