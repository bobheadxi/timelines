package monitoring

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/errorreporting"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/bobheadxi/timelines/config"
)

// AttachErrorLogging attaches the GCP Stackdriver Error Reporting client to
// the given zap logger. See https://cloud.google.com/error-reporting/docs/setup/go
func AttachErrorLogging(l *zap.SugaredLogger, service string, meta config.BuildMeta) (*zap.SugaredLogger, error) {
	cloud := config.NewCloudConfig()
	if cloud.Provider() != config.ProviderGCP {
		return nil, errors.New("error logging only supported for GCP")
	}

	l.Info("setting up GCP error reporting")
	opts := config.NewGCPConnectionOptions()
	reporter, err := errorreporting.NewClient(
		context.Background(),
		cloud.GCP.ProjectID, errorreporting.Config{
			ServiceName:    service,
			ServiceVersion: meta.Commit,
			OnError:        func(error) {},
		},
		opts...)
	if err != nil {
		return nil, err
	}

	return l.Desugar().
		WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return &zapcoreGCPErrors{c, reporter, zapcore.NewJSONEncoder(zapcore.EncoderConfig{})}
		})).
		Sugar(), nil
}

type zapcoreGCPErrors struct {
	zapcore.Core

	reporter *errorreporting.Client
	encoder  zapcore.Encoder
}

func (z *zapcoreGCPErrors) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	if entry.Level == zapcore.ErrorLevel {
		buf, err := z.encoder.EncodeEntry(entry, fields)
		if err != nil {
			return fmt.Errorf("failed to encode entry for gcp error reporter: %v", err)
		}
		z.reporter.Report(errorreporting.Entry{
			Error: errors.New(buf.String()),
			Stack: []byte(entry.Stack),
		})
	}
	return z.Core.Write(entry, fields)
}

func (z *zapcoreGCPErrors) Sync() error {
	z.reporter.Flush()
	return z.Core.Sync()
}
