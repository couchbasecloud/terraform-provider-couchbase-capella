package acceptance_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Schema-validator coverage; every case is rejected at plan time and spells out all its inputs.

const (
	// notAUUID fails the UUID regex but satisfies the length validator, isolating the regex.
	notAUUID = "not-a-uuid"

	// placeholderUUID is well formed, so it fills the fields a case is not exercising.
	placeholderUUID = "11111111-1111-1111-1111-111111111111"
)

// bucketBackupExportIDCase is one invalid-input case; exportId is unused by the resource.
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

// TestAccBucketBackupExportResourceInvalidIDs covers the resource validators not already tested.
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

// TestAccDatasourceBucketBackupExportInvalidIDs covers the data source's separate validators.
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
