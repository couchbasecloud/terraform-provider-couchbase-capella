package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
)

// Alert configuration is a singleton at each level (organization, project, cluster) and the
// tests run against the shared global org, project and cluster. So:
//   - TestAccAlertConfiguration runs everything that writes as sequential subtests, so two
//     configurations never overwrite each other.
//   - each subtest snapshots the level and restores it on cleanup: recipients and mutes, which
//     destroy clears, and the values of configs that existed before the test. Configs the test
//     added stay, because Capella cannot remove them.
//   - threshold and frequency overrides cannot be removed, so they are only set at the project
//     and cluster levels, with fixed values so that reruns are idempotent. The organization
//     level is exercised with recipients and mutes only, which destroy cleans up.
//   - recipients come from ACC_ALERT_RECIPIENTS (see alertTestRecipients), because Capella only
//     accepts organization users with access to the cluster.
func TestAccAlertConfiguration(t *testing.T) {
	// Allow this test to run in parallel with other top-level tests, but ensure the subtests
	// run sequentially. This is normally set by resource.ParallelTest.
	t.Parallel()

	clusterLevel := alertLevelPath(globalOrgId, globalProjectId, globalClusterId)
	projectLevel := alertLevelPath(globalOrgId, globalProjectId, "")
	orgLevel := alertLevelPath(globalOrgId, "", "")

	t.Run("Cluster Lifecycle", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			CheckDestroy:             testAccCheckAlertConfigurationMutesAndRecipientsCleared(clusterLevel),
			Steps:                    clusterAlertConfigurationLifecycleSteps(t, clusterLevel),
		})
	})

	t.Run("Cluster Mute Added Outside Terraform Is Cleared", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    clusterAlertConfigurationOutOfBandMuteSteps(t, clusterLevel),
		})
	})

	t.Run("Cluster Destroy Leaves Thresholds", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    clusterAlertConfigurationDestroySteps(t, clusterLevel),
		})
	})

	t.Run("Project Lifecycle", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, projectLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			CheckDestroy:             testAccCheckAlertConfigurationMutesAndRecipientsCleared(projectLevel),
			Steps:                    projectAlertConfigurationLifecycleSteps(t, projectLevel),
		})
	})

	t.Run("Organization Recipients And Mutes", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, orgLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			CheckDestroy:             testAccCheckAlertConfigurationMutesAndRecipientsCleared(orgLevel),
			Steps:                    organizationAlertConfigurationSteps(t, orgLevel),
		})
	})

	t.Run("Datasources", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    alertConfigurationDatasourceSteps(t),
		})
	})

	runAlertConfigurationE2ESubtests(t, clusterLevel, projectLevel, orgLevel)
}

