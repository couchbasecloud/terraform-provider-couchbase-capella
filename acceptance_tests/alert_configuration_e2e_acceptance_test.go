package acceptance_tests

// Terraform coverage of the Capella Alert Management E2E test plan
// (Alert Management E2E Test Plan, Nov 2026), its PRD requirements, and the catalog and guard
// rails in the Capella Alerts sheet. Each test names the plan IDs it covers.
//
// What Terraform can reach is configuration state: thresholds, frequency, mutes and disables,
// recipients and inheritance, at organization, project and cluster scope, plus the read-only
// alert list and overrides. Everything below is checked against what Capella stores, read back
// through the API, not only against Terraform state.
//
// Plan cases Terraform cannot reach, and why:
//   - Delivery and timing (RES-F-*, RES-N-*, FRQ-06..13, FRQ-19..25, MUT-03/04/07, THR-08..14,
//     THR-19, NFR-*, RES-P-*): need firing alerts, channels and clocks; no provider surface.
//   - UI, bell and accessibility (RES-U-*, LST-01/02/09/16/17, HIE-16..18, MUT-24/28/29, DIS-24).
//   - Feature flags and fault injection (RES-E-01..07, RES-E-15..23).
//   - App Service scope (DIS-05, MUT-19, FRQ-12/18, HIE-04): the provider has no App Service
//     level. App Service alerts are configured at project scope in the hierarchy test instead.
//   - Thresholds and frequency at organization scope (THR-15, FRQ-14 at org): Capella cannot
//     remove an override, so a test value would stay on the shared organization for good.
//     Inheritance is tested from a project the test creates and deletes instead.
//   - Clearing an override (THR-18, FRQ-17): blocked in Capella itself; the provider's warning
//     and "leave it in place" behaviour is covered by TestAccAlertConfiguration/Cluster_Lifecycle.
//   - THR-04 (warning >= critical): the warning and critical rails in the Capella Alerts sheet
//     never overlap, so any such pair is already out of range and fails as THR-03.
//   - THR-20 (change on a deploying cluster): the cluster resource waits for a healthy cluster
//     before anything that depends on it is created.
//   - HIE-12 (re-saving an unchanged value): the provider sends only differences, so an
//     unchanged value is never re-sent; every step's empty post-apply plan proves it.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
)

// alertCatalog is every alert that Alert Configuration offers, by the identifier the API uses
// (Capella Alerts sheet, and the v4 API testing proof's bulk mute).
var alertCatalog = []string{
	"cpu_usage_high",
	"disk_usage_high",
	"disk_iops_usage",
	"cb_high_disk_throughput_usage",
	"cb_high_disk_write_queue",
	"indexer_low_resident_percentage",
	"NumDivergingReplicaIndexes",
	"CBHardOutOfMemoryErrors",
	"DataEncryptionDEKGenerationFailures",
	"DataEncryptionKeyRotationFailures",
	"syncgateway_cpu_usage_high",
	"syncgateway_mem_usage",
	"SGWAccessErrorWarning",
	"SGWDataSyncErrorWarning",
	"SGWImportErrorWarning",
}

// alertCamelCaseCatalog is the part of the catalog whose identifiers are not snake_case. Capella
// stores a mute of these under the lowercased name, which the provider does not yet match back
// to the configured key; see the "Mute Keys Keep Their Spelling" subtest.
var alertCamelCaseCatalog = []string{
	"NumDivergingReplicaIndexes",
	"CBHardOutOfMemoryErrors",
	"DataEncryptionDEKGenerationFailures",
	"DataEncryptionKeyRotationFailures",
	"SGWAccessErrorWarning",
	"SGWDataSyncErrorWarning",
	"SGWImportErrorWarning",
}

// alertSnakeCaseCatalog is the rest of the catalog.
var alertSnakeCaseCatalog = func() []string {
	var names []string
	for _, name := range alertCatalog {
		if !containsString(alertCamelCaseCatalog, name) {
			names = append(names, name)
		}
	}
	return names
}()

// alertAccessDenied is Capella's 403 (code 1002) for a role without access to the scope.
var alertAccessDenied = regexp.MustCompile(`Access\s+Denied\.`)

// alertApiRejected matches a write that Capella refused with a 4xx, for cases whose exact
// message is not documented. The plan's "nothing stored" half is checked separately by
// testAccCheckAlertConfigurationUnchanged.
func alertApiRejected(level string) *regexp.Regexp {
	return regexp.MustCompile(`(?s)Could\s+not\s+set\s+the\s+` + level +
		`\s+alert\s+configuration,\s+unexpected\s+error:.*"httpStatusCode":\s*4\d\d`)
}

// alertRangeError matches Capella's guard-rail message, for example:
// threshold 75 for "cpu_usage_high (warning)" is out of allowed range [80, 89].
func alertRangeError(value int, alert, severity string, minimum, maximum int) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(
		`(?s)threshold\s+%d\s+for\s+\\"%s\s+\(%s\)\\"\s+is\s+out\s+of\s+allowed\s+range\s+\[%d,\s+%d\]`,
		value, regexp.QuoteMeta(alert), severity, minimum, maximum))
}

// runAlertConfigurationE2ESubtests runs the plan's configuration cases that write to the shared
// organization, project and cluster. They are called from TestAccAlertConfiguration so that they
// never run alongside another test writing the same singleton levels.
func runAlertConfigurationE2ESubtests(t *testing.T, clusterLevel, projectLevel, orgLevel string) {
	t.Run("API Validation Rejects And Stores Nothing", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		restoreAlertConfigurationAfterTest(t, projectLevel)
		restoreAlertConfigurationAfterTest(t, orgLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    alertConfigurationApiValidationSteps(clusterLevel, projectLevel, orgLevel),
		})
	})

	// A resource with only its IDs manages nothing: create sends nothing, state holds no
	// settings, and destroy clears nothing, at every level.
	t.Run("Ids Only Manages Nothing", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		restoreAlertConfigurationAfterTest(t, projectLevel)
		restoreAlertConfigurationAfterTest(t, orgLevel)
		resourceName := randomStringWithPrefix("tf_acc_alert_ids_only_")
		var clusterBefore, projectBefore, orgBefore *alertsapi.GetAlertConfigurationResponse
		unchanged := resource.ComposeAggregateTestCheckFunc(
			testAccCheckAlertConfigurationUnchanged(clusterLevel, &clusterBefore),
			testAccCheckAlertConfigurationUnchanged(projectLevel, &projectBefore),
			testAccCheckAlertConfigurationUnchanged(orgLevel, &orgBefore),
		)
		config := globalProviderBlock
		var checks []resource.TestCheckFunc
		for _, level := range []string{"organization", "project", "cluster"} {
			config += fmt.Sprintf("\nresource \"couchbase-capella_%s_alert_configuration\" %q {\n%s}\n",
				level, resourceName, alertRBACIds(level, globalProjectId, globalClusterId))
			reference := "couchbase-capella_" + level + "_alert_configuration." + resourceName
			checks = append(checks,
				resource.TestCheckNoResourceAttr(reference, "configs.%"),
				resource.TestCheckNoResourceAttr(reference, "recipients.#"),
				resource.TestCheckNoResourceAttr(reference, "muted_alerts.%"),
			)
		}
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			CheckDestroy:             unchanged,
			Steps: []resource.TestStep{{
				PreConfig: func() {
					var err error
					if clusterBefore, err = getAlertConfigurationForTest(clusterLevel); err != nil {
						t.Fatal(err)
					}
					if projectBefore, err = getAlertConfigurationForTest(projectLevel); err != nil {
						t.Fatal(err)
					}
					if orgBefore, err = getAlertConfigurationForTest(orgLevel); err != nil {
						t.Fatal(err)
					}
				},
				Config: config,
				Check:  resource.ComposeAggregateTestCheckFunc(append(checks, unchanged)...),
			}},
		})
	})

	t.Run("Thresholds And Frequency At The Guard Rails", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    alertConfigurationGuardRailSteps(clusterLevel),
		})
	})

	t.Run("Mute Disable And Enable Lifecycle", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			CheckDestroy: testAccCheckAlertConfigurationOnServer(clusterLevel, func(c *alertsapi.GetAlertConfigurationResponse) error {
				return expectAlertMutes(c, map[string]string{})
			}),
			Steps: alertConfigurationMuteLifecycleSteps(clusterLevel),
		})
	})

	t.Run("Mute Precedence Across Scopes", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		restoreAlertConfigurationAfterTest(t, projectLevel)
		restoreAlertConfigurationAfterTest(t, orgLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    alertConfigurationMutePrecedenceSteps(clusterLevel, projectLevel, orgLevel),
		})
	})

	t.Run("Recipients Lifecycle And Ownership", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps:                    alertConfigurationRecipientsSteps(t, clusterLevel),
		})
	})

	// Capella stores a mute of a CamelCase alert under its lowercased name, and a severity name
	// under the alert identifier. The provider keeps the server's key in state, so the apply
	// fails with "Provider produced inconsistent result" and taints the resource. This asserts
	// the correct behaviour and fails until the provider matches keys the way Capella does.
	t.Run("Mute Keys Keep Their Spelling", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		resourceName := randomStringWithPrefix("tf_acc_alert_mute_keys_")
		resourceReference := "couchbase-capella_cluster_alert_configuration." + resourceName
		until := alertMuteTime(24 * time.Hour)
		var lines []string
		for _, name := range alertCamelCaseCatalog {
			lines = append(lines, fmt.Sprintf("    %s = { action = \"mute\", muted_until = %q }", name, until))
		}
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps: []resource.TestStep{{
				Config: testAccClusterAlertConfigurationConfig(resourceName, "  muted_alerts = {\n"+strings.Join(lines, "\n")+"\n  }"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", fmt.Sprint(len(alertCamelCaseCatalog))),
					resource.TestCheckResourceAttr(resourceReference, "muted_alerts.CBHardOutOfMemoryErrors.action", "mute"),
				),
			}},
		})
	})

	// DIS-06: an alert that is not configurable (support only) must be refused. Plan run #35
	// found the customer API accepting a disable of CNGAuthenticationCheck at cluster scope.
	t.Run("Support Only Alert Rejected", func(t *testing.T) {
		t.Skip("DIS-06: Capella accepts a mute of the support-only CNGAuthenticationCheck (E2E plan run #35); unskip once it is refused")
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps: []resource.TestStep{{
				Config: testAccClusterAlertConfigurationConfig(randomStringWithPrefix("tf_acc_alert_support_only_"), `
  muted_alerts = {
    CNGAuthenticationCheck = {
      action = "mute"
    }
  }`),
				ExpectError: alertApiRejected("cluster"),
			}},
		})
	})

	t.Run("RBAC", func(t *testing.T) {
		restoreAlertConfigurationAfterTest(t, clusterLevel)
		restoreAlertConfigurationAfterTest(t, projectLevel)
		runAlertConfigurationRBACSubtests(t, clusterLevel, projectLevel)
	})
}

