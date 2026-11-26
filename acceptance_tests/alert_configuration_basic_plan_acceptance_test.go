package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

// basicPlanAlertConfigurationMessage is the error Capella returns (422) when alert configuration is
// set on a cluster without the Developer Pro or Enterprise support plan, as seen on 2026-11-26.
// It is matched case-insensitively, and line breaks that Terraform adds when printing are
// tolerated.
const basicPlanAlertConfigurationMessage = "Alert configuration is available for clusters on the Developer Pro or Enterprise support plan."

// TestAccClusterAlertConfigurationBasicPlan covers alert configuration on a cluster with the
// Basic support plan, which cannot configure alerts. It deploys its own single-node cluster,
// because changing the shared cluster's plan would break the tests that need Enterprise, and it
// reuses that cluster for both cases to avoid a second deploy:
//
//  1. a cluster configured on Developer Pro and then downgraded to Basic, which is what an
//     existing customer hits: refresh still works, a change to the configuration is rejected,
//     and destroy still succeeds so the resource can be removed from state.
//  2. configuring alerts on the now-Basic cluster, with each of configs, recipients and
//     muted_alerts on its own, in case Capella gates them separately.
func TestAccClusterAlertConfigurationBasicPlan(t *testing.T) {
	recipient1, recipient2 := alertTestRecipients(t)
	if strings.HasPrefix(basicPlanAlertConfigurationMessage, "TODO") {
		t.Skip("set basicPlanAlertConfigurationMessage to the error Capella returns for alert configuration on a Basic plan cluster")
	}
	basicPlanError := regexp.MustCompile("(?i)" + wrapTolerantPattern(basicPlanAlertConfigurationMessage))

	clusterName := randomStringWithPrefix("tf_acc_alert_basic_")
	alertName := randomStringWithPrefix("tf_acc_alert_basic_cfg_")
	clusterReference := "couchbase-capella_cluster." + clusterName
	alertReference := "couchbase-capella_cluster_alert_configuration." + alertName
	cidr := generateRandomCIDR()
	mutedUntil := alertMuteTime(24 * time.Hour)

	alertBody := `
  configs = {
    cpu_usage_high = {
      warning_threshold = 80
    }
  }

  recipients = ["` + recipient1 + `"]

  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "` + mutedUntil + `"
    }
  }`

	changedAlertBody := strings.Replace(alertBody,
		`["`+recipient1+`"]`, `["`+recipient1+`", "`+recipient2+`"]`, 1)

	// If destroy fails, as it does while removing the alert configuration on a Basic cluster
	// (step 4), Terraform never deletes the cluster. Delete it here so a failed run does not
	// leave a cluster behind.
	var clusterId string
	captureClusterId := func(s *terraform.State) error {
		attributes, err := datasourceAttributes(s, clusterReference)
		if err == nil {
			clusterId = attributes["id"]
		}
		return err
	}
	t.Cleanup(func() { deleteClusterIfPresent(t, clusterId) })

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckClusterDestroyed(t, clusterReference),
		Steps: []resource.TestStep{
			// Configure alerts while the cluster is on a plan that allows it.
			{
				Config: testAccBasicPlanAlertConfig(clusterName, cidr, "developer pro", alertName, alertBody),
				Check: resource.ComposeAggregateTestCheckFunc(
					captureClusterId,
					resource.TestCheckResourceAttr(clusterReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(clusterReference, "current_state", string(clusterapi.Healthy)),
					resource.TestCheckResourceAttr(alertReference, "configs.cpu_usage_high.warning_threshold", "80"),
					resource.TestCheckResourceAttr(alertReference, "recipients.#", "1"),
					resource.TestCheckResourceAttr(alertReference, "muted_alerts.%", "1"),
				),
			},
			// Downgrade to Basic without touching the alert configuration. Refresh must still
			// read it, and the step's post-apply plan check proves the downgrade causes no drift.
			{
				Config: testAccBasicPlanAlertConfig(clusterName, cidr, "basic", alertName, alertBody),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterReference, "support.plan", "basic"),
					resource.TestCheckResourceAttr(alertReference, "configs.cpu_usage_high.warning_threshold", "80"),
					resource.TestCheckResourceAttr(alertReference, "recipients.#", "1"),
					resource.TestCheckResourceAttr(alertReference, "muted_alerts.%", "1"),
				),
			},
			// Changing the configuration on a Basic cluster is rejected.
			{
				Config:      testAccBasicPlanAlertConfig(clusterName, cidr, "basic", alertName, changedAlertBody),
				ExpectError: basicPlanError,
			},
			// Removing the resource must still succeed, or a downgraded customer could never get
			// it out of state. Destroy clears recipients and mutes, which is itself a change.
			{
				Config: testAccBasicPlanClusterConfig(clusterName, cidr, "basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(clusterReference, "support.plan", "basic"),
					testAccCheckNoResourceInState(alertReference),
				),
			},
			// Configuring alerts on a Basic cluster is rejected, for each attribute on its own. Each
			// uses a value Capella does not already hold: the provider sends only differences, and
			// once AV-153495 is fixed destroy leaves the earlier values in place, so reusing them
			// would send nothing and never reach the plan check.
			{
				Config: testAccBasicPlanAlertConfig(clusterName, cidr, "basic", alertName, `
  configs = {
    cpu_usage_high = {
      warning_threshold = 85
    }
  }`),
				ExpectError: basicPlanError,
			},
			{
				Config: testAccBasicPlanAlertConfig(clusterName, cidr, "basic", alertName, `
  recipients = ["`+recipient2+`"]`),
				ExpectError: basicPlanError,
			},
			{
				Config: testAccBasicPlanAlertConfig(clusterName, cidr, "basic", alertName, `
  muted_alerts = {
    cb_high_disk_write_queue = {
      action      = "mute"
      muted_until = "`+mutedUntil+`"
    }
  }`),
				ExpectError: basicPlanError,
			},
			// DIS-07: disabling an alert, a mute with no end, is rejected too.
			{
				Config: testAccBasicPlanAlertConfig(clusterName, cidr, "basic", alertName, `
  muted_alerts = {
    cpu_usage_high = {
      action = "mute"
    }
  }`),
				ExpectError: basicPlanError,
			},
			// A rejected create leaves nothing behind: not in state, and not in Capella.
			{
				Config: testAccBasicPlanClusterConfig(clusterName, cidr, "basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNoResourceInState(alertReference),
					testAccCheckBasicClusterRejectedCreatesStoredNothing(clusterReference, recipient2),
				),
			},
		},
	})
}

