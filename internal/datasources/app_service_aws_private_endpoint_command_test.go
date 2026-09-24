package datasources

import "testing"

// TestExtractAWSServiceNameFromCommand verifies that the AWS VPC endpoint service name is
// correctly parsed out of the `--service-name` flag in a real example command from Couchbase's
// documentation for the App Service private endpoint connection-command API.
func TestExtractAWSServiceNameFromCommand(t *testing.T) {
	command := "aws ec2 create-vpc-endpoint --vpc-id vpc-0e4c66e70f63b51e0 --region us-east-1 " +
		"--service-name com.amazonaws.vpce.us-east-1.vpce-svc-0823b61a6d8cee231 " +
		"--vpc-endpoint-type Interface --subnet-ids subnet-01423b12bd81bb116"

	got, ok := extractAWSServiceNameFromCommand(command)
	if !ok {
		t.Fatalf("expected a match, got none")
	}

	want := "com.amazonaws.vpce.us-east-1.vpce-svc-0823b61a6d8cee231"
	if got != want {
		t.Errorf("extractAWSServiceNameFromCommand() = %q, want %q", got, want)
	}
}

// TestExtractAWSServiceNameFromCommandNoMatch verifies that a command string which doesn't
// contain a `--service-name` flag (e.g. because Couchbase changed the command's phrasing) is
// handled as a non-match rather than an error, per the best-effort contract documented on
// extractAWSServiceNameFromCommand.
func TestExtractAWSServiceNameFromCommandNoMatch(t *testing.T) {
	command := "echo this command does not look like an aws ec2 create-vpc-endpoint invocation at all"

	got, ok := extractAWSServiceNameFromCommand(command)
	if ok {
		t.Fatalf("expected no match, got %q", got)
	}
	if got != "" {
		t.Errorf("extractAWSServiceNameFromCommand() = %q, want empty string when ok is false", got)
	}
}
