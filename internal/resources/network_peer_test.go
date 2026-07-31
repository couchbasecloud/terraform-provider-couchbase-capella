package resources

import (
	"strings"
	"testing"

	stderrors "errors"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func TestValidateProviderConfigMatchesProviderType(t *testing.T) {
	awsConfig := &providerschema.AWSConfig{}
	gcpConfig := &providerschema.GCPConfig{}
	azureConfig := &providerschema.AzureConfig{}

	tests := []struct {
		name           string
		providerType   string
		providerConfig *providerschema.ProviderConfig
		wantErr        string
	}{
		{
			name:           "aws type with aws config",
			providerType:   "aws",
			providerConfig: &providerschema.ProviderConfig{AWSConfig: awsConfig},
		},
		{
			name:           "gcp type with gcp config",
			providerType:   "gcp",
			providerConfig: &providerschema.ProviderConfig{GCPConfig: gcpConfig},
		},
		{
			name:           "azure type with azure config",
			providerType:   "azure",
			providerConfig: &providerschema.ProviderConfig{AzureConfig: azureConfig},
		},
		{
			name:           "azure type with aws config",
			providerType:   "azure",
			providerConfig: &providerschema.ProviderConfig{AWSConfig: awsConfig},
			wantErr:        `provider_type "azure" requires provider_config.azure_config, but provider_config.aws_config was configured`,
		},
		{
			name:           "aws type with gcp config",
			providerType:   "aws",
			providerConfig: &providerschema.ProviderConfig{GCPConfig: gcpConfig},
			wantErr:        `provider_type "aws" requires provider_config.aws_config, but provider_config.gcp_config was configured`,
		},
		{
			name:           "gcp type with azure config",
			providerType:   "gcp",
			providerConfig: &providerschema.ProviderConfig{AzureConfig: azureConfig},
			wantErr:        `provider_type "gcp" requires provider_config.gcp_config, but provider_config.azure_config was configured`,
		},
		{
			name:           "empty provider config",
			providerType:   "aws",
			providerConfig: &providerschema.ProviderConfig{},
			wantErr:        errors.ErrProviderConfigCannotBeEmpty.Error(),
		},
		{
			name:           "multiple config blocks",
			providerType:   "aws",
			providerConfig: &providerschema.ProviderConfig{AWSConfig: awsConfig, AzureConfig: azureConfig},
			wantErr:        "only one of aws_config, gcp_config or azure_config may be set in provider_config, but got: aws_config, azure_config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProviderConfigMatchesProviderType(tt.providerType, tt.providerConfig)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got: %v", tt.wantErr, err)
			}
			if tt.wantErr == errors.ErrProviderConfigCannotBeEmpty.Error() && !stderrors.Is(err, errors.ErrProviderConfigCannotBeEmpty) {
				t.Fatalf("expected error to wrap ErrProviderConfigCannotBeEmpty, got: %v", err)
			}
		})
	}
}
