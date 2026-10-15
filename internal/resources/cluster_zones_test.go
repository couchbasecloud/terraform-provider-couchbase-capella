package resources

import (
	"context"
	"errors"
	"testing"

	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providererrors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func zonesSet(zones ...string) types.Set {
	return providerschema.NewZonesSet(zones)
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
		{name: "empty zones", cloudProviderType: aws, availabilityType: multi, zones: types.SetValueMust(types.StringType, []attr.Value{}), wantErr: true},
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

// TestClusterSchemaZones guards against a plan that replaces every cluster whose zones were
// chosen by Capella: zones must be computed, and UseStateForUnknown must run before
// RequiresReplace so an unconfigured zones keeps its state value instead of becoming unknown.
func TestClusterSchemaZones(t *testing.T) {
	ctx := context.Background()
	zones, ok := ClusterSchema().Attributes["zones"].(*schema.SetAttribute)
	if !ok {
		t.Fatalf("zones is not a set attribute")
	}
	if !zones.Optional || !zones.Computed {
		t.Fatalf("zones must be optional and computed, got optional=%v computed=%v", zones.Optional, zones.Computed)
	}

	want := []string{
		setplanmodifier.UseStateForUnknown().Description(ctx),
		setplanmodifier.RequiresReplace().Description(ctx),
	}
	if len(zones.PlanModifiers) != len(want) {
		t.Fatalf("expected %d plan modifiers, got %d", len(want), len(zones.PlanModifiers))
	}
	for i, modifier := range zones.PlanModifiers {
		if got := modifier.Description(ctx); got != want[i] {
			t.Fatalf("plan modifier %d: expected %q, got %q", i, want[i], got)
		}
	}
}
