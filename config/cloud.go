package config

import (
	"os"

	"google.golang.org/api/option"
)

// CloudProvider denotes supported providers
type CloudProvider string

const (
	// ProviderNone indicates no cloud provider has been set
	ProviderNone CloudProvider = "none"

	// ProviderGCP indicates Google Cloud configuration has been set
	ProviderGCP CloudProvider = "gcp"
)

// Cloud denotes configuration for different cloud providers
type Cloud struct {
	GCP CloudGoogle
}

// CloudGoogle denotes GCP variables
type CloudGoogle struct {
	ProjectID string
}

// NewCloudConfig instantiates cloud configuration from environment
func NewCloudConfig() Cloud {
	return Cloud{
		GCP: CloudGoogle{
			ProjectID: os.Getenv("GCP_PROJECT_ID"),
		},
	}
}

// Provider returns the configured cloud provider
func (c Cloud) Provider() CloudProvider {
	if c.GCP.ProjectID != "" {
		return ProviderGCP
	}
	return ProviderNone
}

const envGCPCredentialsRaw = "GOOGLE_APPLICATION_CREDENTIALS_RAW"

// NewGCPConnectionOptions returns options needed to connect to GCP services
func NewGCPConnectionOptions() []option.ClientOption {
	var opts []option.ClientOption
	rawCredentials := os.Getenv(envGCPCredentialsRaw)
	if rawCredentials != "" {
		opts = []option.ClientOption{
			option.WithCredentialsJSON([]byte(rawCredentials)),
		}
	}
	return opts
}
