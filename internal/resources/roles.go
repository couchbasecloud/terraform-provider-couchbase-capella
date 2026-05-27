package resources

const (
	resourceTypeProject = "project"

	organizationRoleMember  = "organizationMember"
	organizationRoleOwner   = "organizationOwner"
	organizationRoleCreator = "projectCreator"

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
}

var validProjectRoles = []string{
	projectRoleOwner,
	projectRoleManager,
	projectRoleViewer,
	projectRoleDataReaderWriter,
	projectRoleDataReader,
}
