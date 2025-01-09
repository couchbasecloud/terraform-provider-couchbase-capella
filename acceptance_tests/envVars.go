package acceptance_tests

import "os"

func GetEnvVars() error {
	Host = os.Getenv("TF_VAR_host")
	if Host == "" {
		return ErrHostMissing
	}
	Token = os.Getenv("TF_VAR_auth_token")
	if Token == "" {
		return ErrTokenMissing
	}
	OrgId = os.Getenv("TF_VAR_organization_id")
	if OrgId == "" {
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
