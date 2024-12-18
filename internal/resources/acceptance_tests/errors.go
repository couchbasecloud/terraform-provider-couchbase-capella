package acceptance_tests

import "emperror.dev/errors"

const (
	ErrHostMissing     = errors.Sentinel("Capella host is missing.  Set this in TF_VAR_host env var.")
	ErrTokenMissing    = errors.Sentinel("Capella host is missing.  Set this in TF_VAR_auth_token env var.")
	ErrOrgIdMissing    = errors.Sentinel("Capella organization ID is missing.  Set this in TF_VAR_organization_id env var.")
	ErrUsernameMissing = errors.Sentinel("Username is missing.  Set this in CAPELLA_USERNAME env var.")
	ErrPasswordMissing = errors.Sentinel("Username is missing.  Set this in CAPELLA_PASSWORD env var.")
	ErrNoJWT           = errors.Sentinel("Could not get JWT token.")
	ErrNoCIDR          = errors.Sentinel("Could not get suggested CIDR.")
)
