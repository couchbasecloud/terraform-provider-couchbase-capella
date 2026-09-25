package acceptance_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Schema-validator coverage for the bucket backup export resource and data source.
//
// Every case here is rejected while planning, so none of them reaches the V4 API - which is
// the point. Without these, a dropped validator would only surface as a 4xx at apply time, or
// as a request built from a malformed ID and sent anyway.
//
// Each case spells out all of its inputs rather than defaulting the ones it is not exercising.
// That is deliberate: a case that left another field empty would trip that field's validator
// too, and the test would then pass on the wrong diagnostic.
//
// Note that bucket_id carries a length validator but no UUID validator, in both the resource
// and the data source: it is the URL-compatible base64 encoding of the bucket name, not a GUID.

const (
	// notAUUID fails the UUID regex while satisfying the length validator, so it isolates the
	// regex.
	notAUUID = "not-a-uuid"

	// placeholderUUID is well formed, so it never trips a validator. It fills the fields a
	// given case is not exercising.
	placeholderUUID = "11111111-1111-1111-1111-111111111111"
)

// bucketBackupExportIDCase is one invalid-input case. exportId is ignored by the resource,
// which computes its own id.
type bucketBackupExportIDCase struct {
	name           string
	organizationId string
	projectId      string
	clusterId      string
	bucketId       string
	backupId       string
	exportId       string
	expected       *regexp.Regexp
}

// TestAccBucketBackupExportResourceInvalidIDs covers the validators on the resource's five
// required inputs.
//
// Empty bucket_id and malformed backup_id already have dedicated tests
// (TestAccBucketBackupExportResourceEmptyBucketId and
// TestAccBucketBackupExportResourceMalformedBackupId), so they are not repeated here.
func TestAccBucketBackupExportResourceInvalidIDs(t *testing.T) {
	testCases := []bucketBackupExportIDCase{
		{
			name:           "malformed organization_id",
			organizationId: notAUUID,
			projectId:      globalProjectId,
			clusterId:      globalClusterId,
			bucketId:       globalBucketId,
			backupId:       placeholderUUID,
			expected:       regexp.MustCompile(`(?s)Attribute organization_id must be a valid UUID`),
		},
		{
			name:           "malformed project_id",
			organizationId: globalOrgId,
			projectId:      notAUUID,
			clusterId:      globalClusterId,
			bucketId:       globalBucketId,
			backupId:       placeholderUUID,
			expected:       regexp.MustCompile(`(?s)Attribute project_id must be a valid UUID`),
		},
		{
			name:           "malformed cluster_id",
			organizationId: globalOrgId,
			projectId:      globalProjectId,
			clusterId:      notAUUID,
			bucketId:       globalBucketId,
			backupId:       placeholderUUID,
			expected:       regexp.MustCompile(`(?s)Attribute cluster_id must be a valid UUID`),
		},
		{
			name:           "empty backup_id",
			organizationId: globalOrgId,
			projectId:      globalProjectId,
			clusterId:      globalClusterId,
			bucketId:       globalBucketId,
			backupId:       "",
			expected:       regexp.MustCompile(`(?s)Attribute backup_id string length must be at least 1, got: 0`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resourceName := randomStringWithPrefix("tf_acc_export_bad_id_")

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: testAccBucketBackupExportResourceConfigWithIDs(
							resourceName, tc.organizationId, tc.projectId, tc.clusterId, tc.bucketId, tc.backupId,
						),
						ExpectError: tc.expected,
					},
				},
			})
		})
	}
}

// TestAccDatasourceBucketBackupExportInvalidIDs covers the same validators on the data source.
// They are declared separately from the resource's, in internal/datasources/attributes.go, so
// the two can drift - and the data source additionally validates id, which the resource
// computes.
//
// Malformed id already has a dedicated test (TestAccDatasourceBucketBackupExportMalformedId).
func TestAccDatasourceBucketBackupExportInvalidIDs(t *testing.T) {
	testCases := []bucketBackupExportIDCase{
		{
			name:           "malformed organization_id",
			organizationId: notAUUID,
			projectId:      globalProjectId,
			clusterId:      globalClusterId,
			bucketId:       globalBucketId,
			backupId:       placeholderUUID,
			exportId:       placeholderUUID,
			expected:       regexp.MustCompile(`(?s)Attribute organization_id must be a valid UUID`),
		},
		{
			name:           "malformed project_id",
			organizationId: globalOrgId,
			projectId:      notAUUID,
			clusterId:      globalClusterId,
			bucketId:       globalBucketId,
			backupId:       placeholderUUID,
			exportId:       placeholderUUID,
			expected:       regexp.MustCompile(`(?s)Attribute project_id must be a valid UUID`),
		},
		{
			name:           "malformed cluster_id",
			organizationId: globalOrgId,
			projectId:      globalProjectId,
			clusterId:      notAUUID,
			bucketId:       globalBucketId,
			backupId:       placeholderUUID,
			exportId:       placeholderUUID,
			expected:       regexp.MustCompile(`(?s)Attribute cluster_id must be a valid UUID`),
		},
		{
			name:           "malformed backup_id",
			organizationId: globalOrgId,
			projectId:      globalProjectId,
			clusterId:      globalClusterId,
			bucketId:       globalBucketId,
			backupId:       notAUUID,
			exportId:       placeholderUUID,
			expected:       regexp.MustCompile(`(?s)Attribute backup_id must be a valid UUID`),
		},
		{
			name:           "empty backup_id",
			organizationId: globalOrgId,
			projectId:      globalProjectId,
			clusterId:      globalClusterId,
			bucketId:       globalBucketId,
			backupId:       "",
			exportId:       placeholderUUID,
			expected:       regexp.MustCompile(`(?s)Attribute backup_id string length must be at least 1, got: 0`),
		},
		{
			name:           "empty bucket_id",
			organizationId: globalOrgId,
			projectId:      globalProjectId,
			clusterId:      globalClusterId,
			bucketId:       "",
			backupId:       placeholderUUID,
			exportId:       placeholderUUID,
			expected:       regexp.MustCompile(`(?s)Attribute bucket_id string length must be at least 1, got: 0`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dsName := randomStringWithPrefix("tf_acc_export_ds_bad_id_")

			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Steps: []resource.TestStep{
					{
						Config: testAccBucketBackupExportDatasourceConfigWithIDs(
							dsName, tc.organizationId, tc.projectId, tc.clusterId, tc.bucketId, tc.backupId, tc.exportId,
						),
						ExpectError: tc.expected,
					},
				},
			})
		})
	}
}