// clusterAlertConfigurationLifecycleSteps covers create, import, an update of every attribute,
// and removing a field from configs, which Capella cannot do: the value must stay on the server
// and leave state.
func clusterAlertConfigurationLifecycleSteps(t *testing.T, level string) []resource.TestStep {
	recipient1, recipient2 := alertTestRecipients(t)
	resourceName := randomStringWithPrefix("tf_acc_cluster_alert_config_")
	resourceReference := "couchbase-capella_cluster_alert_configuration." + resourceName
	mutedUntil := alertMuteTime(24 * time.Hour)
	laterMutedUntil := alertMuteTime(48 * time.Hour)

	createConfig := testAccClusterAlertConfigurationConfig(resourceName, `
  configs = {
    cpu_usage_high = {
      frequency_hours    = 3
      warning_threshold  = 80
      critical_threshold = 90
    }
    CBHardOutOfMemoryErrors = {
      frequency_hours = 1
    }
  }

  recipients = ["`+recipient1+`"]

  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "`+mutedUntil+`"
    }
  }`)

	updateConfig := testAccClusterAlertConfigurationConfig(resourceName, `
  configs = {
    cpu_usage_high = {
      frequency_hours    = 6
      warning_threshold  = 85
      critical_threshold = 95
    }
    CBHardOutOfMemoryErrors = {
      frequency_hours = 2
    }
  }

  recipients = ["`+recipient1+`", "`+recipient2+`"]

  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "`+laterMutedUntil+`"
    }
    cpu_usage_high = {
      action = "unmute"
    }
  }`)

	// critical_threshold is removed from cpu_usage_high and CBHardOutOfMemoryErrors from configs.
	removeFieldConfig := testAccClusterAlertConfigurationConfig(resourceName, `
  configs = {
    cpu_usage_high = {
      frequency_hours   = 6
      warning_threshold = 85
    }
  }

  recipients = ["`+recipient1+`", "`+recipient2+`"]

  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "`+laterMutedUntil+`"
    }
    cpu_usage_high = {
      action = "unmute"
    }
  }`)

	return []resource.TestStep{
		{
			Config: createConfig,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(resourceReference, "configs.%", "2"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.frequency_hours", "3"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.warning_threshold", "80"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.critical_threshold", "90"),
				resource.TestCheckNoResourceAttr(resourceReference, "configs.cpu_usage_high.threshold"),
				resource.TestCheckResourceAttr(resourceReference, "configs.CBHardOutOfMemoryErrors.frequency_hours", "1"),
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "1"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "recipients.*", recipient1),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", "1"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.disk_usage_high.action", "mute"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.disk_usage_high.muted_until", mutedUntil),
				resource.TestCheckNoResourceAttr(resourceReference, "muted_alerts.disk_usage_high.muted_from"),
				testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return firstAlertError(
						expectAlertInt(c, "cpu_usage_high", "frequencyHours", 3),
						expectAlertFloat(c, "cpu_usage_high", "warningThreshold", 80),
						expectAlertFloat(c, "cpu_usage_high", "criticalThreshold", 90),
						expectAlertInt(c, "CBHardOutOfMemoryErrors", "frequencyHours", 1),
						expectAlertRecipients(c, recipient1),
						expectAlertMutes(c, map[string]string{"disk_usage_high": "mute"}),
					)
				}),
			),
		},
		// Import sets only the IDs; configs, recipients and mutes are deliberately not imported.
		{
			ResourceName:      resourceReference,
			ImportState:       true,
			ImportStateIdFunc: generateAlertConfigurationImportId(resourceReference, "cluster_id", "project_id", "organization_id"),
			ImportStateCheck: testAccCheckImportedAlertConfiguration(map[string]string{
				"organization_id": globalOrgId,
				"project_id":      globalProjectId,
				"cluster_id":      globalClusterId,
			}),
		},
		{
			Config: updateConfig,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(resourceReference, "configs.%", "2"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.frequency_hours", "6"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.warning_threshold", "85"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.critical_threshold", "95"),
				resource.TestCheckResourceAttr(resourceReference, "configs.CBHardOutOfMemoryErrors.frequency_hours", "2"),
				// Capella stores addresses in lowercase; state keeps the configured spelling.
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "2"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "recipients.*", recipient1),
				resource.TestCheckTypeSetElemAttr(resourceReference, "recipients.*", recipient2),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", "2"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.disk_usage_high.action", "mute"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.disk_usage_high.muted_until", laterMutedUntil),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.cpu_usage_high.action", "unmute"),
				resource.TestCheckNoResourceAttr(resourceReference, "muted_alerts.cpu_usage_high.muted_until"),
				testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return firstAlertError(
						expectAlertInt(c, "cpu_usage_high", "frequencyHours", 6),
						expectAlertFloat(c, "cpu_usage_high", "warningThreshold", 85),
						expectAlertFloat(c, "cpu_usage_high", "criticalThreshold", 95),
						expectAlertInt(c, "CBHardOutOfMemoryErrors", "frequencyHours", 2),
						expectAlertRecipients(c, recipient1, recipient2),
						expectAlertMutes(c, map[string]string{"disk_usage_high": "mute", "cpu_usage_high": "unmute"}),
					)
				}),
			),
		},
		// Capella cannot remove an override: the removed values stay on the server and leave
		// state. The step's post-apply plan check proves they do not come back as drift.
		{
			Config: removeFieldConfig,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "configs.%", "1"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.frequency_hours", "6"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.warning_threshold", "85"),
				resource.TestCheckNoResourceAttr(resourceReference, "configs.cpu_usage_high.critical_threshold"),
				resource.TestCheckNoResourceAttr(resourceReference, "configs.CBHardOutOfMemoryErrors.frequency_hours"),
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "2"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", "2"),
				testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return firstAlertError(
						expectAlertFloat(c, "cpu_usage_high", "criticalThreshold", 95),
						expectAlertInt(c, "CBHardOutOfMemoryErrors", "frequencyHours", 2),
					)
				}),
			),
		},
	}
}

// clusterAlertConfigurationOutOfBandMuteSteps checks the documented ownership rule for
// muted_alerts: once the attribute is set, a mute added outside Terraform shows up in the plan
// and apply clears it, even when the block is empty.
func clusterAlertConfigurationOutOfBandMuteSteps(t *testing.T, level string) []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_cluster_alert_oob_")
	resourceReference := "couchbase-capella_cluster_alert_configuration." + resourceName
	config := testAccClusterAlertConfigurationConfig(resourceName, `
  muted_alerts = {}`)

	return []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", "0"),
				resource.TestCheckNoResourceAttr(resourceReference, "configs.%"),
				resource.TestCheckNoResourceAttr(resourceReference, "recipients.#"),
				testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return expectAlertMutes(c, map[string]string{})
				}),
			),
		},
		{
			PreConfig: func() {
				mutedUntil := alertMuteTime(24 * time.Hour)
				err := putAlertConfigurationForTest(level, alertsapi.UpdateAlertConfigurationRequest{
					Alerts: []alertsapi.AlertMuteAction{{AlertName: "disk_usage_high", Action: alertsapi.AlertMuteActionMute, MutedUntil: &mutedUntil}},
				})
				if err != nil {
					t.Fatalf("could not add a mute outside Terraform: %v", err)
				}
			},
			Config: config,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate)},
			},
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", "0"),
				testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return expectAlertMutes(c, map[string]string{})
				}),
			),
		},
	}
}

