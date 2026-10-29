package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// The functions in this file are shared by the alert configuration resources of every level.
// levelPath is the level's path relative to /v4/organizations, and level names it in messages.

// applyAlertConfiguration makes the level match plan and returns the settings to keep in state.
// Create and Update both use it: the request is worked out against what is set at the level now,
// rather than against state, so that only what actually differs is sent, and so that mutes set
// outside Terraform are cleared even when the muted_alerts block is new.
func applyAlertConfiguration(ctx context.Context, data *providerschema.Data, levelPath string, plan providerschema.AlertConfigurationSettings) (providerschema.AlertConfigurationSettings, error) {
	current, err := getAlertConfiguration(ctx, data, levelPath)
	if err != nil {
		return providerschema.AlertConfigurationSettings{}, err
	}

	update, err := alertConfigurationUpdate(ctx, plan, *current)
	if err != nil {
		return providerschema.AlertConfigurationSettings{}, err
	}

	if !isEmptyAlertConfigurationUpdate(update) {
		if err := putAlertConfiguration(ctx, data, levelPath, update); err != nil {
			return providerschema.AlertConfigurationSettings{}, err
		}
		if current, err = getAlertConfiguration(ctx, data, levelPath); err != nil {
			return providerschema.AlertConfigurationSettings{}, err
		}
	}

	return newAlertConfigurationSettings(ctx, plan, *current)
}

// readAlertConfiguration returns the settings to keep in state, refreshing only what prior tracks.
func readAlertConfiguration(ctx context.Context, data *providerschema.Data, levelPath string, prior providerschema.AlertConfigurationSettings) (providerschema.AlertConfigurationSettings, error) {
	current, err := getAlertConfiguration(ctx, data, levelPath)
	if err != nil {
		return providerschema.AlertConfigurationSettings{}, err
	}
	return newAlertConfigurationSettings(ctx, prior, *current)
}

// deleteAlertConfiguration clears the mutes and recipients that state manages. There is no
// delete endpoint, and thresholds and frequencies cannot be removed, so they are left in place.
// A level that no longer exists is not an error.
func deleteAlertConfiguration(ctx context.Context, data *providerschema.Data, levelPath string, state providerschema.AlertConfigurationSettings) error {
	current, err := getAlertConfiguration(ctx, data, levelPath)
	if err != nil {
		if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
			return nil
		}
		return err
	}

	update := alertConfigurationDeleteUpdate(state, *current)
	if isEmptyAlertConfigurationUpdate(update) {
		return nil
	}

	if err := putAlertConfiguration(ctx, data, levelPath, update); err != nil {
		if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
			return nil
		}
		return err
	}
	return nil
}

func getAlertConfiguration(ctx context.Context, data *providerschema.Data, levelPath string) (*alertsapi.GetAlertConfigurationResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/alertConfiguration", data.HostURL, levelPath)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(ctx, cfg, nil, data.Token, nil)
	if err != nil {
		return nil, err
	}

	var config alertsapi.GetAlertConfigurationResponse
	if err := json.Unmarshal(response.Body, &config); err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}
	return &config, nil
}

func putAlertConfiguration(ctx context.Context, data *providerschema.Data, levelPath string, update alertsapi.UpdateAlertConfigurationRequest) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/alertConfiguration", data.HostURL, levelPath)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := data.ClientV1.ExecuteWithRetry(ctx, cfg, update, data.Token, nil)
	return err
}

func newAlertConfigurationSettings(ctx context.Context, prior providerschema.AlertConfigurationSettings, current alertsapi.GetAlertConfigurationResponse) (providerschema.AlertConfigurationSettings, error) {
	settings, diags := providerschema.NewAlertConfigurationSettings(ctx, prior, current)
	if diags.HasError() {
		return providerschema.AlertConfigurationSettings{}, fmt.Errorf("%s: %v", errors.ErrRefreshingState, diags.Errors())
	}
	return settings, nil
}

