package config

import "os"

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
