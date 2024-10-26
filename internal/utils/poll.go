package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	internalerrors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

const (
	maxAttempts = 5
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

type WatchOptions struct {
	Host       string
	OrgId      string
	ProjectId  string
	ClusterId  string
	Bucket     string
	Scope      string
	Collection string
}

// TODO: get build status for all indexes in 1 API call
func WatchIndexes(
	expectedState string, indexes []string, exec func(cfg api.EndpointCfg) (response *api.Response, err error),
	options WatchOptions,
) error {
	attempts := 0

	for i := 0; i < len(indexes); {
		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/queryService/indexBuildStatus/%s?bucket=%s&scope=%s&collection=%s",
			options.Host,
			options.OrgId,
			options.ProjectId,
			options.ClusterId,
			url.QueryEscape(indexes[i]),
			options.Bucket,
			options.Scope,
			options.Collection,
		)

		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		response, err := exec(cfg)
		// retry even on 404 because the index may not have been created yet
		// so wait and check again
		if err != nil {
			attempts++
			if attempts > maxAttempts {
				return internalerrors.ErrMaxAttemptsExceeded
			}
			time.Sleep(10 * time.Second)
			continue
		}

		status := api.IndexBuildStatusResponse{}
		if err = json.Unmarshal(response.Body, &status); err != nil {
			return err
		}

		if status.Status != expectedState {
			attempts++
			if attempts > maxAttempts {
				return internalerrors.ErrMaxAttemptsExceeded
			}
			time.Sleep(10 * time.Second)
			continue
		}

		i++
	}

	return nil
}