// alertConfigurationApiValidationSteps covers what Capella must refuse, and that a refused
// request stores nothing at any scope.
// Plan: THR-03, THR-05, THR-06, THR-07, FRQ-03, MUT-21, MUT-22, RCP-02, RCP-08, DIS-06 (unknown
// kinds), HIE-08 (one request is all or nothing), RES-A-10 style not-found and mismatched IDs.
func alertConfigurationApiValidationSteps(clusterLevel, projectLevel, orgLevel string) []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_alert_invalid_api_")
	var clusterBefore, projectBefore, orgBefore *alertsapi.GetAlertConfigurationResponse
	snapshot := func() {
		var err error
		if clusterBefore, err = getAlertConfigurationForTest(clusterLevel); err != nil {
			panic(err)
		}
		if projectBefore, err = getAlertConfigurationForTest(projectLevel); err != nil {
			panic(err)
		}
		if orgBefore, err = getAlertConfigurationForTest(orgLevel); err != nil {
			panic(err)
		}
	}
	cluster := func(body string) string { return testAccClusterAlertConfigurationConfig(resourceName, body) }
	thresholds := func(alert, field string, value int) string {
		return cluster(fmt.Sprintf("  configs = {\n    %s = {\n      %s = %d\n    }\n  }", alert, field, value))
	}
	mute := func(alert, from, until string) string {
		fields := `action = "mute"`
		if from != "" {
			fields += fmt.Sprintf(`, muted_from = %q`, from)
		}
		if until != "" {
			fields += fmt.Sprintf(`, muted_until = %q`, until)
		}
		return cluster(fmt.Sprintf("  muted_alerts = {\n    %s = { %s }\n  }", alert, fields))
	}
	recipients := func(addresses ...string) string {
		quoted := make([]string, len(addresses))
		for i, address := range addresses {
			quoted[i] = fmt.Sprintf("%q", address)
		}
		return cluster(fmt.Sprintf("  recipients = [%s]", strings.Join(quoted, ", ")))
	}
	future := alertMuteTime(24 * time.Hour)
	past := alertMuteTime(-1 * time.Hour)
	rejected := alertApiRejected("cluster")

	steps := []resource.TestStep{
		// THR-03: below and above each rail, for warning and critical, and negative.
		{PreConfig: snapshot, Config: thresholds("cpu_usage_high", "warning_threshold", 79), ExpectError: alertRangeError(79, "cpu_usage_high", "warning", 80, 89)},
		{Config: thresholds("cpu_usage_high", "warning_threshold", 90), ExpectError: alertRangeError(90, "cpu_usage_high", "warning", 80, 89)},
		{Config: thresholds("disk_usage_high", "warning_threshold", 74), ExpectError: alertRangeError(74, "disk_usage_high", "warning", 75, 85)},
		{Config: thresholds("disk_usage_high", "warning_threshold", 86), ExpectError: alertRangeError(86, "disk_usage_high", "warning", 75, 85)},
		{Config: thresholds("disk_usage_high", "critical_threshold", 85), ExpectError: alertRangeError(85, "disk_usage_high", "critical", 86, 90)},
		{Config: thresholds("disk_usage_high", "critical_threshold", 91), ExpectError: alertRangeError(91, "disk_usage_high", "critical", 86, 90)},
		{Config: thresholds("disk_iops_usage", "critical_threshold", 96), ExpectError: alertRangeError(96, "disk_iops_usage", "critical", 91, 95)},
		{Config: thresholds("cb_high_disk_throughput_usage", "warning_threshold", 79), ExpectError: alertRangeError(79, "cb_high_disk_throughput_usage", "warning", 80, 90)},
		{Config: thresholds("disk_iops_usage", "warning_threshold", -5), ExpectError: alertRangeError(-5, "disk_iops_usage", "warning", 80, 90)},
		// FRQ-03: frequency must be a whole number of hours from 1 to 24. Decimals and strings
		// fail at plan time (TestAccAlertConfigurationPlanTimeValidation).
		{Config: thresholds("cpu_usage_high", "frequency_hours", 0), ExpectError: rejected},
		{Config: thresholds("cpu_usage_high", "frequency_hours", 25), ExpectError: rejected},
		{Config: thresholds("cpu_usage_high", "frequency_hours", -1), ExpectError: rejected},
		// THR-06, THR-07: alerts whose thresholds are fixed, or that have none, refuse one.
		{Config: thresholds("CBHardOutOfMemoryErrors", "warning_threshold", 80), ExpectError: rejected},
		{Config: thresholds("indexer_low_resident_percentage", "critical_threshold", 5), ExpectError: rejected},
		{Config: thresholds("cb_high_disk_write_queue", "warning_threshold", 600000), ExpectError: rejected},
		{Config: thresholds("DataEncryptionDEKGenerationFailures", "threshold", 1), ExpectError: rejected},
		// DIS-06: kinds that Alert Configuration does not offer, in configs and in mutes.
		{Config: thresholds("not_a_capella_alert", "frequency_hours", 3), ExpectError: rejected},
		{Config: mute("not_a_capella_alert", "", future), ExpectError: rejected},
		// MUT-21, MUT-22: until in the past, until before from, and from equal to until.
		{Config: mute("cpu_usage_high", "", past), ExpectError: rejected},
		{Config: mute("cpu_usage_high", alertMuteTime(48*time.Hour), future), ExpectError: rejected},
		{Config: mute("cpu_usage_high", future, future), ExpectError: rejected},
		// RCP-02: malformed addresses.
		{Config: recipients("user@company"), ExpectError: rejected},
		{Config: recipients("user.company.com"), ExpectError: rejected},
		{Config: recipients("user name@couchbase.com"), ExpectError: rejected},
		{Config: recipients(""), ExpectError: rejected},
		// RCP-08: Capella accepts only organization users with access to the cluster, and refuses
		// some domains outright.
		{
			Config:      recipients("tf-acc-not-an-org-user@couchbase.com"),
			ExpectError: regexp.MustCompile(wrapTolerantPattern("recipients must be organization users with access to this cluster")),
		},
		{
			Config:      recipients("tf-acc@example.com"),
			ExpectError: regexp.MustCompile(wrapTolerantPattern("email accounts using example.com are not accepted.")),
		},
		// HIE-08: a valid mute and an invalid frequency in one request. The final step proves
		// the mute was not stored either.
		{
			Config: cluster(`
  configs = {
    cpu_usage_high = {
      frequency_hours = 25
    }
  }
  muted_alerts = {
    cb_high_disk_write_queue = {
      action      = "mute"
      muted_until = "` + future + `"
    }
  }`),
			ExpectError: rejected,
		},
		// The same guard rails apply at project and organization scope.
		{
			Config:      testAccProjectAlertConfigurationConfig(resourceName, "  configs = {\n    cpu_usage_high = {\n      warning_threshold = 79\n    }\n  }"),
			ExpectError: alertRangeError(79, "cpu_usage_high", "warning", 80, 89),
		},
		{
			Config:      testAccOrganizationAlertConfigurationConfig(resourceName, "  configs = {\n    cpu_usage_high = {\n      warning_threshold = 79\n    }\n  }"),
			ExpectError: alertRangeError(79, "cpu_usage_high", "warning", 80, 89),
		},
		// IDs that do not exist, or do not belong together (v4 API testing proof: codes 4031,
		// 4025 and 2000).
		{
			Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_cluster_alert_configuration" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, resourceName, globalOrgId, uuid.NewString(), globalClusterId),
			ExpectError: regexp.MustCompile(`not\s+valid\s+for\s+the\s+cluster`),
		},
		{
			Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_cluster_alert_configuration" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, uuid.NewString()),
			ExpectError: regexp.MustCompile(`Unable\s+to\s+fetch\s+the\s+cluster\s+details`),
		},
		{
			Config: fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project_alert_configuration" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
}
`, globalProviderBlock, resourceName, globalOrgId, uuid.NewString()),
			ExpectError: regexp.MustCompile(`cannot\s+find\s+a\s+project\s+by\s+its\s+ID`),
		},
		// Nothing above was stored at any scope.
		{
			Config: globalProviderBlock,
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccCheckAlertConfigurationUnchanged(clusterLevel, &clusterBefore),
				testAccCheckAlertConfigurationUnchanged(projectLevel, &projectBefore),
				testAccCheckAlertConfigurationUnchanged(orgLevel, &orgBefore),
			),
		},
	}
	return steps
}

// alertConfigurationGuardRailSteps sets each configurable family at both ends of its rails and
// frequency across its range, then moves every value, and checks Capella stored exactly that.
// Plan: THR-02, THR-05 (value exactly at the rail), THR-12 (each family, config side), FRQ-02,
// FRQ-04 (frequency accepted where thresholds are fixed), HIE-03 (fields set separately).
func alertConfigurationGuardRailSteps(clusterLevel string) []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_alert_rails_")
	resourceReference := "couchbase-capella_cluster_alert_configuration." + resourceName

	type values struct {
		cpuWarning, cpuFrequency                int64
		diskWarning, diskCritical               int64
		iopsWarning, iopsCritical               int64
		throughputWarning, throughputCritical   int64
		oomFrequency, divergingFrequency, rrrFq int64
	}
	config := func(v values) string {
		return testAccClusterAlertConfigurationConfig(resourceName, fmt.Sprintf(`
  configs = {
    cpu_usage_high = {
      warning_threshold = %d
      frequency_hours   = %d
    }
    disk_usage_high = {
      warning_threshold  = %d
      critical_threshold = %d
    }
    disk_iops_usage = {
      warning_threshold  = %d
      critical_threshold = %d
    }
    cb_high_disk_throughput_usage = {
      warning_threshold  = %d
      critical_threshold = %d
    }
    CBHardOutOfMemoryErrors = {
      frequency_hours = %d
    }
    NumDivergingReplicaIndexes = {
      frequency_hours = %d
    }
    indexer_low_resident_percentage = {
      frequency_hours = %d
    }
  }`, v.cpuWarning, v.cpuFrequency, v.diskWarning, v.diskCritical, v.iopsWarning, v.iopsCritical,
			v.throughputWarning, v.throughputCritical, v.oomFrequency, v.divergingFrequency, v.rrrFq))
	}
	check := func(v values) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.warning_threshold", fmt.Sprint(v.cpuWarning)),
			resource.TestCheckResourceAttr(resourceReference, "configs.cpu_usage_high.frequency_hours", fmt.Sprint(v.cpuFrequency)),
			resource.TestCheckResourceAttr(resourceReference, "configs.disk_usage_high.warning_threshold", fmt.Sprint(v.diskWarning)),
			resource.TestCheckResourceAttr(resourceReference, "configs.disk_usage_high.critical_threshold", fmt.Sprint(v.diskCritical)),
			resource.TestCheckResourceAttr(resourceReference, "configs.CBHardOutOfMemoryErrors.frequency_hours", fmt.Sprint(v.oomFrequency)),
			resource.TestCheckNoResourceAttr(resourceReference, "configs.CBHardOutOfMemoryErrors.warning_threshold"),
			testAccCheckAlertConfigurationOnServer(clusterLevel, func(c *alertsapi.GetAlertConfigurationResponse) error {
				return firstAlertError(
					expectAlertFloat(c, "cpu_usage_high", "warningThreshold", float64(v.cpuWarning)),
					expectAlertInt(c, "cpu_usage_high", "frequencyHours", v.cpuFrequency),
					expectAlertFloat(c, "disk_usage_high", "warningThreshold", float64(v.diskWarning)),
					expectAlertFloat(c, "disk_usage_high", "criticalThreshold", float64(v.diskCritical)),
					expectAlertFloat(c, "disk_iops_usage", "warningThreshold", float64(v.iopsWarning)),
					expectAlertFloat(c, "disk_iops_usage", "criticalThreshold", float64(v.iopsCritical)),
					expectAlertFloat(c, "cb_high_disk_throughput_usage", "warningThreshold", float64(v.throughputWarning)),
					expectAlertFloat(c, "cb_high_disk_throughput_usage", "criticalThreshold", float64(v.throughputCritical)),
					expectAlertInt(c, "CBHardOutOfMemoryErrors", "frequencyHours", v.oomFrequency),
					expectAlertInt(c, "NumDivergingReplicaIndexes", "frequencyHours", v.divergingFrequency),
					expectAlertInt(c, "indexer_low_resident_percentage", "frequencyHours", v.rrrFq),
					// The cluster's own values are what is in effect.
					expectResolvedAlertFloat(c, "cpu_usage_high", "warningThreshold", float64(v.cpuWarning)),
					expectResolvedAlertInt(c, "cpu_usage_high", "frequencyHours", v.cpuFrequency),
				)
			}),
		)
	}

	lowerRails := values{89, 1, 75, 86, 80, 91, 90, 95, 24, 12, 6}
	upperRails := values{80, 24, 85, 90, 90, 95, 80, 91, 3, 1, 24}
	return []resource.TestStep{
		{Config: config(lowerRails), Check: check(lowerRails)},
		{Config: config(upperRails), Check: check(upperRails)},
	}
}

// alertConfigurationMuteLifecycleSteps drives one alert through every mute transition at
// cluster scope, plus a bulk mute of the whole catalog.
// Plan: MUT-02 (config side), MUT-05, MUT-06, MUT-08, MUT-09, MUT-23, MUT-25, MUT-30 (API side),
// DIS-01, DIS-02 (config side), DIS-20, DIS-21, DIS-22, DIS-23 (API side), AUD-01 and AUD-04.
func alertConfigurationMuteLifecycleSteps(clusterLevel string) []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_alert_mutes_")
	resourceReference := "couchbase-capella_cluster_alert_configuration." + resourceName
	const kind = "cb_high_disk_throughput_usage"
	var beforeRemute time.Time

	scheduledFrom := alertMuteTime(time.Hour)
	scheduledUntil := alertMuteTime(24 * time.Hour)
	longerUntil := alertMuteTime(48 * time.Hour)
	pastFrom := alertMuteTime(-1 * time.Hour)
	shortUntil := alertMuteTime(2 * time.Minute)
	bulkUntil := alertMuteTime(24 * time.Hour)

	muteConfig := func(entries map[string]string) string {
		var lines []string
		for _, name := range alertSnakeCaseCatalog {
			if entry, ok := entries[name]; ok {
				lines = append(lines, fmt.Sprintf("    %s = { %s }", name, entry))
			}
		}
		return testAccClusterAlertConfigurationConfig(resourceName, "  muted_alerts = {\n"+strings.Join(lines, "\n")+"\n  }")
	}
	timed := func(from, until string) string {
		entry := `action = "mute"`
		if from != "" {
			entry += fmt.Sprintf(`, muted_from = %q`, from)
		}
		return entry + fmt.Sprintf(`, muted_until = %q`, until)
	}
	const disable = `action = "mute"`
	const enable = `action = "unmute"`
	only := func(entry string) map[string]string { return map[string]string{kind: entry} }
	all := func(entry string) map[string]string {
		m := make(map[string]string, len(alertSnakeCaseCatalog))
		for _, name := range alertSnakeCaseCatalog {
			m[name] = entry
		}
		return m
	}
	onServer := func(check func(map[string]alertsapi.AlertMuteState) error) resource.TestCheckFunc {
		return testAccCheckAlertConfigurationOnServer(clusterLevel, func(c *alertsapi.GetAlertConfigurationResponse) error {
			return check(c.MutedAlerts)
		})
	}

	mixed := all(timed("", bulkUntil))
	for _, name := range []string{"cpu_usage_high", "disk_usage_high", "syncgateway_mem_usage"} {
		mixed[name] = disable
	}

	return []resource.TestStep{
		// A mute that expires while Terraform manages it must not cause drift, and a later
		// window can then be set (MUT-03, config side). These run first, because the window is
		// fixed when the steps are built.
		{Config: muteConfig(only(timed("", shortUntil)))},
		{
			PreConfig: func() { sleepUntil(shortUntil, 30*time.Second) },
			Config:    muteConfig(only(timed("", shortUntil))),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
			},
		},
		{
			Config: muteConfig(only(timed("", scheduledUntil))),
			Check: onServer(func(m map[string]alertsapi.AlertMuteState) error {
				return expectAlertMute(m, kind, "mute", "", scheduledUntil)
			}),
		},
		// MUT-05: a scheduled mute, From in the future.
		{
			Config: muteConfig(only(timed(scheduledFrom, scheduledUntil))),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts."+kind+".muted_from", scheduledFrom),
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts."+kind+".muted_until", scheduledUntil),
				onServer(func(m map[string]alertsapi.AlertMuteState) error {
					return expectAlertMute(m, kind, "mute", scheduledFrom, scheduledUntil)
				}),
			),
		},
		// MUT-08: muting again with a longer window replaces the window. AUD-01, AUD-04: the
		// change through the v4 API appears in the cluster's Activity Log.
		{
			PreConfig: func() { beforeRemute = time.Now().UTC().Add(-5 * time.Second) },
			Config:    muteConfig(only(timed(scheduledFrom, longerUntil))),
			Check: resource.ComposeAggregateTestCheckFunc(
				onServer(func(m map[string]alertsapi.AlertMuteState) error {
					return expectAlertMute(m, kind, "mute", scheduledFrom, longerUntil)
				}),
				testAccCheckAlertConfigurationActivityLoggedSince(globalClusterId, &beforeRemute),
			),
		},
		// DIS-22: a timed mute turned into a disable (a mute with no end); DIS-23, MUT-30: the
		// overrides list shows the disable.
		{
			Config: muteConfig(only(disable)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckNoResourceAttr(resourceReference, "muted_alerts."+kind+".muted_until"),
				onServer(func(m map[string]alertsapi.AlertMuteState) error {
					return expectAlertMute(m, kind, "mute", "", "")
				}),
				testAccCheckAlertOverrides(func(overrides []alertsapi.AlertConfigurationOverride) error {
					o := findAlertOverride(overrides, "cluster", globalClusterId, kind)
					if o == nil || o.MutedAlert == nil || o.MutedAlert.Action != "mute" || o.MutedAlert.MutedUntil != nil {
						return fmt.Errorf("overrides do not list the disable of %s on the cluster: %+v", kind, o)
					}
					return nil
				}),
			),
		},
		// MUT-25, DIS-22: enabling leaves the alert unmuted, with no window left over, and the
		// disable leaves the overrides list.
		{
			Config: muteConfig(only(enable)),
			Check: resource.ComposeAggregateTestCheckFunc(
				onServer(func(m map[string]alertsapi.AlertMuteState) error {
					return expectAlertMute(m, kind, "unmute", "", "")
				}),
				testAccCheckAlertOverrides(func(overrides []alertsapi.AlertConfigurationOverride) error {
					if o := findAlertOverride(overrides, "cluster", globalClusterId, kind); o != nil && o.MutedAlert != nil && o.MutedAlert.Action == "mute" {
						return fmt.Errorf("overrides still list %s as disabled on the cluster after enabling it", kind)
					}
					return nil
				}),
			),
		},
		// MUT-23: From in the past is treated as starting now. Whatever Capella stores for From,
		// the apply must succeed and the next plan must be empty.
		{
			Config: muteConfig(only(timed(pastFrom, scheduledUntil))),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
			},
			Check: onServer(func(m map[string]alertsapi.AlertMuteState) error {
				return expectAlertMute(m, kind, "mute", "*", scheduledUntil)
			}),
		},
		// MUT-09, DIS-20: every alert in the catalog muted in one apply.
		{
			Config: muteConfig(all(timed("", bulkUntil))),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "muted_alerts.%", fmt.Sprint(len(alertSnakeCaseCatalog))),
				onServer(func(m map[string]alertsapi.AlertMuteState) error {
					for _, name := range alertSnakeCaseCatalog {
						if err := expectAlertMute(m, name, "mute", "", bulkUntil); err != nil {
							return err
						}
					}
					return nil
				}),
			),
		},
		// DIS-21: a mixed change disables only the selected alerts; the others keep their window.
		{
			Config: muteConfig(mixed),
			Check: onServer(func(m map[string]alertsapi.AlertMuteState) error {
				for name, entry := range mixed {
					until := bulkUntil
					if entry == disable {
						until = ""
					}
					if err := expectAlertMute(m, name, "mute", "", until); err != nil {
						return err
					}
				}
				return nil
			}),
		},
		// DIS-20, MUT-09: every alert enabled in one apply.
		{
			Config: muteConfig(all(enable)),
			Check: onServer(func(m map[string]alertsapi.AlertMuteState) error {
				for _, name := range alertSnakeCaseCatalog {
					if err := expectAlertMute(m, name, "unmute", "", ""); err != nil {
						return err
					}
				}
				return nil
			}),
		},
	}
}

// alertConfigurationMutePrecedenceSteps checks which mute is in effect when scopes disagree, by
// reading the resolved mutes at cluster and project scope. Only mutes are used, because they can
// be cleared again.
// Plan: MUT-15, MUT-16, MUT-17, MUT-18, MUT-20, DIS-03, DIS-04, DIS-17, DIS-19, HIE-02 (mutes).
func alertConfigurationMutePrecedenceSteps(clusterLevel, projectLevel, orgLevel string) []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_alert_precedence_")
	const kind = "cb_high_disk_write_queue"
	day := alertMuteTime(24 * time.Hour)
	threeDays := alertMuteTime(72 * time.Hour)

	level := func(resourceType, body string) string {
		if body == "" {
			return ""
		}
		ids := fmt.Sprintf("  organization_id = %q\n", globalOrgId)
		if resourceType != "organization" {
			ids += fmt.Sprintf("  project_id      = %q\n", globalProjectId)
		}
		if resourceType == "cluster" {
			ids += fmt.Sprintf("  cluster_id      = %q\n", globalClusterId)
		}
		return fmt.Sprintf("\nresource \"couchbase-capella_%s_alert_configuration\" %q {\n%s  muted_alerts = {\n%s\n  }\n}\n",
			resourceType, resourceName, ids, body)
	}
	entry := func(name, fields string) string { return fmt.Sprintf("    %s = { %s }", name, fields) }
	config := func(org, project, cluster string) string {
		return globalProviderBlock + level("organization", org) + level("project", project) + level("cluster", cluster)
	}
	resolved := func(levelPath string, check func(map[string]alertsapi.AlertMuteState) error) resource.TestCheckFunc {
		return testAccCheckAlertConfigurationOnServer(levelPath, func(c *alertsapi.GetAlertConfigurationResponse) error {
			if err := check(c.ResolvedMutedAlerts); err != nil {
				return fmt.Errorf("resolved mutes at %s: %w", levelPath, err)
			}
			return nil
		})
	}
	muted := func(levelPath, name, until string) resource.TestCheckFunc {
		return resolved(levelPath, func(m map[string]alertsapi.AlertMuteState) error {
			return expectAlertMute(m, name, "mute", "*", until)
		})
	}
	notMuted := func(levelPath, name string) resource.TestCheckFunc {
		return resolved(levelPath, func(m map[string]alertsapi.AlertMuteState) error {
			if mute, ok := m[name]; ok && mute.Action == "mute" {
				return fmt.Errorf("%s is muted (%+v), want not muted", name, mute)
			}
			return nil
		})
	}

	return []resource.TestStep{
		// DIS-03, DIS-04: disabled at organization scope, so every project and cluster is.
		{
			Config: config(entry(kind, `action = "mute"`), "", ""),
			Check:  resource.ComposeAggregateTestCheckFunc(muted(projectLevel, kind, ""), muted(clusterLevel, kind, "")),
		},
		// MUT-16, DIS-17: organization muted, cluster explicitly enabled: the cluster delivers,
		// and the project, standing for the other clusters, stays muted.
		{
			Config: config(entry(kind, `action = "mute"`), "", entry(kind, `action = "unmute"`)),
			Check:  resource.ComposeAggregateTestCheckFunc(notMuted(clusterLevel, kind), muted(projectLevel, kind, "")),
		},
		// MUT-18: organization muted, project enabled, cluster muted: the most specific wins.
		{
			Config: config(entry(kind, `action = "mute"`), entry(kind, `action = "unmute"`), entry(kind, fmt.Sprintf(`action = "mute", muted_until = %q`, day))),
			Check:  resource.ComposeAggregateTestCheckFunc(muted(clusterLevel, kind, day), notMuted(projectLevel, kind)),
		},
		// MUT-20: cluster muted with a later end than the organization: the cluster's window
		// governs the cluster.
		{
			Config: config(entry(kind, fmt.Sprintf(`action = "mute", muted_until = %q`, day)), "", entry(kind, fmt.Sprintf(`action = "mute", muted_until = %q`, threeDays))),
			Check:  resource.ComposeAggregateTestCheckFunc(muted(clusterLevel, kind, threeDays), muted(projectLevel, kind, day)),
		},
		// DIS-19: cluster disabled, then the organization enables the alert: the cluster stays
		// disabled.
		{
			Config: config(entry(kind, `action = "unmute"`), "", entry(kind, `action = "mute"`)),
			Check:  resource.ComposeAggregateTestCheckFunc(muted(clusterLevel, kind, ""), notMuted(projectLevel, kind)),
		},
		// MUT-17: project muted, cluster enabled: the cluster delivers, the project stays muted.
		{
			Config: config("", entry(kind, `action = "mute"`), entry(kind, `action = "unmute"`)),
			Check:  resource.ComposeAggregateTestCheckFunc(notMuted(clusterLevel, kind), muted(projectLevel, kind, "")),
		},
		// MUT-15: the organization mutes several alerts and enables one of them.
		{
			Config: config(strings.Join([]string{
				entry(kind, `action = "mute"`),
				entry("disk_iops_usage", `action = "mute"`),
				entry("cb_high_disk_throughput_usage", `action = "unmute"`),
			}, "\n"), "", ""),
			Check: resource.ComposeAggregateTestCheckFunc(
				muted(orgLevel, kind, ""),
				muted(orgLevel, "disk_iops_usage", ""),
				notMuted(orgLevel, "cb_high_disk_throughput_usage"),
			),
		},
	}
}

// alertConfigurationRecipientsSteps covers the recipient list's lifecycle and what Terraform
// leaves alone when it does not manage recipients.
// Plan: RCP-01, RCP-03 (case), RCP-06 (config side), HIE-08.
func alertConfigurationRecipientsSteps(t *testing.T, clusterLevel string) []resource.TestStep {
	recipient1, recipient2 := alertTestRecipients(t)
	resourceName := randomStringWithPrefix("tf_acc_alert_recipients_")
	resourceReference := "couchbase-capella_cluster_alert_configuration." + resourceName
	config := func(body string) string { return testAccClusterAlertConfigurationConfig(resourceName, body) }
	var configsBefore *alertsapi.GetAlertConfigurationResponse

	return []resource.TestStep{
		{
			PreConfig: func() {
				var err error
				if configsBefore, err = getAlertConfigurationForTest(clusterLevel); err != nil {
					panic(err)
				}
			},
			Config: config(fmt.Sprintf("  recipients = [%q, %q]", recipient1, recipient2)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "2"),
				testAccCheckAlertConfigurationOnServer(clusterLevel, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return expectAlertRecipients(c, recipient1, recipient2)
				}),
			),
		},
		// RCP-03: the same address in a different case is the same recipient. Whatever case
		// Capella stores, the apply must succeed without drift.
		{
			Config: config(fmt.Sprintf("  recipients = [%q, %q]", strings.ToUpper(recipient1), recipient2)),
			Check: testAccCheckAlertConfigurationOnServer(clusterLevel, func(c *alertsapi.GetAlertConfigurationResponse) error {
				return expectAlertRecipients(c, recipient1, recipient2)
			}),
		},
		// RCP-06: removing a recipient.
		{
			Config: config(fmt.Sprintf("  recipients = [%q]", recipient1)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "1"),
				testAccCheckAlertConfigurationOnServer(clusterLevel, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return expectAlertRecipients(c, recipient1)
				}),
			),
		},
		// An empty list clears the cluster's own recipients, so it inherits again.
		{
			Config: config("  recipients = []"),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceReference, "recipients.#", "0"),
				testAccCheckAlertConfigurationOnServer(clusterLevel, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return expectAlertRecipients(c)
				}),
			),
		},
		{Config: config(fmt.Sprintf("  recipients = [%q]", recipient1))},
		// HIE-08: a resource that manages only mutes leaves the recipients and the configs that
		// Capella already has untouched.
		{
			Config: config(`
  muted_alerts = {
    cb_high_disk_write_queue = {
      action      = "mute"
      muted_until = "` + alertMuteTime(24*time.Hour) + `"
    }
  }`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckNoResourceAttr(resourceReference, "recipients.#"),
				testAccCheckAlertConfigurationOnServer(clusterLevel, func(c *alertsapi.GetAlertConfigurationResponse) error {
					return firstAlertError(
						expectAlertRecipients(c, recipient1),
						sameAlertConfigs(configsBefore.Configs, c.Configs),
					)
				}),
			),
		},
	}
}

// alertRBACCase is one role's expected access, from the PRD's access table and the v4 API
// testing proof.
type alertRBACCase struct {
	name              string
	organizationRoles []string
	projectRoles      []string
	readOrganization  bool
	readProject       bool
	readCluster       bool
	writeProject      bool
	writeCluster      bool
	// seesOnlyProject means every alert the role lists is in the shared project (RBAC-10).
	// listAlertsDenied means the role may not list alerts at all (RBAC-05, RES-A-09).
	seesOnlyProject  bool
	listAlertsDenied bool
}

// runAlertConfigurationRBACSubtests tries every role against every scope with a real API key.
// Successful writes are mutes, which destroy clears.
// Plan: RBAC-01, RBAC-02, RBAC-03, RBAC-04, RBAC-05, RBAC-06, RBAC-10, RES-A-09, HIE-14.
func runAlertConfigurationRBACSubtests(t *testing.T, clusterLevel, projectLevel string) {
	cases := []alertRBACCase{
		// RBAC-05, RES-A-09: a role with no project may read the organization's configuration
		// but not list alerts; Capella answers 403 (confirmed as intended, 2026-11-26).
		{name: "Organization Member", organizationRoles: []string{"organizationMember"}, readOrganization: true, listAlertsDenied: true},
		{name: "Project Creator", organizationRoles: []string{"projectCreator"}, readOrganization: true, listAlertsDenied: true},
		{name: "Project Owner", organizationRoles: []string{"organizationMember"}, projectRoles: []string{"projectOwner"},
			readOrganization: true, readProject: true, readCluster: true, writeProject: true, writeCluster: true, seesOnlyProject: true},
		// RBAC-04: a project manager updates both the project and its clusters (confirmed as
		// intended, 2026-11-26).
		{name: "Project Manager", organizationRoles: []string{"organizationMember"}, projectRoles: []string{"projectManager"},
			readOrganization: true, readProject: true, readCluster: true, writeProject: true, writeCluster: true, seesOnlyProject: true},
		{name: "Project Viewer", organizationRoles: []string{"organizationMember"}, projectRoles: []string{"projectViewer"},
			readOrganization: true, readProject: true, readCluster: true, seesOnlyProject: true},
		{name: "Project Data Reader Writer", organizationRoles: []string{"organizationMember"}, projectRoles: []string{"projectDataReaderWriter"},
			readOrganization: true, readProject: true, readCluster: true, seesOnlyProject: true},
		{name: "Project Data Reader", organizationRoles: []string{"organizationMember"}, projectRoles: []string{"projectDataReader"},
			readOrganization: true, readProject: true, readCluster: true, seesOnlyProject: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			projectRoles := map[string][]string{}
			if len(c.projectRoles) > 0 {
				projectRoles[globalProjectId] = c.projectRoles
			}
			key := createAlertTestAPIKey(t, c.organizationRoles, projectRoles)
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps:                    alertRBACSteps(key, c, clusterLevel, projectLevel),
			})
		})
	}

	// HIE-14, AV-147556, AV-153487: a change made with an API key records the key as who made it,
	// for mutes and for configs, on the level and in the overrides list. Run 2026-11-26:
	// updatedBy was empty for both.
	t.Run("Mute Records The API Key That Set It", func(t *testing.T) {
		key := createAlertTestAPIKey(t, []string{"organizationMember"}, map[string][]string{globalProjectId: {"projectOwner"}})
		resourceName := randomStringWithPrefix("tf_acc_alert_rbac_setby_")
		// Pick a frequency different from the current one, so that a request is sent.
		frequency := int64(5)
		if current, err := getAlertConfigurationForTest(clusterLevel); err == nil {
			if entry, ok := current.Configs["CBHardOutOfMemoryErrors"]; ok && entry.FrequencyHours != nil && *entry.FrequencyHours == frequency {
				frequency = 7
			}
		}
		config := alertRBACProviderBlock(key.token) + fmt.Sprintf(`
resource "couchbase-capella_cluster_alert_configuration" %q {
  provider = couchbase-capella.rbac
%s
  configs = {
    CBHardOutOfMemoryErrors = {
      frequency_hours = %d
    }
  }
  muted_alerts = {
    cb_high_disk_write_queue = {
      action      = "mute"
      muted_until = %q
    }
  }
}
`, resourceName, alertRBACIds("cluster", globalProjectId, globalClusterId), frequency, alertMuteTime(24*time.Hour))
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps: []resource.TestStep{{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAlertConfigurationOnServer(clusterLevel, func(cfg *alertsapi.GetAlertConfigurationResponse) error {
						return expectAlertMuteSetBy(cfg.MutedAlerts, "cb_high_disk_write_queue", key)
					}),
					testAccCheckAlertConfigurationOnServer(clusterLevel, func(cfg *alertsapi.GetAlertConfigurationResponse) error {
						entry := cfg.Configs["CBHardOutOfMemoryErrors"]
						if entry.UpdatedBy == nil || (*entry.UpdatedBy != key.id && *entry.UpdatedBy != key.name) {
							return fmt.Errorf("CBHardOutOfMemoryErrors config updated_by = %v, want the API key %s (%s)", derefString(entry.UpdatedBy), key.name, key.id)
						}
						return nil
					}),
					testAccCheckAlertOverrides(func(overrides []alertsapi.AlertConfigurationOverride) error {
						o := findAlertOverride(overrides, "cluster", globalClusterId, "cb_high_disk_write_queue")
						if o == nil || o.MutedAlert == nil {
							return nil // the overrides list shows disables only, not timed mutes
						}
						return expectAlertMuteSetBy(map[string]alertsapi.AlertMuteState{o.AlertKind: *o.MutedAlert}, o.AlertKind, key)
					}),
				),
			}},
		})
	})

	// RBAC-03: a project creator updates the projects it created, and only those.
	t.Run("Project Creator On Its Own Project", func(t *testing.T) {
		key := createAlertTestAPIKey(t, []string{"projectCreator"}, nil)
		ownProjectId := createAlertTestProject(t, key.token)
		resourceName := randomStringWithPrefix("tf_acc_alert_rbac_pc_")
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps: []resource.TestStep{
				{
					Config: alertRBACProviderBlock(key.token) + alertRBACWriteConfig(resourceName, "project", ownProjectId, ""),
					Check: testAccCheckAlertConfigurationOnServer(alertLevelPath(globalOrgId, ownProjectId, ""), func(cfg *alertsapi.GetAlertConfigurationResponse) error {
						return expectAlertMute(cfg.MutedAlerts, "cb_high_disk_write_queue", "mute", "*", "*")
					}),
				},
				{
					Config:      alertRBACProviderBlock(key.token) + alertRBACWriteConfig(resourceName, "project", globalProjectId, ""),
					ExpectError: alertAccessDenied,
				},
			},
		})
	})

	// RBAC-06: the highest privilege applies per resource: owner of one project and data
	// reader on another.
	t.Run("Owner On One Project Reader On Another", func(t *testing.T) {
		otherProjectId := createAlertTestProject(t, globalToken)
		key := createAlertTestAPIKey(t, []string{"organizationMember"}, map[string][]string{
			globalProjectId: {"projectOwner"},
			otherProjectId:  {"projectDataReader"},
		})
		resourceName := randomStringWithPrefix("tf_acc_alert_rbac_mixed_")
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
			Steps: []resource.TestStep{
				{
					Config: alertRBACProviderBlock(key.token) + alertRBACWriteConfig(resourceName, "project", globalProjectId, ""),
					Check: testAccCheckAlertConfigurationOnServer(projectLevel, func(cfg *alertsapi.GetAlertConfigurationResponse) error {
						return expectAlertMute(cfg.MutedAlerts, "cb_high_disk_write_queue", "mute", "*", "*")
					}),
				},
				{
					Config:      alertRBACProviderBlock(key.token) + alertRBACWriteConfig(resourceName, "project", otherProjectId, ""),
					ExpectError: alertAccessDenied,
				},
				{
					Config: alertRBACProviderBlock(key.token) + alertRBACReadConfig(resourceName, "project", otherProjectId, ""),
					Check:  resource.TestCheckResourceAttr("data.couchbase-capella_project_alert_configuration."+resourceName, "project_id", otherProjectId),
				},
			},
		})
	})
}

func alertRBACSteps(key alertTestAPIKey, c alertRBACCase, clusterLevel, projectLevel string) []resource.TestStep {
	resourceName := randomStringWithPrefix("tf_acc_alert_rbac_")
	provider := alertRBACProviderBlock(key.token)
	read := func(level string, allowed bool) resource.TestStep {
		step := resource.TestStep{Config: provider + alertRBACReadConfig(resourceName, level, globalProjectId, globalClusterId)}
		if allowed {
			step.Check = resource.TestCheckResourceAttr("data.couchbase-capella_"+level+"_alert_configuration."+resourceName, "organization_id", globalOrgId)
		} else {
			step.ExpectError = alertAccessDenied
		}
		return step
	}
	write := func(level, levelPath string, allowed bool) resource.TestStep {
		step := resource.TestStep{Config: provider + alertRBACWriteConfig(resourceName, level, globalProjectId, globalClusterId)}
		if allowed {
			step.Check = testAccCheckAlertConfigurationOnServer(levelPath, func(cfg *alertsapi.GetAlertConfigurationResponse) error {
				return expectAlertMute(cfg.MutedAlerts, "cb_high_disk_write_queue", "mute", "*", "*")
			})
		} else {
			step.ExpectError = alertAccessDenied
		}
		return step
	}

	steps := []resource.TestStep{
		read("organization", c.readOrganization),
		read("project", c.readProject),
		read("cluster", c.readCluster),
		// No role but organization owner updates the organization.
		write("organization", "", false),
		write("project", projectLevel, c.writeProject),
		write("cluster", clusterLevel, c.writeCluster),
	}

	alertsReference := "data.couchbase-capella_alerts." + resourceName
	alertsStep := resource.TestStep{Config: provider + fmt.Sprintf(`
data "couchbase-capella_alerts" %q {
  provider        = couchbase-capella.rbac
  organization_id = %q
}
`, resourceName, globalOrgId)}
	switch {
	case c.listAlertsDenied:
		alertsStep.ExpectError = alertAccessDenied
	case c.seesOnlyProject:
		alertsStep.Check = testAccCheckEveryListedAlert(alertsReference, map[string][]string{"project_id": {globalProjectId}})
	}
	return append(steps, alertsStep)
}

func alertRBACProviderBlock(token string) string {
	return globalProviderBlock + fmt.Sprintf(`
provider "couchbase-capella" {
  alias                = "rbac"
  host                 = var.host
  authentication_token = %q
}
`, token)
}

// alertRBACIds returns the ID attributes of a level's resource or data source.
func alertRBACIds(level, projectId, clusterId string) string {
	ids := fmt.Sprintf("  organization_id = %q\n", globalOrgId)
	if level != "organization" {
		ids += fmt.Sprintf("  project_id      = %q\n", projectId)
	}
	if level == "cluster" {
		ids += fmt.Sprintf("  cluster_id      = %q\n", clusterId)
	}
	return ids
}

func alertRBACReadConfig(resourceName, level, projectId, clusterId string) string {
	return fmt.Sprintf("\ndata \"couchbase-capella_%s_alert_configuration\" %q {\n  provider = couchbase-capella.rbac\n%s}\n",
		level, resourceName, alertRBACIds(level, projectId, clusterId))
}

func alertRBACWriteConfig(resourceName, level, projectId, clusterId string) string {
	return fmt.Sprintf(`
resource "couchbase-capella_%s_alert_configuration" %q {
  provider = couchbase-capella.rbac
%s
  muted_alerts = {
    cb_high_disk_write_queue = {
      action      = "mute"
      muted_until = %q
    }
  }
}
`, level, resourceName, alertRBACIds(level, projectId, clusterId), alertMuteTime(24*time.Hour))
}

// TestAccAlertConfigurationHierarchy checks inheritance in a project and cluster the test
// creates, so that threshold and frequency overrides, which Capella cannot remove, go away with
// them.
// Plan: HIE-01, HIE-02, HIE-03, HIE-05, HIE-06, HIE-07, HIE-09, HIE-10, THR-16, THR-17, FRQ-15,
// FRQ-16, DIS-04, DIS-18, DIS-27, and App Service alerts configured at project scope (THR-02).
func TestAccAlertConfigurationHierarchy(t *testing.T) {
	projectName := randomStringWithPrefix("tf_acc_alert_hier_")
	clusterName := randomStringWithPrefix("tf_acc_alert_hier_cluster_")
	name := randomStringWithPrefix("tf_acc_alert_hier_cfg_")
	cidr := generateRandomCIDR()
	clusterReference := "couchbase-capella_cluster." + clusterName
	clusterDatasource := "data.couchbase-capella_cluster_alert_configuration." + name
	projectDatasource := "data.couchbase-capella_project_alert_configuration." + name
	overridesDatasource := "data.couchbase-capella_alert_configuration_overrides." + name
	var projectId, clusterId string

	type projectValues struct{ cpuWarning int }
	type clusterValues struct {
		present      bool
		cpuFrequency int
		diskCritical int
	}

	config := func(p projectValues, c clusterValues) string {
		clusterAlerts := ""
		if c.present {
			clusterAlerts = fmt.Sprintf(`
resource "couchbase-capella_cluster_alert_configuration" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = couchbase-capella_project.%[3]s.id
  cluster_id      = couchbase-capella_cluster.%[4]s.id

  configs = {
    cpu_usage_high = {
      frequency_hours = %[5]d
    }
    disk_usage_high = {
      critical_threshold = %[6]d
    }
  }

  muted_alerts = {
    indexer_low_resident_percentage = {
      action = "unmute"
    }
  }
}
`, name, globalOrgId, projectName, clusterName, c.cpuFrequency, c.diskCritical)
		}
		dependsOn := "couchbase-capella_project_alert_configuration." + name
		if c.present {
			dependsOn += ", couchbase-capella_cluster_alert_configuration." + name
		}
		return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[2]s" {
  organization_id = "%[3]s"
  name            = "%[2]s"
}

resource "couchbase-capella_project_alert_configuration" "%[4]s" {
  organization_id = "%[3]s"
  project_id      = couchbase-capella_project.%[2]s.id

  configs = {
    cpu_usage_high = {
      warning_threshold  = %[5]d
      critical_threshold = 95
      frequency_hours    = 6
    }
    disk_usage_high = {
      warning_threshold = 80
    }
    syncgateway_cpu_usage_high = {
      warning_threshold  = 80
      critical_threshold = 98
    }
    syncgateway_mem_usage = {
      warning_threshold  = 89
      critical_threshold = 90
    }
  }

  muted_alerts = {
    indexer_low_resident_percentage = {
      action = "mute"
    }
  }
}

# HIE-05: the cluster is deployed into a project that already has overrides.
resource "couchbase-capella_cluster" "%[6]s" {
  organization_id = "%[3]s"
  project_id      = couchbase-capella_project.%[2]s.id
  name            = "%[6]s"
  couchbase_server = {
    version = var.server_version
  }
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[7]s"
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
    plan     = "developer pro"
    timezone = "PT"
  }
  depends_on = [couchbase-capella_project_alert_configuration.%[4]s]
}
%[8]s
data "couchbase-capella_cluster_alert_configuration" "%[4]s" {
  organization_id = "%[3]s"
  project_id      = couchbase-capella_project.%[2]s.id
  cluster_id      = couchbase-capella_cluster.%[6]s.id
  depends_on      = [%[9]s]
}

data "couchbase-capella_project_alert_configuration" "%[4]s" {
  organization_id = "%[3]s"
  project_id      = couchbase-capella_project.%[2]s.id
  depends_on      = [%[9]s]
}

data "couchbase-capella_alert_configuration_overrides" "%[4]s" {
  organization_id = "%[3]s"
  depends_on      = [%[9]s]
}
`, globalProviderBlock, projectName, globalOrgId, name, p.cpuWarning, clusterName, cidr, clusterAlerts, dependsOn)
	}

	captureIds := func(s *terraform.State) error {
		project, err := datasourceAttributes(s, "couchbase-capella_project."+projectName)
		if err != nil {
			return err
		}
		cluster, err := datasourceAttributes(s, clusterReference)
		if err != nil {
			return err
		}
		projectId, clusterId = project["id"], cluster["id"]
		return nil
	}
	resolved := func(alert, field, value string) resource.TestCheckFunc {
		return resource.TestCheckResourceAttr(clusterDatasource, "resolved_configs."+alert+"."+field, value)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			func(s *terraform.State) error {
				if _, ok := s.RootModule().Resources[clusterReference]; !ok {
					return nil // the cluster was never created
				}
				return testAccCheckClusterDestroyed(t, clusterReference)(s)
			},
			// HIE-07, DIS-27: deleting the project and cluster removes their overrides.
			testAccCheckAlertOverrides(func(overrides []alertsapi.AlertConfigurationOverride) error {
				for _, o := range overrides {
					if (o.ProjectId != nil && *o.ProjectId == projectId) || (o.ResourceId != nil && (*o.ResourceId == clusterId || *o.ResourceId == projectId)) {
						return fmt.Errorf("override of %s at %s scope outlived its deleted project or cluster", o.AlertKind, o.Scope)
					}
				}
				return nil
			}),
		),
		Steps: []resource.TestStep{
			// HIE-01, HIE-05, THR-16, FRQ-15, DIS-04: a new cluster has no overrides of its own
			// and takes everything from its project.
			{
				Config: config(projectValues{cpuWarning: 85}, clusterValues{}),
				Check: resource.ComposeAggregateTestCheckFunc(
					captureIds,
					resource.TestCheckResourceAttr(projectDatasource, "configs.cpu_usage_high.warning_threshold", "85"),
					resource.TestCheckResourceAttr(projectDatasource, "configs.syncgateway_cpu_usage_high.critical_threshold", "98"),
					resource.TestCheckResourceAttr(projectDatasource, "configs.syncgateway_mem_usage.warning_threshold", "89"),
					resource.TestCheckResourceAttr(clusterDatasource, "configs.%", "0"),
					resource.TestCheckResourceAttr(clusterDatasource, "muted_alerts.%", "0"),
					resolved("cpu_usage_high", "warning_threshold", "85"),
					resolved("cpu_usage_high", "critical_threshold", "95"),
					resolved("cpu_usage_high", "frequency_hours", "6"),
					resolved("disk_usage_high", "warning_threshold", "80"),
					resource.TestCheckResourceAttr(clusterDatasource, "resolved_muted_alerts.indexer_low_resident_percentage.action", "mute"),
				),
			},
			// HIE-02, HIE-03, FRQ-16, THR-17, DIS-18: the cluster sets only some fields; each
			// field resolves to its most specific value, and the cluster re-enables an alert its
			// project disabled. HIE-09: the cluster's new override is listed.
			{
				Config: config(projectValues{cpuWarning: 85}, clusterValues{present: true, cpuFrequency: 3, diskCritical: 88}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resolved("cpu_usage_high", "frequency_hours", "3"),
					resolved("cpu_usage_high", "warning_threshold", "85"),
					resolved("cpu_usage_high", "critical_threshold", "95"),
					resolved("disk_usage_high", "warning_threshold", "80"),
					resolved("disk_usage_high", "critical_threshold", "88"),
					testAccCheckNotMutedInState(clusterDatasource, "resolved_muted_alerts", "indexer_low_resident_percentage"),
					resource.TestCheckResourceAttr(projectDatasource, "resolved_muted_alerts.indexer_low_resident_percentage.action", "mute"),
					func(s *terraform.State) error {
						return testAccCheckAlertOverrideListed(overridesDatasource, clusterId, "cpu_usage_high")(s)
					},
					func(s *terraform.State) error {
						prefix, err := findAlertOverrideInState(s, overridesDatasource, clusterId, "cpu_usage_high")
						if err != nil {
							return err
						}
						attributes, _ := datasourceAttributes(s, overridesDatasource)
						if got := attributes[prefix+"config.frequency_hours"]; got != "3" {
							return fmt.Errorf("cluster override of cpu_usage_high lists frequency_hours = %q, want 3", got)
						}
						return nil
					},
				),
			},
			// HIE-06: the project changes a value the cluster inherits; the cluster follows it and
			// keeps its own frequency. HIE-09: changing the cluster's own override lists the
			// previous and the current value.
			{
				Config: config(projectValues{cpuWarning: 82}, clusterValues{present: true, cpuFrequency: 3, diskCritical: 89}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resolved("cpu_usage_high", "warning_threshold", "82"),
					resolved("cpu_usage_high", "frequency_hours", "3"),
					resolved("disk_usage_high", "critical_threshold", "89"),
					func(s *terraform.State) error {
						prefix, err := findAlertOverrideInState(s, overridesDatasource, clusterId, "disk_usage_high")
						if err != nil {
							return err
						}
						attributes, _ := datasourceAttributes(s, overridesDatasource)
						current, previous := attributes[prefix+"config.critical_threshold"], attributes[prefix+"previous_config.critical_threshold"]
						if current != "89" || previous != "88" {
							return fmt.Errorf("cluster override of disk_usage_high lists critical_threshold %q, previously %q; want 89, previously 88", current, previous)
						}
						return nil
					},
				),
			},
			// HIE-10: setting the cluster back to the value it would inherit removes the override
			// from the overrides list.
			{
				Config: config(projectValues{cpuWarning: 82}, clusterValues{present: true, cpuFrequency: 6, diskCritical: 89}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resolved("cpu_usage_high", "frequency_hours", "6"),
					func(s *terraform.State) error {
						if _, err := findAlertOverrideInState(s, overridesDatasource, clusterId, "cpu_usage_high"); err == nil {
							return fmt.Errorf("the cluster's cpu_usage_high override is still listed after setting it to the inherited value")
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccAlertConfigurationPlanTimeValidation covers input the provider rejects before anything
// is sent to Capella. Every message was taken from terraform validate against the provider.
// Plan: FRQ-03 (decimals, text), THR-03 (non-numeric, decimal), MUT-22 (malformed times),
// RES-A-03 and the alert list filters.
func TestAccAlertConfigurationPlanTimeValidation(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_alert_plan_invalid_")
	cluster := func(body string) string { return testAccClusterAlertConfigurationConfig(resourceName, body) }
	alerts := func(body string) string {
		return fmt.Sprintf("%s\ndata \"couchbase-capella_alerts\" %q {\n  organization_id = %q\n%s\n}\n", globalProviderBlock, resourceName, globalOrgId, body)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Decimals are refused, though with a confusing "error in the provider" message.
			{Config: cluster("  configs = { cpu_usage_high = { frequency_hours = 1.5 } }"), ExpectError: regexp.MustCompile(`is\s+not\s+an\s+integer`)},
			{Config: cluster("  configs = { cpu_usage_high = { warning_threshold = 82.5 } }"), ExpectError: regexp.MustCompile(`is\s+not\s+an\s+integer`)},
			{Config: cluster(`  configs = { cpu_usage_high = { frequency_hours = "abc" } }`), ExpectError: regexp.MustCompile(`a\s+number\s+is\s+required`)},
			{Config: cluster(`  configs = { cpu_usage_high = { warning_threshold = "high" } }`), ExpectError: regexp.MustCompile(`a\s+number\s+is\s+required`)},
			{Config: cluster(`  muted_alerts = { cpu_usage_high = { muted_until = "2030-01-01T00:00:00Z" } }`), ExpectError: regexp.MustCompile(`attribute\s+"action"\s+is\s+required`)},
			{Config: cluster(`  muted_alerts = { cpu_usage_high = { action = "mute", muted_from = "2030-13-01T00:00:00Z" } }`), ExpectError: regexp.MustCompile(`Invalid\s+RFC3339\s+String\s+Value`)},
			{Config: cluster(`  muted_alerts = { cpu_usage_high = { action = "mute", muted_until = "2030-01-01 00:00:00" } }`), ExpectError: regexp.MustCompile(`Invalid\s+RFC3339\s+String\s+Value`)},
			// The alert list's filters. "open" is the UI's word for firing.
			{Config: alerts(`  status = ["open"]`), ExpectError: regexp.MustCompile(`value\s+must\s+be\s+one\s+of:\s+\["firing"\s+"resolved"\]`)},
			{Config: alerts(`  severity_levels = ["major"]`), ExpectError: regexp.MustCompile(`value\s+must\s+be\s+one\s+of`)},
			{Config: alerts(`  cluster_ids = ["not-a-uuid"]`), ExpectError: regexp.MustCompile(`must\s+be\s+a\s+valid\s+UUID`)},
			{Config: alerts(`  project_ids = ["not-a-uuid"]`), ExpectError: regexp.MustCompile(`must\s+be\s+a\s+valid\s+UUID`)},
			{Config: alerts(`  from = "yesterday"`), ExpectError: regexp.MustCompile(`Invalid\s+RFC3339\s+String\s+Value`)},
		},
	})
}

// TestAccDatasourceAlertsContract checks the alert list against the rules every row must follow,
// whatever is firing in the organization, and the limits on the time window.
// Plan: RES-A-01, RES-A-02, RES-A-04, RES-A-05, RES-A-06, RES-A-10, RES-U-04 (API side), LST-03,
// LST-04, LST-06.
func TestAccDatasourceAlertsContract(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_alerts_contract_")
	reference := "data.couchbase-capella_alerts." + resourceName
	now := time.Now().UTC().Truncate(time.Second)
	window := func(days int) (string, string) {
		return now.Add(-time.Duration(days) * 24 * time.Hour).Format(time.RFC3339), now.Format(time.RFC3339)
	}
	config := func(organizationId, body string) string {
		return fmt.Sprintf("%s\ndata \"couchbase-capella_alerts\" %q {\n  organization_id = %q\n%s\n}\n", globalProviderBlock, resourceName, organizationId, body)
	}
	// The 31-day limit is measured when Capella receives the request, so stay an hour inside it.
	from31 := now.Add(-31*24*time.Hour + time.Hour).Format(time.RFC3339)
	to31 := now.Format(time.RFC3339)
	from32, to32 := window(32)
	readFailed := regexp.MustCompile(`(?s)Could\s+not\s+read\s+alerts\s+in\s+organization`)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// RES-A-02, RES-A-04, RES-A-05, LST-03: both statuses by default, one row per alert,
			// only configurable kinds, and endsAt only on resolved rows.
			{
				Config: config(globalOrgId, ""),
				Check:  testAccCheckAlertListContract(reference, []string{"firing", "resolved"}),
			},
			// RES-A-01: resolved only.
			{
				Config: config(globalOrgId, `  status = ["resolved"]`),
				Check:  testAccCheckAlertListContract(reference, []string{"resolved"}),
			},
			{
				Config: config(globalOrgId, `  status = ["firing"]`),
				Check:  testAccCheckAlertListContract(reference, []string{"firing"}),
			},
			// RES-A-06, RES-U-04: a 31-day window is the most allowed.
			{
				Config: config(globalOrgId, fmt.Sprintf("  from = %q\n  to   = %q", from31, to31)),
				Check:  testAccCheckAlertListContract(reference, []string{"firing", "resolved"}),
			},
			{
				Config:      config(globalOrgId, fmt.Sprintf("  from = %q\n  to   = %q", from32, to32)),
				ExpectError: regexp.MustCompile(wrapTolerantPattern("Alerts are only available for the last 31 days.")),
			},
			// From after to.
			{
				Config:      config(globalOrgId, fmt.Sprintf("  from = %q\n  to   = %q", to31, from31)),
				ExpectError: readFailed,
			},
			// RES-A-10: another tenant's alerts.
			{
				Config:      config(uuid.NewString(), ""),
				ExpectError: readFailed,
			},
		},
	})
}

