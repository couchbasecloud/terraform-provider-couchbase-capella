package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestBackupScheduleSchema_NestedEnumValidators verifies the generated enum
// OneOf validators auto-attach to weekly_schedule's nested fields.
func TestBackupScheduleSchema_NestedEnumValidators(t *testing.T) {
	weekly, ok := BackupScheduleSchema().Attributes["weekly_schedule"].(*schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("weekly_schedule is %T, want *SingleNestedAttribute", BackupScheduleSchema().Attributes["weekly_schedule"])
	}

	t.Run("string enum fields reject invalid values", func(t *testing.T) {
		for name, bad := range map[string]string{"day_of_week": "funday", "retention_time": "7days"} {
			attr, ok := weekly.Attributes[name].(*schema.StringAttribute)
			if !ok {
				t.Fatalf("%s is %T, want *StringAttribute", name, weekly.Attributes[name])
			}
			if len(attr.Validators) == 0 {
				t.Fatalf("weekly_schedule.%s has no validators; OneOf enum validator not attached", name)
			}
			resp := &validator.StringResponse{}
			for _, v := range attr.Validators {
				v.ValidateString(context.Background(), validator.StringRequest{ConfigValue: types.StringValue(bad)}, resp)
			}
			if !resp.Diagnostics.HasError() {
				t.Errorf("weekly_schedule.%s accepted invalid value %q", name, bad)
			}
		}
	})

	t.Run("integer enum fields reject invalid values", func(t *testing.T) {
		for name, bad := range map[string]int64{"start_at": 24, "incremental_every": 3} {
			attr, ok := weekly.Attributes[name].(*schema.Int64Attribute)
			if !ok {
				t.Fatalf("%s is %T, want *Int64Attribute", name, weekly.Attributes[name])
			}
			if len(attr.Validators) == 0 {
				t.Fatalf("weekly_schedule.%s has no validators; OneOf enum validator not attached", name)
			}
			resp := &validator.Int64Response{}
			for _, v := range attr.Validators {
				v.ValidateInt64(context.Background(), validator.Int64Request{ConfigValue: types.Int64Value(bad)}, resp)
			}
			if !resp.Diagnostics.HasError() {
				t.Errorf("weekly_schedule.%s accepted invalid value %d", name, bad)
			}
		}
	})
}
