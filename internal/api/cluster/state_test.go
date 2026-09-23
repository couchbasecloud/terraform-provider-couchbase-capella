package cluster

import "testing"

func TestIsTerminalState(t *testing.T) {
	tests := []struct {
		name  string
		state State
		want  bool
	}{
		{name: "healthy is not terminal", state: Healthy, want: false},
		{name: "deploying is not terminal", state: Deploying, want: false},
		{name: "rebalancing is not terminal", state: Rebalancing, want: false},
		{name: "deploymentFailed is terminal", state: DeploymentFailed, want: true},
		{name: "destroyFailed is terminal", state: DestroyFailed, want: true},
		{name: "peeringFailed is terminal", state: PeeringFailed, want: true},
		{name: "rebalanceFailed is terminal", state: RebalanceFailed, want: true},
		{name: "scaleFailed is terminal", state: ScaleFailed, want: true},
		{name: "upgradeFailed is terminal", state: UpgradeFailed, want: true},
		{name: "degraded is terminal", state: Degraded, want: true},
		{name: "unknown state is not terminal", state: State("somethingElse"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTerminalState(tt.state); got != tt.want {
				t.Errorf("IsTerminalState(%q) = %v, want %v", tt.state, got, tt.want)
			}
		})
	}
}

func TestState_Equal(t *testing.T) {
	tests := []struct {
		name string
		s1   State
		s2   State
		want bool
	}{
		{name: "exact match", s1: Healthy, s2: Healthy, want: true},
		{name: "case insensitive match", s1: State("Healthy"), s2: Healthy, want: true},
		{name: "mismatch", s1: Healthy, s2: DeploymentFailed, want: false},
		{name: "empty states are equal", s1: State(""), s2: State(""), want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s1.Equal(tt.s2); got != tt.want {
				t.Errorf("%q.Equal(%q) = %v, want %v", tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}
