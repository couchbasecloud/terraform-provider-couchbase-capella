package snapshot_backup

import "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"

type Progress struct {
	Status State  `json:"status"`
	Time   string `json:"time"`
}

type CMEK struct {
	ID         string `json:"id"`
	ProviderID string `json:"providerId"`
}

type CMEKs []CMEK

type Server struct {
	Version string `json:"version"`
}

type CrossRegionCopy struct {
	RegionCode string `json:"regionCode"`
	Status     State  `json:"status"`
	Time       string `json:"time"`
}

type SnapshotBackup struct {
	ClusterID         string            `json:"clusterId"`
	CreatedAt         string            `json:"createdAt"`
	Expiration        string            `json:"expiration"`
	ID                string            `json:"id"`
	Progress          Progress          `json:"progress"`
	ProjectID         string            `json:"projectId"`
	Retention         int64             `json:"retention"`
	CrossRegionCopies []CrossRegionCopy `json:"crossRegionCopies"`
	CMEK              CMEKs             `json:"cmek"`
	Server            Server            `json:"server"`
	Size              int               `json:"size"`
	OrganizationID    string            `json:"tenantId"`
	Type              string            `json:"type"`
}

type CreateSnapshotBackupRequest struct {
	Retention     int64    `json:"retention"`
	RegionsToCopy []string `json:"regionsToCopy"`
}

type CreateSnapshotBackupResponse struct {
	ID string `json:"backupId"`
}

type ListSnapshotBackupsResponse struct {
	Data []SnapshotBackup `json:"data"`
}

type EditBackupRetentionRequest struct {
	Retention int64 `json:"retention"`
}

type ProjectSnapshot struct {
	ClusterID         string            `json:"clusterId"`
	CreatedAt         string            `json:"createdAt"`
	Expiration        string            `json:"expiration"`
	ID                string            `json:"id"`
	Progress          Progress          `json:"progress"`
	ProjectID         string            `json:"projectId"`
	Retention         int64             `json:"retention"`
	AppService        string            `json:"appService"`
	CMEK              CMEKs             `json:"cmek"`
	CrossRegionCopies []CrossRegionCopy `json:"crossRegionCopies"`
	Server            Server            `json:"server"`
	DatabaseSize      int               `json:"databaseSize"`
	OrganizationID    string            `json:"tenantId"`
	Type              string            `json:"type"`
}

type ProjectSnapshotBackupData struct {
	ClusterID          string          `json:"clusterId"`
	ClusterName        string          `json:"clusterName"`
	CreationDateTime   string          `json:"creationDateTime"`
	CreatedBy          string          `json:"createdBy"`
	CurrentStatus      string          `json:"currentStatus"`
	CloudProvider      string          `json:"cloudProvider"`
	Region             string          `json:"region"`
	MostRecentSnapshot ProjectSnapshot `json:"mostRecentSnapshot"`
	OldestSnapshot     ProjectSnapshot `json:"oldestSnapshot"`
}

type ListProjectSnapshotBackupsResponse struct {
	Data   []ProjectSnapshotBackupData `json:"data"`
	Cursor *api.Cursor                 `json:"cursor"`
}