// TestAccClusterAlertConfigurationFreeTier checks that alert configuration is refused on a free
// tier cluster, like Basic (DIS-07, HIE-19). An organization may have one free tier cluster, so
// this runs only when ACC_ALERT_FREE_TIER is set, and fails if the organization already has one.
func TestAccClusterAlertConfigurationFreeTier(t *testing.T) {
	if os.Getenv("ACC_ALERT_FREE_TIER") == "" {
		t.Skip("set ACC_ALERT_FREE_TIER=1 to create a free tier cluster for this test")
	}
	clusterName := randomStringWithPrefix("tf_acc_alert_free_")
	name := randomStringWithPrefix("tf_acc_alert_free_cfg_")
	cluster := fmt.Sprintf(`
%[1]s

resource "couchbase-capella_free_tier_cluster" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  name            = "%[2]s"
  description     = "Alert configuration acceptance test"

  cloud_provider = {
    type   = "aws"
    region = "us-east-2"
    cidr   = "%[5]s"
  }
}
`, globalProviderBlock, clusterName, globalOrgId, globalProjectId, generateRandomCIDR())
	alerts := func(body string) string {
		return cluster + fmt.Sprintf(`
resource "couchbase-capella_cluster_alert_configuration" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = couchbase-capella_free_tier_cluster.%[4]s.id
%[5]s
}
`, name, globalOrgId, globalProjectId, clusterName, body)
	}
	refused := regexp.MustCompile("(?i)" + wrapTolerantPattern(basicPlanAlertConfigurationMessage))
	var freeTierId string

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		// The organization has one free tier slot, so confirm with Capella that it was freed.
		CheckDestroy: func(_ *terraform.State) error {
			if freeTierId == "" {
				return nil
			}
			endpoint := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/freeTier/%s", globalHost, globalOrgId, globalProjectId, freeTierId)
			cfg := api.EndpointCfg{Url: endpoint, Method: http.MethodGet, SuccessStatus: http.StatusOK}
			_, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
			if err == nil {
				return fmt.Errorf("free tier cluster %s still exists after destroy", freeTierId)
			}
			if notFound, message := api.CheckResourceNotFoundError(err); !notFound {
				return fmt.Errorf("could not confirm free tier cluster %s was deleted: %s", freeTierId, message)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: cluster,
				Check: func(s *terraform.State) error {
					attributes, err := datasourceAttributes(s, "couchbase-capella_free_tier_cluster."+clusterName)
					if err == nil {
						freeTierId = attributes["id"]
					}
					return err
				},
			},
			{Config: alerts("  configs = { cpu_usage_high = { warning_threshold = 85 } }"), ExpectError: refused},
			{Config: alerts(`  muted_alerts = { cpu_usage_high = { action = "mute" } }`), ExpectError: refused},
		},
	})
}

