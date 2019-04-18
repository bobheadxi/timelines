package monitoring

import (
	"errors"
	"fmt"

	"cloud.google.com/go/profiler"
	"go.uber.org/zap"

	"github.com/bobheadxi/timelines/config"
)

// StartProfiler starts a suitable background profiler based on environment
// variables
func StartProfiler(l *zap.SugaredLogger, service string, meta config.BuildMeta, dev bool) error {
	cloud := config.NewCloudConfig()
	provider := cloud.Provider()
	switch provider {
	case config.ProviderGCP:
		l.Infow("setting up GCP profiling",
			"project_id", cloud.GCP.ProjectID)
		opts := config.NewGCPConnectionOptions()
		if err := profiler.Start(profiler.Config{
			Service:        service,
			ServiceVersion: meta.AnnotatedCommit(dev),
			ProjectID:      cloud.GCP.ProjectID,
			DebugLogging:   dev,
		}, opts...); err != nil {
			return err
		}

	case config.ProviderNone:
		return errors.New("cloudless provider not yet implemented")

	default:
		return fmt.Errorf("unsupported profiling provider '%s'", provider)
	}
	return nil
}
