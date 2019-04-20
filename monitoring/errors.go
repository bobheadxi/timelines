package monitoring

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"runtime"
	"strings"

	"cloud.google.com/go/errorreporting"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/bobheadxi/timelines/config"
	"github.com/bobheadxi/timelines/log"
)

// AttachErrorLogging attaches the GCP Stackdriver Error Reporting client to
// the given zap logger. See https://cloud.google.com/error-reporting/docs/setup/go
func AttachErrorLogging(l *zap.Logger, service string, meta config.BuildMeta, dev bool) (*zap.SugaredLogger, error) {
	cloud := config.NewCloudConfig()
	if cloud.Provider() != config.ProviderGCP {
		return nil, errors.New("error logging only supported for GCP")
	}

	l.Info("setting up GCP error reporting",
		zap.String("project_id", cloud.GCP.ProjectID))

	errHandler := func(error) {}
	if dev {
		errHandler = func(e error) { stdlog.Printf("gcp.error-reporter: %v\n", e) }
	}
	opts := config.NewGCPConnectionOptions()
	reporter, err := errorreporting.NewClient(
		context.Background(),
		cloud.GCP.ProjectID, errorreporting.Config{
			ServiceName:    service,
			ServiceVersion: meta.AnnotatedCommit(dev),
			OnError:        errHandler,
		},
		opts...)
	if err != nil {
		return nil, err
	}

	return l.
		WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, gcpErrorsWrapCore(reporter))
		})).
		Sugar(), nil
}

type gcpErrorReportingZapCore struct {
	reporter *errorreporting.Client
	enc      zapcore.Encoder
}

func gcpErrorsWrapCore(reporter *errorreporting.Client) zapcore.Core {
	return &gcpErrorReportingZapCore{reporter, zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		NameKey:    "logger",
		MessageKey: "msg",
		LevelKey:   "level",

		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})}
}

func (z *gcpErrorReportingZapCore) Enabled(l zapcore.Level) bool {
	return l >= zapcore.WarnLevel
}

func (z *gcpErrorReportingZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	if entry.Stack == "" {
		return nil
	}

	// extract relevant values from fields
	var requestID string
	for _, f := range fields {
		if f.Key == log.LogKeyRID && f.Type == zapcore.StringType {
			requestID = f.String
		}
	}

	// encode everything as the message
	buf, err := z.enc.EncodeEntry(entry, fields)
	if err != nil {
		return fmt.Errorf("failed to encode entry for gcp error reporter: %v", err)
	}

	// report to GCP
	z.reporter.Report(errorreporting.Entry{
		Error: errors.New(buf.String()),
		User:  requestID,

		// GCP Error Reporting does not like Zap's custom stacktraces (from entry.Stack),
		// so a custom stacktrace must be taken that conforms to the standard. Ugh.
		// See stacktrace() for details.
		Stack: stacktrace(),
	})

	return nil
}

func (z *gcpErrorReportingZapCore) Sync() error {
	z.reporter.Flush()
	return nil
}

func (z *gcpErrorReportingZapCore) With(fields []zapcore.Field) zapcore.Core {
	clone := &gcpErrorReportingZapCore{z.reporter, z.enc.Clone()}
	for i := range fields {
		fields[i].AddTo(clone.enc)
	}
	return clone
}

func (z *gcpErrorReportingZapCore) Check(e zapcore.Entry, c *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if z.Enabled(e.Level) {
		return c.AddCore(e, z)
	}
	return c
}

// stackSkip denotes the number of stack levels we need to skip over:
// * stacktrace()
// * Core.Write()
// * CheckedEntry.Write()
// * logger.Warn() or equivalent
const stackSkip = 4

// stacktrace captures the calling goroutine's stack and trims out irrelevant
// levels (as described in documentation for stackSkip)
// TODO: this performs quite poorly, refactor to avoid string conversions
func stacktrace() []byte {
	var buf [16 * 1024]byte
	stack := buf[0:runtime.Stack(buf[:], false)]
	lines := strings.Split(string(stack), "\n")
	lines = append(lines[:1], lines[2*stackSkip+1:]...)
	return []byte(strings.Join(lines, "\n"))
}