// alertActivityLogEvent matches the key or summary of an Activity Log entry for a change to
// alert configuration. The event key is not documented; this is the plan's AUD-01 wording.
var alertActivityLogEvent = regexp.MustCompile(`(?i)(mute|alert[ _-]?config)`)

// testAccCheckAlertConfigurationActivityLoggedSince requires an alert configuration entry in the
// cluster's Activity Log emitted after *since, so that an earlier entry cannot satisfy it. The
// PRD's target is one second; this polls for up to a minute.
func testAccCheckAlertConfigurationActivityLoggedSince(clusterId string, since *time.Time) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		query := url.Values{}
		query.Set("clusterIds", clusterId)
		query.Set("from", since.Format(time.RFC3339))
		query.Set("perPage", "100")
		endpoint := fmt.Sprintf("%s/v4/organizations/%s/events?%s", globalHost, globalOrgId, query.Encode())
		cfg := api.EndpointCfg{Url: endpoint, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		deadline := time.Now().Add(time.Minute)
		for {
			response, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
			if err != nil {
				return fmt.Errorf("could not read the cluster's Activity Log: %w", err)
			}
			var events api.GetEventsResponse
			if err := json.Unmarshal(response.Body, &events); err != nil {
				return fmt.Errorf("could not read the cluster's Activity Log: %w", err)
			}
			for _, event := range events.Data {
				summary := ""
				if event.Summary != nil {
					summary = *event.Summary
				}
				if event.Timestamp.After(*since) && (alertActivityLogEvent.MatchString(event.Key) || alertActivityLogEvent.MatchString(summary)) {
					return nil
				}
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("no alert configuration entry in the cluster's Activity Log since %s", since.Format(time.RFC3339))
			}
			time.Sleep(5 * time.Second)
		}
	}
}

