package resources

import (
	"strings"
	"testing"

	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// awsClusterPlan builds a minimal valid AWS cluster plan whose disk storage and iops the
// caller controls, so a test can make either attribute unknown the way an omitted
// Optional+Computed attribute arrives at create time.
func awsClusterPlan(storage, iops types.Int64) providerschema.Cluster {
	return providerschema.Cluster{
		Name:           types.StringValue("tf-acc-aws-disk"),
		OrganizationId: types.StringValue("00000000-0000-0000-0000-000000000001"),
		ProjectId:      types.StringValue("00000000-0000-0000-0000-000000000002"),
		CloudProvider: &providerschema.CloudProvider{
			Type:   types.StringValue(string(clusterapi.Aws)),
			Region: types.StringValue("us-east-1"),
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
						Type:    types.StringValue("gp3"),
						Storage: storage,
						IOPS:    iops,
					},
				},
				NumOfNodes: types.Int64Value(3),
				Services:   []types.String{types.StringValue("data")},
			},
		},
	}
}

// TestMorphToApiServiceGroupsAWSDisk pins the part of the AWS disk mapping that is correct
// and will stay correct however AV-145015 is fixed: a configured storage and iops reach the
// request body unchanged.
//
// It deliberately does not assert what happens to an unknown attribute. Today that produces
// a 0 on the wire, which is the AV-145015 defect - asserting it here would encode the bug as
// the expected contract and turn the eventual fix into a test failure. The intended
// behaviour is stated instead by TestValidateCreateClusterAwsDiskRequired_AV_145015 below.
func TestMorphToApiServiceGroupsAWSDisk(t *testing.T) {
	c := &Cluster{}
	serviceGroups, err := c.morphToApiServiceGroups(awsClusterPlan(types.Int64Value(50), types.Int64Value(3000)))
	if err != nil {
		t.Fatalf("morphToApiServiceGroups returned error: %v", err)
	}

	disk, err := serviceGroups[0].Node.AsDiskAWS()
	if err != nil {
		t.Fatalf("AsDiskAWS returned error: %v", err)
	}

	if disk.Storage != 50 {
		t.Errorf("Storage = %d, want 50", disk.Storage)
	}
	if disk.Iops != 3000 {
		t.Errorf("Iops = %d, want 3000", disk.Iops)
	}
	if string(disk.Type) != "gp3" {
		t.Errorf("Type = %q, want %q", disk.Type, "gp3")
	}
}

// TestValidateCreateClusterAwsDiskRequired_AV_145015 states the behaviour the provider
// should have and does not yet: an AWS cluster config that omits `storage` or `iops` must be
// rejected by the provider, naming the attribute, rather than being sent to Capella with a
// silent 0.
//
// Why rejection is the right contract rather than omitting the field: `storage` and `iops`
// are plain ints with no `omitempty` in the API contract on both sides - the provider's
// generated spec (internal/generated/api/openapi.gen.go, DiskAWS) and the Capella server's
// own oapi types. Leaving them out of the JSON therefore still unmarshals to 0 server-side
// and produces the identical 422. There is no AWS default for Capella to fall back on, so
// these attributes are mandatory in practice despite being declared Optional
// (internal/resources/cluster_schema.go:63-64), and the only place the practitioner can be
// told that clearly is the provider.
//
// The expected shape follows the sibling checks already in validateCreateCluster
// (internal/resources/cluster.go:833-849), which reject `autoexpansion` on AWS/GCP and
// `iops` on GCP.
func TestValidateCreateClusterAwsDiskRequired_AV_145015(t *testing.T) {
	t.Skip("AV-145015: AWS storage/iops are Optional+Computed, so an omitted attribute is " +
		"unknown and gets serialised as 0; Capella then rejects the create with a range error " +
		"naming a value the practitioner never wrote. Unskip once validateCreateCluster " +
		"rejects an AWS config that omits either attribute.")

	cases := []struct {
		name    string
		storage types.Int64
		iops    types.Int64
		wantIn  string
	}{
		{"storage omitted is unknown at create", types.Int64Unknown(), types.Int64Value(3000), "storage"},
		{"iops omitted is unknown at create", types.Int64Value(50), types.Int64Unknown(), "iops"},
		{"storage explicitly null", types.Int64Null(), types.Int64Value(3000), "storage"},
		{"iops explicitly null", types.Int64Value(50), types.Int64Null(), "iops"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Cluster{}
			err := c.validateCreateCluster(awsClusterPlan(tc.storage, tc.iops))
			if err == nil {
				t.Fatalf("validateCreateCluster returned nil, want an error naming %q", tc.wantIn)
			}
			if !strings.Contains(strings.ToLower(err.Error()), tc.wantIn) {
				t.Errorf("validateCreateCluster error = %q, want it to name %q", err, tc.wantIn)
			}
		})
	}

	// Both configured must still pass, so the new check cannot be satisfied by rejecting
	// every AWS cluster.
	t.Run("both configured is accepted", func(t *testing.T) {
		c := &Cluster{}
		if err := c.validateCreateCluster(awsClusterPlan(types.Int64Value(50), types.Int64Value(3000))); err != nil {
			t.Errorf("validateCreateCluster returned %v, want nil", err)
		}
	})
}
