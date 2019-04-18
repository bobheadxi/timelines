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
func StartProfiler(l *zap.SugaredLogger, meta config.BuildMeta) error {
	cloud := config.NewCloudConfig()
	provider := cloud.Provider()
	switch provider {
	case config.ProviderGCP:
		l.Info("starting server with GCP profiling")
		if err := profiler.Start(profiler.Config{
			Service:        "timelines-server",
			ServiceVersion: meta.Commit,
			ProjectID:      cloud.GCP.ProjectID,
		}); err != nil {
			return err
		}

	case config.ProviderNone:
		l.Errorw("cloudless provider not yet implemented")
		return errors.New("cloudless provider not yet implemented")

	default:
		return fmt.Errorf("unsupported profiling provider '%s'", provider)
	}
	return nil
}
