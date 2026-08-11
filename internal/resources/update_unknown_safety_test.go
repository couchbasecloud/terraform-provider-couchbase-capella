package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestUpdateBodyAttributesResolveUnknownsToState covers the Optional+Computed attributes
// that Update methods read into a request body. The framework marks a null-in-config
// Computed attribute as unknown unless it has a Default, and ValueString()/ValueInt64()/
// ValueBool() unwrap an unknown to "" / 0 / false, so such an attribute must resolve back
// to prior state via useStateForUnknown or the provider transmits an unrequested zero.
//
// Adding an Optional+Computed attribute to an Update request body means adding it here.
func TestUpdateBodyAttributesResolveUnknownsToState(t *testing.T) {
	cases := []struct {
		name       string
		attributes map[string]schema.Attribute
		attrNames  []string
	}{
		{
			// PutBucketRequest. Only durability_level surfaces - the API rejects "" with
			// a 422. The rest apply silently: 0 replicas, no TTL, no RAM quota.
			name:       "bucket",
			attributes: BucketSchema().Attributes,
			attrNames: []string{
				"memory_allocation_in_mb",
				"durability_level",
				"replicas",
				"time_to_live_in_seconds",
				"flush",
			},
		},
		{
			// PutProjectRequest. An empty description wipes the existing one.
			name:       "project",
			attributes: ProjectSchema().Attributes,
			attrNames:  []string{"description"},
		},
		{
			// UpdateClusterAuditSettingsRequest. A false audit_enabled disables auditing.
			// The sets are read into native slices, which cannot hold an unknown.
			name:       "audit_log_settings",
			attributes: AuditLogSettingsSchema().Attributes,
			attrNames:  []string{"audit_enabled", "enabled_event_ids", "disabled_users"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, attrName := range tc.attrNames {
				attr, ok := tc.attributes[attrName]
				if !ok {
					t.Errorf("%s: attribute not present in schema", attrName)
					continue
				}
				if !attr.IsOptional() || !attr.IsComputed() {
					// Only Optional+Computed attributes can be marked unknown this way.
					continue
				}
				if !resolvesUnknownToState(t, attr) {
					t.Errorf(
						"%s.%s is Optional+Computed and is read into an Update request body, "+
							"but an unconfigured value does not resolve to prior state: it would be "+
							"transmitted as a zero value. Add useStateForUnknown to it.",
						tc.name, attrName,
					)
				}
			}
		})
	}
}

// resolvesUnknownToState reports whether an unconfigured value (null config, unknown plan)
// is restored to prior state by the plan modifiers. A Default is exempt: it is a value
// declared in the schema, not a zero value from an unwrapped unknown.
//
// Each modifier sees the plan value the previous one produced, matching
// fwserver.AttributeModifyPlan.
func resolvesUnknownToState(t *testing.T, attr schema.Attribute) bool {
	t.Helper()

	ctx := context.Background()

	switch a := attr.(type) {
	case *schema.StringAttribute:
		if a.Default != nil {
			return true
		}
		req := planmodifier.StringRequest{
			StateValue:  types.StringValue("prior"),
			PlanValue:   types.StringUnknown(),
			ConfigValue: types.StringNull(),
		}
		for _, m := range a.PlanModifiers {
			resp := planmodifier.StringResponse{PlanValue: req.PlanValue}
			m.PlanModifyString(ctx, req, &resp)
			req.PlanValue = resp.PlanValue
		}
		return req.PlanValue.Equal(types.StringValue("prior"))

	case *schema.Int64Attribute:
		if a.Default != nil {
			return true
		}
		req := planmodifier.Int64Request{
			StateValue:  types.Int64Value(42),
			PlanValue:   types.Int64Unknown(),
			ConfigValue: types.Int64Null(),
		}
		for _, m := range a.PlanModifiers {
			resp := planmodifier.Int64Response{PlanValue: req.PlanValue}
			m.PlanModifyInt64(ctx, req, &resp)
			req.PlanValue = resp.PlanValue
		}
		return req.PlanValue.Equal(types.Int64Value(42))

	case *schema.BoolAttribute:
		if a.Default != nil {
			return true
		}
		req := planmodifier.BoolRequest{
			StateValue:  types.BoolValue(true),
			PlanValue:   types.BoolUnknown(),
			ConfigValue: types.BoolNull(),
		}
		for _, m := range a.PlanModifiers {
			resp := planmodifier.BoolResponse{PlanValue: req.PlanValue}
			m.PlanModifyBool(ctx, req, &resp)
			req.PlanValue = resp.PlanValue
		}
		return req.PlanValue.Equal(types.BoolValue(true))

	case *schema.SetAttribute:
		if a.Default != nil {
			return true
		}
		return setResolvesUnknownToState(ctx, a.ElementType, a.PlanModifiers)

	case *schema.SetNestedAttribute:
		if a.Default != nil {
			return true
		}
		return setResolvesUnknownToState(ctx, a.NestedObject.Type(), a.PlanModifiers)

	default:
		t.Fatalf("unhandled attribute type %T - extend resolvesUnknownToState", attr)
		return false
	}
}

// setResolvesUnknownToState is the set case of resolvesUnknownToState. Prior state is an
// empty set, which is representable for any element type including nested objects.
func setResolvesUnknownToState(ctx context.Context, elementType attr.Type, modifiers []planmodifier.Set) bool {
	priorState := types.SetValueMust(elementType, nil)

	req := planmodifier.SetRequest{
		StateValue:  priorState,
		PlanValue:   types.SetUnknown(elementType),
		ConfigValue: types.SetNull(elementType),
	}
	for _, m := range modifiers {
		resp := planmodifier.SetResponse{PlanValue: req.PlanValue}
		m.PlanModifySet(ctx, req, &resp)
		req.PlanValue = resp.PlanValue
	}

	return req.PlanValue.Equal(priorState)
}
