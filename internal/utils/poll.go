package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	internalerrors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

const (
	maxDuration = 20
	timeout     = time.Minute * 60 // 60 min is arbitrary
)

func PollIndex(
	ctx context.Context, exec func() (response *api.Response, err error),
) error {

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	attempt := 0
	timer := time.NewTimer((1 << attempt) * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return internalerrors.ErrLongIndexBuildTime
		case <-timer.C:
			response, err := exec()
			if err != nil {
				return err
			}

			status := api.IndexBuildStatusResponse{}
			if err = json.Unmarshal(response.Body, &status); err != nil {
				return err
			}

			if status.Status == "Ready" {
				return nil
			}

			attempt++
			// exponential backoff upto a max of 20 min.
			d := min(maxDuration, 1<<attempt)
			timer.Reset(time.Duration(d) * time.Minute)
		}

	}
}
