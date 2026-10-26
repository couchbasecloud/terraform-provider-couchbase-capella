package resources

import "slices"

const (
	resourceTypeProject = "project"

	organizationRoleMember  = "organizationMember"
	organizationRoleOwner   = "organizationOwner"
	organizationRoleCreator = "projectCreator"

	organizationRoleBillingAdmin  = "organizationBillingAdmin"
	organizationRoleBillingViewer = "organizationBillingViewer"
	organizationRoleReadOnly      = "organizationReadOnly"

	projectRoleOwner            = "projectOwner"
	projectRoleManager          = "projectManager"
	projectRoleViewer           = "projectViewer"
	projectRoleDataReaderWriter = "projectDataReaderWriter"
	projectRoleDataReader       = "projectDataReader"
)

var validOrganizationRoles = []string{
	organizationRoleMember,
	organizationRoleOwner,
	organizationRoleCreator,
	organizationRoleBillingAdmin,
	organizationRoleBillingViewer,
	organizationRoleReadOnly,
}

// ValidOrganizationRoles returns a copy of the organization roles accepted by the user and API key resources.
func ValidOrganizationRoles() []string {
	return slices.Clone(validOrganizationRoles)
}

var validProjectRoles = []string{
	projectRoleOwner,
	projectRoleManager,
	projectRoleViewer,
	projectRoleDataReaderWriter,
	projectRoleDataReader,
}
