package resources

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// appEndpointDeletedResult represents the outcome of checking whether an App Endpoint
// was deleted outside of Terraform.
type appEndpointDeletedResult int

const (
	// appEndpointDeleted indicates the App Endpoint no longer exists and should be removed from state.
	appEndpointDeleted appEndpointDeletedResult = iota
	// appEndpointExists indicates the App Endpoint still exists (possible race condition).
	appEndpointExists
	// appEndpointCheckFailed indicates the check itself failed and the error should be returned to the user.
	appEndpointCheckFailed
)

// checkAppEndpointDeletedOrForbidden determines whether a 403 response from a Get App Endpoint
// call is caused by the App Endpoint being deleted outside of Terraform, or by a genuine
// permission error.
//
// It does this by listing all App Endpoints for the App Service (using paginated fetches)
// and checking whether the target endpoint appears in any page. If the list call itself
// returns 403, the original error is treated as a genuine permission issue.
func checkAppEndpointDeletedOrForbidden(
	ctx context.Context,
	data *providerschema.Data,
	organizationId, projectId, clusterId, appServiceId, appEndpointName string,
) (appEndpointDeletedResult, string) {
	listURL := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints",
		data.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
	)

	cfg := api.EndpointCfg{
		Url:           listURL,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	allEndpoints, err := api.GetPaginated[[]app_endpoints.GetAppEndpointResponse](
		ctx,
		data.ClientV1,
		data.Token,
		cfg,
		api.SortByName,
	)
	if err != nil {
		var apiError *api.Error
		if errors.As(err, &apiError) && apiError.HttpStatusCode == http.StatusForbidden {
			return appEndpointCheckFailed, "permission denied: unable to access app endpoints for this app service"
		}
		return appEndpointCheckFailed, api.ParseError(err)
	}

	for _, ep := range allEndpoints {
		if ep.Name == appEndpointName {
			tflog.Warn(ctx, "App Endpoint exists in list but Get returned 403; possible race condition, please retry")
			return appEndpointExists, "App Endpoint exists but the Get request returned 403. This may be a transient issue, please retry."
		}
	}

	tflog.Info(ctx, "App Endpoint not found in list response; it has been deleted outside of Terraform")
	return appEndpointDeleted, ""
}

// isForbiddenError checks whether the given error is an api.Error with HTTP status 403.
func isForbiddenError(err error) bool {
	var apiError *api.Error
	return errors.As(err, &apiError) && apiError.HttpStatusCode == http.StatusForbidden
}
