package schema

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestNewZonesListIsSorted guards the property that makes exposing zones as a list safe.
//
// Capella returns availability zones in an order the provider does not control, so a list
// built in API order would reorder between reads and churn anything that consumes it.
// NewZonesList sorts, which makes the order deterministic and the list type harmless -
// see AV-149541, where the alternative was making the attribute a set and breaking every
// configuration that indexes it.
//
// The input here is deliberately NOT in sorted order. A test that feeds already-sorted
// zones proves nothing, because the sorted and unsorted results are identical - do not
// "tidy" these values.
func TestNewZonesListIsSorted(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		zones []string
		want  []string
	}{
		{
			name:  "reverse order is sorted",
			zones: []string{"use1-az4", "use1-az2", "use1-az1"},
			want:  []string{"use1-az1", "use1-az2", "use1-az4"},
		},
		{
			name:  "arbitrary order is sorted",
			zones: []string{"use1-az2", "use1-az5", "use1-az1"},
			want:  []string{"use1-az1", "use1-az2", "use1-az5"},
		},
		{
			name:  "single zone",
			zones: []string{"use1-az1"},
			want:  []string{"use1-az1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got []string
			assert.False(t, NewZonesList(test.zones).ElementsAs(ctx, &got, false).HasError())
			// Equal, not ElementsMatch: the order is the point of this test.
			assert.Equal(t, test.want, got)
		})
	}

	t.Run("two API orderings give the same list", func(t *testing.T) {
		first := NewZonesList([]string{"use1-az4", "use1-az1", "use1-az2"})
		second := NewZonesList([]string{"use1-az2", "use1-az4", "use1-az1"})
		assert.True(t, first.Equal(second),
			"the same zones in a different API order must produce the same list, "+
				"otherwise a refresh churns every consumer of the attribute")
	})

	t.Run("absent and empty are null, not an empty list", func(t *testing.T) {
		// Terraform treats a null collection and an empty one as different values, and the
		// data source reports null when Capella does not supply zones at all.
		assert.True(t, NewZonesList(nil).IsNull())
		assert.True(t, NewZonesList([]string{}).IsNull())
		assert.Equal(t, types.StringType, NewZonesList(nil).ElementType(context.Background()))
	})
}
