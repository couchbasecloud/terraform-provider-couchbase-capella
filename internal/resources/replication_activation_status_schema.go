package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var replicationActivationStatusBuilder = capellaschema.NewSchemaBuilder("replicationActivationStatus")

// ReplicationActivationStatusSchema defines the Terraform schema for the replication activation status resource.
func ReplicationActivationStatusSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", replicationActivationStatusBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", replicationActivationStatusBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", replicationActivationStatusBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "replication_id", replicationActivationStatusBuilder, requiredNonEmptyStringAttribute())

	state := stringAttribute([]string{required},
		stringvalidator.OneOf(string(replicationapi.StatusRunning), string(replicationapi.StatusPaused)))
	state.MarkdownDescription = capellaschema.ReplicationActivationStatusStateDescription
	capellaschema.AddAttr(attrs, "state", replicationActivationStatusBuilder, state)

	return schema.Schema{
		MarkdownDescription: "This resource allows you to pause or resume an XDCR replication. Destroying this resource " +
			"does not change the replication.\n\nPause and resume change one direction only. For a `twoWay` replication, " +
			"use one resource for each direction: the replication `id` with the source cluster ID, and the " +
			"`reverse_replication_id` with the target cluster ID.",
		Attributes: attrs,
	}
}