// deleteClusterIfPresent deletes a cluster in the shared project that a failed test left behind.
func deleteClusterIfPresent(t *testing.T, clusterId string) {
	if clusterId == "" {
		return
	}
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", globalHost, globalOrgId, globalProjectId, clusterId)
	get := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	if _, err := globalClient.ExecuteWithRetry(context.Background(), get, nil, globalToken, nil); err != nil {
		if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
			return // already deleted
		}
		// Anything else does not prove the cluster is gone, so still try to delete it.
		t.Logf("could not check whether cluster %s still exists: %v", clusterId, err)
	}
	t.Logf("deleting cluster %s, which destroy left behind", clusterId)
	del := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	if _, err := globalClient.ExecuteWithRetry(context.Background(), del, nil, globalToken, nil); err != nil {
		t.Errorf("could not delete cluster %s, which the test left behind; delete it by hand: %v", clusterId, err)
	}
}

// wrapTolerantPattern quotes message for a regexp and lets each space match any whitespace, since
// Terraform wraps long error messages across lines.
func wrapTolerantPattern(message string) string {
	words := strings.Fields(message)
	for i, word := range words {
		words[i] = regexp.QuoteMeta(word)
	}
	return strings.Join(words, `\s+`)
}

func testAccCheckNoResourceInState(reference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, m := range s.Modules {
			if _, ok := m.Resources[reference]; ok {
				return fmt.Errorf("%s is still in state", reference)
			}
		}
		return nil
	}
}

// testAccCheckBasicClusterRejectedCreatesStoredNothing checks that none of the values the
// rejected creates tried to set reached the cluster. Values from before the downgrade may remain,
// so only the attempted ones are checked.
func testAccCheckBasicClusterRejectedCreatesStoredNothing(clusterReference, rejectedRecipient string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := datasourceAttributes(s, clusterReference)
		if err != nil {
			return err
		}
		level := alertLevelPath(globalOrgId, globalProjectId, attributes["id"])
		config, err := getAlertConfigurationForTest(level)
		if err != nil {
			return err
		}
		if entry, ok := config.Configs["cpu_usage_high"]; ok && entry.WarningThreshold != nil && *entry.WarningThreshold == 85 {
			return fmt.Errorf("a rejected create stored cpu_usage_high warning_threshold 85")
		}
		for _, address := range config.Recipients {
			if strings.EqualFold(address, rejectedRecipient) {
				return fmt.Errorf("a rejected create stored recipient %s", rejectedRecipient)
			}
		}
		for _, name := range []string{"cb_high_disk_write_queue", "cpu_usage_high"} {
			if mute, ok := config.MutedAlerts[name]; ok && mute.Action == "mute" {
				return fmt.Errorf("a rejected create stored a mute of %s", name)
			}
		}
		return nil
	}
}

func testAccBasicPlanClusterConfig(clusterName, cidr, supportPlan string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  name            = "%[4]s"
  couchbase_server = {
    version = var.server_version
  }
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 50
          type    = "gp3"
          iops    = 3000
        }
      }
      num_of_nodes = 1
      services     = ["data", "index", "query"]
    }
  ]
  availability = {
    "type" : "single"
  }
  support = {
    plan     = "%[6]s"
    timezone = "PT"
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, clusterName, cidr, supportPlan)
}

func testAccBasicPlanAlertConfig(clusterName, cidr, supportPlan, alertName, alertBody string) string {
	return testAccBasicPlanClusterConfig(clusterName, cidr, supportPlan) + fmt.Sprintf(`
resource "couchbase-capella_cluster_alert_configuration" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_cluster.%[4]s.id
%[5]s
}
`, alertName, globalOrgId, globalProjectId, clusterName, alertBody)
}
