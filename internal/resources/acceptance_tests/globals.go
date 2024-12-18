package acceptance_tests

import (
	"time"
)

var (
	Host  string
	Token string

	Username string
	Password string

	OrgId         string
	ProjectId     string
	ClusterId     string
	ProviderBlock string
)

const (
	timeout = 60 * time.Second
)
