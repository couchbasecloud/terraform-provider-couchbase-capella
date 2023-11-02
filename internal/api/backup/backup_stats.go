package backup

// BackupStats represents various backup level data that couchbase provides.
type BackupStats struct {
	// SizeInMB represents backup size in megabytes.
	SizeInMB float64 `json:"sizeInMb"`

	// Items is the number of items saved during the backup.
	Items int64 `json:"items"`

	// Mutations is the number of mutations saved during the backup.
	Mutations int64 `json:"mutations"`

	// Tombstones is the number of tombstones saved during the backup.
	Tombstones int64 `json:"tombstones"`

	// GSI is the number of global secondary indexes saved during the backup.
	GSI int64 `json:"gsi"`

	// FTS is the number of full text search entities saved during the backup.
	FTS int64 `json:"fts"`

	// CBAS is the number of analytics entities saved during the backup.
	CBAS int64 `json:"cbas"`

	// Event represents the number of event entities saved during the backup.
	Event int64 `json:"event"`
}