// testAccCheckAlertListContract checks every listed alert: a status from statuses, a kind from
// the catalog, one row per alert ID, and endsAt set, after startsAt, exactly on resolved rows.
func testAccCheckAlertListContract(reference string, statuses []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := datasourceAttributes(s, reference)
		if err != nil {
			return err
		}
		var count int
		if _, err := fmt.Sscan(attributes["data.#"], &count); err != nil {
			return fmt.Errorf("could not read data.#: %w", err)
		}
		seen := make(map[string]bool, count)
		for i := 0; i < count; i++ {
			get := func(field string) string { return attributes[fmt.Sprintf("data.%d.%s", i, field)] }
			id, status, kind := get("alert_id"), get("status"), get("kind")
			for _, field := range []string{"alert_id", "event_id", "kind", "resource_id", "resource_type", "severity", "status", "starts_at", "event_time"} {
				if get(field) == "" {
					return fmt.Errorf("alert %d (%s) has no %s", i, id, field)
				}
			}
			if seen[id] {
				return fmt.Errorf("alert %s is listed more than once", id)
			}
			seen[id] = true
			if !containsString(statuses, status) {
				return fmt.Errorf("alert %s has status %q, want one of %v", id, status, statuses)
			}
			if !alertKindConfigurable(kind) {
				return fmt.Errorf("alert %s has kind %q, which Alert Configuration does not offer", id, kind)
			}
			startsAt, endsAt := get("starts_at"), get("ends_at")
			switch status {
			case "firing":
				if endsAt != "" {
					return fmt.Errorf("firing alert %s has ends_at %s", id, endsAt)
				}
			case "resolved":
				start, err1 := time.Parse(time.RFC3339Nano, startsAt)
				end, err2 := time.Parse(time.RFC3339Nano, endsAt)
				if err1 != nil || err2 != nil || end.Before(start) {
					return fmt.Errorf("resolved alert %s has starts_at %q and ends_at %q", id, startsAt, endsAt)
				}
			}
		}
		return nil
	}
}

