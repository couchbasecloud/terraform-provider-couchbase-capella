package datasources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestClusterDatasourceZonesIsList pins the declared type of the zones attributes on both
// cluster data sources.
//
// This needs a schema test rather than an acceptance assertion: terraform-plugin-testing
// flattens lists and sets alike by numeric index, so TestCheckResourceAttr on "zones.0" and
// TestCheckTypeSetElemAttr on "zones.*" both pass whichever type is declared. The type is
// only observable to a practitioner through whether `zones[0]` parses, which is a
// configuration-time property of the schema.
//
// AV-149541: AV-147905 briefly made these sets, which broke configurations indexing them -
// and neither a length guard nor try() suppresses that, because indexing a set is rejected
// when the configuration is type-checked rather than when it is evaluated. They are lists
// again, and NewZonesList sorts, so the ordering churn a list would otherwise risk cannot
// happen. See TestNewZonesListIsSorted in internal/schema for the other half of that pair:
// the type here is only safe because the order there is deterministic.
func TestClusterDatasourceZonesIsList(t *testing.T) {
	assertStringList := func(t *testing.T, attribute schema.Attribute, where string) {
		t.Helper()

		zones, ok := attribute.(schema.ListAttribute)
		if !ok {
			pointer, isPointer := attribute.(*schema.ListAttribute)
			if !isPointer {
				t.Fatalf("%s must be a ListAttribute so configurations can keep indexing it "+
					"(AV-149541), got %T", where, attribute)
			}
			zones = *pointer
		}

		if zones.ElementType != types.StringType {
			t.Fatalf("%s element type is %v, want string", where, zones.ElementType)
		}
		if !zones.Computed {
			t.Fatalf("%s must be computed on a data source", where)
		}
	}

	t.Run("couchbase-capella_cluster", func(t *testing.T) {
		attribute, ok := OneClusterSchema().Attributes["zones"]
		if !ok {
			t.Fatal("zones attribute is missing from the cluster data source schema")
		}
		assertStringList(t, attribute, "zones")
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
		assertStringList(t, attribute, "data[].zones")
	})
}