// alertConfigurationUpdate works out the request that makes the level match plan:
//   - configs: each alert with a field that differs from what is set now. Alerts and fields not
//     in plan are left out, since an override cannot be removed.
//   - recipients: the full list, when plan sets it and it differs from the current list.
//   - mutes: when plan has the muted_alerts block, a mute or unmute for each mute that differs,
//     and a clear for each mute at the level that plan does not list.
func alertConfigurationUpdate(ctx context.Context, plan providerschema.AlertConfigurationSettings, current alertsapi.GetAlertConfigurationResponse) (alertsapi.UpdateAlertConfigurationRequest, error) {
	var update alertsapi.UpdateAlertConfigurationRequest

	for name, setting := range plan.Configs {
		if !alertConfigurationMatches(setting, current.Configs[name]) {
			if update.Configs == nil {
				update.Configs = make(map[string]alertsapi.AlertConfigurationUpdate)
			}
			update.Configs[name] = alertsapi.AlertConfigurationUpdate{
				FrequencyHours:    setting.FrequencyHours.ValueInt64Pointer(),
				Threshold:         setting.Threshold.ValueInt64Pointer(),
				WarningThreshold:  setting.WarningThreshold.ValueInt64Pointer(),
				CriticalThreshold: setting.CriticalThreshold.ValueInt64Pointer(),
			}
		}
	}

	if !plan.Recipients.IsNull() && !plan.Recipients.IsUnknown() {
		recipients := []string{}
		if diags := plan.Recipients.ElementsAs(ctx, &recipients, false); diags.HasError() {
			return update, fmt.Errorf("could not read recipients: %v", diags.Errors())
		}
		if !sameAddresses(recipients, current.Recipients) {
			sort.Strings(recipients)
			update.Recipients = &recipients
		}
	}

	if plan.MutedAlerts != nil {
		for _, name := range sortedKeys(plan.MutedAlerts) {
			mute := plan.MutedAlerts[name]
			if currentMute, ok := current.MutedAlerts[name]; ok && alertMuteMatches(mute, currentMute) {
				continue
			}
			update.Alerts = append(update.Alerts, alertsapi.AlertMuteAction{
				AlertName:  name,
				Action:     mute.Action.ValueString(),
				MutedFrom:  timePointer(mute.MutedFrom),
				MutedUntil: timePointer(mute.MutedUntil),
			})
		}
		for _, name := range sortedKeys(current.MutedAlerts) {
			if _, ok := plan.MutedAlerts[name]; !ok {
				update.Alerts = append(update.Alerts, alertsapi.AlertMuteAction{AlertName: name, Action: alertsapi.AlertMuteActionClear})
			}
		}
	}

	return update, nil
}

// alertConfigurationDeleteUpdate clears every mute at the level when state manages mutes, and the
// recipients when state manages them.
func alertConfigurationDeleteUpdate(state providerschema.AlertConfigurationSettings, current alertsapi.GetAlertConfigurationResponse) alertsapi.UpdateAlertConfigurationRequest {
	var update alertsapi.UpdateAlertConfigurationRequest

	if state.MutedAlerts != nil {
		for _, name := range sortedKeys(current.MutedAlerts) {
			update.Alerts = append(update.Alerts, alertsapi.AlertMuteAction{AlertName: name, Action: alertsapi.AlertMuteActionClear})
		}
	}

	if !state.Recipients.IsNull() && len(current.Recipients) > 0 {
		update.Recipients = &[]string{}
	}

	return update
}

func isEmptyAlertConfigurationUpdate(update alertsapi.UpdateAlertConfigurationRequest) bool {
	return len(update.Configs) == 0 && update.Recipients == nil && len(update.Alerts) == 0
}

func alertConfigurationMatches(setting providerschema.AlertConfigurationSetting, current alertsapi.AlertConfigurationEntry) bool {
	return int64Matches(setting.FrequencyHours, current.FrequencyHours) &&
		thresholdMatches(setting.Threshold, current.Threshold) &&
		thresholdMatches(setting.WarningThreshold, current.WarningThreshold) &&
		thresholdMatches(setting.CriticalThreshold, current.CriticalThreshold)
}

