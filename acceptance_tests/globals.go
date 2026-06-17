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
	globalProjectId      string
	globalClusterId      string
	globalClusterName    = "tf_acc_test_cluster_common"
	globalBucketName     = "default"
	globalScopeName      = "_default"
	globalCollectionName = "_default"
	globalBucketId       string

	// globalMetadataBucketName is a second bucket on the global cluster used as the metadata
	// storage keyspace for eventing function tests. It must be a different keyspace from the
	// event source (the global bucket's _default/_default).
	globalMetadataBucketName = "metadata"
	globalMetadataBucketId   string

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

	// Fixture buckets for app_endpoint resource tests. Each test gets its
	// own bucket because Capella only permits one endpoint per bucket/scope/collection.
	// These buckets are not created or managed by Terraform; instead, test helpers
	// manage their lifecycle through the API during test runs, including cleanup
	// when needed.
	globalEPBucketName                   = "tf_acc_ep_bkt"
	globalNoCorsEPBucketName             = "tf_acc_ep_nocors_bkt"
	globalCorsFullEPBucketName           = "tf_acc_ep_cors_full_bkt"
	globalCorsSpecificEPBucketName       = "tf_acc_ep_cors_spec_bkt"
	globalCorsMaxAge0EPBucketName        = "tf_acc_ep_cors_ma0_bkt"
	globalOIDCFullEPBucketName           = "tf_acc_ep_oidc_full_bkt"
	globalOIDCDiscEPBucketName           = "tf_acc_ep_oidc_disc_bkt"
	globalCorsExpandEPBucketName         = "tf_acc_ep_cors_exp_bkt"
	globalCorsWildEPBucketName           = "tf_acc_ep_cors_wld_bkt"
	globalAddOIDCEPBucketName            = "tf_acc_ep_add_oidc_bkt"
	globalACFUpdateEPBucketName          = "tf_acc_ep_acf_bkt"
	globalCorsMaxAgeZeroEPBucketName     = "tf_acc_ep_maz_bkt"
	globalCorsMaxAgeFromZeroEPBucketName = "tf_acc_ep_mafz_bkt"
	// Buckets for currently-skipped tests — not provisioned while skipped, but
	// named so the tests are ready to run once the underlying bugs are fixed.
	globalCorsDisabledFalseEPBucketName = "tf_acc_ep_cors_df_bkt"
	globalMultipleOIDCEPBucketName      = "tf_acc_ep_moidc_bkt"
	globalRemoveCorsEPBucketName        = "tf_acc_ep_rmcors_bkt"
	globalCorsDisableToggleEPBucketName = "tf_acc_ep_cdtgl_bkt"
	globalRemoveOIDCEPBucketName        = "tf_acc_ep_rmoidc_bkt"
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
