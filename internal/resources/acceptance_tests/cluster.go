package acceptance_tests

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

func getCIDR(ctx context.Context, client api.Client, CSP string) (string, error) {
	hostName := ""
	switch {
	case strings.Contains(Host, "localhost"):
		hostName = "http://localhost:8080"
	case strings.Contains(Host, "cloudapi"):
		hostName = strings.Replace(Host, "cloudapi", "api", 1)
	default:
		const msg = "unknown host"
		log.Print(msg, Host)
		return "", ErrUnknownHost
	}

	jwt, err := getJWT(ctx, client, hostName)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"%s/v2/organizations/%s/clusters/deployment-options?provider=%s",
		hostName,
		OrgId,
		CSP,
	)

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+jwt)

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	options := struct {
		SuggestedCidr string `json:"suggestedCidr"`
	}{}
	if err = json.Unmarshal(body, &options); err != nil {
		return "", err
	}

	if options.SuggestedCidr == "" {
		const msg = "no CIDR"
		log.Print(msg, string(body))
		return "", ErrNoCIDR
	}
	return options.SuggestedCidr, nil
}

func getJWT(ctx context.Context, client api.Client, hostName string) (string, error) {
	url := hostName + "/sessions"

	authToken := createBasicAuthToken(Username, Password)

	request, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("Authorization", "Basic "+authToken)

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	res := struct {
		Jwt string `json:"jwt"`
	}{}
	if err = json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	if res.Jwt == "" {
		const msg = "no JWT token"
		log.Print(msg, string(body))
		return "", ErrNoJWT
	}

	return res.Jwt, nil
}

func createBasicAuthToken(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func wait(ctx context.Context, client api.Client) error {
	const maxWaitTime = 60 * time.Minute

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, maxWaitTime)
	defer cancel()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ErrTimeoutWaitingForCluster
		case <-ticker.C:
			url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", Host, OrgId, ProjectId, ClusterId)
			cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
			response, err := client.ExecuteWithRetry(
				ctx,
				cfg,
				nil,
				Token,
				nil,
			)
			if err != nil {
				return err
			}

			clusterResp := clusterapi.GetClusterResponse{}
			err = json.Unmarshal(response.Body, &clusterResp)
			if err != nil {
				return err
			}

			if clusterapi.IsFinalState(clusterResp.CurrentState) {
				return nil
			}
		}
	}
}