// alertKindAliases maps alert-list kinds that do not normalise to a catalog identifier. Add to it
// when a run reports a configurable alert under a new name. App Service kinds, which use SGW where
// identifiers say syncgateway (SGWCPUUsageHigh), are handled by the normalisation.
var alertKindAliases = map[string]string{
	"CBHighDiskUsage": "disk_usage_high",
}

// alertKindConfigurable reports whether an alert-list kind is one Alert Configuration offers. The
// list names alerts by severity, such as CPUUsageHighWarning for cpu_usage_high, so names are
// compared lowercased, without separators and without a trailing severity.
func alertKindConfigurable(kind string) bool {
	if alias, ok := alertKindAliases[kind]; ok {
		kind = alias
	}
	normalise := func(name string) string {
		name = strings.ToLower(regexp.MustCompile(`[^A-Za-z0-9]`).ReplaceAllString(name, ""))
		for _, severity := range []string{"warning", "critical"} {
			name = strings.TrimSuffix(name, severity)
		}
		if strings.HasPrefix(name, "syncgateway") {
			name = "sgw" + strings.TrimPrefix(name, "syncgateway")
		}
		return name
	}
	for _, name := range alertCatalog {
		if normalise(name) == normalise(kind) {
			return true
		}
	}
	return false
}