// int64Matches reports whether a planned field already has its value. A field the plan does not
// set always matches, since it is left as it is.
func int64Matches(planned types.Int64, current *int64) bool {
	return planned.IsNull() || (current != nil && *current == planned.ValueInt64())
}

func thresholdMatches(planned types.Int64, current *float64) bool {
	return planned.IsNull() || (current != nil && *current == float64(planned.ValueInt64()))
}

func alertMuteMatches(mute providerschema.AlertMuteSetting, current alertsapi.AlertMuteState) bool {
	return mute.Action.ValueString() == current.Action &&
		timeMatches(mute.MutedFrom, current.MutedFrom) &&
		timeMatches(mute.MutedUntil, current.MutedUntil)
}

func timeMatches(planned timetypes.RFC3339, current *string) bool {
	if planned.IsNull() || current == nil {
		return planned.IsNull() && current == nil
	}
	plannedTime, diags := planned.ValueRFC3339Time()
	if diags.HasError() {
		return false
	}
	currentTime, err := time.Parse(time.RFC3339Nano, *current)
	return err == nil && plannedTime.Equal(currentTime)
}

func timePointer(t timetypes.RFC3339) *string {
	if t.IsNull() || t.IsUnknown() {
		return nil
	}
	return t.ValueStringPointer()
}

// sameAddresses compares two recipient lists as sets, ignoring case, since Capella stores
// addresses in lowercase.
func sameAddresses(a, b []string) bool {
	set := func(addresses []string) map[string]struct{} {
		out := make(map[string]struct{}, len(addresses))
		for _, address := range addresses {
			out[strings.ToLower(address)] = struct{}{}
		}
		return out
	}
	as, bs := set(a), set(b)
	if len(as) != len(bs) {
		return false
	}
	for address := range as {
		if _, ok := bs[address]; !ok {
			return false
		}
	}
	return true
}

// alertConfigurationRemovalWarning describes the alerts and fields that state manages and plan
// no longer does. Capella keeps their values, so without a warning the removal would look like a
// reset that silently did nothing. It returns "" when nothing was removed.
func alertConfigurationRemovalWarning(level string, state, plan providerschema.AlertConfigurationSettings) string {
	var lines []string
	for _, name := range sortedKeys(state.Configs) {
		planned, ok := plan.Configs[name]
		var kept []string
		for i, field := range state.Configs[name].Fields() {
			if field.Value.IsNull() {
				continue
			}
			if !ok || planned.Fields()[i].Value.IsNull() {
				kept = append(kept, fmt.Sprintf("%s = %d", field.Name, field.Value.ValueInt64()))
			}
		}
		if len(kept) > 0 {
			lines = append(lines, fmt.Sprintf("  - %s: %s", name, strings.Join(kept, ", ")))
		}
	}
	if len(lines) == 0 {
		return ""
	}

	return fmt.Sprintf("Capella cannot remove a threshold or frequency override, so the %s keeps these values and Terraform no longer manages them:\n%s\n\n"+
		"To change a value, set the value you want in configs instead of removing it.", level, strings.Join(lines, "\n"))
}

// alertConfigurationLeftBehindWarning describes the thresholds and frequencies that stay at the
// level after destroy. It returns "" when state sets none.
func alertConfigurationLeftBehindWarning(level string, state providerschema.AlertConfigurationSettings) string {
	var lines []string
	for _, name := range sortedKeys(state.Configs) {
		var kept []string
		for _, field := range state.Configs[name].Fields() {
			if !field.Value.IsNull() {
				kept = append(kept, fmt.Sprintf("%s = %d", field.Name, field.Value.ValueInt64()))
			}
		}
		if len(kept) > 0 {
			lines = append(lines, fmt.Sprintf("  - %s: %s", name, strings.Join(kept, ", ")))
		}
	}
	if len(lines) == 0 {
		return ""
	}

	return fmt.Sprintf("Capella cannot remove a threshold or frequency override, so the %s keeps these values after destroy:\n%s", level, strings.Join(lines, "\n"))
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
