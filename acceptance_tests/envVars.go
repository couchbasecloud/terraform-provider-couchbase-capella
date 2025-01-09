package acceptance_tests

import "os"

func getEnvVars() error {
	globalHost = os.Getenv("TF_VAR_host")
	if globalHost == "" {
		return ErrHostMissing
	}
	globalToken = os.Getenv("TF_VAR_auth_token")
	if globalToken == "" {
		return ErrTokenMissing
	}
	globalOrgId = os.Getenv("TF_VAR_organization_id")
	if globalOrgId == "" {
		return ErrOrgIdMissing
	}
	Username = os.Getenv("CAPELLA_USERNAME")
	if Username == "" {
		return ErrUsernameMissing
	}
	Password = os.Getenv("CAPELLA_PASSWORD")
	if Password == "" {
		return ErrPasswordMissing
	}

	return nil
}
