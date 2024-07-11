package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConvertToList converts a types.Set into a slice of strings.
func ConvertToList(ctx context.Context, inputSet types.Set) ([]string, error) {
	elements := make([]types.String, 0, len(inputSet.Elements()))
	diags := inputSet.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, fmt.Errorf("error while extracting list[ elements")
	}

	var convertedList []string
	for _, ele := range elements {
		convertedList = append(convertedList, ele.ValueString())
	}
	return convertedList, nil
}

// BuildQueryParams builds a query string from the given parameters.
func BuildQueryParams(params map[string][]string) string {
	query := url.Values{}

	for key, values := range params {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	return "?" + query.Encode()
}

// MapResponseEventsBody maps the response body from a call to events API to a
// slice of EventItem.
func MapResponseEventsBody(
	ctx context.Context,
	events []api.GetEventResponse,
) ([]providerschema.EventItem, error) {
	var newEventItems []providerschema.EventItem
	for _, event := range events {
		incidentIdsSet, err := ConvertIncidents(ctx, event.IncidentIds)
		if err != nil {
			return nil, err
		}

		kvString, err := ConvertKV(event.Kv)
		if err != nil {
			return nil, err
		}

		newEventItem, err := providerschema.NewEventItem(&event, incidentIdsSet, kvString)
		if err != nil {
			return nil, err
		}
		newEventItems = append(newEventItems, *newEventItem)
	}
	return newEventItems, nil
}

// MapEventResponseBody maps the response body from a call to the event API to an Event.
func MapEventResponseBody(ctx context.Context, event api.GetEventResponse, state providerschema.Event) (*providerschema.Event, error) {
	incidentIdsSet, err := ConvertIncidents(ctx, event.IncidentIds)
	if err != nil {
		return nil, err
	}
	kvString, err := ConvertKV(event.Kv)

	if err != nil {
		return nil, err
	}
	newEventState, err := providerschema.NewEvent(&event, state.OrganizationId, incidentIdsSet, kvString)
	if err != nil {
		return nil, err
	}
	return newEventState, nil
}

// ConvertIncidents converts a slice of UUID incident IDs into a types.Set of strings.
func ConvertIncidents(ctx context.Context, incidentIds *[]uuid.UUID) (types.Set, error) {
	incidentIdsString := make([]string, 0)
	if incidentIds != nil {
		for _, incidentId := range *incidentIds {
			incidentIdsString = append(incidentIdsString, incidentId.String())
		}
	}
	incidentIdsSet, diag := types.SetValueFrom(ctx, types.StringType, incidentIdsString)
	if diag.HasError() {
		return types.Set{}, fmt.Errorf("incidentIds set error")
	}

	return incidentIdsSet, nil
}

// ConvertKV converts a map to a JSON string
func ConvertKV(mp *map[string]interface{}) (types.String, error) {
	jsonData, err := json.Marshal(mp)
	if err != nil {
		return types.StringNull(), err
	}
	return types.StringValue(string(jsonData)), nil
}
