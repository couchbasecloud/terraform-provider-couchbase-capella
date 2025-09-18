package acceptance_tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/apigen"
	"github.com/google/uuid"
)

// cluster is created with enterprise plan as some features require this.
func createCluster(ctx context.Context, clientV2 *apigen.ClientWithResponses) error {
	n := apigen.Node{}
	_ = n.Disk.FromDiskAWS(apigen.DiskAWS{Type: apigen.DiskAWSType("gp3")})

	clusterRequest := apigen.CreateClusterRequest{
		Name:          "tf_acc_test_cluster_common",
		Availability:  apigen.Availability{Type: apigen.AvailabilityType("multi")},
		CloudProvider: apigen.CloudProvider{Region: "us-east-1", Type: apigen.CloudProviderType("aws")},
		ServiceGroups: []apigen.ServiceGroup{
			{
				Node: &apigen.Node{
					Compute: apigen.Compute{Cpu: 4, Ram: 16},
					Disk:    n.Disk,
				},
				Services:   &[]apigen.Service{apigen.Service("data"), apigen.Service("index"), apigen.Service("query")},
				NumOfNodes: func() *int { v := 3; return &v }(),
			},
		},
		Support: apigen.Support{Plan: apigen.SupportPlan("enterprise"), Timezone: func() *apigen.SupportTimezone { v := apigen.SupportTimezone("PT"); return &v }()},
	}

	orgUUID, _ := uuid.Parse(globalOrgId)
	projUUID, _ := uuid.Parse(globalProjectId)

	res, err := clientV2.PostClusterWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), clusterRequest)
	if err != nil {
		return err
	}
	if res.JSON202 == nil {
		return fmt.Errorf("unexpected status: %s", res.Status())
	}
	globalClusterId = res.JSON202.Id.String()
	return nil
}

func destroyCluster(ctx context.Context, clientV2 *apigen.ClientWithResponses) error {
	orgUUID, _ := uuid.Parse(globalOrgId)
	projUUID, _ := uuid.Parse(globalProjectId)
	cluUUID, _ := uuid.Parse(globalClusterId)
	_, err := clientV2.DeleteClusterWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), nil)
	return err
}

func clusterWait(ctx context.Context, clientV2 *apigen.ClientWithResponses, destroy bool) error {
	const maxWaitTime = 30 * time.Minute
	const checkInterval = 1 * time.Minute

	deadline := time.Now().Add(maxWaitTime)

	orgUUID, _ := uuid.Parse(globalOrgId)
	projUUID, _ := uuid.Parse(globalProjectId)
	cluUUID, _ := uuid.Parse(globalClusterId)

	for time.Now().Before(deadline) {
		res, err := clientV2.GetClusterWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID))
		if err != nil {
			return err
		}

		if destroy {
			if res.StatusCode() == 404 {
				log.Print("cluster destroyed")
				return nil
			}
		} else {
			if res.JSON200 != nil && res.JSON200.CurrentState == apigen.CurrentState("healthy") {
				log.Print("cluster created")
				return nil
			}
		}

		time.Sleep(checkInterval)
	}

	return errors.New("timeout waiting for cluster to be created or destroyed")
}
