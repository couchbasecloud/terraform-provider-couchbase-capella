package custom_plan_modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type blockCreateWhenEnabledSetToFalse struct{}

func BlockCreateWhenEnabledSetToFalse() planmodifier.Bool {
	return blockCreateWhenEnabledSetToFalse{}
}

func (b blockCreateWhenEnabledSetToFalse) Description(_ context.Context) string {
	return "when first enabling private endpoint service, disallow setting enabled to false"
}

func (b blockCreateWhenEnabledSetToFalse) MarkdownDescription(_ context.Context) string {
	return "when first enabling private endpoint service, disallow setting enabled to false"
}

func (b blockCreateWhenEnabledSetToFalse) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.StateValue.IsNull() && !req.ConfigValue.IsNull() && !req.ConfigValue.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot set enabled to false when first enabling private endpoint service",
			"Cannot set enabled to false when first enabling private endpoint service",
		)
		return
	}
}
