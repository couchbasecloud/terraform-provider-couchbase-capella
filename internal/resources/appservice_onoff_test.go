package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
)

func Test_appServiceStateSatisfies(t *testing.T) {
	tests := []struct {
		name         string
		state        string
		currentState appservice.State
		expected     bool
	}{
		{name: "on is satisfied by healthy", state: "on", currentState: appservice.Healthy, expected: true},
		{name: "on is satisfied by turningOn", state: "on", currentState: appservice.TurningOn, expected: true},
		{name: "on is not satisfied by turnedOff", state: "on", currentState: appservice.TurnedOff, expected: false},
		{name: "on is not satisfied by degraded", state: "on", currentState: appservice.Degraded, expected: false},
		{name: "on is not satisfied by turnOnFailed", state: "on", currentState: appservice.TurnOnFailed, expected: false},
		{name: "off is satisfied by turnedOff", state: "off", currentState: appservice.TurnedOff, expected: true},
		{name: "off is satisfied by turningOff", state: "off", currentState: appservice.TurningOff, expected: true},
		{name: "off is not satisfied by healthy", state: "off", currentState: appservice.Healthy, expected: false},
		{name: "off is not satisfied by turnOffFailed", state: "off", currentState: appservice.TurnOffFailed, expected: false},
		{name: "unknown requested state is never satisfied", state: "frozen", currentState: appservice.Healthy, expected: false},
		{name: "empty current state is never satisfied", state: "on", currentState: "", expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, appServiceStateSatisfies(test.state, test.currentState))
		})
	}
}