// clusterAlertConfigurationDestroySteps checks what destroy does to the level: mutes and
// recipients are cleared, and thresholds and frequencies stay because Capella cannot remove them.
func clusterAlertConfigurationDestroySteps(t *testing.T, level string) []resource.TestStep {
	recipient1, _ := alertTestRecipients(t)
	resourceName := randomStringWithPrefix("tf_acc_cluster_alert_destroy_")
	resourceReference := "couchbase-capella_cluster_alert_configuration." + resourceName

	return []resource.TestStep{
		{
			Config: testAccClusterAlertConfigurationConfig(resourceName, `
  configs = {
    cpu_usage_high = {
      warning_threshold = 80
    }
  }

  recipients = ["`+recipient1+`"]

  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "`+alertMuteTime(24*time.Hour)+`"
    }
  }`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.warning_threshold", "80"),
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "1"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", "1"),
			),
		},
		{
			Config: globalProviderBlock,
			Check: testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
				return firstAlertError(
					expectAlertRecipients(c),
					expectAlertMutes(c, map[string]string{}),
					expectAlertFloat(c, "cpu_usage_high", "warningThreshold", 80),
				)
			}),
		},
	}
}

func projectAlertConfigurationLifecycleSteps(t *testing.T, level string) []resource.TestStep {
	recipient1, _ := alertTestRecipients(t)
	resourceName := randomStringWithPrefix("tf_acc_project_alert_config_")
	resourceReference := "couchbase-capella_project_alert_configuration." + resourceName
	mutedUntil := alertMuteTime(24 * time.Hour)

	config := func(frequencyHours int, recipients string) string {
		return testAccProjectAlertConfigurationConfig(resourceName, fmt.Sprintf(`
  configs = {
    cpu_usage_high = {
      frequency_hours   = %d
      warning_threshold = 82 # cpu_usage_high warning must be within [80, 89]
    }
  }

  recipients = %s

  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "%s"
    }
  }`, frequencyHours, recipients, mutedUntil))
	}

	return []resource.TestStep{
		{
			Config: config(4, `["`+recipient1+`"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "configs.%", "1"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.frequency_hours", "4"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.warning_threshold", "82"),
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "1"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "recipients.*", recipient1),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", "1"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.disk_usage_high.action", "mute"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.disk_usage_high.muted_until", mutedUntil),
				testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return firstAlertError(
						expectAlertInt(c, "cpu_usage_high", "frequencyHours", 4),
						expectAlertFloat(c, "cpu_usage_high", "warningThreshold", 82),
						expectAlertRecipients(c, recipient1),
						expectAlertMutes(c, map[string]string{"disk_usage_high": "mute"}),
					)
				}),
			),
		},
		{
			ResourceName:      resourceReference,
			ImportState:       true,
			ImportStateIdFunc: generateAlertConfigurationImportId(resourceReference, "project_id", "organization_id"),
			ImportStateCheck: testAccCheckImportedAlertConfiguration(map[string]string{
				"organization_id": globalOrgId,
				"project_id":      globalProjectId,
			}),
		},
		// Updating only the frequency must leave everything else as it was.
		{
			Config: config(8, `["`+recipient1+`"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "configs.%", "1"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.frequency_hours", "8"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.warning_threshold", "82"),
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "1"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "recipients.*", recipient1),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", "1"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.disk_usage_high.muted_until", mutedUntil),
				testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return firstAlertError(
						expectAlertInt(c, "cpu_usage_high", "frequencyHours", 8),
						expectAlertFloat(c, "cpu_usage_high", "warningThreshold", 82),
						expectAlertRecipients(c, recipient1),
					)
				}),
			),
		},
		// An empty list sets no recipients at the level, which is how it inherits its organization's.
		{
			Config: config(8, `[]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "0"),
				resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.frequency_hours", "8"),
				testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return expectAlertRecipients(c)
				}),
			),
		},
	}
}

// organizationAlertConfigurationSteps manages only recipients and mutes, because a threshold or
// frequency override at the organization level could never be removed from the shared test org.
func organizationAlertConfigurationSteps(t *testing.T, level string) []resource.TestStep {
	recipient1, recipient2 := alertTestRecipients(t)
	resourceName := randomStringWithPrefix("tf_acc_org_alert_config_")
	resourceReference := "couchbase-capella_organization_alert_configuration." + resourceName
	// Keep the organization-wide mute short in case cleanup cannot run.
	mutedUntil := alertMuteTime(2 * time.Hour)

	config := func(recipients string) string {
		return testAccOrganizationAlertConfigurationConfig(resourceName, fmt.Sprintf(`
  recipients = %s

  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "%s"
    }
  }`, recipients, mutedUntil))
	}

	return []resource.TestStep{
		{
			Config: config(`["` + recipient1 + `"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckNoResourceAttr(resourceReference, "configs.%"),
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "1"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "recipients.*", recipient1),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", "1"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.disk_usage_high.action", "mute"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.disk_usage_high.muted_until", mutedUntil),
				testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return firstAlertError(
						expectAlertRecipients(c, recipient1),
						expectAlertMutes(c, map[string]string{"disk_usage_high": "mute"}),
					)
				}),
			),
		},
		{
			ResourceName:      resourceReference,
			ImportState:       true,
			ImportStateIdFunc: generateAlertConfigurationImportId(resourceReference, "organization_id"),
			ImportStateCheck: testAccCheckImportedAlertConfiguration(map[string]string{
				"organization_id": globalOrgId,
			}),
		},
		{
			Config: config(`["` + recipient1 + `", "` + recipient2 + `"]`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "2"),
				resource.TestCheckTypeSetElemAttr(resourceReference, "recipients.*", recipient1),
				resource.TestCheckTypeSetElemAttr(resourceReference, "recipients.*", recipient2),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", "1"),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.disk_usage_high.muted_until", mutedUntil),
				testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return expectAlertRecipients(c, recipient1, recipient2)
				}),
			),
		},
	}
}

