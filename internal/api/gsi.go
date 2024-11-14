package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	internalerrors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// IndexDDLRequest is a request to run an index DDL statement.
type IndexDDLRequest struct {
	// Definition The index DDL statement.  This can be a CREATE/DROP/ALTER/BUILD statement.
	// Multiple delimited queries are not allowed.
	Definition string `json:"definition"`
}

// QueryError is the error message returned by query service.
type QueryError struct {
	// Msg The error message.
	Msg string `json:"msg"`
}

// IndexDDLResponse has an array of errors returned by query service.
type IndexDDLResponse struct {
	Errors []QueryError `json:"errors,omitempty"`
}

// IndexDefinitionResponse represents a single index definition.
type IndexDefinitionResponse struct {
	Bucket       string   `json:"bucket"`
	Scope        string   `json:"scope"`
	Collection   string   `json:"collection"`
	IsPrimary    bool     `json:"is_primary"`
	IndexName    string   `json:"indexName"`
	SecExprs     []string `json:"secExprs"` // these are index keys
	PartitionBy  string   `json:"partition_by"`
	Where        string   `json:"where"`
	NumReplica   int      `json:"numReplica"`
	NumPartition int      `json:"numPartition"`
}

type IndexDefinition struct {
	IndexName  string `json:"indexName"`
	Definition string `json:"definition"`
}

// ListIndexDefinitionsResponse represents a list of index definitions.
type ListIndexDefinitionsResponse struct {
	Definitions []IndexDefinition `json:"definitions"`
}

// IndexBuildStatusResponse is the build status for an index.
type IndexBuildStatusResponse struct {
	Status string `json:"status"`
}

// TODO: get build status for all indexes in 1 API call

type Options struct {
	Host       string
	OrgId      string
	ProjectId  string
	ClusterId  string
	Bucket     string
	Scope      string
	Collection string
}

func WatchIndexes(
	ctx context.Context, desiredState string, indexes []string,
	exec func(cfg EndpointCfg) (response *Response, err error),
	options Options,
) error {

	const (
		maxDuration = 20
		timeout     = time.Minute * 60 // 60 min is arbitrary
	)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	attempt := 0
	timer := time.NewTimer((1 << attempt) * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return internalerrors.ErrMonitorTimeout
		case <-timer.C:
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

				cfg := EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
				response, err := exec(cfg)
				// retry even on 404 because the index may not have been created yet.
				// so wait and check again.
				if err != nil {
					// exponential backoff upto a max of 20 min.
					d := min(maxDuration, 1<<attempt)
					timer.Reset(time.Duration(d) * time.Minute)
					attempt++
					continue
				}

				status := IndexBuildStatusResponse{}
				if err = json.Unmarshal(response.Body, &status); err != nil {
					return err
				}

				if status.Status != desiredState {
					// exponential backoff upto a max of 20 min.
					d := min(maxDuration, 1<<attempt)
					timer.Reset(time.Duration(d) * time.Minute)
					attempt++
					continue
				}

				i++
			}

			// all indexes are created/ready
			return nil
		}
	}
}
