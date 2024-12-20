package acceptance_tests

import (
	"time"
)

var (
	ProviderBlock string
	Host          string
	Token         string

	Username string
	Password string

	OrgId     string
	ProjectId string
	ClusterId string

	Bucket string
)

const (
	timeout = 60 * time.Second
)