// alertConfigurationDatasourceSteps sets a cluster override and reads it back through the
// cluster, project and organization data sources and the overrides data source.
func alertConfigurationDatasourceSteps(t *testing.T) []resource.TestStep {
	recipient1, _ := alertTestRecipients(t)
	resourceName := randomStringWithPrefix("tf_acc_alert_config_ds_")
	mutedUntil := alertMuteTime(24 * time.Hour)
	clusterDatasource := "data.couchbase-capella_cluster_alert_configuration." + resourceName
	projectDatasource := "data.couchbase-capella_project_alert_configuration." + resourceName
	orgDatasource := "data.couchbase-capella_organization_alert_configuration." + resourceName
	overridesDatasource := "data.couchbase-capella_alert_configuration_overrides." + resourceName

	config := fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_alert_configuration" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"

  configs = {
    cpu_usage_high = {
      frequency_hours   = 3
      warning_threshold = 80
    }
  }

  recipients = ["`+recipient1+`"]

  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "%[6]s"
    }
  }
}

data "couchbase-capella_cluster_alert_configuration" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  depends_on      = [couchbase-capella_cluster_alert_configuration.%[2]s]
}

data "couchbase-capella_project_alert_configuration" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  depends_on      = [couchbase-capella_cluster_alert_configuration.%[2]s]
}

data "couchbase-capella_organization_alert_configuration" "%[2]s" {
  organization_id = "%[3]s"
  depends_on      = [couchbase-capella_cluster_alert_configuration.%[2]s]
}

data "couchbase-capella_alert_configuration_overrides" "%[2]s" {
  organization_id = "%[3]s"
  depends_on      = [couchbase-capella_cluster_alert_configuration.%[2]s]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, mutedUntil)

	return []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(clusterDatasource, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(clusterDatasource, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(clusterDatasource, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttr(clusterDatasource, "configs.cpu_usage_high.frequency_hours", "3"),
				resource.TestCheckResourceAttr(clusterDatasource, "configs.cpu_usage_high.warning_threshold", "80"),
				resource.TestCheckResourceAttrSet(clusterDatasource, "configs.cpu_usage_high.updated_at"),
				// updated_by is empty for changes made with an API key (AV-153487); covered by
				// RBAC/Mute_Records_The_API_Key_That_Set_It.
				// The cluster's own values win over anything it inherits.
				resource.TestCheckResourceAttr(clusterDatasource, "resolved_configs.cpu_usage_high.frequency_hours", "3"),
				resource.TestCheckResourceAttr(clusterDatasource, "resolved_configs.cpu_usage_high.warning_threshold", "80"),
				resource.TestCheckResourceAttr(clusterDatasource, "recipients.#", "1"),
				resource.TestCheckTypeSetElemAttr(clusterDatasource, "recipients.*", recipient1),
				resource.TestCheckResourceAttr(clusterDatasource, "muted_alerts.disk_usage_high.action", "mute"),
				resource.TestCheckResourceAttrSet(clusterDatasource, "muted_alerts.disk_usage_high.muted_until"),
				resource.TestCheckResourceAttr(clusterDatasource, "resolved_muted_alerts.disk_usage_high.action", "mute"),

				resource.TestCheckResourceAttr(projectDatasource, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(projectDatasource, "project_id", globalProjectId),
				resource.TestCheckResourceAttrSet(projectDatasource, "resolved_configs.%"),

				resource.TestCheckResourceAttr(orgDatasource, "organization_id", globalOrgId),
				resource.TestCheckResourceAttrSet(orgDatasource, "resolved_configs.%"),

				resource.TestCheckResourceAttr(overridesDatasource, "organization_id", globalOrgId),
				testAccCheckAlertOverrideListed(overridesDatasource, globalClusterId, "cpu_usage_high"),
			),
		},
	}
}