// testAccCheckAlertConfigurationUnchanged compares a level with the snapshot *before.
func testAccCheckAlertConfigurationUnchanged(level string, before **alertsapi.GetAlertConfigurationResponse) resource.TestCheckFunc {
	return testAccCheckAlertConfigurationOnServer(level, func(now *alertsapi.GetAlertConfigurationResponse) error {
		if *before == nil {
			return fmt.Errorf("no snapshot of %s was taken", level)
		}
		b := *before
		if err := sameAlertConfigs(b.Configs, now.Configs); err != nil {
			return fmt.Errorf("%s: %w", level, err)
		}
		if err := expectAlertRecipients(now, b.Recipients...); err != nil {
			return fmt.Errorf("%s: %w", level, err)
		}
		if len(now.MutedAlerts) != len(b.MutedAlerts) {
			return fmt.Errorf("%s: muted alerts = %v, want %v", level, now.MutedAlerts, b.MutedAlerts)
		}
		for name, mute := range b.MutedAlerts {
			if !sameAlertMute(mute, now.MutedAlerts[name]) {
				return fmt.Errorf("%s: mute of %s = %+v, want %+v", level, name, now.MutedAlerts[name], mute)
			}
		}
		return nil
	})
}

// sameAlertConfigs compares the values of two levels' configs, ignoring who set them and when.
func sameAlertConfigs(want, got map[string]alertsapi.AlertConfigurationEntry) error {
	if len(want) != len(got) {
		return fmt.Errorf("configs = %v, want %v", got, want)
	}
	for name, w := range want {
		g, ok := got[name]
		if !ok || fmt.Sprint(derefInt64(w.FrequencyHours), derefFloat64(w.Threshold), derefFloat64(w.WarningThreshold), derefFloat64(w.CriticalThreshold)) !=
			fmt.Sprint(derefInt64(g.FrequencyHours), derefFloat64(g.Threshold), derefFloat64(g.WarningThreshold), derefFloat64(g.CriticalThreshold)) {
			return fmt.Errorf("config of %s = %+v, want %+v", name, g, w)
		}
	}
	return nil
}

