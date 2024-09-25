package custom_plan_modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type immutableInt64Attribute struct{}

func ImmutableInt64Attribute() planmodifier.Int64 {
	return immutableInt64Attribute{}
}

func (n immutableInt64Attribute) Description(ctx context.Context) string {
	return "property cannot be changed"
}

func (n immutableInt64Attribute) MarkdownDescription(ctx context.Context) string {
	return "property cannot be changed"
}

func (n immutableInt64Attribute) PlanModifyInt64(
	_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response,
) {
	if !req.StateValue.IsNull() &&
		!req.ConfigValue.IsNull() &&
		!req.ConfigValue.Equal(req.StateValue) {
		resp.Diagnostics.AddError(
			"property cannot be changed",
			fmt.Sprintf("%s cannot be changed", req.Path.String()),
		)
		return
	}
}