// TestAccDatasourceAlerts reads alerts with every filter set and checks that each returned alert
// satisfies them, which proves the filters reach the API. The list itself may be empty.
func TestAccDatasourceAlerts(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_alerts_")
	datasourceReference := "data.couchbase-capella_alerts." + resourceName
	from := time.Now().UTC().Add(-7 * 24 * time.Hour).Truncate(time.Second).Format(time.RFC3339)
	to := time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_alerts" "%[2]s" {
  organization_id = "%[3]s"
  project_ids     = ["%[4]s"]
  cluster_ids     = ["%[5]s"]
  status          = ["firing", "resolved"]
  severity_levels = ["warning", "critical"]
  from            = "%[6]s"
  to              = "%[7]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, from, to),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(datasourceReference, "from", from),
					resource.TestCheckResourceAttr(datasourceReference, "to", to),
					resource.TestCheckResourceAttrSet(datasourceReference, "data.#"),
					testAccCheckEveryListedAlert(datasourceReference, map[string][]string{
						"project_id": {globalProjectId},
						"status":     {"firing", "resolved"},
						"severity":   {"warning", "critical"},
					}),
					testAccCheckListedAlertsInClusterAndWindow(datasourceReference, globalClusterId, from, to),
				),
			},
		},
	})
}

// TestAccAlertConfigurationInvalidConfig covers the checks that reject a config at plan time,
// before anything is sent to Capella.
func TestAccAlertConfigurationInvalidConfig(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_alert_config_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterAlertConfigurationConfig(resourceName, `
  recipients = ["tf-acc-dup@couchbase.com", "TF-ACC-DUP@couchbase.com"]`),
				ExpectError: regexp.MustCompile("Duplicate Recipient"),
			},
			{
				Config: testAccClusterAlertConfigurationConfig(resourceName, `
  configs = {
    cpu_usage_high = {}
  }`),
				ExpectError: regexp.MustCompile("Empty Alert Configuration"),
			},
			{
				Config: testAccClusterAlertConfigurationConfig(resourceName, `
  configs = {
    cpu_usage_high = {
      warning_threshold = 0
    }
  }`),
				ExpectError: regexp.MustCompile(`value\s+must\s+be\s+none\s+of`),
			},
			{
				Config: testAccClusterAlertConfigurationConfig(resourceName, `
  muted_alerts = {
    disk_usage_high = {
      action = "silence"
    }
  }`),
				ExpectError: regexp.MustCompile(`value\s+must\s+be\s+one\s+of`),
			},
			{
				Config: testAccClusterAlertConfigurationConfig(resourceName, `
  muted_alerts = {
    disk_usage_high = {
      action      = "mute"
      muted_until = "tomorrow"
    }
  }`),
				ExpectError: regexp.MustCompile(`Invalid RFC3339 String Value`),
			},
		},
	})
}

// TestAccOrganizationAlertConfigurationInvalidImportId checks that the organization level, which
// has its own import parser, rejects an import ID not in the organization_id=<id> form.
func TestAccOrganizationAlertConfigurationInvalidImportId(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_org_alert_import_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:        testAccOrganizationAlertConfigurationConfig(resourceName, ""),
				ResourceName:  "couchbase-capella_organization_alert_configuration." + resourceName,
				ImportState:   true,
				ImportStateId: globalOrgId,
				ExpectError:   regexp.MustCompile(`expected\s+organization_id=<organization_id>`),
			},
		},
	})
}

func testAccClusterAlertConfigurationConfig(resourceName, body string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_alert_configuration" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
%[6]s
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, body)
}

func testAccProjectAlertConfigurationConfig(resourceName, body string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project_alert_configuration" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
%[5]s
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, body)
}