func expectResolvedAlertInt(c *alertsapi.GetAlertConfigurationResponse, alert, field string, want int64) error {
	return expectAlertInt(&alertsapi.GetAlertConfigurationResponse{Configs: c.ResolvedConfigs}, alert, field, want)
}

func expectResolvedAlertFloat(c *alertsapi.GetAlertConfigurationResponse, alert, field string, want float64) error {
	return expectAlertFloat(&alertsapi.GetAlertConfigurationResponse{Configs: c.ResolvedConfigs}, alert, field, want)
}

// expectAlertMute checks one mute. An empty from or until must be absent; "*" accepts any value.
func expectAlertMute(mutes map[string]alertsapi.AlertMuteState, alert, action, from, until string) error {
	mute, ok := mutes[alert]
	if !ok {
		return fmt.Errorf("%s has no mute, want %s", alert, action)
	}
	if mute.Action != action {
		return fmt.Errorf("%s mute action = %q, want %q", alert, mute.Action, action)
	}
	if !sameAlertInstant(mute.MutedFrom, from) {
		return fmt.Errorf("%s muted_from = %v, want %q", alert, derefString(mute.MutedFrom), from)
	}
	if !sameAlertInstant(mute.MutedUntil, until) {
		return fmt.Errorf("%s muted_until = %v, want %q", alert, derefString(mute.MutedUntil), until)
	}
	return nil
}

// expectAlertMuteSetBy checks that a mute exists and names the API key that set it (HIE-14).
func expectAlertMuteSetBy(mutes map[string]alertsapi.AlertMuteState, alert string, key alertTestAPIKey) error {
	mute, ok := mutes[alert]
	if !ok || mute.Action != "mute" {
		return fmt.Errorf("%s is not muted: %+v", alert, mute)
	}
	if mute.UpdatedBy == nil || (*mute.UpdatedBy != key.id && *mute.UpdatedBy != key.name) {
		return fmt.Errorf("%s updated_by = %v, want the API key %s (%s)", alert, derefString(mute.UpdatedBy), key.name, key.id)
	}
	return nil
}

// testAccCheckNotMutedInState checks that a mutes map in state does not mute alert: it is absent
// or set to unmute.
func testAccCheckNotMutedInState(reference, attribute, alert string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := datasourceAttributes(s, reference)
		if err != nil {
			return err
		}
		if action := attributes[attribute+"."+alert+".action"]; action == "mute" {
			return fmt.Errorf("%s.%s.%s is muted, want not muted", reference, attribute, alert)
		}
		return nil
	}
}

func sameAlertInstant(got *string, want string) bool {
	switch {
	case want == "*":
		return true
	case want == "":
		return got == nil || *got == ""
	case got == nil:
		return false
	}
	g, err1 := time.Parse(time.RFC3339Nano, *got)
	w, err2 := time.Parse(time.RFC3339Nano, want)
	return err1 == nil && err2 == nil && g.Equal(w)
}

func derefString(v *string) any {
	if v == nil {
		return nil
	}
	return *v
}

// sleepUntil waits until the RFC 3339 instant has passed by margin.
func sleepUntil(instant string, margin time.Duration) {
	at, err := time.Parse(time.RFC3339, instant)
	if err != nil {
		panic(err)
	}
	if wait := time.Until(at.Add(margin)); wait > 0 {
		time.Sleep(wait)
	}
}

func listAlertOverridesForTest() ([]alertsapi.AlertConfigurationOverride, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/alertConfiguration/overrides", globalHost, globalOrgId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	return api.GetPaginated[[]alertsapi.AlertConfigurationOverride](context.Background(), globalClient, globalToken, cfg, "")
}

func testAccCheckAlertOverrides(check func([]alertsapi.AlertConfigurationOverride) error) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		overrides, err := listAlertOverridesForTest()
		if err != nil {
			return fmt.Errorf("could not list alert configuration overrides: %w", err)
		}
		return check(overrides)
	}
}

// findAlertOverride returns the override of alertKind at scope on the project or cluster with id.
func findAlertOverride(overrides []alertsapi.AlertConfigurationOverride, scope, id, alertKind string) *alertsapi.AlertConfigurationOverride {
	for i, o := range overrides {
		if o.Scope != scope || o.AlertKind != alertKind {
			continue
		}
		if (o.ResourceId != nil && *o.ResourceId == id) || (scope == "project" && o.ProjectId != nil && *o.ProjectId == id) {
			return &overrides[i]
		}
	}
	return nil
}

// findAlertOverrideInState returns the attribute prefix, such as "data.3.", of the overrides
// data source's row for alertKind on the resource with resourceId.
func findAlertOverrideInState(s *terraform.State, reference, resourceId, alertKind string) (string, error) {
	attributes, err := datasourceAttributes(s, reference)
	if err != nil {
		return "", err
	}
	for key, value := range attributes {
		if !strings.HasSuffix(key, ".resource_id") || value != resourceId {
			continue
		}
		prefix := strings.TrimSuffix(key, "resource_id")
		if attributes[prefix+"alert_kind"] == alertKind {
			return prefix, nil
		}
	}
	return "", fmt.Errorf("%s lists no override of %s on %s", reference, alertKind, resourceId)
}

type alertTestAPIKey struct {
	id, name, token string
}

// createAlertTestAPIKey creates an API key with the roles given, for the RBAC cases, and deletes
// it when the test ends.
func createAlertTestAPIKey(t *testing.T, organizationRoles []string, projectRoles map[string][]string) alertTestAPIKey {
	t.Helper()
	name := randomStringWithPrefix("tf_acc_alert_rbac_")
	description := "Alert configuration RBAC acceptance test"
	allowedCIDRs := []string{"0.0.0.0/0"}
	expiry := float32(1)
	request := api.CreateApiKeyRequest{
		Name:              name,
		Description:       &description,
		OrganizationRoles: organizationRoles,
		AllowedCIDRs:      &allowedCIDRs,
		Expiry:            &expiry,
	}
	if len(projectRoles) > 0 {
		resources := make(api.Resources, 0, len(projectRoles))
		for projectId, roles := range projectRoles {
			projectType := "project"
			resources = append(resources, api.ResourcesItems{Type: &projectType, Roles: roles, Id: uuid.MustParse(projectId)})
		}
		request.Resources = &resources
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/apikeys", globalHost, globalOrgId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := globalClient.ExecuteWithRetry(context.Background(), cfg, request, globalToken, nil)
	if err != nil {
		t.Fatalf("could not create an API key with roles %v %v: %v", organizationRoles, projectRoles, err)
	}
	var created api.CreateApiKeyResponse
	if err := json.Unmarshal(response.Body, &created); err != nil {
		t.Fatalf("could not read the created API key: %v", err)
	}
	t.Cleanup(func() {
		url := fmt.Sprintf("%s/v4/organizations/%s/apikeys/%s", globalHost, globalOrgId, created.Id)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
		if _, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil); err != nil {
			t.Errorf("could not delete API key %s: %v", created.Id, err)
		}
	})
	return alertTestAPIKey{id: created.Id, name: name, token: created.Token}
}

// createAlertTestProject creates a project with token, and deletes it with the suite's own key
// when the test ends.
func createAlertTestProject(t *testing.T, token string) string {
	t.Helper()
	url := fmt.Sprintf("%s/v4/organizations/%s/projects", globalHost, globalOrgId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	request := api.CreateProjectRequest{Name: randomStringWithPrefix("tf_acc_alert_rbac_project_")}
	response, err := globalClient.ExecuteWithRetry(context.Background(), cfg, request, token, nil)
	if err != nil {
		t.Fatalf("could not create a project: %v", err)
	}
	var created api.CreateProjectResponse
	if err := json.Unmarshal(response.Body, &created); err != nil {
		t.Fatalf("could not read the created project: %v", err)
	}
	projectId := created.Id.String()
	t.Cleanup(func() {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s", globalHost, globalOrgId, projectId)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
		if _, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil); err != nil {
			t.Errorf("could not delete project %s: %v", projectId, err)
		}
	})
	return projectId
}
