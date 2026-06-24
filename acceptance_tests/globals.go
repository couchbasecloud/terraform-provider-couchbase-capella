package acceptance_tests

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/provider"
)

var (
	// globalClient is initialised in TestMain and used by lazy endpoint setup.
	globalClient *api.Client

	// these global variables are set by env vars.
	globalHost  string
	globalToken string
	globalOrgId string

	// these global variables are set by setup().
	globalProjectId          string
	globalClusterId          string
	globalClusterName        = "tf_acc_test_cluster_common"
	globalBucketName         = "default"
	globalScopeName          = "_default"
	globalCollectionName     = "_default"
	globalBucketId           string
	globalBucketCreated      bool
	globalAppServiceId       string
	globalAppServiceName     = "tf_acc_test_app_service_common"
	globalAppEndpointName    = "tf_acc_test_app_endpoint_common"
	globalAppEndpointCreated bool
	globalProjectCreated     bool
	globalClusterCreated     bool
	globalAppServiceCreated  bool

	// dmClusterId is a dedicated cluster for data management (bucket/scope/collection) tests,
	// kept separate from the global cluster to avoid rebalance contention with app service tests.
	dmClusterId      string
	dmClusterCreated bool
	dmClusterName    = "tf_acc_test_cluster_dm"
	dmBucketId       string
	dmBucketCreated  bool
	dmBucketName     = "default"

	appEndpointClusterId          string
	appEndpointClusterName        string
	appEndpointBucketId           string
	appEndpointBucketName         = "default"
	appEndpointAppServiceId       string
	appEndpointAppServiceName     string
	appEndpointCommonEndpointName = "tf_acc_test_app_endpoint_common"
	appEndpointClusterCreated     bool
	appEndpointBucketCreated      bool
	appEndpointAppServiceCreated  bool
	appEndpointCreated            bool

	// Pre-created endpoints for sub-resource tests. One dedicated endpoint per
	// test prevents write conflicts without concurrent creation 500s — they are
	// created lazily and sequentially on first use via ensureXxxEndpoint().
	// Each endpoint uses its own bucket because Capella only allows one endpoint
	// per bucket/scope/collection combination.
	globalACFEndpointName             = "tf_acc_test_acf_endpoint"
	globalACFBucketName               = "tf_acc_acf_bkt"
	globalIFEndpointName              = "tf_acc_test_if_endpoint"
	globalIFBucketName                = "tf_acc_if_bkt"
	globalCORSEndpointName            = "tf_acc_test_cors_endpoint"
	globalCORSBucketName              = "tf_acc_cors_bkt"
	globalCORSOriginOnlyEndpointName  = "tf_acc_test_cors_ori_endpoint"
	globalCORSOriginOnlyBucketName    = "tf_acc_cors_ori_bkt"
	globalOIDCEndpointName            = "tf_acc_test_oidc_endpoint"
	globalOIDCBucketName              = "tf_acc_oidc_bkt"
	globalDefaultOIDCEndpointName     = "tf_acc_test_doidc_endpoint"
	globalDefaultOIDCBucketName       = "tf_acc_doidc_bkt"
	appEndpointActivationEndpointName = "tf_acc_test_activation_endpoint"
	appEndpointActivationBucketName   = "tf_acc_activation_bkt"
	appEndpointLoggingEndpointName    = "tf_acc_test_logging_endpoint"
	appEndpointLoggingBucketName      = "tf_acc_logging_bkt"

	// Fixture collections for app_endpoint resource tests. Each test gets its
	// own collection in the shared appEndpointBucketName bucket (under the
	// _default scope) rather than a dedicated bucket: Capella permits only one
	// endpoint per bucket/scope/collection, and a cluster's bucket cap (1 per
	// 0.2 cores, so 20 here) is far smaller than its collection cap. Sharing the
	// bucket keeps parallel runs well under the bucket limit that previously
	// caused flaky "maximum number of buckets reached" failures. These
	// collections are created via the API by ensureFixtureCollection and torn
	// down with the shared bucket at the end of the suite.
	globalEPCollectionName                   = "tf_acc_ep_col"
	globalNoCorsEPCollectionName             = "tf_acc_ep_nocors_col"
	globalCorsFullEPCollectionName           = "tf_acc_ep_cors_full_col"
	globalCorsSpecificEPCollectionName       = "tf_acc_ep_cors_spec_col"
	globalCorsMaxAge0EPCollectionName        = "tf_acc_ep_cors_ma0_col"
	globalOIDCFullEPCollectionName           = "tf_acc_ep_oidc_full_col"
	globalOIDCDiscEPCollectionName           = "tf_acc_ep_oidc_disc_col"
	globalCorsExpandEPCollectionName         = "tf_acc_ep_cors_exp_col"
	globalCorsWildEPCollectionName           = "tf_acc_ep_cors_wld_col"
	globalAddOIDCEPCollectionName            = "tf_acc_ep_add_oidc_col"
	globalACFUpdateEPCollectionName          = "tf_acc_ep_acf_col"
	globalCorsMaxAgeZeroEPCollectionName     = "tf_acc_ep_maz_col"
	globalCorsMaxAgeFromZeroEPCollectionName = "tf_acc_ep_mafz_col"
	// Collections for "deleted externally" tests — validates the 403 → List fallback
	// that removes resources from state when the App Endpoint is deleted outside TF.
	globalDeletedExternallyEPCollectionName = "tf_acc_ep_del_ext_col"
	globalACFDeletedExtEPCollectionName     = "tf_acc_ep_acf_dex_col"
	// Collections for currently-skipped tests — not provisioned while skipped, but
	// named so the tests are ready to run once the underlying bugs are fixed. 
	globalCorsDisabledFalseEPCollectionName = "tf_acc_ep_cors_df_col"
	globalMultipleOIDCEPCollectionName      = "tf_acc_ep_moidc_col"
	globalRemoveCorsEPCollectionName        = "tf_acc_ep_rmcors_col"
	globalCorsDisableToggleEPCollectionName = "tf_acc_ep_cdtgl_col"
	globalRemoveOIDCEPCollectionName        = "tf_acc_ep_rmoidc_col"
	// TestAccAppEndpointInexistentCollection reuses globalBucketName: it tries
	// _default/INVALID_COLLECTION which does not conflict with the common
	// endpoint on _default/_default, so no dedicated bucket is needed.

	// this global variable is set in TestMain.
	globalProviderBlock string

	// globalProtoV6ProviderFactory are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	globalProtoV6ProviderFactory = map[string]func() (tfprotov6.ProviderServer, error){
		"couchbase-capella": providerserver.NewProtocol6WithError(provider.New()()),
	}
)

const (
	timeout = 60 * time.Second
)
