package datasources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestClusterDatasourceZonesIsSet pins the declared type of the zones attributes.
//
// This needs a schema test rather than an acceptance assertion: terraform-plugin-testing
// flattens lists and sets alike by numeric index, so TestCheckResourceAttr on "zones.0" and
// TestCheckTypeSetElemAttr on "zones.*" both pass whichever type is declared. The type is
// only observable to a practitioner through whether `zones[0]` parses, which is a
// configuration-time property of the schema.
//
// AV-149541: zones on the singular cluster data source changed from a list to a set with
// AV-147905, which is correct - Capella returns these in an order the provider does not
// control - but it breaks configurations that index the attribute. Changing it back is a
// breaking change in the other direction, so it should only happen deliberately.
func TestClusterDatasourceZonesIsSet(t *testing.T) {
	t.Run("couchbase-capella_cluster", func(t *testing.T) {
		attribute, ok := OneClusterSchema().Attributes["zones"]
		if !ok {
			t.Fatal("zones attribute is missing from the cluster data source schema")
		}

		zones, ok := attribute.(schema.SetAttribute)
		if !ok {
			pointer, isPointer := attribute.(*schema.SetAttribute)
			if !isPointer {
				t.Fatalf("zones must be a SetAttribute so a reordering API response does not "+
					"churn the plan, got %T", attribute)
			}
			zones = *pointer
		}

		if zones.ElementType != types.StringType {
			t.Fatalf("zones element type is %v, want string", zones.ElementType)
		}
		if !zones.Computed {
			t.Fatal("zones must be computed on a data source")
		}
	})

	t.Run("couchbase-capella_clusters", func(t *testing.T) {
		data, ok := ClusterSchema().Attributes["data"]
		if !ok {
			t.Fatal("data attribute is missing from the clusters data source schema")
		}

		nested, ok := data.(schema.ListNestedAttribute)
		if !ok {
			pointer, isPointer := data.(*schema.ListNestedAttribute)
			if !isPointer {
				t.Fatalf("data is %T, expected a ListNestedAttribute", data)
			}
			nested = *pointer
		}

		attribute, ok := nested.NestedObject.Attributes["zones"]
		if !ok {
			t.Fatal("zones attribute is missing from the clusters data source schema")
		}

		if _, isSet := attribute.(schema.SetAttribute); !isSet {
			if _, isPointer := attribute.(*schema.SetAttribute); !isPointer {
				t.Fatalf("data[].zones must be a SetAttribute, got %T", attribute)
			}
		}
	})
}
