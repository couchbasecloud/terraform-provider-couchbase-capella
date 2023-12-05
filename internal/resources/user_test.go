package resources

import (
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"gotest.tools/assert"
)

func Test_ConstructPatch(t *testing.T) {
	var (
		organizationOwner       = basetypes.NewStringValue("organizationOwner")
		organizationMember      = basetypes.NewStringValue("organizationMember")
		projectViewer           = basetypes.NewStringValue("projectViewer")
		projectDataReaderWriter = basetypes.NewStringValue("projectDataReaderWriter")
		add                     = "add"
		remove                  = "remove"
		orgRolesPath            = "/organizationRoles"
		projectType             = "project"
	)

	type test struct {
		name           string
		existingSchema providerschema.User
		proposedSchema providerschema.User
		expectedPatch  []api.PatchEntry
	}

	tests := []test{
		{
			name: "[POSITIVE] - Successfully add an organization role",
			existingSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources: []providerschema.Resource{
					{
						Id: basetypes.NewStringValue("100"),
						Roles: []basetypes.StringValue{
							projectViewer,
							projectDataReaderWriter,
						},
					},
				},
			},
			proposedSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{
					organizationOwner,
					organizationMember,
				},
				Resources: []providerschema.Resource{
					{
						Id: basetypes.NewStringValue("100"),
						Roles: []basetypes.StringValue{
							projectViewer,
							projectDataReaderWriter,
						},
					},
				},
			},
			expectedPatch: []api.PatchEntry{
				{
					Op:   add,
					Path: orgRolesPath,
					Value: []string{
						organizationOwner.ValueString(),
					},
				},
			},
		},
		{
			name: "[POSITIVE] - Successfully replace an organization role",
			existingSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources: []providerschema.Resource{
					{
						Id: basetypes.NewStringValue("100"),
						Roles: []basetypes.StringValue{
							projectViewer,
							projectDataReaderWriter,
						},
					},
				},
			},
			proposedSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationOwner},
				Resources: []providerschema.Resource{
					{
						Id: basetypes.NewStringValue("100"),
						Roles: []basetypes.StringValue{
							projectViewer,
							projectDataReaderWriter,
						},
					},
				},
			},
			expectedPatch: []api.PatchEntry{
				{
					Op:   add,
					Path: orgRolesPath,
					Value: []string{
						organizationOwner.ValueString(),
					},
				},
				{
					Op:   remove,
					Path: orgRolesPath,
					Value: []string{
						organizationMember.ValueString(),
					},
				},
			},
		},
		{
			name: "[POSITIVE] - Successfully add a project role",
			existingSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources: []providerschema.Resource{
					{
						Id:   basetypes.NewStringValue("100"),
						Type: basetypes.NewStringValue("project"),
						Roles: []basetypes.StringValue{
							projectViewer,
						},
					},
				},
			},
			proposedSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources: []providerschema.Resource{
					{
						Id:   basetypes.NewStringValue("100"),
						Type: basetypes.NewStringValue("project"),
						Roles: []basetypes.StringValue{
							projectViewer,
							projectDataReaderWriter,
						},
					},
				},
			},
			expectedPatch: []api.PatchEntry{
				{
					Op:   add,
					Path: "/resources/100/roles",
					Value: []string{
						projectDataReaderWriter.ValueString(),
					},
				},
			},
		},
		{
			name: "[POSITIVE] - Successfully remove a project role",
			existingSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources: []providerschema.Resource{
					{
						Id:   basetypes.NewStringValue("100"),
						Type: basetypes.NewStringValue("project"),
						Roles: []basetypes.StringValue{
							projectViewer,
							projectDataReaderWriter,
						},
					},
				},
			},
			proposedSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources: []providerschema.Resource{
					{
						Id:   basetypes.NewStringValue("100"),
						Type: basetypes.NewStringValue("project"),
						Roles: []basetypes.StringValue{
							projectViewer,
						},
					},
				},
			},
			expectedPatch: []api.PatchEntry{
				{
					Op:   remove,
					Path: "/resources/100/roles",
					Value: []string{
						projectDataReaderWriter.ValueString(),
					},
				},
			},
		},
		{
			name: "[POSITIVE] - Remove a project role with type omitted",
			existingSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources: []providerschema.Resource{
					{
						Id:   basetypes.NewStringValue("100"),
						Type: basetypes.NewStringValue("project"),
						Roles: []basetypes.StringValue{
							projectViewer,
							projectDataReaderWriter,
						},
					},
				},
			},
			proposedSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources: []providerschema.Resource{
					{
						Id: basetypes.NewStringValue("100"),
						Roles: []basetypes.StringValue{
							projectViewer,
						},
					},
				},
			},
			expectedPatch: []api.PatchEntry{
				{
					Op:   remove,
					Path: "/resources/100/roles",
					Value: []string{
						projectDataReaderWriter.ValueString(),
					},
				},
			},
		},
		{
			name: "[POSITIVE] - Successfully add a resource",
			existingSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources: []providerschema.Resource{
					{
						Id:   basetypes.NewStringValue("100"),
						Type: basetypes.NewStringValue("project"),
						Roles: []basetypes.StringValue{
							projectViewer,
							projectDataReaderWriter,
						},
					},
				},
			},
			proposedSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources: []providerschema.Resource{
					{
						Id:   basetypes.NewStringValue("100"),
						Type: basetypes.NewStringValue("project"),
						Roles: []basetypes.StringValue{
							projectViewer,
							projectDataReaderWriter,
						},
					},
					{
						Id:   basetypes.NewStringValue("200"),
						Type: basetypes.NewStringValue("project"),
						Roles: []basetypes.StringValue{
							projectViewer,
							projectDataReaderWriter,
						},
					},
				},
			},
			expectedPatch: []api.PatchEntry{
				{
					Op:   add,
					Path: "/resources/200",
					Value: api.Resource{
						Id:   "200",
						Type: &projectType,
						Roles: []string{
							projectViewer.ValueString(),
							projectDataReaderWriter.ValueString(),
						},
					},
				},
			},
		},
		{
			name: "[POSITIVE] - Successfully remove a resource",
			existingSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources: []providerschema.Resource{
					{
						Id:   basetypes.NewStringValue("100"),
						Type: basetypes.NewStringValue("project"),
						Roles: []basetypes.StringValue{
							projectViewer,
							projectDataReaderWriter,
						},
					},
				},
			},
			proposedSchema: providerschema.User{
				OrganizationRoles: []basetypes.StringValue{organizationMember},
				Resources:         []providerschema.Resource{},
			},
			expectedPatch: []api.PatchEntry{
				{
					Op:   remove,
					Path: "/resources/100",
					Value: api.Resource{
						Id:   "100",
						Type: &projectType,
						Roles: []string{
							projectViewer.ValueString(),
							projectDataReaderWriter.ValueString(),
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			patch := constructPatch(test.existingSchema, test.proposedSchema)

			assert.Equal(t, len(test.expectedPatch), len(patch))

			for i, v := range test.expectedPatch {
				assert.Equal(t, v.Op, patch[i].Op)
				assert.Equal(t, v.Path, patch[i].Path)

				// Assert value
				switch expectedValue := v.Value.(type) {
				case []string:
					actualVal, _ := patch[i].Value.([]string)
					assert.DeepEqual(t, expectedValue, actualVal)
				case api.Resource:
					actualVal, _ := patch[i].Value.(api.Resource)
					assert.DeepEqual(t, expectedValue, actualVal)
				}
			}
		})
	}
}
