package acceptance_tests

import (
	"context"
	"fmt"
	"log"
	"time"

	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/apigen"
	"github.com/google/uuid"
)

func createBucket(ctx context.Context, clientV2 *apigen.ClientWithResponses) error {
	bucketRequest := apigen.CreateBucketRequest{
		Name: globalBucketName,
	}

	orgUUID, _ := uuid.Parse(globalOrgId)
	projUUID, _ := uuid.Parse(globalProjectId)
	cluUUID, _ := uuid.Parse(globalClusterId)

	res, err := clientV2.PostBucketWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), bucketRequest)
	if err != nil {
		return err
	}
	if res.JSON201 == nil {
		return fmt.Errorf("unexpected status: %s", res.Status())
	}
	globalBucketId = res.JSON201.Id
	return nil
}

func bucketWait(ctx context.Context, clientV2 *apigen.ClientWithResponses) error {
	const maxWaitTime = 5 * time.Minute

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, maxWaitTime)
	defer cancel()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	orgUUID, _ := uuid.Parse(globalOrgId)
	projUUID, _ := uuid.Parse(globalProjectId)
	cluUUID, _ := uuid.Parse(globalClusterId)

	for {
		select {
		case <-ctx.Done():
			return ErrTimeoutWaitingForBucket
		case <-ticker.C:
			res, err := clientV2.GetBucketByIDWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), apigen.BucketId(globalBucketId))
			if err != nil {
				return err
			}
			if res.JSON200 != nil {
				log.Print("bucket created")
				return nil
			}
			if res.StatusCode() != 404 {
				return fmt.Errorf("unexpected status while waiting bucket: %s", res.Status())
			}
		}
	}
}
