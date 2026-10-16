package resources

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

// ValidOrganizationRoles lists the organization roles accepted by the user and API key resources.
var ValidOrganizationRoles = []string{
	organizationRoleMember,
	organizationRoleOwner,
	organizationRoleCreator,
	organizationRoleBillingAdmin,
	organizationRoleBillingViewer,
	organizationRoleReadOnly,
}

var validProjectRoles = []string{
	projectRoleOwner,
	projectRoleManager,
	projectRoleViewer,
	projectRoleDataReaderWriter,
	projectRoleDataReader,
}
