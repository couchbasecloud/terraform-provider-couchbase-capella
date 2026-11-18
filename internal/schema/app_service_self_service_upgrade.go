package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_self_service_upgrade"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// Fallback descriptions for the App Service self-service upgrade fields, used only until the published
// OpenAPI spec carries them. The leading "\n - " matches the bullet format buildEnhancedDescription
// emits, so the generated docs do not change once the spec takes over.
const (
	SelfServiceUpgradeIdDescription        = "\n - The ID of the self-service upgrade. It is the same value as the `id` of the matching App Service maintenance."
	SelfServiceUpgradeVersionIdDescription = "\n - The ID of the App Service version to upgrade to, as listed by the `couchbase-capella_app_service_self_service_upgrade_available_versions` data source. It can be changed while the upgrade is `pending`."
	SelfServiceUpgradeQueueAtDescription   = "\n - The RFC3339 timestamp at which the upgrade is queued. It must be in the future and no more than 7 days ahead. Omit it to upgrade the App Service now."

	AvailableVersionIdDescription           = "\n - The ID of the App Service version. Use it as the `version_id` of the `couchbase-capella_app_service_self_service_upgrade` resource."
	AvailableVersionTitleDescription        = "\n - The title of the version."
	AvailableVersionDescriptionDescription  = "\n - A short description of the version."
	AvailableVersionReleaseNotesDescription = "\n - Detailed notes on what the version offers."
	AvailableVersionPriorityDescription     = "\n - The priority of the version, for example `critical` or `severe`."
)

// AppServiceSelfServiceUpgrade is a customer-initiated upgrade of an App Service.
type AppServiceSelfServiceUpgrade struct {
	OrganizationId types.String      `tfsdk:"organization_id"`
	ProjectId      types.String      `tfsdk:"project_id"`
	ClusterId      types.String      `tfsdk:"cluster_id"`
	AppServiceId   types.String      `tfsdk:"app_service_id"`
	Id             types.String      `tfsdk:"id"`
	VersionId      types.String      `tfsdk:"version_id"`
	QueueAt        timetypes.RFC3339 `tfsdk:"queue_at"`
	CurrentState   types.String      `tfsdk:"current_state"`
	CompletedAt    timetypes.RFC3339 `tfsdk:"completed_at"`
	Audit          types.Object      `tfsdk:"audit"`
}

// Validate is used to verify that IDs have been properly imported.
func (a AppServiceSelfServiceUpgrade) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
		Id:             a.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

// AppServiceSelfServiceUpgrades is the list of self-service upgrades of an App Service. The API
// returns them in the App Service maintenance shape.
type AppServiceSelfServiceUpgrades struct {
	OrganizationId types.String                `tfsdk:"organization_id"`
	ProjectId      types.String                `tfsdk:"project_id"`
	ClusterId      types.String                `tfsdk:"cluster_id"`
	AppServiceId   types.String                `tfsdk:"app_service_id"`
	Data           []AppServiceMaintenanceData `tfsdk:"data"`
}

// AppServiceSelfServiceUpgradeAvailableVersions is the list of versions an App Service can be upgraded to.
type AppServiceSelfServiceUpgradeAvailableVersions struct {
	OrganizationId types.String                                   `tfsdk:"organization_id"`
	ProjectId      types.String                                   `tfsdk:"project_id"`
	ClusterId      types.String                                   `tfsdk:"cluster_id"`
	AppServiceId   types.String                                   `tfsdk:"app_service_id"`
	Data           []AppServiceSelfServiceUpgradeAvailableVersion `tfsdk:"data"`
}

// AppServiceSelfServiceUpgradeAvailableVersion is a single version an App Service can be upgraded to.
type AppServiceSelfServiceUpgradeAvailableVersion struct {
	Id           types.String `tfsdk:"id"`
	Title        types.String `tfsdk:"title"`
	Description  types.String `tfsdk:"description"`
	ReleaseNotes types.List   `tfsdk:"release_notes"`
	Priority     types.String `tfsdk:"priority"`
}

// NewAppServiceSelfServiceUpgradeAvailableVersion maps an available version from the V4 API to its Terraform model.
func NewAppServiceSelfServiceUpgradeAvailableVersion(ctx context.Context, version appservice_self_service_upgrade.AvailableVersion) (AppServiceSelfServiceUpgradeAvailableVersion, diag.Diagnostics) {
	releaseNotes, diags := types.ListValueFrom(ctx, types.StringType, version.ReleaseNotes)
	if diags.HasError() {
		return AppServiceSelfServiceUpgradeAvailableVersion{}, diags
	}

	return AppServiceSelfServiceUpgradeAvailableVersion{
		Id:           types.StringValue(version.Id),
		Title:        types.StringValue(version.Title),
		Description:  types.StringValue(version.Description),
		ReleaseNotes: releaseNotes,
		Priority:     types.StringValue(version.Priority),
	}, nil
}
