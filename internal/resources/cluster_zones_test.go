package resources

import (
	"context"
	"errors"
	"testing"

	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providererrors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func zonesSet(zones ...string) types.Set {
	return providerschema.NewZonesSet(zones)
}

func emptyZones() types.Set {
	return types.SetValueMust(types.StringType, []attr.Value{})
}

func TestZonesForState(t *testing.T) {
	apiZones := zonesSet("use1-az1")
	tests := []struct {
		name       string
		apiZones   types.Set
		knownZones types.Set
		want       types.Set
	}{
		{name: "API zones replace null", apiZones: apiZones, knownZones: types.SetNull(types.StringType), want: apiZones},
		{name: "API zones replace unknown", apiZones: apiZones, knownZones: types.SetUnknown(types.StringType), want: apiZones},
		{name: "API zones match configured", apiZones: apiZones, knownZones: zonesSet("use1-az1"), want: apiZones},
		{name: "API omits zones, keep configured", apiZones: types.SetNull(types.StringType), knownZones: zonesSet("use1-az1"), want: zonesSet("use1-az1")},
		{name: "API omits zones, unknown stays null", apiZones: types.SetNull(types.StringType), knownZones: types.SetUnknown(types.StringType), want: types.SetNull(types.StringType)},
		{name: "configured empty is kept", apiZones: apiZones, knownZones: emptyZones(), want: emptyZones()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := zonesForState(tt.apiZones, tt.knownZones); !got.Equal(tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestValidateZones(t *testing.T) {
	aws := types.StringValue(string(clusterapi.Aws))
	single := types.StringValue(string(clusterapi.Single))
	multi := types.StringValue(string(clusterapi.Multi))

	tests := []struct {
		name              string
		cloudProviderType types.String
		availabilityType  types.String
		zones             types.Set
		wantErr           bool
	}{
		{name: "null zones", cloudProviderType: types.StringValue("gcp"), availabilityType: multi, zones: types.SetNull(types.StringType)},
		{name: "unknown zones", cloudProviderType: aws, availabilityType: multi, zones: types.SetUnknown(types.StringType)},
		{name: "gcp with zones", cloudProviderType: types.StringValue(string(clusterapi.Gcp)), availabilityType: single, zones: zonesSet("us-east1-b"), wantErr: true},
		{name: "azure with zones", cloudProviderType: types.StringValue(string(clusterapi.Azure)), availabilityType: multi, zones: zonesSet("1", "2", "3"), wantErr: true},
		{name: "empty zones on single", cloudProviderType: aws, availabilityType: single, zones: emptyZones()},
		{name: "empty zones on multi", cloudProviderType: aws, availabilityType: multi, zones: emptyZones()},
		{name: "empty zones on gcp", cloudProviderType: types.StringValue(string(clusterapi.Gcp)), availabilityType: multi, zones: emptyZones()},
		{name: "single with 1", cloudProviderType: aws, availabilityType: single, zones: zonesSet("use1-az1")},
		{name: "single with 2", cloudProviderType: aws, availabilityType: single, zones: zonesSet("use1-az1", "use1-az2"), wantErr: true},
		{name: "multi with 2", cloudProviderType: aws, availabilityType: multi, zones: zonesSet("use1-az1", "use1-az2"), wantErr: true},
		{name: "multi with 3", cloudProviderType: aws, availabilityType: multi, zones: zonesSet("use1-az1", "use1-az2", "use1-az4")},
		{name: "multi with 4", cloudProviderType: aws, availabilityType: multi, zones: zonesSet("use1-az1", "use1-az2", "use1-az4", "use1-az5"), wantErr: true},
		{name: "unknown cloud provider", cloudProviderType: types.StringUnknown(), availabilityType: multi, zones: zonesSet("use1-az1", "use1-az2", "use1-az4")},
		{name: "unknown availability", cloudProviderType: aws, availabilityType: types.StringUnknown(), zones: zonesSet("use1-az1", "use1-az2")},
		{
			name:              "unknown element",
			cloudProviderType: aws,
			availabilityType:  multi,
			zones:             types.SetValueMust(types.StringType, []attr.Value{types.StringValue("use1-az1"), types.StringUnknown()}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateZones(tt.cloudProviderType, tt.availabilityType, tt.zones)
			if tt.wantErr {
				if !errors.Is(err, providererrors.ErrInvalidZones) {
					t.Fatalf("expected ErrInvalidZones, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidateCreateClusterZones(t *testing.T) {
	c := &Cluster{}

	plan := awsClusterPlan(types.Int64Value(50), types.Int64Value(3000))
	plan.Availability = &providerschema.Availability{Type: types.StringValue(string(clusterapi.Multi))}

	plan.Zones = zonesSet("use1-az1", "use1-az2", "use1-az4")
	if err := c.validateCreateCluster(plan); err != nil {
		t.Fatalf("expected 3 zones for multi AZ to be valid, got %v", err)
	}

	plan.Zones = zonesSet("use1-az1", "use1-az2")
	if err := c.validateCreateCluster(plan); !errors.Is(err, providererrors.ErrInvalidZones) {
		t.Fatalf("expected ErrInvalidZones for 2 zones on multi AZ, got %v", err)
	}
}

// TestClusterSchemaZones runs the zones plan modifiers the way the framework does on an
// update, and checks a cluster is only replaced when the user configured different zones.
// An update marks the computed zones unknown, so this guards against replacing clusters whose
// zones were chosen by Capella, or are not reported by the API at all (state is null).
func TestClusterSchemaZones(t *testing.T) {
	ctx := context.Background()
	zones, ok := ClusterSchema().Attributes["zones"].(*schema.SetAttribute)
	if !ok {
		t.Fatalf("zones is not a set attribute")
	}
	if !zones.Optional || !zones.Computed {
		t.Fatalf("zones must be optional and computed, got optional=%v computed=%v", zones.Optional, zones.Computed)
	}

	existingResource := tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{})
	tests := []struct {
		name        string
		config      types.Set
		state       types.Set
		wantReplace bool
	}{
		{name: "not configured, API returns no zones", config: types.SetNull(types.StringType), state: types.SetNull(types.StringType)},
		{name: "not configured, Capella chose zones", config: types.SetNull(types.StringType), state: zonesSet("use1-az4", "use1-az5", "use1-az6")},
		{name: "configured, same zones in another order", config: zonesSet("use1-az2", "use1-az1", "use1-az4"), state: zonesSet("use1-az1", "use1-az2", "use1-az4")},
		{name: "configured, matches zones the old API did not return", config: zonesSet("use1-az1"), state: zonesSet("use1-az1")},
		{name: "configured, different zone", config: zonesSet("use1-az1", "use1-az2", "use1-az5"), state: zonesSet("use1-az1", "use1-az2", "use1-az4"), wantReplace: true},
		{name: "configured empty, state empty", config: emptyZones(), state: emptyZones()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := tt.config
			if plan.IsNull() {
				plan = types.SetUnknown(types.StringType)
			}

			var replace bool
			for _, modifier := range zones.PlanModifiers {
				req := planmodifier.SetRequest{
					Path:        path.Root("zones"),
					ConfigValue: tt.config,
					StateValue:  tt.state,
					PlanValue:   plan,
					State:       tfsdk.State{Raw: existingResource},
					Plan:        tfsdk.Plan{Raw: existingResource},
				}
				resp := &planmodifier.SetResponse{PlanValue: plan}
				modifier.PlanModifySet(ctx, req, resp)
				if resp.Diagnostics.HasError() {
					t.Fatalf("plan modifier returned errors: %v", resp.Diagnostics)
				}
				plan = resp.PlanValue
				replace = replace || resp.RequiresReplace
			}

			if replace != tt.wantReplace {
				t.Fatalf("expected replace=%v, got %v", tt.wantReplace, replace)
			}
		})
	}
}
