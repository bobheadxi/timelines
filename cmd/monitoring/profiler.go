package monitoring

import (
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/profiler"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/bobheadxi/timelines/config"
)

// StartProfiler starts a suitable background profiler based on environment
// variables
func StartProfiler(l *zap.SugaredLogger, service string, meta config.BuildMeta) error {
	cloud := config.NewCloudConfig()
	provider := cloud.Provider()
	switch provider {
	case config.ProviderGCP:
		l.Info("starting server with GCP profiling")
		var opts []option.ClientOption
		if os.Getenv("GOOGLE_APPLICATION_RAW") != "" {
			opts = []option.ClientOption{
				option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_APPLICATION_RAW"))),
			}
		}
		if err := profiler.Start(profiler.Config{
			Service:        service,
			ServiceVersion: meta.Commit,
			ProjectID:      cloud.GCP.ProjectID,
		}, opts...); err != nil {
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
