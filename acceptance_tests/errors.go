package acceptance_tests

import "emperror.dev/errors"

const (
	ErrHostMissing              = errors.Sentinel("Capella host is missing.  Set this in TF_VAR_host env var.")
	ErrTokenMissing             = errors.Sentinel("Capella host is missing.  Set this in TF_VAR_auth_token env var.")
	ErrOrgIdMissing             = errors.Sentinel("Capella organization ID is missing.  Set this in TF_VAR_organization_id env var.")
	ErrNoJWT                    = errors.Sentinel("Could not get JWT token.")
	ErrNoCIDR                   = errors.Sentinel("Could not get suggested CIDR.")
	ErrUnknownHost              = errors.Sentinel("Unknown host name.")
	ErrTimeoutWaitingForCluster = errors.Sentinel("timeout waiting for cluster to be created or destroyed.")
	ErrTimeoutWaitingForBucket  = errors.Sentinel("timeout waiting for bucket to be created.")
)
