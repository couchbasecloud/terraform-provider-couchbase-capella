package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"testing"
)

func TestCluster_Validate(t *testing.T) {
	type fields struct {
		Id              types.String
		AppServiceId    types.String
		Audit           types.Object
		OrganizationId  types.String
		ProjectId       types.String
		Availability    *Availability
		CloudProvider   *CloudProvider
		CouchbaseServer *CouchbaseServer
		Description     types.String
		Name            types.String
		ServiceGroups   []ServiceGroup
		Support         *Support
		CurrentState    types.String
		Etag            types.String
		IfMatch         types.String
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				Id:              tt.fields.Id,
				AppServiceId:    tt.fields.AppServiceId,
				Audit:           tt.fields.Audit,
				OrganizationId:  tt.fields.OrganizationId,
				ProjectId:       tt.fields.ProjectId,
				Availability:    tt.fields.Availability,
				CloudProvider:   tt.fields.CloudProvider,
				CouchbaseServer: tt.fields.CouchbaseServer,
				Description:     tt.fields.Description,
				Name:            tt.fields.Name,
				ServiceGroups:   tt.fields.ServiceGroups,
				Support:         tt.fields.Support,
				CurrentState:    tt.fields.CurrentState,
				Etag:            tt.fields.Etag,
				IfMatch:         tt.fields.IfMatch,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
