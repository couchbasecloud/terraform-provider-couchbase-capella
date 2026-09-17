package resources

import (
	"strconv"
	"testing"

	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestMorphToApiServiceGroupsAzureDisk locks down how morphToApiServiceGroups (cluster.go)
// converts the Azure disk block's Optional+Computed fields (autoexpansion, storage, iops)
// into the create/update request body. An unconfigured (unknown) value must be omitted from
// the request so the server-side default/prior value applies; guards that only check IsNull()
// unwrap an unknown to its zero value and send it explicitly instead. See AV-139797 (Azure
// autoexpansion default), AV-143040 (autoexpansion sent as false on create when omitted) and
// AV-143042 (iops sent as 0 on Ultra disks when only storage, not iops, is unknown).
func TestMorphToApiServiceGroupsAzureDisk(t *testing.T) {
	trueVal := true
	falseVal := false
	storageVal := 100
	iopsVal := 5000

	cases := []struct {
		name          string
		diskType      string
		autoexpansion types.Bool
		storage       types.Int64
		iops          types.Int64
		wantAutoexp   *bool
		wantStorage   *int
		wantIops      *int
	}{
		{
			name:          "autoexpansion unknown (omitted on create) must be omitted, not sent as false",
			diskType:      "P6",
			autoexpansion: types.BoolUnknown(),
			storage:       types.Int64Null(),
			iops:          types.Int64Null(),
			wantAutoexp:   nil,
		},
		{
			name:          "autoexpansion null must be omitted",
			diskType:      "P6",
			autoexpansion: types.BoolNull(),
			storage:       types.Int64Null(),
			iops:          types.Int64Null(),
			wantAutoexp:   nil,
		},
		{
			name:          "autoexpansion explicit true must be sent as true",
			diskType:      "P6",
			autoexpansion: types.BoolValue(true),
			storage:       types.Int64Null(),
			iops:          types.Int64Null(),
			wantAutoexp:   &trueVal,
		},
		{
			name:          "autoexpansion explicit false must be sent as false",
			diskType:      "P6",
			autoexpansion: types.BoolValue(false),
			storage:       types.Int64Null(),
			iops:          types.Int64Null(),
			wantAutoexp:   &falseVal,
		},
		{
			name:          "Ultra disk: storage unknown must be omitted",
			diskType:      "Ultra",
			autoexpansion: types.BoolNull(),
			storage:       types.Int64Unknown(),
			iops:          types.Int64Value(int64(iopsVal)),
			wantStorage:   nil,
			wantIops:      &iopsVal,
		},
		{
			name:          "Ultra disk: iops unknown must be omitted even when storage is known",
			diskType:      "Ultra",
			autoexpansion: types.BoolNull(),
			storage:       types.Int64Value(int64(storageVal)),
			iops:          types.Int64Unknown(),
			wantStorage:   &storageVal,
			wantIops:      nil,
		},
		{
			name:          "Ultra disk: storage and iops both explicit must both be sent",
			diskType:      "Ultra",
			autoexpansion: types.BoolNull(),
			storage:       types.Int64Value(int64(storageVal)),
			iops:          types.Int64Value(int64(iopsVal)),
			wantStorage:   &storageVal,
			wantIops:      &iopsVal,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan := providerschema.Cluster{
				CloudProvider: &providerschema.CloudProvider{
					Type:   types.StringValue(string(clusterapi.Azure)),
					Region: types.StringValue("eastus"),
					Cidr:   types.StringValue("10.0.0.0/23"),
				},
				ServiceGroups: []providerschema.ServiceGroup{
					{
						Node: &providerschema.Node{
							Compute: providerschema.Compute{
								Cpu: types.Int64Value(4),
								Ram: types.Int64Value(16),
							},
							Disk: providerschema.Node_Disk{
								Type:          types.StringValue(tc.diskType),
								Storage:       tc.storage,
								IOPS:          tc.iops,
								Autoexpansion: tc.autoexpansion,
							},
						},
						NumOfNodes: types.Int64Value(3),
						Services:   []types.String{types.StringValue("data")},
					},
				},
			}

			c := &Cluster{}
			serviceGroups, err := c.morphToApiServiceGroups(plan)
			if err != nil {
				t.Fatalf("morphToApiServiceGroups returned error: %v", err)
			}

			disk, err := serviceGroups[0].Node.AsDiskAzure()
			if err != nil {
				t.Fatalf("AsDiskAzure returned error: %v", err)
			}

			if !boolPtrEqual(disk.Autoexpansion, tc.wantAutoexp) {
				t.Errorf("Autoexpansion = %s, want %s", boolPtrString(disk.Autoexpansion), boolPtrString(tc.wantAutoexp))
			}
			if !intPtrEqual(disk.Storage, tc.wantStorage) {
				t.Errorf("Storage = %s, want %s", intPtrString(disk.Storage), intPtrString(tc.wantStorage))
			}
			if !intPtrEqual(disk.Iops, tc.wantIops) {
				t.Errorf("Iops = %s, want %s", intPtrString(disk.Iops), intPtrString(tc.wantIops))
			}
		})
	}
}

func boolPtrEqual(a, b *bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func intPtrEqual(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func boolPtrString(p *bool) string {
	if p == nil {
		return "<nil>"
	}
	return strconv.FormatBool(*p)
}

func intPtrString(p *int) string {
	if p == nil {
		return "<nil>"
	}
	return strconv.Itoa(*p)
}
