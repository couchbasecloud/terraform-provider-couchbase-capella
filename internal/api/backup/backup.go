package backup

// GetBackupResponse is the response received from the Capella V4 Public API when asked to fetch details of an existing backup.
//
// To learn more about backup and restore, see https://docs.couchbase.com/cloud/clusters/backup-restore.html
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetBackupResponse struct {
	BackupStats          *BackupStats  `json:"stats"`
	ScheduleInfo         *ScheduleInfo `json:"scheduleInfo"`
	Method               string        `json:"method"`
	BucketName           string        `json:"bucketName"`
	CycleId              string        `json:"cycleID"`
	Date                 string        `json:"date"`
	RestoreBefore        string        `json:"restoreBefore"`
	Status               State         `json:"status"`
	Id                   string        `json:"id"`
	ClusterId            string        `json:"clusterID"`
	BucketId             string        `json:"bucketID"`
	Source               string        `json:"source"`
	CloudProvider        string        `json:"provider"`
	ProjectId            string        `json:"projectID"`
	OrganizationId       string        `json:"organizationID"`
	ElapsedTimeInSeconds int64         `json:"elapsedTimeInSeconds"`
}

// CreateBackupRequest is the request payload sent to the Capella V4 Public API in order to create a new backup.
//
// Couchbase supports a robust scheduled backup and retention time policy as part of an overall disaster recovery plan for production data.
// Couchbase Capella supports scheduled and on-demand backups of bucket data. A backup can be restored to the same database where it was created or another database in the same organization.
// An on-demand backup of a bucket is always a Full backup. Capella schedules on-demand backup to start immediately.
// On setting up a backup schedule, the bucket automatically backs up the bucket based on the chosen schedule.
//
// To learn more about backup and restore, see https://docs.couchbase.com/cloud/clusters/backup-restore.html
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
// Organization Owner
// Project Owner
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type CreateBackupRequest struct{}

// GetCycleResponse is the response received from the Capella V4 Public API when asked to fetch details of an existing cycle.
//
// To learn more about backup and restore, see https://docs.couchbase.com/cloud/clusters/backup-restore.html
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetCycleResponse struct {
	// CycleId is the cycleId to the which the backup belongs to.
	CycleId string `json:"cycleID"`
}

// GetCyclesResponse is the response received from the Capella V4 Public API when asked to list all cycles for a bucket in a cluster.
//
// To learn more about backup and restore, see https://docs.couchbase.com/cloud/clusters/backup-restore.html
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetCyclesResponse struct {
	Data []GetCycleResponse `json:"data"`
}

// GetBackupsResponse is the response received from the Capella V4 Public API when asked to list all backups.
//
// In order to access this endpoint, the provided API key must have at least one of the following roles:
//
// Organization Owner
// Project Owner
// To learn more, see https://docs.couchbase.com/cloud/organizations/organization-projects-overview.html
type GetBackupsResponse struct {
	Data []GetBackupResponse `json:"data"`
}

type CreateRestoreRequest struct {
	Services              *[]Service `json:"services"`
	FilterKeys            string     `json:"filterKeys"`
	BackupId              string     `json:"backupID"`
	SourceClusterId       string     `json:"sourceClusterID"`
	TargetClusterId       string     `json:"targetClusterID"`
	FilterValues          string     `json:"filterValues"`
	IncludeData           string     `json:"includeData"`
	ExcludeData           string     `json:"excludeData"`
	MapData               string     `json:"mapData"`
	ReplaceTTL            string     `json:"replaceTTL"`
	ReplaceTTLWith        string     `json:"replaceTTLWith"`
	ForceUpdates          bool       `json:"forceUpdates"`
	AutoRemoveCollections bool       `json:"autoRemoveCollections"`
}

type Service string
