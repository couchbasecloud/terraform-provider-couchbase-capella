package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	clusterAlertConfigurationBuilder = capellaschema.NewSchemaBuilder("clusterAlertConfiguration", "UpdateAlertConfigurationRequest")
	alertConfigurationSettingBuilder = capellaschema.NewSchemaBuilder("alertConfigurationSetting", "AlertConfigurationEntry")
	alertMuteSettingBuilder          = capellaschema.NewSchemaBuilder("alertMuteSetting", "AlertMuteAction")
)

// The OpenAPI loader does not read a description written next to a $ref inside allOf, which is how
// the spec documents configs, and the spec names the mutes "alerts" in requests. These carry a
// fallback in the loader's bullet format; AddAttr prefers the spec as soon as it can read one.
const (
	alertConfigsDescription     = "\n - The alert configs to set at this level, keyed by alert identifier, for example `cpu_usage_high`. Terraform manages only the alerts listed here, and only the fields set on each."
	alertMutedAlertsDescription = "\n - The mute overrides at this level, keyed by alert identifier. When this attribute is set, Terraform manages every mute at this level, including ones added in the Capella UI. Leave it out to leave all mutes alone."
)

// alertConfigurationOwnershipDescription is shared by the resource descriptions of all levels.
const alertConfigurationOwnershipDescription = `

Manage each level with only one resource, so that two resources do not overwrite each other's changes.

Each attribute follows its own rule about what Terraform owns:

- ` + "`configs`" + `: only the alerts listed, and only the fields set on each. Alerts configured outside Terraform, for example in the Capella UI, are left alone.
- ` + "`muted_alerts`" + `: every mute at this level, once the attribute is set. A mute added outside Terraform shows up in the next plan, and ` + "`terraform apply`" + ` clears it. Leave the attribute out to leave all mutes alone.
- ` + "`recipients`" + `: the whole list, if set. Removing the attribute leaves the list as it is.

Capella cannot remove a threshold or frequency override once it is set. Removing an alert from ` + "`configs`" + `, or a field from an alert, leaves the value Terraform last set in place and stops Terraform from managing it; ` + "`terraform apply`" + ` warns when this happens. To change a value, set the value you want instead.

On destroy, Terraform clears the mutes and recipients it manages, so the level goes back to inheriting them. Thresholds and frequencies stay with the values Terraform last set, and ` + "`terraform destroy`" + ` warns which ones.

Mute keys must be the alert identifier used in ` + "`configs`" + `, such as ` + "`cpu_usage_high`" + `, not a severity-specific name such as ` + "`CPUUsageHighWarning`" + `. A mute whose ` + "`muted_until`" + ` has passed stays at this level and in state until it is changed or removed; set a later ` + "`muted_until`" + ` before changing anything else on it.`

// ClusterAlertConfigurationSchema returns the schema for the cluster alert configuration resource.
func ClusterAlertConfigurationSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", clusterAlertConfigurationBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", clusterAlertConfigurationBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", clusterAlertConfigurationBuilder, requiredUUIDStringAttribute())
	addAlertConfigurationSettingsAttrs(attrs, clusterAlertConfigurationBuilder)

	return schema.Schema{
		MarkdownDescription: "Manages the alert configuration that a cluster sets for itself: threshold and frequency overrides, email recipients and mutes. What the cluster inherits from its project, its organization and the Capella defaults is shown by the `couchbase-capella_cluster_alert_configuration` data source." +
			alertConfigurationOwnershipDescription,
		Attributes: attrs,
	}
}

// addAlertConfigurationSettingsAttrs adds the attributes shared by the alert configuration resources of every level.
func addAlertConfigurationSettingsAttrs(attrs map[string]schema.Attribute, builder *capellaschema.SchemaBuilder) {
	configs := &schema.MapNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: alertConfigurationSettingAttrs(),
			Validators: []validator.Object{alertConfigurationSettingNotEmpty{}},
		},
	}
	configs.MarkdownDescription = alertConfigsDescription
	capellaschema.AddAttr(attrs, "configs", builder, configs)

	recipients := stringSetAttribute(optional)
	recipients.Validators = append(recipients.Validators, recipientsUniqueIgnoringCase{})
	capellaschema.AddAttr(attrs, "recipients", builder, recipients)

	mutedAlerts := &schema.MapNestedAttribute{
		Optional:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: alertMuteSettingAttrs()},
	}
	mutedAlerts.MarkdownDescription = alertMutedAlertsDescription
	capellaschema.AddAttr(attrs, "muted_alerts", builder, mutedAlerts)
}

func alertConfigurationSettingAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "frequency_hours", alertConfigurationSettingBuilder, int64Attribute(optional))
	for _, name := range []string{"threshold", "warning_threshold", "critical_threshold"} {
		threshold := int64Attribute(optional)
		threshold.Validators = append(threshold.Validators, int64validator.NoneOf(0))
		capellaschema.AddAttr(attrs, name, alertConfigurationSettingBuilder, threshold)
	}
	return attrs
}

func alertMuteSettingAttrs() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "action", alertMuteSettingBuilder, stringAttribute(
		[]string{required},
		stringvalidator.OneOf(alertsapi.AlertMuteActionMute, alertsapi.AlertMuteActionUnmute),
	), "AlertMuteState")
	capellaschema.AddAttr(attrs, "muted_from", alertMuteSettingBuilder, rfc3339Attribute(optional))
	capellaschema.AddAttr(attrs, "muted_until", alertMuteSettingBuilder, rfc3339Attribute(optional))
	return attrs
}

// alertConfigurationSettingNotEmpty rejects an alert in configs that sets no field.
type alertConfigurationSettingNotEmpty struct{}

func (v alertConfigurationSettingNotEmpty) Description(_ context.Context) string {
	return "at least one of frequency_hours, threshold, warning_threshold or critical_threshold must be set"
}

func (v alertConfigurationSettingNotEmpty) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v alertConfigurationSettingNotEmpty) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for _, value := range req.ConfigValue.Attributes() {
		if !value.IsNull() {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(req.Path, "Empty Alert Configuration", "The alert sets no field, so there is nothing to configure: "+v.Description(ctx)+".")
}

// recipientsUniqueIgnoringCase rejects addresses that differ only by case. Capella lowercases and
// deduplicates recipients, so it would store one address where the plan has two, and the apply
// would fail with an inconsistent result.
type recipientsUniqueIgnoringCase struct{}

func (v recipientsUniqueIgnoringCase) Description(_ context.Context) string {
	return "each address must be unique, ignoring case"
}

func (v recipientsUniqueIgnoringCase) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v recipientsUniqueIgnoringCase) ValidateSet(_ context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	seen := make(map[string]string)
	for _, element := range req.ConfigValue.Elements() {
		address, ok := element.(types.String)
		if !ok || address.IsNull() || address.IsUnknown() {
			continue
		}
		key := strings.ToLower(address.ValueString())
		if first, ok := seen[key]; ok {
			resp.Diagnostics.AddAttributeError(req.Path, "Duplicate Recipient",
				fmt.Sprintf("%q and %q are the same address. Capella stores addresses in lowercase, so list each address once.", first, address.ValueString()))
			return
		}
		seen[key] = address.ValueString()
	}
}
