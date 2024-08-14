package custom_plan_modifier

/*
	this custom plan modifier is used by private endpoint service resource to handle external changes to
	the enabled bool value.

	the enabled bool attribute is computed.  the user cannot enable/disable the service by setting this value.
	the service is enabled with POST (ie add resource block), and disabled with DELETE (ie remove resource block).

	enabled value is set based on API server response. if this is changed outside terraform the provider will
	take no action.  this is correct behavior since there is no desired state from user configuration.  it is
	determined by the provider.  this is why custom plan modifier is needed.

	the provider only needs to handle the case when service is enabled by terraform and user disables it externally
	(for example through the UI).  note that the provider will overwrite the state value to false.

	if the service is disabled by terraform the resource is removed from the state file.

	private state is used so that it's not visible to the user.  for more information about private state see:
	https://developer.hashicorp.com/terraform/plugin/framework/resources/private-state

	for more information on custom plan modifiers see:
	https://developer.hashicorp.com/terraform/plugin/framework/resources/plan-modification#creating-attribute-plan-modifiers
*/

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type triggerUpdateOnExternalStatusChange struct{}

func TriggerUpdateOnExternalStatusChange() planmodifier.Bool {
	return triggerUpdateOnExternalStatusChange{}
}

func (d triggerUpdateOnExternalStatusChange) Description(_ context.Context) string {
	return "trigger update if private endpoint service is disabled outside of terraform"
}

func (d triggerUpdateOnExternalStatusChange) MarkdownDescription(_ context.Context) string {
	return "trigger update if private endpoint service is disabled outside of terraform"
}

func (d triggerUpdateOnExternalStatusChange) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	value, diags := req.Private.GetKey(ctx, "is_enabled")
	resp.Diagnostics.Append(diags...)

	if value != nil {
		if isEnabled, err := strconv.ParseBool(string(value)); err == nil {
			if !req.StateValue.ValueBool() && isEnabled {
				resp.PlanValue = types.BoolValue(true)
			}
		}
	}
}