func testAccOrganizationAlertConfigurationConfig(resourceName, body string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_organization_alert_configuration" "%[2]s" {
  organization_id = "%[3]s"
%[4]s
}
`, globalProviderBlock, resourceName, globalOrgId, body)
}

// alertTestRecipients returns the first two addresses in ACC_ALERT_RECIPIENTS, a comma-separated
// list, and skips the test when fewer than two are set. Capella rejects any recipient that is
// not an organization user with access to the cluster (code 1006), so made-up addresses cannot be
// used. Use organization owners, since they can access every cluster, including ones the tests
// create; they receive real mail if an alert fires while the test runs.
func alertTestRecipients(t *testing.T) (string, string) {
	t.Helper()
	var recipients []string
	for _, address := range strings.Split(os.Getenv("ACC_ALERT_RECIPIENTS"), ",") {
		if address = strings.TrimSpace(address); address != "" {
			recipients = append(recipients, address)
		}
	}
	if len(recipients) < 2 {
		t.Skip("set ACC_ALERT_RECIPIENTS to two or more comma-separated emails of organization owners")
	}
	return recipients[0], recipients[1]
}

// alertMuteTime returns a whole-second UTC time d from now, in the form Terraform is given.
func alertMuteTime(d time.Duration) string {
	return time.Now().UTC().Add(d).Truncate(time.Second).Format(time.RFC3339)
}

// alertLevelPath is the level's path relative to /v4/organizations. Empty IDs select a higher level.
func alertLevelPath(organizationId, projectId, clusterId string) string {
	switch {
	case projectId == "":
		return organizationId
	case clusterId == "":
		return fmt.Sprintf("%s/projects/%s", organizationId, projectId)
	default:
		return fmt.Sprintf("%s/projects/%s/clusters/%s", organizationId, projectId, clusterId)
	}
}

// generateAlertConfigurationImportId builds the key=value import ID from the named attributes.
func generateAlertConfigurationImportId(resourceReference string, attributes ...string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if v, ok := m.Resources[resourceReference]; ok {
				rawState = v.Primary.Attributes
			}
		}
		if rawState == nil {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		parts := make([]string, 0, len(attributes))
		for _, attribute := range attributes {
			parts = append(parts, fmt.Sprintf("%s=%s", attribute, rawState[attribute]))
		}
		return strings.Join(parts, ","), nil
	}
}

// testAccCheckImportedAlertConfiguration checks that import parsed each ID into its own attribute
// and, as documented, imported no configs, recipients or mutes.
func testAccCheckImportedAlertConfiguration(ids map[string]string) resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		if len(states) != 1 {
			return fmt.Errorf("expected 1 imported resource, got %d", len(states))
		}
		attributes := states[0].Attributes
		for name, want := range ids {
			if got := attributes[name]; got != want {
				return fmt.Errorf("imported %s = %q, want %q", name, got, want)
			}
		}
		for _, name := range []string{"configs.%", "recipients.#", "muted_alerts.%"} {
			if got, ok := attributes[name]; ok {
				return fmt.Errorf("import should not set %s, got %q", name, got)
			}
		}
		return nil
	}
}

// testAccCheckAlertOverrideListed checks that the overrides data source lists an override of
// alertKind on the resource with resourceId.
func testAccCheckAlertOverrideListed(datasourceReference, resourceId, alertKind string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := datasourceAttributes(s, datasourceReference)
		if err != nil {
			return err
		}
		for key, value := range attributes {
			if !strings.HasSuffix(key, ".resource_id") || value != resourceId {
				continue
			}
			prefix := strings.TrimSuffix(key, "resource_id")
			if attributes[prefix+"alert_kind"] == alertKind {
				return nil
			}
		}
		return fmt.Errorf("%s does not list an override of %s on %s", datasourceReference, alertKind, resourceId)
	}
}

// testAccCheckEveryListedAlert checks that every alert in the data source has, for each
// attribute, one of the allowed values.
func testAccCheckEveryListedAlert(datasourceReference string, allowed map[string][]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := datasourceAttributes(s, datasourceReference)
		if err != nil {
			return err
		}
		var count int
		if _, err := fmt.Sscan(attributes["data.#"], &count); err != nil {
			return fmt.Errorf("could not read data.#: %w", err)
		}
		for i := 0; i < count; i++ {
			for attribute, values := range allowed {
				got := attributes[fmt.Sprintf("data.%d.%s", i, attribute)]
				if !containsString(values, got) {
					return fmt.Errorf("data.%d.%s = %q, which the filter excludes (want one of %v)", i, attribute, got, values)
				}
			}
		}
		return nil
	}
}

// testAccCheckListedAlertsInClusterAndWindow checks the filters testAccCheckEveryListedAlert
// cannot: a cluster-scoped alert belongs to clusterId, and every alert's event time is inside
// [from, to]. An empty list passes, because alerts cannot be raised on demand from Terraform.
func testAccCheckListedAlertsInClusterAndWindow(datasourceReference, clusterId, from, to string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := datasourceAttributes(s, datasourceReference)
		if err != nil {
			return err
		}
		start, err1 := time.Parse(time.RFC3339, from)
		end, err2 := time.Parse(time.RFC3339, to)
		if err1 != nil || err2 != nil {
			return fmt.Errorf("could not parse the window %s to %s", from, to)
		}
		var count int
		if _, err := fmt.Sscan(attributes["data.#"], &count); err != nil {
			return fmt.Errorf("could not read data.#: %w", err)
		}
		for i := 0; i < count; i++ {
			get := func(field string) string { return attributes[fmt.Sprintf("data.%d.%s", i, field)] }
			if get("resource_type") == "cluster" && get("resource_id") != clusterId {
				return fmt.Errorf("alert %s is on cluster %s, which the cluster_ids filter excludes", get("alert_id"), get("resource_id"))
			}
			eventTime, err := time.Parse(time.RFC3339Nano, get("event_time"))
			if err != nil || eventTime.Before(start) || eventTime.After(end) {
				return fmt.Errorf("alert %s has event_time %q, outside the window %s to %s", get("alert_id"), get("event_time"), from, to)
			}
		}
		return nil
	}
}

func datasourceAttributes(s *terraform.State, reference string) (map[string]string, error) {
	for _, m := range s.Modules {
		if v, ok := m.Resources[reference]; ok {
			return v.Primary.Attributes, nil
		}
	}
	return nil, fmt.Errorf("%s not found in state", reference)
}

func containsString(values []string, s string) bool {
	for _, v := range values {
		if v == s {
			return true
		}
	}
	return false
}

// testAccCheckAlertConfigurationOnServer reads the level from Capella and passes it to check, so
// that assertions prove what was stored rather than what was written to state.
func testAccCheckAlertConfigurationOnServer(level string, check func(*alertsapi.GetAlertConfigurationResponse) error) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		config, err := getAlertConfigurationForTest(level)
		if err != nil {
			return err
		}
		return check(config)
	}
}

// testAccCheckAlertConfigurationMutesAndRecipientsCleared is the CheckDestroy for resources that
// manage recipients and mutes: destroy must clear both.
func testAccCheckAlertConfigurationMutesAndRecipientsCleared(level string) resource.TestCheckFunc {
	return testAccCheckAlertConfigurationOnServer(level, func(c *alertsapi.GetAlertConfigurationResponse) error {
		return firstAlertError(expectAlertRecipients(c), expectAlertMutes(c, map[string]string{}))
	})
}

func expectAlertInt(c *alertsapi.GetAlertConfigurationResponse, alert, field string, want int64) error {
	entry, ok := c.Configs[alert]
	if !ok {
		return fmt.Errorf("alert %s has no config at the level", alert)
	}
	var got *int64
	if field == "frequencyHours" {
		got = entry.FrequencyHours
	}
	if got == nil || *got != want {
		return fmt.Errorf("alert %s %s = %v, want %d", alert, field, derefInt64(got), want)
	}
	return nil
}

func expectAlertFloat(c *alertsapi.GetAlertConfigurationResponse, alert, field string, want float64) error {
	entry, ok := c.Configs[alert]
	if !ok {
		return fmt.Errorf("alert %s has no config at the level", alert)
	}
	got := map[string]*float64{
		"threshold":         entry.Threshold,
		"warningThreshold":  entry.WarningThreshold,
		"criticalThreshold": entry.CriticalThreshold,
	}[field]
	if got == nil || *got != want {
		return fmt.Errorf("alert %s %s = %v, want %v", alert, field, derefFloat64(got), want)
	}
	return nil
}

// expectAlertRecipients compares as a set, ignoring case, since Capella stores lowercase.
func expectAlertRecipients(c *alertsapi.GetAlertConfigurationResponse, want ...string) error {
	got := make(map[string]struct{}, len(c.Recipients))
	for _, address := range c.Recipients {
		got[strings.ToLower(address)] = struct{}{}
	}
	mismatch := len(got) != len(want)
	for _, address := range want {
		if _, ok := got[strings.ToLower(address)]; !ok {
			mismatch = true
		}
	}
	if mismatch {
		return fmt.Errorf("recipients = %v, want %v", c.Recipients, want)
	}
	return nil
}

// expectAlertMutes checks the exact set of mutes at the level and each one's action.
func expectAlertMutes(c *alertsapi.GetAlertConfigurationResponse, want map[string]string) error {
	if len(c.MutedAlerts) != len(want) {
		return fmt.Errorf("muted alerts = %v, want %v", c.MutedAlerts, want)
	}
	for name, action := range want {
		mute, ok := c.MutedAlerts[name]
		if !ok || mute.Action != action {
			return fmt.Errorf("muted alerts = %v, want %v", c.MutedAlerts, want)
		}
	}
	return nil
}

func firstAlertError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func derefInt64(v *int64) any {
	if v == nil {
		return nil
	}
	return *v
}

func derefFloat64(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}

func getAlertConfigurationForTest(level string) (*alertsapi.GetAlertConfigurationResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/alertConfiguration", globalHost, level)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
	if err != nil {
		return nil, fmt.Errorf("could not read the alert configuration at %s: %w", level, err)
	}
	var config alertsapi.GetAlertConfigurationResponse
	if err := json.Unmarshal(response.Body, &config); err != nil {
		return nil, fmt.Errorf("could not decode the alert configuration at %s: %w", level, err)
	}
	return &config, nil
}

func putAlertConfigurationForTest(level string, update alertsapi.UpdateAlertConfigurationRequest) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/alertConfiguration", globalHost, level)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := globalClient.ExecuteWithRetry(context.Background(), cfg, update, globalToken, nil)
	return err
}

// restoreAlertConfigurationAfterTest snapshots the level's recipients and mutes and puts them
// back when the test ends, since destroy clears both. Thresholds and frequencies cannot be
// removed, so they are not restored.
func restoreAlertConfigurationAfterTest(t *testing.T, level string) {
	t.Helper()
	before, err := getAlertConfigurationForTest(level)
	if err != nil {
		t.Fatalf("could not snapshot the alert configuration before the test: %v", err)
	}

	t.Cleanup(func() {
		now, err := getAlertConfigurationForTest(level)
		if err != nil {
			t.Errorf("could not restore the alert configuration at %s: %v", level, err)
			return
		}

		// Configs and recipients, and mutes, are restored in separate requests, so that Capella
		// rejecting one cannot leave the other changed.
		var update alertsapi.UpdateAlertConfigurationRequest
		for name, was := range before.Configs {
			if sameAlertConfigs(map[string]alertsapi.AlertConfigurationEntry{name: was}, map[string]alertsapi.AlertConfigurationEntry{name: now.Configs[name]}) == nil {
				continue
			}
			if update.Configs == nil {
				update.Configs = make(map[string]alertsapi.AlertConfigurationUpdate)
			}
			update.Configs[name] = alertsapi.AlertConfigurationUpdate{
				FrequencyHours:    was.FrequencyHours,
				Threshold:         roundAlertThreshold(was.Threshold),
				WarningThreshold:  roundAlertThreshold(was.WarningThreshold),
				CriticalThreshold: roundAlertThreshold(was.CriticalThreshold),
			}
		}
		if expectAlertRecipients(now, before.Recipients...) != nil {
			recipients := append([]string{}, before.Recipients...)
			update.Recipients = &recipients
		}
		if update.Recipients != nil || len(update.Configs) > 0 {
			if err := putAlertConfigurationForTest(level, update); err != nil {
				t.Errorf("could not restore the configs and recipients at %s to their state before the test: %v", level, err)
			}
		}

		// A mute whose window has ended since the snapshot cannot be sent again (Capella refuses a
		// past muted_until) and silences nothing, so it is cleared rather than restored.
		var mutes alertsapi.UpdateAlertConfigurationRequest
		for _, name := range sortedAlertNames(now.MutedAlerts) {
			was, existed := before.MutedAlerts[name]
			if !existed || (alertMuteExpired(was) && !sameAlertMute(was, now.MutedAlerts[name])) {
				mutes.Alerts = append(mutes.Alerts, alertsapi.AlertMuteAction{AlertName: name, Action: alertsapi.AlertMuteActionClear})
			}
		}
		for _, name := range sortedAlertNames(before.MutedAlerts) {
			mute := before.MutedAlerts[name]
			if current, ok := now.MutedAlerts[name]; (ok && sameAlertMute(mute, current)) || alertMuteExpired(mute) {
				continue
			}
			mutes.Alerts = append(mutes.Alerts, alertsapi.AlertMuteAction{
				AlertName: name, Action: mute.Action, MutedFrom: mute.MutedFrom, MutedUntil: mute.MutedUntil,
			})
		}
		if len(mutes.Alerts) > 0 {
			if err := putAlertConfigurationForTest(level, mutes); err != nil {
				t.Errorf("could not restore the mutes at %s to their state before the test: %v", level, err)
			}
		}
	})
}

// alertMuteExpired reports whether a mute's window has already ended.
func alertMuteExpired(mute alertsapi.AlertMuteState) bool {
	if mute.MutedUntil == nil {
		return false
	}
	until, err := time.Parse(time.RFC3339Nano, *mute.MutedUntil)
	return err == nil && !until.After(time.Now())
}

// roundAlertThreshold converts a threshold read from Capella to the whole number it accepts.
func roundAlertThreshold(v *float64) *int64 {
	if v == nil {
		return nil
	}
	rounded := int64(math.Round(*v))
	return &rounded
}

func sameAlertMute(a, b alertsapi.AlertMuteState) bool {
	same := func(x, y *string) bool { return (x == nil && y == nil) || (x != nil && y != nil && *x == *y) }
	return a.Action == b.Action && same(a.MutedFrom, b.MutedFrom) && same(a.MutedUntil, b.MutedUntil)
}

func sortedAlertNames(m map[string]alertsapi.AlertMuteState) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
