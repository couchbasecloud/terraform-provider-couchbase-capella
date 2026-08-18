package appservice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatePredicates(t *testing.T) {
	tests := []struct {
		name        string
		state       State
		wantFinal   bool
		wantFailure bool
	}{
		{name: "pending is transitional", state: Pending, wantFinal: false, wantFailure: false},
		{name: "deploying is transitional", state: Deploying, wantFinal: false, wantFailure: false},
		{name: "destroying is transitional", state: Destroying, wantFinal: false, wantFailure: false},
		{name: "scaling is transitional", state: Scaling, wantFinal: false, wantFailure: false},
		{name: "upgrading is transitional", state: Upgrading, wantFinal: false, wantFailure: false},
		{name: "turningOn is transitional", state: TurningOn, wantFinal: false, wantFailure: false},
		{name: "turningOff is transitional", state: TurningOff, wantFinal: false, wantFailure: false},
		{name: "healthy is a success", state: Healthy, wantFinal: true, wantFailure: false},
		{name: "degraded is a success", state: Degraded, wantFinal: true, wantFailure: false},
		{name: "turnedOff is a success", state: TurnedOff, wantFinal: true, wantFailure: false},
		{name: "deploymentFailed is a failure", state: DeploymentFailed, wantFinal: true, wantFailure: true},
		{name: "destroyFailed is a failure", state: DestroyFailed, wantFinal: true, wantFailure: true},
		{name: "scaleFailed is a failure", state: ScaleFailed, wantFinal: true, wantFailure: true},
		{name: "upgradeFailed is a failure", state: UpgradeFailed, wantFinal: true, wantFailure: true},
		{name: "turnOnFailed is a failure", state: TurnOnFailed, wantFinal: true, wantFailure: true},
		{name: "turnOffFailed is a failure", state: TurnOffFailed, wantFinal: true, wantFailure: true},
		{name: "unrecognised state is transitional", state: State("somethingElse"), wantFinal: false, wantFailure: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantFinal, IsFinalState(tt.state))
			assert.Equal(t, tt.wantFailure, IsFailureState(tt.state))

			// Every failure state must also be a final state
			if tt.wantFailure {
				assert.True(t, IsFinalState(tt.state))
			}
		})
	}
}
